package auth_error

import (
	api_error "github.com/ebobola-dev/socially-app-go-server/internal/errors"
	"github.com/ebobola-dev/socially-app-go-server/internal/response"
)

type InvalidLoginDataError struct {
	msg  string
	code int
	resp *response.ErrorResponse
}

func (e *InvalidLoginDataError) Error() string {
	return e.msg
}

func (e *InvalidLoginDataError) StatusCode() int {
	return e.code
}

func (e *InvalidLoginDataError) Response() *response.ErrorResponse {
	return e.resp
}

func NewInvalidLoginData(serverMsg string) api_error.IApiError {
	return &InvalidLoginDataError{
		msg:  serverMsg,
		code: 400,
		resp: &response.ErrorResponse{
			Message: "Incorrent data",
		},
	}
}
