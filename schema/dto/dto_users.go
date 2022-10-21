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

// UserUpdateRequest struct: For demonstration purposes only
// Email property is missing because in this demo it is the ID and should not be updated
type UserUpdateRequest struct {
	Username   string `json:"username" validate:"required"`
	Passphrase string `json:"passphrase" validate:"required"`
	FirstName  string `json:"firstname" validate:"required"`
	LastName   string `json:"lastname" validate:"required"`
}
