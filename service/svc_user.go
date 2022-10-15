package service

import (
	"dapp/lib"
	"dapp/repo"
	"dapp/schema"
	"dapp/schema/dto"
	"github.com/kataras/iris/v12"
)

// region ======== SETUP =================================================================

// ISvcUser User request service interface
type ISvcUser interface {
	// user functions

	GetUserSvc(userID string) (dto.User, *dto.Problem)
	GetUsersSvc() (*[]any, *dto.Problem)
	PutUserSvc(userID string, request dto.UserUpdateRequest) (any, *dto.Problem)
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

func (s *svcUser) GetUserSvc(userID string) (dto.User, *dto.Problem) {
	res, err := (*s.repoUser).GetUser(userID)
	if err != nil {
		return dto.User{}, lib.NewProblem(iris.StatusExpectationFailed, schema.ErrBuntdb, err.Error())
	}
	return res, nil
}

func (s *svcUser) GetUsersSvc() (*[]any, *dto.Problem) {
	res, err := (*s.repoUser).GetUsers()
	if err != nil {
		return nil, lib.NewProblem(iris.StatusExpectationFailed, schema.ErrBuntdb, err.Error())
	}
	return &res, nil
}

func (s *svcUser) PutUserSvc(userID string, request dto.UserUpdateRequest) (any, *dto.Problem) {
	if _, exists := repo.UsersById[userID]; exists {
		repo.UsersById[userID] = dto.User{
			Username:   request.Username,
			Passphrase: request.Passphrase,
			FirstName:  request.FirstName,
			LastName:   request.LastName,
			Email:      userID,
		}

		return repo.UsersById[userID], nil
	}
	return dto.UserResponse{}, lib.NewProblem(iris.StatusExpectationFailed, schema.ErrBuntdb, "user not exist")
}
