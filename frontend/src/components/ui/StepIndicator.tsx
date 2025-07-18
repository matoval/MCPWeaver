import React from 'react';
import './StepIndicator.scss';

export interface Step {
  id: string;
  title: string;
  description?: string;
  status: 'pending' | 'active' | 'completed' | 'error';
}

export interface StepIndicatorProps {
  steps: Step[];
  currentStep: string;
  className?: string;
  orientation?: 'horizontal' | 'vertical';
  showConnectors?: boolean;
  compact?: boolean;
}

export const StepIndicator: React.FC<StepIndicatorProps> = ({
  steps,
  currentStep,
  className = '',
  orientation = 'horizontal',
  showConnectors = true,
  compact = false
}) => {
  const getStepStatus = (step: Step): 'pending' | 'active' | 'completed' | 'error' => {
    const currentIndex = steps.findIndex(s => s.id === currentStep);
    const stepIndex = steps.findIndex(s => s.id === step.id);
    
    if (step.status === 'error') return 'error';
    if (stepIndex < currentIndex) return 'completed';
    if (stepIndex === currentIndex) return 'active';
    return 'pending';
  };

  const getStepIcon = (status: string): React.ReactNode => {
    switch (status) {
      case 'completed':
        return (
          <svg className="step-indicator__icon" viewBox="0 0 20 20" fill="currentColor">
            <path fillRule="evenodd" d="M16.707 5.293a1 1 0 010 1.414l-8 8a1 1 0 01-1.414 0l-4-4a1 1 0 011.414-1.414L8 12.586l7.293-7.293a1 1 0 011.414 0z" clipRule="evenodd" />
          </svg>
        );
      case 'error':
        return (
          <svg className="step-indicator__icon" viewBox="0 0 20 20" fill="currentColor">
            <path fillRule="evenodd" d="M4.293 4.293a1 1 0 011.414 0L10 8.586l4.293-4.293a1 1 0 111.414 1.414L11.414 10l4.293 4.293a1 1 0 01-1.414 1.414L10 11.414l-4.293 4.293a1 1 0 01-1.414-1.414L8.586 10 4.293 5.707a1 1 0 010-1.414z" clipRule="evenodd" />
          </svg>
        );
      case 'active':
        return (
          <div className="step-indicator__spinner">
            <svg className="step-indicator__spinner-icon" viewBox="0 0 24 24">
              <circle cx="12" cy="12" r="10" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeDasharray="31.416" strokeDashoffset="31.416">
                <animate attributeName="stroke-dasharray" dur="2s" values="0 31.416;15.708 15.708;0 31.416" repeatCount="indefinite"/>
                <animate attributeName="stroke-dashoffset" dur="2s" values="0;-15.708;-31.416" repeatCount="indefinite"/>
              </circle>
            </svg>
          </div>
        );
      default:
        return <div className="step-indicator__number">{steps.findIndex(s => s.id === step.id) + 1}</div>;
    }
  };

  return (
    <div 
      className={`step-indicator ${className} step-indicator--${orientation} ${compact ? 'step-indicator--compact' : ''}`}
      role="progressbar"
      aria-label="Generation progress steps"
    >
      {steps.map((step, index) => {
        const status = getStepStatus(step);
        const isLast = index === steps.length - 1;
        
        return (
          <div key={step.id} className="step-indicator__step-wrapper">
            <div 
              className={`step-indicator__step step-indicator__step--${status}`}
              role="step"
              aria-current={status === 'active' ? 'step' : undefined}
              aria-label={`Step ${index + 1}: ${step.title}`}
            >
              <div className="step-indicator__step-marker">
                {getStepIcon(status)}
              </div>
              
              {!compact && (
                <div className="step-indicator__step-content">
                  <div className="step-indicator__step-title">
                    {step.title}
                  </div>
                  {step.description && (
                    <div className="step-indicator__step-description">
                      {step.description}
                    </div>
                  )}
                </div>
              )}
            </div>
            
            {showConnectors && !isLast && (
              <div className={`step-indicator__connector step-indicator__connector--${status === 'completed' ? 'completed' : 'pending'}`} />
            )}
          </div>
        );
      })}
    </div>
  );
};

export default StepIndicator;