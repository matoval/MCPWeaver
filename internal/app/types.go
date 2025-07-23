package app

import (
	"fmt"
	"time"

	"MCPWeaver/internal/mapping"
)

// Project Management Types
type CreateProjectRequest struct {
	Name       string          `json:"name" validate:"required,min=1,max=100"`
	SpecPath   string          `json:"specPath,omitempty"`
	SpecURL    string          `json:"specUrl,omitempty"`
	OutputPath string          `json:"outputPath" validate:"required"`
	Settings   ProjectSettings `json:"settings"`
}

func (c *CreateProjectRequest) convertValues() {
	// No time.Time fields to convert
}

type ProjectUpdateRequest struct {
	Name       *string          `json:"name,omitempty"`
	SpecPath   *string          `json:"specPath,omitempty"`
	SpecURL    *string          `json:"specUrl,omitempty"`
	OutputPath *string          `json:"outputPath,omitempty"`
	Settings   *ProjectSettings `json:"settings,omitempty"`
}

func (p *ProjectUpdateRequest) convertValues() {
	// No time.Time fields to convert
}

type Project struct {
	ID              string           `json:"id"`
	Name            string           `json:"name"`
	SpecPath        string           `json:"specPath"`
	SpecURL         string           `json:"specUrl"`
	OutputPath      string           `json:"outputPath"`
	Settings        ProjectSettings  `json:"settings"`
	Status          ProjectStatus    `json:"status"`
	CreatedAt       time.Time        `json:"createdAt"`
	UpdatedAt       time.Time        `json:"updatedAt"`
	LastGenerated   *time.Time       `json:"lastGenerated,omitempty"`
	GenerationCount int              `json:"generationCount"`
}

func (p *Project) convertValues() {
	// Convert time.Time fields to string representations for frontend
}

type ProjectSettings struct {
	PackageName     string   `json:"packageName" validate:"required,alphanum"`
	ServerPort      int      `json:"serverPort" validate:"min=1000,max=65535"`
	EnableLogging   bool     `json:"enableLogging"`
	LogLevel        string   `json:"logLevel" validate:"oneof=debug info warn error"`
	CustomTemplates []string `json:"customTemplates,omitempty"`
}

type ProjectStatus string

const (
	ProjectStatusCreated    ProjectStatus = "created"
	ProjectStatusValidating ProjectStatus = "validating"
	ProjectStatusReady      ProjectStatus = "ready"
	ProjectStatusGenerating ProjectStatus = "generating"
	ProjectStatusError      ProjectStatus = "error"
)

// Generation Types
type GenerationJob struct {
	ID          string             `json:"id"`
	ProjectID   string             `json:"projectId"`
	Status      GenerationStatus   `json:"status"`
	Progress    float64            `json:"progress"`
	CurrentStep string             `json:"currentStep"`
	StartTime   time.Time          `json:"startTime"`
	EndTime     *time.Time         `json:"endTime,omitempty"`
	Results     *GenerationResults `json:"results,omitempty"`
	Errors      []GenerationError  `json:"errors,omitempty"`
	Warnings    []string           `json:"warnings,omitempty"`
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

type GenerationResults struct {
	ServerPath     string            `json:"serverPath"`
	GeneratedFiles []GeneratedFile   `json:"generatedFiles"`
	MCPTools       []mapping.MCPTool `json:"mcpTools"`
	Statistics     GenerationStats   `json:"statistics"`
}

type GeneratedFile struct {
	Path        string `json:"path"`
	Type        string `json:"type"`
	Size        int    `json:"size"`
	LinesOfCode int    `json:"linesOfCode"`
}

type GenerationStats struct {
	TotalEndpoints  int           `json:"totalEndpoints"`
	GeneratedTools  int           `json:"generatedTools"`
	ProcessingTime  time.Duration `json:"processingTime"`
	SpecComplexity  string        `json:"specComplexity"`
	TemplateVersion string        `json:"templateVersion"`
}

type GenerationError struct {
	Type        string         `json:"type"`
	Message     string         `json:"message"`
	Details     string         `json:"details,omitempty"`
	Suggestions []string       `json:"suggestions,omitempty"`
	Location    *ErrorLocation `json:"location,omitempty"`
}

type ErrorLocation struct {
	File   string `json:"file"`
	Line   int    `json:"line"`
	Column int    `json:"column"`
}

// Validation Types
type ValidationResult struct {
	Valid          bool                `json:"valid"`
	Errors         []ValidationError   `json:"errors"`
	Warnings       []ValidationWarning `json:"warnings"`
	Suggestions    []string            `json:"suggestions"`
	SpecInfo       *SpecInfo           `json:"specInfo,omitempty"`
	ValidationTime time.Duration       `json:"validationTime"`
	CacheHit       bool                `json:"cacheHit"`
	ValidatedAt    time.Time           `json:"validatedAt"`
}

func (v *ValidationResult) convertValues() {
	// Convert time.Time fields to string representations for frontend
}

type ValidationError struct {
	Type     string         `json:"type"`
	Message  string         `json:"message"`
	Path     string         `json:"path"`
	Line     int            `json:"line,omitempty"`
	Column   int            `json:"column,omitempty"`
	Severity string         `json:"severity"`
	Code     string         `json:"code"`
	Location *ErrorLocation `json:"location,omitempty"`
}

type ValidationWarning struct {
	Type       string `json:"type"`
	Message    string `json:"message"`
	Path       string `json:"path"`
	Suggestion string `json:"suggestion"`
}

type SpecInfo struct {
	Version         string            `json:"version"`
	Title           string            `json:"title"`
	Description     string            `json:"description"`
	OperationCount  int               `json:"operationCount"`
	SchemaCount     int               `json:"schemaCount"`
	SecuritySchemes []SecurityScheme  `json:"securitySchemes"`
	Servers         []ServerInfo      `json:"servers"`
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

// File System Types
type FileFilter struct {
	DisplayName string   `json:"displayName"`
	Pattern     string   `json:"pattern"`
	Extensions  []string `json:"extensions"`
}

// Settings Types
type AppSettings struct {
	Theme                string                `json:"theme"`
	Language             string                `json:"language"`
	AutoSave             bool                  `json:"autoSave"`
	DefaultOutputPath    string                `json:"defaultOutputPath"`
	RecentProjects       []string              `json:"recentProjects"`
	RecentFiles          []string              `json:"recentFiles"`
	WindowSettings       WindowSettings        `json:"windowSettings"`
	EditorSettings       EditorSettings        `json:"editorSettings"`
	GenerationSettings   GenerationSettings    `json:"generationSettings"`
	NotificationSettings NotificationSettings  `json:"notificationSettings"`
	AppearanceSettings   AppearanceSettings    `json:"appearanceSettings"`
	UpdateSettings       UpdateSettings        `json:"updateSettings"`
}

type WindowSettings struct {
	Width     int  `json:"width"`
	Height    int  `json:"height"`
	Maximized bool `json:"maximized"`
	X         int  `json:"x"`
	Y         int  `json:"y"`
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
	PerformanceMode     bool     `json:"performanceMode"`
	MaxWorkers          int      `json:"maxWorkers"`
}

type NotificationSettings struct {
	EnableDesktopNotifications bool   `json:"enableDesktopNotifications"`
	EnableSoundNotifications   bool   `json:"enableSoundNotifications"`
	NotificationPosition       string `json:"notificationPosition"`
	NotificationDuration       int    `json:"notificationDuration"`
	SoundVolume                float64 `json:"soundVolume"`
	ShowGenerationProgress     bool   `json:"showGenerationProgress"`
	ShowErrorNotifications     bool   `json:"showErrorNotifications"`
	ShowSuccessNotifications   bool   `json:"showSuccessNotifications"`
}

type AppearanceSettings struct {
	UITheme            string  `json:"uiTheme"`
	AccentColor        string  `json:"accentColor"`
	WindowOpacity      float64 `json:"windowOpacity"`
	ShowAnimation      bool    `json:"showAnimation"`
	ReducedMotion      bool    `json:"reducedMotion"`
	FontScale          float64 `json:"fontScale"`
	CompactMode        bool    `json:"compactMode"`
	ShowSidebar        bool    `json:"showSidebar"`
	SidebarPosition    string  `json:"sidebarPosition"`
	ShowStatusBar      bool    `json:"showStatusBar"`
	ShowToolbar        bool    `json:"showToolbar"`
}

// Template Types
type Template struct {
	ID          string             `json:"id"`
	Name        string             `json:"name"`
	Description string             `json:"description"`
	Version     string             `json:"version"`
	Author      string             `json:"author"`
	Type        TemplateType       `json:"type"`
	Path        string             `json:"path"`
	IsBuiltIn   bool               `json:"isBuiltIn"`
	Variables   []TemplateVariable `json:"variables"`
	CreatedAt   time.Time          `json:"createdAt"`
	UpdatedAt   time.Time          `json:"updatedAt"`
}

type TemplateType string

const (
	TemplateTypeDefault TemplateType = "default"
	TemplateTypeCustom  TemplateType = "custom"
	TemplateTypePlugin  TemplateType = "plugin"
)

type TemplateVariable struct {
	Name         string   `json:"name"`
	Description  string   `json:"description"`
	Type         string   `json:"type"`
	DefaultValue string   `json:"defaultValue"`
	Required     bool     `json:"required"`
	Options      []string `json:"options,omitempty"` // For enum/select type variables
	Validation   string   `json:"validation,omitempty"` // Validation regex or rules
}

// Template Management Types
type CreateTemplateRequest struct {
	Name        string             `json:"name" validate:"required,min=1,max=100"`
	Description string             `json:"description" validate:"max=500"`
	Version     string             `json:"version" validate:"required,semver"`
	Author      string             `json:"author" validate:"max=100"`
	Type        TemplateType       `json:"type" validate:"required"`
	Path        string             `json:"path" validate:"required"`
	Variables   []TemplateVariable `json:"variables,omitempty"`
}

type UpdateTemplateRequest struct {
	Name        *string             `json:"name,omitempty"`
	Description *string             `json:"description,omitempty"`
	Version     *string             `json:"version,omitempty"`
	Author      *string             `json:"author,omitempty"`
	Type        *TemplateType       `json:"type,omitempty"`
	Path        *string             `json:"path,omitempty"`
	Variables   *[]TemplateVariable `json:"variables,omitempty"`
}

type TemplateValidationResult struct {
	Valid        bool                    `json:"valid"`
	Errors       []TemplateError         `json:"errors,omitempty"`
	Warnings     []TemplateWarning       `json:"warnings,omitempty"`
	Suggestions  []string                `json:"suggestions,omitempty"`
	Performance  *TemplatePerformance    `json:"performance,omitempty"`
	Dependencies []TemplateDependency    `json:"dependencies,omitempty"`
}

type TemplateError struct {
	Type     string `json:"type"`
	Message  string `json:"message"`
	Line     int    `json:"line,omitempty"`
	Column   int    `json:"column,omitempty"`
	Severity string `json:"severity"`
}

type TemplateWarning struct {
	Type       string `json:"type"`
	Message    string `json:"message"`
	Line       int    `json:"line,omitempty"`
	Suggestion string `json:"suggestion,omitempty"`
}

type TemplatePerformance struct {
	RenderTime    time.Duration `json:"renderTime"`
	MemoryUsage   int64         `json:"memoryUsage"`
	Complexity    string        `json:"complexity"` // "low", "medium", "high"
	CacheHit      bool          `json:"cacheHit"`
}

type TemplateDependency struct {
	Name     string `json:"name"`
	Version  string `json:"version"`
	Required bool   `json:"required"`
	Type     string `json:"type"` // "system", "library", "template"
}

type TemplateImportRequest struct {
	Source      string            `json:"source"` // "file", "url", "marketplace"
	Path        string            `json:"path,omitempty"`
	URL         string            `json:"url,omitempty"`
	MarketplaceID string          `json:"marketplaceId,omitempty"`
	ImportOptions TemplateImportOptions `json:"options,omitempty"`
}

type TemplateImportOptions struct {
	OverwriteExisting bool     `json:"overwriteExisting"`
	ValidateOnly      bool     `json:"validateOnly"`
	IncludeDependencies bool   `json:"includeDependencies"`
	TargetType        TemplateType `json:"targetType,omitempty"`
}

type TemplateExportRequest struct {
	TemplateID    string                `json:"templateId"`
	Format        string                `json:"format"` // "zip", "tar", "single"
	TargetPath    string                `json:"targetPath"`
	ExportOptions TemplateExportOptions `json:"options,omitempty"`
}

type TemplateExportOptions struct {
	IncludeDocumentation bool `json:"includeDocumentation"`
	IncludeExamples     bool `json:"includeExamples"`
	IncludeDependencies bool `json:"includeDependencies"`
	Minify              bool `json:"minify"`
}

type TemplateTestRequest struct {
	TemplateID   string                 `json:"templateId"`
	TestData     map[string]interface{} `json:"testData"`
	TestOptions  TemplateTestOptions    `json:"options,omitempty"`
}

type TemplateTestOptions struct {
	ValidateOutput    bool `json:"validateOutput"`
	MeasurePerformance bool `json:"measurePerformance"`
	GenerateReport    bool `json:"generateReport"`
}

type TemplateTestResult struct {
	Success     bool                    `json:"success"`
	Output      string                  `json:"output,omitempty"`
	Errors      []TemplateError         `json:"errors,omitempty"`
	Warnings    []TemplateWarning       `json:"warnings,omitempty"`
	Performance *TemplatePerformance    `json:"performance,omitempty"`
	Report      *TemplateTestReport     `json:"report,omitempty"`
}

type TemplateTestReport struct {
	TemplateID      string        `json:"templateId"`
	TestExecutedAt  time.Time     `json:"testExecutedAt"`
	ExecutionTime   time.Duration `json:"executionTime"`
	OutputSize      int64         `json:"outputSize"`
	VariablesUsed   []string      `json:"variablesUsed"`
	FunctionsUsed   []string      `json:"functionsUsed"`
	Recommendations []string      `json:"recommendations"`
}

// Template Marketplace Types - See marketplace.go for actual definitions

// Error Types
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
	RetryAfter    *time.Duration    `json:"retryAfter,omitempty"`
	Context       *ErrorContext     `json:"context,omitempty"`
}

// Error implements the error interface for APIError
func (e *APIError) Error() string {
	return e.Message
}

// IsRetryable returns true if the error can be retried
func (e *APIError) IsRetryable() bool {
	return e.Recoverable && e.RetryAfter != nil
}

// GetRetryDelay returns the suggested retry delay
func (e *APIError) GetRetryDelay() time.Duration {
	if e.RetryAfter != nil {
		return *e.RetryAfter
	}
	return 0
}

// ErrorSeverity defines the severity level of an error
type ErrorSeverity string

const (
	ErrorSeverityLow      ErrorSeverity = "low"
	ErrorSeverityMedium   ErrorSeverity = "medium"
	ErrorSeverityHigh     ErrorSeverity = "high"
	ErrorSeverityCritical ErrorSeverity = "critical"
)

// ErrorContext provides additional context about where an error occurred
type ErrorContext struct {
	Operation   string            `json:"operation"`
	Component   string            `json:"component"`
	ProjectID   string            `json:"projectId,omitempty"`
	UserID      string            `json:"userId,omitempty"`
	SessionID   string            `json:"sessionId,omitempty"`
	RequestID   string            `json:"requestId,omitempty"`
	StackTrace  string            `json:"stackTrace,omitempty"`
	Metadata    map[string]string `json:"metadata,omitempty"`
}

// ErrorCollection aggregates multiple errors for batch operations
type ErrorCollection struct {
	Errors      []APIError `json:"errors"`
	Warnings    []APIError `json:"warnings"`
	Operation   string     `json:"operation"`
	TotalItems  int        `json:"totalItems"`
	FailedItems int        `json:"failedItems"`
	Timestamp   time.Time  `json:"timestamp"`
}

// HasErrors returns true if there are any errors
func (ec *ErrorCollection) HasErrors() bool {
	return len(ec.Errors) > 0
}

// HasWarnings returns true if there are any warnings
func (ec *ErrorCollection) HasWarnings() bool {
	return len(ec.Warnings) > 0
}

// Error implements the error interface for ErrorCollection
func (ec *ErrorCollection) Error() string {
	if len(ec.Errors) == 0 {
		return "no errors"
	}
	if len(ec.Errors) == 1 {
		return ec.Errors[0].Error()
	}
	return fmt.Sprintf("%d errors occurred during %s", len(ec.Errors), ec.Operation)
}

// RetryPolicy defines retry behavior for operations
type RetryPolicy struct {
	MaxRetries      int           `json:"maxRetries"`
	InitialDelay    time.Duration `json:"initialDelay"`
	MaxDelay        time.Duration `json:"maxDelay"`
	BackoffMultiplier float64     `json:"backoffMultiplier"`
	JitterEnabled   bool          `json:"jitterEnabled"`
	RetryableErrors []string      `json:"retryableErrors"`
}

// DefaultRetryPolicy returns a default retry policy
func DefaultRetryPolicy() RetryPolicy {
	return RetryPolicy{
		MaxRetries:        3,
		InitialDelay:      time.Second,
		MaxDelay:          30 * time.Second,
		BackoffMultiplier: 2.0,
		JitterEnabled:     true,
		RetryableErrors: []string{
			ErrCodeNetworkError,
			ErrCodeInternalError,
			ErrCodeDatabaseError,
		},
	}
}

// Error constants
const (
	ErrCodeValidation      = "VALIDATION_ERROR"
	ErrCodeNotFound        = "NOT_FOUND"
	ErrCodeInternalError   = "INTERNAL_ERROR"
	ErrCodeFileAccess      = "FILE_ACCESS_ERROR"
	ErrCodeNetworkError    = "NETWORK_ERROR"
	ErrCodeParsingError    = "PARSING_ERROR"
	ErrCodeGenerationError = "GENERATION_ERROR"
	ErrCodeDatabaseError   = "DATABASE_ERROR"
	ErrCodePermissionError = "PERMISSION_ERROR"
	ErrCodeRateLimitError  = "RATE_LIMIT_ERROR"
	ErrCodeTimeoutError    = "TIMEOUT_ERROR"
	ErrCodeConfigError     = "CONFIG_ERROR"
	ErrCodeAuthError       = "AUTH_ERROR"
)

// Error type constants
const (
	ErrorTypeValidation   = "validation"
	ErrorTypeSystem       = "system"
	ErrorTypeNetwork      = "network"
	ErrorTypeFileSystem   = "filesystem"
	ErrorTypeDatabase     = "database"
	ErrorTypeGeneration   = "generation"
	ErrorTypePermission   = "permission"
	ErrorTypeConfiguration = "configuration"
	ErrorTypeAuthentication = "authentication"
)

// Event Types
type GenerationProgress struct {
	JobID       string  `json:"jobId"`
	Progress    float64 `json:"progress"`
	CurrentStep string  `json:"currentStep"`
	Message     string  `json:"message"`
}

type Notification struct {
	Type    string   `json:"type"`
	Title   string   `json:"title"`
	Message string   `json:"message"`
	Actions []string `json:"actions,omitempty"`
}

type SystemError struct {
	Type        string    `json:"type"`
	Message     string    `json:"message"`
	Timestamp   time.Time `json:"timestamp"`
	Recoverable bool      `json:"recoverable"`
}

// Import/Export Types
type ImportResult struct {
	Content      string    `json:"content"`
	Valid        bool      `json:"valid"`
	SpecInfo     *SpecInfo `json:"specInfo,omitempty"`
	Errors       []string  `json:"errors,omitempty"`
	Warnings     []string  `json:"warnings,omitempty"`
	ImportedFrom string    `json:"importedFrom"` // "file" or "url"
	FilePath     string    `json:"filePath,omitempty"`
	SourceURL    string    `json:"sourceUrl,omitempty"`
	FileSize     int64     `json:"fileSize"`
	ImportedAt   time.Time `json:"importedAt"`
}

type ExportResult struct {
	ProjectID     string         `json:"projectId"`
	ProjectName   string         `json:"projectName"`
	TargetDir     string         `json:"targetDir"`
	ExportedFiles []ExportedFile `json:"exportedFiles"`
	TotalFiles    int            `json:"totalFiles"`
	TotalSize     int64          `json:"totalSize"`
	ExportedAt    time.Time      `json:"exportedAt"`
}

type ExportedFile struct {
	Name         string    `json:"name"`
	Path         string    `json:"path"`
	Size         int64     `json:"size"`
	ModifiedTime time.Time `json:"modifiedTime"`
}

type FileOperationProgress struct {
	OperationID        string `json:"operationId"`
	Type               string `json:"type"` // "import" or "export"
	Progress           int    `json:"progress"`
	CurrentFile        string `json:"currentFile"`
	TotalFiles         int    `json:"totalFiles"`
	ProcessedFiles     int    `json:"processedFiles"`
	StartTime          string `json:"startTime"`
	ElapsedTime        int64  `json:"elapsedTime"`
	EstimatedRemaining int64  `json:"estimatedRemaining"`
}

type RecentFile struct {
	Path         string `json:"path"`
	Name         string `json:"name"`
	Size         int64  `json:"size"`
	LastAccessed string `json:"lastAccessed"`
	Type         string `json:"type"` // "spec" or "export"
}

// Activity Log and Monitoring Types
type ActivityLogEntry struct {
	ID          string                 `json:"id"`
	Timestamp   time.Time              `json:"timestamp"`
	Level       LogLevel               `json:"level"`
	Component   string                 `json:"component"`
	Operation   string                 `json:"operation"`
	Message     string                 `json:"message"`
	Details     string                 `json:"details,omitempty"`
	Duration    *time.Duration         `json:"duration,omitempty"`
	ProjectID   string                 `json:"projectId,omitempty"`
	UserAction  bool                   `json:"userAction"`
	Metadata    map[string]interface{} `json:"metadata,omitempty"`
}

type LogFilter struct {
	Level      *LogLevel `json:"level,omitempty"`
	Component  *string   `json:"component,omitempty"`
	Operation  *string   `json:"operation,omitempty"`
	ProjectID  *string   `json:"projectId,omitempty"`
	UserAction *bool     `json:"userAction,omitempty"`
	StartTime  *time.Time `json:"startTime,omitempty"`
	EndTime    *time.Time `json:"endTime,omitempty"`
	Search     *string   `json:"search,omitempty"`
	Limit      *int      `json:"limit,omitempty"`
}

type ApplicationStatus struct {
	Status           StatusLevel  `json:"status"`
	Message          string       `json:"message"`
	ActiveOperations int          `json:"activeOperations"`
	LastUpdate       time.Time    `json:"lastUpdate"`
	SystemHealth     SystemHealth `json:"systemHealth"`
}

type StatusLevel string

const (
	StatusIdle    StatusLevel = "idle"
	StatusWorking StatusLevel = "working"  
	StatusError   StatusLevel = "error"
	StatusWarning StatusLevel = "warning"
)

type SystemHealth struct {
	MemoryUsage       float64 `json:"memoryUsage"`     // MB
	CPUUsage          float64 `json:"cpuUsage"`        // Percentage
	DiskSpace         float64 `json:"diskSpace"`       // GB available
	DatabaseSize      float64 `json:"databaseSize"`    // MB
	TemporaryFiles    int     `json:"temporaryFiles"`  // Count
	ActiveConnections int     `json:"activeConnections"`
}

type ErrorReport struct {
	ID           string        `json:"id"`
	Timestamp    time.Time     `json:"timestamp"`
	Type         ErrorType     `json:"type"`
	Severity     ErrorSeverity `json:"severity"`
	Component    string        `json:"component"`
	Operation    string        `json:"operation"`
	Message      string        `json:"message"`
	Details      string        `json:"details,omitempty"`
	StackTrace   string        `json:"stackTrace,omitempty"`
	UserContext  UserContext   `json:"userContext"`
	SystemInfo   SystemInfo    `json:"systemInfo"`
	Recovery     RecoveryInfo  `json:"recovery"`
	Frequency    int           `json:"frequency"`
	FirstSeen    time.Time     `json:"firstSeen"`
	LastSeen     time.Time     `json:"lastSeen"`
}

type ErrorType string

const (
	ErrorTypeValidationErr   ErrorType = "validation"
	ErrorTypeSystemErr       ErrorType = "system"
	ErrorTypeNetworkErr      ErrorType = "network"
	ErrorTypeFileSystemErr   ErrorType = "filesystem"
	ErrorTypeDatabaseErr     ErrorType = "database"
	ErrorTypeGenerationErr   ErrorType = "generation"
	ErrorTypePermissionErr   ErrorType = "permission"
	ErrorTypeConfigurationErr ErrorType = "configuration"
	ErrorTypeAuthenticationErr ErrorType = "authentication"
)

type UserContext struct {
	ProjectID     string            `json:"projectId,omitempty"`
	ProjectName   string            `json:"projectName,omitempty"`
	UserAction    string            `json:"userAction,omitempty"`
	UIState       string            `json:"uiState,omitempty"`
	RecentActions []string          `json:"recentActions,omitempty"`
	Settings      map[string]string `json:"settings,omitempty"`
}

type SystemInfo struct {
	OS             string  `json:"os"`
	Architecture   string  `json:"architecture"`
	GoVersion      string  `json:"goVersion"`
	AppVersion     string  `json:"appVersion"`
	MemoryMB       float64 `json:"memoryMB"`
	CPUUsage       float64 `json:"cpuUsage"`
	DiskSpaceGB    float64 `json:"diskSpaceGB"`
	DatabaseSizeMB float64 `json:"databaseSizeMB"`
}

type RecoveryInfo struct {
	Attempted       bool          `json:"attempted"`
	Successful      bool          `json:"successful"`
	Method          string        `json:"method,omitempty"`
	Duration        time.Duration `json:"duration,omitempty"`
	UserInteraction bool          `json:"userInteraction"`
	DataLoss        bool          `json:"dataLoss"`
}

type LogConfig struct {
	Level         LogLevel      `json:"level"`
	BufferSize    int           `json:"bufferSize"`
	RetentionDays int           `json:"retentionDays"`
	EnableConsole bool          `json:"enableConsole"`
	EnableBuffer  bool          `json:"enableBuffer"`
	FlushInterval time.Duration `json:"flushInterval"`
}

type LogSearchRequest struct {
	Query     string     `json:"query"`
	Filter    LogFilter  `json:"filter"`
	Limit     int        `json:"limit"`
	Offset    int        `json:"offset"`
}

type LogSearchResult struct {
	Entries    []ActivityLogEntry `json:"entries"`
	Total      int                `json:"total"`
	HasMore    bool               `json:"hasMore"`
	SearchTime time.Duration      `json:"searchTime"`
}

type LogExportRequest struct {
	Filter   LogFilter `json:"filter"`
	Format   string    `json:"format"` // "json", "csv", "txt"
	FilePath string    `json:"filePath"`
}

type LogExportResult struct {
	FilePath     string        `json:"filePath"`
	EntriesCount int           `json:"entriesCount"`
	FileSize     int64         `json:"fileSize"`
	ExportTime   time.Duration `json:"exportTime"`
	Format       string        `json:"format"`
}