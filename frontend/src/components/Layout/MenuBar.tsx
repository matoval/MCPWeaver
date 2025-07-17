import React from 'react';
import './MenuBar.scss';

interface MenuBarProps {
  onMenuAction: (action: string) => void;
}

const MenuBar: React.FC<MenuBarProps> = ({ onMenuAction }) => {
  const menuItems = [
    {
      label: 'File',
      items: [
        { label: 'New Project', shortcut: 'Ctrl+N', action: 'newProject' },
        { label: 'Open Project', shortcut: 'Ctrl+O', action: 'openProject' },
        { label: 'Recent Projects', submenu: true },
        { separator: true },
        { label: 'Import OpenAPI Spec', shortcut: 'Ctrl+I', action: 'importSpec' },
        { label: 'Export Server', shortcut: 'Ctrl+E', action: 'exportServer' },
        { separator: true },
        { label: 'Settings', shortcut: 'Ctrl+,', action: 'openSettings' },
        { separator: true },
        { label: 'Exit', shortcut: 'Ctrl+Q', action: 'exit' }
      ]
    },
    {
      label: 'Edit',
      items: [
        { label: 'Undo', shortcut: 'Ctrl+Z', action: 'undo' },
        { label: 'Redo', shortcut: 'Ctrl+Y', action: 'redo' },
        { separator: true },
        { label: 'Cut', shortcut: 'Ctrl+X', action: 'cut' },
        { label: 'Copy', shortcut: 'Ctrl+C', action: 'copy' },
        { label: 'Paste', shortcut: 'Ctrl+V', action: 'paste' },
        { separator: true },
        { label: 'Find', shortcut: 'Ctrl+F', action: 'find' },
        { label: 'Replace', shortcut: 'Ctrl+H', action: 'replace' }
      ]
    },
    {
      label: 'View',
      items: [
        { label: 'Zoom In', shortcut: 'Ctrl++', action: 'zoomIn' },
        { label: 'Zoom Out', shortcut: 'Ctrl+-', action: 'zoomOut' },
        { label: 'Reset Zoom', shortcut: 'Ctrl+0', action: 'resetZoom' },
        { separator: true },
        { label: 'Toggle Sidebar', shortcut: 'Ctrl+B', action: 'toggleSidebar' },
        { label: 'Toggle Status Bar', action: 'toggleStatusBar' },
        { separator: true },
        { label: 'Activity Log', shortcut: 'Ctrl+L', action: 'showActivityLog' },
        { label: 'Performance Metrics', action: 'showMetrics' }
      ]
    },
    {
      label: 'Tools',
      items: [
        { label: 'Validate Spec', shortcut: 'F5', action: 'validateSpec' },
        { label: 'Generate Server', shortcut: 'F6', action: 'generateServer' },
        { label: 'Test Server', shortcut: 'F7', action: 'testServer' },
        { separator: true },
        { label: 'Template Manager', action: 'manageTemplates' },
        { label: 'Clear Cache', action: 'clearCache' }
      ]
    },
    {
      label: 'Help',
      items: [
        { label: 'User Guide', shortcut: 'F1', action: 'showUserGuide' },
        { label: 'API Documentation', action: 'showApiDocs' },
        { label: 'Keyboard Shortcuts', action: 'showShortcuts' },
        { separator: true },
        { label: 'Report Issue', action: 'reportIssue' },
        { label: 'About MCPWeaver', action: 'showAbout' }
      ]
    }
  ];

  return (
    <div className="menu-bar">
      {menuItems.map((menu, index) => (
        <div key={index} className="menu-item">
          <button 
            className="menu-button"
            onClick={() => onMenuAction(menu.label.toLowerCase())}
          >
            {menu.label}
          </button>
          <div className="menu-dropdown">
            {menu.items.map((item, itemIndex) => (
              item.separator ? (
                <div key={itemIndex} className="menu-separator" />
              ) : (
                <button
                  key={itemIndex}
                  className="menu-dropdown-item"
                  onClick={() => item.action && onMenuAction(item.action)}
                >
                  <span className="menu-item-label">{item.label}</span>
                  {item.shortcut && (
                    <span className="menu-item-shortcut">{item.shortcut}</span>
                  )}
                </button>
              )
            ))}
          </div>
        </div>
      ))}
    </div>
  );
};

export default MenuBar;