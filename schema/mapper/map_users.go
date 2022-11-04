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
