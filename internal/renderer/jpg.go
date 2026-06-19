package renderer

import (
	"image"
	"image/jpeg"
	"os"
)

// encodeJPG encodes an image to JPG format with specified quality
func encodeJPG(file *os.File, img image.Image, quality int) error {
	// Ensure quality is within valid range (1-100)
	if quality < 1 || quality > 100 {
		quality = 95
	}

	opts := &jpeg.Options{
		Quality: quality,
	}

	return jpeg.Encode(file, img, opts)
}
