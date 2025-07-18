import React, { useState, useEffect } from 'react';
import { useParams } from 'react-router-dom';
import { wails } from '../../services/wails';
import { ProgressTracker } from '../ui/ProgressTracker';
import { Project } from '../../types';
import './ProjectView.scss';

const ProjectView: React.FC = () => {
  const { id } = useParams<{ id: string }>();
  const [project, setProject] = useState<Project | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  // Load project data
  useEffect(() => {
    const loadProject = async () => {
      if (!id) return;
      
      try {
        setLoading(true);
        setError(null);
        const projectData = await wails.getProjectById(id);
        setProject(projectData);
      } catch (err: any) {
        setError(err.message || 'Failed to load project');
        console.error('Failed to load project:', err);
      } finally {
        setLoading(false);
      }
    };

    loadProject();
  }, [id]);

  // Handle generation completion
  const handleGenerationComplete = () => {
    // Reload project to get updated stats
    if (id) {
      wails.getProjectById(id).then(setProject).catch(console.error);
    }
  };

  if (loading) {
    return (
      <div className="project-view project-view--loading">
        <div className="project-view__loading">
          <div className="project-view__loading-spinner">
            <svg viewBox="0 0 24 24">
              <circle cx="12" cy="12" r="10" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeDasharray="31.416" strokeDashoffset="31.416">
                <animate attributeName="stroke-dasharray" dur="2s" values="0 31.416;15.708 15.708;0 31.416" repeatCount="indefinite"/>
                <animate attributeName="stroke-dashoffset" dur="2s" values="0;-15.708;-31.416" repeatCount="indefinite"/>
              </circle>
            </svg>
          </div>
          <p>Loading project...</p>
        </div>
      </div>
    );
  }

  if (error) {
    return (
      <div className="project-view project-view--error">
        <div className="project-view__error">
          <div className="project-view__error-icon">
            <svg viewBox="0 0 20 20" fill="currentColor">
              <path fillRule="evenodd" d="M18 10a8 8 0 11-16 0 8 8 0 0116 0zm-7 4a1 1 0 11-2 0 1 1 0 012 0zm-1-9a1 1 0 00-1 1v4a1 1 0 102 0V6a1 1 0 00-1-1z" clipRule="evenodd" />
            </svg>
          </div>
          <h2>Error Loading Project</h2>
          <p>{error}</p>
          <button 
            className="project-view__error-retry"
            onClick={() => window.location.reload()}
          >
            Retry
          </button>
        </div>
      </div>
    );
  }

  if (!project) {
    return (
      <div className="project-view project-view--not-found">
        <div className="project-view__not-found">
          <div className="project-view__not-found-icon">
            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z" />
            </svg>
          </div>
          <h2>Project Not Found</h2>
          <p>The project with ID "{id}" could not be found.</p>
        </div>
      </div>
    );
  }

  return (
    <div className="project-view">
      <div className="project-view__header">
        <div className="project-view__title-section">
          <h1 className="project-view__title">{project.name}</h1>
          <div className="project-view__meta">
            <span className={`project-view__status project-view__status--${project.status}`}>
              {project.status}
            </span>
            <span className="project-view__separator">•</span>
            <span className="project-view__date">
              Last updated: {new Date(project.updatedAt).toLocaleDateString()}
            </span>
            {project.lastGenerated && (
              <>
                <span className="project-view__separator">•</span>
                <span className="project-view__generated">
                  Generated: {new Date(project.lastGenerated).toLocaleDateString()}
                </span>
              </>
            )}
          </div>
        </div>
        
        <div className="project-view__stats">
          <div className="project-view__stat">
            <div className="project-view__stat-value">{project.generationCount}</div>
            <div className="project-view__stat-label">Generations</div>
          </div>
        </div>
      </div>

      <div className="project-view__content">
        <div className="project-view__main">
          <div className="project-view__section">
            <h2>Progress Tracking</h2>
            <ProgressTracker
              projectId={project.id}
              showMetrics={true}
              showHistory={true}
              onComplete={handleGenerationComplete}
              className="project-view__progress-tracker"
            />
          </div>

          <div className="project-view__section">
            <h2>Project Details</h2>
            <div className="project-view__details">
              <div className="project-view__detail-group">
                <h3>Source</h3>
                <div className="project-view__detail-item">
                  <span className="project-view__detail-label">Spec Path:</span>
                  <span className="project-view__detail-value">{project.specPath || 'Not specified'}</span>
                </div>
                <div className="project-view__detail-item">
                  <span className="project-view__detail-label">Spec URL:</span>
                  <span className="project-view__detail-value">{project.specUrl || 'Not specified'}</span>
                </div>
              </div>

              <div className="project-view__detail-group">
                <h3>Output</h3>
                <div className="project-view__detail-item">
                  <span className="project-view__detail-label">Output Path:</span>
                  <span className="project-view__detail-value">{project.outputPath}</span>
                </div>
              </div>

              <div className="project-view__detail-group">
                <h3>Settings</h3>
                <div className="project-view__detail-item">
                  <span className="project-view__detail-label">Package Name:</span>
                  <span className="project-view__detail-value">{project.settings.packageName}</span>
                </div>
                <div className="project-view__detail-item">
                  <span className="project-view__detail-label">Server Port:</span>
                  <span className="project-view__detail-value">{project.settings.serverPort}</span>
                </div>
                <div className="project-view__detail-item">
                  <span className="project-view__detail-label">Log Level:</span>
                  <span className="project-view__detail-value">{project.settings.logLevel}</span>
                </div>
                <div className="project-view__detail-item">
                  <span className="project-view__detail-label">Logging Enabled:</span>
                  <span className="project-view__detail-value">
                    {project.settings.enableLogging ? 'Yes' : 'No'}
                  </span>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
};

export default ProjectView;