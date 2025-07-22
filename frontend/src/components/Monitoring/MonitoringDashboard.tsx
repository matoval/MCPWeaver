import React, { useState } from 'react';
import ActivityLog from './ActivityLog';
import SystemStatus from './SystemStatus';
import ErrorReports from './ErrorReports';
import './MonitoringDashboard.scss';

interface MonitoringDashboardProps {
  className?: string;
  defaultView?: 'overview' | 'logs' | 'errors' | 'status';
}

const MonitoringDashboard: React.FC<MonitoringDashboardProps> = ({
  className = '',
  defaultView = 'overview'
}) => {
  const [activeView, setActiveView] = useState(defaultView);

  const views = [
    { id: 'overview', label: 'Overview', icon: '📊' },
    { id: 'logs', label: 'Activity Log', icon: '📋' },
    { id: 'errors', label: 'Error Reports', icon: '🚨' },
    { id: 'status', label: 'System Status', icon: '💻' }
  ];

  const renderContent = () => {
    switch (activeView) {
      case 'overview':
        return (
          <div className="monitoring-dashboard__overview">
            <div className="monitoring-dashboard__row">
              <div className="monitoring-dashboard__panel monitoring-dashboard__panel--full">
                <SystemStatus showDetails={true} />
              </div>
            </div>
            
            <div className="monitoring-dashboard__row">
              <div className="monitoring-dashboard__panel monitoring-dashboard__panel--half">
                <ActivityLog 
                  height="400px"
                  showFilters={true}
                  showSearch={true}
                  showExport={false}
                  initialFilter={{ limit: 50 }}
                />
              </div>
              <div className="monitoring-dashboard__panel monitoring-dashboard__panel--half">
                <ErrorReports 
                  height="400px"
                  showResolved={false}
                  autoRefresh={true}
                />
              </div>
            </div>
          </div>
        );

      case 'logs':
        return (
          <div className="monitoring-dashboard__full-view">
            <ActivityLog 
              height="calc(100vh - 200px)"
              showFilters={true}
              showSearch={true}
              showExport={true}
              className="monitoring-dashboard__full-component"
            />
          </div>
        );

      case 'errors':
        return (
          <div className="monitoring-dashboard__full-view">
            <ErrorReports 
              height="calc(100vh - 200px)"
              showResolved={true}
              autoRefresh={true}
              className="monitoring-dashboard__full-component"
            />
          </div>
        );

      case 'status':
        return (
          <div className="monitoring-dashboard__full-view">
            <div className="monitoring-dashboard__status-grid">
              <div className="monitoring-dashboard__status-main">
                <SystemStatus 
                  showDetails={true}
                  refreshInterval={5000}
                />
              </div>
              
              <div className="monitoring-dashboard__status-side">
                <div className="monitoring-dashboard__status-compact">
                  <h4>Recent Activity</h4>
                  <ActivityLog 
                    height="300px"
                    showFilters={false}
                    showSearch={false}
                    showExport={false}
                    initialFilter={{ limit: 20, userAction: true }}
                  />
                </div>
                
                <div className="monitoring-dashboard__status-compact">
                  <h4>Recent Errors</h4>
                  <ErrorReports 
                    height="300px"
                    showResolved={false}
                    autoRefresh={true}
                  />
                </div>
              </div>
            </div>
          </div>
        );

      default:
        return null;
    }
  };

  return (
    <div className={`monitoring-dashboard ${className}`}>
      <div className="monitoring-dashboard__header">
        <h2 className="monitoring-dashboard__title">Monitoring & Observability</h2>
        
        <div className="monitoring-dashboard__nav">
          {views.map((view) => (
            <button
              key={view.id}
              onClick={() => setActiveView(view.id as any)}
              className={`monitoring-dashboard__nav-item ${
                activeView === view.id ? 'monitoring-dashboard__nav-item--active' : ''
              }`}
            >
              <span className="monitoring-dashboard__nav-icon">{view.icon}</span>
              <span className="monitoring-dashboard__nav-label">{view.label}</span>
            </button>
          ))}
        </div>
      </div>

      <div className="monitoring-dashboard__content">
        {renderContent()}
      </div>
    </div>
  );
};

export default MonitoringDashboard;