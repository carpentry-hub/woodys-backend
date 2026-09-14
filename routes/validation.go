// Package routes proporciona los servicios de la api
package routes

import (
	"errors"
	"strings"

	"github.com/go-playground/validator/v10"
)

// validate es la instancia compartida del validador de structs
var validate = validator.New()

// validationErrorMessage convierte los errores de validator en un mensaje simplificado
// del tipo "required fields missing: Title, Owner, Description"
func validationErrorMessage(err error) string {
	var validationErrors validator.ValidationErrors
	if errors.As(err, &validationErrors) {
		fields := make([]string, len(validationErrors))
		for i, e := range validationErrors {
			fields[i] = e.Field()
		}
		return "required fields missing: " + strings.Join(fields, ", ")
	}
	return err.Error()
}
