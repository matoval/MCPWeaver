package app

import (
	"time"
)

// Notification System Types

// NotificationSystem represents the notification system configuration
type NotificationSystem struct {
	Enabled              bool                      `json:"enabled"`
	ToastEnabled         bool                      `json:"toastEnabled"`
	SystemEnabled        bool                      `json:"systemEnabled"`
	SoundEnabled         bool                      `json:"soundEnabled"`
	DoNotDisturbMode     bool                      `json:"doNotDisturbMode"`
	DoNotDisturbSchedule *DoNotDisturbSchedule     `json:"doNotDisturbSchedule,omitempty"`
	MaxToastNotifications int                      `json:"maxToastNotifications"`
	ToastDuration        time.Duration             `json:"toastDuration"`
	HistoryRetention     time.Duration             `json:"historyRetention"`
	ThrottleSettings     *NotificationThrottle     `json:"throttleSettings"`
	Preferences          *NotificationPreferences  `json:"preferences"`
}

// ToastNotification represents a toast notification
type ToastNotification struct {
	ID           string                   `json:"id"`
	Type         ToastType                `json:"type"`
	Title        string                   `json:"title"`
	Message      string                   `json:"message"`
	Icon         string                   `json:"icon,omitempty"`
	Duration     time.Duration            `json:"duration"`
	Position     ToastPosition            `json:"position"`
	Actions      []NotificationActionBtn  `json:"actions,omitempty"`
	CreatedAt    time.Time                `json:"createdAt"`
	ExpiresAt    time.Time                `json:"expiresAt"`
	Persistent   bool                     `json:"persistent"`
	AutoDismiss  bool                     `json:"autoDismiss"`
	Priority     NotificationPriority     `json:"priority"`
	Category     NotificationCategory     `json:"category"`
	Metadata     map[string]interface{}   `json:"metadata,omitempty"`
	Progress     *NotificationProgress    `json:"progress,omitempty"`
}

// SystemNotification represents a system-level desktop notification
type SystemNotification struct {
	ID         string                   `json:"id"`
	Title      string                   `json:"title"`
	Body       string                   `json:"body"`
	Icon       string                   `json:"icon,omitempty"`
	Sound      string                   `json:"sound,omitempty"`
	Actions    []NotificationActionBtn  `json:"actions,omitempty"`
	Urgency    SystemUrgency            `json:"urgency"`
	Tag        string                   `json:"tag,omitempty"`
	CreatedAt  time.Time                `json:"createdAt"`
	Category   NotificationCategory     `json:"category"`
	Timeout    time.Duration            `json:"timeout"`
	Silent     bool                     `json:"silent"`
	Metadata   map[string]interface{}   `json:"metadata,omitempty"`
}

// NotificationHistory represents a stored notification for history
type NotificationHistory struct {
	ID            string                  `json:"id"`
	Type          string                  `json:"type"` // "toast" or "system"
	Title         string                  `json:"title"`
	Message       string                  `json:"message"`
	Icon          string                  `json:"icon,omitempty"`
	Actions       []NotificationActionBtn `json:"actions,omitempty"`
	Category      NotificationCategory    `json:"category"`
	Priority      NotificationPriority    `json:"priority"`
	CreatedAt     time.Time               `json:"createdAt"`
	ReadAt        *time.Time              `json:"readAt,omitempty"`
	DismissedAt   *time.Time              `json:"dismissedAt,omitempty"`
	InteractedAt  *time.Time              `json:"interactedAt,omitempty"`
	ActionTaken   string                  `json:"actionTaken,omitempty"`
	Source        string                  `json:"source"`
	Metadata      map[string]interface{}  `json:"metadata,omitempty"`
}

// NotificationActionBtn represents an action button on a notification
type NotificationActionBtn struct {
	ID       string     `json:"id"`
	Label    string     `json:"label"`
	Type     ActionType `json:"type"`
	Icon     string     `json:"icon,omitempty"`
	Style    ActionStyle `json:"style"`
	Callback string     `json:"callback,omitempty"`
	Data     map[string]interface{} `json:"data,omitempty"`
}

// NotificationProgress represents progress information for notifications
type NotificationProgress struct {
	Current int    `json:"current"`
	Total   int    `json:"total"`
	Percent int    `json:"percent"`
	Label   string `json:"label,omitempty"`
}

// NotificationPreferences represents user notification preferences
type NotificationPreferences struct {
	Categories map[NotificationCategory]CategoryPreference `json:"categories"`
	Filters    []NotificationFilter                        `json:"filters"`
	Sounds     map[ToastType]string                        `json:"sounds"`
	Volumes    map[ToastType]float64                       `json:"volumes"`
}

// CategoryPreference represents preferences for a notification category
type CategoryPreference struct {
	Enabled       bool                 `json:"enabled"`
	ToastEnabled  bool                 `json:"toastEnabled"`
	SystemEnabled bool                 `json:"systemEnabled"`
	SoundEnabled  bool                 `json:"soundEnabled"`
	MinPriority   NotificationPriority `json:"minPriority"`
	MaxPerHour    int                  `json:"maxPerHour"`
}

// NotificationFilter represents a filter for notifications
type NotificationFilter struct {
	ID        string               `json:"id"`
	Name      string               `json:"name"`
	Enabled   bool                 `json:"enabled"`
	Condition FilterCondition      `json:"condition"`
	Action    FilterAction         `json:"action"`
	Keywords  []string             `json:"keywords,omitempty"`
	Category  NotificationCategory `json:"category,omitempty"`
	Priority  NotificationPriority `json:"priority,omitempty"`
	Source    string               `json:"source,omitempty"`
}

// NotificationThrottle represents throttling settings
type NotificationThrottle struct {
	Enabled           bool          `json:"enabled"`
	MaxPerMinute      int           `json:"maxPerMinute"`
	MaxPerHour        int           `json:"maxPerHour"`
	BurstAllowance    int           `json:"burstAllowance"`
	CooldownPeriod    time.Duration `json:"cooldownPeriod"`
	ByCategory        map[NotificationCategory]ThrottleRule `json:"byCategory"`
	ByPriority        map[NotificationPriority]ThrottleRule `json:"byPriority"`
}

// ThrottleRule represents throttling rules for specific categories or priorities
type ThrottleRule struct {
	MaxPerMinute   int           `json:"maxPerMinute"`
	MaxPerHour     int           `json:"maxPerHour"`
	BurstAllowance int           `json:"burstAllowance"`
	CooldownPeriod time.Duration `json:"cooldownPeriod"`
}

// DoNotDisturbSchedule represents the schedule for do not disturb mode
type DoNotDisturbSchedule struct {
	Enabled     bool              `json:"enabled"`
	StartTime   string            `json:"startTime"`   // HH:MM format
	EndTime     string            `json:"endTime"`     // HH:MM format
	Days        []time.Weekday    `json:"days"`        // Days of week
	Exceptions  []time.Time       `json:"exceptions"`  // Specific dates to override
	AllowUrgent bool              `json:"allowUrgent"` // Allow critical/urgent notifications
}

// NotificationQueue represents a queue of pending notifications
type NotificationQueue struct {
	Notifications []QueuedNotification `json:"notifications"`
	MaxSize       int                  `json:"maxSize"`
	DrainRate     time.Duration        `json:"drainRate"`
	Paused        bool                 `json:"paused"`
}

// QueuedNotification represents a notification in the queue
type QueuedNotification struct {
	Notification interface{} `json:"notification"` // ToastNotification or SystemNotification
	QueuedAt     time.Time   `json:"queuedAt"`
	Priority     int         `json:"priority"`
	Attempts     int         `json:"attempts"`
	MaxAttempts  int         `json:"maxAttempts"`
}

// NotificationTemplate represents a reusable notification template
type NotificationTemplate struct {
	ID          string                 `json:"id"`
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	Type        string                 `json:"type"` // "toast" or "system"
	Title       string                 `json:"title"`
	Message     string                 `json:"message"`
	Icon        string                 `json:"icon,omitempty"`
	Actions     []NotificationActionBtn `json:"actions,omitempty"`
	Category    NotificationCategory   `json:"category"`
	Priority    NotificationPriority   `json:"priority"`
	Variables   []TemplateVariable     `json:"variables,omitempty"`
	CreatedAt   time.Time              `json:"createdAt"`
	UpdatedAt   time.Time              `json:"updatedAt"`
}

// NotificationStats represents statistics about notifications
type NotificationStats struct {
	TotalSent       int64                                    `json:"totalSent"`
	TotalToast      int64                                    `json:"totalToast"`
	TotalSystem     int64                                    `json:"totalSystem"`
	TotalRead       int64                                    `json:"totalRead"`
	TotalDismissed  int64                                    `json:"totalDismissed"`
	TotalInteracted int64                                    `json:"totalInteracted"`
	ByCategory      map[NotificationCategory]CategoryStats  `json:"byCategory"`
	ByPriority      map[NotificationPriority]PriorityStats  `json:"byPriority"`
	ByHour          map[int]int64                            `json:"byHour"`
	ByDay           map[time.Weekday]int64                   `json:"byDay"`
	PeriodStart     time.Time                                `json:"periodStart"`
	PeriodEnd       time.Time                                `json:"periodEnd"`
}

// CategoryStats represents statistics for a notification category
type CategoryStats struct {
	Sent       int64   `json:"sent"`
	Read       int64   `json:"read"`
	Dismissed  int64   `json:"dismissed"`
	Interacted int64   `json:"interacted"`
	ReadRate   float64 `json:"readRate"`
}

// PriorityStats represents statistics for a notification priority level
type PriorityStats struct {
	Sent       int64   `json:"sent"`
	Read       int64   `json:"read"`
	Dismissed  int64   `json:"dismissed"`
	Interacted int64   `json:"interacted"`
	ReadRate   float64 `json:"readRate"`
}

// Notification type enums

// ToastType defines the type of toast notification
type ToastType string

const (
	ToastTypeInfo    ToastType = "info"
	ToastTypeSuccess ToastType = "success"
	ToastTypeWarning ToastType = "warning"
	ToastTypeError   ToastType = "error"
	ToastTypeLoading ToastType = "loading"
	ToastTypeCustom  ToastType = "custom"
)

// ToastPosition defines where toast notifications appear
type ToastPosition string

const (
	ToastPositionTopLeft     ToastPosition = "top-left"
	ToastPositionTopCenter   ToastPosition = "top-center"
	ToastPositionTopRight    ToastPosition = "top-right"
	ToastPositionBottomLeft  ToastPosition = "bottom-left"
	ToastPositionBottomCenter ToastPosition = "bottom-center"
	ToastPositionBottomRight ToastPosition = "bottom-right"
)

// SystemUrgency defines the urgency level for system notifications
type SystemUrgency string

const (
	SystemUrgencyLow      SystemUrgency = "low"
	SystemUrgencyNormal   SystemUrgency = "normal"
	SystemUrgencyHigh     SystemUrgency = "high"
	SystemUrgencyCritical SystemUrgency = "critical"
)

// NotificationCategory defines categories for notifications
type NotificationCategory string

const (
	CategoryGeneral    NotificationCategory = "general"
	CategoryProject    NotificationCategory = "project"
	CategoryGeneration NotificationCategory = "generation"
	CategoryValidation NotificationCategory = "validation"
	CategoryUpdate     NotificationCategory = "update"
	CategoryFile       NotificationCategory = "file"
	CategoryError      NotificationCategory = "error"
	CategorySecurity   NotificationCategory = "security"
	CategorySystem     NotificationCategory = "system"
	CategoryTemplate   NotificationCategory = "template"
)

// ActionStyle defines the visual style of action buttons
type ActionStyle string

const (
	ActionStylePrimary   ActionStyle = "primary"
	ActionStyleSecondary ActionStyle = "secondary"
	ActionStyleDanger    ActionStyle = "danger"
	ActionStyleSuccess   ActionStyle = "success"
	ActionStyleOutline   ActionStyle = "outline"
	ActionStyleText      ActionStyle = "text"
)

// FilterCondition defines how notification filters are evaluated
type FilterCondition string

const (
	FilterConditionContains  FilterCondition = "contains"
	FilterConditionEquals    FilterCondition = "equals"
	FilterConditionStartsWith FilterCondition = "startsWith"
	FilterConditionEndsWith  FilterCondition = "endsWith"
	FilterConditionRegex     FilterCondition = "regex"
)

// FilterAction defines what happens when a filter matches
type FilterAction string

const (
	FilterActionBlock      FilterAction = "block"
	FilterActionAllow      FilterAction = "allow"
	FilterActionModify     FilterAction = "modify"
	FilterActionChangeType FilterAction = "changeType"
	FilterActionDelay      FilterAction = "delay"
)

// Default notification system settings
func DefaultNotificationSystem() NotificationSystem {
	return NotificationSystem{
		Enabled:               true,
		ToastEnabled:          true,
		SystemEnabled:         true,
		SoundEnabled:          true,
		DoNotDisturbMode:      false,
		MaxToastNotifications: 5,
		ToastDuration:         5 * time.Second,
		HistoryRetention:      30 * 24 * time.Hour, // 30 days
		ThrottleSettings: &NotificationThrottle{
			Enabled:        true,
			MaxPerMinute:   10,
			MaxPerHour:     100,
			BurstAllowance: 3,
			CooldownPeriod: time.Minute,
			ByCategory: map[NotificationCategory]ThrottleRule{
				CategoryError: {
					MaxPerMinute:   5,
					MaxPerHour:     30,
					BurstAllowance: 2,
					CooldownPeriod: 30 * time.Second,
				},
				CategorySystem: {
					MaxPerMinute:   2,
					MaxPerHour:     20,
					BurstAllowance: 1,
					CooldownPeriod: 2 * time.Minute,
				},
			},
			ByPriority: map[NotificationPriority]ThrottleRule{
				NotificationPriorityCritical: {
					MaxPerMinute:   20,
					MaxPerHour:     200,
					BurstAllowance: 5,
					CooldownPeriod: 30 * time.Second,
				},
			},
		},
		Preferences: &NotificationPreferences{
			Categories: map[NotificationCategory]CategoryPreference{
				CategoryGeneral: {
					Enabled:       true,
					ToastEnabled:  true,
					SystemEnabled: false,
					SoundEnabled:  false,
					MinPriority:   NotificationPriorityLow,
					MaxPerHour:    50,
				},
				CategoryError: {
					Enabled:       true,
					ToastEnabled:  true,
					SystemEnabled: true,
					SoundEnabled:  true,
					MinPriority:   NotificationPriorityMedium,
					MaxPerHour:    20,
				},
				CategoryUpdate: {
					Enabled:       true,
					ToastEnabled:  true,
					SystemEnabled: true,
					SoundEnabled:  false,
					MinPriority:   NotificationPriorityMedium,
					MaxPerHour:    10,
				},
			},
			Filters: []NotificationFilter{},
			Sounds: map[ToastType]string{
				ToastTypeInfo:    "notification.wav",
				ToastTypeSuccess: "success.wav",
				ToastTypeWarning: "warning.wav",
				ToastTypeError:   "error.wav",
			},
			Volumes: map[ToastType]float64{
				ToastTypeInfo:    0.5,
				ToastTypeSuccess: 0.6,
				ToastTypeWarning: 0.7,
				ToastTypeError:   0.8,
			},
		},
	}
}

// Helper methods

// IsExpired checks if a toast notification has expired
func (n *ToastNotification) IsExpired() bool {
	return !n.Persistent && time.Now().After(n.ExpiresAt)
}

// ShouldAutoDismiss checks if a toast notification should auto-dismiss
func (n *ToastNotification) ShouldAutoDismiss() bool {
	return n.AutoDismiss && !n.Persistent
}

// GetReadRate calculates the read rate for category stats
func (cs *CategoryStats) GetReadRate() float64 {
	if cs.Sent == 0 {
		return 0.0
	}
	return float64(cs.Read) / float64(cs.Sent) * 100.0
}

// GetReadRate calculates the read rate for priority stats
func (ps *PriorityStats) GetReadRate() float64 {
	if ps.Sent == 0 {
		return 0.0
	}
	return float64(ps.Read) / float64(ps.Sent) * 100.0
}

// IsInDoNotDisturbPeriod checks if current time is in do not disturb period
func (dnd *DoNotDisturbSchedule) IsInDoNotDisturbPeriod() bool {
	if !dnd.Enabled {
		return false
	}

	now := time.Now()
	
	// Check if current day is in the schedule
	currentDay := now.Weekday()
	dayIncluded := false
	for _, day := range dnd.Days {
		if day == currentDay {
			dayIncluded = true
			break
		}
	}
	
	if !dayIncluded {
		return false
	}

	// Parse start and end times
	startTime, err := time.Parse("15:04", dnd.StartTime)
	if err != nil {
		return false
	}
	
	endTime, err := time.Parse("15:04", dnd.EndTime)
	if err != nil {
		return false
	}

	// Get current time in same format
	currentTime := time.Date(0, 1, 1, now.Hour(), now.Minute(), 0, 0, time.UTC)
	startTimeOfDay := time.Date(0, 1, 1, startTime.Hour(), startTime.Minute(), 0, 0, time.UTC)
	endTimeOfDay := time.Date(0, 1, 1, endTime.Hour(), endTime.Minute(), 0, 0, time.UTC)

	// Handle overnight periods (e.g., 22:00 to 06:00)
	if endTimeOfDay.Before(startTimeOfDay) {
		return currentTime.After(startTimeOfDay) || currentTime.Before(endTimeOfDay)
	}
	
	// Normal periods (e.g., 09:00 to 17:00)
	return currentTime.After(startTimeOfDay) && currentTime.Before(endTimeOfDay)
}