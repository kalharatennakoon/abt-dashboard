package plugins

import (
	"fmt"
	"sync"

	"abt-dashboard/internal/interfaces"
)

// Registry manages all registered plugins and components
type Registry struct {
	aggregators map[string]interfaces.Aggregatable
	renderers   map[string]interfaces.ChartRenderer
	endpoints   map[string]interfaces.APIEndpoint
	filters     map[string]interfaces.DataFilter
	insights    map[string]interfaces.InsightProvider
	components  map[string]interfaces.DashboardComponent
	mu          sync.RWMutex
}

// NewRegistry creates a new plugin registry
func NewRegistry() *Registry {
	return &Registry{
		aggregators: make(map[string]interfaces.Aggregatable),
		renderers:   make(map[string]interfaces.ChartRenderer),
		endpoints:   make(map[string]interfaces.APIEndpoint),
		filters:     make(map[string]interfaces.DataFilter),
		insights:    make(map[string]interfaces.InsightProvider),
		components:  make(map[string]interfaces.DashboardComponent),
	}
}

// Global registry instance
var GlobalRegistry = NewRegistry()

// RegisterAggregator registers a new data aggregator
func (r *Registry) RegisterAggregator(agg interfaces.Aggregatable) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	agType := agg.GetType()
	if _, exists := r.aggregators[agType]; exists {
		return fmt.Errorf("aggregator type %s already registered", agType)
	}

	r.aggregators[agType] = agg
	return nil
}

// RegisterRenderer registers a new chart renderer
func (r *Registry) RegisterRenderer(renderer interfaces.ChartRenderer) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	chartType := renderer.GetChartType()
	if _, exists := r.renderers[chartType]; exists {
		return fmt.Errorf("renderer type %s already registered", chartType)
	}

	r.renderers[chartType] = renderer
	return nil
}

// RegisterEndpoint registers a new API endpoint
func (r *Registry) RegisterEndpoint(endpoint interfaces.APIEndpoint) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	path := endpoint.GetPath()
	if _, exists := r.endpoints[path]; exists {
		return fmt.Errorf("endpoint %s already registered", path)
	}

	r.endpoints[path] = endpoint
	return nil
}

// RegisterFilter registers a new data filter
func (r *Registry) RegisterFilter(filter interfaces.DataFilter) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	name := filter.GetFilterName()
	if _, exists := r.filters[name]; exists {
		return fmt.Errorf("filter %s already registered", name)
	}

	r.filters[name] = filter
	return nil
}

// RegisterInsightProvider registers a new insight provider
func (r *Registry) RegisterInsightProvider(provider interfaces.InsightProvider) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	insightType := provider.GetInsightType()
	if _, exists := r.insights[insightType]; exists {
		return fmt.Errorf("insight provider %s already registered", insightType)
	}

	r.insights[insightType] = provider
	return nil
}

// RegisterComponent registers a new dashboard component
func (r *Registry) RegisterComponent(component interfaces.DashboardComponent) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	id := component.GetComponentID()
	if _, exists := r.components[id]; exists {
		return fmt.Errorf("component %s already registered", id)
	}

	r.components[id] = component
	return nil
}

// GetAggregator retrieves a registered aggregator by type
func (r *Registry) GetAggregator(agType string) (interfaces.Aggregatable, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	agg, exists := r.aggregators[agType]
	return agg, exists
}

// GetRenderer retrieves a registered renderer by chart type
func (r *Registry) GetRenderer(chartType string) (interfaces.ChartRenderer, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	renderer, exists := r.renderers[chartType]
	return renderer, exists
}

// GetEndpoint retrieves a registered endpoint by path
func (r *Registry) GetEndpoint(path string) (interfaces.APIEndpoint, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	endpoint, exists := r.endpoints[path]
	return endpoint, exists
}

// GetFilter retrieves a registered filter by name
func (r *Registry) GetFilter(name string) (interfaces.DataFilter, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	filter, exists := r.filters[name]
	return filter, exists
}

// GetInsightProvider retrieves a registered insight provider by type
func (r *Registry) GetInsightProvider(insightType string) (interfaces.InsightProvider, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	provider, exists := r.insights[insightType]
	return provider, exists
}

// GetComponent retrieves a registered component by ID
func (r *Registry) GetComponent(id string) (interfaces.DashboardComponent, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	component, exists := r.components[id]
	return component, exists
}

// ListAggregators returns all registered aggregator types
func (r *Registry) ListAggregators() []string {
	r.mu.RLock()
	defer r.mu.RUnlock()

	types := make([]string, 0, len(r.aggregators))
	for agType := range r.aggregators {
		types = append(types, agType)
	}
	return types
}

// ListRenderers returns all registered renderer types
func (r *Registry) ListRenderers() []string {
	r.mu.RLock()
	defer r.mu.RUnlock()

	types := make([]string, 0, len(r.renderers))
	for chartType := range r.renderers {
		types = append(types, chartType)
	}
	return types
}

// ListEndpoints returns all registered endpoint paths
func (r *Registry) ListEndpoints() []string {
	r.mu.RLock()
	defer r.mu.RUnlock()

	paths := make([]string, 0, len(r.endpoints))
	for path := range r.endpoints {
		paths = append(paths, path)
	}
	return paths
}

// ListComponents returns all registered component IDs
func (r *Registry) ListComponents() []string {
	r.mu.RLock()
	defer r.mu.RUnlock()

	ids := make([]string, 0, len(r.components))
	for id := range r.components {
		ids = append(ids, id)
	}
	return ids
}

// GetAvailableInsights returns all available insight types
func (r *Registry) GetAvailableInsights() []string {
	r.mu.RLock()
	defer r.mu.RUnlock()

	types := make([]string, 0, len(r.insights))
	for insightType := range r.insights {
		types = append(types, insightType)
	}
	return types
}
