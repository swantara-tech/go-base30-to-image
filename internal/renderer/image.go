package renderer

import (
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"os"
	"path/filepath"

	"github.com/swantara-tech/go-base30-to-image/internal/models"
)

// RenderSignature renders a signature to an image and saves it to the specified path
func RenderSignature(signature *models.Signature, outputPath string, config models.RenderConfig) error {
	// Ensure output directory exists
	outputDir := filepath.Dir(outputPath)
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return err
	}

	// Calculate canvas size
	canvasWidth, canvasHeight := calculateCanvasSize(signature, config)

	// Create image with white background
	img := image.NewRGBA(image.Rect(0, 0, canvasWidth, canvasHeight))
	white := color.RGBA{255, 255, 255, 255}
	draw.Draw(img, img.Bounds(), &image.Uniform{white}, image.Point{}, draw.Src)

	// Draw signature strokes
	black := color.RGBA{0, 0, 0, 255}
	offsetX := config.Margin - signature.MinX
	offsetY := config.Margin - signature.MinY

	for _, stroke := range signature.Strokes {
		if len(stroke.Points) < 2 {
			continue
		}

		// Draw lines between consecutive points
		for i := 0; i < len(stroke.Points)-1; i++ {
			p1 := stroke.Points[i]
			p2 := stroke.Points[i+1]

			x1 := p1.X + offsetX
			y1 := p1.Y + offsetY
			x2 := p2.X + offsetX
			y2 := p2.Y + offsetY

			drawLine(img, x1, y1, x2, y2, black, config.StrokeWidth)
		}
	}

	// Create output file
	file, err := os.Create(outputPath)
	if err != nil {
		return err
	}
	defer file.Close()

	// Encode based on format
	if config.Format == "jpg" || config.Format == "jpeg" {
		return encodeJPG(file, img, config.Quality)
	}

	// Default to PNG
	return png.Encode(file, img)
}

// calculateCanvasSize determines the optimal canvas size
func calculateCanvasSize(signature *models.Signature, config models.RenderConfig) (int, int) {
	width := config.Width
	height := config.Height

	// Auto-size if not specified
	if width == 0 {
		width = signature.Width + (config.Margin * 2)
	}
	if height == 0 {
		height = signature.Height + (config.Margin * 2)
	}

	return width, height
}

// drawLine draws a line on the image using Bresenham's algorithm
func drawLine(img *image.RGBA, x1, y1, x2, y2 int, c color.Color, thickness float64) {
	dx := abs(x2 - x1)
	dy := abs(y2 - y1)
	
	var sx, sy int
	if x1 < x2 {
		sx = 1
	} else {
		sx = -1
	}
	if y1 < y2 {
		sy = 1
	} else {
		sy = -1
	}

	err := dx - dy

	for {
		// Draw with thickness
		if thickness > 1.0 {
			drawCircle(img, x1, y1, int(thickness), c)
		} else {
			img.Set(x1, y1, c)
		}

		if x1 == x2 && y1 == y2 {
			break
		}

		e2 := 2 * err
		if e2 > -dy {
			err -= dy
			x1 += sx
		}
		if e2 < dx {
			err += dx
			y1 += sy
		}
	}
}

// drawCircle draws a filled circle (for stroke thickness)
func drawCircle(img *image.RGBA, cx, cy, radius int, c color.Color) {
	for y := -radius; y <= radius; y++ {
		for x := -radius; x <= radius; x++ {
			if x*x+y*y <= radius*radius {
				px := cx + x
				py := cy + y
				if px >= 0 && px < img.Bounds().Dx() && py >= 0 && py < img.Bounds().Dy() {
					img.Set(px, py, c)
				}
			}
		}
	}
}

// abs returns the absolute value of an integer
func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}
