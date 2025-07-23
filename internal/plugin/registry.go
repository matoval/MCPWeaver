package plugin

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"
)

// Registry handles plugin discovery and marketplace integration
type Registry struct {
	config     *ManagerConfig
	httpClient *http.Client
	cache      map[string]*MarketplacePlugin
}

// MarketplacePlugin represents a plugin in the marketplace
type MarketplacePlugin struct {
	*PluginInfo
	DownloadURL   string            `json:"downloadUrl"`
	Screenshots   []string          `json:"screenshots,omitempty"`
	Documentation string            `json:"documentation,omitempty"`
	Reviews       []Review          `json:"reviews,omitempty"`
	Stats         *MarketplaceStats `json:"stats,omitempty"`
	UpdatedAt     time.Time         `json:"updatedAt"`
	Featured      bool              `json:"featured"`
	Verified      bool              `json:"verified"`
	Category      string            `json:"category"`
	Price         *Price            `json:"price,omitempty"`
}

// Review represents a user review
type Review struct {
	UserID    string    `json:"userId"`
	UserName  string    `json:"userName"`
	Rating    int       `json:"rating"` // 1-5
	Comment   string    `json:"comment"`
	CreatedAt time.Time `json:"createdAt"`
	Helpful   int       `json:"helpful"`
}

// MarketplaceStats tracks plugin statistics
type MarketplaceStats struct {
	Downloads     int64   `json:"downloads"`
	Rating        float64 `json:"rating"`
	ReviewCount   int     `json:"reviewCount"`
	LastUpdated   time.Time `json:"lastUpdated"`
	Compatibility []string  `json:"compatibility"`
}

// Price represents plugin pricing information
type Price struct {
	Amount   float64 `json:"amount"`
	Currency string  `json:"currency"`
	Type     string  `json:"type"` // "free", "one-time", "subscription"
	Trial    *Trial  `json:"trial,omitempty"`
}

// Trial represents trial information
type Trial struct {
	Duration time.Duration `json:"duration"`
	Features []string      `json:"features,omitempty"`
}

// SearchRequest represents a plugin search request
type SearchRequest struct {
	Query      string   `json:"query"`
	Category   string   `json:"category,omitempty"`
	Tags       []string `json:"tags,omitempty"`
	MinRating  float64  `json:"minRating,omitempty"`
	MaxPrice   float64  `json:"maxPrice,omitempty"`
	Free       bool     `json:"free,omitempty"`
	Featured   bool     `json:"featured,omitempty"`
	Verified   bool     `json:"verified,omitempty"`
	Limit      int      `json:"limit,omitempty"`
	Offset     int      `json:"offset,omitempty"`
}

// SearchResponse represents search results
type SearchResponse struct {
	Plugins    []*MarketplacePlugin `json:"plugins"`
	Total      int                  `json:"total"`
	Categories []string             `json:"categories"`
	Tags       []string             `json:"tags"`
	Page       int                  `json:"page"`
	PerPage    int                  `json:"perPage"`
}

// NewRegistry creates a new plugin registry
func NewRegistry(config *ManagerConfig) *Registry {
	return &Registry{
		config: config,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
		cache: make(map[string]*MarketplacePlugin),
	}
}

// Search searches for plugins in the marketplace
func (r *Registry) Search(ctx context.Context, req *SearchRequest) (*SearchResponse, error) {
	if !r.config.EnableMarketplace {
		return nil, fmt.Errorf("marketplace is disabled")
	}
	
	// Build request URL
	url := fmt.Sprintf("%s/api/v1/plugins/search", r.config.MarketplaceURL)
	
	// Create HTTP request
	httpReq, err := http.NewRequestWithContext(ctx, "POST", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	
	// Add request body
	reqBody, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}
	
	httpReq.Body = io.NopCloser(strings.NewReader(string(reqBody)))
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("User-Agent", "MCPWeaver/1.0.0")
	
	// Execute request
	resp, err := r.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("marketplace returned status %d", resp.StatusCode)
	}
	
	// Parse response
	var searchResp SearchResponse
	if err := json.NewDecoder(resp.Body).Decode(&searchResp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}
	
	// Cache results
	for _, plugin := range searchResp.Plugins {
		r.cache[plugin.ID] = plugin
	}
	
	return &searchResp, nil
}

// GetPlugin retrieves a specific plugin from marketplace
func (r *Registry) GetPlugin(ctx context.Context, pluginID string) (*MarketplacePlugin, error) {
	if !r.config.EnableMarketplace {
		return nil, fmt.Errorf("marketplace is disabled")
	}
	
	// Check cache first
	if plugin, exists := r.cache[pluginID]; exists {
		return plugin, nil
	}
	
	// Build request URL
	url := fmt.Sprintf("%s/api/v1/plugins/%s", r.config.MarketplaceURL, pluginID)
	
	// Create HTTP request
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	
	req.Header.Set("User-Agent", "MCPWeaver/1.0.0")
	
	// Execute request
	resp, err := r.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()
	
	if resp.StatusCode == http.StatusNotFound {
		return nil, fmt.Errorf("plugin not found: %s", pluginID)
	}
	
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("marketplace returned status %d", resp.StatusCode)
	}
	
	// Parse response
	var plugin MarketplacePlugin
	if err := json.NewDecoder(resp.Body).Decode(&plugin); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}
	
	// Cache result
	r.cache[pluginID] = &plugin
	
	return &plugin, nil
}

// Download downloads a plugin from the marketplace
func (r *Registry) Download(ctx context.Context, pluginID, downloadPath string) error {
	plugin, err := r.GetPlugin(ctx, pluginID)
	if err != nil {
		return fmt.Errorf("failed to get plugin info: %w", err)
	}
	
	// Create download request
	req, err := http.NewRequestWithContext(ctx, "GET", plugin.DownloadURL, nil)
	if err != nil {
		return fmt.Errorf("failed to create download request: %w", err)
	}
	
	req.Header.Set("User-Agent", "MCPWeaver/1.0.0")
	
	// Execute download
	resp, err := r.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to download plugin: %w", err)
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("download failed with status %d", resp.StatusCode)
	}
	
	// Save to file
	file, err := os.Create(downloadPath)
	if err != nil {
		return fmt.Errorf("failed to create download file: %w", err)
	}
	defer file.Close()
	
	_, err = io.Copy(file, resp.Body)
	if err != nil {
		return fmt.Errorf("failed to write download: %w", err)
	}
	
	return nil
}

// GetCategories retrieves available plugin categories
func (r *Registry) GetCategories(ctx context.Context) ([]string, error) {
	if !r.config.EnableMarketplace {
		return nil, fmt.Errorf("marketplace is disabled")
	}
	
	// Build request URL
	url := fmt.Sprintf("%s/api/v1/categories", r.config.MarketplaceURL)
	
	// Create HTTP request
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	
	req.Header.Set("User-Agent", "MCPWeaver/1.0.0")
	
	// Execute request
	resp, err := r.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("marketplace returned status %d", resp.StatusCode)
	}
	
	// Parse response
	var categories []string
	if err := json.NewDecoder(resp.Body).Decode(&categories); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}
	
	return categories, nil
}

// GetFeatured retrieves featured plugins from marketplace
func (r *Registry) GetFeatured(ctx context.Context, limit int) ([]*MarketplacePlugin, error) {
	if !r.config.EnableMarketplace {
		return nil, fmt.Errorf("marketplace is disabled")
	}
	
	// Build request URL
	url := fmt.Sprintf("%s/api/v1/plugins/featured?limit=%d", r.config.MarketplaceURL, limit)
	
	// Create HTTP request
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	
	req.Header.Set("User-Agent", "MCPWeaver/1.0.0")
	
	// Execute request
	resp, err := r.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("marketplace returned status %d", resp.StatusCode)
	}
	
	// Parse response
	var plugins []*MarketplacePlugin
	if err := json.NewDecoder(resp.Body).Decode(&plugins); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}
	
	// Cache results
	for _, plugin := range plugins {
		r.cache[plugin.ID] = plugin
	}
	
	return plugins, nil
}

// CheckUpdates checks for updates to installed plugins
func (r *Registry) CheckUpdates(ctx context.Context, installedPlugins map[string]*PluginInstance) (map[string]*MarketplacePlugin, error) {
	if !r.config.EnableMarketplace {
		return nil, fmt.Errorf("marketplace is disabled")
	}
	
	updates := make(map[string]*MarketplacePlugin)
	
	for pluginID, instance := range installedPlugins {
		// Get latest version from marketplace
		latest, err := r.GetPlugin(ctx, pluginID)
		if err != nil {
			// Plugin not found in marketplace, skip
			continue
		}
		
		// Compare versions (simplified version comparison)
		if latest.Version != instance.Info.Version {
			updates[pluginID] = latest
		}
	}
	
	return updates, nil
}