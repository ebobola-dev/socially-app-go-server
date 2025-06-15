package validation

import (
	"unicode"

	"github.com/go-playground/validator/v10"
)

func validatePassword(fl validator.FieldLevel) bool {
	password := fl.Field().String()

	if len(password) < 8 || len(password) > 32 {
		return false
	}

	var hasLetter, hasDigit bool
	for _, c := range password {
		switch {
		case unicode.IsLetter(c):
			hasLetter = true
		case unicode.IsDigit(c):
			hasDigit = true
		}
		if hasLetter && hasDigit {
			return true
		}
	}
	return false
}
