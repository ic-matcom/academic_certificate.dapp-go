package dto

// User struct
type User struct {
	Username   string `json:"username" validate:"required"`
	Passphrase string `json:"passphrase" validate:"required"`
	FirstName  string `json:"firstname" validate:"required"`
	LastName   string `json:"lastname" validate:"required"`
	Email      string `json:"e-mail" validate:"required,email"`
}

type UserResponse struct {
	Username  string `json:"username" validate:"required"`
	FirstName string `json:"firstname" validate:"required"`
	LastName  string `json:"lastname" validate:"required"`
	Email     string `json:"e-mail" validate:"required,email"`
}
