import React, { Component, ErrorInfo, ReactNode } from 'react';
import { wails } from '../services/wails';
import './ErrorBoundary.scss';

interface ErrorBoundaryProps {
  children: ReactNode;
  fallback?: ReactNode;
  onError?: (error: Error, errorInfo: ErrorInfo) => void;
  context?: string;
}

interface ErrorBoundaryState {
  hasError: boolean;
  error: Error | null;
  errorInfo: ErrorInfo | null;
  errorId: string | null;
}

export class ErrorBoundary extends Component<ErrorBoundaryProps, ErrorBoundaryState> {
  private retryCount = 0;
  private maxRetries = 3;

  constructor(props: ErrorBoundaryProps) {
    super(props);
    this.state = {
      hasError: false,
      error: null,
      errorInfo: null,
      errorId: null,
    };
  }

  static getDerivedStateFromError(error: Error): Partial<ErrorBoundaryState> {
    return {
      hasError: true,
      error,
      errorId: `error-${Date.now()}-${Math.random().toString(36).substr(2, 9)}`,
    };
  }

  componentDidCatch(error: Error, errorInfo: ErrorInfo) {
    this.setState({ errorInfo });

    // Log error to console for development
    console.error('ErrorBoundary caught an error:', error, errorInfo);

    // Report error to backend if available
    this.reportError(error, errorInfo);

    // Call custom error handler if provided
    if (this.props.onError) {
      this.props.onError(error, errorInfo);
    }
  }

  private reportError = async (error: Error, errorInfo: ErrorInfo) => {
    try {
      const errorReport = {
        message: error.message,
        stack: error.stack,
        componentStack: errorInfo.componentStack,
        context: this.props.context || 'unknown',
        timestamp: new Date().toISOString(),
        userAgent: navigator.userAgent,
        url: window.location.href,
        errorId: this.state.errorId,
      };

      // Send error report to backend
      await wails.reportError(errorReport);
    } catch (reportError) {
      console.error('Failed to report error:', reportError);
    }
  };

  private handleRetry = () => {
    if (this.retryCount < this.maxRetries) {
      this.retryCount++;
      this.setState({
        hasError: false,
        error: null,
        errorInfo: null,
        errorId: null,
      });
    }
  };

  private handleReload = () => {
    window.location.reload();
  };

  private handleReportIssue = () => {
    const errorDetails = {
      message: this.state.error?.message || 'Unknown error',
      stack: this.state.error?.stack || 'No stack trace',
      componentStack: this.state.errorInfo?.componentStack || 'No component stack',
      context: this.props.context || 'unknown',
      errorId: this.state.errorId,
    };

    const githubUrl = `https://github.com/matoval/MCPWeaver/issues/new?title=${encodeURIComponent(
      `Error: ${errorDetails.message}`
    )}&body=${encodeURIComponent(
      `## Error Details\n\n**Error ID**: ${errorDetails.errorId}\n**Context**: ${errorDetails.context}\n**Message**: ${errorDetails.message}\n\n**Stack Trace**:\n\`\`\`\n${errorDetails.stack}\n\`\`\`\n\n**Component Stack**:\n\`\`\`\n${errorDetails.componentStack}\n\`\`\`\n\n**Environment**:\n- Browser: ${navigator.userAgent}\n- URL: ${window.location.href}\n- Timestamp: ${new Date().toISOString()}\n\n## Steps to Reproduce\n\n1. [Please describe the steps that led to this error]\n2. \n3. \n\n## Expected Behavior\n\n[Please describe what you expected to happen]\n\n## Additional Context\n\n[Please add any additional context about the error here]`
    )}`;

    window.open(githubUrl, '_blank');
  };

  render() {
    if (this.state.hasError) {
      // Custom fallback UI
      if (this.props.fallback) {
        return this.props.fallback;
      }

      // Default error UI
      return (
        <div className="error-boundary">
          <div className="error-boundary__container">
            <div className="error-boundary__icon">
              <svg viewBox="0 0 24 24" fill="currentColor">
                <path d="M12 2C6.48 2 2 6.48 2 12s4.48 10 10 10 10-4.48 10-10S17.52 2 12 2zm-2 15l-5-5 1.41-1.41L10 14.17l7.59-7.59L19 8l-9 9z" />
              </svg>
            </div>
            
            <div className="error-boundary__content">
              <h1 className="error-boundary__title">Something went wrong</h1>
              <p className="error-boundary__message">
                We're sorry, but something unexpected happened. This error has been reported to help us improve the application.
              </p>
              
              <div className="error-boundary__details">
                <details className="error-boundary__error-details">
                  <summary>Error Details</summary>
                  <div className="error-boundary__error-info">
                    <p><strong>Error ID:</strong> {this.state.errorId}</p>
                    <p><strong>Context:</strong> {this.props.context || 'Application'}</p>
                    <p><strong>Message:</strong> {this.state.error?.message}</p>
                    {this.state.error?.stack && (
                      <div className="error-boundary__stack-trace">
                        <strong>Stack Trace:</strong>
                        <pre>{this.state.error.stack}</pre>
                      </div>
                    )}
                  </div>
                </details>
              </div>
              
              <div className="error-boundary__actions">
                {this.retryCount < this.maxRetries && (
                  <button
                    onClick={this.handleRetry}
                    className="error-boundary__button error-boundary__button--primary"
                  >
                    Try Again ({this.maxRetries - this.retryCount} attempts left)
                  </button>
                )}
                
                <button
                  onClick={this.handleReload}
                  className="error-boundary__button error-boundary__button--secondary"
                >
                  Reload Page
                </button>
                
                <button
                  onClick={this.handleReportIssue}
                  className="error-boundary__button error-boundary__button--outline"
                >
                  Report Issue
                </button>
              </div>
              
              <div className="error-boundary__help">
                <p>
                  If this error persists, please try the following:
                </p>
                <ul>
                  <li>Refresh the page</li>
                  <li>Clear your browser cache</li>
                  <li>Check for application updates</li>
                  <li>Report the issue using the button above</li>
                </ul>
              </div>
            </div>
          </div>
        </div>
      );
    }

    return this.props.children;
  }
}

// Higher-order component for wrapping components with error boundary
export function withErrorBoundary<P extends object>(
  Component: React.ComponentType<P>,
  errorBoundaryProps?: Partial<ErrorBoundaryProps>
) {
  const WrappedComponent = (props: P) => (
    <ErrorBoundary {...errorBoundaryProps}>
      <Component {...props} />
    </ErrorBoundary>
  );

  WrappedComponent.displayName = `withErrorBoundary(${Component.displayName || Component.name})`;

  return WrappedComponent;
}

export default ErrorBoundary;