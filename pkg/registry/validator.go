package registry

import (
	"github.com/hareku/emosearch-api/pkg/domain/validator"
	"github.com/hareku/emosearch-api/pkg/infrastructure/pgvalidator"
)

func (r *registry) NewValidator() validator.Validator {
	return pgvalidator.NewPgValidator()
}
