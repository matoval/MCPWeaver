package app

import (
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// RollbackManager handles update rollback operations
type RollbackManager struct {
	backupDirectory string
	maxBackups      int
}

// NewRollbackManager creates a new rollback manager
func NewRollbackManager(backupDir string) *RollbackManager {
	return &RollbackManager{
		backupDirectory: backupDir,
		maxBackups:      5, // Keep last 5 backups
	}
}

// CreateBackup creates a backup of the current executable before update
func (rm *RollbackManager) CreateBackup(currentExePath, version string) (*RollbackInfo, error) {
	// Ensure backup directory exists
	if err := os.MkdirAll(rm.backupDirectory, 0755); err != nil {
		return nil, &APIError{
			Type:    ErrorTypeFileSystem,
			Code:    "BACKUP_DIR_CREATION_FAILED",
			Message: fmt.Sprintf("Failed to create backup directory: %v", err),
		}
	}

	// Get current executable info
	stat, err := os.Stat(currentExePath)
	if err != nil {
		return nil, &APIError{
			Type:    ErrorTypeFileSystem,
			Code:    "CURRENT_EXE_STAT_FAILED",
			Message: fmt.Sprintf("Failed to stat current executable: %v", err),
		}
	}

	// Create timestamped backup filename
	timestamp := time.Now().Format("20060102-150405")
	backupFileName := fmt.Sprintf("mcpweaver_%s_%s.backup", version, timestamp)
	backupPath := filepath.Join(rm.backupDirectory, backupFileName)

	// Copy current executable to backup location
	err = rm.copyFile(currentExePath, backupPath)
	if err != nil {
		return nil, &APIError{
			Type:    ErrorTypeFileSystem,
			Code:    "BACKUP_COPY_FAILED",
			Message: fmt.Sprintf("Failed to create backup copy: %v", err),
		}
	}

	// Set executable permissions on backup
	if err := os.Chmod(backupPath, stat.Mode()); err != nil {
		return nil, &APIError{
			Type:    ErrorTypeFileSystem,
			Code:    "BACKUP_PERMISSIONS_FAILED",
			Message: fmt.Sprintf("Failed to set backup permissions: %v", err),
		}
	}

	// Clean up old backups
	rm.cleanupOldBackups()

	return &RollbackInfo{
		Available:       true,
		BackupPath:      backupPath,
		BackupVersion:   version,
		BackupCreatedAt: time.Now(),
		BackupSize:      stat.Size(),
	}, nil
}

// PerformRollback performs a rollback using the specified backup
func (rm *RollbackManager) PerformRollback(rollbackInfo *RollbackInfo, currentExePath string) error {
	if rollbackInfo == nil || !rollbackInfo.Available {
		return &APIError{
			Type:    ErrorTypeSystem,
			Code:    "NO_ROLLBACK_AVAILABLE",
			Message: "No rollback information available",
		}
	}

	// Verify backup file exists
	if !rm.fileExists(rollbackInfo.BackupPath) {
		return &APIError{
			Type:    ErrorTypeFileSystem,
			Code:    "BACKUP_FILE_NOT_FOUND",
			Message: fmt.Sprintf("Backup file not found: %s", rollbackInfo.BackupPath),
		}
	}

	// Create a backup of the current (failed) version before rollback
	failedBackupPath := fmt.Sprintf("%s.failed_%s", currentExePath, time.Now().Format("20060102-150405"))
	if err := rm.copyFile(currentExePath, failedBackupPath); err != nil {
		return &APIError{
			Type:    ErrorTypeFileSystem,
			Code:    "FAILED_VERSION_BACKUP_ERROR",
			Message: fmt.Sprintf("Failed to backup current version before rollback: %v", err),
		}
	}

	// Get original file permissions
	stat, err := os.Stat(currentExePath)
	if err != nil {
		return &APIError{
			Type:    ErrorTypeFileSystem,
			Code:    "CURRENT_EXE_STAT_FAILED",
			Message: fmt.Sprintf("Failed to get current executable info: %v", err),
		}
	}

	// Replace current executable with backup
	err = rm.copyFile(rollbackInfo.BackupPath, currentExePath)
	if err != nil {
		// Try to restore from failed backup if rollback copy failed
		if restoreErr := rm.copyFile(failedBackupPath, currentExePath); restoreErr != nil {
			return &APIError{
				Type:    ErrorTypeSystem,
				Code:    "ROLLBACK_AND_RESTORE_FAILED",
				Message: fmt.Sprintf("Rollback failed and unable to restore: rollback error: %v, restore error: %v", err, restoreErr),
				Severity: ErrorSeverityCritical,
			}
		}
		return &APIError{
			Type:    ErrorTypeFileSystem,
			Code:    "ROLLBACK_COPY_FAILED",
			Message: fmt.Sprintf("Failed to copy backup to current location: %v", err),
		}
	}

	// Restore original permissions
	if err := os.Chmod(currentExePath, stat.Mode()); err != nil {
		return &APIError{
			Type:    ErrorTypeFileSystem,
			Code:    "ROLLBACK_PERMISSIONS_FAILED",
			Message: fmt.Sprintf("Rollback succeeded but failed to set permissions: %v", err),
			Severity: ErrorSeverityMedium,
		}
	}

	// Clean up failed version backup (keep it for forensics but rename)
	forensicPath := filepath.Join(rm.backupDirectory, fmt.Sprintf("failed_update_%s.forensic", time.Now().Format("20060102-150405")))
	os.Rename(failedBackupPath, forensicPath)

	return nil
}

// ListAvailableBackups returns a list of available backups
func (rm *RollbackManager) ListAvailableBackups() ([]BackupInfo, error) {
	backups := make([]BackupInfo, 0)

	// Read backup directory
	files, err := os.ReadDir(rm.backupDirectory)
	if err != nil {
		if os.IsNotExist(err) {
			return backups, nil // Return empty list if directory doesn't exist
		}
		return nil, &APIError{
			Type:    ErrorTypeFileSystem,
			Code:    "BACKUP_DIR_READ_FAILED",
			Message: fmt.Sprintf("Failed to read backup directory: %v", err),
		}
	}

	// Process backup files
	for _, file := range files {
		if file.IsDir() {
			continue
		}

		if filepath.Ext(file.Name()) != ".backup" {
			continue
		}

		filePath := filepath.Join(rm.backupDirectory, file.Name())
		stat, err := os.Stat(filePath)
		if err != nil {
			continue
		}

		backup := BackupInfo{
			Path:      filePath,
			Name:      file.Name(),
			Size:      stat.Size(),
			CreatedAt: stat.ModTime(),
		}

		// Try to extract version from filename
		// Expected format: mcpweaver_VERSION_TIMESTAMP.backup
		if version := rm.extractVersionFromFilename(file.Name()); version != "" {
			backup.Version = version
		}

		backups = append(backups, backup)
	}

	return backups, nil
}

// DeleteBackup deletes a specific backup file
func (rm *RollbackManager) DeleteBackup(backupPath string) error {
	if !rm.fileExists(backupPath) {
		return &APIError{
			Type:    ErrorTypeFileSystem,
			Code:    "BACKUP_FILE_NOT_FOUND",
			Message: fmt.Sprintf("Backup file not found: %s", backupPath),
		}
	}

	err := os.Remove(backupPath)
	if err != nil {
		return &APIError{
			Type:    ErrorTypeFileSystem,
			Code:    "BACKUP_DELETE_FAILED",
			Message: fmt.Sprintf("Failed to delete backup: %v", err),
		}
	}

	return nil
}

// ValidateBackup validates that a backup file is intact and usable
func (rm *RollbackManager) ValidateBackup(backupPath string) (*BackupValidation, error) {
	validation := &BackupValidation{
		Path:        backupPath,
		ValidatedAt: time.Now(),
	}

	// Check if file exists
	if !rm.fileExists(backupPath) {
		validation.Valid = false
		validation.Errors = append(validation.Errors, "Backup file does not exist")
		return validation, nil
	}

	// Check file size
	stat, err := os.Stat(backupPath)
	if err != nil {
		validation.Valid = false
		validation.Errors = append(validation.Errors, fmt.Sprintf("Failed to stat backup file: %v", err))
		return validation, nil
	}

	if stat.Size() == 0 {
		validation.Valid = false
		validation.Errors = append(validation.Errors, "Backup file is empty")
		return validation, nil
	}

	validation.Size = stat.Size()

	// Check if file is executable (basic validation)
	if stat.Mode()&0111 == 0 {
		validation.Warnings = append(validation.Warnings, "Backup file is not executable")
	}

	// TODO: Add more sophisticated validation (checksum, signature, etc.)
	validation.Valid = len(validation.Errors) == 0

	return validation, nil
}

// GetRollbackCapabilities returns information about rollback capabilities
func (rm *RollbackManager) GetRollbackCapabilities() *RollbackCapabilities {
	backups, _ := rm.ListAvailableBackups()
	
	return &RollbackCapabilities{
		Available:    len(backups) > 0,
		BackupCount:  len(backups),
		MaxBackups:   rm.maxBackups,
		BackupDir:    rm.backupDirectory,
		Features: []string{
			"automatic_backup",
			"multiple_versions",
			"backup_validation",
			"forensic_preservation",
		},
	}
}

// Private helper methods

func (rm *RollbackManager) copyFile(src, dst string) error {
	sourceFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer sourceFile.Close()

	destFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer destFile.Close()

	_, err = destFile.ReadFrom(sourceFile)
	return err
}

func (rm *RollbackManager) fileExists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}

func (rm *RollbackManager) cleanupOldBackups() {
	backups, err := rm.ListAvailableBackups()
	if err != nil || len(backups) <= rm.maxBackups {
		return
	}

	// Sort by creation time (oldest first)
	// Simple bubble sort for small arrays
	for i := 0; i < len(backups)-1; i++ {
		for j := 0; j < len(backups)-i-1; j++ {
			if backups[j].CreatedAt.After(backups[j+1].CreatedAt) {
				backups[j], backups[j+1] = backups[j+1], backups[j]
			}
		}
	}

	// Delete oldest backups
	excessCount := len(backups) - rm.maxBackups
	for i := 0; i < excessCount; i++ {
		os.Remove(backups[i].Path)
	}
}

func (rm *RollbackManager) extractVersionFromFilename(filename string) string {
	// TODO: Implement version extraction from filename
	// Expected format: mcpweaver_VERSION_TIMESTAMP.backup
	return ""
}

// Supporting types for rollback functionality

// BackupInfo contains information about a backup file
type BackupInfo struct {
	Path      string    `json:"path"`
	Name      string    `json:"name"`
	Version   string    `json:"version"`
	Size      int64     `json:"size"`
	CreatedAt time.Time `json:"createdAt"`
}

// BackupValidation contains the result of validating a backup
type BackupValidation struct {
	Path        string    `json:"path"`
	Valid       bool      `json:"valid"`
	Size        int64     `json:"size"`
	Errors      []string  `json:"errors"`
	Warnings    []string  `json:"warnings"`
	ValidatedAt time.Time `json:"validatedAt"`
}

// RollbackCapabilities describes the rollback capabilities of the system
type RollbackCapabilities struct {
	Available   bool     `json:"available"`
	BackupCount int      `json:"backupCount"`
	MaxBackups  int      `json:"maxBackups"`
	BackupDir   string   `json:"backupDir"`
	Features    []string `json:"features"`
}