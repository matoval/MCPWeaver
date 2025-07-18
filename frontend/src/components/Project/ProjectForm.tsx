import React, { useState, useEffect } from 'react';
import { 
  Save, 
  X, 
  Upload, 
  Link, 
  Folder, 
  FileText, 
  AlertCircle, 
  CheckCircle,
  Settings,
  Info
} from 'lucide-react';
import { 
  CreateProject, 
  UpdateProject, 
  SelectFile, 
  SelectDirectory,
  ValidateSpec,
  ValidateURL,
  GetDefaultOpenAPIFilters
} from '../../wailsjs/go/app/App';
import { app } from '../../wailsjs/go/models';
import './ProjectForm.scss';

interface ProjectFormProps {
  project?: app.Project;
  onSave: (project: app.Project) => void;
  onCancel: () => void;
}

const ProjectForm: React.FC<ProjectFormProps> = ({ project, onSave, onCancel }) => {
  const [formData, setFormData] = useState({
    name: '',
    specPath: '',
    specUrl: '',
    outputPath: '',
    settings: {
      packageName: 'generated-server',
      serverPort: 8080,
      enableLogging: true,
      logLevel: 'info',
      customTemplates: [] as string[]
    }
  });

  const [validation, setValidation] = useState<app.ValidationResult | null>(null);
  const [isValidating, setIsValidating] = useState(false);
  const [isSaving, setIsSaving] = useState(false);
  const [errors, setErrors] = useState<Record<string, string>>({});
  const [specSource, setSpecSource] = useState<'file' | 'url'>('file');
  const [showAdvancedSettings, setShowAdvancedSettings] = useState(false);

  useEffect(() => {
    if (project) {
      setFormData({
        name: project.name,
        specPath: project.specPath,
        specUrl: project.specUrl,
        outputPath: project.outputPath,
        settings: {
          packageName: project.settings.packageName,
          serverPort: project.settings.serverPort,
          enableLogging: project.settings.enableLogging,
          logLevel: project.settings.logLevel,
          customTemplates: project.settings.customTemplates || []
        }
      });
      setSpecSource(project.specPath ? 'file' : 'url');
    }
  }, [project]);

  const validateForm = () => {
    const newErrors: Record<string, string> = {};

    if (!formData.name.trim()) {
      newErrors.name = 'Project name is required';
    }

    if (!formData.outputPath.trim()) {
      newErrors.outputPath = 'Output path is required';
    }

    if (specSource === 'file' && !formData.specPath.trim()) {
      newErrors.specPath = 'OpenAPI specification file is required';
    }

    if (specSource === 'url' && !formData.specUrl.trim()) {
      newErrors.specUrl = 'OpenAPI specification URL is required';
    }

    if (!formData.settings.packageName.trim()) {
      newErrors.packageName = 'Package name is required';
    }

    if (formData.settings.serverPort < 1000 || formData.settings.serverPort > 65535) {
      newErrors.serverPort = 'Server port must be between 1000 and 65535';
    }

    setErrors(newErrors);
    return Object.keys(newErrors).length === 0;
  };

  const handleSelectFile = async () => {
    try {
      const filters = await GetDefaultOpenAPIFilters();
      const filePath = await SelectFile(filters);
      if (filePath) {
        setFormData(prev => ({ ...prev, specPath: filePath }));
        validateSpec(filePath);
      }
    } catch (error) {
      console.error('Failed to select file:', error);
    }
  };

  const handleSelectDirectory = async () => {
    try {
      const dirPath = await SelectDirectory('Select Output Directory');
      if (dirPath) {
        setFormData(prev => ({ ...prev, outputPath: dirPath }));
      }
    } catch (error) {
      console.error('Failed to select directory:', error);
    }
  };

  const validateSpec = async (pathOrUrl: string) => {
    if (!pathOrUrl.trim()) return;

    setIsValidating(true);
    setValidation(null);

    try {
      let result: app.ValidationResult;
      if (specSource === 'file') {
        result = await ValidateSpec(pathOrUrl);
      } else {
        result = await ValidateURL(pathOrUrl);
      }
      setValidation(result);
    } catch (error) {
      console.error('Failed to validate spec:', error);
      setValidation({
        valid: false,
        errors: [{
          type: 'validation',
          message: 'Failed to validate specification',
          path: '',
          severity: 'error',
          code: 'VALIDATION_FAILED'
        }],
        warnings: [],
        suggestions: [],
        validationTime: 0
      } as app.ValidationResult);
    } finally {
      setIsValidating(false);
    }
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    
    if (!validateForm()) return;

    setIsSaving(true);

    try {
      const requestData = {
        name: formData.name,
        specPath: specSource === 'file' ? formData.specPath : '',
        specUrl: specSource === 'url' ? formData.specUrl : '',
        outputPath: formData.outputPath,
        settings: formData.settings
      };

      let savedProject: app.Project;

      if (project) {
        savedProject = await UpdateProject(project.id, requestData);
      } else {
        savedProject = await CreateProject(requestData);
      }

      onSave(savedProject);
    } catch (error) {
      console.error('Failed to save project:', error);
      setErrors({ submit: 'Failed to save project. Please try again.' });
    } finally {
      setIsSaving(false);
    }
  };

  const handleInputChange = (field: string, value: any) => {
    setFormData(prev => ({ ...prev, [field]: value }));
    if (errors[field]) {
      setErrors(prev => ({ ...prev, [field]: '' }));
    }
  };

  const handleSettingsChange = (field: string, value: any) => {
    setFormData(prev => ({
      ...prev,
      settings: { ...prev.settings, [field]: value }
    }));
    if (errors[field]) {
      setErrors(prev => ({ ...prev, [field]: '' }));
    }
  };

  const handleSpecSourceChange = (source: 'file' | 'url') => {
    setSpecSource(source);
    setValidation(null);
    setFormData(prev => ({ ...prev, specPath: '', specUrl: '' }));
  };

  const handleSpecValueChange = (value: string) => {
    if (specSource === 'file') {
      setFormData(prev => ({ ...prev, specPath: value }));
    } else {
      setFormData(prev => ({ ...prev, specUrl: value }));
    }
    
    // Debounce validation
    clearTimeout(window.validateTimeout);
    window.validateTimeout = setTimeout(() => {
      validateSpec(value);
    }, 1000);
  };

  return (
    <div className="project-form-overlay">
      <div className="project-form">
        <div className="project-form-header">
          <h2>{project ? 'Edit Project' : 'Create New Project'}</h2>
          <button className="close-button" onClick={onCancel}>
            <X size={20} />
          </button>
        </div>

        <form onSubmit={handleSubmit} className="project-form-content">
          <div className="form-section">
            <h3>Project Information</h3>
            
            <div className="form-group">
              <label htmlFor="name">Project Name *</label>
              <input
                id="name"
                type="text"
                value={formData.name}
                onChange={(e) => handleInputChange('name', e.target.value)}
                placeholder="Enter project name"
                className={errors.name ? 'error' : ''}
              />
              {errors.name && <span className="error-message">{errors.name}</span>}
            </div>

            <div className="form-group">
              <label htmlFor="outputPath">Output Directory *</label>
              <div className="input-with-button">
                <input
                  id="outputPath"
                  type="text"
                  value={formData.outputPath}
                  onChange={(e) => handleInputChange('outputPath', e.target.value)}
                  placeholder="Select output directory"
                  className={errors.outputPath ? 'error' : ''}
                />
                <button type="button" onClick={handleSelectDirectory}>
                  <Folder size={16} />
                  Browse
                </button>
              </div>
              {errors.outputPath && <span className="error-message">{errors.outputPath}</span>}
            </div>
          </div>

          <div className="form-section">
            <h3>OpenAPI Specification</h3>
            
            <div className="spec-source-tabs">
              <button
                type="button"
                className={specSource === 'file' ? 'active' : ''}
                onClick={() => handleSpecSourceChange('file')}
              >
                <FileText size={16} />
                File
              </button>
              <button
                type="button"
                className={specSource === 'url' ? 'active' : ''}
                onClick={() => handleSpecSourceChange('url')}
              >
                <Link size={16} />
                URL
              </button>
            </div>

            {specSource === 'file' ? (
              <div className="form-group">
                <label htmlFor="specPath">OpenAPI File *</label>
                <div className="input-with-button">
                  <input
                    id="specPath"
                    type="text"
                    value={formData.specPath}
                    onChange={(e) => handleSpecValueChange(e.target.value)}
                    placeholder="Select OpenAPI specification file"
                    className={errors.specPath ? 'error' : ''}
                  />
                  <button type="button" onClick={handleSelectFile}>
                    <Upload size={16} />
                    Browse
                  </button>
                </div>
                {errors.specPath && <span className="error-message">{errors.specPath}</span>}
              </div>
            ) : (
              <div className="form-group">
                <label htmlFor="specUrl">OpenAPI URL *</label>
                <input
                  id="specUrl"
                  type="url"
                  value={formData.specUrl}
                  onChange={(e) => handleSpecValueChange(e.target.value)}
                  placeholder="https://api.example.com/openapi.json"
                  className={errors.specUrl ? 'error' : ''}
                />
                {errors.specUrl && <span className="error-message">{errors.specUrl}</span>}
              </div>
            )}

            {isValidating && (
              <div className="validation-status validating">
                <div className="spinner"></div>
                <span>Validating specification...</span>
              </div>
            )}

            {validation && (
              <div className={`validation-status ${validation.valid ? 'valid' : 'invalid'}`}>
                {validation.valid ? (
                  <>
                    <CheckCircle size={16} />
                    <span>Specification is valid</span>
                  </>
                ) : (
                  <>
                    <AlertCircle size={16} />
                    <span>Specification has errors</span>
                  </>
                )}
              </div>
            )}

            {validation && validation.specInfo && (
              <div className="spec-info">
                <h4>Specification Details</h4>
                <div className="spec-details">
                  <div className="spec-detail">
                    <strong>Title:</strong> {validation.specInfo.title}
                  </div>
                  <div className="spec-detail">
                    <strong>Version:</strong> {validation.specInfo.version}
                  </div>
                  <div className="spec-detail">
                    <strong>Operations:</strong> {validation.specInfo.operationCount}
                  </div>
                  <div className="spec-detail">
                    <strong>Schemas:</strong> {validation.specInfo.schemaCount}
                  </div>
                </div>
              </div>
            )}

            {validation && validation.errors && validation.errors.length > 0 && (
              <div className="validation-errors">
                <h4>Validation Errors</h4>
                {validation.errors.map((error, index) => (
                  <div key={index} className="validation-error">
                    <AlertCircle size={14} />
                    <span>{error.message}</span>
                    {error.path && <span className="error-path">at {error.path}</span>}
                  </div>
                ))}
              </div>
            )}

            {validation && validation.warnings && validation.warnings.length > 0 && (
              <div className="validation-warnings">
                <h4>Warnings</h4>
                {validation.warnings.map((warning, index) => (
                  <div key={index} className="validation-warning">
                    <Info size={14} />
                    <span>{warning.message}</span>
                    {warning.suggestion && <span className="warning-suggestion">{warning.suggestion}</span>}
                  </div>
                ))}
              </div>
            )}
          </div>

          <div className="form-section">
            <div className="section-header">
              <h3>Server Settings</h3>
              <button
                type="button"
                className="toggle-advanced"
                onClick={() => setShowAdvancedSettings(!showAdvancedSettings)}
              >
                <Settings size={16} />
                {showAdvancedSettings ? 'Hide Advanced' : 'Show Advanced'}
              </button>
            </div>

            <div className="form-group">
              <label htmlFor="packageName">Package Name *</label>
              <input
                id="packageName"
                type="text"
                value={formData.settings.packageName}
                onChange={(e) => handleSettingsChange('packageName', e.target.value)}
                placeholder="generated-server"
                className={errors.packageName ? 'error' : ''}
              />
              {errors.packageName && <span className="error-message">{errors.packageName}</span>}
            </div>

            <div className="form-group">
              <label htmlFor="serverPort">Server Port *</label>
              <input
                id="serverPort"
                type="number"
                min="1000"
                max="65535"
                value={formData.settings.serverPort}
                onChange={(e) => handleSettingsChange('serverPort', parseInt(e.target.value))}
                className={errors.serverPort ? 'error' : ''}
              />
              {errors.serverPort && <span className="error-message">{errors.serverPort}</span>}
            </div>

            {showAdvancedSettings && (
              <>
                <div className="form-group">
                  <label className="checkbox-label">
                    <input
                      type="checkbox"
                      checked={formData.settings.enableLogging}
                      onChange={(e) => handleSettingsChange('enableLogging', e.target.checked)}
                    />
                    Enable Logging
                  </label>
                </div>

                <div className="form-group">
                  <label htmlFor="logLevel">Log Level</label>
                  <select
                    id="logLevel"
                    value={formData.settings.logLevel}
                    onChange={(e) => handleSettingsChange('logLevel', e.target.value)}
                  >
                    <option value="debug">Debug</option>
                    <option value="info">Info</option>
                    <option value="warn">Warning</option>
                    <option value="error">Error</option>
                  </select>
                </div>
              </>
            )}
          </div>

          <div className="form-actions">
            <button type="button" className="cancel-button" onClick={onCancel}>
              Cancel
            </button>
            <button 
              type="submit" 
              className="save-button"
              disabled={isSaving || (validation && !validation.valid)}
            >
              {isSaving ? (
                <>
                  <div className="spinner"></div>
                  Saving...
                </>
              ) : (
                <>
                  <Save size={16} />
                  {project ? 'Update Project' : 'Create Project'}
                </>
              )}
            </button>
          </div>

          {errors.submit && (
            <div className="submit-error">
              <AlertCircle size={16} />
              <span>{errors.submit}</span>
            </div>
          )}
        </form>
      </div>
    </div>
  );
};

export default ProjectForm;