package common_error

import (
	"github.com/ebobola-dev/socially-app-go-server/internal/response"
)

type MissingDeviceIdError struct {
	msg  string
	code int
	resp *response.ErrorResponse
}

func (e *MissingDeviceIdError) Error() string {
	return e.msg
}

func (e *MissingDeviceIdError) StatusCode() int {
	return e.code
}

func (e *MissingDeviceIdError) Response() *response.ErrorResponse {
	return e.resp
}

var ErrMissingDeviceId = &MissingDeviceIdError{
	msg:  "Missing device id",
	code: 400,
	resp: &response.ErrorResponse{
		Message: "device_id must be specified in request headers",
	},
}
