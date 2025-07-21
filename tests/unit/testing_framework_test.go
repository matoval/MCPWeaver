package unit

import (
	"context"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/matoval/MCPWeaver/internal/testing"
)

// Test helper functions
func createTestConfig() *testing.TestConfig {
	return &testing.TestConfig{
		Timeout:                  1 * time.Minute,
		MaxConcurrentTests:       2,
		EnableParallelTesting:    false,
		ContinueOnFailure:        true,
		EnableSecurityScanning:   false,
		EnableLinting:           false,
		EnablePerformanceTesting: false,
		EnableIntegrationTesting: false,
		MCPProtocolVersion:      "2024-11-05",
		RequiredMethods:         []string{"initialize", "tools/list"},
		RequiredCapabilities:    []string{"tools"},
		MaxResponseTime:         time.Second,
		MaxMemoryUsage:          100 * 1024 * 1024,
		GenerateReport:          false,
		ReportFormat:           "json",
		LogLevel:               "error",
		RetryAttempts:          1,
		RetryDelay:             100 * time.Millisecond,
	}
}

func createTestServer(t *testing.T, validCode bool) string {
	tmpDir, err := os.MkdirTemp("", "test_server_*")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}

	var mainGoContent string
	var goModContent string

	if validCode {
		mainGoContent = `package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/sourcegraph/jsonrpc2"
)

type MCPServer struct{}

type Tool struct {
	Name        string     ` + "`json:\"name\"`" + `
	Description string     ` + "`json:\"description\"`" + `
	InputSchema InputSchema ` + "`json:\"inputSchema\"`" + `
}

type InputSchema struct {
	Type       string               ` + "`json:\"type\"`" + `
	Properties map[string]Property  ` + "`json:\"properties\"`" + `
	Required   []string            ` + "`json:\"required,omitempty\"`" + `
}

type Property struct {
	Type        string ` + "`json:\"type\"`" + `
	Description string ` + "`json:\"description\"`" + `
}

type ToolRequest struct {
	Name      string                 ` + "`json:\"name\"`" + `
	Arguments map[string]interface{} ` + "`json:\"arguments\"`" + `
}

type ToolResponse struct {
	Content []Content ` + "`json:\"content\"`" + `
}

type Content struct {
	Type string ` + "`json:\"type\"`" + `
	Text string ` + "`json:\"text\"`" + `
}

func (s *MCPServer) initialize(ctx context.Context, params interface{}) (interface{}, error) {
	return map[string]interface{}{
		"protocolVersion": "2024-11-05",
		"capabilities": map[string]interface{}{
			"tools": map[string]interface{}{},
		},
		"serverInfo": map[string]interface{}{
			"name":    "Test Server",
			"version": "1.0.0",
		},
	}, nil
}

func (s *MCPServer) toolsList(ctx context.Context, params interface{}) (interface{}, error) {
	tools := []Tool{
		{
			Name:        "test_tool",
			Description: "A test tool",
			InputSchema: InputSchema{
				Type: "object",
				Properties: map[string]Property{
					"message": {
						Type:        "string",
						Description: "Test message",
					},
				},
				Required: []string{"message"},
			},
		},
	}
	return map[string]interface{}{"tools": tools}, nil
}

func (s *MCPServer) toolsCall(ctx context.Context, params interface{}) (interface{}, error) {
	var req ToolRequest
	if err := json.Unmarshal(params.(json.RawMessage), &req); err != nil {
		return nil, err
	}

	if req.Name != "test_tool" {
		return nil, fmt.Errorf("Tool not found: %s", req.Name)
	}

	message, ok := req.Arguments["message"].(string)
	if !ok {
		return nil, fmt.Errorf("Missing required argument: message")
	}

	return ToolResponse{
		Content: []Content{
			{
				Type: "text",
				Text: fmt.Sprintf("Hello, %s!", message),
			},
		},
	}, nil
}

func main() {
	server := &MCPServer{}
	
	handler := jsonrpc2.HandlerWithError(func(ctx context.Context, conn *jsonrpc2.Conn, req *jsonrpc2.Request) (interface{}, error) {
		switch req.Method {
		case "initialize":
			return server.initialize(ctx, req.Params)
		case "tools/list":
			return server.toolsList(ctx, req.Params)
		case "tools/call":
			return server.toolsCall(ctx, req.Params)
		default:
			return nil, fmt.Errorf("method not found: %s", req.Method)
		}
	})

	conn := jsonrpc2.NewConn(context.Background(), jsonrpc2.NewBufferedStream(os.Stdin, os.Stdout), handler)
	<-conn.DisconnectNotify()
}`

		goModContent = `module test-mcp-server

go 1.21

require github.com/sourcegraph/jsonrpc2 v0.2.0
`
	} else {
		// Invalid code with syntax errors
		mainGoContent = `package main

import (
	"context"
	"fmt"
	// Missing imports

func main() {
	// Syntax error: missing opening brace
	if true
		fmt.Println("Hello")
	}
	// Missing closing brace
`

		goModContent = `module test-mcp-server

go 1.21

// Missing required dependencies
`
	}

	// Write main.go
	if err := os.WriteFile(filepath.Join(tmpDir, "main.go"), []byte(mainGoContent), 0644); err != nil {
		t.Fatalf("Failed to write main.go: %v", err)
	}

	// Write go.mod
	if err := os.WriteFile(filepath.Join(tmpDir, "go.mod"), []byte(goModContent), 0644); err != nil {
		t.Fatalf("Failed to write go.mod: %v", err)
	}

	return tmpDir
}

func TestTestSuite_Creation(t *testing.T) {
	config := createTestConfig()
	suite := testing.NewTestSuite(config)

	if suite == nil {
		t.Fatal("TestSuite creation failed")
	}

	// Test configuration retrieval
	retrievedConfig := suite.GetConfiguration()
	if retrievedConfig.Timeout != config.Timeout {
		t.Errorf("Expected timeout %v, got %v", config.Timeout, retrievedConfig.Timeout)
	}
}

func TestCompilationValidator_ValidCode(t *testing.T) {
	config := createTestConfig()
	validator := testing.NewCompilationValidator(config)

	if validator.Name() != "compilation" {
		t.Errorf("Expected validator name 'compilation', got '%s'", validator.Name())
	}

	if !validator.SupportsAsync() {
		t.Error("CompilationValidator should support async execution")
	}

	// Test with valid code
	serverPath := createTestServer(t, true)
	defer os.RemoveAll(serverPath)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	result, err := validator.Validate(ctx, serverPath)
	if err != nil {
		t.Fatalf("Validation failed: %v", err)
	}

	if !result.Success {
		t.Errorf("Expected validation to succeed, but got errors: %v", result.Errors)
	}

	if result.FilesValidated != 2 {
		t.Errorf("Expected 2 files validated, got %d", result.FilesValidated)
	}
}

func TestCompilationValidator_InvalidCode(t *testing.T) {
	config := createTestConfig()
	validator := testing.NewCompilationValidator(config)

	// Test with invalid code
	serverPath := createTestServer(t, false)
	defer os.RemoveAll(serverPath)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	result, err := validator.Validate(ctx, serverPath)
	if err != nil {
		t.Fatalf("Validation execution failed: %v", err)
	}

	if result.Success {
		t.Error("Expected validation to fail for invalid code")
	}

	if len(result.Errors) == 0 {
		t.Error("Expected validation errors for invalid code")
	}
}

func TestSyntaxValidator(t *testing.T) {
	config := createTestConfig()
	validator := testing.NewSyntaxValidator(config)

	if validator.Name() != "syntax" {
		t.Errorf("Expected validator name 'syntax', got '%s'", validator.Name())
	}

	if validator.SupportsAsync() {
		t.Error("SyntaxValidator should not support async execution")
	}

	// Test with valid code
	serverPath := createTestServer(t, true)
	defer os.RemoveAll(serverPath)

	ctx := context.Background()
	result, err := validator.Validate(ctx, serverPath)
	if err != nil {
		t.Fatalf("Syntax validation failed: %v", err)
	}

	if !result.Success {
		t.Errorf("Expected syntax validation to succeed, but got errors: %v", result.Errors)
	}
}

func TestConfigManager(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "config_test_*")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	configPath := filepath.Join(tmpDir, "test-config.json")
	configManager := testing.NewConfigManager(configPath)

	// Test initial configuration
	currentConfig := configManager.GetCurrentConfig()
	if currentConfig == nil {
		t.Fatal("Expected default configuration, got nil")
	}

	// Test profile creation
	customConfig := createTestConfig()
	customConfig.Timeout = 2 * time.Minute

	err = configManager.CreateProfile("test", "Test profile", customConfig)
	if err != nil {
		t.Fatalf("Failed to create profile: %v", err)
	}

	// Test profile retrieval
	retrievedConfig, err := configManager.GetProfile("test")
	if err != nil {
		t.Fatalf("Failed to get profile: %v", err)
	}

	if retrievedConfig.Timeout != customConfig.Timeout {
		t.Errorf("Expected timeout %v, got %v", customConfig.Timeout, retrievedConfig.Timeout)
	}

	// Test setting profile
	err = configManager.SetProfile("test")
	if err != nil {
		t.Fatalf("Failed to set profile: %v", err)
	}

	if configManager.GetCurrentProfile() != "test" {
		t.Errorf("Expected current profile 'test', got '%s'", configManager.GetCurrentProfile())
	}

	// Test configuration save/load
	err = configManager.SaveConfiguration()
	if err != nil {
		t.Fatalf("Failed to save configuration: %v", err)
	}

	// Create new manager and load
	newManager := testing.NewConfigManager(configPath)
	err = newManager.LoadConfiguration()
	if err != nil {
		t.Fatalf("Failed to load configuration: %v", err)
	}

	profiles := newManager.ListProfiles()
	if len(profiles) == 0 {
		t.Error("Expected profiles after loading, got none")
	}

	if _, exists := profiles["test"]; !exists {
		t.Error("Expected 'test' profile after loading")
	}
}

func TestTestPipeline_Creation(t *testing.T) {
	config := createTestConfig()
	pipeline := testing.NewTestPipeline(config)

	if pipeline == nil {
		t.Fatal("TestPipeline creation failed")
	}

	if pipeline.IsRunning() {
		t.Error("New pipeline should not be running")
	}

	status := pipeline.GetPipelineStatus()
	if status == nil {
		t.Fatal("Expected pipeline status, got nil")
	}

	if status["running"].(bool) {
		t.Error("Pipeline should not be running initially")
	}

	totalStages := status["totalStages"].(int)
	if totalStages == 0 {
		t.Error("Pipeline should have stages defined")
	}
}

func TestTestPipeline_Execution(t *testing.T) {
	config := createTestConfig()
	// Disable features that require external tools for faster testing
	config.EnableSecurityScanning = false
	config.EnableLinting = false
	config.Timeout = 30 * time.Second

	pipeline := testing.NewTestPipeline(config)

	// Test with valid server
	serverPath := createTestServer(t, true)
	defer os.RemoveAll(serverPath)

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
	defer cancel()

	result, err := pipeline.ExecutePipeline(ctx, serverPath)
	if err != nil {
		t.Fatalf("Pipeline execution failed: %v", err)
	}

	if result == nil {
		t.Fatal("Expected pipeline result, got nil")
	}

	if result.TotalStages == 0 {
		t.Error("Expected pipeline to have stages")
	}

	if result.CompletedStages == 0 {
		t.Error("Expected at least some stages to complete")
	}

	// Check that critical stages completed successfully
	if stageResult, exists := result.StageResults["pre_validation"]; exists {
		if !stageResult.Success {
			t.Errorf("Pre-validation stage should succeed: %s", stageResult.ErrorMessage)
		}
	}

	if stageResult, exists := result.StageResults["compilation_validation"]; exists {
		if !stageResult.Success {
			t.Errorf("Compilation validation stage should succeed: %s", stageResult.ErrorMessage)
		}
	}
}

func TestBatchTestRunner(t *testing.T) {
	config := createTestConfig()
	config.Timeout = 30 * time.Second
	config.EnableSecurityScanning = false
	config.EnableLinting = false

	batchRunner := testing.NewBatchTestRunner(config)

	// Create multiple test servers
	validServer := createTestServer(t, true)
	defer os.RemoveAll(validServer)

	invalidServer := createTestServer(t, false)
	defer os.RemoveAll(invalidServer)

	request := &testing.BatchTestRequest{
		RequestID:     "test_batch",
		ServerPaths:   []string{validServer, invalidServer},
		Parallel:      false, // Sequential for predictable testing
		MaxWorkers:    1,
		StopOnFailure: false,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()

	result, err := batchRunner.RunBatchTests(ctx, request)
	if err != nil {
		t.Fatalf("Batch test execution failed: %v", err)
	}

	if result.TotalServers != 2 {
		t.Errorf("Expected 2 servers, got %d", result.TotalServers)
	}

	if result.CompletedTests+result.FailedTests != result.TotalServers {
		t.Error("Completed + Failed tests should equal total servers")
	}

	// Check that both servers were processed
	if len(result.ServerResults) != 2 {
		t.Errorf("Expected 2 server results, got %d", len(result.ServerResults))
	}

	// Valid server should have better results than invalid server
	validResult := result.ServerResults[validServer]
	invalidResult := result.ServerResults[invalidServer]

	if validResult != nil && invalidResult != nil {
		if validResult.Success && invalidResult.Success {
			t.Error("Expected valid server to pass and invalid server to fail")
		}
	}
}

func TestTestReporter(t *testing.T) {
	config := createTestConfig()
	config.GenerateReport = true
	config.ReportFormat = "json"

	tmpDir, err := os.MkdirTemp("", "report_test_*")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	config.ReportOutputPath = filepath.Join(tmpDir, "test-report.json")

	reporter := testing.NewTestReporter(config)

	// Create mock test result
	testResult := &testing.TestResult{
		TestID:       "test_001",
		ServerPath:   "/test/server",
		Timestamp:    time.Now(),
		Duration:     30 * time.Second,
		Success:      true,
		TotalTests:   5,
		PassedTests:  4,
		FailedTests:  1,
		SkippedTests: 0,
		Errors:       []string{"Minor test error"},
		Warnings:     []string{"Test warning"},
		Recommendations: []string{"Consider optimization"},
	}

	// Test report generation
	err = reporter.GenerateReport(testResult)
	if err != nil {
		t.Fatalf("Report generation failed: %v", err)
	}

	// Verify report file was created
	if _, err := os.Stat(config.ReportOutputPath); os.IsNotExist(err) {
		t.Error("Report file was not created")
	}

	// Test metrics generation
	metrics, err := reporter.GenerateMetricsReport(testResult)
	if err != nil {
		t.Fatalf("Metrics generation failed: %v", err)
	}

	if metrics.TestID != testResult.TestID {
		t.Errorf("Expected test ID %s, got %s", testResult.TestID, metrics.TestID)
	}

	if metrics.OverallScore < 0 || metrics.OverallScore > 100 {
		t.Errorf("Overall score should be between 0 and 100, got %f", metrics.OverallScore)
	}
}

func TestDiagnosticsEngine(t *testing.T) {
	config := createTestConfig()
	diagnostics := testing.NewDiagnosticsEngine(config)

	// Create mock failed test result
	testResult := &testing.TestResult{
		TestID:    "failed_test",
		Success:   false,
		Errors:    []string{"compilation failed: syntax error", "undefined: fmt"},
		Warnings:  []string{"go fmt issues"},
	}

	// Create mock failed pipeline result
	pipelineResult := &testing.PipelineResult{
		Success: false,
		StageResults: map[string]*testing.StageResult{
			"compilation_validation": {
				StageName:    "compilation_validation",
				Success:      false,
				ErrorMessage: "go build failed",
			},
		},
		Errors: []string{"Pipeline stage failed"},
	}

	ctx := context.Background()
	serverPath := createTestServer(t, false)
	defer os.RemoveAll(serverPath)

	// Analyze failures
	report, err := diagnostics.AnalyzeFailures(ctx, serverPath, testResult, pipelineResult)
	if err != nil {
		t.Fatalf("Failure analysis failed: %v", err)
	}

	if report == nil {
		t.Fatal("Expected diagnostic report, got nil")
	}

	if report.FailureAnalysis == nil {
		t.Fatal("Expected failure analysis, got nil")
	}

	// Check that failure type was categorized
	if report.FailureAnalysis.FailureType == "" {
		t.Error("Expected failure type to be categorized")
	}

	// Check that recommendations were generated
	if len(report.Recommendations) == 0 {
		t.Error("Expected recommendations for failed tests")
	}

	// Check that troubleshooting guide was generated
	if report.TroubleshootingGuide == nil {
		t.Error("Expected troubleshooting guide")
	}

	// Test saving diagnostic report
	tmpDir, err := os.MkdirTemp("", "diag_test_*")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	reportPath := filepath.Join(tmpDir, "diagnostic-report.json")
	err = diagnostics.SaveDiagnosticReport(report, reportPath)
	if err != nil {
		t.Fatalf("Failed to save diagnostic report: %v", err)
	}

	// Verify report file was created
	if _, err := os.Stat(reportPath); os.IsNotExist(err) {
		t.Error("Diagnostic report file was not created")
	}
}

// Benchmark tests for performance validation
func BenchmarkTestSuite_ValidCode(b *testing.B) {
	config := createTestConfig()
	config.EnableSecurityScanning = false
	config.EnableLinting = false
	config.EnablePerformanceTesting = false
	config.EnableIntegrationTesting = false

	suite := testing.NewTestSuite(config)
	serverPath := createTestServer(b, true)
	defer os.RemoveAll(serverPath)

	ctx := context.Background()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		result, err := suite.RunTests(ctx, serverPath)
		if err != nil {
			b.Fatalf("Test execution failed: %v", err)
		}
		if !result.Success {
			b.Errorf("Expected tests to pass, got errors: %v", result.Errors)
		}
	}
}

func BenchmarkCompilationValidator(b *testing.B) {
	config := createTestConfig()
	validator := testing.NewCompilationValidator(config)
	serverPath := createTestServer(b, true)
	defer os.RemoveAll(serverPath)

	ctx := context.Background()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		result, err := validator.Validate(ctx, serverPath)
		if err != nil {
			b.Fatalf("Validation failed: %v", err)
		}
		if !result.Success {
			b.Errorf("Expected validation to pass, got errors: %v", result.Errors)
		}
	}
}

// Test edge cases and error conditions
func TestTestSuite_MissingFiles(t *testing.T) {
	config := createTestConfig()
	suite := testing.NewTestSuite(config)

	// Test with non-existent directory
	ctx := context.Background()
	result, err := suite.RunTests(ctx, "/nonexistent/path")
	if err == nil {
		t.Error("Expected error for non-existent path")
	}

	// Test with empty directory
	tmpDir, err := os.MkdirTemp("", "empty_test_*")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	result, err = suite.RunTests(ctx, tmpDir)
	if err == nil {
		t.Error("Expected error for empty directory")
	}
	if result != nil && result.Success {
		t.Error("Expected test to fail for missing files")
	}
}

func TestConfigManager_Validation(t *testing.T) {
	configManager := testing.NewConfigManager("test-config.json")

	// Test with nil config
	err := configManager.ValidateConfig(nil)
	if err == nil {
		t.Error("Expected error for nil config")
	}

	// Test with invalid timeout
	invalidConfig := createTestConfig()
	invalidConfig.Timeout = -1 * time.Second
	err = configManager.ValidateConfig(invalidConfig)
	if err == nil {
		t.Error("Expected error for negative timeout")
	}

	// Test with invalid max concurrent tests
	invalidConfig = createTestConfig()
	invalidConfig.MaxConcurrentTests = 0
	err = configManager.ValidateConfig(invalidConfig)
	if err == nil {
		t.Error("Expected error for zero max concurrent tests")
	}

	// Test with invalid report format
	invalidConfig = createTestConfig()
	invalidConfig.ReportFormat = "invalid"
	err = configManager.ValidateConfig(invalidConfig)
	if err == nil {
		t.Error("Expected error for invalid report format")
	}

	// Test with valid config
	validConfig := createTestConfig()
	err = configManager.ValidateConfig(validConfig)
	if err != nil {
		t.Errorf("Expected no error for valid config, got: %v", err)
	}
}

func TestTestPipeline_ConcurrentExecution(t *testing.T) {
	config := createTestConfig()
	config.EnableSecurityScanning = false
	config.EnableLinting = false

	pipeline1 := testing.NewTestPipeline(config)
	pipeline2 := testing.NewTestPipeline(config)

	serverPath := createTestServer(t, true)
	defer os.RemoveAll(serverPath)

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
	defer cancel()

	// Try to run both pipelines concurrently on the same instance
	go func() {
		pipeline1.ExecutePipeline(ctx, serverPath)
	}()

	// Second execution should not fail due to concurrency
	result, err := pipeline2.ExecutePipeline(ctx, serverPath)
	if err != nil {
		t.Fatalf("Concurrent pipeline execution failed: %v", err)
	}

	if result == nil {
		t.Fatal("Expected pipeline result")
	}
}