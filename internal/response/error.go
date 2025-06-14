package response

import (
	"strings"

	"github.com/go-playground/validator/v10"
)

type ErrorResponse struct {
	Message string            `json:"_message"`
	Fields  map[string]string `json:"-"`
}

func (e ErrorResponse) ToJSON() map[string]any {
	resp := map[string]any{
		"_message": e.Message,
	}
	if len(e.Fields) > 0 {
		for k, v := range e.Fields {
			resp[k] = v
		}
	}
	return resp
}

func ParseValidationErrors(err error) ErrorResponse {
	validationErrors, ok := err.(validator.ValidationErrors)
	if !ok {
		return ErrorResponse{Message: "invalid input"}
	}

	fieldErrors := make(map[string]string)
	for _, e := range validationErrors {
		field := strings.ToLower(e.Field())
		switch e.Tag() {
		case "required":
			fieldErrors[field] = "is required"
		case "email":
			fieldErrors[field] = "invalid email format"
		case "min":
			fieldErrors[field] = "too short"
		case "max":
			fieldErrors[field] = "too long"
		case "len":
			fieldErrors[field] = "length must be " + e.Param()
		case "otp_value":
			fieldErrors[field] = "otp value must be array of 4 numbers 0-9"
		default:
			fieldErrors[field] = "invalid value"
		}
	}

	return ErrorResponse{
		Message: "validation error",
		Fields:  fieldErrors,
	}
}
