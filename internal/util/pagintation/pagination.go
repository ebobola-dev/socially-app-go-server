package pagination

import (
	common_error "github.com/ebobola-dev/socially-app-go-server/internal/errors/common"
	"github.com/gofiber/fiber/v2"
)

var defaultLimit = 10

type Pagination struct {
	Offset int `json:"offset"`
	Limit  int `json:"limit"`
}

func (p *Pagination) ToMap() map[string]int {
	return map[string]int{
		"offset": p.Offset,
		"limit":  p.Limit,
	}
}

func FromFiberCtx(c *fiber.Ctx) (Pagination, error) {
	offset := c.QueryInt("offset", 0)
	limit := c.QueryInt("limit", defaultLimit)

	if offset < 0 || limit < 0 {
		return Pagination{}, common_error.ErrInvalidPagintation
	}

	return Pagination{Offset: offset, Limit: limit}, nil
}

func Default() Pagination {
	return Pagination{Offset: 0, Limit: defaultLimit}
}
