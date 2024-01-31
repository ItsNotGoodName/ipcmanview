package core

import (
	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	en_translations "github.com/go-playground/validator/v10/translations/en"
)

var (
	Validate   *validator.Validate
	Translator ut.Translator
)

func init() {
	Validate = validator.New()
	Validate.RegisterValidation("host", validateHost)

	// Translate
	en := en.New()
	uni := ut.New(en, en)
	Translator, _ = uni.GetTranslator("en")
	en_translations.RegisterDefaultTranslations(Validate, Translator)
}

func validateHost(fl validator.FieldLevel) bool {
	err := Validate.Var(fl.Field().String(), "hostname")
	if err == nil {
		return true
	}

	err = Validate.Var(fl.Field().String(), "ip")
	if err == nil {
		return true
	}

	return false
}
