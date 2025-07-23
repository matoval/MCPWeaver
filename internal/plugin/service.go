package plugin

import (
	"context"
	"encoding/json"
	"fmt"
)

// Service provides plugin functionality to the main application
type Service struct {
	manager  *Manager
	registry *Registry
	security *SecurityManager
}

// NewService creates a new plugin service
func NewService(config *ManagerConfig) *Service {
	if config == nil {
		config = DefaultManagerConfig()
	}
	
	return &Service{
		manager:  NewManager(config),
		registry: NewRegistry(config),
		security: NewSecurityManager(config),
	}
}

// Initialize initializes the plugin service
func (s *Service) Initialize() error {
	if err := s.manager.Initialize(); err != nil {
		return fmt.Errorf("failed to initialize plugin manager: %w", err)
	}
	
	return nil
}

// Shutdown shuts down the plugin service
func (s *Service) Shutdown() error {
	return s.manager.Shutdown()
}

// Plugin Management API Methods

// GetPlugins returns all loaded plugins
func (s *Service) GetPlugins() (map[string]*PluginInstanceAPI, error) {
	plugins := s.manager.GetPlugins()
	apiPlugins := make(map[string]*PluginInstanceAPI)
	
	for id, instance := range plugins {
		apiPlugins[id] = ToPluginInstanceAPI(instance)
	}
	
	return apiPlugins, nil
}

// GetPlugin returns a specific plugin
func (s *Service) GetPlugin(pluginID string) (*PluginInstanceAPI, error) {
	instance, exists := s.manager.GetPlugin(pluginID)
	if !exists {
		return nil, fmt.Errorf("plugin not found: %s", pluginID)
	}
	
	return ToPluginInstanceAPI(instance), nil
}

// LoadPlugin loads a plugin from file path
func (s *Service) LoadPlugin(pluginPath string) error {
	return s.manager.LoadPlugin(pluginPath)
}

// UnloadPlugin unloads a plugin
func (s *Service) UnloadPlugin(pluginID string) error {
	return s.manager.UnloadPlugin(pluginID)
}

// EnablePlugin enables a disabled plugin
func (s *Service) EnablePlugin(pluginID string) error {
	return s.manager.EnablePlugin(pluginID)
}

// DisablePlugin disables an active plugin
func (s *Service) DisablePlugin(pluginID string) error {
	return s.manager.DisablePlugin(pluginID)
}

// ConfigurePlugin updates plugin configuration
func (s *Service) ConfigurePlugin(pluginID string, config json.RawMessage) error {
	return s.manager.ConfigurePlugin(pluginID, config)
}

// GetPluginsByCapability returns plugins with specific capability
func (s *Service) GetPluginsByCapability(capability string) ([]*PluginInstanceAPI, error) {
	cap := Capability(capability)
	instances := s.manager.GetPluginsByCapability(cap)
	
	apiInstances := make([]*PluginInstanceAPI, len(instances))
	for i, instance := range instances {
		apiInstances[i] = ToPluginInstanceAPI(instance)
	}
	
	return apiInstances, nil
}

// Marketplace API Methods

// SearchPlugins searches for plugins in marketplace
func (s *Service) SearchPlugins(ctx context.Context, query string, category string, tags []string, limit int) (*SearchResponse, error) {
	searchReq := &SearchRequest{
		Query:    query,
		Category: category,
		Tags:     tags,
		Limit:    limit,
	}
	
	return s.registry.Search(ctx, searchReq)
}

// GetMarketplacePlugin gets plugin info from marketplace
func (s *Service) GetMarketplacePlugin(ctx context.Context, pluginID string) (*MarketplacePluginAPI, error) {
	plugin, err := s.registry.GetPlugin(ctx, pluginID)
	if err != nil {
		return nil, err
	}
	
	return ToMarketplacePluginAPI(plugin), nil
}

// GetFeaturedPlugins gets featured plugins from marketplace
func (s *Service) GetFeaturedPlugins(ctx context.Context, limit int) ([]*MarketplacePluginAPI, error) {
	plugins, err := s.registry.GetFeatured(ctx, limit)
	if err != nil {
		return nil, err
	}
	
	apiPlugins := make([]*MarketplacePluginAPI, len(plugins))
	for i, plugin := range plugins {
		apiPlugins[i] = ToMarketplacePluginAPI(plugin)
	}
	
	return apiPlugins, nil
}

// GetCategories gets available plugin categories
func (s *Service) GetCategories(ctx context.Context) ([]string, error) {
	return s.registry.GetCategories(ctx)
}

// InstallPlugin installs a plugin from marketplace
func (s *Service) InstallPlugin(ctx context.Context, pluginID string) error {
	// Download plugin
	downloadPath := fmt.Sprintf("%s/%s.plugin", s.manager.config.TempDir, pluginID)
	if err := s.registry.Download(ctx, pluginID, downloadPath); err != nil {
		return fmt.Errorf("failed to download plugin: %w", err)
	}
	
	// Load plugin
	if err := s.manager.LoadPlugin(downloadPath); err != nil {
		return fmt.Errorf("failed to load plugin: %w", err)
	}
	
	return nil
}

// CheckForUpdates checks for plugin updates
func (s *Service) CheckForUpdates(ctx context.Context) (map[string]*MarketplacePluginAPI, error) {
	installedPlugins := s.manager.GetPlugins()
	updates, err := s.registry.CheckUpdates(ctx, installedPlugins)
	if err != nil {
		return nil, err
	}
	
	apiUpdates := make(map[string]*MarketplacePluginAPI)
	for pluginID, plugin := range updates {
		apiUpdates[pluginID] = ToMarketplacePluginAPI(plugin)
	}
	
	return apiUpdates, nil
}

// Plugin Execution API Methods

// ExecuteTemplateProcessor executes a template processor plugin
func (s *Service) ExecuteTemplateProcessor(ctx context.Context, pluginID string, template string, data map[string]interface{}) (string, error) {
	instance, exists := s.manager.GetPlugin(pluginID)
	if !exists {
		return "", fmt.Errorf("plugin not found: %s", pluginID)
	}
	
	processor, ok := instance.Plugin.(TemplateProcessor)
	if !ok {
		return "", fmt.Errorf("plugin is not a template processor: %s", pluginID)
	}
	
	return processor.ProcessTemplate(ctx, template, data)
}

// ExecuteValidator executes a validator plugin
func (s *Service) ExecuteValidator(ctx context.Context, pluginID string, specData map[string]interface{}) (*ValidationResultAPI, error) {
	instance, exists := s.manager.GetPlugin(pluginID)
	if !exists {
		return nil, fmt.Errorf("plugin not found: %s", pluginID)
	}
	
	validator, ok := instance.Plugin.(Validator)
	if !ok {
		return nil, fmt.Errorf("plugin is not a validator: %s", pluginID)
	}
	
	// Convert spec data to ParsedAPI (simplified)
	// In a real implementation, you'd properly convert the data
	result, err := validator.ValidateSpec(ctx, nil)
	if err != nil {
		return nil, err
	}
	
	return ToValidationResultAPI(result), nil
}

// ExecuteOutputConverter executes an output converter plugin
func (s *Service) ExecuteOutputConverter(ctx context.Context, pluginID string, input []byte, inputFormat, outputFormat string) ([]byte, error) {
	instance, exists := s.manager.GetPlugin(pluginID)
	if !exists {
		return nil, fmt.Errorf("plugin not found: %s", pluginID)
	}
	
	converter, ok := instance.Plugin.(OutputConverter)
	if !ok {
		return nil, fmt.Errorf("plugin is not an output converter: %s", pluginID)
	}
	
	return converter.ConvertOutput(ctx, input, inputFormat, outputFormat)
}

// ExecuteIntegration executes an integration plugin action
func (s *Service) ExecuteIntegration(ctx context.Context, pluginID string, action string, params map[string]interface{}) (interface{}, error) {
	instance, exists := s.manager.GetPlugin(pluginID)
	if !exists {
		return nil, fmt.Errorf("plugin not found: %s", pluginID)
	}
	
	integration, ok := instance.Plugin.(Integration)
	if !ok {
		return nil, fmt.Errorf("plugin is not an integration: %s", pluginID)
	}
	
	return integration.ExecuteAction(ctx, action, params)
}

// ExecuteTests runs tests using a testing plugin
func (s *Service) ExecuteTests(ctx context.Context, pluginID string, serverPath string, config map[string]interface{}) (*TestResultAPI, error) {
	instance, exists := s.manager.GetPlugin(pluginID)
	if !exists {
		return nil, fmt.Errorf("plugin not found: %s", pluginID)
	}
	
	tester, ok := instance.Plugin.(Testing)
	if !ok {
		return nil, fmt.Errorf("plugin is not a testing plugin: %s", pluginID)
	}
	
	result, err := tester.RunTests(ctx, serverPath, config)
	if err != nil {
		return nil, err
	}
	
	return ToTestResultAPI(result), nil
}

// Utility API Methods

// ValidatePluginConfig validates plugin configuration
func (s *Service) ValidatePluginConfig(pluginID string, config json.RawMessage) error {
	instance, exists := s.manager.GetPlugin(pluginID)
	if !exists {
		return fmt.Errorf("plugin not found: %s", pluginID)
	}
	
	// Validate configuration against plugin schema
	return s.manager.validatePluginConfig(instance, config)
}

// GetPluginCapabilities returns available plugin capabilities
func (s *Service) GetPluginCapabilities() []string {
	return []string{
		string(CapabilityTemplateProcessor),
		string(CapabilityOutputConverter),
		string(CapabilityValidator),
		string(CapabilityUIComponent),
		string(CapabilityIntegration),
		string(CapabilityTesting),
		string(CapabilityParser),
		string(CapabilityGenerator),
		string(CapabilityMiddleware),
	}
}

// GetPluginPermissions returns available plugin permissions
func (s *Service) GetPluginPermissions() []string {
	return []string{
		string(PermissionFileSystem),
		string(PermissionNetwork),
		string(PermissionDatabase),
		string(PermissionSettings),
		string(PermissionProjects),
		string(PermissionTemplates),
		string(PermissionExec),
		string(PermissionSystemInfo),
		string(PermissionClipboard),
		string(PermissionNotifications),
	}
}

// Internal conversion functions

// ToValidationResultAPI converts ValidationResult to API-safe format
func ToValidationResultAPI(result *ValidationResult) *ValidationResultAPI {
	if result == nil {
		return nil
	}
	
	apiResult := &ValidationResultAPI{
		Valid: result.Valid,
	}
	
	// Convert errors
	apiResult.Errors = make([]ValidationErrorAPI, len(result.Errors))
	for i, err := range result.Errors {
		apiResult.Errors[i] = ValidationErrorAPI{
			Type:     err.Type,
			Message:  err.Message,
			Path:     err.Path,
			Line:     err.Line,
			Column:   err.Column,
			Severity: err.Severity,
			Code:     err.Code,
			Fix:      err.Fix,
			Context:  err.Context,
		}
	}
	
	// Convert warnings
	apiResult.Warnings = make([]ValidationErrorAPI, len(result.Warnings))
	for i, warn := range result.Warnings {
		apiResult.Warnings[i] = ValidationErrorAPI{
			Type:     warn.Type,
			Message:  warn.Message,
			Path:     warn.Path,
			Line:     warn.Line,
			Column:   warn.Column,
			Severity: warn.Severity,
			Code:     warn.Code,
			Fix:      warn.Fix,
			Context:  warn.Context,
		}
	}
	
	// Convert info
	apiResult.Info = make([]ValidationErrorAPI, len(result.Info))
	for i, info := range result.Info {
		apiResult.Info[i] = ValidationErrorAPI{
			Type:     info.Type,
			Message:  info.Message,
			Path:     info.Path,
			Line:     info.Line,
			Column:   info.Column,
			Severity: info.Severity,
			Code:     info.Code,
			Fix:      info.Fix,
			Context:  info.Context,
		}
	}
	
	// Convert stats
	if result.Stats != nil {
		apiResult.Stats = &ValidationStatsAPI{
			TotalChecks:  result.Stats.TotalChecks,
			Duration:     int64(result.Stats.Duration),
			RulesApplied: result.Stats.RulesApplied,
			FilesChecked: result.Stats.FilesChecked,
			LinesChecked: result.Stats.LinesChecked,
		}
	}
	
	return apiResult
}

// ToTestResultAPI converts TestResult to API-safe format
func ToTestResultAPI(result *TestResult) *TestResultAPI {
	if result == nil {
		return nil
	}
	
	apiResult := &TestResultAPI{
		Passed:   result.Passed,
		Duration: int64(result.Duration),
		Summary:  result.Summary,
	}
	
	// Convert test cases
	apiResult.Tests = make([]TestCaseAPI, len(result.Tests))
	for i, test := range result.Tests {
		apiResult.Tests[i] = TestCaseAPI{
			Name:     test.Name,
			Status:   test.Status,
			Duration: int64(test.Duration),
			Message:  test.Message,
			Error:    test.Error,
		}
	}
	
	// Convert coverage
	if result.Coverage != nil {
		apiResult.Coverage = &CoverageAPI{
			Lines:     result.Coverage.Lines,
			Covered:   result.Coverage.Covered,
			Percent:   result.Coverage.Percent,
			Files:     result.Coverage.Files,
			Functions: result.Coverage.Functions,
		}
	}
	
	// Convert performance
	if result.Performance != nil {
		apiResult.Performance = &PerformanceAPI{
			RequestsPerSecond: result.Performance.RequestsPerSecond,
			AverageLatency:    int64(result.Performance.AverageLatency),
			MaxLatency:        int64(result.Performance.MaxLatency),
			MinLatency:        int64(result.Performance.MinLatency),
			MemoryUsage:       result.Performance.MemoryUsage,
			CPUUsage:          result.Performance.CPUUsage,
		}
	}
	
	return apiResult
}