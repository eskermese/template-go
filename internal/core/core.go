package core

import (
	"errors"
	"strings"

	"github.com/go-playground/validator/v10"
)

var validate *validator.Validate

func init() {
	validate = validator.New()
}

func validateStruct(st interface{}) error {
	var res Errors

	if err := validate.Struct(st); err != nil {
		var validationErrors validator.ValidationErrors
		if errors.As(err, &validationErrors) {
			for _, e := range validationErrors {
				field := strings.ToLower(e.StructField())
				msg := getErrorMsg(e)

				res = append(res, FieldError{Field: field, Message: msg})
			}
		}
	}

	if len(res) == 0 {
		return nil
	}

	return res
}

func getErrorMsg(fe validator.FieldError) string {
	switch fe.Tag() {
	case "required":
		return "This field is required"
	case "lte":
		return "Should be less than " + fe.Param()
	case "gte":
		return "Should be greater than " + fe.Param()
	case "min":
		return "Should be minimum " + fe.Param()
	case "max":
		return "Should be maximum " + fe.Param()
	}

	return fe.Error() // default error
}
