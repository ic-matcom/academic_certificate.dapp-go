package utils

import (
	"fmt"

	"dapp/lib"
	"dapp/schema"
	"github.com/tkanos/gonfig"
)

// region ======== TYPES =================================================================

// conf unexported configuration schema holder struct
type conf struct { //nolint:maligned
	// ENVIRONMENT
	Debug    bool
	APIDocIP string
	DappPort string

	// Cryptographic conf
	JWTSignKey string
	TkMaxAge   uint8

	// STORE DB
	StoreDBPath string

	// CRON JOB
	CronEnabled bool
	LogDBPath   string
	EveryTime   int

	// HLF Network & Crypto Materials
	CppPath           string
	WalletFolder      string
	DappIdentityUser  string
	DappIdentityAdmin string
}

// SvcConfig exported configuration service struct
type SvcConfig struct {
	Path string `string:"Path to the config YAML file"`
	conf `conf:"Configuration object"`
}

// endregion =============================================================================

// NewSvcConfig create a new configuration service.
func NewSvcConfig() *SvcConfig {
	c := conf{}

	var configPath = lib.GetEnvOrError(schema.EnvConfigPath)
	var jwtSignKey = "secret__sample__with__32__chars_" // lib.GetEnvOrError(schema.EnvJWTSignKey)

	exist, err := lib.FileExists(configPath)
	if err != nil || !exist {
		panic(fmt.Errorf("server config file not found, check the %s environment variable", schema.EnvConfigPath))
	}

	err = gonfig.GetConf(configPath, &c) // getting the conf
	if err != nil {
		panic(err)
	} // error check

	c.JWTSignKey = jwtSignKey // saving the sign key into the configuration object

	return &SvcConfig{configPath, c} // We are using struct composition here. Hence, the anonymous field (https://golangbot.com/inheritance/)
}
