# ABT Dashboard API Documentation

## Overview
The ABT Dashboard provides RESTful API endpoints for accessing analytics data. All endpoints return JSON responses with appropriate caching headers and gzip compression.

## Base URL
```
http://localhost:8080
```

## Response Format
All API responses follow this structure:
```json
{
  "data": [...],           // Array or object containing the requested data
  "status": "success",     // Status indicator
  "timestamp": "2025-08-22T19:00:00Z"  // Response timestamp
}
```

## Error Handling
Error responses include:
```json
{
  "error": "Error message",
  "status": "error",
  "code": 400,
  "timestamp": "2025-08-22T19:00:00Z"
}
```

## Performance Headers
All responses include performance optimization headers:
```
Cache-Control: public, max-age=300
Content-Encoding: gzip
ETag: "dashboard-data"
Content-Type: application/json; charset=utf-8
```

---

## Endpoints

### 1. Country Revenue Analysis

#### GET `/api/revenue/countries`
Retrieves revenue data aggregated by country and product.

**Parameters:**
- `limit` (integer, optional): Maximum number of records to return (default: 100, max: 1000)
- `offset` (integer, optional): Number of records to skip for pagination (default: 0)

**Example Request:**
```bash
curl "http://localhost:8080/api/revenue/countries?limit=50&offset=0"
```

**Response:**
```json
[
  {
    "country": "Sri Lanka",
    "product_name": "Widget A",
    "total_revenue_cents": 1250000,
    "number_of_transactions": 145
  },
  {
    "country": "India",
    "product_name": "Gadget B",
    "total_revenue_cents": 980000,
    "number_of_transactions": 89
  }
]
```

**Response Fields:**
- `country` (string): Country name
- `product_name` (string): Product name
- `total_revenue_cents` (integer): Total revenue in cents
- `number_of_transactions` (integer): Number of transactions

**Performance:**
- Typical response time: 200-500ms
- Supports pagination for large datasets
- Data sorted by revenue (descending)

---

### 2. Top Products Analysis

#### GET `/api/products/top`
Retrieves the most popular products by purchase frequency.

**Parameters:**
- `limit` (integer, optional): Number of top products to return (default: 20)
- `by` (string, optional): Sort criteria - "units" or "transactions" (default: "units")

**Example Request:**
```bash
curl "http://localhost:8080/api/products/top?limit=20&by=units"
```

**Response:**
```json
[
  {
    "product_name": "Widget Pro",
    "tx_count": 234,
    "units_sold": 1456,
    "stock_qty": 500
  },
  {
    "product_name": "Gadget Max",
    "tx_count": 198,
    "units_sold": 1234,
    "stock_qty": 750
  }
]
```

**Response Fields:**
- `product_name` (string): Product name
- `tx_count` (integer): Number of transactions
- `units_sold` (integer): Total units sold
- `stock_qty` (integer): Current stock quantity

**Performance:**
- Typical response time: 100-300ms
- Data sorted by units sold (descending)

---

### 3. Monthly Sales Trends

#### GET `/api/sales/by-month`
Retrieves sales data aggregated by month for trend analysis.

**Parameters:** None

**Example Request:**
```bash
curl "http://localhost:8080/api/sales/by-month"
```

**Response:**
```json
[
  {
    "year_month": "2024-01",
    "units_sold": 5670,
    "tx_count": 432,
    "revenue_cents": 12450000
  },
  {
    "year_month": "2024-02",
    "units_sold": 6234,
    "tx_count": 489,
    "revenue_cents": 13890000
  }
]
```

**Response Fields:**
- `year_month` (string): Month in YYYY-MM format
- `units_sold` (integer): Total units sold in the month
- `tx_count` (integer): Number of transactions
- `revenue_cents` (integer): Total revenue in cents

**Performance:**
- Typical response time: 150-400ms
- Data sorted chronologically

---

### 4. Regional Performance Analysis

#### GET `/api/regions/top`
Retrieves performance data for top regions by revenue.

**Parameters:**
- `limit` (integer, optional): Number of top regions to return (default: 30)

**Example Request:**
```bash
curl "http://localhost:8080/api/regions/top?limit=30"
```

**Response:**
```json
[
  {
    "region": "Western",
    "total_revenue_cents": 25670000,
    "items_sold": 12345,
    "number_of_transactions": 1567
  },
  {
    "region": "Central",
    "total_revenue_cents": 18940000,
    "items_sold": 9876,
    "number_of_transactions": 1234
  }
]
```

**Response Fields:**
- `region` (string): Region name
- `total_revenue_cents` (integer): Total revenue in cents
- `items_sold` (integer): Total items sold
- `number_of_transactions` (integer): Number of transactions

**Performance:**
- Typical response time: 200-600ms
- Data sorted by revenue (descending)

---

## Data Types and Formats

### Currency
All monetary values are represented in cents (integer) to avoid floating-point precision issues:
```json
{
  "total_revenue_cents": 1250000  // Represents $12,500.00
}
```

### Dates
Dates are formatted as ISO 8601 strings:
```json
{
  "year_month": "2024-03",           // YYYY-MM format for monthly data
  "timestamp": "2025-08-22T19:00:00Z" // Full ISO 8601 for timestamps
}
```

### Pagination
For endpoints supporting pagination:
```json
{
  "limit": 50,      // Records per page
  "offset": 100,    // Starting position
  "total": 1500     // Total available records (in future versions)
}
```

---

## Performance Characteristics

### Response Times
| Endpoint | Typical | Large Dataset | Maximum |
|----------|---------|---------------|---------|
| `/api/revenue/countries` | 200-500ms | 800ms | 2s |
| `/api/products/top` | 100-300ms | 400ms | 1s |
| `/api/sales/by-month` | 150-400ms | 600ms | 1.5s |
| `/api/regions/top` | 200-600ms | 800ms | 2s |

### Caching
- **Browser Cache**: 5 minutes (`max-age=300`)
- **ETag Support**: Conditional requests supported
- **Compression**: Gzip compression reduces payload by 70-80%

### Rate Limiting
Currently no rate limiting is implemented. For production deployment, consider:
- 100 requests per minute per IP
- Burst allowance of 20 requests
- 429 status code for exceeded limits

---

## Error Codes

### HTTP Status Codes
- `200 OK`: Successful request
- `400 Bad Request`: Invalid parameters
- `404 Not Found`: Endpoint not found
- `500 Internal Server Error`: Server error
- `503 Service Unavailable`: Temporary server overload

### Common Error Scenarios

#### Invalid Parameters
```bash
curl "http://localhost:8080/api/revenue/countries?limit=invalid"
```
```json
{
  "error": "Invalid limit parameter",
  "status": "error",
  "code": 400
}
```

#### Limit Exceeded
```bash
curl "http://localhost:8080/api/revenue/countries?limit=2000"
```
```json
{
  "error": "Limit exceeds maximum allowed value (1000)",
  "status": "error", 
  "code": 400
}
```

#### Server Error
```json
{
  "error": "Internal server error processing request",
  "status": "error",
  "code": 500
}
```

---

## Usage Examples

### JavaScript Fetch
```javascript
// Fetch country revenue data
async function getCountryRevenue() {
  try {
    const response = await fetch('/api/revenue/countries?limit=50');
    const data = await response.json();
    console.log(data);
  } catch (error) {
    console.error('Error:', error);
  }
}

// Fetch with timeout
async function fetchWithTimeout(url, timeout = 8000) {
  const controller = new AbortController();
  const timeoutId = setTimeout(() => controller.abort(), timeout);
  
  try {
    const response = await fetch(url, {
      signal: controller.signal,
      cache: 'default'
    });
    clearTimeout(timeoutId);
    return response;
  } catch (error) {
    clearTimeout(timeoutId);
    throw error;
  }
}
```

### Python Requests
```python
import requests

# Fetch top products
response = requests.get('http://localhost:8080/api/products/top?limit=20')
if response.status_code == 200:
    data = response.json()
    print(data)
else:
    print(f"Error: {response.status_code}")

# With error handling
def fetch_api_data(endpoint, params=None):
    try:
        response = requests.get(f'http://localhost:8080{endpoint}', 
                              params=params, timeout=8)
        response.raise_for_status()
        return response.json()
    except requests.RequestException as e:
        print(f"API Error: {e}")
        return None
```

### cURL Examples
```bash
# Get first 50 country revenue records
curl -H "Accept-Encoding: gzip" \
     "http://localhost:8080/api/revenue/countries?limit=50"

# Get top 10 products
curl "http://localhost:8080/api/products/top?limit=10"

# Get monthly sales data
curl "http://localhost:8080/api/sales/by-month"

# Get top 20 regions
curl "http://localhost:8080/api/regions/top?limit=20"

# Test with compression
curl -H "Accept-Encoding: gzip" -I \
     "http://localhost:8080/api/revenue/countries"
```

---

## Development and Testing

### Health Check
```bash
# Check if server is running
curl http://localhost:8080/api/revenue/countries?limit=1
```

### Performance Testing
```bash
# Test response time
curl -w "@curl-format.txt" http://localhost:8080/api/revenue/countries

# Where curl-format.txt contains:
#     time_namelookup:  %{time_namelookup}\n
#        time_connect:  %{time_connect}\n
#     time_appconnect:  %{time_appconnect}\n
#    time_pretransfer:  %{time_pretransfer}\n
#       time_redirect:  %{time_redirect}\n
#  time_starttransfer:  %{time_starttransfer}\n
#                     ----------\n
#          time_total:  %{time_total}\n
```

### Load Testing
```bash
# Apache Bench
ab -n 1000 -c 10 http://localhost:8080/api/revenue/countries

# Or with wrk
wrk -t12 -c400 -d30s http://localhost:8080/api/revenue/countries
```

---

## Future API Enhancements

### Planned Endpoints
- `GET /api/customers/ltv` - Customer lifetime value analysis
- `GET /api/inventory/optimization` - Inventory optimization insights
- `GET /api/market/penetration` - Market penetration analysis
- `GET /api/trends/forecasting` - Demand forecasting data
- `GET /api/performance/kpis` - Key performance indicators

### Extensibility
The API is designed for easy extension using the plugin system:
```go
// Add new endpoint
func (api *API) CustomAnalysis(w http.ResponseWriter, r *http.Request) {
    data := api.Agg.GetCustomAnalysis()
    api.writeJSON(w, data)
}

// Register route
mux.Handle("GET /api/custom/analysis", 
    gzipMiddleware(http.HandlerFunc(api.CustomAnalysis)))
```

---

## Security Considerations

### Current Implementation
- Input validation for parameters
- SQL injection prevention (no SQL, in-memory data)
- XSS prevention through JSON responses
- CORS headers can be configured

### Production Recommendations
- Add authentication (JWT tokens)
- Implement rate limiting
- Add request logging and monitoring
- Use HTTPS in production
- Validate and sanitize all inputs
- Add API versioning

---

**ABT Dashboard API - Comprehensive Analytics Data Access**
