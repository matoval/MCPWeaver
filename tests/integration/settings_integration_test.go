package integration

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
	"time"

	"MCPWeaver/internal/app"
	"MCPWeaver/internal/database"
)

func TestSettingsIntegrationFullWorkflow(t *testing.T) {
	// Create temporary directory for test
	tempDir, err := os.MkdirTemp("", "mcpweaver_integration_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create test database
	dbPath := filepath.Join(tempDir, "test.db")
	db, err := database.NewDB(dbPath)
	if err != nil {
		t.Fatalf("Failed to create test database: %v", err)
	}
	defer db.Close()

	// Create app instance with test database
	testApp := &app.App{}
	// Note: In a real integration test, you would inject the database dependency

	// Test 1: Load default settings
	settings, err := testApp.GetSettings()
	if err != nil {
		t.Fatalf("Failed to get initial settings: %v", err)
	}

	originalTheme := settings.Theme
	if originalTheme != "light" {
		t.Errorf("Expected default theme 'light', got '%s'", originalTheme)
	}

	// Test 2: Update settings
	newTheme := "dark"
	err = testApp.UpdateTheme(newTheme)
	if err != nil {
		t.Fatalf("Failed to update theme: %v", err)
	}

	// Verify settings were updated
	updatedSettings, err := testApp.GetSettings()
	if err != nil {
		t.Fatalf("Failed to get updated settings: %v", err)
	}

	if updatedSettings.Theme != newTheme {
		t.Errorf("Expected theme '%s' after update, got '%s'", newTheme, updatedSettings.Theme)
	}

	// Test 3: Update complex settings (editor settings)
	newEditorSettings := app.EditorSettings{
		FontSize:        16,
		FontFamily:      "Fira Code",
		TabSize:         2,
		WordWrap:        false,
		LineNumbers:     true,
		SyntaxHighlight: true,
	}

	err = testApp.UpdateEditorSettings(newEditorSettings)
	if err != nil {
		t.Fatalf("Failed to update editor settings: %v", err)
	}

	// Verify editor settings were updated
	settingsAfterEditor, err := testApp.GetSettings()
	if err != nil {
		t.Fatalf("Failed to get settings after editor update: %v", err)
	}

	if settingsAfterEditor.EditorSettings.FontSize != newEditorSettings.FontSize {
		t.Errorf("Expected font size %d, got %d", newEditorSettings.FontSize, settingsAfterEditor.EditorSettings.FontSize)
	}

	if settingsAfterEditor.EditorSettings.FontFamily != newEditorSettings.FontFamily {
		t.Errorf("Expected font family '%s', got '%s'", newEditorSettings.FontFamily, settingsAfterEditor.EditorSettings.FontFamily)
	}

	// Test 4: Settings persistence
	// Simulate app restart by creating new app instance
	testApp2 := &app.App{}
	
	persistedSettings, err := testApp2.GetSettings()
	if err != nil {
		t.Fatalf("Failed to get persisted settings: %v", err)
	}

	if persistedSettings.Theme != newTheme {
		t.Errorf("Theme not persisted: expected '%s', got '%s'", newTheme, persistedSettings.Theme)
	}

	if persistedSettings.EditorSettings.FontSize != newEditorSettings.FontSize {
		t.Errorf("Editor settings not persisted: expected font size %d, got %d", newEditorSettings.FontSize, persistedSettings.EditorSettings.FontSize)
	}

	// Test 5: Settings reset
	err = testApp2.ResetSettings()
	if err != nil {
		t.Fatalf("Failed to reset settings: %v", err)
	}

	resetSettings, err := testApp2.GetSettings()
	if err != nil {
		t.Fatalf("Failed to get reset settings: %v", err)
	}

	if resetSettings.Theme != "light" {
		t.Errorf("Settings not reset: expected theme 'light', got '%s'", resetSettings.Theme)
	}

	if resetSettings.EditorSettings.FontSize != 14 {
		t.Errorf("Settings not reset: expected font size 14, got %d", resetSettings.EditorSettings.FontSize)
	}
}

func TestSettingsExportImportIntegration(t *testing.T) {
	// Create temporary directory for test
	tempDir, err := os.MkdirTemp("", "mcpweaver_export_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	testApp := &app.App{}

	// Modify some settings
	err = testApp.UpdateTheme("dark")
	if err != nil {
		t.Fatalf("Failed to update theme: %v", err)
	}

	err = testApp.UpdateLanguage("es")
	if err != nil {
		t.Fatalf("Failed to update language: %v", err)
	}

	// Get current settings for comparison
	originalSettings, err := testApp.GetSettings()
	if err != nil {
		t.Fatalf("Failed to get original settings: %v", err)
	}

	// Note: In a real integration test, you would test the actual export/import
	// functionality with file dialogs. For now, we test the underlying logic.

	// Simulate export by manually creating an export file
	exportData, err := json.MarshalIndent(originalSettings, "", "  ")
	if err != nil {
		t.Fatalf("Failed to marshal settings: %v", err)
	}

	exportFile := filepath.Join(tempDir, "exported_settings.json")
	err = os.WriteFile(exportFile, exportData, 0644)
	if err != nil {
		t.Fatalf("Failed to write export file: %v", err)
	}

	// Reset settings
	err = testApp.ResetSettings()
	if err != nil {
		t.Fatalf("Failed to reset settings: %v", err)
	}

	// Verify settings were reset
	resetSettings, err := testApp.GetSettings()
	if err != nil {
		t.Fatalf("Failed to get reset settings: %v", err)
	}

	if resetSettings.Theme == originalSettings.Theme {
		t.Error("Settings should be different after reset")
	}

	// Simulate import by manually reading and applying the export file
	importData, err := os.ReadFile(exportFile)
	if err != nil {
		t.Fatalf("Failed to read import file: %v", err)
	}

	var importedSettings app.AppSettings
	err = json.Unmarshal(importData, &importedSettings)
	if err != nil {
		t.Fatalf("Failed to unmarshal imported settings: %v", err)
	}

	// Apply imported settings
	err = testApp.UpdateSettings(importedSettings)
	if err != nil {
		t.Fatalf("Failed to apply imported settings: %v", err)
	}

	// Verify settings were restored
	restoredSettings, err := testApp.GetSettings()
	if err != nil {
		t.Fatalf("Failed to get restored settings: %v", err)
	}

	if restoredSettings.Theme != originalSettings.Theme {
		t.Errorf("Theme not restored: expected '%s', got '%s'", originalSettings.Theme, restoredSettings.Theme)
	}

	if restoredSettings.Language != originalSettings.Language {
		t.Errorf("Language not restored: expected '%s', got '%s'", originalSettings.Language, restoredSettings.Language)
	}
}

func TestSettingsMigrationIntegration(t *testing.T) {
	// Create temporary directory for test
	tempDir, err := os.MkdirTemp("", "mcpweaver_migration_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create an "old" settings file without notification and appearance settings
	oldSettings := map[string]interface{}{
		"theme":             "light",
		"language":          "en",
		"autoSave":          true,
		"defaultOutputPath": "./output",
		"recentProjects":    []string{},
		"recentFiles":       []string{},
		"windowSettings": map[string]interface{}{
			"width":     1200,
			"height":    800,
			"maximized": false,
			"x":         100,
			"y":         100,
		},
		"editorSettings": map[string]interface{}{
			"fontSize":        14,
			"fontFamily":      "Monaco",
			"tabSize":         4,
			"wordWrap":        true,
			"lineNumbers":     true,
			"syntaxHighlight": true,
		},
		"generationSettings": map[string]interface{}{
			"defaultTemplate":     "default",
			"enableValidation":    true,
			"autoOpenOutput":      true,
			"showAdvancedOptions": false,
			"backupOnGenerate":    true,
			"customTemplates":     []string{},
			// Missing: performanceMode, maxWorkers
		},
		// Missing: notificationSettings, appearanceSettings
	}

	oldSettingsData, err := json.MarshalIndent(oldSettings, "", "  ")
	if err != nil {
		t.Fatalf("Failed to marshal old settings: %v", err)
	}

	// Write old settings to a file that would be read by the migration system
	settingsFile := filepath.Join(tempDir, "settings.json")
	err = os.WriteFile(settingsFile, oldSettingsData, 0644)
	if err != nil {
		t.Fatalf("Failed to write old settings file: %v", err)
	}

	// Create app and run migration
	testApp := &app.App{}
	
	// Test migration
	migratedSettings, err := testApp.MigrateAndLoadSettings()
	if err != nil {
		t.Fatalf("Failed to migrate settings: %v", err)
	}

	// Verify that missing settings were added with defaults
	if migratedSettings.NotificationSettings == (app.NotificationSettings{}) {
		t.Error("Migration should have added notification settings")
	}

	if migratedSettings.AppearanceSettings == (app.AppearanceSettings{}) {
		t.Error("Migration should have added appearance settings")
	}

	if migratedSettings.GenerationSettings.MaxWorkers == 0 {
		t.Error("Migration should have added max workers to generation settings")
	}

	// Verify that existing settings were preserved
	if migratedSettings.Theme != "light" {
		t.Errorf("Migration should preserve existing theme: expected 'light', got '%s'", migratedSettings.Theme)
	}

	if migratedSettings.EditorSettings.FontFamily != "Monaco" {
		t.Errorf("Migration should preserve existing editor settings: expected 'Monaco', got '%s'", migratedSettings.EditorSettings.FontFamily)
	}

	// Verify default values for new settings
	if !migratedSettings.NotificationSettings.EnableDesktopNotifications {
		t.Error("New notification settings should have desktop notifications enabled by default")
	}

	if migratedSettings.AppearanceSettings.UITheme != "system" {
		t.Errorf("New appearance settings should have 'system' theme by default, got '%s'", migratedSettings.AppearanceSettings.UITheme)
	}
}

func TestSettingsValidationIntegration(t *testing.T) {
	testApp := &app.App{}

	// Test validation of complex settings updates
	testCases := []struct {
		name     string
		settings app.AppSettings
		expectError bool
	}{
		{
			name: "Valid settings",
			settings: app.AppSettings{
				Theme:    "dark",
				Language: "fr",
				WindowSettings: app.WindowSettings{
					Width:  1000,
					Height: 700,
				},
				EditorSettings: app.EditorSettings{
					FontSize:   16,
					FontFamily: "Fira Code",
					TabSize:    2,
				},
			},
			expectError: false,
		},
		{
			name: "Invalid theme but valid other settings",
			settings: app.AppSettings{
				Theme:    "invalid-theme",
				Language: "en",
			},
			expectError: true,
		},
		{
			name: "Empty language",
			settings: app.AppSettings{
				Theme:    "light",
				Language: "",
			},
			expectError: false, // Should be auto-corrected
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := testApp.UpdateSettings(tc.settings)
			
			if tc.expectError && err == nil {
				t.Error("Expected validation error but got none")
			}
			
			if !tc.expectError && err != nil {
				t.Errorf("Expected no error but got: %v", err)
			}
		})
	}
}

func TestSettingsPerformanceIntegration(t *testing.T) {
	testApp := &app.App{}

	// Test that settings operations complete within reasonable time
	start := time.Now()
	
	// Perform multiple settings operations
	for i := 0; i < 100; i++ {
		settings, err := testApp.GetSettings()
		if err != nil {
			t.Fatalf("Failed to get settings on iteration %d: %v", i, err)
		}
		
		// Modify a setting
		if i%2 == 0 {
			settings.Theme = "dark"
		} else {
			settings.Theme = "light"
		}
		
		err = testApp.UpdateSettings(*settings)
		if err != nil {
			t.Fatalf("Failed to update settings on iteration %d: %v", i, err)
		}
	}
	
	elapsed := time.Since(start)
	
	// Settings operations should complete quickly (under 1 second for 100 operations)
	if elapsed > time.Second {
		t.Errorf("Settings operations took too long: %v", elapsed)
	}
}

func TestRecentProjectsIntegration(t *testing.T) {
	testApp := &app.App{}

	// Add multiple recent projects
	projectIDs := []string{"proj1", "proj2", "proj3", "proj4", "proj5"}
	
	for _, id := range projectIDs {
		err := testApp.AddRecentProject(id)
		if err != nil {
			t.Fatalf("Failed to add recent project %s: %v", id, err)
		}
	}

	// Verify recent projects were added
	settings, err := testApp.GetSettings()
	if err != nil {
		t.Fatalf("Failed to get settings: %v", err)
	}

	if len(settings.RecentProjects) != len(projectIDs) {
		t.Errorf("Expected %d recent projects, got %d", len(projectIDs), len(settings.RecentProjects))
	}

	// Verify order (most recent first)
	for i, expectedID := range []string{"proj5", "proj4", "proj3", "proj2", "proj1"} {
		if i < len(settings.RecentProjects) && settings.RecentProjects[i] != expectedID {
			t.Errorf("Expected recent project at index %d to be '%s', got '%s'", i, expectedID, settings.RecentProjects[i])
		}
	}

	// Test clearing recent projects
	err = testApp.ClearRecentProjects()
	if err != nil {
		t.Fatalf("Failed to clear recent projects: %v", err)
	}

	// Verify recent projects were cleared
	settingsAfterClear, err := testApp.GetSettings()
	if err != nil {
		t.Fatalf("Failed to get settings after clear: %v", err)
	}

	if len(settingsAfterClear.RecentProjects) != 0 {
		t.Errorf("Expected 0 recent projects after clear, got %d", len(settingsAfterClear.RecentProjects))
	}
}