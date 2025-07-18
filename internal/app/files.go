package app

import (
	"encoding/json"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/wailsapp/wails/v2/pkg/runtime"
	"gopkg.in/yaml.v3"
)

// fileExists checks if a file exists and returns appropriate error
func (a *App) fileExists(path string) error {
	if path == "" {
		return a.createAPIError("validation", ErrCodeValidation, "File path is required", nil)
	}

	if _, err := os.Stat(path); os.IsNotExist(err) {
		return a.createAPIError("file_system", ErrCodeFileAccess, "File does not exist", map[string]string{
			"path": path,
		})
	} else if err != nil {
		return a.createAPIError("file_system", ErrCodeFileAccess, "Failed to check file existence", map[string]string{
			"path": path,
			"error": err.Error(),
		})
	}

	return nil
}

// dirExists checks if a directory exists and is writable
func (a *App) dirExists(path string) error {
	if path == "" {
		return a.createAPIError("validation", ErrCodeValidation, "Directory path is required", nil)
	}

	if _, err := os.Stat(path); os.IsNotExist(err) {
		return a.createAPIError("file_system", ErrCodeFileAccess, "Directory does not exist", map[string]string{
			"path": path,
		})
	}

	// Test if directory is writable
	testFile := filepath.Join(path, ".mcpweaver_test")
	if err := os.WriteFile(testFile, []byte("test"), 0644); err != nil {
		return a.createAPIError("file_system", ErrCodeFileAccess, "Directory is not writable", map[string]string{
			"path": path,
			"error": err.Error(),
		})
	}
	os.Remove(testFile) // Clean up test file

	return nil
}

// ensureDir ensures a directory exists, creating it if necessary
func (a *App) ensureDir(path string) error {
	if path == "" {
		return a.createAPIError("validation", ErrCodeValidation, "Directory path is required", nil)
	}

	if err := os.MkdirAll(path, 0755); err != nil {
		return a.createAPIError("file_system", ErrCodeFileAccess, "Failed to create directory", map[string]string{
			"path": path,
			"error": err.Error(),
		})
	}

	return nil
}

// convertFilters converts FileFilter slice to Wails runtime.FileFilter slice
func convertFilters(filters []FileFilter) []runtime.FileFilter {
	wailsFilters := make([]runtime.FileFilter, len(filters))
	for i, filter := range filters {
		wailsFilters[i] = runtime.FileFilter{
			DisplayName: filter.DisplayName,
			Pattern:     filter.Pattern,
		}
	}
	return wailsFilters
}

// SelectFile opens a file selection dialog
func (a *App) SelectFile(filters []FileFilter) (string, error) {
	if a.ctx == nil {
		return "", a.createAPIError("internal", ErrCodeInternalError, "Application context not initialized", nil)
	}

	// Open file dialog
	filePath, err := runtime.OpenFileDialog(a.ctx, runtime.OpenDialogOptions{
		Title:   "Select OpenAPI Specification",
		Filters: convertFilters(filters),
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
	if err := a.fileExists(filePath); err != nil {
		return "", err
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
	if err := a.dirExists(dirPath); err != nil {
		return "", err
	}

	return dirPath, nil
}

// SaveFile opens a save file dialog
func (a *App) SaveFile(content string, defaultPath string, filters []FileFilter) (string, error) {
	if a.ctx == nil {
		return "", a.createAPIError("internal", ErrCodeInternalError, "Application context not initialized", nil)
	}

	// Open save dialog
	filePath, err := runtime.SaveFileDialog(a.ctx, runtime.SaveDialogOptions{
		Title:           "Save File",
		DefaultFilename: defaultPath,
		Filters:         convertFilters(filters),
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
	if err := a.ensureDir(dir); err != nil {
		return "", err
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
	// Check if file exists
	if err := a.fileExists(path); err != nil {
		return "", err
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
	if err := a.ensureDir(dir); err != nil {
		return err
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
			Pattern:     "*.json;*.yaml;*.yml;*.openapi",
			Extensions:  []string{".json", ".yaml", ".yml", ".openapi"},
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
			DisplayName: "OpenAPI Files",
			Pattern:     "*.openapi",
			Extensions:  []string{".openapi"},
		},
		{
			DisplayName: "All Files",
			Pattern:     "*.*",
			Extensions:  []string{"*"},
		},
	}
}

// GetSupportedFileFormats returns all supported file formats for OpenAPI specs
func (a *App) GetSupportedFileFormats() []string {
	return []string{
		"application/json",
		"application/x-yaml",
		"text/yaml",
		"text/x-yaml",
		"application/yaml",
		"text/plain",
	}
}

// DetectFileFormat detects the format of a file based on content and extension
func (a *App) DetectFileFormat(content string, filename string) (string, error) {
	// First, try to detect based on content
	content = strings.TrimSpace(content)
	
	if len(content) == 0 {
		return "", a.createAPIError("validation", ErrCodeValidation, "File is empty", nil)
	}

	// Check if it's JSON
	if strings.HasPrefix(content, "{") && strings.HasSuffix(content, "}") {
		var jsonData map[string]interface{}
		if err := json.Unmarshal([]byte(content), &jsonData); err == nil {
			return "json", nil
		}
	}

	// Check if it's YAML
	if !strings.HasPrefix(content, "{") {
		var yamlData map[string]interface{}
		if err := yaml.Unmarshal([]byte(content), &yamlData); err == nil {
			return "yaml", nil
		}
	}

	// Fall back to extension-based detection
	ext := strings.ToLower(filepath.Ext(filename))
	switch ext {
	case ".json":
		return "json", nil
	case ".yaml", ".yml":
		return "yaml", nil
	case ".openapi":
		// OpenAPI files can be either JSON or YAML
		if strings.HasPrefix(content, "{") {
			return "json", nil
		}
		return "yaml", nil
	default:
		// Try to parse as JSON first, then YAML
		var jsonData map[string]interface{}
		if err := json.Unmarshal([]byte(content), &jsonData); err == nil {
			return "json", nil
		}
		
		var yamlData map[string]interface{}
		if err := yaml.Unmarshal([]byte(content), &yamlData); err == nil {
			return "yaml", nil
		}
		
		return "", a.createAPIError("validation", ErrCodeValidation, "Unable to detect file format", map[string]string{
			"filename": filename,
			"extension": ext,
		})
	}
}

// ImportOpenAPISpec imports an OpenAPI specification from a file
func (a *App) ImportOpenAPISpec(filePath string) (*ImportResult, error) {
	if filePath == "" {
		return nil, a.createAPIError("validation", ErrCodeValidation, "File path is required", nil)
	}

	// Check if file exists
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return nil, a.createAPIError("file_system", ErrCodeFileAccess, "File does not exist", map[string]string{
			"path": filePath,
		})
	}

	// Validate file extension
	if !a.isValidOpenAPIFile(filePath) {
		return nil, a.createAPIError("validation", ErrCodeValidation, "Invalid file format. Please select a JSON or YAML file", map[string]string{
			"path": filePath,
		})
	}

	// Read file content
	content, err := os.ReadFile(filePath)
	if err != nil {
		return nil, a.createAPIError("file_system", ErrCodeFileAccess, "Failed to read file", map[string]string{
			"path": filePath,
			"error": err.Error(),
		})
	}

	// Get file info
	fileInfo, err := os.Stat(filePath)
	if err != nil {
		return nil, a.createAPIError("file_system", ErrCodeFileAccess, "Failed to get file info", map[string]string{
			"path": filePath,
			"error": err.Error(),
		})
	}

	// Parse and validate the OpenAPI spec
	result, err := a.parseOpenAPIContent(string(content), filePath)
	if err != nil {
		return nil, err
	}

	result.ImportedFrom = "file"
	result.FilePath = filePath
	result.FileSize = fileInfo.Size()
	result.ImportedAt = time.Now()

	// Add to recent files
	if result.Valid {
		a.AddRecentFile(filePath, "spec")
	}

	return result, nil
}

// ImportOpenAPISpecFromURL imports an OpenAPI specification from a URL
func (a *App) ImportOpenAPISpecFromURL(url string) (*ImportResult, error) {
	if url == "" {
		return nil, a.createAPIError("validation", ErrCodeValidation, "URL is required", nil)
	}

	// Validate URL format
	if !strings.HasPrefix(url, "http://") && !strings.HasPrefix(url, "https://") {
		return nil, a.createAPIError("validation", ErrCodeValidation, "Invalid URL format. URL must start with http:// or https://", map[string]string{
			"url": url,
		})
	}

	// Create progress tracking
	progressID := a.generateProgressID()
	a.emitFileProgress(progressID, "import", 0, "Starting URL import...", 1, 0)

	// Create HTTP client with timeout
	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	// Make HTTP request
	a.emitFileProgress(progressID, "import", 25, "Connecting to URL...", 1, 0)
	resp, err := client.Get(url)
	if err != nil {
		a.emitFileProgress(progressID, "import", 0, "Failed to connect", 1, 0)
		return nil, a.createAPIError("network", ErrCodeNetworkError, "Failed to fetch URL", map[string]string{
			"url": url,
			"error": err.Error(),
		})
	}
	defer resp.Body.Close()

	// Check response status
	if resp.StatusCode != http.StatusOK {
		a.emitFileProgress(progressID, "import", 0, "HTTP error: "+resp.Status, 1, 0)
		return nil, a.createAPIError("network", ErrCodeNetworkError, "Failed to fetch URL", map[string]string{
			"url": url,
			"status": resp.Status,
		})
	}

	// Read response body with progress tracking
	a.emitFileProgress(progressID, "import", 50, "Downloading content...", 1, 0)
	content, err := io.ReadAll(resp.Body)
	if err != nil {
		a.emitFileProgress(progressID, "import", 0, "Failed to read response", 1, 0)
		return nil, a.createAPIError("network", ErrCodeNetworkError, "Failed to read response", map[string]string{
			"url": url,
			"error": err.Error(),
		})
	}

	// Parse and validate the OpenAPI spec
	a.emitFileProgress(progressID, "import", 75, "Parsing and validating content...", 1, 0)
	result, err := a.parseOpenAPIContent(string(content), url)
	if err != nil {
		a.emitFileProgress(progressID, "import", 0, "Validation failed", 1, 0)
		return nil, err
	}

	result.ImportedFrom = "url"
	result.SourceURL = url
	result.FileSize = int64(len(content))
	result.ImportedAt = time.Now()

	a.emitFileProgress(progressID, "import", 100, "Import completed successfully", 1, 1)
	return result, nil
}

// ExportGeneratedServer exports a generated server to a specified directory
func (a *App) ExportGeneratedServer(projectID, targetDir string) (*ExportResult, error) {
	if projectID == "" {
		return nil, a.createAPIError("validation", ErrCodeValidation, "Project ID is required", nil)
	}

	if targetDir == "" {
		return nil, a.createAPIError("validation", ErrCodeValidation, "Target directory is required", nil)
	}

	// Create progress tracking
	progressID := a.generateProgressID()
	a.emitFileProgress(progressID, "export", 0, "Starting export...", 1, 0)

	// Get project
	project, err := a.GetProject(projectID)
	if err != nil {
		a.emitFileProgress(progressID, "export", 0, "Failed to get project", 1, 0)
		return nil, err
	}

	// Check if project has been generated
	if project.Status != ProjectStatusReady || project.GenerationCount == 0 {
		a.emitFileProgress(progressID, "export", 0, "Project not ready for export", 1, 0)
		return nil, a.createAPIError("validation", ErrCodeValidation, "Project has not been generated yet", nil)
	}

	// Create target directory if it doesn't exist
	a.emitFileProgress(progressID, "export", 10, "Creating target directory...", 1, 0)
	if err := os.MkdirAll(targetDir, 0755); err != nil {
		a.emitFileProgress(progressID, "export", 0, "Failed to create directory", 1, 0)
		return nil, a.createAPIError("file_system", ErrCodeFileAccess, "Failed to create target directory", map[string]string{
			"path": targetDir,
			"error": err.Error(),
		})
	}

	// Copy generated files
	sourceDir := project.OutputPath
	copiedFiles := []ExportedFile{}

	// List of files to copy
	filesToCopy := []string{
		"main.go",
		"go.mod",
		"README.md",
		"Dockerfile",
		"Makefile",
		".gitignore",
	}

	totalFiles := len(filesToCopy)
	processedFiles := 0

	for _, fileName := range filesToCopy {
		sourcePath := filepath.Join(sourceDir, fileName)
		targetPath := filepath.Join(targetDir, fileName)

		if _, err := os.Stat(sourcePath); os.IsNotExist(err) {
			processedFiles++
			continue // Skip if file doesn't exist
		}

		// Update progress
		progress := 20 + (processedFiles * 60 / totalFiles)
		a.emitFileProgress(progressID, "export", progress, "Copying "+fileName+"...", totalFiles, processedFiles)

		// Copy file
		if err := a.copyFile(sourcePath, targetPath); err != nil {
			a.emitFileProgress(progressID, "export", 0, "Failed to copy "+fileName, totalFiles, processedFiles)
			return nil, a.createAPIError("file_system", ErrCodeFileAccess, "Failed to copy file", map[string]string{
				"source": sourcePath,
				"target": targetPath,
				"error": err.Error(),
			})
		}

		// Get file info
		fileInfo, err := os.Stat(targetPath)
		if err != nil {
			processedFiles++
			continue
		}

		copiedFiles = append(copiedFiles, ExportedFile{
			Name:         fileName,
			Path:         targetPath,
			Size:         fileInfo.Size(),
			ModifiedTime: fileInfo.ModTime(),
		})
		processedFiles++
	}

	// Copy additional directories if they exist
	a.emitFileProgress(progressID, "export", 80, "Copying additional directories...", totalFiles, processedFiles)
	dirsToCopy := []string{"docs", "examples", "scripts"}
	for _, dirName := range dirsToCopy {
		sourceDirPath := filepath.Join(sourceDir, dirName)
		targetDirPath := filepath.Join(targetDir, dirName)

		if _, err := os.Stat(sourceDirPath); os.IsNotExist(err) {
			continue // Skip if directory doesn't exist
		}

		if err := a.copyDirectory(sourceDirPath, targetDirPath); err != nil {
			// Log error but don't fail the export
			continue
		}
	}

	result := &ExportResult{
		ProjectID:     projectID,
		ProjectName:   project.Name,
		TargetDir:     targetDir,
		ExportedFiles: copiedFiles,
		TotalFiles:    len(copiedFiles),
		TotalSize:     a.calculateTotalSize(copiedFiles),
		ExportedAt:    time.Now(),
	}

	a.emitFileProgress(progressID, "export", 100, "Export completed successfully", totalFiles, processedFiles)
	return result, nil
}

// parseOpenAPIContent parses and validates OpenAPI content
func (a *App) parseOpenAPIContent(content, source string) (*ImportResult, error) {
	result := &ImportResult{
		Content:   content,
		Valid:     false,
		Errors:    []string{},
		Warnings:  []string{},
	}

	// Validate content size
	if len(content) == 0 {
		result.Errors = append(result.Errors, "File is empty")
		return result, nil
	}

	if len(content) > 10*1024*1024 { // 10MB limit
		result.Errors = append(result.Errors, "File is too large (maximum 10MB)")
		return result, nil
	}

	// Basic validation - check if content is valid JSON or YAML
	var specData map[string]interface{}
	var parseError error
	
	// Try parsing as JSON first
	if err := json.Unmarshal([]byte(content), &specData); err != nil {
		// Try parsing as YAML
		if err := yaml.Unmarshal([]byte(content), &specData); err != nil {
			result.Errors = append(result.Errors, "Invalid JSON or YAML format: "+err.Error())
			return result, nil
		}
		parseError = err
	}

	// Check if it's an OpenAPI spec
	var openApiVersion string
	if version, hasOpenAPI := specData["openapi"].(string); hasOpenAPI {
		openApiVersion = version
	} else if version, hasSwagger := specData["swagger"].(string); hasSwagger {
		openApiVersion = version
		result.Warnings = append(result.Warnings, "Using older Swagger 2.0 specification. Consider upgrading to OpenAPI 3.0+")
	} else {
		result.Errors = append(result.Errors, "Not a valid OpenAPI specification. Missing 'openapi' or 'swagger' field")
		return result, nil
	}

	// Validate OpenAPI version
	if openApiVersion != "" {
		if !a.isValidOpenAPIVersion(openApiVersion) {
			result.Warnings = append(result.Warnings, "Unsupported OpenAPI version: "+openApiVersion)
		}
	}

	// Extract and validate basic info
	info := &SpecInfo{
		Version: openApiVersion,
		Title: "Unknown",
		Description: "",
	}

	if infoData, ok := specData["info"].(map[string]interface{}); ok {
		if version, ok := infoData["version"].(string); ok {
			info.Version = version
		} else {
			result.Errors = append(result.Errors, "Missing required 'info.version' field")
			return result, nil
		}
		
		if title, ok := infoData["title"].(string); ok {
			info.Title = title
		} else {
			result.Errors = append(result.Errors, "Missing required 'info.title' field")
			return result, nil
		}
		
		if description, ok := infoData["description"].(string); ok {
			info.Description = description
		}
	} else {
		result.Errors = append(result.Errors, "Missing required 'info' section")
		return result, nil
	}

	// Validate and count operations
	operationCount := 0
	schemaCount := 0
	
	if paths, ok := specData["paths"].(map[string]interface{}); ok {
		if len(paths) == 0 {
			result.Warnings = append(result.Warnings, "No paths defined in the specification")
		}
		
		for pathName, pathData := range paths {
			if pathMethods, ok := pathData.(map[string]interface{}); ok {
				for method := range pathMethods {
					if a.isValidHTTPMethod(method) {
						operationCount++
					}
				}
			}
			
			// Validate path format
			if !strings.HasPrefix(pathName, "/") {
				result.Warnings = append(result.Warnings, "Path '"+pathName+"' should start with '/'")
			}
		}
	} else {
		result.Warnings = append(result.Warnings, "No paths section found")
	}

	// Count schemas/components
	if components, ok := specData["components"].(map[string]interface{}); ok {
		if schemas, ok := components["schemas"].(map[string]interface{}); ok {
			schemaCount = len(schemas)
		}
	} else if definitions, ok := specData["definitions"].(map[string]interface{}); ok {
		// Swagger 2.0 format
		schemaCount = len(definitions)
	}

	info.OperationCount = operationCount
	info.SchemaCount = schemaCount

	// Extract and validate servers
	servers := []ServerInfo{}
	if serverData, ok := specData["servers"].([]interface{}); ok {
		for _, server := range serverData {
			if serverMap, ok := server.(map[string]interface{}); ok {
				serverInfo := ServerInfo{}
				if url, ok := serverMap["url"].(string); ok {
					serverInfo.URL = url
					if !a.isValidServerURL(url) {
						result.Warnings = append(result.Warnings, "Invalid server URL: "+url)
					}
				}
				if desc, ok := serverMap["description"].(string); ok {
					serverInfo.Description = desc
				}
				servers = append(servers, serverInfo)
			}
		}
	} else if host, ok := specData["host"].(string); ok {
		// Swagger 2.0 format
		scheme := "https"
		if schemes, ok := specData["schemes"].([]interface{}); ok && len(schemes) > 0 {
			if s, ok := schemes[0].(string); ok {
				scheme = s
			}
		}
		basePath := ""
		if bp, ok := specData["basePath"].(string); ok {
			basePath = bp
		}
		
		serverInfo := ServerInfo{
			URL: scheme + "://" + host + basePath,
			Description: "Generated from Swagger 2.0 host",
		}
		servers = append(servers, serverInfo)
	}
	
	if len(servers) == 0 {
		result.Warnings = append(result.Warnings, "No servers defined in the specification")
	}
	
	info.Servers = servers

	// Extract security schemes
	securitySchemes := []SecurityScheme{}
	if components, ok := specData["components"].(map[string]interface{}); ok {
		if secSchemes, ok := components["securitySchemes"].(map[string]interface{}); ok {
			for name, scheme := range secSchemes {
				if schemeMap, ok := scheme.(map[string]interface{}); ok {
					secScheme := SecurityScheme{Name: name}
					if schemeType, ok := schemeMap["type"].(string); ok {
						secScheme.Type = schemeType
					}
					if desc, ok := schemeMap["description"].(string); ok {
						secScheme.Description = desc
					}
					securitySchemes = append(securitySchemes, secScheme)
				}
			}
		}
	}
	info.SecuritySchemes = securitySchemes

	// Additional validation checks
	if operationCount == 0 {
		result.Warnings = append(result.Warnings, "No operations found in the specification")
	}

	if operationCount > 1000 {
		result.Warnings = append(result.Warnings, "Large specification with "+string(rune(operationCount))+" operations. Generation may take longer.")
	}

	if parseError != nil {
		result.Warnings = append(result.Warnings, "File parsed as YAML but contains JSON syntax errors: "+parseError.Error())
	}

	result.Valid = len(result.Errors) == 0
	result.SpecInfo = info

	return result, nil
}

// isValidOpenAPIVersion checks if the OpenAPI version is supported
func (a *App) isValidOpenAPIVersion(version string) bool {
	supportedVersions := []string{"3.0.0", "3.0.1", "3.0.2", "3.0.3", "3.1.0", "2.0"}
	for _, v := range supportedVersions {
		if strings.HasPrefix(version, v) {
			return true
		}
	}
	return false
}

// isValidHTTPMethod checks if a method is a valid HTTP method
func (a *App) isValidHTTPMethod(method string) bool {
	validMethods := []string{"get", "post", "put", "patch", "delete", "head", "options", "trace"}
	method = strings.ToLower(method)
	for _, m := range validMethods {
		if method == m {
			return true
		}
	}
	return false
}

// isValidServerURL performs basic URL validation
func (a *App) isValidServerURL(url string) bool {
	if url == "" {
		return false
	}
	
	// Basic validation - should be a valid URL or template
	if strings.HasPrefix(url, "http://") || strings.HasPrefix(url, "https://") {
		return true
	}
	
	// Allow relative URLs and templates
	if strings.HasPrefix(url, "/") || strings.Contains(url, "{") {
		return true
	}
	
	return false
}

// isValidOpenAPIFile checks if a file has a valid OpenAPI extension
func (a *App) isValidOpenAPIFile(filePath string) bool {
	ext := strings.ToLower(filepath.Ext(filePath))
	return ext == ".json" || ext == ".yaml" || ext == ".yml" || ext == ".openapi"
}

// copyFile copies a file from source to target
func (a *App) copyFile(source, target string) error {
	sourceFile, err := os.Open(source)
	if err != nil {
		return err
	}
	defer sourceFile.Close()

	// Create target directory if it doesn't exist
	targetDir := filepath.Dir(target)
	if err := os.MkdirAll(targetDir, 0755); err != nil {
		return err
	}

	targetFile, err := os.Create(target)
	if err != nil {
		return err
	}
	defer targetFile.Close()

	_, err = io.Copy(targetFile, sourceFile)
	return err
}

// copyDirectory recursively copies a directory
func (a *App) copyDirectory(source, target string) error {
	return filepath.Walk(source, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Calculate relative path
		relPath, err := filepath.Rel(source, path)
		if err != nil {
			return err
		}

		targetPath := filepath.Join(target, relPath)

		if info.IsDir() {
			return os.MkdirAll(targetPath, info.Mode())
		}

		return a.copyFile(path, targetPath)
	})
}

// calculateTotalSize calculates the total size of exported files
func (a *App) calculateTotalSize(files []ExportedFile) int64 {
	var total int64
	for _, file := range files {
		total += file.Size
	}
	return total
}

// generateProgressID generates a unique ID for progress tracking
func (a *App) generateProgressID() string {
	return "file_op_" + string(rune(time.Now().UnixNano()))
}

// emitFileProgress emits file operation progress via Wails events
func (a *App) emitFileProgress(operationID string, operationType string, progress int, message string, totalFiles int, processedFiles int) {
	if a.ctx == nil {
		return
	}

	progressData := FileOperationProgress{
		OperationID:        operationID,
		Type:               operationType,
		Progress:           progress,
		CurrentFile:        message,
		TotalFiles:         totalFiles,
		ProcessedFiles:     processedFiles,
		StartTime:          time.Now().Format(time.RFC3339),
		ElapsedTime:        0,
		EstimatedRemaining: 0,
	}

	// Emit progress event
	runtime.EventsEmit(a.ctx, "file:progress", progressData)
}

// AddRecentFile adds a file to the recent files list
func (a *App) AddRecentFile(filePath string, fileType string) error {
	if filePath == "" {
		return a.createAPIError("validation", ErrCodeValidation, "File path is required", nil)
	}

	// Get file info
	fileInfo, err := os.Stat(filePath)
	if err != nil {
		return a.createAPIError("file_system", ErrCodeFileAccess, "Failed to get file info", map[string]string{
			"path": filePath,
			"error": err.Error(),
		})
	}

	recentFile := RecentFile{
		Path:         filePath,
		Name:         filepath.Base(filePath),
		Size:         fileInfo.Size(),
		LastAccessed: time.Now().Format(time.RFC3339),
		Type:         fileType,
	}

	// Get current settings
	settings, err := a.GetSettings()
	if err != nil {
		return err
	}

	// Initialize recent files if not exists
	recentFiles := []RecentFile{}
	if len(settings.RecentFiles) > 0 {
		// Parse existing recent files (stored as JSON strings)
		for _, recentFileData := range settings.RecentFiles {
			var rf RecentFile
			if err := json.Unmarshal([]byte(recentFileData), &rf); err == nil {
				recentFiles = append(recentFiles, rf)
			}
		}
	}

	// Remove duplicate if exists
	for i, rf := range recentFiles {
		if rf.Path == filePath {
			recentFiles = append(recentFiles[:i], recentFiles[i+1:]...)
			break
		}
	}

	// Add new file to the beginning
	recentFiles = append([]RecentFile{recentFile}, recentFiles...)

	// Limit to 10 recent files
	if len(recentFiles) > 10 {
		recentFiles = recentFiles[:10]
	}

	// Convert back to JSON strings
	recentFileStrings := []string{}
	for _, rf := range recentFiles {
		if data, err := json.Marshal(rf); err == nil {
			recentFileStrings = append(recentFileStrings, string(data))
		}
	}

	// Update settings
	settings.RecentFiles = recentFileStrings
	return a.UpdateSettings(*settings)
}

// GetRecentFiles returns the list of recent files
func (a *App) GetRecentFiles() ([]RecentFile, error) {
	settings, err := a.GetSettings()
	if err != nil {
		return nil, err
	}

	recentFiles := []RecentFile{}
	for _, recentFileData := range settings.RecentFiles {
		var rf RecentFile
		if err := json.Unmarshal([]byte(recentFileData), &rf); err == nil {
			// Check if file still exists
			if _, err := os.Stat(rf.Path); err == nil {
				recentFiles = append(recentFiles, rf)
			}
		}
	}

	return recentFiles, nil
}

// ClearRecentFiles clears the recent files list
func (a *App) ClearRecentFiles() error {
	settings, err := a.GetSettings()
	if err != nil {
		return err
	}

	settings.RecentFiles = []string{}
	return a.UpdateSettings(*settings)
}

// RemoveRecentFile removes a specific file from the recent files list
func (a *App) RemoveRecentFile(filePath string) error {
	if filePath == "" {
		return a.createAPIError("validation", ErrCodeValidation, "File path is required", nil)
	}

	settings, err := a.GetSettings()
	if err != nil {
		return err
	}

	// Parse and filter recent files
	recentFiles := []RecentFile{}
	for _, recentFileData := range settings.RecentFiles {
		var rf RecentFile
		if err := json.Unmarshal([]byte(recentFileData), &rf); err == nil {
			if rf.Path != filePath {
				recentFiles = append(recentFiles, rf)
			}
		}
	}

	// Convert back to JSON strings
	recentFileStrings := []string{}
	for _, rf := range recentFiles {
		if data, err := json.Marshal(rf); err == nil {
			recentFileStrings = append(recentFileStrings, string(data))
		}
	}

	// Update settings
	settings.RecentFiles = recentFileStrings
	return a.UpdateSettings(*settings)
}