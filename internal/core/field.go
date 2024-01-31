package core

import "fmt"

type FieldError struct {
	Field   string
	Message string
}

func (e FieldError) Error() string {
	return fmt.Sprintf("%s: %s", e.Field, e.Message)
}

type FieldErrors []FieldError

func (e FieldErrors) Error() string {
	return "Multiple field errors."
}
