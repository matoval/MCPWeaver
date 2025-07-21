package testing

import (
	"encoding/json"
	"fmt"
	"html/template"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// TestReporter handles test result reporting and metrics
type TestReporter struct {
	config *TestConfig
}

// NewTestReporter creates a new test reporter
func NewTestReporter(config *TestConfig) *TestReporter {
	return &TestReporter{
		config: config,
	}
}

// GenerateReport generates a comprehensive test report
func (tr *TestReporter) GenerateReport(result *TestResult) error {
	if !tr.config.GenerateReport {
		return nil
	}

	switch tr.config.ReportFormat {
	case "json":
		return tr.generateJSONReport(result)
	case "html":
		return tr.generateHTMLReport(result)
	case "xml":
		return tr.generateXMLReport(result)
	default:
		return tr.generateJSONReport(result)
	}
}

// generateJSONReport generates a JSON format report
func (tr *TestReporter) generateJSONReport(result *TestResult) error {
	reportData := tr.prepareReportData(result)
	
	jsonData, err := json.MarshalIndent(reportData, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal report data: %w", err)
	}

	outputPath := tr.getReportPath("json")
	if err := os.WriteFile(outputPath, jsonData, 0644); err != nil {
		return fmt.Errorf("failed to write JSON report: %w", err)
	}

	return nil
}

// generateHTMLReport generates an HTML format report
func (tr *TestReporter) generateHTMLReport(result *TestResult) error {
	reportData := tr.prepareReportData(result)
	
	htmlTemplate := `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>MCP Server Test Report</title>
    <style>
        body { font-family: Arial, sans-serif; margin: 20px; background-color: #f5f5f5; }
        .container { max-width: 1200px; margin: 0 auto; background: white; padding: 20px; border-radius: 8px; box-shadow: 0 2px 4px rgba(0,0,0,0.1); }
        .header { border-bottom: 2px solid #007acc; padding-bottom: 20px; margin-bottom: 30px; }
        .header h1 { color: #007acc; margin: 0; }
        .header .subtitle { color: #666; margin-top: 5px; }
        .summary { display: grid; grid-template-columns: repeat(auto-fit, minmax(200px, 1fr)); gap: 20px; margin-bottom: 30px; }
        .summary-card { background: #f8f9fa; padding: 20px; border-radius: 6px; border-left: 4px solid #007acc; }
        .summary-card h3 { margin: 0 0 10px 0; color: #333; font-size: 14px; text-transform: uppercase; }
        .summary-card .value { font-size: 24px; font-weight: bold; color: #007acc; }
        .summary-card .label { font-size: 12px; color: #666; margin-top: 5px; }
        .success { border-left-color: #28a745; }
        .success .value { color: #28a745; }
        .failure { border-left-color: #dc3545; }
        .failure .value { color: #dc3545; }
        .warning { border-left-color: #ffc107; }
        .warning .value { color: #e65100; }
        .section { margin-bottom: 30px; }
        .section h2 { color: #333; border-bottom: 1px solid #ddd; padding-bottom: 10px; }
        .test-grid { display: grid; gap: 15px; }
        .test-item { background: #f8f9fa; padding: 15px; border-radius: 6px; border-left: 4px solid #ddd; }
        .test-item.success { border-left-color: #28a745; }
        .test-item.failure { border-left-color: #dc3545; }
        .test-item.warning { border-left-color: #ffc107; }
        .test-name { font-weight: bold; margin-bottom: 5px; }
        .test-duration { font-size: 12px; color: #666; }
        .test-details { margin-top: 10px; font-size: 14px; }
        .error-list { background: #ffe6e6; padding: 10px; border-radius: 4px; margin-top: 10px; }
        .error-list ul { margin: 0; padding-left: 20px; }
        .warning-list { background: #fff3cd; padding: 10px; border-radius: 4px; margin-top: 10px; }
        .performance-metrics { display: grid; grid-template-columns: repeat(auto-fit, minmax(150px, 1fr)); gap: 15px; }
        .metric { text-align: center; padding: 15px; background: #f8f9fa; border-radius: 6px; }
        .metric .label { font-size: 12px; color: #666; text-transform: uppercase; }
        .metric .value { font-size: 18px; font-weight: bold; color: #007acc; margin-top: 5px; }
        .recommendations { background: #e7f3ff; padding: 15px; border-radius: 6px; border-left: 4px solid #007acc; }
        .recommendations h3 { margin-top: 0; color: #007acc; }
        .recommendations ul { margin: 10px 0; padding-left: 20px; }
        .load-test-results { display: grid; gap: 15px; }
        .load-test { background: #f8f9fa; padding: 15px; border-radius: 6px; }
        .load-test h4 { margin: 0 0 10px 0; color: #333; }
        .load-metrics { display: grid; grid-template-columns: repeat(auto-fit, minmax(120px, 1fr)); gap: 10px; }
        .load-metric { text-align: center; }
        .load-metric .label { font-size: 11px; color: #666; }
        .load-metric .value { font-size: 14px; font-weight: bold; color: #007acc; }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h1>MCP Server Test Report</h1>
            <div class="subtitle">{{.ServerPath}} - {{.Timestamp.Format "January 2, 2006 15:04:05"}}</div>
        </div>

        <div class="summary">
            <div class="summary-card {{if .Success}}success{{else}}failure{{end}}">
                <h3>Overall Status</h3>
                <div class="value">{{if .Success}}PASS{{else}}FAIL{{end}}</div>
                <div class="label">Test Result</div>
            </div>
            <div class="summary-card">
                <h3>Total Tests</h3>
                <div class="value">{{.TotalTests}}</div>
                <div class="label">Test Cases</div>
            </div>
            <div class="summary-card success">
                <h3>Passed</h3>
                <div class="value">{{.PassedTests}}</div>
                <div class="label">Successful Tests</div>
            </div>
            <div class="summary-card failure">
                <h3>Failed</h3>
                <div class="value">{{.FailedTests}}</div>
                <div class="label">Failed Tests</div>
            </div>
            <div class="summary-card">
                <h3>Duration</h3>
                <div class="value">{{.Duration.Truncate 1000000}}</div>
                <div class="label">Total Time</div>
            </div>
        </div>

        {{if .ValidationResults}}
        <div class="section">
            <h2>Validation Results</h2>
            <div class="test-grid">
                {{range $name, $result := .ValidationResults}}
                <div class="test-item {{if $result.Success}}success{{else}}failure{{end}}">
                    <div class="test-name">{{$name}} Validation</div>
                    <div class="test-duration">Duration: {{$result.Duration.Truncate 1000000}}</div>
                    <div class="test-details">Files Validated: {{$result.FilesValidated}}</div>
                    {{if $result.Errors}}
                    <div class="error-list">
                        <strong>Errors:</strong>
                        <ul>{{range $result.Errors}}<li>{{.}}</li>{{end}}</ul>
                    </div>
                    {{end}}
                    {{if $result.Warnings}}
                    <div class="warning-list">
                        <strong>Warnings:</strong>
                        <ul>{{range $result.Warnings}}<li>{{.}}</li>{{end}}</ul>
                    </div>
                    {{end}}
                </div>
                {{end}}
            </div>
        </div>
        {{end}}

        {{if .ProtocolResults}}
        <div class="section">
            <h2>Protocol Compliance</h2>
            <div class="test-item {{if .ProtocolResults.Success}}success{{else}}failure{{end}}">
                <div class="test-name">MCP Protocol Compliance</div>
                <div class="test-duration">Duration: {{.ProtocolResults.Duration.Truncate 1000000}}</div>
                <div class="test-details">
                    <div>Protocol Version: {{.ProtocolResults.ProtocolVersion}}</div>
                    <div>Supported Methods: {{len .ProtocolResults.SupportedMethods}}</div>
                    <div>Supported Capabilities: {{len .ProtocolResults.SupportedCapabilities}}</div>
                </div>
                {{if .ProtocolResults.Errors}}
                <div class="error-list">
                    <strong>Protocol Errors:</strong>
                    <ul>{{range .ProtocolResults.Errors}}<li>{{.}}</li>{{end}}</ul>
                </div>
                {{end}}
            </div>
        </div>
        {{end}}

        {{if .PerformanceResults}}
        <div class="section">
            <h2>Performance Metrics</h2>
            <div class="performance-metrics">
                <div class="metric">
                    <div class="label">Avg Response Time</div>
                    <div class="value">{{.PerformanceResults.AverageResponseTime.Truncate 1000000}}</div>
                </div>
                <div class="metric">
                    <div class="label">P95 Response Time</div>
                    <div class="value">{{.PerformanceResults.P95ResponseTime.Truncate 1000000}}</div>
                </div>
                <div class="metric">
                    <div class="label">Max Response Time</div>
                    <div class="value">{{.PerformanceResults.MaxResponseTime.Truncate 1000000}}</div>
                </div>
                <div class="metric">
                    <div class="label">Peak Memory</div>
                    <div class="value">{{.PerformanceResults.PeakMemoryUsage | formatBytes}}</div>
                </div>
                <div class="metric">
                    <div class="label">Requests/sec</div>
                    <div class="value">{{printf "%.1f" .PerformanceResults.RequestsPerSecond}}</div>
                </div>
                <div class="metric">
                    <div class="label">Success Rate</div>
                    <div class="value">{{printf "%.1f%%" (div (mul .PerformanceResults.SuccessfulRequests 100.0) (add .PerformanceResults.SuccessfulRequests .PerformanceResults.FailedRequests))}}</div>
                </div>
            </div>

            {{if .PerformanceResults.LoadTestResults}}
            <div class="section">
                <h3>Load Test Results</h3>
                <div class="load-test-results">
                    {{range $name, $result := .PerformanceResults.LoadTestResults}}
                    <div class="load-test">
                        <h4>{{$name | title}}</h4>
                        <div class="load-metrics">
                            <div class="load-metric">
                                <div class="label">Duration</div>
                                <div class="value">{{$result.Duration.Truncate 1000000}}</div>
                            </div>
                            <div class="load-metric">
                                <div class="label">Total Requests</div>
                                <div class="value">{{$result.TotalRequests}}</div>
                            </div>
                            <div class="load-metric">
                                <div class="label">Success Rate</div>
                                <div class="value">{{printf "%.1f%%" (mul (sub 1.0 $result.ErrorRate) 100.0)}}</div>
                            </div>
                            <div class="load-metric">
                                <div class="label">Avg Response</div>
                                <div class="value">{{$result.AverageResponseTime.Truncate 1000000}}</div>
                            </div>
                            <div class="load-metric">
                                <div class="label">RPS</div>
                                <div class="value">{{printf "%.1f" $result.RequestsPerSecond}}</div>
                            </div>
                        </div>
                    </div>
                    {{end}}
                </div>
            </div>
            {{end}}
        </div>
        {{end}}

        {{if .Recommendations}}
        <div class="section">
            <div class="recommendations">
                <h3>Recommendations</h3>
                <ul>
                    {{range .Recommendations}}<li>{{.}}</li>{{end}}
                </ul>
            </div>
        </div>
        {{end}}

        {{if .Errors}}
        <div class="section">
            <h2>Errors</h2>
            <div class="error-list">
                <ul>
                    {{range .Errors}}<li>{{.}}</li>{{end}}
                </ul>
            </div>
        </div>
        {{end}}
    </div>
</body>
</html>`

	tmpl, err := template.New("report").Funcs(template.FuncMap{
		"formatBytes": func(bytes int64) string {
			if bytes < 1024 {
				return fmt.Sprintf("%d B", bytes)
			} else if bytes < 1024*1024 {
				return fmt.Sprintf("%.1f KB", float64(bytes)/1024)
			} else {
				return fmt.Sprintf("%.1f MB", float64(bytes)/(1024*1024))
			}
		},
		"title": strings.Title,
		"div": func(a, b float64) float64 {
			if b == 0 {
				return 0
			}
			return a / b
		},
		"mul": func(a, b float64) float64 {
			return a * b
		},
		"add": func(a, b int) int {
			return a + b
		},
		"sub": func(a, b float64) float64 {
			return a - b
		},
	}).Parse(htmlTemplate)
	if err != nil {
		return fmt.Errorf("failed to parse HTML template: %w", err)
	}

	outputPath := tr.getReportPath("html")
	file, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("failed to create HTML report file: %w", err)
	}
	defer file.Close()

	if err := tmpl.Execute(file, reportData); err != nil {
		return fmt.Errorf("failed to execute HTML template: %w", err)
	}

	return nil
}

// generateXMLReport generates an XML format report
func (tr *TestReporter) generateXMLReport(result *TestResult) error {
	reportData := tr.prepareReportData(result)
	
	xmlTemplate := `<?xml version="1.0" encoding="UTF-8"?>
<testReport>
    <metadata>
        <testId>{{.TestID}}</testId>
        <serverPath>{{.ServerPath}}</serverPath>
        <timestamp>{{.Timestamp.Format "2006-01-02T15:04:05Z07:00"}}</timestamp>
        <duration>{{.Duration}}</duration>
        <success>{{.Success}}</success>
    </metadata>
    <summary>
        <totalTests>{{.TotalTests}}</totalTests>
        <passedTests>{{.PassedTests}}</passedTests>
        <failedTests>{{.FailedTests}}</failedTests>
        <skippedTests>{{.SkippedTests}}</skippedTests>
    </summary>
    {{if .ValidationResults}}
    <validationResults>
        {{range $name, $result := .ValidationResults}}
        <validation name="{{$name}}">
            <success>{{$result.Success}}</success>
            <duration>{{$result.Duration}}</duration>
            <filesValidated>{{$result.FilesValidated}}</filesValidated>
            {{if $result.Errors}}
            <errors>
                {{range $result.Errors}}<error>{{.}}</error>{{end}}
            </errors>
            {{end}}
            {{if $result.Warnings}}
            <warnings>
                {{range $result.Warnings}}<warning>{{.}}</warning>{{end}}
            </warnings>
            {{end}}
        </validation>
        {{end}}
    </validationResults>
    {{end}}
    {{if .ProtocolResults}}
    <protocolResults>
        <success>{{.ProtocolResults.Success}}</success>
        <duration>{{.ProtocolResults.Duration}}</duration>
        <protocolVersion>{{.ProtocolResults.ProtocolVersion}}</protocolVersion>
        <supportedMethods>
            {{range .ProtocolResults.SupportedMethods}}<method>{{.}}</method>{{end}}
        </supportedMethods>
        <supportedCapabilities>
            {{range .ProtocolResults.SupportedCapabilities}}<capability>{{.}}</capability>{{end}}
        </supportedCapabilities>
    </protocolResults>
    {{end}}
    {{if .PerformanceResults}}
    <performanceResults>
        <success>{{.PerformanceResults.Success}}</success>
        <duration>{{.PerformanceResults.Duration}}</duration>
        <averageResponseTime>{{.PerformanceResults.AverageResponseTime}}</averageResponseTime>
        <peakMemoryUsage>{{.PerformanceResults.PeakMemoryUsage}}</peakMemoryUsage>
        <requestsPerSecond>{{.PerformanceResults.RequestsPerSecond}}</requestsPerSecond>
        <successfulRequests>{{.PerformanceResults.SuccessfulRequests}}</successfulRequests>
        <failedRequests>{{.PerformanceResults.FailedRequests}}</failedRequests>
    </performanceResults>
    {{end}}
    {{if .Errors}}
    <errors>
        {{range .Errors}}<error>{{.}}</error>{{end}}
    </errors>
    {{end}}
    {{if .Recommendations}}
    <recommendations>
        {{range .Recommendations}}<recommendation>{{.}}</recommendation>{{end}}
    </recommendations>
    {{end}}
</testReport>`

	tmpl, err := template.New("xmlReport").Parse(xmlTemplate)
	if err != nil {
		return fmt.Errorf("failed to parse XML template: %w", err)
	}

	outputPath := tr.getReportPath("xml")
	file, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("failed to create XML report file: %w", err)
	}
	defer file.Close()

	if err := tmpl.Execute(file, reportData); err != nil {
		return fmt.Errorf("failed to execute XML template: %w", err)
	}

	return nil
}

// prepareReportData prepares data for report generation
func (tr *TestReporter) prepareReportData(result *TestResult) *TestResult {
	// Return a copy of the result with any necessary transformations
	return result
}

// getReportPath generates the output path for the report
func (tr *TestReporter) getReportPath(format string) string {
	if tr.config.ReportOutputPath != "" {
		// If a specific path is configured, use it
		if strings.HasSuffix(tr.config.ReportOutputPath, "."+format) {
			return tr.config.ReportOutputPath
		}
		// Add extension if not present
		return tr.config.ReportOutputPath + "." + format
	}

	// Generate default path
	timestamp := time.Now().Format("20060102_150405")
	filename := fmt.Sprintf("mcp_test_report_%s.%s", timestamp, format)
	
	// Use current directory by default
	return filepath.Join(".", filename)
}

// GenerateMetricsReport generates a metrics-focused report
func (tr *TestReporter) GenerateMetricsReport(result *TestResult) (*TestMetrics, error) {
	metrics := &TestMetrics{
		TestID:        result.TestID,
		Timestamp:     result.Timestamp,
		Duration:      result.Duration,
		OverallScore:  tr.calculateOverallScore(result),
		QualityScore:  tr.calculateQualityScore(result),
		PerformanceScore: tr.calculatePerformanceScore(result),
		ComplianceScore: tr.calculateComplianceScore(result),
		Categories:    make(map[string]*CategoryMetrics),
	}

	// Calculate category metrics
	if len(result.ValidationResults) > 0 {
		metrics.Categories["validation"] = tr.calculateValidationMetrics(result.ValidationResults)
	}

	if result.ProtocolResults != nil {
		metrics.Categories["protocol"] = tr.calculateProtocolMetrics(result.ProtocolResults)
	}

	if result.PerformanceResults != nil {
		metrics.Categories["performance"] = tr.calculatePerformanceMetrics(result.PerformanceResults)
	}

	if result.IntegrationResults != nil {
		metrics.Categories["integration"] = tr.calculateIntegrationMetrics(result.IntegrationResults)
	}

	return metrics, nil
}

// TestMetrics represents comprehensive test metrics
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

// CategoryMetrics represents metrics for a specific test category
type CategoryMetrics struct {
	Score        float64                 `json:"score"`
	TotalTests   int                     `json:"totalTests"`
	PassedTests  int                     `json:"passedTests"`
	FailedTests  int                     `json:"failedTests"`
	Duration     time.Duration           `json:"duration"`
	Details      map[string]interface{}  `json:"details,omitempty"`
}

// Calculate various scores and metrics

func (tr *TestReporter) calculateOverallScore(result *TestResult) float64 {
	if result.TotalTests == 0 {
		return 0.0
	}
	
	// Basic success rate
	baseScore := float64(result.PassedTests) / float64(result.TotalTests) * 100
	
	// Apply penalties for errors
	errorPenalty := float64(len(result.Errors)) * 5.0
	
	// Apply bonus for comprehensive testing
	comprehensiveBonus := 0.0
	if result.ProtocolResults != nil && result.PerformanceResults != nil {
		comprehensiveBonus = 5.0
	}
	
	score := baseScore - errorPenalty + comprehensiveBonus
	
	// Clamp between 0 and 100
	if score < 0 {
		return 0.0
	}
	if score > 100 {
		return 100.0
	}
	
	return score
}

func (tr *TestReporter) calculateQualityScore(result *TestResult) float64 {
	if len(result.ValidationResults) == 0 {
		return 0.0
	}
	
	totalScore := 0.0
	validationCount := 0
	
	for _, validation := range result.ValidationResults {
		if validation.Success {
			totalScore += 100.0
		} else {
			// Partial score based on files validated vs errors
			if validation.FilesValidated > 0 {
				errorRatio := float64(len(validation.Errors)) / float64(validation.FilesValidated)
				totalScore += (1.0 - errorRatio) * 100.0
			}
		}
		validationCount++
	}
	
	if validationCount == 0 {
		return 0.0
	}
	
	return totalScore / float64(validationCount)
}

func (tr *TestReporter) calculatePerformanceScore(result *TestResult) float64 {
	if result.PerformanceResults == nil {
		return 0.0
	}
	
	perf := result.PerformanceResults
	score := 100.0
	
	// Response time scoring (assuming target is 1 second)
	targetResponseTime := time.Second
	if perf.AverageResponseTime > targetResponseTime {
		penalty := float64(perf.AverageResponseTime) / float64(targetResponseTime) * 20
		score -= penalty
	}
	
	// Memory usage scoring (assuming target is 100MB)
	targetMemory := int64(100 * 1024 * 1024)
	if perf.PeakMemoryUsage > targetMemory {
		penalty := float64(perf.PeakMemoryUsage) / float64(targetMemory) * 15
		score -= penalty
	}
	
	// Memory leak penalty
	if perf.MemoryLeakDetected {
		score -= 30.0
	}
	
	// Success rate bonus
	if perf.SuccessfulRequests > 0 {
		successRate := float64(perf.SuccessfulRequests) / float64(perf.SuccessfulRequests + perf.FailedRequests)
		score = score * successRate
	}
	
	if score < 0 {
		return 0.0
	}
	if score > 100 {
		return 100.0
	}
	
	return score
}

func (tr *TestReporter) calculateComplianceScore(result *TestResult) float64 {
	if result.ProtocolResults == nil {
		return 0.0
	}
	
	protocol := result.ProtocolResults
	if !protocol.Success {
		return 0.0
	}
	
	score := 100.0
	
	// Check required methods support
	requiredMethods := tr.config.RequiredMethods
	supportedCount := 0
	for _, required := range requiredMethods {
		for _, supported := range protocol.SupportedMethods {
			if required == supported {
				supportedCount++
				break
			}
		}
	}
	
	if len(requiredMethods) > 0 {
		methodScore := float64(supportedCount) / float64(len(requiredMethods)) * 50
		score = 50 + methodScore
	}
	
	// Penalty for protocol errors
	errorPenalty := float64(len(protocol.Errors)) * 10
	score -= errorPenalty
	
	if score < 0 {
		return 0.0
	}
	if score > 100 {
		return 100.0
	}
	
	return score
}

func (tr *TestReporter) calculateValidationMetrics(validationResults map[string]*ValidationResult) *CategoryMetrics {
	metrics := &CategoryMetrics{
		Details: make(map[string]interface{}),
	}
	
	totalDuration := time.Duration(0)
	totalFiles := 0
	
	for _, result := range validationResults {
		metrics.TotalTests++
		if result.Success {
			metrics.PassedTests++
		} else {
			metrics.FailedTests++
		}
		totalDuration += result.Duration
		totalFiles += result.FilesValidated
	}
	
	metrics.Duration = totalDuration
	metrics.Score = tr.calculateQualityScore(&TestResult{ValidationResults: validationResults})
	metrics.Details["totalFilesValidated"] = totalFiles
	
	return metrics
}

func (tr *TestReporter) calculateProtocolMetrics(protocolResults *ProtocolTestResult) *CategoryMetrics {
	metrics := &CategoryMetrics{
		TotalTests:  1,
		Duration:    protocolResults.Duration,
		Details:     make(map[string]interface{}),
	}
	
	if protocolResults.Success {
		metrics.PassedTests = 1
	} else {
		metrics.FailedTests = 1
	}
	
	metrics.Score = tr.calculateComplianceScore(&TestResult{ProtocolResults: protocolResults})
	metrics.Details["supportedMethods"] = len(protocolResults.SupportedMethods)
	metrics.Details["supportedCapabilities"] = len(protocolResults.SupportedCapabilities)
	
	return metrics
}

func (tr *TestReporter) calculatePerformanceMetrics(performanceResults *PerformanceTestResult) *CategoryMetrics {
	metrics := &CategoryMetrics{
		TotalTests:  1,
		Duration:    performanceResults.Duration,
		Details:     make(map[string]interface{}),
	}
	
	if performanceResults.Success {
		metrics.PassedTests = 1
	} else {
		metrics.FailedTests = 1
	}
	
	metrics.Score = tr.calculatePerformanceScore(&TestResult{PerformanceResults: performanceResults})
	metrics.Details["averageResponseTime"] = performanceResults.AverageResponseTime.String()
	metrics.Details["peakMemoryUsage"] = performanceResults.PeakMemoryUsage
	metrics.Details["requestsPerSecond"] = performanceResults.RequestsPerSecond
	
	return metrics
}

func (tr *TestReporter) calculateIntegrationMetrics(integrationResults *IntegrationTestResult) *CategoryMetrics {
	metrics := &CategoryMetrics{
		TotalTests:  len(integrationResults.ScenarioResults),
		Duration:    integrationResults.Duration,
		Details:     make(map[string]interface{}),
	}
	
	for _, scenario := range integrationResults.ScenarioResults {
		if scenario.Success {
			metrics.PassedTests++
		} else {
			metrics.FailedTests++
		}
	}
	
	if metrics.TotalTests > 0 {
		metrics.Score = float64(metrics.PassedTests) / float64(metrics.TotalTests) * 100
	}
	
	metrics.Details["scenarioCount"] = len(integrationResults.ScenarioResults)
	metrics.Details["clientCompatibility"] = integrationResults.ClientCompatibility
	
	return metrics
}