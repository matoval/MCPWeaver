package app

import (
	"fmt"
)

// Auto-Update API Methods

// CheckForUpdates checks for available updates
func (a *App) CheckForUpdates() (*UpdateInfo, error) {
	if a.updateService == nil {
		return nil, a.createAPIError(ErrorTypeSystem, "UPDATE_SERVICE_NOT_AVAILABLE", "Update service is not available", nil)
	}

	updateInfo, err := a.updateService.CheckForUpdates()
	if err != nil {
		a.emitError(err.(*APIError))
		return nil, err
	}

	return updateInfo, nil
}

// GetUpdateStatus returns the current update status and progress
func (a *App) GetUpdateStatus() *UpdateProgress {
	if a.updateService == nil {
		return &UpdateProgress{
			Status:   UpdateStatusIdle,
			Progress: 0.0,
			Error: &APIError{
				Type:    ErrorTypeSystem,
				Code:    "UPDATE_SERVICE_NOT_AVAILABLE",
				Message: "Update service is not available",
			},
		}
	}

	return a.updateService.GetUpdateProgress()
}

// GetUpdateSettings returns the current update settings
func (a *App) GetUpdateSettings() *UpdateSettings {
	if a.updateService == nil {
		return &UpdateSettings{
			Enabled: false,
		}
	}

	return a.updateService.GetUpdateSettings()
}

// UpdateUpdateSettings updates the update settings
func (a *App) UpdateUpdateSettings(settings *UpdateSettings) error {
	if a.updateService == nil {
		return a.createAPIError(ErrorTypeSystem, "UPDATE_SERVICE_NOT_AVAILABLE", "Update service is not available", nil)
	}

	err := a.updateService.UpdateSettings(settings)
	if err != nil {
		apiErr := a.createAPIError(ErrorTypeSystem, "UPDATE_SETTINGS_FAILED",
			fmt.Sprintf("Failed to update settings: %v", err), nil)
		a.emitError(apiErr)
		return apiErr
	}

	// Emit settings updated event
	a.emitNotification("info", "Settings Updated", "Update settings have been successfully updated")

	return nil
}

// DownloadUpdate downloads the available update
func (a *App) DownloadUpdate(updateInfo *UpdateInfo) error {
	if a.updateService == nil {
		return a.createAPIError(ErrorTypeSystem, "UPDATE_SERVICE_NOT_AVAILABLE", "Update service is not available", nil)
	}

	if updateInfo == nil {
		return a.createAPIError(ErrorTypeValidation, "INVALID_UPDATE_INFO", "Update info is required", nil)
	}

	err := a.updateService.DownloadUpdate(updateInfo)
	if err != nil {
		a.emitError(err.(*APIError))
		return err
	}

	return nil
}

// InstallUpdate installs the downloaded update
func (a *App) InstallUpdate() error {
	if a.updateService == nil {
		return a.createAPIError(ErrorTypeSystem, "UPDATE_SERVICE_NOT_AVAILABLE", "Update service is not available", nil)
	}

	err := a.updateService.InstallUpdate()
	if err != nil {
		a.emitError(err.(*APIError))
		return err
	}

	return nil
}

// GetLastUpdateCheck returns information about the last update check
func (a *App) GetLastUpdateCheck() *UpdateCheck {
	if a.updateService == nil {
		return &UpdateCheck{
			Success: false,
			Error: &APIError{
				Type:    ErrorTypeSystem,
				Code:    "UPDATE_SERVICE_NOT_AVAILABLE",
				Message: "Update service is not available",
			},
		}
	}

	return a.updateService.GetLastUpdateCheck()
}

// SkipUpdateVersion marks an update version as skipped
func (a *App) SkipUpdateVersion(version string) error {
	if version == "" {
		return a.createAPIError(ErrorTypeValidation, "INVALID_VERSION", "Version is required", nil)
	}

	// TODO: Implement version skipping persistence
	// For now, just emit an event
	a.emitNotification("info", "Version Skipped",
		fmt.Sprintf("Version %s has been skipped", version))

	return nil
}

// ScheduleUpdate schedules an update for later installation
func (a *App) ScheduleUpdate(updateInfo *UpdateInfo, schedule *UpdateSchedule) error {
	if a.updateService == nil || a.updateService.scheduler == nil {
		return a.createAPIError(ErrorTypeSystem, "UPDATE_SERVICE_NOT_AVAILABLE", "Update service is not available", nil)
	}

	if updateInfo == nil {
		return a.createAPIError(ErrorTypeValidation, "INVALID_UPDATE_INFO", "Update info is required", nil)
	}

	if schedule == nil {
		return a.createAPIError(ErrorTypeValidation, "INVALID_SCHEDULE", "Update schedule is required", nil)
	}

	// Schedule the update job
	job, err := a.updateService.scheduler.ScheduleJob(ScheduledJobTypeInstall, schedule, updateInfo)
	if err != nil {
		a.emitError(err.(*APIError))
		return err
	}

	a.emitNotification("info", "Update Scheduled",
		fmt.Sprintf("Update to version %s has been scheduled for %s", updateInfo.Version, job.NextRun.Format("2006-01-02 15:04:05")))

	return nil
}

// CancelScheduledUpdate cancels a scheduled update
func (a *App) CancelScheduledUpdate() error {
	if a.updateService == nil || a.updateService.scheduler == nil {
		return a.createAPIError(ErrorTypeSystem, "UPDATE_SERVICE_NOT_AVAILABLE", "Update service is not available", nil)
	}

	err := a.updateService.scheduler.CancelScheduledJob()
	if err != nil {
		a.emitError(err.(*APIError))
		return err
	}

	a.emitNotification("info", "Update Cancelled", "Scheduled update has been cancelled")

	return nil
}

// RollbackUpdate rolls back to the previous version
func (a *App) RollbackUpdate() error {
	if a.updateService == nil {
		return a.createAPIError(ErrorTypeSystem, "UPDATE_SERVICE_NOT_AVAILABLE", "Update service is not available", nil)
	}

	// Get rollback capabilities first
	capabilities := a.GetRollbackCapabilities()
	if !capabilities.Available {
		return a.createAPIError(ErrorTypeSystem, "NO_ROLLBACK_AVAILABLE", "No rollback available", nil)
	}

	// Get available backups
	backups, err := a.GetAvailableBackups()
	if err != nil {
		return err
	}

	if len(backups) == 0 {
		return a.createAPIError(ErrorTypeSystem, "NO_BACKUPS_AVAILABLE", "No backups available for rollback", nil)
	}

	// Use the most recent backup
	latestBackup := backups[0]
	for _, backup := range backups {
		if backup.CreatedAt.After(latestBackup.CreatedAt) {
			latestBackup = backup
		}
	}

	// Create rollback info from backup
	rollbackInfo := &RollbackInfo{
		Available:       true,
		BackupPath:      latestBackup.Path,
		BackupVersion:   latestBackup.Version,
		BackupCreatedAt: latestBackup.CreatedAt,
		BackupSize:      latestBackup.Size,
	}

	// Perform rollback
	err = a.updateService.rollbackManager.PerformRollback(rollbackInfo, "")
	if err != nil {
		a.emitError(err.(*APIError))
		return err
	}

	a.emitNotification("info", "Rollback Completed", "Application has been rolled back to the previous version")
	return nil
}

// GetUpdateHistory returns the update history
func (a *App) GetUpdateHistory() ([]UpdateResult, error) {
	if a.updateService == nil {
		return nil, a.createAPIError(ErrorTypeSystem, "UPDATE_SERVICE_NOT_AVAILABLE", "Update service is not available", nil)
	}

	// Get update history from the service
	history, err := a.updateService.GetUpdateHistory()
	if err != nil {
		a.emitError(err.(*APIError))
		return nil, err
	}

	return history, nil
}

// GetUpdateAnalytics returns update analytics data
func (a *App) GetUpdateAnalytics() ([]UpdateAnalytics, error) {
	if a.updateService == nil {
		return nil, a.createAPIError(ErrorTypeSystem, "UPDATE_SERVICE_NOT_AVAILABLE", "Update service is not available", nil)
	}

	// Get analytics data from the service
	analytics, err := a.updateService.GetUpdateAnalytics()
	if err != nil {
		a.emitError(err.(*APIError))
		return nil, err
	}

	return analytics, nil
}

// TestUpdateConnection tests the connection to the update server
func (a *App) TestUpdateConnection() (*UpdateConnectionTest, error) {
	if a.updateService == nil {
		return nil, a.createAPIError(ErrorTypeSystem, "UPDATE_SERVICE_NOT_AVAILABLE", "Update service is not available", nil)
	}

	// TODO: Implement connection test
	// For now, return a mock successful test
	return &UpdateConnectionTest{
		Success:      true,
		ResponseTime: 100, // milliseconds
		ServerInfo: &UpdateServerInfo{
			Version: "1.0.0",
			Status:  "online",
		},
	}, nil
}

// PauseUpdateDownload pauses an ongoing update download
func (a *App) PauseUpdateDownload() error {
	if a.updateService == nil {
		return a.createAPIError(ErrorTypeSystem, "UPDATE_SERVICE_NOT_AVAILABLE", "Update service is not available", nil)
	}

	// Get current progress to check if download is in progress
	progress := a.updateService.GetUpdateProgress()
	if progress.Status != UpdateStatusDownloading {
		return a.createAPIError(ErrorTypeValidation, "INVALID_DOWNLOAD_STATE", 
			"No download is currently in progress to pause", map[string]string{
				"current_status": string(progress.Status),
			})
	}

	// Pause the download by changing status
	err := a.updateService.PauseDownload()
	if err != nil {
		a.emitError(err.(*APIError))
		return err
	}

	// Emit notification
	a.emitNotification("info", "Download Paused", "Update download has been paused")

	// Emit progress event
	pausedProgress := a.updateService.GetUpdateProgress()
	a.emitEvent("update:download_paused", pausedProgress)

	return nil
}

// ResumeUpdateDownload resumes a paused update download
func (a *App) ResumeUpdateDownload() error {
	if a.updateService == nil {
		return a.createAPIError(ErrorTypeSystem, "UPDATE_SERVICE_NOT_AVAILABLE", "Update service is not available", nil)
	}

	// Get current progress to check if download is paused
	progress := a.updateService.GetUpdateProgress()
	if progress.Status != UpdateStatusPaused {
		return a.createAPIError(ErrorTypeValidation, "INVALID_DOWNLOAD_STATE", 
			"No paused download to resume", map[string]string{
				"current_status": string(progress.Status),
			})
	}

	// Resume the download
	err := a.updateService.ResumeDownload()
	if err != nil {
		a.emitError(err.(*APIError))
		return err
	}

	// Emit notification
	a.emitNotification("info", "Download Resumed", "Update download has been resumed")

	// Emit progress event
	resumedProgress := a.updateService.GetUpdateProgress()
	a.emitEvent("update:download_resumed", resumedProgress)

	return nil
}

// CancelUpdateDownload cancels an ongoing update download
func (a *App) CancelUpdateDownload() error {
	if a.updateService == nil {
		return a.createAPIError(ErrorTypeSystem, "UPDATE_SERVICE_NOT_AVAILABLE", "Update service is not available", nil)
	}

	// Get current progress to check if download can be cancelled
	progress := a.updateService.GetUpdateProgress()
	if progress.Status != UpdateStatusDownloading && progress.Status != UpdateStatusPaused {
		return a.createAPIError(ErrorTypeValidation, "INVALID_DOWNLOAD_STATE", 
			"No active or paused download to cancel", map[string]string{
				"current_status": string(progress.Status),
			})
	}

	// Cancel the download
	err := a.updateService.CancelDownload()
	if err != nil {
		a.emitError(err.(*APIError))
		return err
	}

	// Emit notification
	a.emitNotification("info", "Download Cancelled", "Update download has been cancelled")

	// Emit progress event
	cancelledProgress := a.updateService.GetUpdateProgress()
	a.emitEvent("update:download_cancelled", cancelledProgress)

	return nil
}

// GetUpdateConfiguration returns the current update configuration
func (a *App) GetUpdateConfiguration() *UpdateConfiguration {
	if a.updateService == nil {
		defaultConfig := DefaultUpdateConfiguration()
		return &defaultConfig
	}

	// TODO: Expose configuration safely (without sensitive data)
	defaultConfig := DefaultUpdateConfiguration()
	return &defaultConfig
}

// ValidateUpdateConfiguration validates the update configuration
func (a *App) ValidateUpdateConfiguration(config *UpdateConfiguration) (*UpdateConfigurationValidation, error) {
	if config == nil {
		return nil, a.createAPIError(ErrorTypeValidation, "INVALID_CONFIG", "Configuration is required", nil)
	}

	// TODO: Implement configuration validation
	return &UpdateConfigurationValidation{
		Valid:    true,
		Errors:   []string{},
		Warnings: []string{},
	}, nil
}

// GetRollbackCapabilities returns information about rollback capabilities
func (a *App) GetRollbackCapabilities() *RollbackCapabilities {
	if a.updateService == nil || a.updateService.rollbackManager == nil {
		return &RollbackCapabilities{
			Available:   false,
			BackupCount: 0,
			Features:    []string{},
		}
	}

	return a.updateService.rollbackManager.GetRollbackCapabilities()
}

// GetAvailableBackups returns a list of available backups
func (a *App) GetAvailableBackups() ([]BackupInfo, error) {
	if a.updateService == nil || a.updateService.rollbackManager == nil {
		return nil, a.createAPIError(ErrorTypeSystem, "UPDATE_SERVICE_NOT_AVAILABLE", "Update service is not available", nil)
	}

	backups, err := a.updateService.rollbackManager.ListAvailableBackups()
	if err != nil {
		a.emitError(err.(*APIError))
		return nil, err
	}

	return backups, nil
}

// ValidateBackup validates a specific backup file
func (a *App) ValidateBackup(backupPath string) (*BackupValidation, error) {
	if a.updateService == nil || a.updateService.rollbackManager == nil {
		return nil, a.createAPIError(ErrorTypeSystem, "UPDATE_SERVICE_NOT_AVAILABLE", "Update service is not available", nil)
	}

	if backupPath == "" {
		return nil, a.createAPIError(ErrorTypeValidation, "INVALID_BACKUP_PATH", "Backup path is required", nil)
	}

	validation, err := a.updateService.rollbackManager.ValidateBackup(backupPath)
	if err != nil {
		a.emitError(err.(*APIError))
		return nil, err
	}

	return validation, nil
}

// DeleteBackup deletes a specific backup file
func (a *App) DeleteBackup(backupPath string) error {
	if a.updateService == nil || a.updateService.rollbackManager == nil {
		return a.createAPIError(ErrorTypeSystem, "UPDATE_SERVICE_NOT_AVAILABLE", "Update service is not available", nil)
	}

	if backupPath == "" {
		return a.createAPIError(ErrorTypeValidation, "INVALID_BACKUP_PATH", "Backup path is required", nil)
	}

	err := a.updateService.rollbackManager.DeleteBackup(backupPath)
	if err != nil {
		a.emitError(err.(*APIError))
		return err
	}

	a.emitNotification("info", "Backup Deleted", "Backup file has been successfully deleted")
	return nil
}

// RollbackToVersion rolls back to a specific version backup
func (a *App) RollbackToVersion(version string) error {
	if a.updateService == nil || a.updateService.rollbackManager == nil {
		return a.createAPIError(ErrorTypeSystem, "UPDATE_SERVICE_NOT_AVAILABLE", "Update service is not available", nil)
	}

	if version == "" {
		return a.createAPIError(ErrorTypeValidation, "INVALID_VERSION", "Version is required", nil)
	}

	// Get available backups
	backups, err := a.GetAvailableBackups()
	if err != nil {
		return err
	}

	// Find backup for the specified version
	var targetBackup *BackupInfo
	for _, backup := range backups {
		if backup.Version == version {
			targetBackup = &backup
			break
		}
	}

	if targetBackup == nil {
		return a.createAPIError(ErrorTypeSystem, "BACKUP_NOT_FOUND",
			fmt.Sprintf("No backup found for version %s", version), nil)
	}

	// Create rollback info
	rollbackInfo := &RollbackInfo{
		Available:       true,
		BackupPath:      targetBackup.Path,
		BackupVersion:   targetBackup.Version,
		BackupCreatedAt: targetBackup.CreatedAt,
		BackupSize:      targetBackup.Size,
	}

	// Perform rollback
	err = a.updateService.rollbackManager.PerformRollback(rollbackInfo, "")
	if err != nil {
		a.emitError(err.(*APIError))
		return err
	}

	a.emitNotification("info", "Rollback Completed",
		fmt.Sprintf("Application has been rolled back to version %s", version))
	return nil
}

// GetScheduledJob returns information about the current scheduled job
func (a *App) GetScheduledJob() *ScheduledJob {
	if a.updateService == nil || a.updateService.scheduler == nil {
		return nil
	}

	return a.updateService.scheduler.GetScheduledJob()
}

// PauseScheduledUpdate pauses the current scheduled update
func (a *App) PauseScheduledUpdate() error {
	if a.updateService == nil || a.updateService.scheduler == nil {
		return a.createAPIError(ErrorTypeSystem, "UPDATE_SERVICE_NOT_AVAILABLE", "Update service is not available", nil)
	}

	err := a.updateService.scheduler.PauseScheduledJob()
	if err != nil {
		a.emitError(err.(*APIError))
		return err
	}

	a.emitNotification("info", "Update Paused", "Scheduled update has been paused")
	return nil
}

// ResumeScheduledUpdate resumes a paused scheduled update
func (a *App) ResumeScheduledUpdate() error {
	if a.updateService == nil || a.updateService.scheduler == nil {
		return a.createAPIError(ErrorTypeSystem, "UPDATE_SERVICE_NOT_AVAILABLE", "Update service is not available", nil)
	}

	err := a.updateService.scheduler.ResumeScheduledJob()
	if err != nil {
		a.emitError(err.(*APIError))
		return err
	}

	a.emitNotification("info", "Update Resumed", "Scheduled update has been resumed")
	return nil
}

// UpdateScheduledUpdate updates the schedule of a scheduled update
func (a *App) UpdateScheduledUpdate(newSchedule *UpdateSchedule) error {
	if a.updateService == nil || a.updateService.scheduler == nil {
		return a.createAPIError(ErrorTypeSystem, "UPDATE_SERVICE_NOT_AVAILABLE", "Update service is not available", nil)
	}

	if newSchedule == nil {
		return a.createAPIError(ErrorTypeValidation, "INVALID_SCHEDULE", "New schedule is required", nil)
	}

	err := a.updateService.scheduler.UpdateSchedule(newSchedule)
	if err != nil {
		a.emitError(err.(*APIError))
		return err
	}

	a.emitNotification("info", "Schedule Updated", "Update schedule has been successfully updated")
	return nil
}

// ScheduleUpdateCheck schedules automatic update checks
func (a *App) ScheduleUpdateCheck(schedule *UpdateSchedule) error {
	if a.updateService == nil || a.updateService.scheduler == nil {
		return a.createAPIError(ErrorTypeSystem, "UPDATE_SERVICE_NOT_AVAILABLE", "Update service is not available", nil)
	}

	if schedule == nil {
		return a.createAPIError(ErrorTypeValidation, "INVALID_SCHEDULE", "Update schedule is required", nil)
	}

	// Schedule the check job
	job, err := a.updateService.scheduler.ScheduleJob(ScheduledJobTypeCheck, schedule, nil)
	if err != nil {
		a.emitError(err.(*APIError))
		return err
	}

	a.emitNotification("info", "Update Check Scheduled",
		fmt.Sprintf("Automatic update checks have been scheduled for %s", job.NextRun.Format("2006-01-02 15:04:05")))

	return nil
}

// Supporting types for API methods

// UpdateConnectionTest represents the result of testing update server connection
type UpdateConnectionTest struct {
	Success      bool              `json:"success"`
	ResponseTime int64             `json:"responseTime"` // milliseconds
	Error        *APIError         `json:"error,omitempty"`
	ServerInfo   *UpdateServerInfo `json:"serverInfo,omitempty"`
	TestedAt     string            `json:"testedAt"`
}

// UpdateServerInfo contains information about the update server
type UpdateServerInfo struct {
	Version     string `json:"version"`
	Status      string `json:"status"`
	LastUpdated string `json:"lastUpdated,omitempty"`
}

// UpdateConfigurationValidation represents the result of validating update configuration
type UpdateConfigurationValidation struct {
	Valid    bool     `json:"valid"`
	Errors   []string `json:"errors"`
	Warnings []string `json:"warnings"`
}

