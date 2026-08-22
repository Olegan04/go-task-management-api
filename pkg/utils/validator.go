package utils

import (
	"fmt"

	"github.com/go-playground/validator/v10"
)

var validate = validator.New()

func ValidateStruct(s interface{}) map[string]string {
	err := validate.Struct(s)
	if err == nil {
		return nil
	}
	errorsMap := make(map[string]string)

	for _, err := range err.(validator.ValidationErrors) {
		field := err.Field()
		tag := err.Tag()
		param := err.Param()

		var msg string
		switch tag {
		case "required":
			msg = fmt.Sprintf("%s is required", field)
		case "email":
			msg = fmt.Sprintf("%s must be a valid email address", field)
		case "min":
			msg = fmt.Sprintf("%s must be at least %s characters", field, param)
		case "max":
			msg = fmt.Sprintf("%s must be at most %s characters", field, param)
		case "oneof":
			msg = fmt.Sprintf("%s must be one of: %s", field, param)
		default:
			msg = fmt.Sprintf("%s is invalid (validation: %s)", field, tag)
		}
		errorsMap[field] = msg
	}
	return errorsMap
}
