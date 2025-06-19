package common_error

import (
	"fmt"

	api_error "github.com/ebobola-dev/socially-app-go-server/internal/errors"
	"github.com/ebobola-dev/socially-app-go-server/internal/response"
)

type BadContentTypeError struct {
	msg  string
	resp *response.ErrorResponse
}

func (e *BadContentTypeError) Error() string {
	return e.msg
}

func (e *BadContentTypeError) StatusCode() int {
	return 400
}

func (e *BadContentTypeError) Response() *response.ErrorResponse {
	return e.resp
}

func NewBadContentTypeErr(recievedType, requiredType string) api_error.IApiError {
	return &BadContentTypeError{
		msg: fmt.Sprintf("Bad Content-Type: %s, required: %s", recievedType, requiredType),
		resp: &response.ErrorResponse{
			Message: fmt.Sprintf("Content-Type must be '%s'", requiredType),
		},
	}
}
