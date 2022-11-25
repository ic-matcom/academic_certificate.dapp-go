package mapper

import (
	"dapp/schema/dto"
	"dapp/schema/models"
)

func MapModelUser2DtoUserResponse(user models.User) dto.UserResponse {
	return dto.UserResponse{
		ID:        user.ID,
		Username:  user.Username,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Email:     user.Email,
		Role:      roleLabel2RoleName(user.Role),
	}
}

func MapDtoUser2DtoUserResponse(user dto.User) dto.UserResponse {
	return dto.UserResponse{
		ID:        user.ID,
		Username:  user.Username,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Email:     user.Email,
		Role:      user.Role,
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
		Role:       roleLabel2RoleName(user.Role),
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
		Role:       user.Role,
	}
}

func MapUserData2ModelUser(userID int, user dto.UserData) models.User {
	return models.User{
		ID:         userID,
		Username:   user.Username,
		Passphrase: user.Passphrase,
		FirstName:  user.FirstName,
		LastName:   user.LastName,
		Email:      user.Email,
		Role:       user.Role,
	}
}

func MapUserData2UserResponse(userID int, user dto.UserData) dto.UserResponse {
	return dto.UserResponse{
		ID:        userID,
		Username:  user.Username,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Email:     user.Email,
		Role:      user.Role,
	}
}

func roleLabel2RoleName(roleLabel string) string {
	if roleLabel == models.Role_Invalid {
		return "Usuario Invalidado"
	}
	if roleLabel == models.Role_SystemAdmin {
		return "Administrador de Sistemas"
	}
	if roleLabel == models.Role_CertificateAdmin {
		return "Administrador de Certificados"
	}
	if roleLabel == models.Role_Secretary {
		return "Secretario General"
	}
	if roleLabel == models.Role_Dean {
		return "Decano de Facultad"
	}
	if roleLabel == models.Role_Rector {
		return "Rector de Universidad"
	}
	return "Rol no encontrado"
}
