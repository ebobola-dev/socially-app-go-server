package pagination

import (
	common_error "github.com/ebobola-dev/socially-app-go-server/internal/errors/common"
	"github.com/gofiber/fiber/v2"
)

var defaultLimit = 10

type Pagintation struct {
	Offset int `json:"offset"`
	Limit  int `json:"limit"`
}

func FromFiberCtx(c *fiber.Ctx) (*Pagintation, error) {
	offset := c.QueryInt("offset", 0)
	limit := c.QueryInt("limit", defaultLimit)

	if offset < 0 || limit < 0 {
		return &Pagintation{}, common_error.ErrInvalidPagintation
	}

	return &Pagintation{Offset: offset, Limit: limit}, nil
}

func Default() *Pagintation {
	return &Pagintation{Offset: 0, Limit: defaultLimit}
}
