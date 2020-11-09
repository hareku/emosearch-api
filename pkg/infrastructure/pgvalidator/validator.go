package pgvalidator

import (
	"context"

	"github.com/go-playground/locales/en"
	"github.com/go-playground/locales/ja"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	en_translations "github.com/go-playground/validator/v10/translations/en"
	ja_translations "github.com/go-playground/validator/v10/translations/ja"
	dvalidator "github.com/hareku/emosearch-api/pkg/domain/validator"
)

type errValidation struct {
	errs  validator.ValidationErrors
	trans ut.Translator
}

func (e *errValidation) Error() string {
	return e.errs.Error()
}

func (e *errValidation) Unwrap() error {
	return e.errs
}

func (e *errValidation) ToMap() map[string]string {
	m := map[string]string{}
	for _, err := range e.errs {
		m[err.Field()] = err.Translate(e.trans)
	}
	return m
}

func NewErrValidation(errs validator.ValidationErrors, trans ut.Translator) dvalidator.ErrValidation {
	return &errValidation{errs, trans}
}

type pgValidator struct {
	validate *validator.Validate
	uni      *ut.UniversalTranslator
}

// NewPgValidator creates Validator which is implemented by go-playground validator.
func NewPgValidator() dvalidator.Validator {
	en := en.New()
	uni := ut.New(en, en, ja.New())

	trans, _ := uni.GetTranslator("en")

	validate := validator.New()
	en_translations.RegisterDefaultTranslations(validate, trans)
	ja_translations.RegisterDefaultTranslations(validate, trans)

	return &pgValidator{validate, uni}
}

func (v *pgValidator) StructCtx(ctx context.Context, s interface{}) error {
	err := v.validate.StructCtx(ctx, s)
	if errs, ok := err.(validator.ValidationErrors); ok {
		return NewErrValidation(errs, v.getTranslator(ctx))
	}

	return nil
}

func (v *pgValidator) getTranslator(ctx context.Context) ut.Translator {
	// TODO: fetch from accept-language hader
	trans, _ := v.uni.GetTranslator("en")
	return trans
}
