package validator

import (
	"goapi/models"

	"github.com/go-playground/validator/v10"
)

var validate = validator.New()

// ValidateTeacher applies the `validate` tags on models.Teacher.
func ValidateTeacher(t *models.Teacher) error {
	return validate.Struct(t)
}
