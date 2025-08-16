package models

type Manfactory struct {
	ID                string `json:"id,omitempty"` // UUID
	ManfactoryName    string `json:"manfactory_name" validate:"required,min=2,max=50,alpha"`
	ManfactoryCountry string `json:"manfactory_country" validate:"required,min=2,max=50,alpha"`
	ContactEmail      string `json:"email" validate:"required,email"`
	ContactNumber     string `json:"phone_number,omitempty" validate:"omitempty,alphanum"`
	LicenseNumber     string `json:"license_number,omitempty" validate:"omitempty,alpha"`
}

func CreateManfactoryValidator(m Manfactory) (string, bool) {
	err := validate.Struct(m)
	return FirstError(err)
}
