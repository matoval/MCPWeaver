import { useCallback, useRef, useState, useEffect } from 'react';

export interface DragAndDropConfig {
  accept?: string[];
  multiple?: boolean;
  maxSize?: number;
  disabled?: boolean;
  onDrop?: (files: File[]) => void;
  onDragEnter?: (event: DragEvent) => void;
  onDragOver?: (event: DragEvent) => void;
  onDragLeave?: (event: DragEvent) => void;
  onError?: (error: string) => void;
}

export interface DragAndDropState {
  isDragging: boolean;
  isOver: boolean;
  files: File[];
}

export const useDragAndDrop = (config: DragAndDropConfig = {}) => {
  const {
    accept = [],
    multiple = false,
    maxSize = 10 * 1024 * 1024, // 10MB default
    disabled = false,
    onDrop,
    onDragEnter,
    onDragOver,
    onDragLeave,
    onError
  } = config;

  const [state, setState] = useState<DragAndDropState>({
    isDragging: false,
    isOver: false,
    files: []
  });

  const dragCounterRef = useRef(0);
  const dropZoneRef = useRef<HTMLElement | null>(null);

  const validateFile = useCallback((file: File): boolean => {
    // Check file type
    if (accept.length > 0) {
      const fileType = file.type;
      const fileName = file.name;
      const fileExtension = fileName.split('.').pop()?.toLowerCase();

      const isValidType = accept.some(acceptType => {
        if (acceptType.startsWith('.')) {
          // Extension check
          return acceptType.slice(1).toLowerCase() === fileExtension;
        } else if (acceptType.includes('*')) {
          // MIME type with wildcard
          const [type] = acceptType.split('/');
          return fileType.startsWith(type);
        } else {
          // Exact MIME type
          return fileType === acceptType;
        }
      });

      if (!isValidType) {
        onError?.(`File type not supported: ${fileType || fileExtension}`);
        return false;
      }
    }

    // Check file size
    if (file.size > maxSize) {
      onError?.(`File too large: ${(file.size / 1024 / 1024).toFixed(2)}MB. Maximum size: ${(maxSize / 1024 / 1024).toFixed(2)}MB`);
      return false;
    }

    return true;
  }, [accept, maxSize, onError]);

  const processFiles = useCallback((fileList: FileList | File[]): File[] => {
    const files = Array.from(fileList);
    
    if (!multiple && files.length > 1) {
      onError?.('Multiple files not allowed');
      return [];
    }

    const validFiles = files.filter(validateFile);
    
    if (validFiles.length === 0 && files.length > 0) {
      return [];
    }

    return validFiles;
  }, [multiple, validateFile, onError]);

  const handleDragEnter = useCallback((event: DragEvent) => {
    event.preventDefault();
    event.stopPropagation();

    if (disabled) return;

    dragCounterRef.current++;

    if (event.dataTransfer?.items) {
      const items = Array.from(event.dataTransfer.items);
      const hasFiles = items.some(item => item.kind === 'file');
      
      if (hasFiles) {
        setState(prev => ({ ...prev, isDragging: true }));
        onDragEnter?.(event);
      }
    }
  }, [disabled, onDragEnter]);

  const handleDragOver = useCallback((event: DragEvent) => {
    event.preventDefault();
    event.stopPropagation();

    if (disabled) return;

    // Set the appropriate drop effect
    if (event.dataTransfer) {
      event.dataTransfer.dropEffect = 'copy';
    }

    setState(prev => ({ ...prev, isOver: true }));
    onDragOver?.(event);
  }, [disabled, onDragOver]);

  const handleDragLeave = useCallback((event: DragEvent) => {
    event.preventDefault();
    event.stopPropagation();

    if (disabled) return;

    dragCounterRef.current--;

    if (dragCounterRef.current === 0) {
      setState(prev => ({ ...prev, isDragging: false, isOver: false }));
      onDragLeave?.(event);
    }
  }, [disabled, onDragLeave]);

  const handleDrop = useCallback((event: DragEvent) => {
    event.preventDefault();
    event.stopPropagation();

    if (disabled) return;

    dragCounterRef.current = 0;
    setState(prev => ({ ...prev, isDragging: false, isOver: false }));

    const files = event.dataTransfer?.files;
    if (files && files.length > 0) {
      const validFiles = processFiles(files);
      
      if (validFiles.length > 0) {
        setState(prev => ({ ...prev, files: validFiles }));
        onDrop?.(validFiles);
      }
    }
  }, [disabled, processFiles, onDrop]);

  const bindDropZone = useCallback((element: HTMLElement | null) => {
    if (dropZoneRef.current) {
      // Remove event listeners from previous element
      dropZoneRef.current.removeEventListener('dragenter', handleDragEnter);
      dropZoneRef.current.removeEventListener('dragover', handleDragOver);
      dropZoneRef.current.removeEventListener('dragleave', handleDragLeave);
      dropZoneRef.current.removeEventListener('drop', handleDrop);
    }

    dropZoneRef.current = element;

    if (element) {
      // Add event listeners to new element
      element.addEventListener('dragenter', handleDragEnter);
      element.addEventListener('dragover', handleDragOver);
      element.addEventListener('dragleave', handleDragLeave);
      element.addEventListener('drop', handleDrop);
    }
  }, [handleDragEnter, handleDragOver, handleDragLeave, handleDrop]);

  const getInputProps = useCallback(() => ({
    type: 'file' as const,
    accept: accept.join(','),
    multiple,
    onChange: (event: React.ChangeEvent<HTMLInputElement>) => {
      const files = event.target.files;
      if (files && files.length > 0) {
        const validFiles = processFiles(files);
        if (validFiles.length > 0) {
          setState(prev => ({ ...prev, files: validFiles }));
          onDrop?.(validFiles);
        }
      }
      // Reset input value to allow selecting the same file again
      event.target.value = '';
    }
  }), [accept, multiple, processFiles, onDrop]);

  const getDropZoneProps = useCallback(() => ({
    ref: bindDropZone,
    'data-dragging': state.isDragging,
    'data-over': state.isOver,
    'aria-label': 'Drop files here or click to select',
    role: 'button',
    tabIndex: 0,
    onKeyDown: (event: React.KeyboardEvent) => {
      if (event.key === 'Enter' || event.key === ' ') {
        event.preventDefault();
        // Trigger file input click
        const input = document.createElement('input');
        Object.assign(input, getInputProps());
        input.click();
      }
    }
  }), [bindDropZone, state.isDragging, state.isOver, getInputProps]);

  const openFileDialog = useCallback(() => {
    const input = document.createElement('input');
    Object.assign(input, getInputProps());
    input.click();
  }, [getInputProps]);

  const clearFiles = useCallback(() => {
    setState(prev => ({ ...prev, files: [] }));
  }, []);

  const removeFile = useCallback((index: number) => {
    setState(prev => ({
      ...prev,
      files: prev.files.filter((_, i) => i !== index)
    }));
  }, []);

  // Cleanup on unmount
  useEffect(() => {
    return () => {
      if (dropZoneRef.current) {
        dropZoneRef.current.removeEventListener('dragenter', handleDragEnter);
        dropZoneRef.current.removeEventListener('dragover', handleDragOver);
        dropZoneRef.current.removeEventListener('dragleave', handleDragLeave);
        dropZoneRef.current.removeEventListener('drop', handleDrop);
      }
    };
  }, [handleDragEnter, handleDragOver, handleDragLeave, handleDrop]);

  // Global drag events for visual feedback
  useEffect(() => {
    if (disabled) return;

    const handleGlobalDragEnter = (event: DragEvent) => {
      if (event.dataTransfer?.items) {
        const items = Array.from(event.dataTransfer.items);
        const hasFiles = items.some(item => item.kind === 'file');
        
        if (hasFiles) {
          setState(prev => ({ ...prev, isDragging: true }));
        }
      }
    };

    const handleGlobalDrop = () => {
      dragCounterRef.current = 0;
      setState(prev => ({ ...prev, isDragging: false, isOver: false }));
    };

    document.addEventListener('dragenter', handleGlobalDragEnter);
    document.addEventListener('drop', handleGlobalDrop);
    document.addEventListener('dragend', handleGlobalDrop);

    return () => {
      document.removeEventListener('dragenter', handleGlobalDragEnter);
      document.removeEventListener('drop', handleGlobalDrop);
      document.removeEventListener('dragend', handleGlobalDrop);
    };
  }, [disabled]);

  const getFileInfo = useCallback((file: File) => ({
    name: file.name,
    size: file.size,
    type: file.type,
    lastModified: file.lastModified,
    sizeString: formatFileSize(file.size),
    extension: file.name.split('.').pop()?.toLowerCase() || ''
  }), []);

  return {
    // State
    ...state,
    
    // Methods
    getInputProps,
    getDropZoneProps,
    openFileDialog,
    clearFiles,
    removeFile,
    getFileInfo,
    
    // Utilities
    validateFile,
    processFiles
  };
};

// Utility function to format file size
const formatFileSize = (bytes: number): string => {
  if (bytes === 0) return '0 Bytes';
  
  const k = 1024;
  const sizes = ['Bytes', 'KB', 'MB', 'GB', 'TB'];
  const i = Math.floor(Math.log(bytes) / Math.log(k));
  
  return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i];
};

export default useDragAndDrop;