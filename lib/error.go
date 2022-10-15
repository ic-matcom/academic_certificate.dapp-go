package lib

import (
	"fmt"
	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	enTranslations "github.com/go-playground/validator/v10/translations/en"
	"github.com/kataras/iris/v12"
	"dapp/schema/dto"
)

const DefaultErrorLocale = "en"

func InitTranslations(validate *validator.Validate) *ut.UniversalTranslator {
	english := en.New()
	uni := ut.New(english, english)
	if trans, found := uni.GetTranslator(DefaultErrorLocale); found {
		_ = enTranslations.RegisterDefaultTranslations(validate, trans)
	}

	return uni
}

// NewProblem construct a new api error struct and return a pointer to it
//
// - s [uint] ~ HTTP status tu respond
//
// - t [string] ~ Title of the error
//
// - d [string] ~ Description or detail of the error
func NewProblem(s uint, t string, d string) *dto.Problem {
	return &dto.Problem{Status: s, Title: t, Detail: d}
}

// HandleError the error, below you will find the right way to do that...
func HandleError(ctx iris.Context, ut *ut.UniversalTranslator, err error, code int) {
	if errs, ok := err.(validator.ValidationErrors); ok {
		t, _ := ut.GetTranslator(DefaultErrorLocale)
		// Wrap the errors with JSON format, the underline library returns the errors as interface.
		validationErrors := wrapValidationTranslateErrors(errs, t)

		// Fire an application/json+problem response and stop the handlers chain.
		ctx.StopWithProblem(code, iris.NewProblem().
			Title("Validation error").
			Detail("One or more fields failed to be validated").
			Type(ctx.RouteName()).
			Key("errors", validationErrors))
		return
	}

	// It's probably an internal JSON error, let's dont give more info here.
	ctx.StopWithStatus(iris.StatusInternalServerError)
	return
}

func wrapValidationErrors(validatorErrs validator.ValidationErrors) []dto.ValidationError {
	validationErrors := make([]dto.ValidationError, 0, len(validatorErrs))
	for _, validationErr := range validatorErrs {
		validationErrors = append(validationErrors, dto.ValidationError{
			ActualTag: validationErr.ActualTag(),
			Namespace: validationErr.Namespace(),
			Kind:      validationErr.Kind().String(),
			Type:      validationErr.Type().String(),
			Value:     fmt.Sprintf("%v", validationErr.Value()),
			Param:     validationErr.Param(),
			Message:   validationErr.Error(),
		})
	}

	return validationErrors
}

func wrapValidationTranslateErrors(validatorErrs validator.ValidationErrors, trans ut.Translator) []dto.ValidationError {
	validationErrors := make([]dto.ValidationError, 0, len(validatorErrs))
	for _, validationErr := range validatorErrs {
		validationErrors = append(validationErrors, dto.ValidationError{
			ActualTag: validationErr.ActualTag(),
			Namespace: validationErr.Namespace(),
			Kind:      validationErr.Kind().String(),
			Type:      validationErr.Type().String(),
			Value:     fmt.Sprintf("%v", validationErr.Value()),
			Param:     validationErr.Param(),
			Message:   validationErr.Translate(trans),
		})
	}
	return validationErrors
}

// addTranslation Add Your Own Error Message
func addTranslation(validate *validator.Validate, trans *ut.Translator, tag string, errMessage string) {
	registerFn := func(ut ut.Translator) error {
		return ut.Add(tag, errMessage, false)
	}

	transFn := func(ut ut.Translator, fe validator.FieldError) string {
		param := fe.Param()
		tag := fe.Tag()

		t, err := ut.T(tag, fe.Field(), param)
		if err != nil {
			return fe.(error).Error()
		}
		return t
	}

	_ = validate.RegisterTranslation(tag, *trans, registerFn, transFn)
}
