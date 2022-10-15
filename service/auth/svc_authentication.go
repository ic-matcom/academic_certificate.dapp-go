package auth

import (
	"dapp/repo"
)

type SvcAuthentication struct {
	AuthProviders map[string]Provider // similar to slices, maps are reference types.
}

// NewSvcAuthentication creates the authentication service. It provides the methods to make the
// authentication intent with the register providers.
//
// - providers [Array] ~ Maps of providers string token / identifiers
//
// - conf [*SvcConfig] ~ App conf instance pointer
func NewSvcAuthentication(providers map[string]bool, repoUser *repo.RepoUser) *SvcAuthentication {
	k := &SvcAuthentication{AuthProviders: make(map[string]Provider)}

	for v := range providers {
		k.AuthProviders[v] = &ProviderDrone{
			repo: repoUser,
		}
	}

	return k
}
