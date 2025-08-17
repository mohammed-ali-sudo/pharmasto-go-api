// models/teacher.go
package models

import (
	"fmt"
	"reflect"

	"github.com/go-playground/validator/v10"
)

type Teacher struct {
	ID        int    `json:"id,omitempty"        db:"id"`
	FirstName string `json:"first_name,omitempty" db:"first_name" validate:"required,min=2,max=50,alpha"`
	LastName  string `json:"last_name,omitempty"  db:"last_name"  validate:"required,min=2,max=50,alpha"`
	Email     string `json:"email,omitempty"      db:"email"      validate:"required,email"`
	Class     string `json:"class,omitempty"      db:"class"      validate:"omitempty,alphanum"`
	Subject   string `json:"subject,omitempty"    db:"subject"    validate:"omitempty,alpha"`
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

func CreateTeacherValidator(t Teacher) (string, bool) {
	err := validate.Struct(t)
	return FirstError(err)
}

func generateInsertQuery(model any) string {
	modelType := reflect.TypeOf(model)
	if modelType.Kind() == reflect.Ptr {
		modelType = modelType.Elem()
	}

	var columns, placeholders string
	for i := 0; i < modelType.NumField(); i++ {
		field := modelType.Field(i)
		dbTag := field.Tag.Get("db") // assuming you use `db:"column_name"`
		if dbTag == "" {
			dbTag = field.Name
		}
		if i > 0 {
			columns += ", "
			placeholders += ", "
		}
		columns += dbTag
		placeholders += fmt.Sprintf("$%d", i+1) // PostgreSQL-style placeholders
	}

	return fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s)", modelType.Name(), columns, placeholders)
}
