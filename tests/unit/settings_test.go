package unit

import (
	"encoding/json"
	"os"
	"testing"

	"MCPWeaver/internal/app"
)

func TestGetDefaultSettings(t *testing.T) {
	// Create a test app instance
	testApp := &app.App{}

	// Test getting default settings
	settings, err := testApp.GetSettings()
	if err != nil {
		t.Fatalf("Failed to get default settings: %v", err)
	}

	// Verify default values
	if settings.Theme != "light" {
		t.Errorf("Expected default theme 'light', got '%s'", settings.Theme)
	}

	if settings.Language != "en" {
		t.Errorf("Expected default language 'en', got '%s'", settings.Language)
	}

	if !settings.AutoSave {
		t.Error("Expected auto-save to be enabled by default")
	}

	if settings.DefaultOutputPath != "./output" {
		t.Errorf("Expected default output path './output', got '%s'", settings.DefaultOutputPath)
	}

	// Test window settings
	if settings.WindowSettings.Width != 1200 {
		t.Errorf("Expected default window width 1200, got %d", settings.WindowSettings.Width)
	}

	if settings.WindowSettings.Height != 800 {
		t.Errorf("Expected default window height 800, got %d", settings.WindowSettings.Height)
	}

	// Test editor settings
	if settings.EditorSettings.FontSize != 14 {
		t.Errorf("Expected default font size 14, got %d", settings.EditorSettings.FontSize)
	}

	if settings.EditorSettings.FontFamily != "Monaco" {
		t.Errorf("Expected default font family 'Monaco', got '%s'", settings.EditorSettings.FontFamily)
	}

	// Test generation settings
	if settings.GenerationSettings.DefaultTemplate != "default" {
		t.Errorf("Expected default template 'default', got '%s'", settings.GenerationSettings.DefaultTemplate)
	}

	if !settings.GenerationSettings.EnableValidation {
		t.Error("Expected validation to be enabled by default")
	}

	if settings.GenerationSettings.MaxWorkers != 4 {
		t.Errorf("Expected default max workers 4, got %d", settings.GenerationSettings.MaxWorkers)
	}

	// Test notification settings
	if !settings.NotificationSettings.EnableDesktopNotifications {
		t.Error("Expected desktop notifications to be enabled by default")
	}

	if settings.NotificationSettings.NotificationPosition != "top-right" {
		t.Errorf("Expected default notification position 'top-right', got '%s'", settings.NotificationSettings.NotificationPosition)
	}

	// Test appearance settings
	if settings.AppearanceSettings.UITheme != "system" {
		t.Errorf("Expected default UI theme 'system', got '%s'", settings.AppearanceSettings.UITheme)
	}

	if settings.AppearanceSettings.WindowOpacity != 1.0 {
		t.Errorf("Expected default window opacity 1.0, got %f", settings.AppearanceSettings.WindowOpacity)
	}
}

func TestSettingsValidation(t *testing.T) {
	testApp := &app.App{}
	
	// Test valid theme values
	validThemes := []string{"light", "dark", "auto"}
	for _, theme := range validThemes {
		err := testApp.UpdateTheme(theme)
		if err != nil {
			t.Errorf("Valid theme '%s' should not produce error: %v", theme, err)
		}
	}

	// Test invalid theme value
	err := testApp.UpdateTheme("invalid")
	if err == nil {
		t.Error("Invalid theme should produce validation error")
	}

	// Test valid language
	err = testApp.UpdateLanguage("es")
	if err != nil {
		t.Errorf("Valid language should not produce error: %v", err)
	}

	// Test empty language
	err = testApp.UpdateLanguage("")
	if err == nil {
		t.Error("Empty language should produce validation error")
	}

	// Test window settings validation
	validWindow := app.WindowSettings{
		Width:  800,
		Height: 600,
	}
	err = testApp.UpdateWindowSettings(validWindow)
	if err != nil {
		t.Errorf("Valid window settings should not produce error: %v", err)
	}

	// Test invalid window settings (too small)
	invalidWindow := app.WindowSettings{
		Width:  500,  // Too small
		Height: 300,  // Too small
	}
	err = testApp.UpdateWindowSettings(invalidWindow)
	if err != nil {
		t.Errorf("Small window settings should be auto-corrected, not produce error: %v", err)
	}
}

func TestSettingsSerialization(t *testing.T) {
	testApp := &app.App{}
	
	// Get default settings
	settings, err := testApp.GetSettings()
	if err != nil {
		t.Fatalf("Failed to get settings: %v", err)
	}

	// Test JSON marshaling
	data, err := json.Marshal(settings)
	if err != nil {
		t.Fatalf("Failed to marshal settings: %v", err)
	}

	// Test JSON unmarshaling
	var unmarshaled app.AppSettings
	err = json.Unmarshal(data, &unmarshaled)
	if err != nil {
		t.Fatalf("Failed to unmarshal settings: %v", err)
	}

	// Verify that unmarshaled settings match original
	if unmarshaled.Theme != settings.Theme {
		t.Errorf("Theme mismatch after serialization: expected '%s', got '%s'", settings.Theme, unmarshaled.Theme)
	}

	if unmarshaled.Language != settings.Language {
		t.Errorf("Language mismatch after serialization: expected '%s', got '%s'", settings.Language, unmarshaled.Language)
	}

	if unmarshaled.AutoSave != settings.AutoSave {
		t.Errorf("AutoSave mismatch after serialization: expected %t, got %t", settings.AutoSave, unmarshaled.AutoSave)
	}
}

func TestSettingsFileOperations(t *testing.T) {
	// Create temporary directory for test
	tempDir, err := os.MkdirTemp("", "mcpweaver_settings_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create test app with custom settings path
	testApp := &app.App{}
	
	// Test settings file path generation
	settingsPath := testApp.GetSettingsFilePath()
	if settingsPath == "" {
		t.Error("Settings file path should not be empty")
	}

	// Test backup path generation
	backupPath := testApp.GetSettingsBackupPath()
	expectedBackupPath := settingsPath + ".backup"
	if backupPath != expectedBackupPath {
		t.Errorf("Expected backup path '%s', got '%s'", expectedBackupPath, backupPath)
	}

	// Test backup existence check (should be false initially)
	if testApp.HasSettingsBackup() {
		t.Error("Should not have backup initially")
	}
}

func TestSettingsMigration(t *testing.T) {
	// Test that migration system doesn't crash
	testApp := &app.App{}
	
	// Test loading settings with migration
	settings, err := testApp.MigrateAndLoadSettings()
	if err != nil {
		t.Fatalf("Failed to migrate and load settings: %v", err)
	}

	// Verify that migrated settings have all required fields
	if settings.NotificationSettings == (app.NotificationSettings{}) {
		t.Error("Migration should initialize notification settings")
	}

	if settings.AppearanceSettings == (app.AppearanceSettings{}) {
		t.Error("Migration should initialize appearance settings")
	}

	if settings.GenerationSettings.MaxWorkers == 0 {
		t.Error("Migration should set max workers to default value")
	}
}

func TestSettingsValidationEdgeCases(t *testing.T) {
	testApp := &app.App{}
	
	// Test editor settings validation
	editorSettings := app.EditorSettings{
		FontSize:   5,  // Too small
		TabSize:    1,  // Too small
		FontFamily: "", // Empty
	}
	
	err := testApp.UpdateEditorSettings(editorSettings)
	if err != nil {
		t.Errorf("Editor settings validation should auto-correct, not error: %v", err)
	}

	// Test generation settings validation
	generationSettings := app.GenerationSettings{
		MaxWorkers: 0, // Too small
	}
	
	err = testApp.UpdateGenerationSettings(generationSettings)
	if err != nil {
		t.Errorf("Generation settings validation should auto-correct, not error: %v", err)
	}

	// Test notification settings validation
	notificationSettings := app.NotificationSettings{
		NotificationDuration: 500,  // Too small
		SoundVolume:         -0.5,  // Invalid
	}
	
	err = testApp.UpdateNotificationSettings(notificationSettings)
	if err != nil {
		t.Errorf("Notification settings validation should auto-correct, not error: %v", err)
	}

	// Test appearance settings validation
	appearanceSettings := app.AppearanceSettings{
		WindowOpacity: 2.0,  // Too large
		FontScale:     0.1,  // Too small
	}
	
	err = testApp.UpdateAppearanceSettings(appearanceSettings)
	if err != nil {
		t.Errorf("Appearance settings validation should auto-correct, not error: %v", err)
	}
}

func TestSettingsUpdateEvents(t *testing.T) {
	testApp := &app.App{}
	
	// Note: In a real test environment, you would need to set up the Wails context
	// and event listener to properly test events. For now, we just test that
	// the methods don't crash when called.
	
	// Test theme update
	err := testApp.UpdateTheme("dark")
	if err != nil {
		t.Errorf("Theme update should not error: %v", err)
	}

	// Test language update
	err = testApp.UpdateLanguage("fr")
	if err != nil {
		t.Errorf("Language update should not error: %v", err)
	}

	// Test window settings update
	windowSettings := app.WindowSettings{
		Width:  1000,
		Height: 700,
	}
	err = testApp.UpdateWindowSettings(windowSettings)
	if err != nil {
		t.Errorf("Window settings update should not error: %v", err)
	}
}

func TestRecentProjectsManagement(t *testing.T) {
	testApp := &app.App{}
	
	// Test adding recent project
	err := testApp.AddRecentProject("test-project-1")
	if err != nil {
		t.Errorf("Adding recent project should not error: %v", err)
	}

	// Test adding empty project ID
	err = testApp.AddRecentProject("")
	if err == nil {
		t.Error("Adding empty project ID should produce validation error")
	}

	// Test clearing recent projects
	err = testApp.ClearRecentProjects()
	if err != nil {
		t.Errorf("Clearing recent projects should not error: %v", err)
	}
}

// Benchmark tests
func BenchmarkGetSettings(b *testing.B) {
	testApp := &app.App{}
	
	for i := 0; i < b.N; i++ {
		_, err := testApp.GetSettings()
		if err != nil {
			b.Fatalf("Failed to get settings: %v", err)
		}
	}
}

func BenchmarkSettingsValidation(b *testing.B) {
	testApp := &app.App{}
	settings, _ := testApp.GetSettings()
	
	for i := 0; i < b.N; i++ {
		err := testApp.UpdateSettings(*settings)
		if err != nil {
			b.Fatalf("Failed to update settings: %v", err)
		}
	}
}

func BenchmarkSettingsSerialization(b *testing.B) {
	testApp := &app.App{}
	settings, _ := testApp.GetSettings()
	
	for i := 0; i < b.N; i++ {
		_, err := json.Marshal(settings)
		if err != nil {
			b.Fatalf("Failed to marshal settings: %v", err)
		}
	}
}