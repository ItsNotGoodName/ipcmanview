package core

import (
	"context"
	"errors"
	"fmt"

	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	en_translations "github.com/go-playground/validator/v10/translations/en"
)

// ---------- Validator

var (
	validate   *validator.Validate
	translator ut.Translator
)

func init() {
	validate = validator.New()
	validate.RegisterValidation("host", validationHost)

	// Translate
	en := en.New()
	uni := ut.New(en, en)
	translator, _ = uni.GetTranslator("en")
	en_translations.RegisterDefaultTranslations(validate, translator)
}

func validationHost(fl validator.FieldLevel) bool {
	err := validate.Var(fl.Field().String(), "hostname")
	if err == nil {
		return true
	}

	err = validate.Var(fl.Field().String(), "ip")
	if err == nil {
		return true
	}

	return false
}

// ---------- Field

func NewFieldError(field, message string) FieldError {
	return FieldError{
		Field:          field,
		message:        message,
		validatorError: nil,
	}
}

func ToFieldError(field, err validator.FieldError) FieldError {
	return FieldError{
		Field:          err.Field(),
		message:        "",
		validatorError: err,
	}
}

type FieldError struct {
	Field          string
	message        string
	validatorError validator.FieldError
}

func (e FieldError) Message() string {
	if e.validatorError == nil {
		return e.message
	}
	return e.validatorError.Translate(translator)
}

func (e FieldError) Error() string {
	return fmt.Sprintf("'%s': '%s'", e.Field, e.Message())
}

type FieldErrors []FieldError

func (e FieldErrors) Error() string {
	var b []byte
	for i, err := range e {
		if i > 0 {
			b = append(b, '\n')
		}
		b = append(b, err.Error()...)
	}
	return string(b)
}

func IsFieldErrors(err error) bool {
	return errors.Is(err, FieldErrors{}) || errors.Is(err, FieldError{})
}

func AsFieldErrors(err error) (FieldErrors, bool) {
	var fieldErrs FieldErrors
	if ok := errors.As(err, &fieldErrs); ok && len(fieldErrs) > 0 {
		return fieldErrs, ok
	}

	var fieldErr FieldError
	if ok := errors.As(err, &fieldErr); ok {
		return FieldErrors{fieldErr}, ok
	}

	return nil, false
}

// ---------- Validate

func checkValidateError(err error) error {
	if err == nil {
		return err
	}

	errs, ok := err.(validator.ValidationErrors)
	if !ok {
		return err
	}

	fieldErrs := make(FieldErrors, 0, len(errs))
	for _, e := range errs {
		fieldErrs = append(fieldErrs, FieldError{
			Field:          e.Field(),
			message:        "",
			validatorError: e,
		})
	}

	return fieldErrs
}

func ValidateStruct(ctx context.Context, s any) error {
	return checkValidateError(validate.StructCtx(ctx, s))
}

func ValidateStructPartial(ctx context.Context, s any, fields ...string) error {
	return checkValidateError(validate.StructPartialCtx(ctx, s, fields...))
}
