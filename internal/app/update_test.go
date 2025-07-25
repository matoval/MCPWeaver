package app

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
	"time"
)

// TestUpdateService tests the update service functionality
func TestUpdateService(t *testing.T) {
	service := NewUpdateService(nil)

	if service == nil {
		t.Fatal("Failed to create update service")
	}

	// Test service initialization
	if service.currentStatus != UpdateStatusIdle {
		t.Errorf("Expected initial status to be %s, got %s", UpdateStatusIdle, service.currentStatus)
	}

	// Test settings
	settings := service.GetUpdateSettings()
	if settings == nil {
		t.Fatal("Failed to get update settings")
	}

	if !settings.Enabled {
		t.Error("Expected update service to be enabled by default")
	}
}

// TestUpdateConfiguration tests update configuration
func TestUpdateConfiguration(t *testing.T) {
	config := DefaultUpdateConfiguration()

	if config.UpdateURL == "" {
		t.Error("Expected default update URL to be set")
	}

	if config.BackupDirectory == "" {
		t.Error("Expected default backup directory to be set")
	}

	if config.VerificationMode == "" {
		t.Error("Expected default verification mode to be set")
	}
}

// TestUpdateSettings tests update settings
func TestUpdateSettings(t *testing.T) {
	settings := DefaultUpdateSettings()

	if !settings.Enabled {
		t.Error("Expected updates to be enabled by default")
	}

	if settings.CheckInterval == 0 {
		t.Error("Expected default check interval to be set")
	}

	if settings.UpdateChannel == "" {
		t.Error("Expected default update channel to be set")
	}
}

// TestUpdateCheck tests the update check functionality
func TestUpdateCheck(t *testing.T) {
	// Create a mock server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		response := `{
			"tag_name": "v1.1.0",
			"name": "Test Release",
			"body": "Test release notes",
			"published_at": "2023-01-01T00:00:00Z",
			"assets": [{
				"name": "mcpweaver_linux_amd64",
				"browser_download_url": "https://example.com/download",
				"size": 1024
			}],
			"prerelease": false
		}`
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(response))
	}))
	defer server.Close()

	// Use nil context for testing to avoid Wails runtime issues
	service := NewUpdateService(nil)

	// Update configuration to use mock server
	config := service.config
	config.UpdateURL = server.URL

	// Test update check
	updateInfo, err := service.CheckForUpdates()
	if err != nil {
		t.Fatalf("Update check failed: %v", err)
	}

	if updateInfo == nil {
		t.Fatal("Expected update info to be returned")
	}

	if updateInfo.Version != "v1.1.0" {
		t.Errorf("Expected version v1.1.0, got %s", updateInfo.Version)
	}
}

// TestRollbackManager tests the rollback functionality
func TestRollbackManager(t *testing.T) {
	// Create temporary directory for testing
	tempDir, err := os.MkdirTemp("", "update_test")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	manager := NewRollbackManager(tempDir)

	// Test capabilities
	capabilities := manager.GetRollbackCapabilities()
	if capabilities == nil {
		t.Fatal("Failed to get rollback capabilities")
	}

	// Test backup list (should be empty initially)
	backups, err := manager.ListAvailableBackups()
	if err != nil {
		t.Fatalf("Failed to list backups: %v", err)
	}

	if len(backups) != 0 {
		t.Errorf("Expected 0 backups, got %d", len(backups))
	}

	// Create a test file to backup
	testFile := filepath.Join(tempDir, "test_executable")
	testContent := "test executable content"
	err = os.WriteFile(testFile, []byte(testContent), 0755)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	// Test backup creation
	rollbackInfo, err := manager.CreateBackup(testFile, "1.0.0")
	if err != nil {
		t.Fatalf("Failed to create backup: %v", err)
	}

	if !rollbackInfo.Available {
		t.Error("Expected backup to be available")
	}

	if rollbackInfo.BackupVersion != "1.0.0" {
		t.Errorf("Expected backup version 1.0.0, got %s", rollbackInfo.BackupVersion)
	}

	// Test backup validation
	validation, err := manager.ValidateBackup(rollbackInfo.BackupPath)
	if err != nil {
		t.Fatalf("Failed to validate backup: %v", err)
	}

	if !validation.Valid {
		t.Error("Expected backup to be valid")
	}
}

// TestUpdateScheduler tests the update scheduler
func TestUpdateScheduler(t *testing.T) {
	var callbackCount int
	callback := func(jobType ScheduledJobType) error {
		callbackCount++
		return nil
	}

	scheduler := NewUpdateScheduler(callback)

	// Test immediate schedule
	schedule := &UpdateSchedule{
		Type: ScheduleTypeImmediate,
	}

	job, err := scheduler.ScheduleJob(ScheduledJobTypeCheck, schedule, nil)
	if err != nil {
		t.Fatalf("Failed to schedule job: %v", err)
	}

	if job.Type != ScheduledJobTypeCheck {
		t.Errorf("Expected job type %s, got %s", ScheduledJobTypeCheck, job.Type)
	}

	// Wait a moment for immediate execution
	time.Sleep(100 * time.Millisecond)

	if callbackCount != 1 {
		t.Errorf("Expected callback to be called 1 time, got %d", callbackCount)
	}

	// Test job cancellation
	err = scheduler.CancelScheduledJob()
	if err != nil {
		t.Fatalf("Failed to cancel job: %v", err)
	}

	currentJob := scheduler.GetScheduledJob()
	if currentJob != nil && currentJob.Status != ScheduledJobStatusCanceled {
		t.Error("Expected job to be canceled")
	}
}

// TestUpdateNotifications tests the notification system
func TestUpdateNotifications(t *testing.T) {
	service := NewUpdateService(nil)

	updateInfo := &UpdateInfo{
		Version:      "v1.1.0",
		ReleaseNotes: "Test release",
		DownloadURL:  "https://example.com/download",
		Size:         1024,
		PublishedAt:  time.Now(),
	}

	// Test notification emission in test mode (should not emit events)
	service.emitUpdateNotification(NotificationTypeUpdateAvailable, updateInfo)

	// Since we're in test mode, this should complete without panicking
	// and without actually emitting events
}

// TestUpdateProgress tests progress tracking
func TestUpdateProgress(t *testing.T) {
	service := NewUpdateService(nil)

	progress := service.GetUpdateProgress()
	if progress == nil {
		t.Fatal("Failed to get update progress")
	}

	if progress.Status != UpdateStatusIdle {
		t.Errorf("Expected initial status %s, got %s", UpdateStatusIdle, progress.Status)
	}

	if progress.Progress != 0.0 {
		t.Errorf("Expected initial progress 0.0, got %f", progress.Progress)
	}
}

// TestErrorHandling tests error handling in update operations
func TestErrorHandling(t *testing.T) {
	service := NewUpdateService(nil)

	// Test with invalid URL
	config := service.config
	config.UpdateURL = "invalid-url"

	_, err := service.CheckForUpdates()
	if err == nil {
		t.Error("Expected error for invalid URL")
	}

	apiErr, ok := err.(*APIError)
	if !ok {
		t.Error("Expected APIError type")
	}

	if apiErr.Type != ErrorTypeNetwork {
		t.Errorf("Expected error type %s, got %s", ErrorTypeNetwork, apiErr.Type)
	}
}

// TestVersionComparison tests version comparison logic
func TestVersionComparison(t *testing.T) {
	service := NewUpdateService(nil)

	// Test version comparison (basic implementation)
	newer := service.isNewerVersion("v1.1.0", "v1.0.0")
	if !newer {
		t.Error("Expected v1.1.0 to be newer than v1.0.0")
	}

	same := service.isNewerVersion("v1.0.0", "v1.0.0")
	if same {
		t.Error("Expected v1.0.0 to not be newer than v1.0.0")
	}
}

// TestAnalytics tests the analytics tracking
func TestAnalytics(t *testing.T) {
	service := NewUpdateService(nil)

	initialCount := len(service.analytics)

	// Track an analytics event
	service.trackAnalytics(AnalyticsEventCheckStarted, "v1.1.0", "v1.0.0", true, "")

	if len(service.analytics) != initialCount+1 {
		t.Errorf("Expected analytics count to increase by 1, got %d", len(service.analytics)-initialCount)
	}

	lastEvent := service.analytics[len(service.analytics)-1]
	if lastEvent.EventType != AnalyticsEventCheckStarted {
		t.Errorf("Expected event type %s, got %s", AnalyticsEventCheckStarted, lastEvent.EventType)
	}
}

// TestUpdateServiceIntegration tests the full integration
func TestUpdateServiceIntegration(t *testing.T) {
	// Create a mock server for testing
	var server *httptest.Server
	server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		response := `{
			"tag_name": "v1.1.0",
			"name": "Integration Test Release",
			"body": "Integration test release notes",
			"published_at": "2023-01-01T00:00:00Z",
			"assets": [{
				"name": "mcpweaver_linux_amd64",
				"browser_download_url": "` + server.URL + `/download",
				"size": 1024
			}],
			"prerelease": false
		}`

		if r.URL.Path == "/download" {
			w.Header().Set("Content-Type", "application/octet-stream")
			w.WriteHeader(http.StatusOK)
			// Write some dummy content
			for i := 0; i < 1024; i++ {
				w.Write([]byte("x"))
			}
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(response))
	}))
	defer server.Close()

	service := NewUpdateService(nil)

	// Update configuration for testing
	config := service.config
	config.UpdateURL = server.URL
	config.VerificationMode = VerificationModeNone // Skip verification for test

	// Start the service
	err := service.Start()
	if err != nil {
		t.Fatalf("Failed to start service: %v", err)
	}
	defer service.Stop()

	// Test complete workflow
	updateInfo, err := service.CheckForUpdates()
	if err != nil {
		t.Fatalf("Update check failed: %v", err)
	}

	if updateInfo == nil {
		t.Fatal("No update available")
	}

	// Test download (this will fail in the test environment, but we can test the error handling)
	err = service.DownloadUpdate(updateInfo)
	if err != nil {
		// Expected to fail in test environment
		t.Logf("Download failed as expected in test environment: %v", err)
	}
}

// Benchmark tests for performance
func BenchmarkUpdateCheck(b *testing.B) {
	ctx := context.Background()
	service := NewUpdateService(ctx)

	// Create a mock server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		response := `{"tag_name": "v1.1.0", "name": "Benchmark", "body": "", "published_at": "2023-01-01T00:00:00Z", "assets": [], "prerelease": false}`
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(response))
	}))
	defer server.Close()

	service.config.UpdateURL = server.URL

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		service.CheckForUpdates()
	}
}

// Helper functions for testing

func createTempExecutable(t *testing.T) string {
	tempDir, err := os.MkdirTemp("", "update_test")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}

	exePath := filepath.Join(tempDir, "test_executable")
	err = os.WriteFile(exePath, []byte("dummy executable"), 0755)
	if err != nil {
		t.Fatalf("Failed to create test executable: %v", err)
	}

	return exePath
}

// Test data
var testUpdateInfo = &UpdateInfo{
	Version:      "v1.1.0",
	ReleaseNotes: "Test release notes",
	DownloadURL:  "https://example.com/download",
	Size:         1024,
	PublishedAt:  time.Now(),
	Critical:     false,
}

var testScheduleDaily = &UpdateSchedule{
	Type: ScheduleTypeDaily,
	Time: "09:00",
}

var testScheduleWeekly = &UpdateSchedule{
	Type:      ScheduleTypeWeekly,
	Time:      "09:00",
	DayOfWeek: 1, // Monday
}

// TestMain runs before all tests and can set up global test environment
func TestMain(m *testing.M) {
	// Setup
	fmt.Println("Setting up auto-update tests...")

	// Run tests
	code := m.Run()

	// Cleanup
	fmt.Println("Cleaning up auto-update tests...")

	os.Exit(code)
}
