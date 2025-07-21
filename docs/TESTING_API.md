# Testing Framework API Reference

This document provides a complete API reference for the MCPWeaver testing framework, including all types, interfaces, and functions available for programmatic use.

## Table of Contents

1. [Core Types](#core-types)
2. [Configuration Types](#configuration-types)
3. [Result Types](#result-types)
4. [Validator Interfaces](#validator-interfaces)
5. [Testing Components](#testing-components)
6. [Pipeline Components](#pipeline-components)
7. [Reporting Components](#reporting-components)
8. [Diagnostics Components](#diagnostics-components)
9. [Factory Functions](#factory-functions)
10. [Constants and Enums](#constants-and-enums)

## Core Types

### TestConfig

Main configuration structure for the testing framework.

```go
type TestConfig struct {
    // Basic execution settings
    Timeout               time.Duration // Overall test timeout
    MaxConcurrentTests    int          // Maximum parallel tests
    EnableParallelTesting bool         // Enable parallel execution
    ContinueOnFailure     bool         // Continue after failures

    // Feature flags
    EnableSecurityScanning   bool // Enable security vulnerability scanning
    EnableLinting           bool // Enable code linting
    EnablePerformanceTesting bool // Enable performance tests
    EnableIntegrationTesting bool // Enable integration tests

    // MCP protocol settings
    MCPProtocolVersion   string   // Expected MCP protocol version
    RequiredMethods      []string // Required MCP methods to test
    RequiredCapabilities []string // Required MCP capabilities

    // Performance thresholds
    MaxResponseTime      time.Duration // Maximum acceptable response time
    MaxMemoryUsage       int64         // Maximum acceptable memory usage

    // External tools and paths
    TestDataPath      string // Path to test data directory
    MCPClientPath     string // Path to MCP client for integration tests

    // Report generation
    GenerateReport    bool   // Whether to generate test reports
    ReportFormat      string // Report format: "html", "json", "xml"
    ReportOutputPath  string // Output path for reports
    LogLevel         string // Logging level: "debug", "info", "warn", "error"

    // Retry configuration
    RetryAttempts     int           // Number of retry attempts for failed operations
    RetryDelay        time.Duration // Delay between retry attempts
}
```

### TestSuite

Main orchestrator for test execution.

```go
type TestSuite struct {
    // Private fields, access through methods
}

// Constructor
func NewTestSuite(config *TestConfig) *TestSuite

// Main testing method
func (ts *TestSuite) RunTests(ctx context.Context, serverPath string) (*TestResult, error)

// Individual test components
func (ts *TestSuite) RunValidation(ctx context.Context, serverPath string) (map[string]*ValidationResult, error)
func (ts *TestSuite) RunProtocolTests(ctx context.Context, serverPath string) (*ProtocolTestResult, error)
func (ts *TestSuite) RunIntegrationTests(ctx context.Context, serverPath string) (*IntegrationTestResult, error)
func (ts *TestSuite) RunPerformanceTests(ctx context.Context, serverPath string) (*PerformanceTestResult, error)

// Utility methods
func (ts *TestSuite) ValidatePrerequisites(ctx context.Context, serverPath string) error
func (ts *TestSuite) GetConfiguration() *TestConfig
```

## Configuration Types

### ConfigManager

Manages test configuration profiles.

```go
type ConfigManager struct {
    // Private fields
}

// Constructor
func NewConfigManager(configPath string) *ConfigManager

// Configuration management
func (cm *ConfigManager) LoadConfiguration() error
func (cm *ConfigManager) SaveConfiguration() error
func (cm *ConfigManager) GetCurrentConfig() *TestConfig
func (cm *ConfigManager) SetProfile(profileName string) error

// Profile management
func (cm *ConfigManager) CreateProfile(name, description string, config *TestConfig) error
func (cm *ConfigManager) DeleteProfile(name string) error
func (cm *ConfigManager) ListProfiles() map[string]*TestConfig
func (cm *ConfigManager) GetProfile(name string) (*TestConfig, error)
func (cm *ConfigManager) UpdateProfile(name string, config *TestConfig) error

// Built-in profiles
func (cm *ConfigManager) CreateBuiltinProfiles() error

// Validation and utilities
func (cm *ConfigManager) ValidateConfig(config *TestConfig) error
func (cm *ConfigManager) GetCurrentProfile() string
func (cm *ConfigManager) ExportProfile(profileName, outputPath string) error
func (cm *ConfigManager) ImportProfile(filePath string) error
func (cm *ConfigManager) GetConfigSummary() map[string]interface{}
```

### TestProfile

Represents a named configuration profile.

```go
type TestProfile struct {
    Name        string      `json:"name"`
    Description string      `json:"description"`
    Config      *TestConfig `json:"config"`
    CreatedAt   time.Time   `json:"createdAt"`
    UpdatedAt   time.Time   `json:"updatedAt"`
}
```

## Result Types

### TestResult

Complete test execution results.

```go
type TestResult struct {
    TestID            string        `json:"testId"`
    ServerPath        string        `json:"serverPath"`
    Timestamp         time.Time     `json:"timestamp"`
    Duration          time.Duration `json:"duration"`
    Success           bool          `json:"success"`
    TotalTests        int           `json:"totalTests"`
    PassedTests       int           `json:"passedTests"`
    FailedTests       int           `json:"failedTests"`
    SkippedTests      int           `json:"skippedTests"`

    // Detailed results by category
    ValidationResults   map[string]*ValidationResult `json:"validationResults,omitempty"`
    ProtocolResults     *ProtocolTestResult         `json:"protocolResults,omitempty"`
    IntegrationResults  *IntegrationTestResult      `json:"integrationResults,omitempty"`
    PerformanceResults  *PerformanceTestResult      `json:"performanceResults,omitempty"`

    // Summary information
    Errors              []string `json:"errors"`
    Warnings            []string `json:"warnings"`
    Recommendations     []string `json:"recommendations"`
}
```

### ValidationResult

Results from validation tests (compilation, syntax, linting, security).

```go
type ValidationResult struct {
    ValidatorName   string        `json:"validatorName"`
    Success         bool          `json:"success"`
    Duration        time.Duration `json:"duration"`
    FilesValidated  int           `json:"filesValidated"`
    Errors          []string      `json:"errors"`
    Warnings        []string      `json:"warnings"`
    Details         map[string]interface{} `json:"details,omitempty"`
}
```

### ProtocolTestResult

Results from MCP protocol compliance testing.

```go
type ProtocolTestResult struct {
    Success               bool                         `json:"success"`
    Duration              time.Duration                `json:"duration"`
    ProtocolVersion       string                       `json:"protocolVersion"`
    SupportedMethods      []string                     `json:"supportedMethods"`
    SupportedCapabilities []string                     `json:"supportedCapabilities"`
    MethodTests           map[string]*MethodTest       `json:"methodTests"`
    CapabilityTests       map[string]*CapabilityTest   `json:"capabilityTests"`
    Errors                []string                     `json:"errors"`
}

type MethodTest struct {
    Method        string        `json:"method"`
    Success       bool          `json:"success"`
    ResponseTime  time.Duration `json:"responseTime"`
    Request       interface{}   `json:"request"`
    Response      interface{}   `json:"response"`
    ErrorMessage  string        `json:"errorMessage,omitempty"`
}

type CapabilityTest struct {
    Capability    string                 `json:"capability"`
    Success       bool                   `json:"success"`
    Supported     bool                   `json:"supported"`
    TestDetails   map[string]interface{} `json:"testDetails,omitempty"`
    ErrorMessage  string                 `json:"errorMessage,omitempty"`
}
```

### IntegrationTestResult

Results from integration testing.

```go
type IntegrationTestResult struct {
    Success             bool                              `json:"success"`
    Duration            time.Duration                     `json:"duration"`
    ScenarioResults     map[string]*ScenarioTestResult    `json:"scenarioResults"`
    ClientCompatibility map[string]bool                   `json:"clientCompatibility"`
    Errors              []string                          `json:"errors"`
}

type ScenarioTestResult struct {
    Scenario     string       `json:"scenario"`
    Success      bool         `json:"success"`
    Duration     time.Duration `json:"duration"`
    Steps        []StepResult  `json:"steps"`
    ErrorMessage string        `json:"errorMessage,omitempty"`
}

type StepResult struct {
    Step         string                 `json:"step"`
    Success      bool                   `json:"success"`
    Duration     time.Duration          `json:"duration"`
    Details      map[string]interface{} `json:"details,omitempty"`
    ErrorMessage string                 `json:"errorMessage,omitempty"`
}
```

### PerformanceTestResult

Results from performance testing.

```go
type PerformanceTestResult struct {
    Success              bool                        `json:"success"`
    Duration             time.Duration               `json:"duration"`
    AverageResponseTime  time.Duration               `json:"averageResponseTime"`
    MedianResponseTime   time.Duration               `json:"medianResponseTime"`
    P95ResponseTime      time.Duration               `json:"p95ResponseTime"`
    P99ResponseTime      time.Duration               `json:"p99ResponseTime"`
    MaxResponseTime      time.Duration               `json:"maxResponseTime"`
    AverageMemoryUsage   int64                       `json:"averageMemoryUsage"`
    PeakMemoryUsage      int64                       `json:"peakMemoryUsage"`
    MemoryLeakDetected   bool                        `json:"memoryLeakDetected"`
    RequestsPerSecond    float64                     `json:"requestsPerSecond"`
    SuccessfulRequests   int                         `json:"successfulRequests"`
    FailedRequests       int                         `json:"failedRequests"`
    LoadTestResults      map[string]*LoadTestMetric  `json:"loadTestResults,omitempty"`
    Errors               []string                    `json:"errors"`
}

type LoadTestMetric struct {
    Scenario            string        `json:"scenario"`
    Duration            time.Duration `json:"duration"`
    TotalRequests       int           `json:"totalRequests"`
    SuccessfulRequests  int           `json:"successfulRequests"`
    FailedRequests      int           `json:"failedRequests"`
    AverageResponseTime time.Duration `json:"averageResponseTime"`
    RequestsPerSecond   float64       `json:"requestsPerSecond"`
    ErrorRate           float64       `json:"errorRate"`
}
```

## Validator Interfaces

### Validator

Base interface for all validators.

```go
type Validator interface {
    Name() string
    SupportsAsync() bool
    Validate(ctx context.Context, serverPath string) (*ValidationResult, error)
}
```

### Individual Validators

Each validator implements the Validator interface:

```go
// Compilation validation
type CompilationValidator struct { /* ... */ }
func NewCompilationValidator(config *TestConfig) *CompilationValidator

// Syntax validation
type SyntaxValidator struct { /* ... */ }
func NewSyntaxValidator(config *TestConfig) *SyntaxValidator

// Linting validation
type LintValidator struct { /* ... */ }
func NewLintValidator(config *TestConfig) *LintValidator

// Security validation
type SecurityValidator struct { /* ... */ }
func NewSecurityValidator(config *TestConfig) *SecurityValidator

// Dependency validation
type DependencyValidator struct { /* ... */ }
func NewDependencyValidator(config *TestConfig) *DependencyValidator
```

## Testing Components

### ProtocolTester

Tests MCP protocol compliance.

```go
type ProtocolTester struct { /* ... */ }

func NewProtocolTester(config *TestConfig) *ProtocolTester
func (pt *ProtocolTester) TestCompliance(ctx context.Context, serverPath string) (*ProtocolTestResult, error)
```

### IntegrationTester

Tests integration scenarios.

```go
type IntegrationTester struct { /* ... */ }

func NewIntegrationTester(config *TestConfig) *IntegrationTester
func (it *IntegrationTester) TestIntegration(ctx context.Context, serverPath string) (*IntegrationTestResult, error)
```

### PerformanceTester

Tests performance and load.

```go
type PerformanceTester struct { /* ... */ }

func NewPerformanceTester(config *TestConfig) *PerformanceTester
func (pt *PerformanceTester) TestPerformance(ctx context.Context, serverPath string) (*PerformanceTestResult, error)
```

### JSONRPCClient

JSON-RPC client for MCP communication.

```go
type JSONRPCClient struct { /* ... */ }

func NewJSONRPCClient(stdin io.WriteCloser, stdout io.ReadCloser, stderr io.ReadCloser) *JSONRPCClient
func (c *JSONRPCClient) Call(ctx context.Context, method string, params interface{}) (interface{}, error)
func (c *JSONRPCClient) Close() error
```

## Pipeline Components

### TestPipeline

Automated test execution pipeline.

```go
type TestPipeline struct { /* ... */ }

func NewTestPipeline(config *TestConfig) *TestPipeline
func (tp *TestPipeline) ExecutePipeline(ctx context.Context, serverPath string) (*PipelineResult, error)
func (tp *TestPipeline) IsRunning() bool
func (tp *TestPipeline) GetPipelineStatus() map[string]interface{}
```

### PipelineStage

Individual pipeline stage configuration.

```go
type PipelineStage struct {
    Name        string
    Description string
    Enabled     bool
    Parallel    bool
    Execute     func(ctx context.Context, serverPath string) error
    Validate    func(ctx context.Context, result *TestResult) error
    Timeout     time.Duration
    Retries     int
    OnFailure   string // "continue", "stop", "retry"
}
```

### PipelineResult

Results from pipeline execution.

```go
type PipelineResult struct {
    PipelineID      string                       `json:"pipelineId"`
    StartTime       time.Time                    `json:"startTime"`
    EndTime         time.Time                    `json:"endTime"`
    Duration        time.Duration                `json:"duration"`
    Success         bool                         `json:"success"`
    StageResults    map[string]*StageResult      `json:"stageResults"`
    TestResult      *TestResult                  `json:"testResult,omitempty"`
    Errors          []string                     `json:"errors"`
    TotalStages     int                          `json:"totalStages"`
    CompletedStages int                          `json:"completedStages"`
    SkippedStages   int                          `json:"skippedStages"`
    FailedStages    int                          `json:"failedStages"`
}

type StageResult struct {
    StageName    string                 `json:"stageName"`
    StartTime    time.Time              `json:"startTime"`
    EndTime      time.Time              `json:"endTime"`
    Duration     time.Duration          `json:"duration"`
    Success      bool                   `json:"success"`
    Skipped      bool                   `json:"skipped"`
    RetryCount   int                    `json:"retryCount"`
    ErrorMessage string                 `json:"errorMessage,omitempty"`
    Details      map[string]interface{} `json:"details,omitempty"`
}
```

### BatchTestRunner

Runs tests on multiple servers.

```go
type BatchTestRunner struct { /* ... */ }

func NewBatchTestRunner(config *TestConfig) *BatchTestRunner
func (btr *BatchTestRunner) RunBatchTests(ctx context.Context, request *BatchTestRequest) (*BatchTestResult, error)
```

### BatchTestRequest

Configuration for batch testing.

```go
type BatchTestRequest struct {
    RequestID     string   `json:"requestId"`
    ServerPaths   []string `json:"serverPaths"`
    Parallel      bool     `json:"parallel"`
    MaxWorkers    int      `json:"maxWorkers"`
    StopOnFailure bool     `json:"stopOnFailure"`
}
```

### BatchTestResult

Results from batch testing.

```go
type BatchTestResult struct {
    RequestID       string                      `json:"requestId"`
    StartTime       time.Time                   `json:"startTime"`
    EndTime         time.Time                   `json:"endTime"`
    Duration        time.Duration               `json:"duration"`
    Success         bool                        `json:"success"`
    TotalServers    int                         `json:"totalServers"`
    CompletedTests  int                         `json:"completedTests"`
    FailedTests     int                         `json:"failedTests"`
    SkippedTests    int                         `json:"skippedTests"`
    ServerResults   map[string]*PipelineResult  `json:"serverResults"`
    Errors          []string                    `json:"errors"`
    Summary         *BatchTestSummary           `json:"summary"`
}

type BatchTestSummary struct {
    AverageTestDuration    time.Duration            `json:"averageTestDuration"`
    FastestTest            time.Duration            `json:"fastestTest"`
    SlowestTest            time.Duration            `json:"slowestTest"`
    SuccessRate            float64                  `json:"successRate"`
    CommonFailures         map[string]int           `json:"commonFailures"`
    StageSuccessRates      map[string]float64       `json:"stageSuccessRates"`
    ResourceUsageSummary   map[string]interface{}   `json:"resourceUsageSummary"`
}
```

## Reporting Components

### TestReporter

Generates test reports in various formats.

```go
type TestReporter struct { /* ... */ }

func NewTestReporter(config *TestConfig) *TestReporter
func (tr *TestReporter) GenerateReport(result *TestResult) error
func (tr *TestReporter) GenerateMetricsReport(result *TestResult) (*TestMetrics, error)
```

### TestMetrics

Comprehensive test metrics and scores.

```go
type TestMetrics struct {
    TestID           string                        `json:"testId"`
    Timestamp        time.Time                     `json:"timestamp"`
    Duration         time.Duration                 `json:"duration"`
    OverallScore     float64                       `json:"overallScore"`
    QualityScore     float64                       `json:"qualityScore"`
    PerformanceScore float64                       `json:"performanceScore"`
    ComplianceScore  float64                       `json:"complianceScore"`
    Categories       map[string]*CategoryMetrics   `json:"categories"`
}

type CategoryMetrics struct {
    Score        float64                 `json:"score"`
    TotalTests   int                     `json:"totalTests"`
    PassedTests  int                     `json:"passedTests"`
    FailedTests  int                     `json:"failedTests"`
    Duration     time.Duration           `json:"duration"`
    Details      map[string]interface{}  `json:"details,omitempty"`
}
```

## Diagnostics Components

### DiagnosticsEngine

Analyzes test failures and provides debugging guidance.

```go
type DiagnosticsEngine struct { /* ... */ }

func NewDiagnosticsEngine(config *TestConfig) *DiagnosticsEngine
func (de *DiagnosticsEngine) AnalyzeFailures(ctx context.Context, serverPath string, testResult *TestResult, pipelineResult *PipelineResult) (*DiagnosticReport, error)
func (de *DiagnosticsEngine) SaveDiagnosticReport(report *DiagnosticReport, outputPath string) error
```

### DiagnosticReport

Comprehensive failure analysis and recommendations.

```go
type DiagnosticReport struct {
    ReportID          string                    `json:"reportId"`
    Timestamp         time.Time                 `json:"timestamp"`
    ServerPath        string                    `json:"serverPath"`
    FailureAnalysis   *FailureAnalysis          `json:"failureAnalysis"`
    EnvironmentInfo   *EnvironmentInfo          `json:"environmentInfo"`
    CodeAnalysis      *CodeAnalysis             `json:"codeAnalysis"`
    Dependencies      *DependencyAnalysis       `json:"dependencies"`
    Recommendations   []DiagnosticRecommendation `json:"recommendations"`
    TroubleshootingGuide *TroubleshootingGuide  `json:"troubleshootingGuide"`
    RelatedIssues     []RelatedIssue            `json:"relatedIssues"`
    Severity          string                    `json:"severity"`
    EstimatedFixTime  string                    `json:"estimatedFixTime"`
}
```

### FailureAnalysis

Detailed analysis of test failures.

```go
type FailureAnalysis struct {
    FailureType       string                 `json:"failureType"`
    FailureCategory   string                 `json:"failureCategory"`
    PrimaryError      string                 `json:"primaryError"`
    SecondaryErrors   []string               `json:"secondaryErrors"`
    FailedStages      []StageFailureInfo     `json:"failedStages"`
    ErrorPatterns     []ErrorPattern         `json:"errorPatterns"`
    RootCause         *RootCauseAnalysis     `json:"rootCause"`
    Impact            string                 `json:"impact"`
    Reproducibility   string                 `json:"reproducibility"`
}
```

### DiagnosticRecommendation

Actionable recommendations for fixing issues.

```go
type DiagnosticRecommendation struct {
    ID          string   `json:"id"`
    Priority    string   `json:"priority"`
    Category    string   `json:"category"`
    Title       string   `json:"title"`
    Description string   `json:"description"`
    Actions     []Action `json:"actions"`
    Resources   []string `json:"resources"`
    EstimatedTime string `json:"estimatedTime"`
}

type Action struct {
    Step        int    `json:"step"`
    Description string `json:"description"`
    Command     string `json:"command,omitempty"`
    Example     string `json:"example,omitempty"`
    Warning     string `json:"warning,omitempty"`
}
```

## Factory Functions

### Core Components

```go
// Main testing components
func NewTestSuite(config *TestConfig) *TestSuite
func NewTestPipeline(config *TestConfig) *TestPipeline
func NewConfigManager(configPath string) *ConfigManager
func NewDiagnosticsEngine(config *TestConfig) *DiagnosticsEngine

// Reporting
func NewTestReporter(config *TestConfig) *TestReporter

// Batch testing
func NewBatchTestRunner(config *TestConfig) *BatchTestRunner
```

### Validators

```go
func NewCompilationValidator(config *TestConfig) *CompilationValidator
func NewSyntaxValidator(config *TestConfig) *SyntaxValidator
func NewLintValidator(config *TestConfig) *LintValidator
func NewSecurityValidator(config *TestConfig) *SecurityValidator
func NewDependencyValidator(config *TestConfig) *DependencyValidator
```

### Testers

```go
func NewProtocolTester(config *TestConfig) *ProtocolTester
func NewIntegrationTester(config *TestConfig) *IntegrationTester
func NewPerformanceTester(config *TestConfig) *PerformanceTester
```

### Utilities

```go
func NewJSONRPCClient(stdin io.WriteCloser, stdout io.ReadCloser, stderr io.ReadCloser) *JSONRPCClient
func NewProtocolConformanceValidator(config *TestConfig) *ProtocolConformanceValidator
```

## Constants and Enums

### Default Values

```go
const (
    DefaultTimeout            = 5 * time.Minute
    DefaultMaxConcurrentTests = 3
    DefaultMaxResponseTime    = time.Second
    DefaultMaxMemoryUsage     = 100 * 1024 * 1024 // 100MB
    DefaultRetryAttempts      = 2
    DefaultRetryDelay         = time.Second
)
```

### Report Formats

```go
const (
    ReportFormatHTML = "html"
    ReportFormatJSON = "json" 
    ReportFormatXML  = "xml"
)
```

### Log Levels

```go
const (
    LogLevelDebug = "debug"
    LogLevelInfo  = "info"
    LogLevelWarn  = "warn"
    LogLevelError = "error"
)
```

### Built-in Profiles

```go
const (
    ProfileDefault     = "default"
    ProfileFast        = "fast"
    ProfileThorough    = "thorough"
    ProfileDevelopment = "development"
    ProfileCI          = "ci"
    ProfileSecurity    = "security"
    ProfilePerformance = "performance"
)
```

### Test Categories

```go
const (
    CategoryValidation  = "validation"
    CategoryProtocol    = "protocol"
    CategoryIntegration = "integration"
    CategoryPerformance = "performance"
)
```

### Failure Types

```go
const (
    FailureTypeCompilation = "compilation_error"
    FailureTypeSyntax      = "syntax_error"
    FailureTypeTimeout     = "timeout_error"
    FailureTypeNetwork     = "network_error"
    FailureTypePermission  = "permission_error"
    FailureTypeMemory      = "memory_error"
    FailureTypeDependency  = "dependency_error"
    FailureTypeProtocol    = "protocol_error"
    FailureTypeUnknown     = "unknown_error"
)
```

### Severity Levels

```go
const (
    SeverityCritical = "critical"
    SeverityMajor    = "major"
    SeverityMinor    = "minor"
)
```

This API reference provides complete documentation for all public types, interfaces, and functions in the MCPWeaver testing framework. Use this reference when integrating the testing framework into your applications or when extending the framework with custom validators or testers.