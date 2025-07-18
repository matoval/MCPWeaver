import { useState, useEffect, useCallback, useRef } from 'react';
import { wails } from '../services/wails';
import { 
  GenerationJob, 
  GenerationProgress, 
  ProgressMetrics, 
  ProgressHistoryEntry,
  GenerationStatus,
  GenerationStartedEvent,
  GenerationProgressEvent,
  GenerationCompletedEvent,
  GenerationFailedEvent,
  GenerationCancelledEvent
} from '../types';

export interface UseProgressReturn {
  // Current job state
  currentJob: GenerationJob | null;
  progress: number;
  currentStep: string;
  isGenerating: boolean;
  
  // Progress details
  metrics: ProgressMetrics | null;
  history: ProgressHistoryEntry[];
  
  // Actions
  startGeneration: (projectId: string) => Promise<void>;
  cancelGeneration: (jobId: string) => Promise<void>;
  clearHistory: () => void;
  
  // Error state
  error: string | null;
  clearError: () => void;
}

export function useProgress(): UseProgressReturn {
  const [currentJob, setCurrentJob] = useState<GenerationJob | null>(null);
  const [progress, setProgress] = useState(0);
  const [currentStep, setCurrentStep] = useState('');
  const [isGenerating, setIsGenerating] = useState(false);
  const [metrics, setMetrics] = useState<ProgressMetrics | null>(null);
  const [history, setHistory] = useState<ProgressHistoryEntry[]>([]);
  const [error, setError] = useState<string | null>(null);
  
  const progressStartTimeRef = useRef<number | null>(null);
  const progressDataRef = useRef<GenerationProgress[]>([]);
  const cleanupFunctionsRef = useRef<(() => void)[]>([]);

  // Calculate metrics based on progress data
  const calculateMetrics = useCallback((progressData: GenerationProgress[], job: GenerationJob | null): ProgressMetrics | null => {
    if (!progressData.length || !job) return null;

    const now = new Date().getTime();
    const startTime = new Date(job.startTime).getTime();
    const elapsedTime = now - startTime;
    const currentProgress = progressData[progressData.length - 1];
    
    // Calculate processing rate (progress per second)
    const processingRate = currentProgress.progress / (elapsedTime / 1000);
    
    // Estimate time remaining
    const remainingProgress = 1 - currentProgress.progress;
    const estimatedTimeRemaining = processingRate > 0 ? remainingProgress / processingRate : 0;
    
    return {
      startTime: job.startTime,
      currentTime: new Date().toISOString(),
      elapsedTime,
      processingRate,
      estimatedTimeRemaining,
      filesGenerated: currentProgress.filesGenerated || 0,
      totalFiles: 3, // main.go, go.mod, README.md as baseline
      errorCount: currentProgress.errorCount || 0,
      warningCount: currentProgress.warningCount || 0,
      memoryUsage: currentProgress.memoryUsage || 0,
      cpuUsage: 0 // Would need system monitoring for accurate CPU usage
    };
  }, []);

  // Update metrics when progress changes
  useEffect(() => {
    if (currentJob && progressDataRef.current.length > 0) {
      const newMetrics = calculateMetrics(progressDataRef.current, currentJob);
      setMetrics(newMetrics);
    }
  }, [currentJob, calculateMetrics]);

  // Event handlers
  const handleGenerationStarted = useCallback((event: GenerationStartedEvent) => {
    const job = event.data;
    setCurrentJob(job);
    setProgress(0);
    setCurrentStep('Initializing generation');
    setIsGenerating(true);
    setError(null);
    progressStartTimeRef.current = Date.now();
    progressDataRef.current = [];
    
    console.log('Generation started:', job);
  }, []);

  const handleGenerationProgress = useCallback((event: GenerationProgressEvent) => {
    const progressData = event.data;
    setProgress(progressData.progress);
    setCurrentStep(progressData.currentStep);
    
    // Store progress data for metrics calculation
    progressDataRef.current.push(progressData);
    
    console.log('Generation progress:', progressData);
  }, []);

  const handleGenerationCompleted = useCallback((event: GenerationCompletedEvent) => {
    const job = event.data;
    setCurrentJob(job);
    setProgress(1);
    setCurrentStep('Generation completed successfully');
    setIsGenerating(false);
    
    // Add to history
    const historyEntry: ProgressHistoryEntry = {
      jobId: job.id,
      projectId: job.projectId,
      projectName: `Project ${job.projectId}`, // Would need to fetch project name
      status: job.status,
      progress: 1,
      startTime: job.startTime,
      endTime: job.endTime,
      duration: job.endTime ? new Date(job.endTime).getTime() - new Date(job.startTime).getTime() : 0,
      success: true,
      statistics: job.results?.statistics
    };
    
    setHistory(prev => [historyEntry, ...prev].slice(0, 50)); // Keep last 50 entries
    
    // Reset after a delay
    setTimeout(() => {
      setCurrentJob(null);
      setProgress(0);
      setCurrentStep('');
      setMetrics(null);
      progressDataRef.current = [];
    }, 3000);
    
    console.log('Generation completed:', job);
  }, []);

  const handleGenerationFailed = useCallback((event: GenerationFailedEvent) => {
    const { jobId, message } = event.data;
    setIsGenerating(false);
    setError(message);
    
    // Add to history
    const historyEntry: ProgressHistoryEntry = {
      jobId,
      projectId: currentJob?.projectId || 'unknown',
      projectName: `Project ${currentJob?.projectId || 'unknown'}`,
      status: 'failed',
      progress: progress,
      startTime: currentJob?.startTime || new Date().toISOString(),
      endTime: new Date().toISOString(),
      duration: currentJob ? Date.now() - new Date(currentJob.startTime).getTime() : 0,
      success: false,
      errorMessage: message
    };
    
    setHistory(prev => [historyEntry, ...prev].slice(0, 50));
    
    console.error('Generation failed:', message);
  }, [currentJob, progress]);

  const handleGenerationCancelled = useCallback((event: GenerationCancelledEvent) => {
    const job = event.data;
    setCurrentJob(job);
    setIsGenerating(false);
    setCurrentStep('Generation cancelled');
    
    // Add to history
    const historyEntry: ProgressHistoryEntry = {
      jobId: job.id,
      projectId: job.projectId,
      projectName: `Project ${job.projectId}`,
      status: 'cancelled',
      progress: progress,
      startTime: job.startTime,
      endTime: job.endTime,
      duration: job.endTime ? new Date(job.endTime).getTime() - new Date(job.startTime).getTime() : 0,
      success: false,
      errorMessage: 'Generation was cancelled by user'
    };
    
    setHistory(prev => [historyEntry, ...prev].slice(0, 50));
    
    // Reset after a delay
    setTimeout(() => {
      setCurrentJob(null);
      setProgress(0);
      setCurrentStep('');
      setMetrics(null);
      progressDataRef.current = [];
    }, 2000);
    
    console.log('Generation cancelled:', job);
  }, [progress]);

  // Set up event listeners
  useEffect(() => {
    const cleanup1 = wails.onEvent('generation:started', handleGenerationStarted);
    const cleanup2 = wails.onEvent('generation:progress', handleGenerationProgress);
    const cleanup3 = wails.onEvent('generation:completed', handleGenerationCompleted);
    const cleanup4 = wails.onEvent('generation:failed', handleGenerationFailed);
    const cleanup5 = wails.onEvent('generation:cancelled', handleGenerationCancelled);
    
    cleanupFunctionsRef.current = [cleanup1, cleanup2, cleanup3, cleanup4, cleanup5];
    
    return () => {
      cleanupFunctionsRef.current.forEach(cleanup => cleanup());
    };
  }, [
    handleGenerationStarted,
    handleGenerationProgress,
    handleGenerationCompleted,
    handleGenerationFailed,
    handleGenerationCancelled
  ]);

  // Load persisted history on mount
  useEffect(() => {
    const savedHistory = localStorage.getItem('mcpweaver-progress-history');
    if (savedHistory) {
      try {
        const parsed = JSON.parse(savedHistory);
        setHistory(parsed);
      } catch (error) {
        console.error('Failed to load progress history:', error);
      }
    }
  }, []);

  // Persist history to localStorage
  useEffect(() => {
    if (history.length > 0) {
      localStorage.setItem('mcpweaver-progress-history', JSON.stringify(history));
    }
  }, [history]);

  // Actions
  const startGeneration = useCallback(async (projectId: string) => {
    try {
      setError(null);
      await wails.generateServer(projectId);
    } catch (error: any) {
      setError(error.message || 'Failed to start generation');
      setIsGenerating(false);
    }
  }, []);

  const cancelGeneration = useCallback(async (jobId: string) => {
    try {
      setError(null);
      await wails.cancelGenerationJob(jobId);
    } catch (error: any) {
      setError(error.message || 'Failed to cancel generation');
    }
  }, []);

  const clearHistory = useCallback(() => {
    setHistory([]);
    localStorage.removeItem('mcpweaver-progress-history');
  }, []);

  const clearError = useCallback(() => {
    setError(null);
  }, []);

  return {
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
  };
}