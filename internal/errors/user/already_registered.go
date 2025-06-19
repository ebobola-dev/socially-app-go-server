package user_error

import (
	"fmt"

	api_error "github.com/ebobola-dev/socially-app-go-server/internal/errors"
	"github.com/ebobola-dev/socially-app-go-server/internal/response"
)

type AlreadyRegisteredError struct {
	msg  string
	code int
	resp *response.ErrorResponse
}

func (e *AlreadyRegisteredError) Error() string {
	return e.msg
}

func (e *AlreadyRegisteredError) StatusCode() int {
	return e.code
}

func (e *AlreadyRegisteredError) Response() *response.ErrorResponse {
	return e.resp
}

func NewAlreadyRegisteredError(email string) api_error.IApiError {
	return &AlreadyRegisteredError{
		msg:  "Email already registered",
		code: 400,
		resp: &response.ErrorResponse{
			Message: fmt.Sprintf("User with email '%s' already registered", email),
		},
	}
}
