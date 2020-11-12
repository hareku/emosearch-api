package comprehend

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/comprehend"
	"github.com/hareku/emosearch-api/pkg/domain/sentiment"
)

type comprehendDetector struct {
	client *comprehend.Comprehend
}

// NewComprehendDetector creates Detector which implemented by AWS Comprehend.
func NewComprehendDetector(client *comprehend.Comprehend) sentiment.Detector {
	return &comprehendDetector{client}
}

func (d *comprehendDetector) Detect(ctx context.Context, text string) (*sentiment.DetectOutput, error) {
	output, err := d.client.DetectSentimentWithContext(ctx, &comprehend.DetectSentimentInput{
		LanguageCode: aws.String(comprehend.LanguageCodeJa),
		Text:         aws.String(text),
	})
	if err != nil {
		return nil, fmt.Errorf("aws comprehend error: %w", err)
	}

	res := &sentiment.DetectOutput{
		Score: sentiment.Score{
			Mixed:    output.SentimentScore.Mixed,
			Neutral:  output.SentimentScore.Neutral,
			Negative: output.SentimentScore.Negative,
			Positive: output.SentimentScore.Positive,
		},
		Label: determineLabel(output.Sentiment),
	}

	return res, nil
}

func (d *comprehendDetector) BatchDetect(ctx context.Context, textList []*string) ([]sentiment.DetectOutput, error) {
	output, err := d.client.BatchDetectSentimentWithContext(ctx, &comprehend.BatchDetectSentimentInput{
		LanguageCode: aws.String(comprehend.LanguageCodeJa),
		TextList:     textList,
	})
	if err != nil {
		return nil, fmt.Errorf("aws comprehend error: %w", err)
	}

	res := []sentiment.DetectOutput{}
	for _, result := range output.ResultList {
		res = append(res, sentiment.DetectOutput{
			Score: sentiment.Score{
				Mixed:    result.SentimentScore.Mixed,
				Neutral:  result.SentimentScore.Neutral,
				Negative: result.SentimentScore.Negative,
				Positive: result.SentimentScore.Positive,
			},
			Label: determineLabel(result.Sentiment),
		})
	}

	return res, nil
}

func determineLabel(comprehendLabel *string) sentiment.Label {
	switch *comprehendLabel {
	case comprehend.SentimentTypePositive:
		return sentiment.LabelPositive
	case comprehend.SentimentTypeNegative:
		return sentiment.LabelNegative
	case comprehend.SentimentTypeMixed:
		return sentiment.LabelNeutral
	case comprehend.SentimentTypeNeutral:
		return sentiment.LabelNeutral
	default:
		return sentiment.LabelUnknown
	}
}
