package app

import (
	"encoding/json"
	"testing"
	"time"
)

func TestNotificationService_Basic(t *testing.T) {
	// Create notification service in test mode
	service := NewNotificationService(nil, nil)
	defer service.Stop()

	// Test service creation
	if service == nil {
		t.Fatal("Failed to create notification service")
	}

	// Test start
	err := service.Start()
	if err != nil {
		t.Fatalf("Failed to start notification service: %v", err)
	}

	// Test that service is properly initialized
	if service.config == nil {
		t.Error("Service config not initialized")
	}

	if service.activeToasts == nil {
		t.Error("Active toasts map not initialized")
	}

	if service.history == nil {
		t.Error("History slice not initialized")
	}
}

func TestNotificationService_ShowToast(t *testing.T) {
	service := NewNotificationService(nil, nil)
	defer service.Stop()
	service.Start()

	tests := []struct {
		name        string
		toast       *ToastNotification
		expectError bool
	}{
		{
			name: "Valid info toast",
			toast: &ToastNotification{
				Type:     ToastTypeInfo,
				Title:    "Test Info",
				Message:  "This is a test info message",
				Category: CategoryGeneral,
				Priority: NotificationPriorityMedium,
			},
			expectError: false,
		},
		{
			name: "Valid success toast with actions",
			toast: &ToastNotification{
				Type:     ToastTypeSuccess,
				Title:    "Success",
				Message:  "Operation completed successfully",
				Category: CategoryProject,
				Priority: NotificationPriorityMedium,
				Actions: []NotificationActionBtn{
					{ID: "view", Label: "View Details", Type: ActionTypeViewDetails},
					{ID: "dismiss", Label: "Dismiss", Type: ActionTypeDismiss},
				},
			},
			expectError: false,
		},
		{
			name: "Error toast with high priority",
			toast: &ToastNotification{
				Type:     ToastTypeError,
				Title:    "Error",
				Message:  "An error occurred",
				Category: CategoryError,
				Priority: NotificationPriorityHigh,
			},
			expectError: false,
		},
		{
			name: "Loading toast with progress",
			toast: &ToastNotification{
				Type:     ToastTypeLoading,
				Title:    "Processing",
				Message:  "Please wait...",
				Category: CategoryGeneration,
				Priority: NotificationPriorityMedium,
				Progress: &NotificationProgress{
					Current: 50,
					Total:   100,
					Percent: 50,
					Label:   "Generating MCP server...",
				},
			},
			expectError: false,
		},
		{
			name: "Persistent critical toast",
			toast: &ToastNotification{
				Type:       ToastTypeError,
				Title:      "Critical Error",
				Message:    "System failure detected",
				Category:   CategoryError,
				Priority:   NotificationPriorityCritical,
				Persistent: true,
			},
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := service.ShowToast(tt.toast)

			if tt.expectError && err == nil {
				t.Error("Expected error but got none")
			}

			if !tt.expectError && err != nil {
				t.Errorf("Unexpected error: %v", err)
			}

			if !tt.expectError {
				// Check that toast was added to active toasts
				activeToasts := service.GetActiveToasts()
				found := false
				for _, active := range activeToasts {
					if active.Title == tt.toast.Title && active.Message == tt.toast.Message {
						found = true
						// Verify ID was generated
						if active.ID == "" {
							t.Error("Toast ID was not generated")
						}
						// Verify timestamps
						if active.CreatedAt.IsZero() {
							t.Error("Toast CreatedAt was not set")
						}
						break
					}
				}
				if !found {
					t.Error("Toast not found in active toasts")
				}

				// Check that toast was added to history
				history := service.GetNotificationHistory(100, 0, "")
				found = false
				for _, hist := range history {
					if hist.Title == tt.toast.Title && hist.Message == tt.toast.Message {
						found = true
						break
					}
				}
				if !found {
					t.Error("Toast not found in history")
				}
			}
		})
	}
}

func TestNotificationService_SystemNotification(t *testing.T) {
	service := NewNotificationService(nil, nil)
	defer service.Stop()
	service.Start()

	tests := []struct {
		name         string
		notification *SystemNotification
		expectError  bool
	}{
		{
			name: "Valid system notification",
			notification: &SystemNotification{
				Title:    "System Alert",
				Body:     "Important system message",
				Urgency:  SystemUrgencyNormal,
				Category: CategorySystem,
			},
			expectError: false,
		},
		{
			name: "Critical system notification with actions",
			notification: &SystemNotification{
				Title:    "Critical System Alert",
				Body:     "Immediate attention required",
				Urgency:  SystemUrgencyCritical,
				Category: CategorySecurity,
				Actions: []NotificationActionBtn{
					{ID: "fix", Label: "Fix Now", Type: ActionTypeInstallNow},
					{ID: "later", Label: "Fix Later", Type: ActionTypeInstallLater},
				},
			},
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := service.ShowSystemNotification(tt.notification)

			if tt.expectError && err == nil {
				t.Error("Expected error but got none")
			}

			if !tt.expectError && err != nil {
				t.Errorf("Unexpected error: %v", err)
			}

			if !tt.expectError {
				// Check that notification was added to history
				history := service.GetNotificationHistory(100, 0, "")
				found := false
				for _, hist := range history {
					if hist.Title == tt.notification.Title && hist.Message == tt.notification.Body {
						found = true
						if hist.Type != "system" {
							t.Error("Notification type should be 'system'")
						}
						break
					}
				}
				if !found {
					t.Error("System notification not found in history")
				}
			}
		})
	}
}

func TestNotificationService_DismissToast(t *testing.T) {
	service := NewNotificationService(nil, nil)
	defer service.Stop()
	service.Start()

	// Create a test toast
	toast := &ToastNotification{
		Type:     ToastTypeInfo,
		Title:    "Test Toast",
		Message:  "Test message",
		Category: CategoryGeneral,
		Priority: NotificationPriorityMedium,
	}

	err := service.ShowToast(toast)
	if err != nil {
		t.Fatalf("Failed to show toast: %v", err)
	}

	// Get the ID of the created toast
	activeToasts := service.GetActiveToasts()
	if len(activeToasts) == 0 {
		t.Fatal("No active toasts found")
	}

	toastID := activeToasts[0].ID

	// Test dismissing the toast
	err = service.DismissToast(toastID)
	if err != nil {
		t.Errorf("Failed to dismiss toast: %v", err)
	}

	// Check that toast is no longer active
	activeToasts = service.GetActiveToasts()
	for _, active := range activeToasts {
		if active.ID == toastID {
			t.Error("Toast still active after dismissal")
		}
	}

	// Check that history shows dismissal
	history := service.GetNotificationHistory(100, 0, "")
	found := false
	for _, hist := range history {
		if hist.ID == toastID {
			if hist.DismissedAt == nil {
				t.Error("Toast dismissal time not recorded in history")
			}
			found = true
			break
		}
	}
	if !found {
		t.Error("Toast not found in history")
	}

	// Test dismissing non-existent toast
	err = service.DismissToast("non-existent-id")
	if err == nil {
		t.Error("Expected error when dismissing non-existent toast")
	}
}

func TestNotificationService_DismissAllToasts(t *testing.T) {
	service := NewNotificationService(nil, nil)
	defer service.Stop()
	service.Start()

	// Create multiple test toasts
	toasts := []*ToastNotification{
		{Type: ToastTypeInfo, Title: "Toast 1", Message: "Message 1", Category: CategoryGeneral, Priority: NotificationPriorityMedium},
		{Type: ToastTypeSuccess, Title: "Toast 2", Message: "Message 2", Category: CategoryProject, Priority: NotificationPriorityMedium},
		{Type: ToastTypeWarning, Title: "Toast 3", Message: "Message 3", Category: CategoryValidation, Priority: NotificationPriorityHigh},
	}

	for _, toast := range toasts {
		err := service.ShowToast(toast)
		if err != nil {
			t.Fatalf("Failed to show toast: %v", err)
		}
	}

	// Verify toasts are active
	activeToasts := service.GetActiveToasts()
	if len(activeToasts) != len(toasts) {
		t.Errorf("Expected %d active toasts, got %d", len(toasts), len(activeToasts))
	}

	// Dismiss all toasts
	err := service.DismissAllToasts()
	if err != nil {
		t.Errorf("Failed to dismiss all toasts: %v", err)
	}

	// Check that no toasts are active
	activeToasts = service.GetActiveToasts()
	if len(activeToasts) != 0 {
		t.Errorf("Expected 0 active toasts after dismissing all, got %d", len(activeToasts))
	}
}

func TestNotificationService_MarkAsRead(t *testing.T) {
	service := NewNotificationService(nil, nil)
	defer service.Stop()
	service.Start()

	// Create a test toast
	toast := &ToastNotification{
		Type:     ToastTypeInfo,
		Title:    "Test Toast",
		Message:  "Test message",
		Category: CategoryGeneral,
		Priority: NotificationPriorityMedium,
	}

	err := service.ShowToast(toast)
	if err != nil {
		t.Fatalf("Failed to show toast: %v", err)
	}

	// Get the ID of the created toast
	activeToasts := service.GetActiveToasts()
	if len(activeToasts) == 0 {
		t.Fatal("No active toasts found")
	}

	toastID := activeToasts[0].ID

	// Mark as read
	err = service.MarkAsRead(toastID)
	if err != nil {
		t.Errorf("Failed to mark notification as read: %v", err)
	}

	// Check that history shows read time
	history := service.GetNotificationHistory(100, 0, "")
	found := false
	for _, hist := range history {
		if hist.ID == toastID {
			if hist.ReadAt == nil {
				t.Error("Notification read time not recorded in history")
			}
			found = true
			break
		}
	}
	if !found {
		t.Error("Notification not found in history")
	}

	// Test marking non-existent notification as read
	err = service.MarkAsRead("non-existent-id")
	if err == nil {
		t.Error("Expected error when marking non-existent notification as read")
	}
}

func TestNotificationService_ExecuteAction(t *testing.T) {
	service := NewNotificationService(nil, nil)
	defer service.Stop()
	service.Start()

	// Create a test toast with actions
	toast := &ToastNotification{
		Type:     ToastTypeInfo,
		Title:    "Test Toast",
		Message:  "Test message",
		Category: CategoryGeneral,
		Priority: NotificationPriorityMedium,
		Actions: []NotificationActionBtn{
			{ID: "view", Label: "View Details", Type: ActionTypeViewDetails},
			{ID: "dismiss", Label: "Dismiss", Type: ActionTypeDismiss},
		},
	}

	err := service.ShowToast(toast)
	if err != nil {
		t.Fatalf("Failed to show toast: %v", err)
	}

	// Get the ID of the created toast
	activeToasts := service.GetActiveToasts()
	if len(activeToasts) == 0 {
		t.Fatal("No active toasts found")
	}

	toastID := activeToasts[0].ID

	// Execute action
	err = service.ExecuteAction(toastID, "view", map[string]interface{}{"test": "data"})
	if err != nil {
		t.Errorf("Failed to execute action: %v", err)
	}

	// Check that history shows interaction
	history := service.GetNotificationHistory(100, 0, "")
	found := false
	for _, hist := range history {
		if hist.ID == toastID {
			if hist.InteractedAt == nil {
				t.Error("Notification interaction time not recorded in history")
			}
			if hist.ActionTaken != "view" {
				t.Errorf("Expected action 'view', got '%s'", hist.ActionTaken)
			}
			found = true
			break
		}
	}
	if !found {
		t.Error("Notification not found in history")
	}

	// Test executing non-existent action
	err = service.ExecuteAction(toastID, "non-existent", nil)
	if err == nil {
		t.Error("Expected error when executing non-existent action")
	}

	// Test executing action on non-existent notification
	err = service.ExecuteAction("non-existent-id", "view", nil)
	if err == nil {
		t.Error("Expected error when executing action on non-existent notification")
	}
}

func TestNotificationService_FilteringAndThrottling(t *testing.T) {
	service := NewNotificationService(nil, nil)
	defer service.Stop()

	// Configure throttling
	config := DefaultNotificationSystem()
	config.ThrottleSettings.Enabled = true
	config.ThrottleSettings.MaxPerMinute = 2
	config.ThrottleSettings.ByCategory = map[NotificationCategory]ThrottleRule{
		CategoryError: {
			MaxPerMinute: 1,
			MaxPerHour:   5,
		},
	}
	service.UpdateConfig(&config)
	service.Start()

	// Test throttling by sending multiple notifications quickly
	toast := &ToastNotification{
		Type:     ToastTypeError,
		Title:    "Error",
		Message:  "Test error",
		Category: CategoryError,
		Priority: NotificationPriorityHigh,
	}

	// First notification should succeed
	err := service.ShowToast(toast)
	if err != nil {
		t.Errorf("First notification failed: %v", err)
	}

	// Second notification should be throttled (queued)
	err = service.ShowToast(toast)
	if err != nil {
		t.Errorf("Second notification failed: %v", err)
	}

	// Check active toasts (should only have one due to throttling)
	activeToasts := service.GetActiveToasts()
	if len(activeToasts) > 1 {
		t.Errorf("Expected at most 1 active toast due to throttling, got %d", len(activeToasts))
	}
}

func TestNotificationService_DoNotDisturbMode(t *testing.T) {
	service := NewNotificationService(nil, nil)
	defer service.Stop()

	// Configure do not disturb mode
	config := DefaultNotificationSystem()
	config.DoNotDisturbMode = true
	config.DoNotDisturbSchedule = &DoNotDisturbSchedule{
		Enabled:     true,
		AllowUrgent: true,
		StartTime:   "00:00",
		EndTime:     "23:59",
		Days:        []time.Weekday{time.Sunday, time.Monday, time.Tuesday, time.Wednesday, time.Thursday, time.Friday, time.Saturday},
	}
	service.UpdateConfig(&config)
	service.Start()

	// Test normal priority notification (should be queued)
	toast := &ToastNotification{
		Type:     ToastTypeInfo,
		Title:    "Info",
		Message:  "Normal info",
		Category: CategoryGeneral,
		Priority: NotificationPriorityMedium,
	}

	err := service.ShowToast(toast)
	if err != nil {
		t.Errorf("Normal notification failed: %v", err)
	}

	// Should have no active toasts (queued due to DND)
	activeToasts := service.GetActiveToasts()
	if len(activeToasts) > 0 {
		t.Errorf("Expected no active toasts in DND mode, got %d", len(activeToasts))
	}

	// Test critical notification (should bypass DND)
	criticalToast := &ToastNotification{
		Type:     ToastTypeError,
		Title:    "Critical Error",
		Message:  "Critical error message",
		Category: CategoryError,
		Priority: NotificationPriorityCritical,
	}

	err = service.ShowToast(criticalToast)
	if err != nil {
		t.Errorf("Critical notification failed: %v", err)
	}

	// Should have one active toast (critical bypassed DND)
	activeToasts = service.GetActiveToasts()
	if len(activeToasts) != 1 {
		t.Errorf("Expected 1 active toast for critical notification, got %d", len(activeToasts))
	}
}

func TestNotificationService_ConfigurationAndPersistence(t *testing.T) {
	service := NewNotificationService(nil, nil)
	defer service.Stop()
	service.Start()

	// Test getting default config
	config := service.GetConfig()
	if config == nil {
		t.Error("Config should not be nil")
	}

	if !config.Enabled {
		t.Error("Default config should have notifications enabled")
	}

	// Test updating config
	newConfig := *config
	newConfig.ToastDuration = 10 * time.Second
	newConfig.MaxToastNotifications = 3

	err := service.UpdateConfig(&newConfig)
	if err != nil {
		t.Errorf("Failed to update config: %v", err)
	}

	// Verify config was updated
	updatedConfig := service.GetConfig()
	if updatedConfig.ToastDuration != 10*time.Second {
		t.Error("Toast duration was not updated")
	}

	if updatedConfig.MaxToastNotifications != 3 {
		t.Error("Max toast notifications was not updated")
	}
}

func TestNotificationService_Statistics(t *testing.T) {
	service := NewNotificationService(nil, nil)
	defer service.Stop()
	service.Start()

	// Create some test notifications
	notifications := []*ToastNotification{
		{Type: ToastTypeInfo, Title: "Info 1", Message: "Message 1", Category: CategoryGeneral, Priority: NotificationPriorityMedium},
		{Type: ToastTypeSuccess, Title: "Success 1", Message: "Message 2", Category: CategoryProject, Priority: NotificationPriorityMedium},
		{Type: ToastTypeError, Title: "Error 1", Message: "Message 3", Category: CategoryError, Priority: NotificationPriorityHigh},
		{Type: ToastTypeError, Title: "Error 2", Message: "Message 4", Category: CategoryError, Priority: NotificationPriorityHigh},
	}

	for _, notification := range notifications {
		err := service.ShowToast(notification)
		if err != nil {
			t.Fatalf("Failed to show notification: %v", err)
		}
	}

	// Get statistics
	stats := service.GetNotificationStats()
	if stats == nil {
		t.Error("Stats should not be nil")
	}

	if stats.TotalSent != int64(len(notifications)) {
		t.Errorf("Expected total sent %d, got %d", len(notifications), stats.TotalSent)
	}

	if stats.TotalToast != int64(len(notifications)) {
		t.Errorf("Expected total toast %d, got %d", len(notifications), stats.TotalToast)
	}

	// Check category stats
	if stats.ByCategory == nil {
		t.Error("Category stats should not be nil")
	}

	errorStats, exists := stats.ByCategory[CategoryError]
	if !exists {
		t.Error("Error category stats should exist")
	}

	if errorStats.Sent != 2 {
		t.Errorf("Expected 2 error notifications sent, got %d", errorStats.Sent)
	}

	// Check priority stats
	if stats.ByPriority == nil {
		t.Error("Priority stats should not be nil")
	}

	highPriorityStats, exists := stats.ByPriority[NotificationPriorityHigh]
	if !exists {
		t.Error("High priority stats should exist")
	}

	if highPriorityStats.Sent != 2 {
		t.Errorf("Expected 2 high priority notifications sent, got %d", highPriorityStats.Sent)
	}
}

func TestNotificationService_HistoryManagement(t *testing.T) {
	service := NewNotificationService(nil, nil)
	defer service.Stop()
	service.Start()

	// Create test notifications
	notifications := []*ToastNotification{
		{Type: ToastTypeInfo, Title: "Info 1", Message: "Message 1", Category: CategoryGeneral, Priority: NotificationPriorityMedium},
		{Type: ToastTypeSuccess, Title: "Success 1", Message: "Message 2", Category: CategoryProject, Priority: NotificationPriorityMedium},
		{Type: ToastTypeError, Title: "Error 1", Message: "Message 3", Category: CategoryError, Priority: NotificationPriorityHigh},
	}

	for _, notification := range notifications {
		err := service.ShowToast(notification)
		if err != nil {
			t.Fatalf("Failed to show notification: %v", err)
		}
	}

	// Test getting all history
	history := service.GetNotificationHistory(100, 0, "")
	if len(history) != len(notifications) {
		t.Errorf("Expected %d notifications in history, got %d", len(notifications), len(history))
	}

	// Test pagination
	firstPage := service.GetNotificationHistory(2, 0, "")
	if len(firstPage) != 2 {
		t.Errorf("Expected 2 notifications in first page, got %d", len(firstPage))
	}

	secondPage := service.GetNotificationHistory(2, 2, "")
	if len(secondPage) != 1 {
		t.Errorf("Expected 1 notification in second page, got %d", len(secondPage))
	}

	// Test category filtering
	errorHistory := service.GetNotificationHistory(100, 0, CategoryError)
	if len(errorHistory) != 1 {
		t.Errorf("Expected 1 error notification in filtered history, got %d", len(errorHistory))
	}

	if errorHistory[0].Category != CategoryError {
		t.Error("Filtered history should only contain error notifications")
	}

	// Test empty category filtering
	emptyHistory := service.GetNotificationHistory(100, 0, "nonexistent")
	if len(emptyHistory) != 0 {
		t.Errorf("Expected 0 notifications for non-existent category, got %d", len(emptyHistory))
	}
}

func TestNotificationService_JSONSerialization(t *testing.T) {
	// Test that all notification types can be properly serialized/deserialized
	toast := &ToastNotification{
		ID:       "test-id",
		Type:     ToastTypeInfo,
		Title:    "Test Toast",
		Message:  "Test message",
		Category: CategoryGeneral,
		Priority: NotificationPriorityMedium,
		Actions: []NotificationActionBtn{
			{ID: "action1", Label: "Action 1", Type: ActionTypeViewDetails},
		},
		CreatedAt: time.Now(),
		Metadata:  map[string]interface{}{"key": "value"},
	}

	// Test JSON marshaling
	jsonData, err := json.Marshal(toast)
	if err != nil {
		t.Errorf("Failed to marshal toast to JSON: %v", err)
	}

	// Test JSON unmarshaling
	var unmarshaledToast ToastNotification
	err = json.Unmarshal(jsonData, &unmarshaledToast)
	if err != nil {
		t.Errorf("Failed to unmarshal toast from JSON: %v", err)
	}

	// Verify data integrity
	if unmarshaledToast.ID != toast.ID {
		t.Error("ID mismatch after JSON roundtrip")
	}

	if unmarshaledToast.Title != toast.Title {
		t.Error("Title mismatch after JSON roundtrip")
	}

	if len(unmarshaledToast.Actions) != len(toast.Actions) {
		t.Error("Actions count mismatch after JSON roundtrip")
	}
}

// Helper function to create a test app with notification service
func createTestApp() *App {
	app := NewApp()
	app.notificationService = NewNotificationService(nil, nil)
	app.notificationService.Start()
	return app
}

func TestNotificationAPI_BasicOperations(t *testing.T) {
	app := createTestApp()
	defer app.notificationService.Stop()

	// Test showing toast via API
	toast := &ToastNotification{
		Type:     ToastTypeInfo,
		Title:    "API Test",
		Message:  "Testing API",
		Category: CategoryGeneral,
		Priority: NotificationPriorityMedium,
	}

	err := app.ShowToastNotification(toast)
	if err != nil {
		t.Errorf("Failed to show toast via API: %v", err)
	}

	// Test getting active toasts via API
	activeToasts := app.GetActiveToasts()
	if len(activeToasts) == 0 {
		t.Error("No active toasts returned from API")
	}

	// Test dismissing toast via API
	if len(activeToasts) > 0 {
		err = app.DismissToast(activeToasts[0].ID)
		if err != nil {
			t.Errorf("Failed to dismiss toast via API: %v", err)
		}
	}

	// Test getting history via API
	history := app.GetNotificationHistory(100, 0, "")
	if len(history) == 0 {
		t.Error("No history returned from API")
	}

	// Test getting stats via API
	stats := app.GetNotificationStats()
	if stats == nil {
		t.Error("No stats returned from API")
	}
}

func TestNotificationAPI_ConfigurationManagement(t *testing.T) {
	app := createTestApp()
	defer app.notificationService.Stop()

	// Test getting config via API
	config := app.GetNotificationConfig()
	if config == nil {
		t.Error("No config returned from API")
	}

	// Test updating config via API
	newConfig := *config
	newConfig.ToastDuration = 15 * time.Second

	err := app.UpdateNotificationConfig(&newConfig)
	if err != nil {
		t.Errorf("Failed to update config via API: %v", err)
	}

	// Verify config was updated
	updatedConfig := app.GetNotificationConfig()
	if updatedConfig.ToastDuration != 15*time.Second {
		t.Error("Config was not updated via API")
	}
}

func TestNotificationAPI_DoNotDisturbMode(t *testing.T) {
	app := createTestApp()
	defer app.notificationService.Stop()

	// Test enabling DND mode via API
	schedule := &DoNotDisturbSchedule{
		Enabled:     true,
		AllowUrgent: true,
	}

	err := app.EnableDoNotDisturbMode(schedule)
	if err != nil {
		t.Errorf("Failed to enable DND mode via API: %v", err)
	}

	// Verify DND mode is enabled
	config := app.GetNotificationConfig()
	if !config.DoNotDisturbMode {
		t.Error("DND mode was not enabled via API")
	}

	// Test disabling DND mode via API
	err = app.DisableDoNotDisturbMode()
	if err != nil {
		t.Errorf("Failed to disable DND mode via API: %v", err)
	}

	// Verify DND mode is disabled
	config = app.GetNotificationConfig()
	if config.DoNotDisturbMode {
		t.Error("DND mode was not disabled via API")
	}
}

func TestNotificationAPI_ErrorHandling(t *testing.T) {
	app := NewApp() // App without notification service

	// Test API methods with nil service
	err := app.ShowToastNotification(&ToastNotification{Title: "Test", Message: "Test"})
	if err == nil {
		t.Error("Expected error when showing toast with nil service")
	}

	activeToasts := app.GetActiveToasts()
	if len(activeToasts) != 0 {
		t.Error("Expected empty slice when getting active toasts with nil service")
	}

	history := app.GetNotificationHistory(100, 0, "")
	if len(history) != 0 {
		t.Error("Expected empty slice when getting history with nil service")
	}

	stats := app.GetNotificationStats()
	if stats.TotalSent != 0 {
		t.Error("Expected empty stats when getting stats with nil service")
	}

	// Test API validation
	app = createTestApp()
	defer app.notificationService.Stop()

	// Test nil toast
	err = app.ShowToastNotification(nil)
	if err == nil {
		t.Error("Expected error when showing nil toast")
	}

	// Test empty title
	err = app.ShowToastNotification(&ToastNotification{Message: "Test"})
	if err == nil {
		t.Error("Expected error when showing toast with empty title")
	}

	// Test empty message
	err = app.ShowToastNotification(&ToastNotification{Title: "Test"})
	if err == nil {
		t.Error("Expected error when showing toast with empty message")
	}
}
