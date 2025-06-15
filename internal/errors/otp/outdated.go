package otp_error

import (
	api_error "github.com/ebobola-dev/socially-app-go-server/internal/errors"
	"github.com/ebobola-dev/socially-app-go-server/internal/response"
)

type OtdIsOutdatedError struct {
	msg  string
	code int
	resp *response.ErrorResponse
}

func (e *OtdIsOutdatedError) Error() string {
	return e.msg
}

func (e *OtdIsOutdatedError) StatusCode() int {
	return e.code
}

func (e *OtdIsOutdatedError) Response() *response.ErrorResponse {
	return e.resp
}

func NewOtdIsOutdatedError() api_error.ApiError {
	return &OtdIsOutdatedError{
		msg:  "Outdated",
		code: 400,
		resp: &response.ErrorResponse{
			Message: "OTP code is outdated, resend new otp code",
		},
	}
}
