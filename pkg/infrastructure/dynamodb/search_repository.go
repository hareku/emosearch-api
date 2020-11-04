package dynamodb

import (
	"github.com/hareku/emosearch-api/pkg/domain/model"
	"github.com/hareku/emosearch-api/pkg/domain/repository"
)

type dynamoDBSearchRepository struct{}

// NewDynamoDBSearchRepository creates SearchRepository which is implemented by DynamoDB.
func NewDynamoDBSearchRepository() repository.SearchRepository {
	return &dynamoDBSearchRepository{}
}

func (r *dynamoDBSearchRepository) Create(model.Search) (bool, error) {
	return true, nil
}

func (r *dynamoDBSearchRepository) FindByID(model.SearchID) (model.Search, error) {
	return model.Search{}, nil
}
