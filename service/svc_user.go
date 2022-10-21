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

	GetUserSvc(userID string) (dto.UserResponse, *dto.Problem)
	GetUsersSvc() (*[]dto.UserResponse, *dto.Problem)
	PutUserSvc(userID string, request dto.UserUpdateRequest) (dto.UserResponse, *dto.Problem)
	PostUserSvc(user dto.User) (dto.UserResponse, *dto.Problem)
	DeleteUserSvc(userID string) (dto.UserResponse, *dto.Problem)
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

func (s *svcUser) GetUserSvc(userID string) (dto.UserResponse, *dto.Problem) {
	res, err := (*s.repoUser).GetUser(userID)
	if err != nil {
		return dto.UserResponse{}, lib.NewProblem(iris.StatusExpectationFailed, schema.ErrBuntdb, err.Error())
	}
	return mapper.MapUser2UserResponse(res), nil
}

func (s *svcUser) GetUsersSvc() (*[]dto.UserResponse, *dto.Problem) {
	res, err := (*s.repoUser).GetUsers()
	if err != nil {
		return nil, lib.NewProblem(iris.StatusExpectationFailed, schema.ErrBuntdb, err.Error())
	}
	usersResponse := []dto.UserResponse{}
	for i := 0; i < len(res); i++ {
		usersResponse = append(usersResponse, mapper.MapUser2UserResponse(res[i]))
	}
	return &usersResponse, nil
}

func (s *svcUser) PutUserSvc(userID string, request dto.UserUpdateRequest) (dto.UserResponse, *dto.Problem) {
	res, err := s.repoUser.UpdateUser(userID, request)
	if err != nil {
		return dto.UserResponse{}, lib.NewProblem(iris.StatusExpectationFailed, schema.ErrBuntdb, err.Error())
	}
	return mapper.MapUser2UserResponse(res), nil
}

func (s *svcUser) PostUserSvc(user dto.User) (dto.UserResponse, *dto.Problem) {
	passphraseEncoded, _ := lib.Checksum("SHA256", []byte(user.Passphrase))
	user.Passphrase = passphraseEncoded
	res, err := s.repoUser.AddUser(user)
	if err != nil {
		return dto.UserResponse{}, lib.NewProblem(iris.StatusExpectationFailed, schema.ErrBuntdb, err.Error())
	}
	return mapper.MapUser2UserResponse(res), nil
}

func (s *svcUser) DeleteUserSvc(userID string) (dto.UserResponse, *dto.Problem) {
	res, err := s.repoUser.RemoveUser(userID)
	if err != nil {
		return dto.UserResponse{}, lib.NewProblem(iris.StatusExpectationFailed, schema.ErrBuntdb, err.Error())
	}
	return mapper.MapUser2UserResponse(res), nil
}
