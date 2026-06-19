package decoder

import (
	"fmt"
	"log/slog"
	"strings"

	"github.com/swantara-tech/go-base30-to-image/internal/models"
)

// ParseSignature parses jSignature base30 format into a Signature model.
// Format: each stroke is xLeg_yLeg, strokes separated by '_'.
// Within each leg, Z = flip to negative delta, Y = flip to positive delta.
func ParseSignature(base30Data string) (*models.Signature, error) {
	data := base30Data
	lower := strings.ToLower(data)
	if strings.HasPrefix(lower, "image/jsignature;base30,") {
		data = data[len("image/jsignature;base30,"):]
	}

	if len(strings.TrimSpace(data)) == 0 {
		return nil, fmt.Errorf("empty signature data")
	}

	rawStrokes := UncompressStrokes(base30Data)
	if len(rawStrokes) == 0 {
		return nil, fmt.Errorf("no valid strokes found in signature")
	}

	signature := &models.Signature{
		Strokes: make([]models.Stroke, 0, len(rawStrokes)),
	}

	firstPoint := true

	for _, raw := range rawStrokes {
		length := len(raw.X)
		if len(raw.Y) < length {
			length = len(raw.Y)
		}
		if length == 0 {
			continue
		}

		stroke := models.Stroke{
			Points: make([]models.Point, 0, length),
		}

		for i := 0; i < length; i++ {
			pt := models.Point{X: raw.X[i], Y: raw.Y[i]}
			stroke.Points = append(stroke.Points, pt)

			// Update bounding box
			if firstPoint {
				signature.MinX = pt.X
				signature.MinY = pt.Y
				signature.MaxX = pt.X
				signature.MaxY = pt.Y
				firstPoint = false
			} else {
				if pt.X < signature.MinX {
					signature.MinX = pt.X
				}
				if pt.Y < signature.MinY {
					signature.MinY = pt.Y
				}
				if pt.X > signature.MaxX {
					signature.MaxX = pt.X
				}
				if pt.Y > signature.MaxY {
					signature.MaxY = pt.Y
				}
			}
		}

		if len(stroke.Points) > 0 {
			signature.Strokes = append(signature.Strokes, stroke)
		}
	}

	if len(signature.Strokes) == 0 {
		return nil, fmt.Errorf("no valid strokes found in signature")
	}

	signature.Width = signature.MaxX - signature.MinX
	signature.Height = signature.MaxY - signature.MinY

	slog.Debug("Signature parsed",
		"strokes", len(signature.Strokes),
		"width", signature.Width,
		"height", signature.Height,
		"minX", signature.MinX,
		"minY", signature.MinY,
		"maxX", signature.MaxX,
		"maxY", signature.MaxY,
	)

	return signature, nil
}

// ValidateSignature checks if the base30 data is valid
func ValidateSignature(base30Data string) error {
	_, err := ParseSignature(base30Data)
	return err
}
