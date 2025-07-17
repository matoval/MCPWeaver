import React, { useState } from 'react';
import { 
  Folder, 
  Clock, 
  FileText, 
  Activity, 
  ChevronDown, 
  ChevronRight 
} from 'lucide-react';
import './Sidebar.scss';

interface SidebarSection {
  id: string;
  title: string;
  icon: React.ComponentType<any>;
  collapsible: boolean;
  defaultExpanded: boolean;
  content: React.ComponentType<any>;
}

interface SidebarProps {
  isOpen: boolean;
  onToggle: () => void;
}

// Placeholder components for sidebar content
const ProjectList: React.FC = () => (
  <div className="sidebar-content">
    <div className="sidebar-item">
      <span>No projects yet</span>
    </div>
  </div>
);

const RecentProjectsList: React.FC = () => (
  <div className="sidebar-content">
    <div className="sidebar-item">
      <span>No recent projects</span>
    </div>
  </div>
);

const TemplateList: React.FC = () => (
  <div className="sidebar-content">
    <div className="sidebar-item">
      <span>Default Template</span>
    </div>
  </div>
);

const ActivityLog: React.FC = () => (
  <div className="sidebar-content">
    <div className="sidebar-item">
      <span>No recent activity</span>
    </div>
  </div>
);

const Sidebar: React.FC<SidebarProps> = ({ isOpen, onToggle }) => {
  const [expandedSections, setExpandedSections] = useState<Set<string>>(
    new Set(['projects'])
  );

  const sidebarSections: SidebarSection[] = [
    {
      id: 'projects',
      title: 'Projects',
      icon: Folder,
      collapsible: true,
      defaultExpanded: true,
      content: ProjectList
    },
    {
      id: 'recent',
      title: 'Recent',
      icon: Clock,
      collapsible: true,
      defaultExpanded: false,
      content: RecentProjectsList
    },
    {
      id: 'templates',
      title: 'Templates',
      icon: FileText,
      collapsible: true,
      defaultExpanded: false,
      content: TemplateList
    },
    {
      id: 'activity',
      title: 'Activity',
      icon: Activity,
      collapsible: true,
      defaultExpanded: false,
      content: ActivityLog
    }
  ];

  const toggleSection = (sectionId: string) => {
    setExpandedSections(prev => {
      const newSet = new Set(prev);
      if (newSet.has(sectionId)) {
        newSet.delete(sectionId);
      } else {
        newSet.add(sectionId);
      }
      return newSet;
    });
  };

  return (
    <div className={`sidebar ${isOpen ? 'open' : 'closed'}`}>
      {sidebarSections.map((section) => {
        const isExpanded = expandedSections.has(section.id);
        const ContentComponent = section.content;
        const IconComponent = section.icon;
        
        return (
          <div key={section.id} className="sidebar-section">
            <button
              className="sidebar-section-header"
              onClick={() => section.collapsible && toggleSection(section.id)}
            >
              <div className="sidebar-section-title">
                <IconComponent className="sidebar-section-icon" size={16} />
                <span>{section.title}</span>
              </div>
              {section.collapsible && (
                <div className="sidebar-section-toggle">
                  {isExpanded ? (
                    <ChevronDown size={14} />
                  ) : (
                    <ChevronRight size={14} />
                  )}
                </div>
              )}
            </button>
            {isExpanded && (
              <div className="sidebar-section-content">
                <ContentComponent />
              </div>
            )}
          </div>
        );
      })}
    </div>
  );
};

export default Sidebar;