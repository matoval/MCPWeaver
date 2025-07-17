# MCPWeaver Observability and Monitoring Specification

## Overview

This document defines the observability and monitoring strategy for MCPWeaver, focusing on user-facing visibility into application state, performance metrics, and system health while maintaining simplicity and lightweight operation.

## Observability Principles

### User-Centric Observability
- **Transparency**: Users can see what the application is doing at all times
- **Clarity**: Status information is clear and actionable
- **Non-Intrusive**: Observability doesn't impact performance or user experience
- **Actionable**: Users can take appropriate actions based on the information provided

### Lightweight Design
- **Minimal Overhead**: < 1% performance impact
- **Local Storage**: All metrics stored locally (no external dependencies)
- **Configurable**: Users can adjust verbosity and retention
- **Privacy-First**: No external telemetry unless explicitly enabled

## User-Facing Observability

### Real-Time Status Indicators

#### Application Status Bar
```go
type ApplicationStatus struct {
    Status           StatusLevel    `json:"status"`
    Message          string         `json:"message"`
    ActiveOperations int            `json:"activeOperations"`
    LastUpdate       time.Time      `json:"lastUpdate"`
    SystemHealth     SystemHealth   `json:"systemHealth"`
}

type StatusLevel string

const (
    StatusIdle    StatusLevel = "idle"
    StatusWorking StatusLevel = "working"  
    StatusError   StatusLevel = "error"
    StatusWarning StatusLevel = "warning"
)

type SystemHealth struct {
    MemoryUsage     float64 `json:"memoryUsage"`     // MB
    CPUUsage        float64 `json:"cpuUsage"`        // Percentage
    DiskSpace       float64 `json:"diskSpace"`       // GB available
    DatabaseSize    float64 `json:"databaseSize"`    // MB
    TemporaryFiles  int     `json:"temporaryFiles"`  // Count
    ActiveConnections int   `json:"activeConnections"`
}
```

#### Generation Progress Tracking
```go
type GenerationProgress struct {
    JobID           string            `json:"jobId"`
    ProjectName     string            `json:"projectName"`
    Status          GenerationStatus  `json:"status"`
    Progress        float64           `json:"progress"`       // 0.0 to 1.0
    CurrentStep     string            `json:"currentStep"`
    StepDetails     string            `json:"stepDetails"`
    ElapsedTime     time.Duration     `json:"elapsedTime"`
    EstimatedTime   time.Duration     `json:"estimatedTime"`
    ProcessingRate  float64           `json:"processingRate"` // endpoints/sec
    Metrics         ProgressMetrics   `json:"metrics"`
}

type ProgressMetrics struct {
    EndpointsProcessed int     `json:"endpointsProcessed"`
    EndpointsTotal     int     `json:"endpointsTotal"`
    FilesGenerated     int     `json:"filesGenerated"`
    LinesGenerated     int     `json:"linesGenerated"`
    ErrorsEncountered  int     `json:"errorsEncountered"`
    WarningsIssued     int     `json:"warningsIssued"`
    MemoryUsed         float64 `json:"memoryUsed"`        // MB
    TempFilesCreated   int     `json:"tempFilesCreated"`
}
```

#### Real-Time Activity Log
```go
type ActivityLogEntry struct {
    ID          string        `json:"id"`
    Timestamp   time.Time     `json:"timestamp"`
    Level       LogLevel      `json:"level"`
    Component   string        `json:"component"`
    Operation   string        `json:"operation"`
    Message     string        `json:"message"`
    Details     string        `json:"details"`
    Duration    time.Duration `json:"duration"`
    ProjectID   string        `json:"projectId"`
    UserAction  bool          `json:"userAction"`
    Metadata    map[string]interface{} `json:"metadata"`
}

type LogLevel string

const (
    LogLevelTrace LogLevel = "trace"
    LogLevelDebug LogLevel = "debug"
    LogLevelInfo  LogLevel = "info"
    LogLevelWarn  LogLevel = "warn"
    LogLevelError LogLevel = "error"
)
```

### Performance Metrics Dashboard

#### Application Performance Metrics
```go
type PerformanceMetrics struct {
    // Application Performance
    StartupTime         time.Duration     `json:"startupTime"`
    AverageResponseTime time.Duration     `json:"averageResponseTime"`
    MemoryUsage         MemoryStats       `json:"memoryUsage"`
    CPUUsage            CPUStats          `json:"cpuUsage"`
    
    // Generation Performance
    GenerationStats     GenerationPerfStats `json:"generationStats"`
    
    // File Operations
    FileOperations      FileOpStats         `json:"fileOperations"`
    
    // Database Performance
    DatabaseStats       DatabaseStats       `json:"databaseStats"`
    
    // UI Performance
    UIPerformance       UIStats             `json:"uiPerformance"`
    
    // Collection Period
    CollectionStart     time.Time           `json:"collectionStart"`
    CollectionEnd       time.Time           `json:"collectionEnd"`
}

type MemoryStats struct {
    AllocatedMB      float64 `json:"allocatedMB"`
    SystemMB         float64 `json:"systemMB"`
    HeapMB           float64 `json:"heapMB"`
    StackMB          float64 `json:"stackMB"`
    GCCount          int     `json:"gcCount"`
    GCPauseTotal     time.Duration `json:"gcPauseTotal"`
    GCPauseAverage   time.Duration `json:"gcPauseAverage"`
}

type CPUStats struct {
    UsagePercent     float64 `json:"usagePercent"`
    Cores            int     `json:"cores"`
    ProcessTime      time.Duration `json:"processTime"`
    SystemTime       time.Duration `json:"systemTime"`
    IdleTime         time.Duration `json:"idleTime"`
}

type GenerationPerfStats struct {
    TotalGenerations    int             `json:"totalGenerations"`
    SuccessfulGenerations int           `json:"successfulGenerations"`
    FailedGenerations   int             `json:"failedGenerations"`
    AverageTime         time.Duration   `json:"averageTime"`
    FastestTime         time.Duration   `json:"fastestTime"`
    SlowestTime         time.Duration   `json:"slowestTime"`
    ThroughputEPS       float64         `json:"throughputEPS"` // endpoints per second
    ComplexityBreakdown map[string]int  `json:"complexityBreakdown"`
}

type FileOpStats struct {
    ReadsTotal      int64         `json:"readsTotal"`
    WritesTotal     int64         `json:"writesTotal"`
    ReadsMB         float64       `json:"readsMB"`
    WrittenMB       float64       `json:"writtenMB"`
    AverageReadTime time.Duration `json:"averageReadTime"`
    AverageWriteTime time.Duration `json:"averageWriteTime"`
    ErrorCount      int           `json:"errorCount"`
}

type DatabaseStats struct {
    QueriesTotal      int64         `json:"queriesTotal"`
    QueryTimeAverage  time.Duration `json:"queryTimeAverage"`
    QueryTimeSlowest  time.Duration `json:"queryTimeSlowest"`
    ConnectionsActive int           `json:"connectionsActive"`
    ConnectionsMax    int           `json:"connectionsMax"`
    CacheHitRate      float64       `json:"cacheHitRate"`
    SizeMB           float64       `json:"sizeMB"`
}

type UIStats struct {
    RenderTimeAverage  time.Duration `json:"renderTimeAverage"`
    RenderTimeSlowest  time.Duration `json:"renderTimeSlowest"`
    ComponentUpdates   int           `json:"componentUpdates"`
    EventsProcessed    int           `json:"eventsProcessed"`
    MemoryLeaks        int           `json:"memoryLeaks"`
    ErrorBoundaryHits  int           `json:"errorBoundaryHits"`
}
```

### Error Tracking and Reporting

#### Error Classification
```go
type ErrorReport struct {
    ID           string        `json:"id"`
    Timestamp    time.Time     `json:"timestamp"`
    Type         ErrorType     `json:"type"`
    Severity     ErrorSeverity `json:"severity"`
    Component    string        `json:"component"`
    Operation    string        `json:"operation"`
    Message      string        `json:"message"`
    Details      string        `json:"details"`
    StackTrace   string        `json:"stackTrace"`
    UserContext  UserContext   `json:"userContext"`
    SystemInfo   SystemInfo    `json:"systemInfo"`
    Recovery     RecoveryInfo  `json:"recovery"`
    Frequency    int           `json:"frequency"`
    FirstSeen    time.Time     `json:"firstSeen"`
    LastSeen     time.Time     `json:"lastSeen"`
}

type ErrorSeverity string

const (
    SeverityCritical ErrorSeverity = "critical"
    SeverityHigh     ErrorSeverity = "high"
    SeverityMedium   ErrorSeverity = "medium"
    SeverityLow      ErrorSeverity = "low"
)

type UserContext struct {
    ProjectID       string            `json:"projectId"`
    ProjectName     string            `json:"projectName"`
    UserAction      string            `json:"userAction"`
    UIState         string            `json:"uiState"`
    RecentActions   []string          `json:"recentActions"`
    Settings        map[string]string `json:"settings"`
}

type SystemInfo struct {
    OS              string  `json:"os"`
    Architecture    string  `json:"architecture"`
    GoVersion       string  `json:"goVersion"`
    AppVersion      string  `json:"appVersion"`
    MemoryMB        float64 `json:"memoryMB"`
    CPUUsage        float64 `json:"cpuUsage"`
    DiskSpaceGB     float64 `json:"diskSpaceGB"`
    DatabaseSizeMB  float64 `json:"databaseSizeMB"`
}

type RecoveryInfo struct {
    Attempted       bool          `json:"attempted"`
    Successful      bool          `json:"successful"`
    Method          string        `json:"method"`
    Duration        time.Duration `json:"duration"`
    UserInteraction bool          `json:"userInteraction"`
    DataLoss        bool          `json:"dataLoss"`
}
```

## Internal Observability

### Structured Logging

#### Log Configuration
```go
type LogConfig struct {
    Level           LogLevel          `json:"level"`
    Output          LogOutput         `json:"output"`
    Format          LogFormat         `json:"format"`
    Rotation        LogRotation       `json:"rotation"`
    Components      map[string]LogLevel `json:"components"`
    EnableConsole   bool              `json:"enableConsole"`
    EnableFile      bool              `json:"enableFile"`
    EnableBuffer    bool              `json:"enableBuffer"`
    BufferSize      int               `json:"bufferSize"`
    FlushInterval   time.Duration     `json:"flushInterval"`
}

type LogOutput string

const (
    LogOutputConsole LogOutput = "console"
    LogOutputFile    LogOutput = "file"
    LogOutputBuffer  LogOutput = "buffer"
    LogOutputAll     LogOutput = "all"
)

type LogFormat string

const (
    LogFormatJSON LogFormat = "json"
    LogFormatText LogFormat = "text"
    LogFormatHuman LogFormat = "human"
)

type LogRotation struct {
    MaxSizeMB    int           `json:"maxSizeMB"`
    MaxAge       time.Duration `json:"maxAge"`
    MaxBackups   int           `json:"maxBackups"`
    Compress     bool          `json:"compress"`
}
```

#### Structured Log Entry
```go
type LogEntry struct {
    Timestamp   time.Time              `json:"timestamp"`
    Level       LogLevel               `json:"level"`
    Component   string                 `json:"component"`
    Operation   string                 `json:"operation"`
    Message     string                 `json:"message"`
    Fields      map[string]interface{} `json:"fields"`
    Duration    time.Duration          `json:"duration,omitempty"`
    Error       string                 `json:"error,omitempty"`
    StackTrace  string                 `json:"stackTrace,omitempty"`
    TraceID     string                 `json:"traceId,omitempty"`
    SpanID      string                 `json:"spanId,omitempty"`
    UserID      string                 `json:"userId,omitempty"`
    ProjectID   string                 `json:"projectId,omitempty"`
    RequestID   string                 `json:"requestId,omitempty"`
}
```

### Metrics Collection

#### Metrics Registry
```go
type MetricsRegistry struct {
    counters   map[string]*Counter
    gauges     map[string]*Gauge
    histograms map[string]*Histogram
    timers     map[string]*Timer
    mutex      sync.RWMutex
}

type Counter struct {
    name        string
    value       int64
    labels      map[string]string
    description string
    created     time.Time
}

type Gauge struct {
    name        string
    value       float64
    labels      map[string]string
    description string
    updated     time.Time
}

type Histogram struct {
    name        string
    buckets     []float64
    counts      []int64
    sum         float64
    count       int64
    labels      map[string]string
    description string
}

type Timer struct {
    name        string
    duration    time.Duration
    count       int64
    sum         time.Duration
    min         time.Duration
    max         time.Duration
    labels      map[string]string
    description string
}
```

#### Key Metrics
```go
// Application Metrics
const (
    MetricAppStartupTime     = "app_startup_time"
    MetricAppMemoryUsage     = "app_memory_usage"
    MetricAppCPUUsage        = "app_cpu_usage"
    MetricAppActiveUsers     = "app_active_users"
    MetricAppUptime          = "app_uptime"
)

// Generation Metrics
const (
    MetricGenerationTotal     = "generation_total"
    MetricGenerationSuccess   = "generation_success"
    MetricGenerationFailure   = "generation_failure"
    MetricGenerationDuration  = "generation_duration"
    MetricGenerationQueueSize = "generation_queue_size"
)

// File Operation Metrics
const (
    MetricFileReads          = "file_reads_total"
    MetricFileWrites         = "file_writes_total"
    MetricFileReadDuration   = "file_read_duration"
    MetricFileWriteDuration  = "file_write_duration"
    MetricFileErrors         = "file_errors_total"
)

// Database Metrics
const (
    MetricDatabaseQueries     = "database_queries_total"
    MetricDatabaseConnections = "database_connections_active"
    MetricDatabaseQueryTime   = "database_query_duration"
    MetricDatabaseErrors      = "database_errors_total"
)

// UI Metrics
const (
    MetricUIRenderTime       = "ui_render_duration"
    MetricUIComponentUpdates = "ui_component_updates"
    MetricUIEventProcessing  = "ui_event_processing"
    MetricUIErrors           = "ui_errors_total"
)
```

### Health Checks

#### Health Check System
```go
type HealthChecker struct {
    checks map[string]HealthCheck
    mutex  sync.RWMutex
}

type HealthCheck interface {
    Name() string
    Check(ctx context.Context) HealthResult
    Interval() time.Duration
    Timeout() time.Duration
}

type HealthResult struct {
    Status    HealthStatus          `json:"status"`
    Message   string                `json:"message"`
    Details   map[string]interface{} `json:"details"`
    Timestamp time.Time             `json:"timestamp"`
    Duration  time.Duration         `json:"duration"`
}

type HealthStatus string

const (
    HealthStatusHealthy   HealthStatus = "healthy"
    HealthStatusUnhealthy HealthStatus = "unhealthy"
    HealthStatusDegraded  HealthStatus = "degraded"
    HealthStatusUnknown   HealthStatus = "unknown"
)

// Built-in Health Checks
type DatabaseHealthCheck struct {
    db *sql.DB
}

type FileSystemHealthCheck struct {
    paths []string
}

type MemoryHealthCheck struct {
    threshold float64
}

type DiskSpaceHealthCheck struct {
    path      string
    threshold float64
}
```

### Alerting System

#### Alert Configuration
```go
type AlertConfig struct {
    Enabled         bool                    `json:"enabled"`
    Rules           []AlertRule             `json:"rules"`
    Channels        []AlertChannel          `json:"channels"`
    Throttling      AlertThrottling         `json:"throttling"`
    Maintenance     MaintenanceWindow       `json:"maintenance"`
}

type AlertRule struct {
    ID          string        `json:"id"`
    Name        string        `json:"name"`
    Description string        `json:"description"`
    Condition   string        `json:"condition"`
    Threshold   float64       `json:"threshold"`
    Duration    time.Duration `json:"duration"`
    Severity    AlertSeverity `json:"severity"`
    Enabled     bool          `json:"enabled"`
    Labels      map[string]string `json:"labels"`
}

type AlertSeverity string

const (
    AlertSeverityInfo     AlertSeverity = "info"
    AlertSeverityWarning  AlertSeverity = "warning"
    AlertSeverityError    AlertSeverity = "error"
    AlertSeverityCritical AlertSeverity = "critical"
)

type AlertChannel struct {
    Type    string                 `json:"type"`
    Config  map[string]interface{} `json:"config"`
    Enabled bool                   `json:"enabled"`
}

type AlertThrottling struct {
    Enabled     bool          `json:"enabled"`
    Window      time.Duration `json:"window"`
    MaxAlerts   int           `json:"maxAlerts"`
    Cooldown    time.Duration `json:"cooldown"`
}

type MaintenanceWindow struct {
    Enabled   bool      `json:"enabled"`
    StartTime time.Time `json:"startTime"`
    EndTime   time.Time `json:"endTime"`
    Recurring bool      `json:"recurring"`
}
```

## Data Retention and Storage

### Log Storage
```go
type LogStorage struct {
    MaxSizeMB      int           `json:"maxSizeMB"`
    MaxAge         time.Duration `json:"maxAge"`
    MaxFiles       int           `json:"maxFiles"`
    CompressionEnabled bool      `json:"compressionEnabled"`
    EncryptionEnabled  bool      `json:"encryptionEnabled"`
    RetentionPolicy    RetentionPolicy `json:"retentionPolicy"`
}

type RetentionPolicy struct {
    TraceLevel time.Duration `json:"traceLevel"`
    DebugLevel time.Duration `json:"debugLevel"`
    InfoLevel  time.Duration `json:"infoLevel"`
    WarnLevel  time.Duration `json:"warnLevel"`
    ErrorLevel time.Duration `json:"errorLevel"`
}
```

### Metrics Storage
```go
type MetricsStorage struct {
    Interval        time.Duration `json:"interval"`
    RetentionPeriod time.Duration `json:"retentionPeriod"`
    Aggregation     bool          `json:"aggregation"`
    Compression     bool          `json:"compression"`
    MaxPointsPerSeries int        `json:"maxPointsPerSeries"`
}
```

## User Interface Integration

### Status Components
```typescript
// React components for observability UI
interface StatusBarProps {
  status: ApplicationStatus;
  onStatusClick: () => void;
}

interface ProgressIndicatorProps {
  progress: GenerationProgress;
  showDetails: boolean;
}

interface ActivityLogProps {
  entries: ActivityLogEntry[];
  filter: LogFilter;
  onFilterChange: (filter: LogFilter) => void;
}

interface MetricsDashboardProps {
  metrics: PerformanceMetrics;
  timeRange: TimeRange;
  refreshInterval: number;
}
```

### Real-time Updates
```typescript
// WebSocket-like updates using Wails events
const useObservability = () => {
  const [status, setStatus] = useState<ApplicationStatus>();
  const [metrics, setMetrics] = useState<PerformanceMetrics>();
  const [logs, setLogs] = useState<ActivityLogEntry[]>([]);

  useEffect(() => {
    EventsOn('observability:status', setStatus);
    EventsOn('observability:metrics', setMetrics);
    EventsOn('observability:log', (entry: ActivityLogEntry) => {
      setLogs(prev => [entry, ...prev].slice(0, 1000));
    });

    return () => {
      EventsOff('observability:status');
      EventsOff('observability:metrics');
      EventsOff('observability:log');
    };
  }, []);

  return { status, metrics, logs };
};
```

## Configuration and Customization

### Observability Settings
```go
type ObservabilityConfig struct {
    Enabled            bool              `json:"enabled"`
    LogLevel           LogLevel          `json:"logLevel"`
    MetricsInterval    time.Duration     `json:"metricsInterval"`
    HealthCheckInterval time.Duration    `json:"healthCheckInterval"`
    RetentionPeriod    time.Duration     `json:"retentionPeriod"`
    AlertsEnabled      bool              `json:"alertsEnabled"`
    UIRefreshInterval  time.Duration     `json:"uiRefreshInterval"`
    ExportEnabled      bool              `json:"exportEnabled"`
    PrivacyMode        bool              `json:"privacyMode"`
    CustomMetrics      []CustomMetric    `json:"customMetrics"`
}

type CustomMetric struct {
    Name        string            `json:"name"`
    Type        string            `json:"type"`
    Description string            `json:"description"`
    Query       string            `json:"query"`
    Labels      map[string]string `json:"labels"`
    Enabled     bool              `json:"enabled"`
}
```

This observability specification provides comprehensive monitoring capabilities while maintaining the lightweight and user-focused nature of MCPWeaver. The system prioritizes user visibility and actionable insights while keeping operational complexity minimal.