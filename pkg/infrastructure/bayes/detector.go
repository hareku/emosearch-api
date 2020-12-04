package bayes

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/hareku/emosearch-api/pkg/domain/sentiment"
)

type bayesDetector struct{}

// NewBayesDetector creates Detector which implemented by Naive Bayes API.
// https://github.com/hareku/sentiment-analysis-api
func NewBayesDetector() sentiment.Detector {
	return &bayesDetector{}
}

type response struct {
	Result []score `json:"result"`
}

type score struct {
	Positive float64
	Negative float64
	Neutral  float64
}

func (s *score) toOutput() sentiment.DetectOutput {
	return sentiment.DetectOutput{
		Score: sentiment.Score{
			Positive: &s.Positive,
			Negative: &s.Negative,
			Neutral:  &s.Neutral,
		},
		Label: s.determineLabel(),
	}
}

func (s *score) determineLabel() sentiment.Label {
	if s.Positive > s.Negative && s.Positive > s.Neutral {
		return sentiment.LabelPositive
	}

	if s.Negative > s.Positive && s.Negative > s.Neutral {
		return sentiment.LabelNegative
	}

	return sentiment.LabelNeutral
}

func (d *bayesDetector) BatchDetect(ctx context.Context, textList []*string) ([]sentiment.DetectOutput, error) {
	reqBody := struct {
		TextList []string
	}{
		TextList: []string{},
	}
	for _, text := range textList {
		reqBody.TextList = append(reqBody.TextList, *text)
	}

	b, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal json request body: %w", err)
	}

	resp, err := http.Post("https://jxbe3mkwui.execute-api.ap-northeast-1.amazonaws.com/Prod/", "application/json", bytes.NewBuffer(b))
	if err != nil {
		return nil, fmt.Errorf("http request error: %w", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("response body reading error: %w", err)
	}

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("status code is %d, response body: %s", resp.StatusCode, body)
	}

	var res response
	err = json.Unmarshal(body, &res)
	if err != nil {
		return nil, fmt.Errorf("json unmarshal error: %w", err)
	}

	output := []sentiment.DetectOutput{}
	for _, score := range res.Result {
		output = append(output, score.toOutput())
	}

	return output, nil
}
