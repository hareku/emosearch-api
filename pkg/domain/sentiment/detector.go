package sentiment

import (
	"context"
)

// Score is the sentiment score of text.
type Score struct {
	Mixed    *float64
	Negative *float64
	Neutral  *float64
	Positive *float64
}

// Detector provides sentiment detections.
type Detector interface {
	Detect(ctx context.Context, text string) (*Score, error)
}
