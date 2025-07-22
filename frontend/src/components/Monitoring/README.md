# Monitoring & Observability Components

This directory contains React components for the activity log and monitoring system implemented for issue #22.

## Components

### 1. ActivityLog
A comprehensive activity log viewer with real-time updates, filtering, search, and export capabilities.

**Features:**
- Real-time log streaming
- Filtering by level, component, operation, time range, and user actions
- Full-text search with pagination
- Export to JSON, CSV, and TXT formats
- Auto-refresh with configurable intervals
- Log clearing functionality
- Responsive design with dark theme support

**Usage:**
```tsx
import { ActivityLog } from '../components/Monitoring';

<ActivityLog 
  height="400px"
  showFilters={true}
  showSearch={true}
  showExport={true}
  initialFilter={{ level: 'info', limit: 100 }}
/>
```

### 2. SystemStatus
A system health monitoring component showing application status and resource usage.

**Features:**
- Real-time system health metrics (CPU, memory, disk space)
- Application status indicators (idle, working, warning, error)
- Active operations tracking
- Automatic refresh
- Compact and detailed view modes
- Health thresholds with color coding

**Usage:**
```tsx
import { SystemStatus } from '../components/Monitoring';

<SystemStatus 
  refreshInterval={10000}
  showDetails={true}
  compact={false}
/>
```

### 3. ErrorReports
An error tracking and reporting component for debugging and support.

**Features:**
- Error report listing with expandable details
- Filtering by severity and error type
- Frequency tracking for recurring errors
- Detailed system information capture
- Recovery attempt tracking
- User context preservation

**Usage:**
```tsx
import { ErrorReports } from '../components/Monitoring';

<ErrorReports 
  height="400px"
  showResolved={false}
  autoRefresh={true}
  refreshInterval={30000}
/>
```

### 4. MonitoringDashboard
A comprehensive dashboard combining all monitoring components.

**Features:**
- Multi-view layout (Overview, Activity Log, Error Reports, System Status)
- Responsive grid layout
- Navigation between views
- Integrated component management

**Usage:**
```tsx
import { MonitoringDashboard } from '../components/Monitoring';

<MonitoringDashboard 
  defaultView="overview"
  className="custom-monitoring"
/>
```

## Integration Example

To integrate monitoring into your existing application:

```tsx
// Add to your main Router or App component
import { MonitoringDashboard } from './components/Monitoring';

// Add route
<Route path="/monitoring" element={<MonitoringDashboard />} />

// Or embed in existing layouts
<div className="app-layout">
  <main className="main-content">
    {/* Your existing content */}
  </main>
  <aside className="monitoring-sidebar">
    <SystemStatus compact={true} />
    <ActivityLog 
      height="300px" 
      showFilters={false}
      showSearch={false}
      initialFilter={{ limit: 10 }}
    />
  </aside>
</div>
```

## Status Bar Integration

For a compact status display in your status bar:

```tsx
import { SystemStatus } from './components/Monitoring';

// In your StatusBar component
<div className="status-bar">
  <SystemStatus compact={true} />
  {/* Other status bar items */}
</div>
```

## API Requirements

The components expect these API methods to be available via the Wails service:

- `getActivityLogs(filter)` - Retrieve activity logs with filtering
- `searchActivityLogs(request)` - Search through activity logs
- `exportActivityLogs(request)` - Export logs to file
- `getApplicationStatus()` - Get current application status
- `updateLogConfig(config)` - Update logging configuration
- `clearActivityLogs(olderThanHours)` - Clear old log entries
- `createErrorReport(...)` - Create new error report
- `getErrorReports(includeResolved)` - Get error reports

## Real-time Updates

The components support real-time updates through Wails events:

- `activity:log` - New activity log entry
- `system:status` - Application status change
- `error:report` - New error report

## Theming

All components support dark theme through CSS custom properties:

```css
[data-theme="dark"] {
  --bg-primary: #1a202c;
  --bg-secondary: #2d3748;
  --text-primary: #f7fafc;
  --text-secondary: #a0aec0;
  /* ... other theme variables */
}
```

## Accessibility

Components include:
- Keyboard navigation support
- Screen reader friendly markup
- High contrast mode support
- Reduced motion preferences
- Focus management

## Performance

- Virtualized scrolling for large log lists
- Debounced search input
- Memoized filtering and sorting
- Efficient re-rendering with React hooks
- Configurable refresh intervals

## Browser Support

- Modern browsers (Chrome 90+, Firefox 88+, Safari 14+)
- Progressive enhancement for older browsers
- Responsive design for mobile devices