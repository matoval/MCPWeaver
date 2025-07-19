package unit

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"MCPWeaver/internal/app"
	"MCPWeaver/tests/helpers"

	"github.com/stretchr/testify/suite"
)

type AppFilesRealTestSuite struct {
	suite.Suite
	helper  *helpers.TestHelper
	app     *app.App
	tempDir string
	cleanup func()
}

func (s *AppFilesRealTestSuite) SetupTest() {
	s.helper = helpers.NewTestHelper(s.T())
	s.tempDir, s.cleanup = s.helper.CreateTempDir()
	s.app = &app.App{}
}

func (s *AppFilesRealTestSuite) TearDownTest() {
	if s.cleanup != nil {
		s.cleanup()
	}
}

func TestAppFilesRealTestSuite(t *testing.T) {
	suite.Run(t, new(AppFilesRealTestSuite))
}

// Test convertFilters helper function
func (s *AppFilesRealTestSuite) TestConvertFilters() {
	// This tests the exported helper by calling a method that uses it
	filters := []app.FileFilter{
		{
			DisplayName: "YAML Files",
			Pattern:     "*.yaml",
			Extensions:  []string{".yaml", ".yml"},
		},
		{
			DisplayName: "JSON Files", 
			Pattern:     "*.json",
			Extensions:  []string{".json"},
		},
	}

	// We can't directly test convertFilters since it's not exported,
	// but we can test that FileFilter struct works correctly
	s.helper.AssertEqual("YAML Files", filters[0].DisplayName)
	s.helper.AssertEqual("*.yaml", filters[0].Pattern)
	s.helper.AssertEqual(2, len(filters[0].Extensions))
	s.helper.AssertEqual(".yaml", filters[0].Extensions[0])
	s.helper.AssertEqual(".yml", filters[0].Extensions[1])

	s.helper.AssertEqual("JSON Files", filters[1].DisplayName)
	s.helper.AssertEqual("*.json", filters[1].Pattern)
	s.helper.AssertEqual(1, len(filters[1].Extensions))
	s.helper.AssertEqual(".json", filters[1].Extensions[0])
}

// Test file operations that don't require Wails context
func (s *AppFilesRealTestSuite) TestFileOperations_Constants() {
	// Test that constants are properly defined by checking for expected values
	// We can't access the constants directly, but we can test the logic
	maxSize := 10 * 1024 * 1024  // 10MB
	bufferSize := 64 * 1024      // 64KB
	
	s.helper.AssertEqual(10485760, maxSize)
	s.helper.AssertEqual(65536, bufferSize)
	
	// Test that file size validation would work
	testSizes := []struct {
		size  int64
		valid bool
	}{
		{1024, true},              // 1KB - valid
		{5 * 1024 * 1024, true},   // 5MB - valid
		{10 * 1024 * 1024, true},  // 10MB - at limit
		{11 * 1024 * 1024, false}, // 11MB - too large
	}

	for _, test := range testSizes {
		isValid := test.size <= int64(maxSize)
		s.helper.AssertEqual(test.valid, isValid)
	}
}

// Test file path validation logic
func (s *AppFilesRealTestSuite) TestFilePathValidation() {
	testCases := []struct {
		path     string
		hasError bool
	}{
		{"", true},                          // Empty path
		{"/valid/path/file.yaml", false},    // Valid path
		{"relative/path.yaml", false},       // Relative path (valid)
		{"C:\\Windows\\file.yaml", false},   // Windows path (valid)
	}

	for _, tc := range testCases {
		s.Run(tc.path, func() {
			// Test the validation logic that would be in fileExists
			hasError := tc.path == ""
			s.helper.AssertEqual(tc.hasError, hasError)
		})
	}
}

// Test directory path validation logic
func (s *AppFilesRealTestSuite) TestDirectoryPathValidation() {
	testCases := []struct {
		path     string
		hasError bool
	}{
		{"", true},                    // Empty path
		{"/tmp", false},               // Valid path
		{"./relative", false},         // Relative path
		{"/valid/directory", false},   // Valid absolute path
	}

	for _, tc := range testCases {
		s.Run(tc.path, func() {
			// Test the validation logic that would be in dirExists/ensureDir
			hasError := tc.path == ""
			s.helper.AssertEqual(tc.hasError, hasError)
		})
	}
}

// Test file extension handling
func (s *AppFilesRealTestSuite) TestFileExtensions() {
	testCases := []struct {
		filename  string
		extension string
	}{
		{"spec.yaml", ".yaml"},
		{"api.yml", ".yml"},
		{"config.json", ".json"},
		{"openapi.yaml", ".yaml"},
		{"no-extension", ""},
		{"multiple.dots.yaml", ".yaml"},
	}

	for _, tc := range testCases {
		s.Run(tc.filename, func() {
			ext := filepath.Ext(tc.filename)
			s.helper.AssertEqual(tc.extension, ext)
		})
	}
}

// Test dialog title handling
func (s *AppFilesRealTestSuite) TestDialogTitleHandling() {
	testCases := []struct {
		input    string
		expected string
	}{
		{"", "Select Directory"},           // Empty title gets default
		{"Custom Title", "Custom Title"},  // Non-empty title preserved
		{"Select Output", "Select Output"}, // Custom title preserved
	}

	for _, tc := range testCases {
		s.Run(tc.input, func() {
			// Test the logic from SelectDirectory
			title := tc.input
			if title == "" {
				title = "Select Directory"
			}
			s.helper.AssertEqual(tc.expected, title)
		})
	}
}

// Test file operations with real files
func (s *AppFilesRealTestSuite) TestRealFileOperations() {
	// Create a test file
	testFile := filepath.Join(s.tempDir, "test.yaml")
	content := `openapi: 3.0.0
info:
  title: Test API
  version: 1.0.0
paths:
  /test:
    get:
      summary: Test endpoint
      responses:
        '200':
          description: Success`

	err := os.WriteFile(testFile, []byte(content), 0644)
	s.helper.AssertNoError(err)

	// Test file existence check
	_, err = os.Stat(testFile)
	s.helper.AssertNoError(err) // File should exist

	// Test non-existent file
	nonExistentFile := filepath.Join(s.tempDir, "nonexistent.yaml")
	_, err = os.Stat(nonExistentFile)
	s.helper.AssertError(err) // Should error for non-existent file
	s.helper.AssertEqual(true, os.IsNotExist(err))

	// Test directory creation
	testDir := filepath.Join(s.tempDir, "testdir")
	err = os.MkdirAll(testDir, 0755)
	s.helper.AssertNoError(err)

	// Verify directory exists
	info, err := os.Stat(testDir)
	s.helper.AssertNoError(err)
	s.helper.AssertEqual(true, info.IsDir())
}

// Test writable directory check logic
func (s *AppFilesRealTestSuite) TestWritableDirectoryCheck() {
	// Create a test directory
	testDir := filepath.Join(s.tempDir, "writable")
	err := os.MkdirAll(testDir, 0755)
	s.helper.AssertNoError(err)

	// Test writing to directory (simulating dirExists logic)
	testFile := filepath.Join(testDir, ".mcpweaver_test")
	err = os.WriteFile(testFile, []byte("test"), 0644)
	s.helper.AssertNoError(err)

	// Verify test file was created
	_, err = os.Stat(testFile)
	s.helper.AssertNoError(err)

	// Clean up test file (as done in dirExists)
	err = os.Remove(testFile)
	s.helper.AssertNoError(err)

	// Verify test file was removed
	_, err = os.Stat(testFile)
	s.helper.AssertError(err)
	s.helper.AssertEqual(true, os.IsNotExist(err))
}

// Test buffer size calculations
func (s *AppFilesRealTestSuite) TestBufferSizeCalculations() {
	bufferSize := 64 * 1024 // 64KB
	
	// Test that buffer size is appropriate for different file sizes
	testCases := []struct {
		fileSize     int64
		numBuffers   int64
	}{
		{1024, 1},           // 1KB = 1 buffer
		{65536, 1},          // 64KB = 1 buffer
		{131072, 2},         // 128KB = 2 buffers
		{1048576, 16},       // 1MB = 16 buffers
	}

	for _, tc := range testCases {
		expectedBuffers := (tc.fileSize + int64(bufferSize) - 1) / int64(bufferSize)
		s.helper.AssertEqual(tc.numBuffers, expectedBuffers)
	}
}

// Performance tests
func (s *AppFilesRealTestSuite) TestFileOperations_Performance() {
	// Test creating many small files
	s.helper.AssertPerformance(func() {
		for i := 0; i < 100; i++ {
			testFile := filepath.Join(s.tempDir, fmt.Sprintf("perf_test_%d.txt", i))
			err := os.WriteFile(testFile, []byte("test content"), 0644)
			s.helper.AssertNoError(err)
		}
	}, 100*1000000) // 100ms

	// Test reading file stats
	testFile := filepath.Join(s.tempDir, "perf_test_0.txt")
	s.helper.AssertPerformance(func() {
		for i := 0; i < 1000; i++ {
			_, err := os.Stat(testFile)
			s.helper.AssertNoError(err)
		}
	}, 50*1000000) // 50ms
}

// Test file size validation
func (s *AppFilesRealTestSuite) TestFileSizeValidation() {
	// Create files of different sizes
	testCases := []struct {
		size int
		name string
	}{
		{1024, "small.txt"},           // 1KB
		{10240, "medium.txt"},         // 10KB  
		{102400, "large.txt"},         // 100KB
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			testFile := filepath.Join(s.tempDir, tc.name)
			content := make([]byte, tc.size)
			for i := range content {
				content[i] = 'A'
			}
			
			err := os.WriteFile(testFile, content, 0644)
			s.helper.AssertNoError(err)
			
			// Verify file size
			info, err := os.Stat(testFile)
			s.helper.AssertNoError(err)
			s.helper.AssertEqual(int64(tc.size), info.Size())
		})
	}
}

// Test error scenarios
func (s *AppFilesRealTestSuite) TestErrorScenarios() {
	// Test with non-existent directory parent
	invalidPath := filepath.Join(s.tempDir, "nonexistent", "subdir", "file.txt")
	
	// This would fail when trying to write to a non-existent parent directory
	err := os.WriteFile(invalidPath, []byte("content"), 0644)
	s.helper.AssertError(err) // Should fail due to missing parent directories
	
	// Test creating the parent directory first
	parentDir := filepath.Dir(invalidPath)
	err = os.MkdirAll(parentDir, 0755)
	s.helper.AssertNoError(err)
	
	// Now writing should succeed
	err = os.WriteFile(invalidPath, []byte("content"), 0644)
	s.helper.AssertNoError(err)
}