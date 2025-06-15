package validation

import "github.com/go-playground/validator/v10"

func NewValidator() *validator.Validate {
	v := validator.New()
	v.RegisterValidation("otp_value", validateOtpValue)
	v.RegisterValidation("date", validateDate)
	v.RegisterValidation("gender", validateGender)
	v.RegisterValidation("avatar_type", validateAvatarType)

	v.RegisterValidation("username_length", validateUsernameLength)
	v.RegisterValidation("username_charset", validateUsernameCharset)
	v.RegisterValidation("username_start_digit", validateUsernameStartDigit)
	v.RegisterValidation("username_start_dot", validateUsernameStartDot)
	v.RegisterValidation("password", validatePassword)
	return v
}
