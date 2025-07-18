import React from 'react';
import { APIError } from '../types';
import './ErrorDisplay.scss';

interface ErrorDisplayProps {
  error: APIError | Error | string;
  onRetry?: () => void;
  onDismiss?: () => void;
  showDetails?: boolean;
  className?: string;
}

export const ErrorDisplay: React.FC<ErrorDisplayProps> = ({
  error,
  onRetry,
  onDismiss,
  showDetails = false,
  className = ''
}) => {
  const getErrorInfo = () => {
    if (typeof error === 'string') {
      return {
        message: error,
        type: 'error',
        severity: 'medium' as const,
        suggestions: [],
        recoverable: false,
        code: 'UNKNOWN_ERROR'
      };
    }

    if (error instanceof Error) {
      return {
        message: error.message,
        type: 'error',
        severity: 'medium' as const,
        suggestions: [],
        recoverable: false,
        code: 'JAVASCRIPT_ERROR'
      };
    }

    // APIError
    return {
      message: error.message,
      type: error.type || 'error',
      severity: error.severity || 'medium',
      suggestions: error.suggestions || [],
      recoverable: error.recoverable || false,
      code: error.code || 'UNKNOWN_ERROR',
      details: error.details,
      correlationId: error.correlationId
    };
  };

  const errorInfo = getErrorInfo();

  const getSeverityIcon = () => {
    switch (errorInfo.severity) {
      case 'low':
        return (
          <svg viewBox="0 0 24 24" fill="currentColor">
            <path d="M12 2C6.48 2 2 6.48 2 12s4.48 10 10 10 10-4.48 10-10S17.52 2 12 2zm-2 15l-5-5 1.41-1.41L10 14.17l7.59-7.59L19 8l-9 9z"/>
          </svg>
        );
      case 'medium':
        return (
          <svg viewBox="0 0 24 24" fill="currentColor">
            <path d="M1 21h22L12 2 1 21zm12-3h-2v-2h2v2zm0-4h-2v-4h2v4z"/>
          </svg>
        );
      case 'high':
        return (
          <svg viewBox="0 0 24 24" fill="currentColor">
            <path d="M12 2C6.48 2 2 6.48 2 12s4.48 10 10 10 10-4.48 10-10S17.52 2 12 2zm-2 15l-5-5 1.41-1.41L10 14.17l7.59-7.59L19 8l-9 9z"/>
          </svg>
        );
      case 'critical':
        return (
          <svg viewBox="0 0 24 24" fill="currentColor">
            <path d="M12 2C6.48 2 2 6.48 2 12s4.48 10 10 10 10-4.48 10-10S17.52 2 12 2zm-2 15l-5-5 1.41-1.41L10 14.17l7.59-7.59L19 8l-9 9z"/>
          </svg>
        );
      default:
        return null;
    }
  };

  const getCategoryTitle = () => {
    switch (errorInfo.type) {
      case 'validation':
        return 'Input Error';
      case 'network':
        return 'Network Error';
      case 'filesystem':
        return 'File System Error';
      case 'database':
        return 'Database Error';
      case 'generation':
        return 'Generation Error';
      case 'system':
        return 'System Error';
      case 'permission':
        return 'Permission Error';
      case 'configuration':
        return 'Configuration Error';
      case 'authentication':
        return 'Authentication Error';
      default:
        return 'Error';
    }
  };

  const getDefaultSuggestions = () => {
    switch (errorInfo.type) {
      case 'validation':
        return [
          'Check your input for errors',
          'Ensure all required fields are filled',
          'Verify the format of your data'
        ];
      case 'network':
        return [
          'Check your internet connection',
          'Try again in a few moments',
          'Verify the service is available'
        ];
      case 'filesystem':
        return [
          'Check file permissions',
          'Ensure the file exists',
          'Verify the file path is correct'
        ];
      case 'generation':
        return [
          'Check the OpenAPI specification',
          'Verify project configuration',
          'Try regenerating the project'
        ];
      default:
        return [
          'Try refreshing the page',
          'Contact support if the issue persists'
        ];
    }
  };

  const suggestions = errorInfo.suggestions.length > 0 ? errorInfo.suggestions : getDefaultSuggestions();

  return (
    <div className={`error-display error-display--${errorInfo.severity} ${className}`}>
      <div className="error-display__header">
        <div className="error-display__icon">
          {getSeverityIcon()}
        </div>
        <div className="error-display__title-area">
          <h3 className="error-display__title">{getCategoryTitle()}</h3>
          {errorInfo.correlationId && (
            <span className="error-display__correlation-id">
              ID: {errorInfo.correlationId}
            </span>
          )}
        </div>
        {onDismiss && (
          <button
            onClick={onDismiss}
            className="error-display__dismiss"
            aria-label="Dismiss error"
          >
            <svg viewBox="0 0 24 24" fill="currentColor">
              <path d="M19 6.41L17.59 5 12 10.59 6.41 5 5 6.41 10.59 12 5 17.59 6.41 19 12 13.41 17.59 19 19 17.59 13.41 12z"/>
            </svg>
          </button>
        )}
      </div>

      <div className="error-display__content">
        <p className="error-display__message">{errorInfo.message}</p>

        {suggestions.length > 0 && (
          <div className="error-display__suggestions">
            <h4>What you can do:</h4>
            <ul>
              {suggestions.map((suggestion, index) => (
                <li key={index}>{suggestion}</li>
              ))}
            </ul>
          </div>
        )}

        {showDetails && errorInfo.details && Object.keys(errorInfo.details).length > 0 && (
          <details className="error-display__details">
            <summary>Technical Details</summary>
            <div className="error-display__details-content">
              {Object.entries(errorInfo.details).map(([key, value]) => (
                <div key={key} className="error-display__detail-item">
                  <strong>{key}:</strong> {value}
                </div>
              ))}
            </div>
          </details>
        )}
      </div>

      {(onRetry || errorInfo.recoverable) && (
        <div className="error-display__actions">
          {onRetry && (
            <button
              onClick={onRetry}
              className="error-display__retry-button"
            >
              <svg viewBox="0 0 24 24" fill="currentColor">
                <path d="M17.65 6.35C16.2 4.9 14.21 4 12 4c-4.42 0-7.99 3.58-7.99 8s3.57 8 7.99 8c3.73 0 6.84-2.55 7.73-6h-2.08c-.82 2.33-3.04 4-5.65 4-3.31 0-6-2.69-6-6s2.69-6 6-6c1.66 0 3.14.69 4.22 1.78L13 11h7V4l-2.35 2.35z"/>
              </svg>
              Try Again
            </button>
          )}
        </div>
      )}
    </div>
  );
};

// Inline error display for form fields
interface InlineErrorProps {
  error: string | APIError | null;
  className?: string;
}

export const InlineError: React.FC<InlineErrorProps> = ({ error, className = '' }) => {
  if (!error) return null;

  const message = typeof error === 'string' ? error : error.message;

  return (
    <div className={`inline-error ${className}`}>
      <svg viewBox="0 0 24 24" fill="currentColor">
        <path d="M12 2C6.48 2 2 6.48 2 12s4.48 10 10 10 10-4.48 10-10S17.52 2 12 2zm-2 15l-5-5 1.41-1.41L10 14.17l7.59-7.59L19 8l-9 9z"/>
      </svg>
      <span>{message}</span>
    </div>
  );
};

// Error list component for multiple errors
interface ErrorListProps {
  errors: (APIError | Error | string)[];
  onRetry?: (index: number) => void;
  onDismiss?: (index: number) => void;
  onDismissAll?: () => void;
  className?: string;
}

export const ErrorList: React.FC<ErrorListProps> = ({
  errors,
  onRetry,
  onDismiss,
  onDismissAll,
  className = ''
}) => {
  if (errors.length === 0) return null;

  return (
    <div className={`error-list ${className}`}>
      <div className="error-list__header">
        <h3>Errors ({errors.length})</h3>
        {onDismissAll && (
          <button
            onClick={onDismissAll}
            className="error-list__dismiss-all"
          >
            Dismiss All
          </button>
        )}
      </div>
      <div className="error-list__items">
        {errors.map((error, index) => (
          <ErrorDisplay
            key={index}
            error={error}
            onRetry={onRetry ? () => onRetry(index) : undefined}
            onDismiss={onDismiss ? () => onDismiss(index) : undefined}
            className="error-list__item"
          />
        ))}
      </div>
    </div>
  );
};

export default ErrorDisplay;