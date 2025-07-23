package plugin

import (
	"context"
	"encoding/json"
	"time"

	"MCPWeaver/internal/mapping"
	"MCPWeaver/internal/parser"
)

// Plugin represents a loadable plugin
type Plugin interface {
	// GetInfo returns plugin metadata
	GetInfo() *PluginInfo
	
	// Initialize initializes the plugin with configuration
	Initialize(ctx context.Context, config json.RawMessage) error
	
	// Shutdown cleanly shuts down the plugin
	Shutdown(ctx context.Context) error
	
	// GetCapabilities returns what the plugin can do
	GetCapabilities() []Capability
}

// PluginInfo contains metadata about a plugin
type PluginInfo struct {
	ID          string            `json:"id"`
	Name        string            `json:"name"`
	Version     string            `json:"version"`
	Description string            `json:"description"`
	Author      string            `json:"author"`
	Homepage    string            `json:"homepage,omitempty"`
	Repository  string            `json:"repository,omitempty"`
	License     string            `json:"license"`
	Tags        []string          `json:"tags,omitempty"`
	MinVersion  string            `json:"minVersion"` // Minimum MCPWeaver version
	MaxVersion  string            `json:"maxVersion"` // Maximum MCPWeaver version
	Config      *PluginConfig     `json:"config,omitempty"`
	Permissions []Permission      `json:"permissions,omitempty"`
	Dependencies []Dependency     `json:"dependencies,omitempty"`
	Metadata    map[string]string `json:"metadata,omitempty"`
}

// PluginConfig describes plugin configuration schema
type PluginConfig struct {
	Schema      json.RawMessage   `json:"schema"`      // JSON Schema for config
	Default     json.RawMessage   `json:"default"`     // Default configuration
	Required    []string          `json:"required"`    // Required config fields
	Examples    []json.RawMessage `json:"examples"`    // Example configurations
}

// Capability represents what a plugin can do
type Capability string

const (
	CapabilityTemplateProcessor   Capability = "template_processor"
	CapabilityOutputConverter     Capability = "output_converter"
	CapabilityValidator           Capability = "validator"
	CapabilityUIComponent         Capability = "ui_component"
	CapabilityIntegration         Capability = "integration"
	CapabilityTesting             Capability = "testing"
	CapabilityParser              Capability = "parser"
	CapabilityGenerator           Capability = "generator"
	CapabilityMiddleware          Capability = "middleware"
)

// Permission represents what system resources a plugin needs
type Permission string

const (
	PermissionFileSystem     Permission = "filesystem"
	PermissionNetwork        Permission = "network"
	PermissionDatabase       Permission = "database"
	PermissionSettings       Permission = "settings"
	PermissionProjects       Permission = "projects"
	PermissionTemplates      Permission = "templates"
	PermissionExec           Permission = "exec"
	PermissionSystemInfo     Permission = "system_info"
	PermissionClipboard      Permission = "clipboard"
	PermissionNotifications  Permission = "notifications"
)

// Dependency represents plugin dependencies
type Dependency struct {
	Name       string `json:"name"`
	Version    string `json:"version"`
	Type       string `json:"type"` // "plugin", "system", "binary"
	Optional   bool   `json:"optional"`
	Repository string `json:"repository,omitempty"`
}

// PluginStatus represents the current state of a plugin
type PluginStatus string

const (
	PluginStatusUnloaded   PluginStatus = "unloaded"
	PluginStatusLoading    PluginStatus = "loading"
	PluginStatusLoaded     PluginStatus = "loaded"
	PluginStatusActive     PluginStatus = "active"
	PluginStatusError      PluginStatus = "error"
	PluginStatusDisabled   PluginStatus = "disabled"
	PluginStatusUnloading  PluginStatus = "unloading"
)

// PluginInstance represents an instance of a loaded plugin
type PluginInstance struct {
	Plugin     Plugin           `json:"-"`
	Info       *PluginInfo      `json:"info"`
	Status     PluginStatus     `json:"status"`
	Config     json.RawMessage  `json:"config,omitempty"`
	LoadedAt   time.Time        `json:"loadedAt"`
	LastError  string           `json:"lastError,omitempty"`
	Stats      *PluginStats     `json:"stats"`
	Manifest   *PluginManifest  `json:"manifest"`
}

// PluginStats tracks plugin usage statistics
type PluginStats struct {
	CallCount       int64         `json:"callCount"`
	TotalDuration   time.Duration `json:"totalDuration"`
	AverageDuration time.Duration `json:"averageDuration"`
	ErrorCount      int64         `json:"errorCount"`
	LastUsed        time.Time     `json:"lastUsed"`
	MemoryUsage     int64         `json:"memoryUsage"` // bytes
}

// PluginManifest describes plugin installation package
type PluginManifest struct {
	*PluginInfo
	Files       []PluginFile   `json:"files"`
	Checksum    string         `json:"checksum"`
	Size        int64          `json:"size"`
	InstallPath string         `json:"installPath,omitempty"`
	Verified    bool           `json:"verified"`
	Signature   string         `json:"signature,omitempty"`
}

// PluginFile describes a file in the plugin package
type PluginFile struct {
	Path     string `json:"path"`
	Size     int64  `json:"size"`
	Checksum string `json:"checksum"`
	Type     string `json:"type"` // "binary", "config", "template", "documentation"
	Platform string `json:"platform,omitempty"` // "windows", "darwin", "linux", "all"
	Arch     string `json:"arch,omitempty"`     // "amd64", "arm64", "all"
}

// TemplateProcessor plugins can process and modify templates
type TemplateProcessor interface {
	Plugin
	ProcessTemplate(ctx context.Context, template string, data map[string]interface{}) (string, error)
	GetSupportedFormats() []string
}

// OutputConverter plugins can convert generated output to different formats
type OutputConverter interface {
	Plugin
	ConvertOutput(ctx context.Context, input []byte, inputFormat, outputFormat string) ([]byte, error)
	GetSupportedFormats() []FormatPair
}

// FormatPair represents input/output format combination
type FormatPair struct {
	Input  string `json:"input"`
	Output string `json:"output"`
}

// Validator plugins can validate OpenAPI specs or generated code
type Validator interface {
	Plugin
	ValidateSpec(ctx context.Context, spec *parser.ParsedAPI) (*ValidationResult, error)
	ValidateGenerated(ctx context.Context, files map[string][]byte) (*ValidationResult, error)
	GetValidationRules() []ValidationRule
}

// ValidationResult represents validation output
type ValidationResult struct {
	Valid    bool               `json:"valid"`
	Errors   []ValidationError  `json:"errors,omitempty"`
	Warnings []ValidationError  `json:"warnings,omitempty"`
	Info     []ValidationError  `json:"info,omitempty"`
	Stats    *ValidationStats   `json:"stats,omitempty"`
}

// ValidationError represents a validation issue
type ValidationError struct {
	Type     string            `json:"type"`
	Message  string            `json:"message"`
	Path     string            `json:"path,omitempty"`
	Line     int               `json:"line,omitempty"`
	Column   int               `json:"column,omitempty"`
	Severity string            `json:"severity"`
	Code     string            `json:"code"`
	Fix      string            `json:"fix,omitempty"`
	Context  map[string]string `json:"context,omitempty"`
}

// ValidationRule describes a validation rule
type ValidationRule struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Severity    string `json:"severity"`
	Category    string `json:"category"`
	Enabled     bool   `json:"enabled"`
}

// ValidationStats provides validation statistics
type ValidationStats struct {
	TotalChecks   int           `json:"totalChecks"`
	Duration      time.Duration `json:"duration"`
	RulesApplied  int           `json:"rulesApplied"`
	FilesChecked  int           `json:"filesChecked"`
	LinesChecked  int           `json:"linesChecked"`
}

// UIComponent plugins can provide custom UI components
type UIComponent interface {
	Plugin
	GetComponentDefinition() *ComponentDefinition
	RenderComponent(ctx context.Context, props map[string]interface{}) (string, error)
}

// ComponentDefinition describes a UI component
type ComponentDefinition struct {
	Name        string                 `json:"name"`
	Type        string                 `json:"type"` // "react", "vue", "html"
	Props       map[string]interface{} `json:"props"`
	Events      []string               `json:"events"`
	Slots       []string               `json:"slots"`
	Styles      string                 `json:"styles,omitempty"`
	Script      string                 `json:"script,omitempty"`
	Template    string                 `json:"template"`
	Dependencies []string              `json:"dependencies,omitempty"`
}

// Integration plugins can integrate with external services
type Integration interface {
	Plugin
	Connect(ctx context.Context, config map[string]interface{}) error
	Disconnect(ctx context.Context) error
	GetConnectionStatus() ConnectionStatus
	ExecuteAction(ctx context.Context, action string, params map[string]interface{}) (interface{}, error)
	GetAvailableActions() []IntegrationAction
}

// ConnectionStatus represents integration connection state
type ConnectionStatus struct {
	Connected bool      `json:"connected"`
	LastCheck time.Time `json:"lastCheck"`
	Message   string    `json:"message,omitempty"`
	Latency   int64     `json:"latency"` // milliseconds
}

// IntegrationAction describes an available action
type IntegrationAction struct {
	ID          string                 `json:"id"`
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	Parameters  map[string]interface{} `json:"parameters"`
	Returns     string                 `json:"returns"`
	Category    string                 `json:"category"`
}

// Testing plugins can test generated MCP servers
type Testing interface {
	Plugin
	RunTests(ctx context.Context, serverPath string, config map[string]interface{}) (*TestResult, error)
	GetTestSuites() []TestSuite
	ValidateTestConfig(config map[string]interface{}) error
}

// TestResult represents test execution results
type TestResult struct {
	Passed       bool          `json:"passed"`
	Duration     time.Duration `json:"duration"`
	Tests        []TestCase    `json:"tests"`
	Coverage     *Coverage     `json:"coverage,omitempty"`
	Performance  *Performance  `json:"performance,omitempty"`
	Summary      string        `json:"summary"`
}

// TestCase represents an individual test
type TestCase struct {
	Name     string        `json:"name"`
	Status   string        `json:"status"` // "passed", "failed", "skipped"
	Duration time.Duration `json:"duration"`
	Message  string        `json:"message,omitempty"`
	Error    string        `json:"error,omitempty"`
}

// Coverage represents code coverage information
type Coverage struct {
	Lines     int     `json:"lines"`
	Covered   int     `json:"covered"`
	Percent   float64 `json:"percent"`
	Files     int     `json:"files"`
	Functions int     `json:"functions"`
}

// Performance represents performance metrics
type Performance struct {
	RequestsPerSecond float64       `json:"requestsPerSecond"`
	AverageLatency    time.Duration `json:"averageLatency"`
	MaxLatency        time.Duration `json:"maxLatency"`
	MinLatency        time.Duration `json:"minLatency"`
	MemoryUsage       int64         `json:"memoryUsage"`
	CPUUsage          float64       `json:"cpuUsage"`
}

// TestSuite describes a collection of tests
type TestSuite struct {
	ID          string            `json:"id"`
	Name        string            `json:"name"`
	Description string            `json:"description"`
	Tests       []string          `json:"tests"`
	Config      map[string]interface{} `json:"config,omitempty"`
	Tags        []string          `json:"tags,omitempty"`
}

// ParserPlugin plugins can extend OpenAPI parsing capabilities
type ParserPlugin interface {
	Plugin
	ParseExtension(ctx context.Context, spec map[string]interface{}, extension string) (interface{}, error)
	GetSupportedExtensions() []string
	ValidateExtension(extension string, value interface{}) error
}

// GeneratorPlugin plugins can extend code generation capabilities  
type GeneratorPlugin interface {
	Plugin
	GenerateCode(ctx context.Context, api *parser.ParsedAPI, tools []mapping.MCPTool, config map[string]interface{}) (map[string][]byte, error)
	GetSupportedLanguages() []string
	GetConfigSchema() json.RawMessage
}

// MiddlewarePlugin plugins can intercept and modify requests/responses
type MiddlewarePlugin interface {
	Plugin
	HandleRequest(ctx context.Context, req *PluginRequest, next MiddlewareHandler) (*PluginResponse, error)
	GetPriority() int
}

// MiddlewareHandler represents the next handler in the chain
type MiddlewareHandler func(ctx context.Context, req *PluginRequest) (*PluginResponse, error)

// PluginRequest represents a request that can be processed by middleware
type PluginRequest struct {
	ID      string                 `json:"id"`
	Type    string                 `json:"type"`
	Method  string                 `json:"method"`
	Path    string                 `json:"path"`
	Headers map[string]string      `json:"headers,omitempty"`
	Body    json.RawMessage        `json:"body,omitempty"`
	Context map[string]interface{} `json:"context,omitempty"`
}

// PluginResponse represents a response from middleware processing
type PluginResponse struct {
	Status  int                    `json:"status"`
	Headers map[string]string      `json:"headers,omitempty"`
	Body    json.RawMessage        `json:"body,omitempty"`
	Error   string                 `json:"error,omitempty"`
	Context map[string]interface{} `json:"context,omitempty"`
}

// PluginEvent represents events that plugins can emit or listen to
type PluginEvent struct {
	Type      string                 `json:"type"`
	Source    string                 `json:"source"`    // Plugin ID
	Target    string                 `json:"target"`    // Target plugin ID or "*" for broadcast
	Timestamp time.Time              `json:"timestamp"`
	Data      map[string]interface{} `json:"data,omitempty"`
	Context   context.Context        `json:"-"`
}

// EventHandler processes plugin events
type EventHandler func(ctx context.Context, event *PluginEvent) error

// EventListener plugins can listen to system events
type EventListener interface {
	Plugin
	HandleEvent(ctx context.Context, event *PluginEvent) error
	GetSubscriptions() []string
}

// EventEmitter plugins can emit custom events
type EventEmitter interface {
	Plugin
	EmitEvent(ctx context.Context, event *PluginEvent) error
	GetEmittedEvents() []string
}