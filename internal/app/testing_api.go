package app

import (
	"context"
	"fmt"
	"path/filepath"

	"MCPWeaver/internal/testing"
)

// RunServerTests runs comprehensive tests on a generated MCP server
func (a *App) RunServerTests(serverPath string, configOptions *TestConfigOptions) (*TestResult, error) {
	if serverPath == "" {
		return nil, a.createAPIError("validation", ErrCodeValidation, "Server path is required", nil)
	}

	// Validate server path exists
	if err := a.fileExists(serverPath); err != nil {
		return nil, a.createAPIError("file_system", ErrCodeFileAccess, "Server path does not exist", map[string]string{
			"path": serverPath,
		})
	}

	// Create test configuration
	config := testing.GetDefaultConfig()
	if configOptions != nil {
		a.applyTestConfigOptions(config, configOptions)
	}

	// Create test suite
	testSuite := testing.NewTestSuite(config)

	// Set up context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), config.Timeout)
	defer cancel()

	// Run comprehensive tests
	result, err := testSuite.RunTests(ctx, serverPath)
	if err != nil {
		return nil, a.createAPIError("internal", ErrCodeInternalError, "Test execution failed", map[string]string{
			"error": err.Error(),
		})
	}

	// Convert internal result to API result
	apiResult := a.convertTestResult(result)

	// Emit test completion event
	a.emitEvent("test:completed", map[string]interface{}{
		"serverPath": serverPath,
		"success":    apiResult.Success,
		"duration":   apiResult.Duration,
		"totalTests": apiResult.TotalTests,
	})

	// Send notification
	if apiResult.Success {
		a.emitNotification("success", "Tests Passed", 
			fmt.Sprintf("All %d tests passed for server", apiResult.TotalTests))
	} else {
		a.emitNotification("warning", "Tests Failed", 
			fmt.Sprintf("%d of %d tests failed", apiResult.FailedTests, apiResult.TotalTests))
	}

	return apiResult, nil
}

// ValidateGeneratedServer performs quick validation of a generated server
func (a *App) ValidateGeneratedServer(serverPath string) (*ValidationResult, error) {
	if serverPath == "" {
		return nil, a.createAPIError("validation", ErrCodeValidation, "Server path is required", nil)
	}

	// Validate server path exists
	if err := a.fileExists(serverPath); err != nil {
		return nil, a.createAPIError("file_system", ErrCodeFileAccess, "Server path does not exist", map[string]string{
			"path": serverPath,
		})
	}

	// Create test configuration for validation only
	config := testing.GetDefaultConfig()
	config.EnablePerformanceTesting = false
	config.EnableIntegrationTesting = false

	// Create test suite
	testSuite := testing.NewTestSuite(config)

	// Set up context with shorter timeout for validation
	ctx, cancel := context.WithTimeout(context.Background(), config.Timeout/2)
	defer cancel()

	// Run quick validation
	result, err := testSuite.ValidateServer(ctx, serverPath)
	if err != nil {
		return nil, a.createAPIError("internal", ErrCodeInternalError, "Validation failed", map[string]string{
			"error": err.Error(),
		})
	}

	// Convert internal result to API result
	apiResult := a.convertValidationResult(result)

	// Emit validation event
	a.emitEvent("validation:completed", map[string]interface{}{
		"serverPath": serverPath,
		"success":    apiResult.Success,
		"duration":   apiResult.Duration,
	})

	return apiResult, nil
}

// RunPerformanceTests runs only performance tests on a server
func (a *App) RunPerformanceTests(serverPath string, configOptions *PerformanceTestOptions) (*PerformanceTestResult, error) {
	if serverPath == "" {
		return nil, a.createAPIError("validation", ErrCodeValidation, "Server path is required", nil)
	}

	// Validate server path exists
	if err := a.fileExists(serverPath); err != nil {
		return nil, a.createAPIError("file_system", ErrCodeFileAccess, "Server path does not exist", map[string]string{
			"path": serverPath,
		})
	}

	// Create test configuration for performance testing
	config := testing.GetDefaultConfig()
	config.EnablePerformanceTesting = true
	config.EnableIntegrationTesting = false
	config.EnableSecurityScanning = false
	config.EnableLinting = false

	if configOptions != nil {
		a.applyPerformanceTestOptions(config, configOptions)
	}

	// Create performance tester
	performanceTester := testing.NewPerformanceTester(config)

	// Set up context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), config.Timeout)
	defer cancel()

	// Run performance tests
	result, err := performanceTester.TestPerformance(ctx, serverPath)
	if err != nil {
		return nil, a.createAPIError("internal", ErrCodeInternalError, "Performance testing failed", map[string]string{
			"error": err.Error(),
		})
	}

	// Convert internal result to API result
	apiResult := a.convertPerformanceTestResult(result)

	// Emit performance test completion event
	a.emitEvent("performance_test:completed", map[string]interface{}{
		"serverPath":           serverPath,
		"success":              apiResult.Success,
		"duration":             apiResult.Duration,
		"averageResponseTime":  apiResult.AverageResponseTime,
		"peakMemoryUsage":      apiResult.PeakMemoryUsage,
		"requestsPerSecond":    apiResult.RequestsPerSecond,
	})

	return apiResult, nil
}

// RunProtocolComplianceTests runs only MCP protocol compliance tests
func (a *App) RunProtocolComplianceTests(serverPath string) (*ProtocolTestResult, error) {
	if serverPath == "" {
		return nil, a.createAPIError("validation", ErrCodeValidation, "Server path is required", nil)
	}

	// Validate server path exists
	if err := a.fileExists(serverPath); err != nil {
		return nil, a.createAPIError("file_system", ErrCodeFileAccess, "Server path does not exist", map[string]string{
			"path": serverPath,
		})
	}

	// Create test configuration for protocol testing
	config := testing.GetDefaultConfig()
	config.EnablePerformanceTesting = false
	config.EnableIntegrationTesting = false
	config.EnableSecurityScanning = false
	config.EnableLinting = false

	// Create protocol tester
	protocolTester := testing.NewProtocolTester(config)

	// Set up context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), config.Timeout)
	defer cancel()

	// Run protocol compliance tests
	result, err := protocolTester.TestCompliance(ctx, serverPath)
	if err != nil {
		return nil, a.createAPIError("internal", ErrCodeInternalError, "Protocol testing failed", map[string]string{
			"error": err.Error(),
		})
	}

	// Convert internal result to API result
	apiResult := a.convertProtocolTestResult(result)

	// Emit protocol test completion event
	a.emitEvent("protocol_test:completed", map[string]interface{}{
		"serverPath":             serverPath,
		"success":                apiResult.Success,
		"duration":               apiResult.Duration,
		"protocolVersion":        apiResult.ProtocolVersion,
		"supportedMethods":       apiResult.SupportedMethods,
		"supportedCapabilities":  apiResult.SupportedCapabilities,
	})

	return apiResult, nil
}

// GetTestConfiguration returns the default test configuration
func (a *App) GetTestConfiguration() *TestConfig {
	config := testing.GetDefaultConfig()
	return a.convertTestConfig(config)
}

// UpdateTestConfiguration updates the test configuration
func (a *App) UpdateTestConfiguration(config *TestConfig) error {
	if config == nil {
		return a.createAPIError("validation", ErrCodeValidation, "Test configuration is required", nil)
	}

	// Validate configuration
	if err := a.validateTestConfig(config); err != nil {
		return err
	}

	// Store configuration (in a real implementation, this would be persisted)
	// For now, we just validate it
	a.emitNotification("info", "Configuration Updated", "Test configuration has been updated")

	return nil
}

// Helper methods for conversion between internal and API types

func (a *App) convertTestResult(internal *testing.TestResult) *TestResult {
	return &TestResult{
		TestID:     internal.TestID,
		ServerPath: internal.ServerPath,
		Timestamp:  internal.Timestamp,
		Duration:   internal.Duration,
		Success:    internal.Success,

		ValidationResults:  a.convertValidationResults(internal.ValidationResults),
		ProtocolResults:    a.convertProtocolTestResult(internal.ProtocolResults),
		IntegrationResults: a.convertIntegrationTestResult(internal.IntegrationResults),
		PerformanceResults: a.convertPerformanceTestResult(internal.PerformanceResults),

		TotalTests:      internal.TotalTests,
		PassedTests:     internal.PassedTests,
		FailedTests:     internal.FailedTests,
		SkippedTests:    internal.SkippedTests,
		Errors:          internal.Errors,
		Warnings:        internal.Warnings,
		Recommendations: internal.Recommendations,
	}
}

func (a *App) convertValidationResult(internal *testing.ValidationResult) *ValidationResult {
	return &ValidationResult{
		ValidatorName:  internal.ValidatorName,
		Success:        internal.Success,
		Duration:       internal.Duration,
		Errors:         internal.Errors,
		Warnings:       internal.Warnings,
		FilesValidated: internal.FilesValidated,
		Details:        internal.Details,
	}
}

func (a *App) convertValidationResults(internal map[string]*testing.ValidationResult) map[string]*ValidationResult {
	results := make(map[string]*ValidationResult)
	for name, result := range internal {
		results[name] = a.convertValidationResult(result)
	}
	return results
}

func (a *App) convertProtocolTestResult(internal *testing.ProtocolTestResult) *ProtocolTestResult {
	if internal == nil {
		return nil
	}

	return &ProtocolTestResult{
		Success:               internal.Success,
		Duration:              internal.Duration,
		ProtocolVersion:       internal.ProtocolVersion,
		SupportedMethods:      internal.SupportedMethods,
		SupportedCapabilities: internal.SupportedCapabilities,
		MethodTests:           a.convertMethodTests(internal.MethodTests),
		CapabilityTests:       a.convertCapabilityTests(internal.CapabilityTests),
		Errors:                internal.Errors,
	}
}

func (a *App) convertMethodTests(internal map[string]*testing.MethodTest) map[string]*MethodTest {
	tests := make(map[string]*MethodTest)
	for name, test := range internal {
		tests[name] = &MethodTest{
			Method:       test.Method,
			Success:      test.Success,
			ResponseTime: test.ResponseTime,
			Request:      test.Request,
			Response:     test.Response,
			ErrorMessage: test.ErrorMessage,
		}
	}
	return tests
}

func (a *App) convertCapabilityTests(internal map[string]*testing.CapabilityTest) map[string]*CapabilityTest {
	tests := make(map[string]*CapabilityTest)
	for name, test := range internal {
		tests[name] = &CapabilityTest{
			Capability:   test.Capability,
			Success:      test.Success,
			Supported:    test.Supported,
			TestDetails:  test.TestDetails,
			ErrorMessage: test.ErrorMessage,
		}
	}
	return tests
}

func (a *App) convertIntegrationTestResult(internal *testing.IntegrationTestResult) *IntegrationTestResult {
	if internal == nil {
		return nil
	}

	return &IntegrationTestResult{
		Success:             internal.Success,
		Duration:            internal.Duration,
		ScenarioResults:     a.convertScenarioResults(internal.ScenarioResults),
		ClientCompatibility: internal.ClientCompatibility,
		Errors:              internal.Errors,
	}
}

func (a *App) convertScenarioResults(internal map[string]*testing.ScenarioTestResult) map[string]*ScenarioTestResult {
	results := make(map[string]*ScenarioTestResult)
	for name, result := range internal {
		results[name] = &ScenarioTestResult{
			Scenario:     result.Scenario,
			Success:      result.Success,
			Duration:     result.Duration,
			Steps:        a.convertStepResults(result.Steps),
			ErrorMessage: result.ErrorMessage,
		}
	}
	return results
}

func (a *App) convertStepResults(internal []testing.StepResult) []StepResult {
	results := make([]StepResult, len(internal))
	for i, step := range internal {
		results[i] = StepResult{
			Step:         step.Step,
			Success:      step.Success,
			Duration:     step.Duration,
			Details:      step.Details,
			ErrorMessage: step.ErrorMessage,
		}
	}
	return results
}

func (a *App) convertPerformanceTestResult(internal *testing.PerformanceTestResult) *PerformanceTestResult {
	if internal == nil {
		return nil
	}

	return &PerformanceTestResult{
		Success:  internal.Success,
		Duration: internal.Duration,

		AverageResponseTime: internal.AverageResponseTime,
		MedianResponseTime:  internal.MedianResponseTime,
		P95ResponseTime:     internal.P95ResponseTime,
		P99ResponseTime:     internal.P99ResponseTime,
		MaxResponseTime:     internal.MaxResponseTime,

		AverageMemoryUsage: internal.AverageMemoryUsage,
		PeakMemoryUsage:    internal.PeakMemoryUsage,
		MemoryLeakDetected: internal.MemoryLeakDetected,

		RequestsPerSecond:     internal.RequestsPerSecond,
		ConcurrentConnections: internal.ConcurrentConnections,
		SuccessfulRequests:    internal.SuccessfulRequests,
		FailedRequests:        internal.FailedRequests,

		LoadTestResults: a.convertLoadTestResults(internal.LoadTestResults),
		Errors:          internal.Errors,
	}
}

func (a *App) convertLoadTestResults(internal map[string]*testing.LoadTestMetric) map[string]*LoadTestMetric {
	results := make(map[string]*LoadTestMetric)
	for name, metric := range internal {
		results[name] = &LoadTestMetric{
			Scenario:            metric.Scenario,
			Duration:            metric.Duration,
			TotalRequests:       metric.TotalRequests,
			SuccessfulRequests:  metric.SuccessfulRequests,
			FailedRequests:      metric.FailedRequests,
			AverageResponseTime: metric.AverageResponseTime,
			RequestsPerSecond:   metric.RequestsPerSecond,
			ErrorRate:           metric.ErrorRate,
		}
	}
	return results
}

func (a *App) convertTestConfig(internal *testing.TestConfig) *TestConfig {
	return &TestConfig{
		Timeout:                  internal.Timeout,
		MaxConcurrentTests:       internal.MaxConcurrentTests,
		EnableParallelTesting:    internal.EnableParallelTesting,
		ContinueOnFailure:        internal.ContinueOnFailure,
		EnableSecurityScanning:   internal.EnableSecurityScanning,
		EnableLinting:            internal.EnableLinting,
		EnablePerformanceTesting: internal.EnablePerformanceTesting,
		EnableIntegrationTesting: internal.EnableIntegrationTesting,
		MCPProtocolVersion:       internal.MCPProtocolVersion,
		RequiredMethods:          internal.RequiredMethods,
		RequiredCapabilities:     internal.RequiredCapabilities,
		MaxResponseTime:          internal.MaxResponseTime,
		MaxMemoryUsage:           internal.MaxMemoryUsage,
		TestDataPath:             internal.TestDataPath,
		MCPClientPath:            internal.MCPClientPath,
		GenerateReport:           internal.GenerateReport,
		ReportFormat:             internal.ReportFormat,
		ReportOutputPath:         internal.ReportOutputPath,
		LogLevel:                 internal.LogLevel,
		RetryAttempts:            internal.RetryAttempts,
		RetryDelay:               internal.RetryDelay,
	}
}

func (a *App) applyTestConfigOptions(config *testing.TestConfig, options *TestConfigOptions) {
	if options.Timeout != nil {
		config.Timeout = *options.Timeout
	}
	if options.EnablePerformanceTesting != nil {
		config.EnablePerformanceTesting = *options.EnablePerformanceTesting
	}
	if options.EnableIntegrationTesting != nil {
		config.EnableIntegrationTesting = *options.EnableIntegrationTesting
	}
	if options.EnableSecurityScanning != nil {
		config.EnableSecurityScanning = *options.EnableSecurityScanning
	}
	if options.EnableLinting != nil {
		config.EnableLinting = *options.EnableLinting
	}
	if options.GenerateReport != nil {
		config.GenerateReport = *options.GenerateReport
	}
	if options.ReportFormat != nil {
		config.ReportFormat = *options.ReportFormat
	}
}

func (a *App) applyPerformanceTestOptions(config *testing.TestConfig, options *PerformanceTestOptions) {
	if options.MaxResponseTime != nil {
		config.MaxResponseTime = *options.MaxResponseTime
	}
	if options.MaxMemoryUsage != nil {
		config.MaxMemoryUsage = *options.MaxMemoryUsage
	}
}

func (a *App) validateTestConfig(config *TestConfig) error {
	if config.Timeout <= 0 {
		return a.createAPIError("validation", ErrCodeValidation, "Timeout must be positive", nil)
	}
	if config.MaxConcurrentTests <= 0 {
		return a.createAPIError("validation", ErrCodeValidation, "MaxConcurrentTests must be positive", nil)
	}
	if config.MaxResponseTime <= 0 {
		return a.createAPIError("validation", ErrCodeValidation, "MaxResponseTime must be positive", nil)
	}
	if config.MaxMemoryUsage <= 0 {
		return a.createAPIError("validation", ErrCodeValidation, "MaxMemoryUsage must be positive", nil)
	}

	return nil
}

// fileExists helper (reused from existing code)
func (a *App) fileExists(path string) error {
	if path == "" {
		return a.createAPIError("validation", ErrCodeValidation, "File path is required", nil)
	}

	// Check if the path is absolute, if not make it relative to current directory
	if !filepath.IsAbs(path) {
		return a.createAPIError("validation", ErrCodeValidation, "File path must be absolute", nil)
	}

	// Check if file exists
	if _, err := filepath.Stat(path); err != nil {
		return a.createAPIError("file_system", ErrCodeFileAccess, "File does not exist", map[string]string{
			"path": path,
		})
	}

	return nil
}