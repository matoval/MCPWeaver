package validator

import (
	"context"
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/getkin/kin-openapi/openapi3"
)

// Service handles validation operations
type Service struct {
	loader *openapi3.Loader
}

// ValidationResult represents comprehensive validation results
type ValidationResult struct {
	Valid           bool                `json:"valid"`
	Errors          []ValidationError   `json:"errors"`
	Warnings        []ValidationWarning `json:"warnings"`
	Suggestions     []string            `json:"suggestions"`
	SpecInfo        *SpecInfo           `json:"specInfo,omitempty"`
	ValidationTime  time.Duration       `json:"validationTime"`
	CacheHit        bool                `json:"cacheHit"`
	ValidatedAt     time.Time           `json:"validatedAt"`
}

// ValidationError represents a validation error with detailed information
type ValidationError struct {
	Type        string         `json:"type"`
	Code        string         `json:"code"`
	Message     string         `json:"message"`
	Path        string         `json:"path"`
	Line        int            `json:"line,omitempty"`
	Column      int            `json:"column,omitempty"`
	Severity    SeverityLevel  `json:"severity"`
	Location    *ErrorLocation `json:"location,omitempty"`
	Context     string         `json:"context,omitempty"`
	Suggestion  string         `json:"suggestion,omitempty"`
}

// ValidationWarning represents a validation warning
type ValidationWarning struct {
	Type        string `json:"type"`
	Code        string `json:"code"`
	Message     string `json:"message"`
	Path        string `json:"path"`
	Suggestion  string `json:"suggestion"`
	Context     string `json:"context,omitempty"`
}

// SeverityLevel represents the severity of a validation issue
type SeverityLevel string

const (
	SeverityError   SeverityLevel = "error"
	SeverityWarning SeverityLevel = "warning"
	SeverityInfo    SeverityLevel = "info"
)

// ErrorLocation provides detailed location information for errors
type ErrorLocation struct {
	File   string `json:"file"`
	Line   int    `json:"line"`
	Column int    `json:"column"`
	Path   string `json:"path"`
}

// SpecInfo provides information about the validated specification
type SpecInfo struct {
	Version         string              `json:"version"`
	Title           string              `json:"title"`
	Description     string              `json:"description"`
	OperationCount  int                 `json:"operationCount"`
	SchemaCount     int                 `json:"schemaCount"`
	SecuritySchemes []SecuritySchemeInfo `json:"securitySchemes"`
	Servers         []ServerInfo        `json:"servers"`
	Tags            []TagInfo           `json:"tags"`
	Complexity      ComplexityLevel     `json:"complexity"`
	EstimatedSize   string              `json:"estimatedSize"`
}

// SecuritySchemeInfo provides information about security schemes
type SecuritySchemeInfo struct {
	Type         string `json:"type"`
	Name         string `json:"name"`
	Description  string `json:"description"`
	In           string `json:"in,omitempty"`
	Scheme       string `json:"scheme,omitempty"`
	BearerFormat string `json:"bearerFormat,omitempty"`
	OpenIDConnectURL string `json:"openIdConnectUrl,omitempty"`
}

// ServerInfo provides information about servers
type ServerInfo struct {
	URL         string                    `json:"url"`
	Description string                    `json:"description"`
	Variables   map[string]*ServerVariable `json:"variables,omitempty"`
}

// ServerVariable provides information about server variables
type ServerVariable struct {
	Default     string   `json:"default"`
	Description string   `json:"description"`
	Enum        []string `json:"enum,omitempty"`
}

// TagInfo provides information about tags
type TagInfo struct {
	Name         string            `json:"name"`
	Description  string            `json:"description"`
	ExternalDocs *ExternalDocsInfo `json:"externalDocs,omitempty"`
}

// ExternalDocsInfo provides information about external documentation
type ExternalDocsInfo struct {
	Description string `json:"description"`
	URL         string `json:"url"`
}

// ComplexityLevel represents the complexity of the specification
type ComplexityLevel string

const (
	ComplexityLow    ComplexityLevel = "low"
	ComplexityMedium ComplexityLevel = "medium"
	ComplexityHigh   ComplexityLevel = "high"
)

// New creates a new validator service
func New() *Service {
	return &Service{
		loader: openapi3.NewLoader(),
	}
}

// ValidateSpec validates an OpenAPI specification document
func (s *Service) ValidateSpec(ctx context.Context, spec *openapi3.T) (*ValidationResult, error) {
	startTime := time.Now()
	
	result := &ValidationResult{
		Valid:          true,
		Errors:         []ValidationError{},
		Warnings:       []ValidationWarning{},
		Suggestions:    []string{},
		ValidationTime: 0,
		CacheHit:       false,
		ValidatedAt:    startTime,
	}

	// Basic specification validation
	if err := spec.Validate(ctx); err != nil {
		result.Valid = false
		s.parseValidationError(err, result)
	}

	// Extract specification information
	s.extractSpecInfo(spec, result)

	// Perform detailed validation
	s.validateOperations(spec, result)
	s.validateSchemas(spec, result)
	s.validateSecuritySchemes(spec, result)
	s.validateServers(spec, result)
	s.validateExamples(spec, result)
	s.validateReferences(spec, result)

	// Generate suggestions
	s.generateSuggestions(spec, result)

	// Set final validation time
	result.ValidationTime = time.Since(startTime)

	return result, nil
}

// ValidateFile validates an OpenAPI specification file
func (s *Service) ValidateFile(ctx context.Context, filePath string) (*ValidationResult, error) {
	// Check if file exists
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return &ValidationResult{
			Valid: false,
			Errors: []ValidationError{
				{
					Type:     "file",
					Code:     "FILE_NOT_FOUND",
					Message:  fmt.Sprintf("File not found: %s", filePath),
					Path:     filePath,
					Severity: SeverityError,
				},
			},
			Warnings:    []ValidationWarning{},
			Suggestions: []string{},
		}, nil
	}

	// Load the specification
	spec, err := s.loader.LoadFromFile(filePath)
	if err != nil {
		return &ValidationResult{
			Valid: false,
			Errors: []ValidationError{
				{
					Type:     "parsing",
					Code:     "PARSING_ERROR",
					Message:  fmt.Sprintf("Failed to parse OpenAPI specification: %v", err),
					Path:     filePath,
					Severity: SeverityError,
					Location: &ErrorLocation{
						File: filepath.Base(filePath),
					},
				},
			},
			Warnings:    []ValidationWarning{},
			Suggestions: s.getParsingErrorSuggestions(err),
		}, nil
	}

	// Validate the specification
	return s.ValidateSpec(ctx, spec)
}

// ValidateURL validates an OpenAPI specification from a URL
func (s *Service) ValidateURL(ctx context.Context, specURL string) (*ValidationResult, error) {
	// Validate URL format
	if _, err := url.Parse(specURL); err != nil {
		return &ValidationResult{
			Valid: false,
			Errors: []ValidationError{
				{
					Type:     "url",
					Code:     "INVALID_URL",
					Message:  fmt.Sprintf("Invalid URL format: %v", err),
					Path:     specURL,
					Severity: SeverityError,
				},
			},
			Warnings:    []ValidationWarning{},
			Suggestions: []string{"Please provide a valid URL"},
		}, nil
	}

	// Load the specification from URL
	parsedURL, err := url.Parse(specURL)
	if err != nil {
		return &ValidationResult{
			Valid: false,
			Errors: []ValidationError{
				{
					Type:     "url",
					Code:     "INVALID_URL",
					Message:  fmt.Sprintf("Invalid URL format: %v", err),
					Path:     specURL,
					Severity: SeverityError,
				},
			},
			Warnings:    []ValidationWarning{},
			Suggestions: []string{"Please provide a valid URL"},
		}, nil
	}
	
	spec, err := s.loader.LoadFromURI(parsedURL)
	if err != nil {
		return &ValidationResult{
			Valid: false,
			Errors: []ValidationError{
				{
					Type:     "network",
					Code:     "NETWORK_ERROR",
					Message:  fmt.Sprintf("Failed to fetch OpenAPI specification: %v", err),
					Path:     specURL,
					Severity: SeverityError,
				},
			},
			Warnings:    []ValidationWarning{},
			Suggestions: s.getNetworkErrorSuggestions(err),
		}, nil
	}

	// Validate the specification
	return s.ValidateSpec(ctx, spec)
}

// parseValidationError parses validation errors from kin-openapi
func (s *Service) parseValidationError(err error, result *ValidationResult) {
	errorMsg := err.Error()
	
	// Try to extract line and column information
	lineRegex := regexp.MustCompile(`line (\d+)`)
	columnRegex := regexp.MustCompile(`column (\d+)`)
	
	var line, column int
	if matches := lineRegex.FindStringSubmatch(errorMsg); len(matches) > 1 {
		fmt.Sscanf(matches[1], "%d", &line)
	}
	if matches := columnRegex.FindStringSubmatch(errorMsg); len(matches) > 1 {
		fmt.Sscanf(matches[1], "%d", &column)
	}

	result.Errors = append(result.Errors, ValidationError{
		Type:     "validation",
		Code:     "SPEC_VALIDATION_ERROR",
		Message:  errorMsg,
		Path:     "/",
		Line:     line,
		Column:   column,
		Severity: SeverityError,
		Location: &ErrorLocation{
			Line:   line,
			Column: column,
		},
	})
}

// extractSpecInfo extracts information from the specification
func (s *Service) extractSpecInfo(spec *openapi3.T, result *ValidationResult) {
	info := &SpecInfo{
		Version:         spec.OpenAPI,
		Title:           spec.Info.Title,
		Description:     spec.Info.Description,
		OperationCount:  s.countOperations(spec),
		SchemaCount:     s.countSchemas(spec),
		SecuritySchemes: s.extractSecuritySchemes(spec),
		Servers:         s.extractServers(spec),
		Tags:            s.extractTags(spec),
		Complexity:      s.assessComplexity(spec),
		EstimatedSize:   s.estimateSize(spec),
	}
	
	result.SpecInfo = info
}

// validateOperations validates all operations in the specification
func (s *Service) validateOperations(spec *openapi3.T, result *ValidationResult) {
	if spec.Paths == nil {
		result.Warnings = append(result.Warnings, ValidationWarning{
			Type:       "operations",
			Code:       "NO_PATHS",
			Message:    "No paths defined in the specification",
			Path:       "/paths",
			Suggestion: "Add API endpoints to generate MCP tools",
		})
		return
	}

	operationIds := make(map[string]bool)
	
	for path, pathItem := range spec.Paths.Map() {
		if pathItem == nil {
			continue
		}
		
		for method, operation := range pathItem.Operations() {
			if operation == nil {
				continue
			}
			
			opPath := fmt.Sprintf("/paths%s/%s", path, strings.ToLower(method))
			
			// Check for duplicate operation IDs
			if operation.OperationID != "" {
				if operationIds[operation.OperationID] {
					result.Errors = append(result.Errors, ValidationError{
						Type:     "operations",
						Code:     "DUPLICATE_OPERATION_ID",
						Message:  fmt.Sprintf("Duplicate operation ID: %s", operation.OperationID),
						Path:     opPath,
						Severity: SeverityError,
					})
					result.Valid = false
				}
				operationIds[operation.OperationID] = true
			}

			// Check for missing descriptions
			if operation.Description == "" && operation.Summary == "" {
				result.Warnings = append(result.Warnings, ValidationWarning{
					Type:       "operations",
					Code:       "MISSING_DESCRIPTION",
					Message:    fmt.Sprintf("Operation %s %s has no description", method, path),
					Path:       opPath,
					Suggestion: "Add description or summary to improve generated tool documentation",
				})
			}

			// Validate parameters
			s.validateParameters(operation.Parameters, opPath, result)

			// Validate responses
			s.validateResponses(operation.Responses, opPath, result)
		}
	}
}

// validateSchemas validates all schemas in the specification
func (s *Service) validateSchemas(spec *openapi3.T, result *ValidationResult) {
	if spec.Components == nil || spec.Components.Schemas == nil {
		result.Warnings = append(result.Warnings, ValidationWarning{
			Type:       "schemas",
			Code:       "NO_SCHEMAS",
			Message:    "No schemas defined in the specification",
			Path:       "/components/schemas",
			Suggestion: "Define reusable schemas to improve specification maintainability",
		})
		return
	}

	for name, schemaRef := range spec.Components.Schemas {
		if schemaRef.Value == nil {
			result.Errors = append(result.Errors, ValidationError{
				Type:     "schemas",
				Code:     "EMPTY_SCHEMA",
				Message:  fmt.Sprintf("Schema '%s' has no definition", name),
				Path:     fmt.Sprintf("/components/schemas/%s", name),
				Severity: SeverityError,
			})
			result.Valid = false
			continue
		}

		schema := schemaRef.Value
		schemaPath := fmt.Sprintf("/components/schemas/%s", name)

		// Check for schemas without descriptions
		if schema.Description == "" {
			result.Warnings = append(result.Warnings, ValidationWarning{
				Type:       "schemas",
				Code:       "MISSING_DESCRIPTION",
				Message:    fmt.Sprintf("Schema '%s' has no description", name),
				Path:       schemaPath,
				Suggestion: "Add schema descriptions to improve documentation",
			})
		}

		// Validate schema properties
		s.validateSchemaProperties(schema, schemaPath, result)
	}
}

// validateSecuritySchemes validates security schemes
func (s *Service) validateSecuritySchemes(spec *openapi3.T, result *ValidationResult) {
	if spec.Components == nil || spec.Components.SecuritySchemes == nil {
		result.Suggestions = append(result.Suggestions, 
			"Consider adding security schemes if your API requires authentication")
		return
	}

	for name, securitySchemeRef := range spec.Components.SecuritySchemes {
		if securitySchemeRef.Value == nil {
			result.Errors = append(result.Errors, ValidationError{
				Type:     "security",
				Code:     "EMPTY_SECURITY_SCHEME",
				Message:  fmt.Sprintf("Security scheme '%s' has no definition", name),
				Path:     fmt.Sprintf("/components/securitySchemes/%s", name),
				Severity: SeverityError,
			})
			result.Valid = false
		}
	}
}

// validateServers validates server configurations
func (s *Service) validateServers(spec *openapi3.T, result *ValidationResult) {
	if len(spec.Servers) == 0 {
		result.Suggestions = append(result.Suggestions, 
			"Consider adding server information to specify the base URL")
		return
	}

	for i, server := range spec.Servers {
		if server.URL == "" {
			result.Warnings = append(result.Warnings, ValidationWarning{
				Type:       "servers",
				Code:       "EMPTY_SERVER_URL",
				Message:    fmt.Sprintf("Server #%d has no URL", i+1),
				Path:       fmt.Sprintf("/servers[%d]", i),
				Suggestion: "Provide a valid server URL",
			})
		}
	}
}

// validateExamples validates examples in the specification
func (s *Service) validateExamples(spec *openapi3.T, result *ValidationResult) {
	// This is a placeholder for example validation
	// In a real implementation, you would validate examples against schemas
}

// validateReferences validates all $ref references in the specification
func (s *Service) validateReferences(spec *openapi3.T, result *ValidationResult) {
	// This is a placeholder for reference validation
	// In a real implementation, you would check that all references are valid
}

// Helper methods for extracting information
func (s *Service) countOperations(spec *openapi3.T) int {
	if spec.Paths == nil {
		return 0
	}
	
	count := 0
	for _, pathItem := range spec.Paths.Map() {
		if pathItem != nil {
			count += len(pathItem.Operations())
		}
	}
	return count
}

func (s *Service) countSchemas(spec *openapi3.T) int {
	if spec.Components == nil || spec.Components.Schemas == nil {
		return 0
	}
	return len(spec.Components.Schemas)
}

func (s *Service) extractSecuritySchemes(spec *openapi3.T) []SecuritySchemeInfo {
	var schemes []SecuritySchemeInfo
	
	if spec.Components == nil || spec.Components.SecuritySchemes == nil {
		return schemes
	}
	
	for name, schemeRef := range spec.Components.SecuritySchemes {
		if schemeRef.Value == nil {
			continue
		}
		
		scheme := schemeRef.Value
		info := SecuritySchemeInfo{
			Type:        scheme.Type,
			Name:        name,
			Description: scheme.Description,
		}
		
		if scheme.In != "" {
			info.In = scheme.In
		}
		if scheme.Scheme != "" {
			info.Scheme = scheme.Scheme
		}
		if scheme.BearerFormat != "" {
			info.BearerFormat = scheme.BearerFormat
		}
		if scheme.OpenIdConnectUrl != "" {
			info.OpenIDConnectURL = scheme.OpenIdConnectUrl
		}
		
		schemes = append(schemes, info)
	}
	
	return schemes
}

func (s *Service) extractServers(spec *openapi3.T) []ServerInfo {
	var servers []ServerInfo
	
	for _, server := range spec.Servers {
		info := ServerInfo{
			URL:         server.URL,
			Description: server.Description,
			Variables:   make(map[string]*ServerVariable),
		}
		
		for name, variable := range server.Variables {
			info.Variables[name] = &ServerVariable{
				Default:     variable.Default,
				Description: variable.Description,
				Enum:        variable.Enum,
			}
		}
		
		servers = append(servers, info)
	}
	
	return servers
}

func (s *Service) extractTags(spec *openapi3.T) []TagInfo {
	var tags []TagInfo
	
	for _, tag := range spec.Tags {
		info := TagInfo{
			Name:        tag.Name,
			Description: tag.Description,
		}
		
		if tag.ExternalDocs != nil {
			info.ExternalDocs = &ExternalDocsInfo{
				Description: tag.ExternalDocs.Description,
				URL:         tag.ExternalDocs.URL,
			}
		}
		
		tags = append(tags, info)
	}
	
	return tags
}

func (s *Service) assessComplexity(spec *openapi3.T) ComplexityLevel {
	operationCount := s.countOperations(spec)
	schemaCount := s.countSchemas(spec)
	
	totalComplexity := operationCount + schemaCount
	
	if totalComplexity > 100 {
		return ComplexityHigh
	} else if totalComplexity > 20 {
		return ComplexityMedium
	}
	return ComplexityLow
}

func (s *Service) estimateSize(spec *openapi3.T) string {
	operationCount := s.countOperations(spec)
	
	if operationCount > 100 {
		return "Large"
	} else if operationCount > 20 {
		return "Medium"
	}
	return "Small"
}

// Additional validation helper methods
func (s *Service) validateParameters(params openapi3.Parameters, basePath string, result *ValidationResult) {
	for i, paramRef := range params {
		if paramRef.Value == nil {
			continue
		}
		
		param := paramRef.Value
		paramPath := fmt.Sprintf("%s/parameters[%d]", basePath, i)
		
		if param.Description == "" {
			result.Warnings = append(result.Warnings, ValidationWarning{
				Type:       "parameters",
				Code:       "MISSING_DESCRIPTION",
				Message:    fmt.Sprintf("Parameter '%s' has no description", param.Name),
				Path:       paramPath,
				Suggestion: "Add parameter descriptions to improve usability",
			})
		}
	}
}

func (s *Service) validateResponses(responses *openapi3.Responses, basePath string, result *ValidationResult) {
	if responses == nil {
		result.Warnings = append(result.Warnings, ValidationWarning{
			Type:       "responses",
			Code:       "NO_RESPONSES",
			Message:    "No responses defined",
			Path:       fmt.Sprintf("%s/responses", basePath),
			Suggestion: "Define response schemas to improve API documentation",
		})
		return
	}

	// Check for missing success responses
	hasSuccessResponse := false
	for status := range responses.Map() {
		if strings.HasPrefix(status, "2") {
			hasSuccessResponse = true
			break
		}
	}
	
	if !hasSuccessResponse {
		result.Warnings = append(result.Warnings, ValidationWarning{
			Type:       "responses",
			Code:       "NO_SUCCESS_RESPONSE",
			Message:    "No success response (2xx) defined",
			Path:       fmt.Sprintf("%s/responses", basePath),
			Suggestion: "Add success response definitions",
		})
	}
}

func (s *Service) validateSchemaProperties(schema *openapi3.Schema, schemaPath string, result *ValidationResult) {
	// Validate required properties exist
	for _, required := range schema.Required {
		if schema.Properties == nil || schema.Properties[required] == nil {
			result.Errors = append(result.Errors, ValidationError{
				Type:     "schemas",
				Code:     "MISSING_REQUIRED_PROPERTY",
				Message:  fmt.Sprintf("Required property '%s' is not defined", required),
				Path:     fmt.Sprintf("%s/properties", schemaPath),
				Severity: SeverityError,
			})
			result.Valid = false
		}
	}
}

// generateSuggestions generates helpful suggestions based on the specification
func (s *Service) generateSuggestions(spec *openapi3.T, result *ValidationResult) {
	// Check for large specifications
	if s.countOperations(spec) > 50 {
		result.Suggestions = append(result.Suggestions, 
			"Large specifications may result in many MCP tools. Consider grouping related operations")
	}
	
	// Check for missing base URL
	if len(spec.Servers) == 0 {
		result.Suggestions = append(result.Suggestions, 
			"Consider adding server information to specify the base URL")
	}
	
	// Check for missing OpenAPI version
	if spec.OpenAPI == "" {
		result.Suggestions = append(result.Suggestions, 
			"Specify the OpenAPI version for better compatibility")
	}
}

// getParsingErrorSuggestions returns suggestions for parsing errors
func (s *Service) getParsingErrorSuggestions(err error) []string {
	errorStr := err.Error()
	var suggestions []string
	
	if strings.Contains(errorStr, "array schema") {
		suggestions = append(suggestions, 
			"Consider simplifying array schemas to use standard 'items' object definitions")
	}
	if strings.Contains(errorStr, "regex") {
		suggestions = append(suggestions, 
			"Remove or simplify regex patterns with unsupported features like lookaheads")
	}
	if strings.Contains(errorStr, "OpenAPI 2.0") {
		suggestions = append(suggestions, 
			"Convert your specification to OpenAPI 3.0+ format")
	}
	if strings.Contains(errorStr, "YAML") {
		suggestions = append(suggestions, 
			"Check YAML syntax and indentation")
	}
	if strings.Contains(errorStr, "JSON") {
		suggestions = append(suggestions, 
			"Check JSON syntax and structure")
	}
	
	return suggestions
}

// getNetworkErrorSuggestions returns suggestions for network errors
func (s *Service) getNetworkErrorSuggestions(err error) []string {
	errorStr := err.Error()
	var suggestions []string
	
	if strings.Contains(errorStr, "HTTP error") {
		suggestions = append(suggestions, 
			"Check if the URL is accessible and returns a valid OpenAPI specification")
	}
	if strings.Contains(errorStr, "timeout") {
		suggestions = append(suggestions, 
			"The server may be slow to respond. Try again or check the URL")
	}
	if strings.Contains(errorStr, "connection") {
		suggestions = append(suggestions, 
			"Check your internet connection and firewall settings")
	}
	
	return suggestions
}