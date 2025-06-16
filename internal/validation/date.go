package validation

import (
	"time"

	"github.com/go-playground/validator/v10"
)

func validateDate(fl validator.FieldLevel) bool {
	dateStr := fl.Field().String()
	_, err := time.Parse("02.01.2006", dateStr)
	return err == nil
}
