import { useEffect, useCallback, useRef } from 'react';

export interface FocusableElement {
  element: HTMLElement;
  priority: number;
  group?: string;
}

export interface FocusGroup {
  name: string;
  elements: HTMLElement[];
  current: number;
}

export const useFocusManagement = (enabled: boolean = true) => {
  const focusableElementsRef = useRef<Map<string, FocusGroup>>(new Map());
  const currentGroupRef = useRef<string | null>(null);
  const enabledRef = useRef(enabled);

  useEffect(() => {
    enabledRef.current = enabled;
  }, [enabled]);

  const getFocusableElements = useCallback((container?: HTMLElement): HTMLElement[] => {
    const root = container || document.body;
    const selector = [
      'a[href]:not([disabled])',
      'button:not([disabled])',
      'textarea:not([disabled])',
      'input:not([disabled]):not([type="hidden"])',
      'select:not([disabled])',
      '[tabindex]:not([tabindex="-1"]):not([disabled])',
      '[contenteditable]:not([contenteditable="false"])'
    ].join(',');

    return Array.from(root.querySelectorAll(selector))
      .filter(el => {
        const element = el as HTMLElement;
        return isElementVisible(element) && !isElementInert(element);
      }) as HTMLElement[];
  }, []);

  const isElementVisible = useCallback((element: HTMLElement): boolean => {
    const style = window.getComputedStyle(element);
    return (
      style.display !== 'none' &&
      style.visibility !== 'hidden' &&
      style.opacity !== '0' &&
      element.offsetWidth > 0 &&
      element.offsetHeight > 0
    );
  }, []);

  const isElementInert = useCallback((element: HTMLElement): boolean => {
    // Check if element or any parent has inert attribute
    let current: HTMLElement | null = element;
    while (current) {
      if (current.hasAttribute('inert') || current.getAttribute('aria-hidden') === 'true') {
        return true;
      }
      current = current.parentElement;
    }
    return false;
  }, []);

  const registerFocusGroup = useCallback((
    groupName: string,
    container?: HTMLElement,
    options?: { circular?: boolean; autoFocus?: boolean }
  ) => {
    const elements = getFocusableElements(container);
    const group: FocusGroup = {
      name: groupName,
      elements,
      current: 0
    };

    focusableElementsRef.current.set(groupName, group);

    if (options?.autoFocus && elements.length > 0) {
      elements[0].focus();
      currentGroupRef.current = groupName;
    }

    return () => {
      focusableElementsRef.current.delete(groupName);
      if (currentGroupRef.current === groupName) {
        currentGroupRef.current = null;
      }
    };
  }, [getFocusableElements]);

  const focusNext = useCallback((groupName?: string) => {
    if (!enabledRef.current) return false;

    const targetGroup = groupName || currentGroupRef.current;
    if (!targetGroup) return false;

    const group = focusableElementsRef.current.get(targetGroup);
    if (!group || group.elements.length === 0) return false;

    group.current = (group.current + 1) % group.elements.length;
    group.elements[group.current].focus();
    currentGroupRef.current = targetGroup;
    
    return true;
  }, []);

  const focusPrevious = useCallback((groupName?: string) => {
    if (!enabledRef.current) return false;

    const targetGroup = groupName || currentGroupRef.current;
    if (!targetGroup) return false;

    const group = focusableElementsRef.current.get(targetGroup);
    if (!group || group.elements.length === 0) return false;

    group.current = group.current === 0 ? group.elements.length - 1 : group.current - 1;
    group.elements[group.current].focus();
    currentGroupRef.current = targetGroup;
    
    return true;
  }, []);

  const focusFirst = useCallback((groupName?: string) => {
    if (!enabledRef.current) return false;

    const targetGroup = groupName || currentGroupRef.current;
    if (!targetGroup) return false;

    const group = focusableElementsRef.current.get(targetGroup);
    if (!group || group.elements.length === 0) return false;

    group.current = 0;
    group.elements[group.current].focus();
    currentGroupRef.current = targetGroup;
    
    return true;
  }, []);

  const focusLast = useCallback((groupName?: string) => {
    if (!enabledRef.current) return false;

    const targetGroup = groupName || currentGroupRef.current;
    if (!targetGroup) return false;

    const group = focusableElementsRef.current.get(targetGroup);
    if (!group || group.elements.length === 0) return false;

    group.current = group.elements.length - 1;
    group.elements[group.current].focus();
    currentGroupRef.current = targetGroup;
    
    return true;
  }, []);

  const switchGroup = useCallback((groupName: string) => {
    if (!enabledRef.current) return false;

    const group = focusableElementsRef.current.get(groupName);
    if (!group || group.elements.length === 0) return false;

    currentGroupRef.current = groupName;
    group.elements[group.current].focus();
    
    return true;
  }, []);

  const trapFocus = useCallback((
    container: HTMLElement,
    options?: { initialFocus?: HTMLElement; returnFocus?: HTMLElement }
  ) => {
    const focusableElements = getFocusableElements(container);
    if (focusableElements.length === 0) return () => {};

    const firstElement = focusableElements[0];
    const lastElement = focusableElements[focusableElements.length - 1];
    const previousActiveElement = document.activeElement as HTMLElement;

    // Focus initial element
    const initialFocus = options?.initialFocus || firstElement;
    initialFocus.focus();

    const handleKeyDown = (e: KeyboardEvent) => {
      if (e.key !== 'Tab') return;

      if (e.shiftKey) {
        // Shift + Tab
        if (document.activeElement === firstElement) {
          e.preventDefault();
          lastElement.focus();
        }
      } else {
        // Tab
        if (document.activeElement === lastElement) {
          e.preventDefault();
          firstElement.focus();
        }
      }
    };

    document.addEventListener('keydown', handleKeyDown);

    return () => {
      document.removeEventListener('keydown', handleKeyDown);
      if (options?.returnFocus && previousActiveElement) {
        previousActiveElement.focus();
      }
    };
  }, [getFocusableElements]);

  const handleArrowNavigation = useCallback((e: KeyboardEvent) => {
    if (!enabledRef.current) return;

    const currentGroup = currentGroupRef.current;
    if (!currentGroup) return;

    switch (e.key) {
      case 'ArrowDown':
      case 'ArrowRight':
        if (focusNext()) {
          e.preventDefault();
          e.stopPropagation();
        }
        break;
      case 'ArrowUp':
      case 'ArrowLeft':
        if (focusPrevious()) {
          e.preventDefault();
          e.stopPropagation();
        }
        break;
      case 'Home':
        if (focusFirst()) {
          e.preventDefault();
          e.stopPropagation();
        }
        break;
      case 'End':
        if (focusLast()) {
          e.preventDefault();
          e.stopPropagation();
        }
        break;
    }
  }, [focusNext, focusPrevious, focusFirst, focusLast]);

  useEffect(() => {
    if (enabled) {
      document.addEventListener('keydown', handleArrowNavigation);
      return () => document.removeEventListener('keydown', handleArrowNavigation);
    }
  }, [handleArrowNavigation, enabled]);

  return {
    registerFocusGroup,
    focusNext,
    focusPrevious,
    focusFirst,
    focusLast,
    switchGroup,
    trapFocus,
    getFocusableElements,
    currentGroup: currentGroupRef.current,
    enabled: enabledRef.current
  };
};

export default useFocusManagement;