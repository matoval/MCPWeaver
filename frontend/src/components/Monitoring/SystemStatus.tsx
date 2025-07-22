import React, { useState, useEffect, useCallback } from 'react';
import { ApplicationStatus, SystemHealth, StatusLevel } from '../../types';
import { wails } from '../../services/wails';
import './SystemStatus.scss';

interface SystemStatusProps {
  refreshInterval?: number;
  showDetails?: boolean;
  compact?: boolean;
  className?: string;
}

const SystemStatus: React.FC<SystemStatusProps> = ({
  refreshInterval = 10000, // 10 seconds
  showDetails = true,
  compact = false,
  className = ''
}) => {
  const [status, setStatus] = useState<ApplicationStatus | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [lastUpdated, setLastUpdated] = useState<Date | null>(null);

  // Status level colors and icons
  const statusConfig: Record<StatusLevel, { color: string; icon: string; label: string }> = {
    idle: { color: '#10b981', icon: '✓', label: 'Idle' },
    working: { color: '#3b82f6', icon: '⚡', label: 'Working' },
    warning: { color: '#f59e0b', icon: '⚠', label: 'Warning' },
    error: { color: '#ef4444', icon: '✗', label: 'Error' }
  };

  // Load application status
  const loadStatus = useCallback(async () => {
    try {
      setError(null);
      const data = await wails.getApplicationStatus();
      setStatus(data);
      setLastUpdated(new Date());
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to load status');
      console.error('Failed to load application status:', err);
    } finally {
      setLoading(false);
    }
  }, []);

  // Format memory/storage sizes
  const formatBytes = (bytes: number, unit: string = 'MB'): string => {
    return `${bytes.toFixed(1)} ${unit}`;
  };

  // Format percentage
  const formatPercentage = (value: number): string => {
    return `${value.toFixed(1)}%`;
  };

  // Get health color based on value and thresholds
  const getHealthColor = (value: number, type: 'memory' | 'cpu' | 'disk'): string => {
    const thresholds = {
      memory: { warning: 70, critical: 90 },
      cpu: { warning: 70, critical: 90 },
      disk: { warning: 10, critical: 5 } // For disk, low values are bad
    };

    const threshold = thresholds[type];
    
    if (type === 'disk') {
      // For disk space, low values are problematic
      if (value <= threshold.critical) return '#ef4444';
      if (value <= threshold.warning) return '#f59e0b';
      return '#10b981';
    } else {
      // For memory and CPU, high values are problematic
      if (value >= threshold.critical) return '#ef4444';
      if (value >= threshold.warning) return '#f59e0b';
      return '#10b981';
    }
  };

  // Auto-refresh effect
  useEffect(() => {
    loadStatus();
    
    if (refreshInterval > 0) {
      const interval = setInterval(loadStatus, refreshInterval);
      return () => clearInterval(interval);
    }
  }, [refreshInterval, loadStatus]);

  if (loading && !status) {
    return (
      <div className={`system-status system-status--loading ${className}`}>
        <div className="system-status__loading">Loading status...</div>
      </div>
    );
  }

  if (error) {
    return (
      <div className={`system-status system-status--error ${className}`}>
        <div className="system-status__error">
          <span className="system-status__error-icon">⚠</span>
          <span>Failed to load system status</span>
          <button onClick={loadStatus} className="system-status__retry">
            Retry
          </button>
        </div>
      </div>
    );
  }

  if (!status) {
    return null;
  }

  const currentStatusConfig = statusConfig[status.status];
  const health = status.systemHealth;

  if (compact) {
    return (
      <div className={`system-status system-status--compact ${className}`}>
        <div className="system-status__compact-content">
          <div 
            className="system-status__indicator"
            style={{ color: currentStatusConfig.color }}
          >
            {currentStatusConfig.icon}
          </div>
          <span className="system-status__label">{currentStatusConfig.label}</span>
          <span className="system-status__memory">
            {formatBytes(health.memoryUsage)} RAM
          </span>
          <span className="system-status__cpu">
            {formatPercentage(health.cpuUsage)} CPU
          </span>
        </div>
      </div>
    );
  }

  return (
    <div className={`system-status ${className}`}>
      <div className="system-status__header">
        <div className="system-status__main">
          <div 
            className="system-status__indicator"
            style={{ color: currentStatusConfig.color }}
          >
            {currentStatusConfig.icon}
          </div>
          <div className="system-status__info">
            <div className="system-status__label">{currentStatusConfig.label}</div>
            <div className="system-status__message">{status.message}</div>
          </div>
        </div>
        
        <div className="system-status__summary">
          <div className="system-status__operations">
            {status.activeOperations} active operations
          </div>
          {lastUpdated && (
            <div className="system-status__updated">
              Last updated: {lastUpdated.toLocaleTimeString()}
            </div>
          )}
        </div>
      </div>

      {showDetails && (
        <div className="system-status__details">
          <div className="system-status__metrics">
            <div className="system-status__metric">
              <div className="system-status__metric-label">Memory Usage</div>
              <div className="system-status__metric-value">
                <span 
                  className="system-status__metric-number"
                  style={{ color: getHealthColor(health.memoryUsage, 'memory') }}
                >
                  {formatBytes(health.memoryUsage)}
                </span>
              </div>
              <div className="system-status__metric-bar">
                <div 
                  className="system-status__metric-fill"
                  style={{ 
                    width: `${Math.min(health.memoryUsage, 100)}%`,
                    backgroundColor: getHealthColor(health.memoryUsage, 'memory')
                  }}
                />
              </div>
            </div>

            <div className="system-status__metric">
              <div className="system-status__metric-label">CPU Usage</div>
              <div className="system-status__metric-value">
                <span 
                  className="system-status__metric-number"
                  style={{ color: getHealthColor(health.cpuUsage, 'cpu') }}
                >
                  {formatPercentage(health.cpuUsage)}
                </span>
              </div>
              <div className="system-status__metric-bar">
                <div 
                  className="system-status__metric-fill"
                  style={{ 
                    width: `${Math.min(health.cpuUsage, 100)}%`,
                    backgroundColor: getHealthColor(health.cpuUsage, 'cpu')
                  }}
                />
              </div>
            </div>

            <div className="system-status__metric">
              <div className="system-status__metric-label">Disk Space</div>
              <div className="system-status__metric-value">
                <span 
                  className="system-status__metric-number"
                  style={{ color: getHealthColor(health.diskSpace, 'disk') }}
                >
                  {formatBytes(health.diskSpace, 'GB')} available
                </span>
              </div>
            </div>

            <div className="system-status__metric">
              <div className="system-status__metric-label">Database Size</div>
              <div className="system-status__metric-value">
                <span className="system-status__metric-number">
                  {formatBytes(health.databaseSize)}
                </span>
              </div>
            </div>

            <div className="system-status__metric">
              <div className="system-status__metric-label">Connections</div>
              <div className="system-status__metric-value">
                <span className="system-status__metric-number">
                  {health.activeConnections}
                </span>
              </div>
            </div>

            {health.temporaryFiles > 0 && (
              <div className="system-status__metric">
                <div className="system-status__metric-label">Temp Files</div>
                <div className="system-status__metric-value">
                  <span 
                    className="system-status__metric-number"
                    style={{ color: health.temporaryFiles > 10 ? '#f59e0b' : '#10b981' }}
                  >
                    {health.temporaryFiles}
                  </span>
                </div>
              </div>
            )}
          </div>
        </div>
      )}

      <div className="system-status__actions">
        <button 
          onClick={loadStatus}
          className="system-status__refresh"
          disabled={loading}
        >
          {loading ? 'Refreshing...' : 'Refresh'}
        </button>
      </div>
    </div>
  );
};

export default SystemStatus;