package main

import (
	"flag"
	"fmt"
	"os"

	"chat-assistant-backend/internal/config"
	"chat-assistant-backend/internal/importer"
	"chat-assistant-backend/internal/importer/parsers"
	"chat-assistant-backend/internal/logger"
)

func main() {
	var (
		file     = flag.String("file", "", "Path to the JSON file to import (required)")
		platform = flag.String("platform", "", "Platform type: chatgpt, claude, gemini (required)")
		userID   = flag.String("user-id", "", "User ID to associate with imported data (required)")
		dryRun   = flag.Bool("dry-run", false, "Perform a dry run without writing to database")
		verbose  = flag.Bool("verbose", false, "Enable verbose logging")
	)
	flag.Parse()

	// Validate required flags
	if *file == "" || *platform == "" || *userID == "" {
		fmt.Fprintf(os.Stderr, "Error: --file, --platform, and --user-id are required\n")
		flag.Usage()
		os.Exit(1)
	}

	// Validate file exists
	if _, err := os.Stat(*file); os.IsNotExist(err) {
		fmt.Fprintf(os.Stderr, "Error: file %s does not exist\n", *file)
		os.Exit(1)
	}

	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to load configuration: %v\n", err)
		os.Exit(1)
	}

	// Initialize logger
	logLevel := "info"
	if *verbose {
		logLevel = "debug"
	}

	if err := logger.Init(logLevel, "console", "stdout"); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to initialize logger: %v\n", err)
		os.Exit(1)
	}
	defer logger.Sync()

	// Register all parsers
	parsers.RegisterAll()

	// Execute import
	importerService := importer.NewService(cfg)
	result, err := importerService.Import(*file, *platform, *userID, *dryRun)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Import failed: %v\n", err)
		os.Exit(1)
	}

	// Print results
	printResults(result)
	fmt.Println("Import completed successfully!")
}

func printResults(result *importer.ImportResult) {
	fmt.Printf("\n=== Import Results ===\n")
	fmt.Printf("Platform: %s\n", result.Platform)
	fmt.Printf("Conversations: %d\n", result.ConversationCount)
	fmt.Printf("Messages: %d\n", result.MessageCount)
	fmt.Printf("Success: %d\n", result.SuccessCount)
	fmt.Printf("Errors: %d\n", result.ErrorCount)
	fmt.Printf("Duration: %s\n", result.Duration)

	if len(result.Errors) > 0 {
		fmt.Printf("\n=== Errors ===\n")
		for _, err := range result.Errors {
			fmt.Printf("- %s\n", err)
		}
	}
}
