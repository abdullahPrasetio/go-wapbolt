package models

type UserV2Request struct {
	Fullname string `json:"full_name" validate:"required" description:"User full legal name"`
	Email    string `json:"email" validate:"required,email" description:"Contact email"`
	Phone    string `json:"phone,omitempty" validate:"numeric" description:"Mobile phone number"`
}
