package middleware

import (
	pagination "github.com/ebobola-dev/socially-app-go-server/internal/util/pagintation"

	"github.com/gofiber/fiber/v2"
)

func Pagination() fiber.Handler {
	return func(c *fiber.Ctx) error {
		pagination, err := pagination.FromFiberCtx(c)
		if err != nil {
			return err
		}
		c.Locals("pagintation", pagination)
		return c.Next()
	}
}

func GetPagination(c *fiber.Ctx) pagination.Pagination {
	tx, ok := c.Locals("pagintation").(pagination.Pagination)
	if !ok {
		panic("pagintation not found in context")
	}
	return tx
}
