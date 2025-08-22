package factory

import (
	"fmt"
	"strings"

	"abt-dashboard/internal/interfaces"
	"abt-dashboard/internal/models"
)

// AggregatorFactory creates aggregators based on type
type AggregatorFactory struct{}

// CreateAggregator creates a new aggregator of the specified type
func (f *AggregatorFactory) CreateAggregator(agType string, config map[string]interface{}) (interfaces.Aggregatable, error) {
	switch strings.ToLower(agType) {
	case "revenue-by-country":
		return NewRevenueByCountryAggregator(config), nil
	case "product-popularity":
		return NewProductPopularityAggregator(config), nil
	case "monthly-trends":
		return NewMonthlyTrendsAggregator(config), nil
	case "regional-performance":
		return NewRegionalPerformanceAggregator(config), nil
	case "time-series":
		return NewTimeSeriesAggregator(config), nil
	case "category-breakdown":
		return NewCategoryBreakdownAggregator(config), nil
	default:
		return nil, fmt.Errorf("unknown aggregator type: %s", agType)
	}
}

// ChartRendererFactory creates chart renderers based on type
type ChartRendererFactory struct{}

// CreateRenderer creates a new chart renderer of the specified type
func (f *ChartRendererFactory) CreateRenderer(chartType string, config map[string]interface{}) (interfaces.ChartRenderer, error) {
	switch strings.ToLower(chartType) {
	case "table":
		return NewTableRenderer(config), nil
	case "bar-chart":
		return NewBarChartRenderer(config), nil
	case "horizontal-bar-chart":
		return NewHorizontalBarChartRenderer(config), nil
	case "vertical-bar-chart":
		return NewVerticalBarChartRenderer(config), nil
	case "dual-bar-chart":
		return NewDualBarChartRenderer(config), nil
	case "line-chart":
		return NewLineChartRenderer(config), nil
	case "pie-chart":
		return NewPieChartRenderer(config), nil
	case "donut-chart":
		return NewDonutChartRenderer(config), nil
	case "area-chart":
		return NewAreaChartRenderer(config), nil
	case "scatter-plot":
		return NewScatterPlotRenderer(config), nil
	case "heatmap":
		return NewHeatmapRenderer(config), nil
	default:
		return nil, fmt.Errorf("unknown chart type: %s", chartType)
	}
}

// InsightProviderFactory creates insight providers
type InsightProviderFactory struct{}

// CreateInsightProvider creates a new insight provider of the specified type
func (f *InsightProviderFactory) CreateInsightProvider(providerType string, config map[string]interface{}) (interfaces.InsightProvider, error) {
	switch strings.ToLower(providerType) {
	case "trend-analysis":
		return NewTrendAnalysisProvider(config), nil
	case "anomaly-detection":
		return NewAnomalyDetectionProvider(config), nil
	case "performance-insights":
		return NewPerformanceInsightsProvider(config), nil
	case "recommendation-engine":
		return NewRecommendationEngineProvider(config), nil
	case "seasonal-analysis":
		return NewSeasonalAnalysisProvider(config), nil
	default:
		return nil, fmt.Errorf("unknown insight provider type: %s", providerType)
	}
}

// DashboardComponentFactory creates dashboard components
type DashboardComponentFactory struct{}

// CreateComponent creates a new dashboard component
func (f *DashboardComponentFactory) CreateComponent(componentType string, config map[string]interface{}) (interfaces.DashboardComponent, error) {
	switch strings.ToLower(componentType) {
	case "chart-component":
		return NewChartComponent(config), nil
	case "table-component":
		return NewTableComponent(config), nil
	case "kpi-component":
		return NewKPIComponent(config), nil
	case "text-component":
		return NewTextComponent(config), nil
	case "filter-component":
		return NewFilterComponent(config), nil
	case "insight-component":
		return NewInsightComponent(config), nil
	default:
		return nil, fmt.Errorf("unknown component type: %s", componentType)
	}
}

// Example implementations for demonstration (these would be separate files in production)

// BaseAggregator provides common functionality for aggregators
type BaseAggregator struct {
	agType string
	config map[string]interface{}
	data   interface{}
}

func (ba *BaseAggregator) GetType() string {
	return ba.agType
}

func (ba *BaseAggregator) Reset() {
	ba.data = nil
}

func (ba *BaseAggregator) GetResults() interface{} {
	return ba.data
}

func (ba *BaseAggregator) Aggregate(transaction models.Transaction) {
	// Default implementation - can be overridden by specific aggregators
	// This is a placeholder that does nothing
}

// Example aggregator implementations
func NewRevenueByCountryAggregator(config map[string]interface{}) interfaces.Aggregatable {
	return &RevenueByCountryAggregator{
		BaseAggregator: BaseAggregator{agType: "revenue-by-country", config: config},
		countries:      make(map[string]*models.CountryProductAgg),
	}
}

type RevenueByCountryAggregator struct {
	BaseAggregator
	countries map[string]*models.CountryProductAgg
}

func (r *RevenueByCountryAggregator) Aggregate(transaction models.Transaction) {
	key := transaction.Country + "-" + transaction.ProductName
	if agg, exists := r.countries[key]; exists {
		agg.TotalRevenue += transaction.UnitPriceCents * transaction.Quantity
		agg.NumberOfTx++
	} else {
		r.countries[key] = &models.CountryProductAgg{
			Country:      transaction.Country,
			ProductName:  transaction.ProductName,
			TotalRevenue: transaction.UnitPriceCents * transaction.Quantity,
			NumberOfTx:   1,
		}
	}

	// Convert map to slice for results
	results := make([]models.CountryProductAgg, 0, len(r.countries))
	for _, agg := range r.countries {
		results = append(results, *agg)
	}
	r.data = results
}

// Additional factory implementations would go here...
func NewProductPopularityAggregator(config map[string]interface{}) interfaces.Aggregatable {
	// Implementation would go here
	return &BaseAggregator{agType: "product-popularity", config: config}
}

func NewMonthlyTrendsAggregator(config map[string]interface{}) interfaces.Aggregatable {
	return &BaseAggregator{agType: "monthly-trends", config: config}
}

func NewRegionalPerformanceAggregator(config map[string]interface{}) interfaces.Aggregatable {
	return &BaseAggregator{agType: "regional-performance", config: config}
}

func NewTimeSeriesAggregator(config map[string]interface{}) interfaces.Aggregatable {
	return &BaseAggregator{agType: "time-series", config: config}
}

func NewCategoryBreakdownAggregator(config map[string]interface{}) interfaces.Aggregatable {
	return &BaseAggregator{agType: "category-breakdown", config: config}
}

// Example chart renderer implementations
func NewTableRenderer(config map[string]interface{}) interfaces.ChartRenderer {
	return &BaseRenderer{chartType: "table", config: config}
}

func NewBarChartRenderer(config map[string]interface{}) interfaces.ChartRenderer {
	return &BaseRenderer{chartType: "bar-chart", config: config}
}

func NewHorizontalBarChartRenderer(config map[string]interface{}) interfaces.ChartRenderer {
	return &BaseRenderer{chartType: "horizontal-bar-chart", config: config}
}

func NewVerticalBarChartRenderer(config map[string]interface{}) interfaces.ChartRenderer {
	return &BaseRenderer{chartType: "vertical-bar-chart", config: config}
}

func NewDualBarChartRenderer(config map[string]interface{}) interfaces.ChartRenderer {
	return &BaseRenderer{chartType: "dual-bar-chart", config: config}
}

func NewLineChartRenderer(config map[string]interface{}) interfaces.ChartRenderer {
	return &BaseRenderer{chartType: "line-chart", config: config}
}

func NewPieChartRenderer(config map[string]interface{}) interfaces.ChartRenderer {
	return &BaseRenderer{chartType: "pie-chart", config: config}
}

func NewDonutChartRenderer(config map[string]interface{}) interfaces.ChartRenderer {
	return &BaseRenderer{chartType: "donut-chart", config: config}
}

func NewAreaChartRenderer(config map[string]interface{}) interfaces.ChartRenderer {
	return &BaseRenderer{chartType: "area-chart", config: config}
}

func NewScatterPlotRenderer(config map[string]interface{}) interfaces.ChartRenderer {
	return &BaseRenderer{chartType: "scatter-plot", config: config}
}

func NewHeatmapRenderer(config map[string]interface{}) interfaces.ChartRenderer {
	return &BaseRenderer{chartType: "heatmap", config: config}
}

// Base renderer implementation
type BaseRenderer struct {
	chartType string
	config    map[string]interface{}
}

func (br *BaseRenderer) GetChartType() string {
	return br.chartType
}

func (br *BaseRenderer) GetContainerID() string {
	if id, ok := br.config["container_id"].(string); ok {
		return id
	}
	return br.chartType + "-container"
}

func (br *BaseRenderer) Render(data interface{}) string {
	// This would be implemented by specific renderers
	return fmt.Sprintf("<!-- %s chart would be rendered here -->", br.chartType)
}

// Insight provider implementations
func NewTrendAnalysisProvider(config map[string]interface{}) interfaces.InsightProvider {
	return &BaseInsightProvider{insightType: "trend-analysis", config: config}
}

func NewAnomalyDetectionProvider(config map[string]interface{}) interfaces.InsightProvider {
	return &BaseInsightProvider{insightType: "anomaly-detection", config: config}
}

func NewPerformanceInsightsProvider(config map[string]interface{}) interfaces.InsightProvider {
	return &BaseInsightProvider{insightType: "performance-insights", config: config}
}

func NewRecommendationEngineProvider(config map[string]interface{}) interfaces.InsightProvider {
	return &BaseInsightProvider{insightType: "recommendation-engine", config: config}
}

func NewSeasonalAnalysisProvider(config map[string]interface{}) interfaces.InsightProvider {
	return &BaseInsightProvider{insightType: "seasonal-analysis", config: config}
}

type BaseInsightProvider struct {
	insightType string
	config      map[string]interface{}
}

func (bip *BaseInsightProvider) GetInsightType() string {
	return bip.insightType
}

func (bip *BaseInsightProvider) GetPriority() int {
	if priority, ok := bip.config["priority"].(int); ok {
		return priority
	}
	return 1 // Default priority
}

func (bip *BaseInsightProvider) GenerateInsight(data interface{}) models.Insight {
	// This would be implemented by specific providers
	return models.Insight{
		Type:        bip.insightType,
		Title:       fmt.Sprintf("%s Insight", bip.insightType),
		Description: "Generated insight placeholder",
		Confidence:  0.8,
	}
}

// Component factory implementations
func NewChartComponent(config map[string]interface{}) interfaces.DashboardComponent {
	return &BaseComponent{componentType: "chart-component", config: config}
}

func NewTableComponent(config map[string]interface{}) interfaces.DashboardComponent {
	return &BaseComponent{componentType: "table-component", config: config}
}

func NewKPIComponent(config map[string]interface{}) interfaces.DashboardComponent {
	return &BaseComponent{componentType: "kpi-component", config: config}
}

func NewTextComponent(config map[string]interface{}) interfaces.DashboardComponent {
	return &BaseComponent{componentType: "text-component", config: config}
}

func NewFilterComponent(config map[string]interface{}) interfaces.DashboardComponent {
	return &BaseComponent{componentType: "filter-component", config: config}
}

func NewInsightComponent(config map[string]interface{}) interfaces.DashboardComponent {
	return &BaseComponent{componentType: "insight-component", config: config}
}

type BaseComponent struct {
	componentType string
	config        map[string]interface{}
}

func (bc *BaseComponent) GetComponentID() string {
	if id, ok := bc.config["id"].(string); ok {
		return id
	}
	return bc.componentType
}

func (bc *BaseComponent) GetTitle() string {
	if title, ok := bc.config["title"].(string); ok {
		return title
	}
	return bc.componentType
}

func (bc *BaseComponent) GetHTML() string {
	return fmt.Sprintf("<div id='%s'>%s component</div>", bc.GetComponentID(), bc.componentType)
}

func (bc *BaseComponent) GetCSS() string {
	return fmt.Sprintf(".%s { /* styles for %s */ }", bc.componentType, bc.componentType)
}

func (bc *BaseComponent) GetJavaScript() string {
	return fmt.Sprintf("// JavaScript for %s", bc.componentType)
}
