package privilege_error

import (
	"fmt"

	api_error "github.com/ebobola-dev/socially-app-go-server/internal/errors"
	"github.com/ebobola-dev/socially-app-go-server/internal/response"
)

type DuplicateNameError struct {
	msg  string
	resp *response.ErrorResponse
}

func (e *DuplicateNameError) Error() string {
	return e.msg
}

func (e *DuplicateNameError) StatusCode() int {
	return 400
}

func (e *DuplicateNameError) Response() *response.ErrorResponse {
	return e.resp
}

func NewDuplicateNameError(name string) api_error.ApiError {
	return &DuplicateNameError{
		msg: fmt.Sprintf("Duplicate privilege name: %s", name),
		resp: &response.ErrorResponse{
			Message: fmt.Sprintf("Privilege '%s' already exists", name),
		},
	}
}
