import React, { useState, useEffect } from 'react';
import { Template, TemplateType, TemplateVariable, CreateTemplateRequest, UpdateTemplateRequest } from '../../types';
import './TemplateForm.scss';

interface TemplateFormProps {
  template?: Template;
  onSave: (template: Template) => void;
  onCancel: () => void;
  isEditing?: boolean;
}

const TemplateForm: React.FC<TemplateFormProps> = ({
  template,
  onSave,
  onCancel,
  isEditing = false,
}) => {
  const [formData, setFormData] = useState({
    name: '',
    description: '',
    version: '1.0.0',
    author: '',
    type: 'custom' as TemplateType,
    path: '',
    variables: [] as TemplateVariable[],
  });

  const [errors, setErrors] = useState<Record<string, string>>({});
  const [saving, setSaving] = useState(false);

  useEffect(() => {
    if (template) {
      setFormData({
        name: template.name,
        description: template.description,
        version: template.version,
        author: template.author,
        type: template.type,
        path: template.path,
        variables: [...template.variables],
      });
    }
  }, [template]);

  const validateForm = () => {
    const newErrors: Record<string, string> = {};

    if (!formData.name.trim()) {
      newErrors.name = 'Template name is required';
    } else if (formData.name.length > 100) {
      newErrors.name = 'Template name must be 100 characters or less';
    }

    if (!formData.version.trim()) {
      newErrors.version = 'Version is required';
    } else if (!/^\d+\.\d+\.\d+$/.test(formData.version)) {
      newErrors.version = 'Version must be in format major.minor.patch (e.g., 1.0.0)';
    }

    if (!formData.path.trim()) {
      newErrors.path = 'Template file path is required';
    }

    if (formData.description.length > 500) {
      newErrors.description = 'Description must be 500 characters or less';
    }

    if (formData.author.length > 100) {
      newErrors.author = 'Author must be 100 characters or less';
    }

    // Validate variables
    formData.variables.forEach((variable, index) => {
      if (!variable.name.trim()) {
        newErrors[`variable_${index}_name`] = 'Variable name is required';
      }
      if (!variable.type.trim()) {
        newErrors[`variable_${index}_type`] = 'Variable type is required';
      }
    });

    setErrors(newErrors);
    return Object.keys(newErrors).length === 0;
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    
    if (!validateForm()) {
      return;
    }

    setSaving(true);
    try {
      let savedTemplate: Template;

      if (isEditing && template) {
        const updateRequest: UpdateTemplateRequest = {
          name: formData.name !== template.name ? formData.name : undefined,
          description: formData.description !== template.description ? formData.description : undefined,
          version: formData.version !== template.version ? formData.version : undefined,
          author: formData.author !== template.author ? formData.author : undefined,
          type: formData.type !== template.type ? formData.type : undefined,
          path: formData.path !== template.path ? formData.path : undefined,
          variables: JSON.stringify(formData.variables) !== JSON.stringify(template.variables) ? formData.variables : undefined,
        };

        savedTemplate = await window.go.app.App.UpdateTemplate(template.id, updateRequest);
      } else {
        const createRequest: CreateTemplateRequest = {
          name: formData.name,
          description: formData.description,
          version: formData.version,
          author: formData.author,
          type: formData.type,
          path: formData.path,
          variables: formData.variables,
        };

        savedTemplate = await window.go.app.App.CreateTemplate(createRequest);
      }

      onSave(savedTemplate);
    } catch (err) {
      setErrors({
        submit: err instanceof Error ? err.message : 'Failed to save template',
      });
    } finally {
      setSaving(false);
    }
  };

  const handleInputChange = (field: string, value: any) => {
    setFormData(prev => ({ ...prev, [field]: value }));
    
    // Clear error when user starts typing
    if (errors[field]) {
      setErrors(prev => ({ ...prev, [field]: '' }));
    }
  };

  const handleSelectFile = async () => {
    try {
      const filters = [
        { displayName: 'Template Files', pattern: '*.tmpl;*.template', extensions: ['tmpl', 'template'] },
        { displayName: 'All Files', pattern: '*', extensions: ['*'] },
      ];
      
      const filePath = await window.go.app.App.SelectFile(filters);
      if (filePath) {
        handleInputChange('path', filePath);
      }
    } catch (err) {
      setErrors(prev => ({
        ...prev,
        path: err instanceof Error ? err.message : 'Failed to select file',
      }));
    }
  };

  const addVariable = () => {
    const newVariable: TemplateVariable = {
      name: '',
      description: '',
      type: 'string',
      defaultValue: '',
      required: false,
      options: [],
      validation: '',
    };

    setFormData(prev => ({
      ...prev,
      variables: [...prev.variables, newVariable],
    }));
  };

  const updateVariable = (index: number, field: keyof TemplateVariable, value: any) => {
    setFormData(prev => ({
      ...prev,
      variables: prev.variables.map((variable, i) =>
        i === index ? { ...variable, [field]: value } : variable
      ),
    }));

    // Clear variable-specific errors
    const errorKey = `variable_${index}_${field}`;
    if (errors[errorKey]) {
      setErrors(prev => ({ ...prev, [errorKey]: '' }));
    }
  };

  const removeVariable = (index: number) => {
    setFormData(prev => ({
      ...prev,
      variables: prev.variables.filter((_, i) => i !== index),
    }));
  };

  const parseOptions = (optionsText: string): string[] => {
    return optionsText.split(',').map(opt => opt.trim()).filter(opt => opt.length > 0);
  };

  const formatOptions = (options: string[]): string => {
    return options.join(', ');
  };

  return (
    <div className="template-form">
      <div className="template-form__header">
        <h2>{isEditing ? 'Edit Template' : 'Create New Template'}</h2>
      </div>

      <form onSubmit={handleSubmit} className="template-form__form">
        {errors.submit && (
          <div className="template-form__error template-form__error--global">
            {errors.submit}
          </div>
        )}

        <div className="template-form__section">
          <h3>Basic Information</h3>
          
          <div className="template-form__row">
            <div className="template-form__field">
              <label htmlFor="name">Template Name *</label>
              <input
                id="name"
                type="text"
                value={formData.name}
                onChange={(e) => handleInputChange('name', e.target.value)}
                className={errors.name ? 'template-form__input--error' : ''}
                placeholder="Enter template name"
              />
              {errors.name && <span className="template-form__error">{errors.name}</span>}
            </div>

            <div className="template-form__field">
              <label htmlFor="version">Version *</label>
              <input
                id="version"
                type="text"
                value={formData.version}
                onChange={(e) => handleInputChange('version', e.target.value)}
                className={errors.version ? 'template-form__input--error' : ''}
                placeholder="1.0.0"
              />
              {errors.version && <span className="template-form__error">{errors.version}</span>}
            </div>
          </div>

          <div className="template-form__row">
            <div className="template-form__field">
              <label htmlFor="author">Author</label>
              <input
                id="author"
                type="text"
                value={formData.author}
                onChange={(e) => handleInputChange('author', e.target.value)}
                className={errors.author ? 'template-form__input--error' : ''}
                placeholder="Enter author name"
              />
              {errors.author && <span className="template-form__error">{errors.author}</span>}
            </div>

            <div className="template-form__field">
              <label htmlFor="type">Type</label>
              <select
                id="type"
                value={formData.type}
                onChange={(e) => handleInputChange('type', e.target.value as TemplateType)}
              >
                <option value="custom">Custom</option>
                <option value="plugin">Plugin</option>
                {isEditing && template?.isBuiltIn && <option value="default">Default</option>}
              </select>
            </div>
          </div>

          <div className="template-form__field">
            <label htmlFor="description">Description</label>
            <textarea
              id="description"
              value={formData.description}
              onChange={(e) => handleInputChange('description', e.target.value)}
              className={errors.description ? 'template-form__input--error' : ''}
              placeholder="Enter template description"
              rows={3}
            />
            {errors.description && <span className="template-form__error">{errors.description}</span>}
          </div>

          <div className="template-form__field">
            <label htmlFor="path">Template File Path *</label>
            <div className="template-form__file-input">
              <input
                id="path"
                type="text"
                value={formData.path}
                onChange={(e) => handleInputChange('path', e.target.value)}
                className={errors.path ? 'template-form__input--error' : ''}
                placeholder="Enter or select template file path"
              />
              <button
                type="button"
                onClick={handleSelectFile}
                className="template-form__file-button"
              >
                Browse
              </button>
            </div>
            {errors.path && <span className="template-form__error">{errors.path}</span>}
          </div>
        </div>

        <div className="template-form__section">
          <div className="template-form__section-header">
            <h3>Template Variables</h3>
            <button
              type="button"
              onClick={addVariable}
              className="template-form__add-button"
            >
              Add Variable
            </button>
          </div>

          {formData.variables.length === 0 ? (
            <p className="template-form__empty">No variables defined. Add variables to make your template configurable.</p>
          ) : (
            <div className="template-form__variables">
              {formData.variables.map((variable, index) => (
                <div key={index} className="template-form__variable">
                  <div className="template-form__variable-header">
                    <h4>Variable {index + 1}</h4>
                    <button
                      type="button"
                      onClick={() => removeVariable(index)}
                      className="template-form__remove-button"
                    >
                      Remove
                    </button>
                  </div>

                  <div className="template-form__row">
                    <div className="template-form__field">
                      <label>Name *</label>
                      <input
                        type="text"
                        value={variable.name}
                        onChange={(e) => updateVariable(index, 'name', e.target.value)}
                        className={errors[`variable_${index}_name`] ? 'template-form__input--error' : ''}
                        placeholder="Variable name"
                      />
                      {errors[`variable_${index}_name`] && (
                        <span className="template-form__error">{errors[`variable_${index}_name`]}</span>
                      )}
                    </div>

                    <div className="template-form__field">
                      <label>Type *</label>
                      <select
                        value={variable.type}
                        onChange={(e) => updateVariable(index, 'type', e.target.value)}
                      >
                        <option value="string">String</option>
                        <option value="int">Integer</option>
                        <option value="bool">Boolean</option>
                        <option value="float">Float</option>
                        <option value="array">Array</option>
                        <option value="object">Object</option>
                        <option value="enum">Enum</option>
                      </select>
                    </div>
                  </div>

                  <div className="template-form__field">
                    <label>Description</label>
                    <input
                      type="text"
                      value={variable.description}
                      onChange={(e) => updateVariable(index, 'description', e.target.value)}
                      placeholder="Variable description"
                    />
                  </div>

                  <div className="template-form__row">
                    <div className="template-form__field">
                      <label>Default Value</label>
                      <input
                        type="text"
                        value={variable.defaultValue}
                        onChange={(e) => updateVariable(index, 'defaultValue', e.target.value)}
                        placeholder="Default value"
                      />
                    </div>

                    <div className="template-form__field template-form__field--checkbox">
                      <label>
                        <input
                          type="checkbox"
                          checked={variable.required}
                          onChange={(e) => updateVariable(index, 'required', e.target.checked)}
                        />
                        Required
                      </label>
                    </div>
                  </div>

                  {variable.type === 'enum' && (
                    <div className="template-form__field">
                      <label>Options (comma-separated)</label>
                      <input
                        type="text"
                        value={formatOptions(variable.options || [])}
                        onChange={(e) => updateVariable(index, 'options', parseOptions(e.target.value))}
                        placeholder="option1, option2, option3"
                      />
                    </div>
                  )}

                  <div className="template-form__field">
                    <label>Validation (regex pattern)</label>
                    <input
                      type="text"
                      value={variable.validation}
                      onChange={(e) => updateVariable(index, 'validation', e.target.value)}
                      placeholder="^[a-zA-Z0-9]+$"
                    />
                  </div>
                </div>
              ))}
            </div>
          )}
        </div>

        <div className="template-form__actions">
          <button
            type="button"
            onClick={onCancel}
            className="template-form__button template-form__button--secondary"
            disabled={saving}
          >
            Cancel
          </button>
          <button
            type="submit"
            className="template-form__button template-form__button--primary"
            disabled={saving}
          >
            {saving ? 'Saving...' : (isEditing ? 'Update Template' : 'Create Template')}
          </button>
        </div>
      </form>
    </div>
  );
};

export default TemplateForm;