package app

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"sync"
	"time"

	wailsruntime "github.com/wailsapp/wails/v2/pkg/runtime"
)

// UpdateService manages the auto-update functionality
type UpdateService struct {
	ctx             context.Context
	config          *UpdateConfiguration
	settings        *UpdateSettings
	currentStatus   UpdateStatus
	progress        *UpdateProgress
	lastCheck       *UpdateCheck
	rollbackManager *RollbackManager
	scheduler       *UpdateScheduler
	isTestMode      bool
	mutex           sync.RWMutex
	httpClient      *http.Client
	analytics       []UpdateAnalytics
	subscribers     []UpdateSubscriber
	ticker          *time.Ticker
	stopChan        chan bool
}

// UpdateSubscriber defines the interface for update event subscribers
type UpdateSubscriber interface {
	OnUpdateAvailable(info *UpdateInfo)
	OnUpdateProgress(progress *UpdateProgress)
	OnUpdateCompleted(result *UpdateResult)
	OnUpdateFailed(err *APIError)
}

// NewUpdateService creates a new update service instance
func NewUpdateService(ctx context.Context) *UpdateService {
	config := DefaultUpdateConfiguration()
	settings := DefaultUpdateSettings()

	// Detect test mode (when nil context provided)
	isTestMode := ctx == nil

	// Use background context if none provided (for testing)
	if ctx == nil {
		ctx = context.Background()
	}

	service := &UpdateService{
		ctx:             ctx,
		config:          &config,
		settings:        &settings,
		currentStatus:   UpdateStatusIdle,
		rollbackManager: NewRollbackManager(config.BackupDirectory),
		scheduler:       nil, // Will be initialized below
		isTestMode:      isTestMode,
		progress: &UpdateProgress{
			Status:     UpdateStatusIdle,
			Progress:   0.0,
			LastUpdate: time.Now(),
		},
		httpClient: &http.Client{
			Timeout: config.Timeout,
		},
		analytics:   make([]UpdateAnalytics, 0),
		subscribers: make([]UpdateSubscriber, 0),
		stopChan:    make(chan bool),
	}

	// Initialize scheduler with callback
	service.scheduler = NewUpdateScheduler(service.handleScheduledJob)

	return service
}

// Start initializes and starts the update service
func (u *UpdateService) Start() error {
	u.mutex.Lock()
	defer u.mutex.Unlock()

	if !u.settings.Enabled {
		return nil
	}

	// Start periodic update checks if enabled
	if u.settings.AutoCheck {
		u.startPeriodicChecks()
	}

	// Emit service started event
	u.emitEvent("update:service_started", map[string]interface{}{
		"settings": u.settings,
		"config":   u.config,
	})

	return nil
}

// Stop stops the update service
func (u *UpdateService) Stop() error {
	u.mutex.Lock()
	defer u.mutex.Unlock()

	if u.ticker != nil {
		u.ticker.Stop()
		u.ticker = nil
	}

	if u.scheduler != nil {
		u.scheduler.Stop()
	}

	select {
	case u.stopChan <- true:
	default:
	}

	u.emitEvent("update:service_stopped", map[string]interface{}{
		"timestamp": time.Now(),
	})

	return nil
}

// CheckForUpdates checks for available updates
func (u *UpdateService) CheckForUpdates() (*UpdateInfo, error) {
	u.mutex.Lock()
	defer u.mutex.Unlock()

	if u.currentStatus != UpdateStatusIdle {
		return nil, &APIError{
			Type:    ErrorTypeSystem,
			Code:    "UPDATE_IN_PROGRESS",
			Message: "Update operation already in progress",
		}
	}

	u.setStatus(UpdateStatusChecking)
	u.trackAnalytics(AnalyticsEventCheckStarted, "", "", true, "")

	check := &UpdateCheck{
		ID:         generateID(),
		CheckedAt:  time.Now(),
		Source:     UpdateSourceGitHub,
		UserAgent:  u.config.UserAgent,
		ClientInfo: u.getClientInfo(),
	}

	updateInfo, err := u.performUpdateCheck()
	if err != nil {
		check.Success = false
		check.Error = err.(*APIError)
		u.lastCheck = check
		u.setStatus(UpdateStatusIdle)
		u.trackAnalytics(AnalyticsEventCheckFailed, "", "", false, err.Error())
		return nil, err
	}

	check.Success = true
	check.UpdateInfo = updateInfo
	u.lastCheck = check

	if updateInfo != nil {
		u.setStatus(UpdateStatusAvailable)
		u.trackAnalytics(AnalyticsEventCheckCompleted, updateInfo.Version, "", true, "")
		u.notifySubscribers("OnUpdateAvailable", updateInfo)
		u.emitUpdateNotification(NotificationTypeUpdateAvailable, updateInfo)
	} else {
		u.setStatus(UpdateStatusIdle)
		u.trackAnalytics(AnalyticsEventCheckCompleted, "", "", true, "no_update")
	}

	return updateInfo, nil
}

// DownloadUpdate downloads the available update
func (u *UpdateService) DownloadUpdate(updateInfo *UpdateInfo) error {
	u.mutex.Lock()
	defer u.mutex.Unlock()

	if u.currentStatus != UpdateStatusAvailable {
		return &APIError{
			Type:    ErrorTypeSystem,
			Code:    "INVALID_UPDATE_STATE",
			Message: "No update available for download",
		}
	}

	u.setStatus(UpdateStatusDownloading)
	u.trackAnalytics(AnalyticsEventDownloadStarted, updateInfo.Version, "", true, "")

	// Create temp directory for download
	tempDir := filepath.Join(u.config.TempDirectory, "updates")
	if err := os.MkdirAll(tempDir, 0755); err != nil {
		u.setStatus(UpdateStatusFailed)
		apiErr := &APIError{
			Type:    ErrorTypeFileSystem,
			Code:    "TEMP_DIR_CREATION_FAILED",
			Message: fmt.Sprintf("Failed to create temp directory: %v", err),
		}
		u.trackAnalytics(AnalyticsEventDownloadFailed, updateInfo.Version, "", false, err.Error())
		return apiErr
	}

	// Download the update file
	downloadPath := filepath.Join(tempDir, fmt.Sprintf("mcpweaver_%s_%s_%s", updateInfo.Version, runtime.GOOS, runtime.GOARCH))
	if runtime.GOOS == "windows" {
		downloadPath += ".exe"
	}

	err := u.downloadFile(updateInfo.DownloadURL, downloadPath, updateInfo.Size)
	if err != nil {
		u.setStatus(UpdateStatusFailed)
		u.trackAnalytics(AnalyticsEventDownloadFailed, updateInfo.Version, "", false, err.Error())
		return err
	}

	// Verify the downloaded file
	u.setStatus(UpdateStatusVerifying)
	verificationResult, err := u.verifyUpdate(downloadPath, updateInfo)
	if err != nil {
		u.setStatus(UpdateStatusFailed)
		u.trackAnalytics(AnalyticsEventDownloadFailed, updateInfo.Version, "", false, err.Error())
		return err
	}

	if !verificationResult.ChecksumValid {
		u.setStatus(UpdateStatusFailed)
		apiErr := &APIError{
			Type:    ErrorTypeSystem,
			Code:    "VERIFICATION_FAILED",
			Message: "Update verification failed: checksum mismatch",
		}
		u.trackAnalytics(AnalyticsEventDownloadFailed, updateInfo.Version, "", false, "verification_failed")
		return apiErr
	}

	u.setStatus(UpdateStatusReady)
	u.trackAnalytics(AnalyticsEventDownloadCompleted, updateInfo.Version, "", true, "")

	// Store download path for installation
	u.progress.CurrentStep = "Ready for installation"
	u.progress.Progress = 100.0

	u.emitUpdateNotification(NotificationTypeUpdateReady, updateInfo)

	return nil
}

// InstallUpdate installs the downloaded update
func (u *UpdateService) InstallUpdate() error {
	u.mutex.Lock()
	defer u.mutex.Unlock()

	if u.currentStatus != UpdateStatusReady {
		return &APIError{
			Type:    ErrorTypeSystem,
			Code:    "INVALID_UPDATE_STATE",
			Message: "No update ready for installation",
		}
	}

	u.setStatus(UpdateStatusInstalling)
	u.trackAnalytics(AnalyticsEventInstallStarted, "", "", true, "")

	// Create backup of current executable
	backupInfo, err := u.createBackup()
	if err != nil {
		u.setStatus(UpdateStatusFailed)
		u.trackAnalytics(AnalyticsEventInstallFailed, "", "", false, err.Error())
		return err
	}

	// Install the update (this would typically restart the application)
	err = u.performInstallation()
	if err != nil {
		u.setStatus(UpdateStatusRollingBack)
		// Attempt rollback
		rollbackErr := u.performRollback(backupInfo)
		if rollbackErr != nil {
			u.setStatus(UpdateStatusFailed)
			u.trackAnalytics(AnalyticsEventInstallFailed, "", "", false, fmt.Sprintf("install_failed_rollback_failed: %v, %v", err, rollbackErr))
			return &APIError{
				Type:    ErrorTypeSystem,
				Code:    "INSTALLATION_AND_ROLLBACK_FAILED",
				Message: fmt.Sprintf("Installation failed and rollback failed: %v", rollbackErr),
				Details: map[string]string{
					"install_error":  err.Error(),
					"rollback_error": rollbackErr.Error(),
				},
				Severity: ErrorSeverityCritical,
			}
		}
		u.setStatus(UpdateStatusRollbackComplete)
		u.trackAnalytics(AnalyticsEventRollbackCompleted, "", "", true, "")
		return err
	}

	u.setStatus(UpdateStatusCompleted)
	u.trackAnalytics(AnalyticsEventInstallCompleted, "", "", true, "")

	result := &UpdateResult{
		Success:      true,
		UpdatedAt:    time.Now(),
		RollbackInfo: backupInfo,
	}

	u.notifySubscribers("OnUpdateCompleted", result)
	u.emitUpdateNotification(NotificationTypeUpdateCompleted, nil)

	return nil
}

// GetUpdateSettings returns the current update settings
func (u *UpdateService) GetUpdateSettings() *UpdateSettings {
	u.mutex.RLock()
	defer u.mutex.RUnlock()
	return u.settings
}

// UpdateSettings updates the update settings
func (u *UpdateService) UpdateSettings(settings *UpdateSettings) error {
	u.mutex.Lock()
	defer u.mutex.Unlock()

	u.settings = settings

	// Restart periodic checks if auto-check setting changed
	if settings.AutoCheck && u.ticker == nil {
		u.startPeriodicChecks()
	} else if !settings.AutoCheck && u.ticker != nil {
		u.ticker.Stop()
		u.ticker = nil
	}

	u.emitEvent("update:settings_changed", settings)

	return nil
}

// GetUpdateProgress returns the current update progress
func (u *UpdateService) GetUpdateProgress() *UpdateProgress {
	u.mutex.RLock()
	defer u.mutex.RUnlock()
	return u.progress
}

// GetLastUpdateCheck returns the last update check result
func (u *UpdateService) GetLastUpdateCheck() *UpdateCheck {
	u.mutex.RLock()
	defer u.mutex.RUnlock()
	return u.lastCheck
}

// Subscribe adds an update subscriber
func (u *UpdateService) Subscribe(subscriber UpdateSubscriber) {
	u.mutex.Lock()
	defer u.mutex.Unlock()
	u.subscribers = append(u.subscribers, subscriber)
}

// Unsubscribe removes an update subscriber
func (u *UpdateService) Unsubscribe(subscriber UpdateSubscriber) {
	u.mutex.Lock()
	defer u.mutex.Unlock()
	for i, sub := range u.subscribers {
		if sub == subscriber {
			u.subscribers = append(u.subscribers[:i], u.subscribers[i+1:]...)
			break
		}
	}
}

// Private methods

func (u *UpdateService) startPeriodicChecks() {
	if u.ticker != nil {
		u.ticker.Stop()
	}

	u.ticker = time.NewTicker(u.settings.CheckInterval)
	go func() {
		for {
			select {
			case <-u.ticker.C:
				if u.currentStatus == UpdateStatusIdle {
					go u.CheckForUpdates()
				}
			case <-u.stopChan:
				return
			}
		}
	}()
}

func (u *UpdateService) performUpdateCheck() (*UpdateInfo, error) {
	req, err := http.NewRequestWithContext(u.ctx, "GET", u.config.UpdateURL, nil)
	if err != nil {
		return nil, &APIError{
			Type:    ErrorTypeNetwork,
			Code:    "REQUEST_CREATION_FAILED",
			Message: fmt.Sprintf("Failed to create update request: %v", err),
		}
	}

	req.Header.Set("User-Agent", u.config.UserAgent)
	for key, value := range u.config.CustomHeaders {
		req.Header.Set(key, value)
	}

	resp, err := u.httpClient.Do(req)
	if err != nil {
		return nil, &APIError{
			Type:    ErrorTypeNetwork,
			Code:    "UPDATE_CHECK_FAILED",
			Message: fmt.Sprintf("Failed to check for updates: %v", err),
		}
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, &APIError{
			Type:    ErrorTypeNetwork,
			Code:    "UPDATE_CHECK_HTTP_ERROR",
			Message: fmt.Sprintf("Update check failed with status: %d", resp.StatusCode),
		}
	}

	var release struct {
		TagName     string `json:"tag_name"`
		Name        string `json:"name"`
		Body        string `json:"body"`
		PublishedAt string `json:"published_at"`
		Assets      []struct {
			Name        string `json:"name"`
			DownloadURL string `json:"browser_download_url"`
			Size        int64  `json:"size"`
		} `json:"assets"`
		Prerelease bool `json:"prerelease"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&release); err != nil {
		return nil, &APIError{
			Type:    ErrorTypeSystem,
			Code:    "RESPONSE_PARSE_FAILED",
			Message: fmt.Sprintf("Failed to parse update response: %v", err),
		}
	}

	// Skip pre-releases if not enabled
	if release.Prerelease && !u.settings.PreReleaseEnabled {
		return nil, nil
	}

	// Find the appropriate asset for current platform
	var asset *struct {
		Name        string `json:"name"`
		DownloadURL string `json:"browser_download_url"`
		Size        int64  `json:"size"`
	}

	expectedName := fmt.Sprintf("mcpweaver_%s_%s", runtime.GOOS, runtime.GOARCH)
	if runtime.GOOS == "windows" {
		expectedName += ".exe"
	}

	for _, a := range release.Assets {
		if a.Name == expectedName {
			asset = &a
			break
		}
	}

	if asset == nil {
		return nil, &APIError{
			Type:    ErrorTypeSystem,
			Code:    "NO_COMPATIBLE_ASSET",
			Message: fmt.Sprintf("No compatible asset found for platform %s/%s", runtime.GOOS, runtime.GOARCH),
		}
	}

	publishedAt, err := time.Parse(time.RFC3339, release.PublishedAt)
	if err != nil {
		publishedAt = time.Now()
	}

	updateInfo := &UpdateInfo{
		Version:      release.TagName,
		ReleaseNotes: release.Body,
		DownloadURL:  asset.DownloadURL,
		Size:         asset.Size,
		PublishedAt:  publishedAt,
		Critical:     false, // TODO: Determine from release notes or metadata
	}

	// Check if this is actually a newer version
	currentVersion := u.getCurrentVersion()
	if !u.isNewerVersion(updateInfo.Version, currentVersion) {
		return nil, nil
	}

	return updateInfo, nil
}

func (u *UpdateService) downloadFile(url, filepath string, expectedSize int64) error {
	req, err := http.NewRequestWithContext(u.ctx, "GET", url, nil)
	if err != nil {
		return &APIError{
			Type:    ErrorTypeNetwork,
			Code:    "DOWNLOAD_REQUEST_FAILED",
			Message: fmt.Sprintf("Failed to create download request: %v", err),
		}
	}

	resp, err := u.httpClient.Do(req)
	if err != nil {
		return &APIError{
			Type:    ErrorTypeNetwork,
			Code:    "DOWNLOAD_FAILED",
			Message: fmt.Sprintf("Failed to download update: %v", err),
		}
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return &APIError{
			Type:    ErrorTypeNetwork,
			Code:    "DOWNLOAD_HTTP_ERROR",
			Message: fmt.Sprintf("Download failed with status: %d", resp.StatusCode),
		}
	}

	file, err := os.Create(filepath)
	if err != nil {
		return &APIError{
			Type:    ErrorTypeFileSystem,
			Code:    "FILE_CREATION_FAILED",
			Message: fmt.Sprintf("Failed to create download file: %v", err),
		}
	}
	defer file.Close()

	// Track download progress
	var bytesReceived int64
	startTime := time.Now()
	lastUpdate := startTime

	buffer := make([]byte, 32*1024) // 32KB buffer
	for {
		n, err := resp.Body.Read(buffer)
		if n > 0 {
			if _, writeErr := file.Write(buffer[:n]); writeErr != nil {
				return &APIError{
					Type:    ErrorTypeFileSystem,
					Code:    "FILE_WRITE_FAILED",
					Message: fmt.Sprintf("Failed to write to file: %v", writeErr),
				}
			}
			bytesReceived += int64(n)

			// Update progress
			now := time.Now()
			if now.Sub(lastUpdate) >= time.Second {
				progress := float64(bytesReceived) / float64(expectedSize) * 100
				speed := bytesReceived / int64(now.Sub(startTime).Seconds())

				u.progress.Progress = progress
				u.progress.BytesReceived = bytesReceived
				u.progress.BytesTotal = expectedSize
				u.progress.Speed = speed
				u.progress.LastUpdate = now

				if speed > 0 {
					remaining := (expectedSize - bytesReceived) / speed
					estimatedTime := time.Duration(remaining) * time.Second
					u.progress.EstimatedTime = &estimatedTime
				}

				u.notifySubscribers("OnUpdateProgress", u.progress)
				u.emitEvent("update:progress", u.progress)
				lastUpdate = now
			}
		}
		if err == io.EOF {
			break
		}
		if err != nil {
			return &APIError{
				Type:    ErrorTypeNetwork,
				Code:    "DOWNLOAD_READ_FAILED",
				Message: fmt.Sprintf("Failed to read download data: %v", err),
			}
		}
	}

	return nil
}

func (u *UpdateService) verifyUpdate(filePath string, updateInfo *UpdateInfo) (*VerificationResult, error) {
	result := &VerificationResult{
		VerifiedAt: time.Now(),
		Algorithm:  "SHA256",
	}

	if u.config.VerificationMode == VerificationModeNone {
		result.ChecksumValid = true
		result.TrustedSource = true
		return result, nil
	}

	// Verify checksum
	if u.config.VerificationMode >= VerificationModeChecksum {
		if updateInfo.ChecksumURL == "" {
			return result, &APIError{
				Type:    ErrorTypeSystem,
				Code:    "NO_CHECKSUM_URL",
				Message: "Checksum verification required but no checksum URL provided",
			}
		}

		expectedChecksum, err := u.downloadChecksum(updateInfo.ChecksumURL)
		if err != nil {
			return result, err
		}

		actualChecksum, err := u.calculateChecksum(filePath)
		if err != nil {
			return result, err
		}

		result.ChecksumValid = expectedChecksum == actualChecksum
	}

	return result, nil
}

func (u *UpdateService) downloadChecksum(url string) (string, error) {
	resp, err := u.httpClient.Get(url)
	if err != nil {
		return "", &APIError{
			Type:    ErrorTypeNetwork,
			Code:    "CHECKSUM_DOWNLOAD_FAILED",
			Message: fmt.Sprintf("Failed to download checksum: %v", err),
		}
	}
	defer resp.Body.Close()

	checksum, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", &APIError{
			Type:    ErrorTypeNetwork,
			Code:    "CHECKSUM_READ_FAILED",
			Message: fmt.Sprintf("Failed to read checksum: %v", err),
		}
	}

	return string(checksum), nil
}

func (u *UpdateService) calculateChecksum(filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", &APIError{
			Type:    ErrorTypeFileSystem,
			Code:    "FILE_OPEN_FAILED",
			Message: fmt.Sprintf("Failed to open file for checksum: %v", err),
		}
	}
	defer file.Close()

	hash := sha256.New()
	if _, err := io.Copy(hash, file); err != nil {
		return "", &APIError{
			Type:    ErrorTypeFileSystem,
			Code:    "CHECKSUM_CALCULATION_FAILED",
			Message: fmt.Sprintf("Failed to calculate checksum: %v", err),
		}
	}

	return hex.EncodeToString(hash.Sum(nil)), nil
}

func (u *UpdateService) createBackup() (*RollbackInfo, error) {
	currentExePath, err := os.Executable()
	if err != nil {
		return nil, &APIError{
			Type:    ErrorTypeSystem,
			Code:    "GET_EXECUTABLE_PATH_FAILED",
			Message: fmt.Sprintf("Failed to get current executable path: %v", err),
		}
	}

	return u.rollbackManager.CreateBackup(currentExePath, u.getCurrentVersion())
}

func (u *UpdateService) performInstallation() error {
	// TODO: Implement actual installation
	// This would typically involve replacing the current executable
	// and restarting the application
	return nil
}

func (u *UpdateService) performRollback(backupInfo *RollbackInfo) error {
	currentExePath, err := os.Executable()
	if err != nil {
		return &APIError{
			Type:    ErrorTypeSystem,
			Code:    "GET_EXECUTABLE_PATH_FAILED",
			Message: fmt.Sprintf("Failed to get current executable path: %v", err),
		}
	}

	return u.rollbackManager.PerformRollback(backupInfo, currentExePath)
}

func (u *UpdateService) setStatus(status UpdateStatus) {
	u.currentStatus = status
	u.progress.Status = status
	u.progress.LastUpdate = time.Now()
}

func (u *UpdateService) getCurrentVersion() string {
	// TODO: Get actual current version
	return "1.0.0"
}

func (u *UpdateService) isNewerVersion(newVersion, currentVersion string) bool {
	// TODO: Implement proper semantic version comparison
	return newVersion != currentVersion
}

func (u *UpdateService) getClientInfo() ClientInfo {
	return ClientInfo{
		Version:      u.getCurrentVersion(),
		Platform:     runtime.GOOS,
		Architecture: runtime.GOARCH,
		OS:           runtime.GOOS,
	}
}

func (u *UpdateService) trackAnalytics(eventType AnalyticsEvent, version, previousVersion string, success bool, errorMsg string) {
	analytics := UpdateAnalytics{
		EventType:       eventType,
		Version:         version,
		PreviousVersion: previousVersion,
		UpdateChannel:   u.settings.UpdateChannel,
		Success:         success,
		Error:           errorMsg,
		ClientInfo:      u.getClientInfo(),
		Timestamp:       time.Now(),
	}

	u.analytics = append(u.analytics, analytics)
	u.emitEvent("update:analytics", analytics)
}

func (u *UpdateService) notifySubscribers(method string, data interface{}) {
	for _, subscriber := range u.subscribers {
		switch method {
		case "OnUpdateAvailable":
			subscriber.OnUpdateAvailable(data.(*UpdateInfo))
		case "OnUpdateProgress":
			subscriber.OnUpdateProgress(data.(*UpdateProgress))
		case "OnUpdateCompleted":
			subscriber.OnUpdateCompleted(data.(*UpdateResult))
		case "OnUpdateFailed":
			subscriber.OnUpdateFailed(data.(*APIError))
		}
	}
}

func (u *UpdateService) emitEvent(eventName string, data interface{}) {
	// Skip events in test mode to avoid Wails runtime issues
	if u.isTestMode {
		return
	}

	// Only emit events if we have a valid Wails context
	if u.ctx != nil {
		defer func() {
			if r := recover(); r != nil {
				// Silently handle context-related panics during testing
			}
		}()
		wailsruntime.EventsEmit(u.ctx, eventName, data)
	}
}

func (u *UpdateService) emitUpdateNotification(notificationType NotificationType, updateInfo *UpdateInfo) {
	// Skip notifications in test mode
	if u.isTestMode {
		return
	}

	var title, message string
	var actions []NotificationAction
	var priority NotificationPriority

	switch notificationType {
	case NotificationTypeUpdateAvailable:
		title = "Update Available"
		message = fmt.Sprintf("MCPWeaver %s is available. Would you like to download it?", updateInfo.Version)
		actions = []NotificationAction{
			{ID: "download", Label: "Download Now", Type: ActionTypeInstallNow},
			{ID: "later", Label: "Remind Me Later", Type: ActionTypeInstallLater},
			{ID: "skip", Label: "Skip This Version", Type: ActionTypeSkipVersion},
		}
		priority = NotificationPriorityMedium

	case NotificationTypeUpdateReady:
		title = "Update Ready"
		message = "The update has been downloaded and verified. Restart to apply the update."
		actions = []NotificationAction{
			{ID: "install", Label: "Restart and Update", Type: ActionTypeInstallNow},
			{ID: "later", Label: "Install Later", Type: ActionTypeInstallLater},
		}
		priority = NotificationPriorityHigh

	case NotificationTypeUpdateCompleted:
		title = "Update Completed"
		message = "MCPWeaver has been successfully updated."
		priority = NotificationPriorityLow

	case NotificationTypeCriticalUpdate:
		title = "Critical Update Available"
		message = "A critical security update is available. Please update immediately."
		actions = []NotificationAction{
			{ID: "install", Label: "Install Now", Type: ActionTypeInstallNow},
		}
		priority = NotificationPriorityCritical
	}

	notification := &UpdateNotification{
		Type:       notificationType,
		Title:      title,
		Message:    message,
		Actions:    actions,
		UpdateInfo: updateInfo,
		Timestamp:  time.Now(),
		Priority:   priority,
		Persistent: priority == NotificationPriorityCritical,
	}

	u.emitEvent("update:notification", notification)
}

func generateID() string {
	return fmt.Sprintf("%d", time.Now().UnixNano())
}

// handleScheduledJob handles scheduled update operations
func (u *UpdateService) handleScheduledJob(jobType ScheduledJobType) error {
	switch jobType {
	case ScheduledJobTypeCheck:
		_, err := u.CheckForUpdates()
		return err
	case ScheduledJobTypeDownload:
		// Get the last available update info
		if u.lastCheck != nil && u.lastCheck.UpdateInfo != nil {
			return u.DownloadUpdate(u.lastCheck.UpdateInfo)
		}
		return &APIError{
			Type:    ErrorTypeSystem,
			Code:    "NO_UPDATE_INFO",
			Message: "No update information available for download",
		}
	case ScheduledJobTypeInstall:
		return u.InstallUpdate()
	default:
		return &APIError{
			Type:    ErrorTypeValidation,
			Code:    "INVALID_JOB_TYPE",
			Message: fmt.Sprintf("Invalid scheduled job type: %s", jobType),
		}
	}
}
