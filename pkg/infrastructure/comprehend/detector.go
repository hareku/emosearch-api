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
		neutral := *result.SentimentScore.Neutral + *result.SentimentScore.Mixed
		res = append(res, sentiment.DetectOutput{
			Score: sentiment.Score{
				Positive: result.SentimentScore.Positive,
				Negative: result.SentimentScore.Negative,
				Neutral:  &neutral,
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
