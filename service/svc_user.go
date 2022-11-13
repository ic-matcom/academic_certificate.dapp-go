package service

import (
	"dapp/lib"
	"dapp/repo"
	"dapp/schema"
	"dapp/schema/dto"
	"dapp/schema/mapper"

	"github.com/kataras/iris/v12"
)

// region ======== SETUP =================================================================

// ISvcUser User request service interface
type ISvcUser interface {
	// user functions

	GetUserSvc(userID int) (dto.UserResponse, *dto.Problem)
	GetUsersSvc() (*[]dto.UserResponse, *dto.Problem)
	PutUserSvc(userID int, user dto.UserUpdate) (dto.UserResponse, *dto.Problem)
	PostUserSvc(user dto.User) (dto.UserResponse, *dto.Problem)
	DeleteUserSvc(userID int) (dto.UserResponse, *dto.Problem)
}

type svcUser struct {
	repoUser *repo.RepoUser
}

// endregion =============================================================================

// NewSvcUserReqs instantiate the User request services
func NewSvcUserReqs(repoUser *repo.RepoUser) ISvcUser {
	return &svcUser{repoUser}
}

// region ======== METHODS ======================================================

func (s *svcUser) GetUserSvc(userID int) (dto.UserResponse, *dto.Problem) {
	res, err := (*s.repoUser).GetUser(userID)
	if err != nil {
		return dto.UserResponse{}, lib.NewProblem(iris.StatusExpectationFailed, schema.ErrBuntdb, err.Error())
	}
	return mapper.MapDtoUser2DtoUserResponse(res), nil
}

func (s *svcUser) GetUsersSvc() (*[]dto.UserResponse, *dto.Problem) {
	res, err := (*s.repoUser).GetUsers()
	if err != nil {
		return nil, lib.NewProblem(iris.StatusExpectationFailed, schema.ErrBuntdb, err.Error())
	}
	var usersResponse []dto.UserResponse
	for i := 0; i < len(res); i++ {
		usersResponse = append(usersResponse, mapper.MapDtoUser2DtoUserResponse(res[i]))
	}
	return &usersResponse, nil
}

func (s *svcUser) PutUserSvc(userID int, user dto.UserUpdate) (dto.UserResponse, *dto.Problem) {
	passphraseEncoded, _ := lib.Checksum("SHA256", []byte(user.Passphrase))
	user.Passphrase = passphraseEncoded
	dtoUser, err := s.repoUser.UpdateUser(userID, user)
	if err != nil {
		return dto.UserResponse{}, lib.NewProblem(iris.StatusExpectationFailed, schema.ErrBuntdb, err.Error())
	}
	return mapper.MapDtoUser2DtoUserResponse(dtoUser), nil
}

func (s *svcUser) PostUserSvc(user dto.User) (dto.UserResponse, *dto.Problem) {
	passphraseEncoded, _ := lib.Checksum("SHA256", []byte(user.Passphrase))
	user.Passphrase = passphraseEncoded
	_, err := s.repoUser.AddUser(user)
	if err != nil {
		return dto.UserResponse{}, lib.NewProblem(iris.StatusExpectationFailed, schema.ErrBuntdb, err.Error())
	}
	return mapper.MapDtoUser2DtoUserResponse(user), nil
}

func (s *svcUser) DeleteUserSvc(userID int) (dto.UserResponse, *dto.Problem) {
	user, err := s.repoUser.RemoveUser(userID)
	if err != nil {
		return dto.UserResponse{}, lib.NewProblem(iris.StatusExpectationFailed, schema.ErrBuntdb, err.Error())
	}
	return mapper.MapDtoUser2DtoUserResponse(user), nil
}
