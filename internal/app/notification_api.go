package app

import (
	"fmt"
	"time"
)

// Notification API Methods

// ShowToastNotification displays a toast notification
func (a *App) ShowToastNotification(toast *ToastNotification) error {
	if a.notificationService == nil {
		return a.createAPIError(ErrorTypeSystem, "NOTIFICATION_SERVICE_NOT_AVAILABLE",
			"Notification service is not available", nil)
	}

	if toast == nil {
		return a.createAPIError(ErrorTypeValidation, "INVALID_TOAST",
			"Toast notification is required", nil)
	}

	// Validate required fields
	if toast.Title == "" {
		return a.createAPIError(ErrorTypeValidation, "MISSING_TITLE",
			"Toast title is required", nil)
	}

	if toast.Message == "" {
		return a.createAPIError(ErrorTypeValidation, "MISSING_MESSAGE",
			"Toast message is required", nil)
	}

	// Set defaults if not provided
	if toast.Type == "" {
		toast.Type = ToastTypeInfo
	}

	if toast.Position == "" {
		toast.Position = ToastPositionTopRight
	}

	if toast.Priority == "" {
		toast.Priority = NotificationPriorityMedium
	}

	if toast.Category == "" {
		toast.Category = CategoryGeneral
	}

	err := a.notificationService.ShowToast(toast)
	if err != nil {
		a.emitError(err.(*APIError))
		return err
	}

	return nil
}

// ShowSystemNotification displays a system-level desktop notification
func (a *App) ShowSystemNotification(notification *SystemNotification) error {
	if a.notificationService == nil {
		return a.createAPIError(ErrorTypeSystem, "NOTIFICATION_SERVICE_NOT_AVAILABLE",
			"Notification service is not available", nil)
	}

	if notification == nil {
		return a.createAPIError(ErrorTypeValidation, "INVALID_NOTIFICATION",
			"System notification is required", nil)
	}

	// Validate required fields
	if notification.Title == "" {
		return a.createAPIError(ErrorTypeValidation, "MISSING_TITLE",
			"Notification title is required", nil)
	}

	if notification.Body == "" {
		return a.createAPIError(ErrorTypeValidation, "MISSING_BODY",
			"Notification body is required", nil)
	}

	// Set defaults if not provided
	if notification.Urgency == "" {
		notification.Urgency = SystemUrgencyNormal
	}

	if notification.Category == "" {
		notification.Category = CategoryGeneral
	}

	if notification.Timeout == 0 {
		notification.Timeout = 5 * time.Second
	}

	err := a.notificationService.ShowSystemNotification(notification)
	if err != nil {
		a.emitError(err.(*APIError))
		return err
	}

	return nil
}

// DismissToast dismisses a specific toast notification
func (a *App) DismissToast(id string) error {
	if a.notificationService == nil {
		return a.createAPIError(ErrorTypeSystem, "NOTIFICATION_SERVICE_NOT_AVAILABLE",
			"Notification service is not available", nil)
	}

	if id == "" {
		return a.createAPIError(ErrorTypeValidation, "MISSING_ID",
			"Notification ID is required", nil)
	}

	err := a.notificationService.DismissToast(id)
	if err != nil {
		a.emitError(err.(*APIError))
		return err
	}

	return nil
}

// DismissAllToasts dismisses all active toast notifications
func (a *App) DismissAllToasts() error {
	if a.notificationService == nil {
		return a.createAPIError(ErrorTypeSystem, "NOTIFICATION_SERVICE_NOT_AVAILABLE",
			"Notification service is not available", nil)
	}

	err := a.notificationService.DismissAllToasts()
	if err != nil {
		a.emitError(err.(*APIError))
		return err
	}

	return nil
}

// GetActiveToasts returns all currently active toast notifications
func (a *App) GetActiveToasts() []*ToastNotification {
	if a.notificationService == nil {
		return []*ToastNotification{}
	}

	return a.notificationService.GetActiveToasts()
}

// MarkNotificationAsRead marks a notification as read
func (a *App) MarkNotificationAsRead(id string) error {
	if a.notificationService == nil {
		return a.createAPIError(ErrorTypeSystem, "NOTIFICATION_SERVICE_NOT_AVAILABLE",
			"Notification service is not available", nil)
	}

	if id == "" {
		return a.createAPIError(ErrorTypeValidation, "MISSING_ID",
			"Notification ID is required", nil)
	}

	err := a.notificationService.MarkAsRead(id)
	if err != nil {
		a.emitError(err.(*APIError))
		return err
	}

	return nil
}

// ExecuteNotificationAction executes a notification action
func (a *App) ExecuteNotificationAction(notificationID, actionID string, data map[string]interface{}) error {
	if a.notificationService == nil {
		return a.createAPIError(ErrorTypeSystem, "NOTIFICATION_SERVICE_NOT_AVAILABLE",
			"Notification service is not available", nil)
	}

	if notificationID == "" {
		return a.createAPIError(ErrorTypeValidation, "MISSING_NOTIFICATION_ID",
			"Notification ID is required", nil)
	}

	if actionID == "" {
		return a.createAPIError(ErrorTypeValidation, "MISSING_ACTION_ID",
			"Action ID is required", nil)
	}

	err := a.notificationService.ExecuteAction(notificationID, actionID, data)
	if err != nil {
		a.emitError(err.(*APIError))
		return err
	}

	return nil
}

// GetNotificationHistory returns notification history with optional filtering
func (a *App) GetNotificationHistory(limit int, offset int, category string) []NotificationHistory {
	if a.notificationService == nil {
		return []NotificationHistory{}
	}

	// Validate and set defaults
	if limit <= 0 || limit > 1000 {
		limit = 100
	}

	if offset < 0 {
		offset = 0
	}

	var filterCategory NotificationCategory
	if category != "" {
		filterCategory = NotificationCategory(category)
	}

	return a.notificationService.GetNotificationHistory(limit, offset, filterCategory)
}

// GetNotificationStats returns notification statistics
func (a *App) GetNotificationStats() *NotificationStats {
	if a.notificationService == nil {
		return &NotificationStats{}
	}

	return a.notificationService.GetNotificationStats()
}

// GetNotificationConfig returns the current notification system configuration
func (a *App) GetNotificationConfig() *NotificationSystem {
	if a.notificationService == nil {
		defaultConfig := DefaultNotificationSystem()
		return &defaultConfig
	}

	return a.notificationService.GetConfig()
}

// UpdateNotificationConfig updates the notification system configuration
func (a *App) UpdateNotificationConfig(config *NotificationSystem) error {
	if a.notificationService == nil {
		return a.createAPIError(ErrorTypeSystem, "NOTIFICATION_SERVICE_NOT_AVAILABLE",
			"Notification service is not available", nil)
	}

	if config == nil {
		return a.createAPIError(ErrorTypeValidation, "INVALID_CONFIG",
			"Configuration is required", nil)
	}

	err := a.notificationService.UpdateConfig(config)
	if err != nil {
		apiErr := a.createAPIError(ErrorTypeSystem, "CONFIG_UPDATE_FAILED",
			fmt.Sprintf("Failed to update configuration: %v", err), nil)
		a.emitError(apiErr)
		return apiErr
	}

	a.emitNotification("info", "Configuration Updated",
		"Notification settings have been successfully updated")

	return nil
}

// EnableDoNotDisturbMode enables do not disturb mode
func (a *App) EnableDoNotDisturbMode(schedule *DoNotDisturbSchedule) error {
	if a.notificationService == nil {
		return a.createAPIError(ErrorTypeSystem, "NOTIFICATION_SERVICE_NOT_AVAILABLE",
			"Notification service is not available", nil)
	}

	config := a.notificationService.GetConfig()
	config.DoNotDisturbMode = true
	if schedule != nil {
		config.DoNotDisturbSchedule = schedule
	}

	err := a.notificationService.UpdateConfig(config)
	if err != nil {
		apiErr := a.createAPIError(ErrorTypeSystem, "DND_ENABLE_FAILED",
			fmt.Sprintf("Failed to enable do not disturb mode: %v", err), nil)
		a.emitError(apiErr)
		return apiErr
	}

	a.emitNotification("info", "Do Not Disturb Enabled",
		"Do not disturb mode has been enabled")

	return nil
}

// DisableDoNotDisturbMode disables do not disturb mode
func (a *App) DisableDoNotDisturbMode() error {
	if a.notificationService == nil {
		return a.createAPIError(ErrorTypeSystem, "NOTIFICATION_SERVICE_NOT_AVAILABLE",
			"Notification service is not available", nil)
	}

	config := a.notificationService.GetConfig()
	config.DoNotDisturbMode = false

	err := a.notificationService.UpdateConfig(config)
	if err != nil {
		apiErr := a.createAPIError(ErrorTypeSystem, "DND_DISABLE_FAILED",
			fmt.Sprintf("Failed to disable do not disturb mode: %v", err), nil)
		a.emitError(apiErr)
		return apiErr
	}

	a.emitNotification("info", "Do Not Disturb Disabled",
		"Do not disturb mode has been disabled")

	return nil
}

// AddNotificationFilter adds a new notification filter
func (a *App) AddNotificationFilter(filter *NotificationFilter) error {
	if a.notificationService == nil {
		return a.createAPIError(ErrorTypeSystem, "NOTIFICATION_SERVICE_NOT_AVAILABLE",
			"Notification service is not available", nil)
	}

	if filter == nil {
		return a.createAPIError(ErrorTypeValidation, "INVALID_FILTER",
			"Filter is required", nil)
	}

	// Validate filter
	if filter.Name == "" {
		return a.createAPIError(ErrorTypeValidation, "MISSING_FILTER_NAME",
			"Filter name is required", nil)
	}

	if filter.Condition == "" {
		return a.createAPIError(ErrorTypeValidation, "MISSING_CONDITION",
			"Filter condition is required", nil)
	}

	if filter.Action == "" {
		return a.createAPIError(ErrorTypeValidation, "MISSING_ACTION",
			"Filter action is required", nil)
	}

	// Generate ID if not provided
	if filter.ID == "" {
		filter.ID = fmt.Sprintf("filter_%d", time.Now().UnixNano())
	}

	// Get current config and add filter
	config := a.notificationService.GetConfig()
	config.Preferences.Filters = append(config.Preferences.Filters, *filter)

	err := a.notificationService.UpdateConfig(config)
	if err != nil {
		apiErr := a.createAPIError(ErrorTypeSystem, "FILTER_ADD_FAILED",
			fmt.Sprintf("Failed to add filter: %v", err), nil)
		a.emitError(apiErr)
		return apiErr
	}

	a.emitNotification("info", "Filter Added",
		fmt.Sprintf("Notification filter '%s' has been added", filter.Name))

	return nil
}

// RemoveNotificationFilter removes a notification filter
func (a *App) RemoveNotificationFilter(filterID string) error {
	if a.notificationService == nil {
		return a.createAPIError(ErrorTypeSystem, "NOTIFICATION_SERVICE_NOT_AVAILABLE",
			"Notification service is not available", nil)
	}

	if filterID == "" {
		return a.createAPIError(ErrorTypeValidation, "MISSING_FILTER_ID",
			"Filter ID is required", nil)
	}

	// Get current config and remove filter
	config := a.notificationService.GetConfig()
	originalCount := len(config.Preferences.Filters)

	filtered := make([]NotificationFilter, 0)
	for _, filter := range config.Preferences.Filters {
		if filter.ID != filterID {
			filtered = append(filtered, filter)
		}
	}

	if len(filtered) == originalCount {
		return a.createAPIError(ErrorTypeValidation, "FILTER_NOT_FOUND",
			fmt.Sprintf("Filter with ID %s not found", filterID), nil)
	}

	config.Preferences.Filters = filtered

	err := a.notificationService.UpdateConfig(config)
	if err != nil {
		apiErr := a.createAPIError(ErrorTypeSystem, "FILTER_REMOVE_FAILED",
			fmt.Sprintf("Failed to remove filter: %v", err), nil)
		a.emitError(apiErr)
		return apiErr
	}

	a.emitNotification("info", "Filter Removed", "Notification filter has been removed")

	return nil
}

// UpdateNotificationFilter updates an existing notification filter
func (a *App) UpdateNotificationFilter(filter *NotificationFilter) error {
	if a.notificationService == nil {
		return a.createAPIError(ErrorTypeSystem, "NOTIFICATION_SERVICE_NOT_AVAILABLE",
			"Notification service is not available", nil)
	}

	if filter == nil || filter.ID == "" {
		return a.createAPIError(ErrorTypeValidation, "INVALID_FILTER",
			"Filter with valid ID is required", nil)
	}

	// Get current config and update filter
	config := a.notificationService.GetConfig()
	found := false

	for i, existingFilter := range config.Preferences.Filters {
		if existingFilter.ID == filter.ID {
			config.Preferences.Filters[i] = *filter
			found = true
			break
		}
	}

	if !found {
		return a.createAPIError(ErrorTypeValidation, "FILTER_NOT_FOUND",
			fmt.Sprintf("Filter with ID %s not found", filter.ID), nil)
	}

	err := a.notificationService.UpdateConfig(config)
	if err != nil {
		apiErr := a.createAPIError(ErrorTypeSystem, "FILTER_UPDATE_FAILED",
			fmt.Sprintf("Failed to update filter: %v", err), nil)
		a.emitError(apiErr)
		return apiErr
	}

	a.emitNotification("info", "Filter Updated",
		fmt.Sprintf("Notification filter '%s' has been updated", filter.Name))

	return nil
}

// TestNotificationFilter tests a notification filter against sample data
func (a *App) TestNotificationFilter(filter *NotificationFilter, sampleNotification *ToastNotification) (*NotificationFilterTestResult, error) {
	if filter == nil {
		return nil, a.createAPIError(ErrorTypeValidation, "INVALID_FILTER",
			"Filter is required", nil)
	}

	if sampleNotification == nil {
		return nil, a.createAPIError(ErrorTypeValidation, "INVALID_SAMPLE",
			"Sample notification is required", nil)
	}

	// Create a temporary notification service for testing
	tempService := NewNotificationService(nil, nil)
	tempConfig := DefaultNotificationSystem()
	tempConfig.Preferences.Filters = []NotificationFilter{*filter}
	tempService.UpdateConfig(&tempConfig)

	// Test the filter
	passes := tempService.passesFilters(sampleNotification)

	result := &NotificationFilterTestResult{
		FilterID:   filter.ID,
		FilterName: filter.Name,
		Passes:     passes,
		Action:     filter.Action,
		TestResult: "Filter applied successfully",
		TestedAt:   time.Now(),
	}

	if !passes && filter.Action == FilterActionBlock {
		result.TestResult = "Notification would be blocked by this filter"
	} else if passes && filter.Action == FilterActionAllow {
		result.TestResult = "Notification would be allowed by this filter"
	}

	return result, nil
}

// ClearNotificationHistory clears all notification history
func (a *App) ClearNotificationHistory() error {
	if a.notificationService == nil {
		return a.createAPIError(ErrorTypeSystem, "NOTIFICATION_SERVICE_NOT_AVAILABLE",
			"Notification service is not available", nil)
	}

	// Clear history by updating with empty slice
	a.notificationService.historyMutex.Lock()
	a.notificationService.history = make([]NotificationHistory, 0)
	a.notificationService.historyMutex.Unlock()

	// Clear from database if available
	if a.notificationService.db != nil {
		go func() {
			a.notificationService.db.Exec("DELETE FROM notification_history")
		}()
	}

	a.emitNotification("info", "History Cleared", "Notification history has been cleared")

	return nil
}

// ExportNotificationHistory exports notification history to a file
func (a *App) ExportNotificationHistory(filePath string, format string) error {
	if a.notificationService == nil {
		return a.createAPIError(ErrorTypeSystem, "NOTIFICATION_SERVICE_NOT_AVAILABLE",
			"Notification service is not available", nil)
	}

	if filePath == "" {
		return a.createAPIError(ErrorTypeValidation, "MISSING_FILE_PATH",
			"File path is required", nil)
	}

	if format == "" {
		format = "json"
	}

	// TODO: Implement actual export functionality
	// For now, just emit an event
	a.emitNotification("info", "Export Started",
		fmt.Sprintf("Exporting notification history to %s", filePath))

	return nil
}

// PauseNotifications pauses all notifications temporarily
func (a *App) PauseNotifications() error {
	if a.notificationService == nil {
		return a.createAPIError(ErrorTypeSystem, "NOTIFICATION_SERVICE_NOT_AVAILABLE",
			"Notification service is not available", nil)
	}

	a.notificationService.queueMutex.Lock()
	a.notificationService.queue.Paused = true
	a.notificationService.queueMutex.Unlock()

	a.emitNotification("info", "Notifications Paused",
		"All notifications have been paused")

	return nil
}

// ResumeNotifications resumes paused notifications
func (a *App) ResumeNotifications() error {
	if a.notificationService == nil {
		return a.createAPIError(ErrorTypeSystem, "NOTIFICATION_SERVICE_NOT_AVAILABLE",
			"Notification service is not available", nil)
	}

	a.notificationService.queueMutex.Lock()
	a.notificationService.queue.Paused = false
	a.notificationService.queueMutex.Unlock()

	a.emitNotification("info", "Notifications Resumed",
		"Notifications have been resumed")

	return nil
}

// GetNotificationQueue returns the current notification queue status
func (a *App) GetNotificationQueue() *NotificationQueueStatus {
	if a.notificationService == nil {
		return &NotificationQueueStatus{
			Size:   0,
			Paused: false,
		}
	}

	a.notificationService.queueMutex.RLock()
	defer a.notificationService.queueMutex.RUnlock()

	return &NotificationQueueStatus{
		Size:        len(a.notificationService.queue.Notifications),
		MaxSize:     a.notificationService.queue.MaxSize,
		Paused:      a.notificationService.queue.Paused,
		DrainRate:   a.notificationService.queue.DrainRate,
		QueuedItems: len(a.notificationService.queue.Notifications),
	}
}

// ClearNotificationQueue clears all queued notifications
func (a *App) ClearNotificationQueue() error {
	if a.notificationService == nil {
		return a.createAPIError(ErrorTypeSystem, "NOTIFICATION_SERVICE_NOT_AVAILABLE",
			"Notification service is not available", nil)
	}

	a.notificationService.queueMutex.Lock()
	a.notificationService.queue.Notifications = make([]QueuedNotification, 0)
	a.notificationService.queueMutex.Unlock()

	a.emitNotification("info", "Queue Cleared", "Notification queue has been cleared")

	return nil
}

// Helper method to show quick notification messages
func (a *App) ShowInfoToast(title, message string) error {
	return a.ShowToastNotification(&ToastNotification{
		Type:     ToastTypeInfo,
		Title:    title,
		Message:  message,
		Category: CategoryGeneral,
		Priority: NotificationPriorityMedium,
		Position: ToastPositionTopRight,
	})
}

func (a *App) ShowSuccessToast(title, message string) error {
	return a.ShowToastNotification(&ToastNotification{
		Type:     ToastTypeSuccess,
		Title:    title,
		Message:  message,
		Category: CategoryGeneral,
		Priority: NotificationPriorityMedium,
		Position: ToastPositionTopRight,
	})
}

func (a *App) ShowWarningToast(title, message string) error {
	return a.ShowToastNotification(&ToastNotification{
		Type:     ToastTypeWarning,
		Title:    title,
		Message:  message,
		Category: CategoryGeneral,
		Priority: NotificationPriorityHigh,
		Position: ToastPositionTopRight,
	})
}

func (a *App) ShowErrorToast(title, message string) error {
	return a.ShowToastNotification(&ToastNotification{
		Type:     ToastTypeError,
		Title:    title,
		Message:  message,
		Category: CategoryError,
		Priority: NotificationPriorityHigh,
		Position: ToastPositionTopRight,
	})
}

// Supporting types for API methods

// NotificationFilterTestResult represents the result of testing a filter
type NotificationFilterTestResult struct {
	FilterID   string       `json:"filterId"`
	FilterName string       `json:"filterName"`
	Passes     bool         `json:"passes"`
	Action     FilterAction `json:"action"`
	TestResult string       `json:"testResult"`
	TestedAt   time.Time    `json:"testedAt"`
}

// NotificationQueueStatus represents the status of the notification queue
type NotificationQueueStatus struct {
	Size        int           `json:"size"`
	MaxSize     int           `json:"maxSize"`
	Paused      bool          `json:"paused"`
	DrainRate   time.Duration `json:"drainRate"`
	QueuedItems int           `json:"queuedItems"`
}
