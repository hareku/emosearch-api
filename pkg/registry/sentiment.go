package registry

import (
	"github.com/hareku/emosearch-api/pkg/domain/sentiment"
	"github.com/hareku/emosearch-api/pkg/infrastructure/bayes"
)

func (r *registry) NewSentimentDetector() sentiment.Detector {
	return bayes.NewBayesDetector()
}
