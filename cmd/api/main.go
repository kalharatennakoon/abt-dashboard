package main

import (
	"flag"
	"log"
	"os"

	"abt-dashboard/internal/handlers"
	"abt-dashboard/internal/ingest"
	"abt-dashboard/internal/metrics"
	"abt-dashboard/internal/models"
	"abt-dashboard/internal/server"
	"abt-dashboard/internal/transform"
)

func main() {
	var (
		dataPath      string
		inventoryPath string
		staticDir     string
		addr          string
		configPath    string
		useFlexible   bool
	)

	// Command-line flags for file names
	flag.StringVar(&dataPath, "data", "dataset.csv", "path to transactions data file")
	flag.StringVar(&inventoryPath, "inventory", "dataset.csv", "path to inventory CSV")
	flag.StringVar(&staticDir, "static", "web", "path to static files directory")
	flag.StringVar(&addr, "addr", ":8080", "server listen address")
	flag.StringVar(&configPath, "config", "config/data_transformation.yaml", "path to transformation config")
	flag.BoolVar(&useFlexible, "flexible", true, "use flexible data handling system")
	flag.Parse()

	var transactions []models.Transaction

	if useFlexible {
		// Use flexible data handling system
		log.Printf("Using flexible data handling system with config: %s", configPath)

		// Load transformation configuration
		config, err := transform.LoadTransformationConfigFromPath(configPath)
		if err != nil {
			log.Printf("Failed to load config, using defaults: %v", err)
			config = transform.LoadDefaultTransformationConfig()
		}

		// Create flexible data handler
		dataHandler := transform.NewFlexibleDataHandler(config)

		// Process data file with automatic format detection and transformation
		var result *transform.TransformationResult
		transactions, result, err = dataHandler.ProcessDataFile(dataPath)
		if err != nil {
			log.Fatalf("Failed to process data file: %v", err)
		}

		// Log processing results
		log.Printf("Data processing completed:")
		log.Printf("  - Original records: %d", result.OriginalRecords)
		log.Printf("  - Transformed records: %d", result.TransformedRecords)
		log.Printf("  - Skipped records: %d", result.SkippedRecords)
		log.Printf("  - Processing time: %v", result.ProcessingTime)
		log.Printf("  - Data quality score: %.2f%%", result.DataQuality.Completeness*100)
		log.Printf("  - Transformations applied: %v", result.Transformations)

		// Log warnings and errors if any
		if len(result.Warnings) > 0 {
			log.Printf("Warnings encountered:")
			for _, warning := range result.Warnings {
				log.Printf("  - %s", warning)
			}
		}

		if len(result.Errors) > 0 {
			log.Printf("Errors encountered:")
			for _, error := range result.Errors {
				log.Printf("  - %s", error)
			}
		}

		// Generate and log data quality report
		qualityReport := dataHandler.GetDataQualityReport(transactions)
		if len(qualityReport.Issues) > 0 {
			log.Printf("Data quality issues detected:")
			for _, issue := range qualityReport.Issues {
				log.Printf("  - %s: %s (%d records)", issue.Type, issue.Description, issue.Count)
			}
		}

		if len(qualityReport.Recommendations) > 0 {
			log.Printf("Data quality recommendations:")
			for _, rec := range qualityReport.Recommendations {
				log.Printf("  - %s", rec)
			}
		}
	} else {
		// Use traditional data handling
		log.Printf("Using traditional data handling system")

		// Open dataset CSV
		transReader, err := os.Open(dataPath)
		if err != nil {
			log.Fatalf("failed to open transactions: %v", err)
		}
		defer transReader.Close()

		// Parse CSVs using traditional method
		transactions, err = ingest.ParseTransactionsCSV(transReader)
		if err != nil {
			log.Fatalf("failed to parse transactions: %v", err)
		}
	}

	// Handle inventory data (common for both approaches)
	invMap := map[string]models.Inventory{}
	var invReader *os.File
	invReader, err := os.Open(inventoryPath)
	if err != nil {
		log.Printf("inventory file missing: %v", err)
	}
	if invReader != nil {
		inv, err := ingest.ParseInventoryCSV(invReader)
		if err != nil {
			log.Printf("failed to parse inventory (continuing anyway): %v", err)
		} else {
			invMap = inv
		}
	}

	// Aggregate
	agg := metrics.NewAggregator()
	agg.Ingest(transactions, invMap)

	// Start HTTP server
	api := &handlers.API{Agg: agg}
	srv := server.New(api, staticDir)
	log.Printf("Server listening on %s", addr)
	if err := srv.Listen(addr); err != nil {
		log.Fatal(err)
	}
}
