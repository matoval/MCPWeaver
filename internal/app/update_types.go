package app

import (
	"crypto"
	"time"
)

// Auto-Update System Types

// UpdateInfo represents information about an available update
type UpdateInfo struct {
	Version      string            `json:"version"`
	ReleaseNotes string            `json:"releaseNotes"`
	DownloadURL  string            `json:"downloadUrl"`
	ChecksumURL  string            `json:"checksumUrl"`
	SignatureURL string            `json:"signatureUrl"`
	Size         int64             `json:"size"`
	PublishedAt  time.Time         `json:"publishedAt"`
	Critical     bool              `json:"critical"`
	Metadata     map[string]string `json:"metadata,omitempty"`
}

// UpdateStatus represents the current status of the update system
type UpdateStatus string

const (
	UpdateStatusIdle         UpdateStatus = "idle"
	UpdateStatusChecking     UpdateStatus = "checking"
	UpdateStatusAvailable    UpdateStatus = "available"
	UpdateStatusDownloading  UpdateStatus = "downloading"
	UpdateStatusVerifying    UpdateStatus = "verifying"
	UpdateStatusReady        UpdateStatus = "ready"
	UpdateStatusInstalling   UpdateStatus = "installing"
	UpdateStatusCompleted    UpdateStatus = "completed"
	UpdateStatusFailed       UpdateStatus = "failed"
	UpdateStatusRollingBack  UpdateStatus = "rolling_back"
	UpdateStatusRollbackComplete UpdateStatus = "rollback_complete"
)

// UpdateProgress tracks the progress of an update operation
type UpdateProgress struct {
	Status         UpdateStatus `json:"status"`
	Progress       float64      `json:"progress"`
	CurrentStep    string       `json:"currentStep"`
	BytesTotal     int64        `json:"bytesTotal"`
	BytesReceived  int64        `json:"bytesReceived"`
	Speed          int64        `json:"speed"` // bytes per second
	EstimatedTime  *time.Duration `json:"estimatedTime,omitempty"`
	Error          *APIError    `json:"error,omitempty"`
	LastUpdate     time.Time    `json:"lastUpdate"`
}

// UpdateSettings configures the auto-update behavior
type UpdateSettings struct {
	Enabled           bool              `json:"enabled"`
	AutoCheck         bool              `json:"autoCheck"`
	CheckInterval     time.Duration     `json:"checkInterval"`
	AutoDownload      bool              `json:"autoDownload"`
	AutoInstall       bool              `json:"autoInstall"`
	PromptUser        bool              `json:"promptUser"`
	Schedule          *UpdateSchedule   `json:"schedule,omitempty"`
	UpdateChannel     UpdateChannel     `json:"updateChannel"`
	PreReleaseEnabled bool              `json:"preReleaseEnabled"`
	BandwidthLimit    int64             `json:"bandwidthLimit"` // bytes per second, 0 = unlimited
	RetryPolicy       UpdateRetryPolicy `json:"retryPolicy"`
}

// UpdateSchedule defines when updates should be checked/installed
type UpdateSchedule struct {
	Type        ScheduleType `json:"type"`
	Time        string       `json:"time,omitempty"`        // HH:MM format
	DayOfWeek   int          `json:"dayOfWeek,omitempty"`   // 0-6, Sunday=0
	DayOfMonth  int          `json:"dayOfMonth,omitempty"`  // 1-31
	NextCheck   time.Time    `json:"nextCheck"`
}

// ScheduleType defines the type of update schedule
type ScheduleType string

const (
	ScheduleTypeImmediate ScheduleType = "immediate"
	ScheduleTypeDaily     ScheduleType = "daily"
	ScheduleTypeWeekly    ScheduleType = "weekly"
	ScheduleTypeMonthly   ScheduleType = "monthly"
	ScheduleTypeManual    ScheduleType = "manual"
)

// UpdateChannel defines the update channel to follow
type UpdateChannel string

const (
	UpdateChannelStable    UpdateChannel = "stable"
	UpdateChannelBeta      UpdateChannel = "beta"
	UpdateChannelAlpha     UpdateChannel = "alpha"
	UpdateChannelNightly   UpdateChannel = "nightly"
)

// UpdateRetryPolicy defines retry behavior for update operations
type UpdateRetryPolicy struct {
	MaxRetries        int           `json:"maxRetries"`
	InitialDelay      time.Duration `json:"initialDelay"`
	MaxDelay          time.Duration `json:"maxDelay"`
	BackoffMultiplier float64       `json:"backoffMultiplier"`
	RetryOnNetworkError bool        `json:"retryOnNetworkError"`
	RetryOnVerificationError bool   `json:"retryOnVerificationError"`
}

// UpdateResult represents the result of an update operation
type UpdateResult struct {
	Success        bool              `json:"success"`
	Version        string            `json:"version"`
	PreviousVersion string           `json:"previousVersion"`
	UpdatedAt      time.Time         `json:"updatedAt"`
	Duration       time.Duration     `json:"duration"`
	Error          *APIError         `json:"error,omitempty"`
	RollbackInfo   *RollbackInfo     `json:"rollbackInfo,omitempty"`
	VerificationResult *VerificationResult `json:"verificationResult,omitempty"`
}

// RollbackInfo contains information about rollback capabilities
type RollbackInfo struct {
	Available       bool      `json:"available"`
	BackupPath      string    `json:"backupPath"`
	BackupVersion   string    `json:"backupVersion"`
	BackupCreatedAt time.Time `json:"backupCreatedAt"`
	BackupSize      int64     `json:"backupSize"`
}

// VerificationResult contains the result of update verification
type VerificationResult struct {
	ChecksumValid    bool      `json:"checksumValid"`
	SignatureValid   bool      `json:"signatureValid"`
	CertificateValid bool      `json:"certificateValid"`
	Algorithm        string    `json:"algorithm"`
	VerifiedAt       time.Time `json:"verifiedAt"`
	TrustedSource    bool      `json:"trustedSource"`
}

// UpdateCheck represents a check for updates
type UpdateCheck struct {
	ID          string          `json:"id"`
	CheckedAt   time.Time       `json:"checkedAt"`
	Success     bool            `json:"success"`
	UpdateInfo  *UpdateInfo     `json:"updateInfo,omitempty"`
	Error       *APIError       `json:"error,omitempty"`
	Source      UpdateSource    `json:"source"`
	UserAgent   string          `json:"userAgent"`
	ClientInfo  ClientInfo      `json:"clientInfo"`
}

// UpdateSource defines where updates are checked from
type UpdateSource string

const (
	UpdateSourceGitHub     UpdateSource = "github"
	UpdateSourceCustom     UpdateSource = "custom"
	UpdateSourceEnterprise UpdateSource = "enterprise"
)

// ClientInfo contains information about the current client
type ClientInfo struct {
	Version     string `json:"version"`
	Platform    string `json:"platform"`
	Architecture string `json:"architecture"`
	OS          string `json:"os"`
	OSVersion   string `json:"osVersion"`
}

// UpdateConfiguration holds the update system configuration
type UpdateConfiguration struct {
	UpdateURL         string                 `json:"updateUrl"`
	PublicKey         string                 `json:"publicKey"`
	CertificatePath   string                 `json:"certificatePath"`
	BackupDirectory   string                 `json:"backupDirectory"`
	TempDirectory     string                 `json:"tempDirectory"`
	UserAgent         string                 `json:"userAgent"`
	Timeout           time.Duration          `json:"timeout"`
	VerificationMode  VerificationMode       `json:"verificationMode"`
	HashAlgorithm     crypto.Hash            `json:"hashAlgorithm"`
	DeltaUpdates      bool                   `json:"deltaUpdates"`
	CompressionLevel  int                    `json:"compressionLevel"`
	CustomHeaders     map[string]string      `json:"customHeaders,omitempty"`
}

// VerificationMode defines the level of verification required
type VerificationMode string

const (
	VerificationModeNone      VerificationMode = "none"
	VerificationModeChecksum  VerificationMode = "checksum"
	VerificationModeSignature VerificationMode = "signature"
	VerificationModeFull      VerificationMode = "full"
)

// UpdateNotification represents a notification about updates
type UpdateNotification struct {
	Type        NotificationType `json:"type"`
	Title       string           `json:"title"`
	Message     string           `json:"message"`
	Actions     []NotificationAction `json:"actions"`
	UpdateInfo  *UpdateInfo      `json:"updateInfo,omitempty"`
	Timestamp   time.Time        `json:"timestamp"`
	Persistent  bool             `json:"persistent"`
	Priority    NotificationPriority `json:"priority"`
}

// NotificationType defines the type of update notification
type NotificationType string

const (
	NotificationTypeUpdateAvailable NotificationType = "update_available"
	NotificationTypeUpdateReady     NotificationType = "update_ready"
	NotificationTypeUpdateFailed    NotificationType = "update_failed"
	NotificationTypeUpdateCompleted NotificationType = "update_completed"
	NotificationTypeCriticalUpdate  NotificationType = "critical_update"
	NotificationTypeRollbackNeeded  NotificationType = "rollback_needed"
)

// NotificationAction represents an action the user can take
type NotificationAction struct {
	ID    string `json:"id"`
	Label string `json:"label"`
	Type  ActionType `json:"type"`
}

// ActionType defines the type of notification action
type ActionType string

const (
	ActionTypeInstallNow   ActionType = "install_now"
	ActionTypeInstallLater ActionType = "install_later"
	ActionTypeSkipVersion  ActionType = "skip_version"
	ActionTypeViewDetails  ActionType = "view_details"
	ActionTypeDismiss      ActionType = "dismiss"
	ActionTypeRollback     ActionType = "rollback"
)

// NotificationPriority defines the priority of a notification
type NotificationPriority string

const (
	NotificationPriorityLow      NotificationPriority = "low"
	NotificationPriorityMedium   NotificationPriority = "medium"
	NotificationPriorityHigh     NotificationPriority = "high"
	NotificationPriorityCritical NotificationPriority = "critical"
)

// UpdateAnalytics tracks analytics for update operations
type UpdateAnalytics struct {
	UserID           string        `json:"userId,omitempty"`
	SessionID        string        `json:"sessionId"`
	EventType        AnalyticsEvent `json:"eventType"`
	Version          string        `json:"version"`
	PreviousVersion  string        `json:"previousVersion,omitempty"`
	UpdateChannel    UpdateChannel `json:"updateChannel"`
	Duration         time.Duration `json:"duration,omitempty"`
	Success          bool          `json:"success"`
	Error            string        `json:"error,omitempty"`
	ClientInfo       ClientInfo    `json:"clientInfo"`
	Timestamp        time.Time     `json:"timestamp"`
	Size             int64         `json:"size,omitempty"`
	DownloadSpeed    int64         `json:"downloadSpeed,omitempty"`
	UserAction       ActionType    `json:"userAction,omitempty"`
}

// AnalyticsEvent defines the type of analytics event
type AnalyticsEvent string

const (
	AnalyticsEventCheckStarted    AnalyticsEvent = "check_started"
	AnalyticsEventCheckCompleted  AnalyticsEvent = "check_completed"
	AnalyticsEventCheckFailed     AnalyticsEvent = "check_failed"
	AnalyticsEventDownloadStarted AnalyticsEvent = "download_started"
	AnalyticsEventDownloadCompleted AnalyticsEvent = "download_completed"
	AnalyticsEventDownloadFailed  AnalyticsEvent = "download_failed"
	AnalyticsEventInstallStarted  AnalyticsEvent = "install_started"
	AnalyticsEventInstallCompleted AnalyticsEvent = "install_completed"
	AnalyticsEventInstallFailed   AnalyticsEvent = "install_failed"
	AnalyticsEventUserAction      AnalyticsEvent = "user_action"
	AnalyticsEventRollbackStarted AnalyticsEvent = "rollback_started"
	AnalyticsEventRollbackCompleted AnalyticsEvent = "rollback_completed"
)

// DeltaUpdate represents information about a delta update
type DeltaUpdate struct {
	FromVersion string `json:"fromVersion"`
	ToVersion   string `json:"toVersion"`
	PatchURL    string `json:"patchUrl"`
	PatchSize   int64  `json:"patchSize"`
	ChecksumURL string `json:"checksumUrl"`
	Algorithm   string `json:"algorithm"`
}

// DefaultUpdateSettings returns default update settings
func DefaultUpdateSettings() UpdateSettings {
	return UpdateSettings{
		Enabled:           true,
		AutoCheck:         true,
		CheckInterval:     24 * time.Hour,
		AutoDownload:      false,
		AutoInstall:       false,
		PromptUser:        true,
		UpdateChannel:     UpdateChannelStable,
		PreReleaseEnabled: false,
		BandwidthLimit:    0, // unlimited
		RetryPolicy: UpdateRetryPolicy{
			MaxRetries:               3,
			InitialDelay:             5 * time.Second,
			MaxDelay:                 5 * time.Minute,
			BackoffMultiplier:        2.0,
			RetryOnNetworkError:      true,
			RetryOnVerificationError: false,
		},
	}
}

// DefaultUpdateConfiguration returns default update configuration
func DefaultUpdateConfiguration() UpdateConfiguration {
	return UpdateConfiguration{
		UpdateURL:         "https://api.github.com/repos/matoval/MCPWeaver/releases/latest",
		BackupDirectory:   "./backups",
		TempDirectory:     "./temp",
		UserAgent:         "MCPWeaver-AutoUpdater/1.0",
		Timeout:           30 * time.Second,
		VerificationMode:  VerificationModeChecksum,
		HashAlgorithm:     crypto.SHA256,
		DeltaUpdates:      false,
		CompressionLevel:  6,
	}
}