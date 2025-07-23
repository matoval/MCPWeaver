import React, { useState, useRef, useCallback, useEffect } from 'react';
import { 
  Layout,
  PanelLeft,
  PanelRight,
  Split,
  GitCompare,
  Save,
  Download,
  Settings,
  Eye,
  Code,
  FileText,
  Terminal,
  AlertCircle,
  CheckCircle,
  Clock,
  Maximize2,
  Minimize2
} from 'lucide-react';
import CodeEditor, { CodeEditorHandle } from './CodeEditor';
import FileTree, { FileTreeNode } from './FileTree';
import DiffViewer from './DiffViewer';
import { wails } from '../../services/wails';
import './CodeEditorDashboard.scss';

export interface CodeEditorDashboardProps {
  projectId: string;
  generatedFiles?: FileTreeNode[];
  onFileChange?: (file: FileTreeNode, content: string) => void;
  onFileSave?: (file: FileTreeNode, content: string) => void;
  onFileValidate?: (file: FileTreeNode, isValid: boolean, errors: string[]) => void;
  className?: string;
}

interface EditorTab {
  id: string;
  file: FileTreeNode;
  content: string;
  originalContent: string;
  isDirty: boolean;
  isValidating?: boolean;
  validationErrors?: string[];
}

interface ValidationResult {
  isValid: boolean;
  errors: string[];
  warnings: string[];
}

const CodeEditorDashboard: React.FC<CodeEditorDashboardProps> = ({
  projectId,
  generatedFiles = [],
  onFileChange,
  onFileSave,
  onFileValidate,
  className = ''
}) => {
  const [files, setFiles] = useState<FileTreeNode[]>(generatedFiles);
  const [selectedFile, setSelectedFile] = useState<string | null>(null);
  const [tabs, setTabs] = useState<EditorTab[]>([]);
  const [activeTab, setActiveTab] = useState<string | null>(null);
  const [showFileTree, setShowFileTree] = useState(true);
  const [showDiff, setShowDiff] = useState(false);
  const [layout, setLayout] = useState<'horizontal' | 'vertical'>('horizontal');
  const [isFullscreen, setIsFullscreen] = useState(false);
  const [searchQuery, setSearchQuery] = useState('');
  const [validationResults, setValidationResults] = useState<Map<string, ValidationResult>>(new Map());
  
  const editorRef = useRef<CodeEditorHandle>(null);
  const dashboardRef = useRef<HTMLDivElement>(null);

  // Load generated files
  useEffect(() => {
    const loadGeneratedFiles = async () => {
      try {
        // In a real implementation, this would load the generated files from the backend
        // For now, we'll use the provided generatedFiles
        setFiles(generatedFiles);
      } catch (error) {
        console.error('Failed to load generated files:', error);
      }
    };

    if (projectId) {
      loadGeneratedFiles();
    }
  }, [projectId, generatedFiles]);

  const findFileById = useCallback((files: FileTreeNode[], id: string): FileTreeNode | null => {
    for (const file of files) {
      if (file.id === id) return file;
      if (file.children) {
        const found = findFileById(file.children, id);
        if (found) return found;
      }
    }
    return null;
  }, []);

  const openFile = useCallback(async (file: FileTreeNode) => {
    if (file.type === 'directory') return;

    // Check if file is already open
    const existingTab = tabs.find(tab => tab.id === file.id);
    if (existingTab) {
      setActiveTab(file.id);
      setSelectedFile(file.id);
      return;
    }

    try {
      // Load file content (in a real implementation, this would fetch from backend)
      const content = file.content || `// Generated file: ${file.name}\n// Content would be loaded from the server\n`;
      
      const newTab: EditorTab = {
        id: file.id,
        file,
        content,
        originalContent: content,
        isDirty: false
      };

      setTabs(prev => [...prev, newTab]);
      setActiveTab(file.id);
      setSelectedFile(file.id);
    } catch (error) {
      console.error('Failed to load file:', error);
    }
  }, [tabs]);

  const closeTab = useCallback((tabId: string) => {
    const tab = tabs.find(t => t.id === tabId);
    if (tab?.isDirty) {
      const shouldClose = window.confirm('This file has unsaved changes. Do you want to close it anyway?');
      if (!shouldClose) return;
    }

    setTabs(prev => prev.filter(t => t.id !== tabId));
    
    if (activeTab === tabId) {
      const remainingTabs = tabs.filter(t => t.id !== tabId);
      setActiveTab(remainingTabs.length > 0 ? remainingTabs[remainingTabs.length - 1].id : null);
      setSelectedFile(remainingTabs.length > 0 ? remainingTabs[remainingTabs.length - 1].id : null);
    }
  }, [tabs, activeTab]);

  const handleContentChange = useCallback((content: string | undefined) => {
    if (!activeTab || !content) return;

    setTabs(prev => prev.map(tab => {
      if (tab.id === activeTab) {
        const updatedTab = {
          ...tab,
          content,
          isDirty: content !== tab.originalContent
        };
        
        if (onFileChange) {
          onFileChange(tab.file, content);
        }
        
        return updatedTab;
      }
      return tab;
    }));
  }, [activeTab, onFileChange]);

  const handleSave = useCallback(async (content?: string) => {
    if (!activeTab) return;

    const tab = tabs.find(t => t.id === activeTab);
    if (!tab) return;

    const contentToSave = content || tab.content;

    try {
      // In a real implementation, this would save to the backend
      if (onFileSave) {
        onFileSave(tab.file, contentToSave);
      }

      setTabs(prev => prev.map(t => {
        if (t.id === activeTab) {
          return {
            ...t,
            originalContent: contentToSave,
            isDirty: false
          };
        }
        return t;
      }));
    } catch (error) {
      console.error('Failed to save file:', error);
    }
  }, [activeTab, tabs, onFileSave]);

  const validateFile = useCallback(async (file: FileTreeNode, content: string): Promise<ValidationResult> => {
    try {
      // In a real implementation, this would validate using the backend
      // For Go files, we might check syntax, imports, etc.
      const errors: string[] = [];
      const warnings: string[] = [];

      // Simple validation example
      if (file.name.endsWith('.go')) {
        if (!content.includes('package ')) {
          errors.push('Missing package declaration');
        }
        if (content.includes('func main()') && !content.includes('package main')) {
          errors.push('main function requires package main');
        }
      }

      const result = {
        isValid: errors.length === 0,
        errors,
        warnings
      };

      setValidationResults(prev => new Map(prev).set(file.id, result));
      
      if (onFileValidate) {
        onFileValidate(file, result.isValid, result.errors);
      }

      return result;
    } catch (error) {
      console.error('Validation failed:', error);
      return { isValid: false, errors: ['Validation failed'], warnings: [] };
    }
  }, [onFileValidate]);

  const formatCode = useCallback(() => {
    if (editorRef.current) {
      editorRef.current.format();
    }
  }, []);

  const toggleDiff = useCallback(() => {
    setShowDiff(!showDiff);
  }, [showDiff]);

  const toggleLayout = useCallback(() => {
    setLayout(prev => prev === 'horizontal' ? 'vertical' : 'horizontal');
  }, []);

  const toggleFullscreen = useCallback(() => {
    setIsFullscreen(!isFullscreen);
  }, []);

  const activeTabData = tabs.find(tab => tab.id === activeTab);
  const validationResult = activeTabData ? validationResults.get(activeTabData.id) : null;

  return (
    <div 
      ref={dashboardRef}
      className={`code-editor-dashboard ${isFullscreen ? 'code-editor-dashboard--fullscreen' : ''} ${className}`}
    >
      <div className="code-editor-dashboard__header">
        <div className="code-editor-dashboard__title">
          <Code size={20} />
          <h2>Code Editor</h2>
          <span className="code-editor-dashboard__project-id">Project: {projectId}</span>
        </div>

        <div className="code-editor-dashboard__toolbar">
          <button
            className={`code-editor-dashboard__tool ${showFileTree ? 'active' : ''}`}
            onClick={() => setShowFileTree(!showFileTree)}
            title="Toggle File Explorer"
          >
            <PanelLeft size={16} />
          </button>

          <button
            className={`code-editor-dashboard__tool ${showDiff ? 'active' : ''}`}
            onClick={toggleDiff}
            title="Toggle Diff View"
            disabled={!activeTabData}
          >
            <GitCompare size={16} />
          </button>

          <button
            className="code-editor-dashboard__tool"
            onClick={toggleLayout}
            title="Toggle Layout"
          >
            {layout === 'horizontal' ? <Split size={16} /> : <Layout size={16} />}
          </button>

          <div className="code-editor-dashboard__separator" />

          <button
            className="code-editor-dashboard__tool"
            onClick={formatCode}
            title="Format Code"
            disabled={!activeTabData}
          >
            <FileText size={16} />
          </button>

          <button
            className="code-editor-dashboard__tool"
            onClick={() => handleSave()}
            title="Save File"
            disabled={!activeTabData?.isDirty}
          >
            <Save size={16} />
          </button>

          <div className="code-editor-dashboard__separator" />

          <button
            className="code-editor-dashboard__tool"
            onClick={toggleFullscreen}
            title="Toggle Fullscreen"
          >
            {isFullscreen ? <Minimize2 size={16} /> : <Maximize2 size={16} />}
          </button>
        </div>
      </div>

      <div className="code-editor-dashboard__content">
        {showFileTree && (
          <div className="code-editor-dashboard__sidebar">
            <FileTree
              files={files}
              onFileSelect={openFile}
              selectedFile={selectedFile}
              searchQuery={searchQuery}
              onSearchChange={setSearchQuery}
              showSearch={true}
              allowEdit={false}
            />
          </div>
        )}

        <div className="code-editor-dashboard__main">
          {tabs.length > 0 && (
            <div className="code-editor-dashboard__tabs">
              {tabs.map(tab => (
                <div
                  key={tab.id}
                  className={`code-editor-dashboard__tab ${activeTab === tab.id ? 'active' : ''}`}
                  onClick={() => setActiveTab(tab.id)}
                >
                  <div className="code-editor-dashboard__tab-content">
                    <FileText size={14} />
                    <span className="code-editor-dashboard__tab-name">{tab.file.name}</span>
                    {tab.isDirty && <div className="code-editor-dashboard__tab-dirty" />}
                    {validationResults.get(tab.id) && (
                      <div className={`code-editor-dashboard__tab-status ${
                        validationResults.get(tab.id)?.isValid ? 'valid' : 'invalid'
                      }`}>
                        {validationResults.get(tab.id)?.isValid ? 
                          <CheckCircle size={12} /> : 
                          <AlertCircle size={12} />
                        }
                      </div>
                    )}
                  </div>
                  <button
                    className="code-editor-dashboard__tab-close"
                    onClick={(e) => {
                      e.stopPropagation();
                      closeTab(tab.id);
                    }}
                  >
                    ×
                  </button>
                </div>
              ))}
            </div>
          )}

          <div className={`code-editor-dashboard__editors ${layout}`}>
            {activeTabData ? (
              <div className="code-editor-dashboard__editor-container">
                <CodeEditor
                  ref={editorRef}
                  value={activeTabData.content}
                  onChange={handleContentChange}
                  language={activeTabData.file.language || 'go'}
                  filename={activeTabData.file.name}
                  onSave={handleSave}
                  onFormat={formatCode}
                  onValidate={(markers) => {
                    const errors = markers
                      .filter(m => m.severity === 8) // Error severity
                      .map(m => m.message);
                    validateFile(activeTabData.file, activeTabData.content);
                  }}
                  className="code-editor-dashboard__editor"
                />

                {validationResult && !validationResult.isValid && (
                  <div className="code-editor-dashboard__validation">
                    <div className="code-editor-dashboard__validation-header">
                      <AlertCircle size={16} />
                      <span>Validation Errors</span>
                    </div>
                    <div className="code-editor-dashboard__validation-errors">
                      {validationResult.errors.map((error, index) => (
                        <div key={index} className="code-editor-dashboard__validation-error">
                          {error}
                        </div>
                      ))}
                    </div>
                  </div>
                )}
              </div>
            ) : (
              <div className="code-editor-dashboard__empty">
                <Eye size={48} />
                <h3>No File Selected</h3>
                <p>Select a file from the explorer to start editing</p>
              </div>
            )}

            {showDiff && activeTabData && (
              <div className="code-editor-dashboard__diff-container">
                <DiffViewer
                  originalContent={activeTabData.originalContent}
                  modifiedContent={activeTabData.content}
                  originalTitle="Original"
                  modifiedTitle="Current"
                  language={activeTabData.file.language || 'go'}
                  onContentChange={handleContentChange}
                  onAcceptChange={() => handleSave()}
                  onRejectChange={() => {
                    handleContentChange(activeTabData.originalContent);
                  }}
                  readOnly={false}
                  metadata={{
                    modifiedDate: new Date(),
                    author: 'User',
                    message: 'Local changes'
                  }}
                />
              </div>
            )}
          </div>
        </div>
      </div>

      <div className="code-editor-dashboard__status">
        <div className="code-editor-dashboard__status-left">
          {activeTabData && (
            <>
              <span>{activeTabData.file.language || 'text'}</span>
              <span>•</span>
              <span>{activeTabData.content.split('\n').length} lines</span>
              {activeTabData.isDirty && (
                <>
                  <span>•</span>
                  <span className="code-editor-dashboard__status-dirty">Unsaved changes</span>
                </>
              )}
            </>
          )}
        </div>

        <div className="code-editor-dashboard__status-right">
          {validationResult && (
            <div className={`code-editor-dashboard__status-validation ${
              validationResult.isValid ? 'valid' : 'invalid'
            }`}>
              {validationResult.isValid ? (
                <><CheckCircle size={14} /> Valid</>
              ) : (
                <><AlertCircle size={14} /> {validationResult.errors.length} error(s)</>
              )}
            </div>
          )}
          <span>{tabs.length} file(s) open</span>
        </div>
      </div>
    </div>
  );
};

export default CodeEditorDashboard;