import React, { useRef } from 'react';
import useDragAndDrop from '../../hooks/useDragAndDrop';
import useAnimations from '../../hooks/useAnimations';

interface DragDropZoneProps {
  accept?: string[];
  multiple?: boolean;
  maxSize?: number;
  disabled?: boolean;
  className?: string;
  children?: React.ReactNode;
  placeholder?: string;
  onFilesSelected?: (files: File[]) => void;
  onError?: (error: string) => void;
  showPreview?: boolean;
  allowClick?: boolean;
}

const DragDropZone: React.FC<DragDropZoneProps> = ({
  accept = ['.json', '.yaml', '.yml', 'application/json', 'text/yaml'],
  multiple = false,
  maxSize = 10 * 1024 * 1024,
  disabled = false,
  className = '',
  children,
  placeholder = 'Drop OpenAPI specification files here or click to browse',
  onFilesSelected,
  onError,
  showPreview = true,
  allowClick = true
}) => {
  const fileInputRef = useRef<HTMLInputElement>(null);
  const animations = useAnimations();

  const {
    isDragging,
    isOver,
    files,
    getInputProps,
    getDropZoneProps,
    openFileDialog,
    clearFiles,
    removeFile,
    getFileInfo
  } = useDragAndDrop({
    accept,
    multiple,
    maxSize,
    disabled,
    onDrop: onFilesSelected,
    onError
  });

  const handleClick = () => {
    if (allowClick && !disabled) {
      openFileDialog();
    }
  };

  const handleKeyDown = (event: React.KeyboardEvent) => {
    if ((event.key === 'Enter' || event.key === ' ') && allowClick && !disabled) {
      event.preventDefault();
      openFileDialog();
    }
  };

  const getDropZoneClasses = () => {
    const baseClasses = [
      'drag-drop-zone',
      className
    ];

    if (isDragging) baseClasses.push('dragging');
    if (isOver) baseClasses.push('drag-over');
    if (disabled) baseClasses.push('disabled');
    if (files.length > 0) baseClasses.push('has-files');

    return baseClasses.filter(Boolean).join(' ');
  };

  return (
    <div className="drag-drop-container">
      <div
        {...getDropZoneProps()}
        className={getDropZoneClasses()}
        onClick={handleClick}
        onKeyDown={handleKeyDown}
      >
        <input
          ref={fileInputRef}
          {...getInputProps()}
          style={{ display: 'none' }}
        />

        {children || (
          <div className="drop-zone-content">
            <div className="drop-zone-icon">
              {isDragging || isOver ? (
                <svg width="48" height="48" viewBox="0 0 24 24" fill="none" stroke="currentColor">
                  <path d="M21 15v4a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2v-4" />
                  <polyline points="7,10 12,15 17,10" />
                  <line x1="12" y1="15" x2="12" y2="3" />
                </svg>
              ) : (
                <svg width="48" height="48" viewBox="0 0 24 24" fill="none" stroke="currentColor">
                  <path d="M14 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V8z" />
                  <polyline points="14,2 14,8 20,8" />
                  <line x1="16" y1="13" x2="8" y2="13" />
                  <line x1="16" y1="17" x2="8" y2="17" />
                  <polyline points="10,9 9,9 8,9" />
                </svg>
              )}
            </div>
            
            <div className="drop-zone-text">
              <p className="primary-text">
                {isDragging || isOver ? 'Drop files here' : placeholder}
              </p>
              
              {!disabled && (
                <p className="secondary-text">
                  {accept.length > 0 && (
                    <span>Supported formats: {accept.join(', ')}</span>
                  )}
                  {maxSize && (
                    <span> • Max size: {(maxSize / 1024 / 1024).toFixed(0)}MB</span>
                  )}
                </p>
              )}
            </div>
          </div>
        )}
      </div>

      {showPreview && files.length > 0 && (
        <div className="file-preview">
          <div className="preview-header">
            <h4>Selected Files</h4>
            <button
              type="button"
              onClick={clearFiles}
              className="clear-button"
              aria-label="Clear all files"
            >
              Clear All
            </button>
          </div>
          
          <div className="file-list">
            {files.map((file, index) => {
              const fileInfo = getFileInfo(file);
              return (
                <div key={`${file.name}-${index}`} className="file-item">
                  <div className="file-icon">
                    <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor">
                      <path d="M14 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V8z" />
                      <polyline points="14,2 14,8 20,8" />
                    </svg>
                  </div>
                  
                  <div className="file-details">
                    <div className="file-name">{fileInfo.name}</div>
                    <div className="file-meta">
                      {fileInfo.sizeString} • {fileInfo.extension?.toUpperCase()}
                    </div>
                  </div>
                  
                  <button
                    type="button"
                    onClick={() => removeFile(index)}
                    className="remove-button"
                    aria-label={`Remove ${fileInfo.name}`}
                  >
                    <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor">
                      <line x1="18" y1="6" x2="6" y2="18" />
                      <line x1="6" y1="6" x2="18" y2="18" />
                    </svg>
                  </button>
                </div>
              );
            })}
          </div>
        </div>
      )}

      <style jsx>{`
        .drag-drop-container {
          width: 100%;
        }

        .drag-drop-zone {
          border: 2px dashed var(--border-color);
          border-radius: 8px;
          padding: 2rem;
          text-align: center;
          background: var(--bg-secondary);
          cursor: pointer;
          transition: all 0.2s ease;
          position: relative;
          min-height: 200px;
          display: flex;
          align-items: center;
          justify-content: center;
        }

        .drag-drop-zone:hover:not(.disabled) {
          border-color: var(--accent-color);
          background: var(--bg-tertiary);
        }

        .drag-drop-zone:focus {
          outline: 2px solid var(--accent-color);
          outline-offset: 2px;
        }

        .drag-drop-zone.dragging {
          border-color: var(--accent-color);
          background: var(--accent-color-alpha);
          transform: scale(1.02);
        }

        .drag-drop-zone.drag-over {
          border-color: var(--success-color);
          background: var(--success-color-alpha);
          border-style: solid;
        }

        .drag-drop-zone.disabled {
          opacity: 0.5;
          cursor: not-allowed;
          pointer-events: none;
        }

        .drag-drop-zone.has-files {
          border-color: var(--success-color);
          background: var(--success-color-light);
        }

        .drop-zone-content {
          display: flex;
          flex-direction: column;
          align-items: center;
          gap: 1rem;
        }

        .drop-zone-icon {
          color: var(--text-secondary);
          transition: color 0.2s ease;
        }

        .drag-drop-zone:hover .drop-zone-icon,
        .drag-drop-zone.dragging .drop-zone-icon,
        .drag-drop-zone.drag-over .drop-zone-icon {
          color: var(--accent-color);
        }

        .drop-zone-text {
          max-width: 400px;
        }

        .primary-text {
          font-size: 1.1rem;
          font-weight: 500;
          color: var(--text-primary);
          margin: 0 0 0.5rem 0;
        }

        .secondary-text {
          font-size: 0.9rem;
          color: var(--text-secondary);
          margin: 0;
          line-height: 1.4;
        }

        .file-preview {
          margin-top: 1.5rem;
          padding: 1rem;
          background: var(--bg-secondary);
          border-radius: 6px;
          border: 1px solid var(--border-color);
        }

        .preview-header {
          display: flex;
          align-items: center;
          justify-content: space-between;
          margin-bottom: 1rem;
        }

        .preview-header h4 {
          margin: 0;
          font-size: 1rem;
          font-weight: 600;
          color: var(--text-primary);
        }

        .clear-button {
          background: none;
          border: 1px solid var(--border-color);
          color: var(--text-secondary);
          padding: 0.5rem 1rem;
          border-radius: 4px;
          font-size: 0.85rem;
          cursor: pointer;
          transition: all 0.2s ease;
        }

        .clear-button:hover {
          background: var(--bg-tertiary);
          color: var(--text-primary);
        }

        .file-list {
          display: flex;
          flex-direction: column;
          gap: 0.75rem;
        }

        .file-item {
          display: flex;
          align-items: center;
          gap: 0.75rem;
          padding: 0.75rem;
          background: var(--bg-primary);
          border: 1px solid var(--border-color);
          border-radius: 4px;
          transition: background 0.2s ease;
        }

        .file-item:hover {
          background: var(--bg-tertiary);
        }

        .file-icon {
          color: var(--text-secondary);
          flex-shrink: 0;
        }

        .file-details {
          flex: 1;
          min-width: 0;
        }

        .file-name {
          font-weight: 500;
          color: var(--text-primary);
          truncate: ellipsis;
          overflow: hidden;
          white-space: nowrap;
        }

        .file-meta {
          font-size: 0.8rem;
          color: var(--text-secondary);
          margin-top: 0.25rem;
        }

        .remove-button {
          background: none;
          border: none;
          color: var(--text-secondary);
          cursor: pointer;
          padding: 0.25rem;
          border-radius: 2px;
          transition: all 0.2s ease;
          flex-shrink: 0;
        }

        .remove-button:hover {
          background: var(--error-color-light);
          color: var(--error-color);
        }

        @media (max-width: 768px) {
          .drag-drop-zone {
            padding: 1.5rem 1rem;
            min-height: 150px;
          }

          .drop-zone-icon svg {
            width: 36px;
            height: 36px;
          }

          .primary-text {
            font-size: 1rem;
          }

          .secondary-text {
            font-size: 0.85rem;
          }

          .file-preview {
            margin-top: 1rem;
            padding: 0.75rem;
          }

          .file-item {
            padding: 0.5rem;
            gap: 0.5rem;
          }
        }
      `}</style>
    </div>
  );
};

export default DragDropZone;