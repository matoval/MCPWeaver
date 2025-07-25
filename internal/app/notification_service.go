package app

import (
	"context"
	"crypto/rand"
	"database/sql"
	"encoding/json"
	"fmt"
	"math/big"
	"sort"
	"strings"
	"sync"
	"time"

	wailsruntime "github.com/wailsapp/wails/v2/pkg/runtime"
)

// NotificationService manages all notification functionality
type NotificationService struct {
	ctx                 context.Context
	db                  *sql.DB
	config              *NotificationSystem
	activeToasts        map[string]*ToastNotification
	history             []NotificationHistory
	queue               *NotificationQueue
	throttleTracker     map[string][]time.Time // Key: category/priority, Value: timestamps
	filters             []NotificationFilter
	templates           map[string]*NotificationTemplate
	stats               *NotificationStats
	subscribers         []NotificationSubscriber
	isTestMode          bool
	mutex               sync.RWMutex
	historyMutex        sync.RWMutex
	queueMutex          sync.RWMutex
	ticker              *time.Ticker
	stopChan            chan bool
	queueProcessor      *time.Ticker
	historyCleanupTimer *time.Timer
}

// NotificationSubscriber defines the interface for notification event subscribers
type NotificationSubscriber interface {
	OnNotificationSent(notification interface{})
	OnNotificationRead(id string)
	OnNotificationDismissed(id string)
	OnNotificationInteracted(id string, action string)
}

// NewNotificationService creates a new notification service instance
func NewNotificationService(ctx context.Context, db *sql.DB) *NotificationService {
	config := DefaultNotificationSystem()

	// Detect test mode (when nil context provided)
	isTestMode := ctx == nil

	// Use background context if none provided (for testing)
	if ctx == nil {
		ctx = context.Background()
	}

	service := &NotificationService{
		ctx:             ctx,
		db:              db,
		config:          &config,
		activeToasts:    make(map[string]*ToastNotification),
		history:         make([]NotificationHistory, 0),
		throttleTracker: make(map[string][]time.Time),
		filters:         config.Preferences.Filters,
		templates:       make(map[string]*NotificationTemplate),
		stats:           &NotificationStats{},
		subscribers:     make([]NotificationSubscriber, 0),
		isTestMode:      isTestMode,
		stopChan:        make(chan bool),
		queue: &NotificationQueue{
			Notifications: make([]QueuedNotification, 0),
			MaxSize:       100,
			DrainRate:     time.Second,
			Paused:        false,
		},
	}

	// Initialize database tables if not in test mode
	if !isTestMode && db != nil {
		service.initializeTables()
	}

	// Load existing data from database
	service.loadFromDatabase()

	// Initialize default templates
	service.initializeDefaultTemplates()

	return service
}

// Start initializes and starts the notification service
func (ns *NotificationService) Start() error {
	ns.mutex.Lock()
	defer ns.mutex.Unlock()

	if !ns.config.Enabled {
		return nil
	}

	// Start queue processor
	ns.startQueueProcessor()

	// Start periodic cleanup
	ns.startPeriodicCleanup()

	// Emit service started event
	ns.emitEvent("notification:service_started", map[string]interface{}{
		"config": ns.config,
	})

	return nil
}

// Stop stops the notification service
func (ns *NotificationService) Stop() error {
	ns.mutex.Lock()
	defer ns.mutex.Unlock()

	// Stop timers
	if ns.ticker != nil {
		ns.ticker.Stop()
		ns.ticker = nil
	}

	if ns.queueProcessor != nil {
		ns.queueProcessor.Stop()
		ns.queueProcessor = nil
	}

	if ns.historyCleanupTimer != nil {
		ns.historyCleanupTimer.Stop()
		ns.historyCleanupTimer = nil
	}

	// Signal stop
	select {
	case ns.stopChan <- true:
	default:
	}

	// Save data to database
	ns.saveToDatabase()

	ns.emitEvent("notification:service_stopped", map[string]interface{}{
		"timestamp": time.Now(),
	})

	return nil
}

// ShowToast displays a toast notification
func (ns *NotificationService) ShowToast(toast *ToastNotification) error {
	ns.mutex.Lock()
	defer ns.mutex.Unlock()

	// Check if notifications are enabled
	if !ns.config.Enabled || !ns.config.ToastEnabled {
		return &APIError{
			Type:    ErrorTypeSystem,
			Code:    "NOTIFICATIONS_DISABLED",
			Message: "Toast notifications are disabled",
		}
	}

	// Check do not disturb mode
	if ns.isInDoNotDisturbMode(toast.Priority) {
		return ns.queueNotification(toast, "toast")
	}

	// Apply filters
	if !ns.passesFilters(toast) {
		return nil // Filtered out, not an error
	}

	// Check throttling
	if ns.isThrottled(toast.Category, toast.Priority) {
		return ns.queueNotification(toast, "toast")
	}

	// Generate ID if not provided
	if toast.ID == "" {
		toast.ID = ns.generateID()
	}

	// Set default values
	if toast.CreatedAt.IsZero() {
		toast.CreatedAt = time.Now()
	}

	if toast.Duration == 0 {
		toast.Duration = ns.config.ToastDuration
	}

	if !toast.Persistent {
		toast.ExpiresAt = toast.CreatedAt.Add(toast.Duration)
		toast.AutoDismiss = true
	}

	// Check maximum active toasts
	if len(ns.activeToasts) >= ns.config.MaxToastNotifications {
		// Remove oldest non-persistent toast
		ns.removeOldestToast()
	}

	// Add to active toasts
	ns.activeToasts[toast.ID] = toast

	// Add to history
	ns.addToHistory(toast, "toast")

	// Update statistics
	ns.updateStats(toast.Category, toast.Priority)

	// Track for throttling
	ns.trackNotification(toast.Category, toast.Priority)

	// Notify subscribers
	ns.notifySubscribers("OnNotificationSent", toast)

	// Emit to frontend
	ns.emitEvent("notification:toast", toast)

	// Schedule auto-dismiss if enabled
	if toast.AutoDismiss && !toast.Persistent {
		go ns.scheduleToastDismissal(toast.ID, toast.Duration)
	}

	return nil
}

// ShowSystemNotification displays a system-level desktop notification
func (ns *NotificationService) ShowSystemNotification(notification *SystemNotification) error {
	ns.mutex.Lock()
	defer ns.mutex.Unlock()

	// Check if notifications are enabled
	if !ns.config.Enabled || !ns.config.SystemEnabled {
		return &APIError{
			Type:    ErrorTypeSystem,
			Code:    "SYSTEM_NOTIFICATIONS_DISABLED",
			Message: "System notifications are disabled",
		}
	}

	// Convert priority for do not disturb check
	priority := ns.urgencyToPriority(notification.Urgency)

	// Check do not disturb mode
	if ns.isInDoNotDisturbMode(priority) {
		return ns.queueNotification(notification, "system")
	}

	// Apply filters (convert to toast-like structure for filtering)
	filterNotification := &ToastNotification{
		Title:    notification.Title,
		Message:  notification.Body,
		Category: notification.Category,
		Priority: priority,
	}
	if !ns.passesFilters(filterNotification) {
		return nil // Filtered out, not an error
	}

	// Check throttling
	if ns.isThrottled(notification.Category, priority) {
		return ns.queueNotification(notification, "system")
	}

	// Generate ID if not provided
	if notification.ID == "" {
		notification.ID = ns.generateID()
	}

	// Set default values
	if notification.CreatedAt.IsZero() {
		notification.CreatedAt = time.Now()
	}

	// Add to history
	ns.addToHistory(notification, "system")

	// Update statistics
	ns.updateStats(notification.Category, priority)

	// Track for throttling
	ns.trackNotification(notification.Category, priority)

	// Notify subscribers
	ns.notifySubscribers("OnNotificationSent", notification)

	// Send system notification via Wails
	if !ns.isTestMode {
		go func() {
			defer func() {
				if r := recover(); r != nil {
					// Silently handle Wails context-related panics during testing
				}
			}()

			// Use Wails runtime to show system notification
			wailsruntime.MessageDialog(ns.ctx, wailsruntime.MessageDialogOptions{
				Type:    wailsruntime.InfoDialog,
				Title:   notification.Title,
				Message: notification.Body,
			})
		}()
	}

	// Emit to frontend
	ns.emitEvent("notification:system", notification)

	return nil
}

// DismissToast dismisses a specific toast notification
func (ns *NotificationService) DismissToast(id string) error {
	ns.mutex.Lock()
	defer ns.mutex.Unlock()

	_, exists := ns.activeToasts[id]
	if !exists {
		return &APIError{
			Type:    ErrorTypeValidation,
			Code:    "TOAST_NOT_FOUND",
			Message: fmt.Sprintf("Toast notification with ID %s not found", id),
		}
	}

	// Remove from active toasts
	delete(ns.activeToasts, id)

	// Update history
	ns.updateHistoryDismissal(id)

	// Notify subscribers
	ns.notifySubscribers("OnNotificationDismissed", id)

	// Emit to frontend
	ns.emitEvent("notification:toast_dismissed", map[string]interface{}{
		"id":        id,
		"timestamp": time.Now(),
	})

	return nil
}

// DismissAllToasts dismisses all active toast notifications
func (ns *NotificationService) DismissAllToasts() error {
	ns.mutex.Lock()
	defer ns.mutex.Unlock()

	dismissedIDs := make([]string, 0, len(ns.activeToasts))

	for id := range ns.activeToasts {
		dismissedIDs = append(dismissedIDs, id)
		ns.updateHistoryDismissal(id)
		ns.notifySubscribers("OnNotificationDismissed", id)
	}

	// Clear all active toasts
	ns.activeToasts = make(map[string]*ToastNotification)

	// Emit to frontend
	ns.emitEvent("notification:all_toasts_dismissed", map[string]interface{}{
		"dismissedIDs": dismissedIDs,
		"timestamp":    time.Now(),
	})

	return nil
}

// MarkAsRead marks a notification as read
func (ns *NotificationService) MarkAsRead(id string) error {
	ns.historyMutex.Lock()
	defer ns.historyMutex.Unlock()

	for i := range ns.history {
		if ns.history[i].ID == id {
			if ns.history[i].ReadAt == nil {
				now := time.Now()
				ns.history[i].ReadAt = &now

				// Update statistics
				ns.stats.TotalRead++

				// Notify subscribers
				ns.notifySubscribers("OnNotificationRead", id)

				// Emit to frontend
				ns.emitEvent("notification:marked_read", map[string]interface{}{
					"id":        id,
					"timestamp": now,
				})
			}
			return nil
		}
	}

	return &APIError{
		Type:    ErrorTypeValidation,
		Code:    "NOTIFICATION_NOT_FOUND",
		Message: fmt.Sprintf("Notification with ID %s not found", id),
	}
}

// ExecuteAction executes a notification action
func (ns *NotificationService) ExecuteAction(notificationID, actionID string, data map[string]interface{}) error {
	ns.mutex.Lock()
	defer ns.mutex.Unlock()

	// Find the notification (check active toasts first, then history)
	var actions []NotificationActionBtn

	if toast, exists := ns.activeToasts[notificationID]; exists {
		actions = toast.Actions
	} else {
		// Search in history
		ns.historyMutex.RLock()
		for _, hist := range ns.history {
			if hist.ID == notificationID {
				actions = hist.Actions
				break
			}
		}
		ns.historyMutex.RUnlock()
	}

	if len(actions) == 0 {
		return &APIError{
			Type:    ErrorTypeValidation,
			Code:    "NOTIFICATION_NOT_FOUND",
			Message: fmt.Sprintf("Notification with ID %s not found or has no actions", notificationID),
		}
	}

	// Find the action
	var targetAction *NotificationActionBtn
	for _, action := range actions {
		if action.ID == actionID {
			targetAction = &action
			break
		}
	}

	if targetAction == nil {
		return &APIError{
			Type:    ErrorTypeValidation,
			Code:    "ACTION_NOT_FOUND",
			Message: fmt.Sprintf("Action with ID %s not found", actionID),
		}
	}

	// Update history with interaction
	ns.updateHistoryInteraction(notificationID, actionID)

	// Notify subscribers
	ns.notifySubscribers("OnNotificationInteracted", notificationID, actionID)

	// Execute the action callback if provided
	if targetAction.Callback != "" {
		ns.emitEvent("notification:action_executed", map[string]interface{}{
			"notificationId": notificationID,
			"actionId":       actionID,
			"callback":       targetAction.Callback,
			"data":           data,
			"timestamp":      time.Now(),
		})
	}

	// Handle built-in actions
	switch targetAction.Type {
	case ActionTypeDismiss:
		if _, exists := ns.activeToasts[notificationID]; exists {
			return ns.DismissToast(notificationID)
		}
	case ActionTypeInstallNow:
		// Emit action for update handling
		ns.emitEvent("notification:install_now", map[string]interface{}{
			"notificationId": notificationID,
			"data":           data,
		})
	case ActionTypeInstallLater:
		// Dismiss the notification and schedule reminder
		if _, exists := ns.activeToasts[notificationID]; exists {
			ns.DismissToast(notificationID)
		}
		ns.emitEvent("notification:install_later", map[string]interface{}{
			"notificationId": notificationID,
			"data":           data,
		})
	}

	return nil
}

// GetActiveToasts returns all currently active toast notifications
func (ns *NotificationService) GetActiveToasts() []*ToastNotification {
	ns.mutex.RLock()
	defer ns.mutex.RUnlock()

	toasts := make([]*ToastNotification, 0, len(ns.activeToasts))
	for _, toast := range ns.activeToasts {
		toasts = append(toasts, toast)
	}

	// Sort by creation time (newest first)
	sort.Slice(toasts, func(i, j int) bool {
		return toasts[i].CreatedAt.After(toasts[j].CreatedAt)
	})

	return toasts
}

// GetNotificationHistory returns notification history with optional filtering
func (ns *NotificationService) GetNotificationHistory(limit int, offset int, category NotificationCategory) []NotificationHistory {
	ns.historyMutex.RLock()
	defer ns.historyMutex.RUnlock()

	filtered := make([]NotificationHistory, 0)

	// Filter by category if specified
	for _, hist := range ns.history {
		if category == "" || hist.Category == category {
			filtered = append(filtered, hist)
		}
	}

	// Sort by creation time (newest first)
	sort.Slice(filtered, func(i, j int) bool {
		return filtered[i].CreatedAt.After(filtered[j].CreatedAt)
	})

	// Apply pagination
	if offset >= len(filtered) {
		return []NotificationHistory{}
	}

	end := offset + limit
	if end > len(filtered) {
		end = len(filtered)
	}

	return filtered[offset:end]
}

// GetNotificationStats returns notification statistics
func (ns *NotificationService) GetNotificationStats() *NotificationStats {
	ns.mutex.RLock()
	defer ns.mutex.RUnlock()

	// Create a copy to avoid race conditions
	statsCopy := *ns.stats
	return &statsCopy
}

// UpdateConfig updates the notification system configuration
func (ns *NotificationService) UpdateConfig(config *NotificationSystem) error {
	ns.mutex.Lock()
	defer ns.mutex.Unlock()

	ns.config = config
	ns.filters = config.Preferences.Filters

	// Restart services if needed
	if config.Enabled {
		ns.startQueueProcessor()
		ns.startPeriodicCleanup()
	}

	// Save to database
	ns.saveToDatabase()

	ns.emitEvent("notification:config_updated", config)

	return nil
}

// GetConfig returns the current notification system configuration
func (ns *NotificationService) GetConfig() *NotificationSystem {
	ns.mutex.RLock()
	defer ns.mutex.RUnlock()

	// Create a copy to avoid race conditions
	configCopy := *ns.config
	return &configCopy
}

// Private helper methods

func (ns *NotificationService) initializeTables() {
	// Create notification history table
	historyTable := `
	CREATE TABLE IF NOT EXISTS notification_history (
		id TEXT PRIMARY KEY,
		type TEXT NOT NULL,
		title TEXT NOT NULL,
		message TEXT NOT NULL,
		icon TEXT,
		actions TEXT, -- JSON
		category TEXT NOT NULL,
		priority TEXT NOT NULL,
		created_at DATETIME NOT NULL,
		read_at DATETIME,
		dismissed_at DATETIME,
		interacted_at DATETIME,
		action_taken TEXT,
		source TEXT,
		metadata TEXT -- JSON
	)`

	// Create notification templates table
	templatesTable := `
	CREATE TABLE IF NOT EXISTS notification_templates (
		id TEXT PRIMARY KEY,
		name TEXT NOT NULL,
		description TEXT,
		type TEXT NOT NULL,
		title TEXT NOT NULL,
		message TEXT NOT NULL,
		icon TEXT,
		actions TEXT, -- JSON
		category TEXT NOT NULL,
		priority TEXT NOT NULL,
		variables TEXT, -- JSON
		created_at DATETIME NOT NULL,
		updated_at DATETIME NOT NULL
	)`

	// Create notification config table
	configTable := `
	CREATE TABLE IF NOT EXISTS notification_config (
		id INTEGER PRIMARY KEY CHECK (id = 1),
		config TEXT NOT NULL -- JSON
	)`

	if ns.db != nil {
		ns.db.Exec(historyTable)
		ns.db.Exec(templatesTable)
		ns.db.Exec(configTable)
	}
}

func (ns *NotificationService) loadFromDatabase() {
	if ns.db == nil {
		return
	}

	// Load configuration
	var configJSON string
	err := ns.db.QueryRow("SELECT config FROM notification_config WHERE id = 1").Scan(&configJSON)
	if err == nil {
		var config NotificationSystem
		if json.Unmarshal([]byte(configJSON), &config) == nil {
			ns.config = &config
			ns.filters = config.Preferences.Filters
		}
	}

	// Load notification history
	rows, err := ns.db.Query("SELECT id, type, title, message, icon, actions, category, priority, created_at, read_at, dismissed_at, interacted_at, action_taken, source, metadata FROM notification_history ORDER BY created_at DESC LIMIT 1000")
	if err != nil {
		return
	}
	defer rows.Close()

	ns.history = make([]NotificationHistory, 0)
	for rows.Next() {
		var hist NotificationHistory
		var actionsJSON, metadataJSON sql.NullString
		var readAt, dismissedAt, interactedAt sql.NullTime

		err := rows.Scan(&hist.ID, &hist.Type, &hist.Title, &hist.Message, &hist.Icon, &actionsJSON, &hist.Category, &hist.Priority, &hist.CreatedAt, &readAt, &dismissedAt, &interactedAt, &hist.ActionTaken, &hist.Source, &metadataJSON)
		if err != nil {
			continue
		}

		if readAt.Valid {
			hist.ReadAt = &readAt.Time
		}
		if dismissedAt.Valid {
			hist.DismissedAt = &dismissedAt.Time
		}
		if interactedAt.Valid {
			hist.InteractedAt = &interactedAt.Time
		}

		if actionsJSON.Valid {
			json.Unmarshal([]byte(actionsJSON.String), &hist.Actions)
		}
		if metadataJSON.Valid {
			json.Unmarshal([]byte(metadataJSON.String), &hist.Metadata)
		}

		ns.history = append(ns.history, hist)
	}

	// Load templates
	templateRows, err := ns.db.Query("SELECT id, name, description, type, title, message, icon, actions, category, priority, variables, created_at, updated_at FROM notification_templates")
	if err != nil {
		return
	}
	defer templateRows.Close()

	for templateRows.Next() {
		var template NotificationTemplate
		var actionsJSON, variablesJSON sql.NullString

		err := templateRows.Scan(&template.ID, &template.Name, &template.Description, &template.Type, &template.Title, &template.Message, &template.Icon, &actionsJSON, &template.Category, &template.Priority, &variablesJSON, &template.CreatedAt, &template.UpdatedAt)
		if err != nil {
			continue
		}

		if actionsJSON.Valid {
			json.Unmarshal([]byte(actionsJSON.String), &template.Actions)
		}
		if variablesJSON.Valid {
			json.Unmarshal([]byte(variablesJSON.String), &template.Variables)
		}

		ns.templates[template.ID] = &template
	}
}

func (ns *NotificationService) saveToDatabase() {
	if ns.db == nil {
		return
	}

	// Save configuration
	configJSON, err := json.Marshal(ns.config)
	if err == nil {
		ns.db.Exec("INSERT OR REPLACE INTO notification_config (id, config) VALUES (1, ?)", string(configJSON))
	}

	// Note: History is saved as notifications are created/updated
	// Templates are saved when they are created/updated
}

func (ns *NotificationService) initializeDefaultTemplates() {
	// Project creation success template
	ns.templates["project_created"] = &NotificationTemplate{
		ID:          "project_created",
		Name:        "Project Created",
		Description: "Notification shown when a new project is created",
		Type:        "toast",
		Title:       "Project Created",
		Message:     "Project '{{projectName}}' has been created successfully",
		Icon:        "success",
		Category:    CategoryProject,
		Priority:    NotificationPriorityMedium,
		Variables: []TemplateVariable{
			{Name: "projectName", Type: "string", Required: true},
		},
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// Generation completed template
	ns.templates["generation_completed"] = &NotificationTemplate{
		ID:          "generation_completed",
		Name:        "Generation Completed",
		Description: "Notification shown when code generation is completed",
		Type:        "system",
		Title:       "Generation Completed",
		Message:     "MCP server generation completed for '{{projectName}}'",
		Icon:        "success",
		Category:    CategoryGeneration,
		Priority:    NotificationPriorityHigh,
		Actions: []NotificationActionBtn{
			{ID: "view_output", Label: "View Output", Type: ActionTypeViewDetails, Style: ActionStylePrimary},
			{ID: "dismiss", Label: "Dismiss", Type: ActionTypeDismiss, Style: ActionStyleSecondary},
		},
		Variables: []TemplateVariable{
			{Name: "projectName", Type: "string", Required: true},
		},
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// Error template
	ns.templates["error_occurred"] = &NotificationTemplate{
		ID:          "error_occurred",
		Name:        "Error Occurred",
		Description: "Notification shown when an error occurs",
		Type:        "toast",
		Title:       "Error",
		Message:     "{{errorMessage}}",
		Icon:        "error",
		Category:    CategoryError,
		Priority:    NotificationPriorityHigh,
		Variables: []TemplateVariable{
			{Name: "errorMessage", Type: "string", Required: true},
		},
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}

func (ns *NotificationService) isInDoNotDisturbMode(priority NotificationPriority) bool {
	if !ns.config.DoNotDisturbMode {
		return false
	}

	// Allow critical notifications if configured
	if priority == NotificationPriorityCritical {
		if ns.config.DoNotDisturbSchedule != nil && ns.config.DoNotDisturbSchedule.AllowUrgent {
			return false
		}
	}

	// Check schedule if configured
	if ns.config.DoNotDisturbSchedule != nil && ns.config.DoNotDisturbSchedule.Enabled {
		return ns.config.DoNotDisturbSchedule.IsInDoNotDisturbPeriod()
	}

	// If DND is enabled but no active schedule is configured, always in DND mode
	return true
}

func (ns *NotificationService) passesFilters(notification *ToastNotification) bool {
	for _, filter := range ns.filters {
		if !filter.Enabled {
			continue
		}

		matches := false

		// Check filter conditions
		switch filter.Condition {
		case FilterConditionContains:
			for _, keyword := range filter.Keywords {
				if strings.Contains(strings.ToLower(notification.Title), strings.ToLower(keyword)) ||
					strings.Contains(strings.ToLower(notification.Message), strings.ToLower(keyword)) {
					matches = true
					break
				}
			}
		case FilterConditionEquals:
			for _, keyword := range filter.Keywords {
				if strings.EqualFold(notification.Title, keyword) ||
					strings.EqualFold(notification.Message, keyword) {
					matches = true
					break
				}
			}
		}

		// Check category filter
		if filter.Category != "" && filter.Category == notification.Category {
			matches = true
		}

		// Check priority filter
		if filter.Priority != "" && filter.Priority == notification.Priority {
			matches = true
		}

		// Apply filter action if matches
		if matches {
			switch filter.Action {
			case FilterActionBlock:
				return false
			case FilterActionAllow:
				return true
			}
		}
	}

	return true // Default: allow all notifications
}

func (ns *NotificationService) isThrottled(category NotificationCategory, priority NotificationPriority) bool {
	if !ns.config.ThrottleSettings.Enabled {
		return false
	}

	now := time.Now()
	cutoffMinute := now.Add(-time.Minute)
	cutoffHour := now.Add(-time.Hour)

	// Check category-specific throttling
	if rule, exists := ns.config.ThrottleSettings.ByCategory[category]; exists {
		categoryKey := string(category)
		timestamps := ns.throttleTracker[categoryKey]

		// Count notifications in the last minute and hour
		minuteCount := 0
		hourCount := 0
		for _, ts := range timestamps {
			if ts.After(cutoffMinute) {
				minuteCount++
			}
			if ts.After(cutoffHour) {
				hourCount++
			}
		}

		if minuteCount >= rule.MaxPerMinute || hourCount >= rule.MaxPerHour {
			return true
		}
	}

	// Check priority-specific throttling
	if rule, exists := ns.config.ThrottleSettings.ByPriority[priority]; exists {
		priorityKey := string(priority)
		timestamps := ns.throttleTracker[priorityKey]

		minuteCount := 0
		hourCount := 0
		for _, ts := range timestamps {
			if ts.After(cutoffMinute) {
				minuteCount++
			}
			if ts.After(cutoffHour) {
				hourCount++
			}
		}

		if minuteCount >= rule.MaxPerMinute || hourCount >= rule.MaxPerHour {
			return true
		}
	}

	// Check global throttling
	totalMinuteCount := 0
	totalHourCount := 0
	for _, timestamps := range ns.throttleTracker {
		for _, ts := range timestamps {
			if ts.After(cutoffMinute) {
				totalMinuteCount++
			}
			if ts.After(cutoffHour) {
				totalHourCount++
			}
		}
	}

	return totalMinuteCount >= ns.config.ThrottleSettings.MaxPerMinute ||
		totalHourCount >= ns.config.ThrottleSettings.MaxPerHour
}

func (ns *NotificationService) trackNotification(category NotificationCategory, priority NotificationPriority) {
	now := time.Now()

	// Track by category
	categoryKey := string(category)
	ns.throttleTracker[categoryKey] = append(ns.throttleTracker[categoryKey], now)

	// Track by priority
	priorityKey := string(priority)
	ns.throttleTracker[priorityKey] = append(ns.throttleTracker[priorityKey], now)

	// Clean old timestamps (keep only last hour)
	cutoff := now.Add(-time.Hour)
	for key, timestamps := range ns.throttleTracker {
		filtered := make([]time.Time, 0)
		for _, ts := range timestamps {
			if ts.After(cutoff) {
				filtered = append(filtered, ts)
			}
		}
		ns.throttleTracker[key] = filtered
	}
}

func (ns *NotificationService) queueNotification(notification interface{}, notificationType string) error {
	ns.queueMutex.Lock()
	defer ns.queueMutex.Unlock()

	if len(ns.queue.Notifications) >= ns.queue.MaxSize {
		// Remove oldest notification
		ns.queue.Notifications = ns.queue.Notifications[1:]
	}

	priority := 0
	if notificationType == "toast" {
		if toast, ok := notification.(*ToastNotification); ok {
			priority = int(ns.priorityToInt(toast.Priority))
		}
	} else if notificationType == "system" {
		if sysNotif, ok := notification.(*SystemNotification); ok {
			priority = int(ns.priorityToInt(ns.urgencyToPriority(sysNotif.Urgency)))
		}
	}

	queued := QueuedNotification{
		Notification: notification,
		QueuedAt:     time.Now(),
		Priority:     priority,
		Attempts:     0,
		MaxAttempts:  3,
	}

	ns.queue.Notifications = append(ns.queue.Notifications, queued)

	// Sort by priority (higher priority first)
	sort.Slice(ns.queue.Notifications, func(i, j int) bool {
		return ns.queue.Notifications[i].Priority > ns.queue.Notifications[j].Priority
	})

	return nil
}

func (ns *NotificationService) startQueueProcessor() {
	if ns.queueProcessor != nil {
		ns.queueProcessor.Stop()
	}

	// Skip in test mode or if service is not properly initialized
	if ns.isTestMode || ns.stopChan == nil {
		return
	}

	ns.queueProcessor = time.NewTicker(ns.queue.DrainRate)
	go func() {
		defer func() {
			if r := recover(); r != nil {
				// Silently handle panics during shutdown
			}
		}()

		for {
			select {
			case <-ns.queueProcessor.C:
				ns.processQueue()
			case <-ns.stopChan:
				return
			}
		}
	}()
}

func (ns *NotificationService) processQueue() {
	ns.queueMutex.Lock()
	defer ns.queueMutex.Unlock()

	if ns.queue.Paused || len(ns.queue.Notifications) == 0 {
		return
	}

	// Process one notification from the queue
	queued := ns.queue.Notifications[0]
	ns.queue.Notifications = ns.queue.Notifications[1:]

	// Try to show the notification
	var err error
	if toast, ok := queued.Notification.(*ToastNotification); ok {
		err = ns.ShowToast(toast)
	} else if sysNotif, ok := queued.Notification.(*SystemNotification); ok {
		err = ns.ShowSystemNotification(sysNotif)
	}

	// If failed and haven't exceeded max attempts, re-queue
	if err != nil {
		queued.Attempts++
		if queued.Attempts < queued.MaxAttempts {
			ns.queue.Notifications = append(ns.queue.Notifications, queued)
		}
	}
}

func (ns *NotificationService) startPeriodicCleanup() {
	// Skip in test mode or if service is not properly initialized
	if ns.isTestMode || ns.stopChan == nil {
		return
	}

	// Clean up expired toasts every minute
	ns.ticker = time.NewTicker(time.Minute)
	go func() {
		defer func() {
			if r := recover(); r != nil {
				// Silently handle panics during shutdown
			}
		}()

		for {
			select {
			case <-ns.ticker.C:
				ns.cleanupExpiredToasts()
			case <-ns.stopChan:
				return
			}
		}
	}()

	// Clean up old history based on retention period
	if ns.historyCleanupTimer != nil {
		ns.historyCleanupTimer.Stop()
	}
	ns.historyCleanupTimer = time.AfterFunc(time.Hour, ns.cleanupOldHistory)
}

func (ns *NotificationService) cleanupExpiredToasts() {
	ns.mutex.Lock()
	defer ns.mutex.Unlock()

	now := time.Now()
	for id, toast := range ns.activeToasts {
		if toast.IsExpired() {
			delete(ns.activeToasts, id)
			ns.updateHistoryDismissal(id)
			ns.emitEvent("notification:toast_expired", map[string]interface{}{
				"id":        id,
				"timestamp": now,
			})
		}
	}
}

func (ns *NotificationService) cleanupOldHistory() {
	ns.historyMutex.Lock()
	defer ns.historyMutex.Unlock()

	cutoff := time.Now().Add(-ns.config.HistoryRetention)
	filtered := make([]NotificationHistory, 0)

	for _, hist := range ns.history {
		if hist.CreatedAt.After(cutoff) {
			filtered = append(filtered, hist)
		}
	}

	ns.history = filtered

	// Schedule next cleanup
	ns.historyCleanupTimer = time.AfterFunc(time.Hour, ns.cleanupOldHistory)
}

func (ns *NotificationService) scheduleToastDismissal(id string, duration time.Duration) {
	time.Sleep(duration)
	ns.DismissToast(id)
}

func (ns *NotificationService) removeOldestToast() {
	var oldestID string
	var oldestTime time.Time

	for id, toast := range ns.activeToasts {
		if !toast.Persistent && (oldestID == "" || toast.CreatedAt.Before(oldestTime)) {
			oldestID = id
			oldestTime = toast.CreatedAt
		}
	}

	if oldestID != "" {
		delete(ns.activeToasts, oldestID)
		ns.updateHistoryDismissal(oldestID)
		ns.emitEvent("notification:toast_auto_removed", map[string]interface{}{
			"id":        oldestID,
			"timestamp": time.Now(),
		})
	}
}

func (ns *NotificationService) addToHistory(notification interface{}, notificationType string) {
	ns.historyMutex.Lock()
	defer ns.historyMutex.Unlock()

	var history NotificationHistory

	if toast, ok := notification.(*ToastNotification); ok {
		history = NotificationHistory{
			ID:        toast.ID,
			Type:      "toast",
			Title:     toast.Title,
			Message:   toast.Message,
			Icon:      toast.Icon,
			Actions:   toast.Actions,
			Category:  toast.Category,
			Priority:  toast.Priority,
			CreatedAt: toast.CreatedAt,
			Source:    "toast",
			Metadata:  toast.Metadata,
		}
	} else if sysNotif, ok := notification.(*SystemNotification); ok {
		history = NotificationHistory{
			ID:        sysNotif.ID,
			Type:      "system",
			Title:     sysNotif.Title,
			Message:   sysNotif.Body,
			Icon:      sysNotif.Icon,
			Actions:   sysNotif.Actions,
			Category:  sysNotif.Category,
			Priority:  ns.urgencyToPriority(sysNotif.Urgency),
			CreatedAt: sysNotif.CreatedAt,
			Source:    "system",
			Metadata:  sysNotif.Metadata,
		}
	}

	ns.history = append(ns.history, history)

	// Save to database if available
	if ns.db != nil {
		go ns.saveHistoryToDatabase(history)
	}

	// Update statistics
	ns.stats.TotalSent++
	if notificationType == "toast" {
		ns.stats.TotalToast++
	} else {
		ns.stats.TotalSystem++
	}
}

func (ns *NotificationService) saveHistoryToDatabase(history NotificationHistory) {
	if ns.db == nil {
		return
	}

	actionsJSON, _ := json.Marshal(history.Actions)
	metadataJSON, _ := json.Marshal(history.Metadata)

	query := `
	INSERT INTO notification_history 
	(id, type, title, message, icon, actions, category, priority, created_at, read_at, dismissed_at, interacted_at, action_taken, source, metadata)
	VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

	ns.db.Exec(query, history.ID, history.Type, history.Title, history.Message, history.Icon,
		string(actionsJSON), history.Category, history.Priority, history.CreatedAt,
		history.ReadAt, history.DismissedAt, history.InteractedAt, history.ActionTaken,
		history.Source, string(metadataJSON))
}

func (ns *NotificationService) updateHistoryDismissal(id string) {
	ns.historyMutex.Lock()
	defer ns.historyMutex.Unlock()

	for i := range ns.history {
		if ns.history[i].ID == id && ns.history[i].DismissedAt == nil {
			now := time.Now()
			ns.history[i].DismissedAt = &now
			ns.stats.TotalDismissed++

			// Update database
			if ns.db != nil {
				go func() {
					ns.db.Exec("UPDATE notification_history SET dismissed_at = ? WHERE id = ?", now, id)
				}()
			}
			break
		}
	}
}

func (ns *NotificationService) updateHistoryInteraction(id, actionID string) {
	ns.historyMutex.Lock()
	defer ns.historyMutex.Unlock()

	for i := range ns.history {
		if ns.history[i].ID == id {
			now := time.Now()
			ns.history[i].InteractedAt = &now
			ns.history[i].ActionTaken = actionID
			ns.stats.TotalInteracted++

			// Update database
			if ns.db != nil {
				go func() {
					ns.db.Exec("UPDATE notification_history SET interacted_at = ?, action_taken = ? WHERE id = ?", now, actionID, id)
				}()
			}
			break
		}
	}
}

func (ns *NotificationService) updateStats(category NotificationCategory, priority NotificationPriority) {
	// Update category stats
	if ns.stats.ByCategory == nil {
		ns.stats.ByCategory = make(map[NotificationCategory]CategoryStats)
	}
	categoryStats := ns.stats.ByCategory[category]
	categoryStats.Sent++
	ns.stats.ByCategory[category] = categoryStats

	// Update priority stats
	if ns.stats.ByPriority == nil {
		ns.stats.ByPriority = make(map[NotificationPriority]PriorityStats)
	}
	priorityStats := ns.stats.ByPriority[priority]
	priorityStats.Sent++
	ns.stats.ByPriority[priority] = priorityStats

	// Update hourly stats
	if ns.stats.ByHour == nil {
		ns.stats.ByHour = make(map[int]int64)
	}
	hour := time.Now().Hour()
	ns.stats.ByHour[hour]++

	// Update daily stats
	if ns.stats.ByDay == nil {
		ns.stats.ByDay = make(map[time.Weekday]int64)
	}
	day := time.Now().Weekday()
	ns.stats.ByDay[day]++
}

func (ns *NotificationService) urgencyToPriority(urgency SystemUrgency) NotificationPriority {
	switch urgency {
	case SystemUrgencyLow:
		return NotificationPriorityLow
	case SystemUrgencyNormal:
		return NotificationPriorityMedium
	case SystemUrgencyHigh:
		return NotificationPriorityHigh
	case SystemUrgencyCritical:
		return NotificationPriorityCritical
	default:
		return NotificationPriorityMedium
	}
}

func (ns *NotificationService) priorityToInt(priority NotificationPriority) int {
	switch priority {
	case NotificationPriorityLow:
		return 1
	case NotificationPriorityMedium:
		return 2
	case NotificationPriorityHigh:
		return 3
	case NotificationPriorityCritical:
		return 4
	default:
		return 2
	}
}

func (ns *NotificationService) generateID() string {
	// Use crypto/rand for secure ID generation
	n, err := rand.Int(rand.Reader, big.NewInt(1000000000000))
	if err != nil {
		// Fallback to timestamp-based ID if crypto/rand fails
		return fmt.Sprintf("%d", time.Now().UnixNano())
	}
	return fmt.Sprintf("%d%d", time.Now().UnixNano(), n.Int64())
}

func (ns *NotificationService) emitEvent(eventName string, data interface{}) {
	// Skip events in test mode to avoid Wails runtime issues
	if ns.isTestMode {
		return
	}

	// Only emit events if we have a valid Wails context
	if ns.ctx != nil {
		defer func() {
			if r := recover(); r != nil {
				// Silently handle context-related panics during testing
			}
		}()
		wailsruntime.EventsEmit(ns.ctx, eventName, data)
	}
}

func (ns *NotificationService) notifySubscribers(method string, args ...interface{}) {
	for _, subscriber := range ns.subscribers {
		switch method {
		case "OnNotificationSent":
			if len(args) > 0 {
				subscriber.OnNotificationSent(args[0])
			}
		case "OnNotificationRead":
			if len(args) > 0 {
				if id, ok := args[0].(string); ok {
					subscriber.OnNotificationRead(id)
				}
			}
		case "OnNotificationDismissed":
			if len(args) > 0 {
				if id, ok := args[0].(string); ok {
					subscriber.OnNotificationDismissed(id)
				}
			}
		case "OnNotificationInteracted":
			if len(args) > 1 {
				if id, ok := args[0].(string); ok {
					if action, ok := args[1].(string); ok {
						subscriber.OnNotificationInteracted(id, action)
					}
				}
			}
		}
	}
}

// Subscribe adds a notification subscriber
func (ns *NotificationService) Subscribe(subscriber NotificationSubscriber) {
	ns.mutex.Lock()
	defer ns.mutex.Unlock()
	ns.subscribers = append(ns.subscribers, subscriber)
}

// Unsubscribe removes a notification subscriber
func (ns *NotificationService) Unsubscribe(subscriber NotificationSubscriber) {
	ns.mutex.Lock()
	defer ns.mutex.Unlock()
	for i, sub := range ns.subscribers {
		if sub == subscriber {
			ns.subscribers = append(ns.subscribers[:i], ns.subscribers[i+1:]...)
			break
		}
	}
}
