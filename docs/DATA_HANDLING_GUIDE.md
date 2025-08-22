# Flexible Data Handling & Transformation Guide

## Overview

The ABT Dashboard implements a comprehensive data handling system that provides maximum flexibility in processing various data formats, applying transformations, and ensuring data quality. This system is designed to handle real-world data challenges including format variations, data quality issues, and optimization requirements.

## Core Components

### 1. Data Transformation Engine (`internal/transform/engine.go`)

The transformation engine is the core component that orchestrates all data processing operations:

```go
// Create a transformation engine with custom configuration
config := TransformConfig{
    EnableValidation:   true,
    EnableOptimization: true,
    PriceMultiplier:   100.0, // Convert dollars to cents
    DateFormats: []string{
        "2006-01-02T15:04:05Z",
        "2006-01-02",
        "01/02/2006",
    },
}

engine := NewDataTransformationEngine(config)
```

**Key Features:**
- **Pluggable Architecture**: Support for custom transformations, validators, and optimizations
- **Comprehensive Logging**: Detailed processing logs and quality metrics
- **Error Handling**: Configurable error handling strategies
- **Performance Monitoring**: Processing time and memory usage tracking

### 2. Format Converter (`internal/transform/format_converter.go`)

Handles conversion between multiple data formats with automatic format detection:

**Supported Input Formats:**
- **CSV** (Comma-Separated Values)
- **TSV** (Tab-Separated Values)  
- **JSON** (JavaScript Object Notation)
- **YAML** (YAML Ain't Markup Language)
- **XML** (eXtensible Markup Language) - planned

**Format Detection Logic:**
```go
converter := NewFormatConverter(config)

// Automatic format detection
format := converter.DetectFormat(data)

// Manual format specification
transactions, err := converter.ConvertToTransactions(reader, FormatJSON)
```

### 3. Flexible Data Handler (`internal/transform/data_handler.go`)

Provides high-level interface for complete data processing workflows:

```go
handler := NewFlexibleDataHandler(config)

// Process any supported file format
transactions, result, err := handler.ProcessDataFile("data.csv")

// Generate data quality report
qualityReport := handler.GetDataQualityReport(transactions)
```

## Data Transformation Capabilities

### 1. Built-in Transformations

#### Currency Normalization
- Removes currency symbols ($, €, £, ¥, etc.)
- Handles thousand separators (commas, spaces)
- Converts to standard cents format
- Supports multiple currency formats

```yaml
# Example currency inputs supported:
- "$1,234.56"    → 123456 cents
- "€999.99"      → 99999 cents  
- "1234.50"      → 123450 cents
- "£50"          → 5000 cents
```

#### Date Normalization
- Supports 15+ date formats
- Automatic timezone handling
- Flexible parsing with fallback formats
- Unix timestamp support

```yaml
# Supported date formats:
- "2024-03-15T10:30:45Z"     # ISO 8601
- "2024-03-15"               # YYYY-MM-DD
- "03/15/2024"               # MM/DD/YYYY
- "15-03-2024"               # DD-MM-YYYY
- "Mar 15, 2024"             # Month DD, YYYY
- "1710503445"               # Unix timestamp
```

#### String Cleaning
- Removes extra whitespace
- Eliminates non-printable characters
- Normalizes encoding issues
- Standardizes case formatting

#### Geographic Standardization
- Country name normalization (USA → United States)
- Region mapping (N → North, SW → Southwest)
- Custom mapping support via configuration
- Default value assignment for missing data

#### Product Name Normalization
- Consistent capitalization
- Brand/model standardization
- Custom product mapping
- Duplicate product detection

### 2. Custom Transformations

The system supports custom transformations through the `Transformation` interface:

```go
type CustomTransformation struct {
    config TransformConfig
}

func (ct *CustomTransformation) Name() string {
    return "CustomTransformation"
}

func (ct *CustomTransformation) Description() string {
    return "Applies custom business logic"
}

func (ct *CustomTransformation) Transform(data interface{}) (interface{}, error) {
    // Custom transformation logic
    return data, nil
}

// Register custom transformation
engine.RegisterTransformation(&CustomTransformation{})
```

## Data Validation System

### 1. Built-in Validators

#### Required Field Validator
- Ensures critical fields are present
- Configurable required field list
- Custom error messages

#### Data Type Validator
- Validates field formats and types
- Range checking for numeric values
- Pattern matching for strings
- Business rule validation

#### Range Validator
- Price range validation ($0.01 - $500,000)
- Quantity limits (1 - 100,000)
- Date range checking (reasonable business dates)

#### Uniqueness Validator
- Transaction ID uniqueness
- Composite key validation
- Duplicate detection across multiple fields

### 2. Custom Validators

```go
type CustomValidator struct{}

func (cv *CustomValidator) Name() string {
    return "CustomValidator"
}

func (cv *CustomValidator) Validate(data interface{}) error {
    // Custom validation logic
    return nil
}

// Register custom validator
engine.RegisterValidator(&CustomValidator{})
```

## Data Quality Monitoring

### Quality Metrics

The system automatically calculates comprehensive data quality metrics:

```go
type DataQualityMetrics struct {
    Completeness    float64            // % of required fields present
    Consistency     float64            // % of data following patterns
    Validity        float64            // % of data passing validation
    Uniqueness      float64            // % of unique identifiers
    FieldMetrics    map[string]float64 // Per-field quality scores
}
```

### Quality Thresholds

Configure quality thresholds in `configs/data_transformation.yaml`:

```yaml
quality:
  thresholds:
    completeness: 0.95    # 95% completeness required
    validity: 0.90        # 90% validity required
    uniqueness: 0.99      # 99% uniqueness required
    consistency: 0.95     # 95% consistency required
    
  quality_actions:
    warn_threshold: 0.85   # Warn below 85%
    reject_threshold: 0.70 # Reject below 70%
```

### Data Quality Report

```go
type DataQualityReport struct {
    Metrics         DataQualityMetrics
    TotalRecords    int
    Timestamp       time.Time
    Issues          []DataQualityIssue
    Recommendations []string
}

// Generate comprehensive quality report
report := handler.GetDataQualityReport(transactions)
```

## Optimization Strategies

### 1. Built-in Optimizations

#### Duplicate Removal
- Simple ID-based deduplication
- Fast hash-based detection
- Preserves first occurrence

#### Advanced Deduplication  
- Composite key matching
- Multiple field comparison
- Fuzzy matching capabilities

#### Index Optimization
- Data structure optimization
- Query performance improvement
- Memory usage optimization

### 2. Performance Configuration

```yaml
performance:
  batch_size: 10000           # Process in batches
  parallel_processing: true   # Enable parallel processing
  max_workers: 4             # Number of worker goroutines
  memory_limit: "1GB"        # Maximum memory usage
  gc_frequency: 1000         # Garbage collection frequency
```

## Configuration Management

### Complete Configuration Example

```yaml
# configs/data_transformation.yaml
transformation:
  enable_validation: true
  enable_optimization: true
  
  date_formats:
    - "2006-01-02T15:04:05Z"
    - "2006-01-02"
    - "01/02/2006"
    - "02-01-2006"
    
  null_values:
    - ""
    - "NULL"
    - "N/A"
    - "-"
    
  defaults:
    country: "Unknown"
    region: "Unspecified"
    price_multiplier: 100.0
    
  custom_mappings:
    "usa": "United States"
    "uk": "United Kingdom"
    "region_n": "North"
    "region_s": "South"

validation:
  required_fields:
    - "transaction_id"
    - "country"
    - "product_name"
    - "price"
    - "quantity"
    - "transaction_date"
    
  field_rules:
    price:
      min_value: 0.01
      max_value: 500000.00
    quantity:
      min_value: 1
      max_value: 100000

quality:
  thresholds:
    completeness: 0.95
    validity: 0.90
    uniqueness: 0.99
    consistency: 0.95
```

## Usage Examples

### 1. Basic File Processing

```go
// Load configuration
config := LoadTransformationConfig("configs/data_transformation.yaml")

// Create handler
handler := NewFlexibleDataHandler(config)

// Process file (automatic format detection)
transactions, result, err := handler.ProcessDataFile("data/sales_data.csv")
if err != nil {
    log.Fatal("Processing failed:", err)
}

// Display results
fmt.Printf("Processed %d transactions in %v\n", 
    result.TransformedRecords, result.ProcessingTime)
fmt.Printf("Data quality score: %.2f%%\n", 
    result.DataQuality.Completeness * 100)
```

### 2. Stream Processing

```go
// Process data from HTTP request
func processUploadedData(w http.ResponseWriter, r *http.Request) {
    file, header, err := r.FormFile("datafile")
    if err != nil {
        http.Error(w, "Failed to read file", http.StatusBadRequest)
        return
    }
    defer file.Close()
    
    // Detect format from filename
    format := detectFormatFromFilename(header.Filename)
    
    // Process stream
    transactions, result, err := handler.ProcessDataStream(file, format)
    if err != nil {
        http.Error(w, "Processing failed", http.StatusInternalServerError)
        return
    }
    
    // Return results
    json.NewEncoder(w).Encode(map[string]interface{}{
        "transactions": len(transactions),
        "quality": result.DataQuality,
        "processing_time": result.ProcessingTime.String(),
    })
}
```

### 3. Format Conversion

```go
// Convert CSV to JSON
csvFile, _ := os.Open("data.csv")
defer csvFile.Close()

jsonFile, _ := os.Create("data.json")
defer jsonFile.Close()

// Read and transform
transactions, _, err := handler.ProcessDataStream(csvFile, FormatCSV)
if err != nil {
    log.Fatal("Failed to process CSV:", err)
}

// Export to JSON
err = handler.ExportTransactions(transactions, FormatJSON, jsonFile)
if err != nil {
    log.Fatal("Failed to export JSON:", err)
}
```

### 4. Data Quality Analysis

```go
// Load data
transactions, _, _ := handler.ProcessDataFile("data.csv")

// Generate quality report
report := handler.GetDataQualityReport(transactions)

// Display issues
for _, issue := range report.Issues {
    fmt.Printf("%s: %s (%d records)\n", 
        issue.Type, issue.Description, issue.Count)
}

// Display recommendations
for _, rec := range report.Recommendations {
    fmt.Printf("Recommendation: %s\n", rec)
}
```

### 5. Custom Transformation

```go
// Define custom transformation for tax calculation
type TaxCalculator struct {
    config TransformConfig
}

func (tc *TaxCalculator) Name() string {
    return "TaxCalculator"
}

func (tc *TaxCalculator) Description() string {
    return "Calculates tax based on country and product type"
}

func (tc *TaxCalculator) Transform(data interface{}) (interface{}, error) {
    if tx, ok := data.(*models.Transaction); ok {
        // Add tax calculation logic
        taxRate := tc.getTaxRate(tx.Country, tx.ProductName)
        tx.UnitPriceCents = int64(float64(tx.UnitPriceCents) * (1 + taxRate))
        return tx, nil
    }
    return data, nil
}

func (tc *TaxCalculator) getTaxRate(country, product string) float64 {
    // Custom tax logic
    return 0.10 // 10% default tax
}

// Register and use
handler.engine.RegisterTransformation(&TaxCalculator{config: config})
```

## Integration with Existing System

### 1. Updating Ingest Module

The flexible data handler integrates seamlessly with the existing ingest system:

```go
// internal/ingest/ingest.go - Enhanced version
func ParseTransactionsWithFlexibleHandling(r io.Reader) ([]models.Transaction, error) {
    // Load configuration
    config := LoadDefaultTransformationConfig()
    
    // Create handler
    handler := NewFlexibleDataHandler(config)
    
    // Process data with automatic format detection
    transactions, result, err := handler.ProcessDataStream(r, FormatCSV)
    if err != nil {
        return nil, err
    }
    
    // Log processing results
    log.Printf("Processed %d transactions with %.2f%% quality score",
        len(transactions), result.DataQuality.Completeness*100)
    
    return transactions, nil
}
```

### 2. API Endpoints for Data Management

```go
// Add data management endpoints
func (s *Server) setupDataManagementRoutes() {
    s.mux.HandleFunc("POST /api/data/upload", s.handleDataUpload)
    s.mux.HandleFunc("GET /api/data/quality", s.handleDataQuality)
    s.mux.HandleFunc("POST /api/data/transform", s.handleDataTransform)
    s.mux.HandleFunc("GET /api/data/formats", s.handleSupportedFormats)
}

func (s *Server) handleDataUpload(w http.ResponseWriter, r *http.Request) {
    // Handle file upload with flexible format support
    // Apply transformations and return quality metrics
}

func (s *Server) handleDataQuality(w http.ResponseWriter, r *http.Request) {
    // Return current data quality metrics
    report := s.dataHandler.GetDataQualityReport(s.transactions)
    json.NewEncoder(w).Encode(report)
}
```

## Performance Considerations

### 1. Memory Management

```go
// Large file processing with streaming
func ProcessLargeFile(filePath string) error {
    file, _ := os.Open(filePath)
    defer file.Close()
    
    // Process in chunks to manage memory
    const chunkSize = 10000
    
    scanner := bufio.NewScanner(file)
    var chunk []string
    
    for scanner.Scan() {
        chunk = append(chunk, scanner.Text())
        
        if len(chunk) >= chunkSize {
            processChunk(chunk)
            chunk = chunk[:0] // Reset slice
        }
    }
    
    // Process remaining chunk
    if len(chunk) > 0 {
        processChunk(chunk)
    }
    
    return nil
}
```

### 2. Parallel Processing

```go
// Enable parallel processing for large datasets
config := TransformConfig{
    EnableOptimization: true,
}

// Configure performance settings
performance := PerformanceConfig{
    BatchSize:           10000,
    ParallelProcessing:  true,
    MaxWorkers:         runtime.NumCPU(),
    MemoryLimit:        "2GB",
}
```

### 3. Caching Strategy

```go
// Cache transformation results for repeated operations
type CachedDataHandler struct {
    handler *FlexibleDataHandler
    cache   map[string][]models.Transaction
}

func (cdh *CachedDataHandler) ProcessFile(filePath string) ([]models.Transaction, error) {
    // Check cache first
    if cached, exists := cdh.cache[filePath]; exists {
        return cached, nil
    }
    
    // Process and cache
    transactions, _, err := cdh.handler.ProcessDataFile(filePath)
    if err == nil {
        cdh.cache[filePath] = transactions
    }
    
    return transactions, err
}
```

## Error Handling & Logging

### 1. Error Strategies

```yaml
error_handling:
  strategies:
    validation_error: "warn"     # "skip", "warn", "fail"
    transformation_error: "warn" # "skip", "warn", "fail"
    data_quality_error: "warn"   # "skip", "warn", "fail"
  
  max_errors: 1000
  continue_on_error: true
```

### 2. Detailed Logging

```go
// Configure logging levels
monitoring:
  enable_detailed_logging: true
  
  log_levels:
    transformation: "info"
    validation: "warn" 
    optimization: "info"
    quality: "info"
```

### 3. Error Recovery

```go
// Implement error recovery strategies
func (fdh *FlexibleDataHandler) ProcessWithRecovery(filePath string) ([]models.Transaction, error) {
    transactions, result, err := fdh.ProcessDataFile(filePath)
    
    if err != nil {
        // Try alternative processing strategy
        log.Printf("Primary processing failed, trying alternative approach: %v", err)
        
        // Fallback to basic CSV parsing
        return fdh.processWithBasicParser(filePath)
    }
    
    // Check data quality and apply recovery if needed
    if result.DataQuality.Completeness < 0.70 {
        log.Printf("Low data quality detected (%.2f%%), applying recovery", 
            result.DataQuality.Completeness)
        
        return fdh.processWithRelaxedValidation(filePath)
    }
    
    return transactions, nil
}
```

## Testing & Validation

### 1. Unit Tests

```go
func TestDataTransformationEngine(t *testing.T) {
    config := TransformConfig{
        EnableValidation: true,
        PriceMultiplier: 100.0,
    }
    
    engine := NewDataTransformationEngine(config)
    
    // Test transformation
    testData := `transaction_id,country,price
    TX001,USA,$10.50
    TX002,UK,£8.99`
    
    transactions, result, err := engine.TransformCSVData(strings.NewReader(testData))
    
    assert.NoError(t, err)
    assert.Equal(t, 2, len(transactions))
    assert.Equal(t, int64(1050), transactions[0].UnitPriceCents)
}
```

### 2. Integration Tests

```go
func TestCompleteDataProcessing(t *testing.T) {
    // Test with various file formats
    testFiles := []string{
        "testdata/sample.csv",
        "testdata/sample.json", 
        "testdata/sample.yaml",
    }
    
    handler := NewFlexibleDataHandler(LoadTestConfig())
    
    for _, file := range testFiles {
        transactions, result, err := handler.ProcessDataFile(file)
        
        assert.NoError(t, err)
        assert.True(t, len(transactions) > 0)
        assert.True(t, result.DataQuality.Completeness > 0.80)
    }
}
```

### 3. Performance Tests

```go
func BenchmarkDataProcessing(b *testing.B) {
    handler := NewFlexibleDataHandler(LoadDefaultConfig())
    testData := generateTestData(10000) // 10K records
    
    b.ResetTimer()
    
    for i := 0; i < b.N; i++ {
        _, _, err := handler.ProcessDataStream(
            strings.NewReader(testData), FormatCSV)
        if err != nil {
            b.Fatal(err)
        }
    }
}
```

## Deployment Considerations

### 1. Configuration Management

```bash
# Environment-specific configurations
configs/
├── data_transformation.yaml         # Default configuration
├── data_transformation.dev.yaml     # Development overrides
├── data_transformation.prod.yaml    # Production overrides
└── data_transformation.test.yaml    # Testing overrides
```

### 2. Monitoring Integration

```go
// Add Prometheus metrics
var (
    processedRecords = prometheus.NewCounterVec(
        prometheus.CounterOpts{
            Name: "data_processing_records_total",
            Help: "Total number of processed records",
        },
        []string{"format", "status"},
    )
    
    processingDuration = prometheus.NewHistogramVec(
        prometheus.HistogramOpts{
            Name: "data_processing_duration_seconds",
            Help: "Time spent processing data",
        },
        []string{"format"},
    )
    
    dataQualityScore = prometheus.NewGaugeVec(
        prometheus.GaugeOpts{
            Name: "data_quality_score",
            Help: "Current data quality score",
        },
        []string{"metric_type"},
    )
)
```

### 3. Health Checks

```go
// Add health check endpoint for data processing system
func (s *Server) handleDataHealthCheck(w http.ResponseWriter, r *http.Request) {
    status := map[string]interface{}{
        "data_handler_ready": s.dataHandler != nil,
        "supported_formats":  s.dataHandler.GetSupportedFormats(),
        "configuration_valid": s.dataHandler.ValidateConfiguration() == nil,
        "statistics":         s.dataHandler.GetTransformationStatistics(),
    }
    
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(status)
}
```

## Future Enhancements

### 1. Machine Learning Integration

```go
// Planned: ML-based data quality prediction
type MLValidator struct {
    model *tensorflow.SavedModel
}

func (ml *MLValidator) PredictQuality(transaction *models.Transaction) float64 {
    // Use ML model to predict data quality issues
    return 0.95
}
```

### 2. Real-time Stream Processing

```go
// Planned: Kafka/Redis stream processing
func (fdh *FlexibleDataHandler) ProcessStream(stream kafka.Reader) {
    for {
        msg, err := stream.ReadMessage(context.Background())
        if err != nil {
            break
        }
        
        // Process message with flexible handler
        transactions, _, _ := fdh.ProcessDataStream(
            bytes.NewReader(msg.Value), FormatJSON)
        
        // Send to processing pipeline
        fdh.sendToProcessingPipeline(transactions)
    }
}
```

### 3. Advanced Analytics

```go
// Planned: Advanced data profiling and anomaly detection
type DataProfiler struct {
    handler *FlexibleDataHandler
}

func (dp *DataProfiler) ProfileData(transactions []models.Transaction) DataProfile {
    return DataProfile{
        PatternAnalysis:   dp.analyzePatterns(transactions),
        AnomalyDetection: dp.detectAnomalies(transactions),
        TrendAnalysis:    dp.analyzeTrends(transactions),
        Recommendations:  dp.generateRecommendations(transactions),
    }
}
```

---

**This flexible data handling system provides comprehensive capabilities for processing various data formats while maintaining high data quality and performance standards. The modular architecture ensures easy extensibility and customization for specific business requirements.**
