import React, { useState } from 'react';
import { wails } from '../../services/wails';
import { ExportResult } from '../../types';
import './FileExportDialog.scss';

export interface FileExportDialogProps {
  isOpen: boolean;
  onClose: () => void;
  onExport: (result: ExportResult) => void;
  projectId: string;
  projectName: string;
  className?: string;
}

export const FileExportDialog: React.FC<FileExportDialogProps> = ({
  isOpen,
  onClose,
  onExport,
  projectId,
  projectName,
  className = ''
}) => {
  const [targetDir, setTargetDir] = useState('');
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [exportResult, setExportResult] = useState<ExportResult | null>(null);

  // Handle directory selection
  const handleSelectDirectory = async () => {
    try {
      setError(null);
      const dirPath = await wails.selectDirectory('Select Export Directory');
      
      if (dirPath) {
        setTargetDir(dirPath);
      }
    } catch (err: any) {
      setError(err.message || 'Failed to select directory');
    }
  };

  // Handle export
  const handleExport = async () => {
    if (!targetDir) {
      setError('Please select a target directory');
      return;
    }

    try {
      setLoading(true);
      setError(null);
      
      const result = await wails.exportGeneratedServer(projectId, targetDir);
      
      setExportResult(result);
      onExport(result);
    } catch (err: any) {
      setError(err.message || 'Failed to export server');
    } finally {
      setLoading(false);
    }
  };

  // Handle close
  const handleClose = () => {
    if (loading) return;
    setError(null);
    setTargetDir('');
    setExportResult(null);
    onClose();
  };

  // Format file size
  const formatFileSize = (bytes: number): string => {
    if (bytes < 1024) return `${bytes} B`;
    if (bytes < 1024 * 1024) return `${(bytes / 1024).toFixed(1)} KB`;
    return `${(bytes / (1024 * 1024)).toFixed(1)} MB`;
  };

  if (!isOpen) return null;

  return (
    <div className={`file-export-dialog ${className}`}>
      <div className="file-export-dialog__overlay" onClick={handleClose} />
      
      <div className="file-export-dialog__content">
        <div className="file-export-dialog__header">
          <h2>Export Generated Server</h2>
          <button 
            className="file-export-dialog__close"
            onClick={handleClose}
            disabled={loading}
          >
            <svg viewBox="0 0 20 20" fill="currentColor">
              <path fillRule="evenodd" d="M4.293 4.293a1 1 0 011.414 0L10 8.586l4.293-4.293a1 1 0 111.414 1.414L11.414 10l4.293 4.293a1 1 0 01-1.414 1.414L10 11.414l-4.293 4.293a1 1 0 01-1.414-1.414L8.586 10 4.293 5.707a1 1 0 010-1.414z" clipRule="evenodd" />
            </svg>
          </button>
        </div>

        {/* Error Display */}
        {error && (
          <div className="file-export-dialog__error">
            <svg viewBox="0 0 20 20" fill="currentColor">
              <path fillRule="evenodd" d="M18 10a8 8 0 11-16 0 8 8 0 0116 0zm-7 4a1 1 0 11-2 0 1 1 0 012 0zm-1-9a1 1 0 00-1 1v4a1 1 0 102 0V6a1 1 0 00-1-1z" clipRule="evenodd" />
            </svg>
            {error}
          </div>
        )}

        <div className="file-export-dialog__body">
          {!exportResult ? (
            <div className="file-export-dialog__setup">
              <div className="file-export-dialog__project-info">
                <div className="file-export-dialog__project-icon">
                  <svg viewBox="0 0 20 20" fill="currentColor">
                    <path fillRule="evenodd" d="M4 4a2 2 0 00-2 2v8a2 2 0 002 2h12a2 2 0 002-2V6a2 2 0 00-2-2h-5L9 2H4z" clipRule="evenodd" />
                  </svg>
                </div>
                <div className="file-export-dialog__project-details">
                  <h3>{projectName}</h3>
                  <p>Export the generated MCP server to a directory of your choice</p>
                </div>
              </div>

              <div className="file-export-dialog__directory-selection">
                <label htmlFor="target-dir">Export Directory</label>
                <div className="file-export-dialog__directory-input-group">
                  <input
                    id="target-dir"
                    type="text"
                    value={targetDir}
                    onChange={(e) => setTargetDir(e.target.value)}
                    placeholder="Select a directory to export to..."
                    className="file-export-dialog__directory-input"
                    disabled={loading}
                  />
                  <button
                    type="button"
                    onClick={handleSelectDirectory}
                    disabled={loading}
                    className="file-export-dialog__browse-button"
                  >
                    <svg viewBox="0 0 20 20" fill="currentColor">
                      <path fillRule="evenodd" d="M4 4a2 2 0 00-2 2v8a2 2 0 002 2h12a2 2 0 002-2V6a2 2 0 00-2-2h-5L9 2H4z" clipRule="evenodd" />
                    </svg>
                    Browse
                  </button>
                </div>
              </div>

              <div className="file-export-dialog__export-info">
                <h4>What will be exported:</h4>
                <ul>
                  <li>
                    <svg viewBox="0 0 20 20" fill="currentColor">
                      <path fillRule="evenodd" d="M16.707 5.293a1 1 0 010 1.414l-8 8a1 1 0 01-1.414 0l-4-4a1 1 0 011.414-1.414L8 12.586l7.293-7.293a1 1 0 011.414 0z" clipRule="evenodd" />
                    </svg>
                    <code>main.go</code> - MCP server implementation
                  </li>
                  <li>
                    <svg viewBox="0 0 20 20" fill="currentColor">
                      <path fillRule="evenodd" d="M16.707 5.293a1 1 0 010 1.414l-8 8a1 1 0 01-1.414 0l-4-4a1 1 0 011.414-1.414L8 12.586l7.293-7.293a1 1 0 011.414 0z" clipRule="evenodd" />
                    </svg>
                    <code>go.mod</code> - Go module definition
                  </li>
                  <li>
                    <svg viewBox="0 0 20 20" fill="currentColor">
                      <path fillRule="evenodd" d="M16.707 5.293a1 1 0 010 1.414l-8 8a1 1 0 01-1.414 0l-4-4a1 1 0 011.414-1.414L8 12.586l7.293-7.293a1 1 0 011.414 0z" clipRule="evenodd" />
                    </svg>
                    <code>README.md</code> - Documentation
                  </li>
                  <li>
                    <svg viewBox="0 0 20 20" fill="currentColor">
                      <path fillRule="evenodd" d="M16.707 5.293a1 1 0 010 1.414l-8 8a1 1 0 01-1.414 0l-4-4a1 1 0 011.414-1.414L8 12.586l7.293-7.293a1 1 0 011.414 0z" clipRule="evenodd" />
                    </svg>
                    <code>Dockerfile</code> - Container configuration
                  </li>
                  <li>
                    <svg viewBox="0 0 20 20" fill="currentColor">
                      <path fillRule="evenodd" d="M16.707 5.293a1 1 0 010 1.414l-8 8a1 1 0 01-1.414 0l-4-4a1 1 0 011.414-1.414L8 12.586l7.293-7.293a1 1 0 011.414 0z" clipRule="evenodd" />
                    </svg>
                    <code>Makefile</code> - Build automation
                  </li>
                  <li>
                    <svg viewBox="0 0 20 20" fill="currentColor">
                      <path fillRule="evenodd" d="M16.707 5.293a1 1 0 010 1.414l-8 8a1 1 0 01-1.414 0l-4-4a1 1 0 011.414-1.414L8 12.586l7.293-7.293a1 1 0 011.414 0z" clipRule="evenodd" />
                    </svg>
                    <code>.gitignore</code> - Git ignore rules
                  </li>
                </ul>
              </div>

              <div className="file-export-dialog__actions">
                <button
                  type="button"
                  onClick={handleClose}
                  disabled={loading}
                  className="file-export-dialog__cancel-button"
                >
                  Cancel
                </button>
                <button
                  type="button"
                  onClick={handleExport}
                  disabled={loading || !targetDir}
                  className="file-export-dialog__export-button"
                >
                  {loading ? (
                    <>
                      <div className="file-export-dialog__spinner" />
                      Exporting...
                    </>
                  ) : (
                    <>
                      <svg viewBox="0 0 20 20" fill="currentColor">
                        <path fillRule="evenodd" d="M3 17a1 1 0 011-1h12a1 1 0 110 2H4a1 1 0 01-1-1zM6.293 6.707a1 1 0 010-1.414l3-3a1 1 0 011.414 0l3 3a1 1 0 01-1.414 1.414L11 5.414V13a1 1 0 11-2 0V5.414L7.707 6.707a1 1 0 01-1.414 0z" clipRule="evenodd" />
                      </svg>
                      Export Server
                    </>
                  )}
                </button>
              </div>
            </div>
          ) : (
            <div className="file-export-dialog__result">
              <div className="file-export-dialog__success-icon">
                <svg viewBox="0 0 20 20" fill="currentColor">
                  <path fillRule="evenodd" d="M16.707 5.293a1 1 0 010 1.414l-8 8a1 1 0 01-1.414 0l-4-4a1 1 0 011.414-1.414L8 12.586l7.293-7.293a1 1 0 011.414 0z" clipRule="evenodd" />
                </svg>
              </div>
              
              <h3>Export Successful!</h3>
              <p>Your MCP server has been exported successfully.</p>
              
              <div className="file-export-dialog__export-details">
                <div className="file-export-dialog__export-stat">
                  <span className="file-export-dialog__stat-label">Location:</span>
                  <span className="file-export-dialog__stat-value">{exportResult.targetDir}</span>
                </div>
                <div className="file-export-dialog__export-stat">
                  <span className="file-export-dialog__stat-label">Files:</span>
                  <span className="file-export-dialog__stat-value">{exportResult.totalFiles}</span>
                </div>
                <div className="file-export-dialog__export-stat">
                  <span className="file-export-dialog__stat-label">Size:</span>
                  <span className="file-export-dialog__stat-value">{formatFileSize(exportResult.totalSize)}</span>
                </div>
              </div>
              
              <div className="file-export-dialog__exported-files">
                <h4>Exported Files:</h4>
                <ul>
                  {exportResult.exportedFiles.map((file, index) => (
                    <li key={index}>
                      <svg viewBox="0 0 20 20" fill="currentColor">
                        <path fillRule="evenodd" d="M4 4a2 2 0 00-2 2v8a2 2 0 002 2h12a2 2 0 002-2V6a2 2 0 00-2-2h-5L9 2H4z" clipRule="evenodd" />
                      </svg>
                      <span className="file-export-dialog__file-name">{file.name}</span>
                      <span className="file-export-dialog__file-size">{formatFileSize(file.size)}</span>
                    </li>
                  ))}
                </ul>
              </div>
              
              <div className="file-export-dialog__actions">
                <button
                  type="button"
                  onClick={handleClose}
                  className="file-export-dialog__done-button"
                >
                  Done
                </button>
              </div>
            </div>
          )}
        </div>
      </div>
    </div>
  );
};

export default FileExportDialog;