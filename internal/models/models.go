package models

import "time"

// Transaction represents a single sale record.
type Transaction struct {
	ID             string    // unique transaction ID
	Country        string    // e.g., "Sri Lanka"
	Region         string    // e.g., "Western"
	ProductName    string    // e.g., "Widget A"
	UnitPriceCents int64     // price per unit, stored in cents to avoid float issues
	Quantity       int64     // number of units sold
	TxTime         time.Time // transaction timestamp
}

// Inventory represents available stock for a product.
type Inventory struct {
	ProductName string
	StockQty    int64
}

// Aggregated view: revenue by country/product
type CountryProductAgg struct {
	Country      string `json:"country"`
	ProductName  string `json:"product_name"`
	TotalRevenue int64  `json:"total_revenue_cents"`
	NumberOfTx   int64  `json:"number_of_transactions"`
}

// Aggregated view: product popularity
type ProductAgg struct {
	ProductName string `json:"product_name"`
	TxCount     int64  `json:"tx_count"`
	UnitsSold   int64  `json:"units_sold"`
	StockQty    int64  `json:"stock_qty"`
}

// Aggregated view: monthly sales trends
type MonthAgg struct {
	YearMonth    string `json:"year_month"` // format: YYYY-MM
	UnitsSold    int64  `json:"units_sold"`
	TxCount      int64  `json:"tx_count"`
	RevenueCents int64  `json:"revenue_cents"`
}

// Aggregated view: regional performance
type RegionAgg struct {
	Region       string `json:"region"`
	TotalRevenue int64  `json:"total_revenue_cents"`
	ItemsSold    int64  `json:"items_sold"`
	NumberOfTx   int64  `json:"number_of_transactions"`
}

// Insight represents a business insight generated from data analysis
type Insight struct {
	ID          string                 `json:"id"`
	Type        string                 `json:"type"` // e.g., "trend", "anomaly", "recommendation"
	Title       string                 `json:"title"`
	Description string                 `json:"description"`
	Severity    string                 `json:"severity"`   // "low", "medium", "high", "critical"
	Confidence  float64                `json:"confidence"` // 0.0 to 1.0
	Data        map[string]interface{} `json:"data"`       // Supporting data for the insight
	CreatedAt   time.Time              `json:"created_at"`
	ExpiresAt   *time.Time             `json:"expires_at,omitempty"`
}
