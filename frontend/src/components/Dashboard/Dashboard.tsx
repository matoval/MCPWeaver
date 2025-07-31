import React, { useState, useEffect } from 'react';
import { Plus, FolderOpen, Download } from 'lucide-react';
import ProjectList from '../Project/ProjectList';
import ProjectForm from '../Project/ProjectForm';
import { app } from '../../../wailsjs/go/models';
import { SelectFile, ImportProject, GetRecentProjects, GenerateServer, GetDefaultOpenAPIFilters, ReadFile } from '../../../wailsjs/go/app/App';
import './Dashboard.scss';

const Dashboard: React.FC = () => {
  const [showProjectForm, setShowProjectForm] = useState(false);
  const [editingProject, setEditingProject] = useState<app.Project | undefined>();
  const [currentView, setCurrentView] = useState<'dashboard' | 'projects'>('dashboard');
  const [recentProjects, setRecentProjects] = useState<app.Project[]>([]);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  // Load recent projects on mount
  useEffect(() => {
    loadRecentProjects();
  }, []);

  const loadRecentProjects = async () => {
    try {
      const projects = await GetRecentProjects();
      setRecentProjects(projects || []);
    } catch (error) {
      console.error('Failed to load recent projects:', error);
    }
  };

  const handleNewProject = () => {
    setEditingProject(undefined);
    setShowProjectForm(true);
  };

  const handleEditProject = (project: app.Project) => {
    setEditingProject(project);
    setShowProjectForm(true);
  };

  const handleOpenProject = (project: app.Project) => {
    try {
      setError(null);
      // For now, we'll show the project form in edit mode as a "project view"
      // In a full implementation, we'd navigate to a dedicated project view page
      setEditingProject(project);
      setShowProjectForm(true);
      console.log('Opening project:', project.name);
    } catch (error) {
      setError('Failed to open project');
      console.error('Failed to open project:', error);
    }
  };

  const handleGenerateProject = async (project: app.Project) => {
    try {
      setError(null);
      setLoading(true);
      
      const job = await GenerateServer(project.id);
      console.log('Generation started:', job);
      
      // Show success message or navigate to generation progress view
      alert(`Generation started for project "${project.name}". Job ID: ${job.id}`);
      
    } catch (error: any) {
      setError(`Failed to start generation: ${error.message || error}`);
      console.error('Failed to generate project:', error);
    } finally {
      setLoading(false);
    }
  };

  const handleImportProject = async () => {
    try {
      setError(null);
      setLoading(true);

      // Get filters for project files (JSON format)
      const filters = [
        {
          displayName: 'Project Files',
          pattern: '*.json',
          extensions: ['.json']
        },
        {
          displayName: 'All Files',
          pattern: '*.*',
          extensions: ['*']
        }
      ];

      // Open file dialog
      const filePath = await SelectFile(filters);
      if (!filePath) {
        // User cancelled
        return;
      }

      // Read the file content
      const fileContent = await ReadFile(filePath);

      // Import the project
      const importedProject = await ImportProject(fileContent);
      
      // Refresh recent projects list
      await loadRecentProjects();
      
      // Show success message
      alert(`Project "${importedProject.name}" imported successfully!`);
      
    } catch (error: any) {
      setError(`Failed to import project: ${error.message || error}`);
      console.error('Failed to import project:', error);
    } finally {
      setLoading(false);
    }
  };

  const handleSaveProject = (project: app.Project) => {
    setShowProjectForm(false);
    setEditingProject(undefined);
    // Refresh recent projects list
    loadRecentProjects();
  };

  const handleCancelForm = () => {
    setShowProjectForm(false);
    setEditingProject(undefined);
  };

  const handleBrowseProjects = () => {
    setCurrentView('projects');
  };

  if (currentView === 'projects') {
    return (
      <>
        <ProjectList
          onNewProject={handleNewProject}
          onEditProject={handleEditProject}
          onOpenProject={handleOpenProject}
          onGenerateProject={handleGenerateProject}
        />
        {showProjectForm && (
          <ProjectForm
            project={editingProject}
            onSave={handleSaveProject}
            onCancel={handleCancelForm}
          />
        )}
      </>
    );
  }

  return (
    <>
      <div className="dashboard">
        <div className="dashboard-header">
          <h1>Welcome to MCPWeaver</h1>
          <p>Transform your OpenAPI specifications into Model Context Protocol (MCP) servers</p>
        </div>

        <div className="dashboard-content">
          <div className="quick-actions">
            <h2>Quick Actions</h2>
            <div className="action-cards">
              <div className="action-card">
                <Plus className="action-icon" size={24} />
                <h3>New Project</h3>
                <p>Create a new MCP server project from scratch</p>
                <button 
                  className="action-button primary" 
                  onClick={handleNewProject}
                  disabled={loading}
                >
                  {loading ? 'Loading...' : 'Get Started'}
                </button>
              </div>
              
              <div className="action-card">
                <FolderOpen className="action-icon" size={24} />
                <h3>Browse Projects</h3>
                <p>View and manage all your existing projects</p>
                <button 
                  className="action-button secondary" 
                  onClick={handleBrowseProjects}
                  disabled={loading}
                >
                  {loading ? 'Loading...' : 'Browse'}
                </button>
              </div>
              
              <div className="action-card">
                <Download className="action-icon" size={24} />
                <h3>Import Project</h3>
                <p>Import a previously exported project</p>
                <button 
                  className="action-button secondary" 
                  onClick={handleImportProject}
                  disabled={loading}
                >
                  {loading ? 'Importing...' : 'Import'}
                </button>
              </div>
            </div>
          </div>

          <div className="recent-projects">
            <h2>Recent Projects</h2>
            <div className="project-list">
              {error && (
                <div className="error-message">
                  <p>{error}</p>
                  <button onClick={() => setError(null)}>Dismiss</button>
                </div>
              )}
              
              {recentProjects.length === 0 ? (
                <div className="empty-state">
                  <p>No recent projects</p>
                  <span>Your recent projects will appear here</span>
                  <button className="empty-action" onClick={handleBrowseProjects}>
                    View All Projects
                  </button>
                </div>
              ) : (
                <div className="recent-project-cards">
                  {recentProjects.slice(0, 3).map((project) => (
                    <div key={project.id} className="recent-project-card">
                      <div className="project-info">
                        <h4>{project.name}</h4>
                        <p>{project.specPath ? `File: ${project.specPath.split('/').pop()}` : 
                            project.specURL ? `URL: ${project.specURL}` : 'No source'}</p>
                        <span className="project-status" data-status={project.status}>{project.status}</span>
                      </div>
                      <div className="project-actions">
                        <button 
                          className="action-btn primary"
                          onClick={() => handleOpenProject(project)}
                        >
                          Open
                        </button>
                        <button 
                          className="action-btn secondary"
                          onClick={() => handleGenerateProject(project)}
                          disabled={loading}
                        >
                          {loading ? 'Generating...' : 'Generate'}
                        </button>
                      </div>
                    </div>
                  ))}
                  {recentProjects.length > 3 && (
                    <button className="view-all-button" onClick={handleBrowseProjects}>
                      View All {recentProjects.length} Projects
                    </button>
                  )}
                </div>
              )}
            </div>
          </div>
        </div>
      </div>

      {showProjectForm && (
        <ProjectForm
          project={editingProject}
          onSave={handleSaveProject}
          onCancel={handleCancelForm}
        />
      )}
    </>
  );
};

export default Dashboard;