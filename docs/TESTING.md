# MCP Server Testing Framework

This document provides comprehensive documentation for the MCPWeaver testing framework, which enables thorough testing of generated MCP servers.

## Table of Contents

1. [Overview](#overview)
2. [Getting Started](#getting-started)
3. [Test Types](#test-types)
4. [Configuration](#configuration)
5. [Running Tests](#running-tests)
6. [Test Results and Reporting](#test-results-and-reporting)
7. [Debugging and Diagnostics](#debugging-and-diagnostics)
8. [API Reference](#api-reference)
9. [Best Practices](#best-practices)
10. [Troubleshooting](#troubleshooting)

## Overview

The MCPWeaver testing framework provides comprehensive testing capabilities for generated MCP servers, including:

- **Server Validation**: Compilation, syntax, and structural validation
- **Protocol Compliance**: MCP protocol conformance testing
- **Integration Testing**: Real client compatibility testing
- **Performance Testing**: Response time, memory usage, and load testing
- **Security Scanning**: Vulnerability and security issue detection
- **Automated Pipeline**: Configurable test execution workflows
- **Detailed Reporting**: Comprehensive test reports with metrics
- **Failure Diagnostics**: Advanced debugging and troubleshooting

### Key Features

- ✅ **Comprehensive Coverage**: Tests all aspects of MCP server functionality
- ✅ **Flexible Configuration**: Multiple test profiles for different scenarios
- ✅ **Parallel Execution**: Efficient parallel test execution
- ✅ **Rich Reporting**: HTML, JSON, and XML report formats
- ✅ **Failure Analysis**: Intelligent error diagnosis and recommendations
- ✅ **Pipeline Automation**: Automated test workflows with retry logic
- ✅ **Performance Monitoring**: Memory leak detection and performance metrics

## Getting Started

### Prerequisites

Before using the testing framework, ensure you have:

1. **Go 1.21+** installed and available in PATH
2. **Generated MCP server** with `main.go` and `go.mod` files
3. **Optional tools** for enhanced testing:
   - `golangci-lint` for advanced linting
   - `gosec` for security scanning
   - `govulncheck` for vulnerability scanning

### Basic Usage

```go
package main

import (
    "context"
    "log"
    "github.com/matoval/MCPWeaver/internal/testing"
)

func main() {
    // Create test configuration
    config := &testing.TestConfig{
        Timeout:               5 * time.Minute,
        EnableSecurityScanning: true,
        EnableLinting:         true,
        GenerateReport:        true,
        ReportFormat:          "html",
    }

    // Create test suite
    suite := testing.NewTestSuite(config)

    // Run tests
    ctx := context.Background()
    result, err := suite.RunTests(ctx, "/path/to/generated/server")
    if err != nil {
        log.Fatalf("Test execution failed: %v", err)
    }

    if result.Success {
        log.Println("All tests passed!")
    } else {
        log.Printf("Tests failed: %v", result.Errors)
    }
}
```

## Test Types

### 1. Server Validation

Validates that the generated server code is correct and compilable.

#### Compilation Validation
- Checks if the server compiles successfully
- Validates `main.go` and `go.mod` files exist
- Ensures all dependencies are available

#### Syntax Validation  
- Validates Go syntax using AST parsing
- Checks for required imports and package structure
- Verifies main function and MCP types are present

#### Linting
- Runs `go fmt`, `go vet`, and `golangci-lint`
- Checks code quality and style issues
- Provides warnings for best practice violations

#### Security Scanning
- Scans for common security vulnerabilities
- Uses `gosec` for advanced security analysis
- Checks for potential password/secret exposure

#### Dependency Analysis
- Validates module dependencies with `go mod verify`
- Scans for known vulnerabilities with `govulncheck`
- Checks for outdated or conflicting dependencies

### 2. Protocol Compliance Testing

Ensures the server conforms to MCP protocol specifications.

```go
// Test protocol compliance
protocolTester := testing.NewProtocolTester(config)
result, err := protocolTester.TestCompliance(ctx, serverPath)
```

Features:
- **Initialization Testing**: Tests the MCP handshake process
- **Method Support**: Validates required methods (`initialize`, `tools/list`, `tools/call`)
- **Response Format**: Ensures JSON-RPC 2.0 compliance
- **Error Handling**: Tests error response handling
- **Capability Detection**: Validates declared capabilities

### 3. Integration Testing

Tests real-world integration scenarios with MCP clients.

```go
// Test integration scenarios
integrationTester := testing.NewIntegrationTester(config)
result, err := integrationTester.TestIntegration(ctx, serverPath)
```

Test scenarios:
- **Server Startup**: Validates server starts correctly
- **Client Connection**: Tests client-server handshake
- **Tools Discovery**: Tests `tools/list` functionality  
- **Tool Execution**: Tests `tools/call` with various parameters
- **Error Handling**: Tests invalid requests and error responses
- **Concurrent Requests**: Tests handling multiple simultaneous requests

### 4. Performance Testing

Measures server performance under various conditions.

```go
// Test performance
perfTester := testing.NewPerformanceTester(config)
result, err := perfTester.TestPerformance(ctx, serverPath)
```

Performance metrics:
- **Response Time**: Average, median, P95, P99, and maximum response times
- **Memory Usage**: Peak memory consumption and leak detection
- **Throughput**: Requests per second under normal load
- **Load Testing**: Performance under light, medium, and heavy loads

### 5. Pipeline Testing

Automated test execution with configurable stages.

```go
// Create and run test pipeline
pipeline := testing.NewTestPipeline(config)
result, err := pipeline.ExecutePipeline(ctx, serverPath)
```

Pipeline stages:
1. **Pre-validation**: Environment and prerequisite checks
2. **Dependency Check**: Dependency resolution and installation
3. **Compilation Validation**: Code compilation verification
4. **Syntax Validation**: Code syntax and structure validation
5. **Security Scan**: Security vulnerability scanning
6. **Lint Check**: Code quality and style checking
7. **Environment Setup**: Test environment preparation
8. **Test Execution**: Comprehensive test suite execution

## Configuration

### Test Configuration

The framework uses a flexible configuration system with support for multiple profiles:

```go
type TestConfig struct {
    // Basic settings
    Timeout               time.Duration
    MaxConcurrentTests    int
    EnableParallelTesting bool
    ContinueOnFailure     bool

    // Feature flags
    EnableSecurityScanning   bool
    EnableLinting           bool
    EnablePerformanceTesting bool
    EnableIntegrationTesting bool

    // MCP settings
    MCPProtocolVersion   string
    RequiredMethods      []string
    RequiredCapabilities []string

    // Performance limits
    MaxResponseTime      time.Duration
    MaxMemoryUsage       int64

    // Paths and tools
    TestDataPath      string
    MCPClientPath     string

    // Reporting
    GenerateReport    bool
    ReportFormat      string
    ReportOutputPath  string
    LogLevel         string

    // Retry behavior
    RetryAttempts    int
    RetryDelay       time.Duration
}
```

### Configuration Profiles

The framework includes several built-in profiles:

#### Default Profile
```go
config := testing.GetDefaultTestConfig()
```
- Balanced settings for general testing
- 5-minute timeout, parallel execution enabled
- All validation types enabled
- HTML reporting

#### Fast Profile  
```go
configManager := testing.NewConfigManager("config.json")
config, _ := configManager.GetProfile("fast")
```
- Quick validation for development
- 2-minute timeout, reduced test scope
- Security and performance testing disabled
- JSON reporting

#### Thorough Profile
```go
config, _ := configManager.GetProfile("thorough")
```
- Comprehensive testing with strict limits
- 15-minute timeout, all validations enabled
- Enhanced security scanning
- Detailed HTML reporting

#### CI Profile
```go
config, _ := configManager.GetProfile("ci")
```
- Optimized for continuous integration
- 8-minute timeout, parallel execution
- XML reporting for CI integration
- Minimal retries for fast feedback

### Custom Configuration

Create custom profiles for specific needs:

```go
// Create configuration manager
configManager := testing.NewConfigManager("test-config.json")

// Create custom profile
customConfig := &testing.TestConfig{
    Timeout:                  10 * time.Minute,
    MaxConcurrentTests:       2,
    EnableParallelTesting:    true,
    EnableSecurityScanning:   true,
    EnableLinting:           false,
    EnablePerformanceTesting: true,
    MaxResponseTime:         500 * time.Millisecond,
    MaxMemoryUsage:          50 * 1024 * 1024, // 50MB
    ReportFormat:            "html",
    LogLevel:                "info",
}

// Save custom profile
err := configManager.CreateProfile("custom", "Custom testing profile", customConfig)
err = configManager.SetProfile("custom")
err = configManager.SaveConfiguration()
```

## Running Tests

### Single Server Testing

```go
// Basic test execution
suite := testing.NewTestSuite(config)
result, err := suite.RunTests(ctx, "/path/to/server")

// Check results
if result.Success {
    fmt.Printf("All tests passed! Duration: %v\n", result.Duration)
} else {
    fmt.Printf("Tests failed: %d errors\n", len(result.Errors))
    for _, err := range result.Errors {
        fmt.Printf("- %s\n", err)
    }
}
```

### Batch Testing

Test multiple servers simultaneously:

```go
// Create batch test runner
batchRunner := testing.NewBatchTestRunner(config)

// Define batch request
request := &testing.BatchTestRequest{
    RequestID:     "batch_001",
    ServerPaths:   []string{"/server1", "/server2", "/server3"},
    Parallel:      true,
    MaxWorkers:    3,
    StopOnFailure: false,
}

// Run batch tests
result, err := batchRunner.RunBatchTests(ctx, request)

// Analyze results
fmt.Printf("Batch test completed: %d/%d servers passed\n", 
    result.CompletedTests, result.TotalServers)
fmt.Printf("Success rate: %.1f%%\n", result.Summary.SuccessRate)
```

### Pipeline Execution

Use the automated pipeline for comprehensive testing:

```go
// Create pipeline
pipeline := testing.NewTestPipeline(config)

// Check pipeline status
status := pipeline.GetPipelineStatus()
fmt.Printf("Pipeline has %d stages, %d enabled\n", 
    status["totalStages"], status["enabledStages"])

// Execute pipeline
result, err := pipeline.ExecutePipeline(ctx, serverPath)

// Analyze stage results
for stageName, stageResult := range result.StageResults {
    if stageResult.Success {
        fmt.Printf("✅ %s completed in %v\n", stageName, stageResult.Duration)
    } else {
        fmt.Printf("❌ %s failed: %s\n", stageName, stageResult.ErrorMessage)
    }
}
```

## Test Results and Reporting

### Test Result Structure

```go
type TestResult struct {
    TestID            string
    ServerPath        string
    Timestamp         time.Time
    Duration          time.Duration
    Success           bool
    TotalTests        int
    PassedTests       int
    FailedTests       int
    SkippedTests      int
    
    // Detailed results
    ValidationResults   map[string]*ValidationResult
    ProtocolResults     *ProtocolTestResult
    IntegrationResults  *IntegrationTestResult
    PerformanceResults  *PerformanceTestResult
    
    Errors              []string
    Warnings            []string
    Recommendations     []string
}
```

### Report Generation

The framework supports multiple report formats:

#### HTML Reports (Default)
```go
config.ReportFormat = "html"
config.ReportOutputPath = "test-report.html"
```
- Rich visual presentation with charts and metrics
- Interactive sections for detailed analysis
- CSS styling for professional appearance
- Embedded test metrics and recommendations

#### JSON Reports
```go
config.ReportFormat = "json"
config.ReportOutputPath = "test-report.json"
```
- Machine-readable format for automation
- Complete test data for programmatic analysis
- Integration with CI/CD pipelines
- API consumption and data processing

#### XML Reports
```go
config.ReportFormat = "xml"
config.ReportOutputPath = "test-report.xml"
```
- Compatible with CI systems (Jenkins, GitLab CI)
- Standard format for test reporting tools
- Integration with test management systems

### Metrics and Analytics

Generate detailed metrics:

```go
// Generate metrics report
reporter := testing.NewTestReporter(config)
metrics, err := reporter.GenerateMetricsReport(testResult)

// Access metrics
fmt.Printf("Overall Score: %.1f\n", metrics.OverallScore)
fmt.Printf("Quality Score: %.1f\n", metrics.QualityScore)
fmt.Printf("Performance Score: %.1f\n", metrics.PerformanceScore)
fmt.Printf("Compliance Score: %.1f\n", metrics.ComplianceScore)

// Category breakdown
for category, categoryMetrics := range metrics.Categories {
    fmt.Printf("%s: %.1f%% success rate\n", 
        category, float64(categoryMetrics.PassedTests)/float64(categoryMetrics.TotalTests)*100)
}
```

## Debugging and Diagnostics

### Failure Analysis

When tests fail, use the diagnostics engine for detailed analysis:

```go
// Create diagnostics engine
diagnostics := testing.NewDiagnosticsEngine(config)

// Analyze failures
report, err := diagnostics.AnalyzeFailures(ctx, serverPath, testResult, pipelineResult)

// Review analysis
fmt.Printf("Failure Type: %s\n", report.FailureAnalysis.FailureType)
fmt.Printf("Root Cause: %s (confidence: %.0f%%)\n", 
    report.FailureAnalysis.RootCause.ProbableCause,
    report.FailureAnalysis.RootCause.ConfidenceLevel*100)

// Get recommendations
for _, rec := range report.Recommendations {
    fmt.Printf("\n%s (Priority: %s)\n", rec.Title, rec.Priority)
    fmt.Printf("%s\n", rec.Description)
    for _, action := range rec.Actions {
        fmt.Printf("  %d. %s\n", action.Step, action.Description)
    }
}
```

### Troubleshooting Guide

Access step-by-step troubleshooting:

```go
// Get troubleshooting guide
guide := report.TroubleshootingGuide

// Try quick fixes first
for _, quickFix := range guide.QuickFixes {
    fmt.Printf("Quick Fix: %s\n", quickFix.Name)
    fmt.Printf("Command: %s\n", quickFix.Command)
    fmt.Printf("Expected: %s\n", quickFix.Expected)
}

// Follow detailed steps if needed
for _, step := range guide.DetailedSteps {
    fmt.Printf("\nStep %d: %s\n", step.Step, step.Title)
    fmt.Printf("%s\n", step.Description)
    for _, cmd := range step.Commands {
        fmt.Printf("  $ %s\n", cmd)
    }
}
```

### Common Issues and Solutions

#### Compilation Errors
```
Error: cannot find package "github.com/sourcegraph/jsonrpc2"
Solution: Run 'go mod tidy' to resolve dependencies
```

#### Timeout Issues
```
Error: test timeout after 5 minutes
Solutions:
1. Increase timeout in configuration
2. Optimize server performance
3. Check system resources
```

#### Protocol Compliance Failures
```
Error: missing required method 'tools/list'
Solution: Ensure all required MCP methods are implemented
```

## API Reference

### Core Types

#### TestConfig
Configuration for test execution with all available options.

#### TestSuite
Main testing orchestrator that coordinates all test types.

#### TestResult  
Complete test results with detailed breakdown by category.

#### ValidationResult
Results from code validation (compilation, syntax, linting, security).

#### ProtocolTestResult
Results from MCP protocol compliance testing.

#### IntegrationTestResult
Results from integration testing with scenarios and steps.

#### PerformanceTestResult
Results from performance testing with metrics and load tests.

### Key Interfaces

#### Validator
```go
type Validator interface {
    Name() string
    SupportsAsync() bool
    Validate(ctx context.Context, serverPath string) (*ValidationResult, error)
}
```

#### TestReporter
```go
type TestReporter interface {
    GenerateReport(result *TestResult) error
    GenerateMetricsReport(result *TestResult) (*TestMetrics, error)
}
```

### Factory Functions

```go
// Create components
func NewTestSuite(config *TestConfig) *TestSuite
func NewTestPipeline(config *TestConfig) *TestPipeline
func NewConfigManager(configPath string) *ConfigManager
func NewDiagnosticsEngine(config *TestConfig) *DiagnosticsEngine

// Create validators
func NewCompilationValidator(config *TestConfig) *CompilationValidator
func NewSyntaxValidator(config *TestConfig) *SyntaxValidator
func NewLintValidator(config *TestConfig) *LintValidator
func NewSecurityValidator(config *TestConfig) *SecurityValidator
func NewDependencyValidator(config *TestConfig) *DependencyValidator

// Create testers
func NewProtocolTester(config *TestConfig) *ProtocolTester
func NewIntegrationTester(config *TestConfig) *IntegrationTester
func NewPerformanceTester(config *TestConfig) *PerformanceTester
```

## Best Practices

### 1. Configuration Management

- Use appropriate profiles for different scenarios
- Create custom profiles for specific requirements
- Store configuration in version control
- Document configuration choices

### 2. Test Organization

- Run fast tests first (syntax, compilation)
- Use parallel execution for independent tests
- Group related tests in pipeline stages
- Implement proper retry logic for flaky tests

### 3. Performance Testing

- Set realistic performance thresholds
- Test under expected load conditions
- Monitor memory usage and leaks
- Use consistent test environments

### 4. Error Handling

- Use the diagnostics engine for failure analysis
- Review troubleshooting guides before manual debugging
- Document common issues and solutions
- Implement proper logging for debugging

### 5. Reporting

- Generate reports appropriate for your audience
- Use HTML reports for human review
- Use JSON/XML reports for automation
- Archive reports for historical analysis

### 6. CI/CD Integration

- Use the CI profile for automated pipelines
- Set appropriate timeouts for CI environments
- Generate machine-readable reports
- Implement proper failure notifications

## Troubleshooting

### Common Issues

#### Go Not Found
```
Error: exec: "go": executable file not found in $PATH
Solution: Install Go and ensure it's in your PATH
```

#### Module Issues
```
Error: go.mod file not found
Solution: Ensure the server directory contains a valid go.mod file
```

#### Permission Errors
```
Error: permission denied when writing report
Solution: Check file permissions and disk space
```

#### Memory Issues
```
Error: cannot allocate memory
Solution: Increase available memory or reduce concurrent tests
```

#### Network Timeouts
```
Error: timeout downloading dependencies
Solution: Check network connectivity and proxy settings
```

### Getting Help

1. **Check the diagnostic report** for specific guidance
2. **Review troubleshooting guides** generated by the diagnostics engine
3. **Consult the logs** with appropriate log level (`debug`, `info`, `warn`, `error`)
4. **Verify environment** requirements and tool availability
5. **Test with different profiles** to isolate issues

### Debug Mode

Enable detailed logging for troubleshooting:

```go
config.LogLevel = "debug"
```

This will provide verbose output including:
- Detailed test execution steps
- Command output and errors
- Performance metrics during execution
- Environment information
- Configuration values used

For additional support, refer to the project documentation and issue tracker.