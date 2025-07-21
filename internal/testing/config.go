package testing

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// ConfigManager handles test configuration management
type ConfigManager struct {
	configPath     string
	defaultConfig  *TestConfig
	profiles       map[string]*TestConfig
	currentProfile string
}

// NewConfigManager creates a new configuration manager
func NewConfigManager(configPath string) *ConfigManager {
	return &ConfigManager{
		configPath:     configPath,
		defaultConfig:  getDefaultTestConfig(),
		profiles:       make(map[string]*TestConfig),
		currentProfile: "default",
	}
}

// getDefaultTestConfig returns the default test configuration
func getDefaultTestConfig() *TestConfig {
	return &TestConfig{
		Timeout:               5 * time.Minute,
		MaxConcurrentTests:    3,
		EnableParallelTesting: true,
		ContinueOnFailure:     false,
		EnableSecurityScanning: true,
		EnableLinting:         true,
		EnablePerformanceTesting: true,
		EnableIntegrationTesting: true,
		MCPProtocolVersion:    "2024-11-05",
		RequiredMethods:       []string{"initialize", "tools/list", "tools/call"},
		RequiredCapabilities:  []string{"tools"},
		MaxResponseTime:       time.Second,
		MaxMemoryUsage:        100 * 1024 * 1024, // 100MB
		TestDataPath:          "",
		MCPClientPath:         "",
		GenerateReport:        true,
		ReportFormat:          "html",
		ReportOutputPath:      "",
		LogLevel:              "info",
		RetryAttempts:         2,
		RetryDelay:           time.Second,
	}
}

// TestProfile represents a named test configuration profile
type TestProfile struct {
	Name        string      `json:"name"`
	Description string      `json:"description"`
	Config      *TestConfig `json:"config"`
	CreatedAt   time.Time   `json:"createdAt"`
	UpdatedAt   time.Time   `json:"updatedAt"`
}

// ConfigurationSet represents a collection of test profiles
type ConfigurationSet struct {
	Version        string                  `json:"version"`
	DefaultProfile string                  `json:"defaultProfile"`
	Profiles       map[string]*TestProfile `json:"profiles"`
	LastUpdated    time.Time               `json:"lastUpdated"`
}

// LoadConfiguration loads test configurations from file
func (cm *ConfigManager) LoadConfiguration() error {
	if _, err := os.Stat(cm.configPath); os.IsNotExist(err) {
		// Create default configuration if file doesn't exist
		return cm.SaveConfiguration()
	}

	data, err := os.ReadFile(cm.configPath)
	if err != nil {
		return fmt.Errorf("failed to read config file: %w", err)
	}

	var configSet ConfigurationSet
	if err := json.Unmarshal(data, &configSet); err != nil {
		return fmt.Errorf("failed to parse config file: %w", err)
	}

	// Load profiles
	for name, profile := range configSet.Profiles {
		cm.profiles[name] = profile.Config
	}

	// Set current profile
	if configSet.DefaultProfile != "" {
		cm.currentProfile = configSet.DefaultProfile
	}

	return nil
}

// SaveConfiguration saves current configurations to file
func (cm *ConfigManager) SaveConfiguration() error {
	// Ensure config directory exists
	configDir := filepath.Dir(cm.configPath)
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	// Prepare configuration set
	configSet := &ConfigurationSet{
		Version:        "1.0",
		DefaultProfile: cm.currentProfile,
		Profiles:       make(map[string]*TestProfile),
		LastUpdated:    time.Now(),
	}

	// Add default profile if not exists
	if len(cm.profiles) == 0 {
		cm.profiles["default"] = cm.defaultConfig
	}

	// Convert profiles
	for name, config := range cm.profiles {
		configSet.Profiles[name] = &TestProfile{
			Name:        name,
			Description: getProfileDescription(name),
			Config:      config,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		}
	}

	// Serialize to JSON
	data, err := json.MarshalIndent(configSet, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to serialize config: %w", err)
	}

	// Write to file
	if err := os.WriteFile(cm.configPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	return nil
}

// getProfileDescription returns description for built-in profiles
func getProfileDescription(name string) string {
	descriptions := map[string]string{
		"default":     "Default testing configuration with balanced settings",
		"fast":        "Fast testing configuration for quick validation",
		"thorough":    "Comprehensive testing with all validations enabled",
		"development": "Development-friendly configuration with relaxed timeouts",
		"ci":          "Continuous integration optimized configuration",
		"security":    "Security-focused testing with enhanced scanning",
		"performance": "Performance-focused testing with load testing",
	}

	if desc, exists := descriptions[name]; exists {
		return desc
	}
	return "Custom test configuration profile"
}

// GetCurrentConfig returns the current active configuration
func (cm *ConfigManager) GetCurrentConfig() *TestConfig {
	if config, exists := cm.profiles[cm.currentProfile]; exists {
		return config
	}
	return cm.defaultConfig
}

// SetProfile sets the active configuration profile
func (cm *ConfigManager) SetProfile(profileName string) error {
	if _, exists := cm.profiles[profileName]; !exists {
		return fmt.Errorf("profile '%s' not found", profileName)
	}
	cm.currentProfile = profileName
	return nil
}

// CreateProfile creates a new configuration profile
func (cm *ConfigManager) CreateProfile(name, description string, config *TestConfig) error {
	if name == "" {
		return fmt.Errorf("profile name cannot be empty")
	}

	if config == nil {
		config = cm.defaultConfig
	}

	cm.profiles[name] = config
	return nil
}

// DeleteProfile removes a configuration profile
func (cm *ConfigManager) DeleteProfile(name string) error {
	if name == "default" {
		return fmt.Errorf("cannot delete default profile")
	}

	if _, exists := cm.profiles[name]; !exists {
		return fmt.Errorf("profile '%s' not found", name)
	}

	delete(cm.profiles, name)

	// Reset to default if current profile was deleted
	if cm.currentProfile == name {
		cm.currentProfile = "default"
	}

	return nil
}

// ListProfiles returns all available profiles
func (cm *ConfigManager) ListProfiles() map[string]*TestConfig {
	profiles := make(map[string]*TestConfig)
	for name, config := range cm.profiles {
		profiles[name] = config
	}
	return profiles
}

// GetProfile returns a specific profile configuration
func (cm *ConfigManager) GetProfile(name string) (*TestConfig, error) {
	if config, exists := cm.profiles[name]; exists {
		return config, nil
	}
	return nil, fmt.Errorf("profile '%s' not found", name)
}

// UpdateProfile updates an existing profile
func (cm *ConfigManager) UpdateProfile(name string, config *TestConfig) error {
	if _, exists := cm.profiles[name]; !exists {
		return fmt.Errorf("profile '%s' not found", name)
	}

	cm.profiles[name] = config
	return nil
}

// CreateBuiltinProfiles creates standard built-in configuration profiles
func (cm *ConfigManager) CreateBuiltinProfiles() error {
	profiles := map[string]*TestConfig{
		"fast": {
			Timeout:                  2 * time.Minute,
			MaxConcurrentTests:       5,
			EnableParallelTesting:    true,
			ContinueOnFailure:        true,
			EnableSecurityScanning:   false,
			EnableLinting:           false,
			EnablePerformanceTesting: false,
			EnableIntegrationTesting: true,
			MCPProtocolVersion:      "2024-11-05",
			RequiredMethods:         []string{"initialize", "tools/list"},
			RequiredCapabilities:    []string{"tools"},
			MaxResponseTime:         2 * time.Second,
			MaxMemoryUsage:          200 * 1024 * 1024,
			GenerateReport:          true,
			ReportFormat:            "json",
			LogLevel:                "warn",
			RetryAttempts:           1,
			RetryDelay:             500 * time.Millisecond,
		},
		"thorough": {
			Timeout:                  15 * time.Minute,
			MaxConcurrentTests:       2,
			EnableParallelTesting:    true,
			ContinueOnFailure:        false,
			EnableSecurityScanning:   true,
			EnableLinting:           true,
			EnablePerformanceTesting: true,
			EnableIntegrationTesting: true,
			MCPProtocolVersion:      "2024-11-05",
			RequiredMethods:         []string{"initialize", "tools/list", "tools/call"},
			RequiredCapabilities:    []string{"tools"},
			MaxResponseTime:         500 * time.Millisecond,
			MaxMemoryUsage:          50 * 1024 * 1024,
			GenerateReport:          true,
			ReportFormat:            "html",
			LogLevel:                "debug",
			RetryAttempts:           3,
			RetryDelay:             2 * time.Second,
		},
		"development": {
			Timeout:                  10 * time.Minute,
			MaxConcurrentTests:       3,
			EnableParallelTesting:    true,
			ContinueOnFailure:        true,
			EnableSecurityScanning:   false,
			EnableLinting:           true,
			EnablePerformanceTesting: false,
			EnableIntegrationTesting: true,
			MCPProtocolVersion:      "2024-11-05",
			RequiredMethods:         []string{"initialize", "tools/list"},
			RequiredCapabilities:    []string{"tools"},
			MaxResponseTime:         3 * time.Second,
			MaxMemoryUsage:          150 * 1024 * 1024,
			GenerateReport:          true,
			ReportFormat:            "html",
			LogLevel:                "info",
			RetryAttempts:           2,
			RetryDelay:             time.Second,
		},
		"ci": {
			Timeout:                  8 * time.Minute,
			MaxConcurrentTests:       4,
			EnableParallelTesting:    true,
			ContinueOnFailure:        false,
			EnableSecurityScanning:   true,
			EnableLinting:           true,
			EnablePerformanceTesting: true,
			EnableIntegrationTesting: true,
			MCPProtocolVersion:      "2024-11-05",
			RequiredMethods:         []string{"initialize", "tools/list", "tools/call"},
			RequiredCapabilities:    []string{"tools"},
			MaxResponseTime:         time.Second,
			MaxMemoryUsage:          100 * 1024 * 1024,
			GenerateReport:          true,
			ReportFormat:            "xml",
			LogLevel:                "warn",
			RetryAttempts:           1,
			RetryDelay:             500 * time.Millisecond,
		},
		"security": {
			Timeout:                  20 * time.Minute,
			MaxConcurrentTests:       2,
			EnableParallelTesting:    true,
			ContinueOnFailure:        false,
			EnableSecurityScanning:   true,
			EnableLinting:           true,
			EnablePerformanceTesting: false,
			EnableIntegrationTesting: true,
			MCPProtocolVersion:      "2024-11-05",
			RequiredMethods:         []string{"initialize", "tools/list", "tools/call"},
			RequiredCapabilities:    []string{"tools"},
			MaxResponseTime:         time.Second,
			MaxMemoryUsage:          100 * 1024 * 1024,
			GenerateReport:          true,
			ReportFormat:            "html",
			LogLevel:                "info",
			RetryAttempts:           2,
			RetryDelay:             time.Second,
		},
		"performance": {
			Timeout:                  30 * time.Minute,
			MaxConcurrentTests:       1,
			EnableParallelTesting:    false,
			ContinueOnFailure:        false,
			EnableSecurityScanning:   false,
			EnableLinting:           false,
			EnablePerformanceTesting: true,
			EnableIntegrationTesting: true,
			MCPProtocolVersion:      "2024-11-05",
			RequiredMethods:         []string{"initialize", "tools/list", "tools/call"},
			RequiredCapabilities:    []string{"tools"},
			MaxResponseTime:         200 * time.Millisecond,
			MaxMemoryUsage:          50 * 1024 * 1024,
			GenerateReport:          true,
			ReportFormat:            "html",
			LogLevel:                "info",
			RetryAttempts:           3,
			RetryDelay:             2 * time.Second,
		},
	}

	for name, config := range profiles {
		cm.profiles[name] = config
	}

	return nil
}

// ValidateConfig validates a test configuration
func (cm *ConfigManager) ValidateConfig(config *TestConfig) error {
	if config == nil {
		return fmt.Errorf("config cannot be nil")
	}

	if config.Timeout <= 0 {
		return fmt.Errorf("timeout must be positive")
	}

	if config.MaxConcurrentTests <= 0 {
		return fmt.Errorf("maxConcurrentTests must be positive")
	}

	if config.MCPProtocolVersion == "" {
		return fmt.Errorf("MCPProtocolVersion cannot be empty")
	}

	if len(config.RequiredMethods) == 0 {
		return fmt.Errorf("at least one required method must be specified")
	}

	if config.MaxResponseTime <= 0 {
		return fmt.Errorf("maxResponseTime must be positive")
	}

	if config.MaxMemoryUsage <= 0 {
		return fmt.Errorf("maxMemoryUsage must be positive")
	}

	validReportFormats := map[string]bool{
		"json": true,
		"html": true,
		"xml":  true,
	}

	if !validReportFormats[config.ReportFormat] {
		return fmt.Errorf("invalid report format: %s", config.ReportFormat)
	}

	validLogLevels := map[string]bool{
		"debug": true,
		"info":  true,
		"warn":  true,
		"error": true,
	}

	if !validLogLevels[config.LogLevel] {
		return fmt.Errorf("invalid log level: %s", config.LogLevel)
	}

	if config.RetryAttempts < 0 {
		return fmt.Errorf("retryAttempts cannot be negative")
	}

	if config.RetryDelay < 0 {
		return fmt.Errorf("retryDelay cannot be negative")
	}

	return nil
}

// GetCurrentProfile returns the name of the current active profile
func (cm *ConfigManager) GetCurrentProfile() string {
	return cm.currentProfile
}

// ExportProfile exports a profile to a standalone JSON file
func (cm *ConfigManager) ExportProfile(profileName, outputPath string) error {
	config, err := cm.GetProfile(profileName)
	if err != nil {
		return err
	}

	profile := &TestProfile{
		Name:        profileName,
		Description: getProfileDescription(profileName),
		Config:      config,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	data, err := json.MarshalIndent(profile, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to serialize profile: %w", err)
	}

	if err := os.WriteFile(outputPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write profile file: %w", err)
	}

	return nil
}

// ImportProfile imports a profile from a JSON file
func (cm *ConfigManager) ImportProfile(filePath string) error {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("failed to read profile file: %w", err)
	}

	var profile TestProfile
	if err := json.Unmarshal(data, &profile); err != nil {
		return fmt.Errorf("failed to parse profile file: %w", err)
	}

	if err := cm.ValidateConfig(profile.Config); err != nil {
		return fmt.Errorf("invalid profile configuration: %w", err)
	}

	cm.profiles[profile.Name] = profile.Config
	return nil
}

// GetConfigSummary returns a summary of the current configuration
func (cm *ConfigManager) GetConfigSummary() map[string]interface{} {
	config := cm.GetCurrentConfig()
	
	return map[string]interface{}{
		"currentProfile":         cm.currentProfile,
		"totalProfiles":          len(cm.profiles),
		"timeout":                config.Timeout.String(),
		"maxConcurrentTests":     config.MaxConcurrentTests,
		"enableParallelTesting":  config.EnableParallelTesting,
		"continueOnFailure":      config.ContinueOnFailure,
		"enableSecurityScanning": config.EnableSecurityScanning,
		"enableLinting":          config.EnableLinting,
		"enablePerformanceTesting": config.EnablePerformanceTesting,
		"enableIntegrationTesting": config.EnableIntegrationTesting,
		"mcpProtocolVersion":     config.MCPProtocolVersion,
		"requiredMethods":        config.RequiredMethods,
		"requiredCapabilities":   config.RequiredCapabilities,
		"maxResponseTime":        config.MaxResponseTime.String(),
		"maxMemoryUsage":         config.MaxMemoryUsage,
		"generateReport":         config.GenerateReport,
		"reportFormat":           config.ReportFormat,
		"logLevel":               config.LogLevel,
		"retryAttempts":          config.RetryAttempts,
		"retryDelay":             config.RetryDelay.String(),
	}
}