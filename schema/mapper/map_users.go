package mapper

import (
	"dapp/schema/dto"
	"dapp/schema/models"
)

func MapDtoUser2DtoUserResponse(user dto.User) dto.UserResponse {
	return dto.UserResponse{
		ID:        user.ID,
		Username:  user.Username,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Email:     user.Email,
	}
}

func MapModelUser2DtoUser(user models.User) dto.User {
	return dto.User{
		ID:         user.ID,
		Username:   user.Username,
		Passphrase: user.Passphrase,
		FirstName:  user.FirstName,
		LastName:   user.LastName,
		Email:      user.Email,
	}
}

func MapDtoUser2ModelUser(user dto.User) models.User {
	return models.User{
		ID:         user.ID,
		Username:   user.Username,
		Passphrase: user.Passphrase,
		FirstName:  user.FirstName,
		LastName:   user.LastName,
		Email:      user.Email,
	}
}

func MapUserUpd2ModelUser(userID int, user dto.UserUpdate) models.User {
	return models.User{
		ID:         userID,
		Username:   user.Username,
		Passphrase: user.Passphrase,
		FirstName:  user.FirstName,
		LastName:   user.LastName,
		Email:      user.Email,
	}
}

func MapUserUpd2UserResponse(userID int, user dto.UserUpdate) dto.UserResponse {
	return dto.UserResponse{
		ID:        userID,
		Username:  user.Username,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Email:     user.Email,
	}
}
