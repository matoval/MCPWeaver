package app

import (
	"context"
	"os"
	"path/filepath"
	"testing"
	"time"

	"MCPWeaver/internal/database"
	"MCPWeaver/internal/validator"
)

func TestApp_ValidateSpec(t *testing.T) {
	// Create temporary directory for database
	tempDir := t.TempDir()

	// Create app instance
	app := NewApp()

	// Initialize database
	dbPath := filepath.Join(tempDir, "test.db")
	db, err := database.Open(dbPath)
	if err != nil {
		t.Fatalf("Failed to open database: %v", err)
	}
	defer db.Close()

	app.db = db.GetConn()
	app.validationCacheRepo = database.NewValidationCacheRepository(db)
	app.validatorService = validator.New()
	app.ctx = context.Background()

	// Create test OpenAPI spec file
	specContent := `
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
`

	specPath := filepath.Join(tempDir, "test_spec.yaml")
	err = os.WriteFile(specPath, []byte(specContent), 0644)
	if err != nil {
		t.Fatalf("Failed to write spec file: %v", err)
	}

	// Test validation
	result, err := app.ValidateSpec(specPath)
	if err != nil {
		t.Fatalf("ValidateSpec() error = %v", err)
	}

	if !result.Valid {
		t.Errorf("ValidateSpec() should return valid for good spec")
	}

	if result.CacheHit {
		t.Errorf("ValidateSpec() should not be cache hit on first call")
	}

	if result.SpecInfo == nil {
		t.Error("ValidateSpec() should return spec info")
	}

	if result.SpecInfo.Title != "Test API" {
		t.Errorf("ValidateSpec() spec title = %s, want 'Test API'", result.SpecInfo.Title)
	}

	// Test cache hit
	result2, err := app.ValidateSpec(specPath)
	if err != nil {
		t.Fatalf("ValidateSpec() error = %v", err)
	}

	if !result2.CacheHit {
		t.Errorf("ValidateSpec() should be cache hit on second call")
	}
}

func TestApp_ValidateURL(t *testing.T) {
	// Create temporary directory for database
	tempDir := t.TempDir()

	// Create app instance
	app := NewApp()

	// Initialize database
	dbPath := filepath.Join(tempDir, "test.db")
	db, err := database.Open(dbPath)
	if err != nil {
		t.Fatalf("Failed to open database: %v", err)
	}
	defer db.Close()

	app.db = db.GetConn()
	app.validationCacheRepo = database.NewValidationCacheRepository(db)
	app.validatorService = validator.New()
	app.ctx = context.Background()

	// Test invalid URL
	result, err := app.ValidateURL("invalid-url")
	if err != nil {
		t.Fatalf("ValidateURL() error = %v", err)
	}

	if result.Valid {
		t.Error("ValidateURL() should return invalid for invalid URL")
	}

	if len(result.Errors) == 0 {
		t.Error("ValidateURL() should return errors for invalid URL")
	}
}

func TestApp_ValidateSpec_FileNotFound(t *testing.T) {
	// Create temporary directory for database
	tempDir := t.TempDir()

	// Create app instance
	app := NewApp()

	// Initialize database
	dbPath := filepath.Join(tempDir, "test.db")
	db, err := database.Open(dbPath)
	if err != nil {
		t.Fatalf("Failed to open database: %v", err)
	}
	defer db.Close()

	app.db = db.GetConn()
	app.validationCacheRepo = database.NewValidationCacheRepository(db)
	app.validatorService = validator.New()
	app.ctx = context.Background()

	// Test non-existent file
	result, err := app.ValidateSpec("/non/existent/file.yaml")
	if err != nil {
		t.Fatalf("ValidateSpec() error = %v", err)
	}

	if result.Valid {
		t.Error("ValidateSpec() should return invalid for non-existent file")
	}

	if len(result.Errors) == 0 {
		t.Error("ValidateSpec() should return errors for non-existent file")
	}
}

func TestApp_ExportValidationResult(t *testing.T) {
	app := NewApp()

	// Create test validation result
	result := &ValidationResult{
		Valid:          true,
		Errors:         []ValidationError{},
		Warnings:       []ValidationWarning{},
		Suggestions:    []string{"Consider adding more documentation"},
		ValidationTime: 100 * time.Millisecond,
		ValidatedAt:    time.Now(),
	}

	// Test export
	exported, err := app.ExportValidationResult(result)
	if err != nil {
		t.Fatalf("ExportValidationResult() error = %v", err)
	}

	if exported == "" {
		t.Error("ExportValidationResult() should return non-empty string")
	}

	// Should contain JSON structure
	if !contains(exported, "\"version\"") {
		t.Error("ExportValidationResult() should contain version field")
	}

	if !contains(exported, "\"validationResult\"") {
		t.Error("ExportValidationResult() should contain validationResult field")
	}

	if !contains(exported, "\"exportedAt\"") {
		t.Error("ExportValidationResult() should contain exportedAt field")
	}
}

func TestApp_ExportValidationResult_Nil(t *testing.T) {
	app := NewApp()

	// Test export with nil result
	_, err := app.ExportValidationResult(nil)
	if err == nil {
		t.Error("ExportValidationResult() should return error for nil result")
	}
}

func TestApp_GetValidationCacheStats(t *testing.T) {
	// Create temporary directory for database
	tempDir := t.TempDir()

	// Create app instance
	app := NewApp()

	// Initialize database
	dbPath := filepath.Join(tempDir, "test.db")
	db, err := database.Open(dbPath)
	if err != nil {
		t.Fatalf("Failed to open database: %v", err)
	}
	defer db.Close()

	app.db = db.GetConn()
	app.validationCacheRepo = database.NewValidationCacheRepository(db)
	app.ctx = context.Background()

	// Test getting stats
	stats, err := app.GetValidationCacheStats()
	if err != nil {
		t.Fatalf("GetValidationCacheStats() error = %v", err)
	}

	if stats == nil {
		t.Error("GetValidationCacheStats() should return stats")
	}

	if stats.TotalEntries < 0 {
		t.Error("GetValidationCacheStats() should return non-negative total entries")
	}
}

func TestApp_ClearValidationCache(t *testing.T) {
	// Create temporary directory for database
	tempDir := t.TempDir()

	// Create app instance
	app := NewApp()

	// Initialize database
	dbPath := filepath.Join(tempDir, "test.db")
	db, err := database.Open(dbPath)
	if err != nil {
		t.Fatalf("Failed to open database: %v", err)
	}
	defer db.Close()

	app.db = db.GetConn()
	app.validationCacheRepo = database.NewValidationCacheRepository(db)
	app.ctx = context.Background()

	// Test clearing cache
	err = app.ClearValidationCache()
	if err != nil {
		t.Fatalf("ClearValidationCache() error = %v", err)
	}
}

func TestApp_ValidateSpec_EmptyPath(t *testing.T) {
	app := NewApp()

	// Test with empty path
	_, err := app.ValidateSpec("")
	if err == nil {
		t.Error("ValidateSpec() should return error for empty path")
	}
}

func TestApp_ValidateURL_EmptyURL(t *testing.T) {
	app := NewApp()

	// Test with empty URL
	_, err := app.ValidateURL("")
	if err == nil {
		t.Error("ValidateURL() should return error for empty URL")
	}
}

func TestApp_ValidationIntegration(t *testing.T) {
	// Create temporary directory for database
	tempDir := t.TempDir()

	// Create app instance
	app := NewApp()

	// Initialize database
	dbPath := filepath.Join(tempDir, "test.db")
	db, err := database.Open(dbPath)
	if err != nil {
		t.Fatalf("Failed to open database: %v", err)
	}
	defer db.Close()

	app.db = db.GetConn()
	app.validationCacheRepo = database.NewValidationCacheRepository(db)
	app.validatorService = validator.New()
	app.ctx = context.Background()

	// Create test OpenAPI spec file with issues
	specContent := `
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
  /test2:
    get:
      operationId: getTest
      responses:
        '200':
          description: Success
`

	specPath := filepath.Join(tempDir, "test_spec_with_issues.yaml")
	err = os.WriteFile(specPath, []byte(specContent), 0644)
	if err != nil {
		t.Fatalf("Failed to write spec file: %v", err)
	}

	// Test validation
	result, err := app.ValidateSpec(specPath)
	if err != nil {
		t.Fatalf("ValidateSpec() error = %v", err)
	}

	if result.Valid {
		t.Error("ValidateSpec() should return invalid for spec with duplicate operation IDs")
	}

	if len(result.Errors) == 0 {
		t.Error("ValidateSpec() should return errors for spec with duplicate operation IDs")
	}

	if len(result.Warnings) == 0 {
		t.Error("ValidateSpec() should return warnings for spec with missing descriptions")
	}

	// Test export
	exported, err := app.ExportValidationResult(result)
	if err != nil {
		t.Fatalf("ExportValidationResult() error = %v", err)
	}

	if exported == "" {
		t.Error("ExportValidationResult() should return non-empty string")
	}

	// Test cache stats
	stats, err := app.GetValidationCacheStats()
	if err != nil {
		t.Fatalf("GetValidationCacheStats() error = %v", err)
	}

	if stats.TotalEntries == 0 {
		t.Error("GetValidationCacheStats() should show cached entry")
	}
}

// Helper function to check if string contains substring
func contains(s, substr string) bool {
	return len(s) >= len(substr) &&
		(s == substr ||
			s[:len(substr)] == substr ||
			s[len(s)-len(substr):] == substr ||
			containsAt(s, substr, 1))
}

func containsAt(s, substr string, start int) bool {
	if start >= len(s) {
		return false
	}
	if start+len(substr) > len(s) {
		return containsAt(s, substr, start+1)
	}
	if s[start:start+len(substr)] == substr {
		return true
	}
	return containsAt(s, substr, start+1)
}
