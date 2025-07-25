package mapping

import (
	"fmt"
	"strings"
	"unicode"

	"MCPWeaver/internal/parser"
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

// MCPTool represents an MCP tool definition
type MCPTool struct {
	Name        string      `json:"name"`
	Description string      `json:"description"`
	InputSchema InputSchema `json:"inputSchema"`
	Method      string      `json:"-"` // HTTP method for implementation
	Path        string      `json:"-"` // API path for implementation
	BaseURL     string      `json:"-"` // Base URL for API calls
}

// InputSchema represents the JSON schema for tool input parameters
type InputSchema struct {
	Type       string              `json:"type"`
	Properties map[string]Property `json:"properties"`
	Required   []string            `json:"required"`
}

// Property represents a property in the input schema
type Property struct {
	Type        string    `json:"type"`
	Description string    `json:"description,omitempty"`
	Example     any       `json:"example,omitempty"`
	Enum        []string  `json:"enum,omitempty"`
	Format      string    `json:"format,omitempty"`
	Items       *Property `json:"items,omitempty"`
}

// Service handles the conversion from OpenAPI operations to MCP tools
type Service struct {
	baseURL string
}

// NewService creates a new endpoint mapping service
func NewService(baseURL string) *Service {
	return &Service{
		baseURL: baseURL,
	}
}

// MapOperationsToTools converts OpenAPI operations to MCP tools
func (s *Service) MapOperationsToTools(operations []parser.Operation) ([]MCPTool, error) {
	var tools []MCPTool
	toolNames := make(map[string]int) // Track duplicate names

	for _, op := range operations {
		tool, err := s.mapOperationToTool(op)
		if err != nil {
			return nil, fmt.Errorf("failed to map operation %s: %w", op.ID, err)
		}

		// Handle duplicate tool names
		if count, exists := toolNames[tool.Name]; exists {
			toolNames[tool.Name] = count + 1
			tool.Name = fmt.Sprintf("%s_%d", tool.Name, count+1)
		} else {
			toolNames[tool.Name] = 1
		}

		tools = append(tools, tool)
	}

	return tools, nil
}

// mapOperationToTool converts a single OpenAPI operation to an MCP tool
func (s *Service) mapOperationToTool(op parser.Operation) (MCPTool, error) {
	tool := MCPTool{
		Name:        s.generateToolName(op),
		Description: s.generateToolDescription(op),
		Method:      op.Method,
		Path:        op.Path,
		BaseURL:     s.baseURL,
		InputSchema: InputSchema{
			Type:       "object",
			Properties: make(map[string]Property),
			Required:   []string{},
		},
	}

	// Map parameters to input schema properties
	for _, param := range op.Parameters {
		property, err := s.mapParameterToProperty(param)
		if err != nil {
			return tool, fmt.Errorf("failed to map parameter %s: %w", param.Name, err)
		}

		tool.InputSchema.Properties[param.Name] = property

		if param.Required {
			tool.InputSchema.Required = append(tool.InputSchema.Required, param.Name)
		}
	}

	// Handle request body for POST/PUT operations
	if op.RequestBody != nil && (op.Method == "POST" || op.Method == "PUT" || op.Method == "PATCH") {
		bodyProperty, err := s.mapRequestBodyToProperty(op.RequestBody)
		if err != nil {
			return tool, fmt.Errorf("failed to map request body: %w", err)
		}

		tool.InputSchema.Properties["body"] = bodyProperty

		if op.RequestBody.Required {
			tool.InputSchema.Required = append(tool.InputSchema.Required, "body")
		}
	}

	return tool, nil
}

// generateToolName creates a tool name from the operation
func (s *Service) generateToolName(op parser.Operation) string {
	if op.ID != "" {
		return op.ID
	}

	// Generate from method and path
	parts := strings.Split(op.Path, "/")
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
	return strings.ToLower(op.Method) + toTitle(pathPart)
}

// generateToolDescription creates a description for the tool
func (s *Service) generateToolDescription(op parser.Operation) string {
	if op.Description != "" {
		return op.Description
	}

	if op.Summary != "" {
		return op.Summary
	}

	// Generate generic description
	action := s.getActionFromMethod(op.Method)
	resource := s.getResourceFromPath(op.Path)

	return fmt.Sprintf("%s %s", action, resource)
}

// getActionFromMethod returns a human-readable action for the HTTP method
func (s *Service) getActionFromMethod(method string) string {
	switch strings.ToUpper(method) {
	case "GET":
		return "Retrieve"
	case "POST":
		return "Create"
	case "PUT":
		return "Update"
	case "PATCH":
		return "Modify"
	case "DELETE":
		return "Delete"
	default:
		return "Perform operation on"
	}
}

// getResourceFromPath extracts the resource name from the path
func (s *Service) getResourceFromPath(path string) string {
	parts := strings.Split(path, "/")
	for i := len(parts) - 1; i >= 0; i-- {
		part := parts[i]
		if part != "" && !strings.Contains(part, "{") {
			return part
		}
	}
	return "resource"
}

// mapParameterToProperty converts an OpenAPI parameter to a JSON schema property
func (s *Service) mapParameterToProperty(param parser.Parameter) (Property, error) {
	property := Property{
		Description: param.Description,
		Example:     param.Example,
	}

	// Handle missing schema gracefully
	if param.Schema == nil || param.Schema.Value == nil {
		property.Type = "string"
		return property, nil
	}

	schema := param.Schema.Value

	// Map basic type with error handling
	if schema.Type != nil && len(*schema.Type) > 0 {
		schemaType := (*schema.Type)[0]
		// Validate type is supported
		supportedTypes := map[string]bool{
			"string": true, "number": true, "integer": true,
			"boolean": true, "array": true, "object": true,
		}
		if !supportedTypes[schemaType] {
			return property, fmt.Errorf("unsupported parameter type: %s", schemaType)
		}
		property.Type = schemaType
	} else {
		property.Type = "string"
	}

	// Map format with validation
	if schema.Format != "" {
		property.Format = schema.Format
	}

	// Map enum values with type safety
	if len(schema.Enum) > 0 {
		for _, enumVal := range schema.Enum {
			if str, ok := enumVal.(string); ok {
				property.Enum = append(property.Enum, str)
			}
		}
	}

	// Handle array types with proper validation
	if property.Type == "array" && schema.Items != nil && schema.Items.Value != nil {
		items := &Property{}
		if schema.Items.Value.Type != nil && len(*schema.Items.Value.Type) > 0 {
			items.Type = (*schema.Items.Value.Type)[0]
		} else {
			items.Type = "string"
		}
		if schema.Items.Value.Format != "" {
			items.Format = schema.Items.Value.Format
		}
		property.Items = items
	}

	return property, nil
}

// mapRequestBodyToProperty converts a request body to a JSON schema property
func (s *Service) mapRequestBodyToProperty(reqBody *parser.RequestBody) (Property, error) {
	property := Property{
		Type:        "object",
		Description: reqBody.Description,
	}

	// For now, we'll treat all request bodies as generic objects
	// In a more sophisticated implementation, we would parse the actual schema
	if len(reqBody.Content) > 0 {
		// Look for JSON content first
		if jsonContent, exists := reqBody.Content["application/json"]; exists {
			if jsonContent.Schema != nil && jsonContent.Schema.Value != nil {
				// This would require more complex schema parsing
				// For now, keep as generic object
				property.Description = fmt.Sprintf("%s (JSON object)", property.Description)
			}
		}
	}

	return property, nil
}
