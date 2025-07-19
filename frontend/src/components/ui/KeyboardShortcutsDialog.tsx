import React, { useState, useEffect } from 'react';
import { useKeyboardShortcuts, ShortcutCategory } from '../../hooks/useKeyboardShortcuts';

interface KeyboardShortcutsDialogProps {
  isOpen: boolean;
  onClose: () => void;
}

const KeyboardShortcutsDialog: React.FC<KeyboardShortcutsDialogProps> = ({
  isOpen,
  onClose
}) => {
  const { getShortcutsByCategory, getShortcutText } = useKeyboardShortcuts();
  const [categories, setCategories] = useState<ShortcutCategory[]>([]);
  const [searchQuery, setSearchQuery] = useState('');

  useEffect(() => {
    if (isOpen) {
      setCategories(getShortcutsByCategory());
    }
  }, [isOpen, getShortcutsByCategory]);

  useEffect(() => {
    const handleEscape = (e: KeyboardEvent) => {
      if (e.key === 'Escape' && isOpen) {
        onClose();
      }
    };

    document.addEventListener('keydown', handleEscape);
    return () => document.removeEventListener('keydown', handleEscape);
  }, [isOpen, onClose]);

  const filteredCategories = categories.map(category => ({
    ...category,
    shortcuts: category.shortcuts.filter(shortcut =>
      shortcut.description.toLowerCase().includes(searchQuery.toLowerCase()) ||
      getShortcutText(shortcut).toLowerCase().includes(searchQuery.toLowerCase())
    )
  })).filter(category => category.shortcuts.length > 0);

  if (!isOpen) return null;

  return (
    <div className="keyboard-shortcuts-overlay" onClick={onClose}>
      <div className="keyboard-shortcuts-dialog" onClick={e => e.stopPropagation()}>
        <div className="dialog-header">
          <h2>Keyboard Shortcuts</h2>
          <button 
            className="close-button" 
            onClick={onClose}
            aria-label="Close dialog"
          >
            ×
          </button>
        </div>

        <div className="search-container">
          <input
            type="text"
            placeholder="Search shortcuts..."
            value={searchQuery}
            onChange={(e) => setSearchQuery(e.target.value)}
            className="search-input"
            autoFocus
          />
        </div>

        <div className="shortcuts-content">
          {filteredCategories.map(category => (
            <div key={category.name} className="shortcut-category">
              <h3 className="category-title">{category.name}</h3>
              <div className="shortcuts-grid">
                {category.shortcuts.map((shortcut, index) => (
                  <div key={index} className="shortcut-item">
                    <span className="shortcut-description">
                      {shortcut.description}
                    </span>
                    <kbd className="shortcut-keys">
                      {getShortcutText(shortcut)}
                    </kbd>
                  </div>
                ))}
              </div>
            </div>
          ))}
          
          {filteredCategories.length === 0 && searchQuery && (
            <div className="no-results">
              No shortcuts found for "{searchQuery}"
            </div>
          )}
        </div>

        <div className="dialog-footer">
          <p className="help-text">
            Press <kbd>Esc</kbd> to close this dialog
          </p>
        </div>
      </div>

      <style jsx>{`
        .keyboard-shortcuts-overlay {
          position: fixed;
          top: 0;
          left: 0;
          right: 0;
          bottom: 0;
          background: rgba(0, 0, 0, 0.6);
          display: flex;
          align-items: center;
          justify-content: center;
          z-index: 1000;
          backdrop-filter: blur(4px);
        }

        .keyboard-shortcuts-dialog {
          background: var(--bg-primary);
          border: 1px solid var(--border-color);
          border-radius: 8px;
          box-shadow: 0 10px 30px rgba(0, 0, 0, 0.3);
          width: 90%;
          max-width: 800px;
          max-height: 80vh;
          display: flex;
          flex-direction: column;
          color: var(--text-primary);
        }

        .dialog-header {
          padding: 20px 24px 16px;
          border-bottom: 1px solid var(--border-color);
          display: flex;
          align-items: center;
          justify-content: space-between;
        }

        .dialog-header h2 {
          margin: 0;
          font-size: 1.5rem;
          font-weight: 600;
        }

        .close-button {
          background: none;
          border: none;
          font-size: 1.5rem;
          color: var(--text-secondary);
          cursor: pointer;
          padding: 4px 8px;
          border-radius: 4px;
          transition: all 0.2s ease;
        }

        .close-button:hover {
          background: var(--bg-secondary);
          color: var(--text-primary);
        }

        .search-container {
          padding: 16px 24px;
          border-bottom: 1px solid var(--border-color);
        }

        .search-input {
          width: 100%;
          padding: 8px 12px;
          border: 1px solid var(--border-color);
          border-radius: 4px;
          background: var(--bg-secondary);
          color: var(--text-primary);
          font-size: 0.9rem;
        }

        .search-input:focus {
          outline: none;
          border-color: var(--accent-color);
          box-shadow: 0 0 0 2px var(--accent-color-alpha);
        }

        .shortcuts-content {
          flex: 1;
          overflow-y: auto;
          padding: 20px 24px;
        }

        .shortcut-category {
          margin-bottom: 24px;
        }

        .shortcut-category:last-child {
          margin-bottom: 0;
        }

        .category-title {
          font-size: 1.1rem;
          font-weight: 600;
          margin: 0 0 12px 0;
          color: var(--accent-color);
        }

        .shortcuts-grid {
          display: grid;
          gap: 8px;
        }

        .shortcut-item {
          display: flex;
          align-items: center;
          justify-content: space-between;
          padding: 8px 12px;
          background: var(--bg-secondary);
          border-radius: 4px;
          transition: background 0.2s ease;
        }

        .shortcut-item:hover {
          background: var(--bg-tertiary);
        }

        .shortcut-description {
          font-size: 0.9rem;
          color: var(--text-primary);
        }

        .shortcut-keys {
          font-family: 'SF Mono', Monaco, Consolas, monospace;
          font-size: 0.8rem;
          background: var(--bg-tertiary);
          color: var(--text-secondary);
          padding: 4px 8px;
          border-radius: 4px;
          border: 1px solid var(--border-color);
          white-space: nowrap;
        }

        .no-results {
          text-align: center;
          color: var(--text-secondary);
          padding: 40px 20px;
          font-style: italic;
        }

        .dialog-footer {
          padding: 16px 24px;
          border-top: 1px solid var(--border-color);
          text-align: center;
        }

        .help-text {
          margin: 0;
          font-size: 0.85rem;
          color: var(--text-secondary);
        }

        .help-text kbd {
          font-family: inherit;
          background: var(--bg-tertiary);
          padding: 2px 6px;
          border-radius: 3px;
          border: 1px solid var(--border-color);
          font-size: 0.8rem;
        }

        @media (max-width: 768px) {
          .keyboard-shortcuts-dialog {
            width: 95%;
            max-height: 90vh;
          }
          
          .shortcuts-grid {
            gap: 6px;
          }
          
          .shortcut-item {
            padding: 6px 8px;
          }
          
          .shortcut-description {
            font-size: 0.85rem;
          }
          
          .shortcut-keys {
            font-size: 0.75rem;
            padding: 3px 6px;
          }
        }
      `}</style>
    </div>
  );
};

export default KeyboardShortcutsDialog;