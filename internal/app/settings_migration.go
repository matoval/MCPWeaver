package app

import (
	"encoding/json"
	"fmt"
	"log"
)

// SettingsMigration represents a settings migration
type SettingsMigration struct {
	Version     int
	Name        string
	Description string
	Migrate     func(*AppSettings) error
}

// settingsMigrations contains all settings migrations in order
var settingsMigrations = []SettingsMigration{
	{
		Version:     1,
		Name:        "add_notification_settings",
		Description: "Add notification and appearance settings",
		Migrate: func(settings *AppSettings) error {
			// Initialize notification settings if missing
			if settings.NotificationSettings == (NotificationSettings{}) {
				settings.NotificationSettings = NotificationSettings{
					EnableDesktopNotifications: true,
					EnableSoundNotifications:   false,
					NotificationPosition:       "top-right",
					NotificationDuration:       5000,
					SoundVolume:                0.5,
					ShowGenerationProgress:     true,
					ShowErrorNotifications:     true,
					ShowSuccessNotifications:   true,
				}
			}

			// Initialize appearance settings if missing
			if settings.AppearanceSettings == (AppearanceSettings{}) {
				settings.AppearanceSettings = AppearanceSettings{
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
				}
			}

			// Add new generation settings if missing
			if settings.GenerationSettings.MaxWorkers == 0 {
				settings.GenerationSettings.MaxWorkers = 4
			}

			return nil
		},
	},
}

// migrateSettings runs all pending settings migrations
func (a *App) migrateSettings(settings *AppSettings) error {
	// Get current settings version (stored in the settings file)
	currentVersion := a.getSettingsVersion(settings)

	log.Printf("Current settings version: %d", currentVersion)

	// Apply pending migrations
	for _, migration := range settingsMigrations {
		if migration.Version > currentVersion {
			log.Printf("Applying settings migration %d: %s", migration.Version, migration.Name)

			if err := migration.Migrate(settings); err != nil {
				return fmt.Errorf("failed to apply settings migration %d: %w", migration.Version, err)
			}

			// Update settings version
			a.setSettingsVersion(settings, migration.Version)

			log.Printf("Successfully applied settings migration %d", migration.Version)
		}
	}

	log.Printf("Settings migrations completed")
	return nil
}

// getSettingsVersion returns the current settings version
func (a *App) getSettingsVersion(settings *AppSettings) int {
	// For now, we'll use a simple heuristic to determine version
	// In the future, this could be stored as a dedicated field

	// Version 0: Original settings without notification/appearance settings
	if settings.NotificationSettings == (NotificationSettings{}) {
		return 0
	}

	// Version 1: Has notification and appearance settings
	return 1
}

// setSettingsVersion sets the settings version (placeholder for future use)
func (a *App) setSettingsVersion(settings *AppSettings, version int) {
	// In the future, we could add a Version field to AppSettings
	// For now, this is a no-op as we determine version by presence of fields
}

// MigrateAndLoadSettings loads settings from file and applies migrations
func (a *App) MigrateAndLoadSettings() (*AppSettings, error) {
	// Load settings from file
	settings, err := a.loadSettingsFromFile()
	if err != nil {
		log.Printf("Failed to load settings from file: %v", err)
		settings = getDefaultSettings()
	}

	// Apply migrations
	if err := a.migrateSettings(settings); err != nil {
		return nil, fmt.Errorf("failed to migrate settings: %w", err)
	}

	// Validate migrated settings
	if err := a.validateSettings(settings); err != nil {
		return nil, fmt.Errorf("failed to validate migrated settings: %w", err)
	}

	// Save migrated settings back to file
	a.settings = settings
	if err := a.saveSettingsToFile(); err != nil {
		log.Printf("Warning: failed to save migrated settings: %v", err)
	}

	return settings, nil
}

// BackupSettings creates a backup of current settings
func (a *App) BackupSettings() error {
	if a.settings == nil {
		return nil // Nothing to backup
	}

	// Create backup path
	backupPath := a.GetSettingsFilePath() + ".backup"

	// Marshal settings to JSON
	data, err := json.MarshalIndent(a.settings, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal settings for backup: %w", err)
	}

	// Write backup file
	return a.WriteFile(backupPath, string(data))
}

// RestoreSettingsFromBackup restores settings from backup
func (a *App) RestoreSettingsFromBackup() error {
	backupPath := a.GetSettingsFilePath() + ".backup"

	// Check if backup exists
	if a.fileExists(backupPath) != nil {
		return a.createAPIError("file_system", ErrCodeFileAccess, "Settings backup not found", map[string]string{
			"path": backupPath,
		})
	}

	// Read backup file
	data, err := a.ReadFile(backupPath)
	if err != nil {
		return a.createAPIError("file_system", ErrCodeFileAccess, "Failed to read settings backup", map[string]string{
			"path":  backupPath,
			"error": err.Error(),
		})
	}

	// Unmarshal settings
	var settings AppSettings
	if err := json.Unmarshal([]byte(data), &settings); err != nil {
		return a.createAPIError("validation", ErrCodeValidation, "Invalid settings backup format", map[string]string{
			"path":  backupPath,
			"error": err.Error(),
		})
	}

	// Validate restored settings
	if err := a.validateSettings(&settings); err != nil {
		return err
	}

	// Update current settings
	a.settings = &settings

	// Save restored settings
	if err := a.saveSettingsToFile(); err != nil {
		return a.createAPIError("internal", ErrCodeInternalError, "Failed to save restored settings", map[string]string{
			"error": err.Error(),
		})
	}

	// Send notification
	a.emitNotification("success", "Settings Restored", "Settings have been restored from backup")

	return nil
}
