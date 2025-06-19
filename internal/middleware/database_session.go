package middleware

import (
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func DatabaseSession(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var handlerErr error

		db.Transaction(func(tx *gorm.DB) error {
			c.Locals("tx", tx)
			handlerErr = c.Next()
			if handlerErr != nil {
				return handlerErr
			}
			return nil
		})

		return handlerErr
	}
}

func GetTX(c *fiber.Ctx) *gorm.DB {
	tx, ok := c.Locals("tx").(*gorm.DB)
	if !ok {
		panic("Database transaction not found in context")
	}
	return tx
}
