package app

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"
)

// LogLevel represents the severity level of a log entry
type LogLevel int

const (
	LogLevelDebug LogLevel = iota
	LogLevelInfo
	LogLevelWarn
	LogLevelError
	LogLevelFatal
)

// String returns the string representation of the log level
func (l LogLevel) String() string {
	switch l {
	case LogLevelDebug:
		return "DEBUG"
	case LogLevelInfo:
		return "INFO"
	case LogLevelWarn:
		return "WARN"
	case LogLevelError:
		return "ERROR"
	case LogLevelFatal:
		return "FATAL"
	default:
		return "UNKNOWN"
	}
}

// LogEntry represents a structured log entry
type LogEntry struct {
	Timestamp     time.Time         `json:"timestamp"`
	Level         LogLevel          `json:"level"`
	Message       string            `json:"message"`
	Context       map[string]string `json:"context,omitempty"`
	Error         string            `json:"error,omitempty"`
	StackTrace    string            `json:"stackTrace,omitempty"`
	CorrelationID string            `json:"correlationId,omitempty"`
	Component     string            `json:"component,omitempty"`
	Operation     string            `json:"operation,omitempty"`
	Duration      time.Duration     `json:"duration,omitempty"`
	UserID        string            `json:"userId,omitempty"`
	ProjectID     string            `json:"projectId,omitempty"`
	SessionID     string            `json:"sessionId,omitempty"`
	RequestID     string            `json:"requestId,omitempty"`
}

// Logger provides structured logging capabilities
type Logger struct {
	level      LogLevel
	outputs    []io.Writer
	context    map[string]string
	component  string
	jsonFormat bool
}

// NewLogger creates a new logger instance
func NewLogger(level LogLevel, component string) *Logger {
	return &Logger{
		level:      level,
		outputs:    []io.Writer{os.Stdout},
		context:    make(map[string]string),
		component:  component,
		jsonFormat: true,
	}
}

// NewFileLogger creates a logger that writes to a file
func NewFileLogger(level LogLevel, component string, logFile string) (*Logger, error) {
	// Create log directory if it doesn't exist
	logDir := filepath.Dir(logFile)
	if err := os.MkdirAll(logDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create log directory: %w", err)
	}

	// Open log file
	file, err := os.OpenFile(logFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return nil, fmt.Errorf("failed to open log file: %w", err)
	}

	return &Logger{
		level:      level,
		outputs:    []io.Writer{file, os.Stdout},
		context:    make(map[string]string),
		component:  component,
		jsonFormat: true,
	}, nil
}

// WithContext adds context to the logger
func (l *Logger) WithContext(key, value string) *Logger {
	newLogger := *l
	newLogger.context = make(map[string]string)
	for k, v := range l.context {
		newLogger.context[k] = v
	}
	newLogger.context[key] = value
	return &newLogger
}

// WithCorrelationID adds a correlation ID to the logger
func (l *Logger) WithCorrelationID(correlationID string) *Logger {
	return l.WithContext("correlationId", correlationID)
}

// WithProjectID adds a project ID to the logger
func (l *Logger) WithProjectID(projectID string) *Logger {
	return l.WithContext("projectId", projectID)
}

// WithUserID adds a user ID to the logger
func (l *Logger) WithUserID(userID string) *Logger {
	return l.WithContext("userId", userID)
}

// WithOperation adds an operation name to the logger
func (l *Logger) WithOperation(operation string) *Logger {
	return l.WithContext("operation", operation)
}

// Debug logs a debug message
func (l *Logger) Debug(message string, args ...interface{}) {
	if l.level <= LogLevelDebug {
		l.log(LogLevelDebug, message, args...)
	}
}

// Info logs an info message
func (l *Logger) Info(message string, args ...interface{}) {
	if l.level <= LogLevelInfo {
		l.log(LogLevelInfo, message, args...)
	}
}

// Warn logs a warning message
func (l *Logger) Warn(message string, args ...interface{}) {
	if l.level <= LogLevelWarn {
		l.log(LogLevelWarn, message, args...)
	}
}

// Error logs an error message
func (l *Logger) Error(message string, args ...interface{}) {
	if l.level <= LogLevelError {
		l.log(LogLevelError, message, args...)
	}
}

// Fatal logs a fatal message and exits
func (l *Logger) Fatal(message string, args ...interface{}) {
	l.log(LogLevelFatal, message, args...)
	os.Exit(1)
}

// LogError logs an APIError with full context
func (l *Logger) LogError(err *APIError) {
	entry := LogEntry{
		Timestamp:     time.Now(),
		Level:         LogLevelError,
		Message:       err.Message,
		Error:         err.Error(),
		CorrelationID: err.CorrelationID,
		Component:     l.component,
		Context:       make(map[string]string),
	}

	// Add error context
	if err.Context != nil {
		entry.Operation = err.Context.Operation
		entry.ProjectID = err.Context.ProjectID
		entry.UserID = err.Context.UserID
		entry.SessionID = err.Context.SessionID
		entry.RequestID = err.Context.RequestID
		entry.StackTrace = err.Context.StackTrace
		
		// Add metadata to context
		for k, v := range err.Context.Metadata {
			entry.Context[k] = v
		}
	}

	// Add error details to context
	for k, v := range err.Details {
		entry.Context[k] = v
	}

	// Add logger context
	for k, v := range l.context {
		entry.Context[k] = v
	}

	l.writeEntry(entry)
}

// LogOperation logs the start and end of an operation
func (l *Logger) LogOperation(operation string, fn func() error) error {
	startTime := time.Now()
	
	l.WithOperation(operation).Info("Operation started")
	
	err := fn()
	
	duration := time.Since(startTime)
	
	if err != nil {
		l.WithOperation(operation).WithContext("duration", duration.String()).Error("Operation failed: %v", err)
	} else {
		l.WithOperation(operation).WithContext("duration", duration.String()).Info("Operation completed")
	}
	
	return err
}

// LogRetry logs retry attempts
func (l *Logger) LogRetry(operation string, attempt int, maxAttempts int, err error, delay time.Duration) {
	l.WithOperation(operation).WithContext("attempt", fmt.Sprintf("%d/%d", attempt, maxAttempts)).
		WithContext("delay", delay.String()).
		Warn("Operation failed, retrying: %v", err)
}

// log is the internal logging method
func (l *Logger) log(level LogLevel, message string, args ...interface{}) {
	if len(args) > 0 {
		message = fmt.Sprintf(message, args...)
	}

	entry := LogEntry{
		Timestamp: time.Now(),
		Level:     level,
		Message:   message,
		Component: l.component,
		Context:   make(map[string]string),
	}

	// Add logger context
	for k, v := range l.context {
		entry.Context[k] = v
	}

	// Add stack trace for errors
	if level >= LogLevelError {
		entry.StackTrace = getStackTrace(3) // Skip this method, log method, and caller
	}

	l.writeEntry(entry)
}

// writeEntry writes a log entry to all outputs
func (l *Logger) writeEntry(entry LogEntry) {
	var output string
	
	if l.jsonFormat {
		jsonData, err := json.Marshal(entry)
		if err != nil {
			// Fallback to simple format if JSON marshaling fails
			output = fmt.Sprintf("[%s] %s %s: %s\n",
				entry.Timestamp.Format(time.RFC3339),
				entry.Level.String(),
				entry.Component,
				entry.Message)
		} else {
			output = string(jsonData) + "\n"
		}
	} else {
		// Human-readable format
		output = fmt.Sprintf("[%s] %s %s: %s",
			entry.Timestamp.Format("2006-01-02 15:04:05"),
			entry.Level.String(),
			entry.Component,
			entry.Message)
		
		if len(entry.Context) > 0 {
			output += " {"
			first := true
			for k, v := range entry.Context {
				if !first {
					output += ", "
				}
				output += fmt.Sprintf("%s: %s", k, v)
				first = false
			}
			output += "}"
		}
		output += "\n"
	}

	for _, writer := range l.outputs {
		writer.Write([]byte(output))
	}
}

// SetJSONFormat sets whether to output logs in JSON format
func (l *Logger) SetJSONFormat(jsonFormat bool) {
	l.jsonFormat = jsonFormat
}

// SetLevel sets the minimum log level
func (l *Logger) SetLevel(level LogLevel) {
	l.level = level
}

// AddOutput adds an output writer
func (l *Logger) AddOutput(writer io.Writer) {
	l.outputs = append(l.outputs, writer)
}

// ErrorReporter handles error reporting to external services
type ErrorReporter struct {
	logger    *Logger
	enabled   bool
	endpoints []string
}

// NewErrorReporter creates a new error reporter
func NewErrorReporter(logger *Logger) *ErrorReporter {
	return &ErrorReporter{
		logger:  logger,
		enabled: true,
	}
}

// ReportError reports an error to external services
func (er *ErrorReporter) ReportError(err *APIError) {
	if !er.enabled {
		return
	}

	// Log the error locally
	er.logger.LogError(err)

	// Here you could send the error to external services like:
	// - Sentry
	// - Rollbar
	// - Custom monitoring endpoints
	// - Email notifications
	// - Slack webhooks
	
	// For now, we'll just log that we would report it
	er.logger.Info("Error reported to monitoring services", 
		"correlationId", err.CorrelationID,
		"errorType", err.Type,
		"errorCode", err.Code)
}

// Enable enables error reporting
func (er *ErrorReporter) Enable() {
	er.enabled = true
}

// Disable disables error reporting
func (er *ErrorReporter) Disable() {
	er.enabled = false
}

// Note: getStackTrace is defined in errors.go