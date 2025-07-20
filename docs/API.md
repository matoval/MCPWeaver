# MCPWeaver API Reference

This document provides comprehensive documentation for the MCPWeaver API, including all available methods, data structures, and usage examples.

## Table of Contents

- [Overview](#overview)
- [Application Lifecycle](#application-lifecycle)
- [File Operations](#file-operations)
- [Project Management](#project-management)
- [OpenAPI Validation](#openapi-validation)
- [MCP Server Generation](#mcp-server-generation)
- [Template Management](#template-management)
- [Settings Management](#settings-management)
- [Performance Monitoring](#performance-monitoring)
- [Error Handling](#error-handling)
- [Data Types](#data-types)
- [Events](#events)

## Overview

MCPWeaver exposes its functionality through a comprehensive API that enables both desktop GUI interaction and programmatic access. The API is built using the Wails framework, allowing seamless communication between the Go backend and the React frontend.

### API Architecture

```
Frontend (React/TypeScript)
        ↕ (Wails Context Binding)
Backend (Go) - App Struct
        ↕
Internal Services:
├── Parser Service
├── Validator Service  
├── Generator Service
├── Mapping Service
└── Database Layer
```

### Error Handling Pattern

All API methods follow a consistent error handling pattern:

```go
type APIError struct {
    Type          string            `json:"type"`
    Code          string            `json:"code"`
    Message       string            `json:"message"`
    Details       map[string]string `json:"details,omitempty"`
    Timestamp     time.Time         `json:"timestamp"`
    Suggestions   []string          `json:"suggestions,omitempty"`
    CorrelationID string            `json:"correlationId,omitempty"`
    Severity      ErrorSeverity     `json:"severity"`
    Recoverable   bool              `json:"recoverable"`
}
```

## Application Lifecycle

### OnStartup

Initializes the application when it starts.

**Method:** `OnStartup(ctx context.Context) error`

**Description:** Called when the app starts, before the frontend is loaded. Initializes database connections, repositories, and application settings.

**Parameters:**
- `ctx` (context.Context): Application context

**Returns:**
- `error`: Error if initialization fails

**Example:**
```go
// Called automatically by Wails framework
err := app.OnStartup(ctx)
if err != nil {
    log.Fatalf("Failed to start application: %v", err)
}
```

### OnShutdown

Cleans up resources when the application shuts down.

**Method:** `OnShutdown(ctx context.Context) error`

**Description:** Called when the app is about to quit. Saves settings, closes database connections, and performs cleanup.

**Parameters:**
- `ctx` (context.Context): Application context

**Returns:**
- `error`: Error if cleanup fails

### OnDomReady

Called after the frontend DOM is ready.

**Method:** `OnDomReady(ctx context.Context)`

**Description:** Emits a system ready event to notify the frontend that the backend is fully initialized.

**Parameters:**
- `ctx` (context.Context): Application context

### OnBeforeClose

Called when the application is about to quit.

**Method:** `OnBeforeClose(ctx context.Context) bool`

**Description:** Performs final checks before application shutdown.

**Parameters:**
- `ctx` (context.Context): Application context

**Returns:**
- `bool`: `false` to allow the app to close, `true` to prevent closing

## File Operations

### SelectFile

Opens a file selection dialog.

**Method:** `SelectFile(filters []FileFilter) (string, error)`

**Description:** Opens a native file selection dialog with specified filters.

**Parameters:**
- `filters` ([]FileFilter): File type filters

**Returns:**
- `string`: Selected file path
- `error`: Error if operation fails

**Example:**
```javascript
// TypeScript/JavaScript usage
const filters = [
    { name: "OpenAPI Files", extensions: ["json", "yaml", "yml"] },
    { name: "All Files", extensions: ["*"] }
];

try {
    const filePath = await window.go.app.App.SelectFile(filters);
    console.log("Selected file:", filePath);
} catch (error) {
    console.error("File selection failed:", error);
}
```

### SelectDirectory

Opens a directory selection dialog.

**Method:** `SelectDirectory(title string) (string, error)`

**Description:** Opens a native directory selection dialog.

**Parameters:**
- `title` (string): Dialog title

**Returns:**
- `string`: Selected directory path
- `error`: Error if operation fails

**Example:**
```javascript
try {
    const dirPath = await window.go.app.App.SelectDirectory("Select Output Directory");
    console.log("Selected directory:", dirPath);
} catch (error) {
    console.error("Directory selection failed:", error);
}
```

### SaveFile

Opens a save file dialog.

**Method:** `SaveFile(content string, defaultPath string, filters []FileFilter) (string, error)`

**Description:** Opens a native save file dialog and writes content to the selected file.

**Parameters:**
- `content` (string): Content to save
- `defaultPath` (string): Default file path/name
- `filters` ([]FileFilter): File type filters

**Returns:**
- `string`: Saved file path
- `error`: Error if operation fails

### ReadFile

Reads content from a file.

**Method:** `ReadFile(path string) (string, error)`

**Description:** Reads and returns the content of a file.

**Parameters:**
- `path` (string): File path to read

**Returns:**
- `string`: File content
- `error`: Error if operation fails

**Example:**
```javascript
try {
    const content = await window.go.app.App.ReadFile("/path/to/openapi.yaml");
    console.log("File content:", content);
} catch (error) {
    console.error("Failed to read file:", error);
}
```

### WriteFile

Writes content to a file.

**Method:** `WriteFile(path string, content string) error`

**Description:** Writes content to a file, creating directories if needed.

**Parameters:**
- `path` (string): File path to write
- `content` (string): Content to write

**Returns:**
- `error`: Error if operation fails

### FileExists

Checks if a file exists.

**Method:** `FileExists(path string) (bool, error)`

**Description:** Checks whether a file exists at the specified path.

**Parameters:**
- `path` (string): File path to check

**Returns:**
- `bool`: `true` if file exists, `false` otherwise
- `error`: Error if operation fails

### File Format Detection

#### GetDefaultOpenAPIFilters

Returns default file filters for OpenAPI specifications.

**Method:** `GetDefaultOpenAPIFilters() []FileFilter`

**Description:** Returns predefined file filters for OpenAPI file selection.

**Returns:**
- `[]FileFilter`: Array of file filters

#### GetSupportedFileFormats

Returns supported file formats.

**Method:** `GetSupportedFileFormats() []string`

**Description:** Returns a list of supported file formats for OpenAPI specifications.

**Returns:**
- `[]string`: Array of supported formats

#### DetectFileFormat

Detects the format of file content.

**Method:** `DetectFileFormat(content string, filename string) (string, error)`

**Description:** Analyzes file content and filename to determine the format.

**Parameters:**
- `content` (string): File content
- `filename` (string): File name

**Returns:**
- `string`: Detected format ("json", "yaml", "yml")
- `error`: Error if detection fails

### OpenAPI Import/Export

#### ImportOpenAPISpec

Imports an OpenAPI specification from a file.

**Method:** `ImportOpenAPISpec(filePath string) (*ImportResult, error)`

**Description:** Imports and validates an OpenAPI specification from a local file.

**Parameters:**
- `filePath` (string): Path to OpenAPI specification file

**Returns:**
- `*ImportResult`: Import result with validation details
- `error`: Error if import fails

**Example:**
```javascript
try {
    const result = await window.go.app.App.ImportOpenAPISpec("/path/to/openapi.yaml");
    console.log("Import successful:", result);
    console.log("Validation errors:", result.validationResult.errors);
} catch (error) {
    console.error("Import failed:", error);
}
```

#### ImportOpenAPISpecFromURL

Imports an OpenAPI specification from a URL.

**Method:** `ImportOpenAPISpecFromURL(url string) (*ImportResult, error)`

**Description:** Downloads and imports an OpenAPI specification from a URL.

**Parameters:**
- `url` (string): URL to OpenAPI specification

**Returns:**
- `*ImportResult`: Import result with validation details
- `error`: Error if import fails

**Example:**
```javascript
try {
    const result = await window.go.app.App.ImportOpenAPISpecFromURL(
        "https://petstore3.swagger.io/api/v3/openapi.json"
    );
    console.log("URL import successful:", result);
} catch (error) {
    console.error("URL import failed:", error);
}
```

#### ExportGeneratedServer

Exports a generated MCP server to a directory.

**Method:** `ExportGeneratedServer(projectID, targetDir string) (*ExportResult, error)`

**Description:** Exports the generated MCP server files for a project to a target directory.

**Parameters:**
- `projectID` (string): Project identifier
- `targetDir` (string): Target export directory

**Returns:**
- `*ExportResult`: Export result with file details
- `error`: Error if export fails

### Recent Files Management

#### AddRecentFile

Adds a file to the recent files list.

**Method:** `AddRecentFile(filePath string, fileType string) error`

**Description:** Adds a file to the recent files list for quick access.

**Parameters:**
- `filePath` (string): File path
- `fileType` (string): File type ("openapi", "project", etc.)

**Returns:**
- `error`: Error if operation fails

#### GetRecentFiles

Returns the list of recent files.

**Method:** `GetRecentFiles() ([]RecentFile, error)`

**Description:** Returns the list of recently accessed files.

**Returns:**
- `[]RecentFile`: Array of recent files
- `error`: Error if operation fails

#### ClearRecentFiles

Clears the recent files list.

**Method:** `ClearRecentFiles() error`

**Description:** Removes all entries from the recent files list.

**Returns:**
- `error`: Error if operation fails

## Project Management

### CreateProject

Creates a new project.

**Method:** `CreateProject(request CreateProjectRequest) (*Project, error)`

**Description:** Creates a new project with the specified configuration.

**Parameters:**
- `request` (CreateProjectRequest): Project creation request

**Returns:**
- `*Project`: Created project
- `error`: Error if creation fails

**Example:**
```javascript
const request = {
    name: "My API Project",
    description: "Convert my REST API to MCP",
    openAPISpec: "...", // OpenAPI specification content
    outputDirectory: "/path/to/output",
    template: "standard"
};

try {
    const project = await window.go.app.App.CreateProject(request);
    console.log("Project created:", project);
} catch (error) {
    console.error("Project creation failed:", error);
}
```

### GetProjects

Returns all projects.

**Method:** `GetProjects() ([]*Project, error)`

**Description:** Returns a list of all projects in the workspace.

**Returns:**
- `[]*Project`: Array of projects
- `error`: Error if operation fails

### GetProject

Returns a specific project by ID.

**Method:** `GetProject(id string) (*Project, error)`

**Description:** Returns project details for the specified ID.

**Parameters:**
- `id` (string): Project identifier

**Returns:**
- `*Project`: Project details
- `error`: Error if project not found

### UpdateProject

Updates an existing project.

**Method:** `UpdateProject(id string, updates ProjectUpdateRequest) (*Project, error)`

**Description:** Updates project properties with the specified changes.

**Parameters:**
- `id` (string): Project identifier
- `updates` (ProjectUpdateRequest): Update request

**Returns:**
- `*Project`: Updated project
- `error`: Error if update fails

### DeleteProject

Deletes a project.

**Method:** `DeleteProject(id string) error`

**Description:** Permanently deletes a project and all associated data.

**Parameters:**
- `id` (string): Project identifier

**Returns:**
- `error`: Error if deletion fails

### GetRecentProjects

Returns recently accessed projects.

**Method:** `GetRecentProjects() ([]*Project, error)`

**Description:** Returns a list of recently accessed projects.

**Returns:**
- `[]*Project`: Array of recent projects
- `error`: Error if operation fails

### SearchProjects

Searches for projects by query.

**Method:** `SearchProjects(query string) ([]*Project, error)`

**Description:** Searches projects by name, description, or tags.

**Parameters:**
- `query` (string): Search query

**Returns:**
- `[]*Project`: Array of matching projects
- `error`: Error if search fails

### ExportProject

Exports a project to JSON.

**Method:** `ExportProject(projectID string) (string, error)`

**Description:** Exports project data as JSON for backup or sharing.

**Parameters:**
- `projectID` (string): Project identifier

**Returns:**
- `string`: JSON representation of project
- `error`: Error if export fails

### ImportProject

Imports a project from JSON.

**Method:** `ImportProject(jsonData string) (*Project, error)`

**Description:** Imports a project from JSON data.

**Parameters:**
- `jsonData` (string): JSON project data

**Returns:**
- `*Project`: Imported project
- `error`: Error if import fails

## OpenAPI Validation

### ValidateSpec

Validates an OpenAPI specification file.

**Method:** `ValidateSpec(specPath string) (*ValidationResult, error)`

**Description:** Validates an OpenAPI specification file and returns detailed results.

**Parameters:**
- `specPath` (string): Path to OpenAPI specification file

**Returns:**
- `*ValidationResult`: Validation results
- `error`: Error if validation fails

**Example:**
```javascript
try {
    const result = await window.go.app.App.ValidateSpec("/path/to/openapi.yaml");
    
    if (result.valid) {
        console.log("Specification is valid!");
    } else {
        console.log("Validation errors:", result.errors);
        console.log("Warnings:", result.warnings);
    }
    
    console.log("Suggestions:", result.suggestions);
} catch (error) {
    console.error("Validation failed:", error);
}
```

### ValidateURL

Validates an OpenAPI specification from a URL.

**Method:** `ValidateURL(url string) (*ValidationResult, error)`

**Description:** Downloads and validates an OpenAPI specification from a URL.

**Parameters:**
- `url` (string): URL to OpenAPI specification

**Returns:**
- `*ValidationResult`: Validation results
- `error`: Error if validation fails

### ExportValidationResult

Exports validation results to a file.

**Method:** `ExportValidationResult(result *ValidationResult) (string, error)`

**Description:** Exports validation results as JSON for external analysis.

**Parameters:**
- `result` (*ValidationResult): Validation result to export

**Returns:**
- `string`: JSON representation of validation result
- `error`: Error if export fails

### Validation Cache Management

#### GetValidationCacheStats

Returns validation cache statistics.

**Method:** `GetValidationCacheStats() (*database.ValidationCacheStats, error)`

**Description:** Returns statistics about the validation cache.

**Returns:**
- `*database.ValidationCacheStats`: Cache statistics
- `error`: Error if operation fails

#### ClearValidationCache

Clears the validation cache.

**Method:** `ClearValidationCache() error`

**Description:** Removes all entries from the validation cache.

**Returns:**
- `error`: Error if operation fails

## MCP Server Generation

### GenerateServer

Starts MCP server generation for a project.

**Method:** `GenerateServer(projectID string) (*GenerationJob, error)`

**Description:** Initiates MCP server generation for the specified project.

**Parameters:**
- `projectID` (string): Project identifier

**Returns:**
- `*GenerationJob`: Generation job details
- `error`: Error if generation fails to start

**Example:**
```javascript
try {
    const job = await window.go.app.App.GenerateServer("project-123");
    console.log("Generation started:", job);
    
    // Monitor progress
    const interval = setInterval(async () => {
        const updatedJob = await window.go.app.App.GetGenerationJob(job.id);
        console.log("Progress:", updatedJob.progress);
        
        if (updatedJob.status === "completed" || updatedJob.status === "failed") {
            clearInterval(interval);
            console.log("Generation finished:", updatedJob);
        }
    }, 1000);
} catch (error) {
    console.error("Generation failed:", error);
}
```

### GetGenerationJob

Returns details of a generation job.

**Method:** `GetGenerationJob(jobID string) (*GenerationJob, error)`

**Description:** Returns current status and details of a generation job.

**Parameters:**
- `jobID` (string): Generation job identifier

**Returns:**
- `*GenerationJob`: Job details
- `error`: Error if job not found

### CancelGeneration

Cancels a running generation job.

**Method:** `CancelGeneration(jobID string) error`

**Description:** Attempts to cancel a running generation job.

**Parameters:**
- `jobID` (string): Generation job identifier

**Returns:**
- `error`: Error if cancellation fails

### GetGenerationHistory

Returns generation history for a project.

**Method:** `GetGenerationHistory(projectID string) ([]*GenerationJob, error)`

**Description:** Returns a list of all generation jobs for a project.

**Parameters:**
- `projectID` (string): Project identifier

**Returns:**
- `[]*GenerationJob`: Array of generation jobs
- `error`: Error if operation fails

## Template Management

*Note: Template management methods are currently disabled but will be restored in future versions.*

### Template Operations

The following template operations are planned:

- `CreateTemplate(request CreateTemplateRequest) (*Template, error)`
- `GetTemplate(id string) (*Template, error)`
- `GetAllTemplates() ([]*Template, error)`
- `UpdateTemplate(id string, request UpdateTemplateRequest) (*Template, error)`
- `DeleteTemplate(id string) error`
- `ValidateTemplate(templatePath string) (*TemplateValidationResult, error)`

## Settings Management

### GetSettings

Returns current application settings.

**Method:** `GetSettings() (*AppSettings, error)`

**Description:** Returns the current application settings configuration.

**Returns:**
- `*AppSettings`: Application settings
- `error`: Error if operation fails

### UpdateSettings

Updates application settings.

**Method:** `UpdateSettings(settings AppSettings) error`

**Description:** Updates application settings with new values.

**Parameters:**
- `settings` (AppSettings): New settings configuration

**Returns:**
- `error`: Error if update fails

**Example:**
```javascript
try {
    const currentSettings = await window.go.app.App.GetSettings();
    
    // Update theme
    currentSettings.theme = "dark";
    currentSettings.language = "en";
    
    await window.go.app.App.UpdateSettings(currentSettings);
    console.log("Settings updated successfully");
} catch (error) {
    console.error("Settings update failed:", error);
}
```

### ResetSettings

Resets settings to defaults.

**Method:** `ResetSettings() error`

**Description:** Resets all application settings to their default values.

**Returns:**
- `error`: Error if reset fails

### GetSettingsFilePath

Returns the settings file path.

**Method:** `GetSettingsFilePath() string`

**Description:** Returns the file system path where settings are stored.

**Returns:**
- `string`: Settings file path

### Theme and Appearance

#### UpdateTheme

Updates the application theme.

**Method:** `UpdateTheme(theme string) error`

**Description:** Updates the application theme.

**Parameters:**
- `theme` (string): Theme name ("light", "dark", "system")

**Returns:**
- `error`: Error if update fails

#### UpdateLanguage

Updates the application language.

**Method:** `UpdateLanguage(language string) error`

**Description:** Updates the application language.

**Parameters:**
- `language` (string): Language code ("en", "es", "fr", etc.)

**Returns:**
- `error`: Error if update fails

#### UpdateWindowSettings

Updates window settings.

**Method:** `UpdateWindowSettings(windowSettings WindowSettings) error`

**Description:** Updates window size, position, and behavior settings.

**Parameters:**
- `windowSettings` (WindowSettings): New window settings

**Returns:**
- `error`: Error if update fails

### Settings Import/Export

#### ExportSettings

Exports settings to JSON.

**Method:** `ExportSettings() (string, error)`

**Description:** Exports current settings as JSON for backup or sharing.

**Returns:**
- `string`: JSON representation of settings
- `error`: Error if export fails

#### ImportSettings

Imports settings from JSON.

**Method:** `ImportSettings() error`

**Description:** Imports settings from JSON data.

**Returns:**
- `error`: Error if import fails

## Performance Monitoring

### GetPerformanceMetrics

Returns current performance metrics.

**Method:** `GetPerformanceMetrics() *PerformanceMetrics`

**Description:** Returns current application performance metrics.

**Returns:**
- `*PerformanceMetrics`: Performance metrics

### GetMemoryUsage

Returns current memory usage.

**Method:** `GetMemoryUsage() float64`

**Description:** Returns current memory usage in MB.

**Returns:**
- `float64`: Memory usage in MB

### ForceGarbageCollection

Forces garbage collection.

**Method:** `ForceGarbageCollection()`

**Description:** Forces Go garbage collection to free memory.

## Error Handling

### ReportError

Reports an error from the frontend.

**Method:** `ReportError(errorReport map[string]interface{}) error`

**Description:** Allows the frontend to report errors to the backend for logging and analysis.

**Parameters:**
- `errorReport` (map[string]interface{}): Error report data

**Returns:**
- `error`: Error if reporting fails

**Example:**
```javascript
try {
    await window.go.app.App.ReportError({
        type: "validation_error",
        message: "Invalid OpenAPI specification",
        details: {
            file: "/path/to/spec.yaml",
            line: "42",
            column: "10"
        },
        timestamp: new Date().toISOString(),
        userAgent: navigator.userAgent
    });
} catch (error) {
    console.error("Error reporting failed:", error);
}
```

## Data Types

### Project

```go
type Project struct {
    ID              string    `json:"id"`
    Name            string    `json:"name"`
    Description     string    `json:"description"`
    OpenAPISpec     string    `json:"openAPISpec"`
    OutputDirectory string    `json:"outputDirectory"`
    Template        string    `json:"template"`
    Status          string    `json:"status"`
    CreatedAt       time.Time `json:"createdAt"`
    UpdatedAt       time.Time `json:"updatedAt"`
    LastGenerated   time.Time `json:"lastGenerated,omitempty"`
}
```

### CreateProjectRequest

```go
type CreateProjectRequest struct {
    Name            string `json:"name"`
    Description     string `json:"description,omitempty"`
    OpenAPISpec     string `json:"openAPISpec"`
    OutputDirectory string `json:"outputDirectory,omitempty"`
    Template        string `json:"template,omitempty"`
}
```

### ValidationResult

```go
type ValidationResult struct {
    Valid           bool                  `json:"valid"`
    Errors          []ValidationError     `json:"errors"`
    Warnings        []ValidationWarning   `json:"warnings"`
    Suggestions     []string             `json:"suggestions"`
    SpecInfo        *SpecInfo            `json:"specInfo,omitempty"`
    ValidationTime  time.Duration        `json:"validationTime"`
    CacheHit        bool                 `json:"cacheHit"`
    ValidatedAt     time.Time            `json:"validatedAt"`
}
```

### GenerationJob

```go
type GenerationJob struct {
    ID          string            `json:"id"`
    ProjectID   string            `json:"projectId"`
    Status      GenerationStatus  `json:"status"`
    Progress    float64          `json:"progress"`
    Step        string           `json:"step"`
    StartedAt   time.Time        `json:"startedAt"`
    CompletedAt *time.Time       `json:"completedAt,omitempty"`
    Error       string           `json:"error,omitempty"`
    OutputPath  string           `json:"outputPath,omitempty"`
}
```

### ImportResult

```go
type ImportResult struct {
    Success          bool               `json:"success"`
    ValidationResult *ValidationResult  `json:"validationResult"`
    SpecInfo         *SpecInfo         `json:"specInfo,omitempty"`
    Format           string            `json:"format"`
    Size             int64             `json:"size"`
    ImportTime       time.Duration     `json:"importTime"`
    Source           string            `json:"source"`
}
```

### FileFilter

```go
type FileFilter struct {
    Name       string   `json:"name"`
    Extensions []string `json:"extensions"`
}
```

### AppSettings

```go
type AppSettings struct {
    Theme             string            `json:"theme"`
    Language          string            `json:"language"`
    AutoSave          bool              `json:"autoSave"`
    DefaultOutputPath string            `json:"defaultOutputPath"`
    RecentProjects    []string          `json:"recentProjects"`
    WindowSettings    WindowSettings    `json:"windowSettings"`
    EditorSettings    EditorSettings    `json:"editorSettings"`
    GenerationSettings GenerationSettings `json:"generationSettings"`
}
```

### RecentFile

```go
type RecentFile struct {
    Path       string    `json:"path"`
    Name       string    `json:"name"`
    Type       string    `json:"type"`
    Size       int64     `json:"size"`
    AccessedAt time.Time `json:"accessedAt"`
}
```

## Events

MCPWeaver emits various events that the frontend can listen to for real-time updates.

### System Events

**system:startup**
- Emitted when the application starts
- Payload: `{ timestamp, version, startup_time }`

**system:ready**
- Emitted when the DOM is ready
- Payload: `{ timestamp }`

**system:shutdown**
- Emitted when the application shuts down
- Payload: `{ timestamp }`

**system:error**
- Emitted when an error occurs
- Payload: `APIError` object

**system:notification**
- Emitted for system notifications
- Payload: `{ type, title, message, timestamp }`

### File Events

**file:progress**
- Emitted during file operations
- Payload: `{ operationId, type, progress, message, totalFiles, processedFiles }`

### Project Events

**project:created**
- Emitted when a project is created
- Payload: `Project` object

**project:updated**
- Emitted when a project is updated
- Payload: `Project` object

**project:deleted**
- Emitted when a project is deleted
- Payload: `{ projectId }`

### Generation Events

**generation:progress**
- Emitted during MCP server generation
- Payload: `GenerationJob` object

**generation:completed**
- Emitted when generation completes
- Payload: `GenerationJob` object

**generation:failed**
- Emitted when generation fails
- Payload: `GenerationJob` object

### Event Listening Example

```javascript
// Listen for generation progress
window.runtime.EventsOn("generation:progress", (job) => {
    console.log(`Generation progress: ${job.progress}%`);
    updateProgressBar(job.progress);
});

// Listen for system notifications
window.runtime.EventsOn("system:notification", (notification) => {
    showNotification(notification.title, notification.message, notification.type);
});

// Listen for errors
window.runtime.EventsOn("system:error", (error) => {
    console.error("System error:", error);
    showErrorDialog(error.message);
});
```

## Usage Examples

### Complete Project Workflow

```javascript
async function completeWorkflow() {
    try {
        // 1. Create a new project
        const project = await window.go.app.App.CreateProject({
            name: "Petstore API",
            description: "Convert Petstore API to MCP",
            openAPISpec: "", // Will be imported
            template: "standard"
        });
        
        // 2. Import OpenAPI specification
        const importResult = await window.go.app.App.ImportOpenAPISpecFromURL(
            "https://petstore3.swagger.io/api/v3/openapi.json"
        );
        
        if (!importResult.success) {
            throw new Error("Failed to import specification");
        }
        
        // 3. Update project with imported spec
        await window.go.app.App.UpdateProject(project.id, {
            openAPISpec: importResult.specInfo.content
        });
        
        // 4. Validate the specification
        const validation = await window.go.app.App.ValidateSpec(importResult.source);
        
        if (!validation.valid) {
            console.warn("Validation warnings:", validation.warnings);
            // Handle validation errors
        }
        
        // 5. Generate MCP server
        const job = await window.go.app.App.GenerateServer(project.id);
        
        // 6. Monitor generation progress
        const checkProgress = async () => {
            const updatedJob = await window.go.app.App.GetGenerationJob(job.id);
            
            if (updatedJob.status === "completed") {
                console.log("Generation completed!");
                console.log("Output path:", updatedJob.outputPath);
                return;
            }
            
            if (updatedJob.status === "failed") {
                throw new Error(`Generation failed: ${updatedJob.error}`);
            }
            
            console.log(`Progress: ${updatedJob.progress}% - ${updatedJob.step}`);
            setTimeout(checkProgress, 1000);
        };
        
        await checkProgress();
        
    } catch (error) {
        console.error("Workflow failed:", error);
        // Report error for analysis
        await window.go.app.App.ReportError({
            type: "workflow_error",
            message: error.message,
            stack: error.stack,
            timestamp: new Date().toISOString()
        });
    }
}
```

### Settings Management

```javascript
async function manageSettings() {
    try {
        // Get current settings
        const settings = await window.go.app.App.GetSettings();
        console.log("Current settings:", settings);
        
        // Update specific settings
        settings.theme = "dark";
        settings.autoSave = true;
        settings.defaultOutputPath = "/home/user/mcp-servers";
        
        // Apply settings
        await window.go.app.App.UpdateSettings(settings);
        
        // Update theme immediately
        await window.go.app.App.UpdateTheme("dark");
        
        console.log("Settings updated successfully");
        
    } catch (error) {
        console.error("Settings management failed:", error);
    }
}
```

### Performance Monitoring

```javascript
function monitorPerformance() {
    setInterval(async () => {
        try {
            const metrics = await window.go.app.App.GetPerformanceMetrics();
            const memoryUsage = await window.go.app.App.GetMemoryUsage();
            
            console.log("Performance metrics:", metrics);
            console.log(`Memory usage: ${memoryUsage.toFixed(2)} MB`);
            
            // Force garbage collection if memory usage is high
            if (memoryUsage > 100) {
                await window.go.app.App.ForceGarbageCollection();
                console.log("Forced garbage collection");
            }
            
        } catch (error) {
            console.error("Performance monitoring failed:", error);
        }
    }, 5000); // Check every 5 seconds
}
```

## Integration Examples

### Wails Frontend Integration

MCPWeaver uses Wails context binding to expose Go methods to the frontend. Here's how to integrate with the API from React/TypeScript:

#### Setup

```typescript
// types/api.ts - TypeScript type definitions
export interface Project {
  id: string;
  name: string;
  description: string;
  openAPISpec: string;
  outputDirectory: string;
  template: string;
  status: string;
  createdAt: string;
  updatedAt: string;
  lastGenerated?: string;
}

export interface ValidationResult {
  valid: boolean;
  errors: ValidationError[];
  warnings: ValidationWarning[];
  suggestions: string[];
  specInfo?: SpecInfo;
  validationTime: number;
  cacheHit: boolean;
  validatedAt: string;
}

export interface GenerationJob {
  id: string;
  projectId: string;
  status: 'pending' | 'running' | 'completed' | 'failed' | 'cancelled';
  progress: number;
  step: string;
  startedAt: string;
  completedAt?: string;
  error?: string;
  outputPath?: string;
}
```

#### React Hook for API Integration

```typescript
// hooks/useAPI.ts
import { useState, useCallback } from 'react';

interface APIError {
  type: string;
  code: string;
  message: string;
  details?: Record<string, string>;
  timestamp: string;
  suggestions?: string[];
  correlationId?: string;
  severity: string;
  recoverable: boolean;
}

interface APIResponse<T> {
  data?: T;
  error?: APIError;
  loading: boolean;
}

export const useAPI = () => {
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<APIError | null>(null);

  const call = useCallback(async <T>(
    apiCall: () => Promise<T>
  ): Promise<APIResponse<T>> => {
    setLoading(true);
    setError(null);

    try {
      const data = await apiCall();
      return { data, loading: false };
    } catch (err: any) {
      const apiError = err as APIError;
      setError(apiError);
      return { error: apiError, loading: false };
    } finally {
      setLoading(false);
    }
  }, []);

  return { call, loading, error };
};
```

#### Project Management Component

```typescript
// components/ProjectManager.tsx
import React, { useState, useEffect } from 'react';
import { useAPI } from '../hooks/useAPI';
import { Project, CreateProjectRequest } from '../types/api';

const ProjectManager: React.FC = () => {
  const { call, loading, error } = useAPI();
  const [projects, setProjects] = useState<Project[]>([]);
  const [selectedProject, setSelectedProject] = useState<Project | null>(null);

  // Load projects on component mount
  useEffect(() => {
    loadProjects();
  }, []);

  const loadProjects = async () => {
    const response = await call(() => window.go.app.App.GetProjects());
    if (response.data) {
      setProjects(response.data);
    }
  };

  const createProject = async (request: CreateProjectRequest) => {
    const response = await call(() => window.go.app.App.CreateProject(request));
    if (response.data) {
      setProjects(prev => [...prev, response.data!]);
      return response.data;
    }
    return null;
  };

  const generateServer = async (projectId: string) => {
    const response = await call(() => window.go.app.App.GenerateServer(projectId));
    if (response.data) {
      // Monitor generation progress
      monitorGeneration(response.data.id);
    }
  };

  const monitorGeneration = async (jobId: string) => {
    const checkProgress = async () => {
      const response = await call(() => window.go.app.App.GetGenerationJob(jobId));
      if (response.data) {
        const job = response.data;
        
        if (job.status === 'completed' || job.status === 'failed') {
          // Generation finished
          console.log('Generation finished:', job);
          return;
        }
        
        // Continue monitoring
        setTimeout(checkProgress, 1000);
      }
    };
    
    checkProgress();
  };

  return (
    <div className="project-manager">
      {/* Project management UI */}
      {loading && <div>Loading...</div>}
      {error && (
        <div className="error">
          <h3>Error: {error.message}</h3>
          {error.suggestions && (
            <ul>
              {error.suggestions.map((suggestion, index) => (
                <li key={index}>{suggestion}</li>
              ))}
            </ul>
          )}
        </div>
      )}
      {/* Rest of the component */}
    </div>
  );
};

export default ProjectManager;
```

#### Event Handling

```typescript
// services/EventService.ts
class EventService {
  private listeners: Map<string, Array<(data: any) => void>> = new Map();

  constructor() {
    this.setupEventListeners();
  }

  private setupEventListeners() {
    // Listen for system events
    window.runtime.EventsOn('system:startup', this.handleSystemStartup.bind(this));
    window.runtime.EventsOn('system:error', this.handleSystemError.bind(this));
    window.runtime.EventsOn('system:notification', this.handleNotification.bind(this));
    
    // Listen for generation events
    window.runtime.EventsOn('generation:progress', this.handleGenerationProgress.bind(this));
    window.runtime.EventsOn('generation:completed', this.handleGenerationCompleted.bind(this));
    window.runtime.EventsOn('generation:failed', this.handleGenerationFailed.bind(this));
    
    // Listen for project events
    window.runtime.EventsOn('project:created', this.handleProjectCreated.bind(this));
    window.runtime.EventsOn('project:updated', this.handleProjectUpdated.bind(this));
  }

  private handleSystemStartup(data: any) {
    console.log('System started:', data);
    this.emit('system:ready', data);
  }

  private handleSystemError(error: any) {
    console.error('System error:', error);
    this.emit('error', error);
  }

  private handleNotification(notification: any) {
    console.log('Notification:', notification);
    this.emit('notification', notification);
  }

  private handleGenerationProgress(job: any) {
    console.log('Generation progress:', job.progress + '%');
    this.emit('generation:progress', job);
  }

  private handleGenerationCompleted(job: any) {
    console.log('Generation completed:', job);
    this.emit('generation:completed', job);
  }

  private handleGenerationFailed(job: any) {
    console.error('Generation failed:', job.error);
    this.emit('generation:failed', job);
  }

  private handleProjectCreated(project: any) {
    console.log('Project created:', project);
    this.emit('project:created', project);
  }

  private handleProjectUpdated(project: any) {
    console.log('Project updated:', project);
    this.emit('project:updated', project);
  }

  public on(event: string, callback: (data: any) => void) {
    if (!this.listeners.has(event)) {
      this.listeners.set(event, []);
    }
    this.listeners.get(event)!.push(callback);
  }

  public off(event: string, callback: (data: any) => void) {
    const listeners = this.listeners.get(event);
    if (listeners) {
      const index = listeners.indexOf(callback);
      if (index > -1) {
        listeners.splice(index, 1);
      }
    }
  }

  private emit(event: string, data: any) {
    const listeners = this.listeners.get(event);
    if (listeners) {
      listeners.forEach(callback => callback(data));
    }
  }
}

export const eventService = new EventService();
```

## SDK Usage

For programmatic access to MCPWeaver functionality, you can use the API directly:

### Go SDK Example

```go
package main

import (
    "context"
    "fmt"
    "log"
    
    "MCPWeaver/internal/app"
)

func main() {
    // Create app instance
    app := app.NewApp()
    
    // Initialize
    ctx := context.Background()
    if err := app.OnStartup(ctx); err != nil {
        log.Fatal(err)
    }
    defer app.OnShutdown(ctx)
    
    // Create a project
    project, err := app.CreateProject(app.CreateProjectRequest{
        Name:        "Example API",
        Description: "Example project",
        OpenAPISpec: openAPIContent,
        Template:    "standard",
    })
    if err != nil {
        log.Fatal(err)
    }
    
    fmt.Printf("Created project: %s\n", project.ID)
    
    // Generate MCP server
    job, err := app.GenerateServer(project.ID)
    if err != nil {
        log.Fatal(err)
    }
    
    fmt.Printf("Generation started: %s\n", job.ID)
}
```

### Error Handling Best Practices

```typescript
// utils/errorHandler.ts
export class ErrorHandler {
  static async handleAPICall<T>(
    apiCall: () => Promise<T>,
    context: string
  ): Promise<T | null> {
    try {
      return await apiCall();
    } catch (error: any) {
      this.logError(error, context);
      this.showUserError(error);
      
      // Report error for analysis
      if (error.type !== 'validation_error') {
        await this.reportError(error, context);
      }
      
      return null;
    }
  }

  private static logError(error: any, context: string) {
    console.error(`[${context}] API Error:`, {
      type: error.type,
      code: error.code,
      message: error.message,
      details: error.details,
      timestamp: error.timestamp
    });
  }

  private static showUserError(error: any) {
    // Show user-friendly error message
    const message = this.getUserFriendlyMessage(error);
    
    // Use your notification system
    // notificationService.showError(message, error.suggestions);
  }

  private static getUserFriendlyMessage(error: any): string {
    switch (error.type) {
      case 'validation_error':
        return 'The OpenAPI specification has validation errors. Please fix them and try again.';
      case 'file_not_found':
        return 'The specified file could not be found. Please check the file path.';
      case 'network_error':
        return 'Unable to connect to the URL. Please check your internet connection.';
      case 'generation_error':
        return 'An error occurred during server generation. Please check the logs for details.';
      default:
        return error.message || 'An unexpected error occurred.';
    }
  }

  private static async reportError(error: any, context: string) {
    try {
      await window.go.app.App.ReportError({
        type: error.type,
        message: error.message,
        context: context,
        timestamp: new Date().toISOString(),
        userAgent: navigator.userAgent,
        url: window.location.href,
        stack: error.stack
      });
    } catch (reportError) {
      console.error('Failed to report error:', reportError);
    }
  }
}
```

## Testing

### Backend Testing

```go
// tests/app_test.go
package tests

import (
    "context"
    "testing"
    
    "MCPWeaver/internal/app"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/require"
)

func TestCreateProject(t *testing.T) {
    app := app.NewApp()
    ctx := context.Background()
    
    err := app.OnStartup(ctx)
    require.NoError(t, err)
    defer app.OnShutdown(ctx)
    
    request := app.CreateProjectRequest{
        Name:        "Test Project",
        Description: "Test Description",
        OpenAPISpec: `{"openapi":"3.0.0","info":{"title":"Test","version":"1.0.0"},"paths":{}}`,
        Template:    "standard",
    }
    
    project, err := app.CreateProject(request)
    require.NoError(t, err)
    assert.NotEmpty(t, project.ID)
    assert.Equal(t, request.Name, project.Name)
}

func TestValidateSpec(t *testing.T) {
    app := app.NewApp()
    ctx := context.Background()
    
    err := app.OnStartup(ctx)
    require.NoError(t, err)
    defer app.OnShutdown(ctx)
    
    // Create a temporary OpenAPI file
    specContent := `{
        "openapi": "3.0.0",
        "info": {
            "title": "Test API",
            "version": "1.0.0"
        },
        "paths": {}
    }`
    
    err = app.WriteFile("/tmp/test-spec.json", specContent)
    require.NoError(t, err)
    
    result, err := app.ValidateSpec("/tmp/test-spec.json")
    require.NoError(t, err)
    assert.True(t, result.Valid)
}
```

### Frontend Testing

```typescript
// tests/api.test.ts
import { renderHook, act } from '@testing-library/react';
import { useAPI } from '../hooks/useAPI';

// Mock the Wails API
const mockAPI = {
  CreateProject: jest.fn(),
  GetProjects: jest.fn(),
  ValidateSpec: jest.fn(),
};

(global as any).window = {
  go: {
    app: {
      App: mockAPI
    }
  }
};

describe('useAPI Hook', () => {
  beforeEach(() => {
    jest.clearAllMocks();
  });

  it('should handle successful API calls', async () => {
    const mockProject = {
      id: '123',
      name: 'Test Project',
      description: 'Test Description',
      openAPISpec: '{}',
      template: 'standard',
      status: 'active',
      createdAt: new Date().toISOString(),
      updatedAt: new Date().toISOString()
    };

    mockAPI.CreateProject.mockResolvedValue(mockProject);

    const { result } = renderHook(() => useAPI());

    let response;
    await act(async () => {
      response = await result.current.call(() => 
        mockAPI.CreateProject({
          name: 'Test Project',
          description: 'Test Description',
          openAPISpec: '{}',
          template: 'standard'
        })
      );
    });

    expect(response.data).toEqual(mockProject);
    expect(response.error).toBeUndefined();
    expect(result.current.loading).toBe(false);
  });

  it('should handle API errors', async () => {
    const mockError = {
      type: 'validation_error',
      code: 'INVALID_SPEC',
      message: 'Invalid OpenAPI specification',
      timestamp: new Date().toISOString(),
      severity: 'error',
      recoverable: true
    };

    mockAPI.ValidateSpec.mockRejectedValue(mockError);

    const { result } = renderHook(() => useAPI());

    let response;
    await act(async () => {
      response = await result.current.call(() => 
        mockAPI.ValidateSpec('/path/to/invalid-spec.yaml')
      );
    });

    expect(response.data).toBeUndefined();
    expect(response.error).toEqual(mockError);
    expect(result.current.error).toEqual(mockError);
  });
});
```

## Performance Optimization

### Caching Strategy

```typescript
// services/CacheService.ts
class CacheService {
  private cache = new Map<string, { data: any; timestamp: number; ttl: number }>();

  set(key: string, data: any, ttl: number = 300000) { // 5 minutes default
    this.cache.set(key, {
      data,
      timestamp: Date.now(),
      ttl
    });
  }

  get(key: string): any | null {
    const entry = this.cache.get(key);
    if (!entry) return null;

    if (Date.now() - entry.timestamp > entry.ttl) {
      this.cache.delete(key);
      return null;
    }

    return entry.data;
  }

  clear() {
    this.cache.clear();
  }

  // Cache API responses
  async getCachedProjects(): Promise<Project[]> {
    const cached = this.get('projects');
    if (cached) return cached;

    const projects = await window.go.app.App.GetProjects();
    this.set('projects', projects, 60000); // 1 minute cache
    return projects;
  }
}

export const cacheService = new CacheService();
```

### Batch Operations

```typescript
// services/BatchService.ts
class BatchService {
  private queue: Array<() => Promise<any>> = [];
  private processing = false;

  async addToBatch<T>(operation: () => Promise<T>): Promise<T> {
    return new Promise((resolve, reject) => {
      this.queue.push(async () => {
        try {
          const result = await operation();
          resolve(result);
        } catch (error) {
          reject(error);
        }
      });

      this.processBatch();
    });
  }

  private async processBatch() {
    if (this.processing || this.queue.length === 0) return;

    this.processing = true;

    while (this.queue.length > 0) {
      const batch = this.queue.splice(0, 5); // Process 5 at a time
      await Promise.all(batch.map(operation => operation()));
    }

    this.processing = false;
  }
}

export const batchService = new BatchService();
```

This API reference provides comprehensive documentation for all MCPWeaver functionality. For additional examples and tutorials, see the [User Guide](USER_GUIDE.md) and [Developer Guide](DEVELOPER.md).