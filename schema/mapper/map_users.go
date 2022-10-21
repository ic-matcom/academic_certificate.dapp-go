package mapper

import (
	"dapp/schema/dto"
)

func MapUser2UserResponse(user dto.User) dto.UserResponse {
	return dto.UserResponse{
		Username:  user.Username,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Email:     user.Email,
	}
}

func MapUserUpd2User(userID string, user dto.UserUpdateRequest) dto.User {
	return dto.User{
		Username:   user.Username,
		Passphrase: user.Passphrase,
		FirstName:  user.FirstName,
		LastName:   user.LastName,
		Email:      userID,
	}
}
