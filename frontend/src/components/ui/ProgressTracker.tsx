import React, { useState, useEffect } from 'react';
import { useProgress } from '../../hooks/useProgress';
import { ProgressBar } from './ProgressBar';
import { StepIndicator, Step } from './StepIndicator';
import { ProgressMetrics, GenerationStatus } from '../../types';
import './ProgressTracker.scss';

export interface ProgressTrackerProps {
  projectId?: string;
  className?: string;
  compact?: boolean;
  showMetrics?: boolean;
  showHistory?: boolean;
  onCancel?: () => void;
  onComplete?: () => void;
}

export const ProgressTracker: React.FC<ProgressTrackerProps> = ({
  projectId,
  className = '',
  compact = false,
  showMetrics = true,
  showHistory = false,
  onCancel,
  onComplete
}) => {
  const {
    currentJob,
    progress,
    currentStep,
    isGenerating,
    metrics,
    history,
    startGeneration,
    cancelGeneration,
    clearHistory,
    error,
    clearError
  } = useProgress();

  const [showCancelDialog, setShowCancelDialog] = useState(false);

  // Define generation steps
  const steps: Step[] = [
    {
      id: 'parsing',
      title: 'Parsing',
      description: 'Analyzing OpenAPI specification',
      status: 'pending'
    },
    {
      id: 'mapping',
      title: 'Mapping',
      description: 'Converting operations to MCP tools',
      status: 'pending'
    },
    {
      id: 'generating',
      title: 'Generating',
      description: 'Creating server code',
      status: 'pending'
    },
    {
      id: 'validating',
      title: 'Validating',
      description: 'Checking generated code',
      status: 'pending'
    },
    {
      id: 'completed',
      title: 'Complete',
      description: 'Generation finished successfully',
      status: 'pending'
    }
  ];

  // Map generation status to step ID
  const getStepIdFromStatus = (status: GenerationStatus): string => {
    switch (status) {
      case 'parsing': return 'parsing';
      case 'mapping': return 'mapping';
      case 'generating': return 'generating';
      case 'validating': return 'validating';
      case 'completed': return 'completed';
      case 'failed': return 'parsing'; // Show error on first step
      case 'cancelled': return 'parsing'; // Show cancelled on first step
      default: return 'parsing';
    }
  };

  // Get current step for indicator
  const getCurrentStep = (): string => {
    if (!currentJob) return 'parsing';
    return getStepIdFromStatus(currentJob.status);
  };

  // Update steps with error status if needed
  const getStepsWithStatus = (): Step[] => {
    if (!currentJob) return steps;
    
    const currentStepId = getCurrentStep();
    const currentStepIndex = steps.findIndex(s => s.id === currentStepId);
    
    return steps.map((step, index) => ({
      ...step,
      status: currentJob.status === 'failed' && index === currentStepIndex ? 'error' : step.status
    }));
  };

  // Handle start generation
  const handleStart = async () => {
    if (!projectId) return;
    
    try {
      clearError();
      await startGeneration(projectId);
    } catch (error) {
      console.error('Failed to start generation:', error);
    }
  };

  // Handle cancel generation
  const handleCancel = async () => {
    if (!currentJob) return;
    
    try {
      await cancelGeneration(currentJob.id);
      setShowCancelDialog(false);
      onCancel?.();
    } catch (error) {
      console.error('Failed to cancel generation:', error);
    }
  };

  // Format duration
  const formatDuration = (ms: number): string => {
    const seconds = Math.floor(ms / 1000);
    const minutes = Math.floor(seconds / 60);
    const hours = Math.floor(minutes / 60);
    
    if (hours > 0) {
      return `${hours}h ${minutes % 60}m ${seconds % 60}s`;
    } else if (minutes > 0) {
      return `${minutes}m ${seconds % 60}s`;
    } else {
      return `${seconds}s`;
    }
  };

  // Format rate
  const formatRate = (rate: number): string => {
    if (rate < 1) {
      return `${(rate * 100).toFixed(1)}%/s`;
    }
    return `${rate.toFixed(2)}/s`;
  };

  // Format file size
  const formatFileSize = (bytes: number): string => {
    if (bytes < 1024) return `${bytes} B`;
    if (bytes < 1024 * 1024) return `${(bytes / 1024).toFixed(1)} KB`;
    return `${(bytes / (1024 * 1024)).toFixed(1)} MB`;
  };

  // Notify completion
  useEffect(() => {
    if (currentJob?.status === 'completed') {
      onComplete?.();
    }
  }, [currentJob?.status, onComplete]);

  if (!isGenerating && !currentJob && !error) {
    return (
      <div className={`progress-tracker ${className} progress-tracker--idle`}>
        <div className="progress-tracker__idle-state">
          <div className="progress-tracker__idle-icon">
            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M13 10V3L4 14h7v7l9-11h-7z" />
            </svg>
          </div>
          <div className="progress-tracker__idle-text">
            <h3>Ready to Generate</h3>
            <p>Click generate to start creating your MCP server</p>
          </div>
          {projectId && (
            <button 
              className="progress-tracker__start-button"
              onClick={handleStart}
              disabled={!projectId}
            >
              Start Generation
            </button>
          )}
        </div>
      </div>
    );
  }

  return (
    <div className={`progress-tracker ${className} ${compact ? 'progress-tracker--compact' : ''}`}>
      {/* Error State */}
      {error && (
        <div className="progress-tracker__error">
          <div className="progress-tracker__error-icon">
            <svg viewBox="0 0 20 20" fill="currentColor">
              <path fillRule="evenodd" d="M18 10a8 8 0 11-16 0 8 8 0 0116 0zm-7 4a1 1 0 11-2 0 1 1 0 012 0zm-1-9a1 1 0 00-1 1v4a1 1 0 102 0V6a1 1 0 00-1-1z" clipRule="evenodd" />
            </svg>
          </div>
          <div className="progress-tracker__error-content">
            <h3>Generation Failed</h3>
            <p>{error}</p>
            <button 
              className="progress-tracker__error-dismiss"
              onClick={clearError}
            >
              Dismiss
            </button>
          </div>
        </div>
      )}

      {/* Progress Header */}
      <div className="progress-tracker__header">
        <div className="progress-tracker__status">
          <h3 className="progress-tracker__title">
            {currentJob?.status === 'completed' ? 'Generation Complete' : 'Generating MCP Server'}
          </h3>
          <p className="progress-tracker__step">{currentStep}</p>
        </div>
        
        {isGenerating && currentJob && (
          <div className="progress-tracker__actions">
            <button 
              className="progress-tracker__cancel-button"
              onClick={() => setShowCancelDialog(true)}
              title="Cancel generation"
            >
              <svg viewBox="0 0 20 20" fill="currentColor">
                <path fillRule="evenodd" d="M4.293 4.293a1 1 0 011.414 0L10 8.586l4.293-4.293a1 1 0 111.414 1.414L11.414 10l4.293 4.293a1 1 0 01-1.414 1.414L10 11.414l-4.293 4.293a1 1 0 01-1.414-1.414L8.586 10 4.293 5.707a1 1 0 010-1.414z" clipRule="evenodd" />
              </svg>
            </button>
          </div>
        )}
      </div>

      {/* Progress Bar */}
      <ProgressBar
        progress={progress}
        variant={currentJob?.status === 'failed' ? 'danger' : currentJob?.status === 'completed' ? 'success' : 'primary'}
        animated={isGenerating}
        className="progress-tracker__progress-bar"
      />

      {/* Step Indicator */}
      {!compact && (
        <StepIndicator
          steps={getStepsWithStatus()}
          currentStep={getCurrentStep()}
          className="progress-tracker__steps"
          orientation="horizontal"
          showConnectors={true}
        />
      )}

      {/* Metrics */}
      {showMetrics && metrics && !compact && (
        <div className="progress-tracker__metrics">
          <div className="progress-tracker__metrics-grid">
            <div className="progress-tracker__metric">
              <div className="progress-tracker__metric-value">
                {formatDuration(metrics.elapsedTime)}
              </div>
              <div className="progress-tracker__metric-label">Elapsed Time</div>
            </div>
            
            <div className="progress-tracker__metric">
              <div className="progress-tracker__metric-value">
                {formatDuration(metrics.estimatedTimeRemaining * 1000)}
              </div>
              <div className="progress-tracker__metric-label">Time Remaining</div>
            </div>
            
            <div className="progress-tracker__metric">
              <div className="progress-tracker__metric-value">
                {formatRate(metrics.processingRate)}
              </div>
              <div className="progress-tracker__metric-label">Processing Rate</div>
            </div>
            
            <div className="progress-tracker__metric">
              <div className="progress-tracker__metric-value">
                {metrics.filesGenerated}/{metrics.totalFiles}
              </div>
              <div className="progress-tracker__metric-label">Files Generated</div>
            </div>
            
            {metrics.errorCount > 0 && (
              <div className="progress-tracker__metric progress-tracker__metric--error">
                <div className="progress-tracker__metric-value">
                  {metrics.errorCount}
                </div>
                <div className="progress-tracker__metric-label">Errors</div>
              </div>
            )}
            
            {metrics.warningCount > 0 && (
              <div className="progress-tracker__metric progress-tracker__metric--warning">
                <div className="progress-tracker__metric-value">
                  {metrics.warningCount}
                </div>
                <div className="progress-tracker__metric-label">Warnings</div>
              </div>
            )}
          </div>
        </div>
      )}

      {/* History */}
      {showHistory && history.length > 0 && !compact && (
        <div className="progress-tracker__history">
          <div className="progress-tracker__history-header">
            <h4>Recent Generations</h4>
            <button 
              className="progress-tracker__history-clear"
              onClick={clearHistory}
            >
              Clear History
            </button>
          </div>
          <div className="progress-tracker__history-list">
            {history.slice(0, 5).map((entry) => (
              <div key={entry.jobId} className="progress-tracker__history-item">
                <div className="progress-tracker__history-item-status">
                  {entry.success ? (
                    <svg viewBox="0 0 20 20" fill="currentColor" className="progress-tracker__history-icon progress-tracker__history-icon--success">
                      <path fillRule="evenodd" d="M16.707 5.293a1 1 0 010 1.414l-8 8a1 1 0 01-1.414 0l-4-4a1 1 0 011.414-1.414L8 12.586l7.293-7.293a1 1 0 011.414 0z" clipRule="evenodd" />
                    </svg>
                  ) : (
                    <svg viewBox="0 0 20 20" fill="currentColor" className="progress-tracker__history-icon progress-tracker__history-icon--error">
                      <path fillRule="evenodd" d="M4.293 4.293a1 1 0 011.414 0L10 8.586l4.293-4.293a1 1 0 111.414 1.414L11.414 10l4.293 4.293a1 1 0 01-1.414 1.414L10 11.414l-4.293 4.293a1 1 0 01-1.414-1.414L8.586 10 4.293 5.707a1 1 0 010-1.414z" clipRule="evenodd" />
                    </svg>
                  )}
                </div>
                <div className="progress-tracker__history-item-content">
                  <div className="progress-tracker__history-item-title">
                    {entry.projectName}
                  </div>
                  <div className="progress-tracker__history-item-meta">
                    {new Date(entry.startTime).toLocaleDateString()} • {formatDuration(entry.duration || 0)}
                  </div>
                  {entry.errorMessage && (
                    <div className="progress-tracker__history-item-error">
                      {entry.errorMessage}
                    </div>
                  )}
                </div>
              </div>
            ))}
          </div>
        </div>
      )}

      {/* Cancel Confirmation Dialog */}
      {showCancelDialog && (
        <div className="progress-tracker__dialog-overlay">
          <div className="progress-tracker__dialog">
            <h3>Cancel Generation?</h3>
            <p>Are you sure you want to cancel the current generation? This action cannot be undone.</p>
            <div className="progress-tracker__dialog-actions">
              <button 
                className="progress-tracker__dialog-button progress-tracker__dialog-button--secondary"
                onClick={() => setShowCancelDialog(false)}
              >
                Keep Generating
              </button>
              <button 
                className="progress-tracker__dialog-button progress-tracker__dialog-button--danger"
                onClick={handleCancel}
              >
                Cancel Generation
              </button>
            </div>
          </div>
        </div>
      )}
    </div>
  );
};

export default ProgressTracker;