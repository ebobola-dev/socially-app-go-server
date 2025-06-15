package otp_error

import (
	"github.com/ebobola-dev/socially-app-go-server/internal/response"
)

type IncorrectValueError struct {
	msg  string
	code int
	resp *response.ErrorResponse
}

func (e *IncorrectValueError) Error() string {
	return e.msg
}

func (e *IncorrectValueError) StatusCode() int {
	return e.code
}

func (e *IncorrectValueError) Response() *response.ErrorResponse {
	return e.resp
}

var ErrIncorect = &IncorrectValueError{
	msg:  "Incorrect",
	code: 400,
	resp: &response.ErrorResponse{
		Message: "Incorrect otp code",
	},
}
