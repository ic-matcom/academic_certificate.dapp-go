package auth

import (
	"dapp/lib"
	"dapp/repo"
	"dapp/schema"
	"dapp/schema/dto"

	"github.com/kataras/iris/v12"
)

type Provider interface {
	GrantIntent(userCredential *dto.UserCredIn, data interface{}) (*dto.GrantIntentResponse, *dto.Problem)
}

// region ======== EVOTE AUTHENTICATION PROVIDER =========================================

type ProviderDrone struct {
	// walletLocations string
	repo *repo.RepoUser
}

func (p *ProviderDrone) GrantIntent(uCred *dto.UserCredIn, options interface{}) (*dto.GrantIntentResponse, *dto.Problem) {
	// getting the users
	user, err := (*p.repo).GetUserByUsername(uCred.Username)
	if err != nil {
		return nil, lib.NewProblem(iris.StatusExpectationFailed, schema.ErrBuntdb, err.Error())
	}
	checksum, _ := lib.Checksum("SHA256", []byte(uCred.Password))
	if user.Passphrase == checksum {
		return &dto.GrantIntentResponse{Identifier: user.Username, Role: user.Role}, nil
	}

	return nil, lib.NewProblem(iris.StatusUnauthorized, schema.ErrFile, schema.ErrCredsNotFound)
}

// endregion =============================================================================
