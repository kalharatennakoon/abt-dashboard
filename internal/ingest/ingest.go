package ingest

import (
	"bufio"
	"encoding/csv"
	"errors"
	"io"
	"strconv"
	"strings"
	"time"

	"abt-dashboard/internal/models"
)

// ParseTransactionsCSV reads transactions.csv into a slice of Transaction.
//
// Expected CSV headers:
//
//	id,country,region,product,unit_price_cents,quantity,tx_time
//
// tx_time supports RFC3339 (e.g. 2024-03-15T12:34:56Z)
// or YYYY-MM-DD (e.g. 2024-03-15).
func ParseTransactionsCSV(r io.Reader) ([]models.Transaction, error) {
	cr := csv.NewReader(bufio.NewReader(r))
	cr.TrimLeadingSpace = true

	// Read header
	header, err := cr.Read()
	if err != nil {
		return nil, err
	}

	// Map headers to indices
	idx := map[string]int{}
	for i, h := range header {
		idx[strings.ToLower(strings.TrimSpace(h))] = i
	}

	required := []string{"transaction_id", "transaction_date", "country", "region", "product_name", "price", "quantity"}
	for _, req := range required {
		if _, ok := idx[req]; !ok {
			return nil, errors.New("missing required column: " + req)
		}
	}

	var out []models.Transaction
	for {
		rec, err := cr.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}

		// Parse price as float and convert to cents
		priceFloat, _ := strconv.ParseFloat(rec[idx["price"]], 64)
		up := int64(priceFloat * 100) // Convert dollars to cents
		qty, _ := strconv.ParseInt(rec[idx["quantity"]], 10, 64)

		// Parse date (try RFC3339 then YYYY-MM-DD)
		tstr := rec[idx["transaction_date"]]
		var tt time.Time
		tt, err = time.Parse(time.RFC3339, tstr)
		if err != nil {
			tt, err = time.Parse("2006-01-02", tstr)
			if err != nil {
				continue // skip bad rows
			}
		}

		out = append(out, models.Transaction{
			ID:             rec[idx["transaction_id"]],
			Country:        rec[idx["country"]],
			Region:         rec[idx["region"]],
			ProductName:    rec[idx["product_name"]],
			UnitPriceCents: up,
			Quantity:       qty,
			TxTime:         tt,
		})
	}
	return out, nil
}

// ParseInventoryCSV reads inventory.csv into a map of ProductName → Inventory.
//
// Expected CSV headers:
//
//	product,stock_qty
func ParseInventoryCSV(r io.Reader) (map[string]models.Inventory, error) {
	cr := csv.NewReader(bufio.NewReader(r))
	cr.TrimLeadingSpace = true

	header, err := cr.Read()
	if err != nil {
		return nil, err
	}

	idx := map[string]int{}
	for i, h := range header {
		idx[strings.ToLower(strings.TrimSpace(h))] = i
	}

	required := []string{"product_name", "stock_quantity"}
	for _, req := range required {
		if _, ok := idx[req]; !ok {
			return nil, errors.New("missing required column: " + req)
		}
	}

	res := make(map[string]models.Inventory, 128)
	for {
		rec, err := cr.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}

		qty, _ := strconv.ParseInt(rec[idx["stock_quantity"]], 10, 64)
		name := rec[idx["product_name"]]
		res[name] = models.Inventory{ProductName: name, StockQty: qty}
	}
	return res, nil
}
