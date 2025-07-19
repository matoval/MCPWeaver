import React, { useState, useEffect } from 'react';
import { TemplateMarketplaceItem } from '../../types';
import './MarketplaceDetails.scss';

interface MarketplaceDetailsProps {
  template: TemplateMarketplaceItem;
  onInstall: (template: TemplateMarketplaceItem) => void;
  onClose: () => void;
}

const MarketplaceDetails: React.FC<MarketplaceDetailsProps> = ({
  template,
  onInstall,
  onClose,
}) => {
  const [activeTab, setActiveTab] = useState<'overview' | 'variables' | 'dependencies' | 'screenshots'>('overview');
  const [installing, setInstalling] = useState(false);

  const handleInstall = async () => {
    setInstalling(true);
    try {
      await onInstall(template);
    } catch (err) {
      alert(err instanceof Error ? err.message : 'Installation failed');
    } finally {
      setInstalling(false);
    }
  };

  const formatDate = (dateString: string) => {
    return new Date(dateString).toLocaleDateString('en-US', {
      year: 'numeric',
      month: 'long',
      day: 'numeric',
    });
  };

  const formatFileSize = (bytes: number) => {
    const sizes = ['Bytes', 'KB', 'MB', 'GB'];
    if (bytes === 0) return '0 Bytes';
    const i = Math.floor(Math.log(bytes) / Math.log(1024));
    return Math.round(bytes / Math.pow(1024, i) * 100) / 100 + ' ' + sizes[i];
  };

  const renderOverview = () => (
    <div className="marketplace-details__overview">
      <div className="marketplace-details__description-section">
        <h3>Description</h3>
        <p>{template.description}</p>
      </div>

      <div className="marketplace-details__info-grid">
        <div className="marketplace-details__info-item">
          <label>Author</label>
          <p>{template.author}</p>
        </div>
        <div className="marketplace-details__info-item">
          <label>Version</label>
          <p>{template.version}</p>
        </div>
        <div className="marketplace-details__info-item">
          <label>Type</label>
          <p className={`marketplace-details__type marketplace-details__type--${template.type}`}>
            {template.type}
          </p>
        </div>
        <div className="marketplace-details__info-item">
          <label>License</label>
          <p>{template.license}</p>
        </div>
        <div className="marketplace-details__info-item">
          <label>Downloads</label>
          <p>{template.downloads.toLocaleString()}</p>
        </div>
        <div className="marketplace-details__info-item">
          <label>Rating</label>
          <p>
            <span className="marketplace-details__rating">
              ⭐ {template.rating.toFixed(1)}
            </span>
          </p>
        </div>
        <div className="marketplace-details__info-item">
          <label>Size</label>
          <p>{formatFileSize(template.size)}</p>
        </div>
        <div className="marketplace-details__info-item">
          <label>Created</label>
          <p>{formatDate(template.createdAt)}</p>
        </div>
        <div className="marketplace-details__info-item">
          <label>Updated</label>
          <p>{formatDate(template.updatedAt)}</p>
        </div>
      </div>

      {template.tags.length > 0 && (
        <div className="marketplace-details__tags-section">
          <h3>Tags</h3>
          <div className="marketplace-details__tags">
            {template.tags.map(tag => (
              <span key={tag} className="marketplace-details__tag">
                {tag}
              </span>
            ))}
          </div>
        </div>
      )}

      {(template.repository || template.homePage) && (
        <div className="marketplace-details__links-section">
          <h3>Links</h3>
          <div className="marketplace-details__links">
            {template.repository && (
              <a
                href={template.repository}
                target="_blank"
                rel="noopener noreferrer"
                className="marketplace-details__link"
              >
                📂 Repository
              </a>
            )}
            {template.homePage && (
              <a
                href={template.homePage}
                target="_blank"
                rel="noopener noreferrer"
                className="marketplace-details__link"
              >
                🏠 Home Page
              </a>
            )}
          </div>
        </div>
      )}
    </div>
  );

  const renderVariables = () => (
    <div className="marketplace-details__variables">
      {template.variables.length === 0 ? (
        <div className="marketplace-details__empty">
          <p>This template has no configurable variables.</p>
        </div>
      ) : (
        <div className="marketplace-details__variables-grid">
          {template.variables.map((variable, index) => (
            <div key={index} className="marketplace-details__variable">
              <div className="marketplace-details__variable-header">
                <h4>{variable.name}</h4>
                <div className="marketplace-details__variable-badges">
                  <span className={`marketplace-details__variable-type marketplace-details__variable-type--${variable.type}`}>
                    {variable.type}
                  </span>
                  {variable.required && (
                    <span className="marketplace-details__variable-required">
                      Required
                    </span>
                  )}
                </div>
              </div>
              
              <p className="marketplace-details__variable-description">
                {variable.description}
              </p>
              
              <div className="marketplace-details__variable-details">
                <div className="marketplace-details__variable-detail">
                  <span>Default:</span>
                  <code>{variable.defaultValue || 'None'}</code>
                </div>
                {variable.options && variable.options.length > 0 && (
                  <div className="marketplace-details__variable-detail">
                    <span>Options:</span>
                    <code>{variable.options.join(', ')}</code>
                  </div>
                )}
                {variable.validation && (
                  <div className="marketplace-details__variable-detail">
                    <span>Validation:</span>
                    <code>{variable.validation}</code>
                  </div>
                )}
              </div>
            </div>
          ))}
        </div>
      )}
    </div>
  );

  const renderDependencies = () => (
    <div className="marketplace-details__dependencies">
      {template.dependencies.length === 0 ? (
        <div className="marketplace-details__empty">
          <p>This template has no dependencies.</p>
        </div>
      ) : (
        <div className="marketplace-details__dependencies-grid">
          {template.dependencies.map((dependency, index) => (
            <div key={index} className="marketplace-details__dependency">
              <div className="marketplace-details__dependency-header">
                <h4>{dependency.name}</h4>
                <div className="marketplace-details__dependency-badges">
                  <span className="marketplace-details__dependency-version">
                    v{dependency.version}
                  </span>
                  <span className={`marketplace-details__dependency-type marketplace-details__dependency-type--${dependency.type}`}>
                    {dependency.type}
                  </span>
                  {dependency.required && (
                    <span className="marketplace-details__dependency-required">
                      Required
                    </span>
                  )}
                </div>
              </div>
            </div>
          ))}
        </div>
      )}
    </div>
  );

  const renderScreenshots = () => (
    <div className="marketplace-details__screenshots">
      {!template.screenshots || template.screenshots.length === 0 ? (
        <div className="marketplace-details__empty">
          <p>No screenshots available for this template.</p>
        </div>
      ) : (
        <div className="marketplace-details__screenshots-grid">
          {template.screenshots.map((screenshot, index) => (
            <div key={index} className="marketplace-details__screenshot">
              <img
                src={screenshot}
                alt={`Screenshot ${index + 1}`}
                loading="lazy"
              />
            </div>
          ))}
        </div>
      )}
    </div>
  );

  return (
    <div className="marketplace-details">
      <div className="marketplace-details__header">
        <div className="marketplace-details__title">
          <h1>{template.name}</h1>
          <div className="marketplace-details__badges">
            <span className={`marketplace-details__type marketplace-details__type--${template.type}`}>
              {template.type}
            </span>
            <span className="marketplace-details__version">v{template.version}</span>
            <span className="marketplace-details__rating">
              ⭐ {template.rating.toFixed(1)}
            </span>
          </div>
        </div>

        <div className="marketplace-details__actions">
          <button
            onClick={handleInstall}
            disabled={installing}
            className="marketplace-details__action marketplace-details__action--install"
          >
            {installing ? '⏳ Installing...' : '📥 Install Template'}
          </button>
          <button
            onClick={onClose}
            className="marketplace-details__action marketplace-details__action--close"
          >
            ✕ Close
          </button>
        </div>
      </div>

      <div className="marketplace-details__tabs">
        <button
          onClick={() => setActiveTab('overview')}
          className={`marketplace-details__tab ${
            activeTab === 'overview' ? 'marketplace-details__tab--active' : ''
          }`}
        >
          📋 Overview
        </button>
        <button
          onClick={() => setActiveTab('variables')}
          className={`marketplace-details__tab ${
            activeTab === 'variables' ? 'marketplace-details__tab--active' : ''
          }`}
        >
          ⚙️ Variables ({template.variables.length})
        </button>
        <button
          onClick={() => setActiveTab('dependencies')}
          className={`marketplace-details__tab ${
            activeTab === 'dependencies' ? 'marketplace-details__tab--active' : ''
          }`}
        >
          📦 Dependencies ({template.dependencies.length})
        </button>
        <button
          onClick={() => setActiveTab('screenshots')}
          className={`marketplace-details__tab ${
            activeTab === 'screenshots' ? 'marketplace-details__tab--active' : ''
          }`}
        >
          🖼️ Screenshots ({template.screenshots?.length || 0})
        </button>
      </div>

      <div className="marketplace-details__content">
        {activeTab === 'overview' && renderOverview()}
        {activeTab === 'variables' && renderVariables()}
        {activeTab === 'dependencies' && renderDependencies()}
        {activeTab === 'screenshots' && renderScreenshots()}
      </div>
    </div>
  );
};

export default MarketplaceDetails;