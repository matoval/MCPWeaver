package app

import (
	"os"
	"path/filepath"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// SelectFile opens a file selection dialog
func (a *App) SelectFile(filters []FileFilter) (string, error) {
	if a.ctx == nil {
		return "", a.createAPIError("internal", ErrCodeInternalError, "Application context not initialized", nil)
	}

	// Convert filters to Wails format
	wailsFilters := make([]runtime.FileFilter, len(filters))
	for i, filter := range filters {
		wailsFilters[i] = runtime.FileFilter{
			DisplayName: filter.DisplayName,
			Pattern:     filter.Pattern,
		}
	}

	// Open file dialog
	filePath, err := runtime.OpenFileDialog(a.ctx, runtime.OpenDialogOptions{
		Title:   "Select OpenAPI Specification",
		Filters: wailsFilters,
	})

	if err != nil {
		return "", a.createAPIError("file_system", ErrCodeFileAccess, "Failed to open file dialog", map[string]string{
			"error": err.Error(),
		})
	}

	// Empty string means user cancelled
	if filePath == "" {
		return "", nil
	}

	// Verify file exists and is readable
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return "", a.createAPIError("file_system", ErrCodeFileAccess, "Selected file does not exist", map[string]string{
			"path": filePath,
		})
	}

	return filePath, nil
}

// SelectDirectory opens a directory selection dialog
func (a *App) SelectDirectory(title string) (string, error) {
	if a.ctx == nil {
		return "", a.createAPIError("internal", ErrCodeInternalError, "Application context not initialized", nil)
	}

	if title == "" {
		title = "Select Directory"
	}

	// Open directory dialog
	dirPath, err := runtime.OpenDirectoryDialog(a.ctx, runtime.OpenDialogOptions{
		Title: title,
	})

	if err != nil {
		return "", a.createAPIError("file_system", ErrCodeFileAccess, "Failed to open directory dialog", map[string]string{
			"error": err.Error(),
		})
	}

	// Empty string means user cancelled
	if dirPath == "" {
		return "", nil
	}

	// Verify directory exists and is writable
	if _, err := os.Stat(dirPath); os.IsNotExist(err) {
		return "", a.createAPIError("file_system", ErrCodeFileAccess, "Selected directory does not exist", map[string]string{
			"path": dirPath,
		})
	}

	// Test if directory is writable
	testFile := filepath.Join(dirPath, ".mcpweaver_test")
	if err := os.WriteFile(testFile, []byte("test"), 0644); err != nil {
		return "", a.createAPIError("file_system", ErrCodeFileAccess, "Directory is not writable", map[string]string{
			"path": dirPath,
			"error": err.Error(),
		})
	}
	os.Remove(testFile) // Clean up test file

	return dirPath, nil
}

// SaveFile opens a save file dialog
func (a *App) SaveFile(content string, defaultPath string, filters []FileFilter) (string, error) {
	if a.ctx == nil {
		return "", a.createAPIError("internal", ErrCodeInternalError, "Application context not initialized", nil)
	}

	// Convert filters to Wails format
	wailsFilters := make([]runtime.FileFilter, len(filters))
	for i, filter := range filters {
		wailsFilters[i] = runtime.FileFilter{
			DisplayName: filter.DisplayName,
			Pattern:     filter.Pattern,
		}
	}

	// Open save dialog
	filePath, err := runtime.SaveFileDialog(a.ctx, runtime.SaveDialogOptions{
		Title:           "Save File",
		DefaultFilename: defaultPath,
		Filters:         wailsFilters,
	})

	if err != nil {
		return "", a.createAPIError("file_system", ErrCodeFileAccess, "Failed to open save dialog", map[string]string{
			"error": err.Error(),
		})
	}

	// Empty string means user cancelled
	if filePath == "" {
		return "", nil
	}

	// Ensure directory exists
	dir := filepath.Dir(filePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return "", a.createAPIError("file_system", ErrCodeFileAccess, "Failed to create directory", map[string]string{
			"path": dir,
			"error": err.Error(),
		})
	}

	// Write content to file
	if err := os.WriteFile(filePath, []byte(content), 0644); err != nil {
		return "", a.createAPIError("file_system", ErrCodeFileAccess, "Failed to save file", map[string]string{
			"path": filePath,
			"error": err.Error(),
		})
	}

	return filePath, nil
}

// ReadFile reads the content of a file
func (a *App) ReadFile(path string) (string, error) {
	if path == "" {
		return "", a.createAPIError("validation", ErrCodeValidation, "File path is required", nil)
	}

	// Check if file exists
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return "", a.createAPIError("file_system", ErrCodeFileAccess, "File does not exist", map[string]string{
			"path": path,
		})
	}

	// Read file content
	content, err := os.ReadFile(path)
	if err != nil {
		return "", a.createAPIError("file_system", ErrCodeFileAccess, "Failed to read file", map[string]string{
			"path": path,
			"error": err.Error(),
		})
	}

	return string(content), nil
}

// WriteFile writes content to a file
func (a *App) WriteFile(path string, content string) error {
	if path == "" {
		return a.createAPIError("validation", ErrCodeValidation, "File path is required", nil)
	}

	// Ensure directory exists
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return a.createAPIError("file_system", ErrCodeFileAccess, "Failed to create directory", map[string]string{
			"path": dir,
			"error": err.Error(),
		})
	}

	// Write content to file
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		return a.createAPIError("file_system", ErrCodeFileAccess, "Failed to write file", map[string]string{
			"path": path,
			"error": err.Error(),
		})
	}

	return nil
}

// FileExists checks if a file exists
func (a *App) FileExists(path string) (bool, error) {
	if path == "" {
		return false, a.createAPIError("validation", ErrCodeValidation, "File path is required", nil)
	}

	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	
	// Other error occurred
	return false, a.createAPIError("file_system", ErrCodeFileAccess, "Failed to check file existence", map[string]string{
		"path": path,
		"error": err.Error(),
	})
}

// GetDefaultOpenAPIFilters returns default file filters for OpenAPI specifications
func (a *App) GetDefaultOpenAPIFilters() []FileFilter {
	return []FileFilter{
		{
			DisplayName: "OpenAPI Specifications",
			Pattern:     "*.json;*.yaml;*.yml",
			Extensions:  []string{".json", ".yaml", ".yml"},
		},
		{
			DisplayName: "JSON Files",
			Pattern:     "*.json",
			Extensions:  []string{".json"},
		},
		{
			DisplayName: "YAML Files",
			Pattern:     "*.yaml;*.yml",
			Extensions:  []string{".yaml", ".yml"},
		},
		{
			DisplayName: "All Files",
			Pattern:     "*.*",
			Extensions:  []string{"*"},
		},
	}
}