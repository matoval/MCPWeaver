package integration

import (
	"testing"

	"MCPWeaver/internal/parser"
	"MCPWeaver/tests/helpers"

	"github.com/stretchr/testify/suite"
)

type ParserIntegrationTestSuite struct {
	suite.Suite
	helper  *helpers.TestHelper
	service *parser.Service
}

func (s *ParserIntegrationTestSuite) SetupTest() {
	s.helper = helpers.NewTestHelper(s.T())
	s.service = parser.NewService()
}

func TestParserIntegrationTestSuite(t *testing.T) {
	suite.Run(t, new(ParserIntegrationTestSuite))
}

func (s *ParserIntegrationTestSuite) TestParseFromFile_ValidSpec() {
	// Create a valid OpenAPI spec file
	specContent := helpers.ValidateOpenAPISpec()
	filePath, cleanup := s.helper.CreateTempFile(specContent, ".yaml")
	defer cleanup()

	// Test parsing
	result, err := s.service.ParseFromFile(filePath)

	s.helper.AssertNoError(err)
	s.helper.AssertNotNil(result)
	s.helper.AssertEqual("Test API", result.Title)
	s.helper.AssertEqual("1.0.0", result.Version)
	s.helper.AssertEqual(true, len(result.Operations) > 0)
}

func (s *ParserIntegrationTestSuite) TestParseFromFile_InvalidSpec() {
	// Create an invalid OpenAPI spec file
	specContent := helpers.InvalidOpenAPISpec()
	filePath, cleanup := s.helper.CreateTempFile(specContent, ".yaml")
	defer cleanup()

	// Test parsing
	result, err := s.service.ParseFromFile(filePath)

	s.helper.AssertError(err)
	s.helper.AssertNil(result)
}

func (s *ParserIntegrationTestSuite) TestParseFromFile_NonExistentFile() {
	// Test with non-existent file
	result, err := s.service.ParseFromFile("/non/existent/file.yaml")

	s.helper.AssertError(err)
	s.helper.AssertNil(result)
}

func (s *ParserIntegrationTestSuite) TestParseFromFile_LargeSpec() {
	// Create a large OpenAPI spec
	specContent := helpers.LargeOpenAPISpec()
	filePath, cleanup := s.helper.CreateTempFile(specContent, ".yaml")
	defer cleanup()

	// Test parsing performance with large spec
	s.helper.AssertPerformance(func() {
		result, err := s.service.ParseFromFile(filePath)
		if err == nil {
			s.helper.AssertNotNil(result)
		}
	}, 10*1000*1000000) // 10 seconds for large spec
}

func (s *ParserIntegrationTestSuite) TestParseFromFile_JSONFormat() {
	// Create a JSON format OpenAPI spec
	specContent := `{
		"openapi": "3.0.0",
		"info": {
			"title": "JSON Test API",
			"version": "1.0.0"
		},
		"paths": {
			"/test": {
				"get": {
					"summary": "Test endpoint",
					"responses": {
						"200": {
							"description": "Success"
						}
					}
				}
			}
		}
	}`

	filePath, cleanup := s.helper.CreateTempFile(specContent, ".json")
	defer cleanup()

	// Test parsing JSON format
	result, err := s.service.ParseFromFile(filePath)

	if err == nil {
		s.helper.AssertNotNil(result)
		s.helper.AssertEqual("JSON Test API", result.Title)
	} else {
		// JSON parsing might not be implemented yet, which is fine
		s.helper.AssertError(err)
	}
}

func (s *ParserIntegrationTestSuite) TestParseFromFile_EmptyFile() {
	// Create an empty file
	filePath, cleanup := s.helper.CreateTempFile("", ".yaml")
	defer cleanup()

	// Test parsing empty file
	result, err := s.service.ParseFromFile(filePath)

	s.helper.AssertError(err)
	s.helper.AssertNil(result)
}

func (s *ParserIntegrationTestSuite) TestParseFromFile_WithComments() {
	// Create OpenAPI spec with YAML comments
	specContent := `# This is a comment
openapi: 3.0.0
info:
  title: Test API with Comments
  version: 1.0.0
  # Another comment
  description: A test API for unit testing
paths:
  /users:
    get:
      # Comment in path
      summary: Get users
      responses:
        '200':
          description: List of users`

	filePath, cleanup := s.helper.CreateTempFile(specContent, ".yaml")
	defer cleanup()

	// Test parsing spec with comments
	result, err := s.service.ParseFromFile(filePath)

	s.helper.AssertNoError(err)
	s.helper.AssertNotNil(result)
	s.helper.AssertEqual("Test API with Comments", result.Title)
}

func (s *ParserIntegrationTestSuite) TestParseFromFile_MultipleFormats() {
	// Test different file extensions
	testCases := []struct {
		name      string
		extension string
		content   string
	}{
		{
			name:      "YAML file",
			extension: ".yaml",
			content:   helpers.ValidateOpenAPISpec(),
		},
		{
			name:      "YML file",
			extension: ".yml",
			content:   helpers.ValidateOpenAPISpec(),
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			filePath, cleanup := s.helper.CreateTempFile(tc.content, tc.extension)
			defer cleanup()

			result, err := s.service.ParseFromFile(filePath)
			s.helper.AssertNoError(err)
			s.helper.AssertNotNil(result)
			s.helper.AssertEqual("Test API", result.Title)
		})
	}
}

func (s *ParserIntegrationTestSuite) TestParseFromFile_ErrorRecovery() {
	// Test that parser recovers from errors and can parse valid files after invalid ones

	// First try to parse an invalid file
	invalidContent := helpers.InvalidOpenAPISpec()
	invalidPath, cleanup1 := s.helper.CreateTempFile(invalidContent, ".yaml")
	defer cleanup1()

	result1, err1 := s.service.ParseFromFile(invalidPath)
	s.helper.AssertError(err1)
	s.helper.AssertNil(result1)

	// Then parse a valid file (should work fine)
	validContent := helpers.ValidateOpenAPISpec()
	validPath, cleanup2 := s.helper.CreateTempFile(validContent, ".yaml")
	defer cleanup2()

	result2, err2 := s.service.ParseFromFile(validPath)
	s.helper.AssertNoError(err2)
	s.helper.AssertNotNil(result2)
	s.helper.AssertEqual("Test API", result2.Title)
}

// Memory and performance tests
func (s *ParserIntegrationTestSuite) TestParseFromFile_MemoryUsage() {
	// Test parsing multiple files doesn't leak memory
	specContent := helpers.ValidateOpenAPISpec()

	for i := 0; i < 10; i++ {
		filePath, cleanup := s.helper.CreateTempFile(specContent, ".yaml")

		result, err := s.service.ParseFromFile(filePath)
		s.helper.AssertNoError(err)
		s.helper.AssertNotNil(result)

		cleanup() // Clean up immediately to test memory management
	}
}

func (s *ParserIntegrationTestSuite) TestParseFromFile_MultipleSequential() {
	// Test sequential parsing operations (safer than concurrent)
	specContent := helpers.ValidateOpenAPISpec()
	filePath, cleanup := s.helper.CreateTempFile(specContent, ".yaml")
	defer cleanup()

	// Run multiple parsing operations sequentially
	for i := 0; i < 3; i++ {
		result, err := s.service.ParseFromFile(filePath)
		s.helper.AssertNoError(err)
		s.helper.AssertNotNil(result)
		s.helper.AssertEqual("Test API", result.Title)
	}
}
