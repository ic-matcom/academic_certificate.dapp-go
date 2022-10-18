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

func MapUser2UserResponse(user User) UserResponse {
	return UserResponse{
		Username:  user.Username,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Email:     user.Email,
	}
}

func MapUserUpd2User(userID string, user UserUpdateRequest) User {
	return User{
		Username:   user.Username,
		Passphrase: user.Passphrase,
		FirstName:  user.FirstName,
		LastName:   user.LastName,
		Email:      userID,
	}
}
