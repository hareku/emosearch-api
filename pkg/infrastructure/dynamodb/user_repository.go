package dynamodb

import (
	"github.com/hareku/emosearch-api/pkg/domain"
	"github.com/hareku/emosearch-api/pkg/repository"
)

type dynamoDbUserRepository struct{}

// NewDynamoDatabaseUserRepository creates UserRepository which implemented by DynamoDB.
func NewDynamoDatabaseUserRepository() repository.UserRepository {
	return &dynamoDbUserRepository{}
}

func (r *dynamoDbUserRepository) Store(domain.User) (bool, error) {
	return true, nil
}

func (r *dynamoDbUserRepository) FindByID(domain.UserID) (domain.User, error) {
	return domain.User{}, nil
}
