package registry

import (
	"github.com/hareku/emosearch-api/pkg/domain/auth"
	"github.com/hareku/emosearch-api/pkg/domain/repository"
	"github.com/hareku/emosearch-api/pkg/usecase"
)

// Registry provides methods to make instances.
type Registry interface {
	NewAuthenticator() auth.Authenticator
	NewUserRepository() repository.UserRepository
	NewSearchRepository() repository.SearchRepository
	NewSearchUsecase() usecase.SearchUsecase
}

type registry struct{}

// NewRegistry returns Registry.
func NewRegistry() Registry {
	return &registry{}
}
