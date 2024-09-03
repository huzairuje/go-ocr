package validator

import (
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
)

var validate = validator.New()

type ErrorResponse struct {
	FailedField string
	Tag         string
	Value       string
}

func ValidateStruct(input interface{}) (errResponse []*ErrorResponse) {
	err := validate.Struct(input)
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			var element ErrorResponse
			element.FailedField = trimStringFromDot(err.StructNamespace())
			element.Tag = err.Tag()
			element.Value = err.Param()
			errResponse = append(errResponse, &element)
		}
	}
	return errResponse
}

func ValidateStructResponseSliceString(input interface{}) (errResponse []string) {
	err := validate.Struct(input)
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			var message string
			var element ErrorResponse
			element.FailedField = trimStringFromDot(err.StructNamespace())
			message = fmt.Sprintf("form %s must filled", element.FailedField)
			errResponse = append(errResponse, message)
		}
	}
	return errResponse
}

func trimStringFromDot(str string) string {
	// find the position of the first dot
	dotPos := strings.Index(str, ".")
	// if there's no dot in the string, return the original string
	if dotPos == -1 {
		return str
	}
	// return the substring from the dot until last char
	return strings.ToLower(str[dotPos+1:])
}
