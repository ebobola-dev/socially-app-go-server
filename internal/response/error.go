package response

import (
	string_utils "github.com/ebobola-dev/socially-app-go-server/internal/util"
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
		field := string_utils.ToSnakeCase(e.Field())
		switch e.Tag() {
		case "required":
			fieldErrors[field] = "is required"
		case "email":
			fieldErrors[field] = "invalid email format"
		case "min":
			fieldErrors[field] = "min length " + e.Param()
		case "max":
			fieldErrors[field] = "max length " + e.Param()
		case "len":
			fieldErrors[field] = "length must be " + e.Param()
		case "gt":
			fieldErrors[field] = "must be greater than " + e.Param()
		case "gte":
			fieldErrors[field] = "must be greater than or equal to" + e.Param()
		case "lt":
			fieldErrors[field] = "must be less than " + e.Param()
		case "lte":
			fieldErrors[field] = "must be less than or equal to" + e.Param()
		case "uuid4":
			fieldErrors[field] = "invalid uuid4 string"
		case "otp_value":
			fieldErrors[field] = "otp value must be array of 4 numbers 0-9"
		case "date":
			fieldErrors[field] = "must be string date in dd.mm.yyyy format"
		case "datebt":
			fieldErrors[field] = "must be string date in dd.mm.yyyy format, before today"
		case "gender":
			fieldErrors[field] = "must be string - male or female"
		case "password":
			fieldErrors[field] = "at least one letter, at least one digit, between 8 and 32 characters"
		case "avatar_type":
			fieldErrors[field] = "must be string - (external, avatar1, avatar2, ..., avatar10)"
		case "username_length":
			fieldErrors[field] = "length must be between 4 and 32 characters"
		case "username_charset":
			fieldErrors[field] = "only lowercase Latin letters, numbers, underscores and dots are allowed"
		case "username_start_digit":
			fieldErrors[field] = "cannot start with a number"
		case "username_start_dot":
			fieldErrors[field] = "cannot start with a dot"
		default:
			fieldErrors[field] = "invalid value"
		}
	}

	return ErrorResponse{
		Message: "validation error",
		Fields:  fieldErrors,
	}
}
