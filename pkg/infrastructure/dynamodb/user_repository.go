package dynamodb

import (
	"github.com/hareku/emosearch-api/pkg/domain/model"
	"github.com/hareku/emosearch-api/pkg/domain/repository"
)

type dynamoDbUserRepository struct{}

// NewDynamoDatabaseUserRepository creates UserRepository which implemented by DynamoDB.
func NewDynamoDatabaseUserRepository() repository.UserRepository {
	return &dynamoDbUserRepository{}
}

func (r *dynamoDbUserRepository) Store(model.User) (bool, error) {
	return true, nil
}

func (r *dynamoDbUserRepository) FindByID(model.UserID) (model.User, error) {
	return model.User{}, nil
}
