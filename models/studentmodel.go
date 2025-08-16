package models

// Student model
type Student struct {
	ID        int    `json:"id,omitempty"` // No validation, DB handles PK
	FirstName string `json:"first_name,omitempty" validate:"required,min=2,max=50,alpha"`
	LastName  string `json:"last_name,omitempty"  validate:"required,min=2,max=50,alpha"`
	Email     string `json:"email,omitempty"      validate:"required,email"`
	Class     string `json:"class,omitempty"      validate:"omitempty,alphanum"`
}

// CreateStudentValidator validates a Student struct
func CreateStudentValidator(s Student) (string, bool) {
	err := validate.Struct(s)
	return FirstError(err)
}
