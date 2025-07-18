import { APIError, RetryPolicy, RetryResult } from '../types';

// Default retry policy for frontend operations
export const defaultRetryPolicy: RetryPolicy = {
  maxRetries: 3,
  initialDelay: 1000,
  maxDelay: 30000,
  backoffMultiplier: 2,
  jitterEnabled: true,
  retryableErrors: [
    'NETWORK_ERROR',
    'INTERNAL_ERROR',
    'TIMEOUT_ERROR',
    'RATE_LIMIT_ERROR'
  ]
};

// Types for async operations
export type AsyncOperation<T> = () => Promise<T>;
export type RetryCallback = (attempt: number, error: Error) => void;

// Retry manager for frontend operations
export class RetryManager {
  private policy: RetryPolicy;

  constructor(policy: RetryPolicy = defaultRetryPolicy) {
    this.policy = policy;
  }

  // Execute an async operation with retry logic
  async executeWithRetry<T>(
    operation: AsyncOperation<T>,
    customPolicy?: Partial<RetryPolicy>,
    onRetry?: RetryCallback
  ): Promise<T> {
    const policy = { ...this.policy, ...customPolicy };
    const result = await this.retryWithPolicy(operation, policy, onRetry);
    
    if (!result.success) {
      throw result.lastError || new Error('Operation failed after retries');
    }
    
    return result.result as T;
  }

  // Execute with retry and return detailed result
  async retryWithPolicy<T>(
    operation: AsyncOperation<T>,
    policy: RetryPolicy,
    onRetry?: RetryCallback
  ): Promise<RetryResult & { result?: T }> {
    const startTime = Date.now();
    let lastError: Error | null = null;
    let totalDelay = 0;
    let delay = policy.initialDelay;

    for (let attempt = 0; attempt <= policy.maxRetries; attempt++) {
      try {
        const result = await operation();
        return {
          success: true,
          attempts: attempt + 1,
          totalDelay,
          startTime: new Date(startTime).toISOString(),
          endTime: new Date().toISOString(),
          result
        };
      } catch (error) {
        lastError = error as Error;
        
        // Don't retry on the last attempt
        if (attempt === policy.maxRetries) {
          break;
        }

        // Check if error is retryable
        if (!this.isRetryableError(error as Error, policy)) {
          break;
        }

        // Calculate delay with exponential backoff and jitter
        const actualDelay = this.calculateDelay(delay, policy);
        totalDelay += actualDelay;

        // Call retry callback if provided
        if (onRetry) {
          onRetry(attempt + 1, lastError);
        }

        // Wait before retry
        await this.sleep(actualDelay);

        // Increase delay for next attempt
        delay = Math.min(delay * policy.backoffMultiplier, policy.maxDelay);
      }
    }

    return {
      success: false,
      attempts: policy.maxRetries + 1,
      lastError: this.wrapError(lastError),
      totalDelay,
      startTime: new Date(startTime).toISOString(),
      endTime: new Date().toISOString()
    };
  }

  // Check if an error is retryable
  private isRetryableError(error: Error, policy: RetryPolicy): boolean {
    // Check for APIError first
    if (this.isAPIError(error)) {
      const apiError = error as APIError;
      return apiError.recoverable && policy.retryableErrors.includes(apiError.code);
    }

    // Check for network errors
    const errorMessage = error.message.toLowerCase();
    const networkErrorPatterns = [
      'network error',
      'fetch error',
      'connection failed',
      'timeout',
      'aborted',
      'rate limit',
      'service unavailable',
      'gateway timeout',
      'connection refused'
    ];

    return networkErrorPatterns.some(pattern => 
      errorMessage.includes(pattern)
    );
  }

  // Calculate delay with jitter
  private calculateDelay(baseDelay: number, policy: RetryPolicy): number {
    if (!policy.jitterEnabled) {
      return baseDelay;
    }

    // Add up to 10% jitter
    const jitter = Math.random() * baseDelay * 0.1;
    return baseDelay + jitter;
  }

  // Sleep utility
  private sleep(ms: number): Promise<void> {
    return new Promise(resolve => setTimeout(resolve, ms));
  }

  // Check if error is APIError
  private isAPIError(error: any): error is APIError {
    return error && 
           typeof error === 'object' && 
           'type' in error && 
           'code' in error && 
           'message' in error;
  }

  // Wrap error in APIError format
  private wrapError(error: Error | null): APIError {
    if (!error) {
      return {
        type: 'system',
        code: 'UNKNOWN_ERROR',
        message: 'Unknown error occurred',
        severity: 'medium',
        recoverable: false,
        timestamp: new Date().toISOString()
      };
    }

    if (this.isAPIError(error)) {
      return error as APIError;
    }

    return {
      type: 'system',
      code: 'JAVASCRIPT_ERROR',
      message: error.message,
      severity: 'medium',
      recoverable: true,
      timestamp: new Date().toISOString(),
      context: {
        stackTrace: error.stack
      }
    };
  }
}

// Global retry manager instance
export const retryManager = new RetryManager();

// Utility functions

// Retry a promise with default policy
export async function retryPromise<T>(
  promise: AsyncOperation<T>,
  maxRetries: number = 3,
  initialDelay: number = 1000
): Promise<T> {
  return retryManager.executeWithRetry(promise, {
    maxRetries,
    initialDelay
  });
}

// Retry with exponential backoff
export async function retryWithBackoff<T>(
  promise: AsyncOperation<T>,
  maxRetries: number = 3,
  initialDelay: number = 1000,
  maxDelay: number = 30000
): Promise<T> {
  return retryManager.executeWithRetry(promise, {
    maxRetries,
    initialDelay,
    maxDelay,
    backoffMultiplier: 2
  });
}

// Retry API calls with specific handling
export async function retryApiCall<T>(
  apiCall: AsyncOperation<T>,
  onRetry?: (attempt: number, error: Error) => void
): Promise<T> {
  const apiRetryPolicy: Partial<RetryPolicy> = {
    maxRetries: 3,
    initialDelay: 1000,
    maxDelay: 10000,
    backoffMultiplier: 2,
    jitterEnabled: true,
    retryableErrors: [
      'NETWORK_ERROR',
      'INTERNAL_ERROR',
      'TIMEOUT_ERROR',
      'RATE_LIMIT_ERROR',
      'DATABASE_ERROR'
    ]
  };

  return retryManager.executeWithRetry(apiCall, apiRetryPolicy, onRetry);
}

// Circuit breaker for frontend
export class CircuitBreaker {
  private failures: number = 0;
  private lastFailureTime: number = 0;
  private state: 'CLOSED' | 'OPEN' | 'HALF_OPEN' = 'CLOSED';

  constructor(
    private threshold: number = 5,
    private timeout: number = 60000 // 1 minute
  ) {}

  async execute<T>(operation: AsyncOperation<T>): Promise<T> {
    if (this.state === 'OPEN') {
      if (Date.now() - this.lastFailureTime > this.timeout) {
        this.state = 'HALF_OPEN';
      } else {
        throw new Error('Circuit breaker is OPEN');
      }
    }

    try {
      const result = await operation();
      this.onSuccess();
      return result;
    } catch (error) {
      this.onFailure();
      throw error;
    }
  }

  private onSuccess(): void {
    this.failures = 0;
    this.state = 'CLOSED';
  }

  private onFailure(): void {
    this.failures++;
    this.lastFailureTime = Date.now();

    if (this.failures >= this.threshold) {
      this.state = 'OPEN';
    }
  }

  getState(): string {
    return this.state;
  }

  reset(): void {
    this.failures = 0;
    this.state = 'CLOSED';
    this.lastFailureTime = 0;
  }
}

// Timeout wrapper
export function withTimeout<T>(
  promise: Promise<T>,
  timeoutMs: number,
  timeoutMessage: string = 'Operation timed out'
): Promise<T> {
  return Promise.race([
    promise,
    new Promise<T>((_, reject) => {
      setTimeout(() => {
        reject(new Error(timeoutMessage));
      }, timeoutMs);
    })
  ]);
}

// Debounced retry - useful for user input
export function debounceRetry<T>(
  operation: AsyncOperation<T>,
  delay: number = 500
): AsyncOperation<T> {
  let timeoutId: NodeJS.Timeout | null = null;

  return () => {
    return new Promise<T>((resolve, reject) => {
      if (timeoutId) {
        clearTimeout(timeoutId);
      }

      timeoutId = setTimeout(async () => {
        try {
          const result = await operation();
          resolve(result);
        } catch (error) {
          reject(error);
        }
      }, delay);
    });
  };
}

// Batch retry - for multiple operations
export async function retryBatch<T>(
  operations: AsyncOperation<T>[],
  policy?: Partial<RetryPolicy>
): Promise<(T | Error)[]> {
  const results = await Promise.allSettled(
    operations.map(op => retryManager.executeWithRetry(op, policy))
  );

  return results.map(result => 
    result.status === 'fulfilled' ? result.value : result.reason
  );
}

export default RetryManager;