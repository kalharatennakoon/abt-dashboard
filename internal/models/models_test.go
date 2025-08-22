package models

import (
	"testing"
	"time"
)

func TestTransaction_Validation(t *testing.T) {
	tests := []struct {
		name        string
		transaction Transaction
		shouldError bool
	}{
		{
			name: "valid transaction",
			transaction: Transaction{
				ID:             "tx-001",
				Country:        "Sri Lanka",
				Region:         "Western",
				ProductName:    "Widget A",
				UnitPriceCents: 2500, // $25.00 in cents
				Quantity:       10,
				TxTime:         time.Date(2023, 1, 15, 0, 0, 0, 0, time.UTC),
			},
			shouldError: false,
		},
		{
			name: "invalid transaction - empty ID",
			transaction: Transaction{
				ID:             "",
				Country:        "Sri Lanka",
				Region:         "Western",
				ProductName:    "Widget A",
				UnitPriceCents: 2500,
				Quantity:       10,
				TxTime:         time.Date(2023, 1, 15, 0, 0, 0, 0, time.UTC),
			},
			shouldError: true,
		},
		{
			name: "invalid transaction - negative quantity",
			transaction: Transaction{
				ID:             "tx-002",
				Country:        "Sri Lanka",
				Region:         "Western",
				ProductName:    "Widget A",
				UnitPriceCents: 2500,
				Quantity:       -5,
				TxTime:         time.Date(2023, 1, 15, 0, 0, 0, 0, time.UTC),
			},
			shouldError: true,
		},
		{
			name: "invalid transaction - negative unit price",
			transaction: Transaction{
				ID:             "tx-003",
				Country:        "Sri Lanka",
				Region:         "Western",
				ProductName:    "Widget A",
				UnitPriceCents: -2500,
				Quantity:       10,
				TxTime:         time.Date(2023, 1, 15, 0, 0, 0, 0, time.UTC),
			},
			shouldError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Validate transaction fields
			isValid := tt.transaction.ID != "" &&
				tt.transaction.Country != "" &&
				tt.transaction.ProductName != "" &&
				tt.transaction.UnitPriceCents >= 0 &&
				tt.transaction.Quantity >= 0 &&
				!tt.transaction.TxTime.IsZero()

			if tt.shouldError && isValid {
				t.Errorf("expected invalid transaction but got valid")
			}
			if !tt.shouldError && !isValid {
				t.Errorf("expected valid transaction but got invalid")
			}
		})
	}
}

func TestCountryProductAgg_Calculation(t *testing.T) {
	agg := CountryProductAgg{
		Country:      "Sri Lanka",
		ProductName:  "Widget A",
		TotalRevenue: 50000, // $500.00 in cents
		NumberOfTx:   20,
	}

	// Test revenue per transaction
	avgRevenue := agg.TotalRevenue / agg.NumberOfTx
	expectedAvg := int64(2500) // $25.00 per transaction

	if avgRevenue != expectedAvg {
		t.Errorf("expected average revenue %d, got %d", expectedAvg, avgRevenue)
	}
}

func TestProductAgg_StockValidation(t *testing.T) {
	tests := []struct {
		name    string
		product ProductAgg
		isValid bool
	}{
		{
			name: "valid product with stock",
			product: ProductAgg{
				ProductName: "Widget A",
				TxCount:     50,
				UnitsSold:   500,
				StockQty:    100,
			},
			isValid: true,
		},
		{
			name: "product with zero stock",
			product: ProductAgg{
				ProductName: "Widget B",
				TxCount:     30,
				UnitsSold:   300,
				StockQty:    0,
			},
			isValid: true, // Zero stock is valid (out of stock)
		},
		{
			name: "product with negative stock",
			product: ProductAgg{
				ProductName: "Widget C",
				TxCount:     10,
				UnitsSold:   100,
				StockQty:    -5,
			},
			isValid: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			isValid := tt.product.ProductName != "" &&
				tt.product.TxCount >= 0 &&
				tt.product.UnitsSold >= 0 &&
				tt.product.StockQty >= 0

			if isValid != tt.isValid {
				t.Errorf("expected validity %v, got %v", tt.isValid, isValid)
			}
		})
	}
}

func TestMonthAgg_TimeComparison(t *testing.T) {
	jan := MonthAgg{YearMonth: "2023-01", UnitsSold: 100, TxCount: 10, RevenueCents: 25000}
	feb := MonthAgg{YearMonth: "2023-02", UnitsSold: 150, TxCount: 15, RevenueCents: 37500}

	// Test that months can be compared
	if jan.YearMonth >= feb.YearMonth {
		t.Error("expected January to be less than February")
	}

	// Test revenue calculation
	expectedJanAvg := jan.RevenueCents / jan.TxCount
	if expectedJanAvg != 2500 {
		t.Errorf("expected January average revenue 2500, got %d", expectedJanAvg)
	}
}

func TestRegionAgg_PerformanceMetrics(t *testing.T) {
	region := RegionAgg{
		Region:       "Western",
		TotalRevenue: 1000000, // $10,000.00 in cents
		ItemsSold:    400,
		NumberOfTx:   100,
	}

	// Test average revenue per item
	avgRevenuePerItem := region.TotalRevenue / region.ItemsSold
	expectedAvgPerItem := int64(2500) // $25.00 per item

	if avgRevenuePerItem != expectedAvgPerItem {
		t.Errorf("expected average revenue per item %d, got %d", expectedAvgPerItem, avgRevenuePerItem)
	}

	// Test average items per transaction
	avgItemsPerTx := region.ItemsSold / region.NumberOfTx
	expectedAvgItems := int64(4) // 4 items per transaction

	if avgItemsPerTx != expectedAvgItems {
		t.Errorf("expected average items per transaction %d, got %d", expectedAvgItems, avgItemsPerTx)
	}
}

func TestInsight_Creation(t *testing.T) {
	createdAt := time.Now()
	expiresAt := createdAt.Add(24 * time.Hour)

	insight := Insight{
		ID:          "insight-001",
		Type:        "trend",
		Title:       "Sales Trending Up",
		Description: "Sales have increased by 15% this month",
		Severity:    "medium",
		Confidence:  0.85,
		Data: map[string]interface{}{
			"trend_percentage": 15.0,
			"period":           "monthly",
		},
		CreatedAt: createdAt,
		ExpiresAt: &expiresAt,
	}

	// Test insight validation
	if insight.ID == "" {
		t.Error("insight ID should not be empty")
	}
	if insight.Confidence < 0 || insight.Confidence > 1 {
		t.Error("insight confidence should be between 0 and 1")
	}
	if insight.CreatedAt.IsZero() {
		t.Error("insight creation time should not be zero")
	}
	if insight.ExpiresAt.Before(insight.CreatedAt) {
		t.Error("insight expiration should be after creation time")
	}
}

func TestInventory_StockManagement(t *testing.T) {
	inventory := Inventory{
		ProductName: "Widget A",
		StockQty:    50,
	}

	// Test stock reduction
	soldQuantity := int64(10)
	inventory.StockQty -= soldQuantity

	expectedStock := int64(40)
	if inventory.StockQty != expectedStock {
		t.Errorf("expected stock %d after sale, got %d", expectedStock, inventory.StockQty)
	}

	// Test low stock detection
	isLowStock := inventory.StockQty < 10
	if isLowStock {
		t.Log("Stock is running low for", inventory.ProductName)
	}
}

// Benchmark tests for performance validation
func BenchmarkTransaction_Creation(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = Transaction{
			ID:             "tx-001",
			Country:        "Sri Lanka",
			Region:         "Western",
			ProductName:    "Widget A",
			UnitPriceCents: 2500,
			Quantity:       10,
			TxTime:         time.Now(),
		}
	}
}

func BenchmarkCountryProductAgg_Calculation(b *testing.B) {
	agg := CountryProductAgg{
		Country:      "Sri Lanka",
		ProductName:  "Widget A",
		TotalRevenue: 50000,
		NumberOfTx:   20,
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = agg.TotalRevenue / agg.NumberOfTx
	}
}
