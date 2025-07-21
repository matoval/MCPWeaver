package app

import (
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// GetSettings returns the current application settings
func (a *App) GetSettings() (*AppSettings, error) {
	if a.settings == nil {
		return getDefaultSettings(), nil
	}
	return a.settings, nil
}

// UpdateSettings updates the application settings
func (a *App) UpdateSettings(settings AppSettings) error {
	// Validate settings
	if err := a.validateSettings(&settings); err != nil {
		return err
	}

	// Update current settings
	a.settings = &settings

	// Save settings to storage
	if err := a.saveSettingsToFile(); err != nil {
		return a.createAPIError("internal", ErrCodeInternalError, "Failed to save settings", map[string]string{
			"error": err.Error(),
		})
	}

	// Emit settings updated event
	runtime.EventsEmit(a.ctx, "settings:updated", settings)

	// Send notification
	a.emitNotification("success", "Settings Updated", "Application settings have been saved successfully")

	return nil
}

// ResetSettings resets all settings to their default values
func (a *App) ResetSettings() error {
	// Reset to default settings
	a.settings = getDefaultSettings()

	// Save default settings to storage
	if err := a.saveSettingsToFile(); err != nil {
		return a.createAPIError("internal", ErrCodeInternalError, "Failed to reset settings", map[string]string{
			"error": err.Error(),
		})
	}

	// Emit settings reset event
	runtime.EventsEmit(a.ctx, "settings:reset", a.settings)

	// Send notification
	a.emitNotification("success", "Settings Reset", "Application settings have been reset to defaults")

	return nil
}

// GetSettingsFilePath returns the path to the settings file
func (a *App) GetSettingsFilePath() string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "./mcpweaver-settings.json"
	}
	return filepath.Join(homeDir, ".mcpweaver", "settings.json")
}

// saveSettingsToFile saves settings to a JSON file
func (a *App) saveSettingsToFile() error {
	settingsPath := a.GetSettingsFilePath()
	
	// Ensure directory exists
	dir := filepath.Dir(settingsPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	// Marshal settings to JSON
	data, err := json.MarshalIndent(a.settings, "", "  ")
	if err != nil {
		return err
	}

	// Write to file
	return os.WriteFile(settingsPath, data, 0644)
}

// loadSettingsFromFile loads settings from a JSON file
func (a *App) loadSettingsFromFile() (*AppSettings, error) {
	settingsPath := a.GetSettingsFilePath()
	
	// Check if file exists
	if _, err := os.Stat(settingsPath); os.IsNotExist(err) {
		return getDefaultSettings(), nil
	}

	// Read file
	data, err := os.ReadFile(settingsPath)
	if err != nil {
		return nil, err
	}

	// Unmarshal JSON
	var settings AppSettings
	if err := json.Unmarshal(data, &settings); err != nil {
		return nil, err
	}

	return &settings, nil
}

// validateSettings validates the settings structure
func (a *App) validateSettings(settings *AppSettings) error {
	// Validate theme
	if settings.Theme != "light" && settings.Theme != "dark" && settings.Theme != "auto" {
		return a.createAPIError("validation", ErrCodeValidation, "Invalid theme value", map[string]string{
			"theme": settings.Theme,
			"valid": "light, dark, auto",
		})
	}

	// Validate language
	if settings.Language == "" {
		settings.Language = "en"
	}

	// Validate window settings
	if settings.WindowSettings.Width < 800 {
		settings.WindowSettings.Width = 800
	}
	if settings.WindowSettings.Height < 600 {
		settings.WindowSettings.Height = 600
	}

	// Validate editor settings
	if settings.EditorSettings.FontSize < 8 || settings.EditorSettings.FontSize > 24 {
		settings.EditorSettings.FontSize = 14
	}
	if settings.EditorSettings.TabSize < 2 || settings.EditorSettings.TabSize > 8 {
		settings.EditorSettings.TabSize = 4
	}
	if settings.EditorSettings.FontFamily == "" {
		settings.EditorSettings.FontFamily = "Monaco"
	}

	// Validate generation settings
	if settings.GenerationSettings.DefaultTemplate == "" {
		settings.GenerationSettings.DefaultTemplate = "default"
	}

	// Validate default output path
	if settings.DefaultOutputPath == "" {
		settings.DefaultOutputPath = "./output"
	}

	// Limit recent projects to 20
	if len(settings.RecentProjects) > 20 {
		settings.RecentProjects = settings.RecentProjects[:20]
	}

	// Validate notification settings
	if settings.NotificationSettings.NotificationPosition == "" {
		settings.NotificationSettings.NotificationPosition = "top-right"
	}
	if settings.NotificationSettings.NotificationDuration < 1000 {
		settings.NotificationSettings.NotificationDuration = 5000
	}
	if settings.NotificationSettings.SoundVolume < 0 || settings.NotificationSettings.SoundVolume > 1 {
		settings.NotificationSettings.SoundVolume = 0.5
	}

	// Validate appearance settings
	if settings.AppearanceSettings.UITheme == "" {
		settings.AppearanceSettings.UITheme = "system"
	}
	if settings.AppearanceSettings.AccentColor == "" {
		settings.AppearanceSettings.AccentColor = "#007acc"
	}
	if settings.AppearanceSettings.WindowOpacity < 0.1 || settings.AppearanceSettings.WindowOpacity > 1.0 {
		settings.AppearanceSettings.WindowOpacity = 1.0
	}
	if settings.AppearanceSettings.FontScale < 0.5 || settings.AppearanceSettings.FontScale > 2.0 {
		settings.AppearanceSettings.FontScale = 1.0
	}
	if settings.AppearanceSettings.SidebarPosition == "" {
		settings.AppearanceSettings.SidebarPosition = "left"
	}

	// Validate generation settings performance options
	if settings.GenerationSettings.MaxWorkers < 1 || settings.GenerationSettings.MaxWorkers > 16 {
		settings.GenerationSettings.MaxWorkers = 4
	}

	return nil
}

// UpdateTheme updates just the theme setting
func (a *App) UpdateTheme(theme string) error {
	if theme != "light" && theme != "dark" && theme != "auto" {
		return a.createAPIError("validation", ErrCodeValidation, "Invalid theme value", map[string]string{
			"theme": theme,
			"valid": "light, dark, auto",
		})
	}

	a.settings.Theme = theme

	// Save settings
	if err := a.saveSettingsToFile(); err != nil {
		return a.createAPIError("internal", ErrCodeInternalError, "Failed to save theme setting", map[string]string{
			"error": err.Error(),
		})
	}

	// Emit theme change event
	runtime.EventsEmit(a.ctx, "theme:changed", theme)

	return nil
}

// UpdateLanguage updates just the language setting
func (a *App) UpdateLanguage(language string) error {
	if language == "" {
		return a.createAPIError("validation", ErrCodeValidation, "Language is required", nil)
	}

	a.settings.Language = language

	// Save settings
	if err := a.saveSettingsToFile(); err != nil {
		return a.createAPIError("internal", ErrCodeInternalError, "Failed to save language setting", map[string]string{
			"error": err.Error(),
		})
	}

	// Emit language change event
	runtime.EventsEmit(a.ctx, "language:changed", language)

	return nil
}

// UpdateWindowSettings updates the window settings
func (a *App) UpdateWindowSettings(windowSettings WindowSettings) error {
	// Validate window settings
	if windowSettings.Width < 800 {
		windowSettings.Width = 800
	}
	if windowSettings.Height < 600 {
		windowSettings.Height = 600
	}

	a.settings.WindowSettings = windowSettings

	// Save settings
	if err := a.saveSettingsToFile(); err != nil {
		return a.createAPIError("internal", ErrCodeInternalError, "Failed to save window settings", map[string]string{
			"error": err.Error(),
		})
	}

	// Emit window settings change event
	runtime.EventsEmit(a.ctx, "window:settings:changed", windowSettings)

	return nil
}

// AddRecentProject adds a project to the recent projects list
func (a *App) AddRecentProject(projectID string) error {
	if projectID == "" {
		return a.createAPIError("validation", ErrCodeValidation, "Project ID is required", nil)
	}

	// Use the helper function from projects.go
	a.addToRecentProjects(projectID)

	// Save settings
	if err := a.saveSettingsToFile(); err != nil {
		return a.createAPIError("internal", ErrCodeInternalError, "Failed to save recent projects", map[string]string{
			"error": err.Error(),
		})
	}

	// Emit recent projects updated event
	runtime.EventsEmit(a.ctx, "recent:projects:updated", a.settings.RecentProjects)

	return nil
}

// ClearRecentProjects clears the recent projects list
func (a *App) ClearRecentProjects() error {
	a.settings.RecentProjects = []string{}

	// Save settings
	if err := a.saveSettingsToFile(); err != nil {
		return a.createAPIError("internal", ErrCodeInternalError, "Failed to clear recent projects", map[string]string{
			"error": err.Error(),
		})
	}

	// Emit recent projects cleared event
	runtime.EventsEmit(a.ctx, "recent:projects:cleared", nil)

	// Send notification
	a.emitNotification("success", "Recent Projects Cleared", "Recent projects list has been cleared")

	return nil
}

// ExportSettings exports settings to a file
func (a *App) ExportSettings() (string, error) {
	// Marshal settings to JSON
	data, err := json.MarshalIndent(a.settings, "", "  ")
	if err != nil {
		return "", a.createAPIError("internal", ErrCodeInternalError, "Failed to export settings", map[string]string{
			"error": err.Error(),
		})
	}

	// Open save dialog
	filePath, err := a.SaveFile(string(data), "mcpweaver-settings.json", []FileFilter{
		{
			DisplayName: "JSON Files",
			Pattern:     "*.json",
			Extensions:  []string{".json"},
		},
	})

	if err != nil {
		return "", err
	}

	if filePath != "" {
		a.emitNotification("success", "Settings Exported", "Settings have been exported successfully")
	}

	return filePath, nil
}

// ImportSettings imports settings from a file
func (a *App) ImportSettings() error {
	// Open file dialog
	filePath, err := a.SelectFile([]FileFilter{
		{
			DisplayName: "JSON Files",
			Pattern:     "*.json",
			Extensions:  []string{".json"},
		},
	})

	if err != nil {
		return err
	}

	if filePath == "" {
		return nil // User cancelled
	}

	// Read file
	data, err := os.ReadFile(filePath)
	if err != nil {
		return a.createAPIError("file_system", ErrCodeFileAccess, "Failed to read settings file", map[string]string{
			"path": filePath,
			"error": err.Error(),
		})
	}

	// Unmarshal JSON
	var settings AppSettings
	if err := json.Unmarshal(data, &settings); err != nil {
		return a.createAPIError("validation", ErrCodeValidation, "Invalid settings file format", map[string]string{
			"path": filePath,
			"error": err.Error(),
		})
	}

	// Validate imported settings
	if err := a.validateSettings(&settings); err != nil {
		return err
	}

	// Update current settings
	a.settings = &settings

	// Save settings
	if err := a.saveSettingsToFile(); err != nil {
		return a.createAPIError("internal", ErrCodeInternalError, "Failed to save imported settings", map[string]string{
			"error": err.Error(),
		})
	}

	// Emit settings imported event
	runtime.EventsEmit(a.ctx, "settings:imported", settings)

	// Send notification
	a.emitNotification("success", "Settings Imported", "Settings have been imported successfully")

	return nil
}

// UpdateNotificationSettings updates notification settings
func (a *App) UpdateNotificationSettings(notificationSettings NotificationSettings) error {
	// Validate notification settings
	if notificationSettings.NotificationPosition == "" {
		notificationSettings.NotificationPosition = "top-right"
	}
	if notificationSettings.NotificationDuration < 1000 {
		notificationSettings.NotificationDuration = 5000
	}
	if notificationSettings.SoundVolume < 0 || notificationSettings.SoundVolume > 1 {
		notificationSettings.SoundVolume = 0.5
	}

	a.settings.NotificationSettings = notificationSettings

	// Save settings
	if err := a.saveSettingsToFile(); err != nil {
		return a.createAPIError("internal", ErrCodeInternalError, "Failed to save notification settings", map[string]string{
			"error": err.Error(),
		})
	}

	// Emit notification settings change event
	runtime.EventsEmit(a.ctx, "notification:settings:changed", notificationSettings)

	return nil
}

// UpdateAppearanceSettings updates appearance settings
func (a *App) UpdateAppearanceSettings(appearanceSettings AppearanceSettings) error {
	// Validate appearance settings
	if appearanceSettings.UITheme == "" {
		appearanceSettings.UITheme = "system"
	}
	if appearanceSettings.AccentColor == "" {
		appearanceSettings.AccentColor = "#007acc"
	}
	if appearanceSettings.WindowOpacity < 0.1 || appearanceSettings.WindowOpacity > 1.0 {
		appearanceSettings.WindowOpacity = 1.0
	}
	if appearanceSettings.FontScale < 0.5 || appearanceSettings.FontScale > 2.0 {
		appearanceSettings.FontScale = 1.0
	}
	if appearanceSettings.SidebarPosition == "" {
		appearanceSettings.SidebarPosition = "left"
	}

	a.settings.AppearanceSettings = appearanceSettings

	// Save settings
	if err := a.saveSettingsToFile(); err != nil {
		return a.createAPIError("internal", ErrCodeInternalError, "Failed to save appearance settings", map[string]string{
			"error": err.Error(),
		})
	}

	// Emit appearance settings change event
	runtime.EventsEmit(a.ctx, "appearance:settings:changed", appearanceSettings)

	return nil
}

// UpdateEditorSettings updates editor settings
func (a *App) UpdateEditorSettings(editorSettings EditorSettings) error {
	// Validate editor settings
	if editorSettings.FontSize < 8 || editorSettings.FontSize > 24 {
		editorSettings.FontSize = 14
	}
	if editorSettings.TabSize < 2 || editorSettings.TabSize > 8 {
		editorSettings.TabSize = 4
	}
	if editorSettings.FontFamily == "" {
		editorSettings.FontFamily = "Monaco"
	}

	a.settings.EditorSettings = editorSettings

	// Save settings
	if err := a.saveSettingsToFile(); err != nil {
		return a.createAPIError("internal", ErrCodeInternalError, "Failed to save editor settings", map[string]string{
			"error": err.Error(),
		})
	}

	// Emit editor settings change event
	runtime.EventsEmit(a.ctx, "editor:settings:changed", editorSettings)

	return nil
}

// UpdateGenerationSettings updates generation settings
func (a *App) UpdateGenerationSettings(generationSettings GenerationSettings) error {
	// Validate generation settings
	if generationSettings.DefaultTemplate == "" {
		generationSettings.DefaultTemplate = "default"
	}
	if generationSettings.MaxWorkers < 1 || generationSettings.MaxWorkers > 16 {
		generationSettings.MaxWorkers = 4
	}

	a.settings.GenerationSettings = generationSettings

	// Save settings
	if err := a.saveSettingsToFile(); err != nil {
		return a.createAPIError("internal", ErrCodeInternalError, "Failed to save generation settings", map[string]string{
			"error": err.Error(),
		})
	}

	// Emit generation settings change event
	runtime.EventsEmit(a.ctx, "generation:settings:changed", generationSettings)

	return nil
}

// BackupSettingsAPI creates a backup of current settings (API method)
func (a *App) BackupSettingsAPI() error {
	if a.settings == nil {
		return nil // Nothing to backup
	}

	// Create backup using the migration file function
	return a.BackupSettings()
}

// RestoreSettingsFromBackupAPI restores settings from backup (API method)
func (a *App) RestoreSettingsFromBackupAPI() error {
	return a.RestoreSettingsFromBackup()
}

// GetSettingsBackupPath returns the path to the settings backup file
func (a *App) GetSettingsBackupPath() string {
	return a.GetSettingsFilePath() + ".backup"
}

// HasSettingsBackup checks if a settings backup exists
func (a *App) HasSettingsBackup() bool {
	backupPath := a.GetSettingsBackupPath()
	return a.fileExists(backupPath) == nil
}