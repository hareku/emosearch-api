package registry

import (
	"github.com/hareku/emosearch-api/pkg/usecase"
)

func (r *registry) NewUserUsecase() usecase.UserUsecase {
	return usecase.NewUserUsecase(r.NewAuthenticator(), r.NewUserRepository())
}

func (r *registry) NewSearchUsecase() usecase.SearchUsecase {
	return usecase.NewSearchUsecase(r.NewAuthenticator(), r.NewSearchRepository())
}
