package testing

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os/exec"
	"strings"
	"sync"
	"time"
)

// ProtocolTester handles MCP protocol compliance testing
type ProtocolTester struct {
	config *TestConfig
	mutex  sync.RWMutex
}

// NewProtocolTester creates a new protocol tester
func NewProtocolTester(config *TestConfig) *ProtocolTester {
	return &ProtocolTester{
		config: config,
	}
}

// TestCompliance tests MCP protocol compliance
func (pt *ProtocolTester) TestCompliance(ctx context.Context, serverPath string) (*ProtocolTestResult, error) {
	startTime := time.Now()
	
	result := &ProtocolTestResult{
		Success:               true,
		ProtocolVersion:       pt.config.MCPProtocolVersion,
		SupportedMethods:      make([]string, 0),
		SupportedCapabilities: make([]string, 0),
		MethodTests:           make(map[string]*MethodTest),
		CapabilityTests:       make(map[string]*CapabilityTest),
		Errors:                make([]string, 0),
	}

	// Start the MCP server
	serverProcess, stdin, stdout, stderr, err := pt.startServer(ctx, serverPath)
	if err != nil {
		result.Success = false
		result.Errors = append(result.Errors, fmt.Sprintf("Failed to start server: %v", err))
		result.Duration = time.Since(startTime)
		return result, err
	}
	defer pt.stopServer(serverProcess)

	// Create JSON-RPC client
	client := NewJSONRPCClient(stdin, stdout, stderr)
	defer client.Close()

	// Test initialization
	if err := pt.testInitialization(ctx, client, result); err != nil {
		result.Success = false
		result.Errors = append(result.Errors, fmt.Sprintf("Initialization test failed: %v", err))
	}

	// Test required methods
	for _, method := range pt.config.RequiredMethods {
		if err := pt.testMethod(ctx, client, method, result); err != nil {
			result.Errors = append(result.Errors, fmt.Sprintf("Method test failed for %s: %v", method, err))
		}
	}

	// Test capabilities
	for _, capability := range pt.config.RequiredCapabilities {
		if err := pt.testCapability(ctx, client, capability, result); err != nil {
			result.Errors = append(result.Errors, fmt.Sprintf("Capability test failed for %s: %v", capability, err))
		}
	}

	// Test error handling
	if err := pt.testErrorHandling(ctx, client, result); err != nil {
		result.Errors = append(result.Errors, fmt.Sprintf("Error handling test failed: %v", err))
	}

	result.Duration = time.Since(startTime)
	result.Success = len(result.Errors) == 0

	return result, nil
}

// startServer starts the MCP server process
func (pt *ProtocolTester) startServer(ctx context.Context, serverPath string) (*exec.Cmd, io.WriteCloser, io.ReadCloser, io.ReadCloser, error) {
	// First compile the server
	compileCmd := exec.CommandContext(ctx, "go", "build", "-o", "test-mcp-server", "main.go")
	compileCmd.Dir = serverPath
	if err := compileCmd.Run(); err != nil {
		return nil, nil, nil, nil, fmt.Errorf("failed to compile server: %w", err)
	}

	// Start the server
	serverBinary := fmt.Sprintf("%s/test-mcp-server", serverPath)
	cmd := exec.CommandContext(ctx, serverBinary)
	cmd.Dir = serverPath

	stdin, err := cmd.StdinPipe()
	if err != nil {
		return nil, nil, nil, nil, fmt.Errorf("failed to create stdin pipe: %w", err)
	}

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		stdin.Close()
		return nil, nil, nil, nil, fmt.Errorf("failed to create stdout pipe: %w", err)
	}

	stderr, err := cmd.StderrPipe()
	if err != nil {
		stdin.Close()
		stdout.Close()
		return nil, nil, nil, nil, fmt.Errorf("failed to create stderr pipe: %w", err)
	}

	if err := cmd.Start(); err != nil {
		stdin.Close()
		stdout.Close()
		stderr.Close()
		return nil, nil, nil, nil, fmt.Errorf("failed to start server: %w", err)
	}

	// Give the server a moment to start
	time.Sleep(100 * time.Millisecond)

	return cmd, stdin, stdout, stderr, nil
}

// stopServer stops the MCP server process
func (pt *ProtocolTester) stopServer(cmd *exec.Cmd) {
	if cmd != nil && cmd.Process != nil {
		cmd.Process.Kill()
		cmd.Wait()
	}
}

// testInitialization tests the MCP initialization process
func (pt *ProtocolTester) testInitialization(ctx context.Context, client *JSONRPCClient, result *ProtocolTestResult) error {
	startTime := time.Now()

	// Send initialize request
	initRequest := map[string]interface{}{
		"protocolVersion": pt.config.MCPProtocolVersion,
		"capabilities": map[string]interface{}{
			"roots": map[string]interface{}{
				"listChanged": true,
			},
		},
		"clientInfo": map[string]interface{}{
			"name":    "MCPWeaver Test Client",
			"version": "1.0.0",
		},
	}

	response, err := client.Call(ctx, "initialize", initRequest)
	if err != nil {
		return fmt.Errorf("initialize request failed: %w", err)
	}

	methodTest := &MethodTest{
		Method:       "initialize",
		Success:      true,
		ResponseTime: time.Since(startTime),
		Request:      initRequest,
		Response:     response,
	}

	// Validate response structure
	if err := pt.validateInitializeResponse(response, result); err != nil {
		methodTest.Success = false
		methodTest.ErrorMessage = err.Error()
		return err
	}

	result.MethodTests["initialize"] = methodTest
	result.SupportedMethods = append(result.SupportedMethods, "initialize")

	return nil
}

// validateInitializeResponse validates the initialize response
func (pt *ProtocolTester) validateInitializeResponse(response interface{}, result *ProtocolTestResult) error {
	responseMap, ok := response.(map[string]interface{})
	if !ok {
		return fmt.Errorf("initialize response is not a valid object")
	}

	// Check protocol version
	if protocolVersion, exists := responseMap["protocolVersion"]; exists {
		if versionStr, ok := protocolVersion.(string); ok {
			result.ProtocolVersion = versionStr
		}
	}

	// Check capabilities
	if capabilities, exists := responseMap["capabilities"]; exists {
		if capMap, ok := capabilities.(map[string]interface{}); ok {
			for capability := range capMap {
				result.SupportedCapabilities = append(result.SupportedCapabilities, capability)
			}
		}
	}

	// Check server info
	if serverInfo, exists := responseMap["serverInfo"]; exists {
		if _, ok := serverInfo.(map[string]interface{}); !ok {
			return fmt.Errorf("serverInfo is not a valid object")
		}
	} else {
		return fmt.Errorf("missing serverInfo in initialize response")
	}

	return nil
}

// testMethod tests a specific MCP method
func (pt *ProtocolTester) testMethod(ctx context.Context, client *JSONRPCClient, method string, result *ProtocolTestResult) error {
	startTime := time.Now()

	var request interface{}
	var expectedFields []string

	// Prepare method-specific requests
	switch method {
	case "tools/list":
		request = map[string]interface{}{}
		expectedFields = []string{"tools"}
	case "tools/call":
		// We'll test this with a dummy call first
		request = map[string]interface{}{
			"name":      "dummy_tool",
			"arguments": map[string]interface{}{},
		}
		expectedFields = []string{"content"}
	default:
		// Generic method test
		request = map[string]interface{}{}
	}

	response, err := client.Call(ctx, method, request)
	
	methodTest := &MethodTest{
		Method:       method,
		Success:      err == nil,
		ResponseTime: time.Since(startTime),
		Request:      request,
		Response:     response,
	}

	if err != nil {
		// For tools/call with dummy tool, we expect an error
		if method == "tools/call" && strings.Contains(err.Error(), "Tool not found") {
			methodTest.Success = true
			methodTest.ErrorMessage = "Expected error for non-existent tool"
		} else {
			methodTest.ErrorMessage = err.Error()
		}
	} else {
		// Validate response structure
		if err := pt.validateMethodResponse(method, response, expectedFields); err != nil {
			methodTest.Success = false
			methodTest.ErrorMessage = err.Error()
		}
	}

	result.MethodTests[method] = methodTest

	if methodTest.Success {
		result.SupportedMethods = append(result.SupportedMethods, method)
	}

	return nil
}

// validateMethodResponse validates method response structure
func (pt *ProtocolTester) validateMethodResponse(method string, response interface{}, expectedFields []string) error {
	responseMap, ok := response.(map[string]interface{})
	if !ok {
		return fmt.Errorf("response is not a valid object")
	}

	// Check for expected fields
	for _, field := range expectedFields {
		if _, exists := responseMap[field]; !exists {
			return fmt.Errorf("missing required field: %s", field)
		}
	}

	// Method-specific validation
	switch method {
	case "tools/list":
		if tools, exists := responseMap["tools"]; exists {
			if toolsArray, ok := tools.([]interface{}); ok {
				// Validate tool structure
				for i, tool := range toolsArray {
					if err := pt.validateToolStructure(tool, i); err != nil {
						return err
					}
				}
			} else {
				return fmt.Errorf("tools field is not an array")
			}
		}
	case "tools/call":
		if content, exists := responseMap["content"]; exists {
			if contentArray, ok := content.([]interface{}); ok {
				// Validate content structure
				for i, item := range contentArray {
					if err := pt.validateContentStructure(item, i); err != nil {
						return err
					}
				}
			} else {
				return fmt.Errorf("content field is not an array")
			}
		}
	}

	return nil
}

// validateToolStructure validates individual tool structure
func (pt *ProtocolTester) validateToolStructure(tool interface{}, index int) error {
	toolMap, ok := tool.(map[string]interface{})
	if !ok {
		return fmt.Errorf("tool at index %d is not a valid object", index)
	}

	requiredFields := []string{"name", "description", "inputSchema"}
	for _, field := range requiredFields {
		if _, exists := toolMap[field]; !exists {
			return fmt.Errorf("tool at index %d missing required field: %s", index, field)
		}
	}

	// Validate inputSchema structure
	if inputSchema, exists := toolMap["inputSchema"]; exists {
		if schemaMap, ok := inputSchema.(map[string]interface{}); ok {
			if schemaType, exists := schemaMap["type"]; !exists || schemaType != "object" {
				return fmt.Errorf("tool at index %d has invalid inputSchema type", index)
			}
		} else {
			return fmt.Errorf("tool at index %d has invalid inputSchema structure", index)
		}
	}

	return nil
}

// validateContentStructure validates content item structure
func (pt *ProtocolTester) validateContentStructure(content interface{}, index int) error {
	contentMap, ok := content.(map[string]interface{})
	if !ok {
		return fmt.Errorf("content at index %d is not a valid object", index)
	}

	requiredFields := []string{"type", "text"}
	for _, field := range requiredFields {
		if _, exists := contentMap[field]; !exists {
			return fmt.Errorf("content at index %d missing required field: %s", index, field)
		}
	}

	return nil
}

// testCapability tests a specific MCP capability
func (pt *ProtocolTester) testCapability(ctx context.Context, client *JSONRPCClient, capability string, result *ProtocolTestResult) error {
	capabilityTest := &CapabilityTest{
		Capability: capability,
		Success:    false,
		Supported:  false,
	}

	// Check if capability was reported during initialization
	for _, supportedCap := range result.SupportedCapabilities {
		if supportedCap == capability {
			capabilityTest.Supported = true
			break
		}
	}

	// Test capability-specific functionality
	switch capability {
	case "tools":
		// Test tools capability by calling tools/list
		response, err := client.Call(ctx, "tools/list", map[string]interface{}{})
		if err == nil {
			capabilityTest.Success = true
			capabilityTest.TestDetails = map[string]interface{}{
				"toolsListResponse": response,
			}
		} else {
			capabilityTest.ErrorMessage = err.Error()
		}
	default:
		// For unknown capabilities, just check if they were declared
		capabilityTest.Success = capabilityTest.Supported
	}

	result.CapabilityTests[capability] = capabilityTest

	return nil
}

// testErrorHandling tests server error handling
func (pt *ProtocolTester) testErrorHandling(ctx context.Context, client *JSONRPCClient, result *ProtocolTestResult) error {
	// Test invalid method
	_, err := client.Call(ctx, "invalid/method", map[string]interface{}{})
	if err == nil {
		return fmt.Errorf("server should return error for invalid method")
	}

	// Test invalid parameters
	_, err = client.Call(ctx, "tools/call", "invalid_params")
	if err == nil {
		return fmt.Errorf("server should return error for invalid parameters")
	}

	// Add error handling test to results
	result.MethodTests["error_handling"] = &MethodTest{
		Method:  "error_handling",
		Success: true,
		Request: map[string]interface{}{"test": "error_handling"},
	}

	return nil
}

// JSONRPCClient handles JSON-RPC communication with the MCP server
type JSONRPCClient struct {
	stdin   io.WriteCloser
	stdout  io.ReadCloser
	stderr  io.ReadCloser
	reqID   int
	mutex   sync.Mutex
	scanner *bufio.Scanner
}

// NewJSONRPCClient creates a new JSON-RPC client
func NewJSONRPCClient(stdin io.WriteCloser, stdout io.ReadCloser, stderr io.ReadCloser) *JSONRPCClient {
	return &JSONRPCClient{
		stdin:   stdin,
		stdout:  stdout,
		stderr:  stderr,
		reqID:   1,
		scanner: bufio.NewScanner(stdout),
	}
}

// Call makes a JSON-RPC call to the server
func (c *JSONRPCClient) Call(ctx context.Context, method string, params interface{}) (interface{}, error) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	// Prepare request
	request := map[string]interface{}{
		"jsonrpc": "2.0",
		"id":      c.reqID,
		"method":  method,
		"params":  params,
	}
	c.reqID++

	// Send request
	requestBytes, err := json.Marshal(request)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	requestBytes = append(requestBytes, '\n')
	if _, err := c.stdin.Write(requestBytes); err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}

	// Read response with timeout
	responseChan := make(chan []byte, 1)
	errorChan := make(chan error, 1)

	go func() {
		if c.scanner.Scan() {
			responseChan <- c.scanner.Bytes()
		} else {
			if err := c.scanner.Err(); err != nil {
				errorChan <- err
			} else {
				errorChan <- fmt.Errorf("no response received")
			}
		}
	}()

	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case responseBytes := <-responseChan:
		// Parse response
		var response map[string]interface{}
		if err := json.Unmarshal(responseBytes, &response); err != nil {
			return nil, fmt.Errorf("failed to parse response: %w", err)
		}

		// Check for error
		if errorObj, exists := response["error"]; exists {
			if errorMap, ok := errorObj.(map[string]interface{}); ok {
				if message, exists := errorMap["message"]; exists {
					return nil, fmt.Errorf("JSON-RPC error: %v", message)
				}
			}
			return nil, fmt.Errorf("JSON-RPC error: %v", errorObj)
		}

		// Return result
		if result, exists := response["result"]; exists {
			return result, nil
		}

		return nil, fmt.Errorf("no result in response")
	case err := <-errorChan:
		return nil, err
	case <-time.After(10 * time.Second):
		return nil, fmt.Errorf("timeout waiting for response")
	}
}

// Close closes the JSON-RPC client
func (c *JSONRPCClient) Close() error {
	if c.stdin != nil {
		c.stdin.Close()
	}
	if c.stdout != nil {
		c.stdout.Close()
	}
	if c.stderr != nil {
		c.stderr.Close()
	}
	return nil
}

// ProtocolConformanceValidator validates against MCP specification
type ProtocolConformanceValidator struct {
	config *TestConfig
}

// NewProtocolConformanceValidator creates a protocol conformance validator
func NewProtocolConformanceValidator(config *TestConfig) *ProtocolConformanceValidator {
	return &ProtocolConformanceValidator{config: config}
}

// ValidateConformance validates full protocol conformance
func (v *ProtocolConformanceValidator) ValidateConformance(ctx context.Context, serverPath string) (*ProtocolTestResult, error) {
	tester := NewProtocolTester(v.config)
	return tester.TestCompliance(ctx, serverPath)
}

// ValidateMessageFormats validates JSON-RPC message formats
func (v *ProtocolConformanceValidator) ValidateMessageFormats(messages []map[string]interface{}) error {
	for i, message := range messages {
		if err := v.validateMessageFormat(message, i); err != nil {
			return err
		}
	}
	return nil
}

// validateMessageFormat validates individual message format
func (v *ProtocolConformanceValidator) validateMessageFormat(message map[string]interface{}, index int) error {
	// Check JSON-RPC version
	if jsonrpc, exists := message["jsonrpc"]; !exists || jsonrpc != "2.0" {
		return fmt.Errorf("message at index %d missing or invalid jsonrpc field", index)
	}

	// Check for required fields based on message type
	if _, hasID := message["id"]; hasID {
		// Request or response
		if method, hasMethod := message["method"]; hasMethod {
			// Request message
			if methodStr, ok := method.(string); !ok || methodStr == "" {
				return fmt.Errorf("message at index %d has invalid method field", index)
			}
		} else if _, hasResult := message["result"]; !hasResult {
			if _, hasError := message["error"]; !hasError {
				return fmt.Errorf("message at index %d missing result or error field", index)
			}
		}
	} else {
		// Notification
		if method, hasMethod := message["method"]; !hasMethod {
			return fmt.Errorf("message at index %d missing method field for notification", index)
		} else if methodStr, ok := method.(string); !ok || methodStr == "" {
			return fmt.Errorf("message at index %d has invalid method field", index)
		}
	}

	return nil
}