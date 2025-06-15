package otp_error

import (
	"fmt"

	api_error "github.com/ebobola-dev/socially-app-go-server/internal/errors"
	"github.com/ebobola-dev/socially-app-go-server/internal/response"
)

type CantUpdateOtpError struct {
	msg  string
	code int
	resp *response.ErrorResponse
}

func (e *CantUpdateOtpError) Error() string {
	return e.msg
}

func (e *CantUpdateOtpError) StatusCode() int {
	return e.code
}

func (e *CantUpdateOtpError) Response() *response.ErrorResponse {
	return e.resp
}

func NewCantUpdateOtpError(secondsDelta int) api_error.ApiError {
	return &CantUpdateOtpError{
		msg:  fmt.Sprintf("Can't update(delta: %ds)", secondsDelta),
		code: 429,
		resp: &response.ErrorResponse{
			Message: "Wait a minute before resend otp code",
		},
	}
}
