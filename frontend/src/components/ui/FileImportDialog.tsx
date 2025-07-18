import React, { useState } from 'react';
import { wails } from '../../services/wails';
import { ImportResult, FileFilter } from '../../types';
import './FileImportDialog.scss';

export interface FileImportDialogProps {
  isOpen: boolean;
  onClose: () => void;
  onImport: (result: ImportResult) => void;
  title?: string;
  className?: string;
}

export const FileImportDialog: React.FC<FileImportDialogProps> = ({
  isOpen,
  onClose,
  onImport,
  title = 'Import OpenAPI Specification',
  className = ''
}) => {
  const [activeTab, setActiveTab] = useState<'file' | 'url'>('file');
  const [url, setUrl] = useState('');
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [dragOver, setDragOver] = useState(false);

  // Handle file selection
  const handleFileSelect = async () => {
    try {
      setLoading(true);
      setError(null);
      
      const filters = await wails.getDefaultOpenAPIFilters();
      const filePath = await wails.selectFile(filters);
      
      if (!filePath) {
        setLoading(false);
        return; // User cancelled
      }
      
      const result = await wails.importOpenAPISpec(filePath);
      
      if (result.valid) {
        onImport(result);
        onClose();
      } else {
        setError(result.errors?.join(', ') || 'Failed to import specification');
      }
    } catch (err: any) {
      setError(err.message || 'Failed to import file');
    } finally {
      setLoading(false);
    }
  };

  // Handle URL import
  const handleUrlImport = async () => {
    if (!url.trim()) {
      setError('Please enter a valid URL');
      return;
    }

    try {
      setLoading(true);
      setError(null);
      
      const result = await wails.importOpenAPISpecFromURL(url.trim());
      
      if (result.valid) {
        onImport(result);
        onClose();
      } else {
        setError(result.errors?.join(', ') || 'Failed to import specification');
      }
    } catch (err: any) {
      setError(err.message || 'Failed to import from URL');
    } finally {
      setLoading(false);
    }
  };

  // Handle drag and drop
  const handleDragOver = (e: React.DragEvent) => {
    e.preventDefault();
    setDragOver(true);
  };

  const handleDragLeave = (e: React.DragEvent) => {
    e.preventDefault();
    setDragOver(false);
  };

  const handleDrop = async (e: React.DragEvent) => {
    e.preventDefault();
    setDragOver(false);
    
    const files = Array.from(e.dataTransfer.files);
    if (files.length === 0) return;
    
    const file = files[0];
    const validExtensions = ['.json', '.yaml', '.yml'];
    const fileExtension = file.name.toLowerCase().substring(file.name.lastIndexOf('.'));
    
    if (!validExtensions.includes(fileExtension)) {
      setError('Please drop a valid OpenAPI specification file (.json, .yaml, .yml)');
      return;
    }

    try {
      setLoading(true);
      setError(null);
      
      const content = await file.text();
      const result = await wails.importOpenAPISpecFromURL('data:application/json;base64,' + btoa(content));
      
      if (result.valid) {
        onImport(result);
        onClose();
      } else {
        setError(result.errors?.join(', ') || 'Failed to import specification');
      }
    } catch (err: any) {
      setError(err.message || 'Failed to import dropped file');
    } finally {
      setLoading(false);
    }
  };

  // Handle close
  const handleClose = () => {
    if (loading) return;
    setError(null);
    setUrl('');
    onClose();
  };

  // Handle key press
  const handleKeyPress = (e: React.KeyboardEvent) => {
    if (e.key === 'Enter' && activeTab === 'url' && !loading) {
      handleUrlImport();
    }
  };

  if (!isOpen) return null;

  return (
    <div className={`file-import-dialog ${className}`}>
      <div className="file-import-dialog__overlay" onClick={handleClose} />
      
      <div className="file-import-dialog__content">
        <div className="file-import-dialog__header">
          <h2>{title}</h2>
          <button 
            className="file-import-dialog__close"
            onClick={handleClose}
            disabled={loading}
          >
            <svg viewBox="0 0 20 20" fill="currentColor">
              <path fillRule="evenodd" d="M4.293 4.293a1 1 0 011.414 0L10 8.586l4.293-4.293a1 1 0 111.414 1.414L11.414 10l4.293 4.293a1 1 0 01-1.414 1.414L10 11.414l-4.293 4.293a1 1 0 01-1.414-1.414L8.586 10 4.293 5.707a1 1 0 010-1.414z" clipRule="evenodd" />
            </svg>
          </button>
        </div>

        <div className="file-import-dialog__tabs">
          <button
            className={`file-import-dialog__tab ${activeTab === 'file' ? 'file-import-dialog__tab--active' : ''}`}
            onClick={() => setActiveTab('file')}
            disabled={loading}
          >
            <svg viewBox="0 0 20 20" fill="currentColor">
              <path fillRule="evenodd" d="M4 4a2 2 0 00-2 2v8a2 2 0 002 2h12a2 2 0 002-2V6a2 2 0 00-2-2h-5L9 2H4z" clipRule="evenodd" />
            </svg>
            From File
          </button>
          <button
            className={`file-import-dialog__tab ${activeTab === 'url' ? 'file-import-dialog__tab--active' : ''}`}
            onClick={() => setActiveTab('url')}
            disabled={loading}
          >
            <svg viewBox="0 0 20 20" fill="currentColor">
              <path fillRule="evenodd" d="M12.586 4.586a2 2 0 112.828 2.828l-3 3a2 2 0 01-2.828 0 1 1 0 00-1.414 1.414 4 4 0 005.656 0l3-3a4 4 0 00-5.656-5.656l-1.5 1.5a1 1 0 101.414 1.414l1.5-1.5zm-5 5a2 2 0 012.828 0 1 1 0 101.414-1.414 4 4 0 00-5.656 0l-3 3a4 4 0 105.656 5.656l1.5-1.5a1 1 0 10-1.414-1.414l-1.5 1.5a2 2 0 11-2.828-2.828l3-3z" clipRule="evenodd" />
            </svg>
            From URL
          </button>
        </div>

        {/* Error Display */}
        {error && (
          <div className="file-import-dialog__error">
            <svg viewBox="0 0 20 20" fill="currentColor">
              <path fillRule="evenodd" d="M18 10a8 8 0 11-16 0 8 8 0 0116 0zm-7 4a1 1 0 11-2 0 1 1 0 012 0zm-1-9a1 1 0 00-1 1v4a1 1 0 102 0V6a1 1 0 00-1-1z" clipRule="evenodd" />
            </svg>
            {error}
          </div>
        )}

        <div className="file-import-dialog__body">
          {activeTab === 'file' ? (
            <div className="file-import-dialog__file-section">
              <div 
                className={`file-import-dialog__drop-zone ${dragOver ? 'file-import-dialog__drop-zone--drag-over' : ''}`}
                onDragOver={handleDragOver}
                onDragLeave={handleDragLeave}
                onDrop={handleDrop}
              >
                <div className="file-import-dialog__drop-zone-content">
                  <svg viewBox="0 0 20 20" fill="currentColor">
                    <path fillRule="evenodd" d="M4 4a2 2 0 00-2 2v8a2 2 0 002 2h12a2 2 0 002-2V6a2 2 0 00-2-2h-5L9 2H4z" clipRule="evenodd" />
                  </svg>
                  <h3>Drop OpenAPI file here</h3>
                  <p>or click to browse files</p>
                  <div className="file-import-dialog__supported-formats">
                    Supported formats: JSON, YAML, YML
                  </div>
                </div>
              </div>
              
              <div className="file-import-dialog__actions">
                <button 
                  className="file-import-dialog__browse-button"
                  onClick={handleFileSelect}
                  disabled={loading}
                >
                  {loading ? (
                    <>
                      <div className="file-import-dialog__spinner" />
                      Importing...
                    </>
                  ) : (
                    <>
                      <svg viewBox="0 0 20 20" fill="currentColor">
                        <path fillRule="evenodd" d="M3 17a1 1 0 011-1h12a1 1 0 110 2H4a1 1 0 01-1-1zm3.293-7.707a1 1 0 011.414 0L9 10.586V3a1 1 0 112 0v7.586l1.293-1.293a1 1 0 111.414 1.414l-3 3a1 1 0 01-1.414 0l-3-3a1 1 0 010-1.414z" clipRule="evenodd" />
                      </svg>
                      Browse Files
                    </>
                  )}
                </button>
              </div>
            </div>
          ) : (
            <div className="file-import-dialog__url-section">
              <div className="file-import-dialog__url-input-group">
                <label htmlFor="spec-url">OpenAPI Specification URL</label>
                <input
                  id="spec-url"
                  type="url"
                  value={url}
                  onChange={(e) => setUrl(e.target.value)}
                  onKeyPress={handleKeyPress}
                  placeholder="https://api.example.com/openapi.json"
                  disabled={loading}
                  className="file-import-dialog__url-input"
                />
                <div className="file-import-dialog__url-examples">
                  <span>Examples:</span>
                  <button 
                    type="button"
                    onClick={() => setUrl('https://petstore.swagger.io/v2/swagger.json')}
                    disabled={loading}
                  >
                    Petstore API
                  </button>
                  <button 
                    type="button"
                    onClick={() => setUrl('https://raw.githubusercontent.com/OAI/OpenAPI-Specification/main/examples/v3.0/petstore.yaml')}
                    disabled={loading}
                  >
                    Petstore V3
                  </button>
                </div>
              </div>
              
              <div className="file-import-dialog__actions">
                <button 
                  className="file-import-dialog__import-button"
                  onClick={handleUrlImport}
                  disabled={loading || !url.trim()}
                >
                  {loading ? (
                    <>
                      <div className="file-import-dialog__spinner" />
                      Importing...
                    </>
                  ) : (
                    <>
                      <svg viewBox="0 0 20 20" fill="currentColor">
                        <path fillRule="evenodd" d="M3 17a1 1 0 011-1h12a1 1 0 110 2H4a1 1 0 01-1-1zm3.293-7.707a1 1 0 011.414 0L9 10.586V3a1 1 0 112 0v7.586l1.293-1.293a1 1 0 111.414 1.414l-3 3a1 1 0 01-1.414 0l-3-3a1 1 0 010-1.414z" clipRule="evenodd" />
                      </svg>
                      Import from URL
                    </>
                  )}
                </button>
              </div>
            </div>
          )}
        </div>
      </div>
    </div>
  );
};

export default FileImportDialog;