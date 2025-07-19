import React, { useState, useEffect } from 'react';
import { Template, TemplateType } from '../../types';
import './TemplateList.scss';

interface TemplateListProps {
  onSelectTemplate?: (template: Template) => void;
  onEditTemplate?: (template: Template) => void;
  onDeleteTemplate?: (template: Template) => void;
  onDuplicateTemplate?: (template: Template) => void;
  onTestTemplate?: (template: Template) => void;
  filterType?: TemplateType;
  searchQuery?: string;
}

const TemplateList: React.FC<TemplateListProps> = ({
  onSelectTemplate,
  onEditTemplate,
  onDeleteTemplate,
  onDuplicateTemplate,
  onTestTemplate,
  filterType,
  searchQuery = '',
}) => {
  const [templates, setTemplates] = useState<Template[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [selectedTemplate, setSelectedTemplate] = useState<string | null>(null);

  useEffect(() => {
    loadTemplates();
  }, [filterType, searchQuery]);

  const loadTemplates = async () => {
    try {
      setLoading(true);
      setError(null);

      let templatesData: Template[];

      if (searchQuery) {
        // Use search if query provided
        templatesData = await window.go.app.App.SearchTemplates(searchQuery);
      } else if (filterType) {
        // Filter by type
        templatesData = await window.go.app.App.GetTemplatesByType(filterType);
      } else {
        // Get all templates
        templatesData = await window.go.app.App.GetAllTemplates();
      }

      setTemplates(templatesData);
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to load templates');
    } finally {
      setLoading(false);
    }
  };

  const handleSelectTemplate = (template: Template) => {
    setSelectedTemplate(template.id);
    onSelectTemplate?.(template);
  };

  const handleEdit = (template: Template, event: React.MouseEvent) => {
    event.stopPropagation();
    onEditTemplate?.(template);
  };

  const handleDelete = async (template: Template, event: React.MouseEvent) => {
    event.stopPropagation();
    
    if (template.isBuiltIn) {
      alert('Cannot delete built-in templates');
      return;
    }

    if (window.confirm(`Are you sure you want to delete template "${template.name}"?`)) {
      try {
        await window.go.app.App.DeleteTemplate(template.id);
        await loadTemplates(); // Refresh list
      } catch (err) {
        setError(err instanceof Error ? err.message : 'Failed to delete template');
      }
    }
  };

  const handleDuplicate = async (template: Template, event: React.MouseEvent) => {
    event.stopPropagation();
    const newName = prompt(`Enter name for duplicate of "${template.name}":`, `${template.name} (Copy)`);
    
    if (newName && newName.trim()) {
      try {
        await window.go.app.App.DuplicateTemplate(template.id, newName.trim());
        await loadTemplates(); // Refresh list
        onDuplicateTemplate?.(template);
      } catch (err) {
        setError(err instanceof Error ? err.message : 'Failed to duplicate template');
      }
    }
  };

  const handleTest = (template: Template, event: React.MouseEvent) => {
    event.stopPropagation();
    onTestTemplate?.(template);
  };

  const getTemplateTypeIcon = (type: TemplateType) => {
    switch (type) {
      case 'default':
        return '🔧';
      case 'custom':
        return '⚡';
      case 'plugin':
        return '🔌';
      default:
        return '📄';
    }
  };

  const getTemplateTypeLabel = (type: TemplateType) => {
    switch (type) {
      case 'default':
        return 'Default';
      case 'custom':
        return 'Custom';
      case 'plugin':
        return 'Plugin';
      default:
        return 'Unknown';
    }
  };

  if (loading) {
    return (
      <div className="template-list template-list--loading">
        <div className="template-list__spinner">Loading templates...</div>
      </div>
    );
  }

  if (error) {
    return (
      <div className="template-list template-list--error">
        <div className="template-list__error">
          <h3>Error loading templates</h3>
          <p>{error}</p>
          <button onClick={loadTemplates} className="button button--primary">
            Retry
          </button>
        </div>
      </div>
    );
  }

  if (templates.length === 0) {
    return (
      <div className="template-list template-list--empty">
        <div className="template-list__empty">
          <h3>No templates found</h3>
          <p>
            {searchQuery
              ? `No templates match "${searchQuery}"`
              : filterType
              ? `No ${filterType} templates available`
              : 'No templates available'}
          </p>
        </div>
      </div>
    );
  }

  return (
    <div className="template-list">
      <div className="template-list__header">
        <h2>Templates ({templates.length})</h2>
      </div>
      
      <div className="template-list__grid">
        {templates.map((template) => (
          <div
            key={template.id}
            className={`template-card ${
              selectedTemplate === template.id ? 'template-card--selected' : ''
            } ${template.isBuiltIn ? 'template-card--builtin' : ''}`}
            onClick={() => handleSelectTemplate(template)}
          >
            <div className="template-card__header">
              <div className="template-card__icon">
                {getTemplateTypeIcon(template.type)}
              </div>
              <div className="template-card__title">
                <h3>{template.name}</h3>
                <span className="template-card__version">v{template.version}</span>
              </div>
              <div className="template-card__actions">
                <button
                  onClick={(e) => handleTest(template, e)}
                  className="template-card__action"
                  title="Test template"
                >
                  🧪
                </button>
                <button
                  onClick={(e) => handleDuplicate(template, e)}
                  className="template-card__action"
                  title="Duplicate template"
                >
                  📋
                </button>
                {!template.isBuiltIn && (
                  <>
                    <button
                      onClick={(e) => handleEdit(template, e)}
                      className="template-card__action"
                      title="Edit template"
                    >
                      ✏️
                    </button>
                    <button
                      onClick={(e) => handleDelete(template, e)}
                      className="template-card__action template-card__action--danger"
                      title="Delete template"
                    >
                      🗑️
                    </button>
                  </>
                )}
              </div>
            </div>

            <div className="template-card__body">
              <p className="template-card__description">
                {template.description || 'No description available'}
              </p>
              
              <div className="template-card__metadata">
                <div className="template-card__meta-item">
                  <span className="template-card__meta-label">Type:</span>
                  <span className={`template-card__badge template-card__badge--${template.type}`}>
                    {getTemplateTypeLabel(template.type)}
                  </span>
                </div>
                
                <div className="template-card__meta-item">
                  <span className="template-card__meta-label">Author:</span>
                  <span>{template.author || 'Unknown'}</span>
                </div>
                
                {template.variables && template.variables.length > 0 && (
                  <div className="template-card__meta-item">
                    <span className="template-card__meta-label">Variables:</span>
                    <span>{template.variables.length}</span>
                  </div>
                )}
              </div>
            </div>

            <div className="template-card__footer">
              <div className="template-card__timestamps">
                <div className="template-card__timestamp">
                  Created: {new Date(template.createdAt).toLocaleDateString()}
                </div>
                {template.updatedAt !== template.createdAt && (
                  <div className="template-card__timestamp">
                    Updated: {new Date(template.updatedAt).toLocaleDateString()}
                  </div>
                )}
              </div>
            </div>
          </div>
        ))}
      </div>
    </div>
  );
};

export default TemplateList;