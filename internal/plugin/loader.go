package plugin

import (
	"context"
	"encoding/json"
	"fmt"
	"plugin"
	"time"
)

// Loader handles loading plugins from files
type Loader struct {
	config  *ManagerConfig
	timeout time.Duration
}

// NewLoader creates a new plugin loader
func NewLoader(config *ManagerConfig) *Loader {
	return &Loader{
		config:  config,
		timeout: config.LoadTimeout,
	}
}

// Load loads a plugin from a file path
func (l *Loader) Load(pluginPath string) (Plugin, error) {
	// Load the plugin using Go's plugin system
	p, err := plugin.Open(pluginPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open plugin: %w", err)
	}
	
	// Look for the NewPlugin function
	newPluginFunc, err := p.Lookup("NewPlugin")
	if err != nil {
		return nil, fmt.Errorf("plugin missing NewPlugin function: %w", err)
	}
	
	// Cast to proper function type
	createPlugin, ok := newPluginFunc.(func() Plugin)
	if !ok {
		return nil, fmt.Errorf("NewPlugin function has wrong signature")
	}
	
	// Create plugin instance
	pluginInstance := createPlugin()
	if pluginInstance == nil {
		return nil, fmt.Errorf("plugin creation returned nil")
	}
	
	return pluginInstance, nil
}

// LoadWithWrapper loads a plugin with a wrapper for sandboxing
func (l *Loader) LoadWithWrapper(pluginPath string, sandbox bool) (Plugin, error) {
	if !sandbox {
		return l.Load(pluginPath)
	}
	
	// Load the raw plugin
	rawPlugin, err := l.Load(pluginPath)
	if err != nil {
		return nil, err
	}
	
	// Wrap in sandbox
	return NewSandboxWrapper(rawPlugin, l.config), nil
}

// SandboxWrapper wraps a plugin to provide sandboxing
type SandboxWrapper struct {
	plugin Plugin
	config *ManagerConfig
}

// NewSandboxWrapper creates a new sandbox wrapper
func NewSandboxWrapper(plugin Plugin, config *ManagerConfig) *SandboxWrapper {
	return &SandboxWrapper{
		plugin: plugin,
		config: config,
	}
}

// GetInfo returns plugin metadata
func (w *SandboxWrapper) GetInfo() *PluginInfo {
	return w.plugin.GetInfo()
}

// Initialize initializes the plugin with configuration
func (w *SandboxWrapper) Initialize(ctx context.Context, config json.RawMessage) error {
	// Add sandbox context
	sandboxCtx := w.createSandboxContext(ctx)
	return w.plugin.Initialize(sandboxCtx, config)
}

// Shutdown cleanly shuts down the plugin
func (w *SandboxWrapper) Shutdown(ctx context.Context) error {
	sandboxCtx := w.createSandboxContext(ctx)
	return w.plugin.Shutdown(sandboxCtx)
}

// GetCapabilities returns what the plugin can do
func (w *SandboxWrapper) GetCapabilities() []Capability {
	return w.plugin.GetCapabilities()
}

func (w *SandboxWrapper) createSandboxContext(ctx context.Context) context.Context {
	// Add sandbox restrictions to context
	sandboxCtx := context.WithValue(ctx, "sandbox", true)
	sandboxCtx = context.WithValue(sandboxCtx, "allowedHosts", w.config.AllowedHosts)
	sandboxCtx = context.WithValue(sandboxCtx, "tempDir", w.config.TempDir)
	
	return sandboxCtx
}