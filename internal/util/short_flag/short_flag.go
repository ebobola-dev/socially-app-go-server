package short_flag

import (
	"strings"

	"github.com/gofiber/fiber/v2"
)

func FromFiberCtx(c *fiber.Ctx) bool {
	short := strings.ToLower(strings.TrimSpace(c.Query("short", "")))
	return short == "1" || short == "true" || short == "yes" || short == "on"
}
