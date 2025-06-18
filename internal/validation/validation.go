package validation

import (
	"reflect"

	"github.com/ebobola-dev/socially-app-go-server/internal/model"
	"github.com/ebobola-dev/socially-app-go-server/internal/util/nullable"
	"github.com/go-playground/validator/v10"
)

func NewValidator() *validator.Validate {
	v := validator.New()
	v.RegisterValidation("otp_value", validateOtpValue)
	v.RegisterValidation("date", validateDate)
	v.RegisterValidation("datebt", validateDateBt)
	v.RegisterValidation("gender", validateGender)
	v.RegisterValidation("avatar_type", validateAvatarType)

	v.RegisterValidation("username_length", validateUsernameLength)
	v.RegisterValidation("username_charset", validateUsernameCharset)
	v.RegisterValidation("username_start_digit", validateUsernameStartDigit)
	v.RegisterValidation("username_start_dot", validateUsernameStartDot)
	v.RegisterValidation("password", validatePassword)

	v.RegisterCustomTypeFunc(func(val reflect.Value) interface{} {
		if !val.IsValid() || val.IsZero() {
			return nil
		}
		field := val.FieldByName("Value")
		if !field.IsValid() {
			return nil
		}
		return field.Interface()
	},
		nullable.Nullable[string]{},
		nullable.Nullable[int]{},
		nullable.Nullable[model.Gender]{},
	)
	return v
}
