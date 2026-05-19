package models

type RegisterRequest struct {
	Fullname string `json:"full_name" validate:"required" description:"User full legal name"`
	Email    string `json:"email" validate:"required,email" description:"Registered email"`
	Password string `json:"password" validate:"required,min=8" description:"Secure password"`
}
