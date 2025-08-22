package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"abt-dashboard/internal/metrics"
)

// API wraps our aggregator so it can serve JSON endpoints.
type API struct {
	Agg *metrics.Aggregator
}

func (api *API) writeJSON(w http.ResponseWriter, v interface{}) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("Cache-Control", "public, max-age=300") // 5 minutes browser caching for better performance
	w.Header().Set("ETag", "\"dashboard-data\"")           // Simple ETag for cache validation
	json.NewEncoder(w).Encode(v)
}

// GET /api/revenue/countries?limit=100&offset=0
func (api *API) CountryRevenue(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	limit, _ := strconv.Atoi(q.Get("limit"))
	if limit <= 0 || limit > 1000 {
		limit = 100 // Default limit
	}
	offset, _ := strconv.Atoi(q.Get("offset"))
	if offset < 0 {
		offset = 0
	}

	all := api.Agg.CountryRevenueTable()

	// Apply pagination
	start := offset
	if start >= len(all) {
		api.writeJSON(w, []interface{}{}) // Empty result
		return
	}

	end := start + limit
	if end > len(all) {
		end = len(all)
	}

	result := all[start:end]
	api.writeJSON(w, result)
}

// GET /api/products/top?limit=20&by=units
func (api *API) TopProducts(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	limit, _ := strconv.Atoi(q.Get("limit"))
	if limit <= 0 {
		limit = 20
	}
	byUnits := q.Get("by") == "units"
	api.writeJSON(w, api.Agg.TopProducts(limit, byUnits))
}

// GET /api/sales/by-month
func (api *API) SalesByMonth(w http.ResponseWriter, r *http.Request) {
	api.writeJSON(w, api.Agg.SalesByMonth())
}

// GET /api/regions/top?limit=30
func (api *API) TopRegions(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	limit, _ := strconv.Atoi(q.Get("limit"))
	if limit <= 0 {
		limit = 30
	}
	api.writeJSON(w, api.Agg.TopRegions(limit))
}
