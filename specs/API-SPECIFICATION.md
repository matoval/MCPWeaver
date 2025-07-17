# MCPWeaver Internal API Specification

## Overview

This document defines the internal API contracts for MCPWeaver's backend services and frontend-backend communication through Wails context binding.

## Wails Context API

### Project Management API

#### CreateProject
```go
func (a *App) CreateProject(request CreateProjectRequest) (*Project, error)
```

**Request Structure:**
```go
type CreateProjectRequest struct {
    Name        string            `json:"name" validate:"required,min=1,max=100"`
    SpecPath    string            `json:"specPath,omitempty"`
    SpecURL     string            `json:"specUrl,omitempty"`
    OutputPath  string            `json:"outputPath" validate:"required"`
    Settings    ProjectSettings   `json:"settings"`
}

type ProjectSettings struct {
    PackageName     string   `json:"packageName" validate:"required,alphanum"`
    ServerPort      int      `json:"serverPort" validate:"min=1000,max=65535"`
    EnableLogging   bool     `json:"enableLogging"`
    LogLevel        string   `json:"logLevel" validate:"oneof=debug info warn error"`
    CustomTemplates []string `json:"customTemplates,omitempty"`
}
```

**Response Structure:**
```go
type Project struct {
    ID             string           `json:"id"`
    Name           string           `json:"name"`
    SpecPath       string           `json:"specPath"`
    SpecURL        string           `json:"specUrl"`
    OutputPath     string           `json:"outputPath"`
    Settings       ProjectSettings  `json:"settings"`
    Status         ProjectStatus    `json:"status"`
    CreatedAt      time.Time        `json:"createdAt"`
    UpdatedAt      time.Time        `json:"updatedAt"`
    LastGenerated  *time.Time       `json:"lastGenerated,omitempty"`
    GenerationCount int             `json:"generationCount"`
}

type ProjectStatus string

const (
    StatusCreated    ProjectStatus = "created"
    StatusValidating ProjectStatus = "validating"
    StatusReady      ProjectStatus = "ready"
    StatusGenerating ProjectStatus = "generating"
    StatusError      ProjectStatus = "error"
)
```

**Error Responses:**
```go
type APIError struct {
    Type      string            `json:"type"`
    Code      string            `json:"code"`
    Message   string            `json:"message"`
    Details   map[string]string `json:"details,omitempty"`
    Timestamp time.Time         `json:"timestamp"`
}

// Common error codes
const (
    ErrCodeValidation    = "VALIDATION_ERROR"
    ErrCodeNotFound      = "NOT_FOUND"
    ErrCodeInternalError = "INTERNAL_ERROR"
    ErrCodeFileAccess    = "FILE_ACCESS_ERROR"
)
```

#### GetProjects
```go
func (a *App) GetProjects() ([]*Project, error)
```

**Response:** Array of `Project` objects

#### GetProject
```go
func (a *App) GetProject(id string) (*Project, error)
```

**Parameters:**
- `id`: Project ID (string)

**Response:** `Project` object

#### UpdateProject
```go
func (a *App) UpdateProject(id string, updates ProjectUpdateRequest) (*Project, error)
```

**Request Structure:**
```go
type ProjectUpdateRequest struct {
    Name        *string          `json:"name,omitempty"`
    SpecPath    *string          `json:"specPath,omitempty"`
    SpecURL     *string          `json:"specUrl,omitempty"`
    OutputPath  *string          `json:"outputPath,omitempty"`
    Settings    *ProjectSettings `json:"settings,omitempty"`
}
```

#### DeleteProject
```go
func (a *App) DeleteProject(id string) error
```

**Parameters:**
- `id`: Project ID (string)

### Generation API

#### GenerateServer
```go
func (a *App) GenerateServer(projectId string) (*GenerationJob, error)
```

**Response Structure:**
```go
type GenerationJob struct {
    ID          string            `json:"id"`
    ProjectID   string            `json:"projectId"`
    Status      GenerationStatus  `json:"status"`
    Progress    float64           `json:"progress"`
    CurrentStep string            `json:"currentStep"`
    StartTime   time.Time         `json:"startTime"`
    EndTime     *time.Time        `json:"endTime,omitempty"`
    Results     *GenerationResults `json:"results,omitempty"`
    Errors      []GenerationError `json:"errors,omitempty"`
}

type GenerationStatus string

const (
    StatusStarted   GenerationStatus = "started"
    StatusParsing   GenerationStatus = "parsing"
    StatusMapping   GenerationStatus = "mapping"
    StatusGenerating GenerationStatus = "generating"
    StatusValidating GenerationStatus = "validating"
    StatusCompleted GenerationStatus = "completed"
    StatusFailed    GenerationStatus = "failed"
    StatusCancelled GenerationStatus = "cancelled"
)

type GenerationResults struct {
    ServerPath      string            `json:"serverPath"`
    GeneratedFiles  []GeneratedFile   `json:"generatedFiles"`
    MCPTools        []MCPTool         `json:"mcpTools"`
    Statistics      GenerationStats   `json:"statistics"`
}

type GeneratedFile struct {
    Path         string `json:"path"`
    Type         string `json:"type"`
    Size         int64  `json:"size"`
    LinesOfCode  int    `json:"linesOfCode"`
}

type GenerationStats struct {
    TotalEndpoints     int           `json:"totalEndpoints"`
    GeneratedTools     int           `json:"generatedTools"`
    ProcessingTime     time.Duration `json:"processingTime"`
    SpecComplexity     string        `json:"specComplexity"`
    TemplateVersion    string        `json:"templateVersion"`
}

type GenerationError struct {
    Type        string            `json:"type"`
    Message     string            `json:"message"`
    Details     string            `json:"details,omitempty"`
    Suggestions []string          `json:"suggestions,omitempty"`
    Location    *ErrorLocation    `json:"location,omitempty"`
}

type ErrorLocation struct {
    File   string `json:"file"`
    Line   int    `json:"line"`
    Column int    `json:"column"`
}
```

#### GetGenerationJob
```go
func (a *App) GetGenerationJob(jobId string) (*GenerationJob, error)
```

#### CancelGeneration
```go
func (a *App) CancelGeneration(jobId string) error
```

#### GetGenerationHistory
```go
func (a *App) GetGenerationHistory(projectId string) ([]*GenerationJob, error)
```

### Validation API

#### ValidateSpec
```go
func (a *App) ValidateSpec(specPath string) (*ValidationResult, error)
```

**Response Structure:**
```go
type ValidationResult struct {
    Valid           bool                `json:"valid"`
    Errors          []ValidationError   `json:"errors"`
    Warnings        []ValidationWarning `json:"warnings"`
    Suggestions     []string            `json:"suggestions"`
    SpecInfo        *SpecInfo           `json:"specInfo,omitempty"`
    ValidationTime  time.Duration       `json:"validationTime"`
}

type ValidationError struct {
    Type        string         `json:"type"`
    Message     string         `json:"message"`
    Path        string         `json:"path"`
    Line        int            `json:"line,omitempty"`
    Column      int            `json:"column,omitempty"`
    Severity    string         `json:"severity"`
    Code        string         `json:"code"`
    Location    *ErrorLocation `json:"location,omitempty"`
}

type ValidationWarning struct {
    Type        string `json:"type"`
    Message     string `json:"message"`
    Path        string `json:"path"`
    Suggestion  string `json:"suggestion"`
}

type SpecInfo struct {
    Version         string              `json:"version"`
    Title           string              `json:"title"`
    Description     string              `json:"description"`
    OperationCount  int                 `json:"operationCount"`
    SchemaCount     int                 `json:"schemaCount"`
    SecuritySchemes []SecurityScheme    `json:"securitySchemes"`
    Servers         []ServerInfo        `json:"servers"`
}

type SecurityScheme struct {
    Type        string `json:"type"`
    Name        string `json:"name"`
    Description string `json:"description"`
}

type ServerInfo struct {
    URL         string `json:"url"`
    Description string `json:"description"`
}
```

#### ValidateURL
```go
func (a *App) ValidateURL(url string) (*ValidationResult, error)
```

### File System API

#### SelectFile
```go
func (a *App) SelectFile(filters []FileFilter) (string, error)
```

**Request Structure:**
```go
type FileFilter struct {
    DisplayName string   `json:"displayName"`
    Pattern     string   `json:"pattern"`
    Extensions  []string `json:"extensions"`
}
```

#### SelectDirectory
```go
func (a *App) SelectDirectory(title string) (string, error)
```

#### SaveFile
```go
func (a *App) SaveFile(content string, defaultPath string, filters []FileFilter) (string, error)
```

#### ReadFile
```go
func (a *App) ReadFile(path string) (string, error)
```

#### WriteFile
```go
func (a *App) WriteFile(path string, content string) error
```

#### FileExists
```go
func (a *App) FileExists(path string) (bool, error)
```

### Settings API

#### GetSettings
```go
func (a *App) GetSettings() (*AppSettings, error)
```

**Response Structure:**
```go
type AppSettings struct {
    Theme           string            `json:"theme"`
    Language        string            `json:"language"`
    AutoSave        bool              `json:"autoSave"`
    DefaultOutputPath string          `json:"defaultOutputPath"`
    RecentProjects  []string          `json:"recentProjects"`
    WindowSettings  WindowSettings    `json:"windowSettings"`
    EditorSettings  EditorSettings    `json:"editorSettings"`
    GenerationSettings GenerationSettings `json:"generationSettings"`
}

type WindowSettings struct {
    Width       int  `json:"width"`
    Height      int  `json:"height"`
    Maximized   bool `json:"maximized"`
    X           int  `json:"x"`
    Y           int  `json:"y"`
}

type EditorSettings struct {
    FontSize        int    `json:"fontSize"`
    FontFamily      string `json:"fontFamily"`
    TabSize         int    `json:"tabSize"`
    WordWrap        bool   `json:"wordWrap"`
    LineNumbers     bool   `json:"lineNumbers"`
    SyntaxHighlight bool   `json:"syntaxHighlight"`
}

type GenerationSettings struct {
    DefaultTemplate     string   `json:"defaultTemplate"`
    EnableValidation    bool     `json:"enableValidation"`
    AutoOpenOutput      bool     `json:"autoOpenOutput"`
    ShowAdvancedOptions bool     `json:"showAdvancedOptions"`
    BackupOnGenerate    bool     `json:"backupOnGenerate"`
    CustomTemplates     []string `json:"customTemplates"`
}
```

#### UpdateSettings
```go
func (a *App) UpdateSettings(settings AppSettings) error
```

#### ResetSettings
```go
func (a *App) ResetSettings() error
```

### Template API

#### GetTemplates
```go
func (a *App) GetTemplates() ([]*Template, error)
```

**Response Structure:**
```go
type Template struct {
    ID          string            `json:"id"`
    Name        string            `json:"name"`
    Description string            `json:"description"`
    Version     string            `json:"version"`
    Author      string            `json:"author"`
    Type        TemplateType      `json:"type"`
    Path        string            `json:"path"`
    IsBuiltIn   bool              `json:"isBuiltIn"`
    Variables   []TemplateVariable `json:"variables"`
    CreatedAt   time.Time         `json:"createdAt"`
    UpdatedAt   time.Time         `json:"updatedAt"`
}

type TemplateType string

const (
    TemplateTypeDefault TemplateType = "default"
    TemplateTypeCustom  TemplateType = "custom"
    TemplateTypePlugin  TemplateType = "plugin"
)

type TemplateVariable struct {
    Name         string `json:"name"`
    Description  string `json:"description"`
    Type         string `json:"type"`
    DefaultValue string `json:"defaultValue"`
    Required     bool   `json:"required"`
}
```

#### InstallTemplate
```go
func (a *App) InstallTemplate(templatePath string) (*Template, error)
```

#### RemoveTemplate
```go
func (a *App) RemoveTemplate(templateId string) error
```

## Event System

### Real-time Events

MCPWeaver uses Wails' event system for real-time updates between backend and frontend:

#### Project Events
```go
// Emitted when a project is created
runtime.EventsEmit(ctx, "project:created", project)

// Emitted when a project is updated
runtime.EventsEmit(ctx, "project:updated", project)

// Emitted when a project is deleted
runtime.EventsEmit(ctx, "project:deleted", projectId)
```

#### Generation Events
```go
// Emitted when generation starts
runtime.EventsEmit(ctx, "generation:started", generationJob)

// Emitted for progress updates
runtime.EventsEmit(ctx, "generation:progress", GenerationProgress{
    JobID:       jobId,
    Progress:    0.75,
    CurrentStep: "Generating server code...",
    Message:     "Processing 15/20 endpoints",
})

// Emitted when generation completes
runtime.EventsEmit(ctx, "generation:completed", generationJob)

// Emitted when generation fails
runtime.EventsEmit(ctx, "generation:failed", GenerationError{
    JobID:   jobId,
    Type:    "validation",
    Message: "Invalid OpenAPI specification",
})
```

#### System Events
```go
// Emitted for system notifications
runtime.EventsEmit(ctx, "system:notification", Notification{
    Type:    "info",
    Title:   "Generation Complete",
    Message: "MCP server generated successfully",
    Actions: []string{"Open Folder", "Run Tests"},
})

// Emitted for application errors
runtime.EventsEmit(ctx, "system:error", SystemError{
    Type:        "critical",
    Message:     "Database connection failed",
    Timestamp:   time.Now(),
    Recoverable: true,
})
```

## Error Handling

### Standard Error Response Format
All API methods return errors in a consistent format:

```go
type APIError struct {
    Type        string            `json:"type"`
    Code        string            `json:"code"`
    Message     string            `json:"message"`
    Details     map[string]string `json:"details,omitempty"`
    Timestamp   time.Time         `json:"timestamp"`
    Suggestions []string          `json:"suggestions,omitempty"`
}
```

### Error Types
- `validation`: Input validation errors
- `not_found`: Resource not found
- `file_system`: File operations errors
- `network`: Network-related errors
- `internal`: Internal server errors
- `parsing`: OpenAPI parsing errors
- `generation`: Code generation errors

### Error Codes
- `VALIDATION_ERROR`: Input validation failed
- `NOT_FOUND`: Resource not found
- `FILE_ACCESS_ERROR`: File system access denied
- `NETWORK_ERROR`: Network connectivity issues
- `PARSING_ERROR`: OpenAPI specification parsing failed
- `GENERATION_ERROR`: Code generation failed
- `INTERNAL_ERROR`: Internal application error

## Data Validation

### Input Validation Rules
All API inputs are validated using struct tags and custom validators:

```go
type CreateProjectRequest struct {
    Name        string `json:"name" validate:"required,min=1,max=100,alphanum_space"`
    SpecPath    string `json:"specPath" validate:"omitempty,file_exists"`
    SpecURL     string `json:"specUrl" validate:"omitempty,url"`
    OutputPath  string `json:"outputPath" validate:"required,dir_writable"`
    Settings    ProjectSettings `json:"settings" validate:"required"`
}
```

### Custom Validators
- `file_exists`: Validates file existence
- `dir_writable`: Validates directory write permissions
- `alphanum_space`: Alphanumeric characters and spaces only
- `openapi_spec`: Valid OpenAPI specification format

## Rate Limiting and Concurrency

### Concurrent Operations
- Maximum 3 concurrent generations per user
- File operations are queued to prevent conflicts
- Background validation runs with lower priority

### Resource Management
- Memory usage monitoring during generation
- Automatic cleanup of temporary files
- Connection pooling for database operations

## API Versioning

### Version Header
All API calls include version information:
```go
const APIVersion = "1.0.0"
```

### Backward Compatibility
- Additive changes only in minor versions
- Breaking changes require major version bump
- Deprecation notices for removed features

This API specification provides a comprehensive contract for all internal component interactions in MCPWeaver, ensuring type safety and consistent behavior across the application.