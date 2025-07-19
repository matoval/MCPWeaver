import { useState, useEffect, useCallback, useRef } from 'react';

export interface AutoSaveConfig {
  delay?: number;
  key: string;
  enabled?: boolean;
  onSave?: (data: any) => Promise<void> | void;
  onRestore?: (data: any) => void;
  onError?: (error: Error) => void;
}

export interface AutoSaveState {
  isSaving: boolean;
  lastSaved: number | null;
  hasUnsavedChanges: boolean;
  error: string | null;
}

export const useAutoSave = <T>(
  data: T,
  config: AutoSaveConfig
) => {
  const {
    delay = 2000,
    key,
    enabled = true,
    onSave,
    onRestore,
    onError
  } = config;

  const [state, setState] = useState<AutoSaveState>({
    isSaving: false,
    lastSaved: null,
    hasUnsavedChanges: false,
    error: null
  });

  const saveTimeoutRef = useRef<NodeJS.Timeout>();
  const lastSavedDataRef = useRef<string>('');
  const isInitialLoadRef = useRef(true);

  // Save to localStorage
  const saveToStorage = useCallback(async (dataToSave: T) => {
    try {
      setState(prev => ({ ...prev, isSaving: true, error: null }));
      
      const serializedData = JSON.stringify(dataToSave);
      localStorage.setItem(key, serializedData);
      
      // Call custom save function if provided
      if (onSave) {
        await onSave(dataToSave);
      }
      
      lastSavedDataRef.current = serializedData;
      
      setState(prev => ({
        ...prev,
        isSaving: false,
        lastSaved: Date.now(),
        hasUnsavedChanges: false
      }));
    } catch (error) {
      const errorMessage = error instanceof Error ? error.message : 'Save failed';
      setState(prev => ({
        ...prev,
        isSaving: false,
        error: errorMessage
      }));
      
      if (onError) {
        onError(error instanceof Error ? error : new Error(errorMessage));
      }
    }
  }, [key, onSave, onError]);

  // Restore from localStorage
  const restoreFromStorage = useCallback(() => {
    try {
      const savedData = localStorage.getItem(key);
      if (savedData) {
        const parsedData = JSON.parse(savedData);
        lastSavedDataRef.current = savedData;
        
        if (onRestore) {
          onRestore(parsedData);
        }
        
        setState(prev => ({
          ...prev,
          lastSaved: Date.now(),
          hasUnsavedChanges: false
        }));
        
        return parsedData;
      }
    } catch (error) {
      const errorMessage = error instanceof Error ? error.message : 'Restore failed';
      setState(prev => ({ ...prev, error: errorMessage }));
      
      if (onError) {
        onError(error instanceof Error ? error : new Error(errorMessage));
      }
    }
    
    return null;
  }, [key, onRestore, onError]);

  // Manual save
  const save = useCallback(() => {
    if (enabled && data !== undefined) {
      saveToStorage(data);
    }
  }, [enabled, data, saveToStorage]);

  // Clear saved data
  const clearSaved = useCallback(() => {
    try {
      localStorage.removeItem(key);
      lastSavedDataRef.current = '';
      setState(prev => ({
        ...prev,
        lastSaved: null,
        hasUnsavedChanges: false,
        error: null
      }));
    } catch (error) {
      const errorMessage = error instanceof Error ? error.message : 'Clear failed';
      setState(prev => ({ ...prev, error: errorMessage }));
      
      if (onError) {
        onError(error instanceof Error ? error : new Error(errorMessage));
      }
    }
  }, [key, onError]);

  // Check if data has changed
  useEffect(() => {
    if (!enabled || data === undefined) return;

    const currentData = JSON.stringify(data);
    
    // Skip on initial load
    if (isInitialLoadRef.current) {
      isInitialLoadRef.current = false;
      lastSavedDataRef.current = currentData;
      return;
    }

    const hasChanges = currentData !== lastSavedDataRef.current;
    
    setState(prev => ({ ...prev, hasUnsavedChanges: hasChanges }));

    if (hasChanges) {
      // Clear existing timeout
      if (saveTimeoutRef.current) {
        clearTimeout(saveTimeoutRef.current);
      }

      // Set new timeout for auto-save
      saveTimeoutRef.current = setTimeout(() => {
        saveToStorage(data);
      }, delay);
    }

    return () => {
      if (saveTimeoutRef.current) {
        clearTimeout(saveTimeoutRef.current);
      }
    };
  }, [data, enabled, delay, saveToStorage]);

  // Save before page unload
  useEffect(() => {
    const handleBeforeUnload = (event: BeforeUnloadEvent) => {
      if (state.hasUnsavedChanges && enabled) {
        event.preventDefault();
        event.returnValue = 'You have unsaved changes. Are you sure you want to leave?';
        
        // Try to save immediately
        if (data !== undefined) {
          try {
            const serializedData = JSON.stringify(data);
            localStorage.setItem(key, serializedData);
          } catch (error) {
            console.warn('Failed to save before unload:', error);
          }
        }
      }
    };

    window.addEventListener('beforeunload', handleBeforeUnload);
    
    return () => {
      window.removeEventListener('beforeunload', handleBeforeUnload);
    };
  }, [state.hasUnsavedChanges, enabled, data, key]);

  // Cleanup timeout on unmount
  useEffect(() => {
    return () => {
      if (saveTimeoutRef.current) {
        clearTimeout(saveTimeoutRef.current);
      }
    };
  }, []);

  // Page visibility API for saving when tab becomes hidden
  useEffect(() => {
    if (!enabled) return;

    const handleVisibilityChange = () => {
      if (document.hidden && state.hasUnsavedChanges && data !== undefined) {
        saveToStorage(data);
      }
    };

    document.addEventListener('visibilitychange', handleVisibilityChange);
    
    return () => {
      document.removeEventListener('visibilitychange', handleVisibilityChange);
    };
  }, [enabled, state.hasUnsavedChanges, data, saveToStorage]);

  const getTimeSinceLastSave = useCallback(() => {
    if (!state.lastSaved) return null;
    return Date.now() - state.lastSaved;
  }, [state.lastSaved]);

  const formatLastSaveTime = useCallback(() => {
    if (!state.lastSaved) return 'Never';
    
    const timeSince = getTimeSinceLastSave();
    if (!timeSince) return 'Never';
    
    if (timeSince < 60000) {
      return 'Just now';
    } else if (timeSince < 3600000) {
      return `${Math.floor(timeSince / 60000)} minute(s) ago`;
    } else {
      return new Date(state.lastSaved).toLocaleTimeString();
    }
  }, [state.lastSaved, getTimeSinceLastSave]);

  return {
    // State
    ...state,
    
    // Actions
    save,
    restoreFromStorage,
    clearSaved,
    
    // Utilities
    getTimeSinceLastSave,
    formatLastSaveTime,
    
    // Configuration
    config: {
      delay,
      key,
      enabled
    }
  };
};

export default useAutoSave;