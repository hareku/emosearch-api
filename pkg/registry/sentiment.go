package registry

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/comprehend"
	"github.com/hareku/emosearch-api/pkg/domain/sentiment"
	internal_comprehend "github.com/hareku/emosearch-api/pkg/infrastructure/comprehend"
)

var comprehendClient *comprehend.Comprehend

func getComprehendClient() *comprehend.Comprehend {
	if comprehendClient == nil {
		awsConf := &aws.Config{
			Region: aws.String("ap-northeast-1"),
		}
		comprehendClient = comprehend.New(session.Must(session.NewSession()), awsConf)
	}

	return comprehendClient
}

func (r *registry) NewSentimentDetector() sentiment.Detector {
	return internal_comprehend.NewComprehendDetector(getComprehendClient())
}
