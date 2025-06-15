package validation

import (
	"unicode"

	"github.com/go-playground/validator/v10"
)

func validateUsernameCharset(fl validator.FieldLevel) bool {
	username := fl.Field().String()
	for _, r := range username {
		if !(unicode.IsLower(r) || unicode.IsDigit(r) || r == '_' || r == '.') {
			return false
		}
	}
	return true
}

func validateUsernameStartDot(fl validator.FieldLevel) bool {
	username := fl.Field().String()
	return !(len(username) > 0 && username[0] == '.')
}

func validateUsernameStartDigit(fl validator.FieldLevel) bool {
	username := fl.Field().String()
	return !(len(username) > 0 && unicode.IsDigit(rune(username[0])))
}

func validateUsernameLength(fl validator.FieldLevel) bool {
	l := len(fl.Field().String())
	return l >= 4 && l <= 16
}
