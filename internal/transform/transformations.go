package transform

import (
	"fmt"
	"regexp"
	"strings"
	"time"

	"abt-dashboard/internal/models"
)

// CurrencyNormalization handles currency format normalization
type CurrencyNormalization struct {
	config TransformConfig
}

func (c *CurrencyNormalization) Name() string {
	return "CurrencyNormalization"
}

func (c *CurrencyNormalization) Description() string {
	return "Normalizes currency values by removing symbols and converting to standard format"
}

func (c *CurrencyNormalization) Transform(data interface{}) (interface{}, error) {
	if tx, ok := data.(*models.Transaction); ok {
		// Currency normalization is handled in parsePrice method
		return tx, nil
	}
	return data, nil
}

// DateNormalization handles date format normalization
type DateNormalization struct {
	config TransformConfig
}

func (d *DateNormalization) Name() string {
	return "DateNormalization"
}

func (d *DateNormalization) Description() string {
	return "Normalizes date formats to standard ISO format"
}

func (d *DateNormalization) Transform(data interface{}) (interface{}, error) {
	if tx, ok := data.(*models.Transaction); ok {
		// Ensure date is in UTC
		if !tx.TxTime.IsZero() {
			tx.TxTime = tx.TxTime.UTC()
		}
		return tx, nil
	}
	return data, nil
}

// StringCleaning handles string normalization and cleaning
type StringCleaning struct {
	config TransformConfig
}

func (s *StringCleaning) Name() string {
	return "StringCleaning"
}

func (s *StringCleaning) Description() string {
	return "Cleans and normalizes string values by trimming whitespace and removing special characters"
}

func (s *StringCleaning) Transform(data interface{}) (interface{}, error) {
	if tx, ok := data.(*models.Transaction); ok {
		tx.ID = s.cleanString(tx.ID)
		tx.Country = s.cleanString(tx.Country)
		tx.Region = s.cleanString(tx.Region)
		tx.ProductName = s.cleanString(tx.ProductName)
		return tx, nil
	}
	return data, nil
}

func (s *StringCleaning) cleanString(str string) string {
	// Remove extra whitespace
	str = strings.TrimSpace(str)

	// Remove non-printable characters
	re := regexp.MustCompile(`[^\p{L}\p{N}\p{P}\p{Z}]`)
	str = re.ReplaceAllString(str, "")

	// Normalize multiple spaces to single space
	re = regexp.MustCompile(`\s+`)
	str = re.ReplaceAllString(str, " ")

	return str
}

// CountryMapping handles country name standardization
type CountryMapping struct {
	config TransformConfig
}

func (c *CountryMapping) Name() string {
	return "CountryMapping"
}

func (c *CountryMapping) Description() string {
	return "Maps country names to standardized values"
}

func (c *CountryMapping) Transform(data interface{}) (interface{}, error) {
	if tx, ok := data.(*models.Transaction); ok {
		tx.Country = c.mapCountry(tx.Country)
		return tx, nil
	}
	return data, nil
}

func (c *CountryMapping) mapCountry(country string) string {
	// Standard country mappings
	mappings := map[string]string{
		"usa":                                   "United States",
		"us":                                    "United States",
		"america":                               "United States",
		"united states of america":              "United States",
		"uk":                                    "United Kingdom",
		"britain":                               "United Kingdom",
		"great britain":                         "United Kingdom",
		"england":                               "United Kingdom",
		"uae":                                   "United Arab Emirates",
		"south korea":                           "South Korea",
		"republic of korea":                     "South Korea",
		"north korea":                           "North Korea",
		"democratic people's republic of korea": "North Korea",
		"prc":                                   "China",
		"people's republic of china":            "China",
		"roc":                                   "Taiwan",
		"republic of china":                     "Taiwan",
	}

	normalizedCountry := strings.ToLower(strings.TrimSpace(country))
	if mapped, exists := mappings[normalizedCountry]; exists {
		return mapped
	}

	// Apply custom mappings from config
	if mapped, exists := c.config.CustomMappings[normalizedCountry]; exists {
		return mapped
	}

	// Return title case version
	return strings.Title(strings.ToLower(country))
}

// RegionMapping handles region name standardization
type RegionMapping struct {
	config TransformConfig
}

func (r *RegionMapping) Name() string {
	return "RegionMapping"
}

func (r *RegionMapping) Description() string {
	return "Maps region names to standardized values"
}

func (r *RegionMapping) Transform(data interface{}) (interface{}, error) {
	if tx, ok := data.(*models.Transaction); ok {
		tx.Region = r.mapRegion(tx.Region)
		return tx, nil
	}
	return data, nil
}

func (r *RegionMapping) mapRegion(region string) string {
	// Standard region mappings
	mappings := map[string]string{
		"n":       "North",
		"s":       "South",
		"e":       "East",
		"w":       "West",
		"ne":      "Northeast",
		"nw":      "Northwest",
		"se":      "Southeast",
		"sw":      "Southwest",
		"central": "Central",
		"centre":  "Central",
		"mid":     "Central",
		"middle":  "Central",
	}

	normalizedRegion := strings.ToLower(strings.TrimSpace(region))
	if mapped, exists := mappings[normalizedRegion]; exists {
		return mapped
	}

	// Apply custom mappings from config
	configKey := "region_" + normalizedRegion
	if mapped, exists := r.config.CustomMappings[configKey]; exists {
		return mapped
	}

	// Return title case version
	return strings.Title(strings.ToLower(region))
}

// ProductNameNormalization handles product name standardization
type ProductNameNormalization struct {
	config TransformConfig
}

func (p *ProductNameNormalization) Name() string {
	return "ProductNameNormalization"
}

func (p *ProductNameNormalization) Description() string {
	return "Normalizes product names for consistency"
}

func (p *ProductNameNormalization) Transform(data interface{}) (interface{}, error) {
	if tx, ok := data.(*models.Transaction); ok {
		tx.ProductName = p.normalizeProductName(tx.ProductName)
		return tx, nil
	}
	return data, nil
}

func (p *ProductNameNormalization) normalizeProductName(productName string) string {
	productName = strings.TrimSpace(productName)

	// Apply custom mappings from config
	configKey := "product_" + strings.ToLower(productName)
	if mapped, exists := p.config.CustomMappings[configKey]; exists {
		return mapped
	}

	// Standard product name patterns
	patterns := map[string]string{
		`(?i)widget`: "Widget",
		`(?i)gadget`: "Gadget",
		`(?i)device`: "Device",
		`(?i)tool`:   "Tool",
		`(?i)kit`:    "Kit",
		`(?i)set`:    "Set",
		`(?i)pack`:   "Pack",
		`(?i)bundle`: "Bundle",
	}

	for pattern, replacement := range patterns {
		re := regexp.MustCompile(pattern)
		productName = re.ReplaceAllString(productName, replacement)
	}

	// Capitalize first letter of each word
	words := strings.Fields(productName)
	for i, word := range words {
		if len(word) > 0 {
			words[i] = strings.ToUpper(word[:1]) + strings.ToLower(word[1:])
		}
	}

	return strings.Join(words, " ")
}

// RequiredFieldValidator validates that required fields are present
type RequiredFieldValidator struct{}

func (r *RequiredFieldValidator) Name() string {
	return "RequiredFieldValidator"
}

func (r *RequiredFieldValidator) Description() string {
	return "Validates that required fields are present and not empty"
}

func (r *RequiredFieldValidator) Validate(data interface{}) error {
	if tx, ok := data.(*models.Transaction); ok {
		if tx.ID == "" {
			return fmt.Errorf("transaction ID is required")
		}
		if tx.Country == "" {
			return fmt.Errorf("country is required")
		}
		if tx.ProductName == "" {
			return fmt.Errorf("product name is required")
		}
		if tx.UnitPriceCents <= 0 {
			return fmt.Errorf("unit price must be positive")
		}
		if tx.Quantity <= 0 {
			return fmt.Errorf("quantity must be positive")
		}
		if tx.TxTime.IsZero() {
			return fmt.Errorf("transaction time is required")
		}
	}
	return nil
}

// DataTypeValidator validates data types and formats
type DataTypeValidator struct {
	config TransformConfig
}

func (d *DataTypeValidator) Name() string {
	return "DataTypeValidator"
}

func (d *DataTypeValidator) Description() string {
	return "Validates data types and formats according to configuration"
}

func (d *DataTypeValidator) Validate(data interface{}) error {
	if tx, ok := data.(*models.Transaction); ok {
		// Validate ID format (should be alphanumeric)
		if matched, _ := regexp.MatchString(`^[a-zA-Z0-9_-]+$`, tx.ID); !matched {
			return fmt.Errorf("transaction ID contains invalid characters")
		}

		// Validate country name (should contain only letters, spaces, and common punctuation)
		if matched, _ := regexp.MatchString(`^[a-zA-Z\s\.\-']+$`, tx.Country); !matched {
			return fmt.Errorf("country name contains invalid characters")
		}

		// Validate price range (reasonable business limits)
		if tx.UnitPriceCents > 10000000 { // $100,000
			return fmt.Errorf("unit price exceeds reasonable maximum")
		}

		// Validate quantity range
		if tx.Quantity > 1000000 {
			return fmt.Errorf("quantity exceeds reasonable maximum")
		}

		// Validate date range (not too far in past or future)
		now := time.Now()
		if tx.TxTime.Before(now.AddDate(-10, 0, 0)) {
			return fmt.Errorf("transaction date is too far in the past")
		}
		if tx.TxTime.After(now.AddDate(1, 0, 0)) {
			return fmt.Errorf("transaction date is in the future")
		}
	}
	return nil
}

// RangeValidator validates that numeric values are within acceptable ranges
type RangeValidator struct{}

func (r *RangeValidator) Name() string {
	return "RangeValidator"
}

func (r *RangeValidator) Description() string {
	return "Validates that numeric values are within acceptable business ranges"
}

func (r *RangeValidator) Validate(data interface{}) error {
	if tx, ok := data.(*models.Transaction); ok {
		// Price range validation
		if tx.UnitPriceCents < 1 || tx.UnitPriceCents > 50000000 { // $0.01 to $500,000
			return fmt.Errorf("unit price %d cents is outside acceptable range", tx.UnitPriceCents)
		}

		// Quantity range validation
		if tx.Quantity < 1 || tx.Quantity > 100000 {
			return fmt.Errorf("quantity %d is outside acceptable range", tx.Quantity)
		}
	}
	return nil
}

// UniquenessValidator validates uniqueness constraints
type UniquenessValidator struct {
	seenIDs map[string]bool
}

func (u *UniquenessValidator) Name() string {
	return "UniquenessValidator"
}

func (u *UniquenessValidator) Description() string {
	return "Validates uniqueness constraints such as transaction ID uniqueness"
}

func (u *UniquenessValidator) Validate(data interface{}) error {
	if u.seenIDs == nil {
		u.seenIDs = make(map[string]bool)
	}

	if tx, ok := data.(*models.Transaction); ok {
		if u.seenIDs[tx.ID] {
			return fmt.Errorf("duplicate transaction ID: %s", tx.ID)
		}
		u.seenIDs[tx.ID] = true
	}
	return nil
}

// DuplicateRemoval removes duplicate transactions
type DuplicateRemoval struct{}

func (d *DuplicateRemoval) Name() string {
	return "DuplicateRemoval"
}

func (d *DuplicateRemoval) Description() string {
	return "Removes duplicate transactions based on transaction ID"
}

func (d *DuplicateRemoval) Optimize(data interface{}) (interface{}, error) {
	if transactions, ok := data.([]models.Transaction); ok {
		seen := make(map[string]bool)
		var unique []models.Transaction

		for _, tx := range transactions {
			if !seen[tx.ID] {
				seen[tx.ID] = true
				unique = append(unique, tx)
			}
		}

		return unique, nil
	}
	return data, nil
}

// DataDeduplication performs advanced deduplication based on multiple fields
type DataDeduplication struct{}

func (d *DataDeduplication) Name() string {
	return "DataDeduplication"
}

func (d *DataDeduplication) Description() string {
	return "Performs advanced deduplication based on multiple transaction fields"
}

func (d *DataDeduplication) Optimize(data interface{}) (interface{}, error) {
	if transactions, ok := data.([]models.Transaction); ok {
		seen := make(map[string]bool)
		var unique []models.Transaction

		for _, tx := range transactions {
			// Create composite key for deduplication
			key := fmt.Sprintf("%s:%s:%s:%s:%d:%d:%d",
				tx.ID, tx.Country, tx.Region, tx.ProductName,
				tx.UnitPriceCents, tx.Quantity, tx.TxTime.Unix())

			if !seen[key] {
				seen[key] = true
				unique = append(unique, tx)
			}
		}

		return unique, nil
	}
	return data, nil
}

// IndexOptimization optimizes data for better query performance
type IndexOptimization struct{}

func (i *IndexOptimization) Name() string {
	return "IndexOptimization"
}

func (i *IndexOptimization) Description() string {
	return "Optimizes data structure for better query performance"
}

func (i *IndexOptimization) Optimize(data interface{}) (interface{}, error) {
	if transactions, ok := data.([]models.Transaction); ok {
		// Sort transactions by date for better time-based query performance
		// Note: In a real implementation, you might create indexes or optimize data structures
		// For now, we'll just return the transactions as-is
		return transactions, nil
	}
	return data, nil
}
