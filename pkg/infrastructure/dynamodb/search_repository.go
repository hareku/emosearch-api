package dynamodb

import (
	"github.com/hareku/emosearch-api/pkg/domain"
	"github.com/hareku/emosearch-api/pkg/repository"
)

type dynamoDBSearchRepository struct{}

// NewDynamoDBSearchRepository creates SearchRepository which is implemented by DynamoDB.
func NewDynamoDBSearchRepository() repository.SearchRepository {
	return &dynamoDBSearchRepository{}
}

func (r *dynamoDBSearchRepository) Create(domain.Search) (bool, error) {
	return true, nil
}

func (r *dynamoDBSearchRepository) FindByID(domain.SearchID) (domain.Search, error) {
	return domain.Search{}, nil
}
