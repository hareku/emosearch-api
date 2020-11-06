package comprehend

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go/service/comprehend"
	"github.com/hareku/emosearch-api/pkg/domain/sentiment"
)

type comprehendDetector struct {
	client *comprehend.Comprehend
}

// NewComprehendDetector creates Detector which implemented by AWS Comprehend.
func NewComprehendDetector(client *comprehend.Comprehend) sentiment.Detector {
	return &comprehendDetector{}
}

func (d *comprehendDetector) Detect(ctx context.Context, text string) (*sentiment.Score, error) {
	lang := comprehend.LanguageCodeJa
	output, err := d.client.DetectSentimentWithContext(ctx, &comprehend.DetectSentimentInput{
		LanguageCode: &lang,
		Text:         &text,
	})
	if err != nil {
		return nil, fmt.Errorf("comprehend error: %w", err)
	}

	score := &sentiment.Score{
		Mixed:    output.SentimentScore.Mixed,
		Neutral:  output.SentimentScore.Neutral,
		Negative: output.SentimentScore.Negative,
		Positive: output.SentimentScore.Positive,
	}

	return score, nil
}
