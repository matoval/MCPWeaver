package helpers

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"testing"
	"time"

	"MCPWeaver/internal/database"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestHelper provides common testing utilities
type TestHelper struct {
	t *testing.T
}

// NewTestHelper creates a new test helper
func NewTestHelper(t *testing.T) *TestHelper {
	return &TestHelper{t: t}
}

// CreateTempDB creates a temporary SQLite database for testing
func (h *TestHelper) CreateTempDB() (*database.DB, string, func()) {
	tempDir := h.t.TempDir()
	dbPath := filepath.Join(tempDir, "test.db")
	
	db, err := database.Open(dbPath)
	require.NoError(h.t, err, "Failed to create temp database")
	
	cleanup := func() {
		if db != nil {
			db.Close()
		}
		os.RemoveAll(tempDir)
	}
	
	return db, dbPath, cleanup
}

// CreateMockDB creates a mock database for testing
func (h *TestHelper) CreateMockDB() (*sql.DB, sqlmock.Sqlmock, func()) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	require.NoError(h.t, err, "Failed to create mock database")
	
	cleanup := func() {
		db.Close()
	}
	
	return db, mock, cleanup
}

// CreateTempFile creates a temporary file with given content
func (h *TestHelper) CreateTempFile(content string, extension string) (string, func()) {
	tempDir := h.t.TempDir()
	fileName := fmt.Sprintf("test_%d%s", time.Now().UnixNano(), extension)
	filePath := filepath.Join(tempDir, fileName)
	
	err := os.WriteFile(filePath, []byte(content), 0644)
	require.NoError(h.t, err, "Failed to create temp file")
	
	cleanup := func() {
		os.RemoveAll(tempDir)
	}
	
	return filePath, cleanup
}

// CreateTempDir creates a temporary directory
func (h *TestHelper) CreateTempDir() (string, func()) {
	tempDir := h.t.TempDir()
	
	cleanup := func() {
		os.RemoveAll(tempDir)
	}
	
	return tempDir, cleanup
}

// AssertNoError asserts that error is nil
func (h *TestHelper) AssertNoError(err error, msgAndArgs ...interface{}) {
	assert.NoError(h.t, err, msgAndArgs...)
}

// AssertError asserts that error is not nil
func (h *TestHelper) AssertError(err error, msgAndArgs ...interface{}) {
	assert.Error(h.t, err, msgAndArgs...)
}

// AssertEqual asserts that two values are equal
func (h *TestHelper) AssertEqual(expected, actual interface{}, msgAndArgs ...interface{}) {
	assert.Equal(h.t, expected, actual, msgAndArgs...)
}

// AssertNotNil asserts that value is not nil
func (h *TestHelper) AssertNotNil(value interface{}, msgAndArgs ...interface{}) {
	assert.NotNil(h.t, value, msgAndArgs...)
}

// AssertNil asserts that value is nil
func (h *TestHelper) AssertNil(value interface{}, msgAndArgs ...interface{}) {
	assert.Nil(h.t, value, msgAndArgs...)
}

// AssertNotEqual asserts that two values are not equal
func (h *TestHelper) AssertNotEqual(expected, actual interface{}, msgAndArgs ...interface{}) {
	assert.NotEqual(h.t, expected, actual, msgAndArgs...)
}

// AssertContains asserts that string contains substring
func (h *TestHelper) AssertContains(str, substr string, msgAndArgs ...interface{}) {
	assert.Contains(h.t, str, substr, msgAndArgs...)
}

// AssertFileExists asserts that file exists
func (h *TestHelper) AssertFileExists(filePath string, msgAndArgs ...interface{}) {
	_, err := os.Stat(filePath)
	assert.NoError(h.t, err, msgAndArgs...)
}

// AssertFileNotExists asserts that file does not exist
func (h *TestHelper) AssertFileNotExists(filePath string, msgAndArgs ...interface{}) {
	_, err := os.Stat(filePath)
	assert.True(h.t, os.IsNotExist(err), msgAndArgs...)
}

// ValidateOpenAPISpec returns a sample valid OpenAPI spec for testing
func ValidateOpenAPISpec() string {
	return `
openapi: 3.0.0
info:
  title: Test API
  version: 1.0.0
  description: A test API for unit testing
servers:
  - url: https://api.test.com
    description: Test server
paths:
  /users:
    get:
      summary: Get users
      responses:
        '200':
          description: List of users
          content:
            application/json:
              schema:
                type: array
                items:
                  $ref: '#/components/schemas/User'
    post:
      summary: Create user
      requestBody:
        required: true
        content:
          application/json:
            schema:
              $ref: '#/components/schemas/User'
      responses:
        '201':
          description: User created
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/User'
  /users/{id}:
    get:
      summary: Get user by ID
      parameters:
        - name: id
          in: path
          required: true
          schema:
            type: string
      responses:
        '200':
          description: User details
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/User'
        '404':
          description: User not found
components:
  schemas:
    User:
      type: object
      required:
        - id
        - name
        - email
      properties:
        id:
          type: string
          description: Unique user ID
        name:
          type: string
          description: User's full name
        email:
          type: string
          format: email
          description: User's email address
        createdAt:
          type: string
          format: date-time
          description: Account creation timestamp
`
}

// InvalidOpenAPISpec returns an invalid OpenAPI spec for testing error scenarios
func InvalidOpenAPISpec() string {
	return `
openapi: 3.0.0
info:
  # Missing title and version
paths:
  /invalid:
    get:
      # Missing responses
`
}

// LargeOpenAPISpec returns a large OpenAPI spec for performance testing
func LargeOpenAPISpec() string {
	spec := `
openapi: 3.0.0
info:
  title: Large Test API
  version: 1.0.0
  description: A large API for performance testing
servers:
  - url: https://api.large.com
paths:
`
	
	// Add many paths for performance testing
	for i := 0; i < 100; i++ {
		spec += fmt.Sprintf(`
  /endpoint%d:
    get:
      summary: Get endpoint %d
      responses:
        '200':
          description: Success
          content:
            application/json:
              schema:
                type: object
                properties:
                  id:
                    type: integer
                  name:
                    type: string
    post:
      summary: Create endpoint %d
      requestBody:
        required: true
        content:
          application/json:
            schema:
              type: object
              properties:
                name:
                  type: string
      responses:
        '201':
          description: Created
`, i, i, i)
	}
	
	return spec
}

// GetTestFixturePath returns the path to test fixtures
func GetTestFixturePath(filename string) string {
	return filepath.Join("fixtures", filename)
}

// WaitForCondition waits for a condition to be true with timeout
func (h *TestHelper) WaitForCondition(condition func() bool, timeout time.Duration, message string) {
	deadline := time.Now().Add(timeout)
	
	for time.Now().Before(deadline) {
		if condition() {
			return
		}
		time.Sleep(10 * time.Millisecond)
	}
	
	h.t.Fatalf("Condition not met within timeout: %s", message)
}

// MeasureTime measures execution time of a function
func (h *TestHelper) MeasureTime(fn func()) time.Duration {
	start := time.Now()
	fn()
	return time.Since(start)
}

// AssertPerformance asserts that function executes within expected time
func (h *TestHelper) AssertPerformance(fn func(), maxDuration time.Duration, msgAndArgs ...interface{}) {
	duration := h.MeasureTime(fn)
	assert.True(h.t, duration <= maxDuration, 
		"Expected execution time <= %v, got %v. %v", 
		maxDuration, duration, fmt.Sprint(msgAndArgs...))
}

// AssertMemoryUsage asserts that function doesn't use more than expected memory
func (h *TestHelper) AssertMemoryUsage(fn func(), maxMemoryMB float64, msgAndArgs ...interface{}) {
	// This is a simplified memory measurement
	// In a real scenario, you might want to use more sophisticated memory profiling
	var m1, m2 runtime.MemStats
	runtime.ReadMemStats(&m1)
	runtime.GC()
	
	fn()
	
	runtime.GC()
	runtime.ReadMemStats(&m2)
	
	memoryUsedMB := float64(m2.Alloc-m1.Alloc) / 1024 / 1024
	assert.True(h.t, memoryUsedMB <= maxMemoryMB,
		"Expected memory usage <= %.2f MB, got %.2f MB. %v",
		maxMemoryMB, memoryUsedMB, fmt.Sprint(msgAndArgs...))
}