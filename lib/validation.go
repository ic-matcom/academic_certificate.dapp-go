package lib

import (
	"github.com/asaskevich/govalidator"
	"github.com/go-playground/validator/v10"
	"reflect"
	reg "regexp"
	"strings"
)

// ValidateString validate a string given a regular expression
func ValidateString(data string, regexp string) bool {
	return reg.MustCompile(regexp).MatchString(data)
}

// ValidateStringCollection validate a string collection given a regular expression
func ValidateStringCollection(data []interface{}, regexp string) bool {
	var fn govalidator.ConditionIterator = func(value interface{}, index int) bool {
		return reg.MustCompile(regexp).MatchString(value.(string))
	}
	return govalidator.ValidateArray(data, fn)
}

// ValidateStringCollectionUsingValidator10 validation into arrays type
// e.g.
// tag = "required,max=10,min=1,dive,max=12"
//
// 			max=10 		 -> Max array len
// 			dive, max=12 -> Max length of every array element
func ValidateStringCollectionUsingValidator10(validate *validator.Validate, data any, tag string) bool {
	// variable must be a slice
	if reflect.TypeOf(data).Kind() != reflect.Slice {
		return false
	}

	errs := validate.Var(data, tag)
	if errs != nil {
		return false
	}
	if reflect.TypeOf(data).Elem().Kind() != reflect.String {
		return false
	}
	return true
}

// InitValidator Activate behavior to require all fields and adding new validators
func InitValidator(validate *validator.Validate) error {
	// Add here your custom validation
	return nil
}

// NotBlank is the validation function for validating if the current field
// has a value or length greater than zero, or is not a space only string.
// example: v.RegisterValidation("notblank", NotBlank)
func NotBlank(fl validator.FieldLevel) bool {
	field := fl.Field()

	switch field.Kind() {
	case reflect.String:
		return len(strings.TrimSpace(field.String())) > 0
	case reflect.Chan, reflect.Map, reflect.Slice, reflect.Array:
		return field.Len() > 0
	case reflect.Ptr, reflect.Interface, reflect.Func:
		return !field.IsNil()
	default:
		return field.IsValid() && field.Interface() != reflect.Zero(field.Type()).Interface()
	}
}
