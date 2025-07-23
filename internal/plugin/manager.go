package plugin

import (
	"context"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"
	"time"
)

// Manager handles plugin lifecycle and orchestration
type Manager struct {
	plugins     map[string]*PluginInstance
	middleware  []MiddlewarePlugin
	validators  []Validator
	converters  []OutputConverter
	processors  []TemplateProcessor
	integrations []Integration
	testers     []Testing
	uiComponents []UIComponent
	parsers     []ParserPlugin
	generators  []GeneratorPlugin
	listeners   []EventListener
	
	config      *ManagerConfig
	loader      *Loader
	registry    *Registry
	security    *SecurityManager
	eventBus    *EventBus
	
	mu          sync.RWMutex
	ctx         context.Context
	cancel      context.CancelFunc
}

// ManagerConfig configures the plugin manager
type ManagerConfig struct {
	PluginDir         string        `json:"pluginDir"`
	MaxPlugins        int           `json:"maxPlugins"`
	LoadTimeout       time.Duration `json:"loadTimeout"`
	EnableSandbox     bool          `json:"enableSandbox"`
	EnableMarketplace bool          `json:"enableMarketplace"`
	MarketplaceURL    string        `json:"marketplaceUrl"`
	SecurityPolicy    string        `json:"securityPolicy"`
	LogLevel          string        `json:"logLevel"`
	UpdateCheck       time.Duration `json:"updateCheck"`
	CacheDir          string        `json:"cacheDir"`
	TempDir           string        `json:"tempDir"`
	AllowedHosts      []string      `json:"allowedHosts"`
}

// DefaultManagerConfig returns default configuration
func DefaultManagerConfig() *ManagerConfig {
	return &ManagerConfig{
		PluginDir:         "./plugins",
		MaxPlugins:        50,
		LoadTimeout:       30 * time.Second,
		EnableSandbox:     true,
		EnableMarketplace: true,
		MarketplaceURL:    "https://plugins.mcpweaver.com",
		SecurityPolicy:    "strict",
		LogLevel:          "info",
		UpdateCheck:       24 * time.Hour,
		CacheDir:          "./cache/plugins",
		TempDir:           "./temp/plugins",
		AllowedHosts:      []string{"localhost", "127.0.0.1"},
	}
}

// NewManager creates a new plugin manager
func NewManager(config *ManagerConfig) *Manager {
	if config == nil {
		config = DefaultManagerConfig()
	}
	
	ctx, cancel := context.WithCancel(context.Background())
	
	return &Manager{
		plugins:      make(map[string]*PluginInstance),
		middleware:   make([]MiddlewarePlugin, 0),
		validators:   make([]Validator, 0),
		converters:   make([]OutputConverter, 0),
		processors:   make([]TemplateProcessor, 0),
		integrations: make([]Integration, 0),
		testers:      make([]Testing, 0),
		uiComponents: make([]UIComponent, 0),
		parsers:      make([]ParserPlugin, 0),
		generators:   make([]GeneratorPlugin, 0),
		listeners:    make([]EventListener, 0),
		
		config:   config,
		loader:   NewLoader(config),
		registry: NewRegistry(config),
		security: NewSecurityManager(config),
		eventBus: NewEventBus(),
		
		ctx:    ctx,
		cancel: cancel,
	}
}

// Initialize starts the plugin manager
func (m *Manager) Initialize() error {
	m.mu.Lock()
	defer m.mu.Unlock()
	
	// Create required directories
	if err := m.createDirectories(); err != nil {
		return fmt.Errorf("failed to create directories: %w", err)
	}
	
	// Initialize security manager
	if err := m.security.Initialize(); err != nil {
		return fmt.Errorf("failed to initialize security: %w", err)
	}
	
	// Initialize event bus
	if err := m.eventBus.Initialize(); err != nil {
		return fmt.Errorf("failed to initialize event bus: %w", err)
	}
	
	// Load plugins from plugin directory
	if err := m.loadInstalledPlugins(); err != nil {
		return fmt.Errorf("failed to load plugins: %w", err)
	}
	
	// Start update checker if enabled
	if m.config.EnableMarketplace && m.config.UpdateCheck > 0 {
		go m.updateChecker()
	}
	
	return nil
}

// Shutdown stops the plugin manager and all plugins
func (m *Manager) Shutdown() error {
	m.cancel()
	
	m.mu.Lock()
	defer m.mu.Unlock()
	
	var errors []error
	
	// Shutdown all plugins
	for id := range m.plugins {
		if err := m.unloadPluginUnsafe(id); err != nil {
			errors = append(errors, fmt.Errorf("failed to unload plugin %s: %w", id, err))
		}
	}
	
	// Shutdown subsystems
	if err := m.eventBus.Shutdown(); err != nil {
		errors = append(errors, fmt.Errorf("failed to shutdown event bus: %w", err))
	}
	
	if len(errors) > 0 {
		return fmt.Errorf("shutdown errors: %v", errors)
	}
	
	return nil
}

// LoadPlugin loads a plugin from a file path
func (m *Manager) LoadPlugin(pluginPath string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	
	return m.loadPluginUnsafe(pluginPath)
}

// UnloadPlugin unloads a plugin by ID
func (m *Manager) UnloadPlugin(pluginID string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	
	return m.unloadPluginUnsafe(pluginID)
}

// GetPlugin returns a plugin instance by ID
func (m *Manager) GetPlugin(pluginID string) (*PluginInstance, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	
	instance, exists := m.plugins[pluginID]
	return instance, exists
}

// GetPlugins returns all loaded plugins
func (m *Manager) GetPlugins() map[string]*PluginInstance {
	m.mu.RLock()
	defer m.mu.RUnlock()
	
	plugins := make(map[string]*PluginInstance)
	for id, instance := range m.plugins {
		plugins[id] = instance
	}
	
	return plugins
}

// GetPluginsByCapability returns plugins with specific capability
func (m *Manager) GetPluginsByCapability(capability Capability) []*PluginInstance {
	m.mu.RLock()
	defer m.mu.RUnlock()
	
	var result []*PluginInstance
	for _, instance := range m.plugins {
		if instance.Status == PluginStatusActive {
			capabilities := instance.Plugin.GetCapabilities()
			for _, cap := range capabilities {
				if cap == capability {
					result = append(result, instance)
					break
				}
			}
		}
	}
	
	return result
}

// EnablePlugin enables a disabled plugin
func (m *Manager) EnablePlugin(pluginID string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	
	instance, exists := m.plugins[pluginID]
	if !exists {
		return fmt.Errorf("plugin not found: %s", pluginID)
	}
	
	if instance.Status != PluginStatusDisabled {
		return fmt.Errorf("plugin is not disabled: %s", pluginID)
	}
	
	return m.activatePluginUnsafe(instance)
}

// DisablePlugin disables an active plugin
func (m *Manager) DisablePlugin(pluginID string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	
	instance, exists := m.plugins[pluginID]
	if !exists {
		return fmt.Errorf("plugin not found: %s", pluginID)
	}
	
	if instance.Status != PluginStatusActive {
		return fmt.Errorf("plugin is not active: %s", pluginID)
	}
	
	return m.deactivatePluginUnsafe(instance)
}

// ConfigurePlugin updates plugin configuration
func (m *Manager) ConfigurePlugin(pluginID string, config json.RawMessage) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	
	instance, exists := m.plugins[pluginID]
	if !exists {
		return fmt.Errorf("plugin not found: %s", pluginID)
	}
	
	// Validate configuration
	if err := m.validatePluginConfig(instance, config); err != nil {
		return fmt.Errorf("invalid configuration: %w", err)
	}
	
	// Store old config for rollback
	oldConfig := instance.Config
	instance.Config = config
	
	// Re-initialize plugin with new config
	if err := instance.Plugin.Initialize(m.ctx, config); err != nil {
		// Rollback on error
		instance.Config = oldConfig
		return fmt.Errorf("failed to apply configuration: %w", err)
	}
	
	return nil
}

// Internal methods

func (m *Manager) createDirectories() error {
	dirs := []string{
		m.config.PluginDir,
		m.config.CacheDir,
		m.config.TempDir,
	}
	
	for _, dir := range dirs {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("failed to create directory %s: %w", dir, err)
		}
	}
	
	return nil
}

func (m *Manager) loadInstalledPlugins() error {
	return filepath.WalkDir(m.config.PluginDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		
		// Look for plugin manifests
		if d.IsDir() || !strings.HasSuffix(path, ".json") {
			return nil
		}
		
		// Check if this looks like a manifest
		if strings.HasSuffix(path, ".manifest.json") {
			if err := m.loadPluginFromManifest(path); err != nil {
				// Log error but continue loading other plugins
				fmt.Printf("Failed to load plugin from %s: %v\n", path, err)
			}
		}
		
		return nil
	})
}

func (m *Manager) loadPluginFromManifest(manifestPath string) error {
	data, err := os.ReadFile(manifestPath)
	if err != nil {
		return fmt.Errorf("failed to read manifest: %w", err)
	}
	
	var manifest PluginManifest
	if err := json.Unmarshal(data, &manifest); err != nil {
		return fmt.Errorf("failed to parse manifest: %w", err)
	}
	
	// Security validation
	if err := m.security.ValidateManifest(&manifest); err != nil {
		return fmt.Errorf("security validation failed: %w", err)
	}
	
	// Load the plugin
	pluginDir := filepath.Dir(manifestPath)
	return m.loadPluginFromDirectory(pluginDir, &manifest)
}

func (m *Manager) loadPluginFromDirectory(pluginDir string, manifest *PluginManifest) error {
	// Find the main plugin file
	var pluginFile string
	for _, file := range manifest.Files {
		if file.Type == "binary" && strings.HasSuffix(file.Path, getPluginExtension()) {
			pluginFile = filepath.Join(pluginDir, file.Path)
			break
		}
	}
	
	if pluginFile == "" {
		return fmt.Errorf("no plugin binary found")
	}
	
	return m.loadPluginUnsafe(pluginFile)
}

func (m *Manager) loadPluginUnsafe(pluginPath string) error {
	// Load plugin using loader
	plugin, err := m.loader.Load(pluginPath)
	if err != nil {
		return fmt.Errorf("failed to load plugin: %w", err)
	}
	
	info := plugin.GetInfo()
	
	// Check if already loaded
	if _, exists := m.plugins[info.ID]; exists {
		return fmt.Errorf("plugin already loaded: %s", info.ID)
	}
	
	// Check plugin limits
	if len(m.plugins) >= m.config.MaxPlugins {
		return fmt.Errorf("maximum plugin limit reached: %d", m.config.MaxPlugins)
	}
	
	// Validate plugin
	if err := m.validatePlugin(plugin); err != nil {
		return fmt.Errorf("plugin validation failed: %w", err)
	}
	
	// Create instance
	instance := &PluginInstance{
		Plugin:   plugin,
		Info:     info,
		Status:   PluginStatusLoaded,
		LoadedAt: time.Now(),
		Stats:    &PluginStats{},
	}
	
	// Store instance
	m.plugins[info.ID] = instance
	
	// Activate plugin
	if err := m.activatePluginUnsafe(instance); err != nil {
		delete(m.plugins, info.ID)
		return fmt.Errorf("failed to activate plugin: %w", err)
	}
	
	return nil
}

func (m *Manager) unloadPluginUnsafe(pluginID string) error {
	instance, exists := m.plugins[pluginID]
	if !exists {
		return fmt.Errorf("plugin not found: %s", pluginID)
	}
	
	instance.Status = PluginStatusUnloading
	
	// Deactivate if active
	if instance.Status == PluginStatusActive {
		if err := m.deactivatePluginUnsafe(instance); err != nil {
			return err
		}
	}
	
	// Shutdown plugin
	if err := instance.Plugin.Shutdown(m.ctx); err != nil {
		return fmt.Errorf("failed to shutdown plugin: %w", err)
	}
	
	// Remove from collections
	m.removePluginFromCollections(instance)
	
	// Remove from registry
	delete(m.plugins, pluginID)
	
	return nil
}

func (m *Manager) activatePluginUnsafe(instance *PluginInstance) error {
	// Initialize plugin
	if err := instance.Plugin.Initialize(m.ctx, instance.Config); err != nil {
		instance.Status = PluginStatusError
		instance.LastError = err.Error()
		return fmt.Errorf("failed to initialize plugin: %w", err)
	}
	
	// Add to appropriate collections based on capabilities
	m.addPluginToCollections(instance)
	
	instance.Status = PluginStatusActive
	instance.LastError = ""
	
	// Emit plugin loaded event
	event := &PluginEvent{
		Type:      "plugin.loaded",
		Source:    "plugin_manager",
		Target:    "*",
		Timestamp: time.Now(),
		Data: map[string]interface{}{
			"pluginId": instance.Info.ID,
			"version":  instance.Info.Version,
		},
	}
	
	m.eventBus.Emit(m.ctx, event)
	
	return nil
}

func (m *Manager) deactivatePluginUnsafe(instance *PluginInstance) error {
	// Remove from collections
	m.removePluginFromCollections(instance)
	
	instance.Status = PluginStatusDisabled
	
	// Emit plugin unloaded event
	event := &PluginEvent{
		Type:      "plugin.unloaded",
		Source:    "plugin_manager",
		Target:    "*",
		Timestamp: time.Now(),
		Data: map[string]interface{}{
			"pluginId": instance.Info.ID,
		},
	}
	
	m.eventBus.Emit(m.ctx, event)
	
	return nil
}

func (m *Manager) addPluginToCollections(instance *PluginInstance) {
	capabilities := instance.Plugin.GetCapabilities()
	
	for _, capability := range capabilities {
		switch capability {
		case CapabilityMiddleware:
			if middleware, ok := instance.Plugin.(MiddlewarePlugin); ok {
				m.middleware = append(m.middleware, middleware)
				m.sortMiddleware()
			}
		case CapabilityValidator:
			if validator, ok := instance.Plugin.(Validator); ok {
				m.validators = append(m.validators, validator)
			}
		case CapabilityOutputConverter:
			if converter, ok := instance.Plugin.(OutputConverter); ok {
				m.converters = append(m.converters, converter)
			}
		case CapabilityTemplateProcessor:
			if processor, ok := instance.Plugin.(TemplateProcessor); ok {
				m.processors = append(m.processors, processor)
			}
		case CapabilityIntegration:
			if integration, ok := instance.Plugin.(Integration); ok {
				m.integrations = append(m.integrations, integration)
			}
		case CapabilityTesting:
			if tester, ok := instance.Plugin.(Testing); ok {
				m.testers = append(m.testers, tester)
			}
		case CapabilityUIComponent:
			if component, ok := instance.Plugin.(UIComponent); ok {
				m.uiComponents = append(m.uiComponents, component)
			}
		case CapabilityParser:
			if parser, ok := instance.Plugin.(ParserPlugin); ok {
				m.parsers = append(m.parsers, parser)
			}
		case CapabilityGenerator:
			if generator, ok := instance.Plugin.(GeneratorPlugin); ok {
				m.generators = append(m.generators, generator)
			}
		}
	}
	
	// Add event listeners
	if listener, ok := instance.Plugin.(EventListener); ok {
		m.listeners = append(m.listeners, listener)
		
		// Subscribe to events
		subscriptions := listener.GetSubscriptions()
		for _, eventType := range subscriptions {
			m.eventBus.Subscribe(eventType, func(ctx context.Context, event *PluginEvent) error {
				return listener.HandleEvent(ctx, event)
			})
		}
	}
}

func (m *Manager) removePluginFromCollections(instance *PluginInstance) {
	capabilities := instance.Plugin.GetCapabilities()
	
	for _, capability := range capabilities {
		switch capability {
		case CapabilityMiddleware:
			m.middleware = removeMiddleware(m.middleware, instance.Plugin)
		case CapabilityValidator:
			m.validators = removeValidator(m.validators, instance.Plugin)
		case CapabilityOutputConverter:
			m.converters = removeConverter(m.converters, instance.Plugin)
		case CapabilityTemplateProcessor:
			m.processors = removeProcessor(m.processors, instance.Plugin)
		case CapabilityIntegration:
			m.integrations = removeIntegration(m.integrations, instance.Plugin)
		case CapabilityTesting:
			m.testers = removeTester(m.testers, instance.Plugin)
		case CapabilityUIComponent:
			m.uiComponents = removeUIComponent(m.uiComponents, instance.Plugin)
		case CapabilityParser:
			m.parsers = removeParserPlugin(m.parsers, instance.Plugin)
		case CapabilityGenerator:
			m.generators = removeGeneratorPlugin(m.generators, instance.Plugin)
		}
	}
	
	// Remove event listeners
	if listener, ok := instance.Plugin.(EventListener); ok {
		m.listeners = removeEventListener(m.listeners, instance.Plugin)
		
		// Unsubscribe from events
		subscriptions := listener.GetSubscriptions()
		for _, eventType := range subscriptions {
			m.eventBus.Unsubscribe(eventType, listener)
		}
	}
}

func (m *Manager) sortMiddleware() {
	sort.Slice(m.middleware, func(i, j int) bool {
		return m.middleware[i].GetPriority() < m.middleware[j].GetPriority()
	})
}

func (m *Manager) validatePlugin(plugin Plugin) error {
	info := plugin.GetInfo()
	
	// Check required fields
	if info.ID == "" || info.Name == "" || info.Version == "" {
		return fmt.Errorf("missing required plugin info")
	}
	
	// Check permissions
	if err := m.security.ValidatePermissions(info.Permissions); err != nil {
		return fmt.Errorf("invalid permissions: %w", err)
	}
	
	// Check capabilities
	capabilities := plugin.GetCapabilities()
	if len(capabilities) == 0 {
		return fmt.Errorf("plugin must have at least one capability")
	}
	
	return nil
}

func (m *Manager) validatePluginConfig(instance *PluginInstance, config json.RawMessage) error {
	if instance.Info.Config == nil {
		return nil // No config schema to validate against
	}
	
	// TODO: Implement JSON Schema validation
	return nil
}

func (m *Manager) updateChecker() {
	ticker := time.NewTicker(m.config.UpdateCheck)
	defer ticker.Stop()
	
	for {
		select {
		case <-m.ctx.Done():
			return
		case <-ticker.C:
			m.checkForUpdates()
		}
	}
}

func (m *Manager) checkForUpdates() {
	// TODO: Implement update checking against marketplace
}

// Utility functions for removing plugins from slices
func removeMiddleware(slice []MiddlewarePlugin, plugin Plugin) []MiddlewarePlugin {
	for i, item := range slice {
		if item == plugin {
			return append(slice[:i], slice[i+1:]...)
		}
	}
	return slice
}

func removeValidator(slice []Validator, plugin Plugin) []Validator {
	for i, item := range slice {
		if item == plugin {
			return append(slice[:i], slice[i+1:]...)
		}
	}
	return slice
}

func removeConverter(slice []OutputConverter, plugin Plugin) []OutputConverter {
	for i, item := range slice {
		if item == plugin {
			return append(slice[:i], slice[i+1:]...)
		}
	}
	return slice
}

func removeProcessor(slice []TemplateProcessor, plugin Plugin) []TemplateProcessor {
	for i, item := range slice {
		if item == plugin {
			return append(slice[:i], slice[i+1:]...)
		}
	}
	return slice
}

func removeIntegration(slice []Integration, plugin Plugin) []Integration {
	for i, item := range slice {
		if item == plugin {
			return append(slice[:i], slice[i+1:]...)
		}
	}
	return slice
}

func removeTester(slice []Testing, plugin Plugin) []Testing {
	for i, item := range slice {
		if item == plugin {
			return append(slice[:i], slice[i+1:]...)
		}
	}
	return slice
}

func removeUIComponent(slice []UIComponent, plugin Plugin) []UIComponent {
	for i, item := range slice {
		if item == plugin {
			return append(slice[:i], slice[i+1:]...)
		}
	}
	return slice
}

func removeParserPlugin(slice []ParserPlugin, plugin Plugin) []ParserPlugin {
	for i, item := range slice {
		if item == plugin {
			return append(slice[:i], slice[i+1:]...)
		}
	}
	return slice
}

func removeGeneratorPlugin(slice []GeneratorPlugin, plugin Plugin) []GeneratorPlugin {
	for i, item := range slice {
		if item == plugin {
			return append(slice[:i], slice[i+1:]...)
		}
	}
	return slice
}

func removeEventListener(slice []EventListener, plugin Plugin) []EventListener {
	for i, item := range slice {
		if item == plugin {
			return append(slice[:i], slice[i+1:]...)
		}
	}
	return slice
}

func getPluginExtension() string {
	// Platform-specific plugin extensions
	switch {
	case strings.Contains(strings.ToLower(os.Getenv("OS")), "windows"):
		return ".dll"
	case strings.Contains(strings.ToLower(os.Getenv("OS")), "darwin"):
		return ".dylib"
	default:
		return ".so"
	}
}

func calculateChecksum(data []byte) string {
	hash := sha256.Sum256(data)
	return fmt.Sprintf("%x", hash)
}