package app

import (
	"context"
	"crypto/rand"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"os"
	"runtime"
	"sort"
	"strings"
	"sync"
	"time"
)

// ActivityLogService manages the activity log system with circular buffer and real-time events
type ActivityLogService struct {
	buffer   []ActivityLogEntry
	config   LogConfig
	mutex    sync.RWMutex
	writeIdx int
	full     bool
	app      *App
	ticker   *time.Ticker
	ctx      context.Context
	cancel   context.CancelFunc
}

// NewActivityLogService creates a new activity log service
func NewActivityLogService(app *App, config LogConfig) *ActivityLogService {
	ctx, cancel := context.WithCancel(context.Background())

	service := &ActivityLogService{
		buffer: make([]ActivityLogEntry, config.BufferSize),
		config: config,
		app:    app,
		ctx:    ctx,
		cancel: cancel,
	}

	// Start periodic flush if enabled
	if config.FlushInterval > 0 {
		service.ticker = time.NewTicker(config.FlushInterval)
		go service.flushLoop()
	}

	return service
}

// Error tracking methods

// ReportError creates an error report and logs it
func (als *ActivityLogService) ReportError(errorType ErrorType, severity ErrorSeverity, component, operation, message string, err error) *ErrorReport {
	report := &ErrorReport{
		ID:          generateLogID(),
		Timestamp:   time.Now(),
		Type:        errorType,
		Severity:    severity,
		Component:   component,
		Operation:   operation,
		Message:     message,
		FirstSeen:   time.Now(),
		LastSeen:    time.Now(),
		Frequency:   1,
		UserContext: UserContext{
			// Would be populated with actual user context
		},
		SystemInfo: SystemInfo{
			OS:           runtime.GOOS,
			Architecture: runtime.GOARCH,
			GoVersion:    runtime.Version(),
			AppVersion:   "1.0.0",
			MemoryMB:     float64(getMemoryUsage()) / 1024 / 1024,
			CPUUsage:     getCPUUsage(),
		},
	}

	if err != nil {
		report.Details = err.Error()
		report.StackTrace = fmt.Sprintf("%+v", err)
	}

	// Log the error
	level := LogLevelError
	if severity == ErrorSeverityCritical {
		level = LogLevelFatal
	}

	als.LogEntry(level, component, operation, message,
		WithLogDetails(report.Details),
		WithMetadata(map[string]interface{}{
			"errorId":   report.ID,
			"errorType": errorType,
			"severity":  severity,
		}))

	return report
}

// GetErrorReports retrieves error reports (simplified implementation)
func (als *ActivityLogService) GetErrorReports(includeResolved bool) []ErrorReport {
	// This is a simplified implementation
	// In a full implementation, we'd maintain a separate error reports collection

	errorLevel := LogLevelError
	limit := 50
	filter := LogFilter{
		Level: &errorLevel,
		Limit: &limit,
	}

	entries := als.GetLogs(filter)
	reports := make([]ErrorReport, 0, len(entries))

	for _, entry := range entries {
		if entry.Level == LogLevelError || entry.Level == LogLevelFatal {
			report := ErrorReport{
				ID:        entry.ID,
				Timestamp: entry.Timestamp,
				Component: entry.Component,
				Operation: entry.Operation,
				Message:   entry.Message,
				Details:   entry.Details,
				Frequency: 1,
				FirstSeen: entry.Timestamp,
				LastSeen:  entry.Timestamp,
				Type:      ErrorTypeSystemErr,  // Simplified
				Severity:  ErrorSeverityMedium, // Simplified
			}
			reports = append(reports, report)
		}
	}

	return reports
}

// LogEntry adds a new entry to the activity log
func (als *ActivityLogService) LogEntry(level LogLevel, component, operation, message string, options ...LogEntryOption) {
	if !als.shouldLog(level) {
		return
	}

	entry := ActivityLogEntry{
		ID:        generateLogID(),
		Timestamp: time.Now(),
		Level:     level,
		Component: component,
		Operation: operation,
		Message:   message,
	}

	// Apply options
	for _, option := range options {
		option(&entry)
	}

	als.addEntry(entry)
}

// LogEntryOption allows customization of log entries
type LogEntryOption func(*ActivityLogEntry)

// WithLogDetails adds details to the log entry
func WithLogDetails(details string) LogEntryOption {
	return func(entry *ActivityLogEntry) {
		entry.Details = details
	}
}

// WithDuration adds duration to the log entry
func WithDuration(duration time.Duration) LogEntryOption {
	return func(entry *ActivityLogEntry) {
		entry.Duration = &duration
	}
}

// WithProjectID adds project ID to the log entry
func WithProjectID(projectID string) LogEntryOption {
	return func(entry *ActivityLogEntry) {
		entry.ProjectID = projectID
	}
}

// WithUserAction marks the entry as a user action
func WithUserAction(userAction bool) LogEntryOption {
	return func(entry *ActivityLogEntry) {
		entry.UserAction = userAction
	}
}

// WithMetadata adds metadata to the log entry
func WithMetadata(metadata map[string]interface{}) LogEntryOption {
	return func(entry *ActivityLogEntry) {
		entry.Metadata = metadata
	}
}

// GetLogs retrieves logs based on filter criteria
func (als *ActivityLogService) GetLogs(filter LogFilter) []ActivityLogEntry {
	als.mutex.RLock()
	defer als.mutex.RUnlock()

	var entries []ActivityLogEntry

	// Get all entries in chronological order
	if als.full {
		// Buffer is full, start from writeIdx and wrap around
		for i := 0; i < len(als.buffer); i++ {
			idx := (als.writeIdx + i) % len(als.buffer)
			entries = append(entries, als.buffer[idx])
		}
	} else {
		// Buffer not full, just take entries up to writeIdx
		entries = make([]ActivityLogEntry, als.writeIdx)
		copy(entries, als.buffer[:als.writeIdx])
	}

	// Apply filters
	filtered := als.applyFilter(entries, filter)

	// Sort by timestamp (newest first)
	sort.Slice(filtered, func(i, j int) bool {
		return filtered[i].Timestamp.After(filtered[j].Timestamp)
	})

	// Apply limit
	if filter.Limit != nil && len(filtered) > *filter.Limit {
		filtered = filtered[:*filter.Limit]
	}

	return filtered
}

// SearchLogs performs a search across log entries
func (als *ActivityLogService) SearchLogs(ctx context.Context, request LogSearchRequest) (*LogSearchResult, error) {
	startTime := time.Now()

	entries := als.GetLogs(request.Filter)

	if request.Query != "" {
		entries = als.searchEntries(entries, request.Query)
	}

	total := len(entries)
	hasMore := false

	// Apply pagination
	if request.Offset > 0 {
		if request.Offset >= len(entries) {
			entries = []ActivityLogEntry{}
		} else {
			entries = entries[request.Offset:]
		}
	}

	if request.Limit > 0 && len(entries) > request.Limit {
		entries = entries[:request.Limit]
		hasMore = true
	}

	return &LogSearchResult{
		Entries:    entries,
		Total:      total,
		HasMore:    hasMore,
		SearchTime: time.Since(startTime),
	}, nil
}

// ExportLogs exports logs to a file in the specified format
func (als *ActivityLogService) ExportLogs(ctx context.Context, request LogExportRequest) (*LogExportResult, error) {
	startTime := time.Now()

	entries := als.GetLogs(request.Filter)

	var fileSize int64
	var err error

	switch strings.ToLower(request.Format) {
	case "json":
		fileSize, err = als.exportToJSON(entries, request.FilePath)
	case "csv":
		fileSize, err = als.exportToCSV(entries, request.FilePath)
	case "txt":
		fileSize, err = als.exportToText(entries, request.FilePath)
	default:
		return nil, fmt.Errorf("unsupported export format: %s", request.Format)
	}

	if err != nil {
		return nil, fmt.Errorf("export failed: %w", err)
	}

	return &LogExportResult{
		FilePath:     request.FilePath,
		EntriesCount: len(entries),
		FileSize:     fileSize,
		ExportTime:   time.Since(startTime),
		Format:       request.Format,
	}, nil
}

// GetApplicationStatus returns the current application status with system health
func (als *ActivityLogService) GetApplicationStatus() *ApplicationStatus {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	return &ApplicationStatus{
		Status:           als.determineStatus(),
		Message:          als.getStatusMessage(),
		ActiveOperations: als.getActiveOperations(),
		LastUpdate:       time.Now(),
		SystemHealth: SystemHealth{
			MemoryUsage:       float64(m.Alloc) / 1024 / 1024, // MB
			CPUUsage:          getCPUUsage(),
			DiskSpace:         als.getDiskSpace(),
			DatabaseSize:      als.getDatabaseSize(),
			TemporaryFiles:    als.getTemporaryFilesCount(),
			ActiveConnections: 1, // Simplified for now
		},
	}
}

// UpdateLogConfig updates the logging configuration
func (als *ActivityLogService) UpdateLogConfig(config LogConfig) error {
	als.mutex.Lock()
	defer als.mutex.Unlock()

	// If buffer size changed, recreate buffer
	if config.BufferSize != als.config.BufferSize {
		als.resizeBuffer(config.BufferSize)
	}

	als.config = config

	// Restart flush timer if needed
	if als.ticker != nil {
		als.ticker.Stop()
	}

	if config.FlushInterval > 0 {
		als.ticker = time.NewTicker(config.FlushInterval)
		go als.flushLoop()
	}

	als.LogEntry(LogLevelInfo, "ActivityLog", "UpdateConfig", "Log configuration updated",
		WithMetadata(map[string]interface{}{
			"bufferSize":    config.BufferSize,
			"retentionDays": config.RetentionDays,
			"level":         config.Level.String(),
		}))

	return nil
}

// ClearLogs clears logs based on age or all logs
func (als *ActivityLogService) ClearLogs(olderThan time.Duration) int {
	als.mutex.Lock()
	defer als.mutex.Unlock()

	if olderThan == 0 {
		// Clear all logs
		cleared := als.getLogCount()
		als.buffer = make([]ActivityLogEntry, len(als.buffer))
		als.writeIdx = 0
		als.full = false
		return cleared
	}

	// This is simplified - in a real implementation we'd need more complex logic
	// to remove entries from the circular buffer while maintaining order
	return 0
}

// Close shuts down the activity log service
func (als *ActivityLogService) Close() {
	if als.cancel != nil {
		als.cancel()
	}
	if als.ticker != nil {
		als.ticker.Stop()
	}
}

// Internal methods

func (als *ActivityLogService) shouldLog(level LogLevel) bool {
	return levelToInt(level) >= levelToInt(als.config.Level)
}

func levelToInt(level LogLevel) int {
	switch level {
	case LogLevelDebug:
		return 0
	case LogLevelInfo:
		return 1
	case LogLevelWarn:
		return 2
	case LogLevelError:
		return 3
	case LogLevelFatal:
		return 4
	default:
		return 1
	}
}

func (als *ActivityLogService) addEntry(entry ActivityLogEntry) {
	als.mutex.Lock()
	defer als.mutex.Unlock()

	als.buffer[als.writeIdx] = entry
	als.writeIdx = (als.writeIdx + 1) % len(als.buffer)

	if als.writeIdx == 0 {
		als.full = true
	}

	// Emit real-time event if app context is available
	if als.app != nil && als.app.ctx != nil {
		// Emit event using Wails runtime (simplified for now)
		// In a full implementation, this would use the Wails events system
	}

	// Log to console if enabled
	if als.config.EnableConsole {
		als.logToConsole(entry)
	}
}

func (als *ActivityLogService) applyFilter(entries []ActivityLogEntry, filter LogFilter) []ActivityLogEntry {
	if isEmptyFilter(filter) {
		return entries
	}

	filtered := make([]ActivityLogEntry, 0, len(entries))

	for _, entry := range entries {
		if als.matchesFilter(entry, filter) {
			filtered = append(filtered, entry)
		}
	}

	return filtered
}

func (als *ActivityLogService) matchesFilter(entry ActivityLogEntry, filter LogFilter) bool {
	if filter.Level != nil && entry.Level != *filter.Level {
		return false
	}

	if filter.Component != nil && !strings.Contains(strings.ToLower(entry.Component), strings.ToLower(*filter.Component)) {
		return false
	}

	if filter.Operation != nil && !strings.Contains(strings.ToLower(entry.Operation), strings.ToLower(*filter.Operation)) {
		return false
	}

	if filter.ProjectID != nil && entry.ProjectID != *filter.ProjectID {
		return false
	}

	if filter.UserAction != nil && entry.UserAction != *filter.UserAction {
		return false
	}

	if filter.StartTime != nil && entry.Timestamp.Before(*filter.StartTime) {
		return false
	}

	if filter.EndTime != nil && entry.Timestamp.After(*filter.EndTime) {
		return false
	}

	if filter.Search != nil {
		searchTerm := strings.ToLower(*filter.Search)
		if !strings.Contains(strings.ToLower(entry.Message), searchTerm) &&
			!strings.Contains(strings.ToLower(entry.Details), searchTerm) {
			return false
		}
	}

	return true
}

func isEmptyFilter(filter LogFilter) bool {
	return filter.Level == nil && filter.Component == nil && filter.Operation == nil &&
		filter.ProjectID == nil && filter.UserAction == nil && filter.StartTime == nil &&
		filter.EndTime == nil && filter.Search == nil
}

func (als *ActivityLogService) searchEntries(entries []ActivityLogEntry, query string) []ActivityLogEntry {
	query = strings.ToLower(query)
	filtered := make([]ActivityLogEntry, 0, len(entries))

	for _, entry := range entries {
		if strings.Contains(strings.ToLower(entry.Message), query) ||
			strings.Contains(strings.ToLower(entry.Details), query) ||
			strings.Contains(strings.ToLower(entry.Component), query) ||
			strings.Contains(strings.ToLower(entry.Operation), query) {
			filtered = append(filtered, entry)
		}
	}

	return filtered
}

func (als *ActivityLogService) exportToJSON(entries []ActivityLogEntry, filePath string) (int64, error) {
	file, err := os.Create(filePath)
	if err != nil {
		return 0, err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")

	if err := encoder.Encode(entries); err != nil {
		return 0, err
	}

	stat, err := file.Stat()
	if err != nil {
		return 0, err
	}

	return stat.Size(), nil
}

func (als *ActivityLogService) exportToCSV(entries []ActivityLogEntry, filePath string) (int64, error) {
	file, err := os.Create(filePath)
	if err != nil {
		return 0, err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Write header
	header := []string{"ID", "Timestamp", "Level", "Component", "Operation", "Message", "Details", "Duration", "ProjectID", "UserAction"}
	if err := writer.Write(header); err != nil {
		return 0, err
	}

	// Write entries
	for _, entry := range entries {
		durationStr := ""
		if entry.Duration != nil {
			durationStr = entry.Duration.String()
		}

		record := []string{
			entry.ID,
			entry.Timestamp.Format(time.RFC3339),
			entry.Level.String(),
			entry.Component,
			entry.Operation,
			entry.Message,
			entry.Details,
			durationStr,
			entry.ProjectID,
			fmt.Sprintf("%t", entry.UserAction),
		}

		if err := writer.Write(record); err != nil {
			return 0, err
		}
	}

	stat, err := file.Stat()
	if err != nil {
		return 0, err
	}

	return stat.Size(), nil
}

func (als *ActivityLogService) exportToText(entries []ActivityLogEntry, filePath string) (int64, error) {
	file, err := os.Create(filePath)
	if err != nil {
		return 0, err
	}
	defer file.Close()

	for _, entry := range entries {
		line := fmt.Sprintf("%s [%s] %s/%s: %s",
			entry.Timestamp.Format("2006-01-02 15:04:05"),
			strings.ToUpper(entry.Level.String()),
			entry.Component,
			entry.Operation,
			entry.Message)

		if entry.Details != "" {
			line += fmt.Sprintf(" - %s", entry.Details)
		}

		line += "\n"

		if _, err := file.WriteString(line); err != nil {
			return 0, err
		}
	}

	stat, err := file.Stat()
	if err != nil {
		return 0, err
	}

	return stat.Size(), nil
}

func (als *ActivityLogService) resizeBuffer(newSize int) {
	if newSize <= 0 {
		return
	}

	newBuffer := make([]ActivityLogEntry, newSize)

	// Copy existing entries to new buffer
	currentEntries := als.GetLogs(LogFilter{})
	copyCount := len(currentEntries)
	if copyCount > newSize {
		copyCount = newSize
		currentEntries = currentEntries[:copyCount]
	}

	copy(newBuffer, currentEntries)

	als.buffer = newBuffer
	als.writeIdx = copyCount % newSize
	als.full = copyCount == newSize
}

func (als *ActivityLogService) getLogCount() int {
	if als.full {
		return len(als.buffer)
	}
	return als.writeIdx
}

func (als *ActivityLogService) flushLoop() {
	for {
		select {
		case <-als.ctx.Done():
			return
		case <-als.ticker.C:
			// Periodic operations like cleanup could go here
			als.performMaintenance()
		}
	}
}

func (als *ActivityLogService) performMaintenance() {
	// Clean up old entries based on retention policy
	if als.config.RetentionDays > 0 {
		cutoff := time.Now().AddDate(0, 0, -als.config.RetentionDays)
		als.ClearLogs(time.Since(cutoff))
	}
}

func (als *ActivityLogService) logToConsole(entry ActivityLogEntry) {
	levelStr := strings.ToUpper(entry.Level.String())
	fmt.Printf("%s [%s] %s/%s: %s\n",
		entry.Timestamp.Format("15:04:05"),
		levelStr,
		entry.Component,
		entry.Operation,
		entry.Message)
}

func (als *ActivityLogService) determineStatus() StatusLevel {
	// Check recent error entries to determine status
	errorLevel := LogLevelError
	startTime := time.Now().Add(-5 * time.Minute)
	limit := 1
	filter := LogFilter{
		Level:     &errorLevel,
		StartTime: &startTime,
		Limit:     &limit,
	}

	errorEntries := als.GetLogs(filter)
	if len(errorEntries) > 0 {
		return StatusError
	}

	// Check for warnings
	warnLevel := LogLevelWarn
	filter.Level = &warnLevel
	warnEntries := als.GetLogs(filter)
	if len(warnEntries) > 0 {
		return StatusWarning
	}

	// Check for active operations
	if als.getActiveOperations() > 0 {
		return StatusWorking
	}

	return StatusIdle
}

func (als *ActivityLogService) getStatusMessage() string {
	status := als.determineStatus()
	switch status {
	case StatusError:
		return "System errors detected"
	case StatusWarning:
		return "System warnings present"
	case StatusWorking:
		return "Operations in progress"
	default:
		return "System idle"
	}
}

func (als *ActivityLogService) getActiveOperations() int {
	// Simplified - count recent user actions as active operations
	userAction := true
	startTime := time.Now().Add(-1 * time.Minute)
	filter := LogFilter{
		UserAction: &userAction,
		StartTime:  &startTime,
	}

	return len(als.GetLogs(filter))
}

func (als *ActivityLogService) getDiskSpace() float64 {
	// Simplified implementation - would need platform-specific code
	return 50.0 // GB
}

func (als *ActivityLogService) getDatabaseSize() float64 {
	// Would query actual database size
	return 10.0 // MB
}

func (als *ActivityLogService) getTemporaryFilesCount() int {
	// Would count temporary files in temp directory
	return 0
}

func generateLogID() string {
	bytes := make([]byte, 8)
	rand.Read(bytes)
	return fmt.Sprintf("log_%x", bytes)
}
