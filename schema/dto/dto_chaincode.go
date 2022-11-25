package dto

type StateValidation uint

const (
	Invalid  StateValidation = iota // invalidated for some reason
	New                             // certificate without signatures
	SignedS                         // signed by Secretary
	SignedSD                        // signed by Secretary and Dean
	Valid                           // signed by Secretary, Dean and Rector
)

type ValidatorType uint

const (
	NoValidator ValidatorType = iota // invalidated for some reason
	Secretary                        // certificate without signatures
	Dean                             // signed by Secretary
	Rector                           // signed by Secretary and Dean
)

//Auxiliary Functions
func (state StateValidation) String() string {
	names := []string{"Invalid", "New", "SignedS", "SignedSD", "Valid"}
	if state < Invalid || state > Valid {
		return "unknown"
	}
	return names[state]
}

// Asset describes basic details of what makes up a simple asset
type Asset struct {
	DocType               string          `json:"docType" validate:"required"`
	ID                    string          `json:"ID" validate:"required"`
	Certification         string          `json:"certification" validate:"required"`
	GoldCertificate       bool            `json:"gold_certificate"`
	Emitter               string          `json:"emitter" validate:"required"`
	Accredited            string          `json:"accredited" validate:"required"`
	Date                  string          `json:"date" validate:"required"`
	CreatedBy             string          `json:"created_by" validate:"required"`
	SecretaryValidating   string          `json:"secretary_validating"`
	DeanValidating        string          `json:"dean_validating"`
	RectorValidating      string          `json:"rector_validating"`
	FacultyVolumeFolio    string          `json:"volume_folio_faculty" validate:"required"`
	UniversityVolumeFolio string          `json:"volume_folio_university" validate:"required"`
	InvalidReason         string          `json:"invalid_reason"`
	Status                StateValidation `json:"certificate_status" validate:"gte=0,lte=4"`
}

type ValidateAsset struct {
	ID         string        `json:"ID"`
	Validator  string        `json:"validator"`
	ValidatorT ValidatorType `json:"validator_type"`
}

type CreateAsset struct {
	Certification         string `json:"certification" validate:"required"`
	GoldCertificate       bool   `json:"gold_certificate" validate:"required"`
	Emitter               string `json:"emitter" validate:"required"`
	Accredited            string `json:"accredited" validate:"required"`
	Date                  string `json:"date" validate:"required"`
	CreatedBy             string `json:"created_by" validate:"required"`
	FacultyVolumeFolio    string `json:"volume_folio_faculty" validate:"required"`
	UniversityVolumeFolio string `json:"volume_folio_university" validate:"required"`
}

type InvalidateAsset struct {
	ID            string `json:"ID" validate:"required"`
	InvalidReason string `json:"invalid_reason" validate:"required"`
}

type SignAsset struct {
	ID       string `json:"ID" validate:"required"`
	SignedBy string `json:"signed_by" validate:"required"`
}

type QueryParamChaincode struct {
	Channel   string `query:"channel"`
	Chaincode string `query:"chaincode"`
	Signer    string `query:"signer"`
	Bookmark  string `query:"bookmark"`
	PageLimit int    `query:"page_limit"`
}
