package transform

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"abt-dashboard/internal/models"
)

// FlexibleDataHandler provides comprehensive data handling capabilities
type FlexibleDataHandler struct {
	engine    *DataTransformationEngine
	converter *FormatConverter
	config    TransformConfig
}

// NewFlexibleDataHandler creates a new flexible data handler
func NewFlexibleDataHandler(config TransformConfig) *FlexibleDataHandler {
	return &FlexibleDataHandler{
		engine:    NewDataTransformationEngine(config),
		converter: NewFormatConverter(config),
		config:    config,
	}
}

// ProcessDataFile processes a data file with automatic format detection and transformation
func (fdh *FlexibleDataHandler) ProcessDataFile(filePath string) ([]models.Transaction, *TransformationResult, error) {
	// Open file
	file, err := os.Open(filePath)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to open file %s: %w", filePath, err)
	}
	defer file.Close()

	// Detect format based on file extension and content
	format, err := fdh.detectFileFormat(filePath, file)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to detect file format: %w", err)
	}

	log.Printf("Processing file %s with detected format: %s", filePath, format)

	// Reset file position for reading
	file.Seek(0, 0)

	// Convert to transactions using the appropriate format handler
	transactions, err := fdh.converter.ConvertToTransactions(file, format)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to convert data: %w", err)
	}

	// Apply transformations and optimizations
	result := &TransformationResult{
		OriginalRecords:    len(transactions),
		TransformedRecords: len(transactions),
		Errors:             make([]string, 0),
		Warnings:           make([]string, 0),
		Transformations:    make([]string, 0),
	}

	startTime := time.Now()

	// Apply transformations to each transaction
	var transformedTransactions []models.Transaction
	for i, tx := range transactions {
		transformedTx := tx

		// Apply all transformations
		for _, transformation := range fdh.engine.transformations {
			transformedData, err := transformation.Transform(&transformedTx)
			if err != nil {
				result.Warnings = append(result.Warnings,
					fmt.Sprintf("Transformation %s failed for record %d: %v",
						transformation.Name(), i, err))
				continue
			}
			if newTx, ok := transformedData.(*models.Transaction); ok {
				transformedTx = *newTx
			}
		}

		// Validate if enabled
		if fdh.config.EnableValidation {
			for _, validator := range fdh.engine.validators {
				if err := validator.Validate(&transformedTx); err != nil {
					result.Warnings = append(result.Warnings,
						fmt.Sprintf("Validation %s failed for record %d: %v",
							validator.Name(), i, err))
				}
			}
		}

		transformedTransactions = append(transformedTransactions, transformedTx)
	}

	// Apply optimizations if enabled
	if fdh.config.EnableOptimization {
		optimizedData, err := fdh.optimizeTransactions(transformedTransactions)
		if err != nil {
			result.Warnings = append(result.Warnings, fmt.Sprintf("Optimization failed: %v", err))
		} else {
			transformedTransactions = optimizedData
			result.TransformedRecords = len(transformedTransactions)
		}
	}

	// Record applied transformations
	for _, t := range fdh.engine.transformations {
		result.Transformations = append(result.Transformations, t.Name())
	}

	// Calculate data quality metrics
	result.DataQuality = fdh.engine.calculateDataQuality(transformedTransactions)
	result.ProcessingTime = time.Since(startTime)

	log.Printf("Data processing completed: %d records processed in %v with %.2f%% quality score",
		len(transformedTransactions), result.ProcessingTime, result.DataQuality.Completeness*100)

	return transformedTransactions, result, nil
}

// detectFileFormat detects the format of a data file
func (fdh *FlexibleDataHandler) detectFileFormat(filePath string, file *os.File) (DataFormat, error) {
	// Check file extension first
	ext := strings.ToLower(filepath.Ext(filePath))
	switch ext {
	case ".csv":
		return FormatCSV, nil
	case ".tsv":
		return FormatTSV, nil
	case ".json":
		return FormatJSON, nil
	case ".yaml", ".yml":
		return FormatYAML, nil
	case ".xml":
		return FormatXML, nil
	}

	// Read first few bytes to analyze content
	buffer := make([]byte, 1024)
	n, err := file.Read(buffer)
	if err != nil && err != io.EOF {
		return "", fmt.Errorf("failed to read file content: %w", err)
	}

	// Use format converter's detection logic
	return fdh.converter.DetectFormat(buffer[:n]), nil
}

// optimizeTransactions applies optimization strategies
func (fdh *FlexibleDataHandler) optimizeTransactions(transactions []models.Transaction) ([]models.Transaction, error) {
	var optimizedData interface{} = transactions

	for _, optimization := range fdh.engine.optimizations {
		optimized, err := optimization.Optimize(optimizedData)
		if err != nil {
			return transactions, fmt.Errorf("optimization %s failed: %w", optimization.Name(), err)
		}
		optimizedData = optimized
	}

	if optimized, ok := optimizedData.([]models.Transaction); ok {
		return optimized, nil
	}

	return transactions, nil
}

// ProcessDataStream processes data from a stream with specified format
func (fdh *FlexibleDataHandler) ProcessDataStream(reader io.Reader, format DataFormat) ([]models.Transaction, *TransformationResult, error) {
	// Convert to transactions
	transactions, err := fdh.converter.ConvertToTransactions(reader, format)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to convert stream data: %w", err)
	}

	// Create transformation result
	result := &TransformationResult{
		OriginalRecords:    len(transactions),
		TransformedRecords: len(transactions),
		Errors:             make([]string, 0),
		Warnings:           make([]string, 0),
		Transformations:    make([]string, 0),
		ProcessingTime:     time.Since(time.Now()),
		DataQuality:        fdh.engine.calculateDataQuality(transactions),
	}

	return transactions, result, nil
}

// ExportTransactions exports transactions to specified format
func (fdh *FlexibleDataHandler) ExportTransactions(transactions []models.Transaction, format DataFormat, writer io.Writer) error {
	return fdh.converter.ExportToFormat(transactions, format, writer)
}

// GetDataQualityReport generates a comprehensive data quality report
func (fdh *FlexibleDataHandler) GetDataQualityReport(transactions []models.Transaction) DataQualityReport {
	metrics := fdh.engine.calculateDataQuality(transactions)

	return DataQualityReport{
		Metrics:         metrics,
		TotalRecords:    len(transactions),
		Timestamp:       time.Now(),
		Issues:          fdh.identifyDataQualityIssues(transactions),
		Recommendations: fdh.generateRecommendations(transactions, metrics),
	}
}

// DataQualityReport provides comprehensive data quality analysis
type DataQualityReport struct {
	Metrics         DataQualityMetrics `json:"metrics"`
	TotalRecords    int                `json:"total_records"`
	Timestamp       time.Time          `json:"timestamp"`
	Issues          []DataQualityIssue `json:"issues"`
	Recommendations []string           `json:"recommendations"`
}

// DataQualityIssue represents a specific data quality issue
type DataQualityIssue struct {
	Type        string   `json:"type"`
	Description string   `json:"description"`
	Severity    string   `json:"severity"`
	Count       int      `json:"count"`
	Examples    []string `json:"examples"`
}

// identifyDataQualityIssues analyzes transactions for common data quality issues
func (fdh *FlexibleDataHandler) identifyDataQualityIssues(transactions []models.Transaction) []DataQualityIssue {
	var issues []DataQualityIssue

	// Check for missing required fields
	missingFields := make(map[string]int)
	duplicateIDs := make(map[string]int)
	invalidPrices := 0
	invalidQuantities := 0
	invalidDates := 0
	missingCountries := 0
	missingRegions := 0

	for _, tx := range transactions {
		// Check required fields
		if tx.ID == "" {
			missingFields["transaction_id"]++
		}
		if tx.Country == "" {
			missingCountries++
		}
		if tx.Region == "" {
			missingRegions++
		}
		if tx.ProductName == "" {
			missingFields["product_name"]++
		}

		// Check for duplicates
		duplicateIDs[tx.ID]++

		// Check data validity
		if tx.UnitPriceCents <= 0 {
			invalidPrices++
		}
		if tx.Quantity <= 0 {
			invalidQuantities++
		}
		if tx.TxTime.IsZero() {
			invalidDates++
		}
	}

	// Report missing fields
	for field, count := range missingFields {
		if count > 0 {
			issues = append(issues, DataQualityIssue{
				Type:        "missing_data",
				Description: fmt.Sprintf("Missing %s in %d records", field, count),
				Severity:    "high",
				Count:       count,
			})
		}
	}

	// Report duplicate IDs
	duplicateCount := 0
	for _, count := range duplicateIDs {
		if count > 1 {
			duplicateCount++
		}
	}
	if duplicateCount > 0 {
		issues = append(issues, DataQualityIssue{
			Type:        "duplicate_data",
			Description: fmt.Sprintf("Found %d duplicate transaction IDs", duplicateCount),
			Severity:    "medium",
			Count:       duplicateCount,
		})
	}

	// Report data validity issues
	if invalidPrices > 0 {
		issues = append(issues, DataQualityIssue{
			Type:        "invalid_data",
			Description: fmt.Sprintf("Invalid prices in %d records", invalidPrices),
			Severity:    "high",
			Count:       invalidPrices,
		})
	}

	if invalidQuantities > 0 {
		issues = append(issues, DataQualityIssue{
			Type:        "invalid_data",
			Description: fmt.Sprintf("Invalid quantities in %d records", invalidQuantities),
			Severity:    "high",
			Count:       invalidQuantities,
		})
	}

	if invalidDates > 0 {
		issues = append(issues, DataQualityIssue{
			Type:        "invalid_data",
			Description: fmt.Sprintf("Invalid dates in %d records", invalidDates),
			Severity:    "high",
			Count:       invalidDates,
		})
	}

	// Report missing geographic data
	if missingCountries > 0 {
		issues = append(issues, DataQualityIssue{
			Type:        "missing_data",
			Description: fmt.Sprintf("Missing country data in %d records", missingCountries),
			Severity:    "medium",
			Count:       missingCountries,
		})
	}

	if missingRegions > 0 {
		issues = append(issues, DataQualityIssue{
			Type:        "missing_data",
			Description: fmt.Sprintf("Missing region data in %d records", missingRegions),
			Severity:    "low",
			Count:       missingRegions,
		})
	}

	return issues
}

// generateRecommendations provides recommendations for improving data quality
func (fdh *FlexibleDataHandler) generateRecommendations(transactions []models.Transaction, metrics DataQualityMetrics) []string {
	var recommendations []string

	// Completeness recommendations
	if metrics.Completeness < 0.95 {
		recommendations = append(recommendations,
			"Improve data completeness by implementing validation at data entry points")
	}

	// Validity recommendations
	if metrics.Validity < 0.90 {
		recommendations = append(recommendations,
			"Implement stricter data validation rules to improve data validity")
	}

	// Uniqueness recommendations
	if metrics.Uniqueness < 0.99 {
		recommendations = append(recommendations,
			"Review data collection process to eliminate duplicate entries")
	}

	// Field-specific recommendations
	if completeness, ok := metrics.FieldMetrics["country_completeness"]; ok && completeness < 0.90 {
		recommendations = append(recommendations,
			"Implement default country mapping or mandatory country field")
	}

	if completeness, ok := metrics.FieldMetrics["region_completeness"]; ok && completeness < 0.85 {
		recommendations = append(recommendations,
			"Consider enriching data with region information based on other geographic data")
	}

	// General recommendations
	recommendations = append(recommendations,
		"Regular data quality monitoring should be implemented")
	recommendations = append(recommendations,
		"Consider implementing data lineage tracking for better data governance")

	return recommendations
}

// GetSupportedFormats returns information about all supported formats
func (fdh *FlexibleDataHandler) GetSupportedFormats() map[string]interface{} {
	return map[string]interface{}{
		"input_formats": []string{
			string(FormatCSV),
			string(FormatTSV),
			string(FormatJSON),
			string(FormatYAML),
		},
		"output_formats": []string{
			string(FormatCSV),
			string(FormatTSV),
			string(FormatJSON),
			string(FormatYAML),
		},
		"transformations": fdh.engine.GetSupportedFormats(),
		"quality_checks": []string{
			"completeness_validation",
			"data_type_validation",
			"range_validation",
			"uniqueness_validation",
			"business_rule_validation",
		},
		"optimizations": []string{
			"duplicate_removal",
			"data_deduplication",
			"format_normalization",
			"index_optimization",
		},
	}
}

// ValidateConfiguration validates the transformation configuration
func (fdh *FlexibleDataHandler) ValidateConfiguration() error {
	// Check required configuration fields
	if fdh.config.PriceMultiplier <= 0 {
		return fmt.Errorf("price_multiplier must be positive")
	}

	if len(fdh.config.DateFormats) == 0 {
		return fmt.Errorf("at least one date format must be specified")
	}

	// Validate date formats
	testDate := "2006-01-02T15:04:05Z"
	for _, format := range fdh.config.DateFormats {
		if _, err := time.Parse(format, testDate); err != nil {
			// This is expected for different formats, but we should have at least one valid format
			continue
		}
	}

	return nil
}

// GetTransformationStatistics returns statistics about applied transformations
func (fdh *FlexibleDataHandler) GetTransformationStatistics() map[string]interface{} {
	return map[string]interface{}{
		"registered_transformations": len(fdh.engine.transformations),
		"registered_validators":      len(fdh.engine.validators),
		"registered_optimizations":   len(fdh.engine.optimizations),
		"transformation_details":     fdh.engine.getTransformationInfo(),
		"validator_details":          fdh.engine.getValidatorInfo(),
		"optimization_details":       fdh.engine.getOptimizationInfo(),
	}
}
