package validation

import (
	"github.com/ebobola-dev/socially-app-go-server/internal/model"
	"github.com/go-playground/validator/v10"
)

func OtpValueValidator(fl validator.FieldLevel) bool {
	val, ok := fl.Field().Interface().(model.OtpValue)
	if !ok {
		return false
	}

	if len(val) != 4 {
		return false
	}
	for _, digit := range val {
		if digit < 0 || digit > 9 {
			return false
		}
	}
	return true
}

func RegisterCustomValidators(v *validator.Validate) {
	v.RegisterValidation("otp_value", OtpValueValidator)
}
