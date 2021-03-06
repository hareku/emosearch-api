package registry

import (
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/guregu/dynamo"
	"github.com/hareku/emosearch-api/pkg/domain/repository"
	"github.com/hareku/emosearch-api/pkg/infrastructure/dynamodb"
)

var dynamoTable *dynamo.Table

func getDynamoTable() *dynamo.Table {
	if dynamoTable == nil {
		awsConf := &aws.Config{
			Region: aws.String("ap-northeast-1"),
		}

		if region := os.Getenv("AWS_REGION"); region != "" {
			awsConf.Region = aws.String(region)
		}

		if endpoint := os.Getenv("AWS_ENDPOINT"); endpoint != "" {
			awsConf.Endpoint = aws.String(endpoint)
		}

		dynamoDB := dynamo.New(session.New(), awsConf)
		table := dynamoDB.Table("EmoSearchAPI")
		dynamoTable = &table
	}

	return dynamoTable
}

func (r *registry) NewUserRepository() repository.UserRepository {
	return dynamodb.NewDynamoDatabaseUserRepository(*getDynamoTable())
}

func (r *registry) NewSearchRepository() repository.SearchRepository {
	return dynamodb.NewDynamoDBSearchRepository(*getDynamoTable())
}

func (r *registry) NewTweetRepository() repository.TweetRepository {
	return dynamodb.NewDynamoDBTweetRepository(*getDynamoTable())
}
