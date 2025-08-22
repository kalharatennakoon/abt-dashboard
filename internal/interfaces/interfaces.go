package interfaces

import (
	"abt-dashboard/internal/models"
	"net/http"
)

// Aggregatable defines the interface for any data aggregation
type Aggregatable interface {
	// Aggregate processes a transaction and updates internal state
	Aggregate(transaction models.Transaction)

	// GetResults returns the aggregated results
	GetResults() interface{}

	// Reset clears the aggregation state
	Reset()

	// GetType returns the aggregation type identifier
	GetType() string
}

// ChartRenderer defines the interface for chart rendering
type ChartRenderer interface {
	// Render generates the HTML/CSS/JS for the chart
	Render(data interface{}) string

	// GetChartType returns the chart type (bar, line, pie, etc.)
	GetChartType() string

	// GetContainerID returns the DOM container ID for this chart
	GetContainerID() string
}

// APIEndpoint defines the interface for API endpoints
type APIEndpoint interface {
	// Handle processes the HTTP request and returns response
	Handle(w http.ResponseWriter, r *http.Request)

	// GetPath returns the API endpoint path
	GetPath() string

	// GetMethod returns the HTTP method (GET, POST, etc.)
	GetMethod() string

	// GetDescription returns a description of what this endpoint does
	GetDescription() string
}

// DataFilter defines the interface for filtering data
type DataFilter interface {
	// Apply filters the data based on the filter criteria
	Apply(data interface{}) interface{}

	// GetFilterName returns the name of this filter
	GetFilterName() string

	// SetParameters sets filter parameters from query string
	SetParameters(params map[string]string) error
}

// InsightProvider defines the interface for business insights
type InsightProvider interface {
	// GenerateInsight analyzes data and provides business insights
	GenerateInsight(data interface{}) models.Insight

	// GetInsightType returns the type of insight this provider generates
	GetInsightType() string

	// GetPriority returns the priority level for this insight
	GetPriority() int
}

// DashboardComponent defines the interface for dashboard components
type DashboardComponent interface {
	// GetHTML returns the HTML structure for this component
	GetHTML() string

	// GetCSS returns the CSS styles for this component
	GetCSS() string

	// GetJavaScript returns the JavaScript code for this component
	GetJavaScript() string

	// GetComponentID returns the unique identifier for this component
	GetComponentID() string

	// GetTitle returns the display title for this component
	GetTitle() string
}

// ConfigurableComponent defines components that can be configured
type ConfigurableComponent interface {
	DashboardComponent

	// Configure sets up the component with given configuration
	Configure(config map[string]interface{}) error

	// GetConfigSchema returns the configuration schema
	GetConfigSchema() map[string]interface{}
}

// RealtimeComponent defines components that support real-time updates
type RealtimeComponent interface {
	DashboardComponent

	// SupportsRealtime returns true if this component can update in real-time
	SupportsRealtime() bool

	// GetUpdateInterval returns the update interval in seconds
	GetUpdateInterval() int

	// ShouldUpdate determines if the component should update based on new data
	ShouldUpdate(newData interface{}) bool
}
