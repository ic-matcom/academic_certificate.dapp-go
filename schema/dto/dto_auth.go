package dto

// UserCredIn Is a example declaring the validation for tech struct. It will be used when the
// struct is in the endpoint parameters
type UserCredIn struct {
	Username string `example:"richard.sargon@meinermail.com" validate:"required,ascii,gte=3,lte=60"`
	Password string `example:"password1" validate:"required,ascii,gte=3,lte=20"`
}

type GrantIntentResponse struct {
	Identifier string // if we use `json:"<source_name>"` we can map any source to a common particular / internal struct field as Identifier used here
}

// AccessTokenData using by this REST Api (HLF client node) to grant access to the resources
type AccessTokenData struct {
	Scope  []string
	Claims InjectedParam
}

// Claims user claims
type Claims struct {
	Sub string
	Rol string
}

type InjectedParam struct {
	ID      string
	Username string
}
