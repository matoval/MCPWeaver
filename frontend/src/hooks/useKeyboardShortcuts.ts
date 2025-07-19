import { useEffect, useCallback, useRef } from 'react';

export interface KeyboardShortcut {
  key: string;
  ctrlKey?: boolean;
  shiftKey?: boolean;
  altKey?: boolean;
  metaKey?: boolean;
  action: () => void;
  description: string;
  category: string;
  preventDefault?: boolean;
}

export interface ShortcutCategory {
  name: string;
  shortcuts: KeyboardShortcut[];
}

const DEFAULT_SHORTCUTS: KeyboardShortcut[] = [
  // File operations
  {
    key: 'n',
    ctrlKey: true,
    action: () => window.dispatchEvent(new CustomEvent('keyboard:new-project')),
    description: 'Create new project',
    category: 'File',
    preventDefault: true
  },
  {
    key: 'o',
    ctrlKey: true,
    action: () => window.dispatchEvent(new CustomEvent('keyboard:open-file')),
    description: 'Open file',
    category: 'File',
    preventDefault: true
  },
  {
    key: 's',
    ctrlKey: true,
    action: () => window.dispatchEvent(new CustomEvent('keyboard:save-project')),
    description: 'Save project',
    category: 'File',
    preventDefault: true
  },
  {
    key: 'e',
    ctrlKey: true,
    action: () => window.dispatchEvent(new CustomEvent('keyboard:export-project')),
    description: 'Export project',
    category: 'File',
    preventDefault: true
  },
  
  // Generation operations
  {
    key: 'F5',
    action: () => window.dispatchEvent(new CustomEvent('keyboard:generate-server')),
    description: 'Generate MCP server',
    category: 'Generation',
    preventDefault: true
  },
  {
    key: 'F9',
    action: () => window.dispatchEvent(new CustomEvent('keyboard:validate-spec')),
    description: 'Validate OpenAPI spec',
    category: 'Generation',
    preventDefault: true
  },
  
  // Navigation
  {
    key: 'Tab',
    ctrlKey: true,
    action: () => window.dispatchEvent(new CustomEvent('keyboard:next-tab')),
    description: 'Next tab',
    category: 'Navigation',
    preventDefault: true
  },
  {
    key: 'Tab',
    ctrlKey: true,
    shiftKey: true,
    action: () => window.dispatchEvent(new CustomEvent('keyboard:prev-tab')),
    description: 'Previous tab',
    category: 'Navigation',
    preventDefault: true
  },
  {
    key: 'Home',
    ctrlKey: true,
    action: () => window.dispatchEvent(new CustomEvent('keyboard:go-dashboard')),
    description: 'Go to dashboard',
    category: 'Navigation',
    preventDefault: true
  },
  
  // View operations
  {
    key: 'F11',
    action: () => window.dispatchEvent(new CustomEvent('keyboard:toggle-fullscreen')),
    description: 'Toggle fullscreen',
    category: 'View',
    preventDefault: true
  },
  {
    key: 'F1',
    action: () => window.dispatchEvent(new CustomEvent('keyboard:show-help')),
    description: 'Show help',
    category: 'Help',
    preventDefault: true
  },
  {
    key: '/',
    ctrlKey: true,
    action: () => window.dispatchEvent(new CustomEvent('keyboard:show-shortcuts')),
    description: 'Show keyboard shortcuts',
    category: 'Help',
    preventDefault: true
  },
  
  // Edit operations
  {
    key: 'z',
    ctrlKey: true,
    action: () => window.dispatchEvent(new CustomEvent('keyboard:undo')),
    description: 'Undo',
    category: 'Edit',
    preventDefault: true
  },
  {
    key: 'y',
    ctrlKey: true,
    action: () => window.dispatchEvent(new CustomEvent('keyboard:redo')),
    description: 'Redo',
    category: 'Edit',
    preventDefault: true
  },
  {
    key: 'z',
    ctrlKey: true,
    shiftKey: true,
    action: () => window.dispatchEvent(new CustomEvent('keyboard:redo')),
    description: 'Redo (alternative)',
    category: 'Edit',
    preventDefault: true
  },
  
  // Search
  {
    key: 'f',
    ctrlKey: true,
    action: () => window.dispatchEvent(new CustomEvent('keyboard:search')),
    description: 'Search',
    category: 'Search',
    preventDefault: true
  },
  {
    key: 'p',
    ctrlKey: true,
    action: () => window.dispatchEvent(new CustomEvent('keyboard:command-palette')),
    description: 'Command palette',
    category: 'Search',
    preventDefault: true
  }
];

export const useKeyboardShortcuts = (
  customShortcuts: KeyboardShortcut[] = [],
  enabled: boolean = true
) => {
  const shortcutsRef = useRef<KeyboardShortcut[]>([...DEFAULT_SHORTCUTS, ...customShortcuts]);
  const enabledRef = useRef(enabled);

  // Update refs when props change
  useEffect(() => {
    shortcutsRef.current = [...DEFAULT_SHORTCUTS, ...customShortcuts];
  }, [customShortcuts]);

  useEffect(() => {
    enabledRef.current = enabled;
  }, [enabled]);

  const handleKeyDown = useCallback((event: KeyboardEvent) => {
    if (!enabledRef.current) return;

    // Skip if user is typing in an input/textarea
    const target = event.target as HTMLElement;
    if (target.tagName === 'INPUT' || target.tagName === 'TEXTAREA' || target.isContentEditable) {
      return;
    }

    const shortcuts = shortcutsRef.current;
    
    for (const shortcut of shortcuts) {
      const keyMatches = shortcut.key.toLowerCase() === event.key.toLowerCase();
      const ctrlMatches = (shortcut.ctrlKey ?? false) === event.ctrlKey;
      const shiftMatches = (shortcut.shiftKey ?? false) === event.shiftKey;
      const altMatches = (shortcut.altKey ?? false) === event.altKey;
      const metaMatches = (shortcut.metaKey ?? false) === event.metaKey;

      if (keyMatches && ctrlMatches && shiftMatches && altMatches && metaMatches) {
        if (shortcut.preventDefault !== false) {
          event.preventDefault();
          event.stopPropagation();
        }
        
        try {
          shortcut.action();
        } catch (error) {
          console.error('Error executing keyboard shortcut:', error);
        }
        
        break;
      }
    }
  }, []);

  useEffect(() => {
    if (enabled) {
      document.addEventListener('keydown', handleKeyDown);
      return () => document.removeEventListener('keydown', handleKeyDown);
    }
  }, [handleKeyDown, enabled]);

  const getShortcutsByCategory = useCallback((): ShortcutCategory[] => {
    const categories = new Map<string, KeyboardShortcut[]>();
    
    shortcutsRef.current.forEach(shortcut => {
      if (!categories.has(shortcut.category)) {
        categories.set(shortcut.category, []);
      }
      categories.get(shortcut.category)!.push(shortcut);
    });

    return Array.from(categories.entries()).map(([name, shortcuts]) => ({
      name,
      shortcuts: shortcuts.sort((a, b) => a.description.localeCompare(b.description))
    }));
  }, []);

  const getShortcutText = useCallback((shortcut: KeyboardShortcut): string => {
    const parts: string[] = [];
    
    if (shortcut.ctrlKey) parts.push('Ctrl');
    if (shortcut.metaKey) parts.push('Cmd');
    if (shortcut.altKey) parts.push('Alt');
    if (shortcut.shiftKey) parts.push('Shift');
    
    parts.push(shortcut.key === ' ' ? 'Space' : shortcut.key);
    
    return parts.join(' + ');
  }, []);

  return {
    shortcuts: shortcutsRef.current,
    getShortcutsByCategory,
    getShortcutText,
    enabled: enabledRef.current
  };
};

export default useKeyboardShortcuts;