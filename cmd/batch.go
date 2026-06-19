package cmd

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/spf13/cobra"

	"github.com/swantara-tech/go-base30-to-image/internal/decoder"
	"github.com/swantara-tech/go-base30-to-image/internal/models"
	"github.com/swantara-tech/go-base30-to-image/internal/renderer"
	"github.com/swantara-tech/go-base30-to-image/pkg/utils"
)

var (
	csvFile    string
	outputDir  string
	batchFormat string
	batchQuality int
)

var batchCmd = &cobra.Command{
	Use:   "batch",
	Short: "Batch convert signatures from CSV",
	Long:  `Batch convert multiple jSignature Base30 signatures from a CSV file`,
	RunE:  runBatch,
}

func init() {
	rootCmd.AddCommand(batchCmd)

	batchCmd.Flags().StringVarP(&csvFile, "csv", "c", "", "Input CSV file")
	batchCmd.Flags().StringVarP(&outputDir, "output-dir", "d", "./results", "Output directory")
	batchCmd.Flags().StringVar(&batchFormat, "format", "png", "Output format (png or jpg)")
	batchCmd.Flags().IntVar(&batchQuality, "quality", 95, "JPG quality (1-100)")
}

func runBatch(cmd *cobra.Command, args []string) error {
	// Configure logging
	verbose, _ := cmd.Flags().GetBool("verbose")
	if verbose {
		slog.SetDefault(slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
			Level: slog.LevelDebug,
		})))
	}

	if csvFile == "" {
		return fmt.Errorf("--csv flag is required")
	}

	slog.Info("Starting batch conversion",
		"csv", csvFile,
		"output_dir", outputDir,
	)

	// Read CSV file
	records, err := utils.ReadCSVFile(csvFile)
	if err != nil {
		return fmt.Errorf("failed to read CSV: %w", err)
	}

	slog.Info("CSV loaded", "records", len(records))

	// Ensure output directory exists
	if err := utils.EnsureDir(outputDir); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	// Process each record
	successCount := 0
	errorCount := 0

	for i, record := range records {
		slog.Info("Processing record",
			"index", i+1,
			"total", len(records),
			"id", record.ID,
		)

		// Parse signature
		signature, err := decoder.ParseSignature(record.Signature)
		if err != nil {
			slog.Error("Failed to parse signature",
				"id", record.ID,
				"error", err,
			)
			errorCount++
			continue
		}

		// Generate output path
		outputPath := utils.GetOutputPath(outputDir, record.ID, batchFormat)

		// Configure rendering
		config := models.DefaultRenderConfig()
		config.Format = batchFormat
		config.Quality = batchQuality

		// Render signature
		if err := renderer.RenderSignature(signature, outputPath, config); err != nil {
			slog.Error("Failed to render signature",
				"id", record.ID,
				"error", err,
			)
			errorCount++
			continue
		}

		slog.Info("Signature saved",
			"id", record.ID,
			"path", outputPath,
		)
		successCount++
	}

	fmt.Printf("\n✓ Batch conversion complete!\n")
	fmt.Printf("  Success: %d\n", successCount)
	fmt.Printf("  Errors:  %d\n", errorCount)
	fmt.Printf("  Total:   %d\n", len(records))
	fmt.Printf("  Output:  %s\n", outputDir)

	if errorCount > 0 {
		os.Exit(1)
	}

	return nil
}
