package app

import (
	"archive/zip"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"text/template"
	"time"

	"github.com/google/uuid"
)

// CreateTemplate creates a new template
func (a *App) CreateTemplate(request CreateTemplateRequest) (*Template, error) {
	if a.ctx == nil {
		return nil, createInternalError("application context is not available", nil)
	}

	// Validate request
	if err := a.validateCreateTemplateRequest(&request); err != nil {
		return nil, err
	}

	// Check if template with same name already exists
	existing, err := a.repo.Templates.GetByName(request.Name)
	if err == nil && existing != nil {
		return nil, createValidationError(
			fmt.Sprintf("template with name '%s' already exists", request.Name),
			map[string]string{"name": request.Name},
			[]string{"Choose a different name", "Update the existing template instead"},
		)
	}

	// Validate template file exists and is readable
	if !fileExists(request.Path) {
		return nil, createFileSystemError(
			"template file does not exist",
			request.Path,
			"create_template",
		)
	}

	// Create template object
	template := &Template{
		ID:          generateTemplateID(),
		Name:        request.Name,
		Description: request.Description,
		Version:     request.Version,
		Author:      request.Author,
		Type:        request.Type,
		Path:        request.Path,
		IsBuiltIn:   false,
		Variables:   request.Variables,
	}

	// Save to database
	if err := a.repo.Templates.Create(template); err != nil {
		return nil, createDatabaseError("failed to create template", "CREATE_TEMPLATE", err)
	}

	// Emit event
	a.emitEvent("template:created", map[string]interface{}{
		"templateId": template.ID,
		"name":       template.Name,
		"type":       template.Type,
	})

	return template, nil
}

// GetTemplate retrieves a template by ID
func (a *App) GetTemplate(id string) (*Template, error) {
	if id == "" {
		return nil, createValidationError("template ID is required", nil, nil)
	}

	template, err := a.repo.Templates.GetByID(id)
	if err != nil {
		return nil, err
	}

	return template, nil
}

// GetAllTemplates retrieves all templates
func (a *App) GetAllTemplates() ([]*Template, error) {
	templates, err := a.repo.Templates.GetAll()
	if err != nil {
		return nil, err
	}

	return templates, nil
}

// GetTemplatesByType retrieves templates by their type
func (a *App) GetTemplatesByType(templateType TemplateType) ([]*Template, error) {
	templates, err := a.repo.Templates.GetByType(templateType)
	if err != nil {
		return nil, err
	}

	return templates, nil
}

// UpdateTemplate updates an existing template
func (a *App) UpdateTemplate(id string, request UpdateTemplateRequest) (*Template, error) {
	if id == "" {
		return nil, createValidationError("template ID is required", nil, nil)
	}

	// Get existing template
	template, err := a.repo.Templates.GetByID(id)
	if err != nil {
		return nil, err
	}

	// Prevent updating built-in templates
	if template.IsBuiltIn {
		return nil, createValidationError(
			"cannot update built-in templates",
			map[string]string{"templateId": id},
			[]string{"Create a copy of the template", "Create a new custom template"},
		)
	}

	// Apply updates
	if request.Name != nil {
		// Check if new name conflicts with existing templates
		if *request.Name != template.Name {
			existing, err := a.repo.Templates.GetByName(*request.Name)
			if err == nil && existing != nil && existing.ID != id {
				return nil, createValidationError(
					fmt.Sprintf("template with name '%s' already exists", *request.Name),
					map[string]string{"name": *request.Name},
					[]string{"Choose a different name"},
				)
			}
		}
		template.Name = *request.Name
	}

	if request.Description != nil {
		template.Description = *request.Description
	}

	if request.Version != nil {
		// Validate semantic version format
		if err := validateSemanticVersion(*request.Version); err != nil {
			return nil, createValidationError(
				"invalid semantic version format",
				map[string]string{"version": *request.Version},
				[]string{"Use format like '1.0.0', '2.1.3', etc."},
			)
		}
		template.Version = *request.Version
	}

	if request.Author != nil {
		template.Author = *request.Author
	}

	if request.Type != nil {
		template.Type = *request.Type
	}

	if request.Path != nil {
		// Validate new template file exists and is readable
		if !fileExists(*request.Path) {
			return nil, createFileSystemError(
				"template file does not exist",
				*request.Path,
				"update_template",
			)
		}
		template.Path = *request.Path
	}

	if request.Variables != nil {
		template.Variables = *request.Variables
	}

	// Update in database
	if err := a.repo.Templates.Update(template); err != nil {
		return nil, err
	}

	// Emit event
	a.emitEvent("template:updated", map[string]interface{}{
		"templateId": template.ID,
		"name":       template.Name,
	})

	return template, nil
}

// DeleteTemplate deletes a template
func (a *App) DeleteTemplate(id string) error {
	if id == "" {
		return createValidationError("template ID is required", nil, nil)
	}

	// Get template to check if it's built-in
	template, err := a.repo.Templates.GetByID(id)
	if err != nil {
		return err
	}

	// Prevent deleting built-in templates
	if template.IsBuiltIn {
		return createValidationError(
			"cannot delete built-in templates",
			map[string]string{"templateId": id},
			[]string{"Built-in templates are read-only"},
		)
	}

	// Delete from database
	if err := a.repo.Templates.Delete(id); err != nil {
		return err
	}

	// Emit event
	a.emitEvent("template:deleted", map[string]interface{}{
		"templateId": template.ID,
		"name":       template.Name,
	})

	return nil
}

// SearchTemplates searches templates by query
func (a *App) SearchTemplates(query string) ([]*Template, error) {
	templates, err := a.repo.Templates.Search(query)
	if err != nil {
		return nil, err
	}

	return templates, nil
}

// GetBuiltInTemplates retrieves all built-in templates
func (a *App) GetBuiltInTemplates() ([]*Template, error) {
	templates, err := a.repo.Templates.GetBuiltInTemplates()
	if err != nil {
		return nil, err
	}

	return templates, nil
}

// ValidateTemplate validates a template file
func (a *App) ValidateTemplate(templatePath string) (*TemplateValidationResult, error) {
	if templatePath == "" {
		return nil, createValidationError("template path is required", nil, nil)
	}

	// Check if file exists
	if !fileExists(templatePath) {
		return &TemplateValidationResult{
			Valid: false,
			Errors: []TemplateError{
				{
					Type:     "file_not_found",
					Message:  "Template file does not exist",
					Severity: "error",
				},
			},
		}, nil
	}

	// Read template content
	content, err := os.ReadFile(templatePath)
	if err != nil {
		return &TemplateValidationResult{
			Valid: false,
			Errors: []TemplateError{
				{
					Type:     "file_read_error",
					Message:  fmt.Sprintf("Failed to read template file: %v", err),
					Severity: "error",
				},
			},
		}, nil
	}

	// Validate template syntax and structure
	result := &TemplateValidationResult{
		Valid:    true,
		Errors:   []TemplateError{},
		Warnings: []TemplateWarning{},
		Suggestions: []string{},
	}

	// Basic template validation
	contentStr := string(content)
	
	// Check for common template syntax issues
	if err := a.validateTemplateSyntax(contentStr, result); err != nil {
		result.Valid = false
	}

	// Check for variable usage
	a.validateTemplateVariables(contentStr, result)

	// Check for performance issues
	a.validateTemplatePerformance(contentStr, result)

	return result, nil
}

// DuplicateTemplate creates a copy of an existing template
func (a *App) DuplicateTemplate(id string, newName string) (*Template, error) {
	if id == "" {
		return nil, createValidationError("template ID is required", nil, nil)
	}

	if newName == "" {
		return nil, createValidationError("new template name is required", nil, nil)
	}

	// Get original template
	original, err := a.repo.Templates.GetByID(id)
	if err != nil {
		return nil, err
	}

	// Check if new name conflicts
	existing, err := a.repo.Templates.GetByName(newName)
	if err == nil && existing != nil {
		return nil, createValidationError(
			fmt.Sprintf("template with name '%s' already exists", newName),
			map[string]string{"name": newName},
			[]string{"Choose a different name"},
		)
	}

	// Create duplicate
	duplicate := &Template{
		ID:          generateTemplateID(),
		Name:        newName,
		Description: original.Description + " (Copy)",
		Version:     "1.0.0", // Reset version for copy
		Author:      original.Author,
		Type:        TemplateTypeCustom, // Copies are always custom
		Path:        original.Path,      // TODO: Should copy the file
		IsBuiltIn:   false,
		Variables:   original.Variables,
	}

	// Save to database
	if err := a.repo.Templates.Create(duplicate); err != nil {
		return nil, err
	}

	// Emit event
	a.emitEvent("template:duplicated", map[string]interface{}{
		"originalId": original.ID,
		"duplicateId": duplicate.ID,
		"name":       duplicate.Name,
	})

	return duplicate, nil
}

// Helper functions

func (a *App) validateCreateTemplateRequest(request *CreateTemplateRequest) error {
	if request.Name == "" {
		return createValidationError("template name is required", nil, nil)
	}

	if len(request.Name) > 100 {
		return createValidationError("template name is too long (max 100 characters)", nil, nil)
	}

	if request.Version == "" {
		return createValidationError("template version is required", nil, nil)
	}

	if err := validateSemanticVersion(request.Version); err != nil {
		return createValidationError(
			"invalid semantic version format",
			map[string]string{"version": request.Version},
			[]string{"Use format like '1.0.0', '2.1.3', etc."},
		)
	}

	if request.Path == "" {
		return createValidationError("template path is required", nil, nil)
	}

	return nil
}

func (a *App) validateTemplateSyntax(content string, result *TemplateValidationResult) error {
	// Check for balanced template delimiters
	openCount := strings.Count(content, "{{")
	closeCount := strings.Count(content, "}}")
	
	if openCount != closeCount {
		result.Errors = append(result.Errors, TemplateError{
			Type:     "syntax_error",
			Message:  "Unbalanced template delimiters",
			Severity: "error",
		})
		return fmt.Errorf("unbalanced template delimiters")
	}

	// Check for common template action syntax
	if strings.Contains(content, "{{") && !strings.Contains(content, "}}") {
		result.Errors = append(result.Errors, TemplateError{
			Type:     "syntax_error",
			Message:  "Template actions not properly closed",
			Severity: "error",
		})
		return fmt.Errorf("template actions not properly closed")
	}

	// Check for forbidden functions or actions
	forbiddenActions := []string{
		"{{.Exec",
		"{{.System",
		"{{.Command",
		"{{.File",
	}

	for _, forbidden := range forbiddenActions {
		if strings.Contains(content, forbidden) {
			result.Errors = append(result.Errors, TemplateError{
				Type:     "security_error",
				Message:  fmt.Sprintf("Forbidden template action detected: %s", forbidden),
				Severity: "error",
			})
			return fmt.Errorf("forbidden template action: %s", forbidden)
		}
	}

	return nil
}

func (a *App) validateTemplateVariables(content string, result *TemplateValidationResult) {
	// Extract variable references
	variableRefs := extractTemplateVariables(content)
	
	if len(variableRefs) == 0 {
		result.Warnings = append(result.Warnings, TemplateWarning{
			Type:       "no_variables",
			Message:    "Template doesn't use any variables",
			Suggestion: "Consider adding variables for customization",
		})
	}

	// Check for undefined variables (this would require template context)
	// For now, just provide suggestions
	if len(variableRefs) > 10 {
		result.Warnings = append(result.Warnings, TemplateWarning{
			Type:       "complex_template",
			Message:    "Template uses many variables",
			Suggestion: "Consider breaking into smaller templates",
		})
	}
}

func (a *App) validateTemplatePerformance(content string, result *TemplateValidationResult) {
	// Check for potential performance issues
	contentSize := len(content)
	
	if contentSize > 50000 { // 50KB
		result.Warnings = append(result.Warnings, TemplateWarning{
			Type:       "large_template",
			Message:    "Template file is very large",
			Suggestion: "Consider breaking into smaller templates or optimizing content",
		})
	}

	// Check for nested loops
	if strings.Contains(content, "{{range") && strings.Count(content, "{{range") > 1 {
		result.Warnings = append(result.Warnings, TemplateWarning{
			Type:       "nested_loops",
			Message:    "Template contains nested loops",
			Suggestion: "Nested loops can impact performance with large datasets",
		})
	}
}

func generateTemplateID() string {
	return uuid.New().String()
}

func validateSemanticVersion(version string) error {
	// Basic semantic version validation (major.minor.patch)
	parts := strings.Split(version, ".")
	if len(parts) != 3 {
		return fmt.Errorf("semantic version must have format major.minor.patch")
	}

	for _, part := range parts {
		if part == "" {
			return fmt.Errorf("version parts cannot be empty")
		}
		// Additional validation could check if parts are valid numbers
	}

	return nil
}

// compareVersions compares two semantic versions
// Returns: -1 if v1 < v2, 0 if v1 == v2, 1 if v1 > v2
func compareVersions(v1, v2 string) (int, error) {
	if err := validateSemanticVersion(v1); err != nil {
		return 0, fmt.Errorf("invalid version v1: %w", err)
	}
	if err := validateSemanticVersion(v2); err != nil {
		return 0, fmt.Errorf("invalid version v2: %w", err)
	}

	parts1 := strings.Split(v1, ".")
	parts2 := strings.Split(v2, ".")

	for i := 0; i < 3; i++ {
		// Convert to integers for proper comparison
		var num1, num2 int
		fmt.Sscanf(parts1[i], "%d", &num1)
		fmt.Sscanf(parts2[i], "%d", &num2)

		if num1 < num2 {
			return -1, nil
		} else if num1 > num2 {
			return 1, nil
		}
	}

	return 0, nil
}

// isNewerVersion checks if version1 is newer than version2
func isNewerVersion(version1, version2 string) (bool, error) {
	result, err := compareVersions(version1, version2)
	if err != nil {
		return false, err
	}
	return result > 0, nil
}

// incrementVersion increments a semantic version based on the type
func incrementVersion(version string, incrementType string) (string, error) {
	if err := validateSemanticVersion(version); err != nil {
		return "", err
	}

	parts := strings.Split(version, ".")
	var major, minor, patch int
	fmt.Sscanf(parts[0], "%d", &major)
	fmt.Sscanf(parts[1], "%d", &minor)
	fmt.Sscanf(parts[2], "%d", &patch)

	switch incrementType {
	case "major":
		major++
		minor = 0
		patch = 0
	case "minor":
		minor++
		patch = 0
	case "patch":
		patch++
	default:
		return "", fmt.Errorf("invalid increment type: %s (use major, minor, or patch)", incrementType)
	}

	return fmt.Sprintf("%d.%d.%d", major, minor, patch), nil
}

func extractTemplateVariables(content string) []string {
	var variables []string
	
	// Simple extraction of {{.Variable}} patterns
	lines := strings.Split(content, "\n")
	for _, line := range lines {
		if strings.Contains(line, "{{.") {
			// This is a simplified extraction
			// A more robust implementation would use proper template parsing
			start := strings.Index(line, "{{.")
			if start != -1 {
				end := strings.Index(line[start:], "}}")
				if end != -1 {
					variable := line[start+3 : start+end]
					if !contains(variables, variable) {
						variables = append(variables, variable)
					}
				}
			}
		}
	}
	
	return variables
}

func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

// InstallBuiltInTemplates installs the default built-in templates
func (a *App) InstallBuiltInTemplates() error {
	builtInTemplates := []Template{
		{
			ID:          "builtin-go-server",
			Name:        "Go MCP Server",
			Description: "Default Go template for generating MCP servers",
			Version:     "1.0.0",
			Author:      "MCPWeaver",
			Type:        TemplateTypeDefault,
			Path:        "templates/server.go.tmpl",
			IsBuiltIn:   true,
			Variables: []TemplateVariable{
				{
					Name:         "PackageName",
					Description:  "Name of the Go package",
					Type:         "string",
					DefaultValue: "mcp-server",
					Required:     true,
				},
				{
					Name:         "APITitle",
					Description:  "Title of the API",
					Type:         "string",
					Required:     true,
				},
				{
					Name:         "BaseURL",
					Description:  "Base URL for API calls",
					Type:         "string",
					Required:     true,
				},
			},
		},
		{
			ID:          "builtin-go-mod",
			Name:        "Go Module",
			Description: "Go module template for MCP servers",
			Version:     "1.0.0",
			Author:      "MCPWeaver",
			Type:        TemplateTypeDefault,
			Path:        "templates/go.mod.tmpl",
			IsBuiltIn:   true,
			Variables: []TemplateVariable{
				{
					Name:         "PackageName",
					Description:  "Name of the Go package",
					Type:         "string",
					DefaultValue: "mcp-server",
					Required:     true,
				},
			},
		},
	}

	for _, template := range builtInTemplates {
		// Check if already exists
		existing, err := a.repo.Templates.GetByID(template.ID)
		if err == nil && existing != nil {
			continue // Skip if already installed
		}

		// Create the built-in template
		if err := a.repo.Templates.Create(&template); err != nil {
			return createDatabaseError(
				fmt.Sprintf("failed to install built-in template '%s'", template.Name),
				"INSTALL_BUILTIN_TEMPLATE",
				err,
			)
		}
	}

	return nil
}

// GetTemplateStatistics returns statistics about templates
func (a *App) GetTemplateStatistics() (map[string]interface{}, error) {
	stats := make(map[string]interface{})

	// Total count
	total, err := a.repo.Templates.Count()
	if err != nil {
		return nil, err
	}
	stats["total"] = total

	// Count by type
	defaultCount, err := a.repo.Templates.CountByType(TemplateTypeDefault)
	if err != nil {
		return nil, err
	}
	stats["default"] = defaultCount

	customCount, err := a.repo.Templates.CountByType(TemplateTypeCustom)
	if err != nil {
		return nil, err
	}
	stats["custom"] = customCount

	pluginCount, err := a.repo.Templates.CountByType(TemplateTypePlugin)
	if err != nil {
		return nil, err
	}
	stats["plugin"] = pluginCount

	// Built-in count
	builtInTemplates, err := a.repo.Templates.GetBuiltInTemplates()
	if err != nil {
		return nil, err
	}
	stats["builtIn"] = len(builtInTemplates)

	return stats, nil
}

// TestTemplate tests a template with provided data
func (a *App) TestTemplate(request TemplateTestRequest) (*TemplateTestResult, error) {
	if request.TemplateID == "" {
		return nil, createValidationError("template ID is required", nil, nil)
	}

	// Get template
	tmplRecord, err := a.repo.Templates.GetByID(request.TemplateID)
	if err != nil {
		return nil, err
	}

	startTime := time.Now()

	// Prepare result
	result := &TemplateTestResult{
		Success: true,
		Errors:  []TemplateError{},
		Warnings: []TemplateWarning{},
	}

	// Load and parse template
	content, err := os.ReadFile(tmplRecord.Path)
	if err != nil {
		result.Success = false
		result.Errors = append(result.Errors, TemplateError{
			Type:     "file_read_error",
			Message:  fmt.Sprintf("Failed to read template file: %v", err),
			Severity: "error",
		})
		return result, nil
	}

	// Parse template with custom functions
	tmpl, err := template.New(tmplRecord.Name).Funcs(template.FuncMap{
		"title": func(s string) string {
			if s == "" {
				return s
			}
			return strings.Title(s)
		},
		"upper": strings.ToUpper,
		"lower": strings.ToLower,
	}).Parse(string(content))

	if err != nil {
		result.Success = false
		result.Errors = append(result.Errors, TemplateError{
			Type:     "parse_error",
			Message:  fmt.Sprintf("Failed to parse template: %v", err),
			Severity: "error",
		})
		return result, nil
	}

	// Execute template
	var output strings.Builder
	err = tmpl.Execute(&output, request.TestData)
	if err != nil {
		result.Success = false
		result.Errors = append(result.Errors, TemplateError{
			Type:     "execution_error",
			Message:  fmt.Sprintf("Template execution failed: %v", err),
			Severity: "error",
		})
		return result, nil
	}

	result.Output = output.String()
	executionTime := time.Since(startTime)

	// Measure performance if requested
	if request.TestOptions.MeasurePerformance {
		result.Performance = &TemplatePerformance{
			RenderTime:  executionTime,
			MemoryUsage: int64(len(result.Output)),
			Complexity:  calculateTemplateComplexity(string(content)),
			CacheHit:    false,
		}
	}

	// Generate report if requested
	if request.TestOptions.GenerateReport {
		result.Report = &TemplateTestReport{
			TemplateID:      tmplRecord.ID,
			TestExecutedAt:  startTime,
			ExecutionTime:   executionTime,
			OutputSize:      int64(len(result.Output)),
			VariablesUsed:   extractUsedVariables(string(content), request.TestData),
			FunctionsUsed:   extractUsedFunctions(string(content)),
			Recommendations: generateTemplateRecommendations(string(content), request.TestData),
		}
	}

	// Validate output if requested
	if request.TestOptions.ValidateOutput {
		a.validateTemplateOutput(result.Output, result)
	}

	return result, nil
}

// ValidateTemplateAdvanced performs advanced template validation
func (a *App) ValidateTemplateAdvanced(templateID string) (*TemplateValidationResult, error) {
	if templateID == "" {
		return nil, createValidationError("template ID is required", nil, nil)
	}

	// Get template
	tmplRecord, err := a.repo.Templates.GetByID(templateID)
	if err != nil {
		return nil, err
	}

	// Use the internal validation method
	return a.ValidateTemplate(tmplRecord.Path)
}

// CompareTemplates compares two templates and returns differences
func (a *App) CompareTemplates(templateID1, templateID2 string) (map[string]interface{}, error) {
	if templateID1 == "" || templateID2 == "" {
		return nil, createValidationError("both template IDs are required", nil, nil)
	}

	// Get both templates
	template1, err := a.repo.Templates.GetByID(templateID1)
	if err != nil {
		return nil, err
	}

	template2, err := a.repo.Templates.GetByID(templateID2)
	if err != nil {
		return nil, err
	}

	comparison := make(map[string]interface{})

	// Compare basic properties
	comparison["name"] = map[string]string{
		"template1": template1.Name,
		"template2": template2.Name,
		"same":      fmt.Sprintf("%t", template1.Name == template2.Name),
	}

	comparison["version"] = map[string]string{
		"template1": template1.Version,
		"template2": template2.Version,
		"same":      fmt.Sprintf("%t", template1.Version == template2.Version),
	}

	comparison["type"] = map[string]string{
		"template1": string(template1.Type),
		"template2": string(template2.Type),
		"same":      fmt.Sprintf("%t", template1.Type == template2.Type),
	}

	// Compare variables
	comparison["variables"] = a.compareTemplateVariables(template1.Variables, template2.Variables)

	// Compare file content if accessible
	if fileExists(template1.Path) && fileExists(template2.Path) {
		content1, err1 := os.ReadFile(template1.Path)
		content2, err2 := os.ReadFile(template2.Path)

		if err1 == nil && err2 == nil {
			comparison["content"] = map[string]interface{}{
				"size1":     len(content1),
				"size2":     len(content2),
				"identical": string(content1) == string(content2),
				"similarity": calculateSimilarity(string(content1), string(content2)),
			}
		}
	}

	return comparison, nil
}

// Helper functions for template testing

func calculateTemplateComplexity(content string) string {
	score := 0
	score += strings.Count(content, "{{range") * 3
	score += strings.Count(content, "{{if") * 2
	score += strings.Count(content, "{{with") * 2
	score += len(strings.Split(content, "\n")) / 10

	if score < 10 {
		return "low"
	} else if score < 25 {
		return "medium"
	} else {
		return "high"
	}
}

func extractUsedVariables(content string, testData map[string]interface{}) []string {
	var variables []string
	
	for key := range testData {
		if strings.Contains(content, "{{."+key) {
			variables = append(variables, key)
		}
	}
	
	return variables
}

func extractUsedFunctions(content string) []string {
	var functions []string
	
	commonFunctions := []string{"title", "upper", "lower", "trim", "replace", "printf", "range", "if", "with"}
	
	for _, fn := range commonFunctions {
		if strings.Contains(content, fn) {
			functions = append(functions, fn)
		}
	}
	
	return functions
}

func generateTemplateRecommendations(content string, testData map[string]interface{}) []string {
	var recommendations []string
	
	// Check for missing error handling
	if strings.Contains(content, "{{range") && !strings.Contains(content, "{{else}}") {
		recommendations = append(recommendations, "Consider adding {{else}} clause to handle empty ranges")
	}
	
	// Check for hardcoded values
	lines := strings.Split(content, "\n")
	for _, line := range lines {
		if strings.Contains(line, "http://") || strings.Contains(line, "https://") {
			if !strings.Contains(line, "{{") {
				recommendations = append(recommendations, "Consider parameterizing hardcoded URLs")
				break
			}
		}
	}
	
	// Check variable usage
	if len(testData) > 5 && strings.Count(content, "{{.") < len(testData) {
		recommendations = append(recommendations, "Template uses only a subset of available variables")
	}
	
	return recommendations
}

func (a *App) validateTemplateOutput(output string, result *TemplateTestResult) {
	// Basic output validation
	if output == "" {
		result.Warnings = append(result.Warnings, TemplateWarning{
			Type:       "empty_output",
			Message:    "Template generated empty output",
			Suggestion: "Check if template variables are properly set",
		})
	}
	
	// Check for common issues
	if strings.Contains(output, "<no value>") {
		result.Warnings = append(result.Warnings, TemplateWarning{
			Type:       "missing_value",
			Message:    "Template output contains '<no value>' indicating missing variables",
			Suggestion: "Ensure all required variables are provided",
		})
	}
	
	// Check for unclosed template actions
	if strings.Contains(output, "{{") || strings.Contains(output, "}}") {
		result.Warnings = append(result.Warnings, TemplateWarning{
			Type:       "unprocessed_template",
			Message:    "Output contains unprocessed template syntax",
			Suggestion: "Check template syntax and variable names",
		})
	}
}

func (a *App) compareTemplateVariables(vars1, vars2 []TemplateVariable) map[string]interface{} {
	comparison := make(map[string]interface{})
	
	// Create maps for easier comparison
	varMap1 := make(map[string]TemplateVariable)
	varMap2 := make(map[string]TemplateVariable)
	
	for _, v := range vars1 {
		varMap1[v.Name] = v
	}
	for _, v := range vars2 {
		varMap2[v.Name] = v
	}
	
	// Find common, unique to each
	var common, unique1, unique2 []string
	
	for name := range varMap1 {
		if _, exists := varMap2[name]; exists {
			common = append(common, name)
		} else {
			unique1 = append(unique1, name)
		}
	}
	
	for name := range varMap2 {
		if _, exists := varMap1[name]; !exists {
			unique2 = append(unique2, name)
		}
	}
	
	comparison["common"] = common
	comparison["unique_to_template1"] = unique1
	comparison["unique_to_template2"] = unique2
	comparison["total_template1"] = len(vars1)
	comparison["total_template2"] = len(vars2)
	
	return comparison
}

func calculateSimilarity(text1, text2 string) float64 {
	// Simple similarity calculation based on common characters
	if len(text1) == 0 && len(text2) == 0 {
		return 1.0
	}
	
	if len(text1) == 0 || len(text2) == 0 {
		return 0.0
	}
	
	// Count common characters (very basic implementation)
	common := 0
	total := len(text1) + len(text2)
	
	shorter := text1
	longer := text2
	if len(text2) < len(text1) {
		shorter = text2
		longer = text1
	}
	
	for _, char := range shorter {
		if strings.ContainsRune(longer, char) {
			common++
		}
	}
	
	return float64(common*2) / float64(total)
}

// ImportTemplate imports a template from various sources
func (a *App) ImportTemplate(request TemplateImportRequest) (*Template, error) {
	if request.Source == "" {
		return nil, createValidationError("import source is required", nil, nil)
	}

	switch request.Source {
	case "file":
		return a.importTemplateFromFile(request)
	case "url":
		return a.importTemplateFromURL(request)
	case "marketplace":
		return a.importTemplateFromMarketplace(request)
	default:
		return nil, createValidationError(
			fmt.Sprintf("unsupported import source: %s", request.Source),
			map[string]string{"source": request.Source},
			[]string{"Use 'file', 'url', or 'marketplace'"},
		)
	}
}

// importTemplateFromFile imports a template from a local file
func (a *App) importTemplateFromFile(request TemplateImportRequest) (*Template, error) {
	if request.Path == "" {
		return nil, createValidationError("file path is required for file import", nil, nil)
	}

	// Check if file exists
	if !fileExists(request.Path) {
		return nil, createFileSystemError("template file does not exist", request.Path, "import_template")
	}

	// Determine if it's a single template file or a template package
	if filepath.Ext(request.Path) == ".zip" {
		return a.importTemplatePackage(request.Path, request.ImportOptions)
	}

	// Import single template file
	return a.importSingleTemplateFile(request.Path, request.ImportOptions)
}

// importTemplateFromURL imports a template from a URL
func (a *App) importTemplateFromURL(request TemplateImportRequest) (*Template, error) {
	if request.URL == "" {
		return nil, createValidationError("URL is required for URL import", nil, nil)
	}

	// Download template file
	tempFile, err := a.downloadTemplateFromURL(request.URL)
	if err != nil {
		return nil, createNetworkError("failed to download template", map[string]string{"url": request.URL})
	}
	defer os.Remove(tempFile)

	// Import the downloaded file
	importRequest := request
	importRequest.Path = tempFile
	importRequest.Source = "file"
	
	return a.importTemplateFromFile(importRequest)
}

// importTemplateFromMarketplace imports a template from marketplace
func (a *App) importTemplateFromMarketplace(request TemplateImportRequest) (*Template, error) {
	if request.MarketplaceID == "" {
		return nil, createValidationError("marketplace ID is required for marketplace import", nil, nil)
	}

	// This would integrate with a template marketplace API
	// For now, return an error indicating it's not implemented
	return nil, createValidationError(
		"marketplace import not yet implemented",
		map[string]string{"marketplaceId": request.MarketplaceID},
		[]string{"Use file or URL import for now"},
	)
}

// importSingleTemplateFile imports a single template file
func (a *App) importSingleTemplateFile(filePath string, options TemplateImportOptions) (*Template, error) {
	// Generate template metadata
	templateName := strings.TrimSuffix(filepath.Base(filePath), filepath.Ext(filePath))
	
	// Check if template with same name exists
	if !options.OverwriteExisting {
		existing, err := a.repo.Templates.GetByName(templateName)
		if err == nil && existing != nil {
			return nil, createValidationError(
				fmt.Sprintf("template with name '%s' already exists", templateName),
				map[string]string{"name": templateName},
				[]string{"Use overwriteExisting option", "Rename the template"},
			)
		}
	}

	// Create template directory in templates folder
	templatesDir := "templates/imported"
	if err := ensureDir(templatesDir); err != nil {
		return nil, createFileSystemError("failed to create templates directory", templatesDir, "import_template")
	}

	// Copy template file to templates directory
	newTemplatePath := filepath.Join(templatesDir, filepath.Base(filePath))
	if err := copyFile(filePath, newTemplatePath); err != nil {
		return nil, createFileSystemError("failed to copy template file", newTemplatePath, "import_template")
	}

	// Create template record
	template := &Template{
		ID:          generateTemplateID(),
		Name:        templateName,
		Description: fmt.Sprintf("Imported from %s", filePath),
		Version:     "1.0.0",
		Author:      "Imported",
		Type:        options.TargetType,
		Path:        newTemplatePath,
		IsBuiltIn:   false,
		Variables:   []TemplateVariable{}, // TODO: Extract variables from template
	}

	if template.Type == "" {
		template.Type = TemplateTypeCustom
	}

	// Validate template if not in validate-only mode
	if !options.ValidateOnly {
		validationResult, err := a.ValidateTemplate(newTemplatePath)
		if err != nil {
			return nil, err
		}

		if !validationResult.Valid {
			// Clean up copied file
			os.Remove(newTemplatePath)
			return nil, createValidationError(
				"imported template failed validation",
				map[string]string{"path": filePath},
				[]string{"Check template syntax", "Fix validation errors"},
			)
		}
	}

	// Save to database if not in validate-only mode
	if !options.ValidateOnly {
		if err := a.repo.Templates.Create(template); err != nil {
			// Clean up copied file
			os.Remove(newTemplatePath)
			return nil, err
		}

		// Emit event
		a.emitEvent("template:imported", map[string]interface{}{
			"templateId": template.ID,
			"name":       template.Name,
			"source":     "file",
		})
	}

	return template, nil
}

// importTemplatePackage imports a template package (ZIP file)
func (a *App) importTemplatePackage(zipPath string, options TemplateImportOptions) (*Template, error) {
	// Open ZIP file
	reader, err := zip.OpenReader(zipPath)
	if err != nil {
		return nil, createFileSystemError("failed to open template package", zipPath, "import_template")
	}
	defer reader.Close()

	// Look for template manifest
	var manifestFile *zip.File
	var templateFiles []*zip.File
	
	for _, file := range reader.File {
		if file.Name == "template.json" || file.Name == "manifest.json" {
			manifestFile = file
		} else if strings.HasSuffix(file.Name, ".tmpl") || strings.HasSuffix(file.Name, ".template") {
			templateFiles = append(templateFiles, file)
		}
	}

	// Parse manifest if available
	var templateMetadata Template
	if manifestFile != nil {
		manifestReader, err := manifestFile.Open()
		if err != nil {
			return nil, createFileSystemError("failed to read manifest file", "manifest.json", "import_template")
		}
		defer manifestReader.Close()

		if err := json.NewDecoder(manifestReader).Decode(&templateMetadata); err != nil {
			return nil, createValidationError("invalid manifest format", nil, []string{"Check manifest.json syntax"})
		}
	} else {
		// Generate default metadata
		templateMetadata = Template{
			Name:        strings.TrimSuffix(filepath.Base(zipPath), ".zip"),
			Description: "Imported template package",
			Version:     "1.0.0",
			Author:      "Imported",
			Type:        options.TargetType,
		}
	}

	if templateMetadata.Type == "" {
		templateMetadata.Type = TemplateTypeCustom
	}

	// Create extraction directory
	extractDir := filepath.Join("templates/imported", templateMetadata.Name)
	if err := ensureDir(extractDir); err != nil {
		return nil, createFileSystemError("failed to create extraction directory", extractDir, "import_template")
	}

	// Extract template files
	var mainTemplatePath string
	for _, file := range templateFiles {
		extractPath := filepath.Join(extractDir, file.Name)
		
		if err := extractZipFile(file, extractPath); err != nil {
			return nil, createFileSystemError("failed to extract template file", file.Name, "import_template")
		}

		// Use first template file as main template
		if mainTemplatePath == "" {
			mainTemplatePath = extractPath
		}
	}

	if mainTemplatePath == "" {
		return nil, createValidationError("no template files found in package", nil, nil)
	}

	// Update template metadata
	templateMetadata.ID = generateTemplateID()
	templateMetadata.Path = mainTemplatePath
	templateMetadata.IsBuiltIn = false

	// Save to database if not in validate-only mode
	if !options.ValidateOnly {
		if err := a.repo.Templates.Create(&templateMetadata); err != nil {
			// Clean up extracted files
			os.RemoveAll(extractDir)
			return nil, err
		}

		// Emit event
		a.emitEvent("template:imported", map[string]interface{}{
			"templateId": templateMetadata.ID,
			"name":       templateMetadata.Name,
			"source":     "package",
		})
	}

	return &templateMetadata, nil
}

// ExportTemplate exports a template to various formats
func (a *App) ExportTemplate(request TemplateExportRequest) (*ExportResult, error) {
	if request.TemplateID == "" {
		return nil, createValidationError("template ID is required", nil, nil)
	}

	if request.TargetPath == "" {
		return nil, createValidationError("target path is required", nil, nil)
	}

	// Get template
	template, err := a.repo.Templates.GetByID(request.TemplateID)
	if err != nil {
		return nil, err
	}

	switch request.Format {
	case "zip":
		return a.exportTemplateAsZip(template, request)
	case "tar":
		return a.exportTemplateAsTar(template, request)
	case "single":
		return a.exportTemplateAsSingle(template, request)
	default:
		return nil, createValidationError(
			fmt.Sprintf("unsupported export format: %s", request.Format),
			map[string]string{"format": request.Format},
			[]string{"Use 'zip', 'tar', or 'single'"},
		)
	}
}

// exportTemplateAsZip exports template as a ZIP package
func (a *App) exportTemplateAsZip(template *Template, request TemplateExportRequest) (*ExportResult, error) {
	// Create target directory
	if err := ensureDir(filepath.Dir(request.TargetPath)); err != nil {
		return nil, createFileSystemError("failed to create target directory", request.TargetPath, "export_template")
	}

	// Create ZIP file
	zipFile, err := os.Create(request.TargetPath)
	if err != nil {
		return nil, createFileSystemError("failed to create export file", request.TargetPath, "export_template")
	}
	defer zipFile.Close()

	zipWriter := zip.NewWriter(zipFile)
	defer zipWriter.Close()

	var exportedFiles []ExportedFile
	var totalSize int64

	// Add template file
	if err := a.addFileToZip(zipWriter, template.Path, "template.tmpl"); err != nil {
		return nil, err
	}

	fileInfo, _ := os.Stat(template.Path)
	exportedFiles = append(exportedFiles, ExportedFile{
		Name:         "template.tmpl",
		Path:         template.Path,
		Size:         fileInfo.Size(),
		ModifiedTime: fileInfo.ModTime(),
	})
	totalSize += fileInfo.Size()

	// Add manifest
	manifestData := map[string]interface{}{
		"id":          template.ID,
		"name":        template.Name,
		"description": template.Description,
		"version":     template.Version,
		"author":      template.Author,
		"type":        template.Type,
		"variables":   template.Variables,
		"exportedAt":  time.Now(),
		"exportedBy":  "MCPWeaver",
	}

	manifestJSON, err := json.MarshalIndent(manifestData, "", "  ")
	if err != nil {
		return nil, createInternalError("failed to create manifest", err)
	}

	manifestWriter, err := zipWriter.Create("manifest.json")
	if err != nil {
		return nil, createInternalError("failed to create manifest in ZIP", err)
	}
	
	if _, err := manifestWriter.Write(manifestJSON); err != nil {
		return nil, createInternalError("failed to write manifest", err)
	}

	exportedFiles = append(exportedFiles, ExportedFile{
		Name:         "manifest.json",
		Path:         "generated",
		Size:         int64(len(manifestJSON)),
		ModifiedTime: time.Now(),
	})
	totalSize += int64(len(manifestJSON))

	// Add documentation if requested
	if request.ExportOptions.IncludeDocumentation {
		readme := a.generateTemplateDocumentation(template)
		readmeWriter, err := zipWriter.Create("README.md")
		if err == nil {
			readmeWriter.Write([]byte(readme))
			exportedFiles = append(exportedFiles, ExportedFile{
				Name:         "README.md",
				Path:         "generated",
				Size:         int64(len(readme)),
				ModifiedTime: time.Now(),
			})
			totalSize += int64(len(readme))
		}
	}

	// Create export result
	result := &ExportResult{
		ProjectID:     "", // Not applicable for template export
		ProjectName:   template.Name,
		TargetDir:     filepath.Dir(request.TargetPath),
		ExportedFiles: exportedFiles,
		TotalFiles:    len(exportedFiles),
		TotalSize:     totalSize,
		ExportedAt:    time.Now(),
	}

	// Emit event
	a.emitEvent("template:exported", map[string]interface{}{
		"templateId": template.ID,
		"format":     "zip",
		"targetPath": request.TargetPath,
	})

	return result, nil
}

// exportTemplateAsSingle exports template as a single file
func (a *App) exportTemplateAsSingle(template *Template, request TemplateExportRequest) (*ExportResult, error) {
	// Create target directory
	if err := ensureDir(filepath.Dir(request.TargetPath)); err != nil {
		return nil, createFileSystemError("failed to create target directory", request.TargetPath, "export_template")
	}

	// Copy template file
	if err := copyFile(template.Path, request.TargetPath); err != nil {
		return nil, createFileSystemError("failed to copy template file", request.TargetPath, "export_template")
	}

	fileInfo, err := os.Stat(request.TargetPath)
	if err != nil {
		return nil, createFileSystemError("failed to get exported file info", request.TargetPath, "export_template")
	}

	// Create export result
	result := &ExportResult{
		ProjectID:   "", // Not applicable for template export
		ProjectName: template.Name,
		TargetDir:   filepath.Dir(request.TargetPath),
		ExportedFiles: []ExportedFile{
			{
				Name:         filepath.Base(request.TargetPath),
				Path:         request.TargetPath,
				Size:         fileInfo.Size(),
				ModifiedTime: fileInfo.ModTime(),
			},
		},
		TotalFiles: 1,
		TotalSize:  fileInfo.Size(),
		ExportedAt: time.Now(),
	}

	// Emit event
	a.emitEvent("template:exported", map[string]interface{}{
		"templateId": template.ID,
		"format":     "single",
		"targetPath": request.TargetPath,
	})

	return result, nil
}

// exportTemplateAsTar exports template as TAR archive (placeholder)
func (a *App) exportTemplateAsTar(template *Template, request TemplateExportRequest) (*ExportResult, error) {
	// TAR export not implemented yet
	return nil, createValidationError("TAR export format not yet implemented", nil, []string{"Use ZIP or single file export"})
}

// Helper functions for import/export

func (a *App) downloadTemplateFromURL(url string) (string, error) {
	// Create temporary file
	tempFile, err := os.CreateTemp("", "template-download-*.tmp")
	if err != nil {
		return "", err
	}
	defer tempFile.Close()

	// Download file
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to download: HTTP %d", resp.StatusCode)
	}

	// Copy response to temporary file
	_, err = io.Copy(tempFile, resp.Body)
	if err != nil {
		os.Remove(tempFile.Name())
		return "", err
	}

	return tempFile.Name(), nil
}

func (a *App) addFileToZip(zipWriter *zip.Writer, sourcePath, targetName string) error {
	sourceFile, err := os.Open(sourcePath)
	if err != nil {
		return err
	}
	defer sourceFile.Close()

	zipFile, err := zipWriter.Create(targetName)
	if err != nil {
		return err
	}

	_, err = io.Copy(zipFile, sourceFile)
	return err
}

func extractZipFile(file *zip.File, targetPath string) error {
	reader, err := file.Open()
	if err != nil {
		return err
	}
	defer reader.Close()

	// Create target directory
	if err := ensureDir(filepath.Dir(targetPath)); err != nil {
		return err
	}

	// Create target file
	targetFile, err := os.Create(targetPath)
	if err != nil {
		return err
	}
	defer targetFile.Close()

	// Copy content
	_, err = io.Copy(targetFile, reader)
	return err
}

func copyFile(src, dst string) error {
	sourceFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer sourceFile.Close()

	destFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer destFile.Close()

	_, err = io.Copy(destFile, sourceFile)
	return err
}

func (a *App) generateTemplateDocumentation(template *Template) string {
	doc := fmt.Sprintf(`# %s

%s

## Version
%s

## Author
%s

## Type
%s

## Variables

`, template.Name, template.Description, template.Version, template.Author, template.Type)

	if len(template.Variables) == 0 {
		doc += "No variables defined for this template.\n"
	} else {
		for _, variable := range template.Variables {
			doc += fmt.Sprintf("### %s\n", variable.Name)
			doc += fmt.Sprintf("- **Type**: %s\n", variable.Type)
			doc += fmt.Sprintf("- **Description**: %s\n", variable.Description)
			doc += fmt.Sprintf("- **Required**: %t\n", variable.Required)
			if variable.DefaultValue != "" {
				doc += fmt.Sprintf("- **Default**: %s\n", variable.DefaultValue)
			}
			doc += "\n"
		}
	}

	doc += fmt.Sprintf(`
## Usage

This template can be used to generate MCP servers. Import it into MCPWeaver and configure the variables as needed.

## Generated by MCPWeaver

Exported on: %s
`, time.Now().Format("2006-01-02 15:04:05"))

	return doc
}

// CreateTemplateVersion creates a new version of an existing template
func (a *App) CreateTemplateVersion(templateID string, incrementType string, changes string) (*Template, error) {
	if templateID == "" {
		return nil, createValidationError("template ID is required", nil, nil)
	}

	if incrementType == "" {
		incrementType = "patch" // Default to patch increment
	}

	// Get current template
	currentTemplate, err := a.repo.Templates.GetByID(templateID)
	if err != nil {
		return nil, err
	}

	// Prevent versioning built-in templates
	if currentTemplate.IsBuiltIn {
		return nil, createValidationError(
			"cannot create versions of built-in templates",
			map[string]string{"templateId": templateID},
			[]string{"Create a copy of the template first"},
		)
	}

	// Increment version
	newVersion, err := incrementVersion(currentTemplate.Version, incrementType)
	if err != nil {
		return nil, createValidationError(
			fmt.Sprintf("failed to increment version: %v", err),
			map[string]string{"currentVersion": currentTemplate.Version, "incrementType": incrementType},
			[]string{"Check version format and increment type"},
		)
	}

	// Create new template version
	newTemplate := &Template{
		ID:          generateTemplateID(),
		Name:        currentTemplate.Name,
		Description: currentTemplate.Description,
		Version:     newVersion,
		Author:      currentTemplate.Author,
		Type:        currentTemplate.Type,
		Path:        currentTemplate.Path, // Will be updated if path changes
		IsBuiltIn:   false,
		Variables:   currentTemplate.Variables,
	}

	// Update description with changes if provided
	if changes != "" {
		newTemplate.Description = fmt.Sprintf("%s\n\nChanges in v%s: %s", 
			currentTemplate.Description, newVersion, changes)
	}

	// Save new version to database
	if err := a.repo.Templates.Create(newTemplate); err != nil {
		return nil, err
	}

	// Emit event
	a.emitEvent("template:version_created", map[string]interface{}{
		"originalTemplateId": templateID,
		"newTemplateId":      newTemplate.ID,
		"oldVersion":         currentTemplate.Version,
		"newVersion":         newVersion,
		"incrementType":      incrementType,
	})

	return newTemplate, nil
}

// GetTemplateVersions gets all versions of a template by name
func (a *App) GetTemplateVersions(templateName string) ([]*Template, error) {
	if templateName == "" {
		return nil, createValidationError("template name is required", nil, nil)
	}

	// Get all templates with the same name
	allTemplates, err := a.repo.Templates.GetAll()
	if err != nil {
		return nil, err
	}

	var versions []*Template
	for _, template := range allTemplates {
		if template.Name == templateName {
			versions = append(versions, template)
		}
	}

	// Sort by version (newest first)
	for i := 0; i < len(versions)-1; i++ {
		for j := i + 1; j < len(versions); j++ {
			isNewer, err := isNewerVersion(versions[i].Version, versions[j].Version)
			if err == nil && !isNewer {
				// Swap if versions[j] is newer
				versions[i], versions[j] = versions[j], versions[i]
			}
		}
	}

	return versions, nil
}

// GetLatestTemplateVersion gets the latest version of a template by name
func (a *App) GetLatestTemplateVersion(templateName string) (*Template, error) {
	versions, err := a.GetTemplateVersions(templateName)
	if err != nil {
		return nil, err
	}

	if len(versions) == 0 {
		return nil, createValidationError(
			fmt.Sprintf("no templates found with name '%s'", templateName),
			map[string]string{"name": templateName},
			nil,
		)
	}

	return versions[0], nil // First item is the latest due to sorting
}

// UpdateTemplateToVersion updates a template to use a specific version
func (a *App) UpdateTemplateToVersion(templateID string, targetVersion string) (*Template, error) {
	if templateID == "" {
		return nil, createValidationError("template ID is required", nil, nil)
	}

	if targetVersion == "" {
		return nil, createValidationError("target version is required", nil, nil)
	}

	// Get current template
	currentTemplate, err := a.repo.Templates.GetByID(templateID)
	if err != nil {
		return nil, err
	}

	// Find the target version
	versions, err := a.GetTemplateVersions(currentTemplate.Name)
	if err != nil {
		return nil, err
	}

	var targetTemplate *Template
	for _, version := range versions {
		if version.Version == targetVersion {
			targetTemplate = version
			break
		}
	}

	if targetTemplate == nil {
		return nil, createValidationError(
			fmt.Sprintf("version '%s' not found for template '%s'", targetVersion, currentTemplate.Name),
			map[string]string{"templateName": currentTemplate.Name, "targetVersion": targetVersion},
			[]string{"Check available versions"},
		)
	}

	// Update current template to match target version
	updateRequest := UpdateTemplateRequest{
		Version:     &targetTemplate.Version,
		Description: &targetTemplate.Description,
		Path:        &targetTemplate.Path,
		Variables:   &targetTemplate.Variables,
	}

	return a.UpdateTemplate(templateID, updateRequest)
}

// GetTemplateVersionHistory gets the version history for a template
func (a *App) GetTemplateVersionHistory(templateName string) (map[string]interface{}, error) {
	versions, err := a.GetTemplateVersions(templateName)
	if err != nil {
		return nil, err
	}

	if len(versions) == 0 {
		return nil, createValidationError(
			fmt.Sprintf("no templates found with name '%s'", templateName),
			map[string]string{"name": templateName},
			nil,
		)
	}

	history := make(map[string]interface{})
	history["templateName"] = templateName
	history["totalVersions"] = len(versions)
	history["latestVersion"] = versions[0].Version
	history["oldestVersion"] = versions[len(versions)-1].Version

	// Create version details
	var versionDetails []map[string]interface{}
	for _, version := range versions {
		detail := map[string]interface{}{
			"id":          version.ID,
			"version":     version.Version,
			"author":      version.Author,
			"createdAt":   version.CreatedAt,
			"updatedAt":   version.UpdatedAt,
			"description": version.Description,
		}
		versionDetails = append(versionDetails, detail)
	}
	history["versions"] = versionDetails

	// Calculate version statistics
	var majorVersions, minorVersions, patchVersions int
	for _, version := range versions {
		parts := strings.Split(version.Version, ".")
		if len(parts) == 3 {
			var major, minor, patch int
			fmt.Sscanf(parts[0], "%d", &major)
			fmt.Sscanf(parts[1], "%d", &minor)
			fmt.Sscanf(parts[2], "%d", &patch)

			if major > 0 {
				majorVersions++
			}
			if minor > 0 {
				minorVersions++
			}
			if patch > 0 {
				patchVersions++
			}
		}
	}

	history["statistics"] = map[string]interface{}{
		"majorVersions": majorVersions,
		"minorVersions": minorVersions,
		"patchVersions": patchVersions,
	}

	return history, nil
}

// CheckForTemplateUpdates checks if there are newer versions available
func (a *App) CheckForTemplateUpdates(templateID string) (map[string]interface{}, error) {
	if templateID == "" {
		return nil, createValidationError("template ID is required", nil, nil)
	}

	// Get current template
	currentTemplate, err := a.repo.Templates.GetByID(templateID)
	if err != nil {
		return nil, err
	}

	// Get latest version
	latestTemplate, err := a.GetLatestTemplateVersion(currentTemplate.Name)
	if err != nil {
		return nil, err
	}

	updateInfo := make(map[string]interface{})
	updateInfo["currentVersion"] = currentTemplate.Version
	updateInfo["latestVersion"] = latestTemplate.Version
	updateInfo["templateName"] = currentTemplate.Name

	// Check if update is available
	isNewer, err := isNewerVersion(latestTemplate.Version, currentTemplate.Version)
	if err != nil {
		return nil, createInternalError("failed to compare versions", err)
	}

	updateInfo["updateAvailable"] = isNewer
	updateInfo["isCurrentVersion"] = !isNewer && latestTemplate.Version == currentTemplate.Version

	if isNewer {
		updateInfo["updateType"] = getUpdateType(currentTemplate.Version, latestTemplate.Version)
		updateInfo["changes"] = getVersionChanges(currentTemplate.Description, latestTemplate.Description)
	}

	return updateInfo, nil
}

// Helper functions for versioning

func getUpdateType(currentVersion, latestVersion string) string {
	currentParts := strings.Split(currentVersion, ".")
	latestParts := strings.Split(latestVersion, ".")

	if len(currentParts) != 3 || len(latestParts) != 3 {
		return "unknown"
	}

	var currentMajor, currentMinor, currentPatch int
	var latestMajor, latestMinor, latestPatch int

	fmt.Sscanf(currentParts[0], "%d", &currentMajor)
	fmt.Sscanf(currentParts[1], "%d", &currentMinor)
	fmt.Sscanf(currentParts[2], "%d", &currentPatch)

	fmt.Sscanf(latestParts[0], "%d", &latestMajor)
	fmt.Sscanf(latestParts[1], "%d", &latestMinor)
	fmt.Sscanf(latestParts[2], "%d", &latestPatch)

	if latestMajor > currentMajor {
		return "major"
	} else if latestMinor > currentMinor {
		return "minor"
	} else if latestPatch > currentPatch {
		return "patch"
	}

	return "unknown"
}

func getVersionChanges(currentDescription, latestDescription string) string {
	// Simple implementation - extract changes from description
	// Look for "Changes in v" pattern
	lines := strings.Split(latestDescription, "\n")
	for _, line := range lines {
		if strings.Contains(line, "Changes in v") {
			return strings.TrimSpace(line)
		}
	}
	return "No change information available"
}