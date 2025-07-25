package security

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"MCPWeaver/internal/app"
)

// SecurityTestSuite contains all security-related tests
type SecurityTestSuite struct {
	app *app.App
}

// NewSecurityTestSuite creates a new security test suite
func NewSecurityTestSuite() *SecurityTestSuite {
	testApp := &app.App{}
	return &SecurityTestSuite{app: testApp}
}

// TestInputValidation tests various input validation scenarios
func TestInputValidation(t *testing.T) {
	suite := NewSecurityTestSuite()
	validator := app.NewSecurityValidator(suite.app)

	tests := []struct {
		name        string
		input       string
		validate    func(string) error
		expectError bool
		errorCode   string
	}{
		{
			name:        "Valid URL",
			input:       "https://api.example.com/openapi.yaml",
			validate:    validator.ValidateURL,
			expectError: false,
		},
		{
			name:        "Invalid URL scheme",
			input:       "ftp://example.com/file.yaml",
			validate:    validator.ValidateURL,
			expectError: true,
			errorCode:   "INVALID_URL_SCHEME",
		},
		{
			name:        "SSRF localhost",
			input:       "http://localhost:8080/admin",
			validate:    validator.ValidateURL,
			expectError: true,
			errorCode:   "SSRF_LOCALHOST_BLOCKED",
		},
		{
			name:        "SSRF private IP",
			input:       "http://192.168.1.1/config",
			validate:    validator.ValidateURL,
			expectError: true,
			errorCode:   "SSRF_PRIVATE_IP_BLOCKED",
		},
		{
			name:        "SSRF metadata endpoint",
			input:       "http://169.254.169.254/metadata",
			validate:    validator.ValidateURL,
			expectError: true,
			errorCode:   "SSRF_METADATA_BLOCKED",
		},
		{
			name:        "Valid file path",
			input:       "project/spec.yaml",
			validate:    validator.ValidateFilePath,
			expectError: false,
		},
		{
			name:        "Path traversal",
			input:       "../../../etc/passwd",
			validate:    validator.ValidateFilePath,
			expectError: true,
			errorCode:   "PATH_TRAVERSAL_DETECTED",
		},
		{
			name:        "Valid template name",
			input:       "server-template",
			validate:    validator.ValidateTemplateName,
			expectError: false,
		},
		{
			name:        "Invalid template name",
			input:       "template/../evil",
			validate:    validator.ValidateTemplateName,
			expectError: true,
			errorCode:   "INVALID_TEMPLATE_NAME",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.validate(tt.input)

			if tt.expectError {
				if err == nil {
					t.Errorf("Expected error for input %q, but got none", tt.input)
					return
				}

				// Check error type if specified
				if tt.errorCode != "" {
					if apiErr, ok := err.(*app.APIError); ok {
						if apiErr.Code != tt.errorCode {
							t.Errorf("Expected error code %q, got %q", tt.errorCode, apiErr.Code)
						}
					} else {
						t.Errorf("Expected APIError, got %T", err)
					}
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error for input %q: %v", tt.input, err)
				}
			}
		})
	}
}

// TestTemplateSanitization tests template variable sanitization
func TestTemplateSanitization(t *testing.T) {
	suite := NewSecurityTestSuite()
	validator := app.NewSecurityValidator(suite.app)

	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "Normal text",
			input:    "Hello World",
			expected: "Hello World",
		},
		{
			name:     "Template injection attempt",
			input:    "{{.System.Execute \"rm -rf /\"}}",
			expected: "&#123;&#123;.System.Execute \"rm -rf /\"&#125;&#125;",
		},
		{
			name:     "JSP injection attempt",
			input:    "<% System.exit(0); %>",
			expected: "&#60;&#37; System.exit(0); &#37;&#62;",
		},
		{
			name:     "Control characters",
			input:    "Hello\x00\x01\x02World",
			expected: "HelloWorld",
		},
		{
			name:     "Mixed injection",
			input:    "{{range .}}<%=evil%>{{end}}",
			expected: "&#123;&#123;range .&#125;&#125;&#60;&#37;=evil&#37;&#62;&#123;&#123;end&#125;&#125;",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := validator.SanitizeTemplateVariable(tt.input)
			if result != tt.expected {
				t.Errorf("Expected %q, got %q", tt.expected, result)
			}
		})
	}
}

// TestFileSystemSecurity tests secure file system operations
func TestFileSystemSecurity(t *testing.T) {
	suite := NewSecurityTestSuite()

	// Create temporary directory for testing
	tempDir, err := os.MkdirTemp("", "mcpweaver-security-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	secureFS := app.NewSecureFileSystem(suite.app)

	// Test allowed paths
	err = secureFS.SetAllowedPaths([]string{tempDir})
	if err != nil {
		t.Fatalf("Failed to set allowed paths: %v", err)
	}

	tests := []struct {
		name        string
		operation   func() error
		expectError bool
		errorCode   string
	}{
		{
			name: "Write to allowed path",
			operation: func() error {
				return secureFS.SecureWriteFile(filepath.Join(tempDir, "test.txt"), []byte("test"), 0644)
			},
			expectError: false,
		},
		{
			name: "Read from allowed path",
			operation: func() error {
				// First write a file
				secureFS.SecureWriteFile(filepath.Join(tempDir, "read-test.txt"), []byte("content"), 0644)
				_, err := secureFS.SecureReadFile(filepath.Join(tempDir, "read-test.txt"))
				return err
			},
			expectError: false,
		},
		{
			name: "Write outside allowed path",
			operation: func() error {
				return secureFS.SecureWriteFile("/tmp/evil.txt", []byte("evil"), 0644)
			},
			expectError: true,
			errorCode:   "PATH_NOT_ALLOWED",
		},
		{
			name: "Path traversal in write",
			operation: func() error {
				return secureFS.SecureWriteFile(filepath.Join(tempDir, "../evil.txt"), []byte("evil"), 0644)
			},
			expectError: true,
			errorCode:   "PATH_TRAVERSAL_DETECTED",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.operation()

			if tt.expectError {
				if err == nil {
					t.Errorf("Expected error for operation %q, but got none", tt.name)
					return
				}

				if tt.errorCode != "" {
					if apiErr, ok := err.(*app.APIError); ok {
						if apiErr.Code != tt.errorCode {
							t.Errorf("Expected error code %q, got %q", tt.errorCode, apiErr.Code)
						}
					}
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error for operation %q: %v", tt.name, err)
				}
			}
		})
	}
}

// TestHTTPClientSecurity tests secure HTTP client
func TestHTTPClientSecurity(t *testing.T) {
	suite := NewSecurityTestSuite()

	// Create test servers
	maliciousServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Try to serve malicious content
		w.Header().Set("Content-Type", "text/html")
		w.Header().Set("X-Frame-Options", "DENY")
		w.Write([]byte("<script>alert('xss')</script>"))
	}))
	defer maliciousServer.Close()

	validServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/yaml")
		w.Write([]byte("openapi: 3.0.0\ninfo:\n  title: Test API"))
	}))
	defer validServer.Close()

	config := app.GetDefaultHTTPClientConfig()
	client := app.NewSecureHTTPClient(suite.app, config)

	tests := []struct {
		name        string
		url         string
		expectError bool
		errorCode   string
	}{
		{
			name:        "Valid YAML response",
			url:         validServer.URL,
			expectError: false,
		},
		{
			name:        "Malicious HTML response",
			url:         maliciousServer.URL,
			expectError: true,
			errorCode:   "SUSPICIOUS_CONTENT",
		},
		{
			name:        "Localhost URL",
			url:         "http://localhost:8080/api",
			expectError: true,
			errorCode:   "SSRF_LOCALHOST_BLOCKED",
		},
		{
			name:        "Private IP URL",
			url:         "http://192.168.1.1/api",
			expectError: true,
			errorCode:   "SSRF_PRIVATE_IP_BLOCKED",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()

			resp, err := client.Get(ctx, tt.url)
			if resp != nil {
				resp.Body.Close()
			}

			if tt.expectError {
				if err == nil {
					t.Errorf("Expected error for URL %q, but got none", tt.url)
					return
				}

				if tt.errorCode != "" {
					if apiErr, ok := err.(*app.APIError); ok {
						if apiErr.Code != tt.errorCode {
							t.Errorf("Expected error code %q, got %q", tt.errorCode, apiErr.Code)
						}
					}
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error for URL %q: %v", tt.url, err)
				}
			}
		})
	}
}

// TestJSONSecurity tests JSON parsing security
func TestJSONSecurity(t *testing.T) {
	suite := NewSecurityTestSuite()
	validator := app.NewSecurityValidator(suite.app)

	tests := []struct {
		name        string
		json        string
		expectError bool
		errorCode   string
	}{
		{
			name:        "Valid JSON",
			json:        `{"name": "test", "value": 123}`,
			expectError: false,
		},
		{
			name:        "Deeply nested JSON (bomb)",
			json:        strings.Repeat(`{"level":`, 60) + `"deep"` + strings.Repeat(`}`, 60),
			expectError: true,
			errorCode:   "JSON_TOO_DEEP",
		},
		{
			name:        "Normal nesting",
			json:        `{"level1": {"level2": {"level3": "value"}}}`,
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validator.ValidateJSONContent(tt.json)

			if tt.expectError {
				if err == nil {
					t.Errorf("Expected error for JSON %q, but got none", tt.name)
					return
				}

				if tt.errorCode != "" {
					if apiErr, ok := err.(*app.APIError); ok {
						if apiErr.Code != tt.errorCode {
							t.Errorf("Expected error code %q, got %q", tt.errorCode, apiErr.Code)
						}
					}
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error for JSON %q: %v", tt.name, err)
				}
			}
		})
	}
}

// TestYAMLSecurity tests YAML parsing security
func TestYAMLSecurity(t *testing.T) {
	suite := NewSecurityTestSuite()
	validator := app.NewSecurityValidator(suite.app)

	tests := []struct {
		name        string
		yaml        string
		expectError bool
		errorCode   string
	}{
		{
			name: "Valid YAML",
			yaml: `
name: test
value: 123
`,
			expectError: false,
		},
		{
			name:        "YAML with excessive references",
			yaml:        strings.Repeat("&ref ", 150) + "value",
			expectError: true,
			errorCode:   "YAML_EXCESSIVE_REFERENCES",
		},
		{
			name: "YAML with Python object",
			yaml: `
test: !!python/object:os.system
args: ["rm -rf /"]
`,
			expectError: true,
			errorCode:   "YAML_DANGEROUS_CONSTRUCT",
		},
		{
			name: "YAML with Java object",
			yaml: `
test: !!java/object:java.lang.ProcessBuilder
args: ["rm", "-rf", "/"]
`,
			expectError: true,
			errorCode:   "YAML_DANGEROUS_CONSTRUCT",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validator.ValidateYAMLContent(tt.yaml)

			if tt.expectError {
				if err == nil {
					t.Errorf("Expected error for YAML %q, but got none", tt.name)
					return
				}

				if tt.errorCode != "" {
					if apiErr, ok := err.(*app.APIError); ok {
						if apiErr.Code != tt.errorCode {
							t.Errorf("Expected error code %q, got %q", tt.errorCode, apiErr.Code)
						}
					}
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error for YAML %q: %v", tt.name, err)
				}
			}
		})
	}
}

// TestRegexSecurity tests regex validation for ReDoS protection
func TestRegexSecurity(t *testing.T) {
	suite := NewSecurityTestSuite()
	validator := app.NewSecurityValidator(suite.app)

	tests := []struct {
		name        string
		pattern     string
		expectError bool
		errorCode   string
	}{
		{
			name:        "Valid regex",
			pattern:     `^[a-zA-Z0-9]+$`,
			expectError: false,
		},
		{
			name:        "ReDoS pattern - nested quantifiers",
			pattern:     `(.*)*`,
			expectError: true,
			errorCode:   "DANGEROUS_REGEX",
		},
		{
			name:        "ReDoS pattern - catastrophic backtracking",
			pattern:     `(.*)*$`,
			expectError: true,
			errorCode:   "DANGEROUS_REGEX",
		},
		{
			name:        "ReDoS pattern - high repetition",
			pattern:     `(.*){10,}`,
			expectError: true,
			errorCode:   "DANGEROUS_REGEX",
		},
		{
			name:        "Invalid regex syntax",
			pattern:     `[unclosed`,
			expectError: true,
			errorCode:   "INVALID_REGEX",
		},
		{
			name:        "Too long regex",
			pattern:     strings.Repeat("a", 1001),
			expectError: true,
			errorCode:   "REGEX_TOO_LONG",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validator.ValidateRegexPattern(tt.pattern)

			if tt.expectError {
				if err == nil {
					t.Errorf("Expected error for pattern %q, but got none", tt.pattern)
					return
				}

				if tt.errorCode != "" {
					if apiErr, ok := err.(*app.APIError); ok {
						if apiErr.Code != tt.errorCode {
							t.Errorf("Expected error code %q, got %q", tt.errorCode, apiErr.Code)
						}
					}
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error for pattern %q: %v", tt.pattern, err)
				}
			}
		})
	}
}

// TestNetworkSecurity tests network-level security controls
func TestNetworkSecurity(t *testing.T) {
	tests := []struct {
		name        string
		host        string
		expectBlock bool
	}{
		{
			name:        "Public host",
			host:        "api.github.com",
			expectBlock: false,
		},
		{
			name:        "Localhost",
			host:        "localhost",
			expectBlock: true,
		},
		{
			name:        "127.0.0.1",
			host:        "127.0.0.1",
			expectBlock: true,
		},
		{
			name:        "Private IP 192.168.x.x",
			host:        "192.168.1.1",
			expectBlock: true,
		},
		{
			name:        "Private IP 10.x.x.x",
			host:        "10.0.0.1",
			expectBlock: true,
		},
		{
			name:        "Private IP 172.16.x.x",
			host:        "172.16.0.1",
			expectBlock: true,
		},
		{
			name:        "AWS metadata",
			host:        "169.254.169.254",
			expectBlock: true,
		},
		{
			name:        "IPv6 localhost",
			host:        "::1",
			expectBlock: true,
		},
	}

	suite := NewSecurityTestSuite()
	validator := app.NewSecurityValidator(suite.app)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testURL := fmt.Sprintf("http://%s/test", tt.host)
			err := validator.ValidateURL(testURL)

			if tt.expectBlock {
				if err == nil {
					t.Errorf("Expected URL %q to be blocked, but it was allowed", testURL)
				}
			} else {
				if err != nil {
					// For public hosts, we expect validation to pass
					// (though actual connection might fail)
					if apiErr, ok := err.(*app.APIError); ok {
						if strings.Contains(apiErr.Code, "SSRF") {
							t.Errorf("Public URL %q was incorrectly blocked: %v", testURL, err)
						}
					}
				}
			}
		})
	}
}

// BenchmarkSecurityValidation benchmarks security validation performance
func BenchmarkSecurityValidation(b *testing.B) {
	suite := NewSecurityTestSuite()
	validator := app.NewSecurityValidator(suite.app)

	benchmarks := []struct {
		name string
		fn   func()
	}{
		{
			name: "ValidateURL",
			fn: func() {
				validator.ValidateURL("https://api.example.com/openapi.yaml")
			},
		},
		{
			name: "ValidateFilePath",
			fn: func() {
				validator.ValidateFilePath("project/spec.yaml")
			},
		},
		{
			name: "SanitizeTemplateVariable",
			fn: func() {
				validator.SanitizeTemplateVariable("{{.Config.Database.Password}}")
			},
		},
		{
			name: "ValidateJSONContent",
			fn: func() {
				validator.ValidateJSONContent(`{"test": {"nested": {"deep": "value"}}}`)
			},
		},
	}

	for _, bm := range benchmarks {
		b.Run(bm.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				bm.fn()
			}
		})
	}
}

// TestSecurityIntegration tests security controls in integration scenarios
func TestSecurityIntegration(t *testing.T) {
	suite := NewSecurityTestSuite()

	// Create a realistic test scenario
	tempDir, err := os.MkdirTemp("", "mcpweaver-integration-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Test complete workflow with security controls
	t.Run("Secure file import workflow", func(t *testing.T) {
		// This would test a complete import workflow with all security controls active
		validator := app.NewSecurityValidator(suite.app)
		secureFS := app.NewSecureFileSystem(suite.app)

		// Set allowed paths
		err := secureFS.SetAllowedPaths([]string{tempDir})
		if err != nil {
			t.Fatalf("Failed to set allowed paths: %v", err)
		}

		// Test valid workflow
		validSpec := `openapi: 3.0.0
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

		// Validate YAML content
		err = validator.ValidateYAMLContent(validSpec)
		if err != nil {
			t.Errorf("Valid YAML content failed validation: %v", err)
		}

		// Write to secure location
		specPath := filepath.Join(tempDir, "test-spec.yaml")
		err = secureFS.SecureWriteFile(specPath, []byte(validSpec), 0644)
		if err != nil {
			t.Errorf("Failed to write valid spec: %v", err)
		}

		// Read back and verify
		content, err := secureFS.SecureReadFile(specPath)
		if err != nil {
			t.Errorf("Failed to read back spec: %v", err)
		}

		if string(content) != validSpec {
			t.Errorf("Content mismatch after secure read/write")
		}
	})

	t.Run("Block malicious workflow", func(t *testing.T) {
		validator := app.NewSecurityValidator(suite.app)
		secureFS := app.NewSecureFileSystem(suite.app)

		// Try to write to unauthorized location
		err := secureFS.SecureWriteFile("/tmp/evil.yaml", []byte("evil content"), 0644)
		if err == nil {
			t.Error("Expected security error when writing to unauthorized location")
		}

		// Try to validate malicious YAML
		maliciousYaml := `
test: !!python/object:os.system
args: ["rm -rf /"]
`
		err = validator.ValidateYAMLContent(maliciousYaml)
		if err == nil {
			t.Error("Expected security error for malicious YAML")
		}

		// Try path traversal
		err = validator.ValidateFilePath("../../../etc/passwd")
		if err == nil {
			t.Error("Expected security error for path traversal")
		}
	})
}

// Helper function to test if a network address is accessible
func isNetworkAccessible(address string) bool {
	conn, err := net.DialTimeout("tcp", address, 1*time.Second)
	if err != nil {
		return false
	}
	conn.Close()
	return true
}

// TestSecurityHeaders tests HTTP security headers
func TestSecurityHeaders(t *testing.T) {
	suite := NewSecurityTestSuite()
	validator := app.NewSecurityValidator(suite.app)

	tests := []struct {
		name        string
		headers     map[string]string
		expectError bool
		errorCode   string
	}{
		{
			name: "Safe headers",
			headers: map[string]string{
				"Content-Type": "application/json",
				"Accept":       "application/json",
				"User-Agent":   "MCPWeaver/1.0",
			},
			expectError: false,
		},
		{
			name: "Dangerous X-Forwarded-For",
			headers: map[string]string{
				"X-Forwarded-For": "127.0.0.1",
				"Content-Type":    "application/json",
			},
			expectError: true,
			errorCode:   "DANGEROUS_HEADER",
		},
		{
			name: "Dangerous X-Real-IP",
			headers: map[string]string{
				"X-Real-IP":    "192.168.1.1",
				"Content-Type": "application/json",
			},
			expectError: true,
			errorCode:   "DANGEROUS_HEADER",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validator.ValidateHTTPHeaders(tt.headers)

			if tt.expectError {
				if err == nil {
					t.Errorf("Expected error for headers %v, but got none", tt.headers)
					return
				}

				if tt.errorCode != "" {
					if apiErr, ok := err.(*app.APIError); ok {
						if apiErr.Code != tt.errorCode {
							t.Errorf("Expected error code %q, got %q", tt.errorCode, apiErr.Code)
						}
					}
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error for headers %v: %v", tt.headers, err)
				}
			}
		})
	}
}
