package core

import "github.com/go-playground/validator/v10"

// use a single instance of Validate, it caches struct info
var validate *validator.Validate

func init() {
	validate = validator.New()
	validate.RegisterValidation("address", validateAddr)
}

// Structure
func validateAddr(fl validator.FieldLevel) bool {
	err := validate.Var(fl.Field().String(), "hostname")
	if err == nil {
		return true
	}

	err = validate.Var(fl.Field().String(), "hostname_port")
	if err == nil {
		return true
	}

	return false
}
