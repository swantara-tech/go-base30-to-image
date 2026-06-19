package utils

import (
	"encoding/csv"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// ReadFileContent reads the entire content of a file
func ReadFileContent(filePath string) (string, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to read file %s: %w", filePath, err)
	}
	return string(data), nil
}

// ReadCSVFile reads a CSV file and returns batch records
// Expected format: id,signature
func ReadCSVFile(filePath string) ([]BatchRecord, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open CSV file: %w", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	reader.FieldsPerRecord = -1 // Allow variable fields
	reader.LazyQuotes = true
	reader.TrimLeadingSpace = true
	
	// Read all records
	records, err := reader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("failed to read CSV: %w", err)
	}

	if len(records) < 2 {
		return nil, fmt.Errorf("CSV file must have header and at least one data row")
	}

	// Skip header row
	var batchRecords []BatchRecord
	for i, record := range records[1:] {
		if len(record) < 2 {
			return nil, fmt.Errorf("invalid record at row %d: expected 2 columns, got %d", i+2, len(record))
		}

		batchRecords = append(batchRecords, BatchRecord{
			ID:        strings.TrimSpace(record[0]),
			Signature: strings.TrimSpace(record[1]),
		})
	}

	return batchRecords, nil
}

// EnsureDir creates directory if it doesn't exist
func EnsureDir(dirPath string) error {
	return os.MkdirAll(dirPath, 0755)
}

// GetOutputPath generates output file path based on ID and format
func GetOutputPath(outputDir, id, format string) string {
	ext := format
	if ext == "jpeg" {
		ext = "jpg"
	}
	if ext == "" {
		ext = "png"
	}
	
	filename := fmt.Sprintf("%s.%s", id, ext)
	return filepath.Join(outputDir, filename)
}

// ValidateOutputPath checks if the output path is valid
func ValidateOutputPath(path string) error {
	dir := filepath.Dir(path)
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		return fmt.Errorf("output directory does not exist: %s", dir)
	}
	return nil
}

// BatchRecord represents a single record from CSV input
type BatchRecord struct {
	ID        string
	Signature string
}
