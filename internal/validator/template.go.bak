package validator

import (
	"fmt"
	"os"
	"regexp"
	"strings"
	"text/template"
	"time"

	"MCPWeaver/internal/app"
)

// TemplateValidator handles template validation
type TemplateValidator struct {
	maxTemplateSize int64
	allowedFunctions map[string]bool
}

// NewTemplateValidator creates a new template validator
func NewTemplateValidator() *TemplateValidator {
	return &TemplateValidator{
		maxTemplateSize: 1024 * 1024, // 1MB max
		allowedFunctions: map[string]bool{
			"title":     true,
			"upper":     true,
			"lower":     true,
			"trim":      true,
			"replace":   true,
			"contains":  true,
			"hasPrefix": true,
			"hasSuffix": true,
			"split":     true,
			"join":      true,
			"printf":    true,
			"len":       true,
			"index":     true,
			"range":     true,
			"if":        true,
			"else":      true,
			"end":       true,
			"with":      true,
			"and":       true,
			"or":        true,
			"not":       true,
			"eq":        true,
			"ne":        true,
			"lt":        true,
			"le":        true,
			"gt":        true,
			"ge":        true,
		},
	}
}

// ValidateTemplate performs comprehensive template validation
func (v *TemplateValidator) ValidateTemplate(templatePath string, variables []app.TemplateVariable) (*app.TemplateValidationResult, error) {
	startTime := time.Now()
	
	result := &app.TemplateValidationResult{
		Valid:        true,
		Errors:       []app.TemplateError{},
		Warnings:     []app.TemplateWarning{},
		Suggestions:  []string{},
		Dependencies: []app.TemplateDependency{},
		Performance: &app.TemplatePerformance{
			CacheHit: false,
		},
	}

	// Check if file exists
	if !fileExists(templatePath) {
		result.Valid = false
		result.Errors = append(result.Errors, app.TemplateError{
			Type:     "file_not_found",
			Message:  "Template file does not exist",
			Severity: "error",
		})
		return result, nil
	}

	// Check file size
	fileInfo, err := os.Stat(templatePath)
	if err != nil {
		result.Valid = false
		result.Errors = append(result.Errors, app.TemplateError{
			Type:     "file_access_error",
			Message:  fmt.Sprintf("Cannot access template file: %v", err),
			Severity: "error",
		})
		return result, nil
	}

	if fileInfo.Size() > v.maxTemplateSize {
		result.Valid = false
		result.Errors = append(result.Errors, app.TemplateError{
			Type:     "file_too_large",
			Message:  fmt.Sprintf("Template file is too large (%d bytes, max %d bytes)", fileInfo.Size(), v.maxTemplateSize),
			Severity: "error",
		})
		return result, nil
	}

	// Read template content
	content, err := os.ReadFile(templatePath)
	if err != nil {
		result.Valid = false
		result.Errors = append(result.Errors, app.TemplateError{
			Type:     "file_read_error",
			Message:  fmt.Sprintf("Failed to read template file: %v", err),
			Severity: "error",
		})
		return result, nil
	}

	contentStr := string(content)

	// Validate syntax
	if err := v.validateSyntax(contentStr, result); err != nil {
		result.Valid = false
	}

	// Validate security
	if err := v.validateSecurity(contentStr, result); err != nil {
		result.Valid = false
	}

	// Validate variables
	v.validateVariables(contentStr, variables, result)

	// Validate structure
	v.validateStructure(contentStr, result)

	// Check performance
	v.checkPerformance(contentStr, result)

	// Check dependencies
	v.checkDependencies(contentStr, result)

	// Set performance metrics
	result.Performance.RenderTime = time.Since(startTime)
	result.Performance.MemoryUsage = int64(len(content))
	result.Performance.Complexity = v.calculateComplexity(contentStr)

	return result, nil
}

// validateSyntax validates Go template syntax
func (v *TemplateValidator) validateSyntax(content string, result *app.TemplateValidationResult) error {
	// Try to parse the template
	tmpl, err := template.New("test").Parse(content)
	if err != nil {
		result.Errors = append(result.Errors, app.TemplateError{
			Type:     "syntax_error",
			Message:  fmt.Sprintf("Template syntax error: %v", err),
			Severity: "error",
		})
		return err
	}

	// Check for balanced delimiters
	openCount := strings.Count(content, "{{")
	closeCount := strings.Count(content, "}}")
	
	if openCount != closeCount {
		result.Errors = append(result.Errors, app.TemplateError{
			Type:     "syntax_error",
			Message:  fmt.Sprintf("Unbalanced template delimiters: %d opening, %d closing", openCount, closeCount),
			Severity: "error",
		})
		return fmt.Errorf("unbalanced delimiters")
	}

	// Validate template can be executed with empty data
	var testBuf strings.Builder
	testData := make(map[string]interface{})
	if err := tmpl.Execute(&testBuf, testData); err != nil {
		result.Warnings = append(result.Warnings, app.TemplateWarning{
			Type:       "execution_warning",
			Message:    fmt.Sprintf("Template may have execution issues: %v", err),
			Suggestion: "Test template with sample data",
		})
	}

	return nil
}

// validateSecurity checks for security issues
func (v *TemplateValidator) validateSecurity(content string, result *app.TemplateValidationResult) error {
	// Check for forbidden functions
	forbiddenPatterns := []struct {
		pattern string
		message string
	}{
		{`\{\{\.?[^}]*exec[^}]*\}\}`, "Template contains exec function"},
		{`\{\{\.?[^}]*system[^}]*\}\}`, "Template contains system function"},
		{`\{\{\.?[^}]*command[^}]*\}\}`, "Template contains command function"},
		{`\{\{\.?[^}]*eval[^}]*\}\}`, "Template contains eval function"},
		{`\{\{\.?[^}]*shell[^}]*\}\}`, "Template contains shell function"},
		{`\{\{\.?[^}]*os\.[^}]*\}\}`, "Template accesses os package"},
		{`\{\{\.?[^}]*filepath\.[^}]*\}\}`, "Template accesses filepath package"},
		{`\{\{\.?[^}]*file[^}]*\}\}`, "Template contains file operations"},
	}

	for _, forbidden := range forbiddenPatterns {
		re := regexp.MustCompile(`(?i)` + forbidden.pattern)
		if re.MatchString(content) {
			result.Errors = append(result.Errors, app.TemplateError{
				Type:     "security_error",
				Message:  forbidden.message,
				Severity: "error",
			})
			return fmt.Errorf("security violation: %s", forbidden.message)
		}
	}

	// Check for potential code injection
	injectionPatterns := []string{
		"javascript:",
		"<script",
		"</script>",
		"eval(",
		"setTimeout(",
		"setInterval(",
	}

	for _, pattern := range injectionPatterns {
		if strings.Contains(strings.ToLower(content), pattern) {
			result.Warnings = append(result.Warnings, app.TemplateWarning{
				Type:       "security_warning",
				Message:    fmt.Sprintf("Potential security risk detected: %s", pattern),
				Suggestion: "Review template for potential injection vulnerabilities",
			})
		}
	}

	return nil
}

// validateVariables checks variable usage
func (v *TemplateValidator) validateVariables(content string, variables []app.TemplateVariable, result *app.TemplateValidationResult) {
	// Extract variables used in template
	usedVars := v.extractVariableUsage(content)
	
	// Check for undefined variables
	definedVars := make(map[string]app.TemplateVariable)
	for _, variable := range variables {
		definedVars[variable.Name] = variable
	}

	for _, usedVar := range usedVars {
		if _, defined := definedVars[usedVar]; !defined {
			result.Warnings = append(result.Warnings, app.TemplateWarning{
				Type:       "undefined_variable",
				Message:    fmt.Sprintf("Variable '%s' is used but not defined", usedVar),
				Suggestion: "Add variable definition or remove usage",
			})
		}
	}

	// Check for unused variables
	for _, variable := range variables {
		found := false
		for _, usedVar := range usedVars {
			if usedVar == variable.Name {
				found = true
				break
			}
		}
		if !found {
			result.Warnings = append(result.Warnings, app.TemplateWarning{
				Type:       "unused_variable",
				Message:    fmt.Sprintf("Variable '%s' is defined but not used", variable.Name),
				Suggestion: "Remove unused variable or add to template",
			})
		}
	}

	// Validate variable types and constraints
	for _, variable := range variables {
		if variable.Required && variable.DefaultValue == "" {
			result.Suggestions = append(result.Suggestions, 
				fmt.Sprintf("Consider providing a default value for required variable '%s'", variable.Name))
		}

		// Validate variable type
		if err := v.validateVariableType(variable); err != nil {
			result.Warnings = append(result.Warnings, app.TemplateWarning{
				Type:       "variable_type_warning",
				Message:    fmt.Sprintf("Variable '%s': %v", variable.Name, err),
				Suggestion: "Check variable type definition",
			})
		}
	}
}

// validateStructure checks template structure
func (v *TemplateValidator) validateStructure(content string, result *app.TemplateValidationResult) {
	lines := strings.Split(content, "\n")
	
	// Check for very long lines
	for i, line := range lines {
		if len(line) > 200 {
			result.Warnings = append(result.Warnings, app.TemplateWarning{
				Type:       "long_line",
				Message:    fmt.Sprintf("Line %d is very long (%d characters)", i+1, len(line)),
				Line:       i + 1,
				Suggestion: "Consider breaking long lines for better readability",
			})
		}
	}

	// Check nesting depth
	maxNesting := v.checkNestingDepth(content)
	if maxNesting > 5 {
		result.Warnings = append(result.Warnings, app.TemplateWarning{
			Type:       "deep_nesting",
			Message:    fmt.Sprintf("Template has deep nesting (level %d)", maxNesting),
			Suggestion: "Consider refactoring to reduce nesting complexity",
		})
	}

	// Check for common patterns
	if strings.Contains(content, "{{range") && !strings.Contains(content, "{{else}}") {
		result.Suggestions = append(result.Suggestions, 
			"Consider adding {{else}} clause to handle empty ranges")
	}
}

// checkPerformance analyzes performance characteristics
func (v *TemplateValidator) checkPerformance(content string, result *app.TemplateValidationResult) {
	// Check for loops
	rangeCount := strings.Count(content, "{{range")
	if rangeCount > 3 {
		result.Warnings = append(result.Warnings, app.TemplateWarning{
			Type:       "performance_warning",
			Message:    fmt.Sprintf("Template contains %d range loops", rangeCount),
			Suggestion: "Multiple loops may impact performance with large datasets",
		})
	}

	// Check for nested loops
	if v.hasNestedLoops(content) {
		result.Warnings = append(result.Warnings, app.TemplateWarning{
			Type:       "performance_warning",
			Message:    "Template contains nested loops",
			Suggestion: "Nested loops can significantly impact performance",
		})
	}

	// Check template size
	if len(content) > 10000 {
		result.Warnings = append(result.Warnings, app.TemplateWarning{
			Type:       "performance_warning",
			Message:    fmt.Sprintf("Large template size (%d bytes)", len(content)),
			Suggestion: "Consider breaking into smaller templates",
		})
	}
}

// checkDependencies identifies template dependencies
func (v *TemplateValidator) checkDependencies(content string, result *app.TemplateValidationResult) {
	// Check for template function usage
	commonFunctions := []string{"title", "upper", "lower", "trim", "replace", "printf"}
	
	for _, fn := range commonFunctions {
		if strings.Contains(content, fn) {
			result.Dependencies = append(result.Dependencies, app.TemplateDependency{
				Name:     fn,
				Type:     "function",
				Required: true,
			})
		}
	}

	// Check for Go standard library usage (if any template functions use them)
	if strings.Contains(content, "printf") || strings.Contains(content, "sprintf") {
		result.Dependencies = append(result.Dependencies, app.TemplateDependency{
			Name:     "fmt",
			Type:     "library",
			Required: true,
		})
	}
}

// Helper functions

func (v *TemplateValidator) extractVariableUsage(content string) []string {
	var variables []string
	
	// Pattern to match {{.VariableName}} or {{.Variable.Field}}
	re := regexp.MustCompile(`\{\{\.([A-Za-z][A-Za-z0-9_]*(?:\.[A-Za-z][A-Za-z0-9_]*)*)\}\}`)
	matches := re.FindAllStringSubmatch(content, -1)
	
	for _, match := range matches {
		if len(match) > 1 {
			variable := strings.Split(match[1], ".")[0] // Get root variable name
			if !contains(variables, variable) {
				variables = append(variables, variable)
			}
		}
	}
	
	return variables
}

func (v *TemplateValidator) validateVariableType(variable app.TemplateVariable) error {
	validTypes := []string{"string", "int", "bool", "float", "array", "object", "enum"}
	
	isValid := false
	for _, validType := range validTypes {
		if variable.Type == validType {
			isValid = true
			break
		}
	}
	
	if !isValid {
		return fmt.Errorf("invalid variable type '%s'", variable.Type)
	}

	// Type-specific validation
	switch variable.Type {
	case "enum":
		if len(variable.Options) == 0 {
			return fmt.Errorf("enum type requires options")
		}
	case "int", "float":
		if variable.DefaultValue != "" {
			// Could validate if default value is actually a number
		}
	}
	
	return nil
}

func (v *TemplateValidator) checkNestingDepth(content string) int {
	maxDepth := 0
	currentDepth := 0
	
	lines := strings.Split(content, "\n")
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		
		// Count opening constructs
		if strings.Contains(trimmed, "{{range") || 
		   strings.Contains(trimmed, "{{if") ||
		   strings.Contains(trimmed, "{{with") {
			currentDepth++
			if currentDepth > maxDepth {
				maxDepth = currentDepth
			}
		}
		
		// Count closing constructs
		if strings.Contains(trimmed, "{{end}}") {
			currentDepth--
		}
	}
	
	return maxDepth
}

func (v *TemplateValidator) hasNestedLoops(content string) bool {
	depth := 0
	rangeDepth := 0
	
	lines := strings.Split(content, "\n")
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		
		if strings.Contains(trimmed, "{{range") {
			depth++
			rangeDepth++
			if rangeDepth > 1 {
				return true
			}
		}
		
		if strings.Contains(trimmed, "{{end}}") && depth > 0 {
			depth--
			if rangeDepth > 0 {
				rangeDepth--
			}
		}
	}
	
	return false
}

func (v *TemplateValidator) calculateComplexity(content string) string {
	score := 0
	
	// Count various complexity factors
	score += strings.Count(content, "{{range") * 3    // Loops add complexity
	score += strings.Count(content, "{{if") * 2       // Conditionals add complexity
	score += strings.Count(content, "{{with") * 2     // With blocks add complexity
	score += len(strings.Split(content, "\n")) / 10   // Lines of code
	
	if score < 10 {
		return "low"
	} else if score < 25 {
		return "medium"
	} else {
		return "high"
	}
}

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}