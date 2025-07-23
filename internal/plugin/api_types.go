package plugin

import (
	"encoding/json"
	"time"
)

// API-safe types for frontend binding (using strings for time fields)

// PluginInfoAPI represents plugin metadata for API
type PluginInfoAPI struct {
	ID           string            `json:"id"`
	Name         string            `json:"name"`
	Version      string            `json:"version"`
	Description  string            `json:"description"`
	Author       string            `json:"author"`
	Homepage     string            `json:"homepage,omitempty"`
	Repository   string            `json:"repository,omitempty"`
	License      string            `json:"license"`
	Tags         []string          `json:"tags,omitempty"`
	MinVersion   string            `json:"minVersion"`
	MaxVersion   string            `json:"maxVersion"`
	Config       *PluginConfigAPI  `json:"config,omitempty"`
	Permissions  []string          `json:"permissions,omitempty"` // Convert Permission to string
	Dependencies []DependencyAPI   `json:"dependencies,omitempty"`
	Metadata     map[string]string `json:"metadata,omitempty"`
}

// PluginConfigAPI represents plugin configuration for API
type PluginConfigAPI struct {
	Schema   json.RawMessage   `json:"schema"`
	Default  json.RawMessage   `json:"default"`
	Required []string          `json:"required"`
	Examples []json.RawMessage `json:"examples"`
}

// DependencyAPI represents plugin dependencies for API
type DependencyAPI struct {
	Name       string `json:"name"`
	Version    string `json:"version"`
	Type       string `json:"type"`
	Optional   bool   `json:"optional"`
	Repository string `json:"repository,omitempty"`
}

// PluginInstanceAPI represents a loaded plugin instance for API
type PluginInstanceAPI struct {
	Info      *PluginInfoAPI `json:"info"`
	Status    string         `json:"status"` // Convert PluginStatus to string
	Config    json.RawMessage `json:"config,omitempty"`
	LoadedAt  string         `json:"loadedAt"` // ISO 8601 string
	LastError string         `json:"lastError,omitempty"`
	Stats     *PluginStatsAPI `json:"stats"`
	Manifest  *PluginManifestAPI `json:"manifest"`
}

// PluginStatsAPI represents plugin statistics for API
type PluginStatsAPI struct {
	CallCount       int64  `json:"callCount"`
	TotalDuration   int64  `json:"totalDuration"`   // nanoseconds
	AverageDuration int64  `json:"averageDuration"` // nanoseconds
	ErrorCount      int64  `json:"errorCount"`
	LastUsed        string `json:"lastUsed"` // ISO 8601 string
	MemoryUsage     int64  `json:"memoryUsage"`
}

// PluginManifestAPI represents plugin manifest for API
type PluginManifestAPI struct {
	*PluginInfoAPI
	Files       []PluginFileAPI `json:"files"`
	Checksum    string          `json:"checksum"`
	Size        int64           `json:"size"`
	InstallPath string          `json:"installPath,omitempty"`
	Verified    bool            `json:"verified"`
	Signature   string          `json:"signature,omitempty"`
}

// PluginFileAPI represents a plugin file for API
type PluginFileAPI struct {
	Path     string `json:"path"`
	Size     int64  `json:"size"`
	Checksum string `json:"checksum"`
	Type     string `json:"type"`
	Platform string `json:"platform,omitempty"`
	Arch     string `json:"arch,omitempty"`
}

// MarketplacePluginAPI represents marketplace plugin for API
type MarketplacePluginAPI struct {
	*PluginInfoAPI
	DownloadURL   string               `json:"downloadUrl"`
	Screenshots   []string             `json:"screenshots,omitempty"`
	Documentation string               `json:"documentation,omitempty"`
	Reviews       []ReviewAPI          `json:"reviews,omitempty"`
	Stats         *MarketplaceStatsAPI `json:"stats,omitempty"`
	UpdatedAt     string               `json:"updatedAt"` // ISO 8601 string
	Featured      bool                 `json:"featured"`
	Verified      bool                 `json:"verified"`
	Category      string               `json:"category"`
	Price         *PriceAPI            `json:"price,omitempty"`
}

// ReviewAPI represents a user review for API
type ReviewAPI struct {
	UserID    string `json:"userId"`
	UserName  string `json:"userName"`
	Rating    int    `json:"rating"`
	Comment   string `json:"comment"`
	CreatedAt string `json:"createdAt"` // ISO 8601 string
	Helpful   int    `json:"helpful"`
}

// MarketplaceStatsAPI represents marketplace statistics for API
type MarketplaceStatsAPI struct {
	Downloads     int64    `json:"downloads"`
	Rating        float64  `json:"rating"`
	ReviewCount   int      `json:"reviewCount"`
	LastUpdated   string   `json:"lastUpdated"` // ISO 8601 string
	Compatibility []string `json:"compatibility"`
}

// PriceAPI represents pricing information for API
type PriceAPI struct {
	Amount   float64  `json:"amount"`
	Currency string   `json:"currency"`
	Type     string   `json:"type"`
	Trial    *TrialAPI `json:"trial,omitempty"`
}

// TrialAPI represents trial information for API
type TrialAPI struct {
	Duration int64    `json:"duration"` // nanoseconds
	Features []string `json:"features,omitempty"`
}

// ValidationResultAPI represents validation results for API
type ValidationResultAPI struct {
	Valid    bool                   `json:"valid"`
	Errors   []ValidationErrorAPI   `json:"errors,omitempty"`
	Warnings []ValidationErrorAPI   `json:"warnings,omitempty"`
	Info     []ValidationErrorAPI   `json:"info,omitempty"`
	Stats    *ValidationStatsAPI    `json:"stats,omitempty"`
}

// ValidationErrorAPI represents validation error for API
type ValidationErrorAPI struct {
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

// ValidationStatsAPI represents validation statistics for API
type ValidationStatsAPI struct {
	TotalChecks   int   `json:"totalChecks"`
	Duration      int64 `json:"duration"` // nanoseconds
	RulesApplied  int   `json:"rulesApplied"`
	FilesChecked  int   `json:"filesChecked"`
	LinesChecked  int   `json:"linesChecked"`
}

// TestResultAPI represents test results for API
type TestResultAPI struct {
	Passed      bool             `json:"passed"`
	Duration    int64            `json:"duration"` // nanoseconds
	Tests       []TestCaseAPI    `json:"tests"`
	Coverage    *CoverageAPI     `json:"coverage,omitempty"`
	Performance *PerformanceAPI  `json:"performance,omitempty"`
	Summary     string           `json:"summary"`
}

// TestCaseAPI represents test case for API
type TestCaseAPI struct {
	Name     string `json:"name"`
	Status   string `json:"status"`
	Duration int64  `json:"duration"` // nanoseconds
	Message  string `json:"message,omitempty"`
	Error    string `json:"error,omitempty"`
}

// CoverageAPI represents coverage information for API
type CoverageAPI struct {
	Lines     int     `json:"lines"`
	Covered   int     `json:"covered"`
	Percent   float64 `json:"percent"`
	Files     int     `json:"files"`
	Functions int     `json:"functions"`
}

// PerformanceAPI represents performance metrics for API
type PerformanceAPI struct {
	RequestsPerSecond float64 `json:"requestsPerSecond"`
	AverageLatency    int64   `json:"averageLatency"` // nanoseconds
	MaxLatency        int64   `json:"maxLatency"`     // nanoseconds
	MinLatency        int64   `json:"minLatency"`     // nanoseconds
	MemoryUsage       int64   `json:"memoryUsage"`
	CPUUsage          float64 `json:"cpuUsage"`
}

// PluginEventAPI represents plugin events for API
type PluginEventAPI struct {
	Type      string                 `json:"type"`
	Source    string                 `json:"source"`
	Target    string                 `json:"target"`
	Timestamp string                 `json:"timestamp"` // ISO 8601 string
	Data      map[string]interface{} `json:"data,omitempty"`
}

// Conversion functions between internal types and API types

// ToPluginInfoAPI converts PluginInfo to API-safe format
func ToPluginInfoAPI(info *PluginInfo) *PluginInfoAPI {
	if info == nil {
		return nil
	}
	
	apiInfo := &PluginInfoAPI{
		ID:          info.ID,
		Name:        info.Name,
		Version:     info.Version,
		Description: info.Description,
		Author:      info.Author,
		Homepage:    info.Homepage,
		Repository:  info.Repository,
		License:     info.License,
		Tags:        info.Tags,
		MinVersion:  info.MinVersion,
		MaxVersion:  info.MaxVersion,
		Metadata:    info.Metadata,
	}
	
	// Convert permissions
	apiInfo.Permissions = make([]string, len(info.Permissions))
	for i, perm := range info.Permissions {
		apiInfo.Permissions[i] = string(perm)
	}
	
	// Convert config
	if info.Config != nil {
		apiInfo.Config = &PluginConfigAPI{
			Schema:   info.Config.Schema,
			Default:  info.Config.Default,
			Required: info.Config.Required,
			Examples: info.Config.Examples,
		}
	}
	
	// Convert dependencies
	apiInfo.Dependencies = make([]DependencyAPI, len(info.Dependencies))
	for i, dep := range info.Dependencies {
		apiInfo.Dependencies[i] = DependencyAPI{
			Name:       dep.Name,
			Version:    dep.Version,
			Type:       dep.Type,
			Optional:   dep.Optional,
			Repository: dep.Repository,
		}
	}
	
	return apiInfo
}

// ToPluginInstanceAPI converts PluginInstance to API-safe format
func ToPluginInstanceAPI(instance *PluginInstance) *PluginInstanceAPI {
	if instance == nil {
		return nil
	}
	
	apiInstance := &PluginInstanceAPI{
		Info:      ToPluginInfoAPI(instance.Info),
		Status:    string(instance.Status),
		Config:    instance.Config,
		LoadedAt:  instance.LoadedAt.Format(time.RFC3339),
		LastError: instance.LastError,
	}
	
	// Convert stats
	if instance.Stats != nil {
		apiInstance.Stats = &PluginStatsAPI{
			CallCount:       instance.Stats.CallCount,
			TotalDuration:   int64(instance.Stats.TotalDuration),
			AverageDuration: int64(instance.Stats.AverageDuration),
			ErrorCount:      instance.Stats.ErrorCount,
			LastUsed:        instance.Stats.LastUsed.Format(time.RFC3339),
			MemoryUsage:     instance.Stats.MemoryUsage,
		}
	}
	
	// Convert manifest
	if instance.Manifest != nil {
		apiInstance.Manifest = ToPluginManifestAPI(instance.Manifest)
	}
	
	return apiInstance
}

// ToPluginManifestAPI converts PluginManifest to API-safe format
func ToPluginManifestAPI(manifest *PluginManifest) *PluginManifestAPI {
	if manifest == nil {
		return nil
	}
	
	apiManifest := &PluginManifestAPI{
		PluginInfoAPI: ToPluginInfoAPI(manifest.PluginInfo),
		Checksum:      manifest.Checksum,
		Size:          manifest.Size,
		InstallPath:   manifest.InstallPath,
		Verified:      manifest.Verified,
		Signature:     manifest.Signature,
	}
	
	// Convert files
	apiManifest.Files = make([]PluginFileAPI, len(manifest.Files))
	for i, file := range manifest.Files {
		apiManifest.Files[i] = PluginFileAPI{
			Path:     file.Path,
			Size:     file.Size,
			Checksum: file.Checksum,
			Type:     file.Type,
			Platform: file.Platform,
			Arch:     file.Arch,
		}
	}
	
	return apiManifest
}

// ToMarketplacePluginAPI converts MarketplacePlugin to API-safe format
func ToMarketplacePluginAPI(plugin *MarketplacePlugin) *MarketplacePluginAPI {
	if plugin == nil {
		return nil
	}
	
	apiPlugin := &MarketplacePluginAPI{
		PluginInfoAPI: ToPluginInfoAPI(plugin.PluginInfo),
		DownloadURL:   plugin.DownloadURL,
		Screenshots:   plugin.Screenshots,
		Documentation: plugin.Documentation,
		UpdatedAt:     plugin.UpdatedAt.Format(time.RFC3339),
		Featured:      plugin.Featured,
		Verified:      plugin.Verified,
		Category:      plugin.Category,
	}
	
	// Convert reviews
	apiPlugin.Reviews = make([]ReviewAPI, len(plugin.Reviews))
	for i, review := range plugin.Reviews {
		apiPlugin.Reviews[i] = ReviewAPI{
			UserID:    review.UserID,
			UserName:  review.UserName,
			Rating:    review.Rating,
			Comment:   review.Comment,
			CreatedAt: review.CreatedAt.Format(time.RFC3339),
			Helpful:   review.Helpful,
		}
	}
	
	// Convert stats
	if plugin.Stats != nil {
		apiPlugin.Stats = &MarketplaceStatsAPI{
			Downloads:     plugin.Stats.Downloads,
			Rating:        plugin.Stats.Rating,
			ReviewCount:   plugin.Stats.ReviewCount,
			LastUpdated:   plugin.Stats.LastUpdated.Format(time.RFC3339),
			Compatibility: plugin.Stats.Compatibility,
		}
	}
	
	// Convert price
	if plugin.Price != nil {
		apiPlugin.Price = &PriceAPI{
			Amount:   plugin.Price.Amount,
			Currency: plugin.Price.Currency,
			Type:     plugin.Price.Type,
		}
		
		if plugin.Price.Trial != nil {
			apiPlugin.Price.Trial = &TrialAPI{
				Duration: int64(plugin.Price.Trial.Duration),
				Features: plugin.Price.Trial.Features,
			}
		}
	}
	
	return apiPlugin
}

// ToPluginEventAPI converts PluginEvent to API-safe format
func ToPluginEventAPI(event *PluginEvent) *PluginEventAPI {
	if event == nil {
		return nil
	}
	
	return &PluginEventAPI{
		Type:      event.Type,
		Source:    event.Source,
		Target:    event.Target,
		Timestamp: event.Timestamp.Format(time.RFC3339),
		Data:      event.Data,
	}
}