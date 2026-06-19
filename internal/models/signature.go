package models

// Point represents a single coordinate point in a signature stroke
type Point struct {
	X int
	Y int
}

// Stroke represents a continuous line drawn in the signature
type Stroke struct {
	Points []Point
}

// Signature represents a complete parsed signature with all strokes
type Signature struct {
	Strokes    []Stroke
	MinX       int
	MinY       int
	MaxX       int
	MaxY       int
	Width      int
	Height     int
}

// RenderConfig holds configuration for rendering the signature image
type RenderConfig struct {
	Width       int
	Height      int
	Margin      int
	StrokeWidth float64
	Format      string // "png" or "jpg"
	Quality     int    // for JPG only (1-100)
}

// DefaultRenderConfig returns a default rendering configuration
func DefaultRenderConfig() RenderConfig {
	return RenderConfig{
		Width:       0, // 0 means auto-size
		Height:      0,
		Margin:      20,
		StrokeWidth: 2.0,
		Format:      "png",
		Quality:     95,
	}
}

// BatchRecord represents a single record from CSV input
type BatchRecord struct {
	ID        string
	Signature string
}
