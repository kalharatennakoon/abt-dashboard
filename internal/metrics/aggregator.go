package metrics

import (
    "sort"
    "strings"
    "sync"
    "time"

    "abt-dashboard/internal/models"
)

// Aggregator holds in-memory aggregations for analytics.
type Aggregator struct {
    countryProduct map[string]map[string]*models.CountryProductAgg // country → product → agg
    productAgg     map[string]*models.ProductAgg                   // product → agg
    monthAgg       map[string]*models.MonthAgg                     // YYYY-MM → agg
    regionAgg      map[string]*models.RegionAgg                    // region → agg

    mu sync.RWMutex
}

// NewAggregator creates a new, empty aggregator.
func NewAggregator() *Aggregator {
    return &Aggregator{
        countryProduct: make(map[string]map[string]*models.CountryProductAgg),
        productAgg:     make(map[string]*models.ProductAgg),
        monthAgg:       make(map[string]*models.MonthAgg),
        regionAgg:      make(map[string]*models.RegionAgg),
    }
}

// Ingest loads transactions and inventory into the aggregator.
func (a *Aggregator) Ingest(trans []models.Transaction, inv map[string]models.Inventory) {
    a.mu.Lock()
    defer a.mu.Unlock()

    for _, t := range trans {
        // Country-Product aggregation
        if _, ok := a.countryProduct[t.Country]; !ok {
            a.countryProduct[t.Country] = make(map[string]*models.CountryProductAgg)
        }
        cp := a.countryProduct[t.Country][t.ProductName]
        if cp == nil {
            cp = &models.CountryProductAgg{
                Country:     t.Country,
                ProductName: t.ProductName,
            }
            a.countryProduct[t.Country][t.ProductName] = cp
        }
        cp.TotalRevenue += t.UnitPriceCents * t.Quantity
        cp.NumberOfTx++

        // Product-level aggregation
        pa := a.productAgg[t.ProductName]
        if pa == nil {
            pa = &models.ProductAgg{ProductName: t.ProductName}
            a.productAgg[t.ProductName] = pa
        }
        pa.TxCount++
        pa.UnitsSold += t.Quantity

        // Month aggregation
        ym := t.TxTime.Format("2006-01")
        ma := a.monthAgg[ym]
        if ma == nil {
            ma = &models.MonthAgg{YearMonth: ym}
            a.monthAgg[ym] = ma
        }
        ma.UnitsSold += t.Quantity
        ma.TxCount++
        ma.RevenueCents += t.UnitPriceCents * t.Quantity

        // Region aggregation
        ra := a.regionAgg[t.Region]
        if ra == nil {
            ra = &models.RegionAgg{Region: t.Region}
            a.regionAgg[t.Region] = ra
        }
        ra.TotalRevenue += t.UnitPriceCents * t.Quantity
        ra.ItemsSold += t.Quantity
        ra.NumberOfTx++
    }

    // Merge inventory stocks into product aggregates
    for name, p := range a.productAgg {
        if invRow, ok := inv[name]; ok {
            p.StockQty = invRow.StockQty
        }
    }
}

// CountryRevenueTable returns all country-product aggregates sorted by revenue desc.
func (a *Aggregator) CountryRevenueTable() []models.CountryProductAgg {
    a.mu.RLock()
    defer a.mu.RUnlock()

    out := make([]models.CountryProductAgg, 0)
    for _, pm := range a.countryProduct {
        for _, agg := range pm {
            out = append(out, *agg)
        }
    }

    sort.Slice(out, func(i, j int) bool {
        if out[i].TotalRevenue == out[j].TotalRevenue {
            // tie-breaker: country then product
            ci := strings.Compare(out[i].Country, out[j].Country)
            if ci == 0 {
                return out[i].ProductName < out[j].ProductName
            }
            return ci < 0
        }
        return out[i].TotalRevenue > out[j].TotalRevenue
    })
    return out
}

// TopProducts returns products sorted by tx count (or units if byUnits=true).
func (a *Aggregator) TopProducts(limit int, byUnits bool) []models.ProductAgg {
    a.mu.RLock()
    defer a.mu.RUnlock()

    out := make([]models.ProductAgg, 0, len(a.productAgg))
    for _, v := range a.productAgg {
        out = append(out, *v)
    }

    sort.Slice(out, func(i, j int) bool {
        if byUnits {
            if out[i].UnitsSold == out[j].UnitsSold {
                return out[i].ProductName < out[j].ProductName
            }
            return out[i].UnitsSold > out[j].UnitsSold
        }
        if out[i].TxCount == out[j].TxCount {
            return out[i].ProductName < out[j].ProductName
        }
        return out[i].TxCount > out[j].TxCount
    })

    if limit > 0 && len(out) > limit {
        out = out[:limit]
    }
    return out
}

// SalesByMonth returns monthly aggregates sorted chronologically.
func (a *Aggregator) SalesByMonth() []models.MonthAgg {
    a.mu.RLock()
    defer a.mu.RUnlock()

    out := make([]models.MonthAgg, 0, len(a.monthAgg))
    for _, v := range a.monthAgg {
        out = append(out, *v)
    }

    sort.Slice(out, func(i, j int) bool {
        ti, _ := time.Parse("2006-01", out[i].YearMonth)
        tj, _ := time.Parse("2006-01", out[j].YearMonth)
        return ti.Before(tj)
    })
    return out
}

// TopRegions returns regions sorted by revenue desc.
func (a *Aggregator) TopRegions(limit int) []models.RegionAgg {
    a.mu.RLock()
    defer a.mu.RUnlock()

    out := make([]models.RegionAgg, 0, len(a.regionAgg))
    for _, v := range a.regionAgg {
        out = append(out, *v)
    }

    sort.Slice(out, func(i, j int) bool {
        if out[i].TotalRevenue == out[j].TotalRevenue {
            return out[i].Region < out[j].Region
        }
        return out[i].TotalRevenue > out[j].TotalRevenue
    })

    if limit > 0 && len(out) > limit {
        out = out[:limit]
    }
    return out
}
