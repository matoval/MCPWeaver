import React, { useState, useEffect, useCallback } from 'react';
import { wails } from '../../services/wails';
import { Notification, SystemNotificationEvent } from '../../types';
import './NotificationSystem.scss';

export interface NotificationSystemProps {
  maxNotifications?: number;
  autoHideDuration?: number;
  position?: 'top-right' | 'top-left' | 'bottom-right' | 'bottom-left';
  className?: string;
}

export const NotificationSystem: React.FC<NotificationSystemProps> = ({
  maxNotifications = 5,
  autoHideDuration = 5000,
  position = 'top-right',
  className = ''
}) => {
  const [notifications, setNotifications] = useState<Notification[]>([]);

  // Handle system notifications
  const handleSystemNotification = useCallback((event: SystemNotificationEvent) => {
    const notification = {
      ...event.data,
      id: event.data.id || `notification-${Date.now()}`,
      timestamp: event.data.timestamp || new Date().toISOString(),
      read: false
    };

    setNotifications(prev => [notification, ...prev].slice(0, maxNotifications));

    // Auto-hide after duration (except for error notifications)
    if (notification.type !== 'error' && autoHideDuration > 0) {
      setTimeout(() => {
        setNotifications(prev => prev.filter(n => n.id !== notification.id));
      }, autoHideDuration);
    }
  }, [maxNotifications, autoHideDuration]);

  // Set up event listener
  useEffect(() => {
    const cleanup = wails.onEvent('system:notification', handleSystemNotification);
    return cleanup;
  }, [handleSystemNotification]);

  // Dismiss notification
  const dismissNotification = useCallback((id: string) => {
    setNotifications(prev => prev.filter(n => n.id !== id));
  }, []);

  // Mark notification as read
  const markAsRead = useCallback((id: string) => {
    setNotifications(prev => prev.map(n => 
      n.id === id ? { ...n, read: true } : n
    ));
  }, []);

  // Clear all notifications
  const clearAll = useCallback(() => {
    setNotifications([]);
  }, []);

  // Get notification icon
  const getNotificationIcon = (type: string): React.ReactNode => {
    switch (type) {
      case 'success':
        return (
          <svg viewBox="0 0 20 20" fill="currentColor" className="notification__icon">
            <path fillRule="evenodd" d="M16.707 5.293a1 1 0 010 1.414l-8 8a1 1 0 01-1.414 0l-4-4a1 1 0 011.414-1.414L8 12.586l7.293-7.293a1 1 0 011.414 0z" clipRule="evenodd" />
          </svg>
        );
      case 'error':
        return (
          <svg viewBox="0 0 20 20" fill="currentColor" className="notification__icon">
            <path fillRule="evenodd" d="M18 10a8 8 0 11-16 0 8 8 0 0116 0zm-7 4a1 1 0 11-2 0 1 1 0 012 0zm-1-9a1 1 0 00-1 1v4a1 1 0 102 0V6a1 1 0 00-1-1z" clipRule="evenodd" />
          </svg>
        );
      case 'warning':
        return (
          <svg viewBox="0 0 20 20" fill="currentColor" className="notification__icon">
            <path fillRule="evenodd" d="M8.257 3.099c.765-1.36 2.722-1.36 3.486 0l5.58 9.92c.75 1.334-.213 2.98-1.742 2.98H4.42c-1.53 0-2.493-1.646-1.743-2.98l5.58-9.92zM11 13a1 1 0 11-2 0 1 1 0 012 0zm-1-8a1 1 0 00-1 1v3a1 1 0 002 0V6a1 1 0 00-1-1z" clipRule="evenodd" />
          </svg>
        );
      default:
        return (
          <svg viewBox="0 0 20 20" fill="currentColor" className="notification__icon">
            <path fillRule="evenodd" d="M18 10a8 8 0 11-16 0 8 8 0 0116 0zm-7-4a1 1 0 11-2 0 1 1 0 012 0zM9 9a1 1 0 000 2v3a1 1 0 001 1h1a1 1 0 100-2v-3a1 1 0 00-1-1H9z" clipRule="evenodd" />
          </svg>
        );
    }
  };

  // Format timestamp
  const formatTimestamp = (timestamp: string): string => {
    const date = new Date(timestamp);
    const now = new Date();
    const diff = now.getTime() - date.getTime();
    
    if (diff < 60000) { // Less than 1 minute
      return 'Just now';
    } else if (diff < 3600000) { // Less than 1 hour
      const minutes = Math.floor(diff / 60000);
      return `${minutes}m ago`;
    } else if (diff < 86400000) { // Less than 1 day
      const hours = Math.floor(diff / 3600000);
      return `${hours}h ago`;
    } else {
      return date.toLocaleDateString();
    }
  };

  if (notifications.length === 0) {
    return null;
  }

  return (
    <div className={`notification-system ${className} notification-system--${position}`}>
      <div className="notification-system__header">
        <h3>Notifications</h3>
        {notifications.length > 1 && (
          <button 
            className="notification-system__clear-all"
            onClick={clearAll}
            title="Clear all notifications"
          >
            Clear all
          </button>
        )}
      </div>
      
      <div className="notification-system__list">
        {notifications.map((notification) => (
          <div
            key={notification.id}
            className={`notification notification--${notification.type} ${notification.read ? 'notification--read' : ''}`}
            onClick={() => markAsRead(notification.id)}
            role="alert"
            aria-live="polite"
          >
            <div className="notification__icon-container">
              {getNotificationIcon(notification.type)}
            </div>
            
            <div className="notification__content">
              <div className="notification__header">
                <h4 className="notification__title">{notification.title}</h4>
                <div className="notification__meta">
                  <span className="notification__timestamp">
                    {formatTimestamp(notification.timestamp)}
                  </span>
                  <button
                    className="notification__dismiss"
                    onClick={(e) => {
                      e.stopPropagation();
                      dismissNotification(notification.id);
                    }}
                    title="Dismiss notification"
                  >
                    <svg viewBox="0 0 20 20" fill="currentColor">
                      <path fillRule="evenodd" d="M4.293 4.293a1 1 0 011.414 0L10 8.586l4.293-4.293a1 1 0 111.414 1.414L11.414 10l4.293 4.293a1 1 0 01-1.414 1.414L10 11.414l-4.293 4.293a1 1 0 01-1.414-1.414L8.586 10 4.293 5.707a1 1 0 010-1.414z" clipRule="evenodd" />
                    </svg>
                  </button>
                </div>
              </div>
              
              <p className="notification__message">{notification.message}</p>
              
              {notification.actions && notification.actions.length > 0 && (
                <div className="notification__actions">
                  {notification.actions.map((action, index) => (
                    <button
                      key={index}
                      className="notification__action"
                      onClick={(e) => {
                        e.stopPropagation();
                        // Handle action click
                        console.log('Action clicked:', action);
                      }}
                    >
                      {action}
                    </button>
                  ))}
                </div>
              )}
            </div>
          </div>
        ))}
      </div>
    </div>
  );
};

export default NotificationSystem;