# ABT Analytics Dashboard

A comprehensive, extensible business intelligence dashboard for ABT revenu│   │   ├── config_loader.go    # Configuration management
│   │   └── transformations.go  # Built-in transformations
│   ├── handlers/               # HTTP request handlers
│   ├── metrics/                # Metrics aggregation
│   ├── models/                 # Data models
│   └── server/                 # HTTP server configuration
├── configs/
│   └── data_transformation.yaml # 🆕 Transformation configuration
├── docs/                        # 📚 Comprehensive documentation
├── testdata/                    # Test data in multiple formats
└── web/                         # Frontend assets (HTML, CSS, JS)
```

## 📖 Documentation

### Complete Documentation Suite
- **[📋 README.md](README.md)** - Project overview and quick start
- **[🔧 API Documentation](docs/API_DOCUMENTATION.md)** - REST API reference
- **[🚀 Deployment Guide](docs/DEPLOYMENT_GUIDE.md)** - Production deployment
- **[🔗 Extensibility Guide](docs/EXTENSIBILITY_GUIDE.md)** - Plugin development
- **[⚡ Performance Optimizations](docs/PERFORMANCE_OPTIMIZATIONS.md)** - Performance tuning
- **[🔄 Data Handling Guide](docs/DATA_HANDLING_GUIDE.md)** - Flexible data processing
- **[📊 Transformation Results](docs/DATA_TRANSFORMATION_RESULTS.md)** - Processing benchmarks

## 🧪 Testing & Quality

### Unit Testing Coverage
Our application includes comprehensive unit testing across core functionality:

#### ✅ **Models Package** (100% Coverage)
- Transaction validation with edge cases
- Aggregation calculations and performance metrics
- Data integrity and type validation
- Benchmark tests for performance validation

#### ✅ **Handlers Package** (100% API Coverage)  
- All REST API endpoints with mock testing
- HTTP request/response validation
- Pagination and parameter handling
- JSON response formatting and caching

#### ⚠️ **Transform Package** (Partial Coverage)
- Data format detection and conversion
- Transformation engine components
- Interface implementations (needs alignment)

### Running Tests
```bash
# Run all tests
go test ./...

# Run specific package tests
go test ./internal/models/... -v
go test ./internal/handlers/... -v

# Generate coverage report
go test -cover ./internal/models/... ./internal/handlers/...
```

### Test Coverage Report
📋 **[Detailed Test Coverage Report](./TEST_COVERAGE_REPORT.md)**

- **Total Test Files**: 4
- **Total Test Code**: 1,266 lines
- **Passing Test Suites**: 2/3 fully operational
- **Test Quality**: Comprehensive with mocks, benchmarks, edge cases

### Data Quality Monitoring
The system provides comprehensive data quality analysis in flexible mode:

```bash
# Run with quality reporting (experimental)
go run cmd/api/main.go -data=dataset.csv -config=config/data_transformation.yaml

# For stable operation, use traditional mode
go run cmd/api/main.go -data=dataset.csv -flexible=false
- Date parsing: 98.2% success
- Currency conversion: 100% success
- Geographic normalization: 94.1% success
```

### Sample Test Results
- **CSV Processing**: 8 records processed in 3.59ms (100% quality score)
- **JSON Processing**: 3 records processed in 526µs (100% quality score)
- **Format Detection**: 100% accuracy across test formats
- **Transformation Pipeline**: All transformations successful

## 🔧 Configuration

### Data Transformation Configuration (YAML)
```yaml
transformations:
  enable_currency_normalization: true
  enable_date_parsing: true
  enable_geographic_mapping: true
  
validation:
  required_fields: ["date", "amount", "country"]
  data_quality_threshold: 0.8
  
format_detection:
  auto_detect: true
  fallback_format: "csv"
```

### Environment Configuration
```bash
# Environment variables
export ABT_PORT=8080
export ABT_LOG_LEVEL=info
export ABT_CACHE_DURATION=5m
export ABT_DATA_PATH=./dataset.csv
export ABT_CONFIG_PATH=./configs/data_transformation.yaml
```

## 🚀 Performance Optimization

### Load Time Metrics
- **Dashboard Load**: < 10 seconds (typically 3-6 seconds)
- **API Response**: < 2 seconds per endpoint
- **Data Processing**: < 5ms per record (flexible handling)
- **Chart Rendering**: < 1 second per visualization

### Optimization Features
- **Gzip Compression**: Reduces payload size by 70-80%
- **Parallel Loading**: All components load simultaneously
- **Smart Caching**: Browser cache with ETag validation
- **Data Streaming**: Large datasets processed in chunks
- **Lazy Loading**: Charts render as they become visible

## 🔌 Extensibility

### Adding New Components
```go
// Register new chart type
chartFactory.Register("custom_chart", func(config ChartConfig) (Chart, error) {
    return &CustomChart{config: config}, nil
})

// Add new data transformation
transformEngine.RegisterTransformation("custom_transform", func(data []Record) []Record {
    // Custom transformation logic
    return transformedData
})
```

### Plugin Development
See [Extensibility Guide](docs/EXTENSIBILITY_GUIDE.md) for detailed plugin development instructions.

## 🛠️ Development

### Prerequisites
- Go 1.19+
- Node.js 14+ (for frontend tooling, optional)
- Make (optional, for build automation)

### Development Commands
```bash
# Install dependencies
go mod download

# Run tests
go test ./...

# Build for production
go build -o abt-dashboard cmd/api/main.go

# Run with hot reload (install air first: go install github.com/cosmtrek/air@latest)
air

# Format code
go fmt ./...

# Lint code
golangci-lint run
```

### Adding New Data Sources
1. Implement the `DataSource` interface in `internal/models/`
2. Register your data source in the factory
3. Add transformation rules in `configs/data_transformation.yaml`
4. Update documentation

## 📊 Data Processing Pipeline

### Flexible Data Handling Workflow
1. **Format Detection**: Automatic format identification (CSV/JSON/YAML/TSV)
2. **Data Parsing**: Format-specific parsing with error handling
3. **Field Mapping**: Flexible column mapping for different naming conventions
4. **Validation**: Configurable data validation rules
5. **Transformation**: Currency, date, and geographic normalization
6. **Quality Analysis**: Comprehensive quality metrics and reporting
7. **Optimization**: Performance optimization and caching

### Supported Data Variations
- **Date Formats**: ISO 8601, MM/DD/YYYY, DD/MM/YYYY, Unix timestamps, and more
- **Currency Formats**: $, €, £, ¥, ₹ with automatic conversion to cents
- **Geographic Variations**: Country name variations, region mapping
- **Field Names**: Flexible mapping for amount/revenue, date/timestamp, country/region

## 🔍 Monitoring & Observability

### Performance Monitoring
- Real-time performance metrics in browser console
- Server-side timing measurements
- Data processing quality scores
- Memory usage optimization

### Quality Monitoring
- Data validation success rates
- Transformation pipeline health
- Error tracking and reporting
- Quality score trends

## 📈 Future Enhancements

### Planned Features
- [ ] **Real-time Data Streaming**: WebSocket-based live updates
- [ ] **Advanced Analytics**: Machine learning insights
- [ ] **Export Capabilities**: PDF/Excel report generation
- [ ] **User Management**: Authentication and authorization
- [ ] **Custom Dashboards**: User-configurable dashboard layouts
- [ ] **Data Connectors**: Direct database and API integrations

### Roadmap
- **Q1 2024**: Real-time data streaming
- **Q2 2024**: Advanced analytics and ML insights
- **Q3 2024**: Enterprise features (auth, multi-tenancy)
- **Q4 2024**: Cloud deployment and scaling

## 🤝 Contributing

### Development Workflow
1. Fork the repository
2. Create a feature branch: `git checkout -b feature/amazing-feature`
3. Commit changes: `git commit -m 'Add amazing feature'`
4. Push to branch: `git push origin feature/amazing-feature`
5. Open a Pull Request

### Code Standards
- Follow Go best practices and idioms
- Maintain test coverage above 80%
- Use meaningful variable and function names
- Document public APIs with Go doc comments
- Follow the existing architectural patterns

## 📞 Support

### Getting Help
- **Documentation**: Check the [docs/](docs/) directory
- **Issues**: Open an issue on GitHub
- **Discussions**: Use GitHub Discussions for questions

### Common Issues
- **Port conflicts**: Change port with `-port=8081` flag
- **Data format issues**: Check [Data Handling Guide](docs/DATA_HANDLING_GUIDE.md)
- **Performance issues**: Review [Performance Guide](docs/PERFORMANCE_OPTIMIZATIONS.md)


## 🚀 Quick Start

### Prerequisites
- Go 1.19 or later
- Web browser (Chrome, Firefox, Safari, Edge)

### Installation & Running
```bash
# Clone the repository
git clone <repository-url>
cd abt-dashboard

# Download dependencies
go mod download

# Run the application (recommended stable mode)
go run cmd/api/main.go -data=dataset.csv -flexible=false

# Alternative: Build and run production binary
go build -o abt-dashboard cmd/api/main.go
./abt-dashboard -flexible=false
```

### 🌐 Accessing the Dashboard

#### **Option 1: Web Dashboard (Recommended)**
Open your browser and navigate to:
```
http://localhost:8080
```

The web interface provides:
- 📊 **Interactive Country Revenue Table** (sortable)
- 📈 **Top 20 Products Chart** (bar chart)
- � **Monthly Sales Trends** (line chart with peaks)
- 🗺️ **Top 30 Regions Chart** (regional performance)

#### **Option 2: Direct API Access**
Access raw JSON data via REST API endpoints:

**Country Revenue Data:**
```bash
# Top 10 countries by revenue
curl "http://localhost:8080/api/revenue/countries?limit=10"

# With pagination
curl "http://localhost:8080/api/revenue/countries?limit=5&offset=10"
```

**Top Products Data:**
```bash
# Top 20 most purchased products
curl "http://localhost:8080/api/products/top?limit=20"

# By transaction count (default) or units sold
curl "http://localhost:8080/api/products/top?limit=10&by=units"
```

**Monthly Sales Trends:**
```bash
# All monthly data
curl "http://localhost:8080/api/sales/by-month"
```

**Regional Performance:**
```bash
# Top 30 regions by revenue
curl "http://localhost:8080/api/regions/top?limit=30"

# Top 10 regions
curl "http://localhost:8080/api/regions/top?limit=10"
```

### 📊 Getting Insights from the Dashboard

#### **1. Country Revenue Analysis**
- **Endpoint**: `GET /api/revenue/countries?limit=N&offset=M`
- **Purpose**: Analyze revenue performance by country and product
- **Insights**: 
  - Which countries generate the most revenue
  - Top-selling products in each country
  - Transaction volume patterns
- **Example Response**:
```json
[
  {
    "country": "Germany",
    "product_name": "Product_281631",
    "total_revenue_cents": 1112395,
    "number_of_transactions": 7
  }
]
```

#### **2. Product Performance Analysis**
- **Endpoint**: `GET /api/products/top?limit=N&by=metric`
- **Purpose**: Identify best-performing products
- **Metrics**: `transactions` (default) or `units`
- **Insights**:
  - Most frequently purchased products
  - Products with highest unit sales
  - Current stock levels
- **Example Response**:
```json
[
  {
    "product_name": "Product_637255",
    "tx_count": 20,
    "units_sold": 39,
    "stock_qty": 77
  }
]
```

#### **3. Sales Trend Analysis**
- **Endpoint**: `GET /api/sales/by-month`
- **Purpose**: Track sales performance over time
- **Insights**:
  - Monthly sales volume trends
  - Revenue patterns and seasonality
  - Peak sales periods identification
- **Example Response**:
```json
[
  {
    "year_month": "2021-01",
    "units_sold": 323215,
    "tx_count": 129059,
    "revenue_cents": 7912119716
  }
]
```

#### **4. Regional Performance Analysis**
- **Endpoint**: `GET /api/regions/top?limit=N`
- **Purpose**: Compare regional market performance
- **Insights**:
  - Top-performing regions by revenue
  - Regional sales volume comparison
  - Market penetration analysis
- **Example Response**:
```json
[
  {
    "region": "California",
    "total_revenue_cents": 13308437808,
    "items_sold": 501656,
    "number_of_transactions": 200340
  }
]
```

### � Configuration Options

#### **Command Line Options**
```bash
# Basic usage
./abt-dashboard -data=your_data.csv -flexible=false

# All available options
./abt-dashboard \
  -data=dataset.csv \
  -inventory=inventory.csv \
  -static=web \
  -addr=:8080 \
  -config=config/data_transformation.yaml \
  -flexible=false
```

#### **Parameters Explained**
- `-data`: Path to your transaction data file (CSV format)
- `-inventory`: Path to inventory data (optional)
- `-static`: Directory containing web assets (default: web)
- `-addr`: Server address and port (default: :8080)
- `-config`: Configuration file for data transformation
- `-flexible`: Enable advanced data processing (experimental)

### 📈 Performance Features
- ✅ **Sub-10 Second Loading**: Optimized for fast dashboard load times
- ✅ **Parallel API Loading**: All components load simultaneously
- ✅ **HTTP Caching**: 5-minute browser cache with ETag validation
- ✅ **Gzip Compression**: 70-80% payload size reduction
- ✅ **Responsive Design**: Works on desktop, tablet, and mobile devices

### 🔍 Troubleshooting

#### **Common Issues & Solutions**

**1. Server Won't Start**
```bash
# Check if port 8080 is already in use
lsof -i :8080

# Use a different port
./abt-dashboard -addr=:8081

# Check if Go is installed
go version
```

**2. Dashboard Not Loading in Browser**
```bash
# Verify server is running - you should see:
# "Server listening on :8080"

# Test server response
curl http://localhost:8080

# Try different browser or incognito mode
# Clear browser cache (Ctrl+F5 or Cmd+Shift+R)
```

**3. API Endpoints Return 404**
```bash
# Correct API endpoint examples:
curl "http://localhost:8080/api/revenue/countries?limit=5"
curl "http://localhost:8080/api/products/top?limit=10"
curl "http://localhost:8080/api/sales/by-month"
curl "http://localhost:8080/api/regions/top?limit=5"

# Note: Endpoints are case-sensitive
```

**4. Data File Issues**
```bash
# Ensure data file exists
ls -la dataset.csv

# Check file permissions
chmod 644 dataset.csv

# Verify CSV format (first few lines)
head -5 dataset.csv

# Use absolute path if needed
./abt-dashboard -data=/full/path/to/dataset.csv
```

**5. Performance Issues**
```bash
# Use traditional mode for better stability
./abt-dashboard -flexible=false

# Check available memory
free -h  # Linux
vm_stat  # macOS

# Monitor server logs for errors
```

#### **Testing API Endpoints**

**Quick API Health Check:**
```bash
# Test all main endpoints
echo "Testing Country Revenue:"
curl -s "http://localhost:8080/api/revenue/countries?limit=3"

echo -e "\n\nTesting Top Products:"
curl -s "http://localhost:8080/api/products/top?limit=3"

echo -e "\n\nTesting Monthly Sales:"
curl -s "http://localhost:8080/api/sales/by-month" | head -c 200

echo -e "\n\nTesting Top Regions:"
curl -s "http://localhost:8080/api/regions/top?limit=3"
```

**Expected Response Format:**
All APIs return JSON arrays with proper HTTP headers:
```
Content-Type: application/json; charset=utf-8
Cache-Control: public, max-age=300
ETag: "dashboard-data"
```

#### **Browser Developer Tools**

**Check for JavaScript Errors:**
1. Open browser developer tools (F12)
2. Go to Console tab
3. Look for any red error messages
4. Refresh the page and monitor network requests

**Performance Monitoring:**
1. Open developer tools → Network tab
2. Refresh the dashboard
3. Check loading times for each API call
4. All requests should complete within 10 seconds

### 💡 Usage Examples

#### **Example 1: Business Performance Review**
```bash
# Start the dashboard
./abt-dashboard -data=q4_sales.csv -flexible=false

# Access insights via browser: http://localhost:8080
# Key insights available:
# - Top revenue-generating countries
# - Best-selling products by transaction count
# - Monthly sales trends for seasonal analysis
# - Regional performance comparison
```

#### **Example 2: API Data Integration**
```bash
# Get top 10 products for inventory planning
curl "http://localhost:8080/api/products/top?limit=10" | jq .

# Get monthly trends for forecasting
curl "http://localhost:8080/api/sales/by-month" | jq '.[] | select(.year_month >= "2023-01")'

# Get regional performance for market analysis
curl "http://localhost:8080/api/regions/top?limit=20" | jq '.[] | {region, revenue: .total_revenue_cents}'
```

#### **Example 3: Custom Data Analysis**
```bash
# Analyze specific country performance
curl "http://localhost:8080/api/revenue/countries" | jq '.[] | select(.country == "Germany")'

# Get products with low stock
curl "http://localhost:8080/api/products/top?limit=100" | jq '.[] | select(.stock_qty < 50)'

# Find high-volume, low-revenue products (efficiency analysis)
curl "http://localhost:8080/api/products/top?limit=50&by=units" | jq '.[] | select(.units_sold > 100)'
```

# JSON data with custom config
./abt-dashboard -data=sales_data.json -config=custom_transform.yaml

# YAML data
./abt-dashboard -data=sales_data.yaml

# Traditional parsing (fallback)
./abt-dashboard -data=sales_data.csv -flexible=false
```

### Data Transformation Features
- **Currency Normalization**: Handles $, €, £, ¥, ₹ and converts to cents
- **Date Format Flexibility**: Supports 15+ date formats including ISO, regional, and Unix timestamps
- **Geographic Standardization**: Maps country/region variations to standard names
- **Quality Validation**: Comprehensive data validation with configurable rules
- **Error Recovery**: Graceful handling of data quality issues

## 🏗️ Architecture

```
abt-dashboard/
├── cmd/api/main.go              # Application entry point with flexible data handling
├── internal/
│   ├── transform/               # 🆕 Flexible data handling system
│   │   ├── engine.go           # Data transformation orchestrator
│   │   ├── format_converter.go # Multi-format parsing
│   │   ├── data_handler.go     # High-level data processing
│   │   ├── transformations.go  # Built-in transformations
│   │   └── config_loader.go    # Configuration management
├── configs/
│   └── data_transformation.yaml # 🆕 Transformation configuration
├── internal/
│   ├── handlers/                # HTTP request handlers
│   ├── ingest/                  # Data ingestion and processing
│   ├── metrics/                 # Data aggregation logic
│   ├── models/                  # Data models and structures
│   ├── server/                  # HTTP server configuration
│   ├── interfaces/              # Extensibility interfaces
│   ├── plugins/                 # Plugin registry system
│   ├── factory/                 # Component factory patterns
│   ├── templates/               # Reusable component templates
│   ├── config/                  # Configuration management
│   └── extensions/              # Extension loading system
├── web/                         # Frontend assets
│   └── index.html              # Dashboard UI
├── config/
│   └── dashboard.json          # Dashboard configuration
├── testdata/                   # Sample data files
├── docs/                       # Documentation
└── tools/                      # Development utilities
```

## 📈 Data Flow

```
CSV Data → Ingest → Aggregation → API Endpoints → Frontend Charts
    ↓         ↓         ↓            ↓             ↓
[dataset.csv] → [Price Parsing] → [In-Memory Maps] → [JSON API] → [Interactive UI]
```

### Data Processing Pipeline
1. **Ingestion**: CSV parsing with decimal price handling
2. **Aggregation**: Real-time aggregation into multiple views
3. **API Layer**: RESTful endpoints with caching and compression
4. **Frontend**: Responsive charts with real-time updates

## 🔧 API Endpoints

### Core Analytics Endpoints
```
GET /api/revenue/countries?limit=100&offset=0
├── Country-level revenue analysis
├── Pagination support
└── Sortable by revenue (descending)

GET /api/products/top?limit=20&by=units
├── Top products by purchase frequency
├── Configurable limit
└── Sort by units or transactions

GET /api/sales/by-month
├── Monthly sales volume trends
├── Revenue and units sold
└── Peak month identification

GET /api/regions/top?limit=30
├── Regional performance analysis
├── Dual metrics (revenue + items)
└── Ranked by total revenue
```

### Response Format
```json
{
  "data": [...],
  "cache_headers": "5min",
  "compression": "gzip",
  "performance": "<2s response time"
}
```

## 🎨 Frontend Components

### Chart Types
- **Data Tables**: Sortable, paginated country/product data
- **Horizontal Bar Charts**: Product popularity rankings
- **Vertical Bar Charts**: Time-series sales trends
- **Dual Bar Charts**: Regional performance comparisons

### UI Features
- **Responsive Design**: Works on desktop, tablet, mobile
- **Loading States**: Visual feedback during data loading
- **Error Handling**: Graceful degradation for failed requests
- **Performance Monitoring**: Real-time load time tracking
- **Interactive Elements**: Hover effects, tooltips, animations

## ⚡ Performance Optimizations

### Backend Optimizations
```go
// Gzip compression middleware
gzipMiddleware(http.HandlerFunc(api.CountryRevenue))

// Enhanced caching headers
w.Header().Set("Cache-Control", "public, max-age=300")
w.Header().Set("ETag", "\"dashboard-data\"")

// Optimized server timeouts
ReadTimeout:       5 * time.Second,
WriteTimeout:      15 * time.Second,
```

### Frontend Optimizations
```javascript
// Parallel API loading
Promise.allSettled([
    loadCountryRevenue(),
    loadTopProducts(),
    loadMonthlySales(),
    loadRegions()
])

// Request timeout management
async function fetchWithTimeout(url, timeout = 8000) {
    // 8-second timeout per request
}
```

### Performance Results
- **Typical Load Time**: 1-3 seconds
- **Large Dataset**: 2-5 seconds maximum
- **Performance Buffer**: 5-7 seconds under 10s requirement
- **Subsequent Loads**: Near-instantaneous (cached)

## 🔌 Extensibility System

### Adding New Insights (4 Methods)

#### 1. Configuration-Only (No Code)
```json
{
  "components": [
    {
      "id": "customer-lifetime-value",
      "type": "kpi-component",
      "title": "Customer LTV",
      "data_source": "/api/customers/ltv"
    }
  ]
}
```

#### 2. Template-Based (Minimal Code)
```go
template := templates.ChartTemplate{
    ID:         "inventory-analysis",
    Title:      "Inventory Levels",
    ChartType:  "horizontal-bar-chart",
    DataSource: "/api/inventory/levels",
}
```

#### 3. Factory Pattern (Structured)
```go
factory := &factory.AggregatorFactory{}
aggregator, err := factory.CreateAggregator("seasonal-analysis", config)
plugins.GlobalRegistry.RegisterAggregator(aggregator)
```

#### 4. Plugin Development (Advanced)
```go
// External plugin: advanced_analytics.so
func NewAdvancedAnalytics() interfaces.InsightProvider {
    return &AdvancedAnalyticsProvider{}
}
```

### Available Component Types
- **11 Aggregation Patterns**: Revenue, product, regional, temporal, etc.
- **11 Chart Types**: Bar, line, pie, table, heatmap, scatter, etc.
- **5 Insight Categories**: Trends, anomalies, performance, recommendations
- **Unlimited Combinations**: Mix and match any components

## 📋 Development Guide

### Adding New Insight (5-Step Process)
1. **Define Data Model** → Add to `internal/models/models.go`
2. **Create Aggregator** → Implement `interfaces.Aggregatable`
3. **Add API Endpoint** → Add to `internal/handlers/handlers.go`
4. **Register Route** → Update `internal/server/server.go`
5. **Configure Frontend** → Update `config/dashboard.json`

### Example: Customer Satisfaction Analysis
```go
// 1. Data Model
type CustomerSatisfactionAgg struct {
    Period    string  `json:"period"`
    Score     float64 `json:"score"`
    Responses int64   `json:"responses"`
}

// 2. API Endpoint
func (api *API) CustomerSatisfaction(w http.ResponseWriter, r *http.Request) {
    data := api.Agg.GetCustomerSatisfaction()
    api.writeJSON(w, data)
}

// 3. Register Route
mux.Handle("GET /api/satisfaction/trends", 
    gzipMiddleware(http.HandlerFunc(api.CustomerSatisfaction)))
```

### Testing New Components
```bash
# Build and test
go build cmd/api/main.go
go test ./internal/...

# Test API endpoint
curl "http://localhost:8080/api/new/endpoint"

# Validate configuration
go run tools/validate_config.go config.json
```

## 🗄️ Data Schema

### Transaction Model
```go
type Transaction struct {
    ID             string    // Unique transaction ID
    Country        string    // e.g., "Sri Lanka"
    Region         string    // e.g., "Western"
    ProductName    string    // e.g., "Widget A"
    UnitPriceCents int64     // Price in cents (avoid float issues)
    Quantity       int64     // Units sold
    TxTime         time.Time // Transaction timestamp
}
```

### Aggregation Models
```go
// Country-Product aggregation
type CountryProductAgg struct {
    Country      string `json:"country"`
    ProductName  string `json:"product_name"`
    TotalRevenue int64  `json:"total_revenue_cents"`
    NumberOfTx   int64  `json:"number_of_transactions"`
}

// Product popularity
type ProductAgg struct {
    ProductName string `json:"product_name"`
    TxCount     int64  `json:"tx_count"`
    UnitsSold   int64  `json:"units_sold"`
    StockQty    int64  `json:"stock_qty"`
}

// Monthly trends
type MonthAgg struct {
    YearMonth    string `json:"year_month"`
    UnitsSold    int64  `json:"units_sold"`
    TxCount      int64  `json:"tx_count"`
    RevenueCents int64  `json:"revenue_cents"`
}

// Regional performance
type RegionAgg struct {
    Region       string `json:"region"`
    TotalRevenue int64  `json:"total_revenue_cents"`
    ItemsSold    int64  `json:"items_sold"`
    NumberOfTx   int64  `json:"number_of_transactions"`
}
```

## 🚀 Deployment

### Production Setup
```bash
# Build for production
go build -ldflags="-s -w" cmd/api/main.go

# Run with environment variables
export PORT=8080
export DATA_FILE=./dataset.csv
./main

# Or with Docker
docker build -t abt-dashboard .
docker run -p 8080:8080 abt-dashboard
```

### Configuration Options
```bash
# Environment Variables
PORT=8080                    # Server port
DATA_FILE=./dataset.csv      # Data source file
CONFIG_FILE=./config.json    # Dashboard configuration
CACHE_DURATION=300           # Cache duration in seconds
COMPRESSION=true             # Enable gzip compression
```

### Monitoring
- **Performance Metrics**: Built-in performance monitoring
- **Health Checks**: `/health` endpoint for load balancers
- **Logging**: Structured logging with request tracing
- **Metrics**: Response time, error rate, cache hit rate


## 🧪 Testing

### Running Tests
```bash
# Run all tests
go test ./...

# Run with coverage
go test -cover ./...

# Run specific test suite
go test ./internal/handlers/
go test ./internal/metrics/
go test ./internal/aggregators/
```

### Performance Testing
```bash
# Load testing with ab
ab -n 1000 -c 10 http://localhost:8080/api/revenue/countries

# API response time testing
./tools/performance_test.js

# Memory profiling
go test -memprofile=mem.prof ./...
go tool pprof mem.prof
```

## 🔧 Troubleshooting

### Common Issues

#### Dashboard Not Loading
```bash
# Check if server is running
lsof -i :8080

# Check logs
go run cmd/api/main.go 2>&1 | tee server.log

# Verify data file
file dataset.csv
head -n 5 dataset.csv
```

#### Performance Issues
```bash
# Monitor resource usage
top -p $(pgrep main)

# Check API response times
curl -w "@curl-format.txt" http://localhost:8080/api/revenue/countries

# Profile memory usage
go tool pprof http://localhost:8080/debug/pprof/heap
```

#### Extension Loading Errors
```bash
# Validate configuration
go run tools/validate_config.go config.json

# Check plugin compatibility
go run tools/check_plugins.go

# Debug extension loading
DEBUG=true go run cmd/api/main.go
```

## � Quick Reference

### 🚀 **Essential Commands**

#### **Start the Application**
```bash
# Recommended (stable mode)
go run cmd/api/main.go -data=dataset.csv -flexible=false

# Production build
go build -o abt-dashboard cmd/api/main.go
./abt-dashboard -flexible=false

# Custom port
./abt-dashboard -addr=:8081 -flexible=false
```

#### **Access Points**
- **Web Dashboard**: `http://localhost:8080`
- **API Base URL**: `http://localhost:8080/api`
- **Health Check**: `curl http://localhost:8080/`

### 🔗 **Complete API Reference**

#### **Country Revenue Analysis**
```bash
# Get all countries (paginated)
GET /api/revenue/countries?limit=100&offset=0

# Top 10 countries
curl "http://localhost:8080/api/revenue/countries?limit=10"

# Specific range
curl "http://localhost:8080/api/revenue/countries?limit=5&offset=20"
```

#### **Product Performance**
```bash
# Top products by transaction count (default)
GET /api/products/top?limit=20

# Top products by units sold
GET /api/products/top?limit=20&by=units

# Examples
curl "http://localhost:8080/api/products/top?limit=10"
curl "http://localhost:8080/api/products/top?limit=15&by=units"
```

#### **Sales Trends**
```bash
# Monthly sales data
GET /api/sales/by-month

# Example
curl "http://localhost:8080/api/sales/by-month"
```

#### **Regional Performance**
```bash
# Top regions by revenue
GET /api/regions/top?limit=30

# Examples
curl "http://localhost:8080/api/regions/top?limit=10"
curl "http://localhost:8080/api/regions/top?limit=5"
```

### 📊 **Data Insights Available**

#### **1. Revenue Insights**
- **What**: Country-wise revenue breakdown with product details
- **Use Case**: Identify top markets and best-selling products per country
- **Data**: Country, Product, Revenue (in cents), Transaction count

#### **2. Product Insights**
- **What**: Product popularity and inventory analysis
- **Use Case**: Inventory planning, product performance evaluation
- **Data**: Product name, Transaction count, Units sold, Stock quantity

#### **3. Trend Insights**
- **What**: Month-over-month sales performance
- **Use Case**: Seasonal analysis, forecasting, trend identification
- **Data**: Year-month, Units sold, Transaction count, Revenue

#### **4. Regional Insights**
- **What**: Geographic performance comparison
- **Use Case**: Market expansion, regional strategy planning
- **Data**: Region, Total revenue, Items sold, Transaction count

### 🛠️ **Command Line Options**

```bash
./abt-dashboard [OPTIONS]

Options:
  -data string          Path to transaction data file (default "dataset.csv")
  -inventory string     Path to inventory CSV file (default "dataset.csv")
  -static string        Path to static files directory (default "web")
  -addr string          Server listen address (default ":8080")
  -config string        Path to transformation config (default "config/data_transformation.yaml")
  -flexible bool        Use flexible data handling system (default true, recommend false)
```

### 🔧 **Testing Commands**

#### **Quick Health Check**
```bash
# Test server
curl http://localhost:8080/

# Test all APIs
curl -s "http://localhost:8080/api/revenue/countries?limit=1" | head -c 100
curl -s "http://localhost:8080/api/products/top?limit=1" | head -c 100
curl -s "http://localhost:8080/api/sales/by-month" | head -c 100
curl -s "http://localhost:8080/api/regions/top?limit=1" | head -c 100
```

#### **Performance Testing**
```bash
# Response time test
time curl -s "http://localhost:8080/api/revenue/countries?limit=100" > /dev/null

# Load test (requires 'ab' tool)
ab -n 100 -c 10 http://localhost:8080/api/revenue/countries
```

### 📋 **Response Format Examples**

#### **Country Revenue Response**
```json
[
  {
    "country": "Germany",
    "product_name": "Product_281631",
    "total_revenue_cents": 1112395,
    "number_of_transactions": 7
  }
]
```

#### **Product Performance Response**
```json
[
  {
    "product_name": "Product_637255",
    "tx_count": 20,
    "units_sold": 39,
    "stock_qty": 77
  }
]
```

#### **Monthly Sales Response**
```json
[
  {
    "year_month": "2021-01",
    "units_sold": 323215,
    "tx_count": 129059,
    "revenue_cents": 7912119716
  }
]
```

#### **Regional Performance Response**
```json
[
  {
    "region": "California",
    "total_revenue_cents": 13308437808,
    "items_sold": 501656,
    "number_of_transactions": 200340
  }
]
```

### ⚡ **Quick Start Checklist**

- [ ] Go 1.19+ installed (`go version`)
- [ ] Repository cloned
- [ ] Dependencies downloaded (`go mod download`)
- [ ] Data file available (`ls dataset.csv`)
- [ ] Server started (`go run cmd/api/main.go -flexible=false`)
- [ ] Browser opened to `http://localhost:8080`
- [ ] All 4 dashboard sections loading
- [ ] API endpoints responding to curl tests

### 🆘 **Common Issues**

| Issue | Solution |
|-------|----------|
| Port 8080 in use | Use `-addr=:8081` |
| Dashboard not loading | Check console for errors, try incognito mode |
| API returns 404 | Verify exact endpoint URLs |
| Slow performance | Use `-flexible=false` mode |
| Data file not found | Check file path and permissions |

---
