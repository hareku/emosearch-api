package mock

import (
	"context"

	"github.com/hareku/emosearch-api/pkg/domain/sentiment"
)

type mockDetector struct{}

// NewMockDetector creates Detector which is not implemented.
func NewMockDetector() sentiment.Detector {
	return &mockDetector{}
}

func (d *mockDetector) Detect(ctx context.Context, text string) (*sentiment.Score, error) {
	val := float64(0.25)
	score := &sentiment.Score{
		Mixed:    &val,
		Neutral:  &val,
		Negative: &val,
		Positive: &val,
	}

	return score, nil
}
