package schema

// region ======== i18n ERROR KEYS =======================================================
const (
	ErrAuth                              = "err.authentication"
	ErrGeneric                           = "err.generic"
	ErrInvalidEnvVar                     = "err.invalid.environment.var"
	ErrRepositoryOps                     = "err.repo_ops"
	ErrNotFound                          = "err.not_found"
	ErrHttpResError                      = "err.http_response"
	ErrDuplicateKey                      = "err.duplicate_key"
	ErrInvalidType                       = "err.wrong_type_assertion"
	ErrNetwork                           = "err.network"
	ErrBadGateway                        = "err.bad_gateway"
	ErrJsonParse                         = "err.json_parse"
	ErrProcParam                         = "err.processing_param"
	ErrJwtGen                            = "err.jwt_generation"
	ErrWrongAuthProvider                 = "err.wrong_auth_provider"
	ErrUnauthorized                      = "err.unauthorized"
	ErrFileProc                          = "err.processing_file"
	ErrFile                              = "err.system_file_related"
	ErrBuntdbItemNotFound                = "err.database_related.item_not_found"
	ErrBuntdb                            = "err.database_related"
	ErrBuntdbPopulated                   = "err.database_populated"
	ErrBuntdbNotPopulated                = "err.database_not_populated"
	ErrDroneMaximumLoadWeightExceededKey = "err.drone_maximum_load_weight_exceeded"
	ErrDroneVeryLowBatteryKey            = "err.drone_very_low_battery"
	ErrDroneBusyKey                      = "err.drone_busy"
	ErrBuntdbIndex                       = "err.database_index_related"
	ErrStorageProc                       = "err.storage_service_processing"
	ErrVal                               = "err.invalid_data"
	ErrBlockchainTxs                     = "err.blockchain_tx"
	ErrUnmarshalBcTxsResponse            = "err.unmarshal_bc_txs_response"
	ErrCryptProc                         = "err.crypt_material_processing"
	ErrCryptProcMissing                  = "err.crypt_material_processing.missing_files"
	ErrParamURL                          = "err.query_parameter"
	ErrValidationField                   = "err.validation_field"
)

// endregion =============================================================================

// region ======== ERROR DETAILS =========================================================
const (
	ErrCredsNotFound       = "The provided credentials don't seems to be valid"
	ErrDetNotFound         = "resource not found"
	ErrDetContractNotFound = "contract function not found"
	ErrDetHttpResError     = "there is an error on http request"
	ErrDetInvalidType      = "invalid interface type (type assertion)"
	ErrDetInvalidCred      = "something was wrong with the provided user credentials"
	ErrDetInvalidProvider  = "wrong or invalid provider"
	ErrDetInvalidFile      = "the given file seems suspicious"
	ErrDetInvalidField     = "the given field is invalid"
	ErrDetWalletProc       = "failed to create wallet"
	ErrEmailProc           = "failed to send email"
	ErrDetIdentityCreate   = "failed to create the x509 identity"
	ErrDetSDKInit          = "failed to initialize a new SDK instance"
)

// endregion =============================================================================

// region ======== SOME STRINGS ==========================================================

const (
	// ENV VARS

	EnvConfigPath = "SERVER_CONFIG"
	EnvJWTSignKey = "SERVER_JWT_SIGN_KEY"

	// CRYPTO MATERIALS

	WalletStr = "wallet"
)

// endregion =============================================================================

// region ========= CERTIFICATES =========================================================

const (
	DocType       = "CERT"
	CreateAsset   = "CreateAsset"
	ValidateAsset = "ValidateAsset"
	ReadAsset     = "ReadAsset"
	DeleteAsset   = "DeleteAsset"
)

// endregion =============================================================================
