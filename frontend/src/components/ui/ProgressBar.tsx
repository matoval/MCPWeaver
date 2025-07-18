import React from 'react';
import './ProgressBar.scss';

export interface ProgressBarProps {
  progress: number; // 0-1
  className?: string;
  showPercentage?: boolean;
  animated?: boolean;
  variant?: 'primary' | 'success' | 'warning' | 'danger';
  size?: 'sm' | 'md' | 'lg';
  label?: string;
  style?: React.CSSProperties;
}

export const ProgressBar: React.FC<ProgressBarProps> = ({
  progress,
  className = '',
  showPercentage = true,
  animated = true,
  variant = 'primary',
  size = 'md',
  label,
  style
}) => {
  const percentage = Math.round(progress * 100);
  const clampedProgress = Math.min(Math.max(progress, 0), 1);

  return (
    <div 
      className={`progress-bar ${className} progress-bar--${variant} progress-bar--${size}`}
      style={style}
      role="progressbar"
      aria-valuenow={percentage}
      aria-valuemin={0}
      aria-valuemax={100}
      aria-label={label || `Progress: ${percentage}%`}
    >
      {label && (
        <div className="progress-bar__label">
          {label}
        </div>
      )}
      
      <div className="progress-bar__track">
        <div 
          className={`progress-bar__fill ${animated ? 'progress-bar__fill--animated' : ''}`}
          style={{ width: `${clampedProgress * 100}%` }}
        />
      </div>
      
      {showPercentage && (
        <div className="progress-bar__percentage">
          {percentage}%
        </div>
      )}
    </div>
  );
};

export default ProgressBar;