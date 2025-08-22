package templates

import (
	"bytes"
	"fmt"
	"html/template"
	"strings"
)

// Template system for generating reusable dashboard components

// ChartTemplate provides a template for creating chart components
type ChartTemplate struct {
	ID         string
	Title      string
	ChartType  string
	DataSource string
	Width      string
	Height     string
	Options    map[string]interface{}
}

// TableTemplate provides a template for creating table components
type TableTemplate struct {
	ID         string
	Title      string
	DataSource string
	Columns    []ColumnConfig
	Options    map[string]interface{}
}

// ColumnConfig defines table column configuration
type ColumnConfig struct {
	Key       string `json:"key"`
	Title     string `json:"title"`
	Type      string `json:"type"` // "text", "number", "currency", "date"
	Sortable  bool   `json:"sortable"`
	Width     string `json:"width"`
	Alignment string `json:"alignment"` // "left", "center", "right"
}

// KPITemplate provides a template for KPI components
type KPITemplate struct {
	ID         string
	Title      string
	Value      string
	DataSource string
	Icon       string
	Color      string
	Format     string // "number", "currency", "percentage"
}

// GenerateChartHTML generates HTML for a chart component
func (ct *ChartTemplate) GenerateHTML() string {
	tmpl := `
<div class="section">
    <h2>{{.Icon}} {{.Title}}</h2>
    <p style="color: #666; margin-bottom: 15px; font-size: 0.9em;">
        <em>{{.Description}}</em>
    </p>
    <div id="{{.ID}}Container" style="width: {{.Width}}; height: {{.Height}};">
        <div class="loading">Loading {{.Title}} data...</div>
    </div>
</div>`

	data := map[string]interface{}{
		"ID":          ct.ID,
		"Title":       ct.Title,
		"Icon":        ct.getIcon(),
		"Description": ct.getDescription(),
		"Width":       ct.getWidth(),
		"Height":      ct.getHeight(),
	}

	return executeTemplate(tmpl, data)
}

// GenerateChartCSS generates CSS for a chart component
func (ct *ChartTemplate) GenerateCSS() string {
	tmpl := `
/* Styles for {{.ID}} chart */
.{{.ID}}-container {
    background: white;
    border-radius: 8px;
    border: 1px solid #ddd;
    padding: 25px;
    margin-bottom: 20px;
}

.{{.ID}}-chart {
    {{.ChartSpecificCSS}}
}

.{{.ID}}-item {
    transition: all 0.3s ease;
}

.{{.ID}}-item:hover {
    background-color: #f8f9fa;
    transform: translateY(-2px);
    box-shadow: 0 4px 8px rgba(0,0,0,0.1);
}`

	data := map[string]interface{}{
		"ID":               ct.ID,
		"ChartSpecificCSS": ct.getChartSpecificCSS(),
	}

	return executeTemplate(tmpl, data)
}

// GenerateChartJavaScript generates JavaScript for a chart component
func (ct *ChartTemplate) GenerateChartJavaScript() string {
	tmpl := `
// Load and display {{.Title}} data
async function load{{.FunctionName}}() {
    try {
        const response = await fetchWithTimeout('{{.DataSource}}{{.QueryParams}}');
        const data = await response.json();
        display{{.FunctionName}}Chart(data);
    } catch (error) {
        document.getElementById('{{.ID}}Container').innerHTML = ` + "`" + `
            <div class="error">
                <strong>Error loading {{.Title}}:</strong> ${error.message}
            </div>
        ` + "`" + `;
    }
}

// Display {{.Title}} chart
function display{{.FunctionName}}Chart(data) {
    if (!data || data.length === 0) {
        document.getElementById('{{.ID}}Container').innerHTML = ` + "`" + `
            <div class="error">No {{.Title}} data available.</div>
        ` + "`" + `;
        return;
    }

    {{.ChartRenderingLogic}}
    
    document.getElementById('{{.ID}}Container').innerHTML = chart;
}`

	data := map[string]interface{}{
		"ID":                  ct.ID,
		"Title":               ct.Title,
		"FunctionName":        ct.getFunctionName(),
		"DataSource":          ct.DataSource,
		"QueryParams":         ct.getQueryParams(),
		"ChartRenderingLogic": ct.getChartRenderingLogic(),
	}

	return executeTemplate(tmpl, data)
}

// GenerateTableHTML generates HTML for a table component
func (tt *TableTemplate) GenerateHTML() string {
	tmpl := `
<div class="section">
    <h2>📊 {{.Title}}</h2>
    <div class="table-controls">
        <div class="search-container">
            <input type="text" id="{{.ID}}Search" placeholder="Search {{.Title}}..." class="search-input">
        </div>
        <div class="table-info">
            <span id="{{.ID}}Info">Loading...</span>
        </div>
    </div>
    <div id="{{.ID}}Container" class="table-container">
        <div class="loading">Loading {{.Title}} data...</div>
    </div>
    <div class="pagination" id="{{.ID}}Pagination" style="display: none;">
        <button id="{{.ID}}PrevBtn" class="pagination-btn">Previous</button>
        <span id="{{.ID}}PageInfo" class="page-info"></span>
        <button id="{{.ID}}NextBtn" class="pagination-btn">Next</button>
    </div>
</div>`

	data := map[string]interface{}{
		"ID":    tt.ID,
		"Title": tt.Title,
	}

	return executeTemplate(tmpl, data)
}

// GenerateKPIHTML generates HTML for a KPI component
func (kt *KPITemplate) GenerateHTML() string {
	tmpl := `
<div class="kpi-card" id="{{.ID}}">
    <div class="kpi-icon" style="color: {{.Color}};">{{.Icon}}</div>
    <div class="kpi-content">
        <div class="kpi-value" id="{{.ID}}Value">Loading...</div>
        <div class="kpi-label">{{.Title}}</div>
    </div>
    <div class="kpi-change" id="{{.ID}}Change"></div>
</div>`

	data := map[string]interface{}{
		"ID":    kt.ID,
		"Title": kt.Title,
		"Icon":  kt.getIcon(),
		"Color": kt.getColor(),
	}

	return executeTemplate(tmpl, data)
}

// Helper methods for ChartTemplate
func (ct *ChartTemplate) getIcon() string {
	iconMap := map[string]string{
		"bar-chart":            "📊",
		"horizontal-bar-chart": "📈",
		"vertical-bar-chart":   "📊",
		"dual-bar-chart":       "📊",
		"line-chart":           "📈",
		"pie-chart":            "🥧",
		"area-chart":           "📈",
		"table":                "📋",
	}
	if icon, ok := iconMap[ct.ChartType]; ok {
		return icon
	}
	return "📊"
}

func (ct *ChartTemplate) getDescription() string {
	if desc, ok := ct.Options["description"].(string); ok {
		return desc
	}
	return fmt.Sprintf("Interactive %s visualization", ct.ChartType)
}

func (ct *ChartTemplate) getWidth() string {
	if ct.Width != "" {
		return ct.Width
	}
	return "100%"
}

func (ct *ChartTemplate) getHeight() string {
	if ct.Height != "" {
		return ct.Height
	}
	return "400px"
}

func (ct *ChartTemplate) getFunctionName() string {
	// Convert kebab-case to PascalCase
	parts := strings.Split(ct.ID, "-")
	var result strings.Builder
	for _, part := range parts {
		if len(part) > 0 {
			result.WriteString(strings.ToUpper(part[:1]) + part[1:])
		}
	}
	return result.String()
}

func (ct *ChartTemplate) getQueryParams() string {
	params := make([]string, 0)

	if limit, ok := ct.Options["limit"].(int); ok {
		params = append(params, fmt.Sprintf("limit=%d", limit))
	}

	if offset, ok := ct.Options["offset"].(int); ok {
		params = append(params, fmt.Sprintf("offset=%d", offset))
	}

	if len(params) > 0 {
		return "?" + strings.Join(params, "&")
	}
	return ""
}

func (ct *ChartTemplate) getChartSpecificCSS() string {
	switch ct.ChartType {
	case "bar-chart", "horizontal-bar-chart":
		return `
    display: flex;
    flex-direction: column;
    gap: 10px;
    max-height: 600px;
    overflow-y: auto;`
	case "vertical-bar-chart":
		return `
    display: flex;
    align-items: end;
    justify-content: space-between;
    height: 300px;
    padding: 20px 0;`
	case "dual-bar-chart":
		return `
    max-height: 800px;
    overflow-y: auto;`
	default:
		return `
    min-height: 300px;
    position: relative;`
	}
}

func (ct *ChartTemplate) getChartRenderingLogic() string {
	switch ct.ChartType {
	case "horizontal-bar-chart":
		return ct.getHorizontalBarLogic()
	case "vertical-bar-chart":
		return ct.getVerticalBarLogic()
	case "dual-bar-chart":
		return ct.getDualBarLogic()
	case "table":
		return ct.getTableLogic()
	default:
		return "// Chart rendering logic would go here"
	}
}

func (ct *ChartTemplate) getHorizontalBarLogic() string {
	return `
    const maxValue = Math.max(...data.map(item => item.value || item.units_sold || item.tx_count));
    
    const chart = \` + "`" + `
        <div class="{{.ID}}-chart-container">
            <div class="{{.ID}}-chart">
                \${data.map((item, index) => {
                    const value = item.value || item.units_sold || item.tx_count;
                    const width = (value / maxValue) * 100;
                    
                    return \` + "`" + `
                        <div class="{{.ID}}-item">
                            <div class="{{.ID}}-bar" style="width: \${width}%" 
                                 title="\${item.name || item.product_name}: \${formatNumber(value)}">
                                <span class="{{.ID}}-label">\${item.name || item.product_name}</span>
                                <span class="{{.ID}}-value">\${formatNumber(value)}</span>
                            </div>
                        </div>
                    \` + "`" + `;
                }).join('')}
            </div>
        </div>
    \` + "`" + `;`
}

func (ct *ChartTemplate) getVerticalBarLogic() string {
	return `
    const maxValue = Math.max(...data.map(item => item.value || item.units_sold));
    
    const chart = \` + "`" + `
        <div class="{{.ID}}-chart-container">
            <div class="{{.ID}}-chart">
                \${data.map((item, index) => {
                    const value = item.value || item.units_sold;
                    const height = (value / maxValue) * 100;
                    
                    return \` + "`" + `
                        <div class="{{.ID}}-bar" style="height: \${height}%">
                            <div class="{{.ID}}-value">\${formatNumber(value)}</div>
                            <div class="{{.ID}}-label">\${item.label || item.year_month}</div>
                        </div>
                    \` + "`" + `;
                }).join('')}
            </div>
        </div>
    \` + "`" + `;`
}

func (ct *ChartTemplate) getDualBarLogic() string {
	return `
    const maxRevenue = Math.max(...data.map(item => item.total_revenue_cents));
    const maxItems = Math.max(...data.map(item => item.items_sold));
    
    const chart = \` + "`" + `
        <div class="{{.ID}}-chart-container">
            <div class="{{.ID}}-chart">
                \${data.map((item, index) => {
                    const revenueWidth = (item.total_revenue_cents / maxRevenue) * 100;
                    const itemsWidth = (item.items_sold / maxItems) * 100;
                    
                    return \` + "`" + `
                        <div class="{{.ID}}-item">
                            <div class="{{.ID}}-rank">\${index + 1}</div>
                            <div class="{{.ID}}-name">\${item.region || item.name}</div>
                            <div class="{{.ID}}-bars">
                                <div class="{{.ID}}-revenue-bar" style="width: \${revenueWidth}%"></div>
                                <div class="{{.ID}}-items-bar" style="width: \${itemsWidth}%"></div>
                            </div>
                            <div class="{{.ID}}-stats">
                                <div class="{{.ID}}-revenue">\${formatCurrency(item.total_revenue_cents)}</div>
                                <div class="{{.ID}}-items">\${formatNumber(item.items_sold)} items</div>
                            </div>
                        </div>
                    \` + "`" + `;
                }).join('')}
            </div>
        </div>
    \` + "`" + `;`
}

func (ct *ChartTemplate) getTableLogic() string {
	return `
    const chart = \` + "`" + `
        <div class="{{.ID}}-table-container">
            <table class="{{.ID}}-table">
                <thead>
                    <tr>
                        \${Object.keys(data[0]).map(key => \` + "`" + `<th>\${key}</th>\` + "`" + `).join('')}
                    </tr>
                </thead>
                <tbody>
                    \${data.map(row => \` + "`" + `
                        <tr>
                            \${Object.values(row).map(value => \` + "`" + `<td>\${value}</td>\` + "`" + `).join('')}
                        </tr>
                    \` + "`" + `).join('')}
                </tbody>
            </table>
        </div>
    \` + "`" + `;`
}

// Helper methods for KPITemplate
func (kt *KPITemplate) getIcon() string {
	if kt.Icon != "" {
		return kt.Icon
	}
	return "📊"
}

func (kt *KPITemplate) getColor() string {
	if kt.Color != "" {
		return kt.Color
	}
	return "#667eea"
}

// Utility function to execute templates
func executeTemplate(tmplStr string, data interface{}) string {
	tmpl, err := template.New("component").Parse(tmplStr)
	if err != nil {
		return fmt.Sprintf("<!-- Template error: %v -->", err)
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return fmt.Sprintf("<!-- Template execution error: %v -->", err)
	}

	return buf.String()
}

// ComponentGenerator provides methods to generate complete components
type ComponentGenerator struct{}

// GenerateChart creates a complete chart component
func (cg *ComponentGenerator) GenerateChart(template ChartTemplate) map[string]string {
	return map[string]string{
		"html":       template.GenerateHTML(),
		"css":        template.GenerateCSS(),
		"javascript": template.GenerateChartJavaScript(),
	}
}

// GenerateTable creates a complete table component
func (cg *ComponentGenerator) GenerateTable(template TableTemplate) map[string]string {
	return map[string]string{
		"html":       template.GenerateHTML(),
		"css":        "/* Table CSS would go here */",
		"javascript": "/* Table JavaScript would go here */",
	}
}

// GenerateKPI creates a complete KPI component
func (cg *ComponentGenerator) GenerateKPI(template KPITemplate) map[string]string {
	return map[string]string{
		"html":       template.GenerateHTML(),
		"css":        "/* KPI CSS would go here */",
		"javascript": "/* KPI JavaScript would go here */",
	}
}
