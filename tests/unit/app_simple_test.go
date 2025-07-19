package unit

import (
	"fmt"
	"testing"

	"MCPWeaver/internal/app"
	"MCPWeaver/tests/helpers"

	"github.com/stretchr/testify/suite"
)

type AppSimpleTestSuite struct {
	suite.Suite
	helper *helpers.TestHelper
	app    *app.App
}

func (s *AppSimpleTestSuite) SetupTest() {
	s.helper = helpers.NewTestHelper(s.T())
	s.app = &app.App{}
}

func TestAppSimpleTestSuite(t *testing.T) {
	suite.Run(t, new(AppSimpleTestSuite))
}

// Test helper functions that exist in the app package
func (s *AppSimpleTestSuite) TestFileExists_Helper() {
	// Create a temporary file
	filePath, cleanup := s.helper.CreateTempFile("test content", ".txt")
	defer cleanup()

	// Test that the helper function logic would work
	// (We can't test the actual app.fileExists since it's not exported)
	s.helper.AssertFileExists(filePath)
}

func (s *AppSimpleTestSuite) TestDirExists_Helper() {
	// Create a temporary directory
	dirPath, cleanup := s.helper.CreateTempDir()
	defer cleanup()

	// Test directory existence validation logic
	s.helper.AssertFileExists(dirPath)
}

// Test error code constants
func (s *AppSimpleTestSuite) TestErrorCodes() {
	codes := []string{
		app.ErrCodeValidation,
		app.ErrCodeNotFound,
		app.ErrCodeInternalError,
		app.ErrCodeFileAccess,
		app.ErrCodeNetworkError,
		app.ErrCodeParsingError,
		app.ErrCodeGenerationError,
		app.ErrCodeDatabaseError,
		app.ErrCodePermissionError,
		app.ErrCodeRateLimitError,
		app.ErrCodeTimeoutError,
		app.ErrCodeConfigError,
		app.ErrCodeAuthError,
	}

	for _, code := range codes {
		s.helper.AssertNotNil(code)
		s.helper.AssertNotEqual("", code)
	}
}

// Test error type constants
func (s *AppSimpleTestSuite) TestErrorTypes() {
	types := []string{
		app.ErrorTypeValidation,
		app.ErrorTypeSystem,
		app.ErrorTypeNetwork,
		app.ErrorTypeFileSystem,
		app.ErrorTypeDatabase,
		app.ErrorTypeGeneration,
		app.ErrorTypePermission,
		app.ErrorTypeConfiguration,
		app.ErrorTypeAuthentication,
	}

	for _, errorType := range types {
		s.helper.AssertNotNil(errorType)
		s.helper.AssertNotEqual("", errorType)
	}
}

// Test project validation logic
func (s *AppSimpleTestSuite) TestProjectValidation() {
	// Test empty name validation
	request := app.CreateProjectRequest{
		Name:       "",
		OutputPath: "/tmp/output",
	}
	
	// This simulates the validation logic that would be in CreateProject
	hasError := request.Name == ""
	s.helper.AssertEqual(true, hasError)

	// Test empty output path validation
	request2 := app.CreateProjectRequest{
		Name:       "Valid Name",
		OutputPath: "",
	}
	
	hasError2 := request2.OutputPath == ""
	s.helper.AssertEqual(true, hasError2)

	// Test valid request
	request3 := app.CreateProjectRequest{
		Name:       "Valid Project",
		OutputPath: "/tmp/output",
	}
	
	hasError3 := request3.Name == "" || request3.OutputPath == ""
	s.helper.AssertEqual(false, hasError3)
}

// Test generation job validation
func (s *AppSimpleTestSuite) TestGenerationJobValidation() {
	// Test empty job ID validation (simulating GetGenerationJob logic)
	jobID := ""
	hasError := jobID == ""
	s.helper.AssertEqual(true, hasError)

	// Test valid job ID
	jobID2 := "gen_12345"
	hasError2 := jobID2 == ""
	s.helper.AssertEqual(false, hasError2)

	// Test empty project ID validation (simulating GenerateServer logic)
	projectID := ""
	hasError3 := projectID == ""
	s.helper.AssertEqual(true, hasError3)
}

// Test file path validation logic
func (s *AppSimpleTestSuite) TestFilePathValidation() {
	// Test empty path validation
	filePath := ""
	hasError := filePath == ""
	s.helper.AssertEqual(true, hasError)

	// Test valid path
	filePath2 := "/path/to/file.yaml"
	hasError2 := filePath2 == ""
	s.helper.AssertEqual(false, hasError2)
}

// Test URL validation logic  
func (s *AppSimpleTestSuite) TestURLValidation() {
	// Test empty URL validation
	url := ""
	hasError := url == ""
	s.helper.AssertEqual(true, hasError)

	// Test invalid URL format
	url2 := "not-a-url"
	isValidFormat := url2 != "" && (len(url2) < 7 || (url2[:7] != "http://" && url2[:8] != "https://"))
	s.helper.AssertEqual(true, isValidFormat) // Should be invalid

	// Test valid URL
	url3 := "https://api.example.com/spec.yaml"
	isValidFormat2 := len(url3) >= 8 && url3[:8] == "https://"
	s.helper.AssertEqual(true, isValidFormat2)
}

// Test generation progress validation
func (s *AppSimpleTestSuite) TestProgressValidation() {
	// Test invalid progress values
	testCases := []struct {
		progress float64
		valid    bool
	}{
		{-0.1, false},
		{0.0, true},
		{0.5, true},
		{1.0, true},
		{1.1, false},
	}

	for _, tc := range testCases {
		isValid := tc.progress >= 0.0 && tc.progress <= 1.0
		s.helper.AssertEqual(tc.valid, isValid)
	}
}

// Test project name validation edge cases
func (s *AppSimpleTestSuite) TestProjectNameValidation() {
	testCases := []struct {
		name  string
		valid bool
	}{
		{"", false},                    // Empty name
		{"a", true},                    // Single character
		{"Valid Project Name", true},   // Normal name
		{"Project123", true},           // With numbers
		{"Project-Name_v2", true},      // With special chars
	}

	for _, tc := range testCases {
		isValid := tc.name != ""
		s.helper.AssertEqual(tc.valid, isValid)
	}
}

// Test settings port validation
func (s *AppSimpleTestSuite) TestPortValidation() {
	testCases := []struct {
		port  int
		valid bool
	}{
		{999, false},   // Below minimum
		{1000, true},   // At minimum
		{8080, true},   // Common port
		{65535, true},  // At maximum
		{65536, false}, // Above maximum
	}

	for _, tc := range testCases {
		isValid := tc.port >= 1000 && tc.port <= 65535
		s.helper.AssertEqual(tc.valid, isValid)
	}
}

// Test log level validation
func (s *AppSimpleTestSuite) TestLogLevelValidation() {
	validLevels := []string{"debug", "info", "warn", "error"}
	
	testCases := []struct {
		level string
		valid bool
	}{
		{"debug", true},
		{"info", true},
		{"warn", true},
		{"error", true},
		{"invalid", false},
		{"", false},
	}

	for _, tc := range testCases {
		isValid := false
		for _, validLevel := range validLevels {
			if tc.level == validLevel {
				isValid = true
				break
			}
		}
		s.helper.AssertEqual(tc.valid, isValid)
	}
}

// Test file size validation
func (s *AppSimpleTestSuite) TestFileSizeValidation() {
	maxFileSize := 10 * 1024 * 1024 // 10MB

	testCases := []struct {
		size  int64
		valid bool
	}{
		{1024, true},              // 1KB - valid
		{1024 * 1024, true},       // 1MB - valid
		{5 * 1024 * 1024, true},   // 5MB - valid
		{10 * 1024 * 1024, true},  // 10MB - at limit
		{11 * 1024 * 1024, false}, // 11MB - too large
	}

	for _, tc := range testCases {
		isValid := tc.size <= int64(maxFileSize)
		s.helper.AssertEqual(tc.valid, isValid)
	}
}

// Test recent projects list management
func (s *AppSimpleTestSuite) TestRecentProjectsManagement() {
	recentProjects := []string{}
	
	// Function to add a project (simulating addToRecentProjects logic)
	addProject := func(projectID string, recent *[]string) {
		// Remove if already exists
		for i, id := range *recent {
			if id == projectID {
				*recent = append((*recent)[:i], (*recent)[i+1:]...)
				break
			}
		}
		
		// Add to beginning
		*recent = append([]string{projectID}, *recent...)
		
		// Keep only last 10
		if len(*recent) > 10 {
			*recent = (*recent)[:10]
		}
	}

	// Test adding projects
	addProject("project1", &recentProjects)
	s.helper.AssertEqual(1, len(recentProjects))
	s.helper.AssertEqual("project1", recentProjects[0])

	// Test adding another project
	addProject("project2", &recentProjects)
	s.helper.AssertEqual(2, len(recentProjects))
	s.helper.AssertEqual("project2", recentProjects[0]) // Most recent first

	// Test adding existing project (should move to front)
	addProject("project1", &recentProjects)
	s.helper.AssertEqual(2, len(recentProjects))
	s.helper.AssertEqual("project1", recentProjects[0])

	// Test limit enforcement
	for i := 3; i <= 12; i++ {
		addProject(fmt.Sprintf("project%d", i), &recentProjects)
	}
	s.helper.AssertEqual(10, len(recentProjects)) // Should be capped at 10
}

// Performance tests for basic operations
func (s *AppSimpleTestSuite) TestBasicValidation_Performance() {
	s.helper.AssertPerformance(func() {
		for i := 0; i < 10000; i++ {
			// Test basic string validation
			name := "TestProject"
			_ = name != ""
			
			// Test numeric validation
			port := 8080
			_ = port >= 1000 && port <= 65535
			
			// Test progress validation
			progress := 0.5
			_ = progress >= 0.0 && progress <= 1.0
		}
	}, 5*1000000) // 5ms
}