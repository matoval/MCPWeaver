package app

import (
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

type ProjectUpdateRequest struct {
	Name       *string          `json:"name,omitempty"`
	SpecPath   *string          `json:"specPath,omitempty"`
	SpecURL    *string          `json:"specUrl,omitempty"`
	OutputPath *string          `json:"outputPath,omitempty"`
	Settings   *ProjectSettings `json:"settings,omitempty"`
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
	Theme              string             `json:"theme"`
	Language           string             `json:"language"`
	AutoSave           bool               `json:"autoSave"`
	DefaultOutputPath  string             `json:"defaultOutputPath"`
	RecentProjects     []string           `json:"recentProjects"`
	WindowSettings     WindowSettings     `json:"windowSettings"`
	EditorSettings     EditorSettings     `json:"editorSettings"`
	GenerationSettings GenerationSettings `json:"generationSettings"`
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
	Name         string `json:"name"`
	Description  string `json:"description"`
	Type         string `json:"type"`
	DefaultValue string `json:"defaultValue"`
	Required     bool   `json:"required"`
}

// Error Types
type APIError struct {
	Type        string            `json:"type"`
	Code        string            `json:"code"`
	Message     string            `json:"message"`
	Details     map[string]string `json:"details,omitempty"`
	Timestamp   time.Time         `json:"timestamp"`
	Suggestions []string          `json:"suggestions,omitempty"`
}

// Error implements the error interface for APIError
func (e *APIError) Error() string {
	return e.Message
}

// Error constants
const (
	ErrCodeValidation    = "VALIDATION_ERROR"
	ErrCodeNotFound      = "NOT_FOUND"
	ErrCodeInternalError = "INTERNAL_ERROR"
	ErrCodeFileAccess    = "FILE_ACCESS_ERROR"
	ErrCodeNetworkError  = "NETWORK_ERROR"
	ErrCodeParsingError  = "PARSING_ERROR"
	ErrCodeGenerationError = "GENERATION_ERROR"
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