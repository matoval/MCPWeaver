import React from 'react';
import { 
  Plus, 
  FolderOpen, 
  Download, 
  CheckCircle, 
  Code, 
  Upload 
} from 'lucide-react';
import './Toolbar.scss';

interface ToolbarButton {
  id: string;
  label: string;
  icon: React.ComponentType<any>;
  tooltip: string;
  action: string;
  shortcut?: string;
  enabled: boolean;
  variant: 'primary' | 'secondary' | 'danger' | 'success';
}

interface ToolbarProps {
  onAction: (action: string) => void;
  projectLoaded?: boolean;
  projectValidated?: boolean;
  generationComplete?: boolean;
}

const Toolbar: React.FC<ToolbarProps> = ({ 
  onAction, 
  projectLoaded = false, 
  projectValidated = false, 
  generationComplete = false 
}) => {
  const toolbarButtons: ToolbarButton[] = [
    {
      id: 'new-project',
      label: 'New Project',
      icon: Plus,
      tooltip: 'Create a new project (Ctrl+N)',
      action: 'newProject',
      shortcut: 'Ctrl+N',
      enabled: true,
      variant: 'primary'
    },
    {
      id: 'open-project',
      label: 'Open',
      icon: FolderOpen,
      tooltip: 'Open existing project (Ctrl+O)',
      action: 'openProject',
      shortcut: 'Ctrl+O',
      enabled: true,
      variant: 'secondary'
    },
    {
      id: 'import-spec',
      label: 'Import',
      icon: Download,
      tooltip: 'Import OpenAPI specification (Ctrl+I)',
      action: 'importSpec',
      shortcut: 'Ctrl+I',
      enabled: true,
      variant: 'secondary'
    },
    {
      id: 'validate-spec',
      label: 'Validate',
      icon: CheckCircle,
      tooltip: 'Validate OpenAPI specification (F5)',
      action: 'validateSpec',
      shortcut: 'F5',
      enabled: projectLoaded,
      variant: 'secondary'
    },
    {
      id: 'generate-server',
      label: 'Generate',
      icon: Code,
      tooltip: 'Generate MCP server (F6)',
      action: 'generateServer',
      shortcut: 'F6',
      enabled: projectValidated,
      variant: 'success'
    },
    {
      id: 'export-server',
      label: 'Export',
      icon: Upload,
      tooltip: 'Export generated server (Ctrl+E)',
      action: 'exportServer',
      shortcut: 'Ctrl+E',
      enabled: generationComplete,
      variant: 'primary'
    }
  ];

  return (
    <div className="toolbar">
      <div className="toolbar-buttons">
        {toolbarButtons.map((button) => (
          <button
            key={button.id}
            className={`toolbar-button ${button.variant} ${!button.enabled ? 'disabled' : ''}`}
            onClick={() => button.enabled && onAction(button.action)}
            disabled={!button.enabled}
            title={button.tooltip}
          >
            <button.icon className="toolbar-button-icon" size={16} />
            <span className="toolbar-button-label">{button.label}</span>
          </button>
        ))}
      </div>
    </div>
  );
};

export default Toolbar;