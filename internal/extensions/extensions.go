package extensions

import (
	"encoding/json"
	"fmt"
	"log"
	"path/filepath"
	"plugin"
	"reflect"

	"abt-dashboard/internal/config"
	"abt-dashboard/internal/models"
	"abt-dashboard/internal/plugins"
)

// ExtensionManager manages the loading and registration of extensions
type ExtensionManager struct {
	registry      *plugins.Registry
	config        *config.DashboardConfig
	loadedPlugins map[string]*plugin.Plugin
}

// NewExtensionManager creates a new extension manager
func NewExtensionManager(registry *plugins.Registry, config *config.DashboardConfig) *ExtensionManager {
	return &ExtensionManager{
		registry:      registry,
		config:        config,
		loadedPlugins: make(map[string]*plugin.Plugin),
	}
}

// LoadExtensionsFromConfig loads all extensions specified in the configuration
func (em *ExtensionManager) LoadExtensionsFromConfig() error {
	if extensionsConfig, exists := em.config.GetExtension("plugins"); exists {
		if pluginList, ok := extensionsConfig.([]interface{}); ok {
			for _, pluginConfig := range pluginList {
				if pluginMap, ok := pluginConfig.(map[string]interface{}); ok {
					if err := em.LoadExtension(pluginMap); err != nil {
						log.Printf("Failed to load extension: %v", err)
					}
				}
			}
		}
	}
	return nil
}

// LoadExtension loads a single extension based on configuration
func (em *ExtensionManager) LoadExtension(config map[string]interface{}) error {
	extensionType, ok := config["type"].(string)
	if !ok {
		return fmt.Errorf("extension type not specified")
	}

	switch extensionType {
	case "chart":
		return em.loadChartExtension(config)
	case "aggregator":
		return em.loadAggregatorExtension(config)
	case "insight":
		return em.loadInsightExtension(config)
	case "plugin":
		return em.loadPluginExtension(config)
	default:
		return fmt.Errorf("unknown extension type: %s", extensionType)
	}
}

// loadChartExtension loads a chart extension
func (em *ExtensionManager) loadChartExtension(config map[string]interface{}) error {
	chartType, ok := config["chart_type"].(string)
	if !ok {
		return fmt.Errorf("chart_type not specified for chart extension")
	}

	// Create a custom chart renderer based on configuration
	renderer := &CustomChartRenderer{
		chartType:    chartType,
		config:       config,
		htmlTemplate: getStringFromConfig(config, "html_template", ""),
		cssTemplate:  getStringFromConfig(config, "css_template", ""),
		jsTemplate:   getStringFromConfig(config, "js_template", ""),
	}

	return em.registry.RegisterRenderer(renderer)
}

// loadAggregatorExtension loads an aggregator extension
func (em *ExtensionManager) loadAggregatorExtension(config map[string]interface{}) error {
	aggregatorType, ok := config["aggregator_type"].(string)
	if !ok {
		return fmt.Errorf("aggregator_type not specified for aggregator extension")
	}

	// Create a custom aggregator based on configuration
	aggregator := &CustomAggregator{
		aggregatorType: aggregatorType,
		config:         config,
		fields:         getStringArrayFromConfig(config, "fields", []string{}),
		groupBy:        getStringFromConfig(config, "group_by", ""),
		aggregateBy:    getStringFromConfig(config, "aggregate_by", "sum"),
	}

	return em.registry.RegisterAggregator(aggregator)
}

// loadInsightExtension loads an insight provider extension
func (em *ExtensionManager) loadInsightExtension(config map[string]interface{}) error {
	insightType, ok := config["insight_type"].(string)
	if !ok {
		return fmt.Errorf("insight_type not specified for insight extension")
	}

	// Create a custom insight provider based on configuration
	provider := &CustomInsightProvider{
		insightType: insightType,
		config:      config,
		rules:       getInsightRulesFromConfig(config),
	}

	return em.registry.RegisterInsightProvider(provider)
}

// loadPluginExtension loads a compiled plugin extension
func (em *ExtensionManager) loadPluginExtension(config map[string]interface{}) error {
	pluginPath, ok := config["path"].(string)
	if !ok {
		return fmt.Errorf("path not specified for plugin extension")
	}

	// Load the plugin
	p, err := plugin.Open(pluginPath)
	if err != nil {
		return fmt.Errorf("failed to open plugin %s: %v", pluginPath, err)
	}

	// Store the loaded plugin
	pluginName := filepath.Base(pluginPath)
	em.loadedPlugins[pluginName] = p

	// Look for standard plugin interfaces
	if err := em.registerFromPlugin(p, "NewAggregator", em.registry.RegisterAggregator); err != nil {
		log.Printf("Failed to register aggregator from plugin %s: %v", pluginName, err)
	}

	if err := em.registerFromPlugin(p, "NewRenderer", em.registry.RegisterRenderer); err != nil {
		log.Printf("Failed to register renderer from plugin %s: %v", pluginName, err)
	}

	if err := em.registerFromPlugin(p, "NewInsightProvider", em.registry.RegisterInsightProvider); err != nil {
		log.Printf("Failed to register insight provider from plugin %s: %v", pluginName, err)
	}

	return nil
}

// registerFromPlugin registers components from a loaded plugin
func (em *ExtensionManager) registerFromPlugin(p *plugin.Plugin, symbolName string, registerFunc interface{}) error {
	symbol, err := p.Lookup(symbolName)
	if err != nil {
		return err // Symbol not found, which is okay
	}

	// Use reflection to call the registration function
	registerValue := reflect.ValueOf(registerFunc)
	symbolValue := reflect.ValueOf(symbol)

	// Call the symbol function to get the component
	results := symbolValue.Call([]reflect.Value{})
	if len(results) != 1 {
		return fmt.Errorf("expected 1 return value from %s", symbolName)
	}

	// Register the component
	registerResults := registerValue.Call([]reflect.Value{results[0]})
	if len(registerResults) == 1 && !registerResults[0].IsNil() {
		if err, ok := registerResults[0].Interface().(error); ok {
			return err
		}
	}

	return nil
}

// Custom implementations for configuration-based extensions

// CustomChartRenderer implements a configurable chart renderer
type CustomChartRenderer struct {
	chartType    string
	config       map[string]interface{}
	htmlTemplate string
	cssTemplate  string
	jsTemplate   string
}

func (ccr *CustomChartRenderer) GetChartType() string {
	return ccr.chartType
}

func (ccr *CustomChartRenderer) GetContainerID() string {
	return getStringFromConfig(ccr.config, "container_id", ccr.chartType+"-container")
}

func (ccr *CustomChartRenderer) Render(data interface{}) string {
	if ccr.htmlTemplate != "" {
		return executeTemplateWithData(ccr.htmlTemplate, data)
	}
	return fmt.Sprintf("<!-- Custom chart %s -->", ccr.chartType)
}

// CustomAggregator implements a configurable aggregator
type CustomAggregator struct {
	aggregatorType string
	config         map[string]interface{}
	fields         []string
	groupBy        string
	aggregateBy    string
	data           map[string]interface{}
}

func (ca *CustomAggregator) GetType() string {
	return ca.aggregatorType
}

func (ca *CustomAggregator) Aggregate(transaction models.Transaction) {
	// Generic aggregation logic based on configuration
	if ca.data == nil {
		ca.data = make(map[string]interface{})
	}

	// This would implement configurable aggregation logic
	// For now, it's a placeholder
}

func (ca *CustomAggregator) GetResults() interface{} {
	return ca.data
}

func (ca *CustomAggregator) Reset() {
	ca.data = make(map[string]interface{})
}

// CustomInsightProvider implements a configurable insight provider
type CustomInsightProvider struct {
	insightType string
	config      map[string]interface{}
	rules       []InsightRule
}

type InsightRule struct {
	Condition string      `json:"condition"`
	Message   string      `json:"message"`
	Severity  string      `json:"severity"`
	Threshold interface{} `json:"threshold"`
}

func (cip *CustomInsightProvider) GetInsightType() string {
	return cip.insightType
}

func (cip *CustomInsightProvider) GetPriority() int {
	return getIntFromConfig(cip.config, "priority", 1)
}

func (cip *CustomInsightProvider) GenerateInsight(data interface{}) models.Insight {
	// Generate insights based on configured rules
	// This would implement rule-based insight generation
	return models.Insight{
		Type:        cip.insightType,
		Title:       fmt.Sprintf("%s Insight", cip.insightType),
		Description: "Generated from custom rules",
		Severity:    "medium",
		Confidence:  0.8,
	}
}

// Utility functions for configuration parsing

func getStringFromConfig(config map[string]interface{}, key, defaultValue string) string {
	if value, ok := config[key].(string); ok {
		return value
	}
	return defaultValue
}

func getIntFromConfig(config map[string]interface{}, key string, defaultValue int) int {
	if value, ok := config[key].(float64); ok {
		return int(value)
	}
	if value, ok := config[key].(int); ok {
		return value
	}
	return defaultValue
}

func getStringArrayFromConfig(config map[string]interface{}, key string, defaultValue []string) []string {
	if value, ok := config[key].([]interface{}); ok {
		result := make([]string, len(value))
		for i, v := range value {
			if str, ok := v.(string); ok {
				result[i] = str
			}
		}
		return result
	}
	return defaultValue
}

func getInsightRulesFromConfig(config map[string]interface{}) []InsightRule {
	if rulesData, ok := config["rules"].([]interface{}); ok {
		rules := make([]InsightRule, 0, len(rulesData))
		for _, ruleData := range rulesData {
			if ruleBytes, err := json.Marshal(ruleData); err == nil {
				var rule InsightRule
				if json.Unmarshal(ruleBytes, &rule) == nil {
					rules = append(rules, rule)
				}
			}
		}
		return rules
	}
	return []InsightRule{}
}

func executeTemplateWithData(template string, data interface{}) string {
	// Simple template execution - in a real implementation,
	// this would use a proper template engine
	return template
}

// ExtensionRegistry provides a global registry for extensions
var ExtensionRegistry = &ExtensionManager{
	registry:      plugins.GlobalRegistry,
	loadedPlugins: make(map[string]*plugin.Plugin),
}

// RegisterExtension provides a simple way to register new extensions
func RegisterExtension(extensionType string, config map[string]interface{}) error {
	config["type"] = extensionType
	return ExtensionRegistry.LoadExtension(config)
}
