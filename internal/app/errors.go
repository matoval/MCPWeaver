package app

import (
	"fmt"
	"math/rand"
	"runtime"
	"time"
)

// ErrorManager handles error creation, logging, and recovery
type ErrorManager struct {
	defaultRetryPolicy RetryPolicy
	correlationIDGen   func() string
}

// NewErrorManager creates a new error manager
func NewErrorManager() *ErrorManager {
	return &ErrorManager{
		defaultRetryPolicy: DefaultRetryPolicy(),
		correlationIDGen:   generateCorrelationID,
	}
}

// CreateError creates a new APIError with comprehensive context
func (em *ErrorManager) CreateError(errType, code, message string, options ...ErrorOption) *APIError {
	apiError := &APIError{
		Type:          errType,
		Code:          code,
		Message:       message,
		Timestamp:     time.Now(),
		CorrelationID: em.correlationIDGen(),
		Severity:      ErrorSeverityMedium,
		Recoverable:   false,
		Details:       make(map[string]string),
	}

	// Apply options
	for _, option := range options {
		option(apiError)
	}

	// Set default retry policy for retryable errors
	if apiError.Recoverable && apiError.RetryAfter == nil {
		delay := em.defaultRetryPolicy.InitialDelay
		apiError.RetryAfter = &delay
	}

	return apiError
}

// CreateValidationError creates a validation error with suggestions
func (em *ErrorManager) CreateValidationError(message string, details map[string]string, suggestions []string) *APIError {
	return em.CreateError(
		ErrorTypeValidation,
		ErrCodeValidation,
		message,
		WithDetails(details),
		WithSuggestions(suggestions),
		WithSeverity(ErrorSeverityMedium),
	)
}

// CreateNetworkError creates a network error that can be retried
func (em *ErrorManager) CreateNetworkError(message string, details map[string]string) *APIError {
	return em.CreateError(
		ErrorTypeNetwork,
		ErrCodeNetworkError,
		message,
		WithDetails(details),
		WithSeverity(ErrorSeverityHigh),
		WithRecoverable(true),
		WithSuggestions([]string{
			"Check your internet connection",
			"Try again in a few moments",
			"Verify the service is available",
		}),
	)
}

// CreateFileSystemError creates a file system error with recovery guidance
func (em *ErrorManager) CreateFileSystemError(message string, filePath string, operation string) *APIError {
	details := map[string]string{
		"file_path":  filePath,
		"operation":  operation,
	}
	
	suggestions := []string{
		"Check that the file or directory exists",
		"Verify you have the necessary permissions",
		"Ensure the path is correct",
	}

	return em.CreateError(
		ErrorTypeFileSystem,
		ErrCodeFileAccess,
		message,
		WithDetails(details),
		WithSuggestions(suggestions),
		WithSeverity(ErrorSeverityMedium),
	)
}

// CreateGenerationError creates a generation error with context
func (em *ErrorManager) CreateGenerationError(message string, projectID string, step string) *APIError {
	details := map[string]string{
		"project_id": projectID,
		"step":       step,
	}
	
	suggestions := []string{
		"Check the OpenAPI specification for errors",
		"Verify the project configuration",
		"Try regenerating the server",
	}

	return em.CreateError(
		ErrorTypeGeneration,
		ErrCodeGenerationError,
		message,
		WithDetails(details),
		WithSuggestions(suggestions),
		WithSeverity(ErrorSeverityHigh),
		WithContext(&ErrorContext{
			Operation: "generation",
			Component: "generator",
			ProjectID: projectID,
		}),
	)
}

// CreateInternalError creates an internal error with debug information
func (em *ErrorManager) CreateInternalError(message string, err error) *APIError {
	details := map[string]string{
		"internal_error": err.Error(),
	}
	
	// Get stack trace
	stackTrace := getStackTrace(2) // Skip this function and the caller
	
	suggestions := []string{
		"This is an internal error. Please contact support if it persists",
		"Try refreshing the page or restarting the application",
	}

	return em.CreateError(
		ErrorTypeSystem,
		ErrCodeInternalError,
		message,
		WithDetails(details),
		WithSuggestions(suggestions),
		WithSeverity(ErrorSeverityCritical),
		WithRecoverable(true),
		WithContext(&ErrorContext{
			Operation:  "internal",
			Component:  "system",
			StackTrace: stackTrace,
		}),
	)
}

// CreateDatabaseError creates a database error with recovery options
func (em *ErrorManager) CreateDatabaseError(message string, operation string, table string) *APIError {
	details := map[string]string{
		"operation": operation,
		"table":     table,
	}
	
	suggestions := []string{
		"Check database connectivity",
		"Verify the operation is valid",
		"Try again in a few moments",
	}

	return em.CreateError(
		ErrorTypeDatabase,
		ErrCodeDatabaseError,
		message,
		WithDetails(details),
		WithSuggestions(suggestions),
		WithSeverity(ErrorSeverityHigh),
		WithRecoverable(true),
	)
}

// CreateErrorCollection creates a collection for batch operations
func (em *ErrorManager) CreateErrorCollection(operation string, totalItems int) *ErrorCollection {
	return &ErrorCollection{
		Errors:      []APIError{},
		Warnings:    []APIError{},
		Operation:   operation,
		TotalItems:  totalItems,
		FailedItems: 0,
		Timestamp:   time.Now(),
	}
}

// ErrorOption is a function that configures an APIError
type ErrorOption func(*APIError)

// WithDetails adds details to an error
func WithDetails(details map[string]string) ErrorOption {
	return func(e *APIError) {
		if e.Details == nil {
			e.Details = make(map[string]string)
		}
		for k, v := range details {
			e.Details[k] = v
		}
	}
}

// WithSuggestions adds suggestions to an error
func WithSuggestions(suggestions []string) ErrorOption {
	return func(e *APIError) {
		e.Suggestions = suggestions
	}
}

// WithSeverity sets the severity of an error
func WithSeverity(severity ErrorSeverity) ErrorOption {
	return func(e *APIError) {
		e.Severity = severity
	}
}

// WithRecoverable marks an error as recoverable
func WithRecoverable(recoverable bool) ErrorOption {
	return func(e *APIError) {
		e.Recoverable = recoverable
	}
}

// WithRetryAfter sets the retry delay for an error
func WithRetryAfter(delay time.Duration) ErrorOption {
	return func(e *APIError) {
		e.RetryAfter = &delay
	}
}

// WithContext adds context to an error
func WithContext(context *ErrorContext) ErrorOption {
	return func(e *APIError) {
		e.Context = context
	}
}

// WithCorrelationID sets a custom correlation ID
func WithCorrelationID(correlationID string) ErrorOption {
	return func(e *APIError) {
		e.CorrelationID = correlationID
	}
}

// Helper functions

// generateCorrelationID generates a unique correlation ID
func generateCorrelationID() string {
	return fmt.Sprintf("mcpw-%d-%d", time.Now().UnixNano(), rand.Intn(1000))
}

// getStackTrace gets the current stack trace
func getStackTrace(skip int) string {
	const depth = 32
	var pcs [depth]uintptr
	n := runtime.Callers(skip+1, pcs[:])
	
	var trace string
	frames := runtime.CallersFrames(pcs[:n])
	for {
		frame, more := frames.Next()
		trace += fmt.Sprintf("%s:%d %s\n", frame.File, frame.Line, frame.Function)
		if !more {
			break
		}
	}
	return trace
}

// IsRetryableError checks if an error can be retried
func IsRetryableError(err error) bool {
	if apiErr, ok := err.(*APIError); ok {
		return apiErr.IsRetryable()
	}
	return false
}

// GetErrorSeverity returns the severity of an error
func GetErrorSeverity(err error) ErrorSeverity {
	if apiErr, ok := err.(*APIError); ok {
		return apiErr.Severity
	}
	return ErrorSeverityMedium
}

// GetErrorSuggestions returns suggestions for an error
func GetErrorSuggestions(err error) []string {
	if apiErr, ok := err.(*APIError); ok {
		return apiErr.Suggestions
	}
	return []string{}
}

// GetErrorContext returns the context of an error
func GetErrorContext(err error) *ErrorContext {
	if apiErr, ok := err.(*APIError); ok {
		return apiErr.Context
	}
	return nil
}

// CategorizeError categorizes an error based on its type and code
func CategorizeError(err error) string {
	if apiErr, ok := err.(*APIError); ok {
		switch apiErr.Type {
		case ErrorTypeValidation:
			return "User Input Error"
		case ErrorTypeNetwork:
			return "Network Error"
		case ErrorTypeFileSystem:
			return "File System Error"
		case ErrorTypeDatabase:
			return "Database Error"
		case ErrorTypeGeneration:
			return "Generation Error"
		case ErrorTypeSystem:
			return "System Error"
		case ErrorTypePermission:
			return "Permission Error"
		case ErrorTypeConfiguration:
			return "Configuration Error"
		case ErrorTypeAuthentication:
			return "Authentication Error"
		default:
			return "Unknown Error"
		}
	}
	return "Unknown Error"
}