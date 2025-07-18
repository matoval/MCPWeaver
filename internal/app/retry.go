package app

import (
	"context"
	"fmt"
	"math/rand"
	"strings"
	"time"
)

// RetryManager handles retry logic with exponential backoff
type RetryManager struct {
	defaultPolicy RetryPolicy
	errorManager  *ErrorManager
}

// NewRetryManager creates a new retry manager
func NewRetryManager(errorManager *ErrorManager) *RetryManager {
	return &RetryManager{
		defaultPolicy: DefaultRetryPolicy(),
		errorManager:  errorManager,
	}
}

// RetryResult contains the result of a retry operation
type RetryResult struct {
	Success     bool          `json:"success"`
	Attempts    int           `json:"attempts"`
	LastError   error         `json:"lastError,omitempty"`
	TotalDelay  time.Duration `json:"totalDelay"`
	StartTime   time.Time     `json:"startTime"`
	EndTime     time.Time     `json:"endTime"`
}

// RetryFunc is a function that can be retried
type RetryFunc func() error

// RetryWithPolicy executes a function with retry logic using the specified policy
func (rm *RetryManager) RetryWithPolicy(ctx context.Context, policy RetryPolicy, fn RetryFunc) *RetryResult {
	result := &RetryResult{
		Success:   false,
		Attempts:  0,
		StartTime: time.Now(),
	}

	var lastError error
	delay := policy.InitialDelay

	for attempt := 0; attempt <= policy.MaxRetries; attempt++ {
		result.Attempts++

		// Execute the function
		err := fn()
		if err == nil {
			result.Success = true
			result.EndTime = time.Now()
			return result
		}

		lastError = err

		// Check if error is retryable
		if !rm.isRetryableError(err, policy) {
			break
		}

		// Don't sleep after the last attempt
		if attempt < policy.MaxRetries {
			// Check context cancellation
			if ctx.Err() != nil {
				lastError = ctx.Err()
				break
			}

			// Calculate delay with exponential backoff
			actualDelay := rm.calculateDelay(delay, policy)
			result.TotalDelay += actualDelay

			// Sleep with context cancellation support
			select {
			case <-ctx.Done():
				lastError = ctx.Err()
				goto done
			case <-time.After(actualDelay):
			}

			// Increase delay for next attempt
			delay = time.Duration(float64(delay) * policy.BackoffMultiplier)
			if delay > policy.MaxDelay {
				delay = policy.MaxDelay
			}
		}
	}

done:
	result.LastError = lastError
	result.EndTime = time.Now()
	return result
}

// Retry executes a function with default retry policy
func (rm *RetryManager) Retry(ctx context.Context, fn RetryFunc) *RetryResult {
	return rm.RetryWithPolicy(ctx, rm.defaultPolicy, fn)
}

// RetryWithTimeout executes a function with retry logic and timeout
func (rm *RetryManager) RetryWithTimeout(timeout time.Duration, fn RetryFunc) *RetryResult {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	
	return rm.Retry(ctx, fn)
}

// RetryAsync executes a function with retry logic asynchronously
func (rm *RetryManager) RetryAsync(ctx context.Context, fn RetryFunc, callback func(*RetryResult)) {
	go func() {
		result := rm.Retry(ctx, fn)
		if callback != nil {
			callback(result)
		}
	}()
}

// RetryWithBackoff creates a custom retry policy with exponential backoff
func (rm *RetryManager) RetryWithBackoff(ctx context.Context, maxRetries int, initialDelay time.Duration, maxDelay time.Duration, fn RetryFunc) *RetryResult {
	policy := RetryPolicy{
		MaxRetries:        maxRetries,
		InitialDelay:      initialDelay,
		MaxDelay:          maxDelay,
		BackoffMultiplier: 2.0,
		JitterEnabled:     true,
		RetryableErrors:   rm.defaultPolicy.RetryableErrors,
	}
	
	return rm.RetryWithPolicy(ctx, policy, fn)
}

// RetryOperation is a high-level wrapper for common retry operations
func (rm *RetryManager) RetryOperation(ctx context.Context, operationName string, fn RetryFunc) error {
	result := rm.Retry(ctx, fn)
	
	if !result.Success {
		// Create a detailed error with retry information
		details := map[string]string{
			"operation":     operationName,
			"attempts":      fmt.Sprintf("%d", result.Attempts),
			"total_delay":   result.TotalDelay.String(),
			"duration":      result.EndTime.Sub(result.StartTime).String(),
		}
		
		if result.LastError != nil {
			details["last_error"] = result.LastError.Error()
		}
		
		return rm.errorManager.CreateError(
			ErrorTypeSystem,
			ErrCodeInternalError,
			fmt.Sprintf("Operation '%s' failed after %d attempts", operationName, result.Attempts),
			WithDetails(details),
			WithSeverity(ErrorSeverityHigh),
			WithSuggestions([]string{
				"The operation was retried automatically but still failed",
				"Check the underlying service availability",
				"Try again later or contact support if the issue persists",
			}),
		)
	}
	
	return nil
}

// Helper functions

// isRetryableError checks if an error should be retried
func (rm *RetryManager) isRetryableError(err error, policy RetryPolicy) bool {
	if err == nil {
		return false
	}

	// Check if it's an APIError and has retry information
	if apiErr, ok := err.(*APIError); ok {
		if !apiErr.Recoverable {
			return false
		}
		
		// Check if error code is in the retryable list
		for _, code := range policy.RetryableErrors {
			if apiErr.Code == code {
				return true
			}
		}
		return false
	}

	// For non-APIError types, check against common retryable errors
	errorStr := err.Error()
	retryablePatterns := []string{
		"connection refused",
		"timeout",
		"temporary failure",
		"service unavailable",
		"network error",
		"database connection",
		"context deadline exceeded",
	}

	for _, pattern := range retryablePatterns {
		if containsSubstring(errorStr, pattern) {
			return true
		}
	}

	return false
}

// calculateDelay calculates the actual delay with jitter
func (rm *RetryManager) calculateDelay(baseDelay time.Duration, policy RetryPolicy) time.Duration {
	if !policy.JitterEnabled {
		return baseDelay
	}

	// Add jitter to prevent thundering herd
	jitter := time.Duration(rand.Float64() * float64(baseDelay) * 0.1) // 10% jitter
	return baseDelay + jitter
}

// containsSubstring checks if a string contains a substring (case-insensitive)
func containsSubstring(s, substr string) bool {
	return strings.Contains(strings.ToLower(s), strings.ToLower(substr))
}

// CircuitBreaker implements the circuit breaker pattern
type CircuitBreaker struct {
	maxFailures     int
	timeout         time.Duration
	failureCount    int
	lastFailureTime time.Time
	state           CircuitState
	errorManager    *ErrorManager
}

// CircuitState represents the state of a circuit breaker
type CircuitState int

const (
	CircuitClosed CircuitState = iota
	CircuitOpen
	CircuitHalfOpen
)

// NewCircuitBreaker creates a new circuit breaker
func NewCircuitBreaker(maxFailures int, timeout time.Duration, errorManager *ErrorManager) *CircuitBreaker {
	return &CircuitBreaker{
		maxFailures:  maxFailures,
		timeout:      timeout,
		state:        CircuitClosed,
		errorManager: errorManager,
	}
}

// Execute executes a function with circuit breaker protection
func (cb *CircuitBreaker) Execute(fn func() error) error {
	if cb.state == CircuitOpen {
		if time.Since(cb.lastFailureTime) > cb.timeout {
			cb.state = CircuitHalfOpen
		} else {
			return cb.errorManager.CreateError(
				ErrorTypeSystem,
				ErrCodeInternalError,
				"Circuit breaker is open",
				WithDetails(map[string]string{
					"state":         "open",
					"failure_count": fmt.Sprintf("%d", cb.failureCount),
					"timeout":       cb.timeout.String(),
				}),
				WithSeverity(ErrorSeverityMedium),
				WithRecoverable(true),
				WithRetryAfter(cb.timeout),
			)
		}
	}

	err := fn()
	if err != nil {
		cb.onFailure()
		return err
	}

	cb.onSuccess()
	return nil
}

// onFailure handles circuit breaker failure
func (cb *CircuitBreaker) onFailure() {
	cb.failureCount++
	cb.lastFailureTime = time.Now()

	if cb.failureCount >= cb.maxFailures {
		cb.state = CircuitOpen
	}
}

// onSuccess handles circuit breaker success
func (cb *CircuitBreaker) onSuccess() {
	cb.failureCount = 0
	cb.state = CircuitClosed
}

// GetState returns the current circuit breaker state
func (cb *CircuitBreaker) GetState() CircuitState {
	return cb.state
}

// IsOpen returns true if the circuit breaker is open
func (cb *CircuitBreaker) IsOpen() bool {
	return cb.state == CircuitOpen
}

// BulkheadManager manages resource isolation using the bulkhead pattern
type BulkheadManager struct {
	semaphores map[string]chan struct{}
	limits     map[string]int
}

// NewBulkheadManager creates a new bulkhead manager
func NewBulkheadManager() *BulkheadManager {
	return &BulkheadManager{
		semaphores: make(map[string]chan struct{}),
		limits:     make(map[string]int),
	}
}

// SetLimit sets the concurrency limit for a resource
func (bm *BulkheadManager) SetLimit(resource string, limit int) {
	bm.limits[resource] = limit
	bm.semaphores[resource] = make(chan struct{}, limit)
}

// Execute executes a function with bulkhead protection
func (bm *BulkheadManager) Execute(ctx context.Context, resource string, fn func() error) error {
	semaphore, exists := bm.semaphores[resource]
	if !exists {
		return fn() // No limit set, execute directly
	}

	select {
	case semaphore <- struct{}{}: // Acquire
		defer func() { <-semaphore }() // Release
		return fn()
	case <-ctx.Done():
		return ctx.Err()
	}
}