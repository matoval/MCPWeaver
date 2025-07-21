package testing

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/exec"
	"sync"
	"time"
)

// TestSuite represents the comprehensive testing framework for generated MCP servers
type TestSuite struct {
	config          *TestConfig
	validators      []Validator
	protocolTester  *ProtocolTester
	integrationTest *IntegrationTester
	performanceTester *PerformanceTester
	reporter        *TestReporter
	mutex           sync.RWMutex
}

// TestConfig holds configuration for the testing framework
type TestConfig struct {
	// Basic execution settings
	Timeout               time.Duration `json:"timeout"`
	MaxConcurrentTests    int          `json:"maxConcurrentTests"`
	EnableParallelTesting bool         `json:"enableParallelTesting"`
	ContinueOnFailure     bool         `json:"continueOnFailure"`

	// Feature flags
	EnableSecurityScanning   bool `json:"enableSecurityScanning"`
	EnableLinting           bool `json:"enableLinting"`
	EnablePerformanceTesting bool `json:"enablePerformanceTesting"`
	EnableIntegrationTesting bool `json:"enableIntegrationTesting"`

	// MCP protocol settings
	MCPProtocolVersion   string   `json:"mcpProtocolVersion"`
	RequiredMethods      []string `json:"requiredMethods"`
	RequiredCapabilities []string `json:"requiredCapabilities"`

	// Performance thresholds
	MaxResponseTime      time.Duration `json:"maxResponseTime"`
	MaxMemoryUsage       int64         `json:"maxMemoryUsage"`

	// External tools and paths
	TestDataPath      string `json:"testDataPath"`
	MCPClientPath     string `json:"mcpClientPath"`

	// Report generation
	GenerateReport    bool   `json:"generateReport"`
	ReportFormat      string `json:"reportFormat"`
	ReportOutputPath  string `json:"reportOutputPath"`
	LogLevel         string `json:"logLevel"`

	// Retry configuration
	RetryAttempts     int           `json:"retryAttempts"`
	RetryDelay        time.Duration `json:"retryDelay"`
}

// TestResult represents the overall test results
type TestResult struct {
	TestID          string                   `json:"testId"`
	ServerPath      string                   `json:"serverPath"`
	Timestamp       time.Time                `json:"timestamp"`
	Duration        time.Duration            `json:"duration"`
	Success         bool                     `json:"success"`
	
	// Validation Results
	ValidationResults   map[string]*ValidationResult `json:"validationResults"`
	
	// Protocol Testing Results
	ProtocolResults     *ProtocolTestResult  `json:"protocolResults"`
	
	// Integration Testing Results
	IntegrationResults  *IntegrationTestResult `json:"integrationResults"`
	
	// Performance Testing Results
	PerformanceResults  *PerformanceTestResult `json:"performanceResults"`
	
	// Summary
	TotalTests         int      `json:"totalTests"`
	PassedTests        int      `json:"passedTests"`
	FailedTests        int      `json:"failedTests"`
	SkippedTests       int      `json:"skippedTests"`
	Errors             []string `json:"errors,omitempty"`
	Warnings           []string `json:"warnings,omitempty"`
	Recommendations    []string `json:"recommendations,omitempty"`
}

// ValidationResult represents validation test results
type ValidationResult struct {
	ValidatorName  string        `json:"validatorName"`
	Success        bool          `json:"success"`
	Duration       time.Duration `json:"duration"`
	Errors         []string      `json:"errors,omitempty"`
	Warnings       []string      `json:"warnings,omitempty"`
	FilesValidated int           `json:"filesValidated"`
	Details        interface{}   `json:"details,omitempty"`
}

// ProtocolTestResult represents MCP protocol compliance test results
type ProtocolTestResult struct {
	Success                bool                     `json:"success"`
	Duration               time.Duration            `json:"duration"`
	ProtocolVersion        string                   `json:"protocolVersion"`
	SupportedMethods       []string                 `json:"supportedMethods"`
	SupportedCapabilities  []string                 `json:"supportedCapabilities"`
	MethodTests            map[string]*MethodTest   `json:"methodTests"`
	CapabilityTests        map[string]*CapabilityTest `json:"capabilityTests"`
	Errors                 []string                 `json:"errors,omitempty"`
}

// MethodTest represents individual method test results
type MethodTest struct {
	Method         string        `json:"method"`
	Success        bool          `json:"success"`
	ResponseTime   time.Duration `json:"responseTime"`
	Request        interface{}   `json:"request"`
	Response       interface{}   `json:"response"`
	ErrorMessage   string        `json:"errorMessage,omitempty"`
}

// CapabilityTest represents capability test results
type CapabilityTest struct {
	Capability     string        `json:"capability"`
	Success        bool          `json:"success"`
	Supported      bool          `json:"supported"`
	TestDetails    interface{}   `json:"testDetails"`
	ErrorMessage   string        `json:"errorMessage,omitempty"`
}

// IntegrationTestResult represents integration test results
type IntegrationTestResult struct {
	Success           bool                            `json:"success"`
	Duration          time.Duration                   `json:"duration"`
	ScenarioResults   map[string]*ScenarioTestResult  `json:"scenarioResults"`
	ClientCompatibility map[string]bool               `json:"clientCompatibility"`
	Errors            []string                        `json:"errors,omitempty"`
}

// ScenarioTestResult represents individual scenario test results
type ScenarioTestResult struct {
	Scenario       string        `json:"scenario"`
	Success        bool          `json:"success"`
	Duration       time.Duration `json:"duration"`
	Steps          []StepResult  `json:"steps"`
	ErrorMessage   string        `json:"errorMessage,omitempty"`
}

// StepResult represents individual test step results
type StepResult struct {
	Step           string        `json:"step"`
	Success        bool          `json:"success"`
	Duration       time.Duration `json:"duration"`
	Details        interface{}   `json:"details,omitempty"`
	ErrorMessage   string        `json:"errorMessage,omitempty"`
}

// PerformanceTestResult represents performance test results
type PerformanceTestResult struct {
	Success                  bool                      `json:"success"`
	Duration                 time.Duration             `json:"duration"`
	
	// Response Time Metrics
	AverageResponseTime      time.Duration             `json:"averageResponseTime"`
	MedianResponseTime       time.Duration             `json:"medianResponseTime"`
	P95ResponseTime          time.Duration             `json:"p95ResponseTime"`
	P99ResponseTime          time.Duration             `json:"p99ResponseTime"`
	MaxResponseTime          time.Duration             `json:"maxResponseTime"`
	
	// Memory Usage Metrics
	AverageMemoryUsage       int64                     `json:"averageMemoryUsage"`
	PeakMemoryUsage          int64                     `json:"peakMemoryUsage"`
	MemoryLeakDetected       bool                      `json:"memoryLeakDetected"`
	
	// Throughput Metrics
	RequestsPerSecond        float64                   `json:"requestsPerSecond"`
	ConcurrentConnections    int                       `json:"concurrentConnections"`
	SuccessfulRequests       int                       `json:"successfulRequests"`
	FailedRequests           int                       `json:"failedRequests"`
	
	// Load Test Results
	LoadTestResults          map[string]*LoadTestMetric `json:"loadTestResults"`
	
	Errors                   []string                  `json:"errors,omitempty"`
}

// LoadTestMetric represents load test metrics for specific scenarios
type LoadTestMetric struct {
	Scenario              string        `json:"scenario"`
	Duration              time.Duration `json:"duration"`
	TotalRequests         int           `json:"totalRequests"`
	SuccessfulRequests    int           `json:"successfulRequests"`
	FailedRequests        int           `json:"failedRequests"`
	AverageResponseTime   time.Duration `json:"averageResponseTime"`
	RequestsPerSecond     float64       `json:"requestsPerSecond"`
	ErrorRate             float64       `json:"errorRate"`
}

// Validator interface for different types of validation
type Validator interface {
	Name() string
	Validate(ctx context.Context, serverPath string) (*ValidationResult, error)
	SupportsAsync() bool
}

// NewTestSuite creates a new comprehensive test suite
func NewTestSuite(config *TestConfig) *TestSuite {
	ts := &TestSuite{
		config:     config,
		validators: make([]Validator, 0),
		reporter:   NewTestReporter(config),
	}
	
	// Initialize components
	ts.protocolTester = NewProtocolTester(config)
	ts.integrationTest = NewIntegrationTester(config)
	ts.performanceTester = NewPerformanceTester(config)
	
	// Register default validators
	ts.RegisterDefaultValidators()
	
	return ts
}

// RegisterValidator adds a validator to the test suite
func (ts *TestSuite) RegisterValidator(validator Validator) {
	ts.mutex.Lock()
	defer ts.mutex.Unlock()
	
	ts.validators = append(ts.validators, validator)
	log.Printf("Registered validator: %s", validator.Name())
}

// RegisterDefaultValidators registers the standard set of validators
func (ts *TestSuite) RegisterDefaultValidators() {
	// Register built-in validators
	ts.RegisterValidator(NewCompilationValidator(ts.config))
	ts.RegisterValidator(NewSyntaxValidator(ts.config))
	ts.RegisterValidator(NewLintValidator(ts.config))
	ts.RegisterValidator(NewSecurityValidator(ts.config))
	ts.RegisterValidator(NewDependencyValidator(ts.config))
}

// RunTests executes the complete test suite
func (ts *TestSuite) RunTests(ctx context.Context, serverPath string) (*TestResult, error) {
	startTime := time.Now()
	
	result := &TestResult{
		TestID:             generateTestID(),
		ServerPath:         serverPath,
		Timestamp:          startTime,
		ValidationResults:  make(map[string]*ValidationResult),
		Errors:             make([]string, 0),
		Warnings:           make([]string, 0),
		Recommendations:    make([]string, 0),
	}
	
	log.Printf("Starting comprehensive test suite for server: %s", serverPath)
	
	// Phase 1: Validation Tests
	if err := ts.runValidationTests(ctx, serverPath, result); err != nil {
		result.Errors = append(result.Errors, fmt.Sprintf("Validation tests failed: %v", err))
	}
	
	// Phase 2: Protocol Compliance Tests
	if err := ts.runProtocolTests(ctx, serverPath, result); err != nil {
		result.Errors = append(result.Errors, fmt.Sprintf("Protocol tests failed: %v", err))
	}
	
	// Phase 3: Integration Tests (only if server is valid)
	if result.Success && ts.config.EnableIntegrationTesting {
		if err := ts.runIntegrationTests(ctx, serverPath, result); err != nil {
			result.Errors = append(result.Errors, fmt.Sprintf("Integration tests failed: %v", err))
		}
	}
	
	// Phase 4: Performance Tests (only if server is functional)
	if result.Success && ts.config.EnablePerformanceTesting {
		if err := ts.runPerformanceTests(ctx, serverPath, result); err != nil {
			result.Errors = append(result.Errors, fmt.Sprintf("Performance tests failed: %v", err))
		}
	}
	
	// Calculate final results
	result.Duration = time.Since(startTime)
	ts.calculateTestSummary(result)
	
	// Generate report if requested
	if ts.config.GenerateReport {
		if err := ts.reporter.GenerateReport(result); err != nil {
			log.Printf("Failed to generate test report: %v", err)
		}
	}
	
	log.Printf("Test suite completed in %v. Success: %t", result.Duration, result.Success)
	
	return result, nil
}

// runValidationTests executes all validation tests
func (ts *TestSuite) runValidationTests(ctx context.Context, serverPath string, result *TestResult) error {
	log.Printf("Running validation tests...")
	
	for _, validator := range ts.validators {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}
		
		log.Printf("Running validator: %s", validator.Name())
		
		validationResult, err := validator.Validate(ctx, serverPath)
		if err != nil {
			validationResult = &ValidationResult{
				ValidatorName: validator.Name(),
				Success:       false,
				Errors:        []string{err.Error()},
			}
		}
		
		result.ValidationResults[validator.Name()] = validationResult
		result.TotalTests++
		
		if validationResult.Success {
			result.PassedTests++
		} else {
			result.FailedTests++
			result.Errors = append(result.Errors, validationResult.Errors...)
		}
		
		if len(validationResult.Warnings) > 0 {
			result.Warnings = append(result.Warnings, validationResult.Warnings...)
		}
	}
	
	return nil
}

// runProtocolTests executes MCP protocol compliance tests
func (ts *TestSuite) runProtocolTests(ctx context.Context, serverPath string, result *TestResult) error {
	log.Printf("Running protocol compliance tests...")
	
	protocolResult, err := ts.protocolTester.TestCompliance(ctx, serverPath)
	if err != nil {
		return fmt.Errorf("protocol testing failed: %w", err)
	}
	
	result.ProtocolResults = protocolResult
	result.TotalTests++
	
	if protocolResult.Success {
		result.PassedTests++
	} else {
		result.FailedTests++
		result.Errors = append(result.Errors, protocolResult.Errors...)
	}
	
	return nil
}

// runIntegrationTests executes integration tests
func (ts *TestSuite) runIntegrationTests(ctx context.Context, serverPath string, result *TestResult) error {
	log.Printf("Running integration tests...")
	
	integrationResult, err := ts.integrationTest.TestIntegration(ctx, serverPath)
	if err != nil {
		return fmt.Errorf("integration testing failed: %w", err)
	}
	
	result.IntegrationResults = integrationResult
	result.TotalTests++
	
	if integrationResult.Success {
		result.PassedTests++
	} else {
		result.FailedTests++
		result.Errors = append(result.Errors, integrationResult.Errors...)
	}
	
	return nil
}

// runPerformanceTests executes performance tests
func (ts *TestSuite) runPerformanceTests(ctx context.Context, serverPath string, result *TestResult) error {
	log.Printf("Running performance tests...")
	
	performanceResult, err := ts.performanceTester.TestPerformance(ctx, serverPath)
	if err != nil {
		return fmt.Errorf("performance testing failed: %w", err)
	}
	
	result.PerformanceResults = performanceResult
	result.TotalTests++
	
	if performanceResult.Success {
		result.PassedTests++
	} else {
		result.FailedTests++
		result.Errors = append(result.Errors, performanceResult.Errors...)
	}
	
	return nil
}

// calculateTestSummary calculates the overall test summary
func (ts *TestSuite) calculateTestSummary(result *TestResult) {
	// Determine overall success
	result.Success = result.FailedTests == 0 && len(result.Errors) == 0
	
	// Generate recommendations based on results
	if result.PerformanceResults != nil {
		if result.PerformanceResults.AverageResponseTime > ts.config.MaxResponseTime {
			result.Recommendations = append(result.Recommendations, 
				"Consider optimizing server response time")
		}
		
		if result.PerformanceResults.PeakMemoryUsage > ts.config.MaxMemoryUsage {
			result.Recommendations = append(result.Recommendations, 
				"Consider optimizing memory usage")
		}
	}
	
	if result.ProtocolResults != nil && !result.ProtocolResults.Success {
		result.Recommendations = append(result.Recommendations, 
			"Fix MCP protocol compliance issues for better client compatibility")
	}
}

// ValidateServer performs quick validation of a generated server
func (ts *TestSuite) ValidateServer(ctx context.Context, serverPath string) (*ValidationResult, error) {
	// Run only essential validation checks
	compilationValidator := NewCompilationValidator(ts.config)
	syntaxValidator := NewSyntaxValidator(ts.config)
	
	// Check compilation first
	compilationResult, err := compilationValidator.Validate(ctx, serverPath)
	if err != nil || !compilationResult.Success {
		return compilationResult, err
	}
	
	// Check syntax
	syntaxResult, err := syntaxValidator.Validate(ctx, serverPath)
	if err != nil || !syntaxResult.Success {
		return syntaxResult, err
	}
	
	return &ValidationResult{
		ValidatorName:  "quick-validation",
		Success:        true,
		FilesValidated: compilationResult.FilesValidated + syntaxResult.FilesValidated,
		Warnings:       append(compilationResult.Warnings, syntaxResult.Warnings...),
	}, nil
}

// generateTestID generates a unique test ID
func generateTestID() string {
	return fmt.Sprintf("test_%d", time.Now().UnixNano())
}

// GetDefaultConfig returns a default test configuration
func GetDefaultConfig() *TestConfig {
	return &TestConfig{
		Timeout:                  5 * time.Minute,
		MaxConcurrentTests:       3,
		EnableParallelTesting:    true,
		ContinueOnFailure:        false,
		EnableSecurityScanning:   true,
		EnableLinting:           true,
		EnablePerformanceTesting: true,
		EnableIntegrationTesting: false, // Disabled by default
		MCPProtocolVersion:      "2024-11-05",
		RequiredMethods:         []string{"initialize", "tools/list", "tools/call"},
		RequiredCapabilities:    []string{"tools"},
		MaxResponseTime:         time.Second,
		MaxMemoryUsage:          100 << 20, // 100MB
		TestDataPath:            "",
		MCPClientPath:           "",
		GenerateReport:          true,
		ReportFormat:            "json",
		ReportOutputPath:        "",
		LogLevel:               "info",
		RetryAttempts:          2,
		RetryDelay:             time.Second,
	}
}

// Helper methods for TestSuite

// fileExists checks if a file exists
func (ts *TestSuite) fileExists(path string) error {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return fmt.Errorf("file does not exist: %s", path)
	}
	return nil
}

// commandExists checks if a command is available in PATH
func (ts *TestSuite) commandExists(command string) bool {
	_, err := exec.LookPath(command)
	return err == nil
}

// runCommand executes a command in the specified directory
func (ts *TestSuite) runCommand(ctx context.Context, dir string, name string, args ...string) error {
	cmd := exec.CommandContext(ctx, name, args...)
	cmd.Dir = dir
	return cmd.Run()
}

// ensureDir creates a directory if it doesn't exist
func (ts *TestSuite) ensureDir(path string) error {
	return os.MkdirAll(path, 0755)
}