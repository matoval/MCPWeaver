import React from 'react';
import { ContextMenuItem } from '../../hooks/useContextMenu';

interface ContextMenuProps {
  isOpen: boolean;
  position: { x: number; y: number };
  items: ContextMenuItem[];
  onClose: () => void;
  onItemClick: (item: ContextMenuItem, event: React.MouseEvent) => void;
  menuRef: React.RefObject<HTMLDivElement>;
}

const ContextMenu: React.FC<ContextMenuProps> = ({
  isOpen,
  position,
  items,
  onClose,
  onItemClick,
  menuRef
}) => {
  if (!isOpen) return null;

  const renderMenuItem = (item: ContextMenuItem, index: number) => {
    if (item.separator) {
      return (
        <div
          key={`separator-${index}`}
          className="context-menu-separator"
          role="separator"
        />
      );
    }

    return (
      <div
        key={item.id || index}
        className={`context-menu-item ${item.disabled ? 'disabled' : ''} ${item.danger ? 'danger' : ''}`}
        role="menuitem"
        aria-disabled={item.disabled}
        tabIndex={0}
        onClick={(event) => onItemClick(item, event)}
        onKeyDown={(event) => {
          if (event.key === 'Enter' || event.key === ' ') {
            event.preventDefault();
            onItemClick(item, event as any);
          }
        }}
      >
        <div className="menu-item-content">
          {item.icon && (
            <div className="menu-item-icon">
              {item.icon}
            </div>
          )}
          
          <span className="menu-item-label">{item.label}</span>
          
          {item.shortcut && (
            <span className="menu-item-shortcut">{item.shortcut}</span>
          )}
          
          {item.submenu && (
            <div className="menu-item-arrow">
              <svg width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor">
                <polyline points="9,18 15,12 9,6" />
              </svg>
            </div>
          )}
        </div>
      </div>
    );
  };

  return (
    <div
      ref={menuRef}
      className="context-menu"
      style={{
        position: 'fixed',
        left: position.x,
        top: position.y,
        zIndex: 10000
      }}
      role="menu"
      aria-orientation="vertical"
    >
      {items.map((item, index) => renderMenuItem(item, index))}

      <style jsx>{`
        .context-menu {
          background: var(--bg-primary);
          border: 1px solid var(--border-color);
          border-radius: 6px;
          box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
          padding: 4px 0;
          min-width: 200px;
          max-width: 300px;
          font-size: 0.9rem;
          user-select: none;
          backdrop-filter: blur(8px);
        }

        .context-menu-item {
          padding: 0;
          cursor: pointer;
          transition: background-color 0.15s ease;
          position: relative;
        }

        .context-menu-item:hover:not(.disabled) {
          background: var(--bg-secondary);
        }

        .context-menu-item:focus {
          outline: none;
          background: var(--bg-secondary);
        }

        .context-menu-item.disabled {
          opacity: 0.5;
          cursor: not-allowed;
          pointer-events: none;
        }

        .context-menu-item.danger:hover {
          background: var(--error-color-light);
          color: var(--error-color);
        }

        .menu-item-content {
          display: flex;
          align-items: center;
          padding: 8px 12px;
          gap: 8px;
        }

        .menu-item-icon {
          display: flex;
          align-items: center;
          justify-content: center;
          width: 16px;
          height: 16px;
          color: var(--text-secondary);
        }

        .context-menu-item.danger .menu-item-icon {
          color: var(--error-color);
        }

        .menu-item-label {
          flex: 1;
          color: var(--text-primary);
          white-space: nowrap;
          overflow: hidden;
          text-overflow: ellipsis;
        }

        .context-menu-item.danger .menu-item-label {
          color: var(--error-color);
        }

        .menu-item-shortcut {
          font-size: 0.8rem;
          color: var(--text-secondary);
          font-family: var(--font-mono);
          margin-left: auto;
          padding-left: 16px;
        }

        .menu-item-arrow {
          color: var(--text-secondary);
          display: flex;
          align-items: center;
        }

        .context-menu-separator {
          height: 1px;
          background: var(--border-color);
          margin: 4px 0;
        }

        /* Dark theme adjustments */
        @media (prefers-color-scheme: dark) {
          .context-menu {
            background: var(--bg-primary);
            border-color: var(--border-color);
            box-shadow: 0 4px 12px rgba(0, 0, 0, 0.3);
          }
        }

        /* Animation */
        .context-menu {
          animation: contextMenuIn 0.15s ease-out;
        }

        @keyframes contextMenuIn {
          from {
            opacity: 0;
            transform: scale(0.95) translateY(-8px);
          }
          to {
            opacity: 1;
            transform: scale(1) translateY(0);
          }
        }

        /* High contrast mode */
        .high-contrast .context-menu {
          border: 2px solid var(--text-primary);
        }

        .high-contrast .context-menu-item:hover:not(.disabled) {
          background: var(--text-primary);
          color: var(--bg-primary);
        }

        /* Mobile adjustments */
        @media (max-width: 768px) {
          .context-menu {
            min-width: 180px;
            font-size: 0.95rem;
          }

          .menu-item-content {
            padding: 12px 16px;
          }

          .menu-item-shortcut {
            display: none;
          }
        }

        /* Reduced motion */
        @media (prefers-reduced-motion: reduce) {
          .context-menu {
            animation: none;
          }

          .context-menu-item {
            transition: none;
          }
        }
      `}</style>
    </div>
  );
};

export default ContextMenu;