package testing

import (
	"context"
	"fmt"
	"path/filepath"
	"sync"
	"time"
)

// TestPipeline manages automated test execution
type TestPipeline struct {
	config    *TestConfig
	framework *TestSuite
	mutex     sync.RWMutex
	running   bool
}

// NewTestPipeline creates a new test pipeline
func NewTestPipeline(config *TestConfig) *TestPipeline {
	return &TestPipeline{
		config:    config,
		framework: NewTestSuite(config),
	}
}

// PipelineStage represents a stage in the test pipeline
type PipelineStage struct {
	Name        string                                              `json:"name"`
	Description string                                              `json:"description"`
	Enabled     bool                                                `json:"enabled"`
	Parallel    bool                                                `json:"parallel"`
	Execute     func(ctx context.Context, serverPath string) error  `json:"-"`
	Validate    func(ctx context.Context, result *TestResult) error `json:"-"`
	Timeout     time.Duration                                       `json:"timeout"`
	Retries     int                                                 `json:"retries"`
	OnFailure   string                                              `json:"onFailure"` // "continue", "stop", "retry"
}

// PipelineResult represents the result of a pipeline execution
type PipelineResult struct {
	PipelineID      string                  `json:"pipelineId"`
	StartTime       time.Time               `json:"startTime"`
	EndTime         time.Time               `json:"endTime"`
	Duration        time.Duration           `json:"duration"`
	Success         bool                    `json:"success"`
	StageResults    map[string]*StageResult `json:"stageResults"`
	TestResult      *TestResult             `json:"testResult,omitempty"`
	Errors          []string                `json:"errors"`
	TotalStages     int                     `json:"totalStages"`
	CompletedStages int                     `json:"completedStages"`
	SkippedStages   int                     `json:"skippedStages"`
	FailedStages    int                     `json:"failedStages"`
}

// StageResult represents the result of a single pipeline stage
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

// ExecutePipeline runs the complete test pipeline
func (tp *TestPipeline) ExecutePipeline(ctx context.Context, serverPath string) (*PipelineResult, error) {
	tp.mutex.Lock()
	if tp.running {
		tp.mutex.Unlock()
		return nil, fmt.Errorf("pipeline is already running")
	}
	tp.running = true
	tp.mutex.Unlock()

	defer func() {
		tp.mutex.Lock()
		tp.running = false
		tp.mutex.Unlock()
	}()

	pipelineID := fmt.Sprintf("pipeline_%d", time.Now().Unix())
	startTime := time.Now()

	result := &PipelineResult{
		PipelineID:   pipelineID,
		StartTime:    startTime,
		Success:      true,
		StageResults: make(map[string]*StageResult),
		Errors:       make([]string, 0),
	}

	// Define pipeline stages
	stages := tp.definePipelineStages()
	result.TotalStages = len(stages)

	// Execute stages
	for _, stage := range stages {
		if !stage.Enabled {
			result.SkippedStages++
			result.StageResults[stage.Name] = &StageResult{
				StageName: stage.Name,
				Success:   true,
				Skipped:   true,
			}
			continue
		}

		stageResult := tp.executeStage(ctx, stage, serverPath)
		result.StageResults[stage.Name] = stageResult

		if stageResult.Success {
			result.CompletedStages++
		} else {
			result.FailedStages++
			result.Success = false
			result.Errors = append(result.Errors, fmt.Sprintf("Stage %s failed: %s", stage.Name, stageResult.ErrorMessage))

			// Handle failure based on stage configuration
			if stage.OnFailure == "stop" {
				break
			}
		}
	}

	// Run comprehensive test if all validation stages passed
	if result.Success || tp.config.ContinueOnFailure {
		testResult, err := tp.framework.RunTests(ctx, serverPath)
		if err != nil {
			result.Success = false
			result.Errors = append(result.Errors, fmt.Sprintf("Test execution failed: %v", err))
		} else {
			result.TestResult = testResult
			if !testResult.Success {
				result.Success = false
			}
		}
	}

	result.EndTime = time.Now()
	result.Duration = result.EndTime.Sub(result.StartTime)

	return result, nil
}

// executeStage executes a single pipeline stage
func (tp *TestPipeline) executeStage(ctx context.Context, stage *PipelineStage, serverPath string) *StageResult {
	startTime := time.Now()
	stageResult := &StageResult{
		StageName: stage.Name,
		StartTime: startTime,
		Success:   true,
		Details:   make(map[string]interface{}),
	}

	// Apply stage timeout
	stageCtx := ctx
	if stage.Timeout > 0 {
		var cancel context.CancelFunc
		stageCtx, cancel = context.WithTimeout(ctx, stage.Timeout)
		defer cancel()
	}

	// Execute stage with retries
	var lastError error
	for attempt := 0; attempt <= stage.Retries; attempt++ {
		if attempt > 0 {
			stageResult.RetryCount = attempt
			// Wait before retry
			time.Sleep(time.Duration(attempt) * time.Second)
		}

		err := stage.Execute(stageCtx, serverPath)
		if err == nil {
			break
		}

		lastError = err
		if attempt == stage.Retries {
			stageResult.Success = false
			stageResult.ErrorMessage = lastError.Error()
		}
	}

	stageResult.EndTime = time.Now()
	stageResult.Duration = stageResult.EndTime.Sub(stageResult.StartTime)

	return stageResult
}

// definePipelineStages defines the stages of the test pipeline
func (tp *TestPipeline) definePipelineStages() []*PipelineStage {
	stages := []*PipelineStage{
		{
			Name:        "pre_validation",
			Description: "Pre-execution validation checks",
			Enabled:     true,
			Parallel:    false,
			Timeout:     30 * time.Second,
			Retries:     1,
			OnFailure:   "stop",
			Execute: func(ctx context.Context, serverPath string) error {
				return tp.validatePrerequisites(ctx, serverPath)
			},
		},
		{
			Name:        "dependency_check",
			Description: "Check and install dependencies",
			Enabled:     true,
			Parallel:    false,
			Timeout:     2 * time.Minute,
			Retries:     2,
			OnFailure:   "continue",
			Execute: func(ctx context.Context, serverPath string) error {
				return tp.checkDependencies(ctx, serverPath)
			},
		},
		{
			Name:        "compilation_validation",
			Description: "Validate server compilation",
			Enabled:     true,
			Parallel:    false,
			Timeout:     1 * time.Minute,
			Retries:     1,
			OnFailure:   "stop",
			Execute: func(ctx context.Context, serverPath string) error {
				validator := NewCompilationValidator(tp.config)
				result, err := validator.Validate(ctx, serverPath)
				if err != nil {
					return err
				}
				if !result.Success {
					return fmt.Errorf("compilation validation failed: %v", result.Errors)
				}
				return nil
			},
		},
		{
			Name:        "syntax_validation",
			Description: "Validate code syntax and structure",
			Enabled:     true,
			Parallel:    true,
			Timeout:     30 * time.Second,
			Retries:     0,
			OnFailure:   "continue",
			Execute: func(ctx context.Context, serverPath string) error {
				validator := NewSyntaxValidator(tp.config)
				result, err := validator.Validate(ctx, serverPath)
				if err != nil {
					return err
				}
				if !result.Success {
					return fmt.Errorf("syntax validation failed: %v", result.Errors)
				}
				return nil
			},
		},
		{
			Name:        "security_scan",
			Description: "Run security vulnerability scans",
			Enabled:     tp.config.EnableSecurityScanning,
			Parallel:    true,
			Timeout:     2 * time.Minute,
			Retries:     1,
			OnFailure:   "continue",
			Execute: func(ctx context.Context, serverPath string) error {
				validator := NewSecurityValidator(tp.config)
				_, err := validator.Validate(ctx, serverPath)
				if err != nil {
					return err
				}
				// Security warnings don't fail the pipeline
				return nil
			},
		},
		{
			Name:        "lint_check",
			Description: "Run code linting and formatting checks",
			Enabled:     tp.config.EnableLinting,
			Parallel:    true,
			Timeout:     1 * time.Minute,
			Retries:     0,
			OnFailure:   "continue",
			Execute: func(ctx context.Context, serverPath string) error {
				validator := NewLintValidator(tp.config)
				_, err := validator.Validate(ctx, serverPath)
				if err != nil {
					return err
				}
				// Linting warnings don't fail the pipeline
				return nil
			},
		},
		{
			Name:        "environment_setup",
			Description: "Set up test environment",
			Enabled:     true,
			Parallel:    false,
			Timeout:     30 * time.Second,
			Retries:     2,
			OnFailure:   "stop",
			Execute: func(ctx context.Context, serverPath string) error {
				return tp.setupTestEnvironment(ctx, serverPath)
			},
		},
	}

	return stages
}

// validatePrerequisites checks if all prerequisites are met
func (tp *TestPipeline) validatePrerequisites(ctx context.Context, serverPath string) error {
	// Check if server path exists
	if err := tp.framework.fileExists(serverPath); err != nil {
		return fmt.Errorf("server path validation failed: %w", err)
	}

	// Check required files
	requiredFiles := []string{"main.go", "go.mod"}
	for _, file := range requiredFiles {
		filePath := filepath.Join(serverPath, file)
		if err := tp.framework.fileExists(filePath); err != nil {
			return fmt.Errorf("required file %s not found: %w", file, err)
		}
	}

	// Check Go installation
	if !tp.framework.commandExists("go") {
		return fmt.Errorf("Go compiler not found")
	}

	return nil
}

// checkDependencies ensures all dependencies are available
func (tp *TestPipeline) checkDependencies(ctx context.Context, serverPath string) error {
	// Run go mod tidy to ensure dependencies are resolved
	if err := tp.framework.runCommand(ctx, serverPath, "go", "mod", "tidy"); err != nil {
		return fmt.Errorf("failed to resolve dependencies: %w", err)
	}

	// Run go mod download to ensure dependencies are cached
	if err := tp.framework.runCommand(ctx, serverPath, "go", "mod", "download"); err != nil {
		return fmt.Errorf("failed to download dependencies: %w", err)
	}

	return nil
}

// setupTestEnvironment prepares the test environment
func (tp *TestPipeline) setupTestEnvironment(ctx context.Context, serverPath string) error {
	// Create temporary directories if needed
	tempDir := filepath.Join(serverPath, ".test_temp")
	if err := tp.framework.ensureDir(tempDir); err != nil {
		return fmt.Errorf("failed to create temp directory: %w", err)
	}

	// Set up test configuration
	if tp.config.TestDataPath != "" {
		if err := tp.framework.ensureDir(tp.config.TestDataPath); err != nil {
			return fmt.Errorf("failed to create test data directory: %w", err)
		}
	}

	return nil
}

// IsRunning returns whether the pipeline is currently running
func (tp *TestPipeline) IsRunning() bool {
	tp.mutex.RLock()
	defer tp.mutex.RUnlock()
	return tp.running
}

// GetPipelineStatus returns the current pipeline configuration
func (tp *TestPipeline) GetPipelineStatus() map[string]interface{} {
	tp.mutex.RLock()
	defer tp.mutex.RUnlock()

	stages := tp.definePipelineStages()
	stageInfo := make([]map[string]interface{}, 0, len(stages))

	for _, stage := range stages {
		stageInfo = append(stageInfo, map[string]interface{}{
			"name":        stage.Name,
			"description": stage.Description,
			"enabled":     stage.Enabled,
			"parallel":    stage.Parallel,
			"timeout":     stage.Timeout.String(),
			"retries":     stage.Retries,
			"onFailure":   stage.OnFailure,
		})
	}

	return map[string]interface{}{
		"running":       tp.running,
		"totalStages":   len(stages),
		"enabledStages": tp.countEnabledStages(stages),
		"stages":        stageInfo,
		"config": map[string]interface{}{
			"continueOnFailure":      tp.config.ContinueOnFailure,
			"enableSecurityScanning": tp.config.EnableSecurityScanning,
			"enableLinting":          tp.config.EnableLinting,
			"timeout":                tp.config.Timeout.String(),
		},
	}
}

// countEnabledStages counts the number of enabled stages
func (tp *TestPipeline) countEnabledStages(stages []*PipelineStage) int {
	count := 0
	for _, stage := range stages {
		if stage.Enabled {
			count++
		}
	}
	return count
}

// BatchTestRunner handles running tests on multiple servers
type BatchTestRunner struct {
	pipeline *TestPipeline
	config   *TestConfig
}

// NewBatchTestRunner creates a new batch test runner
func NewBatchTestRunner(config *TestConfig) *BatchTestRunner {
	return &BatchTestRunner{
		pipeline: NewTestPipeline(config),
		config:   config,
	}
}

// BatchTestRequest represents a batch test request
type BatchTestRequest struct {
	RequestID     string   `json:"requestId"`
	ServerPaths   []string `json:"serverPaths"`
	Parallel      bool     `json:"parallel"`
	MaxWorkers    int      `json:"maxWorkers"`
	StopOnFailure bool     `json:"stopOnFailure"`
}

// BatchTestResult represents the result of a batch test run
type BatchTestResult struct {
	RequestID      string                     `json:"requestId"`
	StartTime      time.Time                  `json:"startTime"`
	EndTime        time.Time                  `json:"endTime"`
	Duration       time.Duration              `json:"duration"`
	Success        bool                       `json:"success"`
	TotalServers   int                        `json:"totalServers"`
	CompletedTests int                        `json:"completedTests"`
	FailedTests    int                        `json:"failedTests"`
	SkippedTests   int                        `json:"skippedTests"`
	ServerResults  map[string]*PipelineResult `json:"serverResults"`
	Errors         []string                   `json:"errors"`
	Summary        *BatchTestSummary          `json:"summary"`
}

// BatchTestSummary provides aggregate statistics
type BatchTestSummary struct {
	AverageTestDuration  time.Duration          `json:"averageTestDuration"`
	FastestTest          time.Duration          `json:"fastestTest"`
	SlowestTest          time.Duration          `json:"slowestTest"`
	SuccessRate          float64                `json:"successRate"`
	CommonFailures       map[string]int         `json:"commonFailures"`
	StageSuccessRates    map[string]float64     `json:"stageSuccessRates"`
	ResourceUsageSummary map[string]interface{} `json:"resourceUsageSummary"`
}

// RunBatchTests executes tests on multiple servers
func (btr *BatchTestRunner) RunBatchTests(ctx context.Context, request *BatchTestRequest) (*BatchTestResult, error) {
	startTime := time.Now()

	result := &BatchTestResult{
		RequestID:     request.RequestID,
		StartTime:     startTime,
		Success:       true,
		TotalServers:  len(request.ServerPaths),
		ServerResults: make(map[string]*PipelineResult),
		Errors:        make([]string, 0),
	}

	if request.Parallel && request.MaxWorkers > 1 {
		return btr.runParallelBatchTests(ctx, request, result)
	}

	return btr.runSequentialBatchTests(ctx, request, result)
}

// runSequentialBatchTests runs tests sequentially
func (btr *BatchTestRunner) runSequentialBatchTests(ctx context.Context, request *BatchTestRequest, result *BatchTestResult) (*BatchTestResult, error) {
	for _, serverPath := range request.ServerPaths {
		pipelineResult, err := btr.pipeline.ExecutePipeline(ctx, serverPath)
		if err != nil {
			result.Errors = append(result.Errors, fmt.Sprintf("Pipeline execution failed for %s: %v", serverPath, err))
			result.FailedTests++
			result.Success = false

			if request.StopOnFailure {
				break
			}
			continue
		}

		result.ServerResults[serverPath] = pipelineResult

		if pipelineResult.Success {
			result.CompletedTests++
		} else {
			result.FailedTests++
			result.Success = false

			if request.StopOnFailure {
				break
			}
		}
	}

	result.EndTime = time.Now()
	result.Duration = result.EndTime.Sub(result.StartTime)
	result.Summary = btr.generateBatchSummary(result)

	return result, nil
}

// runParallelBatchTests runs tests in parallel
func (btr *BatchTestRunner) runParallelBatchTests(ctx context.Context, request *BatchTestRequest, result *BatchTestResult) (*BatchTestResult, error) {
	// Implement parallel execution with worker pool
	workers := request.MaxWorkers
	if workers <= 0 {
		workers = 3 // Default worker count
	}

	jobs := make(chan string, len(request.ServerPaths))
	results := make(chan *workerResult, len(request.ServerPaths))

	// Start workers
	for i := 0; i < workers; i++ {
		go btr.worker(ctx, jobs, results)
	}

	// Send jobs
	for _, serverPath := range request.ServerPaths {
		jobs <- serverPath
	}
	close(jobs)

	// Collect results
	for i := 0; i < len(request.ServerPaths); i++ {
		workerRes := <-results

		if workerRes.Error != nil {
			result.Errors = append(result.Errors, fmt.Sprintf("Pipeline execution failed for %s: %v", workerRes.ServerPath, workerRes.Error))
			result.FailedTests++
			result.Success = false
		} else {
			result.ServerResults[workerRes.ServerPath] = workerRes.PipelineResult

			if workerRes.PipelineResult.Success {
				result.CompletedTests++
			} else {
				result.FailedTests++
				result.Success = false
			}
		}
	}

	result.EndTime = time.Now()
	result.Duration = result.EndTime.Sub(result.StartTime)
	result.Summary = btr.generateBatchSummary(result)

	return result, nil
}

// workerResult represents the result of a worker execution
type workerResult struct {
	ServerPath     string
	PipelineResult *PipelineResult
	Error          error
}

// worker function for parallel execution
func (btr *BatchTestRunner) worker(ctx context.Context, jobs <-chan string, results chan<- *workerResult) {
	for serverPath := range jobs {
		pipeline := NewTestPipeline(btr.config) // Create new pipeline for each worker
		pipelineResult, err := pipeline.ExecutePipeline(ctx, serverPath)

		results <- &workerResult{
			ServerPath:     serverPath,
			PipelineResult: pipelineResult,
			Error:          err,
		}
	}
}

// generateBatchSummary generates summary statistics for batch test results
func (btr *BatchTestRunner) generateBatchSummary(result *BatchTestResult) *BatchTestSummary {
	summary := &BatchTestSummary{
		CommonFailures:       make(map[string]int),
		StageSuccessRates:    make(map[string]float64),
		ResourceUsageSummary: make(map[string]interface{}),
	}

	if len(result.ServerResults) == 0 {
		return summary
	}

	// Calculate duration statistics
	var totalDuration time.Duration
	var fastest, slowest time.Duration
	first := true

	for _, pipelineResult := range result.ServerResults {
		totalDuration += pipelineResult.Duration

		if first {
			fastest = pipelineResult.Duration
			slowest = pipelineResult.Duration
			first = false
		} else {
			if pipelineResult.Duration < fastest {
				fastest = pipelineResult.Duration
			}
			if pipelineResult.Duration > slowest {
				slowest = pipelineResult.Duration
			}
		}

		// Collect failure patterns
		for _, err := range pipelineResult.Errors {
			summary.CommonFailures[err]++
		}
	}

	summary.AverageTestDuration = totalDuration / time.Duration(len(result.ServerResults))
	summary.FastestTest = fastest
	summary.SlowestTest = slowest

	// Calculate success rate
	if result.TotalServers > 0 {
		summary.SuccessRate = float64(result.CompletedTests) / float64(result.TotalServers) * 100
	}

	// Calculate stage success rates
	stageStats := make(map[string]struct{ total, success int })
	for _, pipelineResult := range result.ServerResults {
		for stageName, stageResult := range pipelineResult.StageResults {
			stats := stageStats[stageName]
			stats.total++
			if stageResult.Success {
				stats.success++
			}
			stageStats[stageName] = stats
		}
	}

	for stageName, stats := range stageStats {
		if stats.total > 0 {
			summary.StageSuccessRates[stageName] = float64(stats.success) / float64(stats.total) * 100
		}
	}

	return summary
}
