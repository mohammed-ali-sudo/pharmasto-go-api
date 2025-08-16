// models/teacher.go
package models

import (
	"fmt"
	"github.com/go-playground/validator/v10"
)

type Teacher struct {
	ID        int    `json:"id,omitempty"` // No validation, DB handles PK
	FirstName string `json:"first_name,omitempty" validate:"required,min=2,max=50,alpha"`
	LastName  string `json:"last_name,omitempty"  validate:"required,min=2,max=50,alpha"`
	Email     string `json:"email,omitempty"      validate:"required,email"`
	Class     string `json:"class,omitempty"      validate:"omitempty,alphanum"`
	Subject   string `json:"subject,omitempty"    validate:"omitempty,alpha"`
}



func CreateTeacherValidator(t Teacher) (string, bool) {
	err := validate.Struct(t)
	return FirstError(err)
}

var validate = validator.New()

func FirstError(err error) (string, bool) {
	if err == nil {
		return "", true
	}
	verrs, ok := err.(validator.ValidationErrors)
	if !ok || len(verrs) == 0 {
		return err.Error(), false
	}
	fe := verrs[0] // أول خطأ فقط

	switch fe.Tag() {
	case "required":
		return fmt.Sprintf("%s is required", fe.Field()), false
	case "min":
		return fmt.Sprintf("%s must be at least %s characters", fe.Field(), fe.Param()), false
	case "max":
		return fmt.Sprintf("%s must be at most %s characters", fe.Field()), false
	case "email":
		return "Invalid email format", false
	case "alpha":
		return fmt.Sprintf("%s must contain only letters (no spaces)", fe.Field()), false
	case "alphanum":
		return fmt.Sprintf("%s must be alphanumeric (no spaces)", fe.Field()), false
	case "alpha_space":
		return fmt.Sprintf("%s must contain only letters and spaces", fe.Field()), false
	case "alphanum_space":
		return fmt.Sprintf("%s must be letters, numbers, and spaces", fe.Field()), false
	case "required_with":
		return fmt.Sprintf("%s is required when %s is present", fe.Field(), fe.Param()), false
	case "required_without":
		return fmt.Sprintf("%s is required when %s is missing", fe.Field(), fe.Param()), false
	case "email_domain":
		return "Email must be a @school.edu address", false
	default:
		return fmt.Sprintf("%s is invalid", fe.Field()), false
	}
}