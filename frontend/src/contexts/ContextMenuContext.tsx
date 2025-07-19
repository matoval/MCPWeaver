import React, { createContext, useContext, ReactNode } from 'react';
import useContextMenu, { ContextMenuItem } from '../hooks/useContextMenu';
import ContextMenu from '../components/ui/ContextMenu';

interface ContextMenuContextType {
  show: (event: React.MouseEvent | MouseEvent, items: ContextMenuItem[], targetElement?: Element) => void;
  hide: () => void;
  bindContextMenu: (element: HTMLElement | null, items: ContextMenuItem[] | (() => ContextMenuItem[])) => () => void;
  isOpen: boolean;
}

const ContextMenuContext = createContext<ContextMenuContextType | undefined>(undefined);

export const useContextMenuContext = () => {
  const context = useContext(ContextMenuContext);
  if (!context) {
    throw new Error('useContextMenuContext must be used within a ContextMenuProvider');
  }
  return context;
};

interface ContextMenuProviderProps {
  children: ReactNode;
}

export const ContextMenuProvider: React.FC<ContextMenuProviderProps> = ({ children }) => {
  const {
    isOpen,
    position,
    items,
    show,
    hide,
    bindContextMenu,
    getMenuItemProps,
    menuRef
  } = useContextMenu();

  const handleItemClick = (item: ContextMenuItem, event: React.MouseEvent) => {
    event.stopPropagation();
    
    if (item.disabled || item.separator) return;
    
    if (item.action) {
      item.action();
    }

    hide();
  };

  const contextValue: ContextMenuContextType = {
    show,
    hide,
    bindContextMenu,
    isOpen
  };

  return (
    <ContextMenuContext.Provider value={contextValue}>
      {children}
      
      <ContextMenu
        isOpen={isOpen}
        position={position}
        items={items}
        onClose={hide}
        onItemClick={handleItemClick}
        menuRef={menuRef}
      />
    </ContextMenuContext.Provider>
  );
};

// Pre-defined context menu items for common actions
export const createFileContextMenu = (
  onOpen?: () => void,
  onDelete?: () => void,
  onRename?: () => void,
  onCopy?: () => void,
  onProperties?: () => void
): ContextMenuItem[] => [
  {
    id: 'open',
    label: 'Open',
    icon: (
      <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor">
        <path d="M14 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V8z" />
        <polyline points="14,2 14,8 20,8" />
      </svg>
    ),
    shortcut: 'Enter',
    action: onOpen
  },
  {
    id: 'separator-1',
    separator: true
  },
  {
    id: 'copy',
    label: 'Copy',
    icon: (
      <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor">
        <rect x="9" y="9" width="13" height="13" rx="2" ry="2" />
        <path d="M5 15H4a2 2 0 0 1-2-2V4a2 2 0 0 1 2-2h9a2 2 0 0 1 2 2v1" />
      </svg>
    ),
    shortcut: 'Ctrl+C',
    action: onCopy
  },
  {
    id: 'rename',
    label: 'Rename',
    icon: (
      <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor">
        <path d="M17 3a2.828 2.828 0 1 1 4 4L7.5 20.5 2 22l1.5-5.5L17 3z" />
      </svg>
    ),
    shortcut: 'F2',
    action: onRename
  },
  {
    id: 'separator-2',
    separator: true
  },
  {
    id: 'delete',
    label: 'Delete',
    icon: (
      <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor">
        <polyline points="3,6 5,6 21,6" />
        <path d="M19 6v14a2 2 0 0 1-2 2H7a2 2 0 0 1-2-2V6m3 0V4a2 2 0 0 1 2-2h4a2 2 0 0 1 2 2v2" />
      </svg>
    ),
    shortcut: 'Delete',
    action: onDelete,
    danger: true
  },
  {
    id: 'separator-3',
    separator: true
  },
  {
    id: 'properties',
    label: 'Properties',
    icon: (
      <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor">
        <circle cx="12" cy="12" r="3" />
        <path d="M19.4 15a1.65 1.65 0 0 0 .33 1.82l.06.06a2 2 0 0 1 0 2.83 2 2 0 0 1-2.83 0l-.06-.06a1.65 1.65 0 0 0-1.82-.33 1.65 1.65 0 0 0-1 1.51V21a2 2 0 0 1-2 2 2 2 0 0 1-2-2v-.09A1.65 1.65 0 0 0 9 19.4a1.65 1.65 0 0 0-1.82.33l-.06.06a2 2 0 0 1-2.83 0 2 2 0 0 1 0-2.83l.06-.06a1.65 1.65 0 0 0 .33-1.82 1.65 1.65 0 0 0-1.51-1H3a2 2 0 0 1-2-2 2 2 0 0 1 2-2h.09A1.65 1.65 0 0 0 4.6 9a1.65 1.65 0 0 0-.33-1.82l-.06-.06a2 2 0 0 1 0-2.83 2 2 0 0 1 2.83 0l.06.06a1.65 1.65 0 0 0 1.82.33H9a1.65 1.65 0 0 0 1 1.51V3a2 2 0 0 1 2-2 2 2 0 0 1 2 2v.09a1.65 1.65 0 0 0 1 1.51 1.65 1.65 0 0 0 1.82-.33l.06-.06a2 2 0 0 1 2.83 0 2 2 0 0 1 0 2.83l-.06.06a1.65 1.65 0 0 0-.33 1.82V9a1.65 1.65 0 0 0 1.51 1H21a2 2 0 0 1 2 2 2 2 0 0 1-2 2h-.09a1.65 1.65 0 0 0-1.51 1z" />
      </svg>
    ),
    action: onProperties
  }
];

export const createProjectContextMenu = (
  onGenerate?: () => void,
  onValidate?: () => void,
  onExport?: () => void,
  onDuplicate?: () => void,
  onDelete?: () => void
): ContextMenuItem[] => [
  {
    id: 'generate',
    label: 'Generate MCP Server',
    icon: (
      <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor">
        <polyline points="16,18 22,12 16,6" />
        <polyline points="8,6 2,12 8,18" />
      </svg>
    ),
    shortcut: 'F5',
    action: onGenerate
  },
  {
    id: 'validate',
    label: 'Validate Specification',
    icon: (
      <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor">
        <polyline points="20,6 9,17 4,12" />
      </svg>
    ),
    shortcut: 'F9',
    action: onValidate
  },
  {
    id: 'export',
    label: 'Export Server',
    icon: (
      <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor">
        <path d="M21 15v4a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2v-4" />
        <polyline points="7,10 12,15 17,10" />
        <line x1="12" y1="15" x2="12" y2="3" />
      </svg>
    ),
    shortcut: 'Ctrl+E',
    action: onExport
  },
  {
    id: 'separator-1',
    separator: true
  },
  {
    id: 'duplicate',
    label: 'Duplicate Project',
    icon: (
      <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor">
        <rect x="9" y="9" width="13" height="13" rx="2" ry="2" />
        <path d="M5 15H4a2 2 0 0 1-2-2V4a2 2 0 0 1 2-2h9a2 2 0 0 1 2 2v1" />
      </svg>
    ),
    action: onDuplicate
  },
  {
    id: 'separator-2',
    separator: true
  },
  {
    id: 'delete',
    label: 'Delete Project',
    icon: (
      <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor">
        <polyline points="3,6 5,6 21,6" />
        <path d="M19 6v14a2 2 0 0 1-2 2H7a2 2 0 0 1-2-2V6m3 0V4a2 2 0 0 1 2-2h4a2 2 0 0 1 2 2v2" />
      </svg>
    ),
    action: onDelete,
    danger: true
  }
];

export default ContextMenuProvider;