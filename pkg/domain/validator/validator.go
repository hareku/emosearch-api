package validator

import (
	"context"
)

// ErrValidation is an error of validation.
type ErrValidation interface {
	Error() string
	Unwrap() error
	ToMap() map[string]string
}

// Validator provides validation methods.
type Validator interface {
	StructCtx(ctx context.Context, s interface{}) error
}
