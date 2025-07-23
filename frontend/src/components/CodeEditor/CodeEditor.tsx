import React, { useRef, useCallback, useEffect, useState } from 'react';
import Editor from '@monaco-editor/react';
import type * as monaco from 'monaco-editor';
import { 
  Save, 
  Download, 
  Copy, 
  Search, 
  Settings, 
  FileText,
  Maximize2,
  Minimize2
} from 'lucide-react';
import './CodeEditor.scss';

export interface CodeEditorProps {
  value: string;
  onChange: (value: string | undefined) => void;
  language?: string;
  filename?: string;
  theme?: 'vs-dark' | 'light' | 'vs';
  readOnly?: boolean;
  onSave?: (content: string) => void;
  onFormat?: () => void;
  onValidate?: (markers: monaco.editor.IMarker[]) => void;
  className?: string;
  showMinimap?: boolean;
  wordWrap?: 'on' | 'off' | 'wordWrapColumn' | 'bounded';
  fontSize?: number;
  tabSize?: number;
  showLineNumbers?: boolean;
  automaticLayout?: boolean;
}

export interface CodeEditorHandle {
  format: () => void;
  getContent: () => string;
  setContent: (content: string) => void;
  focus: () => void;
  getSelection: () => string;
  insertText: (text: string) => void;
  findAndReplace: (searchValue: string, replaceValue: string) => void;
  goToLine: (lineNumber: number) => void;
}

const CodeEditor = React.forwardRef<CodeEditorHandle, CodeEditorProps>(({
  value,
  onChange,
  language = 'go',
  filename = 'main.go',
  theme = 'vs-dark',
  readOnly = false,
  onSave,
  onFormat,
  onValidate,
  className = '',
  showMinimap = true,
  wordWrap = 'on',
  fontSize = 14,
  tabSize = 4,
  showLineNumbers = true,
  automaticLayout = true
}, ref) => {
  const editorRef = useRef<monaco.editor.IStandaloneCodeEditor | null>(null);
  const [isFullscreen, setIsFullscreen] = useState(false);
  const [showSettings, setShowSettings] = useState(false);
  const [localSettings, setLocalSettings] = useState({
    fontSize,
    tabSize,
    showMinimap,
    wordWrap,
    showLineNumbers,
    theme
  });

  // Configure Monaco Editor for Go language support
  useEffect(() => {
    if (typeof window !== 'undefined') {
      // Get Monaco from window global
      const monacoInstance = (window as any).monaco;
      if (!monacoInstance) return;
      
      // Configure Go language
      monacoInstance.languages.register({ id: 'go' });
      
      // Set Go language configuration
      monacoInstance.languages.setLanguageConfiguration('go', {
        comments: {
          lineComment: '//',
          blockComment: ['/*', '*/']
        },
        brackets: [
          ['{', '}'],
          ['[', ']'],
          ['(', ')']
        ],
        autoClosingPairs: [
          { open: '{', close: '}' },
          { open: '[', close: ']' },
          { open: '(', close: ')' },
          { open: '"', close: '"' },
          { open: "'", close: "'" },
          { open: '`', close: '`' }
        ],
        surroundingPairs: [
          { open: '{', close: '}' },
          { open: '[', close: ']' },
          { open: '(', close: ')' },
          { open: '"', close: '"' },
          { open: "'", close: "'" },
          { open: '`', close: '`' }
        ],
        folding: {
          markers: {
            start: new RegExp('^\\s*//\\s*#?region\\b'),
            end: new RegExp('^\\s*//\\s*#?endregion\\b')
          }
        }
      });

      // Set Go syntax highlighting
      monacoInstance.languages.setMonarchTokensProvider('go', {
        keywords: [
          'break', 'case', 'chan', 'const', 'continue', 'default', 'defer',
          'else', 'fallthrough', 'for', 'func', 'go', 'goto', 'if', 'import',
          'interface', 'map', 'package', 'range', 'return', 'select', 'struct',
          'switch', 'type', 'var'
        ],
        builtins: [
          'bool', 'byte', 'complex64', 'complex128', 'error', 'float32',
          'float64', 'int', 'int8', 'int16', 'int32', 'int64', 'rune',
          'string', 'uint', 'uint8', 'uint16', 'uint32', 'uint64', 'uintptr',
          'append', 'cap', 'close', 'complex', 'copy', 'delete', 'imag',
          'len', 'make', 'new', 'panic', 'print', 'println', 'real',
          'recover'
        ],
        typeKeywords: ['any', 'comparable'],
        operators: [
          '+', '-', '*', '/', '%', '&', '|', '^', '<<', '>>', '&^',
          '+=', '-=', '*=', '/=', '%=', '&=', '|=', '^=', '<<=', '>>=', '&^=',
          '&&', '||', '<-', '++', '--', '==', '<', '>', '=', '!', '!=', '<=', '>=',
          ':=', '...', '(', ')', '[', ']', '{', '}', ',', ';', '.', ':'
        ],
        symbols: /[=><!~?:&|+\-*\/\^%]+/,
        escapes: /\\(?:[abfnrtv\\"']|x[0-9A-Fa-f]{1,4}|u[0-9A-Fa-f]{4}|U[0-9A-Fa-f]{8})/,
        tokenizer: {
          root: [
            [/[a-zA-Z_]\w*/, {
              cases: {
                '@keywords': 'keyword',
                '@builtins': 'type.identifier',
                '@typeKeywords': 'keyword.type',
                '@default': 'identifier'
              }
            }],
            [/[{}()\[\]]/, '@brackets'],
            [/[<>](?!@symbols)/, '@brackets'],
            [/@symbols/, {
              cases: {
                '@operators': 'operator',
                '@default': ''
              }
            }],
            [/\d*\.\d+([eE][\-+]?\d+)?/, 'number.float'],
            [/0[xX][0-9a-fA-F]+/, 'number.hex'],
            [/\d+/, 'number'],
            [/[;,.]/, 'delimiter'],
            [/"([^"\\]|\\.)*$/, 'string.invalid'],
            [/"/, { token: 'string.quote', bracket: '@open', next: '@string' }],
            [/`/, { token: 'string.quote', bracket: '@open', next: '@rawstring' }],
            [/'[^\\']'/, 'string'],
            [/(')(@escapes)(')/, ['string', 'string.escape', 'string']],
            [/'/, 'string.invalid']
          ],
          string: [
            [/[^\\"]+/, 'string'],
            [/@escapes/, 'string.escape'],
            [/\\./, 'string.escape.invalid'],
            [/"/, { token: 'string.quote', bracket: '@close', next: '@pop' }]
          ],
          rawstring: [
            [/[^`]*/, 'string'],
            [/`/, { token: 'string.quote', bracket: '@close', next: '@pop' }]
          ]
        }
      });
    }
  }, []);

  const handleEditorDidMount = useCallback((editor: monaco.editor.IStandaloneCodeEditor) => {
    editorRef.current = editor;

    // Add validation
    if (onValidate) {
      const model = editor.getModel();
      if (model) {
        const monacoInstance = (window as any).monaco;
        if (!monacoInstance) return;
        
        const disposable = monacoInstance.editor.onDidChangeMarkers((uris: any[]) => {
          const editorUri = model.uri;
          if (uris.find(uri => uri.toString() === editorUri.toString())) {
            const markers = monacoInstance.editor.getModelMarkers({ resource: editorUri });
            onValidate(markers);
          }
        });

        // Clean up on unmount
        return () => disposable.dispose();
      }
    }

    // Add keyboard shortcuts
    const monacoInstance = (window as any).monaco;
    if (monacoInstance) {
      editor.addCommand(monacoInstance.KeyMod.CtrlCmd | monacoInstance.KeyCode.KeyS, () => {
      if (onSave) {
        onSave(editor.getValue());
      }
    });

      editor.addCommand(monacoInstance.KeyMod.Shift | monacoInstance.KeyMod.Alt | monacoInstance.KeyCode.KeyF, () => {
        handleFormat();
      });
    }
  }, [onSave, onValidate]);

  const handleFormat = useCallback(() => {
    if (editorRef.current) {
      editorRef.current.trigger('keyboard', 'editor.action.formatDocument', {});
      if (onFormat) {
        onFormat();
      }
    }
  }, [onFormat]);

  const handleSave = useCallback(() => {
    if (onSave && editorRef.current) {
      onSave(editorRef.current.getValue());
    }
  }, [onSave]);

  const handleCopy = useCallback(() => {
    if (editorRef.current) {
      const selection = editorRef.current.getSelection();
      const selectedText = selection 
        ? editorRef.current.getModel()?.getValueInRange(selection) 
        : editorRef.current.getValue();
      
      if (selectedText) {
        navigator.clipboard.writeText(selectedText);
      }
    }
  }, []);

  const handleDownload = useCallback(() => {
    if (editorRef.current) {
      const content = editorRef.current.getValue();
      const blob = new Blob([content], { type: 'text/plain' });
      const url = URL.createObjectURL(blob);
      const a = document.createElement('a');
      a.href = url;
      a.download = filename;
      document.body.appendChild(a);
      a.click();
      document.body.removeChild(a);
      URL.revokeObjectURL(url);
    }
  }, [filename]);

  const handleSearch = useCallback(() => {
    if (editorRef.current) {
      editorRef.current.trigger('keyboard', 'actions.find', {});
    }
  }, []);

  const toggleFullscreen = useCallback(() => {
    setIsFullscreen(!isFullscreen);
  }, [isFullscreen]);

  const applySettings = useCallback(() => {
    if (editorRef.current) {
      editorRef.current.updateOptions({
        fontSize: localSettings.fontSize,
        tabSize: localSettings.tabSize,
        minimap: { enabled: localSettings.showMinimap },
        wordWrap: localSettings.wordWrap,
        lineNumbers: localSettings.showLineNumbers ? 'on' : 'off'
      });
      
      const monacoInstance = (window as any).monaco;
      if (monacoInstance) {
        monacoInstance.editor.setTheme(localSettings.theme);
      }
    }
    setShowSettings(false);
  }, [localSettings]);

  // Expose handle methods
  React.useImperativeHandle(ref, () => ({
    format: handleFormat,
    getContent: () => editorRef.current?.getValue() || '',
    setContent: (content: string) => editorRef.current?.setValue(content),
    focus: () => editorRef.current?.focus(),
    getSelection: () => {
      const selection = editorRef.current?.getSelection();
      return selection 
        ? editorRef.current?.getModel()?.getValueInRange(selection) || ''
        : '';
    },
    insertText: (text: string) => {
      const editor = editorRef.current;
      if (editor) {
        const selection = editor.getSelection();
        if (selection) {
          editor.executeEdits('', [{
            range: selection,
            text: text
          }]);
        }
      }
    },
    findAndReplace: (searchValue: string, replaceValue: string) => {
      const editor = editorRef.current;
      if (editor) {
        const model = editor.getModel();
        if (model) {
          const matches = model.findMatches(searchValue, true, false, true, null, true);
          editor.executeEdits('', matches.map(match => ({
            range: match.range,
            text: replaceValue
          })));
        }
      }
    },
    goToLine: (lineNumber: number) => {
      const editor = editorRef.current;
      if (editor) {
        editor.revealLineInCenter(lineNumber);
        editor.setPosition({ lineNumber, column: 1 });
      }
    }
  }), [handleFormat]);

  return (
    <div className={`code-editor ${isFullscreen ? 'code-editor--fullscreen' : ''} ${className}`}>
      <div className="code-editor__toolbar">
        <div className="code-editor__file-info">
          <FileText size={16} />
          <span className="code-editor__filename">{filename}</span>
          <span className="code-editor__language">{language}</span>
        </div>
        
        <div className="code-editor__actions">
          {!readOnly && onSave && (
            <button 
              className="code-editor__action" 
              onClick={handleSave}
              title="Save (Ctrl+S)"
            >
              <Save size={16} />
            </button>
          )}
          
          <button 
            className="code-editor__action" 
            onClick={handleFormat}
            title="Format Code (Shift+Alt+F)"
          >
            <FileText size={16} />
          </button>
          
          <button 
            className="code-editor__action" 
            onClick={handleCopy}
            title="Copy to Clipboard"
          >
            <Copy size={16} />
          </button>
          
          <button 
            className="code-editor__action" 
            onClick={handleDownload}
            title="Download File"
          >
            <Download size={16} />
          </button>
          
          <button 
            className="code-editor__action" 
            onClick={handleSearch}
            title="Find & Replace (Ctrl+F)"
          >
            <Search size={16} />
          </button>
          
          <button 
            className="code-editor__action" 
            onClick={() => setShowSettings(!showSettings)}
            title="Editor Settings"
          >
            <Settings size={16} />
          </button>
          
          <button 
            className="code-editor__action" 
            onClick={toggleFullscreen}
            title="Toggle Fullscreen"
          >
            {isFullscreen ? <Minimize2 size={16} /> : <Maximize2 size={16} />}
          </button>
        </div>
      </div>

      {showSettings && (
        <div className="code-editor__settings">
          <div className="code-editor__settings-content">
            <h3>Editor Settings</h3>
            
            <div className="code-editor__setting-row">
              <label>Font Size:</label>
              <input
                type="number"
                min="8"
                max="32"
                value={localSettings.fontSize}
                onChange={(e) => setLocalSettings(prev => ({ 
                  ...prev, 
                  fontSize: parseInt(e.target.value) 
                }))}
              />
            </div>
            
            <div className="code-editor__setting-row">
              <label>Tab Size:</label>
              <input
                type="number"
                min="1"
                max="8"
                value={localSettings.tabSize}
                onChange={(e) => setLocalSettings(prev => ({ 
                  ...prev, 
                  tabSize: parseInt(e.target.value) 
                }))}
              />
            </div>
            
            <div className="code-editor__setting-row">
              <label>Theme:</label>
              <select
                value={localSettings.theme}
                onChange={(e) => setLocalSettings(prev => ({ 
                  ...prev, 
                  theme: e.target.value as any 
                }))}
              >
                <option value="vs-dark">Dark</option>
                <option value="light">Light</option>
                <option value="vs">Classic</option>
              </select>
            </div>
            
            <div className="code-editor__setting-row">
              <label>Word Wrap:</label>
              <select
                value={localSettings.wordWrap}
                onChange={(e) => setLocalSettings(prev => ({ 
                  ...prev, 
                  wordWrap: e.target.value as any 
                }))}
              >
                <option value="on">On</option>
                <option value="off">Off</option>
                <option value="bounded">Bounded</option>
              </select>
            </div>
            
            <div className="code-editor__setting-row">
              <label>
                <input
                  type="checkbox"
                  checked={localSettings.showMinimap}
                  onChange={(e) => setLocalSettings(prev => ({ 
                    ...prev, 
                    showMinimap: e.target.checked 
                  }))}
                />
                Show Minimap
              </label>
            </div>
            
            <div className="code-editor__setting-row">
              <label>
                <input
                  type="checkbox"
                  checked={localSettings.showLineNumbers}
                  onChange={(e) => setLocalSettings(prev => ({ 
                    ...prev, 
                    showLineNumbers: e.target.checked 
                  }))}
                />
                Show Line Numbers
              </label>
            </div>
            
            <div className="code-editor__settings-actions">
              <button onClick={applySettings}>Apply</button>
              <button onClick={() => setShowSettings(false)}>Cancel</button>
            </div>
          </div>
        </div>
      )}

      <div className="code-editor__editor">
        <Editor
          height="100%"
          defaultLanguage={language}
          language={language}
          value={value}
          onChange={onChange}
          onMount={handleEditorDidMount}
          theme={localSettings.theme}
          options={{
            readOnly,
            fontSize: localSettings.fontSize,
            tabSize: localSettings.tabSize,
            minimap: { enabled: localSettings.showMinimap },
            wordWrap: localSettings.wordWrap,
            lineNumbers: localSettings.showLineNumbers ? 'on' : 'off',
            automaticLayout,
            scrollBeyondLastLine: false,
            renderWhitespace: 'boundary',
            renderControlCharacters: true,
            folding: true,
            foldingStrategy: 'auto',
            showFoldingControls: 'mouseover',
            matchBrackets: 'always',
            find: {
              addExtraSpaceOnTop: false,
              autoFindInSelection: 'never',
              seedSearchStringFromSelection: 'always'
            },
            quickSuggestions: {
              other: true,
              comments: false,
              strings: false
            },
            suggestOnTriggerCharacters: true,
            acceptSuggestionOnEnter: 'on',
            acceptSuggestionOnCommitCharacter: true,
            snippetSuggestions: 'top',
            emptySelectionClipboard: false,
            copyWithSyntaxHighlighting: false,
            formatOnPaste: true,
            formatOnType: true,
            autoIndent: 'advanced',
            dragAndDrop: true,
            mouseWheelZoom: true,
            multiCursorModifier: 'ctrlCmd',
            contextmenu: true
          }}
        />
      </div>
    </div>
  );
});

CodeEditor.displayName = 'CodeEditor';

export default CodeEditor;