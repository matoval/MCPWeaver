package parser

import (
	"strings"
	"testing"

	"github.com/getkin/kin-openapi/openapi3"
)

func TestNewService(t *testing.T) {
	service := NewService()
	if service == nil {
		t.Fatal("NewService returned nil")
	}
	if service.loader == nil {
		t.Fatal("Service loader is nil")
	}
}

func TestGenerateOperationID(t *testing.T) {
	service := NewService()
	
	tests := []struct {
		method   string
		path     string
		expected string
	}{
		{"GET", "/users", "getUsers"},
		{"POST", "/users", "postUsers"},
		{"GET", "/users/{id}", "getUsersId"},
		{"PUT", "/users/{id}/profile", "putUsersIdProfile"},
		{"DELETE", "/api/v1/users/{id}", "deleteApiV1UsersId"},
	}

	for _, tt := range tests {
		t.Run(tt.method+"_"+tt.path, func(t *testing.T) {
			result := service.generateOperationID(tt.method, tt.path)
			if result != tt.expected {
				t.Errorf("generateOperationID(%s, %s) = %s, want %s", tt.method, tt.path, result, tt.expected)
			}
		})
	}
}

func TestPreprocessSpecData(t *testing.T) {
	service := NewService()
	
	input := `
openapi: 3.0.0
info:
  title: Test API
  version: 1.0.0
paths:
  /test:
    get:
      parameters:
        - name: test_param
          in: query
          schema:
            type: string
            pattern: '/(?=.*[a-z])(?=.*[A-Z])(?=.*[0-9])/'
`
	
	result := service.preprocessSpecData([]byte(input))
	resultStr := string(result)
	
	// Check that the problematic regex pattern was fixed
	if strings.Contains(resultStr, "(?=") {
		t.Error("Regex pattern was not fixed")
	}
	
	// Check that it was replaced with something reasonable
	if !strings.Contains(resultStr, "^[a-zA-Z0-9]+$") {
		t.Error("Regex pattern was not replaced with expected value")
	}
}

func TestParseDocument(t *testing.T) {
	service := NewService()
	
	// Create a minimal valid OpenAPI document
	doc := &openapi3.T{
		OpenAPI: "3.0.0",
		Info: &openapi3.Info{
			Title:   "Test API",
			Version: "1.0.0",
		},
		Paths: &openapi3.Paths{},
	}
	
	result, err := service.parseDocument(doc)
	if err != nil {
		t.Fatalf("parseDocument failed: %v", err)
	}
	
	if result.Title != "Test API" {
		t.Errorf("Expected title 'Test API', got %s", result.Title)
	}
	
	if result.Version != "1.0.0" {
		t.Errorf("Expected version '1.0.0', got %s", result.Version)
	}
}

func TestExtractOperations(t *testing.T) {
	service := NewService()
	
	// Create a minimal document with one operation
	doc := &openapi3.T{
		OpenAPI: "3.0.0",
		Info: &openapi3.Info{
			Title:   "Test API",
			Version: "1.0.0",
		},
		Paths: &openapi3.Paths{},
	}
	
	// Add a simple GET operation
	pathItem := &openapi3.PathItem{
		Get: &openapi3.Operation{
			OperationID: "getUsers",
			Summary:     "Get all users",
			Description: "Returns a list of all users",
			Responses: &openapi3.Responses{
				Extensions: make(map[string]interface{}),
			},
		},
	}
	
	doc.Paths = &openapi3.Paths{
		Extensions: make(map[string]interface{}),
	}
	doc.Paths.Set("/users", pathItem)
	
	operations, err := service.extractOperations(doc)
	if err != nil {
		t.Fatalf("extractOperations failed: %v", err)
	}
	
	if len(operations) != 1 {
		t.Fatalf("Expected 1 operation, got %d", len(operations))
	}
	
	op := operations[0]
	if op.ID != "getUsers" {
		t.Errorf("Expected operation ID 'getUsers', got %s", op.ID)
	}
	
	if op.Method != "GET" {
		t.Errorf("Expected method 'GET', got %s", op.Method)
	}
	
	if op.Path != "/users" {
		t.Errorf("Expected path '/users', got %s", op.Path)
	}
}