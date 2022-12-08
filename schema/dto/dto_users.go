package dto

// User struct
type User struct {
	ID         int    `json:"id"`
	Username   string `json:"username" validate:"required"`
	Passphrase string `json:"passphrase" validate:"required"`
	FirstName  string `json:"firstname" validate:"required"`
	LastName   string `json:"lastname" validate:"required"`
	Email      string `json:"email" validate:"required,email"`
	Role       string `json:"rol" validate:"required"`
}

type UserResponse struct {
	ID        int    `json:"id"`
	Username  string `json:"username" validate:"required"`
	FirstName string `json:"firstname" validate:"required"`
	LastName  string `json:"lastname" validate:"required"`
	Email     string `json:"email" validate:"required,email"`
	Role      string `json:"rol" validate:"required"`
}

type UserData struct {
	Username   string `json:"username" validate:"required"`
	Passphrase string `json:"passphrase" validate:"required"`
	FirstName  string `json:"firstname" validate:"required"`
	LastName   string `json:"lastname" validate:"required"`
	Email      string `json:"email" validate:"required,email"`
	Role       string `json:"rol" validate:"required"`
}

type EditUserData struct {
	Username   string `json:"username,omitempty"`
	Passphrase string `json:"passphrase,omitempty"`
	FirstName  string `json:"firstname,omitempty"`
	LastName   string `json:"lastname,omitempty"`
	Email      string `json:"email,omitempty"`
	Role       string `json:"rol,omitempty"`
}
