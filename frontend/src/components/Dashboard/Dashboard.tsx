import React, { useState } from 'react';
import { Plus, FolderOpen, Download } from 'lucide-react';
import ProjectList from '../Project/ProjectList';
import ProjectForm from '../Project/ProjectForm';
import { app } from '../../wailsjs/go/models';
import './Dashboard.scss';

const Dashboard: React.FC = () => {
  const [showProjectForm, setShowProjectForm] = useState(false);
  const [editingProject, setEditingProject] = useState<app.Project | undefined>();
  const [currentView, setCurrentView] = useState<'dashboard' | 'projects'>('dashboard');

  const handleNewProject = () => {
    setEditingProject(undefined);
    setShowProjectForm(true);
  };

  const handleEditProject = (project: app.Project) => {
    setEditingProject(project);
    setShowProjectForm(true);
  };

  const handleOpenProject = (project: app.Project) => {
    // TODO: Navigate to project view
    console.log('Opening project:', project.name);
  };

  const handleGenerateProject = (project: app.Project) => {
    // TODO: Start generation process
    console.log('Generating project:', project.name);
  };

  const handleSaveProject = (project: app.Project) => {
    setShowProjectForm(false);
    setEditingProject(undefined);
    // Refresh project list if needed
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
                <button className="action-button primary" onClick={handleNewProject}>
                  Get Started
                </button>
              </div>
              
              <div className="action-card">
                <FolderOpen className="action-icon" size={24} />
                <h3>Browse Projects</h3>
                <p>View and manage all your existing projects</p>
                <button className="action-button secondary" onClick={handleBrowseProjects}>
                  Browse
                </button>
              </div>
              
              <div className="action-card">
                <Download className="action-icon" size={24} />
                <h3>Import Project</h3>
                <p>Import a previously exported project</p>
                <button className="action-button secondary">Import</button>
              </div>
            </div>
          </div>

          <div className="recent-projects">
            <h2>Recent Projects</h2>
            <div className="project-list">
              <div className="empty-state">
                <p>No recent projects</p>
                <span>Your recent projects will appear here</span>
                <button className="empty-action" onClick={handleBrowseProjects}>
                  View All Projects
                </button>
              </div>
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