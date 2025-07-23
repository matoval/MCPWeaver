import React, { useRef, useEffect, useState, useCallback } from 'react';
import { DiffEditor } from '@monaco-editor/react';
import type * as monaco from 'monaco-editor';
import { 
  GitCommit,
  GitBranch,
  Clock,
  Eye,
  EyeOff,
  Download,
  Copy,
  RotateCcw,
  Check,
  X
} from 'lucide-react';
import './DiffViewer.scss';

export interface DiffViewerProps {
  originalContent: string;
  modifiedContent: string;
  originalTitle?: string;
  modifiedTitle?: string;
  language?: string;
  theme?: 'vs-dark' | 'light' | 'vs';
  readOnly?: boolean;
  onContentChange?: (content: string) => void;
  onAcceptChange?: () => void;
  onRejectChange?: () => void;
  showSideBySide?: boolean;
  showMetadata?: boolean;
  metadata?: {
    originalDate?: Date;
    modifiedDate?: Date;
    author?: string;
    message?: string;
  };
  className?: string;
}

export interface DiffStats {
  additions: number;
  deletions: number;
  modifications: number;
}

const DiffViewer: React.FC<DiffViewerProps> = ({
  originalContent,
  modifiedContent,
  originalTitle = 'Original',
  modifiedTitle = 'Modified',
  language = 'go',
  theme = 'vs-dark',
  readOnly = true,
  onContentChange,
  onAcceptChange,
  onRejectChange,
  showSideBySide = true,
  showMetadata = true,
  metadata,
  className = ''
}) => {
  const diffEditorRef = useRef<monaco.editor.IStandaloneDiffEditor | null>(null);
  const [diffStats, setDiffStats] = useState<DiffStats>({ additions: 0, deletions: 0, modifications: 0 });
  const [showWhitespace, setShowWhitespace] = useState(false);
  const [isInlineView, setIsInlineView] = useState(!showSideBySide);

  const calculateDiffStats = useCallback((original: string, modified: string): DiffStats => {
    const originalLines = original.split('\n');
    const modifiedLines = modified.split('\n');
    
    let additions = 0;
    let deletions = 0;
    let modifications = 0;

    // Simple diff calculation (in a real implementation, you'd use a proper diff algorithm)
    const maxLines = Math.max(originalLines.length, modifiedLines.length);
    
    for (let i = 0; i < maxLines; i++) {
      const origLine = originalLines[i] || '';
      const modLine = modifiedLines[i] || '';
      
      if (origLine === '' && modLine !== '') {
        additions++;
      } else if (origLine !== '' && modLine === '') {
        deletions++;
      } else if (origLine !== modLine) {
        modifications++;
      }
    }

    return { additions, deletions, modifications };
  }, []);

  useEffect(() => {
    const stats = calculateDiffStats(originalContent, modifiedContent);
    setDiffStats(stats);
  }, [originalContent, modifiedContent, calculateDiffStats]);

  const handleEditorDidMount = useCallback((editor: monaco.editor.IStandaloneDiffEditor) => {
    diffEditorRef.current = editor;

    // Configure diff editor options
    editor.updateOptions({
      renderSideBySide: !isInlineView,
      renderMarginRevertIcon: !readOnly,
      enableSplitViewResizing: true,
      renderIndicators: true,
      originalEditable: false,
      readOnly: readOnly
    });

    // Listen for content changes
    if (onContentChange && !readOnly) {
      const modifiedEditor = editor.getModifiedEditor();
      const disposable = modifiedEditor.onDidChangeModelContent(() => {
        const content = modifiedEditor.getValue();
        onContentChange(content);
      });

      return () => disposable.dispose();
    }
  }, [isInlineView, readOnly, onContentChange]);

  const handleDownload = useCallback(() => {
    const content = diffEditorRef.current?.getModifiedEditor().getValue() || modifiedContent;
    const blob = new Blob([content], { type: 'text/plain' });
    const url = URL.createObjectURL(blob);
    const a = document.createElement('a');
    a.href = url;
    a.download = `${modifiedTitle.replace(/\s+/g, '_')}.${language}`;
    document.body.appendChild(a);
    a.click();
    document.body.removeChild(a);
    URL.revokeObjectURL(url);
  }, [modifiedContent, modifiedTitle, language]);

  const handleCopyModified = useCallback(() => {
    const content = diffEditorRef.current?.getModifiedEditor().getValue() || modifiedContent;
    navigator.clipboard.writeText(content);
  }, [modifiedContent]);

  const handleCopyOriginal = useCallback(() => {
    const content = diffEditorRef.current?.getOriginalEditor().getValue() || originalContent;
    navigator.clipboard.writeText(content);
  }, [originalContent]);

  const handleRevert = useCallback(() => {
    if (diffEditorRef.current && !readOnly) {
      const modifiedEditor = diffEditorRef.current.getModifiedEditor();
      modifiedEditor.setValue(originalContent);
    }
  }, [originalContent, readOnly]);

  const toggleViewMode = useCallback(() => {
    setIsInlineView(!isInlineView);
    if (diffEditorRef.current) {
      diffEditorRef.current.updateOptions({
        renderSideBySide: isInlineView
      });
    }
  }, [isInlineView]);

  const toggleWhitespace = useCallback(() => {
    setShowWhitespace(!showWhitespace);
    if (diffEditorRef.current) {
      const originalEditor = diffEditorRef.current.getOriginalEditor();
      const modifiedEditor = diffEditorRef.current.getModifiedEditor();
      
      const renderWhitespace = showWhitespace ? 'none' : 'all';
      originalEditor.updateOptions({ renderWhitespace });
      modifiedEditor.updateOptions({ renderWhitespace });
    }
  }, [showWhitespace]);

  const formatDate = useCallback((date: Date): string => {
    return date.toLocaleString('en-US', {
      year: 'numeric',
      month: 'short',
      day: 'numeric',
      hour: '2-digit',
      minute: '2-digit'
    });
  }, []);

  return (
    <div className={`diff-viewer ${className}`}>
      <div className="diff-viewer__header">
        <div className="diff-viewer__title-section">
          <div className="diff-viewer__titles">
            <div className="diff-viewer__title diff-viewer__title--original">
              <GitBranch size={16} />
              <span>{originalTitle}</span>
            </div>
            <div className="diff-viewer__title diff-viewer__title--modified">
              <GitCommit size={16} />
              <span>{modifiedTitle}</span>
            </div>
          </div>
          
          {showMetadata && metadata && (
            <div className="diff-viewer__metadata">
              {metadata.author && (
                <span className="diff-viewer__author">by {metadata.author}</span>
              )}
              {metadata.modifiedDate && (
                <span className="diff-viewer__date">
                  <Clock size={14} />
                  {formatDate(metadata.modifiedDate)}
                </span>
              )}
              {metadata.message && (
                <span className="diff-viewer__message">{metadata.message}</span>
              )}
            </div>
          )}
        </div>

        <div className="diff-viewer__stats">
          <div className="diff-viewer__stat diff-viewer__stat--additions">
            +{diffStats.additions}
          </div>
          <div className="diff-viewer__stat diff-viewer__stat--deletions">
            -{diffStats.deletions}
          </div>
          <div className="diff-viewer__stat diff-viewer__stat--modifications">
            ~{diffStats.modifications}
          </div>
        </div>

        <div className="diff-viewer__actions">
          <button
            className="diff-viewer__action"
            onClick={toggleViewMode}
            title={isInlineView ? 'Switch to Side by Side' : 'Switch to Inline'}
          >
            {isInlineView ? 'Side by Side' : 'Inline'}
          </button>

          <button
            className="diff-viewer__action"
            onClick={toggleWhitespace}
            title={showWhitespace ? 'Hide Whitespace' : 'Show Whitespace'}
          >
            {showWhitespace ? <EyeOff size={16} /> : <Eye size={16} />}
          </button>

          <button
            className="diff-viewer__action"
            onClick={handleCopyOriginal}
            title="Copy Original"
          >
            <Copy size={16} />
            Original
          </button>

          <button
            className="diff-viewer__action"
            onClick={handleCopyModified}
            title="Copy Modified"
          >
            <Copy size={16} />
            Modified
          </button>

          <button
            className="diff-viewer__action"
            onClick={handleDownload}
            title="Download Modified"
          >
            <Download size={16} />
          </button>

          {!readOnly && (
            <>
              <button
                className="diff-viewer__action"
                onClick={handleRevert}
                title="Revert to Original"
              >
                <RotateCcw size={16} />
              </button>

              {onAcceptChange && (
                <button
                  className="diff-viewer__action diff-viewer__action--success"
                  onClick={onAcceptChange}
                  title="Accept Changes"
                >
                  <Check size={16} />
                  Accept
                </button>
              )}

              {onRejectChange && (
                <button
                  className="diff-viewer__action diff-viewer__action--danger"
                  onClick={onRejectChange}
                  title="Reject Changes"
                >
                  <X size={16} />
                  Reject
                </button>
              )}
            </>
          )}
        </div>
      </div>

      <div className="diff-viewer__editor">
        <DiffEditor
          height="100%"
          language={language}
          theme={theme}
          original={originalContent}
          modified={modifiedContent}
          onMount={handleEditorDidMount}
          options={{
            renderSideBySide: !isInlineView,
            renderMarginRevertIcon: !readOnly,
            enableSplitViewResizing: true,
            renderIndicators: true,
            originalEditable: false,
            readOnly: readOnly,
            automaticLayout: true,
            scrollBeyondLastLine: false,
            renderWhitespace: showWhitespace ? 'all' : 'none',
            ignoreTrimWhitespace: false,
            renderLineHighlight: 'all',
            scrollbar: {
              vertical: 'visible',
              horizontal: 'visible',
              useShadows: false,
              verticalHasArrows: false,
              horizontalHasArrows: false,
              verticalScrollbarSize: 10,
              horizontalScrollbarSize: 10
            },
            minimap: {
              enabled: true,
              side: 'right'
            },
            wordWrap: 'off',
            diffCodeLens: true,
            diffAlgorithm: 'advanced'
          }}
        />
      </div>

      {!readOnly && (
        <div className="diff-viewer__footer">
          <div className="diff-viewer__footer-stats">
            <span>
              {diffStats.additions + diffStats.deletions + diffStats.modifications} changes
            </span>
          </div>
          
          <div className="diff-viewer__footer-actions">
            <button
              className="diff-viewer__footer-button diff-viewer__footer-button--secondary"
              onClick={handleRevert}
            >
              Discard Changes
            </button>
            
            {onAcceptChange && (
              <button
                className="diff-viewer__footer-button diff-viewer__footer-button--primary"
                onClick={onAcceptChange}
              >
                Accept All Changes
              </button>
            )}
          </div>
        </div>
      )}
    </div>
  );
};

export default DiffViewer;