package models

type UserRequest struct {
	Username string `json:"username" validate:"required,min=3" description:"Unique username"`
	Email    string `json:"email" validate:"required,email" description:"Valid email address"`
	Age      int    `json:"age,omitempty" validate:"min=18" description:"User age (must be 18+)"`
}
