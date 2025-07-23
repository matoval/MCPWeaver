package plugin

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
)

// ExampleTemplateProcessor demonstrates a simple template processor plugin
type ExampleTemplateProcessor struct {
	config map[string]interface{}
}

// GetInfo returns plugin metadata
func (p *ExampleTemplateProcessor) GetInfo() *PluginInfo {
	return &PluginInfo{
		ID:          "example-template-processor",
		Name:        "Example Template Processor",
		Version:     "1.0.0",
		Description: "A simple example template processor that demonstrates the plugin system",
		Author:      "MCPWeaver Team",
		License:     "MIT",
		Tags:        []string{"template", "example", "demo"},
		MinVersion:  "1.0.0",
		MaxVersion:  "2.0.0",
		Permissions: []Permission{PermissionTemplates},
		Config: &PluginConfig{
			Schema:   json.RawMessage(`{"type": "object", "properties": {"prefix": {"type": "string"}}}`),
			Default:  json.RawMessage(`{"prefix": "Generated"}`),
			Required: []string{},
			Examples: []json.RawMessage{json.RawMessage(`{"prefix": "Custom"}`)},
		},
	}
}

// Initialize initializes the plugin with configuration
func (p *ExampleTemplateProcessor) Initialize(ctx context.Context, config json.RawMessage) error {
	if config != nil {
		if err := json.Unmarshal(config, &p.config); err != nil {
			return fmt.Errorf("failed to parse config: %w", err)
		}
	} else {
		p.config = map[string]interface{}{
			"prefix": "Generated",
		}
	}
	return nil
}

// Shutdown cleanly shuts down the plugin
func (p *ExampleTemplateProcessor) Shutdown(ctx context.Context) error {
	p.config = nil
	return nil
}

// GetCapabilities returns the capabilities this plugin provides
func (p *ExampleTemplateProcessor) GetCapabilities() []Capability {
	return []Capability{CapabilityTemplateProcessor}
}

// ProcessTemplate processes a template with the given data
func (p *ExampleTemplateProcessor) ProcessTemplate(ctx context.Context, template string, data map[string]interface{}) (string, error) {
	prefix := "Generated"
	if p.config != nil {
		if configPrefix, ok := p.config["prefix"].(string); ok {
			prefix = configPrefix
		}
	}
	
	// Simple template processing - replace {{key}} with values from data
	result := template
	for key, value := range data {
		placeholder := fmt.Sprintf("{{%s}}", key)
		replacement := fmt.Sprintf("%v", value)
		result = strings.ReplaceAll(result, placeholder, replacement)
	}
	
	// Add prefix comment
	if strings.TrimSpace(result) != "" {
		result = fmt.Sprintf("// %s by Example Template Processor\n%s", prefix, result)
	}
	
	return result, nil
}

// ExampleValidator demonstrates a simple validator plugin
type ExampleValidator struct {
	config map[string]interface{}
}

// GetInfo returns plugin metadata
func (p *ExampleValidator) GetInfo() *PluginInfo {
	return &PluginInfo{
		ID:          "example-validator",
		Name:        "Example Validator",
		Version:     "1.0.0",
		Description: "A simple example validator that demonstrates validation plugins",
		Author:      "MCPWeaver Team",
		License:     "MIT",
		Tags:        []string{"validator", "example", "demo"},
		MinVersion:  "1.0.0",
		MaxVersion:  "2.0.0",
		Permissions: []Permission{PermissionProjects},
		Config: &PluginConfig{
			Schema:   json.RawMessage(`{"type": "object", "properties": {"strict": {"type": "boolean"}}}`),
			Default:  json.RawMessage(`{"strict": false}`),
			Required: []string{},
			Examples: []json.RawMessage{json.RawMessage(`{"strict": true}`)},
		},
	}
}

// Initialize initializes the plugin with configuration
func (p *ExampleValidator) Initialize(ctx context.Context, config json.RawMessage) error {
	if config != nil {
		if err := json.Unmarshal(config, &p.config); err != nil {
			return fmt.Errorf("failed to parse config: %w", err)
		}
	} else {
		p.config = map[string]interface{}{
			"strict": false,
		}
	}
	return nil
}

// Shutdown cleanly shuts down the plugin
func (p *ExampleValidator) Shutdown(ctx context.Context) error {
	p.config = nil
	return nil
}

// GetCapabilities returns the capabilities this plugin provides
func (p *ExampleValidator) GetCapabilities() []Capability {
	return []Capability{CapabilityValidator}
}

// SimpleSpec represents a simplified OpenAPI spec for example validation
type SimpleSpec struct {
	Title   string `json:"title"`
	Version string `json:"version"`
}

// ValidateSpec validates an OpenAPI specification
func (p *ExampleValidator) ValidateSpec(ctx context.Context, spec interface{}) (*ValidationResult, error) {
	result := &ValidationResult{
		Valid:    true,
		Errors:   []ValidationError{},
		Warnings: []ValidationError{},
		Info:     []ValidationError{},
		Stats: &ValidationStats{
			TotalChecks:  3,
			RulesApplied: 3,
			FilesChecked: 1,
		},
	}
	
	strict := false
	if p.config != nil {
		if configStrict, ok := p.config["strict"].(bool); ok {
			strict = configStrict
		}
	}
	
	// Example validation rules
	if spec == nil {
		result.Valid = false
		result.Errors = append(result.Errors, ValidationError{
			Type:     "null_spec",
			Message:  "Specification is null",
			Severity: "error",
			Code:     "SPEC_NULL",
		})
		return result, nil
	}
	
	// Try to parse as SimpleSpec for validation
	var simpleSpec SimpleSpec
	if specMap, ok := spec.(map[string]interface{}); ok {
		if title, exists := specMap["title"]; exists {
			if titleStr, ok := title.(string); ok {
				simpleSpec.Title = titleStr
			}
		}
		if version, exists := specMap["version"]; exists {
			if versionStr, ok := version.(string); ok {
				simpleSpec.Version = versionStr
			}
		}
	}
	
	// Check if spec has title (example validation)
	if simpleSpec.Title == "" {
		if strict {
			result.Valid = false
			result.Errors = append(result.Errors, ValidationError{
				Type:     "missing_title",
				Message:  "Specification title is required in strict mode",
				Severity: "error",
				Code:     "MISSING_TITLE",
				Fix:      "Add a title field to your OpenAPI specification",
			})
		} else {
			result.Warnings = append(result.Warnings, ValidationError{
				Type:     "missing_title",
				Message:  "Specification title is recommended",
				Severity: "warning",
				Code:     "MISSING_TITLE",
				Fix:      "Consider adding a title field to your OpenAPI specification",
			})
		}
	}
	
	// Check version format (example validation)
	if simpleSpec.Version == "" {
		result.Warnings = append(result.Warnings, ValidationError{
			Type:     "missing_version",
			Message:  "Specification version is recommended",
			Severity: "warning",
			Code:     "MISSING_VERSION",
			Fix:      "Add a version field to your OpenAPI specification",
		})
	}
	
	// Add info message
	result.Info = append(result.Info, ValidationError{
		Type:     "validation_complete",
		Message:  "Example validation completed successfully",
		Severity: "info",
		Code:     "VALIDATION_COMPLETE",
	})
	
	return result, nil
}

// ExampleOutputConverter demonstrates a simple output converter plugin
type ExampleOutputConverter struct {
	config map[string]interface{}
}

// GetInfo returns plugin metadata
func (p *ExampleOutputConverter) GetInfo() *PluginInfo {
	return &PluginInfo{
		ID:          "example-output-converter",
		Name:        "Example Output Converter",
		Version:     "1.0.0",
		Description: "A simple example output converter that demonstrates format conversion",
		Author:      "MCPWeaver Team",
		License:     "MIT",
		Tags:        []string{"converter", "example", "demo"},
		MinVersion:  "1.0.0",
		MaxVersion:  "2.0.0",
		Permissions: []Permission{PermissionFileSystem},
		Config: &PluginConfig{
			Schema:   json.RawMessage(`{"type": "object", "properties": {"addHeader": {"type": "boolean"}}}`),
			Default:  json.RawMessage(`{"addHeader": true}`),
			Required: []string{},
			Examples: []json.RawMessage{json.RawMessage(`{"addHeader": false}`)},
		},
	}
}

// Initialize initializes the plugin with configuration
func (p *ExampleOutputConverter) Initialize(ctx context.Context, config json.RawMessage) error {
	if config != nil {
		if err := json.Unmarshal(config, &p.config); err != nil {
			return fmt.Errorf("failed to parse config: %w", err)
		}
	} else {
		p.config = map[string]interface{}{
			"addHeader": true,
		}
	}
	return nil
}

// Shutdown cleanly shuts down the plugin
func (p *ExampleOutputConverter) Shutdown(ctx context.Context) error {
	p.config = nil
	return nil
}

// GetCapabilities returns the capabilities this plugin provides
func (p *ExampleOutputConverter) GetCapabilities() []Capability {
	return []Capability{CapabilityOutputConverter}
}

// ConvertOutput converts output from one format to another
func (p *ExampleOutputConverter) ConvertOutput(ctx context.Context, input []byte, inputFormat, outputFormat string) ([]byte, error) {
	addHeader := true
	if p.config != nil {
		if configHeader, ok := p.config["addHeader"].(bool); ok {
			addHeader = configHeader
		}
	}
	
	// Simple format conversion example
	var result []byte
	
	switch {
	case inputFormat == "json" && outputFormat == "yaml":
		// Convert JSON to YAML (simplified)
		result = []byte(fmt.Sprintf("# Converted from JSON to YAML\ndata: |\n  %s", string(input)))
		
	case inputFormat == "yaml" && outputFormat == "json":
		// Convert YAML to JSON (simplified)
		result = []byte(fmt.Sprintf(`{"data": "%s"}`, strings.ReplaceAll(string(input), "\n", "\\n")))
		
	case inputFormat == outputFormat:
		// Same format, just pass through
		result = input
		
	default:
		return nil, fmt.Errorf("unsupported conversion: %s to %s", inputFormat, outputFormat)
	}
	
	// Add header if configured
	if addHeader {
		header := fmt.Sprintf("# Generated by Example Output Converter (%s -> %s)\n", inputFormat, outputFormat)
		result = append([]byte(header), result...)
	}
	
	return result, nil
}

// CreateExamplePlugins creates instances of example plugins for testing
func CreateExamplePlugins() []Plugin {
	return []Plugin{
		&ExampleTemplateProcessor{},
		&ExampleValidator{},
		&ExampleOutputConverter{},
	}
}

// GetExamplePluginManifests returns example plugin manifests
func GetExamplePluginManifests() []*PluginManifest {
	examples := CreateExamplePlugins()
	manifests := make([]*PluginManifest, len(examples))
	
	for i, plugin := range examples {
		info := plugin.GetInfo()
		manifests[i] = &PluginManifest{
			PluginInfo: info,
			Files: []PluginFile{{
				Path:     fmt.Sprintf("%s.so", info.ID),
				Size:     1024000, // 1MB example
				Checksum: "example-checksum",
				Type:     "binary/plugin",
			}},
			Checksum:  "manifest-checksum",
			Size:      1024000,
			Verified:  true,
			Signature: "example-signature",
		}
	}
	
	return manifests
}