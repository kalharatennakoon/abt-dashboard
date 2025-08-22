package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"sync"
)

// DashboardConfig holds the configuration for the entire dashboard
type DashboardConfig struct {
	Title       string                 `json:"title"`
	Description string                 `json:"description"`
	Theme       string                 `json:"theme"`
	Components  []ComponentConfig      `json:"components"`
	API         APIConfig              `json:"api"`
	Performance PerformanceConfig      `json:"performance"`
	Extensions  map[string]interface{} `json:"extensions"`
	mu          sync.RWMutex
}

// ComponentConfig defines configuration for a dashboard component
type ComponentConfig struct {
	ID          string                 `json:"id"`
	Type        string                 `json:"type"`
	Title       string                 `json:"title"`
	Enabled     bool                   `json:"enabled"`
	Position    int                    `json:"position"`
	Size        string                 `json:"size"` // "small", "medium", "large", "full-width"
	DataSource  string                 `json:"data_source"`
	RefreshRate int                    `json:"refresh_rate"` // seconds
	Options     map[string]interface{} `json:"options"`
}

// APIConfig defines configuration for API endpoints
type APIConfig struct {
	BaseURL     string            `json:"base_url"`
	Timeout     int               `json:"timeout"`    // seconds
	CacheTime   int               `json:"cache_time"` // seconds
	Compression bool              `json:"compression"`
	Headers     map[string]string `json:"headers"`
}

// PerformanceConfig defines performance-related settings
type PerformanceConfig struct {
	LoadTimeout       int  `json:"load_timeout"` // seconds
	ParallelLoading   bool `json:"parallel_loading"`
	LazyLoading       bool `json:"lazy_loading"`
	ShowIndicators    bool `json:"show_indicators"`
	CacheResponses    bool `json:"cache_responses"`
	CompressResponses bool `json:"compress_responses"`
}

// DefaultConfig returns a default configuration
func DefaultConfig() *DashboardConfig {
	return &DashboardConfig{
		Title:       "ABT Analytics Dashboard",
		Description: "Revenue Analytics & Business Intelligence",
		Theme:       "default",
		Components: []ComponentConfig{
			{
				ID:          "country-revenue",
				Type:        "table",
				Title:       "Country Revenue",
				Enabled:     true,
				Position:    1,
				Size:        "large",
				DataSource:  "/api/revenue/countries",
				RefreshRate: 300,
				Options: map[string]interface{}{
					"limit":      50,
					"pagination": true,
					"sorting":    true,
				},
			},
			{
				ID:          "top-products",
				Type:        "horizontal-bar-chart",
				Title:       "Top Products",
				Enabled:     true,
				Position:    2,
				Size:        "medium",
				DataSource:  "/api/products/top",
				RefreshRate: 300,
				Options: map[string]interface{}{
					"limit":     20,
					"show_bars": true,
					"gradient":  true,
				},
			},
			{
				ID:          "monthly-sales",
				Type:        "vertical-bar-chart",
				Title:       "Monthly Sales Trends",
				Enabled:     true,
				Position:    3,
				Size:        "medium",
				DataSource:  "/api/sales/by-month",
				RefreshRate: 300,
				Options: map[string]interface{}{
					"highlight_peaks": true,
					"show_summary":    true,
				},
			},
			{
				ID:          "regions",
				Type:        "dual-bar-chart",
				Title:       "Top Regions",
				Enabled:     true,
				Position:    4,
				Size:        "large",
				DataSource:  "/api/regions/top",
				RefreshRate: 300,
				Options: map[string]interface{}{
					"limit":      30,
					"dual_bars":  true,
					"scrollable": true,
				},
			},
		},
		API: APIConfig{
			BaseURL:     "http://localhost:8080",
			Timeout:     8,
			CacheTime:   300,
			Compression: true,
			Headers: map[string]string{
				"Accept": "application/json",
			},
		},
		Performance: PerformanceConfig{
			LoadTimeout:       10,
			ParallelLoading:   true,
			LazyLoading:       false,
			ShowIndicators:    true,
			CacheResponses:    true,
			CompressResponses: true,
		},
		Extensions: make(map[string]interface{}),
	}
}

// LoadConfig loads configuration from a JSON file
func LoadConfig(filename string) (*DashboardConfig, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %v", err)
	}

	var config DashboardConfig
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %v", err)
	}

	return &config, nil
}

// SaveConfig saves configuration to a JSON file
func (c *DashboardConfig) SaveConfig(filename string) error {
	c.mu.RLock()
	defer c.mu.RUnlock()

	data, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal config: %v", err)
	}

	if err := ioutil.WriteFile(filename, data, 0644); err != nil {
		return fmt.Errorf("failed to write config file: %v", err)
	}

	return nil
}

// AddComponent adds a new component to the configuration
func (c *DashboardConfig) AddComponent(component ComponentConfig) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.Components = append(c.Components, component)
}

// RemoveComponent removes a component by ID
func (c *DashboardConfig) RemoveComponent(id string) bool {
	c.mu.Lock()
	defer c.mu.Unlock()

	for i, component := range c.Components {
		if component.ID == id {
			c.Components = append(c.Components[:i], c.Components[i+1:]...)
			return true
		}
	}
	return false
}

// UpdateComponent updates an existing component
func (c *DashboardConfig) UpdateComponent(id string, updates ComponentConfig) bool {
	c.mu.Lock()
	defer c.mu.Unlock()

	for i, component := range c.Components {
		if component.ID == id {
			c.Components[i] = updates
			return true
		}
	}
	return false
}

// GetComponent retrieves a component by ID
func (c *DashboardConfig) GetComponent(id string) (*ComponentConfig, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	for _, component := range c.Components {
		if component.ID == id {
			return &component, true
		}
	}
	return nil, false
}

// GetEnabledComponents returns only enabled components, sorted by position
func (c *DashboardConfig) GetEnabledComponents() []ComponentConfig {
	c.mu.RLock()
	defer c.mu.RUnlock()

	var enabled []ComponentConfig
	for _, component := range c.Components {
		if component.Enabled {
			enabled = append(enabled, component)
		}
	}

	// Sort by position
	for i := 0; i < len(enabled)-1; i++ {
		for j := 0; j < len(enabled)-i-1; j++ {
			if enabled[j].Position > enabled[j+1].Position {
				enabled[j], enabled[j+1] = enabled[j+1], enabled[j]
			}
		}
	}

	return enabled
}

// SetExtension sets an extension configuration value
func (c *DashboardConfig) SetExtension(key string, value interface{}) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.Extensions[key] = value
}

// GetExtension retrieves an extension configuration value
func (c *DashboardConfig) GetExtension(key string) (interface{}, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	value, exists := c.Extensions[key]
	return value, exists
}
