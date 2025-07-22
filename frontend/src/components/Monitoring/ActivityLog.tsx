import React, { useState, useEffect, useCallback, useMemo } from 'react';
import { 
  ActivityLogEntry, 
  LogFilter, 
  LogSearchRequest, 
  LogSearchResult, 
  LogLevel,
  LogExportRequest 
} from '../../types';
import { wails } from '../../services/wails';
import './ActivityLog.scss';

interface ActivityLogProps {
  height?: string;
  showFilters?: boolean;
  showSearch?: boolean;
  showExport?: boolean;
  initialFilter?: LogFilter;
  className?: string;
}

const ActivityLog: React.FC<ActivityLogProps> = ({
  height = '400px',
  showFilters = true,
  showSearch = true,
  showExport = true,
  initialFilter = {},
  className = ''
}) => {
  const [entries, setEntries] = useState<ActivityLogEntry[]>([]);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [filter, setFilter] = useState<LogFilter>(initialFilter);
  const [searchQuery, setSearchQuery] = useState('');
  const [searchResults, setSearchResults] = useState<LogSearchResult | null>(null);
  const [isSearchMode, setIsSearchMode] = useState(false);
  const [selectedLevel, setSelectedLevel] = useState<LogLevel | ''>('');
  const [selectedComponent, setSelectedComponent] = useState('');
  const [userActionsOnly, setUserActionsOnly] = useState(false);
  const [autoRefresh, setAutoRefresh] = useState(true);
  const [refreshInterval, setRefreshInterval] = useState(5000);

  // Log level colors for styling
  const levelColors: Record<LogLevel, string> = {
    debug: '#6B7280',
    info: '#3B82F6',
    warn: '#F59E0B',
    error: '#EF4444',
    fatal: '#DC2626'
  };

  // Load activity logs
  const loadLogs = useCallback(async () => {
    setLoading(true);
    setError(null);
    
    try {
      const currentFilter: LogFilter = {
        ...filter,
        level: selectedLevel || undefined,
        component: selectedComponent || undefined,
        userAction: userActionsOnly || undefined,
        limit: 100
      };

      const data = await wails.getActivityLogs(currentFilter);
      setEntries(data || []);
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to load activity logs');
      console.error('Failed to load activity logs:', err);
    } finally {
      setLoading(false);
    }
  }, [filter, selectedLevel, selectedComponent, userActionsOnly]);

  // Search logs
  const searchLogs = useCallback(async () => {
    if (!searchQuery.trim()) {
      setIsSearchMode(false);
      setSearchResults(null);
      return;
    }

    setLoading(true);
    setError(null);
    
    try {
      const searchRequest: LogSearchRequest = {
        query: searchQuery.trim(),
        filter: {
          level: selectedLevel || undefined,
          component: selectedComponent || undefined,
          userAction: userActionsOnly || undefined
        },
        limit: 100,
        offset: 0
      };

      const result = await wails.searchActivityLogs(searchRequest);
      setSearchResults(result);
      setIsSearchMode(true);
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to search logs');
      console.error('Failed to search logs:', err);
    } finally {
      setLoading(false);
    }
  }, [searchQuery, selectedLevel, selectedComponent, userActionsOnly]);

  // Export logs
  const exportLogs = useCallback(async (format: 'json' | 'csv' | 'txt') => {
    try {
      const timestamp = new Date().toISOString().replace(/[:.]/g, '-');
      const filename = `activity-logs-${timestamp}.${format}`;
      
      const exportRequest: LogExportRequest = {
        filter: {
          level: selectedLevel || undefined,
          component: selectedComponent || undefined,
          userAction: userActionsOnly || undefined
        },
        format,
        filePath: filename
      };

      const result = await wails.exportActivityLogs(exportRequest);
      console.log('Logs exported:', result);
      
      // Show success message (could be a notification)
      alert(`Logs exported successfully to: ${result.filePath}\nEntries: ${result.entriesCount}\nSize: ${(result.fileSize / 1024).toFixed(2)} KB`);
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to export logs');
      console.error('Failed to export logs:', err);
    }
  }, [selectedLevel, selectedComponent, userActionsOnly]);

  // Clear logs
  const clearLogs = useCallback(async (olderThanHours: number = 0) => {
    try {
      const cleared = await wails.clearActivityLogs(olderThanHours);
      console.log(`Cleared ${cleared} log entries`);
      await loadLogs(); // Reload after clearing
      
      alert(`Cleared ${cleared} log entries`);
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to clear logs');
      console.error('Failed to clear logs:', err);
    }
  }, [loadLogs]);

  // Format timestamp
  const formatTimestamp = (timestamp: string) => {
    try {
      return new Date(timestamp).toLocaleString();
    } catch {
      return timestamp;
    }
  };

  // Format duration
  const formatDuration = (duration?: string) => {
    if (!duration) return '';
    return ` (${duration})`;
  };

  // Get unique components for filter
  const uniqueComponents = useMemo(() => {
    const components = new Set(entries.map(entry => entry.component));
    return Array.from(components).sort();
  }, [entries]);

  // Auto-refresh effect
  useEffect(() => {
    if (autoRefresh && !isSearchMode) {
      const interval = setInterval(loadLogs, refreshInterval);
      return () => clearInterval(interval);
    }
  }, [autoRefresh, isSearchMode, refreshInterval, loadLogs]);

  // Initial load
  useEffect(() => {
    loadLogs();
  }, [loadLogs]);

  // Handle search input
  const handleSearchChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const value = e.target.value;
    setSearchQuery(value);
    
    if (!value.trim()) {
      setIsSearchMode(false);
      setSearchResults(null);
    }
  };

  // Handle search submit
  const handleSearchSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    searchLogs();
  };

  // Get current entries (search results or normal entries)
  const currentEntries = isSearchMode && searchResults ? searchResults.entries : entries;

  return (
    <div className={`activity-log ${className}`}>
      <div className="activity-log__header">
        <h3 className="activity-log__title">Activity Log</h3>
        
        <div className="activity-log__controls">
          {showSearch && (
            <form onSubmit={handleSearchSubmit} className="activity-log__search">
              <input
                type="text"
                value={searchQuery}
                onChange={handleSearchChange}
                placeholder="Search logs..."
                className="activity-log__search-input"
              />
              <button 
                type="submit" 
                className="activity-log__search-button"
                disabled={loading}
              >
                Search
              </button>
              {isSearchMode && (
                <button 
                  type="button" 
                  onClick={() => {
                    setSearchQuery('');
                    setIsSearchMode(false);
                    setSearchResults(null);
                  }}
                  className="activity-log__clear-search"
                >
                  Clear
                </button>
              )}
            </form>
          )}

          <div className="activity-log__actions">
            <button 
              onClick={loadLogs} 
              disabled={loading}
              className="activity-log__refresh-button"
            >
              {loading ? 'Loading...' : 'Refresh'}
            </button>
            
            <label className="activity-log__auto-refresh">
              <input
                type="checkbox"
                checked={autoRefresh}
                onChange={(e) => setAutoRefresh(e.target.checked)}
              />
              Auto-refresh
            </label>

            {showExport && (
              <div className="activity-log__export">
                <span>Export:</span>
                <button onClick={() => exportLogs('json')}>JSON</button>
                <button onClick={() => exportLogs('csv')}>CSV</button>
                <button onClick={() => exportLogs('txt')}>TXT</button>
              </div>
            )}

            <div className="activity-log__clear">
              <button onClick={() => clearLogs(24)}>Clear 24h+</button>
              <button onClick={() => clearLogs(0)}>Clear All</button>
            </div>
          </div>
        </div>
      </div>

      {showFilters && (
        <div className="activity-log__filters">
          <select
            value={selectedLevel}
            onChange={(e) => setSelectedLevel(e.target.value as LogLevel | '')}
            className="activity-log__filter"
          >
            <option value="">All Levels</option>
            <option value="debug">Debug</option>
            <option value="info">Info</option>
            <option value="warn">Warning</option>
            <option value="error">Error</option>
            <option value="fatal">Fatal</option>
          </select>

          <select
            value={selectedComponent}
            onChange={(e) => setSelectedComponent(e.target.value)}
            className="activity-log__filter"
          >
            <option value="">All Components</option>
            {uniqueComponents.map(component => (
              <option key={component} value={component}>{component}</option>
            ))}
          </select>

          <label className="activity-log__filter-checkbox">
            <input
              type="checkbox"
              checked={userActionsOnly}
              onChange={(e) => setUserActionsOnly(e.target.checked)}
            />
            User Actions Only
          </label>
        </div>
      )}

      {error && (
        <div className="activity-log__error">
          {error}
        </div>
      )}

      {isSearchMode && searchResults && (
        <div className="activity-log__search-info">
          Found {searchResults.total} entries in {searchResults.searchTime}
          {searchResults.hasMore && ' (showing first 100)'}
        </div>
      )}

      <div className="activity-log__content" style={{ height }}>
        {currentEntries.length === 0 ? (
          <div className="activity-log__empty">
            {loading ? 'Loading...' : 'No log entries found'}
          </div>
        ) : (
          <div className="activity-log__entries">
            {currentEntries.map((entry) => (
              <div
                key={entry.id}
                className={`activity-log__entry activity-log__entry--${entry.level}`}
              >
                <div className="activity-log__entry-header">
                  <span 
                    className="activity-log__entry-level"
                    style={{ color: levelColors[entry.level] }}
                  >
                    {entry.level.toUpperCase()}
                  </span>
                  <span className="activity-log__entry-timestamp">
                    {formatTimestamp(entry.timestamp)}
                  </span>
                  <span className="activity-log__entry-component">
                    {entry.component}/{entry.operation}
                  </span>
                  {entry.userAction && (
                    <span className="activity-log__entry-user-action">👤</span>
                  )}
                  {entry.duration && (
                    <span className="activity-log__entry-duration">
                      {formatDuration(entry.duration)}
                    </span>
                  )}
                </div>
                
                <div className="activity-log__entry-message">
                  {entry.message}
                </div>
                
                {entry.details && (
                  <div className="activity-log__entry-details">
                    {entry.details}
                  </div>
                )}
                
                {entry.metadata && Object.keys(entry.metadata).length > 0 && (
                  <div className="activity-log__entry-metadata">
                    <details>
                      <summary>Metadata</summary>
                      <pre>{JSON.stringify(entry.metadata, null, 2)}</pre>
                    </details>
                  </div>
                )}
              </div>
            ))}
          </div>
        )}
      </div>
    </div>
  );
};

export default ActivityLog;