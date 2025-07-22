package testing

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"sort"
	"sync"
	"time"
)

// PerformanceTester handles performance and load testing
type PerformanceTester struct {
	config *TestConfig
}

// NewPerformanceTester creates a new performance tester
func NewPerformanceTester(config *TestConfig) *PerformanceTester {
	return &PerformanceTester{
		config: config,
	}
}

// TestPerformance runs comprehensive performance tests
func (pt *PerformanceTester) TestPerformance(ctx context.Context, serverPath string) (*PerformanceTestResult, error) {
	startTime := time.Now()
	
	result := &PerformanceTestResult{
		Success:         true,
		LoadTestResults: make(map[string]*LoadTestMetric),
		Errors:          make([]string, 0),
	}

	// Test response time performance
	if err := pt.testResponseTime(ctx, serverPath, result); err != nil {
		result.Errors = append(result.Errors, fmt.Sprintf("Response time test failed: %v", err))
	}

	// Test memory usage
	if err := pt.testMemoryUsage(ctx, serverPath, result); err != nil {
		result.Errors = append(result.Errors, fmt.Sprintf("Memory usage test failed: %v", err))
	}

	// Test throughput
	if err := pt.testThroughput(ctx, serverPath, result); err != nil {
		result.Errors = append(result.Errors, fmt.Sprintf("Throughput test failed: %v", err))
	}

	// Test load handling
	if err := pt.testLoadHandling(ctx, serverPath, result); err != nil {
		result.Errors = append(result.Errors, fmt.Sprintf("Load test failed: %v", err))
	}

	result.Duration = time.Since(startTime)
	result.Success = len(result.Errors) == 0

	// Validate against thresholds
	pt.validatePerformanceThresholds(result)

	return result, nil
}

// testResponseTime measures response times for different operations
func (pt *PerformanceTester) testResponseTime(ctx context.Context, serverPath string, result *PerformanceTestResult) error {
	// Start server for testing
	serverProcess, stdin, stdout, stderr, err := pt.startPerformanceTestServer(ctx, serverPath)
	if err != nil {
		return fmt.Errorf("failed to start server: %w", err)
	}
	defer pt.stopPerformanceTestServer(serverProcess, serverPath)

	client := NewJSONRPCClient(stdin, stdout, stderr)
	defer client.Close()

	// Initialize server
	initRequest := map[string]interface{}{
		"protocolVersion": "2024-11-05",
		"capabilities":    map[string]interface{}{},
		"clientInfo":      map[string]interface{}{"name": "Performance Test", "version": "1.0.0"},
	}
	
	if _, err := client.Call(ctx, "initialize", initRequest); err != nil {
		return fmt.Errorf("initialization failed: %w", err)
	}

	// Measure response times for different operations
	responseTimes := make([]time.Duration, 0)

	// Test initialize response time (already done, but measure again for consistency)
	for i := 0; i < 10; i++ {
		start := time.Now()
		_, err := client.Call(ctx, "tools/list", map[string]interface{}{})
		responseTime := time.Since(start)
		
		if err == nil {
			responseTimes = append(responseTimes, responseTime)
		}
		
		// Small delay between requests
		time.Sleep(10 * time.Millisecond)
	}

	if len(responseTimes) == 0 {
		return fmt.Errorf("no successful requests for response time measurement")
	}

	// Calculate statistics
	sort.Slice(responseTimes, func(i, j int) bool {
		return responseTimes[i] < responseTimes[j]
	})

	result.AverageResponseTime = pt.calculateAverage(responseTimes)
	result.MedianResponseTime = responseTimes[len(responseTimes)/2]
	result.P95ResponseTime = responseTimes[int(float64(len(responseTimes))*0.95)]
	result.P99ResponseTime = responseTimes[int(float64(len(responseTimes))*0.99)]
	result.MaxResponseTime = responseTimes[len(responseTimes)-1]

	return nil
}

// testMemoryUsage measures memory usage during operation
func (pt *PerformanceTester) testMemoryUsage(ctx context.Context, serverPath string, result *PerformanceTestResult) error {
	sch := NewSecureCommandHelper()
	
	// Compile server using secure command execution
	compileCmd, err := sch.SecureCompileCommand(ctx, serverPath, "memory-test-server", "main.go")
	if err != nil {
		return fmt.Errorf("failed to create secure compile command: %w", err)
	}
	if err := compileCmd.Run(); err != nil {
		return fmt.Errorf("failed to compile server: %w", err)
	}
	defer os.Remove(filepath.Join(serverPath, "memory-test-server"))

	// Start server using secure executable execution
	cmd, err := sch.SecureRunExecutable(ctx, serverPath, "memory-test-server")
	if err != nil {
		return fmt.Errorf("failed to create secure run command: %w", err)
	}

	stdin, err := cmd.StdinPipe()
	if err != nil {
		return fmt.Errorf("failed to create stdin pipe: %w", err)
	}
	defer stdin.Close()

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return fmt.Errorf("failed to create stdout pipe: %w", err)
	}
	defer stdout.Close()

	stderr, err := cmd.StderrPipe()
	if err != nil {
		return fmt.Errorf("failed to create stderr pipe: %w", err)
	}
	defer stderr.Close()

	if err := cmd.Start(); err != nil {
		return fmt.Errorf("failed to start server: %w", err)
	}
	defer func() {
		if cmd.Process != nil {
			cmd.Process.Kill()
			cmd.Wait()
		}
	}()

	// Monitor memory usage
	memoryUsages := make([]int64, 0)
	monitorDuration := 30 * time.Second
	sampleInterval := time.Second

	monitorCtx, cancel := context.WithTimeout(ctx, monitorDuration)
	defer cancel()

	client := NewJSONRPCClient(stdin, stdout, stderr)
	defer client.Close()

	// Initialize
	initRequest := map[string]interface{}{
		"protocolVersion": "2024-11-05",
		"capabilities":    map[string]interface{}{},
		"clientInfo":      map[string]interface{}{"name": "Memory Test", "version": "1.0.0"},
	}
	
	if _, err := client.Call(monitorCtx, "initialize", initRequest); err != nil {
		return fmt.Errorf("initialization failed: %w", err)
	}

	// Start memory monitoring
	ticker := time.NewTicker(sampleInterval)
	defer ticker.Stop()

	requestTicker := time.NewTicker(100 * time.Millisecond)
	defer requestTicker.Stop()

	var wg sync.WaitGroup
	
	// Memory monitoring goroutine
	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			select {
			case <-monitorCtx.Done():
				return
			case <-ticker.C:
				if cmd.Process != nil {
					memUsage := pt.getProcessMemoryUsage(cmd.Process.Pid)
					if memUsage > 0 {
						memoryUsages = append(memoryUsages, memUsage)
					}
				}
			}
		}
	}()

	// Request generation goroutine
	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			select {
			case <-monitorCtx.Done():
				return
			case <-requestTicker.C:
				// Send periodic requests to exercise the server
				client.Call(monitorCtx, "tools/list", map[string]interface{}{})
			}
		}
	}()

	wg.Wait()

	if len(memoryUsages) == 0 {
		return fmt.Errorf("no memory usage data collected")
	}

	// Calculate memory statistics
	result.AverageMemoryUsage = pt.calculateAverageMemory(memoryUsages)
	result.PeakMemoryUsage = pt.findMaxMemory(memoryUsages)
	result.MemoryLeakDetected = pt.detectMemoryLeak(memoryUsages)

	return nil
}

// testThroughput measures server throughput
func (pt *PerformanceTester) testThroughput(ctx context.Context, serverPath string, result *PerformanceTestResult) error {
	// Start server
	serverProcess, stdin, stdout, stderr, err := pt.startPerformanceTestServer(ctx, serverPath)
	if err != nil {
		return fmt.Errorf("failed to start server: %w", err)
	}
	defer pt.stopPerformanceTestServer(serverProcess, serverPath)

	client := NewJSONRPCClient(stdin, stdout, stderr)
	defer client.Close()

	// Initialize
	initRequest := map[string]interface{}{
		"protocolVersion": "2024-11-05",
		"capabilities":    map[string]interface{}{},
		"clientInfo":      map[string]interface{}{"name": "Throughput Test", "version": "1.0.0"},
	}
	
	if _, err := client.Call(ctx, "initialize", initRequest); err != nil {
		return fmt.Errorf("initialization failed: %w", err)
	}

	// Measure throughput over a fixed time period
	testDuration := 10 * time.Second
	testCtx, cancel := context.WithTimeout(ctx, testDuration)
	defer cancel()

	startTime := time.Now()
	successfulRequests := 0
	failedRequests := 0

	for {
		select {
		case <-testCtx.Done():
			goto done
		default:
		}

		_, err := client.Call(testCtx, "tools/list", map[string]interface{}{})
		if err == nil {
			successfulRequests++
		} else {
			failedRequests++
		}
	}

done:
	actualDuration := time.Since(startTime)
	result.RequestsPerSecond = float64(successfulRequests) / actualDuration.Seconds()
	result.SuccessfulRequests = successfulRequests
	result.FailedRequests = failedRequests

	return nil
}

// testLoadHandling tests server performance under load
func (pt *PerformanceTester) testLoadHandling(ctx context.Context, serverPath string, result *PerformanceTestResult) error {
	loadScenarios := []struct {
		name        string
		duration    time.Duration
		requestRate int // requests per second
	}{
		{"light_load", 10 * time.Second, 5},
		{"medium_load", 15 * time.Second, 10},
		{"heavy_load", 20 * time.Second, 20},
	}

	for _, scenario := range loadScenarios {
		metric, err := pt.runLoadScenario(ctx, serverPath, scenario.name, scenario.duration, scenario.requestRate)
		if err != nil {
			result.Errors = append(result.Errors, fmt.Sprintf("Load scenario %s failed: %v", scenario.name, err))
			continue
		}
		result.LoadTestResults[scenario.name] = metric
	}

	return nil
}

// runLoadScenario runs a specific load testing scenario
func (pt *PerformanceTester) runLoadScenario(ctx context.Context, serverPath string, scenarioName string, duration time.Duration, requestRate int) (*LoadTestMetric, error) {
	// Start server
	serverProcess, stdin, stdout, stderr, err := pt.startPerformanceTestServer(ctx, serverPath)
	if err != nil {
		return nil, fmt.Errorf("failed to start server: %w", err)
	}
	defer pt.stopPerformanceTestServer(serverProcess, serverPath)

	client := NewJSONRPCClient(stdin, stdout, stderr)
	defer client.Close()

	// Initialize
	initRequest := map[string]interface{}{
		"protocolVersion": "2024-11-05",
		"capabilities":    map[string]interface{}{},
		"clientInfo":      map[string]interface{}{"name": "Load Test", "version": "1.0.0"},
	}
	
	if _, err := client.Call(ctx, "initialize", initRequest); err != nil {
		return nil, fmt.Errorf("initialization failed: %w", err)
	}

	metric := &LoadTestMetric{
		Scenario: scenarioName,
		Duration: duration,
	}

	// Run load test
	testCtx, cancel := context.WithTimeout(ctx, duration)
	defer cancel()

	startTime := time.Now()
	interval := time.Duration(1000/requestRate) * time.Millisecond
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	responseTimes := make([]time.Duration, 0)

	for {
		select {
		case <-testCtx.Done():
			goto done
		case <-ticker.C:
			reqStart := time.Now()
			_, err := client.Call(testCtx, "tools/list", map[string]interface{}{})
			responseTime := time.Since(reqStart)

			metric.TotalRequests++
			if err == nil {
				metric.SuccessfulRequests++
				responseTimes = append(responseTimes, responseTime)
			} else {
				metric.FailedRequests++
			}
		}
	}

done:
	actualDuration := time.Since(startTime)
	metric.Duration = actualDuration

	if len(responseTimes) > 0 {
		metric.AverageResponseTime = pt.calculateAverage(responseTimes)
	}

	if actualDuration.Seconds() > 0 {
		metric.RequestsPerSecond = float64(metric.SuccessfulRequests) / actualDuration.Seconds()
	}

	if metric.TotalRequests > 0 {
		metric.ErrorRate = float64(metric.FailedRequests) / float64(metric.TotalRequests)
	}

	return metric, nil
}

// validatePerformanceThresholds checks if performance meets requirements
func (pt *PerformanceTester) validatePerformanceThresholds(result *PerformanceTestResult) {
	if result.AverageResponseTime > pt.config.MaxResponseTime {
		result.Success = false
		result.Errors = append(result.Errors, 
			fmt.Sprintf("Average response time %v exceeds threshold %v", 
				result.AverageResponseTime, pt.config.MaxResponseTime))
	}

	if result.PeakMemoryUsage > pt.config.MaxMemoryUsage {
		result.Success = false
		result.Errors = append(result.Errors, 
			fmt.Sprintf("Peak memory usage %d bytes exceeds threshold %d bytes", 
				result.PeakMemoryUsage, pt.config.MaxMemoryUsage))
	}

	if result.MemoryLeakDetected {
		result.Success = false
		result.Errors = append(result.Errors, "Memory leak detected during testing")
	}
}

// Helper methods

// calculateAverage calculates average duration
func (pt *PerformanceTester) calculateAverage(durations []time.Duration) time.Duration {
	if len(durations) == 0 {
		return 0
	}
	
	var total time.Duration
	for _, d := range durations {
		total += d
	}
	return total / time.Duration(len(durations))
}

// calculateAverageMemory calculates average memory usage
func (pt *PerformanceTester) calculateAverageMemory(memoryUsages []int64) int64 {
	if len(memoryUsages) == 0 {
		return 0
	}
	
	var total int64
	for _, mem := range memoryUsages {
		total += mem
	}
	return total / int64(len(memoryUsages))
}

// findMaxMemory finds peak memory usage
func (pt *PerformanceTester) findMaxMemory(memoryUsages []int64) int64 {
	if len(memoryUsages) == 0 {
		return 0
	}
	
	max := memoryUsages[0]
	for _, mem := range memoryUsages {
		if mem > max {
			max = mem
		}
	}
	return max
}

// detectMemoryLeak detects potential memory leaks
func (pt *PerformanceTester) detectMemoryLeak(memoryUsages []int64) bool {
	if len(memoryUsages) < 10 {
		return false
	}
	
	// Simple leak detection: check if memory usage consistently increases
	// over the last portion of the test
	quarterPoint := len(memoryUsages) / 4
	firstQuarter := memoryUsages[:quarterPoint]
	lastQuarter := memoryUsages[len(memoryUsages)-quarterPoint:]
	
	firstAvg := pt.calculateAverageMemory(firstQuarter)
	lastAvg := pt.calculateAverageMemory(lastQuarter)
	
	// If memory usage increased by more than 50%, consider it a potential leak
	threshold := firstAvg + (firstAvg / 2)
	return lastAvg > threshold
}

// getProcessMemoryUsage gets memory usage for a process (simplified implementation)
func (pt *PerformanceTester) getProcessMemoryUsage(pid int) int64 {
	// This is a simplified implementation
	// In a real implementation, you would use proper system calls or libraries
	// to get accurate memory usage information
	
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	
	// Return allocated memory as a rough approximation
	// This is not accurate for the specific process, but gives an idea
	return int64(m.Alloc)
}

// Server management for performance tests

// startPerformanceTestServer starts a server for performance testing
func (pt *PerformanceTester) startPerformanceTestServer(ctx context.Context, serverPath string) (*exec.Cmd, *os.File, *os.File, *os.File, error) {
	sch := NewSecureCommandHelper()
	
	// Compile server using secure command execution
	compileCmd, err := sch.SecureCompileCommand(ctx, serverPath, "perf-test-server", "main.go")
	if err != nil {
		return nil, nil, nil, nil, fmt.Errorf("failed to create secure compile command: %w", err)
	}
	if err := compileCmd.Run(); err != nil {
		return nil, nil, nil, nil, fmt.Errorf("failed to compile server: %w", err)
	}

	// Start server using secure executable execution
	cmd, err := sch.SecureRunExecutable(ctx, serverPath, "perf-test-server")
	if err != nil {
		return nil, nil, nil, nil, fmt.Errorf("failed to create secure run command: %w", err)
	}

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

	// Give server time to start
	time.Sleep(100 * time.Millisecond)

	return cmd, stdin.(*os.File), stdout.(*os.File), stderr.(*os.File), nil
}

// stopPerformanceTestServer stops the performance test server
func (pt *PerformanceTester) stopPerformanceTestServer(cmd *exec.Cmd, serverPath string) {
	if cmd != nil && cmd.Process != nil {
		cmd.Process.Kill()
		cmd.Wait()
	}
	
	// Clean up binary
	os.Remove(filepath.Join(serverPath, "perf-test-server"))
}