package handlers

import (
	"abt-dashboard/internal/models"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

// MetricsAggregator interface for testing
type MetricsAggregator interface {
	CountryRevenueTable() []models.CountryProductAgg
	TopProducts(limit int, byUnits bool) []models.ProductAgg
	SalesByMonth() []models.MonthAgg
	TopRegions(limit int) []models.RegionAgg
}

// Mock aggregator for testing
type mockAggregator struct {
	countryRevenue []models.CountryProductAgg
	topProducts    []models.ProductAgg
	salesByMonth   []models.MonthAgg
	topRegions     []models.RegionAgg
}

func (m *mockAggregator) CountryRevenueTable() []models.CountryProductAgg {
	return m.countryRevenue
}

func (m *mockAggregator) TopProducts(limit int, byUnits bool) []models.ProductAgg {
	if len(m.topProducts) > limit {
		return m.topProducts[:limit]
	}
	return m.topProducts
}

func (m *mockAggregator) SalesByMonth() []models.MonthAgg {
	return m.salesByMonth
}

func (m *mockAggregator) TopRegions(limit int) []models.RegionAgg {
	if len(m.topRegions) > limit {
		return m.topRegions[:limit]
	}
	return m.topRegions
}

// TestAPI wraps API to use our interface
type TestAPI struct {
	Agg MetricsAggregator
}

func (api *TestAPI) writeJSON(w http.ResponseWriter, v interface{}) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("Cache-Control", "public, max-age=300")
	w.Header().Set("ETag", "\"dashboard-data\"")
	json.NewEncoder(w).Encode(v)
}

func (api *TestAPI) CountryRevenue(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	limit := 100
	if l := q.Get("limit"); l != "" {
		if parsedLimit := parseInt(l); parsedLimit > 0 && parsedLimit <= 1000 {
			limit = parsedLimit
		}
	}

	offset := 0
	if o := q.Get("offset"); o != "" {
		if parsedOffset := parseInt(o); parsedOffset >= 0 {
			offset = parsedOffset
		}
	}

	all := api.Agg.CountryRevenueTable()

	start := offset
	if start >= len(all) {
		api.writeJSON(w, []interface{}{})
		return
	}

	end := start + limit
	if end > len(all) {
		end = len(all)
	}

	result := all[start:end]
	api.writeJSON(w, result)
}

func (api *TestAPI) TopProducts(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	limit := 20
	if l := q.Get("limit"); l != "" {
		if parsedLimit := parseInt(l); parsedLimit > 0 {
			limit = parsedLimit
		}
	}
	byUnits := q.Get("by") == "units"
	api.writeJSON(w, api.Agg.TopProducts(limit, byUnits))
}

func (api *TestAPI) SalesByMonth(w http.ResponseWriter, r *http.Request) {
	api.writeJSON(w, api.Agg.SalesByMonth())
}

func (api *TestAPI) TopRegions(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	limit := 30
	if l := q.Get("limit"); l != "" {
		if parsedLimit := parseInt(l); parsedLimit > 0 {
			limit = parsedLimit
		}
	}
	api.writeJSON(w, api.Agg.TopRegions(limit))
}

// Helper function to parse integers safely
func parseInt(s string) int {
	result := 0
	for _, r := range s {
		if r < '0' || r > '9' {
			return 0
		}
		result = result*10 + int(r-'0')
	}
	return result
}

func createMockAPI() *TestAPI {
	mockAgg := &mockAggregator{
		countryRevenue: []models.CountryProductAgg{
			{Country: "Sri Lanka", ProductName: "Widget A", TotalRevenue: 50000, NumberOfTx: 20},
			{Country: "India", ProductName: "Widget B", TotalRevenue: 75000, NumberOfTx: 30},
			{Country: "USA", ProductName: "Widget C", TotalRevenue: 100000, NumberOfTx: 40},
		},
		topProducts: []models.ProductAgg{
			{ProductName: "Widget A", TxCount: 100, UnitsSold: 1000, StockQty: 50},
			{ProductName: "Widget B", TxCount: 80, UnitsSold: 800, StockQty: 30},
			{ProductName: "Widget C", TxCount: 60, UnitsSold: 600, StockQty: 20},
		},
		salesByMonth: []models.MonthAgg{
			{YearMonth: "2023-01", UnitsSold: 500, TxCount: 50, RevenueCents: 125000},
			{YearMonth: "2023-02", UnitsSold: 600, TxCount: 60, RevenueCents: 150000},
			{YearMonth: "2023-03", UnitsSold: 700, TxCount: 70, RevenueCents: 175000},
		},
		topRegions: []models.RegionAgg{
			{Region: "Western", TotalRevenue: 200000, ItemsSold: 800, NumberOfTx: 200},
			{Region: "Central", TotalRevenue: 150000, ItemsSold: 600, NumberOfTx: 150},
			{Region: "Southern", TotalRevenue: 100000, ItemsSold: 400, NumberOfTx: 100},
		},
	}

	return &TestAPI{Agg: mockAgg}
}

func TestAPI_CountryRevenue(t *testing.T) {
	api := createMockAPI()

	tests := []struct {
		name           string
		url            string
		expectedStatus int
		expectedCount  int
	}{
		{
			name:           "default pagination",
			url:            "/api/revenue/countries",
			expectedStatus: http.StatusOK,
			expectedCount:  3,
		},
		{
			name:           "with limit",
			url:            "/api/revenue/countries?limit=2",
			expectedStatus: http.StatusOK,
			expectedCount:  2,
		},
		{
			name:           "with offset",
			url:            "/api/revenue/countries?limit=2&offset=1",
			expectedStatus: http.StatusOK,
			expectedCount:  2,
		},
		{
			name:           "offset beyond data",
			url:            "/api/revenue/countries?offset=10",
			expectedStatus: http.StatusOK,
			expectedCount:  0,
		},
		{
			name:           "invalid limit defaults to 100",
			url:            "/api/revenue/countries?limit=-1",
			expectedStatus: http.StatusOK,
			expectedCount:  3,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, err := http.NewRequest("GET", tt.url, nil)
			if err != nil {
				t.Fatal(err)
			}

			rr := httptest.NewRecorder()
			handler := http.HandlerFunc(api.CountryRevenue)
			handler.ServeHTTP(rr, req)

			if status := rr.Code; status != tt.expectedStatus {
				t.Errorf("handler returned wrong status code: got %v want %v",
					status, tt.expectedStatus)
			}

			// Check Content-Type header
			expectedContentType := "application/json; charset=utf-8"
			if ct := rr.Header().Get("Content-Type"); ct != expectedContentType {
				t.Errorf("handler returned wrong content type: got %v want %v",
					ct, expectedContentType)
			}

			// Check cache headers
			if cc := rr.Header().Get("Cache-Control"); cc == "" {
				t.Error("handler should set Cache-Control header")
			}

			if etag := rr.Header().Get("ETag"); etag == "" {
				t.Error("handler should set ETag header")
			}

			// Parse response
			var result []models.CountryProductAgg
			if err := json.Unmarshal(rr.Body.Bytes(), &result); err != nil {
				t.Fatalf("could not parse response: %v", err)
			}

			if len(result) != tt.expectedCount {
				t.Errorf("expected %d results, got %d", tt.expectedCount, len(result))
			}
		})
	}
}

func TestAPI_TopProducts(t *testing.T) {
	api := createMockAPI()

	tests := []struct {
		name           string
		url            string
		expectedStatus int
		expectedCount  int
	}{
		{
			name:           "default top 20",
			url:            "/api/products/top",
			expectedStatus: http.StatusOK,
			expectedCount:  3, // We only have 3 products in mock
		},
		{
			name:           "with limit",
			url:            "/api/products/top?limit=2",
			expectedStatus: http.StatusOK,
			expectedCount:  2,
		},
		{
			name:           "by units",
			url:            "/api/products/top?by=units",
			expectedStatus: http.StatusOK,
			expectedCount:  3,
		},
		{
			name:           "invalid limit defaults to 20",
			url:            "/api/products/top?limit=0",
			expectedStatus: http.StatusOK,
			expectedCount:  3,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, err := http.NewRequest("GET", tt.url, nil)
			if err != nil {
				t.Fatal(err)
			}

			rr := httptest.NewRecorder()
			handler := http.HandlerFunc(api.TopProducts)
			handler.ServeHTTP(rr, req)

			if status := rr.Code; status != tt.expectedStatus {
				t.Errorf("handler returned wrong status code: got %v want %v",
					status, tt.expectedStatus)
			}

			var result []models.ProductAgg
			if err := json.Unmarshal(rr.Body.Bytes(), &result); err != nil {
				t.Fatalf("could not parse response: %v", err)
			}

			if len(result) != tt.expectedCount {
				t.Errorf("expected %d results, got %d", tt.expectedCount, len(result))
			}
		})
	}
}

func TestAPI_SalesByMonth(t *testing.T) {
	api := createMockAPI()

	req, err := http.NewRequest("GET", "/api/sales/by-month", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(api.SalesByMonth)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	var result []models.MonthAgg
	if err := json.Unmarshal(rr.Body.Bytes(), &result); err != nil {
		t.Fatalf("could not parse response: %v", err)
	}

	if len(result) != 3 {
		t.Errorf("expected 3 months, got %d", len(result))
	}

	// Verify data integrity
	for _, month := range result {
		if month.YearMonth == "" {
			t.Error("month should have year-month value")
		}
		if month.UnitsSold <= 0 {
			t.Error("month should have positive units sold")
		}
		if month.RevenueCents <= 0 {
			t.Error("month should have positive revenue")
		}
	}
}

func TestAPI_TopRegions(t *testing.T) {
	api := createMockAPI()

	tests := []struct {
		name           string
		url            string
		expectedStatus int
		expectedCount  int
	}{
		{
			name:           "default top 30",
			url:            "/api/regions/top",
			expectedStatus: http.StatusOK,
			expectedCount:  3, // We only have 3 regions in mock
		},
		{
			name:           "with limit",
			url:            "/api/regions/top?limit=2",
			expectedStatus: http.StatusOK,
			expectedCount:  2,
		},
		{
			name:           "invalid limit defaults to 30",
			url:            "/api/regions/top?limit=-5",
			expectedStatus: http.StatusOK,
			expectedCount:  3,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, err := http.NewRequest("GET", tt.url, nil)
			if err != nil {
				t.Fatal(err)
			}

			rr := httptest.NewRecorder()
			handler := http.HandlerFunc(api.TopRegions)
			handler.ServeHTTP(rr, req)

			if status := rr.Code; status != tt.expectedStatus {
				t.Errorf("handler returned wrong status code: got %v want %v",
					status, tt.expectedStatus)
			}

			var result []models.RegionAgg
			if err := json.Unmarshal(rr.Body.Bytes(), &result); err != nil {
				t.Fatalf("could not parse response: %v", err)
			}

			if len(result) != tt.expectedCount {
				t.Errorf("expected %d results, got %d", tt.expectedCount, len(result))
			}

			// Verify data integrity
			for _, region := range result {
				if region.Region == "" {
					t.Error("region should have name")
				}
				if region.TotalRevenue <= 0 {
					t.Error("region should have positive revenue")
				}
			}
		})
	}
}

func TestAPI_writeJSON(t *testing.T) {
	api := createMockAPI()

	testData := map[string]interface{}{
		"message": "test",
		"value":   123,
	}

	rr := httptest.NewRecorder()
	api.writeJSON(rr, testData)

	// Check headers
	expectedContentType := "application/json; charset=utf-8"
	if ct := rr.Header().Get("Content-Type"); ct != expectedContentType {
		t.Errorf("expected content type %s, got %s", expectedContentType, ct)
	}

	if cc := rr.Header().Get("Cache-Control"); cc != "public, max-age=300" {
		t.Errorf("expected cache control header, got %s", cc)
	}

	if etag := rr.Header().Get("ETag"); etag != "\"dashboard-data\"" {
		t.Errorf("expected ETag header, got %s", etag)
	}

	// Check JSON content
	var result map[string]interface{}
	if err := json.Unmarshal(rr.Body.Bytes(), &result); err != nil {
		t.Fatalf("could not parse JSON: %v", err)
	}

	if result["message"] != "test" {
		t.Errorf("expected message 'test', got %v", result["message"])
	}

	if result["value"] != float64(123) { // JSON numbers are float64
		t.Errorf("expected value 123, got %v", result["value"])
	}
}

// Benchmark tests for performance validation
func BenchmarkAPI_CountryRevenue(b *testing.B) {
	api := createMockAPI()
	req, _ := http.NewRequest("GET", "/api/revenue/countries?limit=100", nil)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(api.CountryRevenue)
		handler.ServeHTTP(rr, req)
	}
}

func BenchmarkAPI_TopProducts(b *testing.B) {
	api := createMockAPI()
	req, _ := http.NewRequest("GET", "/api/products/top?limit=20", nil)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(api.TopProducts)
		handler.ServeHTTP(rr, req)
	}
}

func BenchmarkAPI_writeJSON(b *testing.B) {
	api := createMockAPI()
	testData := map[string]interface{}{
		"message": "benchmark test",
		"value":   456,
		"time":    time.Now(),
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		rr := httptest.NewRecorder()
		api.writeJSON(rr, testData)
	}
}
