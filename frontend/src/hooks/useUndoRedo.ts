import { useState, useCallback, useRef } from 'react';

export interface UndoRedoAction<T = any> {
  id: string;
  name: string;
  timestamp: number;
  undo: () => void;
  redo: () => void;
  data?: T;
}

export interface UndoRedoState {
  canUndo: boolean;
  canRedo: boolean;
  currentIndex: number;
  actions: UndoRedoAction[];
}

export const useUndoRedo = <T = any>(maxHistorySize: number = 50) => {
  const [state, setState] = useState<UndoRedoState>({
    canUndo: false,
    canRedo: false,
    currentIndex: -1,
    actions: []
  });

  const actionsRef = useRef<UndoRedoAction<T>[]>([]);
  const currentIndexRef = useRef(-1);

  const updateState = useCallback(() => {
    setState({
      canUndo: currentIndexRef.current >= 0,
      canRedo: currentIndexRef.current < actionsRef.current.length - 1,
      currentIndex: currentIndexRef.current,
      actions: [...actionsRef.current]
    });
  }, []);

  const execute = useCallback((action: Omit<UndoRedoAction<T>, 'timestamp'>) => {
    const fullAction: UndoRedoAction<T> = {
      ...action,
      timestamp: Date.now()
    };

    // Remove any actions after current index (when we're in the middle of history)
    if (currentIndexRef.current < actionsRef.current.length - 1) {
      actionsRef.current = actionsRef.current.slice(0, currentIndexRef.current + 1);
    }

    // Add the new action
    actionsRef.current.push(fullAction);
    currentIndexRef.current = actionsRef.current.length - 1;

    // Limit history size
    if (actionsRef.current.length > maxHistorySize) {
      actionsRef.current = actionsRef.current.slice(-maxHistorySize);
      currentIndexRef.current = actionsRef.current.length - 1;
    }

    updateState();
  }, [maxHistorySize, updateState]);

  const undo = useCallback(() => {
    if (currentIndexRef.current >= 0) {
      const action = actionsRef.current[currentIndexRef.current];
      action.undo();
      currentIndexRef.current--;
      updateState();
    }
  }, [updateState]);

  const redo = useCallback(() => {
    if (currentIndexRef.current < actionsRef.current.length - 1) {
      currentIndexRef.current++;
      const action = actionsRef.current[currentIndexRef.current];
      action.redo();
      updateState();
    }
  }, [updateState]);

  const clear = useCallback(() => {
    actionsRef.current = [];
    currentIndexRef.current = -1;
    updateState();
  }, [updateState]);

  const getHistory = useCallback(() => {
    return actionsRef.current.map((action, index) => ({
      ...action,
      isCurrent: index === currentIndexRef.current,
      canUndoTo: index <= currentIndexRef.current,
      canRedoTo: index > currentIndexRef.current
    }));
  }, []);

  const undoTo = useCallback((index: number) => {
    if (index < 0 || index >= actionsRef.current.length || index > currentIndexRef.current) {
      return;
    }

    // Undo actions from current to target index
    while (currentIndexRef.current > index) {
      undo();
    }
  }, [undo]);

  const redoTo = useCallback((index: number) => {
    if (index < 0 || index >= actionsRef.current.length || index <= currentIndexRef.current) {
      return;
    }

    // Redo actions from current to target index
    while (currentIndexRef.current < index) {
      redo();
    }
  }, [redo]);

  return {
    // State
    ...state,
    
    // Actions
    execute,
    undo,
    redo,
    clear,
    
    // History navigation
    getHistory,
    undoTo,
    redoTo,
    
    // Utilities
    getLastAction: () => actionsRef.current[currentIndexRef.current] || null,
    getNextAction: () => actionsRef.current[currentIndexRef.current + 1] || null
  };
};

export default useUndoRedo;