package registry

import (
	"github.com/hareku/emosearch-api/pkg/domain/repository"
	"github.com/hareku/emosearch-api/pkg/infrastructure/dynamodb"
)

func (r *registry) NewUserRepository() repository.UserRepository {
	return dynamodb.NewDynamoDatabaseUserRepository()
}
