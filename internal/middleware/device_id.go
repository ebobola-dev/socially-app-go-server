package middleware

import (
	common_error "github.com/ebobola-dev/socially-app-go-server/internal/errors/common"
	"github.com/gofiber/fiber/v2"
)

func DeviceId() fiber.Handler {
	return func(c *fiber.Ctx) error {
		deviceIdHeaders := c.Request().Header.Peek("device_id")
		if len(deviceIdHeaders) == 0 {
			return common_error.ErrMissingDeviceId
		}
		device_id := string(deviceIdHeaders)
		if device_id == "" {
			return common_error.ErrMissingDeviceId
		}
		c.Locals("device_id", device_id)
		return c.Next()
	}
}

func GetDeviceId(c *fiber.Ctx) string {
	tx, ok := c.Locals("device_id").(string)
	if !ok {
		panic("device_id not found in context")
	}
	return tx
}
