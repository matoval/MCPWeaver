package app

import (
	"runtime"
	"sync"
	"time"
)

// PerformanceMetrics tracks application performance metrics
type PerformanceMetrics struct {
	StartupTime       time.Duration            `json:"startup_time"`
	MemoryUsage       int64                    `json:"memory_usage"`
	GenerationTimes   map[string]time.Duration `json:"generation_times"`
	DatabaseOpTimes   map[string]time.Duration `json:"database_op_times"`
	FileOpTimes       map[string]time.Duration `json:"file_op_times"`
	LastUpdated       time.Time                `json:"last_updated"`
	mutex             sync.RWMutex
}

// PerformanceMonitor handles performance monitoring
type PerformanceMonitor struct {
	metrics   *PerformanceMetrics
	startTime time.Time
	mutex     sync.RWMutex
}

// NewPerformanceMonitor creates a new performance monitor
func NewPerformanceMonitor() *PerformanceMonitor {
	return &PerformanceMonitor{
		metrics: &PerformanceMetrics{
			GenerationTimes: make(map[string]time.Duration),
			DatabaseOpTimes: make(map[string]time.Duration),
			FileOpTimes:     make(map[string]time.Duration),
		},
		startTime: time.Now(),
	}
}

// RecordStartupTime records the application startup time
func (pm *PerformanceMonitor) RecordStartupTime(duration time.Duration) {
	pm.mutex.Lock()
	defer pm.mutex.Unlock()
	pm.metrics.StartupTime = duration
	pm.metrics.LastUpdated = time.Now()
}

// RecordMemoryUsage records current memory usage
func (pm *PerformanceMonitor) RecordMemoryUsage() {
	pm.mutex.Lock()
	defer pm.mutex.Unlock()
	
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	pm.metrics.MemoryUsage = int64(m.Alloc)
	pm.metrics.LastUpdated = time.Now()
}

// RecordGenerationTime records generation time for a specific operation
func (pm *PerformanceMonitor) RecordGenerationTime(operation string, duration time.Duration) {
	pm.mutex.Lock()
	defer pm.mutex.Unlock()
	pm.metrics.GenerationTimes[operation] = duration
	pm.metrics.LastUpdated = time.Now()
}

// RecordDatabaseOpTime records database operation time
func (pm *PerformanceMonitor) RecordDatabaseOpTime(operation string, duration time.Duration) {
	pm.mutex.Lock()
	defer pm.mutex.Unlock()
	pm.metrics.DatabaseOpTimes[operation] = duration
	pm.metrics.LastUpdated = time.Now()
}

// RecordFileOpTime records file operation time
func (pm *PerformanceMonitor) RecordFileOpTime(operation string, duration time.Duration) {
	pm.mutex.Lock()
	defer pm.mutex.Unlock()
	pm.metrics.FileOpTimes[operation] = duration
	pm.metrics.LastUpdated = time.Now()
}

// GetMetrics returns current performance metrics
func (pm *PerformanceMonitor) GetMetrics() *PerformanceMetrics {
	pm.mutex.RLock()
	defer pm.mutex.RUnlock()
	
	// Update memory usage before returning metrics
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	
	// Create a copy to avoid race conditions
	metrics := &PerformanceMetrics{
		StartupTime:     pm.metrics.StartupTime,
		MemoryUsage:     int64(m.Alloc),
		GenerationTimes: make(map[string]time.Duration),
		DatabaseOpTimes: make(map[string]time.Duration),
		FileOpTimes:     make(map[string]time.Duration),
		LastUpdated:     time.Now(),
	}
	
	// Copy maps
	for k, v := range pm.metrics.GenerationTimes {
		metrics.GenerationTimes[k] = v
	}
	for k, v := range pm.metrics.DatabaseOpTimes {
		metrics.DatabaseOpTimes[k] = v
	}
	for k, v := range pm.metrics.FileOpTimes {
		metrics.FileOpTimes[k] = v
	}
	
	return metrics
}

// StartTimer returns a function that records the duration when called
func (pm *PerformanceMonitor) StartTimer(category, operation string) func() {
	start := time.Now()
	return func() {
		duration := time.Since(start)
		switch category {
		case "generation":
			pm.RecordGenerationTime(operation, duration)
		case "database":
			pm.RecordDatabaseOpTime(operation, duration)
		case "file":
			pm.RecordFileOpTime(operation, duration)
		}
	}
}

// GetMemoryUsageBytes returns current memory usage in bytes
func (pm *PerformanceMonitor) GetMemoryUsageBytes() int64 {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	return int64(m.Alloc)
}

// GetMemoryUsageMB returns current memory usage in MB
func (pm *PerformanceMonitor) GetMemoryUsageMB() float64 {
	return float64(pm.GetMemoryUsageBytes()) / 1024 / 1024
}

// IsMemoryWithinLimit checks if memory usage is within the specified limit
func (pm *PerformanceMonitor) IsMemoryWithinLimit(limitMB float64) bool {
	return pm.GetMemoryUsageMB() <= limitMB
}

// ForceGC forces garbage collection and records memory usage
func (pm *PerformanceMonitor) ForceGC() {
	runtime.GC()
	pm.RecordMemoryUsage()
}

// GetUptimeSeconds returns application uptime in seconds
func (pm *PerformanceMonitor) GetUptimeSeconds() float64 {
	return time.Since(pm.startTime).Seconds()
}