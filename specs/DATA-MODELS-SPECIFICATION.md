# MCPWeaver Data Models and Schemas Specification

## Overview

This document defines the data models, database schemas, and data structures used throughout MCPWeaver for consistent data handling and storage.

## Database Schema (SQLite)

### Projects Table
```sql
CREATE TABLE projects (
    id TEXT PRIMARY KEY,
    name TEXT NOT NULL,
    spec_path TEXT,
    spec_url TEXT,
    output_path TEXT NOT NULL,
    settings TEXT NOT NULL, -- JSON serialized ProjectSettings
    status TEXT NOT NULL DEFAULT 'created',
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    last_generated DATETIME,
    generation_count INTEGER DEFAULT 0
);

-- Indexes for performance
CREATE INDEX idx_projects_name ON projects(name);
CREATE INDEX idx_projects_status ON projects(status);
CREATE INDEX idx_projects_created_at ON projects(created_at);
CREATE INDEX idx_projects_last_generated ON projects(last_generated);
```

### Generations Table
```sql
CREATE TABLE generations (
    id TEXT PRIMARY KEY,
    project_id TEXT NOT NULL,
    status TEXT NOT NULL,
    progress REAL DEFAULT 0.0,
    current_step TEXT,
    start_time DATETIME DEFAULT CURRENT_TIMESTAMP,
    end_time DATETIME,
    results TEXT, -- JSON serialized GenerationResults
    errors TEXT, -- JSON serialized GenerationError array
    processing_time_ms INTEGER,
    FOREIGN KEY (project_id) REFERENCES projects(id) ON DELETE CASCADE
);

-- Indexes for performance
CREATE INDEX idx_generations_project_id ON generations(project_id);
CREATE INDEX idx_generations_status ON generations(status);
CREATE INDEX idx_generations_start_time ON generations(start_time);
```

### Templates Table
```sql
CREATE TABLE templates (
    id TEXT PRIMARY KEY,
    name TEXT NOT NULL,
    description TEXT,
    version TEXT NOT NULL,
    author TEXT,
    type TEXT NOT NULL, -- 'default', 'custom', 'plugin'
    path TEXT NOT NULL,
    is_built_in BOOLEAN DEFAULT FALSE,
    variables TEXT, -- JSON serialized TemplateVariable array
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- Indexes for performance
CREATE INDEX idx_templates_name ON templates(name);
CREATE INDEX idx_templates_type ON templates(type);
CREATE INDEX idx_templates_is_built_in ON templates(is_built_in);
```

### Application Settings Table
```sql
CREATE TABLE app_settings (
    key TEXT PRIMARY KEY,
    value TEXT NOT NULL,
    type TEXT NOT NULL, -- 'string', 'number', 'boolean', 'json'
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);
```

### Validation Cache Table
```sql
CREATE TABLE validation_cache (
    spec_hash TEXT PRIMARY KEY,
    spec_path TEXT,
    spec_url TEXT,
    validation_result TEXT NOT NULL, -- JSON serialized ValidationResult
    cached_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    expires_at DATETIME NOT NULL
);

-- Index for cleanup
CREATE INDEX idx_validation_cache_expires_at ON validation_cache(expires_at);
```

## Core Data Models

### Project Models

#### Project
```go
type Project struct {
    ID              string          `json:"id" db:"id"`
    Name            string          `json:"name" db:"name"`
    SpecPath        string          `json:"specPath" db:"spec_path"`
    SpecURL         string          `json:"specUrl" db:"spec_url"`
    OutputPath      string          `json:"outputPath" db:"output_path"`
    Settings        ProjectSettings `json:"settings" db:"settings"`
    Status          ProjectStatus   `json:"status" db:"status"`
    CreatedAt       time.Time       `json:"createdAt" db:"created_at"`
    UpdatedAt       time.Time       `json:"updatedAt" db:"updated_at"`
    LastGenerated   *time.Time      `json:"lastGenerated" db:"last_generated"`
    GenerationCount int             `json:"generationCount" db:"generation_count"`
}

type ProjectStatus string

const (
    ProjectStatusCreated    ProjectStatus = "created"
    ProjectStatusValidating ProjectStatus = "validating"
    ProjectStatusReady      ProjectStatus = "ready"
    ProjectStatusGenerating ProjectStatus = "generating"
    ProjectStatusError      ProjectStatus = "error"
)
```

#### ProjectSettings
```go
type ProjectSettings struct {
    // Server Configuration
    PackageName     string   `json:"packageName"`
    ServerPort      int      `json:"serverPort"`
    BaseURL         string   `json:"baseUrl"`
    
    // Code Generation Options
    EnableLogging   bool     `json:"enableLogging"`
    LogLevel        string   `json:"logLevel"`
    GenerateTests   bool     `json:"generateTests"`
    GenerateDocs    bool     `json:"generateDocs"`
    
    // Template Configuration
    TemplateID      string   `json:"templateId"`
    CustomTemplates []string `json:"customTemplates"`
    
    // Output Options
    OutputFormat    string   `json:"outputFormat"` // "single", "modular"
    CompressOutput  bool     `json:"compressOutput"`
    
    // Validation Options
    StrictValidation bool    `json:"strictValidation"`
    SkipValidation   bool    `json:"skipValidation"`
    
    // Advanced Options
    CustomHeaders    map[string]string `json:"customHeaders"`
    TimeoutSeconds   int               `json:"timeoutSeconds"`
    MaxRetries       int               `json:"maxRetries"`
    EnableCaching    bool              `json:"enableCaching"`
}
```

### Generation Models

#### GenerationJob
```go
type GenerationJob struct {
    ID              string            `json:"id" db:"id"`
    ProjectID       string            `json:"projectId" db:"project_id"`
    Status          GenerationStatus  `json:"status" db:"status"`
    Progress        float64           `json:"progress" db:"progress"`
    CurrentStep     string            `json:"currentStep" db:"current_step"`
    StartTime       time.Time         `json:"startTime" db:"start_time"`
    EndTime         *time.Time        `json:"endTime" db:"end_time"`
    Results         *GenerationResults `json:"results" db:"results"`
    Errors          []GenerationError `json:"errors" db:"errors"`
    ProcessingTimeMs int64            `json:"processingTimeMs" db:"processing_time_ms"`
}

type GenerationStatus string

const (
    GenerationStatusStarted    GenerationStatus = "started"
    GenerationStatusParsing    GenerationStatus = "parsing"
    GenerationStatusMapping    GenerationStatus = "mapping"
    GenerationStatusGenerating GenerationStatus = "generating"
    GenerationStatusValidating GenerationStatus = "validating"
    GenerationStatusCompleted  GenerationStatus = "completed"
    GenerationStatusFailed     GenerationStatus = "failed"
    GenerationStatusCancelled  GenerationStatus = "cancelled"
)
```

#### GenerationResults
```go
type GenerationResults struct {
    ServerPath       string            `json:"serverPath"`
    GeneratedFiles   []GeneratedFile   `json:"generatedFiles"`
    MCPTools         []MCPTool         `json:"mcpTools"`
    Statistics       GenerationStats   `json:"statistics"`
    ValidationResult *ValidationResult `json:"validationResult"`
    BuildInfo        *BuildInfo        `json:"buildInfo"`
}

type GeneratedFile struct {
    Path         string    `json:"path"`
    RelativePath string    `json:"relativePath"`
    Type         FileType  `json:"type"`
    Size         int64     `json:"size"`
    LinesOfCode  int       `json:"linesOfCode"`
    Hash         string    `json:"hash"`
    CreatedAt    time.Time `json:"createdAt"`
}

type FileType string

const (
    FileTypeMain       FileType = "main"
    FileTypeHandler    FileType = "handler"
    FileTypeModel      FileType = "model"
    FileTypeConfig     FileType = "config"
    FileTypeTest       FileType = "test"
    FileTypeDoc        FileType = "documentation"
    FileTypeTemplate   FileType = "template"
    FileTypeAsset      FileType = "asset"
)

type GenerationStats struct {
    TotalEndpoints     int           `json:"totalEndpoints"`
    GeneratedTools     int           `json:"generatedTools"`
    ProcessingTime     time.Duration `json:"processingTime"`
    SpecComplexity     ComplexityLevel `json:"specComplexity"`
    TemplateVersion    string        `json:"templateVersion"`
    CodeLines          int           `json:"codeLines"`
    TestLines          int           `json:"testLines"`
    DocLines           int           `json:"docLines"`
}

type ComplexityLevel string

const (
    ComplexityLow    ComplexityLevel = "low"
    ComplexityMedium ComplexityLevel = "medium"
    ComplexityHigh   ComplexityLevel = "high"
)

type BuildInfo struct {
    GoVersion       string    `json:"goVersion"`
    MCPWeaverVersion string   `json:"mcpWeaverVersion"`
    BuildTime       time.Time `json:"buildTime"`
    Dependencies    []Dependency `json:"dependencies"`
}

type Dependency struct {
    Name    string `json:"name"`
    Version string `json:"version"`
    Source  string `json:"source"`
}
```

#### GenerationError
```go
type GenerationError struct {
    Type        ErrorType     `json:"type"`
    Code        string        `json:"code"`
    Message     string        `json:"message"`
    Details     string        `json:"details"`
    Suggestions []string      `json:"suggestions"`
    Location    *ErrorLocation `json:"location"`
    Timestamp   time.Time     `json:"timestamp"`
    Recoverable bool          `json:"recoverable"`
}

type ErrorType string

const (
    ErrorTypeValidation   ErrorType = "validation"
    ErrorTypeParsing      ErrorType = "parsing"
    ErrorTypeGeneration   ErrorType = "generation"
    ErrorTypeFileSystem   ErrorType = "filesystem"
    ErrorTypeNetwork      ErrorType = "network"
    ErrorTypeTemplate     ErrorType = "template"
    ErrorTypeInternal     ErrorType = "internal"
)

type ErrorLocation struct {
    File   string `json:"file"`
    Line   int    `json:"line"`
    Column int    `json:"column"`
    Path   string `json:"path"`
}
```

### OpenAPI Models

#### ParsedAPI
```go
type ParsedAPI struct {
    Document        *openapi3.T              `json:"document"`
    Title           string                   `json:"title"`
    Version         string                   `json:"version"`
    Description     string                   `json:"description"`
    BaseURL         string                   `json:"baseUrl"`
    Operations      []Operation              `json:"operations"`
    Schemas         map[string]*openapi3.SchemaRef `json:"schemas"`
    SecuritySchemes map[string]*openapi3.SecuritySchemeRef `json:"securitySchemes"`
    Servers         []ServerInfo             `json:"servers"`
    Tags            []TagInfo                `json:"tags"`
    ExternalDocs    *ExternalDocsInfo        `json:"externalDocs"`
}

type Operation struct {
    ID              string                    `json:"id"`
    Method          string                    `json:"method"`
    Path            string                    `json:"path"`
    Summary         string                    `json:"summary"`
    Description     string                    `json:"description"`
    Tags            []string                  `json:"tags"`
    Parameters      []Parameter               `json:"parameters"`
    RequestBody     *RequestBodyInfo          `json:"requestBody"`
    Responses       map[string]*ResponseInfo  `json:"responses"`
    Security        []SecurityRequirement     `json:"security"`
    Deprecated      bool                      `json:"deprecated"`
    OperationID     string                    `json:"operationId"`
}

type Parameter struct {
    Name            string           `json:"name"`
    In              string           `json:"in"`
    Description     string           `json:"description"`
    Required        bool             `json:"required"`
    Schema          *openapi3.SchemaRef `json:"schema"`
    Style           string           `json:"style"`
    Explode         bool             `json:"explode"`
    Example         interface{}      `json:"example"`
    Examples        map[string]*openapi3.ExampleRef `json:"examples"`
}

type RequestBodyInfo struct {
    Description string                           `json:"description"`
    Required    bool                             `json:"required"`
    Content     map[string]*openapi3.MediaType   `json:"content"`
}

type ResponseInfo struct {
    Description string                           `json:"description"`
    Headers     map[string]*openapi3.HeaderRef   `json:"headers"`
    Content     map[string]*openapi3.MediaType   `json:"content"`
    Links       map[string]*openapi3.LinkRef     `json:"links"`
}

type SecurityRequirement struct {
    Name   string   `json:"name"`
    Scopes []string `json:"scopes"`
}

type ServerInfo struct {
    URL         string                    `json:"url"`
    Description string                    `json:"description"`
    Variables   map[string]*ServerVariable `json:"variables"`
}

type ServerVariable struct {
    Default     string   `json:"default"`
    Description string   `json:"description"`
    Enum        []string `json:"enum"`
}

type TagInfo struct {
    Name         string            `json:"name"`
    Description  string            `json:"description"`
    ExternalDocs *ExternalDocsInfo `json:"externalDocs"`
}

type ExternalDocsInfo struct {
    Description string `json:"description"`
    URL         string `json:"url"`
}
```

### MCP Models

#### MCPTool
```go
type MCPTool struct {
    Name        string      `json:"name"`
    Description string      `json:"description"`
    InputSchema InputSchema `json:"inputSchema"`
    
    // Internal fields for generation
    Method      string      `json:"method"`
    Path        string      `json:"path"`
    BaseURL     string      `json:"baseUrl"`
    Operation   *Operation  `json:"operation"`
    Generated   bool        `json:"generated"`
}

type InputSchema struct {
    Type        string                     `json:"type"`
    Properties  map[string]*SchemaProperty `json:"properties"`
    Required    []string                   `json:"required"`
    Description string                     `json:"description"`
}

type SchemaProperty struct {
    Type        string                     `json:"type"`
    Description string                     `json:"description"`
    Format      string                     `json:"format,omitempty"`
    Enum        []interface{}              `json:"enum,omitempty"`
    Default     interface{}                `json:"default,omitempty"`
    Example     interface{}                `json:"example,omitempty"`
    Minimum     *float64                   `json:"minimum,omitempty"`
    Maximum     *float64                   `json:"maximum,omitempty"`
    MinLength   *int                       `json:"minLength,omitempty"`
    MaxLength   *int                       `json:"maxLength,omitempty"`
    Pattern     string                     `json:"pattern,omitempty"`
    Items       *SchemaProperty            `json:"items,omitempty"`
    Properties  map[string]*SchemaProperty `json:"properties,omitempty"`
    Required    []string                   `json:"required,omitempty"`
}
```

### Validation Models

#### ValidationResult
```go
type ValidationResult struct {
    Valid           bool                `json:"valid"`
    Errors          []ValidationError   `json:"errors"`
    Warnings        []ValidationWarning `json:"warnings"`
    Suggestions     []string            `json:"suggestions"`
    SpecInfo        *SpecInfo           `json:"specInfo"`
    ValidationTime  time.Duration       `json:"validationTime"`
    CacheHit        bool                `json:"cacheHit"`
    ValidatedAt     time.Time           `json:"validatedAt"`
}

type ValidationError struct {
    Type        string         `json:"type"`
    Code        string         `json:"code"`
    Message     string         `json:"message"`
    Path        string         `json:"path"`
    Line        int            `json:"line"`
    Column      int            `json:"column"`
    Severity    SeverityLevel  `json:"severity"`
    Location    *ErrorLocation `json:"location"`
    Context     string         `json:"context"`
    Suggestion  string         `json:"suggestion"`
}

type ValidationWarning struct {
    Type        string `json:"type"`
    Code        string `json:"code"`
    Message     string `json:"message"`
    Path        string `json:"path"`
    Suggestion  string `json:"suggestion"`
    Context     string `json:"context"`
}

type SeverityLevel string

const (
    SeverityError   SeverityLevel = "error"
    SeverityWarning SeverityLevel = "warning"
    SeverityInfo    SeverityLevel = "info"
)

type SpecInfo struct {
    Version         string              `json:"version"`
    Title           string              `json:"title"`
    Description     string              `json:"description"`
    OperationCount  int                 `json:"operationCount"`
    SchemaCount     int                 `json:"schemaCount"`
    SecuritySchemes []SecuritySchemeInfo `json:"securitySchemes"`
    Servers         []ServerInfo        `json:"servers"`
    Tags            []TagInfo           `json:"tags"`
    Complexity      ComplexityLevel     `json:"complexity"`
    EstimatedSize   string              `json:"estimatedSize"`
}

type SecuritySchemeInfo struct {
    Type         string `json:"type"`
    Name         string `json:"name"`
    Description  string `json:"description"`
    In           string `json:"in,omitempty"`
    Scheme       string `json:"scheme,omitempty"`
    BearerFormat string `json:"bearerFormat,omitempty"`
    OpenIDConnectURL string `json:"openIdConnectUrl,omitempty"`
}
```

### Template Models

#### Template
```go
type Template struct {
    ID          string             `json:"id" db:"id"`
    Name        string             `json:"name" db:"name"`
    Description string             `json:"description" db:"description"`
    Version     string             `json:"version" db:"version"`
    Author      string             `json:"author" db:"author"`
    Type        TemplateType       `json:"type" db:"type"`
    Path        string             `json:"path" db:"path"`
    IsBuiltIn   bool               `json:"isBuiltIn" db:"is_built_in"`
    Variables   []TemplateVariable `json:"variables" db:"variables"`
    Files       []TemplateFile     `json:"files"`
    CreatedAt   time.Time          `json:"createdAt" db:"created_at"`
    UpdatedAt   time.Time          `json:"updatedAt" db:"updated_at"`
}

type TemplateType string

const (
    TemplateTypeDefault TemplateType = "default"
    TemplateTypeCustom  TemplateType = "custom"
    TemplateTypePlugin  TemplateType = "plugin"
)

type TemplateVariable struct {
    Name         string      `json:"name"`
    Description  string      `json:"description"`
    Type         string      `json:"type"`
    DefaultValue interface{} `json:"defaultValue"`
    Required     bool        `json:"required"`
    Options      []string    `json:"options,omitempty"`
    Pattern      string      `json:"pattern,omitempty"`
    MinLength    int         `json:"minLength,omitempty"`
    MaxLength    int         `json:"maxLength,omitempty"`
}

type TemplateFile struct {
    Path        string `json:"path"`
    Type        string `json:"type"`
    Description string `json:"description"`
    Required    bool   `json:"required"`
    Executable  bool   `json:"executable"`
}
```

### Application Settings Models

#### AppSettings
```go
type AppSettings struct {
    Theme              string             `json:"theme"`
    Language           string             `json:"language"`
    AutoSave           bool               `json:"autoSave"`
    DefaultOutputPath  string             `json:"defaultOutputPath"`
    RecentProjects     []string           `json:"recentProjects"`
    WindowSettings     WindowSettings     `json:"windowSettings"`
    EditorSettings     EditorSettings     `json:"editorSettings"`
    GenerationSettings GenerationSettings `json:"generationSettings"`
    NotificationSettings NotificationSettings `json:"notificationSettings"`
}

type WindowSettings struct {
    Width       int  `json:"width"`
    Height      int  `json:"height"`
    Maximized   bool `json:"maximized"`
    X           int  `json:"x"`
    Y           int  `json:"y"`
    Theme       string `json:"theme"`
    Sidebar     bool `json:"sidebar"`
    StatusBar   bool `json:"statusBar"`
}

type EditorSettings struct {
    FontSize        int    `json:"fontSize"`
    FontFamily      string `json:"fontFamily"`
    TabSize         int    `json:"tabSize"`
    WordWrap        bool   `json:"wordWrap"`
    LineNumbers     bool   `json:"lineNumbers"`
    SyntaxHighlight bool   `json:"syntaxHighlight"`
    Theme           string `json:"theme"`
    AutoComplete    bool   `json:"autoComplete"`
    AutoFormat      bool   `json:"autoFormat"`
}

type GenerationSettings struct {
    DefaultTemplate        string   `json:"defaultTemplate"`
    EnableValidation       bool     `json:"enableValidation"`
    AutoOpenOutput         bool     `json:"autoOpenOutput"`
    ShowAdvancedOptions    bool     `json:"showAdvancedOptions"`
    BackupOnGenerate       bool     `json:"backupOnGenerate"`
    CustomTemplates        []string `json:"customTemplates"`
    ConcurrentGenerations  int      `json:"concurrentGenerations"`
    GenerationTimeout      int      `json:"generationTimeout"`
    EnableCaching          bool     `json:"enableCaching"`
    CacheExpiration        int      `json:"cacheExpiration"`
}

type NotificationSettings struct {
    EnableDesktop     bool   `json:"enableDesktop"`
    EnableSound       bool   `json:"enableSound"`
    EnableGeneration  bool   `json:"enableGeneration"`
    EnableValidation  bool   `json:"enableValidation"`
    EnableErrors      bool   `json:"enableErrors"`
    SoundFile         string `json:"soundFile"`
    Position          string `json:"position"`
    Duration          int    `json:"duration"`
}
```

### Event Models

#### Event System
```go
type Event struct {
    Type      string      `json:"type"`
    Data      interface{} `json:"data"`
    Timestamp time.Time   `json:"timestamp"`
    ID        string      `json:"id"`
}

type Notification struct {
    Type      NotificationType `json:"type"`
    Title     string          `json:"title"`
    Message   string          `json:"message"`
    Actions   []string        `json:"actions"`
    Duration  int             `json:"duration"`
    Timestamp time.Time       `json:"timestamp"`
    Sticky    bool            `json:"sticky"`
}

type NotificationType string

const (
    NotificationTypeInfo    NotificationType = "info"
    NotificationTypeSuccess NotificationType = "success"
    NotificationTypeWarning NotificationType = "warning"
    NotificationTypeError   NotificationType = "error"
)

type SystemError struct {
    Type        string    `json:"type"`
    Code        string    `json:"code"`
    Message     string    `json:"message"`
    Details     string    `json:"details"`
    Timestamp   time.Time `json:"timestamp"`
    Recoverable bool      `json:"recoverable"`
    Component   string    `json:"component"`
    UserAction  string    `json:"userAction"`
}
```

## Data Validation Rules

### Validation Tags
```go
// Common validation tags used throughout the application
const (
    ValidateRequired    = "required"
    ValidateMin         = "min"
    ValidateMax         = "max"
    ValidateEmail       = "email"
    ValidateURL         = "url"
    ValidateAlphaNum    = "alphanum"
    ValidateFileExists  = "file_exists"
    ValidateDirWritable = "dir_writable"
    ValidateOpenAPISpec = "openapi_spec"
)
```

### Custom Validation Functions
```go
func validateProjectName(name string) error {
    if len(name) < 1 || len(name) > 100 {
        return errors.New("project name must be between 1 and 100 characters")
    }
    if !regexp.MustCompile(`^[a-zA-Z0-9\s\-_]+$`).MatchString(name) {
        return errors.New("project name can only contain letters, numbers, spaces, hyphens, and underscores")
    }
    return nil
}

func validateOutputPath(path string) error {
    if !filepath.IsAbs(path) {
        return errors.New("output path must be absolute")
    }
    if info, err := os.Stat(path); err != nil {
        return fmt.Errorf("output path does not exist: %w", err)
    } else if !info.IsDir() {
        return errors.New("output path must be a directory")
    }
    return nil
}
```

This data models specification provides a comprehensive foundation for all data structures and database schemas used in MCPWeaver, ensuring consistency and type safety across the application.