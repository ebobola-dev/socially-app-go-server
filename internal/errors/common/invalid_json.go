package common_error

import (
	api_error "github.com/ebobola-dev/socially-app-go-server/internal/errors"
	"github.com/ebobola-dev/socially-app-go-server/internal/response"
)

type InvalidJSONError struct {
	msg  string
	code int
	resp *response.ErrorResponse
}

func (e *InvalidJSONError) Error() string {
	return e.msg
}

func (e *InvalidJSONError) StatusCode() int {
	return e.code
}

func (e *InvalidJSONError) Response() *response.ErrorResponse {
	return e.resp
}

func NewInvalidJSONError(message string) api_error.ApiError {
	return &InvalidJSONError{
		msg:  "Invalid JSON",
		code: 400,
		resp: &response.ErrorResponse{
			Message: message,
		},
	}
}
