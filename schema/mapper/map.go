package mapper

import (
	"dapp/schema/dto"
	"encoding/json"
	"strings"
)

// TIP ref https://hellokoding.com/crud-restful-apis-with-go-modules-wire-gin-gorm-and-mysql/
// Is we need it, this method can perform validation and return two values: the mapped struct and the error

// ToAccessTokenDataV region ======== AUTHORIZATION =========================================================
// dto.GrantIntentResponse to dto.AccessTokenData
// TODO: ground the rol idea, according to the DApp app logic
func ToAccessTokenDataV(obj *dto.GrantIntentResponse) *dto.AccessTokenData {
	// claims := dto.Claims{ Sub: obj.Identifier, Rol: "undefined" }
	claims := dto.InjectedParam{Username: obj.Identifier, Role: obj.Identifier}

	return &dto.AccessTokenData{Scope: strings.Fields("dapp.fabric"), Claims: claims}
}

func DecodePayload(payload []byte) interface{} {
	// first attempt to parse for JSON, if not successful then just decode to string
	var structured interface{}
	err := json.Unmarshal(payload, &structured)
	if err != nil {
		return string(payload)
	}
	return structured
}

// endregion =============================================================================
