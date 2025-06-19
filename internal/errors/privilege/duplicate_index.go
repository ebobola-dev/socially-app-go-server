package privilege_error

import (
	"fmt"

	api_error "github.com/ebobola-dev/socially-app-go-server/internal/errors"
	"github.com/ebobola-dev/socially-app-go-server/internal/response"
)

type DuplicateIndexError struct {
	msg  string
	resp *response.ErrorResponse
}

func (e *DuplicateIndexError) Error() string {
	return e.msg
}

func (e *DuplicateIndexError) StatusCode() int {
	return 400
}

func (e *DuplicateIndexError) Response() *response.ErrorResponse {
	return e.resp
}

func NewDuplicateIndexError(orderIndex int) api_error.IApiError {
	return &DuplicateIndexError{
		msg: fmt.Sprintf("Duplicate privilege order_index: %d", orderIndex),
		resp: &response.ErrorResponse{
			Message: fmt.Sprintf("Privilege with order_index: '%d' already exists", orderIndex),
		},
	}
}
