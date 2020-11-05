package registry

import (
	"github.com/hareku/emosearch-api/pkg/usecase"
)

func (r *registry) NewSearchUsecase() usecase.SearchUsecase {
	return usecase.NewSearchUsecase(r.NewSearchRepository())
}
