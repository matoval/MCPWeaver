package app

import (
	"crypto/rand"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// SecureFileSystem provides secure file system operations with additional security controls
type SecureFileSystem struct {
	app          *App
	validator    *SecurityValidator
	allowedPaths []string
	maxFileSize  int64
}

// NewSecureFileSystem creates a new secure file system handler
func NewSecureFileSystem(app *App) *SecureFileSystem {
	return &SecureFileSystem{
		app:       app,
		validator: NewSecurityValidator(app),
		allowedPaths: []string{
			"./",         // Current directory
			"output/",    // Output directory
			"temp/",      // Temporary files
			"projects/",  // Project files
			"templates/", // Template files
		},
		maxFileSize: 10 * 1024 * 1024, // 10MB default
	}
}

// SecureReadFile reads a file with security validation
func (sfs *SecureFileSystem) SecureReadFile(path string) ([]byte, error) {
	// Validate the file path
	if err := sfs.validatePath(path); err != nil {
		return nil, err
	}

	// Check if file exists and get info
	fileInfo, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, sfs.app.createAPIError("file_system", "FILE_NOT_FOUND", "File does not exist", map[string]string{
				"path": path,
			})
		}
		return nil, sfs.app.createAPIError("file_system", "FILE_ACCESS_ERROR", "Cannot access file", map[string]string{
			"path":  path,
			"error": err.Error(),
		})
	}

	// Check if it's a regular file
	if !fileInfo.Mode().IsRegular() {
		return nil, sfs.app.createAPIError("file_system", "NOT_REGULAR_FILE", "Path does not point to a regular file", map[string]string{
			"path": path,
			"mode": fileInfo.Mode().String(),
		})
	}

	// Check file size
	if fileInfo.Size() > sfs.maxFileSize {
		return nil, sfs.app.createAPIError("file_system", "FILE_TOO_LARGE", "File exceeds maximum allowed size", map[string]string{
			"path":    path,
			"size":    fmt.Sprintf("%d", fileInfo.Size()),
			"maxSize": fmt.Sprintf("%d", sfs.maxFileSize),
		})
	}

	// Open file with limited reader
	file, err := os.Open(path)
	if err != nil {
		return nil, sfs.app.createAPIError("file_system", "FILE_OPEN_ERROR", "Cannot open file", map[string]string{
			"path":  path,
			"error": err.Error(),
		})
	}
	defer file.Close()

	// Read with size limit
	limitedReader := io.LimitReader(file, sfs.maxFileSize)
	content, err := io.ReadAll(limitedReader)
	if err != nil {
		return nil, sfs.app.createAPIError("file_system", "FILE_READ_ERROR", "Cannot read file", map[string]string{
			"path":  path,
			"error": err.Error(),
		})
	}

	return content, nil
}

// SecureWriteFile writes a file with security validation and atomic operations
func (sfs *SecureFileSystem) SecureWriteFile(path string, content []byte, perm os.FileMode) error {
	// Validate the file path
	if err := sfs.validatePath(path); err != nil {
		return err
	}

	// Validate content size
	if int64(len(content)) > sfs.maxFileSize {
		return sfs.app.createAPIError("file_system", "CONTENT_TOO_LARGE", "Content exceeds maximum allowed size", map[string]string{
			"path":    path,
			"size":    fmt.Sprintf("%d", len(content)),
			"maxSize": fmt.Sprintf("%d", sfs.maxFileSize),
		})
	}

	// Ensure directory exists
	dir := filepath.Dir(path)
	if err := sfs.ensureSecureDirectory(dir, 0755); err != nil {
		return err
	}

	// Write atomically using temporary file
	return sfs.atomicWrite(path, content, perm)
}

// SecureCreateTempFile creates a temporary file with secure settings
func (sfs *SecureFileSystem) SecureCreateTempFile(dir, pattern string) (*os.File, error) {
	// Validate directory path if provided
	if dir != "" {
		if err := sfs.validatePath(dir); err != nil {
			return nil, err
		}
	} else {
		dir = os.TempDir()
	}

	// Ensure the temp directory exists and is secure
	if err := sfs.ensureSecureDirectory(dir, 0700); err != nil {
		return nil, err
	}

	// Generate a cryptographically secure random suffix
	randomBytes := make([]byte, 8)
	if _, err := rand.Read(randomBytes); err != nil {
		return nil, sfs.app.createAPIError("file_system", "RANDOM_GENERATION_ERROR", "Cannot generate secure random name", map[string]string{
			"error": err.Error(),
		})
	}

	// Create the temp file with secure permissions
	tempName := fmt.Sprintf("%s_%x_%d", pattern, randomBytes, time.Now().UnixNano())
	tempPath := filepath.Join(dir, tempName)

	file, err := os.OpenFile(tempPath, os.O_CREATE|os.O_WRONLY|os.O_EXCL, 0600)
	if err != nil {
		return nil, sfs.app.createAPIError("file_system", "TEMP_FILE_CREATE_ERROR", "Cannot create temporary file", map[string]string{
			"path":  tempPath,
			"error": err.Error(),
		})
	}

	return file, nil
}

// SecureDelete deletes a file with security validation
func (sfs *SecureFileSystem) SecureDelete(path string) error {
	// Validate the file path
	if err := sfs.validatePath(path); err != nil {
		return err
	}

	// Check if file exists
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return nil // File doesn't exist, consider it deleted
	}

	// Secure delete by overwriting with random data first (optional, for sensitive files)
	if err := sfs.secureOverwrite(path); err != nil {
		// If secure overwrite fails, log warning but continue with normal delete
		if sfs.app.ctx != nil {
			sfs.app.emitNotification("warning", "Secure Delete", "Could not securely overwrite file before deletion")
		}
	}

	// Delete the file
	if err := os.Remove(path); err != nil {
		return sfs.app.createAPIError("file_system", "FILE_DELETE_ERROR", "Cannot delete file", map[string]string{
			"path":  path,
			"error": err.Error(),
		})
	}

	return nil
}

// SecureCopyFile copies a file with security validation
func (sfs *SecureFileSystem) SecureCopyFile(src, dst string) error {
	// Validate both paths
	if err := sfs.validatePath(src); err != nil {
		return err
	}
	if err := sfs.validatePath(dst); err != nil {
		return err
	}

	// Read source file securely
	content, err := sfs.SecureReadFile(src)
	if err != nil {
		return err
	}

	// Get source file permissions
	srcInfo, err := os.Stat(src)
	if err != nil {
		return sfs.app.createAPIError("file_system", "SOURCE_STAT_ERROR", "Cannot get source file information", map[string]string{
			"path":  src,
			"error": err.Error(),
		})
	}

	// Write to destination securely
	return sfs.SecureWriteFile(dst, content, srcInfo.Mode())
}

// validatePath validates a file path for security
func (sfs *SecureFileSystem) validatePath(path string) error {
	// Use the security validator
	if err := sfs.validator.ValidateFilePath(path); err != nil {
		return err
	}

	// Clean the path
	cleanPath := filepath.Clean(path)

	// Check if path is within allowed directories
	allowed := false
	absPath, err := filepath.Abs(cleanPath)
	if err != nil {
		return sfs.app.createAPIError("file_system", "PATH_RESOLUTION_ERROR", "Cannot resolve absolute path", map[string]string{
			"path":  path,
			"error": err.Error(),
		})
	}

	for _, allowedPath := range sfs.allowedPaths {
		absAllowed, err := filepath.Abs(allowedPath)
		if err != nil {
			continue
		}

		// Check if the path is within the allowed directory
		if strings.HasPrefix(absPath, absAllowed) {
			allowed = true
			break
		}
	}

	if !allowed {
		return sfs.app.createAPIError("security", "PATH_NOT_ALLOWED", "File path is not within allowed directories", map[string]string{
			"path":         path,
			"absolutePath": absPath,
		})
	}

	return nil
}

// ensureSecureDirectory creates a directory with secure permissions
func (sfs *SecureFileSystem) ensureSecureDirectory(path string, perm os.FileMode) error {
	// Check if directory already exists
	if info, err := os.Stat(path); err == nil {
		if !info.IsDir() {
			return sfs.app.createAPIError("file_system", "NOT_A_DIRECTORY", "Path exists but is not a directory", map[string]string{
				"path": path,
			})
		}

		// Check permissions
		if info.Mode().Perm() != perm {
			if err := os.Chmod(path, perm); err != nil {
				return sfs.app.createAPIError("file_system", "CHMOD_ERROR", "Cannot set directory permissions", map[string]string{
					"path":  path,
					"error": err.Error(),
				})
			}
		}
		return nil
	}

	// Create directory with secure permissions
	if err := os.MkdirAll(path, perm); err != nil {
		return sfs.app.createAPIError("file_system", "MKDIR_ERROR", "Cannot create directory", map[string]string{
			"path":  path,
			"error": err.Error(),
		})
	}

	return nil
}

// atomicWrite performs an atomic file write operation
func (sfs *SecureFileSystem) atomicWrite(path string, content []byte, perm os.FileMode) error {
	dir := filepath.Dir(path)
	base := filepath.Base(path)

	// Create temporary file in the same directory
	tempFile, err := sfs.SecureCreateTempFile(dir, ".tmp_"+base)
	if err != nil {
		return err
	}
	tempPath := tempFile.Name()

	// Ensure cleanup on failure
	defer func() {
		tempFile.Close()
		os.Remove(tempPath)
	}()

	// Write content to temporary file
	if _, err := tempFile.Write(content); err != nil {
		return sfs.app.createAPIError("file_system", "TEMP_WRITE_ERROR", "Cannot write to temporary file", map[string]string{
			"tempPath": tempPath,
			"error":    err.Error(),
		})
	}

	// Sync to ensure data is written to disk
	if err := tempFile.Sync(); err != nil {
		return sfs.app.createAPIError("file_system", "SYNC_ERROR", "Cannot sync temporary file", map[string]string{
			"tempPath": tempPath,
			"error":    err.Error(),
		})
	}

	// Close the temporary file
	if err := tempFile.Close(); err != nil {
		return sfs.app.createAPIError("file_system", "TEMP_CLOSE_ERROR", "Cannot close temporary file", map[string]string{
			"tempPath": tempPath,
			"error":    err.Error(),
		})
	}

	// Set the correct permissions
	if err := os.Chmod(tempPath, perm); err != nil {
		return sfs.app.createAPIError("file_system", "CHMOD_TEMP_ERROR", "Cannot set temporary file permissions", map[string]string{
			"tempPath": tempPath,
			"error":    err.Error(),
		})
	}

	// Atomically move the temporary file to the final location
	if err := os.Rename(tempPath, path); err != nil {
		return sfs.app.createAPIError("file_system", "ATOMIC_MOVE_ERROR", "Cannot atomically move file", map[string]string{
			"tempPath":  tempPath,
			"finalPath": path,
			"error":     err.Error(),
		})
	}

	return nil
}

// secureOverwrite overwrites a file with random data before deletion
func (sfs *SecureFileSystem) secureOverwrite(path string) error {
	// Get file size
	fileInfo, err := os.Stat(path)
	if err != nil {
		return err
	}

	// Don't overwrite very large files (performance consideration)
	if fileInfo.Size() > 1024*1024 { // 1MB limit for secure overwrite
		return nil
	}

	// Open file for writing
	file, err := os.OpenFile(path, os.O_WRONLY, 0)
	if err != nil {
		return err
	}
	defer file.Close()

	// Overwrite with random data
	randomData := make([]byte, fileInfo.Size())
	if _, err := rand.Read(randomData); err != nil {
		return err
	}

	// Write random data
	if _, err := file.WriteAt(randomData, 0); err != nil {
		return err
	}

	// Sync to ensure data is written
	return file.Sync()
}

// GetFilePermissions returns secure file permissions based on file type
func (sfs *SecureFileSystem) GetFilePermissions(fileType string) os.FileMode {
	switch fileType {
	case "config", "sensitive":
		return 0600 // Owner read/write only
	case "executable":
		return 0755 // Owner read/write/execute, group/others read/execute
	case "public":
		return 0644 // Owner read/write, group/others read
	case "directory":
		return 0755 // Owner read/write/execute, group/others read/execute
	case "temp":
		return 0600 // Owner read/write only
	default:
		return 0644 // Default to public read
	}
}

// CheckFileIntegrity performs basic integrity checks on a file
func (sfs *SecureFileSystem) CheckFileIntegrity(path string) error {
	// Check if file exists and is accessible
	fileInfo, err := os.Stat(path)
	if err != nil {
		return sfs.app.createAPIError("file_system", "INTEGRITY_CHECK_ERROR", "Cannot access file for integrity check", map[string]string{
			"path":  path,
			"error": err.Error(),
		})
	}

	// Check if file size is reasonable
	if fileInfo.Size() == 0 {
		return sfs.app.createAPIError("file_system", "EMPTY_FILE", "File is empty", map[string]string{
			"path": path,
		})
	}

	// Check file permissions
	mode := fileInfo.Mode()
	if mode&os.ModeType != 0 && !mode.IsRegular() {
		return sfs.app.createAPIError("file_system", "INVALID_FILE_TYPE", "File is not a regular file", map[string]string{
			"path": path,
			"mode": mode.String(),
		})
	}

	// Additional checks could include checksum verification, etc.

	return nil
}

// SetAllowedPaths updates the list of allowed paths for file operations
func (sfs *SecureFileSystem) SetAllowedPaths(paths []string) error {
	// Validate each path
	for _, path := range paths {
		if err := sfs.validator.ValidateFilePath(path); err != nil {
			return sfs.app.createAPIError("validation", "INVALID_ALLOWED_PATH", "Invalid path in allowed paths list", map[string]string{
				"path":  path,
				"error": err.Error(),
			})
		}
	}

	sfs.allowedPaths = paths
	return nil
}

// AddAllowedPath adds a new allowed path
func (sfs *SecureFileSystem) AddAllowedPath(path string) error {
	if err := sfs.validator.ValidateFilePath(path); err != nil {
		return err
	}

	// Check if path already exists
	for _, existing := range sfs.allowedPaths {
		if existing == path {
			return nil // Already exists
		}
	}

	sfs.allowedPaths = append(sfs.allowedPaths, path)
	return nil
}

// IsPathAllowed checks if a path is within allowed directories
func (sfs *SecureFileSystem) IsPathAllowed(path string) bool {
	err := sfs.validatePath(path)
	return err == nil
}
