package registry

import (
	"github.com/hareku/emosearch-api/pkg/domain/auth"
	"github.com/hareku/emosearch-api/pkg/domain/repository"
)

// Registry provides methods to make instances.
type Registry interface {
	// NewSearchRepository() repository.SearchRepository
	NewUserRepository() repository.UserRepository
	NewAuthenticator() auth.Authenticator
}

type registry struct{}

// NewRegistry returns Registry.
func NewRegistry() Registry {
	return &registry{}
}
