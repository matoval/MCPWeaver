import React from 'react';
import { 
  CheckCircle, 
  AlertCircle, 
  Clock, 
  Activity,
  Cpu,
  HardDrive 
} from 'lucide-react';
import './StatusBar.scss';

interface SystemHealth {
  memoryUsage: number;
  cpuUsage: number;
  diskUsage: number;
  status: 'healthy' | 'warning' | 'error';
}

interface StatusBarProps {
  status: 'ready' | 'working' | 'error';
  activeOperations: number;
  systemHealth: SystemHealth;
  onStatusClick: () => void;
}

const StatusBar: React.FC<StatusBarProps> = ({ 
  status, 
  activeOperations, 
  systemHealth, 
  onStatusClick 
}) => {
  const getStatusIcon = () => {
    switch (status) {
      case 'ready':
        return <CheckCircle className="status-icon ready" size={14} />;
      case 'working':
        return <Clock className="status-icon working" size={14} />;
      case 'error':
        return <AlertCircle className="status-icon error" size={14} />;
      default:
        return <CheckCircle className="status-icon ready" size={14} />;
    }
  };

  const getStatusText = () => {
    switch (status) {
      case 'ready':
        return 'Ready';
      case 'working':
        return 'Working...';
      case 'error':
        return 'Error';
      default:
        return 'Ready';
    }
  };

  const getHealthColor = (health: SystemHealth['status']) => {
    switch (health) {
      case 'healthy':
        return 'var(--success-500)';
      case 'warning':
        return 'var(--warning-500)';
      case 'error':
        return 'var(--error-500)';
      default:
        return 'var(--success-500)';
    }
  };

  return (
    <div className="status-bar">
      <button className="status-indicator" onClick={onStatusClick}>
        {getStatusIcon()}
        <span className="status-text">{getStatusText()}</span>
      </button>

      {activeOperations > 0 && (
        <div className="operation-counter">
          <Activity size={14} />
          <span>{activeOperations} active</span>
        </div>
      )}

      <div className="status-spacer" />

      <div className="resource-usage">
        <div className="resource-item">
          <Cpu size={14} />
          <span>{systemHealth.cpuUsage.toFixed(1)}%</span>
        </div>
        <div className="resource-item">
          <HardDrive size={14} />
          <span>{systemHealth.memoryUsage.toFixed(1)}%</span>
        </div>
      </div>

      <div 
        className="system-health-indicator"
        style={{ backgroundColor: getHealthColor(systemHealth.status) }}
        title={`System Health: ${systemHealth.status}`}
      />

      <div className="app-version">
        v1.0.0
      </div>
    </div>
  );
};

export default StatusBar;