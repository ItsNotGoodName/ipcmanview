package validate

import "github.com/go-playground/validator/v10"

// use a single instance of Validate, it caches struct info
var Validate *validator.Validate

func init() {
	Validate = validator.New()
	Validate.RegisterValidation("address", validateAddress)
	Validate.RegisterValidation("hostname_or_ip", validateHostnameOrIP)
}

func validateAddress(fl validator.FieldLevel) bool {
	err := Validate.Var(fl.Field().String(), "hostname")
	if err == nil {
		return true
	}

	err = Validate.Var(fl.Field().String(), "hostname_port")
	if err == nil {
		return true
	}

	err = Validate.Var(fl.Field().String(), "ip")
	if err == nil {
		return true
	}

	return false
}

func validateHostnameOrIP(fl validator.FieldLevel) bool {
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
