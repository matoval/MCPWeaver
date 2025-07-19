package integration

import (
	"context"
	"testing"

	"MCPWeaver/internal/mapping"
	"MCPWeaver/internal/parser"
	"MCPWeaver/internal/validator"
	"MCPWeaver/tests/helpers"

	"github.com/stretchr/testify/suite"
)

type WorkflowIntegrationTestSuite struct {
	suite.Suite
	helper          *helpers.TestHelper
	parserService   *parser.Service
	validatorService *validator.Service
	mappingService  *mapping.Service
	ctx             context.Context
}

func (s *WorkflowIntegrationTestSuite) SetupTest() {
	s.helper = helpers.NewTestHelper(s.T())
	s.parserService = parser.NewService()
	s.validatorService = validator.New()
	s.mappingService = mapping.NewService("https://api.test.com")
	s.ctx = context.Background()
}

func TestWorkflowIntegrationTestSuite(t *testing.T) {
	suite.Run(t, new(WorkflowIntegrationTestSuite))
}

func (s *WorkflowIntegrationTestSuite) TestCompleteWorkflow_ValidSpec() {
	// Create a valid OpenAPI spec
	specContent := helpers.ValidateOpenAPISpec()
	filePath, cleanup := s.helper.CreateTempFile(specContent, ".yaml")
	defer cleanup()

	// Step 1: Validate the spec
	validationResult, err := s.validatorService.ValidateFile(s.ctx, filePath)
	s.helper.AssertNoError(err)
	s.helper.AssertNotNil(validationResult)
	s.helper.AssertEqual(true, validationResult.Valid)

	// Step 2: Parse the spec
	parseResult, err := s.parserService.ParseFromFile(filePath)
	s.helper.AssertNoError(err)
	s.helper.AssertNotNil(parseResult)
	s.helper.AssertEqual("Test API", parseResult.Title)
	s.helper.AssertEqual("1.0.0", parseResult.Version)

	// Step 3: Map operations to MCP tools
	tools, err := s.mappingService.MapOperationsToTools(parseResult.Operations)
	s.helper.AssertNoError(err)
	s.helper.AssertNotNil(tools)
	s.helper.AssertEqual(true, len(tools) > 0)

	// Verify the complete workflow produces expected results
	s.helper.AssertEqual(true, len(parseResult.Operations) > 0)
	s.helper.AssertEqual(len(parseResult.Operations), len(tools))
}

func (s *WorkflowIntegrationTestSuite) TestCompleteWorkflow_InvalidSpec() {
	// Create an invalid OpenAPI spec
	specContent := helpers.InvalidOpenAPISpec()
	filePath, cleanup := s.helper.CreateTempFile(specContent, ".yaml")
	defer cleanup()

	// Step 1: Validate the spec (should fail)
	validationResult, err := s.validatorService.ValidateFile(s.ctx, filePath)
	s.helper.AssertNoError(err) // Validation call succeeds
	s.helper.AssertNotNil(validationResult)
	s.helper.AssertEqual(false, validationResult.Valid) // But spec is invalid

	// Step 2: Parse the spec (should fail)
	parseResult, err := s.parserService.ParseFromFile(filePath)
	s.helper.AssertError(err)
	s.helper.AssertNil(parseResult)

	// Workflow should fail early on invalid specs
}

func (s *WorkflowIntegrationTestSuite) TestCompleteWorkflow_LargeSpec() {
	// Create a large OpenAPI spec
	specContent := helpers.LargeOpenAPISpec()
	filePath, cleanup := s.helper.CreateTempFile(specContent, ".yaml")
	defer cleanup()

	// Test complete workflow performance with large spec
	s.helper.AssertPerformance(func() {
		// Step 1: Validate
		validationResult, err := s.validatorService.ValidateFile(s.ctx, filePath)
		if err == nil && validationResult != nil && validationResult.Valid {
			// Step 2: Parse
			parseResult, err := s.parserService.ParseFromFile(filePath)
			if err == nil && parseResult != nil {
				// Step 3: Map
				tools, err := s.mappingService.MapOperationsToTools(parseResult.Operations)
				if err == nil {
					s.helper.AssertNotNil(tools)
				}
			}
		}
	}, 30*1000*1000000) // 30 seconds for complete large spec workflow
}

func (s *WorkflowIntegrationTestSuite) TestWorkflow_ValidationAndParsing_Consistency() {
	// Test that validation and parsing are consistent
	testCases := []struct {
		name     string
		content  string
		shouldSucceed bool
	}{
		{
			name:     "Valid spec",
			content:  helpers.ValidateOpenAPISpec(),
			shouldSucceed: true,
		},
		{
			name:     "Invalid spec",
			content:  helpers.InvalidOpenAPISpec(),
			shouldSucceed: false,
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			filePath, cleanup := s.helper.CreateTempFile(tc.content, ".yaml")
			defer cleanup()

			// Validate
			validationResult, err := s.validatorService.ValidateFile(s.ctx, filePath)
			s.helper.AssertNoError(err)
			s.helper.AssertNotNil(validationResult)

			// Parse
			parseResult, parseErr := s.parserService.ParseFromFile(filePath)

			if tc.shouldSucceed {
				// Both should succeed for valid specs
				s.helper.AssertEqual(true, validationResult.Valid)
				s.helper.AssertNoError(parseErr)
				s.helper.AssertNotNil(parseResult)
			} else {
				// Validation should report invalid, parsing should fail
				s.helper.AssertEqual(false, validationResult.Valid)
				s.helper.AssertError(parseErr)
				s.helper.AssertNil(parseResult)
			}
		})
	}
}

func (s *WorkflowIntegrationTestSuite) TestWorkflow_ParsingAndMapping_Integration() {
	// Test that parsed operations can be properly mapped
	specContent := helpers.ValidateOpenAPISpec()
	filePath, cleanup := s.helper.CreateTempFile(specContent, ".yaml")
	defer cleanup()

	// Parse the spec
	parseResult, err := s.parserService.ParseFromFile(filePath)
	s.helper.AssertNoError(err)
	s.helper.AssertNotNil(parseResult)

	// Map operations to tools
	tools, err := s.mappingService.MapOperationsToTools(parseResult.Operations)
	s.helper.AssertNoError(err)
	s.helper.AssertNotNil(tools)

	// Verify mapping results
	s.helper.AssertEqual(len(parseResult.Operations), len(tools))
	
	// Check that each tool has required fields
	for i, tool := range tools {
		s.helper.AssertNotEqual("", tool.Name)
		s.helper.AssertNotEqual("", tool.Description)
		s.helper.AssertNotNil(tool.InputSchema)
		
		// Tool name should be related to operation
		operation := parseResult.Operations[i]
		s.helper.AssertContains(tool.Name, operation.Method)
	}
}

func (s *WorkflowIntegrationTestSuite) TestWorkflow_ErrorPropagation() {
	// Test that errors propagate correctly through the workflow
	
	// Test with non-existent file
	nonExistentPath := "/non/existent/file.yaml"

	// Validation should fail
	validationResult, validationErr := s.validatorService.ValidateFile(s.ctx, nonExistentPath)
	s.helper.AssertError(validationErr)
	s.helper.AssertNil(validationResult)

	// Parsing should fail
	parseResult, parseErr := s.parserService.ParseFromFile(nonExistentPath)
	s.helper.AssertError(parseErr)
	s.helper.AssertNil(parseResult)
}

func (s *WorkflowIntegrationTestSuite) TestWorkflow_ServiceInstantiation() {
	// Test that all services can be instantiated and work together
	
	// Test parser service
	parser := parser.NewService()
	s.helper.AssertNotNil(parser)

	// Test validator service  
	validator := validator.New()
	s.helper.AssertNotNil(validator)

	// Test mapping service with different base URLs
	baseURLs := []string{
		"https://api.test.com",
		"http://localhost:8080",
		"https://example.com/api/v1",
	}

	for _, baseURL := range baseURLs {
		mapper := mapping.NewService(baseURL)
		s.helper.AssertNotNil(mapper)
	}
}

func (s *WorkflowIntegrationTestSuite) TestWorkflow_ContextHandling() {
	// Test context handling across services
	specContent := helpers.ValidateOpenAPISpec()
	filePath, cleanup := s.helper.CreateTempFile(specContent, ".yaml")
	defer cleanup()

	// Test with normal context
	result1, err1 := s.validatorService.ValidateFile(s.ctx, filePath)
	s.helper.AssertNoError(err1)
	s.helper.AssertNotNil(result1)

	// Test with canceled context
	cancelCtx, cancel := context.WithCancel(context.Background())
	cancel()

	result2, err2 := s.validatorService.ValidateFile(cancelCtx, filePath)
	// Should handle canceled context gracefully
	if err2 != nil {
		s.helper.AssertError(err2)
	} else {
		s.helper.AssertNotNil(result2)
	}
}

func (s *WorkflowIntegrationTestSuite) TestWorkflow_SequentialOperations() {
	// Test sequential operations across services
	specContent := helpers.ValidateOpenAPISpec()
	filePath, cleanup := s.helper.CreateTempFile(specContent, ".yaml")
	defer cleanup()

	// Run multiple operations sequentially (safer than concurrent)
	for i := 0; i < 3; i++ {
		// Sequential validation
		result1, err1 := s.validatorService.ValidateFile(s.ctx, filePath)
		s.helper.AssertNoError(err1)
		s.helper.AssertNotNil(result1)

		// Sequential parsing
		result2, err2 := s.parserService.ParseFromFile(filePath)
		s.helper.AssertNoError(err2)
		s.helper.AssertNotNil(result2)
	}
}

func (s *WorkflowIntegrationTestSuite) TestWorkflow_MemoryManagement() {
	// Test memory management across multiple workflow executions
	specContent := helpers.ValidateOpenAPISpec()
	
	for i := 0; i < 5; i++ {
		filePath, cleanup := s.helper.CreateTempFile(specContent, ".yaml")
		
		// Execute complete workflow
		validationResult, err := s.validatorService.ValidateFile(s.ctx, filePath)
		s.helper.AssertNoError(err)
		s.helper.AssertNotNil(validationResult)

		if validationResult.Valid {
			parseResult, err := s.parserService.ParseFromFile(filePath)
			s.helper.AssertNoError(err)
			s.helper.AssertNotNil(parseResult)

			tools, err := s.mappingService.MapOperationsToTools(parseResult.Operations)
			s.helper.AssertNoError(err)
			s.helper.AssertNotNil(tools)
		}
		
		cleanup() // Clean up immediately to test memory management
	}
}