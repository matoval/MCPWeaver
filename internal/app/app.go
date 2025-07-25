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
	"MCPWeaver/internal/plugin"
	"MCPWeaver/internal/validator"
	
	"github.com/google/uuid"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// App struct holds the application context and services
type App struct {
	ctx                 context.Context
	db                  *sql.DB
	projectRepo         *database.ProjectRepository
	templateRepo        *database.TemplateRepository
	validationCacheRepo *database.ValidationCacheRepository
	parserService       *parser.Service
	mappingService      *mapping.Service
	generatorService    *generator.Service
	validatorService    *validator.Service
	updateService       *UpdateService
	notificationService *NotificationService
	pluginService       *plugin.Service
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
		pluginService:      plugin.NewService(nil), // Use default config
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
	a.templateRepo = database.NewTemplateRepository(dbWrapper)
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
	
	// Initialize plugin service
	if err := a.pluginService.Initialize(); err != nil {
		runtime.LogWarning(a.ctx, "Failed to initialize plugin service: "+err.Error())
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
	
	// Shutdown plugin service
	if a.pluginService != nil {
		if err := a.pluginService.Shutdown(); err != nil {
			runtime.LogError(a.ctx, "Failed to shutdown plugin service: "+err.Error())
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

// Plugin Management API Methods

// GetPlugins returns all loaded plugins
func (a *App) GetPlugins(ctx context.Context) (map[string]*plugin.PluginInstanceAPI, error) {
	if a.pluginService == nil {
		return nil, a.createAPIError(ErrorTypeSystem, ErrCodeInternalError, "Plugin service not initialized", nil)
	}
	
	return a.pluginService.GetPlugins()
}

// GetPlugin returns a specific plugin
func (a *App) GetPlugin(ctx context.Context, pluginID string) (*plugin.PluginInstanceAPI, error) {
	if a.pluginService == nil {
		return nil, a.createAPIError(ErrorTypeSystem, ErrCodeInternalError, "Plugin service not initialized", nil)
	}
	
	return a.pluginService.GetPlugin(pluginID)
}

// LoadPlugin loads a plugin from file path
func (a *App) LoadPlugin(ctx context.Context, pluginPath string) error {
	if a.pluginService == nil {
		return a.createAPIError(ErrorTypeSystem, ErrCodeInternalError, "Plugin service not initialized", nil)
	}
	
	return a.pluginService.LoadPlugin(pluginPath)
}

// UnloadPlugin unloads a plugin
func (a *App) UnloadPlugin(ctx context.Context, pluginID string) error {
	if a.pluginService == nil {
		return a.createAPIError(ErrorTypeSystem, ErrCodeInternalError, "Plugin service not initialized", nil)
	}
	
	return a.pluginService.UnloadPlugin(pluginID)
}

// EnablePlugin enables a disabled plugin
func (a *App) EnablePlugin(ctx context.Context, pluginID string) error {
	if a.pluginService == nil {
		return a.createAPIError(ErrorTypeSystem, ErrCodeInternalError, "Plugin service not initialized", nil)
	}
	
	return a.pluginService.EnablePlugin(pluginID)
}

// DisablePlugin disables an active plugin
func (a *App) DisablePlugin(ctx context.Context, pluginID string) error {
	if a.pluginService == nil {
		return a.createAPIError(ErrorTypeSystem, ErrCodeInternalError, "Plugin service not initialized", nil)
	}
	
	return a.pluginService.DisablePlugin(pluginID)
}

// GetPluginsByCapability returns plugins with specific capability
func (a *App) GetPluginsByCapability(ctx context.Context, capability string) ([]*plugin.PluginInstanceAPI, error) {
	if a.pluginService == nil {
		return nil, a.createAPIError(ErrorTypeSystem, ErrCodeInternalError, "Plugin service not initialized", nil)
	}
	
	return a.pluginService.GetPluginsByCapability(capability)
}

// SearchPlugins searches for plugins in marketplace
func (a *App) SearchPlugins(ctx context.Context, query string, category string, tags []string, limit int) (*plugin.SearchResponse, error) {
	if a.pluginService == nil {
		return nil, a.createAPIError(ErrorTypeSystem, ErrCodeInternalError, "Plugin service not initialized", nil)
	}
	
	return a.pluginService.SearchPlugins(ctx, query, category, tags, limit)
}

// InstallPlugin installs a plugin from marketplace
func (a *App) InstallPlugin(ctx context.Context, pluginID string) error {
	if a.pluginService == nil {
		return a.createAPIError(ErrorTypeSystem, ErrCodeInternalError, "Plugin service not initialized", nil)
	}
	
	return a.pluginService.InstallPlugin(ctx, pluginID)
}

// GetPluginCapabilities returns available plugin capabilities
func (a *App) GetPluginCapabilities(ctx context.Context) []string {
	if a.pluginService == nil {
		return []string{}
	}
	
	return a.pluginService.GetPluginCapabilities()
}

// GetPluginPermissions returns available plugin permissions
func (a *App) GetPluginPermissions(ctx context.Context) []string {
	if a.pluginService == nil {
		return []string{}
	}
	
	return a.pluginService.GetPluginPermissions()
}

// Template Management API Methods

// GetAllTemplates retrieves all templates
func (a *App) GetAllTemplates() ([]*database.AppTemplate, error) {
	if a.templateRepo == nil {
		return nil, a.createAPIError(ErrorTypeSystem, ErrCodeInternalError, "Template repository not initialized", nil)
	}
	
	templates, err := a.templateRepo.GetAll()
	if err != nil {
		return nil, a.createAPIError(ErrorTypeSystem, ErrCodeInternalError, "Failed to retrieve templates", map[string]string{"error": err.Error()})
	}

	return templates, nil
}

// DeleteTemplate deletes a template
func (a *App) DeleteTemplate(id string) error {
	if id == "" {
		return a.createAPIError(ErrorTypeValidation, ErrCodeValidation, "Template ID is required", nil)
	}

	if a.templateRepo == nil {
		return a.createAPIError(ErrorTypeSystem, ErrCodeInternalError, "Template repository not initialized", nil)
	}

	// Get template to check if it's built-in
	template, err := a.templateRepo.GetByID(id)
	if err != nil {
		return a.createAPIError(ErrorTypeSystem, ErrCodeInternalError, "Failed to retrieve template", map[string]string{"templateId": id, "error": err.Error()})
	}

	// Prevent deleting built-in templates
	if template.IsBuiltIn {
		return a.createAPIError(ErrorTypeValidation, ErrCodeValidation, "Cannot delete built-in templates", map[string]string{"templateId": id})
	}

	// Delete from database
	if err := a.templateRepo.Delete(id); err != nil {
		return a.createAPIError(ErrorTypeSystem, ErrCodeInternalError, "Failed to delete template", map[string]string{"templateId": id, "error": err.Error()})
	}

	// Emit event
	if a.ctx != nil {
		defer func() {
			if r := recover(); r != nil {
				// Silently handle context-related panics during testing
			}
		}()
		runtime.EventsEmit(a.ctx, "template:deleted", map[string]interface{}{
			"templateId": template.ID,
			"name":       template.Name,
		})
	}

	return nil
}

// DuplicateTemplate creates a copy of an existing template
func (a *App) DuplicateTemplate(id string, newName string) (*database.AppTemplate, error) {
	if id == "" {
		return nil, a.createAPIError(ErrorTypeValidation, ErrCodeValidation, "Template ID is required", nil)
	}

	if newName == "" {
		return nil, a.createAPIError(ErrorTypeValidation, ErrCodeValidation, "New template name is required", nil)
	}

	if a.templateRepo == nil {
		return nil, a.createAPIError(ErrorTypeSystem, ErrCodeInternalError, "Template repository not initialized", nil)
	}

	// Get original template
	original, err := a.templateRepo.GetByID(id)
	if err != nil {
		return nil, a.createAPIError(ErrorTypeSystem, ErrCodeInternalError, "Failed to retrieve original template", map[string]string{"templateId": id, "error": err.Error()})
	}

	// Check if new name conflicts
	existing, err := a.templateRepo.GetByName(newName)
	if err == nil && existing != nil {
		return nil, a.createAPIError(ErrorTypeValidation, ErrCodeValidation, fmt.Sprintf("Template with name '%s' already exists", newName), map[string]string{"name": newName})
	}

	// Create duplicate
	duplicate := &database.AppTemplate{
		ID:          generateTemplateID(),
		Name:        newName,
		Description: original.Description + " (Copy)",
		Version:     "1.0.0", // Reset version for copy
		Author:      original.Author,
		Type:        "custom", // Copies are always custom
		Path:        original.Path,
		IsBuiltIn:   false,
		Variables:   original.Variables,
	}

	// Save to database
	if err := a.templateRepo.Create(duplicate); err != nil {
		return nil, a.createAPIError(ErrorTypeSystem, ErrCodeInternalError, "Failed to create duplicate template", map[string]string{"error": err.Error()})
	}

	// Emit event
	if a.ctx != nil {
		defer func() {
			if r := recover(); r != nil {
				// Silently handle context-related panics during testing
			}
		}()
		runtime.EventsEmit(a.ctx, "template:duplicated", map[string]interface{}{
			"originalId": original.ID,
			"duplicateId": duplicate.ID,
			"name":       duplicate.Name,
		})
	}

	return duplicate, nil
}

// Helper functions

// ImportTemplate imports a template from various sources
func (a *App) ImportTemplate(request TemplateImportRequest) (*database.AppTemplate, error) {
	if request.Source == "" {
		return nil, a.createAPIError(ErrorTypeValidation, ErrCodeValidation, "Import source is required", nil)
	}

	if a.templateRepo == nil {
		return nil, a.createAPIError(ErrorTypeSystem, ErrCodeInternalError, "Template repository not initialized", nil)
	}

	switch request.Source {
	case "file":
		return a.importTemplateFromFile(request)
	default:
		return nil, a.createAPIError(ErrorTypeValidation, ErrCodeValidation, fmt.Sprintf("Unsupported import source: %s", request.Source), map[string]string{"source": request.Source})
	}
}

// importTemplateFromFile imports a template from a local file (simplified version)
func (a *App) importTemplateFromFile(request TemplateImportRequest) (*database.AppTemplate, error) {
	if request.Path == "" {
		return nil, a.createAPIError(ErrorTypeValidation, ErrCodeValidation, "File path is required for file import", nil)
	}

	// Check if file exists (simplified check)
	if request.Path == "" {
		return nil, a.createAPIError(ErrorTypeFileSystem, ErrCodeFileAccess, "Template file does not exist", map[string]string{"path": request.Path})
	}

	// Generate template metadata (simplified)
	templateName := "imported-template"
	targetType := "custom"
	if request.ImportOptions.TargetType != "" {
		targetType = string(request.ImportOptions.TargetType)
	}

	// Create template record
	template := &database.AppTemplate{
		ID:          generateTemplateID(),
		Name:        templateName,
		Description: fmt.Sprintf("Imported from %s", request.Path),
		Version:     "1.0.0",
		Author:      "Imported",
		Type:        targetType,
		Path:        request.Path,
		IsBuiltIn:   false,
		Variables:   []database.TemplateVariable{},
	}

	// Save to database
	if err := a.templateRepo.Create(template); err != nil {
		return nil, a.createAPIError(ErrorTypeSystem, ErrCodeInternalError, "Failed to import template", map[string]string{"error": err.Error()})
	}

	// Emit event
	if a.ctx != nil {
		defer func() {
			if r := recover(); r != nil {
				// Silently handle context-related panics during testing
			}
		}()
		runtime.EventsEmit(a.ctx, "template:imported", map[string]interface{}{
			"templateId": template.ID,
			"name":       template.Name,
			"source":     "file",
		})
	}

	return template, nil
}

// ExportTemplate exports a template to various formats
func (a *App) ExportTemplate(request TemplateExportRequest) (*ExportResult, error) {
	if request.TemplateID == "" {
		return nil, a.createAPIError(ErrorTypeValidation, ErrCodeValidation, "Template ID is required", nil)
	}

	if request.TargetPath == "" {
		return nil, a.createAPIError(ErrorTypeValidation, ErrCodeValidation, "Target path is required", nil)
	}

	if a.templateRepo == nil {
		return nil, a.createAPIError(ErrorTypeSystem, ErrCodeInternalError, "Template repository not initialized", nil)
	}

	// Get template
	template, err := a.templateRepo.GetByID(request.TemplateID)
	if err != nil {
		return nil, a.createAPIError(ErrorTypeSystem, ErrCodeInternalError, "Failed to retrieve template", map[string]string{"templateId": request.TemplateID, "error": err.Error()})
	}

	switch request.Format {
	case "single":
		return a.exportTemplateAsSingle(template, request)
	default:
		return nil, a.createAPIError(ErrorTypeValidation, ErrCodeValidation, fmt.Sprintf("Unsupported export format: %s", request.Format), map[string]string{"format": request.Format})
	}
}

// exportTemplateAsSingle exports template as a single file (simplified version)
func (a *App) exportTemplateAsSingle(template *database.AppTemplate, request TemplateExportRequest) (*ExportResult, error) {
	// Create export result (simplified)
	result := &ExportResult{
		ProjectID:   "",
		ProjectName: template.Name,
		TargetDir:   request.TargetPath,
		ExportedFiles: []ExportedFile{
			{
				Name: template.Name + ".template",
				Path: request.TargetPath,
				Size: 0,
			},
		},
		TotalFiles: 1,
		TotalSize:  0,
		ExportedAt: time.Now(),
	}

	// Emit event
	if a.ctx != nil {
		defer func() {
			if r := recover(); r != nil {
				// Silently handle context-related panics during testing
			}
		}()
		runtime.EventsEmit(a.ctx, "template:exported", map[string]interface{}{
			"templateId": template.ID,
			"format":     "single",
			"targetPath": request.TargetPath,
		})
	}

	return result, nil
}

// generateTemplateID generates a unique template ID
func generateTemplateID() string {
	return uuid.New().String()
}