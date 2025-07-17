package app

import (
	"fmt"
	"strings"
	"time"

	"MCPWeaver/internal/parser"
	"github.com/getkin/kin-openapi/openapi3"
)

// ValidateSpec validates an OpenAPI specification from a file path
func (a *App) ValidateSpec(specPath string) (*ValidationResult, error) {
	if specPath == "" {
		return nil, a.createAPIError("validation", ErrCodeValidation, "Specification path is required", nil)
	}

	startTime := time.Now()
	
	// Parse the specification
	parsedAPI, err := a.parserService.ParseFromFile(specPath)
	validationTime := time.Since(startTime)
	
	result := &ValidationResult{
		Valid:          true,
		Errors:         []ValidationError{},
		Warnings:       []ValidationWarning{},
		Suggestions:    []string{},
		ValidationTime: validationTime,
	}

	if err != nil {
		result.Valid = false
		result.Errors = append(result.Errors, ValidationError{
			Type:     "parsing",
			Message:  fmt.Sprintf("Failed to parse OpenAPI specification: %v", err),
			Path:     specPath,
			Severity: "error",
			Code:     ErrCodeParsingError,
		})
		
		// Add suggestions based on error type
		if strings.Contains(err.Error(), "array schema") {
			result.Suggestions = append(result.Suggestions, 
				"Consider simplifying array schemas to use standard 'items' object definitions")
		}
		if strings.Contains(err.Error(), "regex") {
			result.Suggestions = append(result.Suggestions, 
				"Remove or simplify regex patterns with unsupported features like lookaheads")
		}
		if strings.Contains(err.Error(), "OpenAPI 2.0") {
			result.Suggestions = append(result.Suggestions, 
				"Convert your specification to OpenAPI 3.0+ format")
		}
		
		return result, nil
	}

	// Specification parsed successfully, add spec info
	result.SpecInfo = &SpecInfo{
		Version:        parsedAPI.Version,
		Title:          parsedAPI.Title,
		Description:    parsedAPI.Description,
		OperationCount: len(parsedAPI.Operations),
		SchemaCount:    len(parsedAPI.Schemas),
		Servers:        []ServerInfo{},
	}

	// Add server information
	for _, server := range parsedAPI.Servers {
		result.SpecInfo.Servers = append(result.SpecInfo.Servers, ServerInfo{
			URL:         server,
			Description: "",
		})
	}

	// Validate operations
	a.validateOperations(parsedAPI.Operations, result)
	
	// Validate schemas
	a.validateSchemas(parsedAPI.Schemas, result)
	
	// Add general suggestions
	a.addGeneralSuggestions(parsedAPI, result)

	return result, nil
}

// ValidateURL validates an OpenAPI specification from a URL
func (a *App) ValidateURL(url string) (*ValidationResult, error) {
	if url == "" {
		return nil, a.createAPIError("validation", ErrCodeValidation, "URL is required", nil)
	}

	startTime := time.Now()
	
	// Parse the specification from URL
	parsedAPI, err := a.parserService.ParseFromURL(a.ctx, url)
	validationTime := time.Since(startTime)
	
	result := &ValidationResult{
		Valid:          true,
		Errors:         []ValidationError{},
		Warnings:       []ValidationWarning{},
		Suggestions:    []string{},
		ValidationTime: validationTime,
	}

	if err != nil {
		result.Valid = false
		result.Errors = append(result.Errors, ValidationError{
			Type:     "network",
			Message:  fmt.Sprintf("Failed to fetch or parse OpenAPI specification: %v", err),
			Path:     url,
			Severity: "error",
			Code:     ErrCodeNetworkError,
		})
		
		// Add suggestions based on error type
		if strings.Contains(err.Error(), "HTTP error") {
			result.Suggestions = append(result.Suggestions, 
				"Check if the URL is accessible and returns a valid OpenAPI specification")
		}
		if strings.Contains(err.Error(), "timeout") {
			result.Suggestions = append(result.Suggestions, 
				"The server may be slow to respond. Try again or check the URL")
		}
		
		return result, nil
	}

	// Specification parsed successfully, add spec info
	result.SpecInfo = &SpecInfo{
		Version:        parsedAPI.Version,
		Title:          parsedAPI.Title,
		Description:    parsedAPI.Description,
		OperationCount: len(parsedAPI.Operations),
		SchemaCount:    len(parsedAPI.Schemas),
		Servers:        []ServerInfo{},
	}

	// Add server information
	for _, server := range parsedAPI.Servers {
		result.SpecInfo.Servers = append(result.SpecInfo.Servers, ServerInfo{
			URL:         server,
			Description: "",
		})
	}

	// Validate operations
	a.validateOperations(parsedAPI.Operations, result)
	
	// Validate schemas
	a.validateSchemas(parsedAPI.Schemas, result)
	
	// Add general suggestions
	a.addGeneralSuggestions(parsedAPI, result)

	return result, nil
}

// validateOperations validates the operations in the specification
func (a *App) validateOperations(operations []parser.Operation, result *ValidationResult) {
	if len(operations) == 0 {
		result.Warnings = append(result.Warnings, ValidationWarning{
			Type:       "operations",
			Message:    "No operations found in the specification",
			Path:       "/paths",
			Suggestion: "Add API endpoints to generate MCP tools",
		})
		return
	}

	operationIds := make(map[string]bool)
	
	for _, op := range operations {
		// Check for duplicate operation IDs
		if op.ID != "" {
			if operationIds[op.ID] {
				result.Errors = append(result.Errors, ValidationError{
					Type:     "operations",
					Message:  fmt.Sprintf("Duplicate operation ID: %s", op.ID),
					Path:     fmt.Sprintf("/paths%s", op.Path),
					Severity: "error",
					Code:     ErrCodeValidation,
				})
				result.Valid = false
			}
			operationIds[op.ID] = true
		}

		// Check for missing descriptions
		if op.Description == "" && op.Summary == "" {
			result.Warnings = append(result.Warnings, ValidationWarning{
				Type:       "operations",
				Message:    fmt.Sprintf("Operation %s %s has no description", op.Method, op.Path),
				Path:       fmt.Sprintf("/paths%s/%s", op.Path, strings.ToLower(op.Method)),
				Suggestion: "Add description or summary to improve generated tool documentation",
			})
		}

		// Check for parameters without descriptions
		for _, param := range op.Parameters {
			if param.Description == "" {
				result.Warnings = append(result.Warnings, ValidationWarning{
					Type:       "parameters",
					Message:    fmt.Sprintf("Parameter '%s' in %s %s has no description", param.Name, op.Method, op.Path),
					Path:       fmt.Sprintf("/paths%s/%s/parameters", op.Path, strings.ToLower(op.Method)),
					Suggestion: "Add parameter descriptions to improve usability",
				})
			}
		}
	}
}

// validateSchemas validates the schemas in the specification
func (a *App) validateSchemas(schemas map[string]*openapi3.SchemaRef, result *ValidationResult) {
	if len(schemas) == 0 {
		result.Warnings = append(result.Warnings, ValidationWarning{
			Type:       "schemas",
			Message:    "No schemas defined in the specification",
			Path:       "/components/schemas",
			Suggestion: "Define reusable schemas to improve specification maintainability",
		})
		return
	}

	for name, schemaRef := range schemas {
		if schemaRef.Value == nil {
			result.Errors = append(result.Errors, ValidationError{
				Type:     "schemas",
				Message:  fmt.Sprintf("Schema '%s' has no definition", name),
				Path:     fmt.Sprintf("/components/schemas/%s", name),
				Severity: "error",
				Code:     ErrCodeValidation,
			})
			result.Valid = false
			continue
		}

		schema := schemaRef.Value

		// Check for schemas without descriptions
		if schema.Description == "" {
			result.Warnings = append(result.Warnings, ValidationWarning{
				Type:       "schemas",
				Message:    fmt.Sprintf("Schema '%s' has no description", name),
				Path:       fmt.Sprintf("/components/schemas/%s", name),
				Suggestion: "Add schema descriptions to improve documentation",
			})
		}
	}
}

// addGeneralSuggestions adds general suggestions based on the specification
func (a *App) addGeneralSuggestions(parsedAPI *parser.ParsedAPI, result *ValidationResult) {
	// Check for missing base URL
	if parsedAPI.BaseURL == "" && len(parsedAPI.Servers) == 0 {
		result.Suggestions = append(result.Suggestions, 
			"Consider adding server information to specify the base URL")
	}

	// Check for complex specifications
	if len(parsedAPI.Operations) > 50 {
		result.Suggestions = append(result.Suggestions, 
			"Large specifications may result in many MCP tools. Consider grouping related operations")
	}

	// Check for security schemes
	hasSecuritySchemes := false
	if parsedAPI.Document.Components != nil && len(parsedAPI.Document.Components.SecuritySchemes) > 0 {
		hasSecuritySchemes = true
	}

	if !hasSecuritySchemes {
		result.Suggestions = append(result.Suggestions, 
			"Consider adding security schemes if your API requires authentication")
	}

	// Final validation check
	if len(result.Errors) == 0 {
		result.Valid = true
	}
}