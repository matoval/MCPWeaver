import { APIError } from '../types';

// Types for graceful degradation
export interface DegradationConfig {
  feature: string;
  fallbackBehavior: 'hide' | 'disable' | 'simplify' | 'mock';
  priority: 'low' | 'medium' | 'high' | 'critical';
  retryable: boolean;
  gracePeriod?: number; // ms
  maxFailures?: number;
}

export interface FeatureState {
  available: boolean;
  degraded: boolean;
  failures: number;
  lastFailure?: Date;
  fallbackActive: boolean;
  config: DegradationConfig;
}

// Graceful degradation manager
export class GracefulDegradationManager {
  private features: Map<string, FeatureState> = new Map();
  private listeners: Map<string, Set<(state: FeatureState) => void>> = new Map();

  constructor() {
    this.setupDefaultFeatures();
  }

  // Setup default feature configurations
  private setupDefaultFeatures() {
    const defaultConfigs: DegradationConfig[] = [
      {
        feature: 'project-preview',
        fallbackBehavior: 'hide',
        priority: 'low',
        retryable: true,
        gracePeriod: 5000,
        maxFailures: 3
      },
      {
        feature: 'real-time-validation',
        fallbackBehavior: 'disable',
        priority: 'medium',
        retryable: true,
        gracePeriod: 10000,
        maxFailures: 5
      },
      {
        feature: 'auto-save',
        fallbackBehavior: 'simplify',
        priority: 'high',
        retryable: true,
        gracePeriod: 15000,
        maxFailures: 2
      },
      {
        feature: 'generation-progress',
        fallbackBehavior: 'mock',
        priority: 'medium',
        retryable: true,
        gracePeriod: 8000,
        maxFailures: 3
      },
      {
        feature: 'file-upload',
        fallbackBehavior: 'simplify',
        priority: 'high',
        retryable: true,
        gracePeriod: 12000,
        maxFailures: 2
      },
      {
        feature: 'recent-files',
        fallbackBehavior: 'hide',
        priority: 'low',
        retryable: true,
        gracePeriod: 5000,
        maxFailures: 5
      },
      {
        feature: 'syntax-highlighting',
        fallbackBehavior: 'disable',
        priority: 'low',
        retryable: false,
        maxFailures: 1
      },
      {
        feature: 'theme-switching',
        fallbackBehavior: 'disable',
        priority: 'low',
        retryable: false,
        maxFailures: 1
      }
    ];

    defaultConfigs.forEach(config => {
      this.registerFeature(config);
    });
  }

  // Register a feature for degradation management
  registerFeature(config: DegradationConfig) {
    const state: FeatureState = {
      available: true,
      degraded: false,
      failures: 0,
      fallbackActive: false,
      config
    };

    this.features.set(config.feature, state);
    this.listeners.set(config.feature, new Set());
  }

  // Report a failure for a feature
  reportFailure(feature: string, error: APIError | Error): boolean {
    const state = this.features.get(feature);
    if (!state) {
      console.warn(`Feature "${feature}" not registered for degradation`);
      return false;
    }

    state.failures++;
    state.lastFailure = new Date();

    // Check if we should degrade the feature
    const shouldDegrade = this.shouldDegrade(state, error);
    
    if (shouldDegrade && !state.degraded) {
      this.degradeFeature(feature);
    }

    return shouldDegrade;
  }

  // Check if a feature should be degraded
  private shouldDegrade(state: FeatureState, error: APIError | Error): boolean {
    const { config } = state;
    
    // Critical features should not be degraded easily
    if (config.priority === 'critical') {
      return false;
    }

    // Check max failures threshold
    if (config.maxFailures && state.failures >= config.maxFailures) {
      return true;
    }

    // Check error severity for APIError
    if (this.isAPIError(error)) {
      const apiError = error as APIError;
      
      // Don't degrade for validation errors
      if (apiError.type === 'validation') {
        return false;
      }

      // Degrade for system errors
      if (apiError.type === 'system' && apiError.severity === 'high') {
        return true;
      }

      // Degrade for network errors after multiple failures
      if (apiError.type === 'network' && state.failures >= 2) {
        return true;
      }
    }

    // Check for JavaScript errors
    if (error instanceof Error) {
      const errorMessage = error.message.toLowerCase();
      
      // Degrade for specific error types
      if (errorMessage.includes('network') || 
          errorMessage.includes('timeout') || 
          errorMessage.includes('connection')) {
        return state.failures >= 2;
      }
    }

    return false;
  }

  // Degrade a feature
  private degradeFeature(feature: string) {
    const state = this.features.get(feature);
    if (!state) return;

    state.degraded = true;
    state.available = false;
    state.fallbackActive = true;

    console.warn(`Feature "${feature}" has been degraded using fallback: ${state.config.fallbackBehavior}`);

    // Notify listeners
    this.notifyListeners(feature, state);

    // Schedule recovery attempt if retryable
    if (state.config.retryable && state.config.gracePeriod) {
      setTimeout(() => {
        this.attemptRecovery(feature);
      }, state.config.gracePeriod);
    }
  }

  // Attempt to recover a degraded feature
  private attemptRecovery(feature: string) {
    const state = this.features.get(feature);
    if (!state || !state.degraded) return;

    // Reset failure count and try to recover
    state.failures = Math.max(0, state.failures - 1);
    
    if (state.failures === 0) {
      state.degraded = false;
      state.available = true;
      state.fallbackActive = false;
      
      console.info(`Feature "${feature}" has been recovered`);
      this.notifyListeners(feature, state);
    }
  }

  // Report success for a feature (helps with recovery)
  reportSuccess(feature: string) {
    const state = this.features.get(feature);
    if (!state) return;

    if (state.failures > 0) {
      state.failures = Math.max(0, state.failures - 1);
    }

    if (state.degraded && state.failures === 0) {
      state.degraded = false;
      state.available = true;
      state.fallbackActive = false;
      
      console.info(`Feature "${feature}" has been recovered after success`);
      this.notifyListeners(feature, state);
    }
  }

  // Get current state of a feature
  getFeatureState(feature: string): FeatureState | null {
    return this.features.get(feature) || null;
  }

  // Check if a feature is available
  isFeatureAvailable(feature: string): boolean {
    const state = this.features.get(feature);
    return state ? state.available : false;
  }

  // Check if a feature is degraded
  isFeatureDegraded(feature: string): boolean {
    const state = this.features.get(feature);
    return state ? state.degraded : false;
  }

  // Get fallback behavior for a feature
  getFallbackBehavior(feature: string): string | null {
    const state = this.features.get(feature);
    return state ? state.config.fallbackBehavior : null;
  }

  // Subscribe to feature state changes
  subscribe(feature: string, listener: (state: FeatureState) => void): () => void {
    const listeners = this.listeners.get(feature);
    if (listeners) {
      listeners.add(listener);
    }

    // Return unsubscribe function
    return () => {
      const listeners = this.listeners.get(feature);
      if (listeners) {
        listeners.delete(listener);
      }
    };
  }

  // Notify all listeners of a feature state change
  private notifyListeners(feature: string, state: FeatureState) {
    const listeners = this.listeners.get(feature);
    if (listeners) {
      listeners.forEach(listener => {
        try {
          listener(state);
        } catch (error) {
          console.error('Error in degradation listener:', error);
        }
      });
    }
  }

  // Force degrade a feature (for testing or manual control)
  forceDegrade(feature: string) {
    const state = this.features.get(feature);
    if (state) {
      this.degradeFeature(feature);
    }
  }

  // Force recover a feature (for testing or manual control)
  forceRecover(feature: string) {
    const state = this.features.get(feature);
    if (state) {
      state.failures = 0;
      state.degraded = false;
      state.available = true;
      state.fallbackActive = false;
      this.notifyListeners(feature, state);
    }
  }

  // Get all feature states
  getAllFeatureStates(): Record<string, FeatureState> {
    const states: Record<string, FeatureState> = {};
    this.features.forEach((state, feature) => {
      states[feature] = state;
    });
    return states;
  }

  // Helper method to check if error is APIError
  private isAPIError(error: any): error is APIError {
    return error && 
           typeof error === 'object' && 
           'type' in error && 
           'code' in error && 
           'message' in error;
  }
}

// Global instance
export const gracefulDegradation = new GracefulDegradationManager();

// React hook for using graceful degradation
export const useGracefulDegradation = (feature: string) => {
  const [state, setState] = React.useState<FeatureState | null>(
    gracefulDegradation.getFeatureState(feature)
  );

  React.useEffect(() => {
    const unsubscribe = gracefulDegradation.subscribe(feature, setState);
    return unsubscribe;
  }, [feature]);

  const reportFailure = React.useCallback((error: APIError | Error) => {
    return gracefulDegradation.reportFailure(feature, error);
  }, [feature]);

  const reportSuccess = React.useCallback(() => {
    gracefulDegradation.reportSuccess(feature);
  }, [feature]);

  return {
    state,
    isAvailable: state?.available ?? false,
    isDegraded: state?.degraded ?? false,
    fallbackBehavior: state?.config.fallbackBehavior ?? 'hide',
    reportFailure,
    reportSuccess
  };
};

// HOC for graceful degradation
export const withGracefulDegradation = <P extends object>(
  Component: React.ComponentType<P>,
  feature: string,
  fallbackComponent?: React.ComponentType<P>
) => {
  return React.forwardRef<any, P>((props, ref) => {
    const { isAvailable, isDegraded, fallbackBehavior } = useGracefulDegradation(feature);

    if (!isAvailable || isDegraded) {
      switch (fallbackBehavior) {
        case 'hide':
          return null;
        case 'disable':
          return <Component {...props} ref={ref} disabled />;
        case 'simplify':
          return <Component {...props} ref={ref} simplified />;
        case 'mock':
          return fallbackComponent ? <fallbackComponent {...props} ref={ref} /> : null;
        default:
          return <Component {...props} ref={ref} />;
      }
    }

    return <Component {...props} ref={ref} />;
  });
};

// Utility functions
export const handleFeatureError = (feature: string, error: APIError | Error) => {
  const degraded = gracefulDegradation.reportFailure(feature, error);
  
  if (degraded) {
    console.warn(`Feature "${feature}" has been degraded due to error:`, error);
  }
  
  return degraded;
};

export const handleFeatureSuccess = (feature: string) => {
  gracefulDegradation.reportSuccess(feature);
};

// Error boundary with degradation support
interface DegradationErrorBoundaryProps {
  children: React.ReactNode;
  feature: string;
  fallback?: React.ReactNode;
}

export class DegradationErrorBoundary extends React.Component<
  DegradationErrorBoundaryProps,
  { hasError: boolean }
> {
  constructor(props: DegradationErrorBoundaryProps) {
    super(props);
    this.state = { hasError: false };
  }

  static getDerivedStateFromError(error: Error) {
    return { hasError: true };
  }

  componentDidCatch(error: Error, errorInfo: React.ErrorInfo) {
    console.error('DegradationErrorBoundary caught error:', error, errorInfo);
    
    // Report the error to the degradation manager
    gracefulDegradation.reportFailure(this.props.feature, error);
  }

  render() {
    if (this.state.hasError) {
      return this.props.fallback || <div>Feature temporarily unavailable</div>;
    }

    return this.props.children;
  }
}

export default gracefulDegradation;

// Import React for the hooks
import React from 'react';