import React, { useState, useEffect } from 'react';
import { Template, TemplateValidationResult, TemplateTestRequest, TemplateTestResult } from '../../types';
import './TemplateView.scss';

interface TemplateViewProps {
  template: Template;
  onEdit?: (template: Template) => void;
  onDelete?: (template: Template) => void;
  onDuplicate?: (template: Template) => void;
  onClose?: () => void;
}

const TemplateView: React.FC<TemplateViewProps> = ({
  template,
  onEdit,
  onDelete,
  onDuplicate,
  onClose,
}) => {
  const [activeTab, setActiveTab] = useState<'overview' | 'variables' | 'validation' | 'test' | 'versions'>('overview');
  const [validationResult, setValidationResult] = useState<TemplateValidationResult | null>(null);
  const [testResult, setTestResult] = useState<TemplateTestResult | null>(null);
  const [testData, setTestData] = useState<Record<string, any>>({});
  const [versions, setVersions] = useState<Template[]>([]);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    if (activeTab === 'validation') {
      validateTemplate();
    } else if (activeTab === 'versions') {
      loadVersions();
    }
  }, [activeTab, template.id]);

  useEffect(() => {
    // Initialize test data with default values
    const initialTestData: Record<string, any> = {};
    template.variables.forEach(variable => {
      if (variable.defaultValue) {
        initialTestData[variable.name] = variable.defaultValue;
      } else {
        switch (variable.type) {
          case 'string':
            initialTestData[variable.name] = '';
            break;
          case 'int':
            initialTestData[variable.name] = 0;
            break;
          case 'bool':
            initialTestData[variable.name] = false;
            break;
          case 'float':
            initialTestData[variable.name] = 0.0;
            break;
          case 'array':
            initialTestData[variable.name] = [];
            break;
          case 'object':
            initialTestData[variable.name] = {};
            break;
          default:
            initialTestData[variable.name] = '';
        }
      }
    });
    setTestData(initialTestData);
  }, [template.variables]);

  const validateTemplate = async () => {
    setLoading(true);
    setError(null);
    try {
      const result = await window.go.app.App.ValidateTemplateAdvanced(template.id);
      setValidationResult(result);
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to validate template');
    } finally {
      setLoading(false);
    }
  };

  const testTemplate = async () => {
    setLoading(true);
    setError(null);
    try {
      const testRequest: TemplateTestRequest = {
        templateId: template.id,
        testData,
        options: {
          validateOutput: true,
          measurePerformance: true,
          generateReport: true,
        },
      };

      const result = await window.go.app.App.TestTemplate(testRequest);
      setTestResult(result);
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to test template');
    } finally {
      setLoading(false);
    }
  };

  const loadVersions = async () => {
    setLoading(true);
    setError(null);
    try {
      const versionList = await window.go.app.App.GetTemplateVersions(template.name);
      setVersions(versionList);
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to load versions');
    } finally {
      setLoading(false);
    }
  };

  const handleTestDataChange = (variableName: string, value: any) => {
    setTestData(prev => ({ ...prev, [variableName]: value }));
  };

  const formatFileSize = (bytes: number) => {
    if (bytes === 0) return '0 B';
    const k = 1024;
    const sizes = ['B', 'KB', 'MB', 'GB'];
    const i = Math.floor(Math.log(bytes) / Math.log(k));
    return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i];
  };

  const getValidationStatusIcon = (valid: boolean) => {
    return valid ? '✅' : '❌';
  };

  const getSeverityColor = (severity: string) => {
    switch (severity) {
      case 'error':
        return 'var(--error-color)';
      case 'warning':
        return 'var(--warning-color)';
      case 'info':
        return 'var(--info-color)';
      default:
        return 'var(--text-secondary)';
    }
  };

  return (
    <div className="template-view">
      <div className="template-view__header">
        <div className="template-view__title">
          <h1>{template.name}</h1>
          <span className="template-view__version">v{template.version}</span>
          {template.isBuiltIn && (
            <span className="template-view__badge template-view__badge--builtin">Built-in</span>
          )}
          <span className={`template-view__badge template-view__badge--${template.type}`}>
            {template.type}
          </span>
        </div>
        
        <div className="template-view__actions">
          <button
            onClick={() => onDuplicate?.(template)}
            className="template-view__action"
            title="Duplicate template"
          >
            📋 Duplicate
          </button>
          {!template.isBuiltIn && (
            <>
              <button
                onClick={() => onEdit?.(template)}
                className="template-view__action"
                title="Edit template"
              >
                ✏️ Edit
              </button>
              <button
                onClick={() => onDelete?.(template)}
                className="template-view__action template-view__action--danger"
                title="Delete template"
              >
                🗑️ Delete
              </button>
            </>
          )}
          {onClose && (
            <button
              onClick={onClose}
              className="template-view__action template-view__action--close"
              title="Close"
            >
              ✕
            </button>
          )}
        </div>
      </div>

      <div className="template-view__tabs">
        <button
          className={`template-view__tab ${activeTab === 'overview' ? 'template-view__tab--active' : ''}`}
          onClick={() => setActiveTab('overview')}
        >
          Overview
        </button>
        <button
          className={`template-view__tab ${activeTab === 'variables' ? 'template-view__tab--active' : ''}`}
          onClick={() => setActiveTab('variables')}
        >
          Variables ({template.variables.length})
        </button>
        <button
          className={`template-view__tab ${activeTab === 'validation' ? 'template-view__tab--active' : ''}`}
          onClick={() => setActiveTab('validation')}
        >
          Validation
        </button>
        <button
          className={`template-view__tab ${activeTab === 'test' ? 'template-view__tab--active' : ''}`}
          onClick={() => setActiveTab('test')}
        >
          Test
        </button>
        <button
          className={`template-view__tab ${activeTab === 'versions' ? 'template-view__tab--active' : ''}`}
          onClick={() => setActiveTab('versions')}
        >
          Versions
        </button>
      </div>

      <div className="template-view__content">
        {error && (
          <div className="template-view__error">
            <p>{error}</p>
            <button onClick={() => setError(null)}>Dismiss</button>
          </div>
        )}

        {activeTab === 'overview' && (
          <div className="template-view__overview">
            <div className="template-view__info-grid">
              <div className="template-view__info-item">
                <label>Description</label>
                <p>{template.description || 'No description available'}</p>
              </div>
              
              <div className="template-view__info-item">
                <label>Author</label>
                <p>{template.author || 'Unknown'}</p>
              </div>
              
              <div className="template-view__info-item">
                <label>File Path</label>
                <p className="template-view__path">{template.path}</p>
              </div>
              
              <div className="template-view__info-item">
                <label>Created</label>
                <p>{new Date(template.createdAt).toLocaleString()}</p>
              </div>
              
              <div className="template-view__info-item">
                <label>Last Updated</label>
                <p>{new Date(template.updatedAt).toLocaleString()}</p>
              </div>
            </div>
          </div>
        )}

        {activeTab === 'variables' && (
          <div className="template-view__variables">
            {template.variables.length === 0 ? (
              <p className="template-view__empty">No variables defined for this template.</p>
            ) : (
              <div className="template-view__variables-grid">
                {template.variables.map((variable, index) => (
                  <div key={index} className="template-view__variable">
                    <div className="template-view__variable-header">
                      <h4>{variable.name}</h4>
                      <span className={`template-view__variable-type template-view__variable-type--${variable.type}`}>
                        {variable.type}
                      </span>
                      {variable.required && (
                        <span className="template-view__variable-required">Required</span>
                      )}
                    </div>
                    
                    {variable.description && (
                      <p className="template-view__variable-description">{variable.description}</p>
                    )}
                    
                    <div className="template-view__variable-details">
                      {variable.defaultValue && (
                        <div className="template-view__variable-detail">
                          <span>Default:</span> <code>{variable.defaultValue}</code>
                        </div>
                      )}
                      
                      {variable.options && variable.options.length > 0 && (
                        <div className="template-view__variable-detail">
                          <span>Options:</span> {variable.options.join(', ')}
                        </div>
                      )}
                      
                      {variable.validation && (
                        <div className="template-view__variable-detail">
                          <span>Validation:</span> <code>{variable.validation}</code>
                        </div>
                      )}
                    </div>
                  </div>
                ))}
              </div>
            )}
          </div>
        )}

        {activeTab === 'validation' && (
          <div className="template-view__validation">
            {loading ? (
              <div className="template-view__loading">Validating template...</div>
            ) : validationResult ? (
              <div className="template-view__validation-result">
                <div className="template-view__validation-status">
                  <h3>
                    {getValidationStatusIcon(validationResult.valid)} 
                    Validation {validationResult.valid ? 'Passed' : 'Failed'}
                  </h3>
                  <button onClick={validateTemplate} className="template-view__refresh-button">
                    🔄 Re-validate
                  </button>
                </div>

                {validationResult.errors && validationResult.errors.length > 0 && (
                  <div className="template-view__validation-section">
                    <h4>Errors ({validationResult.errors.length})</h4>
                    <div className="template-view__validation-items">
                      {validationResult.errors.map((error, index) => (
                        <div key={index} className="template-view__validation-item template-view__validation-item--error">
                          <div className="template-view__validation-item-header">
                            <span className="template-view__validation-type">{error.type}</span>
                            {error.line && <span className="template-view__validation-line">Line {error.line}</span>}
                          </div>
                          <p>{error.message}</p>
                        </div>
                      ))}
                    </div>
                  </div>
                )}

                {validationResult.warnings && validationResult.warnings.length > 0 && (
                  <div className="template-view__validation-section">
                    <h4>Warnings ({validationResult.warnings.length})</h4>
                    <div className="template-view__validation-items">
                      {validationResult.warnings.map((warning, index) => (
                        <div key={index} className="template-view__validation-item template-view__validation-item--warning">
                          <div className="template-view__validation-item-header">
                            <span className="template-view__validation-type">{warning.type}</span>
                            {warning.line && <span className="template-view__validation-line">Line {warning.line}</span>}
                          </div>
                          <p>{warning.message}</p>
                          {warning.suggestion && (
                            <p className="template-view__validation-suggestion">💡 {warning.suggestion}</p>
                          )}
                        </div>
                      ))}
                    </div>
                  </div>
                )}

                {validationResult.suggestions && validationResult.suggestions.length > 0 && (
                  <div className="template-view__validation-section">
                    <h4>Suggestions</h4>
                    <ul className="template-view__suggestions">
                      {validationResult.suggestions.map((suggestion, index) => (
                        <li key={index}>{suggestion}</li>
                      ))}
                    </ul>
                  </div>
                )}

                {validationResult.performance && (
                  <div className="template-view__validation-section">
                    <h4>Performance Analysis</h4>
                    <div className="template-view__performance">
                      <div className="template-view__performance-item">
                        <span>Complexity:</span>
                        <span className={`template-view__complexity template-view__complexity--${validationResult.performance.complexity}`}>
                          {validationResult.performance.complexity}
                        </span>
                      </div>
                      <div className="template-view__performance-item">
                        <span>Memory Usage:</span>
                        <span>{formatFileSize(validationResult.performance.memoryUsage)}</span>
                      </div>
                    </div>
                  </div>
                )}
              </div>
            ) : (
              <div className="template-view__validation-empty">
                <p>Click "Re-validate" to run validation checks on this template.</p>
                <button onClick={validateTemplate} className="template-view__action">
                  🔍 Validate Template
                </button>
              </div>
            )}
          </div>
        )}

        {activeTab === 'test' && (
          <div className="template-view__test">
            <div className="template-view__test-section">
              <h3>Test Data</h3>
              {template.variables.length === 0 ? (
                <p>No variables to configure for testing.</p>
              ) : (
                <div className="template-view__test-inputs">
                  {template.variables.map((variable) => (
                    <div key={variable.name} className="template-view__test-input">
                      <label>{variable.name}</label>
                      {variable.type === 'bool' ? (
                        <select
                          value={testData[variable.name] ? 'true' : 'false'}
                          onChange={(e) => handleTestDataChange(variable.name, e.target.value === 'true')}
                        >
                          <option value="true">true</option>
                          <option value="false">false</option>
                        </select>
                      ) : variable.type === 'enum' && variable.options ? (
                        <select
                          value={testData[variable.name] || ''}
                          onChange={(e) => handleTestDataChange(variable.name, e.target.value)}
                        >
                          <option value="">Select option...</option>
                          {variable.options.map((option) => (
                            <option key={option} value={option}>{option}</option>
                          ))}
                        </select>
                      ) : (
                        <input
                          type={variable.type === 'int' || variable.type === 'float' ? 'number' : 'text'}
                          value={testData[variable.name] || ''}
                          onChange={(e) => {
                            let value: any = e.target.value;
                            if (variable.type === 'int') {
                              value = parseInt(value) || 0;
                            } else if (variable.type === 'float') {
                              value = parseFloat(value) || 0;
                            }
                            handleTestDataChange(variable.name, value);
                          }}
                          placeholder={variable.description || `Enter ${variable.name}`}
                        />
                      )}
                    </div>
                  ))}
                </div>
              )}
              
              <button
                onClick={testTemplate}
                disabled={loading}
                className="template-view__test-button"
              >
                {loading ? 'Testing...' : '🧪 Test Template'}
              </button>
            </div>

            {testResult && (
              <div className="template-view__test-result">
                <h3>Test Results</h3>
                <div className="template-view__test-status">
                  {getValidationStatusIcon(testResult.success)} 
                  Test {testResult.success ? 'Passed' : 'Failed'}
                </div>

                {testResult.output && (
                  <div className="template-view__test-output">
                    <h4>Generated Output</h4>
                    <pre><code>{testResult.output}</code></pre>
                  </div>
                )}

                {testResult.errors && testResult.errors.length > 0 && (
                  <div className="template-view__test-errors">
                    <h4>Errors</h4>
                    {testResult.errors.map((error, index) => (
                      <div key={index} className="template-view__test-error">
                        <span className="template-view__error-type">{error.type}</span>
                        <p>{error.message}</p>
                      </div>
                    ))}
                  </div>
                )}

                {testResult.performance && (
                  <div className="template-view__test-performance">
                    <h4>Performance Metrics</h4>
                    <div className="template-view__performance-grid">
                      <div className="template-view__performance-item">
                        <span>Render Time:</span>
                        <span>{testResult.performance.renderTime}ms</span>
                      </div>
                      <div className="template-view__performance-item">
                        <span>Memory Usage:</span>
                        <span>{formatFileSize(testResult.performance.memoryUsage)}</span>
                      </div>
                      <div className="template-view__performance-item">
                        <span>Complexity:</span>
                        <span>{testResult.performance.complexity}</span>
                      </div>
                    </div>
                  </div>
                )}
              </div>
            )}
          </div>
        )}

        {activeTab === 'versions' && (
          <div className="template-view__versions">
            {loading ? (
              <div className="template-view__loading">Loading versions...</div>
            ) : versions.length === 0 ? (
              <p className="template-view__empty">No versions found for this template.</p>
            ) : (
              <div className="template-view__versions-list">
                <h3>Version History ({versions.length})</h3>
                {versions.map((version) => (
                  <div
                    key={version.id}
                    className={`template-view__version ${version.id === template.id ? 'template-view__version--current' : ''}`}
                  >
                    <div className="template-view__version-header">
                      <span className="template-view__version-number">v{version.version}</span>
                      {version.id === template.id && (
                        <span className="template-view__version-current">Current</span>
                      )}
                      <span className="template-view__version-date">
                        {new Date(version.createdAt).toLocaleDateString()}
                      </span>
                    </div>
                    <p className="template-view__version-description">{version.description}</p>
                  </div>
                ))}
              </div>
            )}
          </div>
        )}
      </div>
    </div>
  );
};

export default TemplateView;