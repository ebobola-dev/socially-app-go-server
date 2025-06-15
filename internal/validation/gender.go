package validation

import (
	"github.com/ebobola-dev/socially-app-go-server/internal/model"
	"github.com/go-playground/validator/v10"
)

func validateGender(fl validator.FieldLevel) bool {
	g := fl.Field().String()
	return model.IsValidGender(g)
}
