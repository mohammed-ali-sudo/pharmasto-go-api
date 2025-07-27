// models/teacher.go
package models

type Teacher struct {
	ID        int    `json:"id,omitempty"`
	FirstName string `json:"first_name,omitempty" validate:"required,alpha" msg:"الاسم الأول مطلوب ويجب أن يحتوي على حروف فقط"`
	LastName  string `json:"last_name,omitempty"  validate:"required,alpha" msg:"اسم العائلة مطلوب ويجب أن يحتوي على حروف فقط"`
	Email     string `json:"email,omitempty"      validate:"required,email" msg:"البريد الإلكتروني غير صالح"`
	Class     string `json:"class,omitempty"      validate:"required" msg:"الفصل مطلوب"`
	Subject   string `json:"subject,omitempty"    validate:"required" msg:"المادة مطلوبة"`
}
