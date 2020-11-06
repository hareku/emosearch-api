package main

import (
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/guregu/dynamo"
	"github.com/hareku/emosearch-api/pkg/interfaces/lambda"
	"github.com/hareku/emosearch-api/pkg/registry"
)

func main() {
	awsConf := &aws.Config{
		Region: aws.String("ap-northeast-1"),
	}
	awsEndpoint := os.Getenv("AWS_ENDPOINT")
	if awsEndpoint != "" {
		awsConf.Endpoint = aws.String(awsEndpoint)
	}

	dynamoDB := dynamo.New(session.New(), awsConf)
	table := dynamoDB.Table("EmoSearchAPI")

	registry := registry.NewRegistry(table)
	handler := lambda.NewLambdaHandler(registry)
	handler.Start()
}
