package renderer

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/swantara-tech/go-base30-to-image/internal/models"
)

func TestRenderSignature(t *testing.T) {
	// Create a simple test signature
	signature := &models.Signature{
		Strokes: []models.Stroke{
			{
				Points: []models.Point{
					{X: 10, Y: 10},
					{X: 20, Y: 20},
					{X: 30, Y: 10},
				},
			},
		},
		MinX:   10,
		MinY:   10,
		MaxX:   30,
		MaxY:   20,
		Width:  20,
		Height: 10,
	}

	// Create temporary directory
	tempDir, err := os.MkdirTemp("", "jsign-test-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	tests := []struct {
		name   string
		format string
		ext    string
	}{
		{
			name:   "PNG output",
			format: "png",
			ext:    "png",
		},
		{
			name:   "JPG output",
			format: "jpg",
			ext:    "jpg",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			outputPath := filepath.Join(tempDir, "test."+tt.ext)

			config := models.DefaultRenderConfig()
			config.Format = tt.format
			config.Width = 100
			config.Height = 50

			err := RenderSignature(signature, outputPath, config)
			if err != nil {
				t.Errorf("RenderSignature failed: %v", err)
			}

			// Check file exists
			if _, err := os.Stat(outputPath); os.IsNotExist(err) {
				t.Errorf("output file was not created: %s", outputPath)
			}
		})
	}
}

func TestCalculateCanvasSize(t *testing.T) {
	signature := &models.Signature{
		Width:  200,
		Height: 100,
	}

	tests := []struct {
		name         string
		config       models.RenderConfig
		expectWidth  int
		expectHeight int
	}{
		{
			name: "auto size",
			config: models.RenderConfig{
				Width:  0,
				Height: 0,
				Margin: 20,
			},
			expectWidth:  240, // 200 + 2*20
			expectHeight: 140, // 100 + 2*20
		},
		{
			name: "fixed size",
			config: models.RenderConfig{
				Width:  800,
				Height: 300,
				Margin: 20,
			},
			expectWidth:  800,
			expectHeight: 300,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			width, height := calculateCanvasSize(signature, tt.config)
			if width != tt.expectWidth {
				t.Errorf("expected width %d, got %d", tt.expectWidth, width)
			}
			if height != tt.expectHeight {
				t.Errorf("expected height %d, got %d", tt.expectHeight, height)
			}
		})
	}
}

func TestAbs(t *testing.T) {
	tests := []struct {
		name   string
		input  int
		expect int
	}{
		{"positive", 5, 5},
		{"negative", -5, 5},
		{"zero", 0, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := abs(tt.input)
			if result != tt.expect {
				t.Errorf("abs(%d) = %d, want %d", tt.input, result, tt.expect)
			}
		})
	}
}
