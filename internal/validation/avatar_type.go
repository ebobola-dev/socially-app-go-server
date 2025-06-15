package validation

import (
	"github.com/ebobola-dev/socially-app-go-server/internal/model"
	"github.com/go-playground/validator/v10"
)

func validateAvatarType(fl validator.FieldLevel) bool {
	at := fl.Field().String()
	return model.IsValidAvatarType(at)
}
