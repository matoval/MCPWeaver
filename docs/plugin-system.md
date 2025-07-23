# MCPWeaver Plugin System Documentation

## Overview

The MCPWeaver Plugin System provides a comprehensive architecture for extending MCPWeaver's functionality through custom plugins. This system supports various types of plugins including template processors, validators, output converters, and more.

## Architecture

### Core Components

1. **Plugin Manager**: Central orchestrator for plugin lifecycle management
2. **Plugin Loader**: Handles loading and unloading of plugin binaries
3. **Security Manager**: Manages plugin permissions and sandboxing
4. **Registry**: Handles marketplace integration and plugin discovery
5. **Event Bus**: Enables communication between plugins and the main application

### Plugin Types

#### 1. Template Processor Plugins
Process and transform template files with custom logic.

```go
type TemplateProcessor interface {
    ProcessTemplate(ctx context.Context, template string, data map[string]interface{}) (string, error)
}
```

#### 2. Validator Plugins
Validate OpenAPI specifications with custom rules.

```go
type Validator interface {
    ValidateSpec(ctx context.Context, spec interface{}) (*ValidationResult, error)
}
```

#### 3. Output Converter Plugins
Convert generated output between different formats.

```go
type OutputConverter interface {
    ConvertOutput(ctx context.Context, input []byte, inputFormat, outputFormat string) ([]byte, error)
}
```

#### 4. Integration Plugins
Integrate with external services and tools.

```go
type Integration interface {
    ExecuteAction(ctx context.Context, action string, params map[string]interface{}) (interface{}, error)
}
```

#### 5. Testing Plugins
Provide testing capabilities for generated MCP servers.

```go
type Testing interface {
    RunTests(ctx context.Context, serverPath string, config map[string]interface{}) (*TestResult, error)
}
```

## Plugin Development

### Creating a Plugin

1. **Implement the Plugin Interface**:

```go
package main

import (
    "context"
    "encoding/json"
)

type MyPlugin struct {
    config map[string]interface{}
}

func (p *MyPlugin) GetInfo() *PluginInfo {
    return &PluginInfo{
        ID:          "my-plugin",
        Name:        "My Custom Plugin",
        Version:     "1.0.0",
        Description: "A custom plugin for MCPWeaver",
        Author:      "Your Name",
        License:     "MIT",
        Permissions: []Permission{PermissionFileSystem},
    }
}

func (p *MyPlugin) Initialize(ctx context.Context, config json.RawMessage) error {
    if config != nil {
        return json.Unmarshal(config, &p.config)
    }
    return nil
}

func (p *MyPlugin) Shutdown(ctx context.Context) error {
    return nil
}

func (p *MyPlugin) GetCapabilities() []Capability {
    return []Capability{CapabilityTemplateProcessor}
}

func (p *MyPlugin) ProcessTemplate(ctx context.Context, template string, data map[string]interface{}) (string, error) {
    // Your template processing logic here
    return template, nil
}

// Required export function
func NewPlugin() Plugin {
    return &MyPlugin{}
}
```

2. **Create Plugin Manifest**:

```json
{
    "id": "my-plugin",
    "name": "My Custom Plugin",
    "version": "1.0.0",
    "description": "A custom plugin for MCPWeaver",
    "author": "Your Name",
    "license": "MIT",
    "capabilities": ["template_processor"],
    "permissions": ["file_system"],
    "config": {
        "schema": {
            "type": "object",
            "properties": {
                "option1": {"type": "string"},
                "option2": {"type": "boolean"}
            }
        },
        "default": {
            "option1": "default_value",
            "option2": true
        }
    }
}
```

3. **Build the Plugin**:

```bash
go build -buildmode=plugin -o my-plugin.so my-plugin.go
```

### Plugin Configuration

Plugins can be configured through JSON configuration that matches their schema:

```json
{
    "option1": "custom_value",
    "option2": false,
    "advanced_settings": {
        "timeout": 30,
        "retries": 3
    }
}
```

## Security and Sandboxing

### Permission System

Plugins must declare required permissions in their manifest:

- `file_system`: Access to file system operations
- `network`: Network access for HTTP requests
- `database`: Access to MCPWeaver database
- `settings`: Access to application settings
- `projects`: Access to project data
- `templates`: Access to template system
- `exec`: Execute external commands
- `system_info`: Access to system information
- `clipboard`: Access to system clipboard
- `notifications`: Show notifications

### Sandboxing

The security manager provides sandboxing capabilities:

```go
sandbox, err := securityManager.CreateSandbox(pluginID)
if err != nil {
    return err
}

// Check file access
err = sandbox.CheckFileAccess("/path/to/file")
if err != nil {
    return fmt.Errorf("file access denied: %w", err)
}

// Check network access  
err = sandbox.CheckNetworkAccess("example.com")
if err != nil {
    return fmt.Errorf("network access denied: %w", err)
}
```

## Plugin Marketplace

### Publishing Plugins

1. Create plugin package with manifest and binary
2. Sign plugin with private key
3. Upload to marketplace with metadata
4. Plugin goes through review process
5. Published plugins are available for installation

### Installing Plugins

```go
// Search for plugins
results, err := pluginService.SearchPlugins(ctx, "template", "generator", []string{"go"}, 10)

// Install plugin
err = pluginService.InstallPlugin(ctx, "plugin-id")

// Load installed plugin
err = pluginService.LoadPlugin(ctx, "/path/to/plugin.so")
```

## API Reference

### Plugin Management

```go
// Get all loaded plugins
plugins, err := app.GetPlugins(ctx)

// Get specific plugin
plugin, err := app.GetPlugin(ctx, "plugin-id")

// Load plugin from file
err = app.LoadPlugin(ctx, "/path/to/plugin.so")

// Unload plugin
err = app.UnloadPlugin(ctx, "plugin-id")

// Enable/disable plugin
err = app.EnablePlugin(ctx, "plugin-id")
err = app.DisablePlugin(ctx, "plugin-id")

// Get plugins by capability
processors, err := app.GetPluginsByCapability(ctx, "template_processor")
```

### Marketplace API

```go
// Search marketplace
results, err := app.SearchPlugins(ctx, "query", "category", []string{"tag1", "tag2"}, 20)

// Install from marketplace
err = app.InstallPlugin(ctx, "plugin-id")

// Get plugin capabilities and permissions
capabilities := app.GetPluginCapabilities(ctx)
permissions := app.GetPluginPermissions(ctx)
```

## Testing Framework

### Plugin Testing

```go
// Create test framework
testFramework := plugin.NewTestFramework(manager, nil)

// Test plugin
result, err := testFramework.TestPlugin(ctx, "/path/to/plugin.so")

// Check results
if result.Passed {
    fmt.Printf("All tests passed: %s\n", result.Summary)
} else {
    fmt.Printf("Tests failed: %s\n", result.Summary)
    for _, test := range result.Tests {
        if test.Status != "passed" {
            fmt.Printf("- %s: %s\n", test.Name, test.Error)
        }
    }
}
```

### Test Categories

1. **Plugin Info Validation**: Validates plugin metadata
2. **Initialization Test**: Tests plugin initialization
3. **Capability Tests**: Tests each declared capability
4. **Security Tests**: Validates security compliance
5. **Performance Tests**: Benchmarks plugin performance

## Examples

### Example Template Processor

See `internal/plugin/examples.go` for complete examples including:

- Template processor with custom syntax
- Validator with configurable rules  
- Output converter between formats
- Integration plugin for external services

### Example Usage

```go
// Load example plugins
examples := plugin.CreateExamplePlugins()

for _, example := range examples {
    err := manager.LoadPluginFromInstance("example-"+example.GetInfo().ID, example)
    if err != nil {
        log.Printf("Failed to load example plugin: %v", err)
        continue
    }
    
    // Use the plugin
    if processor, ok := example.(plugin.TemplateProcessor); ok {
        result, err := processor.ProcessTemplate(ctx, "Hello {{name}}!", map[string]interface{}{
            "name": "World",
        })
        fmt.Printf("Result: %s\n", result)
    }
}
```

## Best Practices

### Development

1. **Error Handling**: Always return meaningful errors with context
2. **Resource Cleanup**: Implement proper shutdown logic
3. **Configuration**: Use JSON schema for configuration validation
4. **Logging**: Use structured logging for debugging
5. **Testing**: Write comprehensive tests for all functionality

### Security

1. **Minimal Permissions**: Request only necessary permissions
2. **Input Validation**: Validate all external input
3. **Safe Defaults**: Use secure default configurations
4. **Error Messages**: Don't leak sensitive information in errors
5. **Resource Limits**: Respect memory and CPU limits

### Performance

1. **Lazy Loading**: Load resources only when needed
2. **Caching**: Cache expensive operations appropriately
3. **Async Operations**: Use goroutines for concurrent work
4. **Memory Management**: Clean up resources properly
5. **Metrics**: Provide performance metrics

## Troubleshooting

### Common Issues

1. **Plugin Won't Load**: Check file permissions and binary format
2. **Permission Denied**: Verify required permissions in manifest
3. **Configuration Error**: Validate configuration against schema
4. **Capability Not Found**: Ensure plugin implements required interfaces
5. **Security Violation**: Check sandbox restrictions

### Debug Mode

Enable debug logging to troubleshoot plugin issues:

```go
config := &plugin.ManagerConfig{
    LogLevel: "debug", 
    SecurityPolicy: "permissive",
}
```

### Plugin Validation

Use the testing framework to validate plugins:

```bash
go run cmd/plugin-test/main.go -plugin /path/to/plugin.so
```

## Future Roadmap

1. **Hot Reloading**: Update plugins without restart
2. **Plugin Dependencies**: Automatic dependency resolution
3. **Visual Plugin Builder**: GUI for creating simple plugins
4. **Plugin Analytics**: Usage and performance monitoring
5. **Remote Plugins**: Load plugins from remote URLs
6. **Plugin Clustering**: Distribute plugin execution