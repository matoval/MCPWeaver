import React, { createContext, useContext, useState, useCallback, useEffect } from 'react';
import { APIError } from '../types';
import './NotificationSystem.scss';

export type NotificationType = 'success' | 'error' | 'warning' | 'info';

export interface NotificationData {
  id: string;
  type: NotificationType;
  title: string;
  message: string;
  duration?: number;
  persistent?: boolean;
  actions?: NotificationAction[];
  metadata?: Record<string, any>;
}

export interface NotificationAction {
  label: string;
  handler: () => void;
  style?: 'primary' | 'secondary' | 'danger';
}

interface NotificationContextType {
  notifications: NotificationData[];
  addNotification: (notification: Omit<NotificationData, 'id'>) => string;
  removeNotification: (id: string) => void;
  clearAll: () => void;
  showError: (error: APIError | Error | string, options?: Partial<NotificationData>) => string;
  showSuccess: (message: string, options?: Partial<NotificationData>) => string;
  showWarning: (message: string, options?: Partial<NotificationData>) => string;
  showInfo: (message: string, options?: Partial<NotificationData>) => string;
}

const NotificationContext = createContext<NotificationContextType | undefined>(undefined);

export const useNotifications = () => {
  const context = useContext(NotificationContext);
  if (!context) {
    throw new Error('useNotifications must be used within a NotificationProvider');
  }
  return context;
};

interface NotificationProviderProps {
  children: React.ReactNode;
  maxNotifications?: number;
  defaultDuration?: number;
  position?: 'top-right' | 'top-left' | 'bottom-right' | 'bottom-left' | 'top-center' | 'bottom-center';
}

export const NotificationProvider: React.FC<NotificationProviderProps> = ({
  children,
  maxNotifications = 5,
  defaultDuration = 5000,
  position = 'top-right'
}) => {
  const [notifications, setNotifications] = useState<NotificationData[]>([]);

  const generateId = useCallback(() => {
    return `notification-${Date.now()}-${Math.random().toString(36).substr(2, 9)}`;
  }, []);

  const addNotification = useCallback((notification: Omit<NotificationData, 'id'>) => {
    const id = generateId();
    const newNotification: NotificationData = {
      ...notification,
      id,
      duration: notification.duration ?? defaultDuration
    };

    setNotifications(prev => {
      const updated = [newNotification, ...prev];
      // Keep only the most recent notifications
      return updated.slice(0, maxNotifications);
    });

    // Auto-remove if not persistent
    if (!notification.persistent && newNotification.duration > 0) {
      setTimeout(() => {
        removeNotification(id);
      }, newNotification.duration);
    }

    return id;
  }, [generateId, defaultDuration, maxNotifications]);

  const removeNotification = useCallback((id: string) => {
    setNotifications(prev => prev.filter(n => n.id !== id));
  }, []);

  const clearAll = useCallback(() => {
    setNotifications([]);
  }, []);

  const showError = useCallback((error: APIError | Error | string, options?: Partial<NotificationData>) => {
    let title = 'Error';
    let message = '';
    let actions: NotificationAction[] = [];

    if (typeof error === 'string') {
      message = error;
    } else if (error instanceof Error) {
      message = error.message;
      title = 'Application Error';
    } else {
      // APIError
      message = error.message;
      title = getCategoryTitle(error.type);
      
      if (error.suggestions && error.suggestions.length > 0) {
        message += '\n\nSuggestions:\n' + error.suggestions.map(s => `• ${s}`).join('\n');
      }

      if (error.recoverable) {
        actions.push({
          label: 'Retry',
          handler: () => {
            // Emit retry event that can be handled by the calling component
            window.dispatchEvent(new CustomEvent('retry-error', { detail: error }));
          },
          style: 'primary'
        });
      }
    }

    return addNotification({
      type: 'error',
      title,
      message,
      actions,
      persistent: true,
      ...options
    });
  }, [addNotification]);

  const showSuccess = useCallback((message: string, options?: Partial<NotificationData>) => {
    return addNotification({
      type: 'success',
      title: 'Success',
      message,
      ...options
    });
  }, [addNotification]);

  const showWarning = useCallback((message: string, options?: Partial<NotificationData>) => {
    return addNotification({
      type: 'warning',
      title: 'Warning',
      message,
      ...options
    });
  }, [addNotification]);

  const showInfo = useCallback((message: string, options?: Partial<NotificationData>) => {
    return addNotification({
      type: 'info',
      title: 'Information',
      message,
      ...options
    });
  }, [addNotification]);

  const value: NotificationContextType = {
    notifications,
    addNotification,
    removeNotification,
    clearAll,
    showError,
    showSuccess,
    showWarning,
    showInfo
  };

  return (
    <NotificationContext.Provider value={value}>
      {children}
      <NotificationContainer position={position} />
    </NotificationContext.Provider>
  );
};

interface NotificationContainerProps {
  position: 'top-right' | 'top-left' | 'bottom-right' | 'bottom-left' | 'top-center' | 'bottom-center';
}

const NotificationContainer: React.FC<NotificationContainerProps> = ({ position }) => {
  const { notifications, removeNotification } = useNotifications();

  return (
    <div className={`notification-container notification-container--${position}`}>
      {notifications.map(notification => (
        <NotificationItem
          key={notification.id}
          notification={notification}
          onRemove={removeNotification}
        />
      ))}
    </div>
  );
};

interface NotificationItemProps {
  notification: NotificationData;
  onRemove: (id: string) => void;
}

const NotificationItem: React.FC<NotificationItemProps> = ({ notification, onRemove }) => {
  const [isVisible, setIsVisible] = useState(false);
  const [isRemoving, setIsRemoving] = useState(false);

  useEffect(() => {
    // Trigger entrance animation
    const timer = setTimeout(() => setIsVisible(true), 10);
    return () => clearTimeout(timer);
  }, []);

  const handleRemove = useCallback(() => {
    setIsRemoving(true);
    setTimeout(() => {
      onRemove(notification.id);
    }, 300); // Match CSS transition duration
  }, [notification.id, onRemove]);

  const handleAction = useCallback((action: NotificationAction) => {
    action.handler();
    handleRemove();
  }, [handleRemove]);

  const getIcon = () => {
    switch (notification.type) {
      case 'success':
        return (
          <svg viewBox="0 0 24 24" fill="currentColor">
            <path d="M9 16.17L4.83 12l-1.42 1.41L9 19 21 7l-1.41-1.41z"/>
          </svg>
        );
      case 'error':
        return (
          <svg viewBox="0 0 24 24" fill="currentColor">
            <path d="M12 2C6.48 2 2 6.48 2 12s4.48 10 10 10 10-4.48 10-10S17.52 2 12 2zm-2 15l-5-5 1.41-1.41L10 14.17l7.59-7.59L19 8l-9 9z"/>
          </svg>
        );
      case 'warning':
        return (
          <svg viewBox="0 0 24 24" fill="currentColor">
            <path d="M1 21h22L12 2 1 21zm12-3h-2v-2h2v2zm0-4h-2v-4h2v4z"/>
          </svg>
        );
      case 'info':
        return (
          <svg viewBox="0 0 24 24" fill="currentColor">
            <path d="M12 2C6.48 2 2 6.48 2 12s4.48 10 10 10 10-4.48 10-10S17.52 2 12 2zm1 15h-2v-6h2v6zm0-8h-2V7h2v2z"/>
          </svg>
        );
      default:
        return null;
    }
  };

  return (
    <div 
      className={`notification-item notification-item--${notification.type} ${
        isVisible ? 'notification-item--visible' : ''
      } ${isRemoving ? 'notification-item--removing' : ''}`}
    >
      <div className="notification-item__icon">
        {getIcon()}
      </div>
      
      <div className="notification-item__content">
        <h4 className="notification-item__title">{notification.title}</h4>
        <p className="notification-item__message">{notification.message}</p>
        
        {notification.actions && notification.actions.length > 0 && (
          <div className="notification-item__actions">
            {notification.actions.map((action, index) => (
              <button
                key={index}
                onClick={() => handleAction(action)}
                className={`notification-item__action notification-item__action--${action.style || 'secondary'}`}
              >
                {action.label}
              </button>
            ))}
          </div>
        )}
      </div>
      
      <button
        onClick={handleRemove}
        className="notification-item__close"
        aria-label="Close notification"
      >
        <svg viewBox="0 0 24 24" fill="currentColor">
          <path d="M19 6.41L17.59 5 12 10.59 6.41 5 5 6.41 10.59 12 5 17.59 6.41 19 12 13.41 17.59 19 19 17.59 13.41 12z"/>
        </svg>
      </button>
    </div>
  );
};

// Helper function to get category title from error type
function getCategoryTitle(type: string): string {
  switch (type) {
    case 'validation':
      return 'Validation Error';
    case 'network':
      return 'Network Error';
    case 'filesystem':
      return 'File System Error';
    case 'database':
      return 'Database Error';
    case 'generation':
      return 'Generation Error';
    case 'system':
      return 'System Error';
    case 'permission':
      return 'Permission Error';
    case 'configuration':
      return 'Configuration Error';
    case 'authentication':
      return 'Authentication Error';
    default:
      return 'Error';
  }
}

// Hook for global error handling
export const useErrorHandler = () => {
  const { showError } = useNotifications();

  const handleError = useCallback((error: APIError | Error | string, options?: Partial<NotificationData>) => {
    return showError(error, options);
  }, [showError]);

  const handlePromiseError = useCallback((promise: Promise<any>, options?: Partial<NotificationData>) => {
    promise.catch(error => {
      handleError(error, options);
    });
  }, [handleError]);

  return {
    handleError,
    handlePromiseError
  };
};

export default NotificationProvider;