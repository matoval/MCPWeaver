import React, { useState, useEffect } from 'react';
import { 
  Plus, 
  Search, 
  Filter, 
  Download, 
  Upload, 
  Trash2, 
  Edit3, 
  Play, 
  Clock, 
  CheckCircle,
  AlertCircle,
  MoreVertical,
  FileText,
  Folder
} from 'lucide-react';
import { GetProjects, SearchProjects, DeleteProject, ExportProject, ImportProject, GetRecentProjects } from '../../wailsjs/go/app/App';
import { app } from '../../wailsjs/go/models';
import './ProjectList.scss';

interface ProjectListProps {
  onNewProject: () => void;
  onEditProject: (project: app.Project) => void;
  onOpenProject: (project: app.Project) => void;
  onGenerateProject: (project: app.Project) => void;
}

const ProjectList: React.FC<ProjectListProps> = ({ 
  onNewProject, 
  onEditProject, 
  onOpenProject, 
  onGenerateProject 
}) => {
  const [projects, setProjects] = useState<app.Project[]>([]);
  const [recentProjects, setRecentProjects] = useState<app.Project[]>([]);
  const [loading, setLoading] = useState(false);
  const [searchTerm, setSearchTerm] = useState('');
  const [statusFilter, setStatusFilter] = useState<string>('all');
  const [sortBy, setSortBy] = useState<'name' | 'updated' | 'created'>('updated');
  const [sortOrder, setSortOrder] = useState<'asc' | 'desc'>('desc');
  const [showFilters, setShowFilters] = useState(false);
  const [selectedProjects, setSelectedProjects] = useState<string[]>([]);

  // Load projects on component mount
  useEffect(() => {
    loadProjects();
    loadRecentProjects();
  }, []);

  // Search projects when search term changes
  useEffect(() => {
    if (searchTerm) {
      searchProjects();
    } else {
      loadProjects();
    }
  }, [searchTerm]);

  const loadProjects = async () => {
    setLoading(true);
    try {
      const projectsData = await GetProjects();
      setProjects(projectsData || []);
    } catch (error) {
      console.error('Failed to load projects:', error);
    } finally {
      setLoading(false);
    }
  };

  const loadRecentProjects = async () => {
    try {
      const recentData = await GetRecentProjects();
      setRecentProjects(recentData || []);
    } catch (error) {
      console.error('Failed to load recent projects:', error);
    }
  };

  const searchProjects = async () => {
    setLoading(true);
    try {
      const searchResults = await SearchProjects(searchTerm);
      setProjects(searchResults || []);
    } catch (error) {
      console.error('Failed to search projects:', error);
    } finally {
      setLoading(false);
    }
  };

  const handleDeleteProject = async (projectId: string) => {
    if (window.confirm('Are you sure you want to delete this project? This action cannot be undone.')) {
      try {
        await DeleteProject(projectId);
        loadProjects();
        loadRecentProjects();
      } catch (error) {
        console.error('Failed to delete project:', error);
      }
    }
  };

  const handleExportProject = async (projectId: string) => {
    try {
      const exportData = await ExportProject(projectId);
      const blob = new Blob([exportData], { type: 'application/json' });
      const url = URL.createObjectURL(blob);
      const a = document.createElement('a');
      a.href = url;
      a.download = `project-${projectId}.json`;
      a.click();
      URL.revokeObjectURL(url);
    } catch (error) {
      console.error('Failed to export project:', error);
    }
  };

  const handleImportProject = async () => {
    const input = document.createElement('input');
    input.type = 'file';
    input.accept = '.json';
    input.onchange = async (e) => {
      const file = (e.target as HTMLInputElement).files?.[0];
      if (file) {
        try {
          const text = await file.text();
          await ImportProject(text);
          loadProjects();
          loadRecentProjects();
        } catch (error) {
          console.error('Failed to import project:', error);
        }
      }
    };
    input.click();
  };

  const handleBulkDelete = async () => {
    if (selectedProjects.length === 0) return;
    
    const confirmMessage = `Are you sure you want to delete ${selectedProjects.length} project(s)? This action cannot be undone.`;
    if (window.confirm(confirmMessage)) {
      try {
        await Promise.all(selectedProjects.map(id => DeleteProject(id)));
        setSelectedProjects([]);
        loadProjects();
        loadRecentProjects();
      } catch (error) {
        console.error('Failed to delete projects:', error);
      }
    }
  };

  const toggleProjectSelection = (projectId: string) => {
    setSelectedProjects(prev => 
      prev.includes(projectId) 
        ? prev.filter(id => id !== projectId)
        : [...prev, projectId]
    );
  };

  const selectAllProjects = () => {
    setSelectedProjects(projects.map(p => p.id));
  };

  const clearSelection = () => {
    setSelectedProjects([]);
  };

  const getStatusIcon = (status: string) => {
    switch (status) {
      case 'ready':
        return <CheckCircle className="status-icon ready" size={16} />;
      case 'generating':
        return <Play className="status-icon generating" size={16} />;
      case 'error':
        return <AlertCircle className="status-icon error" size={16} />;
      default:
        return <Clock className="status-icon" size={16} />;
    }
  };

  const getStatusText = (status: string) => {
    switch (status) {
      case 'ready':
        return 'Ready';
      case 'generating':
        return 'Generating';
      case 'error':
        return 'Error';
      case 'validating':
        return 'Validating';
      default:
        return 'Created';
    }
  };

  const filteredProjects = projects.filter(project => {
    if (statusFilter !== 'all' && project.status !== statusFilter) return false;
    return true;
  });

  const sortedProjects = [...filteredProjects].sort((a, b) => {
    let aValue: any, bValue: any;
    
    switch (sortBy) {
      case 'name':
        aValue = a.name.toLowerCase();
        bValue = b.name.toLowerCase();
        break;
      case 'updated':
        aValue = new Date(a.updatedAt);
        bValue = new Date(b.updatedAt);
        break;
      case 'created':
        aValue = new Date(a.createdAt);
        bValue = new Date(b.createdAt);
        break;
      default:
        return 0;
    }
    
    if (sortOrder === 'asc') {
      return aValue > bValue ? 1 : -1;
    } else {
      return aValue < bValue ? 1 : -1;
    }
  });

  const formatDate = (dateString: string) => {
    return new Date(dateString).toLocaleDateString('en-US', {
      year: 'numeric',
      month: 'short',
      day: 'numeric',
      hour: '2-digit',
      minute: '2-digit'
    });
  };

  return (
    <div className="project-list">
      <div className="project-list-header">
        <div className="header-left">
          <h1>Projects</h1>
          <span className="project-count">{projects.length} projects</span>
        </div>
        <div className="header-right">
          <div className="search-bar">
            <Search className="search-icon" size={20} />
            <input
              type="text"
              placeholder="Search projects..."
              value={searchTerm}
              onChange={(e) => setSearchTerm(e.target.value)}
            />
          </div>
          <button 
            className="filter-button"
            onClick={() => setShowFilters(!showFilters)}
          >
            <Filter size={20} />
            Filters
          </button>
          <button className="import-button" onClick={handleImportProject}>
            <Upload size={20} />
            Import
          </button>
          <button className="new-project-button" onClick={onNewProject}>
            <Plus size={20} />
            New Project
          </button>
        </div>
      </div>

      {showFilters && (
        <div className="filters-panel">
          <div className="filter-group">
            <label>Status:</label>
            <select 
              value={statusFilter} 
              onChange={(e) => setStatusFilter(e.target.value)}
            >
              <option value="all">All Status</option>
              <option value="created">Created</option>
              <option value="ready">Ready</option>
              <option value="generating">Generating</option>
              <option value="error">Error</option>
            </select>
          </div>
          <div className="filter-group">
            <label>Sort by:</label>
            <select 
              value={sortBy} 
              onChange={(e) => setSortBy(e.target.value as 'name' | 'updated' | 'created')}
            >
              <option value="updated">Last Updated</option>
              <option value="created">Created Date</option>
              <option value="name">Name</option>
            </select>
          </div>
          <div className="filter-group">
            <label>Order:</label>
            <select 
              value={sortOrder} 
              onChange={(e) => setSortOrder(e.target.value as 'asc' | 'desc')}
            >
              <option value="desc">Descending</option>
              <option value="asc">Ascending</option>
            </select>
          </div>
        </div>
      )}

      {selectedProjects.length > 0 && (
        <div className="bulk-actions">
          <span>{selectedProjects.length} selected</span>
          <button onClick={handleBulkDelete} className="bulk-delete">
            <Trash2 size={16} />
            Delete Selected
          </button>
          <button onClick={clearSelection} className="clear-selection">
            Clear Selection
          </button>
        </div>
      )}

      {recentProjects.length > 0 && (
        <div className="recent-projects-section">
          <h2>Recent Projects</h2>
          <div className="recent-projects-grid">
            {recentProjects.slice(0, 6).map(project => (
              <div 
                key={project.id} 
                className="recent-project-card"
                onClick={() => onOpenProject(project)}
              >
                <div className="project-icon">
                  <FileText size={24} />
                </div>
                <h3>{project.name}</h3>
                <p className="project-path">{project.outputPath}</p>
                <div className="project-status">
                  {getStatusIcon(project.status)}
                  <span>{getStatusText(project.status)}</span>
                </div>
              </div>
            ))}
          </div>
        </div>
      )}

      <div className="all-projects-section">
        <div className="section-header">
          <h2>All Projects</h2>
          {projects.length > 0 && (
            <div className="select-actions">
              <button onClick={selectAllProjects}>Select All</button>
              <button onClick={clearSelection}>Clear</button>
            </div>
          )}
        </div>

        {loading ? (
          <div className="loading-state">
            <div className="spinner"></div>
            <p>Loading projects...</p>
          </div>
        ) : sortedProjects.length === 0 ? (
          <div className="empty-state">
            <FileText size={48} className="empty-icon" />
            <h3>No projects found</h3>
            <p>
              {searchTerm || statusFilter !== 'all' 
                ? 'No projects match your current filters' 
                : 'Create your first project to get started'}
            </p>
            <button onClick={onNewProject} className="empty-action">
              <Plus size={20} />
              Create New Project
            </button>
          </div>
        ) : (
          <div className="projects-grid">
            {sortedProjects.map(project => (
              <div 
                key={project.id} 
                className={`project-card ${selectedProjects.includes(project.id) ? 'selected' : ''}`}
              >
                <div className="project-card-header">
                  <input
                    type="checkbox"
                    checked={selectedProjects.includes(project.id)}
                    onChange={() => toggleProjectSelection(project.id)}
                    onClick={(e) => e.stopPropagation()}
                  />
                  <div className="project-status">
                    {getStatusIcon(project.status)}
                    <span>{getStatusText(project.status)}</span>
                  </div>
                  <div className="project-actions">
                    <button onClick={() => handleExportProject(project.id)}>
                      <Download size={16} />
                    </button>
                    <button onClick={() => onEditProject(project)}>
                      <Edit3 size={16} />
                    </button>
                    <button onClick={() => handleDeleteProject(project.id)}>
                      <Trash2 size={16} />
                    </button>
                  </div>
                </div>

                <div className="project-card-content" onClick={() => onOpenProject(project)}>
                  <div className="project-icon">
                    <FileText size={32} />
                  </div>
                  <h3>{project.name}</h3>
                  <p className="project-description">
                    {project.specPath && (
                      <span className="spec-source">
                        <FileText size={14} />
                        {project.specPath.split('/').pop()}
                      </span>
                    )}
                    {project.specUrl && (
                      <span className="spec-source">
                        <FileText size={14} />
                        {project.specUrl}
                      </span>
                    )}
                  </p>
                  <div className="project-path">
                    <Folder size={14} />
                    <span>{project.outputPath}</span>
                  </div>
                </div>

                <div className="project-card-footer">
                  <div className="project-stats">
                    <span className="generation-count">
                      {project.generationCount} generation{project.generationCount !== 1 ? 's' : ''}
                    </span>
                    <span className="last-updated">
                      Updated {formatDate(project.updatedAt)}
                    </span>
                  </div>
                  <div className="project-actions-footer">
                    <button 
                      className="generate-button"
                      onClick={(e) => {
                        e.stopPropagation();
                        onGenerateProject(project);
                      }}
                      disabled={project.status === 'generating'}
                    >
                      <Play size={16} />
                      Generate
                    </button>
                  </div>
                </div>
              </div>
            ))}
          </div>
        )}
      </div>
    </div>
  );
};

export default ProjectList;