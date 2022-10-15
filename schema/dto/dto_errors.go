package dto

// Problem api documentation
type Problem struct {
	Status uint   `example:"503"`
	Title  string `example:"err_code"`
	Detail string `example:"Some error details"`
}

type ValidationError struct {
	ActualTag string `json:"tag"`
	Namespace string `json:"namespace"`
	Kind      string `json:"kind"`
	Type      string `json:"type"`
	Value     string `json:"value"`
	Param     string `json:"param"`
	Message   string `json:"message"`
}
