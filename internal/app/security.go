package app

import (
	"fmt"
	"net"
	"net/url"
	"path/filepath"
	"regexp"
	"strings"
	"unicode"
)

// Security constants and configuration
const (
	// Maximum allowed sizes
	MaxURLLength         = 2048
	MaxFileNameLength    = 255
	MaxContentLength     = 10 * 1024 * 1024 // 10MB
	MaxRegexLength       = 1000
	MaxTemplateVarLength = 500

	// Rate limiting
	MaxURLRequestsPerMinute = 10

	// Allowed patterns
	AllowedTemplateNamePattern = `^[a-zA-Z0-9_-]+$`
	AllowedFileNamePattern     = `^[a-zA-Z0-9._-]+$`
	AllowedProjectNamePattern  = `^[a-zA-Z0-9\s._-]+$`
)

var (
	// Compiled regexes for validation
	templateNameRegex = regexp.MustCompile(AllowedTemplateNamePattern)
	fileNameRegex     = regexp.MustCompile(AllowedFileNamePattern)
	projectNameRegex  = regexp.MustCompile(AllowedProjectNamePattern)

	// Dangerous patterns to block
	dangerousRegexPatterns = []string{
		`(?=.*(?=.*(?=.*`, // Nested lookaheads (ReDoS risk)
		`(.*){10,}`,       // High repetition (ReDoS risk)
		`(.*)*`,           // Nested quantifiers (ReDoS risk)
		`(.+)+`,           // Nested quantifiers (ReDoS risk)
		`(.*)*$`,          // Catastrophic backtracking
	}

	// Private IP ranges to block in SSRF protection
	privateIPRanges = []string{
		"127.0.0.0/8",    // Loopback
		"10.0.0.0/8",     // Private Class A
		"172.16.0.0/12",  // Private Class B
		"192.168.0.0/16", // Private Class C
		"169.254.0.0/16", // Link-local
		"224.0.0.0/4",    // Multicast
		"240.0.0.0/4",    // Reserved
		"::1/128",        // IPv6 loopback
		"fc00::/7",       // IPv6 private
		"fe80::/10",      // IPv6 link-local
	}

	// Compiled private IP networks
	privateNetworks []*net.IPNet
)

// SecurityValidator provides input validation and sanitization
type SecurityValidator struct {
	app *App
}

// NewSecurityValidator creates a new security validator
func NewSecurityValidator(app *App) *SecurityValidator {
	return &SecurityValidator{app: app}
}

func init() {
	// Pre-compile private IP networks for SSRF protection
	for _, cidr := range privateIPRanges {
		_, network, err := net.ParseCIDR(cidr)
		if err == nil {
			privateNetworks = append(privateNetworks, network)
		}
	}
}

// ValidateURL validates and sanitizes URLs for SSRF protection
func (sv *SecurityValidator) ValidateURL(rawURL string) error {
	if rawURL == "" {
		return sv.app.createAPIError("validation", "INVALID_URL", "URL cannot be empty", nil)
	}

	if len(rawURL) > MaxURLLength {
		return sv.app.createAPIError("validation", "URL_TOO_LONG", "URL exceeds maximum length", map[string]string{
			"maxLength": fmt.Sprintf("%d", MaxURLLength),
			"length":    fmt.Sprintf("%d", len(rawURL)),
		})
	}

	// Parse URL
	parsedURL, err := url.Parse(rawURL)
	if err != nil {
		return sv.app.createAPIError("validation", "INVALID_URL_FORMAT", "Invalid URL format", map[string]string{
			"error": err.Error(),
		})
	}

	// Only allow HTTP and HTTPS schemes
	if parsedURL.Scheme != "http" && parsedURL.Scheme != "https" {
		return sv.app.createAPIError("security", "INVALID_URL_SCHEME", "Only HTTP and HTTPS URLs are allowed", map[string]string{
			"scheme": parsedURL.Scheme,
		})
	}

	// Check for SSRF vulnerabilities
	if err := sv.checkSSRF(parsedURL); err != nil {
		return err
	}

	return nil
}

// checkSSRF performs Server-Side Request Forgery protection
func (sv *SecurityValidator) checkSSRF(parsedURL *url.URL) error {
	host := parsedURL.Hostname()
	if host == "" {
		return sv.app.createAPIError("security", "INVALID_HOST", "URL must contain a valid hostname", nil)
	}

	// Block localhost variations
	localhostPatterns := []string{
		"localhost",
		"127.0.0.1",
		"0.0.0.0",
		"[::1]",
		"::1",
	}

	for _, pattern := range localhostPatterns {
		if strings.EqualFold(host, pattern) {
			return sv.app.createAPIError("security", "SSRF_LOCALHOST_BLOCKED", "Access to localhost is not allowed", map[string]string{
				"host": host,
			})
		}
	}

	// Parse IP address if possible
	ip := net.ParseIP(host)
	if ip != nil {
		// Check against private IP ranges
		for _, network := range privateNetworks {
			if network.Contains(ip) {
				return sv.app.createAPIError("security", "SSRF_PRIVATE_IP_BLOCKED", "Access to private IP addresses is not allowed", map[string]string{
					"ip":      ip.String(),
					"network": network.String(),
				})
			}
		}
	}

	// Block common metadata endpoints
	metadataHosts := []string{
		"169.254.169.254", // AWS, GCP metadata
		"metadata.google.internal",
		"metadata.azure.com",
	}

	for _, metadataHost := range metadataHosts {
		if strings.EqualFold(host, metadataHost) {
			return sv.app.createAPIError("security", "SSRF_METADATA_BLOCKED", "Access to cloud metadata endpoints is not allowed", map[string]string{
				"host": host,
			})
		}
	}

	return nil
}

// ValidateFilePath validates file paths for path traversal protection
func (sv *SecurityValidator) ValidateFilePath(filePath string) error {
	if filePath == "" {
		return sv.app.createAPIError("validation", "EMPTY_FILE_PATH", "File path cannot be empty", nil)
	}

	// Check for path traversal attempts
	if strings.Contains(filePath, "..") {
		return sv.app.createAPIError("security", "PATH_TRAVERSAL_DETECTED", "Path traversal attempts are not allowed", map[string]string{
			"path": filePath,
		})
	}

	// Clean the path
	cleanPath := filepath.Clean(filePath)
	if cleanPath != filePath {
		return sv.app.createAPIError("security", "INVALID_PATH_FORMAT", "Invalid path format detected", map[string]string{
			"original": filePath,
			"cleaned":  cleanPath,
		})
	}

	// Check file name length
	fileName := filepath.Base(filePath)
	if len(fileName) > MaxFileNameLength {
		return sv.app.createAPIError("validation", "FILENAME_TOO_LONG", "File name exceeds maximum length", map[string]string{
			"fileName":  fileName,
			"maxLength": fmt.Sprintf("%d", MaxFileNameLength),
		})
	}

	// Validate file name characters
	if !fileNameRegex.MatchString(fileName) {
		return sv.app.createAPIError("validation", "INVALID_FILENAME", "File name contains invalid characters", map[string]string{
			"fileName": fileName,
			"pattern":  AllowedFileNamePattern,
		})
	}

	return nil
}

// ValidateTemplateName validates template names with whitelist approach
func (sv *SecurityValidator) ValidateTemplateName(templateName string) error {
	if templateName == "" {
		return sv.app.createAPIError("validation", "EMPTY_TEMPLATE_NAME", "Template name cannot be empty", nil)
	}

	if len(templateName) > 100 {
		return sv.app.createAPIError("validation", "TEMPLATE_NAME_TOO_LONG", "Template name is too long", map[string]string{
			"name":      templateName,
			"maxLength": "100",
		})
	}

	// Use whitelist approach for template names
	if !templateNameRegex.MatchString(templateName) {
		return sv.app.createAPIError("validation", "INVALID_TEMPLATE_NAME", "Template name contains invalid characters", map[string]string{
			"name":    templateName,
			"pattern": AllowedTemplateNamePattern,
		})
	}

	// Block reserved names
	reservedNames := []string{
		"con", "prn", "aux", "nul", // Windows reserved
		"com1", "com2", "com3", "com4", "com5", "com6", "com7", "com8", "com9",
		"lpt1", "lpt2", "lpt3", "lpt4", "lpt5", "lpt6", "lpt7", "lpt8", "lpt9",
		"admin", "root", "system", "config", // Common reserved
	}

	lowerName := strings.ToLower(templateName)
	for _, reserved := range reservedNames {
		if lowerName == reserved {
			return sv.app.createAPIError("validation", "RESERVED_TEMPLATE_NAME", "Template name is reserved", map[string]string{
				"name": templateName,
			})
		}
	}

	return nil
}

// ValidateProjectName validates project names
func (sv *SecurityValidator) ValidateProjectName(projectName string) error {
	if projectName == "" {
		return sv.app.createAPIError("validation", "EMPTY_PROJECT_NAME", "Project name cannot be empty", nil)
	}

	if len(projectName) > 100 {
		return sv.app.createAPIError("validation", "PROJECT_NAME_TOO_LONG", "Project name is too long", map[string]string{
			"name":      projectName,
			"maxLength": "100",
		})
	}

	if !projectNameRegex.MatchString(projectName) {
		return sv.app.createAPIError("validation", "INVALID_PROJECT_NAME", "Project name contains invalid characters", map[string]string{
			"name":    projectName,
			"pattern": AllowedProjectNamePattern,
		})
	}

	return nil
}

// ValidateRegexPattern validates user-provided regex patterns to prevent ReDoS
func (sv *SecurityValidator) ValidateRegexPattern(pattern string) error {
	if pattern == "" {
		return nil // Empty patterns are allowed
	}

	if len(pattern) > MaxRegexLength {
		return sv.app.createAPIError("validation", "REGEX_TOO_LONG", "Regex pattern is too long", map[string]string{
			"pattern":   pattern,
			"maxLength": fmt.Sprintf("%d", MaxRegexLength),
		})
	}

	// Check for dangerous patterns that could cause ReDoS
	for _, dangerous := range dangerousRegexPatterns {
		if strings.Contains(pattern, dangerous) {
			return sv.app.createAPIError("security", "DANGEROUS_REGEX", "Regex pattern contains potentially dangerous constructs", map[string]string{
				"pattern":   pattern,
				"dangerous": dangerous,
			})
		}
	}

	// Try to compile the regex to ensure it's valid
	_, err := regexp.Compile(pattern)
	if err != nil {
		return sv.app.createAPIError("validation", "INVALID_REGEX", "Invalid regex pattern", map[string]string{
			"pattern": pattern,
			"error":   err.Error(),
		})
	}

	return nil
}

// SanitizeTemplateVariable sanitizes variables for template injection protection
func (sv *SecurityValidator) SanitizeTemplateVariable(variable string) string {
	if len(variable) > MaxTemplateVarLength {
		variable = variable[:MaxTemplateVarLength]
	}

	// Remove control characters except newline, carriage return, and tab
	sanitized := strings.Map(func(r rune) rune {
		if unicode.IsControl(r) && r != '\n' && r != '\r' && r != '\t' {
			return -1
		}
		return r
	}, variable)

	// Escape template-specific characters to prevent injection
	sanitized = strings.ReplaceAll(sanitized, "{{", "&#123;&#123;")
	sanitized = strings.ReplaceAll(sanitized, "}}", "&#125;&#125;")
	sanitized = strings.ReplaceAll(sanitized, "<%", "&#60;&#37;")
	sanitized = strings.ReplaceAll(sanitized, "%>", "&#37;&#62;")

	return sanitized
}

// SanitizeInput performs general input sanitization
func (sv *SecurityValidator) SanitizeInput(input string) string {
	// Remove null bytes
	input = strings.ReplaceAll(input, "\x00", "")

	// Remove or replace control characters (except common whitespace)
	sanitized := strings.Map(func(r rune) rune {
		if unicode.IsControl(r) && r != '\n' && r != '\r' && r != '\t' {
			return -1
		}
		return r
	}, input)

	// Trim whitespace
	sanitized = strings.TrimSpace(sanitized)

	return sanitized
}

// ValidateContentLength validates content length
func (sv *SecurityValidator) ValidateContentLength(content string, maxLength int) error {
	if maxLength <= 0 {
		maxLength = MaxContentLength
	}

	if len(content) > maxLength {
		return sv.app.createAPIError("validation", "CONTENT_TOO_LARGE", "Content exceeds maximum allowed length", map[string]string{
			"length":    fmt.Sprintf("%d", len(content)),
			"maxLength": fmt.Sprintf("%d", maxLength),
		})
	}

	return nil
}

// IsSecureContext checks if the application is running in a secure context
func (sv *SecurityValidator) IsSecureContext() bool {
	// In a desktop application, we consider the local context secure
	// This could be extended to check for debugger attachment, etc.
	return true
}

// ValidateJSONContent validates JSON content for security issues
func (sv *SecurityValidator) ValidateJSONContent(content string) error {
	// Check for excessively deep nesting (JSON bomb protection)
	maxDepth := 50
	depth := 0
	maxDepthReached := 0

	for _, char := range content {
		switch char {
		case '{', '[':
			depth++
			if depth > maxDepthReached {
				maxDepthReached = depth
			}
			if depth > maxDepth {
				return sv.app.createAPIError("security", "JSON_TOO_DEEP", "JSON content has excessive nesting depth", map[string]string{
					"maxDepth":      fmt.Sprintf("%d", maxDepth),
					"detectedDepth": fmt.Sprintf("%d", depth),
				})
			}
		case '}', ']':
			depth--
		}
	}

	return nil
}

// ValidateYAMLContent validates YAML content for security issues
func (sv *SecurityValidator) ValidateYAMLContent(content string) error {
	// Check for YAML bombs (documents with excessive references)
	if strings.Count(content, "&") > 100 {
		return sv.app.createAPIError("security", "YAML_EXCESSIVE_REFERENCES", "YAML content has too many references", nil)
	}

	// Check for potentially dangerous YAML constructs
	dangerousPatterns := []string{
		"!!python/", // Python object serialization
		"!!java/",   // Java object serialization
		"!!js/",     // JavaScript object serialization
		"!!ruby/",   // Ruby object serialization
		"!<tag:",    // Custom tag URIs
	}

	lowerContent := strings.ToLower(content)
	for _, pattern := range dangerousPatterns {
		if strings.Contains(lowerContent, pattern) {
			return sv.app.createAPIError("security", "YAML_DANGEROUS_CONSTRUCT", "YAML content contains potentially dangerous constructs", map[string]string{
				"pattern": pattern,
			})
		}
	}

	return nil
}

// ValidateHTTPHeaders validates HTTP headers for security
func (sv *SecurityValidator) ValidateHTTPHeaders(headers map[string]string) error {
	// Check for dangerous headers
	dangerousHeaders := []string{
		"x-forwarded-for",
		"x-real-ip",
		"x-forwarded-proto",
		"x-forwarded-host",
	}

	for headerName := range headers {
		lowerName := strings.ToLower(headerName)
		for _, dangerous := range dangerousHeaders {
			if lowerName == dangerous {
				return sv.app.createAPIError("security", "DANGEROUS_HEADER", "HTTP header not allowed", map[string]string{
					"header": headerName,
				})
			}
		}
	}

	return nil
}
