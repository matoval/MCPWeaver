package testing

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

// DiagnosticsEngine handles test failure analysis and debugging
type DiagnosticsEngine struct {
	config *TestConfig
}

// NewDiagnosticsEngine creates a new diagnostics engine
func NewDiagnosticsEngine(config *TestConfig) *DiagnosticsEngine {
	return &DiagnosticsEngine{
		config: config,
	}
}

// DiagnosticReport represents a comprehensive diagnostic analysis
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

// FailureAnalysis contains detailed analysis of test failures
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

// StageFailureInfo contains information about failed pipeline stages
type StageFailureInfo struct {
	StageName      string            `json:"stageName"`
	ErrorMessage   string            `json:"errorMessage"`
	FailureTime    time.Time         `json:"failureTime"`
	Duration       time.Duration     `json:"duration"`
	RetryCount     int               `json:"retryCount"`
	Context        map[string]string `json:"context"`
	Suggestions    []string          `json:"suggestions"`
}

// ErrorPattern represents a recognized error pattern
type ErrorPattern struct {
	Pattern     string   `json:"pattern"`
	Description string   `json:"description"`
	Category    string   `json:"category"`
	Solutions   []string `json:"solutions"`
	Frequency   int      `json:"frequency"`
}

// RootCauseAnalysis provides deep analysis of the root cause
type RootCauseAnalysis struct {
	ProbableCause    string            `json:"probableCause"`
	ConfidenceLevel  float64           `json:"confidenceLevel"`
	ContributingFactors []string       `json:"contributingFactors"`
	Evidence         []string          `json:"evidence"`
	Analysis         string            `json:"analysis"`
}

// EnvironmentInfo contains environment analysis
type EnvironmentInfo struct {
	GoVersion         string            `json:"goVersion"`
	OperatingSystem   string            `json:"operatingSystem"`
	Architecture      string            `json:"architecture"`
	AvailableTools    map[string]bool   `json:"availableTools"`
	EnvironmentVars   map[string]string `json:"environmentVars"`
	PathInfo          []string          `json:"pathInfo"`
	SystemResources   *SystemResources  `json:"systemResources"`
}

// SystemResources contains system resource information
type SystemResources struct {
	AvailableMemory string `json:"availableMemory"`
	DiskSpace       string `json:"diskSpace"`
	CPUInfo         string `json:"cpuInfo"`
	LoadAverage     string `json:"loadAverage"`
}

// CodeAnalysis contains analysis of the generated code
type CodeAnalysis struct {
	FilesAnalyzed     int                 `json:"filesAnalyzed"`
	LinesOfCode       int                 `json:"linesOfCode"`
	SyntaxIssues      []SyntaxIssue       `json:"syntaxIssues"`
	StructuralIssues  []StructuralIssue   `json:"structuralIssues"`
	QualityMetrics    *QualityMetrics     `json:"qualityMetrics"`
	MissingComponents []string            `json:"missingComponents"`
	CodeSmells        []CodeSmell         `json:"codeSmells"`
}

// SyntaxIssue represents a syntax-related issue
type SyntaxIssue struct {
	File        string `json:"file"`
	Line        int    `json:"line"`
	Column      int    `json:"column"`
	Message     string `json:"message"`
	Severity    string `json:"severity"`
	Suggestion  string `json:"suggestion"`
}

// StructuralIssue represents a structural code issue
type StructuralIssue struct {
	Type        string `json:"type"`
	Description string `json:"description"`
	Location    string `json:"location"`
	Severity    string `json:"severity"`
	Impact      string `json:"impact"`
	Solution    string `json:"solution"`
}

// QualityMetrics contains code quality metrics
type QualityMetrics struct {
	CyclomaticComplexity int     `json:"cyclomaticComplexity"`
	CodeDuplication      float64 `json:"codeDuplication"`
	TestCoverage         float64 `json:"testCoverage"`
	TechnicalDebt        string  `json:"technicalDebt"`
	Maintainability      string  `json:"maintainability"`
}

// CodeSmell represents a code smell detection
type CodeSmell struct {
	Type        string `json:"type"`
	Description string `json:"description"`
	Location    string `json:"location"`
	Severity    string `json:"severity"`
	Refactoring string `json:"refactoring"`
}

// DependencyAnalysis contains dependency analysis
type DependencyAnalysis struct {
	TotalDependencies int                    `json:"totalDependencies"`
	DirectDependencies int                   `json:"directDependencies"`
	IndirectDependencies int                 `json:"indirectDependencies"`
	OutdatedDependencies []OutdatedDependency `json:"outdatedDependencies"`
	VulnerableDependencies []VulnerableDependency `json:"vulnerableDependencies"`
	LicenseIssues       []LicenseIssue        `json:"licenseIssues"`
	DependencyConflicts []DependencyConflict  `json:"dependencyConflicts"`
}

// OutdatedDependency represents an outdated dependency
type OutdatedDependency struct {
	Name           string `json:"name"`
	CurrentVersion string `json:"currentVersion"`
	LatestVersion  string `json:"latestVersion"`
	UpdateType     string `json:"updateType"`
	ReleaseNotes   string `json:"releaseNotes"`
}

// VulnerableDependency represents a dependency with security vulnerabilities
type VulnerableDependency struct {
	Name            string   `json:"name"`
	Version         string   `json:"version"`
	Vulnerabilities []string `json:"vulnerabilities"`
	Severity        string   `json:"severity"`
	FixedIn         string   `json:"fixedIn"`
}

// LicenseIssue represents a license compatibility issue
type LicenseIssue struct {
	Dependency string `json:"dependency"`
	License    string `json:"license"`
	Issue      string `json:"issue"`
	Severity   string `json:"severity"`
}

// DependencyConflict represents a dependency version conflict
type DependencyConflict struct {
	Package        string   `json:"package"`
	ConflictingVersions []string `json:"conflictingVersions"`
	Resolution     string   `json:"resolution"`
}

// DiagnosticRecommendation represents an actionable recommendation
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

// Action represents a specific action to take
type Action struct {
	Step        int    `json:"step"`
	Description string `json:"description"`
	Command     string `json:"command,omitempty"`
	Example     string `json:"example,omitempty"`
	Warning     string `json:"warning,omitempty"`
}

// TroubleshootingGuide provides step-by-step troubleshooting
type TroubleshootingGuide struct {
	QuickFixes      []QuickFix      `json:"quickFixes"`
	DetailedSteps   []DetailedStep  `json:"detailedSteps"`
	CommonPitfalls  []CommonPitfall `json:"commonPitfalls"`
	PreventionTips  []string        `json:"preventionTips"`
}

// QuickFix represents a quick solution attempt
type QuickFix struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Command     string `json:"command"`
	Expected    string `json:"expected"`
	Risks       string `json:"risks"`
}

// DetailedStep represents a detailed troubleshooting step
type DetailedStep struct {
	Step        int      `json:"step"`
	Title       string   `json:"title"`
	Description string   `json:"description"`
	Commands    []string `json:"commands"`
	Expected    string   `json:"expected"`
	NextSteps   []string `json:"nextSteps"`
}

// CommonPitfall represents a common mistake or issue
type CommonPitfall struct {
	Title       string   `json:"title"`
	Description string   `json:"description"`
	Signs       []string `json:"signs"`
	Solution    string   `json:"solution"`
	Prevention  string   `json:"prevention"`
}

// RelatedIssue represents a related known issue
type RelatedIssue struct {
	Title       string   `json:"title"`
	Description string   `json:"description"`
	URL         string   `json:"url"`
	Similarity  float64  `json:"similarity"`
	Status      string   `json:"status"`
}

// AnalyzeFailures performs comprehensive failure analysis
func (de *DiagnosticsEngine) AnalyzeFailures(ctx context.Context, serverPath string, testResult *TestResult, pipelineResult *PipelineResult) (*DiagnosticReport, error) {
	reportID := fmt.Sprintf("diag_%d", time.Now().Unix())
	
	report := &DiagnosticReport{
		ReportID:    reportID,
		Timestamp:   time.Now(),
		ServerPath:  serverPath,
		Recommendations: make([]DiagnosticRecommendation, 0),
		RelatedIssues:   make([]RelatedIssue, 0),
	}

	// Analyze failures
	failureAnalysis, err := de.analyzeFailures(testResult, pipelineResult)
	if err != nil {
		return nil, fmt.Errorf("failure analysis failed: %w", err)
	}
	report.FailureAnalysis = failureAnalysis

	// Analyze environment
	envInfo, err := de.analyzeEnvironment(ctx, serverPath)
	if err != nil {
		return nil, fmt.Errorf("environment analysis failed: %w", err)
	}
	report.EnvironmentInfo = envInfo

	// Analyze code
	codeAnalysis, err := de.analyzeCode(ctx, serverPath)
	if err != nil {
		return nil, fmt.Errorf("code analysis failed: %w", err)
	}
	report.CodeAnalysis = codeAnalysis

	// Analyze dependencies
	depAnalysis, err := de.analyzeDependencies(ctx, serverPath)
	if err != nil {
		return nil, fmt.Errorf("dependency analysis failed: %w", err)
	}
	report.Dependencies = depAnalysis

	// Generate recommendations
	recommendations := de.generateRecommendations(report)
	report.Recommendations = recommendations

	// Generate troubleshooting guide
	troubleshootingGuide := de.generateTroubleshootingGuide(report)
	report.TroubleshootingGuide = troubleshootingGuide

	// Determine severity and fix time
	report.Severity = de.determineSeverity(report)
	report.EstimatedFixTime = de.estimateFixTime(report)

	return report, nil
}

// analyzeFailures analyzes test and pipeline failures
func (de *DiagnosticsEngine) analyzeFailures(testResult *TestResult, pipelineResult *PipelineResult) (*FailureAnalysis, error) {
	analysis := &FailureAnalysis{
		SecondaryErrors: make([]string, 0),
		FailedStages:    make([]StageFailureInfo, 0),
		ErrorPatterns:   make([]ErrorPattern, 0),
	}

	// Analyze test result failures
	if testResult != nil && !testResult.Success {
		analysis.PrimaryError = "Test execution failed"
		analysis.SecondaryErrors = append(analysis.SecondaryErrors, testResult.Errors...)
		
		// Categorize failure type
		analysis.FailureType = de.categorizeFailureType(testResult.Errors)
		analysis.FailureCategory = de.categorizeFailure(testResult.Errors)
	}

	// Analyze pipeline stage failures
	if pipelineResult != nil {
		for stageName, stageResult := range pipelineResult.StageResults {
			if !stageResult.Success {
				stageInfo := StageFailureInfo{
					StageName:    stageName,
					ErrorMessage: stageResult.ErrorMessage,
					FailureTime:  stageResult.EndTime,
					Duration:     stageResult.Duration,
					RetryCount:   stageResult.RetryCount,
					Context:      make(map[string]string),
					Suggestions:  de.generateStageSuggestions(stageName, stageResult.ErrorMessage),
				}
				analysis.FailedStages = append(analysis.FailedStages, stageInfo)
			}
		}
	}

	// Detect error patterns
	allErrors := append(testResult.Errors, pipelineResult.Errors...)
	analysis.ErrorPatterns = de.detectErrorPatterns(allErrors)

	// Perform root cause analysis
	analysis.RootCause = de.performRootCauseAnalysis(analysis)

	// Determine impact and reproducibility
	analysis.Impact = de.determineImpact(analysis)
	analysis.Reproducibility = de.determineReproducibility(analysis)

	return analysis, nil
}

// categorizeFailureType categorizes the type of failure
func (de *DiagnosticsEngine) categorizeFailureType(errors []string) string {
	errorText := strings.Join(errors, " ")
	errorText = strings.ToLower(errorText)

	if strings.Contains(errorText, "compilation") || strings.Contains(errorText, "build") {
		return "compilation_error"
	}
	if strings.Contains(errorText, "syntax") || strings.Contains(errorText, "parse") {
		return "syntax_error"
	}
	if strings.Contains(errorText, "timeout") || strings.Contains(errorText, "deadline") {
		return "timeout_error"
	}
	if strings.Contains(errorText, "connection") || strings.Contains(errorText, "network") {
		return "network_error"
	}
	if strings.Contains(errorText, "permission") || strings.Contains(errorText, "access") {
		return "permission_error"
	}
	if strings.Contains(errorText, "memory") || strings.Contains(errorText, "oom") {
		return "memory_error"
	}
	if strings.Contains(errorText, "dependency") || strings.Contains(errorText, "module") {
		return "dependency_error"
	}

	return "unknown_error"
}

// categorizeFailure categorizes failures into broader categories
func (de *DiagnosticsEngine) categorizeFailure(errors []string) string {
	errorText := strings.Join(errors, " ")
	errorText = strings.ToLower(errorText)

	if strings.Contains(errorText, "compilation") || strings.Contains(errorText, "syntax") {
		return "code_issue"
	}
	if strings.Contains(errorText, "timeout") || strings.Contains(errorText, "performance") {
		return "performance_issue"
	}
	if strings.Contains(errorText, "protocol") || strings.Contains(errorText, "mcp") {
		return "protocol_issue"
	}
	if strings.Contains(errorText, "environment") || strings.Contains(errorText, "system") {
		return "environment_issue"
	}
	if strings.Contains(errorText, "dependency") || strings.Contains(errorText, "module") {
		return "dependency_issue"
	}

	return "general_issue"
}

// generateStageSuggestions generates suggestions for specific stage failures
func (de *DiagnosticsEngine) generateStageSuggestions(stageName, errorMessage string) []string {
	suggestions := make([]string, 0)
	errorMsg := strings.ToLower(errorMessage)

	switch stageName {
	case "compilation_validation":
		if strings.Contains(errorMsg, "package") {
			suggestions = append(suggestions, "Check package declarations and imports")
			suggestions = append(suggestions, "Verify go.mod file is correct")
		}
		if strings.Contains(errorMsg, "undefined") {
			suggestions = append(suggestions, "Check for missing function or variable definitions")
			suggestions = append(suggestions, "Verify all required dependencies are imported")
		}
	case "syntax_validation":
		suggestions = append(suggestions, "Run 'go fmt' to fix formatting issues")
		suggestions = append(suggestions, "Check for missing brackets or semicolons")
	case "security_scan":
		suggestions = append(suggestions, "Review security scanner output for specific issues")
		suggestions = append(suggestions, "Update dependencies to secure versions")
	case "dependency_check":
		suggestions = append(suggestions, "Run 'go mod tidy' to resolve dependency issues")
		suggestions = append(suggestions, "Check for module version conflicts")
	}

	if len(suggestions) == 0 {
		suggestions = append(suggestions, "Review error message for specific guidance")
		suggestions = append(suggestions, "Check logs for additional context")
	}

	return suggestions
}

// detectErrorPatterns detects common error patterns
func (de *DiagnosticsEngine) detectErrorPatterns(errors []string) []ErrorPattern {
	patterns := make([]ErrorPattern, 0)
	
	commonPatterns := []struct {
		pattern     string
		description string
		category    string
		solutions   []string
	}{
		{
			pattern:     `cannot find package`,
			description: "Missing Go package dependency",
			category:    "dependency",
			solutions:   []string{"Run 'go mod tidy'", "Check import paths", "Verify package exists"},
		},
		{
			pattern:     `undefined: \w+`,
			description: "Undefined identifier",
			category:    "compilation",
			solutions:   []string{"Check function/variable definitions", "Verify imports", "Check spelling"},
		},
		{
			pattern:     `timeout|deadline exceeded`,
			description: "Operation timeout",
			category:    "performance",
			solutions:   []string{"Increase timeout values", "Optimize performance", "Check system resources"},
		},
		{
			pattern:     `connection refused|no such host`,
			description: "Network connectivity issue",
			category:    "network",
			solutions:   []string{"Check network connection", "Verify hostnames", "Check firewall settings"},
		},
		{
			pattern:     `permission denied`,
			description: "File permission issue",
			category:    "permission",
			solutions:   []string{"Check file permissions", "Run with appropriate privileges", "Verify file ownership"},
		},
	}

	errorText := strings.Join(errors, "\n")
	
	for _, patternInfo := range commonPatterns {
		matched, _ := regexp.MatchString(patternInfo.pattern, errorText)
		if matched {
			pattern := ErrorPattern{
				Pattern:     patternInfo.pattern,
				Description: patternInfo.description,
				Category:    patternInfo.category,
				Solutions:   patternInfo.solutions,
				Frequency:   1, // In real implementation, track frequency
			}
			patterns = append(patterns, pattern)
		}
	}

	return patterns
}

// performRootCauseAnalysis performs deep root cause analysis
func (de *DiagnosticsEngine) performRootCauseAnalysis(analysis *FailureAnalysis) *RootCauseAnalysis {
	rootCause := &RootCauseAnalysis{
		ContributingFactors: make([]string, 0),
		Evidence:           make([]string, 0),
	}

	// Analyze based on failure type
	switch analysis.FailureType {
	case "compilation_error":
		rootCause.ProbableCause = "Generated code contains syntax errors or missing dependencies"
		rootCause.ConfidenceLevel = 0.9
		rootCause.ContributingFactors = []string{
			"Template generation issues",
			"Missing import statements",
			"Incorrect package structure",
		}
	case "timeout_error":
		rootCause.ProbableCause = "Server response time exceeds configured limits"
		rootCause.ConfidenceLevel = 0.8
		rootCause.ContributingFactors = []string{
			"Inefficient algorithm implementation",
			"Large data processing",
			"System resource constraints",
		}
	case "protocol_error":
		rootCause.ProbableCause = "Generated server doesn't conform to MCP protocol"
		rootCause.ConfidenceLevel = 0.85
		rootCause.ContributingFactors = []string{
			"Missing protocol methods",
			"Incorrect message format",
			"Invalid response structure",
		}
	default:
		rootCause.ProbableCause = "Multiple factors contributing to test failure"
		rootCause.ConfidenceLevel = 0.6
	}

	// Add evidence from error patterns
	for _, pattern := range analysis.ErrorPatterns {
		rootCause.Evidence = append(rootCause.Evidence, pattern.Description)
	}

	rootCause.Analysis = de.generateRootCauseAnalysisText(rootCause)

	return rootCause
}

// generateRootCauseAnalysisText generates analysis text
func (de *DiagnosticsEngine) generateRootCauseAnalysisText(rootCause *RootCauseAnalysis) string {
	return fmt.Sprintf(
		"Based on the error patterns and evidence collected, the most probable cause is: %s (confidence: %.0f%%). "+
		"Contributing factors include: %s. "+
		"Evidence supporting this analysis: %s.",
		rootCause.ProbableCause,
		rootCause.ConfidenceLevel*100,
		strings.Join(rootCause.ContributingFactors, ", "),
		strings.Join(rootCause.Evidence, ", "),
	)
}

// determineImpact determines the impact of the failure
func (de *DiagnosticsEngine) determineImpact(analysis *FailureAnalysis) string {
	failedStagesCount := len(analysis.FailedStages)
	
	if strings.Contains(analysis.FailureType, "compilation") {
		return "high" // Compilation failures prevent any functionality
	}
	
	if failedStagesCount >= 3 {
		return "high"
	} else if failedStagesCount >= 2 {
		return "medium"
	}
	
	return "low"
}

// determineReproducibility determines how reproducible the failure is
func (de *DiagnosticsEngine) determineReproducibility(analysis *FailureAnalysis) string {
	if strings.Contains(analysis.FailureType, "compilation") || 
	   strings.Contains(analysis.FailureType, "syntax") {
		return "always" // Syntax/compilation errors are deterministic
	}
	
	if strings.Contains(analysis.FailureType, "timeout") ||
	   strings.Contains(analysis.FailureType, "network") {
		return "intermittent" // Network/timing issues can be intermittent
	}
	
	return "often" // Default
}

// analyzeEnvironment analyzes the testing environment
func (de *DiagnosticsEngine) analyzeEnvironment(ctx context.Context, serverPath string) (*EnvironmentInfo, error) {
	envInfo := &EnvironmentInfo{
		AvailableTools:  make(map[string]bool),
		EnvironmentVars: make(map[string]string),
		PathInfo:        make([]string, 0),
		SystemResources: &SystemResources{},
	}

	// Check available tools
	tools := []string{"go", "git", "golangci-lint", "gosec", "govulncheck"}
	for _, tool := range tools {
		envInfo.AvailableTools[tool] = de.commandExists(tool)
	}

	// Get Go version
	if envInfo.AvailableTools["go"] {
		envInfo.GoVersion = de.getGoVersion()
	}

	// Get OS info
	envInfo.OperatingSystem = de.getOperatingSystem()
	envInfo.Architecture = de.getArchitecture()

	// Get relevant environment variables
	relevantVars := []string{"GOPATH", "GOROOT", "PATH", "GOPROXY", "GOSUMDB"}
	for _, varName := range relevantVars {
		if value := os.Getenv(varName); value != "" {
			envInfo.EnvironmentVars[varName] = value
		}
	}

	return envInfo, nil
}

// analyzeCode analyzes the generated code for issues
func (de *DiagnosticsEngine) analyzeCode(ctx context.Context, serverPath string) (*CodeAnalysis, error) {
	analysis := &CodeAnalysis{
		SyntaxIssues:      make([]SyntaxIssue, 0),
		StructuralIssues:  make([]StructuralIssue, 0),
		MissingComponents: make([]string, 0),
		CodeSmells:        make([]CodeSmell, 0),
		QualityMetrics:    &QualityMetrics{},
	}

	// Analyze main.go
	mainFile := filepath.Join(serverPath, "main.go")
	if content, err := os.ReadFile(mainFile); err == nil {
		analysis.FilesAnalyzed++
		analysis.LinesOfCode += strings.Count(string(content), "\n")
		
		// Check for missing components
		analysis.MissingComponents = de.checkMissingComponents(string(content))
		
		// Detect code smells
		analysis.CodeSmells = de.detectCodeSmells(string(content))
	}

	// Analyze go.mod
	goModFile := filepath.Join(serverPath, "go.mod")
	if _, err := os.Stat(goModFile); err == nil {
		analysis.FilesAnalyzed++
	}

	return analysis, nil
}

// checkMissingComponents checks for missing required components
func (de *DiagnosticsEngine) checkMissingComponents(content string) []string {
	missing := make([]string, 0)
	
	requiredComponents := map[string]string{
		"main function":     `func main\(\)`,
		"MCP server type":   `type.*Server.*struct`,
		"initialize method": `func.*initialize`,
		"tools/list method": `func.*tools.*list`,
		"JSON-RPC handling": `jsonrpc2`,
	}

	for component, pattern := range requiredComponents {
		matched, _ := regexp.MatchString(pattern, content)
		if !matched {
			missing = append(missing, component)
		}
	}

	return missing
}

// detectCodeSmells detects code smells in the generated code
func (de *DiagnosticsEngine) detectCodeSmells(content string) []CodeSmell {
	smells := make([]CodeSmell, 0)

	// Check for long functions (>50 lines)
	functionPattern := `func\s+\w+\([^)]*\)[^{]*\{`
	re := regexp.MustCompile(functionPattern)
	matches := re.FindAllStringIndex(content, -1)
	
	for _, match := range matches {
		// Simple heuristic: count lines in function
		functionStart := match[0]
		remainingContent := content[functionStart:]
		braceCount := 0
		lines := 0
		
		for i, char := range remainingContent {
			if char == '\n' {
				lines++
			}
			if char == '{' {
				braceCount++
			} else if char == '}' {
				braceCount--
				if braceCount == 0 {
					if lines > 50 {
						smells = append(smells, CodeSmell{
							Type:        "long_function",
							Description: "Function is too long (>50 lines)",
							Location:    fmt.Sprintf("Character %d", functionStart),
							Severity:    "medium",
							Refactoring: "Break function into smaller, focused functions",
						})
					}
					break
				}
			}
			// Safety check to avoid infinite loop
			if i > 10000 {
				break
			}
		}
	}

	return smells
}

// analyzeDependencies analyzes project dependencies
func (de *DiagnosticsEngine) analyzeDependencies(ctx context.Context, serverPath string) (*DependencyAnalysis, error) {
	analysis := &DependencyAnalysis{
		OutdatedDependencies:   make([]OutdatedDependency, 0),
		VulnerableDependencies: make([]VulnerableDependency, 0),
		LicenseIssues:         make([]LicenseIssue, 0),
		DependencyConflicts:   make([]DependencyConflict, 0),
	}

	// Read go.mod file
	goModPath := filepath.Join(serverPath, "go.mod")
	content, err := os.ReadFile(goModPath)
	if err != nil {
		return analysis, nil // Not critical if go.mod doesn't exist
	}

	// Count dependencies
	lines := strings.Split(string(content), "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.Contains(line, "require") || 
		   (strings.Contains(line, "/") && !strings.HasPrefix(line, "//")) {
			analysis.TotalDependencies++
		}
	}

	analysis.DirectDependencies = analysis.TotalDependencies
	// Indirect dependencies would require more complex parsing

	return analysis, nil
}

// Helper methods

func (de *DiagnosticsEngine) commandExists(command string) bool {
	// Simplified implementation
	return true // Assume tools exist for this example
}

func (de *DiagnosticsEngine) getGoVersion() string {
	return "go1.21" // Simplified
}

func (de *DiagnosticsEngine) getOperatingSystem() string {
	return "linux" // Simplified
}

func (de *DiagnosticsEngine) getArchitecture() string {
	return "amd64" // Simplified
}

// generateRecommendations generates actionable recommendations
func (de *DiagnosticsEngine) generateRecommendations(report *DiagnosticReport) []DiagnosticRecommendation {
	recommendations := make([]DiagnosticRecommendation, 0)

	// Generate recommendations based on failure analysis
	if report.FailureAnalysis != nil {
		switch report.FailureAnalysis.FailureType {
		case "compilation_error":
			recommendations = append(recommendations, DiagnosticRecommendation{
				ID:          "fix_compilation",
				Priority:    "high",
				Category:    "code",
				Title:       "Fix Compilation Errors",
				Description: "Resolve syntax and compilation issues in the generated code",
				Actions: []Action{
					{Step: 1, Description: "Run 'go build' to see detailed error messages"},
					{Step: 2, Description: "Check import statements and package declarations"},
					{Step: 3, Description: "Verify all required functions are implemented"},
				},
				EstimatedTime: "15-30 minutes",
			})
		case "timeout_error":
			recommendations = append(recommendations, DiagnosticRecommendation{
				ID:          "fix_timeout",
				Priority:    "medium",
				Category:    "performance",
				Title:       "Resolve Timeout Issues",
				Description: "Optimize performance to meet response time requirements",
				Actions: []Action{
					{Step: 1, Description: "Profile application performance"},
					{Step: 2, Description: "Optimize slow operations"},
					{Step: 3, Description: "Consider increasing timeout values if appropriate"},
				},
				EstimatedTime: "30-60 minutes",
			})
		}
	}

	return recommendations
}

// generateTroubleshootingGuide generates comprehensive troubleshooting guide
func (de *DiagnosticsEngine) generateTroubleshootingGuide(report *DiagnosticReport) *TroubleshootingGuide {
	return &TroubleshootingGuide{
		QuickFixes: []QuickFix{
			{
				Name:        "Clean and rebuild",
				Description: "Clean build artifacts and rebuild",
				Command:     "go clean && go build",
				Expected:    "Successful compilation",
				Risks:       "Low risk",
			},
		},
		DetailedSteps: []DetailedStep{
			{
				Step:        1,
				Title:       "Verify Environment",
				Description: "Check that all required tools and dependencies are available",
				Commands:    []string{"go version", "go mod verify"},
				Expected:    "Go compiler available and dependencies verified",
				NextSteps:   []string{"If Go not available, install Go", "If dependencies fail, run 'go mod tidy'"},
			},
		},
		CommonPitfalls: []CommonPitfall{
			{
				Title:       "Missing Go installation",
				Description: "Go compiler not found in PATH",
				Signs:       []string{"'go' command not found", "compilation fails"},
				Solution:    "Install Go from https://golang.org/dl/",
				Prevention:  "Ensure Go is properly installed and added to PATH",
			},
		},
		PreventionTips: []string{
			"Regularly update Go and dependencies",
			"Use consistent coding standards",
			"Test generated code before deployment",
			"Monitor system resources during testing",
		},
	}
}

// determineSeverity determines the severity of the issues
func (de *DiagnosticsEngine) determineSeverity(report *DiagnosticReport) string {
	if report.FailureAnalysis != nil {
		switch report.FailureAnalysis.Impact {
		case "high":
			return "critical"
		case "medium":
			return "major"
		case "low":
			return "minor"
		}
	}
	return "minor"
}

// estimateFixTime estimates the time needed to fix the issues
func (de *DiagnosticsEngine) estimateFixTime(report *DiagnosticReport) string {
	if report.FailureAnalysis != nil {
		switch report.FailureAnalysis.FailureType {
		case "compilation_error":
			return "15-30 minutes"
		case "syntax_error":
			return "5-15 minutes"
		case "timeout_error":
			return "30-60 minutes"
		case "protocol_error":
			return "45-90 minutes"
		default:
			return "30-60 minutes"
		}
	}
	return "30-60 minutes"
}

// SaveDiagnosticReport saves the diagnostic report to a file
func (de *DiagnosticsEngine) SaveDiagnosticReport(report *DiagnosticReport, outputPath string) error {
	data, err := json.MarshalIndent(report, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to serialize diagnostic report: %w", err)
	}

	if err := os.WriteFile(outputPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write diagnostic report: %w", err)
	}

	return nil
}