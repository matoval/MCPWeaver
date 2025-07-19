package app

import (
	"context"
	"runtime"
	"sync"
	"time"
)

// TemplatePerformanceMetrics represents performance metrics for template operations
type TemplatePerformanceMetrics struct {
	TemplateID      string        `json:"templateId"`
	Operation       string        `json:"operation"`
	StartTime       time.Time     `json:"startTime"`
	EndTime         time.Time     `json:"endTime"`
	Duration        time.Duration `json:"duration"`
	MemoryUsage     int64         `json:"memoryUsage"`
	CPUUsage        float64       `json:"cpuUsage"`
	Success         bool          `json:"success"`
	ErrorMessage    string        `json:"errorMessage,omitempty"`
	InputSize       int64         `json:"inputSize,omitempty"`
	OutputSize      int64         `json:"outputSize,omitempty"`
	CacheHit        bool          `json:"cacheHit"`
	Complexity      string        `json:"complexity"`
	VariableCount   int           `json:"variableCount"`
	FunctionCount   int           `json:"functionCount"`
	Timestamp       time.Time     `json:"timestamp"`
}

// SystemMetrics represents overall system performance metrics
type SystemMetrics struct {
	Timestamp         time.Time `json:"timestamp"`
	MemoryTotal       int64     `json:"memoryTotal"`
	MemoryUsed        int64     `json:"memoryUsed"`
	MemoryAvailable   int64     `json:"memoryAvailable"`
	CPUUsage          float64   `json:"cpuUsage"`
	GoRoutines        int       `json:"goRoutines"`
	HeapAlloc         uint64    `json:"heapAlloc"`
	HeapSys           uint64    `json:"heapSys"`
	NumGC             uint32    `json:"numGC"`
	ActiveTemplates   int       `json:"activeTemplates"`
	CacheSize         int       `json:"cacheSize"`
	RequestsPerSecond float64   `json:"requestsPerSecond"`
}

// AggregatedMetrics represents aggregated performance data
type AggregatedMetrics struct {
	Operation        string        `json:"operation"`
	TotalExecutions  int           `json:"totalExecutions"`
	SuccessCount     int           `json:"successCount"`
	FailureCount     int           `json:"failureCount"`
	SuccessRate      float64       `json:"successRate"`
	AverageDuration  time.Duration `json:"averageDuration"`
	MinDuration      time.Duration `json:"minDuration"`
	MaxDuration      time.Duration `json:"maxDuration"`
	P50Duration      time.Duration `json:"p50Duration"`
	P95Duration      time.Duration `json:"p95Duration"`
	P99Duration      time.Duration `json:"p99Duration"`
	AverageMemory    int64         `json:"averageMemory"`
	AverageCPU       float64       `json:"averageCpu"`
	CacheHitRate     float64       `json:"cacheHitRate"`
	LastUpdated      time.Time     `json:"lastUpdated"`
}

// PerformanceAlert represents a performance alert
type PerformanceAlert struct {
	ID          string                 `json:"id"`
	Type        string                 `json:"type"`
	Severity    string                 `json:"severity"`
	Message     string                 `json:"message"`
	Threshold   map[string]interface{} `json:"threshold"`
	ActualValue map[string]interface{} `json:"actualValue"`
	TemplateID  string                 `json:"templateId,omitempty"`
	Operation   string                 `json:"operation,omitempty"`
	Timestamp   time.Time              `json:"timestamp"`
	Resolved    bool                   `json:"resolved"`
}

// MetricsCollector handles performance metrics collection and analysis
type MetricsCollector struct {
	metrics    []TemplatePerformanceMetrics
	mutex      sync.RWMutex
	alerts     []PerformanceAlert
	thresholds map[string]interface{}
}

var metricsCollector = &MetricsCollector{
	metrics: make([]TemplatePerformanceMetrics, 0),
	alerts:  make([]PerformanceAlert, 0),
	thresholds: map[string]interface{}{
		"maxDuration":         10 * time.Second,
		"maxMemoryUsage":      100 * 1024 * 1024, // 100MB
		"maxCPUUsage":         80.0,               // 80%
		"minSuccessRate":      95.0,               // 95%
		"maxErrorRate":        5.0,                // 5%
		"alertRetentionDays":  7,
		"metricsRetentionDays": 30,
	},
}

// StartPerformanceMonitoring starts monitoring template performance
func (a *App) StartPerformanceMonitoring(ctx context.Context, templateID, operation string) *PerformanceMonitor {
	return &PerformanceMonitor{
		templateID: templateID,
		operation:  operation,
		startTime:  time.Now(),
		memStart:   getMemoryUsage(),
		cpuStart:   getCPUUsage(),
	}
}

// PerformanceMonitor tracks performance for a single operation
type PerformanceMonitor struct {
	templateID string
	operation  string
	startTime  time.Time
	memStart   int64
	cpuStart   float64
	inputSize  int64
	outputSize int64
	cacheHit   bool
	complexity string
	varCount   int
	funcCount  int
}

// SetInputSize sets the input size for the operation
func (pm *PerformanceMonitor) SetInputSize(size int64) {
	pm.inputSize = size
}

// SetOutputSize sets the output size for the operation
func (pm *PerformanceMonitor) SetOutputSize(size int64) {
	pm.outputSize = size
}

// SetCacheHit indicates whether the operation was a cache hit
func (pm *PerformanceMonitor) SetCacheHit(hit bool) {
	pm.cacheHit = hit
}

// SetComplexity sets the complexity level of the operation
func (pm *PerformanceMonitor) SetComplexity(complexity string) {
	pm.complexity = complexity
}

// SetCounts sets variable and function counts
func (pm *PerformanceMonitor) SetCounts(variables, functions int) {
	pm.varCount = variables
	pm.funcCount = functions
}

// End completes the performance monitoring and records metrics
func (pm *PerformanceMonitor) End(success bool, errorMessage string) {
	endTime := time.Now()
	duration := endTime.Sub(pm.startTime)
	memEnd := getMemoryUsage()
	cpuEnd := getCPUUsage()

	metrics := TemplatePerformanceMetrics{
		TemplateID:    pm.templateID,
		Operation:     pm.operation,
		StartTime:     pm.startTime,
		EndTime:       endTime,
		Duration:      duration,
		MemoryUsage:   memEnd - pm.memStart,
		CPUUsage:      cpuEnd - pm.cpuStart,
		Success:       success,
		ErrorMessage:  errorMessage,
		InputSize:     pm.inputSize,
		OutputSize:    pm.outputSize,
		CacheHit:      pm.cacheHit,
		Complexity:    pm.complexity,
		VariableCount: pm.varCount,
		FunctionCount: pm.funcCount,
		Timestamp:     time.Now(),
	}

	metricsCollector.RecordMetrics(metrics)
}

// RecordMetrics records performance metrics
func (mc *MetricsCollector) RecordMetrics(metrics TemplatePerformanceMetrics) {
	mc.mutex.Lock()
	defer mc.mutex.Unlock()

	mc.metrics = append(mc.metrics, metrics)

	// Check for performance alerts
	mc.checkAlerts(metrics)

	// Clean up old metrics (keep only last 30 days)
	mc.cleanupOldMetrics()
}

// GetTemplateMetrics retrieves metrics for a specific template
func (a *App) GetTemplateMetrics(ctx context.Context, templateID string, timeRange string) ([]TemplatePerformanceMetrics, error) {
	if templateID == "" {
		return nil, createMonitoringError("Template ID is required", "", "get_metrics")
	}

	metricsCollector.mutex.RLock()
	defer metricsCollector.mutex.RUnlock()

	var filtered []TemplatePerformanceMetrics
	cutoff := time.Now()

	switch timeRange {
	case "1h":
		cutoff = cutoff.Add(-1 * time.Hour)
	case "24h":
		cutoff = cutoff.Add(-24 * time.Hour)
	case "7d":
		cutoff = cutoff.Add(-7 * 24 * time.Hour)
	case "30d":
		cutoff = cutoff.Add(-30 * 24 * time.Hour)
	default:
		cutoff = cutoff.Add(-24 * time.Hour) // Default to 24h
	}

	for _, metric := range metricsCollector.metrics {
		if metric.TemplateID == templateID && metric.Timestamp.After(cutoff) {
			filtered = append(filtered, metric)
		}
	}

	return filtered, nil
}

// GetSystemMetrics retrieves current system performance metrics
func (a *App) GetSystemMetrics(ctx context.Context) (*SystemMetrics, error) {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	return &SystemMetrics{
		Timestamp:         time.Now(),
		MemoryTotal:       int64(m.Sys),
		MemoryUsed:        int64(m.Alloc),
		MemoryAvailable:   int64(m.Sys - m.Alloc),
		CPUUsage:          getCPUUsage(),
		GoRoutines:        runtime.NumGoroutine(),
		HeapAlloc:         m.HeapAlloc,
		HeapSys:           m.HeapSys,
		NumGC:             m.NumGC,
		ActiveTemplates:   a.getActiveTemplateCount(),
		CacheSize:         len(templateCache),
		RequestsPerSecond: a.calculateRequestsPerSecond(),
	}, nil
}

// GetAggregatedMetrics retrieves aggregated performance metrics
func (a *App) GetAggregatedMetrics(ctx context.Context, operation string, timeRange string) (*AggregatedMetrics, error) {
	metricsCollector.mutex.RLock()
	defer metricsCollector.mutex.RUnlock()

	cutoff := time.Now()
	switch timeRange {
	case "1h":
		cutoff = cutoff.Add(-1 * time.Hour)
	case "24h":
		cutoff = cutoff.Add(-24 * time.Hour)
	case "7d":
		cutoff = cutoff.Add(-7 * 24 * time.Hour)
	case "30d":
		cutoff = cutoff.Add(-30 * 24 * time.Hour)
	default:
		cutoff = cutoff.Add(-24 * time.Hour)
	}

	var filtered []TemplatePerformanceMetrics
	for _, metric := range metricsCollector.metrics {
		if (operation == "" || metric.Operation == operation) && metric.Timestamp.After(cutoff) {
			filtered = append(filtered, metric)
		}
	}

	if len(filtered) == 0 {
		return &AggregatedMetrics{
			Operation:   operation,
			LastUpdated: time.Now(),
		}, nil
	}

	return calculateAggregatedMetrics(filtered, operation), nil
}

// GetPerformanceAlerts retrieves active performance alerts
func (a *App) GetPerformanceAlerts(ctx context.Context, includeResolved bool) ([]PerformanceAlert, error) {
	metricsCollector.mutex.RLock()
	defer metricsCollector.mutex.RUnlock()

	var alerts []PerformanceAlert
	for _, alert := range metricsCollector.alerts {
		if includeResolved || !alert.Resolved {
			alerts = append(alerts, alert)
		}
	}

	return alerts, nil
}

// ResolveAlert marks a performance alert as resolved
func (a *App) ResolveAlert(ctx context.Context, alertID string) error {
	if alertID == "" {
		return createMonitoringError("Alert ID is required", "", "resolve_alert")
	}

	metricsCollector.mutex.Lock()
	defer metricsCollector.mutex.Unlock()

	for i, alert := range metricsCollector.alerts {
		if alert.ID == alertID {
			metricsCollector.alerts[i].Resolved = true
			return nil
		}
	}

	return createMonitoringError("Alert not found", alertID, "resolve_alert")
}

// SetPerformanceThresholds updates performance alert thresholds
func (a *App) SetPerformanceThresholds(ctx context.Context, thresholds map[string]interface{}) error {
	metricsCollector.mutex.Lock()
	defer metricsCollector.mutex.Unlock()

	// Validate thresholds
	validKeys := map[string]bool{
		"maxDuration":         true,
		"maxMemoryUsage":      true,
		"maxCPUUsage":         true,
		"minSuccessRate":      true,
		"maxErrorRate":        true,
		"alertRetentionDays":  true,
		"metricsRetentionDays": true,
	}

	for key := range thresholds {
		if !validKeys[key] {
			return createMonitoringError("Invalid threshold key", key, "set_thresholds")
		}
	}

	// Update thresholds
	for key, value := range thresholds {
		metricsCollector.thresholds[key] = value
	}

	return nil
}

// GetPerformanceThresholds retrieves current performance thresholds
func (a *App) GetPerformanceThresholds(ctx context.Context) (map[string]interface{}, error) {
	metricsCollector.mutex.RLock()
	defer metricsCollector.mutex.RUnlock()

	// Return a copy to prevent external modification
	thresholds := make(map[string]interface{})
	for key, value := range metricsCollector.thresholds {
		thresholds[key] = value
	}

	return thresholds, nil
}

// ClearMetrics clears all performance metrics
func (a *App) ClearMetrics(ctx context.Context, olderThan string) error {
	metricsCollector.mutex.Lock()
	defer metricsCollector.mutex.Unlock()

	cutoff := time.Now()
	switch olderThan {
	case "1h":
		cutoff = cutoff.Add(-1 * time.Hour)
	case "24h":
		cutoff = cutoff.Add(-24 * time.Hour)
	case "7d":
		cutoff = cutoff.Add(-7 * 24 * time.Hour)
	case "30d":
		cutoff = cutoff.Add(-30 * 24 * time.Hour)
	case "all":
		metricsCollector.metrics = make([]TemplatePerformanceMetrics, 0)
		return nil
	default:
		return createMonitoringError("Invalid time range", olderThan, "clear_metrics")
	}

	var filtered []TemplatePerformanceMetrics
	for _, metric := range metricsCollector.metrics {
		if metric.Timestamp.After(cutoff) {
			filtered = append(filtered, metric)
		}
	}

	metricsCollector.metrics = filtered
	return nil
}

// Helper functions

func (mc *MetricsCollector) checkAlerts(metrics TemplatePerformanceMetrics) {
	// Check duration threshold
	if maxDuration, ok := mc.thresholds["maxDuration"].(time.Duration); ok {
		if metrics.Duration > maxDuration {
			alert := PerformanceAlert{
				ID:       generateAlertID(),
				Type:     "performance",
				Severity: "warning",
				Message:  "Template operation exceeded maximum duration",
				Threshold: map[string]interface{}{
					"maxDuration": maxDuration.String(),
				},
				ActualValue: map[string]interface{}{
					"duration": metrics.Duration.String(),
				},
				TemplateID: metrics.TemplateID,
				Operation:  metrics.Operation,
				Timestamp:  time.Now(),
				Resolved:   false,
			}
			mc.alerts = append(mc.alerts, alert)
		}
	}

	// Check memory usage threshold
	if maxMemory, ok := mc.thresholds["maxMemoryUsage"].(int); ok {
		if metrics.MemoryUsage > int64(maxMemory) {
			alert := PerformanceAlert{
				ID:       generateAlertID(),
				Type:     "memory",
				Severity: "warning",
				Message:  "Template operation exceeded maximum memory usage",
				Threshold: map[string]interface{}{
					"maxMemoryUsage": maxMemory,
				},
				ActualValue: map[string]interface{}{
					"memoryUsage": metrics.MemoryUsage,
				},
				TemplateID: metrics.TemplateID,
				Operation:  metrics.Operation,
				Timestamp:  time.Now(),
				Resolved:   false,
			}
			mc.alerts = append(mc.alerts, alert)
		}
	}

	// Check CPU usage threshold
	if maxCPU, ok := mc.thresholds["maxCPUUsage"].(float64); ok {
		if metrics.CPUUsage > maxCPU {
			alert := PerformanceAlert{
				ID:       generateAlertID(),
				Type:     "cpu",
				Severity: "warning",
				Message:  "Template operation exceeded maximum CPU usage",
				Threshold: map[string]interface{}{
					"maxCPUUsage": maxCPU,
				},
				ActualValue: map[string]interface{}{
					"cpuUsage": metrics.CPUUsage,
				},
				TemplateID: metrics.TemplateID,
				Operation:  metrics.Operation,
				Timestamp:  time.Now(),
				Resolved:   false,
			}
			mc.alerts = append(mc.alerts, alert)
		}
	}
}

func (mc *MetricsCollector) cleanupOldMetrics() {
	retentionDays := 30
	if days, ok := mc.thresholds["metricsRetentionDays"].(int); ok {
		retentionDays = days
	}

	cutoff := time.Now().AddDate(0, 0, -retentionDays)
	var filtered []TemplatePerformanceMetrics

	for _, metric := range mc.metrics {
		if metric.Timestamp.After(cutoff) {
			filtered = append(filtered, metric)
		}
	}

	mc.metrics = filtered

	// Clean up old alerts
	alertRetentionDays := 7
	if days, ok := mc.thresholds["alertRetentionDays"].(int); ok {
		alertRetentionDays = days
	}

	alertCutoff := time.Now().AddDate(0, 0, -alertRetentionDays)
	var filteredAlerts []PerformanceAlert

	for _, alert := range mc.alerts {
		if alert.Timestamp.After(alertCutoff) || !alert.Resolved {
			filteredAlerts = append(filteredAlerts, alert)
		}
	}

	mc.alerts = filteredAlerts
}

func calculateAggregatedMetrics(metrics []TemplatePerformanceMetrics, operation string) *AggregatedMetrics {
	if len(metrics) == 0 {
		return &AggregatedMetrics{
			Operation:   operation,
			LastUpdated: time.Now(),
		}
	}

	var totalDuration time.Duration
	var totalMemory int64
	var totalCPU float64
	var successCount int
	var cacheHits int
	var durations []time.Duration

	for _, metric := range metrics {
		totalDuration += metric.Duration
		totalMemory += metric.MemoryUsage
		totalCPU += metric.CPUUsage
		durations = append(durations, metric.Duration)

		if metric.Success {
			successCount++
		}
		if metric.CacheHit {
			cacheHits++
		}
	}

	// Sort durations for percentile calculations
	for i := 0; i < len(durations); i++ {
		for j := i + 1; j < len(durations); j++ {
			if durations[i] > durations[j] {
				durations[i], durations[j] = durations[j], durations[i]
			}
		}
	}

	totalExecutions := len(metrics)
	failureCount := totalExecutions - successCount

	return &AggregatedMetrics{
		Operation:        operation,
		TotalExecutions:  totalExecutions,
		SuccessCount:     successCount,
		FailureCount:     failureCount,
		SuccessRate:      float64(successCount) / float64(totalExecutions) * 100,
		AverageDuration:  totalDuration / time.Duration(totalExecutions),
		MinDuration:      durations[0],
		MaxDuration:      durations[totalExecutions-1],
		P50Duration:      durations[totalExecutions/2],
		P95Duration:      durations[int(float64(totalExecutions)*0.95)],
		P99Duration:      durations[int(float64(totalExecutions)*0.99)],
		AverageMemory:    totalMemory / int64(totalExecutions),
		AverageCPU:       totalCPU / float64(totalExecutions),
		CacheHitRate:     float64(cacheHits) / float64(totalExecutions) * 100,
		LastUpdated:      time.Now(),
	}
}

func getMemoryUsage() int64 {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	return int64(m.Alloc)
}

func getCPUUsage() float64 {
	// Simple CPU usage approximation
	// In a real implementation, you might use a more sophisticated method
	return float64(runtime.NumGoroutine()) * 0.1
}

func (a *App) getActiveTemplateCount() int {
	// Return number of templates in cache or database
	return len(templateCache)
}

func (a *App) calculateRequestsPerSecond() float64 {
	// Calculate based on recent metrics
	// This is a simplified implementation
	return 0.0
}

func generateAlertID() string {
	return fmt.Sprintf("alert_%d", time.Now().UnixNano())
}

func createMonitoringError(message, details, operation string) error {
	return createError("monitoring", "MONITORING_ERROR", message, map[string]string{
		"details":   details,
		"operation": operation,
	})
}