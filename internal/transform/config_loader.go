package transform

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v2"
)

// ConfigLoader handles loading and managing transformation configurations
type ConfigLoader struct {
	configPath string
	config     TransformConfig
}

// NewConfigLoader creates a new configuration loader
func NewConfigLoader(configPath string) *ConfigLoader {
	return &ConfigLoader{
		configPath: configPath,
	}
}

// LoadConfig loads transformation configuration from file
func (cl *ConfigLoader) LoadConfig() (TransformConfig, error) {
	// Check if config file exists
	if _, err := os.Stat(cl.configPath); os.IsNotExist(err) {
		// Return default configuration if file doesn't exist
		return cl.getDefaultConfig(), nil
	}

	// Read configuration file
	configData, err := ioutil.ReadFile(cl.configPath)
	if err != nil {
		return TransformConfig{}, fmt.Errorf("failed to read config file: %w", err)
	}

	// Parse YAML configuration
	var yamlConfig struct {
		Transformation struct {
			EnableValidation   bool     `yaml:"enable_validation"`
			EnableOptimization bool     `yaml:"enable_optimization"`
			DateFormats        []string `yaml:"date_formats"`
			CurrencyFormats    []string `yaml:"currency_formats"`
			NullValues         []string `yaml:"null_values"`
			Defaults           struct {
				Country         string  `yaml:"country"`
				Region          string  `yaml:"region"`
				PriceMultiplier float64 `yaml:"price_multiplier"`
			} `yaml:"defaults"`
			CustomMappings map[string]string `yaml:"custom_mappings"`
			DataTypes      map[string]string `yaml:"data_types"`
		} `yaml:"transformation"`
	}

	if err := yaml.Unmarshal(configData, &yamlConfig); err != nil {
		return TransformConfig{}, fmt.Errorf("failed to parse config file: %w", err)
	}

	// Convert to TransformConfig
	config := TransformConfig{
		EnableValidation:   yamlConfig.Transformation.EnableValidation,
		EnableOptimization: yamlConfig.Transformation.EnableOptimization,
		DateFormats:        yamlConfig.Transformation.DateFormats,
		CurrencyFormats:    yamlConfig.Transformation.CurrencyFormats,
		NullValues:         yamlConfig.Transformation.NullValues,
		DefaultCountry:     yamlConfig.Transformation.Defaults.Country,
		DefaultRegion:      yamlConfig.Transformation.Defaults.Region,
		PriceMultiplier:    yamlConfig.Transformation.Defaults.PriceMultiplier,
		CustomMappings:     yamlConfig.Transformation.CustomMappings,
		DataTypes:          yamlConfig.Transformation.DataTypes,
	}

	// Apply defaults for missing values
	config = cl.applyDefaults(config)

	cl.config = config
	return config, nil
}

// getDefaultConfig returns a sensible default configuration
func (cl *ConfigLoader) getDefaultConfig() TransformConfig {
	config := TransformConfig{
		EnableValidation:   true,
		EnableOptimization: true,
		DateFormats: []string{
			"2006-01-02T15:04:05Z",      // RFC3339
			"2006-01-02",                // YYYY-MM-DD
			"2006-01-02 15:04:05",       // YYYY-MM-DD HH:MM:SS
			"01/02/2006",                // MM/DD/YYYY
			"02-01-2006",                // DD-MM-YYYY
			"2006/01/02",                // YYYY/MM/DD
			"Jan 2, 2006",               // Month DD, YYYY
			"2 Jan 2006",                // DD Month YYYY
			"2006-01-02T15:04:05-07:00", // RFC3339 with timezone
		},
		CurrencyFormats: []string{
			"USD", "EUR", "GBP", "JPY", "INR", "AUD", "CAD",
		},
		NullValues: []string{
			"", "NULL", "null", "N/A", "n/a", "NA", "na", "-",
			"None", "none", "undefined", "UNDEFINED",
		},
		DefaultCountry:  "",
		DefaultRegion:   "",
		PriceMultiplier: 100.0, // Convert dollars to cents
		CustomMappings: map[string]string{
			// Country mappings
			"usa":     "United States",
			"us":      "United States",
			"america": "United States",
			"uk":      "United Kingdom",
			"britain": "United Kingdom",
			"england": "United Kingdom",

			// Region mappings
			"region_n":  "North",
			"region_s":  "South",
			"region_e":  "East",
			"region_w":  "West",
			"region_ne": "Northeast",
			"region_nw": "Northwest",
			"region_se": "Southeast",
			"region_sw": "Southwest",
		},
		DataTypes: map[string]string{
			"transaction_id":   "string",
			"country":          "string",
			"region":           "string",
			"product_name":     "string",
			"price":            "float",
			"quantity":         "integer",
			"transaction_date": "datetime",
		},
	}

	return config
}

// applyDefaults ensures all required configuration values are set
func (cl *ConfigLoader) applyDefaults(config TransformConfig) TransformConfig {
	// Set default date formats if none provided
	if len(config.DateFormats) == 0 {
		config.DateFormats = cl.getDefaultConfig().DateFormats
	}

	// Set default null values if none provided
	if len(config.NullValues) == 0 {
		config.NullValues = cl.getDefaultConfig().NullValues
	}

	// Set default price multiplier if not set
	if config.PriceMultiplier == 0 {
		config.PriceMultiplier = 100.0
	}

	// Initialize custom mappings if nil
	if config.CustomMappings == nil {
		config.CustomMappings = make(map[string]string)
	}

	// Add default mappings if not present
	defaultMappings := cl.getDefaultConfig().CustomMappings
	for key, value := range defaultMappings {
		if _, exists := config.CustomMappings[key]; !exists {
			config.CustomMappings[key] = value
		}
	}

	// Initialize data types if nil
	if config.DataTypes == nil {
		config.DataTypes = cl.getDefaultConfig().DataTypes
	}

	return config
}

// LoadEnvironmentSpecificConfig loads configuration based on environment
func (cl *ConfigLoader) LoadEnvironmentSpecificConfig(env string) (TransformConfig, error) {
	// Base configuration
	baseConfig, err := cl.LoadConfig()
	if err != nil {
		return TransformConfig{}, err
	}

	// Environment-specific override file
	envConfigPath := cl.getEnvironmentConfigPath(env)

	if _, err := os.Stat(envConfigPath); os.IsNotExist(err) {
		// No environment-specific config, return base config
		return baseConfig, nil
	}

	// Load environment-specific overrides
	envConfigData, err := ioutil.ReadFile(envConfigPath)
	if err != nil {
		return baseConfig, nil // Fall back to base config
	}

	// Parse environment overrides
	var envConfig TransformConfig
	if err := yaml.Unmarshal(envConfigData, &envConfig); err != nil {
		return baseConfig, nil // Fall back to base config
	}

	// Merge configurations (environment overrides base)
	mergedConfig := cl.mergeConfigs(baseConfig, envConfig)

	return mergedConfig, nil
}

// getEnvironmentConfigPath returns path for environment-specific config
func (cl *ConfigLoader) getEnvironmentConfigPath(env string) string {
	dir := filepath.Dir(cl.configPath)
	basename := strings.TrimSuffix(filepath.Base(cl.configPath), filepath.Ext(cl.configPath))
	ext := filepath.Ext(cl.configPath)

	return filepath.Join(dir, fmt.Sprintf("%s.%s%s", basename, env, ext))
}

// mergeConfigs merges two configurations with override taking precedence
func (cl *ConfigLoader) mergeConfigs(base, override TransformConfig) TransformConfig {
	merged := base

	// Override boolean flags
	if override.EnableValidation != base.EnableValidation {
		merged.EnableValidation = override.EnableValidation
	}
	if override.EnableOptimization != base.EnableOptimization {
		merged.EnableOptimization = override.EnableOptimization
	}

	// Override arrays (replace completely if provided)
	if len(override.DateFormats) > 0 {
		merged.DateFormats = override.DateFormats
	}
	if len(override.CurrencyFormats) > 0 {
		merged.CurrencyFormats = override.CurrencyFormats
	}
	if len(override.NullValues) > 0 {
		merged.NullValues = override.NullValues
	}

	// Override strings
	if override.DefaultCountry != "" {
		merged.DefaultCountry = override.DefaultCountry
	}
	if override.DefaultRegion != "" {
		merged.DefaultRegion = override.DefaultRegion
	}

	// Override numeric values
	if override.PriceMultiplier != 0 {
		merged.PriceMultiplier = override.PriceMultiplier
	}

	// Merge custom mappings
	if override.CustomMappings != nil {
		if merged.CustomMappings == nil {
			merged.CustomMappings = make(map[string]string)
		}
		for key, value := range override.CustomMappings {
			merged.CustomMappings[key] = value
		}
	}

	// Merge data types
	if override.DataTypes != nil {
		if merged.DataTypes == nil {
			merged.DataTypes = make(map[string]string)
		}
		for key, value := range override.DataTypes {
			merged.DataTypes[key] = value
		}
	}

	return merged
}

// ValidateConfig validates the loaded configuration
func (cl *ConfigLoader) ValidateConfig(config TransformConfig) error {
	// Validate required fields
	if config.PriceMultiplier <= 0 {
		return fmt.Errorf("price_multiplier must be positive, got %f", config.PriceMultiplier)
	}

	if len(config.DateFormats) == 0 {
		return fmt.Errorf("at least one date format must be specified")
	}

	// Validate date formats
	validFormats := 0
	for _, format := range config.DateFormats {
		if cl.isValidDateFormat(format) {
			validFormats++
		}
	}

	if validFormats == 0 {
		return fmt.Errorf("no valid date formats found")
	}

	return nil
}

// isValidDateFormat checks if a date format string is valid
func (cl *ConfigLoader) isValidDateFormat(format string) bool {
	// Check if format contains the Go time format reference components
	requiredComponents := []string{"2006", "01", "02"}

	for _, component := range requiredComponents {
		if !strings.Contains(format, component) {
			return false
		}
	}

	return true
}

// SaveConfig saves the current configuration to file
func (cl *ConfigLoader) SaveConfig(config TransformConfig) error {
	// Create YAML structure
	yamlConfig := struct {
		Transformation struct {
			EnableValidation   bool     `yaml:"enable_validation"`
			EnableOptimization bool     `yaml:"enable_optimization"`
			DateFormats        []string `yaml:"date_formats"`
			CurrencyFormats    []string `yaml:"currency_formats"`
			NullValues         []string `yaml:"null_values"`
			Defaults           struct {
				Country         string  `yaml:"country"`
				Region          string  `yaml:"region"`
				PriceMultiplier float64 `yaml:"price_multiplier"`
			} `yaml:"defaults"`
			CustomMappings map[string]string `yaml:"custom_mappings"`
			DataTypes      map[string]string `yaml:"data_types"`
		} `yaml:"transformation"`
	}{}

	yamlConfig.Transformation.EnableValidation = config.EnableValidation
	yamlConfig.Transformation.EnableOptimization = config.EnableOptimization
	yamlConfig.Transformation.DateFormats = config.DateFormats
	yamlConfig.Transformation.CurrencyFormats = config.CurrencyFormats
	yamlConfig.Transformation.NullValues = config.NullValues
	yamlConfig.Transformation.Defaults.Country = config.DefaultCountry
	yamlConfig.Transformation.Defaults.Region = config.DefaultRegion
	yamlConfig.Transformation.Defaults.PriceMultiplier = config.PriceMultiplier
	yamlConfig.Transformation.CustomMappings = config.CustomMappings
	yamlConfig.Transformation.DataTypes = config.DataTypes

	// Marshal to YAML
	configData, err := yaml.Marshal(yamlConfig)
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	// Ensure directory exists
	if err := os.MkdirAll(filepath.Dir(cl.configPath), 0755); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	// Write to file
	if err := ioutil.WriteFile(cl.configPath, configData, 0644); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	return nil
}

// GetConfigInfo returns information about the current configuration
func (cl *ConfigLoader) GetConfigInfo() map[string]interface{} {
	return map[string]interface{}{
		"config_path":          cl.configPath,
		"config_exists":        cl.configFileExists(),
		"supported_formats":    len(cl.config.DateFormats),
		"custom_mappings":      len(cl.config.CustomMappings),
		"validation_enabled":   cl.config.EnableValidation,
		"optimization_enabled": cl.config.EnableOptimization,
	}
}

// configFileExists checks if the configuration file exists
func (cl *ConfigLoader) configFileExists() bool {
	_, err := os.Stat(cl.configPath)
	return !os.IsNotExist(err)
}

// LoadDefaultTransformationConfig loads default configuration for the system
func LoadDefaultTransformationConfig() TransformConfig {
	loader := NewConfigLoader("configs/data_transformation.yaml")
	config, err := loader.LoadConfig()
	if err != nil {
		// Return hardcoded defaults if config loading fails
		return loader.getDefaultConfig()
	}
	return config
}

// LoadTransformationConfigFromPath loads configuration from specific path
func LoadTransformationConfigFromPath(configPath string) (TransformConfig, error) {
	loader := NewConfigLoader(configPath)
	return loader.LoadConfig()
}

// LoadEnvironmentTransformationConfig loads configuration based on environment
func LoadEnvironmentTransformationConfig(env string) (TransformConfig, error) {
	loader := NewConfigLoader("configs/data_transformation.yaml")
	return loader.LoadEnvironmentSpecificConfig(env)
}
