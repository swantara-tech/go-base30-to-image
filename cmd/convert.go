package cmd

import (
	"fmt"
	"log/slog"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"github.com/swantara-tech/go-base30-to-image/internal/decoder"
	"github.com/swantara-tech/go-base30-to-image/internal/models"
	"github.com/swantara-tech/go-base30-to-image/internal/renderer"
	"github.com/swantara-tech/go-base30-to-image/pkg/utils"
)

var (
	inputFile  string
	outputFile string
	base30Str  string
	format     string
	quality    int
	width      int
	height     int
)

var convertCmd = &cobra.Command{
	Use:   "convert",
	Short: "Convert a single signature to image",
	Long:  `Convert a single jSignature Base30 signature to PNG or JPG image`,
	RunE:  runConvert,
}

func init() {
	rootCmd.AddCommand(convertCmd)

	convertCmd.Flags().StringVarP(&inputFile, "input", "i", "", "Input file containing signature data")
	convertCmd.Flags().StringVarP(&outputFile, "output", "o", "signature.png", "Output file path")
	convertCmd.Flags().StringVar(&base30Str, "base30", "", "Base30 signature string")
	convertCmd.Flags().StringVar(&format, "format", "png", "Output format (png or jpg)")
	convertCmd.Flags().IntVar(&quality, "quality", 95, "JPG quality (1-100)")
	convertCmd.Flags().IntVar(&width, "width", 0, "Output width (0 for auto)")
	convertCmd.Flags().IntVar(&height, "height", 0, "Output height (0 for auto)")
}

func runConvert(cmd *cobra.Command, args []string) error {
	// Configure logging
	verbose, _ := cmd.Flags().GetBool("verbose")
	if verbose {
		slog.SetDefault(slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
			Level: slog.LevelDebug,
		})))
	}

	// Get signature data
	var signatureData string
	var err error

	if base30Str != "" {
		signatureData = base30Str
	} else if inputFile != "" {
		signatureData, err = utils.ReadFileContent(inputFile)
		if err != nil {
			return err
		}
		signatureData = strings.TrimSpace(signatureData)
	} else {
		return fmt.Errorf("either --input or --base30 must be specified")
	}

	if signatureData == "" {
		return fmt.Errorf("signature data is empty")
	}

	slog.Info("Converting signature",
		"format", format,
		"output", outputFile,
	)

	// Parse signature
	signature, err := decoder.ParseSignature(signatureData)
	if err != nil {
		return fmt.Errorf("failed to parse signature: %w", err)
	}

	slog.Info("Signature parsed",
		"strokes", len(signature.Strokes),
		"width", signature.Width,
		"height", signature.Height,
	)

	// Configure rendering
	config := models.DefaultRenderConfig()
	config.Format = format
	config.Quality = quality
	config.Width = width
	config.Height = height

	// Render signature
	if err := renderer.RenderSignature(signature, outputFile, config); err != nil {
		return fmt.Errorf("failed to render signature: %w", err)
	}

	fmt.Printf("✓ Signature saved to: %s\n", outputFile)
	return nil
}
