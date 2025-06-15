package common_error

import (
	"github.com/ebobola-dev/socially-app-go-server/internal/response"
)

type NotImplementedError struct {
	msg  string
	code int
	resp *response.ErrorResponse
}

func (e *NotImplementedError) Error() string {
	return e.msg
}

func (e *NotImplementedError) StatusCode() int {
	return e.code
}

func (e *NotImplementedError) Response() *response.ErrorResponse {
	return e.resp
}

var ErrNotImplemented = &NotImplementedError{
	msg:  "Not implemented",
	code: 501,
	resp: &response.ErrorResponse{
		Message: "Not implemented yet",
	},
}
