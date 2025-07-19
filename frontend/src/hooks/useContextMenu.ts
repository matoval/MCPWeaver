import { useState, useCallback, useEffect, useRef } from 'react';

export interface ContextMenuItem {
  id: string;
  label: string;
  icon?: React.ReactNode;
  shortcut?: string;
  disabled?: boolean;
  separator?: boolean;
  submenu?: ContextMenuItem[];
  action?: () => void;
  danger?: boolean;
}

export interface ContextMenuPosition {
  x: number;
  y: number;
}

export interface ContextMenuState {
  isOpen: boolean;
  position: ContextMenuPosition;
  items: ContextMenuItem[];
  targetElement?: Element;
}

export const useContextMenu = () => {
  const [state, setState] = useState<ContextMenuState>({
    isOpen: false,
    position: { x: 0, y: 0 },
    items: []
  });

  const menuRef = useRef<HTMLDivElement>(null);
  const timeoutRef = useRef<NodeJS.Timeout>();

  const show = useCallback((
    event: React.MouseEvent | MouseEvent,
    items: ContextMenuItem[],
    targetElement?: Element
  ) => {
    event.preventDefault();
    event.stopPropagation();

    const { clientX: x, clientY: y } = event;
    
    setState({
      isOpen: true,
      position: { x, y },
      items,
      targetElement: targetElement || (event.target as Element)
    });
  }, []);

  const hide = useCallback(() => {
    setState(prev => ({ ...prev, isOpen: false }));
  }, []);

  const adjustPosition = useCallback((menuElement: HTMLDivElement, position: ContextMenuPosition) => {
    const menuRect = menuElement.getBoundingClientRect();
    const viewportWidth = window.innerWidth;
    const viewportHeight = window.innerHeight;

    let { x, y } = position;

    // Adjust horizontal position
    if (x + menuRect.width > viewportWidth) {
      x = viewportWidth - menuRect.width - 8;
    }
    if (x < 8) {
      x = 8;
    }

    // Adjust vertical position
    if (y + menuRect.height > viewportHeight) {
      y = viewportHeight - menuRect.height - 8;
    }
    if (y < 8) {
      y = 8;
    }

    return { x, y };
  }, []);

  const handleItemClick = useCallback((item: ContextMenuItem, event: React.MouseEvent) => {
    event.stopPropagation();
    
    if (item.disabled || item.separator) return;
    
    if (item.submenu) {
      // Handle submenu (could be expanded in the future)
      return;
    }

    if (item.action) {
      item.action();
    }

    hide();
  }, [hide]);

  const handleKeyDown = useCallback((event: KeyboardEvent) => {
    if (!state.isOpen) return;

    switch (event.key) {
      case 'Escape':
        event.preventDefault();
        hide();
        break;
      case 'ArrowDown':
        event.preventDefault();
        // Focus next item logic could be added here
        break;
      case 'ArrowUp':
        event.preventDefault();
        // Focus previous item logic could be added here
        break;
      case 'Enter':
        event.preventDefault();
        // Execute focused item logic could be added here
        break;
    }
  }, [state.isOpen, hide]);

  const handleClickOutside = useCallback((event: MouseEvent) => {
    if (state.isOpen && menuRef.current && !menuRef.current.contains(event.target as Node)) {
      hide();
    }
  }, [state.isOpen, hide]);

  useEffect(() => {
    if (state.isOpen) {
      document.addEventListener('keydown', handleKeyDown);
      document.addEventListener('click', handleClickOutside);
      document.addEventListener('contextmenu', handleClickOutside);

      // Adjust position after render
      if (menuRef.current) {
        const adjustedPosition = adjustPosition(menuRef.current, state.position);
        if (adjustedPosition.x !== state.position.x || adjustedPosition.y !== state.position.y) {
          setState(prev => ({ ...prev, position: adjustedPosition }));
        }
      }

      return () => {
        document.removeEventListener('keydown', handleKeyDown);
        document.removeEventListener('click', handleClickOutside);
        document.removeEventListener('contextmenu', handleClickOutside);
      };
    }
  }, [state.isOpen, state.position, handleKeyDown, handleClickOutside, adjustPosition]);

  const bindContextMenu = useCallback((
    element: HTMLElement | null,
    items: ContextMenuItem[] | (() => ContextMenuItem[])
  ) => {
    if (!element) return () => {};

    const handleContextMenu = (event: MouseEvent) => {
      const menuItems = typeof items === 'function' ? items() : items;
      show(event, menuItems, element);
    };

    element.addEventListener('contextmenu', handleContextMenu);

    return () => {
      element.removeEventListener('contextmenu', handleContextMenu);
    };
  }, [show]);

  const getContextMenuProps = useCallback(() => ({
    ref: menuRef,
    style: {
      position: 'fixed' as const,
      left: state.position.x,
      top: state.position.y,
      zIndex: 10000
    },
    role: 'menu',
    'aria-orientation': 'vertical' as const
  }), [state.position]);

  const getMenuItemProps = useCallback((item: ContextMenuItem, index: number) => ({
    key: item.id || index,
    role: item.separator ? 'separator' : 'menuitem',
    'aria-disabled': item.disabled,
    tabIndex: item.separator ? -1 : 0,
    onClick: (event: React.MouseEvent) => handleItemClick(item, event),
    onKeyDown: (event: React.KeyboardEvent) => {
      if (event.key === 'Enter' || event.key === ' ') {
        event.preventDefault();
        handleItemClick(item, event as any);
      }
    }
  }), [handleItemClick]);

  return {
    // State
    isOpen: state.isOpen,
    position: state.position,
    items: state.items,
    targetElement: state.targetElement,
    
    // Methods
    show,
    hide,
    bindContextMenu,
    
    // Props helpers
    getContextMenuProps,
    getMenuItemProps,
    
    // Ref
    menuRef
  };
};

export default useContextMenu;