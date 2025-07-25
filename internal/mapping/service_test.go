package mapping

import (
	"testing"

	"MCPWeaver/internal/parser"
	"github.com/getkin/kin-openapi/openapi3"
)

func TestNewService(t *testing.T) {
	service := NewService("https://api.example.com")
	if service == nil {
		t.Fatal("NewService returned nil")
	}
	if service.baseURL != "https://api.example.com" {
		t.Errorf("Expected baseURL 'https://api.example.com', got %s", service.baseURL)
	}
}

func TestGenerateToolName(t *testing.T) {
	service := NewService("https://api.example.com")

	tests := []struct {
		operation parser.Operation
		expected  string
	}{
		{
			operation: parser.Operation{ID: "getUserById", Method: "GET", Path: "/users/{id}"},
			expected:  "getUserById",
		},
		{
			operation: parser.Operation{ID: "", Method: "GET", Path: "/users"},
			expected:  "getUsers",
		},
		{
			operation: parser.Operation{ID: "", Method: "POST", Path: "/users/{id}/profile"},
			expected:  "postUsersIdProfile",
		},
	}

	for _, tt := range tests {
		t.Run(tt.expected, func(t *testing.T) {
			result := service.generateToolName(tt.operation)
			if result != tt.expected {
				t.Errorf("generateToolName() = %s, want %s", result, tt.expected)
			}
		})
	}
}

func TestGenerateToolDescription(t *testing.T) {
	service := NewService("https://api.example.com")

	tests := []struct {
		operation parser.Operation
		expected  string
	}{
		{
			operation: parser.Operation{Description: "Get user by ID", Method: "GET", Path: "/users/{id}"},
			expected:  "Get user by ID",
		},
		{
			operation: parser.Operation{Description: "", Summary: "Get all users", Method: "GET", Path: "/users"},
			expected:  "Get all users",
		},
		{
			operation: parser.Operation{Description: "", Summary: "", Method: "POST", Path: "/users"},
			expected:  "Create users",
		},
	}

	for _, tt := range tests {
		t.Run(tt.expected, func(t *testing.T) {
			result := service.generateToolDescription(tt.operation)
			if result != tt.expected {
				t.Errorf("generateToolDescription() = %s, want %s", result, tt.expected)
			}
		})
	}
}

func TestGetActionFromMethod(t *testing.T) {
	service := NewService("https://api.example.com")

	tests := []struct {
		method   string
		expected string
	}{
		{"GET", "Retrieve"},
		{"POST", "Create"},
		{"PUT", "Update"},
		{"PATCH", "Modify"},
		{"DELETE", "Delete"},
		{"OPTIONS", "Perform operation on"},
	}

	for _, tt := range tests {
		t.Run(tt.method, func(t *testing.T) {
			result := service.getActionFromMethod(tt.method)
			if result != tt.expected {
				t.Errorf("getActionFromMethod(%s) = %s, want %s", tt.method, result, tt.expected)
			}
		})
	}
}

func TestGetResourceFromPath(t *testing.T) {
	service := NewService("https://api.example.com")

	tests := []struct {
		path     string
		expected string
	}{
		{"/users", "users"},
		{"/users/{id}", "users"},
		{"/api/v1/users/{id}/profile", "profile"},
		{"/api/v1/users/{id}/profile/{type}", "profile"},
		{"/", "resource"},
	}

	for _, tt := range tests {
		t.Run(tt.path, func(t *testing.T) {
			result := service.getResourceFromPath(tt.path)
			if result != tt.expected {
				t.Errorf("getResourceFromPath(%s) = %s, want %s", tt.path, result, tt.expected)
			}
		})
	}
}

func TestMapParameterToProperty(t *testing.T) {
	service := NewService("https://api.example.com")

	// Create a test parameter
	param := parser.Parameter{
		Name:        "userId",
		In:          "path",
		Description: "User ID",
		Required:    true,
		Schema: &openapi3.SchemaRef{
			Value: &openapi3.Schema{
				Type:   &openapi3.Types{"string"},
				Format: "uuid",
			},
		},
		Example: "123e4567-e89b-12d3-a456-426614174000",
	}

	result, err := service.mapParameterToProperty(param)
	if err != nil {
		t.Fatalf("mapParameterToProperty failed: %v", err)
	}

	if result.Type != "string" {
		t.Errorf("Expected type 'string', got %s", result.Type)
	}

	if result.Format != "uuid" {
		t.Errorf("Expected format 'uuid', got %s", result.Format)
	}

	if result.Description != "User ID" {
		t.Errorf("Expected description 'User ID', got %s", result.Description)
	}
}

func TestMapOperationToTool(t *testing.T) {
	service := NewService("https://api.example.com")

	// Create a test operation
	operation := parser.Operation{
		ID:          "getUser",
		Method:      "GET",
		Path:        "/users/{id}",
		Summary:     "Get user by ID",
		Description: "Retrieve a user by their ID",
		Parameters: []parser.Parameter{
			{
				Name:        "id",
				In:          "path",
				Description: "User ID",
				Required:    true,
				Schema: &openapi3.SchemaRef{
					Value: &openapi3.Schema{
						Type: &openapi3.Types{"string"},
					},
				},
			},
		},
	}

	result, err := service.mapOperationToTool(operation)
	if err != nil {
		t.Fatalf("mapOperationToTool failed: %v", err)
	}

	if result.Name != "getUser" {
		t.Errorf("Expected name 'getUser', got %s", result.Name)
	}

	if result.Method != "GET" {
		t.Errorf("Expected method 'GET', got %s", result.Method)
	}

	if result.Path != "/users/{id}" {
		t.Errorf("Expected path '/users/{id}', got %s", result.Path)
	}

	if result.BaseURL != "https://api.example.com" {
		t.Errorf("Expected baseURL 'https://api.example.com', got %s", result.BaseURL)
	}

	if len(result.InputSchema.Properties) != 1 {
		t.Errorf("Expected 1 property, got %d", len(result.InputSchema.Properties))
	}

	if len(result.InputSchema.Required) != 1 {
		t.Errorf("Expected 1 required parameter, got %d", len(result.InputSchema.Required))
	}

	if result.InputSchema.Required[0] != "id" {
		t.Errorf("Expected required parameter 'id', got %s", result.InputSchema.Required[0])
	}
}
