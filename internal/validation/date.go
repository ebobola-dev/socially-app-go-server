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

// ? Before today
func validateDateBt(fl validator.FieldLevel) bool {
	dateStr := fl.Field().String()
	date, err := time.Parse("02.01.2006", dateStr)
	if err != nil {
		return false
	}
	now := time.Now().UTC()
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	return date.Before(today)
}
