package transform

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"strconv"
	"strings"
	"time"

	"abt-dashboard/internal/models"

	"gopkg.in/yaml.v2"
)

// DataFormat represents supported data formats
type DataFormat string

const (
	FormatCSV  DataFormat = "csv"
	FormatJSON DataFormat = "json"
	FormatYAML DataFormat = "yaml"
	FormatTSV  DataFormat = "tsv"
	FormatXML  DataFormat = "xml"
)

// FormatConverter handles conversion between different data formats
type FormatConverter struct {
	config TransformConfig
}

// NewFormatConverter creates a new format converter
func NewFormatConverter(config TransformConfig) *FormatConverter {
	return &FormatConverter{config: config}
}

// DetectFormat attempts to detect the format of input data
func (fc *FormatConverter) DetectFormat(data []byte) DataFormat {
	dataStr := strings.TrimSpace(string(data))

	// Check for JSON
	if (strings.HasPrefix(dataStr, "{") && strings.HasSuffix(dataStr, "}")) ||
		(strings.HasPrefix(dataStr, "[") && strings.HasSuffix(dataStr, "]")) {
		return FormatJSON
	}

	// Check for XML
	if strings.HasPrefix(dataStr, "<") && strings.Contains(dataStr, ">") {
		return FormatXML
	}

	// Check for YAML (common indicators)
	if strings.Contains(dataStr, ":") &&
		(strings.Contains(dataStr, "\n-") || strings.Contains(dataStr, "---")) {
		return FormatYAML
	}

	// Check for TSV (tab-separated)
	if strings.Contains(dataStr, "\t") {
		return FormatTSV
	}

	// Default to CSV
	return FormatCSV
}

// ConvertToTransactions converts data from various formats to Transaction slice
func (fc *FormatConverter) ConvertToTransactions(reader io.Reader, format DataFormat) ([]models.Transaction, error) {
	switch format {
	case FormatCSV:
		return fc.parseCSV(reader, ',')
	case FormatTSV:
		return fc.parseCSV(reader, '\t')
	case FormatJSON:
		return fc.parseJSON(reader)
	case FormatYAML:
		return fc.parseYAML(reader)
	default:
		return nil, fmt.Errorf("unsupported format: %s", format)
	}
}

// parseCSV handles CSV and TSV parsing with flexible column mapping
func (fc *FormatConverter) parseCSV(reader io.Reader, delimiter rune) ([]models.Transaction, error) {
	csvReader := csv.NewReader(reader)
	csvReader.Comma = delimiter
	csvReader.TrimLeadingSpace = true
	csvReader.FieldsPerRecord = -1 // Allow variable number of fields

	// Read header
	header, err := csvReader.Read()
	if err != nil {
		return nil, fmt.Errorf("failed to read header: %w", err)
	}

	// Create flexible column mapping
	columnMap := fc.createColumnMapping(header)

	var transactions []models.Transaction
	lineNumber := 1

	for {
		record, err := csvReader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("error reading line %d: %w", lineNumber, err)
		}

		transaction, err := fc.parseRecordToTransaction(record, columnMap)
		if err != nil {
			// Log error but continue processing
			fmt.Printf("Warning: Failed to parse line %d: %v\n", lineNumber, err)
			continue
		}

		transactions = append(transactions, *transaction)
		lineNumber++
	}

	return transactions, nil
}

// parseJSON handles JSON format parsing
func (fc *FormatConverter) parseJSON(reader io.Reader) ([]models.Transaction, error) {
	var data interface{}
	decoder := json.NewDecoder(reader)

	if err := decoder.Decode(&data); err != nil {
		return nil, fmt.Errorf("failed to parse JSON: %w", err)
	}

	return fc.extractTransactionsFromData(data)
}

// parseYAML handles YAML format parsing
func (fc *FormatConverter) parseYAML(reader io.Reader) ([]models.Transaction, error) {
	var data interface{}
	decoder := yaml.NewDecoder(reader)

	if err := decoder.Decode(&data); err != nil {
		return nil, fmt.Errorf("failed to parse YAML: %w", err)
	}

	return fc.extractTransactionsFromData(data)
}

// createColumnMapping creates flexible column mapping for various header formats
func (fc *FormatConverter) createColumnMapping(header []string) map[string]int {
	columnMap := make(map[string]int)

	// Define possible column name variations
	columnVariations := map[string][]string{
		"transaction_id":   {"transaction_id", "id", "trans_id", "txn_id", "transaction", "order_id"},
		"country":          {"country", "nation", "country_name", "country_code"},
		"region":           {"region", "state", "province", "area", "zone", "territory"},
		"product_name":     {"product_name", "product", "item", "item_name", "product_title"},
		"price":            {"price", "unit_price", "cost", "amount", "unit_cost", "price_per_unit"},
		"quantity":         {"quantity", "qty", "amount", "count", "units", "number"},
		"transaction_date": {"transaction_date", "date", "timestamp", "time", "tx_date", "order_date"},
	}

	// Map header columns to standard field names
	for i, col := range header {
		normalizedCol := strings.ToLower(strings.TrimSpace(col))

		// Direct mapping
		columnMap[normalizedCol] = i

		// Find matching standard field
		for standardField, variations := range columnVariations {
			for _, variation := range variations {
				if normalizedCol == variation {
					columnMap[standardField] = i
					break
				}
			}
		}
	}

	return columnMap
}

// parseRecordToTransaction converts a record to Transaction with flexible field mapping
func (fc *FormatConverter) parseRecordToTransaction(record []string, columnMap map[string]int) (*models.Transaction, error) {
	tx := &models.Transaction{}

	// Helper function to safely get field value
	getField := func(fieldName string) string {
		if idx, exists := columnMap[fieldName]; exists && idx < len(record) {
			return strings.TrimSpace(record[idx])
		}
		return ""
	}

	// Transaction ID
	tx.ID = getField("transaction_id")
	if tx.ID == "" {
		return nil, fmt.Errorf("transaction ID is required")
	}

	// Country
	tx.Country = getField("country")
	if tx.Country == "" {
		tx.Country = fc.config.DefaultCountry
	}

	// Region
	tx.Region = getField("region")
	if tx.Region == "" {
		tx.Region = fc.config.DefaultRegion
	}

	// Product Name
	tx.ProductName = getField("product_name")
	if tx.ProductName == "" {
		return nil, fmt.Errorf("product name is required")
	}

	// Price
	priceStr := getField("price")
	if priceStr == "" {
		return nil, fmt.Errorf("price is required")
	}

	price, err := fc.parseFlexiblePrice(priceStr)
	if err != nil {
		return nil, fmt.Errorf("invalid price '%s': %w", priceStr, err)
	}
	tx.UnitPriceCents = price

	// Quantity
	quantityStr := getField("quantity")
	if quantityStr == "" {
		quantityStr = "1" // Default quantity
	}

	quantity, err := strconv.ParseInt(quantityStr, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("invalid quantity '%s': %w", quantityStr, err)
	}
	tx.Quantity = quantity

	// Transaction Date
	dateStr := getField("transaction_date")
	if dateStr == "" {
		return nil, fmt.Errorf("transaction date is required")
	}

	date, err := fc.parseFlexibleDate(dateStr)
	if err != nil {
		return nil, fmt.Errorf("invalid date '%s': %w", dateStr, err)
	}
	tx.TxTime = date

	return tx, nil
}

// parseFlexiblePrice handles various price formats
func (fc *FormatConverter) parseFlexiblePrice(priceStr string) (int64, error) {
	// Remove common currency symbols and formatting
	cleanPrice := priceStr

	// Remove currency symbols
	currencySymbols := []string{"$", "€", "£", "¥", "₹", "₽", "₩", "₪", "₦", "₡", "₵", "₴", "₸", "₱", "₫", "₨"}
	for _, symbol := range currencySymbols {
		cleanPrice = strings.ReplaceAll(cleanPrice, symbol, "")
	}

	// Remove thousand separators
	cleanPrice = strings.ReplaceAll(cleanPrice, ",", "")
	cleanPrice = strings.ReplaceAll(cleanPrice, " ", "")

	// Handle percentage format (remove % sign)
	cleanPrice = strings.ReplaceAll(cleanPrice, "%", "")

	// Parse as float
	price, err := strconv.ParseFloat(cleanPrice, 64)
	if err != nil {
		return 0, err
	}

	// Convert to cents using configured multiplier
	return int64(price * fc.config.PriceMultiplier), nil
}

// parseFlexibleDate handles various date formats
func (fc *FormatConverter) parseFlexibleDate(dateStr string) (time.Time, error) {
	// Try configured date formats
	for _, format := range fc.config.DateFormats {
		if date, err := time.Parse(format, dateStr); err == nil {
			return date, nil
		}
	}

	// Try additional common formats
	additionalFormats := []string{
		"January 2, 2006",
		"Jan 2, 2006 15:04:05",
		"2006-01-02T15:04:05.000Z",
		"2006-01-02T15:04:05.000000Z",
		"Mon Jan 2 15:04:05 MST 2006",
		"Mon, 02 Jan 2006 15:04:05 MST",
	}

	for _, format := range additionalFormats {
		if date, err := time.Parse(format, dateStr); err == nil {
			return date, nil
		}
	}

	// Try parsing as Unix timestamp
	if timestamp, err := strconv.ParseInt(dateStr, 10, 64); err == nil {
		// Check if it's seconds or milliseconds
		if timestamp > 1e10 { // Likely milliseconds
			return time.Unix(timestamp/1000, (timestamp%1000)*1e6), nil
		} else { // Likely seconds
			return time.Unix(timestamp, 0), nil
		}
	}

	return time.Time{}, fmt.Errorf("unable to parse date format")
}

// extractTransactionsFromData extracts transactions from parsed JSON/YAML data
func (fc *FormatConverter) extractTransactionsFromData(data interface{}) ([]models.Transaction, error) {
	var transactions []models.Transaction

	switch v := data.(type) {
	case []interface{}:
		// Array of transaction objects
		for i, item := range v {
			if itemMap, ok := item.(map[interface{}]interface{}); ok {
				// Convert map[interface{}]interface{} to map[string]interface{}
				stringMap := make(map[string]interface{})
				for key, value := range itemMap {
					if keyStr, ok := key.(string); ok {
						stringMap[keyStr] = value
					}
				}
				tx, err := fc.mapToTransaction(stringMap)
				if err != nil {
					fmt.Printf("Warning: Failed to parse transaction %d: %v\n", i, err)
					continue
				}
				transactions = append(transactions, *tx)
			} else if itemMap, ok := item.(map[string]interface{}); ok {
				tx, err := fc.mapToTransaction(itemMap)
				if err != nil {
					fmt.Printf("Warning: Failed to parse transaction %d: %v\n", i, err)
					continue
				}
				transactions = append(transactions, *tx)
			}
		}
	case map[string]interface{}:
		// Single transaction object or object containing transactions
		if txArray, exists := v["transactions"]; exists {
			return fc.extractTransactionsFromData(txArray)
		} else if txArray, exists := v["data"]; exists {
			return fc.extractTransactionsFromData(txArray)
		} else {
			// Treat as single transaction
			tx, err := fc.mapToTransaction(v)
			if err != nil {
				return nil, err
			}
			transactions = append(transactions, *tx)
		}
	case map[interface{}]interface{}:
		// Convert to string map and try again
		stringMap := make(map[string]interface{})
		for key, value := range v {
			if keyStr, ok := key.(string); ok {
				stringMap[keyStr] = value
			}
		}
		return fc.extractTransactionsFromData(stringMap)
	default:
		return nil, fmt.Errorf("unsupported data structure type: %T", data)
	}

	return transactions, nil
}

// mapToTransaction converts a map to Transaction
func (fc *FormatConverter) mapToTransaction(data map[string]interface{}) (*models.Transaction, error) {
	tx := &models.Transaction{}

	// Map fields with flexible naming
	fieldMappings := map[string][]string{
		"id":           {"transaction_id", "id", "trans_id", "txn_id"},
		"country":      {"country", "nation", "country_name"},
		"region":       {"region", "state", "province", "area"},
		"product_name": {"product_name", "product", "item", "item_name"},
		"price":        {"price", "unit_price", "cost", "amount"},
		"quantity":     {"quantity", "qty", "amount", "count"},
		"date":         {"transaction_date", "date", "timestamp", "time"},
	}

	// Transaction ID
	tx.ID = fc.getFieldValue(data, fieldMappings["id"])
	if tx.ID == "" {
		return nil, fmt.Errorf("transaction ID is required")
	}

	// Country
	tx.Country = fc.getFieldValue(data, fieldMappings["country"])
	if tx.Country == "" {
		tx.Country = fc.config.DefaultCountry
	}

	// Region
	tx.Region = fc.getFieldValue(data, fieldMappings["region"])
	if tx.Region == "" {
		tx.Region = fc.config.DefaultRegion
	}

	// Product Name
	tx.ProductName = fc.getFieldValue(data, fieldMappings["product_name"])
	if tx.ProductName == "" {
		return nil, fmt.Errorf("product name is required")
	}

	// Price
	priceVal, err := fc.getFieldNumericValue(data, fieldMappings["price"])
	if err != nil {
		return nil, fmt.Errorf("price is required: %w", err)
	}
	tx.UnitPriceCents = int64(priceVal * fc.config.PriceMultiplier)

	// Quantity
	quantityVal, err := fc.getFieldNumericValue(data, fieldMappings["quantity"])
	if err != nil {
		quantityVal = 1 // Default quantity
	}
	tx.Quantity = int64(quantityVal)

	// Date
	dateStr := fc.getFieldValue(data, fieldMappings["date"])
	if dateStr == "" {
		return nil, fmt.Errorf("transaction date is required")
	}

	date, err := fc.parseFlexibleDate(dateStr)
	if err != nil {
		return nil, fmt.Errorf("invalid date: %w", err)
	}
	tx.TxTime = date

	return tx, nil
}

// getFieldValue gets string value using flexible field naming
func (fc *FormatConverter) getFieldValue(data map[string]interface{}, fieldNames []string) string {
	for _, fieldName := range fieldNames {
		if val, exists := data[fieldName]; exists {
			if str, ok := val.(string); ok {
				return str
			}
			return fmt.Sprintf("%v", val)
		}
	}
	return ""
}

// getFieldNumericValue gets numeric value using flexible field naming
func (fc *FormatConverter) getFieldNumericValue(data map[string]interface{}, fieldNames []string) (float64, error) {
	for _, fieldName := range fieldNames {
		if val, exists := data[fieldName]; exists {
			switch v := val.(type) {
			case float64:
				return v, nil
			case int:
				return float64(v), nil
			case int64:
				return float64(v), nil
			case string:
				if parsed, err := strconv.ParseFloat(v, 64); err == nil {
					return parsed, nil
				}
			}
		}
	}
	return 0, fmt.Errorf("numeric field not found")
}

// ExportToFormat exports transactions to specified format
func (fc *FormatConverter) ExportToFormat(transactions []models.Transaction, format DataFormat, writer io.Writer) error {
	switch format {
	case FormatCSV:
		return fc.exportToCSV(transactions, writer, ',')
	case FormatTSV:
		return fc.exportToCSV(transactions, writer, '\t')
	case FormatJSON:
		return fc.exportToJSON(transactions, writer)
	case FormatYAML:
		return fc.exportToYAML(transactions, writer)
	default:
		return fmt.Errorf("unsupported export format: %s", format)
	}
}

// exportToCSV exports transactions to CSV format
func (fc *FormatConverter) exportToCSV(transactions []models.Transaction, writer io.Writer, delimiter rune) error {
	csvWriter := csv.NewWriter(writer)
	csvWriter.Comma = delimiter
	defer csvWriter.Flush()

	// Write header
	header := []string{
		"transaction_id", "country", "region", "product_name",
		"price", "quantity", "transaction_date",
	}
	if err := csvWriter.Write(header); err != nil {
		return err
	}

	// Write data
	for _, tx := range transactions {
		record := []string{
			tx.ID,
			tx.Country,
			tx.Region,
			tx.ProductName,
			fmt.Sprintf("%.2f", float64(tx.UnitPriceCents)/fc.config.PriceMultiplier),
			fmt.Sprintf("%d", tx.Quantity),
			tx.TxTime.Format("2006-01-02T15:04:05Z"),
		}
		if err := csvWriter.Write(record); err != nil {
			return err
		}
	}

	return nil
}

// exportToJSON exports transactions to JSON format
func (fc *FormatConverter) exportToJSON(transactions []models.Transaction, writer io.Writer) error {
	encoder := json.NewEncoder(writer)
	encoder.SetIndent("", "  ")
	return encoder.Encode(transactions)
}

// exportToYAML exports transactions to YAML format
func (fc *FormatConverter) exportToYAML(transactions []models.Transaction, writer io.Writer) error {
	encoder := yaml.NewEncoder(writer)
	defer encoder.Close()
	return encoder.Encode(transactions)
}
