package app

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"MCPWeaver/internal/database"
	"MCPWeaver/internal/parser"
	"MCPWeaver/internal/validator"
	"github.com/getkin/kin-openapi/openapi3"
)

// ValidateSpec validates an OpenAPI specification from a file path with caching
func (a *App) ValidateSpec(specPath string) (*ValidationResult, error) {
	if specPath == "" {
		return nil, a.createAPIError("validation", ErrCodeValidation, "Specification path is required", nil)
	}

	// Check validation cache first
	if a.validationCacheRepo != nil {
		specHash, err := a.validationCacheRepo.GenerateSpecHash(specPath)
		if err == nil {
			cachedResult, err := a.validationCacheRepo.GetByHash(specHash)
			if err == nil && cachedResult != nil {
				// Cache hit - deserialize result
				var result ValidationResult
				if err := json.Unmarshal([]byte(cachedResult.ValidationResult), &result); err == nil {
					result.CacheHit = true
					return &result, nil
				}
			}
		}
	}

	// Cache miss - perform validation
	startTime := time.Now()

	// Use the comprehensive validator service
	result, err := a.validatorService.ValidateFile(a.ctx, specPath)
	if err != nil {
		return nil, a.createAPIError("validation", ErrCodeInternalError, fmt.Sprintf("Validation failed: %v", err), nil)
	}

	// Convert validator result to app result format
	appResult := a.convertValidatorResult(result)
	appResult.ValidationTime = time.Since(startTime)

	// Cache the result if caching is enabled
	if a.validationCacheRepo != nil {
		a.cacheValidationResult(specPath, "", appResult)
	}

	return appResult, nil
}

// ValidateURL validates an OpenAPI specification from a URL with caching
func (a *App) ValidateURL(url string) (*ValidationResult, error) {
	if url == "" {
		return nil, a.createAPIError("validation", ErrCodeValidation, "URL is required", nil)
	}

	// Check validation cache first (using URL-based hash)
	if a.validationCacheRepo != nil {
		specHash := a.validationCacheRepo.GenerateURLHash(url, time.Now())
		cachedResult, err := a.validationCacheRepo.GetByHash(specHash)
		if err == nil && cachedResult != nil {
			// Cache hit - deserialize result
			var result ValidationResult
			if err := json.Unmarshal([]byte(cachedResult.ValidationResult), &result); err == nil {
				result.CacheHit = true
				return &result, nil
			}
		}
	}

	// Cache miss - perform validation
	startTime := time.Now()

	// Use the comprehensive validator service
	result, err := a.validatorService.ValidateURL(a.ctx, url)
	if err != nil {
		return nil, a.createAPIError("validation", ErrCodeInternalError, fmt.Sprintf("Validation failed: %v", err), nil)
	}

	// Convert validator result to app result format
	appResult := a.convertValidatorResult(result)
	appResult.ValidationTime = time.Since(startTime)

	// Cache the result if caching is enabled
	if a.validationCacheRepo != nil {
		a.cacheValidationResult("", url, appResult)
	}

	return appResult, nil
}

// convertValidatorResult converts validator.ValidationResult to app.ValidationResult
func (a *App) convertValidatorResult(vResult *validator.ValidationResult) *ValidationResult {
	result := &ValidationResult{
		Valid:          vResult.Valid,
		Errors:         make([]ValidationError, len(vResult.Errors)),
		Warnings:       make([]ValidationWarning, len(vResult.Warnings)),
		Suggestions:    vResult.Suggestions,
		ValidationTime: vResult.ValidationTime,
		CacheHit:       vResult.CacheHit,
		ValidatedAt:    vResult.ValidatedAt,
	}

	// Convert errors
	for i, verr := range vResult.Errors {
		result.Errors[i] = ValidationError{
			Type:     verr.Type,
			Message:  verr.Message,
			Path:     verr.Path,
			Line:     verr.Line,
			Column:   verr.Column,
			Severity: string(verr.Severity),
			Code:     verr.Code,
		}

		if verr.Location != nil {
			result.Errors[i].Location = &ErrorLocation{
				File:   verr.Location.File,
				Line:   verr.Location.Line,
				Column: verr.Location.Column,
			}
		}
	}

	// Convert warnings
	for i, vwarn := range vResult.Warnings {
		result.Warnings[i] = ValidationWarning{
			Type:       vwarn.Type,
			Message:    vwarn.Message,
			Path:       vwarn.Path,
			Suggestion: vwarn.Suggestion,
		}
	}

	// Convert spec info
	if vResult.SpecInfo != nil {
		result.SpecInfo = &SpecInfo{
			Version:         vResult.SpecInfo.Version,
			Title:           vResult.SpecInfo.Title,
			Description:     vResult.SpecInfo.Description,
			OperationCount:  vResult.SpecInfo.OperationCount,
			SchemaCount:     vResult.SpecInfo.SchemaCount,
			SecuritySchemes: make([]SecurityScheme, len(vResult.SpecInfo.SecuritySchemes)),
			Servers:         make([]ServerInfo, len(vResult.SpecInfo.Servers)),
		}

		// Convert security schemes
		for i, scheme := range vResult.SpecInfo.SecuritySchemes {
			result.SpecInfo.SecuritySchemes[i] = SecurityScheme{
				Type:        scheme.Type,
				Name:        scheme.Name,
				Description: scheme.Description,
			}
		}

		// Convert servers
		for i, server := range vResult.SpecInfo.Servers {
			result.SpecInfo.Servers[i] = ServerInfo{
				URL:         server.URL,
				Description: server.Description,
			}
		}
	}

	return result
}

// cacheValidationResult caches a validation result
func (a *App) cacheValidationResult(specPath, specURL string, result *ValidationResult) {
	// Serialize the result
	resultJSON, err := json.Marshal(result)
	if err != nil {
		return // Skip caching if serialization fails
	}

	// Create cache entry
	cache := &database.ValidationCache{
		SpecPath:         specPath,
		SpecURL:          specURL,
		ValidationResult: string(resultJSON),
		CachedAt:         time.Now(),
		ExpiresAt:        time.Now().Add(24 * time.Hour), // Cache for 24 hours
	}

	// Generate hash based on source type
	if specPath != "" {
		if hash, err := a.validationCacheRepo.GenerateSpecHash(specPath); err == nil {
			cache.SpecHash = hash
		}
	} else if specURL != "" {
		cache.SpecHash = a.validationCacheRepo.GenerateURLHash(specURL, time.Now())
	}

	// Store in cache
	a.validationCacheRepo.Store(cache)
}

// ExportValidationResult exports a validation result to JSON
func (a *App) ExportValidationResult(result *ValidationResult) (string, error) {
	if result == nil {
		return "", a.createAPIError("validation", ErrCodeValidation, "Validation result is required", nil)
	}

	// Create export data
	exportData := map[string]interface{}{
		"version":          "1.0.0",
		"validationResult": result,
		"exportedAt":       time.Now().Format(time.RFC3339),
	}

	// Marshal to JSON
	jsonData, err := json.MarshalIndent(exportData, "", "  ")
	if err != nil {
		return "", a.createAPIError("internal", ErrCodeInternalError, "Failed to export validation result", map[string]string{
			"error": err.Error(),
		})
	}

	return string(jsonData), nil
}

// GetValidationCacheStats returns cache statistics
func (a *App) GetValidationCacheStats() (*database.ValidationCacheStats, error) {
	if a.validationCacheRepo == nil {
		return nil, a.createAPIError("internal", ErrCodeInternalError, "Validation cache not initialized", nil)
	}

	stats, err := a.validationCacheRepo.GetStats()
	if err != nil {
		return nil, a.createAPIError("internal", ErrCodeInternalError, "Failed to get cache stats", map[string]string{
			"error": err.Error(),
		})
	}

	return stats, nil
}

// ClearValidationCache clears expired validation cache entries
func (a *App) ClearValidationCache() error {
	if a.validationCacheRepo == nil {
		return a.createAPIError("internal", ErrCodeInternalError, "Validation cache not initialized", nil)
	}

	err := a.validationCacheRepo.CleanExpired()
	if err != nil {
		return a.createAPIError("internal", ErrCodeInternalError, "Failed to clear cache", map[string]string{
			"error": err.Error(),
		})
	}

	return nil
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
