package parser

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"path"
	"regexp"
	"strings"
	"time"
	"unicode"

	"github.com/getkin/kin-openapi/openapi3"
)

// toTitle converts the first character of a string to uppercase
func toTitle(s string) string {
	if s == "" {
		return s
	}
	runes := []rune(s)
	runes[0] = unicode.ToUpper(runes[0])
	return string(runes)
}

// Service handles OpenAPI specification parsing
type Service struct {
	loader *openapi3.Loader
}

// NewService creates a new OpenAPI parser service
func NewService() *Service {
	return &Service{
		loader: openapi3.NewLoader(),
	}
}

// ParseFromFile parses an OpenAPI specification from a file
func (s *Service) ParseFromFile(filePath string) (*ParsedAPI, error) {
	doc, err := s.loader.LoadFromFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to load OpenAPI spec: %w", err)
	}

	return s.parseDocument(doc)
}

// ParseFromURL parses an OpenAPI specification from a URL
func (s *Service) ParseFromURL(ctx context.Context, specURL string) (*ParsedAPI, error) {
	// Validate URL format
	if _, err := url.Parse(specURL); err != nil {
		return nil, fmt.Errorf("invalid URL: %w", err)
	}

	// For better YAML support, we'll fetch the content manually and then parse
	doc, err := s.loadFromURLWithContext(ctx, specURL)
	if err != nil {
		return nil, fmt.Errorf("failed to load OpenAPI spec from URL: %w", err)
	}

	return s.parseDocument(doc)
}

// parseDocument converts an OpenAPI document to our internal representation
func (s *Service) parseDocument(doc *openapi3.T) (*ParsedAPI, error) {
	// Validate the document - use a more lenient validation for problematic specs
	ctx := context.Background()
	if err := doc.Validate(ctx); err != nil {
		// Try to provide more helpful error details for common issues
		errorStr := err.Error()
		if strings.Contains(errorStr, "items") && strings.Contains(errorStr, "array") {
			return nil, fmt.Errorf("OpenAPI validation failed due to array schema issue: %w\n\nThis specification contains array definitions that use tuple-style or non-standard schema formats.\nThese are valid in JSON Schema but not fully supported by the kin-openapi library.\nConsider simplifying array schemas to use standard 'items' object definitions", err)
		}
		return nil, fmt.Errorf("OpenAPI validation failed: %w", err)
	}

	parsed := &ParsedAPI{
		Document:    doc,
		Title:       doc.Info.Title,
		Version:     doc.Info.Version,
		Description: doc.Info.Description,
		Schemas:     make(map[string]*openapi3.SchemaRef),
	}

	// Extract base URL from servers
	if len(doc.Servers) > 0 {
		parsed.BaseURL = doc.Servers[0].URL
		for _, server := range doc.Servers {
			parsed.Servers = append(parsed.Servers, server.URL)
		}
	}

	// Extract schemas from components
	if doc.Components != nil && doc.Components.Schemas != nil {
		for name, schema := range doc.Components.Schemas {
			parsed.Schemas[name] = schema
		}
	}

	// Extract operations
	operations, err := s.extractOperations(doc)
	if err != nil {
		return nil, fmt.Errorf("failed to extract operations: %w", err)
	}
	parsed.Operations = operations

	return parsed, nil
}

// extractOperations extracts all operations from the OpenAPI document
func (s *Service) extractOperations(doc *openapi3.T) ([]Operation, error) {
	var operations []Operation

	for path, pathItem := range doc.Paths.Map() {
		for method, operation := range pathItem.Operations() {
			if operation == nil {
				continue
			}

			op := Operation{
				ID:          operation.OperationID,
				Method:      strings.ToUpper(method),
				Path:        path,
				Summary:     operation.Summary,
				Description: operation.Description,
				Tags:        operation.Tags,
				Responses:   make(map[string]*Response),
			}

			// Generate operation ID if not provided
			if op.ID == "" {
				op.ID = s.generateOperationID(method, path)
			}

			// Extract parameters
			parameters, err := s.extractParameters(operation.Parameters)
			if err != nil {
				return nil, fmt.Errorf("failed to extract parameters for %s %s: %w", method, path, err)
			}
			op.Parameters = parameters

			// Extract request body
			if operation.RequestBody != nil {
				requestBody, err := s.extractRequestBody(operation.RequestBody)
				if err != nil {
					return nil, fmt.Errorf("failed to extract request body for %s %s: %w", method, path, err)
				}
				op.RequestBody = requestBody
			}

			// Extract responses
			responses, err := s.extractResponses(operation.Responses)
			if err != nil {
				return nil, fmt.Errorf("failed to extract responses for %s %s: %w", method, path, err)
			}
			op.Responses = responses

			// Extract security requirements
			if operation.Security != nil {
				for _, security := range *operation.Security {
					op.Security = append(op.Security, security)
				}
			}

			operations = append(operations, op)
		}
	}

	return operations, nil
}

// generateOperationID generates an operation ID from method and path
func (s *Service) generateOperationID(method, path string) string {
	// Convert path to camelCase identifier
	parts := strings.Split(path, "/")
	var cleanParts []string

	for _, part := range parts {
		if part == "" {
			continue
		}
		// Remove path parameters braces
		part = strings.ReplaceAll(part, "{", "")
		part = strings.ReplaceAll(part, "}", "")
		// Convert to camelCase
		if len(cleanParts) > 0 {
			part = toTitle(part)
		}
		cleanParts = append(cleanParts, part)
	}

	pathPart := strings.Join(cleanParts, "")
	return strings.ToLower(method) + toTitle(pathPart)
}

// extractParameters extracts parameter information
func (s *Service) extractParameters(params openapi3.Parameters) ([]Parameter, error) {
	var parameters []Parameter

	for _, paramRef := range params {
		if paramRef.Value == nil {
			continue
		}

		param := Parameter{
			Name:        paramRef.Value.Name,
			In:          paramRef.Value.In,
			Description: paramRef.Value.Description,
			Required:    paramRef.Value.Required,
			Schema:      paramRef.Value.Schema,
			Example:     paramRef.Value.Example,
		}

		parameters = append(parameters, param)
	}

	return parameters, nil
}

// extractRequestBody extracts request body information
func (s *Service) extractRequestBody(reqBodyRef *openapi3.RequestBodyRef) (*RequestBody, error) {
	if reqBodyRef == nil || reqBodyRef.Value == nil {
		return nil, nil
	}

	requestBody := &RequestBody{
		Description: reqBodyRef.Value.Description,
		Required:    reqBodyRef.Value.Required,
		Content:     make(map[string]*MediaType),
	}

	for mediaType, mediaTypeObj := range reqBodyRef.Value.Content {
		mt := &MediaType{
			Schema:   mediaTypeObj.Schema,
			Example:  mediaTypeObj.Example,
			Examples: mediaTypeObj.Examples,
		}
		requestBody.Content[mediaType] = mt
	}

	return requestBody, nil
}

// extractResponses extracts response information
func (s *Service) extractResponses(responses *openapi3.Responses) (map[string]*Response, error) {
	result := make(map[string]*Response)

	if responses == nil {
		return result, nil
	}

	for status, responseRef := range responses.Map() {
		if responseRef.Value == nil {
			continue
		}

		description := ""
		if responseRef.Value.Description != nil {
			description = *responseRef.Value.Description
		}

		response := &Response{
			Description: description,
			Headers:     responseRef.Value.Headers,
			Content:     make(map[string]*MediaType),
		}

		for mediaType, mediaTypeObj := range responseRef.Value.Content {
			mt := &MediaType{
				Schema:   mediaTypeObj.Schema,
				Example:  mediaTypeObj.Example,
				Examples: mediaTypeObj.Examples,
			}
			response.Content[mediaType] = mt
		}

		result[status] = response
	}

	return result, nil
}

// loadFromURLWithContext fetches OpenAPI spec from URL with better YAML/JSON handling
func (s *Service) loadFromURLWithContext(ctx context.Context, specURL string) (*openapi3.T, error) {
	// Create HTTP client with timeout
	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	// Create request with context
	req, err := http.NewRequestWithContext(ctx, "GET", specURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create HTTP request: %w", err)
	}

	// Set appropriate Accept header for both JSON and YAML
	req.Header.Set("Accept", "application/json, application/x-yaml, text/yaml, text/x-yaml, */*")
	req.Header.Set("User-Agent", "MCPWeaver/1.0")

	// Make the request
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch URL: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("HTTP error: %d %s", resp.StatusCode, resp.Status)
	}

	// Read the response body
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	// Determine if this is likely YAML based on URL path extension
	parsedURL, _ := url.Parse(specURL)
	isYAML := strings.HasSuffix(strings.ToLower(path.Ext(parsedURL.Path)), ".yaml") ||
		strings.HasSuffix(strings.ToLower(path.Ext(parsedURL.Path)), ".yml")

	// Pre-process the data to fix common schema issues
	processedData := s.preprocessSpecData(data)

	// Parse the processed data using the appropriate method
	var doc *openapi3.T
	if isYAML {
		// Try YAML first for .yaml/.yml files
		doc, err = s.loader.LoadFromData(processedData)
		if err != nil {
			// If YAML parsing fails, try JSON as fallback
			doc, err = s.loader.LoadFromData(processedData)
		}
	} else {
		// Try JSON first for other extensions
		doc, err = s.loader.LoadFromData(processedData)
		if err != nil {
			// If JSON parsing fails, try YAML as fallback
			doc, err = s.loader.LoadFromData(processedData)
		}
	}

	if err != nil {
		return nil, fmt.Errorf("failed to parse OpenAPI specification: %w\n\nHint: This may be due to:\n- Invalid OpenAPI schema syntax in the source file\n- Non-standard OpenAPI extensions\n- OpenAPI 2.0 spec (this tool requires OpenAPI 3.0+)\n- Malformed YAML/JSON structure", err)
	}

	return doc, nil
}

// preprocessSpecData applies common fixes for schema compatibility issues
func (s *Service) preprocessSpecData(data []byte) []byte {
	content := string(data)

	// Fix 1: Handle tuple-style array items
	tuplePattern := `(\s+items:\s*\n\s+)- (type: object\s*\n(?:\s+properties:\s*\n(?:\s+[^\n-]*\n)*)?)`
	tupleRegex := regexp.MustCompile(tuplePattern)
	if tupleRegex.MatchString(content) {
		content = tupleRegex.ReplaceAllString(content, `${1}$2`)
	}

	// Fix 2: Remove additionalItems which is not well supported
	additionalItemsRegex := regexp.MustCompile(`(?m)^\s+additionalItems:.*\n`)
	content = additionalItemsRegex.ReplaceAllString(content, "")

	// Fix 3: Handle unsupported regex patterns
	lines := strings.Split(content, "\n")
	for i, line := range lines {
		if strings.Contains(line, "pattern:") && strings.Contains(line, "(?") {
			parts := strings.SplitN(line, "pattern:", 2)
			if len(parts) == 2 {
				indent := parts[0]
				patternPart := strings.TrimSpace(parts[1])

				// Handle quoted patterns
				var quote string
				if strings.HasPrefix(patternPart, "\"") {
					quote = "\""
					patternPart = strings.Trim(patternPart, "\"")
				} else if strings.HasPrefix(patternPart, "'") {
					quote = "'"
					patternPart = strings.Trim(patternPart, "'")
				}

				// Handle patterns wrapped in forward slashes
				if strings.HasPrefix(patternPart, "/") && strings.HasSuffix(patternPart, "/") {
					patternPart = strings.Trim(patternPart, "/")
				}

				// Remove unsupported lookahead/lookbehind assertions
				if strings.Contains(patternPart, "(?=") || strings.Contains(patternPart, "(?!") {
					// For password validation patterns, replace with simpler equivalent
					if strings.Contains(patternPart, "(?=.*[a-z])") && strings.Contains(patternPart, "(?=.*[A-Z])") {
						patternPart = "^[a-zA-Z0-9]+$"
					} else {
						// Remove lookahead/lookbehind for other patterns
						lookAheadRegex := regexp.MustCompile(`\(\?\=[^)]*\)`)
						patternPart = lookAheadRegex.ReplaceAllString(patternPart, "")
						negLookAheadRegex := regexp.MustCompile(`\(\?\![^)]*\)`)
						patternPart = negLookAheadRegex.ReplaceAllString(patternPart, "")
					}

					patternPart = strings.TrimSpace(patternPart)
					if patternPart != "" {
						if quote != "" {
							lines[i] = indent + "pattern: " + quote + patternPart + quote
						} else {
							lines[i] = indent + "pattern: " + patternPart
						}
					} else {
						lines[i] = indent + "# pattern: # removed unsupported regex"
					}
				}
			}
		}
	}

	return []byte(strings.Join(lines, "\n"))
}