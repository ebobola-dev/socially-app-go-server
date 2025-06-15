package common_error

import (
	"fmt"

	api_error "github.com/ebobola-dev/socially-app-go-server/internal/errors"
	"github.com/ebobola-dev/socially-app-go-server/internal/response"
)

type RecordNotFoundError struct {
	msg  string
	code int
	resp *response.ErrorResponse
}

func (e *RecordNotFoundError) Error() string {
	return e.msg
}

func (e *RecordNotFoundError) StatusCode() int {
	return e.code
}

func (e *RecordNotFoundError) Response() *response.ErrorResponse {
	return e.resp
}

func NewRecordNotFoundError(recordName string) api_error.ApiError {
	return &RecordNotFoundError{
		msg:  fmt.Sprintf("%s not found", recordName),
		code: 400,
		resp: &response.ErrorResponse{
			Message: fmt.Sprintf("%s not found", recordName),
		},
	}
}
