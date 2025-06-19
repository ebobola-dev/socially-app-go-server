package common_error

import (
	api_error "github.com/ebobola-dev/socially-app-go-server/internal/errors"
	"github.com/ebobola-dev/socially-app-go-server/internal/response"
)

type InvalidJSONError struct {
	resp *response.ErrorResponse
}

func (e *InvalidJSONError) Error() string {
	return "Invalid JSON"
}

func (e *InvalidJSONError) StatusCode() int {
	return 400
}

func (e *InvalidJSONError) Response() *response.ErrorResponse {
	return e.resp
}

var ErrInvalidJSON = &InvalidJSONError{
	resp: &response.ErrorResponse{
		Message: "Need JSON body",
	},
}

func NewInvalidJsonErr(message string) api_error.IApiError {
	return &InvalidJSONError{
		resp: &response.ErrorResponse{
			Message: message,
		},
	}
}
