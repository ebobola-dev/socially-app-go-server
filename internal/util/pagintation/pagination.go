package pagination

import (
	common_error "github.com/ebobola-dev/socially-app-go-server/internal/errors/common"
	"github.com/gofiber/fiber/v2"
)

var defaultLimit = 10

type Pagitation struct {
	Offset int `json:"offset"`
	Limit  int `json:"limit"`
}

func FromFiberCtx(c *fiber.Ctx) (Pagitation, error) {
	offset := c.QueryInt("offset", 0)
	limit := c.QueryInt("limit", defaultLimit)

	if offset < 0 || limit < 0 {
		return Pagitation{}, common_error.ErrInvalidPagintation
	}

	return Pagitation{Offset: offset, Limit: limit}, nil
}

func Default() Pagitation {
	return Pagitation{Offset: 0, Limit: defaultLimit}
}
