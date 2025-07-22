package app

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"MCPWeaver/internal/database"
	"MCPWeaver/internal/generator"
	"MCPWeaver/internal/mapping"
	"MCPWeaver/internal/parser"
	"MCPWeaver/internal/validator"
	
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// App struct holds the application context and services
type App struct {
	ctx                 context.Context
	db                  *sql.DB
	projectRepo         *database.ProjectRepository
	validationCacheRepo *database.ValidationCacheRepository
	parserService       *parser.Service
	mappingService      *mapping.Service
	generatorService    *generator.Service
	validatorService    *validator.Service
	updateService       *UpdateService
	notificationService *NotificationService
	settings            *AppSettings
	errorManager        *ErrorManager
	performanceMonitor  *PerformanceMonitor
	activityLogService  *ActivityLogService
}

// NewApp creates a new application instance
func NewApp() *App {
	app := &App{
		parserService:      parser.NewService(),
		validatorService:   validator.New(),
		settings:           getDefaultSettings(),
		errorManager:       NewErrorManager(),
		performanceMonitor: NewPerformanceMonitor(),
	}
	
	// Initialize activity log service with default config
	logConfig := LogConfig{
		Level:         LogLevelInfo,
		BufferSize:    1000,
		RetentionDays: 7,
		EnableConsole: true,
		EnableBuffer:  true,
		FlushInterval: 5 * time.Minute,
	}
	app.activityLogService = NewActivityLogService(app, logConfig)
	
	return app
}

// OnStartup is called when the app starts, before the frontend is loaded
func (a *App) OnStartup(ctx context.Context) error {
	a.ctx = ctx
	
	// Record startup time
	startupStart := time.Now()
	
	// Log startup begin
	a.LogActivity(LogLevelInfo, "System", "Startup", "Application startup initiated",
		WithUserAction(true))
	
	// Initialize database
	dbPath := "./mcpweaver.db"
	dbWrapper, err := database.Open(dbPath)
	if err != nil {
		return fmt.Errorf("failed to initialize database: %w", err)
	}
	a.db = dbWrapper.GetConn()
	
	// Initialize repositories
	a.projectRepo = database.NewProjectRepository(dbWrapper)
	a.validationCacheRepo = database.NewValidationCacheRepository(dbWrapper)
	
	// Initialize update service
	a.updateService = NewUpdateService(ctx)
	
	// Initialize notification service
	a.notificationService = NewNotificationService(ctx, a.db)
	
	// Load settings
	settings, err := a.loadSettings()
	if err != nil {
		// Use default settings if loading fails
		runtime.LogWarning(a.ctx, "Failed to load settings, using defaults: "+err.Error())
		settings = getDefaultSettings()
	}
	a.settings = settings
	
	// Start update service
	if err := a.updateService.Start(); err != nil {
		runtime.LogWarning(a.ctx, "Failed to start update service: "+err.Error())
	}
	
	// Start notification service
	if err := a.notificationService.Start(); err != nil {
		runtime.LogWarning(a.ctx, "Failed to start notification service: "+err.Error())
	}
	
	// Record startup time
	startupDuration := time.Since(startupStart)
	a.performanceMonitor.RecordStartupTime(startupDuration)
	
	// Log startup completion
	a.LogActivity(LogLevelInfo, "System", "Startup", "Application startup completed",
		WithDuration(startupDuration),
		WithMetadata(map[string]interface{}{
			"startupTime": startupDuration.String(),
			"version":     "1.0.0",
		}))
	
	// Emit startup event
	runtime.EventsEmit(a.ctx, "system:startup", map[string]interface{}{
		"timestamp": time.Now(),
		"version":   "1.0.0",
		"startup_time": startupDuration.String(),
	})
	
	return nil
}

// OnShutdown is called when the app is about to quit
func (a *App) OnShutdown(ctx context.Context) error {
	// Log shutdown initiation
	a.LogActivity(LogLevelInfo, "System", "Shutdown", "Application shutdown initiated",
		WithUserAction(true))
		
	// Emit shutdown event
	runtime.EventsEmit(a.ctx, "system:shutdown", map[string]interface{}{
		"timestamp": time.Now(),
	})
	
	// Stop update service
	if a.updateService != nil {
		if err := a.updateService.Stop(); err != nil {
			runtime.LogError(a.ctx, "Failed to stop update service: "+err.Error())
		}
	}
	
	// Stop notification service
	if a.notificationService != nil {
		if err := a.notificationService.Stop(); err != nil {
			runtime.LogError(a.ctx, "Failed to stop notification service: "+err.Error())
		}
	}
	
	// Save settings
	if err := a.saveSettingsToFile(); err != nil {
		runtime.LogError(a.ctx, "Failed to save settings: "+err.Error())
	}
	
	// Close database connection
	if a.db != nil {
		if err := a.db.Close(); err != nil {
			runtime.LogError(a.ctx, "Failed to close database: "+err.Error())
		}
	}
	
	// Close activity log service
	if a.activityLogService != nil {
		a.LogActivity(LogLevelInfo, "System", "Shutdown", "Application shutdown completed")
		a.activityLogService.Close()
	}
	
	return nil
}

// GetPerformanceMetrics returns current performance metrics
func (a *App) GetPerformanceMetrics() *PerformanceMetrics {
	return a.performanceMonitor.GetMetrics()
}

// GetMemoryUsage returns current memory usage in MB
func (a *App) GetMemoryUsage() float64 {
	return a.performanceMonitor.GetMemoryUsageMB()
}

// ForceGarbageCollection forces garbage collection
func (a *App) ForceGarbageCollection() {
	a.performanceMonitor.ForceGC()
}

// OnDomReady is called after the frontend dom is ready
func (a *App) OnDomReady(ctx context.Context) {
	// Emit DOM ready event
	runtime.EventsEmit(a.ctx, "system:ready", map[string]interface{}{
		"timestamp": time.Now(),
	})
}

// OnBeforeClose is called when the application is about to quit
func (a *App) OnBeforeClose(ctx context.Context) bool {
	// Allow the app to close
	return false
}

// createAPIError creates a standardized API error
func (a *App) createAPIError(errorType, code, message string, details map[string]string) *APIError {
	return &APIError{
		Type:      errorType,
		Code:      code,
		Message:   message,
		Details:   details,
		Timestamp: time.Now(),
	}
}

// emitError emits an error event to the frontend
func (a *App) emitError(err *APIError) {
	if a.ctx != nil {
		defer func() {
			if r := recover(); r != nil {
				// Silently handle context-related panics during testing
			}
		}()
		runtime.EventsEmit(a.ctx, "system:error", err)
	}
}

// ReportError reports an error from the frontend
func (a *App) ReportError(errorReport map[string]interface{}) error {
	// Log the error
	fmt.Printf("Frontend Error Report: %+v\n", errorReport)
	
	// Here you could send the error to a logging service, database, or monitoring system
	// For now, we'll just emit it as a system event
	if a.ctx != nil {
		defer func() {
			if r := recover(); r != nil {
				// Silently handle context-related panics during testing
			}
		}()
		runtime.EventsEmit(a.ctx, "system:error-report", errorReport)
	}
	
	return nil
}

// Activity Log API Methods

// GetActivityLogs retrieves activity logs based on filter criteria
func (a *App) GetActivityLogs(ctx context.Context, filter LogFilter) ([]ActivityLogEntry, error) {
	if a.activityLogService == nil {
		return nil, a.createAPIError(ErrorTypeSystem, ErrCodeInternalError, "Activity log service not initialized", nil)
	}
	
	logs := a.activityLogService.GetLogs(filter)
	return logs, nil
}

// SearchActivityLogs performs a search across activity log entries
func (a *App) SearchActivityLogs(ctx context.Context, request LogSearchRequest) (*LogSearchResult, error) {
	if a.activityLogService == nil {
		return nil, a.createAPIError(ErrorTypeSystem, ErrCodeInternalError, "Activity log service not initialized", nil)
	}
	
	return a.activityLogService.SearchLogs(ctx, request)
}

// ExportActivityLogs exports activity logs to a file
func (a *App) ExportActivityLogs(ctx context.Context, request LogExportRequest) (*LogExportResult, error) {
	if a.activityLogService == nil {
		return nil, a.createAPIError(ErrorTypeSystem, ErrCodeInternalError, "Activity log service not initialized", nil)
	}
	
	return a.activityLogService.ExportLogs(ctx, request)
}

// GetApplicationStatus returns the current application status
func (a *App) GetApplicationStatus(ctx context.Context) (*ApplicationStatus, error) {
	if a.activityLogService == nil {
		return nil, a.createAPIError(ErrorTypeSystem, ErrCodeInternalError, "Activity log service not initialized", nil)
	}
	
	return a.activityLogService.GetApplicationStatus(), nil
}

// UpdateLogConfig updates the activity log configuration
func (a *App) UpdateLogConfig(ctx context.Context, config LogConfig) error {
	if a.activityLogService == nil {
		return a.createAPIError(ErrorTypeSystem, ErrCodeInternalError, "Activity log service not initialized", nil)
	}
	
	return a.activityLogService.UpdateLogConfig(config)
}

// ClearActivityLogs clears activity logs based on age
func (a *App) ClearActivityLogs(ctx context.Context, olderThanHours int) (int, error) {
	if a.activityLogService == nil {
		return 0, a.createAPIError(ErrorTypeSystem, ErrCodeInternalError, "Activity log service not initialized", nil)
	}
	
	var olderThan time.Duration
	if olderThanHours > 0 {
		olderThan = time.Duration(olderThanHours) * time.Hour
	}
	
	cleared := a.activityLogService.ClearLogs(olderThan)
	return cleared, nil
}

// LogActivity logs an activity entry (for internal use)
func (a *App) LogActivity(level LogLevel, component, operation, message string, options ...LogEntryOption) {
	if a.activityLogService != nil {
		a.activityLogService.LogEntry(level, component, operation, message, options...)
	}
}

// CreateErrorReport creates an error report and logs it
func (a *App) CreateErrorReport(ctx context.Context, errorType ErrorType, severity ErrorSeverity, component, operation, message string, err error) (*ErrorReport, error) {
	if a.activityLogService == nil {
		return nil, a.createAPIError(ErrorTypeSystem, ErrCodeInternalError, "Activity log service not initialized", nil)
	}
	
	report := a.activityLogService.ReportError(errorType, severity, component, operation, message, err)
	return report, nil
}

// GetErrorReports retrieves error reports
func (a *App) GetErrorReports(ctx context.Context, includeResolved bool) ([]ErrorReport, error) {
	if a.activityLogService == nil {
		return nil, a.createAPIError(ErrorTypeSystem, ErrCodeInternalError, "Activity log service not initialized", nil)
	}
	
	reports := a.activityLogService.GetErrorReports(includeResolved)
	return reports, nil
}

// emitNotification emits a notification event to the frontend
func (a *App) emitNotification(notificationType, title, message string) {
	if a.ctx != nil {
		defer func() {
			if r := recover(); r != nil {
				// Silently handle context-related panics during testing
			}
		}()
		runtime.EventsEmit(a.ctx, "system:notification", map[string]interface{}{
			"type":    notificationType,
			"title":   title,
			"message": message,
			"timestamp": time.Now(),
		})
	}
}

// loadSettings loads application settings from storage
func (a *App) loadSettings() (*AppSettings, error) {
	// Use the new migration-aware settings loader
	return a.MigrateAndLoadSettings()
}

// getDefaultSettings returns default application settings
func getDefaultSettings() *AppSettings {
	return &AppSettings{
		Theme:             "light",
		Language:          "en",
		AutoSave:          true,
		DefaultOutputPath: "./output",
		RecentProjects:    []string{},
		WindowSettings: WindowSettings{
			Width:     1200,
			Height:    800,
			Maximized: false,
			X:         100,
			Y:         100,
		},
		EditorSettings: EditorSettings{
			FontSize:        14,
			FontFamily:      "Monaco",
			TabSize:         4,
			WordWrap:        true,
			LineNumbers:     true,
			SyntaxHighlight: true,
		},
		GenerationSettings: GenerationSettings{
			DefaultTemplate:     "default",
			EnableValidation:    true,
			AutoOpenOutput:      true,
			ShowAdvancedOptions: false,
			BackupOnGenerate:    true,
			CustomTemplates:     []string{},
			PerformanceMode:     false,
			MaxWorkers:          4,
		},
		NotificationSettings: NotificationSettings{
			EnableDesktopNotifications: true,
			EnableSoundNotifications:   false,
			NotificationPosition:       "top-right",
			NotificationDuration:       5000,
			SoundVolume:                0.5,
			ShowGenerationProgress:     true,
			ShowErrorNotifications:     true,
			ShowSuccessNotifications:   true,
		},
		AppearanceSettings: AppearanceSettings{
			UITheme:         "system",
			AccentColor:     "#007acc",
			WindowOpacity:   1.0,
			ShowAnimation:   true,
			ReducedMotion:   false,
			FontScale:       1.0,
			CompactMode:     false,
			ShowSidebar:     true,
			SidebarPosition: "left",
			ShowStatusBar:   true,
			ShowToolbar:     true,
		},
		UpdateSettings: DefaultUpdateSettings(),
	}
}