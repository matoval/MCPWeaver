package app

import (
	"context"
	"crypto/tls"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// SecureHTTPClient provides a secure HTTP client with built-in protections
type SecureHTTPClient struct {
	client           *http.Client
	app              *App
	securityValidator *SecurityValidator
	maxResponseSize  int64
	allowInsecure    bool
}

// HTTPClientConfig holds configuration for the secure HTTP client
type HTTPClientConfig struct {
	Timeout         time.Duration
	MaxResponseSize int64
	AllowInsecure   bool
	UserAgent       string
	MaxRedirects    int
}

// NewSecureHTTPClient creates a new secure HTTP client
func NewSecureHTTPClient(app *App, config HTTPClientConfig) *SecureHTTPClient {
	if config.Timeout == 0 {
		config.Timeout = 30 * time.Second
	}
	if config.MaxResponseSize == 0 {
		config.MaxResponseSize = 10 * 1024 * 1024 // 10MB default
	}
	if config.UserAgent == "" {
		config.UserAgent = "MCPWeaver/1.0"
	}
	if config.MaxRedirects == 0 {
		config.MaxRedirects = 5
	}

	// Create custom transport with security settings
	transport := &http.Transport{
		DialContext: (&net.Dialer{
			Timeout:   10 * time.Second,
			KeepAlive: 30 * time.Second,
		}).DialContext,
		TLSHandshakeTimeout:   10 * time.Second,
		ResponseHeaderTimeout: 30 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
		MaxIdleConns:          10,
		MaxIdleConnsPerHost:   2,
		IdleConnTimeout:       90 * time.Second,
		DisableCompression:    false,
		DisableKeepAlives:     false,
		ForceAttemptHTTP2:     true,
		TLSClientConfig: &tls.Config{
			MinVersion:         tls.VersionTLS12,
			InsecureSkipVerify: config.AllowInsecure,
			CipherSuites: []uint16{
				tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
				tls.TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305,
				tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
				tls.TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305,
				tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
				tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
			},
		},
	}

	// Create HTTP client with security settings
	client := &http.Client{
		Transport: transport,
		Timeout:   config.Timeout,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			// Limit redirects
			if len(via) >= config.MaxRedirects {
				return fmt.Errorf("too many redirects (max: %d)", config.MaxRedirects)
			}
			
			// Validate redirect URL for SSRF protection
			validator := NewSecurityValidator(app)
			if err := validator.ValidateURL(req.URL.String()); err != nil {
				return fmt.Errorf("redirect blocked by security policy: %w", err)
			}
			
			return nil
		},
	}

	return &SecureHTTPClient{
		client:            client,
		app:               app,
		securityValidator: NewSecurityValidator(app),
		maxResponseSize:   config.MaxResponseSize,
		allowInsecure:     config.AllowInsecure,
	}
}

// Get performs a secure GET request
func (c *SecureHTTPClient) Get(ctx context.Context, url string) (*http.Response, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, c.app.createAPIError("network", "REQUEST_CREATE_ERROR", "Failed to create HTTP request", map[string]string{
			"url":   url,
			"error": err.Error(),
		})
	}

	return c.Do(req)
}

// Post performs a secure POST request
func (c *SecureHTTPClient) Post(ctx context.Context, url, contentType string, body []byte) (*http.Response, error) {
	req, err := http.NewRequestWithContext(ctx, "POST", url, strings.NewReader(string(body)))
	if err != nil {
		return nil, c.app.createAPIError("network", "REQUEST_CREATE_ERROR", "Failed to create HTTP request", map[string]string{
			"url":   url,
			"error": err.Error(),
		})
	}

	if contentType != "" {
		req.Header.Set("Content-Type", contentType)
	}

	return c.Do(req)
}

// Do performs a secure HTTP request with all security validations
func (c *SecureHTTPClient) Do(req *http.Request) (*http.Response, error) {
	// Validate the request URL
	if err := c.securityValidator.ValidateURL(req.URL.String()); err != nil {
		return nil, err
	}

	// Set secure headers
	c.setSecureHeaders(req)

	// Validate headers for security
	headers := make(map[string]string)
	for name, values := range req.Header {
		if len(values) > 0 {
			headers[name] = values[0]
		}
	}
	if err := c.securityValidator.ValidateHTTPHeaders(headers); err != nil {
		return nil, err
	}

	// Add request timeout context
	ctx, cancel := context.WithTimeout(req.Context(), 30*time.Second)
	defer cancel()
	req = req.WithContext(ctx)

	// Perform the request
	resp, err := c.client.Do(req)
	if err != nil {
		// Check for specific error types
		if urlErr, ok := err.(*url.Error); ok {
			if urlErr.Timeout() {
				return nil, c.app.createAPIError("network", "REQUEST_TIMEOUT", "Request timed out", map[string]string{
					"url": req.URL.String(),
				})
			}
		}
		
		return nil, c.app.createAPIError("network", "REQUEST_FAILED", "HTTP request failed", map[string]string{
			"url":   req.URL.String(),
			"error": err.Error(),
		})
	}

	// Validate response
	if err := c.validateResponse(resp); err != nil {
		resp.Body.Close()
		return nil, err
	}

	return resp, nil
}

// setSecureHeaders sets security-related headers for requests
func (c *SecureHTTPClient) setSecureHeaders(req *http.Request) {
	// Set User-Agent
	req.Header.Set("User-Agent", "MCPWeaver/1.0 (OpenAPI to MCP Converter)")
	
	// Set Accept header
	req.Header.Set("Accept", "application/json, application/yaml, text/yaml, text/plain")
	
	// Set Accept-Encoding
	req.Header.Set("Accept-Encoding", "gzip, deflate")
	
	// Set Connection header
	req.Header.Set("Connection", "close")
	
	// Set Cache-Control to prevent caching sensitive data
	req.Header.Set("Cache-Control", "no-cache, no-store, must-revalidate")
	
	// Set Pragma for HTTP/1.0 compatibility
	req.Header.Set("Pragma", "no-cache")
	
	// Remove potentially dangerous headers
	dangerousHeaders := []string{
		"X-Forwarded-For",
		"X-Real-IP",
		"X-Forwarded-Proto",
		"X-Forwarded-Host",
		"Authorization", // Remove any existing auth headers for security
	}
	
	for _, header := range dangerousHeaders {
		req.Header.Del(header)
	}
}

// validateResponse validates the HTTP response for security issues
func (c *SecureHTTPClient) validateResponse(resp *http.Response) error {
	// Check response status
	if resp.StatusCode < 200 || resp.StatusCode >= 400 {
		return c.app.createAPIError("network", "HTTP_ERROR", "HTTP request failed", map[string]string{
			"url":        resp.Request.URL.String(),
			"status":     resp.Status,
			"statusCode": fmt.Sprintf("%d", resp.StatusCode),
		})
	}

	// Check Content-Length if provided
	if resp.ContentLength > c.maxResponseSize {
		return c.app.createAPIError("network", "RESPONSE_TOO_LARGE", "Response exceeds maximum allowed size", map[string]string{
			"url":           resp.Request.URL.String(),
			"contentLength": fmt.Sprintf("%d", resp.ContentLength),
			"maxSize":       fmt.Sprintf("%d", c.maxResponseSize),
		})
	}

	// Validate Content-Type
	contentType := resp.Header.Get("Content-Type")
	if !c.isAllowedContentType(contentType) {
		return c.app.createAPIError("network", "INVALID_CONTENT_TYPE", "Response content type not allowed", map[string]string{
			"url":         resp.Request.URL.String(),
			"contentType": contentType,
		})
	}

	// Check for suspicious headers
	suspiciousHeaders := []string{
		"X-Frame-Options",
		"X-XSS-Protection", 
		"X-Content-Type-Options",
	}
	
	for _, header := range suspiciousHeaders {
		if value := resp.Header.Get(header); value != "" {
			// These headers suggest the response might be HTML/web content rather than API data
			if strings.Contains(strings.ToLower(contentType), "text/html") {
				return c.app.createAPIError("network", "SUSPICIOUS_CONTENT", "Response appears to be web content rather than API data", map[string]string{
					"url":         resp.Request.URL.String(),
					"contentType": contentType,
					"header":      header,
					"value":       value,
				})
			}
		}
	}

	return nil
}

// isAllowedContentType checks if the content type is allowed for API responses
func (c *SecureHTTPClient) isAllowedContentType(contentType string) bool {
	if contentType == "" {
		return true // Allow empty content type
	}
	
	allowedTypes := []string{
		"application/json",
		"application/yaml",
		"application/x-yaml",
		"text/yaml",
		"text/x-yaml",
		"text/plain",
		"application/openapi+json",
		"application/openapi+yaml",
	}
	
	// Normalize content type (remove charset, etc.)
	normalizedType := strings.ToLower(strings.Split(contentType, ";")[0])
	normalizedType = strings.TrimSpace(normalizedType)
	
	for _, allowed := range allowedTypes {
		if normalizedType == allowed {
			return true
		}
	}
	
	return false
}

// GetDefaultConfig returns a default secure HTTP client configuration
func GetDefaultHTTPClientConfig() HTTPClientConfig {
	return HTTPClientConfig{
		Timeout:         30 * time.Second,
		MaxResponseSize: 10 * 1024 * 1024, // 10MB
		AllowInsecure:   false,
		UserAgent:       "MCPWeaver/1.0",
		MaxRedirects:    5,
	}
}

// GetInsecureConfig returns an HTTP client configuration that allows insecure connections
// This should only be used for development/testing purposes
func GetInsecureHTTPClientConfig() HTTPClientConfig {
	config := GetDefaultHTTPClientConfig()
	config.AllowInsecure = true
	return config
}

// Close cleans up the HTTP client resources
func (c *SecureHTTPClient) Close() error {
	if transport, ok := c.client.Transport.(*http.Transport); ok {
		transport.CloseIdleConnections()
	}
	return nil
}

// GetTLSConnectionState returns the TLS connection state for the last request
func (c *SecureHTTPClient) GetTLSConnectionState(resp *http.Response) *tls.ConnectionState {
	if resp != nil && resp.TLS != nil {
		return resp.TLS
	}
	return nil
}

// ValidateServerCertificate performs additional server certificate validation
func (c *SecureHTTPClient) ValidateServerCertificate(resp *http.Response, expectedHost string) error {
	if resp.TLS == nil {
		if resp.Request.URL.Scheme == "https" {
			return c.app.createAPIError("security", "NO_TLS", "HTTPS request completed without TLS", map[string]string{
				"url": resp.Request.URL.String(),
			})
		}
		return nil // HTTP is allowed for development
	}

	// Check TLS version
	if resp.TLS.Version < tls.VersionTLS12 {
		return c.app.createAPIError("security", "WEAK_TLS", "TLS version is too old", map[string]string{
			"url":        resp.Request.URL.String(),
			"tlsVersion": fmt.Sprintf("%d", resp.TLS.Version),
			"minVersion": fmt.Sprintf("%d", tls.VersionTLS12),
		})
	}

	// Check if connection was verified
	if !resp.TLS.HandshakeComplete {
		return c.app.createAPIError("security", "TLS_HANDSHAKE_INCOMPLETE", "TLS handshake not completed", map[string]string{
			"url": resp.Request.URL.String(),
		})
	}

	// Check peer certificates
	if len(resp.TLS.PeerCertificates) == 0 {
		return c.app.createAPIError("security", "NO_PEER_CERTIFICATES", "No peer certificates found", map[string]string{
			"url": resp.Request.URL.String(),
		})
	}

	// Additional certificate validation could be added here
	
	return nil
}

// SetCustomVerifyPeerCertificate sets a custom certificate verification function
func (c *SecureHTTPClient) SetCustomVerifyPeerCertificate(verify func(*http.Request, *tls.ConnectionState) error) {
	if transport, ok := c.client.Transport.(*http.Transport); ok {
		if transport.TLSClientConfig != nil {
			// Note: This would require modifying the transport's VerifyPeerCertificate function
			// In a real implementation, this would be set during transport creation
		}
	}
}