import { useEffect, useCallback, useRef } from 'react';
import useFocusManagement from './useFocusManagement';

export interface AccessibilityOptions {
  announcePageChanges?: boolean;
  manageFocus?: boolean;
  enableKeyboardNavigation?: boolean;
  skipLinks?: boolean;
  highContrast?: boolean;
  reducedMotion?: boolean;
}

export const useAccessibility = (options: AccessibilityOptions = {}) => {
  const {
    announcePageChanges = true,
    manageFocus = true,
    enableKeyboardNavigation = true,
    skipLinks = true,
    highContrast = false,
    reducedMotion = false
  } = options;

  const announceRegionRef = useRef<HTMLDivElement | null>(null);
  const focusManagement = useFocusManagement(enableKeyboardNavigation);

  // Create screen reader announcement region
  useEffect(() => {
    if (announcePageChanges && !announceRegionRef.current) {
      const region = document.createElement('div');
      region.setAttribute('aria-live', 'polite');
      region.setAttribute('aria-atomic', 'true');
      region.className = 'sr-only';
      region.style.cssText = `
        position: absolute;
        width: 1px;
        height: 1px;
        padding: 0;
        margin: -1px;
        overflow: hidden;
        clip: rect(0, 0, 0, 0);
        white-space: nowrap;
        border: 0;
      `;
      document.body.appendChild(region);
      announceRegionRef.current = region;
    }

    return () => {
      if (announceRegionRef.current) {
        document.body.removeChild(announceRegionRef.current);
        announceRegionRef.current = null;
      }
    };
  }, [announcePageChanges]);

  // Announce content to screen readers
  const announce = useCallback((message: string, priority: 'polite' | 'assertive' = 'polite') => {
    if (!announceRegionRef.current) return;

    announceRegionRef.current.setAttribute('aria-live', priority);
    announceRegionRef.current.textContent = message;

    // Clear after announcement to ensure repeated messages are announced
    setTimeout(() => {
      if (announceRegionRef.current) {
        announceRegionRef.current.textContent = '';
      }
    }, 1000);
  }, []);

  // Add ARIA attributes to element
  const addAriaAttributes = useCallback((
    element: HTMLElement,
    attributes: Record<string, string | boolean | null>
  ) => {
    Object.entries(attributes).forEach(([key, value]) => {
      if (value === null) {
        element.removeAttribute(key);
      } else {
        element.setAttribute(key, String(value));
      }
    });
  }, []);

  // Create ARIA describedby relationship
  const createAriaDescription = useCallback((
    targetElement: HTMLElement,
    description: string,
    id?: string
  ): string => {
    const descriptionId = id || `desc-${Math.random().toString(36).substr(2, 9)}`;
    
    let descElement = document.getElementById(descriptionId);
    if (!descElement) {
      descElement = document.createElement('div');
      descElement.id = descriptionId;
      descElement.className = 'sr-only';
      descElement.style.cssText = `
        position: absolute;
        width: 1px;
        height: 1px;
        padding: 0;
        margin: -1px;
        overflow: hidden;
        clip: rect(0, 0, 0, 0);
        white-space: nowrap;
        border: 0;
      `;
      document.body.appendChild(descElement);
    }
    
    descElement.textContent = description;
    
    const existingDescribedBy = targetElement.getAttribute('aria-describedby');
    const describedByIds = existingDescribedBy 
      ? existingDescribedBy.split(' ').filter(id => id !== descriptionId)
      : [];
    describedByIds.push(descriptionId);
    
    targetElement.setAttribute('aria-describedby', describedByIds.join(' '));
    
    return descriptionId;
  }, []);

  // Handle skip links
  const createSkipLink = useCallback((
    targetId: string,
    text: string = 'Skip to main content'
  ): HTMLAnchorElement => {
    const skipLink = document.createElement('a');
    skipLink.href = `#${targetId}`;
    skipLink.textContent = text;
    skipLink.className = 'skip-link';
    skipLink.style.cssText = `
      position: absolute;
      top: -40px;
      left: 6px;
      background: var(--bg-primary, #fff);
      color: var(--text-primary, #000);
      padding: 8px;
      text-decoration: none;
      border: 2px solid var(--accent-color, #007acc);
      border-radius: 4px;
      z-index: 10000;
      transition: top 0.3s;
    `;

    skipLink.addEventListener('focus', () => {
      skipLink.style.top = '6px';
    });

    skipLink.addEventListener('blur', () => {
      skipLink.style.top = '-40px';
    });

    skipLink.addEventListener('click', (e) => {
      e.preventDefault();
      const target = document.getElementById(targetId);
      if (target) {
        target.focus();
        target.scrollIntoView({ behavior: 'smooth' });
      }
    });

    return skipLink;
  }, []);

  // Manage landmark regions
  const createLandmark = useCallback((
    element: HTMLElement,
    role: string,
    label?: string
  ) => {
    element.setAttribute('role', role);
    if (label) {
      element.setAttribute('aria-label', label);
    }

    // Ensure landmarks are focusable for screen reader navigation
    if (!element.hasAttribute('tabindex')) {
      element.setAttribute('tabindex', '-1');
    }
  }, []);

  // Handle high contrast mode
  useEffect(() => {
    if (highContrast) {
      document.documentElement.classList.add('high-contrast');
    } else {
      document.documentElement.classList.remove('high-contrast');
    }
  }, [highContrast]);

  // Handle reduced motion preference
  useEffect(() => {
    const mediaQuery = window.matchMedia('(prefers-reduced-motion: reduce)');
    
    const handleMotionPreference = (e: MediaQueryListEvent) => {
      if (e.matches || reducedMotion) {
        document.documentElement.classList.add('reduced-motion');
      } else {
        document.documentElement.classList.remove('reduced-motion');
      }
    };

    // Set initial state
    handleMotionPreference({ matches: mediaQuery.matches } as MediaQueryListEvent);
    
    mediaQuery.addEventListener('change', handleMotionPreference);
    return () => mediaQuery.removeEventListener('change', handleMotionPreference);
  }, [reducedMotion]);

  // Create accessible form controls
  const enhanceFormControl = useCallback((
    input: HTMLInputElement | HTMLTextAreaElement | HTMLSelectElement,
    options: {
      label?: string;
      description?: string;
      required?: boolean;
      invalid?: boolean;
      errorMessage?: string;
    } = {}
  ) => {
    const { label, description, required, invalid, errorMessage } = options;
    
    // Handle label
    if (label) {
      let labelElement = document.querySelector(`label[for="${input.id}"]`) as HTMLLabelElement;
      if (!labelElement && input.id) {
        labelElement = document.createElement('label');
        labelElement.setAttribute('for', input.id);
        labelElement.textContent = label;
        input.parentNode?.insertBefore(labelElement, input);
      }
    }

    // Handle required
    if (required !== undefined) {
      if (required) {
        input.setAttribute('required', '');
        input.setAttribute('aria-required', 'true');
      } else {
        input.removeAttribute('required');
        input.removeAttribute('aria-required');
      }
    }

    // Handle invalid state
    if (invalid !== undefined) {
      input.setAttribute('aria-invalid', String(invalid));
      
      if (invalid && errorMessage) {
        const errorId = `${input.id || 'input'}-error`;
        let errorElement = document.getElementById(errorId);
        
        if (!errorElement) {
          errorElement = document.createElement('div');
          errorElement.id = errorId;
          errorElement.className = 'error-message';
          errorElement.setAttribute('role', 'alert');
          input.parentNode?.appendChild(errorElement);
        }
        
        errorElement.textContent = errorMessage;
        input.setAttribute('aria-describedby', 
          `${input.getAttribute('aria-describedby') || ''} ${errorId}`.trim()
        );
      }
    }

    // Handle description
    if (description) {
      createAriaDescription(input, description);
    }
  }, [createAriaDescription]);

  // Create accessible notifications
  const createNotification = useCallback((
    message: string,
    type: 'info' | 'success' | 'warning' | 'error' = 'info',
    timeout: number = 5000
  ) => {
    const notification = document.createElement('div');
    notification.className = `notification notification--${type}`;
    notification.setAttribute('role', type === 'error' ? 'alert' : 'status');
    notification.setAttribute('aria-live', type === 'error' ? 'assertive' : 'polite');
    notification.textContent = message;

    document.body.appendChild(notification);

    if (timeout > 0) {
      setTimeout(() => {
        if (notification.parentNode) {
          notification.parentNode.removeChild(notification);
        }
      }, timeout);
    }

    return notification;
  }, []);

  return {
    // Announcement functions
    announce,
    
    // ARIA helpers
    addAriaAttributes,
    createAriaDescription,
    createLandmark,
    
    // Form accessibility
    enhanceFormControl,
    
    // Navigation helpers
    createSkipLink,
    
    // Notifications
    createNotification,
    
    // Focus management (delegated)
    ...focusManagement,
    
    // State
    options: {
      announcePageChanges,
      manageFocus,
      enableKeyboardNavigation,
      skipLinks,
      highContrast,
      reducedMotion
    }
  };
};

export default useAccessibility;