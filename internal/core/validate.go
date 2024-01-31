package core

import (
	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	en_translations "github.com/go-playground/validator/v10/translations/en"
)

// use a single instance of Validate, it caches struct info
var (
	uni        *ut.UniversalTranslator
	Validate   *validator.Validate
	Translator ut.Translator
)

func init() {
	en := en.New()
	uni = ut.New(en, en)

	// this is usually know or extracted from http 'Accept-Language' header
	// also see uni.FindTranslator(...)
	Translator, _ = uni.GetTranslator("en")

	Validate = validator.New()
	en_translations.RegisterDefaultTranslations(Validate, Translator)
	Validate.RegisterValidation("host", validateHost)
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
