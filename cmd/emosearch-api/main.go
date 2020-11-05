package main

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/guregu/dynamo"
	"github.com/hareku/emosearch-api/pkg/interfaces/lambda"
	"github.com/hareku/emosearch-api/pkg/registry"
)

func main() {
	dynamoDB := dynamo.New(session.New(), &aws.Config{Region: aws.String("ap-northeast-1")})
	table := dynamoDB.Table("EmoSearchAPI")

	registry := registry.NewRegistry(table)
	handler := lambda.NewLambdaHandler(registry)
	handler.Start()
}
