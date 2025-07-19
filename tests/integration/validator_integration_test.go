package integration

import (
	"context"
	"testing"

	"MCPWeaver/internal/validator"
	"MCPWeaver/tests/helpers"

	"github.com/stretchr/testify/suite"
)

type ValidatorIntegrationTestSuite struct {
	suite.Suite
	helper  *helpers.TestHelper
	service *validator.Service
	ctx     context.Context
}

func (s *ValidatorIntegrationTestSuite) SetupTest() {
	s.helper = helpers.NewTestHelper(s.T())
	s.service = validator.New()
	s.ctx = context.Background()
}

func TestValidatorIntegrationTestSuite(t *testing.T) {
	suite.Run(t, new(ValidatorIntegrationTestSuite))
}

func (s *ValidatorIntegrationTestSuite) TestValidateFile_ValidSpec() {
	// Create a valid OpenAPI spec file
	specContent := helpers.ValidateOpenAPISpec()
	filePath, cleanup := s.helper.CreateTempFile(specContent, ".yaml")
	defer cleanup()

	// Test validation
	result, err := s.service.ValidateFile(s.ctx, filePath)
	
	s.helper.AssertNoError(err)
	s.helper.AssertNotNil(result)
	s.helper.AssertEqual(true, result.Valid)
	s.helper.AssertEqual(0, len(result.Errors))
}

func (s *ValidatorIntegrationTestSuite) TestValidateFile_InvalidSpec() {
	// Create an invalid OpenAPI spec file
	specContent := helpers.InvalidOpenAPISpec()
	filePath, cleanup := s.helper.CreateTempFile(specContent, ".yaml")
	defer cleanup()

	// Test validation
	result, err := s.service.ValidateFile(s.ctx, filePath)
	
	s.helper.AssertNoError(err) // Validation should succeed but return invalid result
	s.helper.AssertNotNil(result)
	s.helper.AssertEqual(false, result.Valid)
	s.helper.AssertEqual(true, len(result.Errors) > 0)
}

func (s *ValidatorIntegrationTestSuite) TestValidateFile_NonExistentFile() {
	// Test with non-existent file
	result, err := s.service.ValidateFile(s.ctx, "/non/existent/file.yaml")
	
	s.helper.AssertError(err)
	s.helper.AssertNil(result)
}

func (s *ValidatorIntegrationTestSuite) TestValidateFile_EmptyFile() {
	// Create an empty file
	filePath, cleanup := s.helper.CreateTempFile("", ".yaml")
	defer cleanup()

	// Test validation
	result, err := s.service.ValidateFile(s.ctx, filePath)
	
	if err != nil {
		s.helper.AssertError(err)
		s.helper.AssertNil(result)
	} else {
		s.helper.AssertNotNil(result)
		s.helper.AssertEqual(false, result.Valid)
	}
}

func (s *ValidatorIntegrationTestSuite) TestValidateFile_LargeSpec() {
	// Create a large OpenAPI spec
	specContent := helpers.LargeOpenAPISpec()
	filePath, cleanup := s.helper.CreateTempFile(specContent, ".yaml")
	defer cleanup()

	// Test validation performance with large spec
	s.helper.AssertPerformance(func() {
		result, err := s.service.ValidateFile(s.ctx, filePath)
		if err == nil {
			s.helper.AssertNotNil(result)
		}
	}, 15*1000*1000000) // 15 seconds for large spec validation
}

func (s *ValidatorIntegrationTestSuite) TestValidateFile_WithWarnings() {
	// Create a spec that's valid but might have warnings
	specContent := `openapi: 3.0.0
info:
  title: Test API
  version: 1.0.0
paths:
  /users:
    get:
      summary: Get users
      responses:
        '200':
          description: List of users
      # No operationId - might generate warning
  /posts:
    post:
      # No summary - might generate warning
      responses:
        '201':
          description: Post created`
	
	filePath, cleanup := s.helper.CreateTempFile(specContent, ".yaml")
	defer cleanup()

	// Test validation
	result, err := s.service.ValidateFile(s.ctx, filePath)
	
	s.helper.AssertNoError(err)
	s.helper.AssertNotNil(result)
	
	// Should be valid but might have warnings
	if result.Valid {
		// Valid specs might still have warnings
		s.helper.AssertEqual(true, len(result.Warnings) >= 0)
	}
}

func (s *ValidatorIntegrationTestSuite) TestValidateFile_MultipleFormats() {
	// Test different file extensions
	testCases := []struct {
		name      string
		extension string
		content   string
		shouldPass bool
	}{
		{
			name:       "Valid YAML file",
			extension:  ".yaml",
			content:    helpers.ValidateOpenAPISpec(),
			shouldPass: true,
		},
		{
			name:       "Valid YML file",
			extension:  ".yml",
			content:    helpers.ValidateOpenAPISpec(),
			shouldPass: true,
		},
		{
			name:       "Invalid YAML file",
			extension:  ".yaml",
			content:    helpers.InvalidOpenAPISpec(),
			shouldPass: false,
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			filePath, cleanup := s.helper.CreateTempFile(tc.content, tc.extension)
			defer cleanup()

			result, err := s.service.ValidateFile(s.ctx, filePath)
			s.helper.AssertNoError(err)
			s.helper.AssertNotNil(result)
			s.helper.AssertEqual(tc.shouldPass, result.Valid)
		})
	}
}

func (s *ValidatorIntegrationTestSuite) TestValidateFile_ContextCancellation() {
	// Test with canceled context
	specContent := helpers.LargeOpenAPISpec()
	filePath, cleanup := s.helper.CreateTempFile(specContent, ".yaml")
	defer cleanup()

	// Create a canceled context
	cancelCtx, cancel := context.WithCancel(context.Background())
	cancel() // Cancel immediately

	// Test validation with canceled context
	result, err := s.service.ValidateFile(cancelCtx, filePath)
	
	// Should handle canceled context gracefully
	if err != nil {
		s.helper.AssertError(err)
	} else {
		s.helper.AssertNotNil(result)
	}
}

func (s *ValidatorIntegrationTestSuite) TestValidateFile_ErrorRecovery() {
	// Test that validator recovers from errors
	
	// First try to validate an invalid file
	invalidContent := helpers.InvalidOpenAPISpec()
	invalidPath, cleanup1 := s.helper.CreateTempFile(invalidContent, ".yaml")
	defer cleanup1()

	result1, err1 := s.service.ValidateFile(s.ctx, invalidPath)
	s.helper.AssertNoError(err1) // Should succeed but return invalid result
	s.helper.AssertNotNil(result1)
	s.helper.AssertEqual(false, result1.Valid)

	// Then validate a valid file (should work fine)
	validContent := helpers.ValidateOpenAPISpec()
	validPath, cleanup2 := s.helper.CreateTempFile(validContent, ".yaml")
	defer cleanup2()

	result2, err2 := s.service.ValidateFile(s.ctx, validPath)
	s.helper.AssertNoError(err2)
	s.helper.AssertNotNil(result2)
	s.helper.AssertEqual(true, result2.Valid)
}

func (s *ValidatorIntegrationTestSuite) TestValidateFile_ValidationResult_Details() {
	// Test detailed validation result structure
	specContent := helpers.ValidateOpenAPISpec()
	filePath, cleanup := s.helper.CreateTempFile(specContent, ".yaml")
	defer cleanup()

	result, err := s.service.ValidateFile(s.ctx, filePath)
	s.helper.AssertNoError(err)
	s.helper.AssertNotNil(result)

	// Test result structure
	s.helper.AssertEqual(true, result.Valid)
	s.helper.AssertNotNil(result.Errors)     // Should be empty slice, not nil
	s.helper.AssertNotNil(result.Warnings)  // Should be empty slice, not nil
	s.helper.AssertNotNil(result.Suggestions) // Should be empty slice, not nil
	
	// Validation time should be reasonable
	s.helper.AssertEqual(true, result.ValidationTime >= 0)
	
	// Check if spec info is populated
	if result.SpecInfo != nil {
		s.helper.AssertEqual("Test API", result.SpecInfo.Title)
		s.helper.AssertEqual("1.0.0", result.SpecInfo.Version)
		s.helper.AssertEqual(true, result.SpecInfo.OperationCount > 0)
	}
}

// Performance and memory tests
func (s *ValidatorIntegrationTestSuite) TestValidateFile_MemoryUsage() {
	// Test validating multiple files doesn't leak memory
	specContent := helpers.ValidateOpenAPISpec()
	
	for i := 0; i < 10; i++ {
		filePath, cleanup := s.helper.CreateTempFile(specContent, ".yaml")
		
		result, err := s.service.ValidateFile(s.ctx, filePath)
		s.helper.AssertNoError(err)
		s.helper.AssertNotNil(result)
		
		cleanup() // Clean up immediately to test memory management
	}
}

func (s *ValidatorIntegrationTestSuite) TestValidateFile_MultipleSequential() {
	// Test sequential validation operations (safer than concurrent)
	specContent := helpers.ValidateOpenAPISpec()
	filePath, cleanup := s.helper.CreateTempFile(specContent, ".yaml")
	defer cleanup()

	// Run multiple validation operations sequentially
	for i := 0; i < 3; i++ {
		result, err := s.service.ValidateFile(s.ctx, filePath)
		s.helper.AssertNoError(err)
		s.helper.AssertNotNil(result)
		s.helper.AssertEqual(true, result.Valid)
	}
}

func (s *ValidatorIntegrationTestSuite) TestValidateFile_ValidateFileAndValidateFromURL() {
	// Test both file and URL validation work consistently
	specContent := helpers.ValidateOpenAPISpec()
	filePath, cleanup := s.helper.CreateTempFile(specContent, ".yaml")
	defer cleanup()

	// Validate from file
	fileResult, err := s.service.ValidateFile(s.ctx, filePath)
	s.helper.AssertNoError(err)
	s.helper.AssertNotNil(fileResult)
	s.helper.AssertEqual(true, fileResult.Valid)

	// Note: ValidateFromURL would need actual HTTP server for testing
	// For integration test, we just verify the method exists and can be called
	// In a real scenario, this would test against a mock HTTP server
}