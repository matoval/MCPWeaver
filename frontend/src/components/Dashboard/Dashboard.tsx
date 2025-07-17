import React from 'react';
import { Plus, FolderOpen, Download } from 'lucide-react';
import './Dashboard.scss';

const Dashboard: React.FC = () => {
  return (
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
              <button className="action-button primary">Get Started</button>
            </div>
            
            <div className="action-card">
              <FolderOpen className="action-icon" size={24} />
              <h3>Open Project</h3>
              <p>Continue working on an existing project</p>
              <button className="action-button secondary">Browse</button>
            </div>
            
            <div className="action-card">
              <Download className="action-icon" size={24} />
              <h3>Import OpenAPI</h3>
              <p>Import an OpenAPI specification file</p>
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
            </div>
          </div>
        </div>
      </div>
    </div>
  );
};

export default Dashboard;