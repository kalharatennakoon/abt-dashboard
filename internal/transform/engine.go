package transform

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"strconv"
	"strings"
	"time"

	"abt-dashboard/internal/models"
)

// DataTransformationEngine provides flexible data transformation capabilities
type DataTransformationEngine struct {
	transformations []Transformation
	validators      []Validator
	optimizations   []Optimization
	config          TransformConfig
}

// TransformConfig defines configuration for data transformations
type TransformConfig struct {
	EnableValidation   bool              `json:"enable_validation"`
	EnableOptimization bool              `json:"enable_optimization"`
	DateFormats        []string          `json:"date_formats"`
	CurrencyFormats    []string          `json:"currency_formats"`
	NullValues         []string          `json:"null_values"`
	DefaultCountry     string            `json:"default_country"`
	DefaultRegion      string            `json:"default_region"`
	PriceMultiplier    float64           `json:"price_multiplier"`
	CustomMappings     map[string]string `json:"custom_mappings"`
	DataTypes          map[string]string `json:"data_types"`
}

// Transformation interface for data transformation operations
type Transformation interface {
	Name() string
	Transform(data interface{}) (interface{}, error)
	Description() string
}

// Validator interface for data validation
type Validator interface {
	Name() string
	Validate(data interface{}) error
	Description() string
}

// Optimization interface for data optimization
type Optimization interface {
	Name() string
	Optimize(data interface{}) (interface{}, error)
	Description() string
}

// TransformationResult contains the result of data transformation
type TransformationResult struct {
	OriginalRecords    int                `json:"original_records"`
	TransformedRecords int                `json:"transformed_records"`
	SkippedRecords     int                `json:"skipped_records"`
	Errors             []string           `json:"errors"`
	Warnings           []string           `json:"warnings"`
	Transformations    []string           `json:"transformations_applied"`
	ProcessingTime     time.Duration      `json:"processing_time"`
	DataQuality        DataQualityMetrics `json:"data_quality"`
}

// DataQualityMetrics provides insights into data quality
type DataQualityMetrics struct {
	Completeness float64            `json:"completeness"`
	Consistency  float64            `json:"consistency"`
	Validity     float64            `json:"validity"`
	Uniqueness   float64            `json:"uniqueness"`
	FieldMetrics map[string]float64 `json:"field_metrics"`
}

// NewDataTransformationEngine creates a new transformation engine
func NewDataTransformationEngine(config TransformConfig) *DataTransformationEngine {
	engine := &DataTransformationEngine{
		config: config,
	}

	// Initialize default configurations
	if len(config.DateFormats) == 0 {
		engine.config.DateFormats = []string{
			time.RFC3339,
			"2006-01-02",
			"2006-01-02 15:04:05",
			"01/02/2006",
			"02-01-2006",
			"2006/01/02",
		}
	}

	if len(config.NullValues) == 0 {
		engine.config.NullValues = []string{"", "NULL", "null", "N/A", "n/a", "-"}
	}

	if config.PriceMultiplier == 0 {
		engine.config.PriceMultiplier = 100 // Default: dollars to cents
	}

	// Register default transformations
	engine.RegisterTransformation(&CurrencyNormalization{config: config})
	engine.RegisterTransformation(&DateNormalization{config: config})
	engine.RegisterTransformation(&StringCleaning{config: config})
	engine.RegisterTransformation(&CountryMapping{config: config})
	engine.RegisterTransformation(&RegionMapping{config: config})
	engine.RegisterTransformation(&ProductNameNormalization{config: config})

	// Register default validators
	engine.RegisterValidator(&RequiredFieldValidator{})
	engine.RegisterValidator(&DataTypeValidator{config: config})
	engine.RegisterValidator(&RangeValidator{})
	engine.RegisterValidator(&UniquenessValidator{})

	// Register default optimizations
	engine.RegisterOptimization(&DuplicateRemoval{})
	engine.RegisterOptimization(&DataDeduplication{})
	engine.RegisterOptimization(&IndexOptimization{})

	return engine
}

// RegisterTransformation adds a new transformation to the engine
func (e *DataTransformationEngine) RegisterTransformation(t Transformation) {
	e.transformations = append(e.transformations, t)
}

// RegisterValidator adds a new validator to the engine
func (e *DataTransformationEngine) RegisterValidator(v Validator) {
	e.validators = append(e.validators, v)
}

// RegisterOptimization adds a new optimization to the engine
func (e *DataTransformationEngine) RegisterOptimization(o Optimization) {
	e.optimizations = append(e.optimizations, o)
}

// TransformCSVData processes CSV data with all registered transformations
func (e *DataTransformationEngine) TransformCSVData(reader io.Reader) ([]models.Transaction, *TransformationResult, error) {
	startTime := time.Now()
	result := &TransformationResult{
		Errors:          make([]string, 0),
		Warnings:        make([]string, 0),
		Transformations: make([]string, 0),
	}

	// Parse CSV
	csvReader := csv.NewReader(reader)
	csvReader.TrimLeadingSpace = true

	// Read header
	header, err := csvReader.Read()
	if err != nil {
		return nil, result, fmt.Errorf("failed to read CSV header: %w", err)
	}

	// Create column mapping
	columnMap := make(map[string]int)
	for i, col := range header {
		columnMap[strings.ToLower(strings.TrimSpace(col))] = i
	}

	var transactions []models.Transaction
	var rawRecords [][]string

	// Read all records
	for {
		record, err := csvReader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			result.Errors = append(result.Errors, fmt.Sprintf("CSV read error: %s", err.Error()))
			continue
		}
		rawRecords = append(rawRecords, record)
		result.OriginalRecords++
	}

	// Process each record
	for _, record := range rawRecords {
		transaction, err := e.transformRecord(record, columnMap)
		if err != nil {
			result.Errors = append(result.Errors, err.Error())
			result.SkippedRecords++
			continue
		}

		// Validate if enabled
		if e.config.EnableValidation {
			if err := e.validateTransaction(transaction); err != nil {
				result.Warnings = append(result.Warnings, err.Error())
			}
		}

		transactions = append(transactions, *transaction)
		result.TransformedRecords++
	}

	// Apply optimizations if enabled
	if e.config.EnableOptimization {
		optimizedTransactions, err := e.optimizeData(transactions)
		if err != nil {
			result.Warnings = append(result.Warnings, fmt.Sprintf("Optimization warning: %s", err.Error()))
		} else {
			transactions = optimizedTransactions
		}
	}

	// Calculate data quality metrics
	result.DataQuality = e.calculateDataQuality(transactions)
	result.ProcessingTime = time.Since(startTime)

	// Record applied transformations
	for _, t := range e.transformations {
		result.Transformations = append(result.Transformations, t.Name())
	}

	log.Printf("Data transformation completed: %d/%d records processed in %v",
		result.TransformedRecords, result.OriginalRecords, result.ProcessingTime)

	return transactions, result, nil
}

// transformRecord transforms a single CSV record into a Transaction
func (e *DataTransformationEngine) transformRecord(record []string, columnMap map[string]int) (*models.Transaction, error) {
	transaction := &models.Transaction{}

	// Extract and transform fields
	if idx, ok := columnMap["transaction_id"]; ok && idx < len(record) {
		transaction.ID = e.cleanString(record[idx])
	}

	if idx, ok := columnMap["country"]; ok && idx < len(record) {
		country := e.cleanString(record[idx])
		transaction.Country = e.mapCountry(country)
	}

	if idx, ok := columnMap["region"]; ok && idx < len(record) {
		region := e.cleanString(record[idx])
		transaction.Region = e.mapRegion(region)
	}

	if idx, ok := columnMap["product_name"]; ok && idx < len(record) {
		transaction.ProductName = e.normalizeProductName(record[idx])
	}

	if idx, ok := columnMap["price"]; ok && idx < len(record) {
		price, err := e.parsePrice(record[idx])
		if err != nil {
			return nil, fmt.Errorf("invalid price '%s': %w", record[idx], err)
		}
		transaction.UnitPriceCents = price
	}

	if idx, ok := columnMap["quantity"]; ok && idx < len(record) {
		quantity, err := strconv.ParseInt(e.cleanString(record[idx]), 10, 64)
		if err != nil {
			return nil, fmt.Errorf("invalid quantity '%s': %w", record[idx], err)
		}
		transaction.Quantity = quantity
	}

	if idx, ok := columnMap["transaction_date"]; ok && idx < len(record) {
		date, err := e.parseDate(record[idx])
		if err != nil {
			return nil, fmt.Errorf("invalid date '%s': %w", record[idx], err)
		}
		transaction.TxTime = date
	}

	// Apply all transformations
	for _, transformation := range e.transformations {
		transformedData, err := transformation.Transform(transaction)
		if err != nil {
			log.Printf("Transformation %s failed: %v", transformation.Name(), err)
			continue
		}
		if t, ok := transformedData.(*models.Transaction); ok {
			transaction = t
		}
	}

	return transaction, nil
}

// Helper methods for data transformation

func (e *DataTransformationEngine) cleanString(s string) string {
	s = strings.TrimSpace(s)
	for _, nullValue := range e.config.NullValues {
		if s == nullValue {
			return ""
		}
	}
	return s
}

func (e *DataTransformationEngine) mapCountry(country string) string {
	if country == "" && e.config.DefaultCountry != "" {
		return e.config.DefaultCountry
	}

	// Apply custom mappings
	if mapped, ok := e.config.CustomMappings[strings.ToLower(country)]; ok {
		return mapped
	}

	// Normalize common country name variations
	countryMappings := map[string]string{
		"usa":         "United States",
		"us":          "United States",
		"america":     "United States",
		"uk":          "United Kingdom",
		"britain":     "United Kingdom",
		"england":     "United Kingdom",
		"uae":         "United Arab Emirates",
		"south korea": "Korea",
		"north korea": "Korea",
	}

	if mapped, ok := countryMappings[strings.ToLower(country)]; ok {
		return mapped
	}

	return country
}

func (e *DataTransformationEngine) mapRegion(region string) string {
	if region == "" && e.config.DefaultRegion != "" {
		return e.config.DefaultRegion
	}

	// Apply custom mappings
	if mapped, ok := e.config.CustomMappings[strings.ToLower("region_"+region)]; ok {
		return mapped
	}

	return region
}

func (e *DataTransformationEngine) normalizeProductName(productName string) string {
	productName = e.cleanString(productName)

	// Apply custom mappings
	if mapped, ok := e.config.CustomMappings[strings.ToLower("product_"+productName)]; ok {
		return mapped
	}

	// Normalize product names
	productName = strings.Title(strings.ToLower(productName))

	return productName
}

func (e *DataTransformationEngine) parsePrice(priceStr string) (int64, error) {
	priceStr = e.cleanString(priceStr)

	// Remove currency symbols
	for _, symbol := range []string{"$", "€", "£", "¥", "₹", "₽"} {
		priceStr = strings.ReplaceAll(priceStr, symbol, "")
	}

	// Remove commas
	priceStr = strings.ReplaceAll(priceStr, ",", "")

	price, err := strconv.ParseFloat(priceStr, 64)
	if err != nil {
		return 0, err
	}

	// Convert to cents using multiplier
	return int64(price * e.config.PriceMultiplier), nil
}

func (e *DataTransformationEngine) parseDate(dateStr string) (time.Time, error) {
	dateStr = e.cleanString(dateStr)

	for _, format := range e.config.DateFormats {
		if date, err := time.Parse(format, dateStr); err == nil {
			return date, nil
		}
	}

	return time.Time{}, fmt.Errorf("unable to parse date with any format")
}

func (e *DataTransformationEngine) validateTransaction(transaction *models.Transaction) error {
	for _, validator := range e.validators {
		if err := validator.Validate(transaction); err != nil {
			return fmt.Errorf("validation %s failed: %w", validator.Name(), err)
		}
	}
	return nil
}

func (e *DataTransformationEngine) optimizeData(transactions []models.Transaction) ([]models.Transaction, error) {
	var optimizedTransactions interface{} = transactions

	for _, optimization := range e.optimizations {
		optimized, err := optimization.Optimize(optimizedTransactions)
		if err != nil {
			return transactions, fmt.Errorf("optimization %s failed: %w", optimization.Name(), err)
		}
		optimizedTransactions = optimized
	}

	if optimized, ok := optimizedTransactions.([]models.Transaction); ok {
		return optimized, nil
	}

	return transactions, nil
}

func (e *DataTransformationEngine) calculateDataQuality(transactions []models.Transaction) DataQualityMetrics {
	if len(transactions) == 0 {
		return DataQualityMetrics{}
	}

	metrics := DataQualityMetrics{
		FieldMetrics: make(map[string]float64),
	}

	// Calculate completeness
	completenessScores := make(map[string]int)
	totalFields := 7 // Number of fields in Transaction

	for _, tx := range transactions {
		if tx.ID != "" {
			completenessScores["id"]++
		}
		if tx.Country != "" {
			completenessScores["country"]++
		}
		if tx.Region != "" {
			completenessScores["region"]++
		}
		if tx.ProductName != "" {
			completenessScores["product_name"]++
		}
		if tx.UnitPriceCents > 0 {
			completenessScores["price"]++
		}
		if tx.Quantity > 0 {
			completenessScores["quantity"]++
		}
		if !tx.TxTime.IsZero() {
			completenessScores["date"]++
		}
	}

	totalCompleteness := 0.0
	for field, count := range completenessScores {
		completeness := float64(count) / float64(len(transactions))
		metrics.FieldMetrics[field+"_completeness"] = completeness
		totalCompleteness += completeness
	}
	metrics.Completeness = totalCompleteness / float64(totalFields)

	// Calculate uniqueness (for transaction IDs)
	uniqueIDs := make(map[string]bool)
	for _, tx := range transactions {
		uniqueIDs[tx.ID] = true
	}
	metrics.Uniqueness = float64(len(uniqueIDs)) / float64(len(transactions))

	// Basic validity check (non-negative prices and quantities)
	validTransactions := 0
	for _, tx := range transactions {
		if tx.UnitPriceCents >= 0 && tx.Quantity >= 0 {
			validTransactions++
		}
	}
	metrics.Validity = float64(validTransactions) / float64(len(transactions))

	// Consistency score (basic implementation)
	metrics.Consistency = 0.95 // Placeholder - could implement more sophisticated consistency checks

	return metrics
}

// ExportTransformationReport exports the transformation results to JSON
func (e *DataTransformationEngine) ExportTransformationReport(result *TransformationResult, writer io.Writer) error {
	report := map[string]interface{}{
		"transformation_summary":     result,
		"engine_config":              e.config,
		"registered_transformations": e.getTransformationInfo(),
		"registered_validators":      e.getValidatorInfo(),
		"registered_optimizations":   e.getOptimizationInfo(),
		"timestamp":                  time.Now().Format(time.RFC3339),
	}

	encoder := json.NewEncoder(writer)
	encoder.SetIndent("", "  ")
	return encoder.Encode(report)
}

func (e *DataTransformationEngine) getTransformationInfo() []map[string]string {
	var info []map[string]string
	for _, t := range e.transformations {
		info = append(info, map[string]string{
			"name":        t.Name(),
			"description": t.Description(),
		})
	}
	return info
}

func (e *DataTransformationEngine) getValidatorInfo() []map[string]string {
	var info []map[string]string
	for _, v := range e.validators {
		info = append(info, map[string]string{
			"name":        v.Name(),
			"description": v.Description(),
		})
	}
	return info
}

func (e *DataTransformationEngine) getOptimizationInfo() []map[string]string {
	var info []map[string]string
	for _, o := range e.optimizations {
		info = append(info, map[string]string{
			"name":        o.Name(),
			"description": o.Description(),
		})
	}
	return info
}

// GetSupportedFormats returns all supported data formats and transformations
func (e *DataTransformationEngine) GetSupportedFormats() map[string]interface{} {
	return map[string]interface{}{
		"date_formats":     e.config.DateFormats,
		"currency_formats": e.config.CurrencyFormats,
		"null_values":      e.config.NullValues,
		"data_types":       e.config.DataTypes,
		"transformations":  e.getTransformationInfo(),
		"validators":       e.getValidatorInfo(),
		"optimizations":    e.getOptimizationInfo(),
	}
}
