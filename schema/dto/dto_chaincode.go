package dto

type StateValidation uint

const (
	Invalid StateValidation = iota
	One
	Two
	Three
)

//Auxiliary Functions
func (state StateValidation) String() string {
	names := []string{"Invalid", "One", "Two", "Three"}
	if state < Invalid || state > Three {
		return "unknown"
	}
	return names[state]
}

// TODO: agregar validaciones usando la bilbioteca github.com/go-playground/validator/v10
// hay un ejemplo en el  dto.UserCredIn en dto/dto_auth.go

// Asset describes basic details of what makes up a simple asset
type Asset struct {
	DocType       string          `json:"docType"`
	ID            string          `json:"ID"`
	Color         string          `json:"color"`
	Size          int             `json:"size"`
	Owner         string          `json:"owner"`
	OperationType StateValidation `json:"operationType"`
}
