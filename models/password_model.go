package models

type User struct {
	ID       int    `json:"id,omitempty"`
	Username string `json:"username" validate:"required,min=3,max=50,alphanum"`
	Password string `json:"password,omitempty" validate:"required,min=6"`
}

func (u *User) Validate() error {
	return validate.Struct(u)
}
