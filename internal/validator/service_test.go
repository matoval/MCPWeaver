package validator

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/getkin/kin-openapi/openapi3"
)

func TestService_ValidateFile(t *testing.T) {
	service := New()
	ctx := context.Background()

	tests := []struct {
		name        string
		specContent string
		filename    string
		wantValid   bool
		wantErrors  int
		wantWarnings int
	}{
		{
			name: "valid OpenAPI spec",
			specContent: `
openapi: 3.0.0
info:
  title: Test API
  version: 1.0.0
  description: A test API
paths:
  /test:
    get:
      summary: Test endpoint
      operationId: getTest
      responses:
        '200':
          description: Success
          content:
            application/json:
              schema:
                type: object
                properties:
                  message:
                    type: string
components:
  schemas:
    TestSchema:
      type: object
      description: A test schema
      properties:
        id:
          type: integer
        name:
          type: string
      required:
        - id
`,
			filename:    "valid_spec.yaml",
			wantValid:   true,
			wantErrors:  0,
			wantWarnings: 0,
		},
		{
			name: "spec with missing description",
			specContent: `
openapi: 3.0.0
info:
  title: Test API
  version: 1.0.0
paths:
  /test:
    get:
      operationId: getTest
      responses:
        '200':
          description: Success
`,
			filename:     "missing_desc_spec.yaml",
			wantValid:    true,
			wantErrors:   0,
			wantWarnings: 2, // Missing operation description and no schemas
		},
		{
			name: "spec with duplicate operation ID",
			specContent: `
openapi: 3.0.0
info:
  title: Test API
  version: 1.0.0
paths:
  /test:
    get:
      operationId: duplicateId
      responses:
        '200':
          description: Success
  /test2:
    get:
      operationId: duplicateId
      responses:
        '200':
          description: Success
`,
			filename:     "duplicate_id_spec.yaml",
			wantValid:    false,
			wantErrors:   2, // Duplicate operation ID might be reported twice
			wantWarnings: 3, // Missing operation descriptions and no schemas
		},
		{
			name: "invalid OpenAPI spec",
			specContent: `
openapi: 3.0.0
info:
  title: Test API
  # Missing version
paths:
  /test:
    get:
      responses:
        '200':
          description: Success
`,
			filename:     "invalid_spec.yaml",
			wantValid:    false,
			wantErrors:   1,
			wantWarnings: 2, // Missing operation description and no schemas
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create temporary file
			tempDir := t.TempDir()
			filePath := filepath.Join(tempDir, tt.filename)
			
			err := os.WriteFile(filePath, []byte(tt.specContent), 0644)
			if err != nil {
				t.Fatalf("Failed to write test file: %v", err)
			}

			// Test validation
			result, err := service.ValidateFile(ctx, filePath)
			if err != nil {
				t.Fatalf("ValidateFile() error = %v", err)
			}

			if result.Valid != tt.wantValid {
				t.Errorf("ValidateFile() valid = %v, want %v", result.Valid, tt.wantValid)
			}

			if len(result.Errors) != tt.wantErrors {
				t.Errorf("ValidateFile() errors = %d, want %d", len(result.Errors), tt.wantErrors)
			}

			if len(result.Warnings) != tt.wantWarnings {
				t.Errorf("ValidateFile() warnings = %d, want %d", len(result.Warnings), tt.wantWarnings)
			}

			// Check that timing information is set
			if result.ValidationTime == 0 {
				t.Error("ValidateFile() validation time should be set")
			}

			if result.ValidatedAt.IsZero() {
				t.Error("ValidateFile() validated at should be set")
			}
		})
	}
}

func TestService_ValidateFile_FileNotFound(t *testing.T) {
	service := New()
	ctx := context.Background()

	result, err := service.ValidateFile(ctx, "/nonexistent/file.yaml")
	if err != nil {
		t.Fatalf("ValidateFile() error = %v", err)
	}

	if result.Valid {
		t.Error("ValidateFile() should return invalid for non-existent file")
	}

	if len(result.Errors) != 1 {
		t.Errorf("ValidateFile() errors = %d, want 1", len(result.Errors))
	}

	if result.Errors[0].Type != "file" {
		t.Errorf("ValidateFile() error type = %s, want 'file'", result.Errors[0].Type)
	}
}

func TestService_ValidateURL(t *testing.T) {
	service := New()
	ctx := context.Background()

	tests := []struct {
		name      string
		url       string
		wantValid bool
		wantError bool
	}{
		{
			name:      "invalid URL format",
			url:       "not-a-url",
			wantValid: false,
			wantError: false,
		},
		{
			name:      "empty URL",
			url:       "",
			wantValid: false,
			wantError: false,
		},
		{
			name:      "valid URL format but unreachable",
			url:       "https://example.com/nonexistent.yaml",
			wantValid: false,
			wantError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := service.ValidateURL(ctx, tt.url)
			
			if (err != nil) != tt.wantError {
				t.Errorf("ValidateURL() error = %v, wantError %v", err, tt.wantError)
			}

			if result.Valid != tt.wantValid {
				t.Errorf("ValidateURL() valid = %v, want %v", result.Valid, tt.wantValid)
			}
		})
	}
}

func TestService_ValidateSpec(t *testing.T) {
	service := New()
	ctx := context.Background()

	// Create a valid OpenAPI spec
	spec := &openapi3.T{
		OpenAPI: "3.0.0",
		Info: &openapi3.Info{
			Title:   "Test API",
			Version: "1.0.0",
		},
		Paths: &openapi3.Paths{},
	}

	result, err := service.ValidateSpec(ctx, spec)
	if err != nil {
		t.Fatalf("ValidateSpec() error = %v", err)
	}

	if !result.Valid {
		t.Error("ValidateSpec() should return valid for basic spec")
	}

	if result.SpecInfo == nil {
		t.Error("ValidateSpec() should return spec info")
	}

	if result.SpecInfo.Title != "Test API" {
		t.Errorf("ValidateSpec() spec title = %s, want 'Test API'", result.SpecInfo.Title)
	}
}

func TestService_ComplexityAssessment(t *testing.T) {
	service := New()

	tests := []struct {
		name           string
		operationCount int
		schemaCount    int
		wantComplexity ComplexityLevel
	}{
		{
			name:           "low complexity",
			operationCount: 5,
			schemaCount:    3,
			wantComplexity: ComplexityLow,
		},
		{
			name:           "medium complexity",
			operationCount: 15,
			schemaCount:    10,
			wantComplexity: ComplexityMedium,
		},
		{
			name:           "high complexity",
			operationCount: 80,
			schemaCount:    50,
			wantComplexity: ComplexityHigh,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a spec with the desired counts
			spec := &openapi3.T{
				OpenAPI: "3.0.0",
				Info: &openapi3.Info{
					Title:   "Test API",
					Version: "1.0.0",
				},
				Paths: &openapi3.Paths{},
				Components: &openapi3.Components{
					Schemas: make(map[string]*openapi3.SchemaRef),
				},
			}

			// Add schemas
			for i := 0; i < tt.schemaCount; i++ {
				schemaName := fmt.Sprintf("Schema%d", i)
				spec.Components.Schemas[schemaName] = &openapi3.SchemaRef{
					Value: &openapi3.Schema{
						Type: &openapi3.Types{"object"},
					},
				}
			}

			// Add operations
			for i := 0; i < tt.operationCount; i++ {
				pathName := fmt.Sprintf("/path%d", i)
				pathItem := &openapi3.PathItem{
					Get: &openapi3.Operation{
						OperationID: fmt.Sprintf("operation%d", i),
						Responses: &openapi3.Responses{},
					},
				}
				successResponse := &openapi3.ResponseRef{
					Value: &openapi3.Response{
						Description: &[]string{"Success"}[0],
					},
				}
				pathItem.Get.Responses.Set("200", successResponse)
				spec.Paths.Set(pathName, pathItem)
			}

			complexity := service.assessComplexity(spec)
			if complexity != tt.wantComplexity {
				t.Errorf("assessComplexity() = %v, want %v", complexity, tt.wantComplexity)
			}
		})
	}
}

func TestService_ErrorSuggestions(t *testing.T) {
	service := New()

	tests := []struct {
		name           string
		errorMsg       string
		wantContains   []string
		wantNotContains []string
	}{
		{
			name:     "parsing error with regex",
			errorMsg: "regex pattern error",
			wantContains: []string{
				"Remove or simplify regex patterns",
			},
		},
		{
			name:     "OpenAPI 2.0 error",
			errorMsg: "OpenAPI 2.0 not supported",
			wantContains: []string{
				"Convert your specification to OpenAPI 3.0+",
			},
		},
		{
			name:     "YAML syntax error",
			errorMsg: "YAML parsing failed",
			wantContains: []string{
				"Check YAML syntax and indentation",
			},
		},
		{
			name:     "JSON syntax error",
			errorMsg: "JSON parsing failed",
			wantContains: []string{
				"Check JSON syntax and structure",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := fmt.Errorf(tt.errorMsg)
			suggestions := service.getParsingErrorSuggestions(err)

			for _, want := range tt.wantContains {
				found := false
				for _, suggestion := range suggestions {
					if strings.Contains(suggestion, want) {
						found = true
						break
					}
				}
				if !found {
					t.Errorf("Expected suggestion containing '%s', but not found in %v", want, suggestions)
				}
			}

			for _, wantNot := range tt.wantNotContains {
				for _, suggestion := range suggestions {
					if strings.Contains(suggestion, wantNot) {
						t.Errorf("Did not expect suggestion containing '%s', but found in %v", wantNot, suggestions)
					}
				}
			}
		})
	}
}

func TestService_NetworkErrorSuggestions(t *testing.T) {
	service := New()

	tests := []struct {
		name         string
		errorMsg     string
		wantContains []string
	}{
		{
			name:     "HTTP error",
			errorMsg: "HTTP error 404",
			wantContains: []string{
				"Check if the URL is accessible",
			},
		},
		{
			name:     "timeout error",
			errorMsg: "timeout exceeded",
			wantContains: []string{
				"The server may be slow to respond",
			},
		},
		{
			name:     "connection error",
			errorMsg: "connection refused",
			wantContains: []string{
				"Check your internet connection",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := fmt.Errorf(tt.errorMsg)
			suggestions := service.getNetworkErrorSuggestions(err)

			for _, want := range tt.wantContains {
				found := false
				for _, suggestion := range suggestions {
					if strings.Contains(suggestion, want) {
						found = true
						break
					}
				}
				if !found {
					t.Errorf("Expected suggestion containing '%s', but not found in %v", want, suggestions)
				}
			}
		})
	}
}

func TestService_PerformanceRequirements(t *testing.T) {
	service := New()
	ctx := context.Background()

	// Create a reasonably complex spec
	specContent := `
openapi: 3.0.0
info:
  title: Performance Test API
  version: 1.0.0
paths:
  /users:
    get:
      operationId: getUsers
      responses:
        '200':
          description: Success
  /users/{id}:
    get:
      operationId: getUser
      responses:
        '200':
          description: Success
  /products:
    get:
      operationId: getProducts
      responses:
        '200':
          description: Success
components:
  schemas:
    User:
      type: object
      properties:
        id:
          type: integer
        name:
          type: string
    Product:
      type: object
      properties:
        id:
          type: integer
        name:
          type: string
`

	// Create temporary file
	tempDir := t.TempDir()
	filePath := filepath.Join(tempDir, "performance_test.yaml")
	
	err := os.WriteFile(filePath, []byte(specContent), 0644)
	if err != nil {
		t.Fatalf("Failed to write test file: %v", err)
	}

	// Test performance requirements (should be under 5 seconds)
	startTime := time.Now()
	result, err := service.ValidateFile(ctx, filePath)
	elapsedTime := time.Since(startTime)

	if err != nil {
		t.Fatalf("ValidateFile() error = %v", err)
	}

	if elapsedTime > 5*time.Second {
		t.Errorf("Validation took %v, should be under 5 seconds", elapsedTime)
	}

	if result.ValidationTime > 5*time.Second {
		t.Errorf("Reported validation time %v, should be under 5 seconds", result.ValidationTime)
	}

	t.Logf("Validation completed in %v", elapsedTime)
}

func BenchmarkService_ValidateFile(b *testing.B) {
	service := New()
	ctx := context.Background()

	// Create a test spec
	specContent := `
openapi: 3.0.0
info:
  title: Benchmark API
  version: 1.0.0
paths:
  /test:
    get:
      operationId: getTest
      responses:
        '200':
          description: Success
components:
  schemas:
    TestSchema:
      type: object
      properties:
        id:
          type: integer
`

	// Create temporary file
	tempDir := b.TempDir()
	filePath := filepath.Join(tempDir, "benchmark_spec.yaml")
	
	err := os.WriteFile(filePath, []byte(specContent), 0644)
	if err != nil {
		b.Fatalf("Failed to write test file: %v", err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := service.ValidateFile(ctx, filePath)
		if err != nil {
			b.Fatalf("ValidateFile() error = %v", err)
		}
	}
}