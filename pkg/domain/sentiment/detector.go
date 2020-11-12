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

// Label represents a label type of sentiment scores.
type Label string

const (
	// LabelPositive labels a sentiment score as positive.
	LabelPositive = Label("POSITIVE")

	// LabelNegative labels a sentiment score as negative.
	LabelNegative = Label("NEGATIVE")

	// LabelNeutral labels a sentiment score as neutral.
	LabelNeutral = Label("NEUTRAL")

	// LabelUnknown labels a sentiment score as unknown.
	LabelUnknown = Label("UNKNOWN")

	// LabelUndetected labels a sentiment score as undetected.
	LabelUndetected = Label("UNDETECTED")
)

// DetectOutput is the type of Detector.Detect method.
type DetectOutput struct {
	Score Score
	Label Label
}

// Detector provides sentiment detections.
type Detector interface {
	Detect(ctx context.Context, text string) (*DetectOutput, error)
}
