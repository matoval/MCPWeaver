import React, { useState, useEffect, useCallback } from 'react';
import { ErrorReport, ErrorSeverity, ErrorType } from '../../types';
import { wails } from '../../services/wails';
import './ErrorReports.scss';

interface ErrorReportsProps {
  height?: string;
  showResolved?: boolean;
  autoRefresh?: boolean;
  refreshInterval?: number;
  className?: string;
}

const ErrorReports: React.FC<ErrorReportsProps> = ({
  height = '400px',
  showResolved = false,
  autoRefresh = true,
  refreshInterval = 30000, // 30 seconds
  className = ''
}) => {
  const [reports, setReports] = useState<ErrorReport[]>([]);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [includeResolved, setIncludeResolved] = useState(showResolved);
  const [selectedSeverity, setSelectedSeverity] = useState<ErrorSeverity | ''>('');
  const [selectedType, setSelectedType] = useState<ErrorType | ''>('');
  const [expandedReport, setExpandedReport] = useState<string | null>(null);

  // Severity colors and icons
  const severityConfig: Record<ErrorSeverity, { color: string; icon: string; label: string }> = {
    low: { color: '#10b981', icon: 'ℹ', label: 'Low' },
    medium: { color: '#f59e0b', icon: '⚠', label: 'Medium' },
    high: { color: '#ef4444', icon: '⚠', label: 'High' },
    critical: { color: '#dc2626', icon: '🚨', label: 'Critical' }
  };

  // Error type labels
  const typeLabels: Record<ErrorType, string> = {
    validation: 'Validation',
    system: 'System',
    network: 'Network',
    filesystem: 'File System',
    database: 'Database',
    generation: 'Generation',
    permission: 'Permission',
    configuration: 'Configuration',
    authentication: 'Authentication'
  };

  // Load error reports
  const loadReports = useCallback(async () => {
    setLoading(true);
    setError(null);
    
    try {
      const data = await wails.getErrorReports(includeResolved);
      setReports(data || []);
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to load error reports');
      console.error('Failed to load error reports:', err);
    } finally {
      setLoading(false);
    }
  }, [includeResolved]);

  // Filter reports based on selected criteria
  const filteredReports = reports.filter(report => {
    if (selectedSeverity && report.severity !== selectedSeverity) return false;
    if (selectedType && report.type !== selectedType) return false;
    return true;
  });

  // Format timestamp
  const formatTimestamp = (timestamp: string) => {
    try {
      return new Date(timestamp).toLocaleString();
    } catch {
      return timestamp;
    }
  };

  // Format time ago
  const formatTimeAgo = (timestamp: string) => {
    try {
      const date = new Date(timestamp);
      const now = new Date();
      const diffMs = now.getTime() - date.getTime();
      const diffMins = Math.floor(diffMs / 60000);
      const diffHours = Math.floor(diffMins / 60);
      const diffDays = Math.floor(diffHours / 24);

      if (diffMins < 1) return 'Just now';
      if (diffMins < 60) return `${diffMins}m ago`;
      if (diffHours < 24) return `${diffHours}h ago`;
      return `${diffDays}d ago`;
    } catch {
      return timestamp;
    }
  };

  // Toggle expanded report
  const toggleExpanded = (reportId: string) => {
    setExpandedReport(expandedReport === reportId ? null : reportId);
  };

  // Get unique error types for filter
  const uniqueTypes = [...new Set(reports.map(report => report.type))].sort();

  // Auto-refresh effect
  useEffect(() => {
    if (autoRefresh) {
      const interval = setInterval(loadReports, refreshInterval);
      return () => clearInterval(interval);
    }
  }, [autoRefresh, refreshInterval, loadReports]);

  // Initial load
  useEffect(() => {
    loadReports();
  }, [loadReports]);

  return (
    <div className={`error-reports ${className}`}>
      <div className="error-reports__header">
        <h3 className="error-reports__title">Error Reports</h3>
        
        <div className="error-reports__controls">
          <div className="error-reports__filters">
            <select
              value={selectedSeverity}
              onChange={(e) => setSelectedSeverity(e.target.value as ErrorSeverity | '')}
              className="error-reports__filter"
            >
              <option value="">All Severities</option>
              <option value="low">Low</option>
              <option value="medium">Medium</option>
              <option value="high">High</option>
              <option value="critical">Critical</option>
            </select>

            <select
              value={selectedType}
              onChange={(e) => setSelectedType(e.target.value as ErrorType | '')}
              className="error-reports__filter"
            >
              <option value="">All Types</option>
              {uniqueTypes.map(type => (
                <option key={type} value={type}>{typeLabels[type]}</option>
              ))}
            </select>

            <label className="error-reports__filter-checkbox">
              <input
                type="checkbox"
                checked={includeResolved}
                onChange={(e) => setIncludeResolved(e.target.checked)}
              />
              Include Resolved
            </label>
          </div>

          <div className="error-reports__actions">
            <button 
              onClick={loadReports} 
              disabled={loading}
              className="error-reports__refresh-button"
            >
              {loading ? 'Loading...' : 'Refresh'}
            </button>
          </div>
        </div>
      </div>

      {error && (
        <div className="error-reports__error">
          {error}
        </div>
      )}

      <div className="error-reports__stats">
        <div className="error-reports__stat">
          <span className="error-reports__stat-value">{filteredReports.length}</span>
          <span className="error-reports__stat-label">Total Reports</span>
        </div>
        {Object.entries(severityConfig).map(([severity, config]) => {
          const count = filteredReports.filter(r => r.severity === severity).length;
          if (count === 0) return null;
          
          return (
            <div key={severity} className="error-reports__stat">
              <span 
                className="error-reports__stat-value"
                style={{ color: config.color }}
              >
                {count}
              </span>
              <span className="error-reports__stat-label">{config.label}</span>
            </div>
          );
        })}
      </div>

      <div className="error-reports__content" style={{ height }}>
        {filteredReports.length === 0 ? (
          <div className="error-reports__empty">
            {loading ? 'Loading...' : 'No error reports found'}
          </div>
        ) : (
          <div className="error-reports__list">
            {filteredReports.map((report) => {
              const severityConf = severityConfig[report.severity];
              const isExpanded = expandedReport === report.id;

              return (
                <div
                  key={report.id}
                  className={`error-reports__item error-reports__item--${report.severity}`}
                >
                  <div 
                    className="error-reports__item-header"
                    onClick={() => toggleExpanded(report.id)}
                  >
                    <div className="error-reports__item-main">
                      <div className="error-reports__item-severity">
                        <span 
                          className="error-reports__severity-icon"
                          style={{ color: severityConf.color }}
                        >
                          {severityConf.icon}
                        </span>
                        <span className="error-reports__severity-label">
                          {severityConf.label}
                        </span>
                      </div>
                      
                      <div className="error-reports__item-info">
                        <div className="error-reports__item-type">
                          {typeLabels[report.type]}
                        </div>
                        <div className="error-reports__item-component">
                          {report.component}/{report.operation}
                        </div>
                      </div>
                      
                      <div className="error-reports__item-message">
                        {report.message}
                      </div>
                    </div>

                    <div className="error-reports__item-meta">
                      <div className="error-reports__item-time">
                        {formatTimeAgo(report.timestamp)}
                      </div>
                      {report.frequency > 1 && (
                        <div className="error-reports__item-frequency">
                          {report.frequency}x
                        </div>
                      )}
                      <div className="error-reports__item-expand">
                        {isExpanded ? '▼' : '▶'}
                      </div>
                    </div>
                  </div>

                  {isExpanded && (
                    <div className="error-reports__item-details">
                      <div className="error-reports__detail-section">
                        <h4>Details</h4>
                        <div className="error-reports__detail-grid">
                          <div className="error-reports__detail-item">
                            <label>Timestamp:</label>
                            <span>{formatTimestamp(report.timestamp)}</span>
                          </div>
                          <div className="error-reports__detail-item">
                            <label>First Seen:</label>
                            <span>{formatTimestamp(report.firstSeen)}</span>
                          </div>
                          <div className="error-reports__detail-item">
                            <label>Last Seen:</label>
                            <span>{formatTimestamp(report.lastSeen)}</span>
                          </div>
                          <div className="error-reports__detail-item">
                            <label>Frequency:</label>
                            <span>{report.frequency}</span>
                          </div>
                        </div>
                      </div>

                      {report.details && (
                        <div className="error-reports__detail-section">
                          <h4>Error Details</h4>
                          <pre className="error-reports__detail-text">
                            {report.details}
                          </pre>
                        </div>
                      )}

                      {report.stackTrace && (
                        <div className="error-reports__detail-section">
                          <h4>Stack Trace</h4>
                          <pre className="error-reports__detail-text">
                            {report.stackTrace}
                          </pre>
                        </div>
                      )}

                      <div className="error-reports__detail-section">
                        <h4>System Information</h4>
                        <div className="error-reports__detail-grid">
                          <div className="error-reports__detail-item">
                            <label>OS:</label>
                            <span>{report.systemInfo.os} ({report.systemInfo.architecture})</span>
                          </div>
                          <div className="error-reports__detail-item">
                            <label>Go Version:</label>
                            <span>{report.systemInfo.goVersion}</span>
                          </div>
                          <div className="error-reports__detail-item">
                            <label>App Version:</label>
                            <span>{report.systemInfo.appVersion}</span>
                          </div>
                          <div className="error-reports__detail-item">
                            <label>Memory:</label>
                            <span>{report.systemInfo.memoryMB.toFixed(1)} MB</span>
                          </div>
                          <div className="error-reports__detail-item">
                            <label>CPU Usage:</label>
                            <span>{report.systemInfo.cpuUsage.toFixed(1)}%</span>
                          </div>
                          <div className="error-reports__detail-item">
                            <label>Disk Space:</label>
                            <span>{report.systemInfo.diskSpaceGB.toFixed(1)} GB</span>
                          </div>
                        </div>
                      </div>

                      {report.userContext && (
                        <div className="error-reports__detail-section">
                          <h4>User Context</h4>
                          <div className="error-reports__detail-grid">
                            {report.userContext.projectName && (
                              <div className="error-reports__detail-item">
                                <label>Project:</label>
                                <span>{report.userContext.projectName}</span>
                              </div>
                            )}
                            {report.userContext.userAction && (
                              <div className="error-reports__detail-item">
                                <label>User Action:</label>
                                <span>{report.userContext.userAction}</span>
                              </div>
                            )}
                            {report.userContext.uiState && (
                              <div className="error-reports__detail-item">
                                <label>UI State:</label>
                                <span>{report.userContext.uiState}</span>
                              </div>
                            )}
                          </div>
                        </div>
                      )}

                      <div className="error-reports__detail-section">
                        <h4>Recovery Information</h4>
                        <div className="error-reports__detail-grid">
                          <div className="error-reports__detail-item">
                            <label>Recovery Attempted:</label>
                            <span>{report.recovery.attempted ? 'Yes' : 'No'}</span>
                          </div>
                          {report.recovery.attempted && (
                            <>
                              <div className="error-reports__detail-item">
                                <label>Recovery Successful:</label>
                                <span>{report.recovery.successful ? 'Yes' : 'No'}</span>
                              </div>
                              {report.recovery.method && (
                                <div className="error-reports__detail-item">
                                  <label>Recovery Method:</label>
                                  <span>{report.recovery.method}</span>
                                </div>
                              )}
                              <div className="error-reports__detail-item">
                                <label>Data Loss:</label>
                                <span>{report.recovery.dataLoss ? 'Yes' : 'No'}</span>
                              </div>
                            </>
                          )}
                        </div>
                      </div>
                    </div>
                  )}
                </div>
              );
            })}
          </div>
        )}
      </div>
    </div>
  );
};

export default ErrorReports;