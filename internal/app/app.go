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
}

// NewApp creates a new application instance
func NewApp() *App {
	return &App{
		parserService:      parser.NewService(),
		validatorService:   validator.New(),
		settings:           getDefaultSettings(),
		errorManager:       NewErrorManager(),
		performanceMonitor: NewPerformanceMonitor(),
	}
}

// OnStartup is called when the app starts, before the frontend is loaded
func (a *App) OnStartup(ctx context.Context) error {
	a.ctx = ctx
	
	// Record startup time
	startupStart := time.Now()
	
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
	if err := a.saveSettings(); err != nil {
		runtime.LogError(a.ctx, "Failed to save settings: "+err.Error())
	}
	
	// Close database connection
	if a.db != nil {
		if err := a.db.Close(); err != nil {
			runtime.LogError(a.ctx, "Failed to close database: "+err.Error())
		}
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
	// TODO: Implement settings loading from file/database
	return getDefaultSettings(), nil
}

// saveSettings saves application settings to storage
func (a *App) saveSettings() error {
	// TODO: Implement settings saving to file/database
	return nil
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
		},
		UpdateSettings: DefaultUpdateSettings(),
	}
}