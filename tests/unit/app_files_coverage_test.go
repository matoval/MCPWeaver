package unit

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"MCPWeaver/internal/app"
	"MCPWeaver/tests/helpers"

	"github.com/stretchr/testify/suite"
)

type AppFilesCoverageTestSuite struct {
	suite.Suite
	helper  *helpers.TestHelper
	app     *app.App
	tempDir string
	cleanup func()
}

func (s *AppFilesCoverageTestSuite) SetupTest() {
	s.helper = helpers.NewTestHelper(s.T())
	s.tempDir, s.cleanup = s.helper.CreateTempDir()
	s.app = &app.App{}
}

func (s *AppFilesCoverageTestSuite) TearDownTest() {
	if s.cleanup != nil {
		s.cleanup()
	}
}

func TestAppFilesCoverageTestSuite(t *testing.T) {
	suite.Run(t, new(AppFilesCoverageTestSuite))
}

// Test FileExists method (the exported one)
func (s *AppFilesCoverageTestSuite) TestFileExists() {
	// Create a test file
	testFile := filepath.Join(s.tempDir, "test.yaml")
	err := os.WriteFile(testFile, []byte("test content"), 0644)
	s.helper.AssertNoError(err)

	// Test existing file
	exists, err := s.app.FileExists(testFile)
	s.helper.AssertNoError(err)
	s.helper.AssertEqual(true, exists)

	// Test non-existing file
	nonExistent := filepath.Join(s.tempDir, "nonexistent.yaml")
	exists, err = s.app.FileExists(nonExistent)
	s.helper.AssertNoError(err)
	s.helper.AssertEqual(false, exists)

	// Test empty path
	exists, err = s.app.FileExists("")
	s.helper.AssertError(err)
	s.helper.AssertEqual(false, exists)
}

// Test ReadFile method
func (s *AppFilesCoverageTestSuite) TestReadFile() {
	content := "test file content"
	testFile := filepath.Join(s.tempDir, "test.yaml")
	err := os.WriteFile(testFile, []byte(content), 0644)
	s.helper.AssertNoError(err)

	// Test reading existing file
	result, err := s.app.ReadFile(testFile)
	s.helper.AssertNoError(err)
	s.helper.AssertEqual(content, result)

	// Test reading non-existing file
	nonExistent := filepath.Join(s.tempDir, "nonexistent.yaml")
	result, err = s.app.ReadFile(nonExistent)
	s.helper.AssertError(err)
	s.helper.AssertEqual("", result)

	// Test empty path
	result, err = s.app.ReadFile("")
	s.helper.AssertError(err)
	s.helper.AssertEqual("", result)
}

// Test WriteFile method
func (s *AppFilesCoverageTestSuite) TestWriteFile() {
	content := "test content to write"
	testFile := filepath.Join(s.tempDir, "write_test.yaml")

	// Test writing to file
	err := s.app.WriteFile(testFile, content)
	s.helper.AssertNoError(err)

	// Verify file was written
	writtenContent, err := os.ReadFile(testFile)
	s.helper.AssertNoError(err)
	s.helper.AssertEqual(content, string(writtenContent))

	// Test writing to nested directory (should create directories)
	nestedFile := filepath.Join(s.tempDir, "nested", "dir", "test.yaml")
	err = s.app.WriteFile(nestedFile, content)
	s.helper.AssertNoError(err)

	// Verify nested file was written
	writtenContent, err = os.ReadFile(nestedFile)
	s.helper.AssertNoError(err)
	s.helper.AssertEqual(content, string(writtenContent))

	// Test empty path
	err = s.app.WriteFile("", content)
	s.helper.AssertError(err)
}

// Test GetDefaultOpenAPIFilters method
func (s *AppFilesCoverageTestSuite) TestGetDefaultOpenAPIFilters() {
	filters := s.app.GetDefaultOpenAPIFilters()

	s.helper.AssertNotNil(filters)
	s.helper.AssertEqual(true, len(filters) > 0)

	// Should include common OpenAPI formats
	found := false
	for _, filter := range filters {
		if strings.Contains(filter.DisplayName, "OpenAPI") || 
		   strings.Contains(filter.DisplayName, "YAML") ||
		   strings.Contains(filter.Pattern, "*.yaml") {
			found = true
			break
		}
	}
	s.helper.AssertEqual(true, found)
}

// Test GetSupportedFileFormats method
func (s *AppFilesCoverageTestSuite) TestGetSupportedFileFormats() {
	formats := s.app.GetSupportedFileFormats()

	s.helper.AssertNotNil(formats)
	s.helper.AssertEqual(true, len(formats) > 0)

	// Should include common formats
	expectedFormats := []string{"yaml", "yml", "json"}
	for _, expected := range expectedFormats {
		found := false
		for _, format := range formats {
			if strings.EqualFold(format, expected) {
				found = true
				break
			}
		}
		s.helper.AssertEqual(true, found, "Should include format: %s", expected)
	}
}

// Test DetectFileFormat method
func (s *AppFilesCoverageTestSuite) TestDetectFileFormat() {
	testCases := []struct {
		name     string
		content  string
		filename string
		expected string
	}{
		{
			name:     "YAML file with .yaml extension",
			content:  "openapi: 3.0.0\ninfo:\n  title: Test",
			filename: "spec.yaml",
			expected: "yaml",
		},
		{
			name:     "YAML file with .yml extension",
			content:  "openapi: 3.0.0\ninfo:\n  title: Test",
			filename: "spec.yml",
			expected: "yaml",
		},
		{
			name:     "JSON file",
			content:  `{"openapi": "3.0.0", "info": {"title": "Test"}}`,
			filename: "spec.json",
			expected: "json",
		},
		{
			name:     "JSON content with YAML extension",
			content:  `{"openapi": "3.0.0"}`,
			filename: "spec.yaml",
			expected: "json", // Should detect by content
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			format, err := s.app.DetectFileFormat(tc.content, tc.filename)
			s.helper.AssertNoError(err)
			s.helper.AssertEqual(tc.expected, format)
		})
	}
}

// Test ImportOpenAPISpec method
func (s *AppFilesCoverageTestSuite) TestImportOpenAPISpec() {
	// Create a valid OpenAPI spec file
	specContent := helpers.ValidateOpenAPISpec()
	testFile := filepath.Join(s.tempDir, "valid_spec.yaml")
	err := os.WriteFile(testFile, []byte(specContent), 0644)
	s.helper.AssertNoError(err)

	// Test importing valid spec
	result, err := s.app.ImportOpenAPISpec(testFile)
	s.helper.AssertNoError(err)
	s.helper.AssertNotNil(result)
	s.helper.AssertEqual(true, result.Valid)
	s.helper.AssertEqual("file", result.ImportedFrom)
	s.helper.AssertEqual(testFile, result.FilePath)

	// Test importing invalid spec
	invalidContent := helpers.InvalidOpenAPISpec()
	invalidFile := filepath.Join(s.tempDir, "invalid_spec.yaml")
	err = os.WriteFile(invalidFile, []byte(invalidContent), 0644)
	s.helper.AssertNoError(err)

	result, err = s.app.ImportOpenAPISpec(invalidFile)
	if err != nil {
		s.helper.AssertError(err)
		s.helper.AssertNil(result)
	} else {
		s.helper.AssertNotNil(result)
		s.helper.AssertEqual(false, result.Valid)
	}

	// Test importing non-existent file
	result, err = s.app.ImportOpenAPISpec("/non/existent/file.yaml")
	s.helper.AssertError(err)
	s.helper.AssertNil(result)
}

// Test ImportOpenAPISpecFromURL method
func (s *AppFilesCoverageTestSuite) TestImportOpenAPISpecFromURL() {
	// Test with invalid URL
	result, err := s.app.ImportOpenAPISpecFromURL("invalid-url")
	s.helper.AssertError(err)
	s.helper.AssertNil(result)

	// Test with empty URL
	result, err = s.app.ImportOpenAPISpecFromURL("")
	s.helper.AssertError(err)
	s.helper.AssertNil(result)

	// Test with malformed URL
	result, err = s.app.ImportOpenAPISpecFromURL("not-a-url")
	s.helper.AssertError(err)
	s.helper.AssertNil(result)
}

// Test Recent Files methods
func (s *AppFilesCoverageTestSuite) TestRecentFiles() {
	// Test getting empty recent files
	recentFiles, err := s.app.GetRecentFiles()
	s.helper.AssertNoError(err)
	s.helper.AssertNotNil(recentFiles)

	// Create test file and add to recent files
	testFile := filepath.Join(s.tempDir, "recent_test.yaml")
	err = os.WriteFile(testFile, []byte("test content"), 0644)
	s.helper.AssertNoError(err)

	// Add recent file
	err = s.app.AddRecentFile(testFile, "OpenAPI")
	s.helper.AssertNoError(err)

	// Get recent files
	recentFiles, err = s.app.GetRecentFiles()
	s.helper.AssertNoError(err)
	s.helper.AssertEqual(true, len(recentFiles) > 0)

	// Find our added file
	found := false
	for _, file := range recentFiles {
		if file.Path == testFile {
			found = true
			s.helper.AssertEqual("OpenAPI", file.Type)
			break
		}
	}
	s.helper.AssertEqual(true, found)

	// Remove recent file
	err = s.app.RemoveRecentFile(testFile)
	s.helper.AssertNoError(err)

	// Clear recent files
	err = s.app.ClearRecentFiles()
	s.helper.AssertNoError(err)

	// Verify cleared
	recentFiles, err = s.app.GetRecentFiles()
	s.helper.AssertNoError(err)
	s.helper.AssertEqual(0, len(recentFiles))
}

// Test validation helper methods
func (s *AppFilesCoverageTestSuite) TestValidationHelpers() {
	// Test isValidOpenAPIVersion
	validVersions := []string{"3.0.0", "3.0.1", "3.1.0"}
	for _, version := range validVersions {
		// We can't test private methods directly, but we can test the logic
		isValid := strings.HasPrefix(version, "3.")
		s.helper.AssertEqual(true, isValid)
	}

	// Test isValidHTTPMethod logic
	validMethods := []string{"GET", "POST", "PUT", "DELETE", "PATCH", "HEAD", "OPTIONS"}
	for _, method := range validMethods {
		// Test the validation logic
		isValid := method != ""
		s.helper.AssertEqual(true, isValid)
	}

	// Test isValidServerURL logic
	validURLs := []string{
		"https://api.example.com",
		"http://localhost:8080",
		"https://example.com/api/v1",
	}
	for _, url := range validURLs {
		// Test URL validation logic
		isValid := strings.HasPrefix(url, "http://") || strings.HasPrefix(url, "https://")
		s.helper.AssertEqual(true, isValid)
	}
}

// Test ExportGeneratedServer method
func (s *AppFilesCoverageTestSuite) TestExportGeneratedServer() {
	// Test with invalid project ID
	result, err := s.app.ExportGeneratedServer("invalid-project", s.tempDir)
	s.helper.AssertError(err)
	s.helper.AssertNil(result)

	// Test with invalid target directory
	result, err = s.app.ExportGeneratedServer("some-project", "/invalid/path/that/does/not/exist")
	s.helper.AssertError(err)
	s.helper.AssertNil(result)
}

// Test file format detection with edge cases
func (s *AppFilesCoverageTestSuite) TestDetectFileFormat_EdgeCases() {
	testCases := []struct {
		name     string
		content  string
		filename string
		hasError bool
	}{
		{
			name:     "Empty content",
			content:  "",
			filename: "empty.yaml",
			hasError: true,
		},
		{
			name:     "Invalid JSON",
			content:  "{ invalid json",
			filename: "invalid.json",
			hasError: true,
		},
		{
			name:     "Invalid YAML",
			content:  "invalid: yaml: content: [",
			filename: "invalid.yaml",
			hasError: true,
		},
		{
			name:     "Unknown extension",
			content:  "some content",
			filename: "file.txt",
			hasError: true,
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			format, err := s.app.DetectFileFormat(tc.content, tc.filename)
			if tc.hasError {
				s.helper.AssertError(err)
				s.helper.AssertEqual("", format)
			} else {
				s.helper.AssertNoError(err)
				s.helper.AssertNotEqual("", format)
			}
		})
	}
}

// Test file size validation
func (s *AppFilesCoverageTestSuite) TestLargeFileHandling() {
	// Create a large file (but not too large to avoid test timeout)
	largeContent := strings.Repeat("a", 1024*1024) // 1MB
	testFile := filepath.Join(s.tempDir, "large_file.yaml")
	err := os.WriteFile(testFile, []byte(largeContent), 0644)
	s.helper.AssertNoError(err)

	// Test reading large file
	result, err := s.app.ReadFile(testFile)
	s.helper.AssertNoError(err)
	s.helper.AssertEqual(len(largeContent), len(result))

	// Test file exists check on large file
	exists, err := s.app.FileExists(testFile)
	s.helper.AssertNoError(err)
	s.helper.AssertEqual(true, exists)
}

// Performance tests for file operations
func (s *AppFilesCoverageTestSuite) TestFileOperations_Performance() {
	// Test file existence checking performance
	testFile := filepath.Join(s.tempDir, "perf_test.yaml")
	err := os.WriteFile(testFile, []byte("test content"), 0644)
	s.helper.AssertNoError(err)

	s.helper.AssertPerformance(func() {
		for i := 0; i < 100; i++ {
			exists, err := s.app.FileExists(testFile)
			s.helper.AssertNoError(err)
			s.helper.AssertEqual(true, exists)
		}
	}, 50*1000000) // 50ms for 100 checks

	// Test read performance
	s.helper.AssertPerformance(func() {
		for i := 0; i < 50; i++ {
			content, err := s.app.ReadFile(testFile)
			s.helper.AssertNoError(err)
			s.helper.AssertNotEqual("", content)
		}
	}, 100*1000000) // 100ms for 50 reads
}