import React, { useState, useEffect } from 'react';
import { Template, TemplateType } from '../../types';
import TemplateList from './TemplateList';
import TemplateForm from './TemplateForm';
import TemplateView from './TemplateView';
import './TemplateManager.scss';

type ViewMode = 'list' | 'create' | 'edit' | 'view';

const TemplateManager: React.FC = () => {
  const [viewMode, setViewMode] = useState<ViewMode>('list');
  const [selectedTemplate, setSelectedTemplate] = useState<Template | null>(null);
  const [searchQuery, setSearchQuery] = useState('');
  const [filterType, setFilterType] = useState<TemplateType | ''>('');
  const [templates, setTemplates] = useState<Template[]>([]);

  useEffect(() => {
    loadTemplates();
  }, []);

  const loadTemplates = async () => {
    try {
      const templatesData = await window.go.app.App.GetAllTemplates();
      setTemplates(templatesData);
    } catch (err) {
      console.error('Failed to load templates:', err);
    }
  };

  const handleCreateTemplate = () => {
    setSelectedTemplate(null);
    setViewMode('create');
  };

  const handleEditTemplate = (template: Template) => {
    setSelectedTemplate(template);
    setViewMode('edit');
  };

  const handleViewTemplate = (template: Template) => {
    setSelectedTemplate(template);
    setViewMode('view');
  };

  const handleDeleteTemplate = async (template: Template) => {
    if (template.isBuiltIn) {
      alert('Cannot delete built-in templates');
      return;
    }

    if (window.confirm(`Are you sure you want to delete template "${template.name}"?`)) {
      try {
        await window.go.app.App.DeleteTemplate(template.id);
        await loadTemplates();
        
        // If we're viewing the deleted template, go back to list
        if (selectedTemplate?.id === template.id) {
          setViewMode('list');
          setSelectedTemplate(null);
        }
      } catch (err) {
        alert(err instanceof Error ? err.message : 'Failed to delete template');
      }
    }
  };

  const handleDuplicateTemplate = async (template: Template) => {
    const newName = prompt(`Enter name for duplicate of "${template.name}":`, `${template.name} (Copy)`);
    
    if (newName && newName.trim()) {
      try {
        await window.go.app.App.DuplicateTemplate(template.id, newName.trim());
        await loadTemplates();
      } catch (err) {
        alert(err instanceof Error ? err.message : 'Failed to duplicate template');
      }
    }
  };

  const handleTemplateTest = (template: Template) => {
    setSelectedTemplate(template);
    setViewMode('view');
    // The TemplateView component will handle switching to the test tab
  };

  const handleSaveTemplate = async (template: Template) => {
    await loadTemplates();
    setViewMode('list');
    setSelectedTemplate(null);
  };

  const handleCancel = () => {
    setViewMode('list');
    setSelectedTemplate(null);
  };

  const handleImportTemplate = async () => {
    try {
      const filters = [
        { displayName: 'Template Files', pattern: '*.tmpl;*.template;*.zip', extensions: ['tmpl', 'template', 'zip'] },
        { displayName: 'All Files', pattern: '*', extensions: ['*'] },
      ];
      
      const filePath = await window.go.app.App.SelectFile(filters);
      if (filePath) {
        const importRequest = {
          source: 'file',
          path: filePath,
          options: {
            overwriteExisting: false,
            validateOnly: false,
            includeDependencies: true,
            targetType: 'custom' as TemplateType,
          },
        };

        await window.go.app.App.ImportTemplate(importRequest);
        await loadTemplates();
      }
    } catch (err) {
      alert(err instanceof Error ? err.message : 'Failed to import template');
    }
  };

  const handleExportTemplate = async (template: Template) => {
    try {
      const targetPath = await window.go.app.App.SaveFile(
        '',
        `${template.name}-v${template.version}.zip`,
        [{ displayName: 'ZIP Files', pattern: '*.zip', extensions: ['zip'] }]
      );

      if (targetPath) {
        const exportRequest = {
          templateId: template.id,
          format: 'zip',
          targetPath,
          options: {
            includeDocumentation: true,
            includeExamples: false,
            includeDependencies: true,
            minify: false,
          },
        };

        await window.go.app.App.ExportTemplate(exportRequest);
        alert(`Template exported successfully to ${targetPath}`);
      }
    } catch (err) {
      alert(err instanceof Error ? err.message : 'Failed to export template');
    }
  };

  const getFilteredTemplates = () => {
    let filtered = templates;

    if (searchQuery) {
      filtered = filtered.filter(template =>
        template.name.toLowerCase().includes(searchQuery.toLowerCase()) ||
        template.description.toLowerCase().includes(searchQuery.toLowerCase()) ||
        template.author.toLowerCase().includes(searchQuery.toLowerCase())
      );
    }

    if (filterType) {
      filtered = filtered.filter(template => template.type === filterType);
    }

    return filtered;
  };

  const renderHeader = () => (
    <div className="template-manager__header">
      <div className="template-manager__title">
        <h1>Template Manager</h1>
        <span className="template-manager__count">
          {getFilteredTemplates().length} of {templates.length} templates
        </span>
      </div>

      {viewMode === 'list' && (
        <div className="template-manager__toolbar">
          <div className="template-manager__search">
            <input
              type="text"
              placeholder="Search templates..."
              value={searchQuery}
              onChange={(e) => setSearchQuery(e.target.value)}
              className="template-manager__search-input"
            />
            <select
              value={filterType}
              onChange={(e) => setFilterType(e.target.value as TemplateType | '')}
              className="template-manager__filter"
            >
              <option value="">All Types</option>
              <option value="default">Default</option>
              <option value="custom">Custom</option>
              <option value="plugin">Plugin</option>
            </select>
          </div>

          <div className="template-manager__actions">
            <button
              onClick={handleImportTemplate}
              className="template-manager__action template-manager__action--secondary"
            >
              📥 Import
            </button>
            <button
              onClick={handleCreateTemplate}
              className="template-manager__action template-manager__action--primary"
            >
              ➕ Create Template
            </button>
          </div>
        </div>
      )}

      {(viewMode === 'create' || viewMode === 'edit') && (
        <div className="template-manager__breadcrumb">
          <button onClick={() => setViewMode('list')} className="template-manager__breadcrumb-link">
            Templates
          </button>
          <span className="template-manager__breadcrumb-separator">›</span>
          <span>{viewMode === 'create' ? 'Create Template' : `Edit ${selectedTemplate?.name}`}</span>
        </div>
      )}

      {viewMode === 'view' && selectedTemplate && (
        <div className="template-manager__breadcrumb">
          <button onClick={() => setViewMode('list')} className="template-manager__breadcrumb-link">
            Templates
          </button>
          <span className="template-manager__breadcrumb-separator">›</span>
          <span>{selectedTemplate.name}</span>
          <div className="template-manager__view-actions">
            <button
              onClick={() => handleExportTemplate(selectedTemplate)}
              className="template-manager__action template-manager__action--secondary"
            >
              📤 Export
            </button>
          </div>
        </div>
      )}
    </div>
  );

  const renderContent = () => {
    switch (viewMode) {
      case 'create':
        return (
          <TemplateForm
            onSave={handleSaveTemplate}
            onCancel={handleCancel}
            isEditing={false}
          />
        );

      case 'edit':
        return selectedTemplate ? (
          <TemplateForm
            template={selectedTemplate}
            onSave={handleSaveTemplate}
            onCancel={handleCancel}
            isEditing={true}
          />
        ) : null;

      case 'view':
        return selectedTemplate ? (
          <TemplateView
            template={selectedTemplate}
            onEdit={handleEditTemplate}
            onDelete={handleDeleteTemplate}
            onDuplicate={handleDuplicateTemplate}
            onClose={() => setViewMode('list')}
          />
        ) : null;

      case 'list':
      default:
        return (
          <TemplateList
            onSelectTemplate={handleViewTemplate}
            onEditTemplate={handleEditTemplate}
            onDeleteTemplate={handleDeleteTemplate}
            onDuplicateTemplate={handleDuplicateTemplate}
            onTestTemplate={handleTemplateTest}
            filterType={filterType as TemplateType}
            searchQuery={searchQuery}
          />
        );
    }
  };

  return (
    <div className="template-manager">
      {renderHeader()}
      <div className="template-manager__content">
        {renderContent()}
      </div>
    </div>
  );
};

export default TemplateManager;