package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "jsign-convert",
	Short: "Convert jSignature Base30 data to PNG/JPG images",
	Long: `A CLI tool to convert jSignature signature data from Base30 format 
to PNG or JPG images. Supports single conversion and batch processing.`,
}

// Execute runs the root command
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func init() {
	// Global flags
	rootCmd.PersistentFlags().BoolP("verbose", "v", false, "Enable verbose logging")
}
