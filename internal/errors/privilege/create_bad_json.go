package privilege_error

import (
	"github.com/ebobola-dev/socially-app-go-server/internal/response"
)

type CreateBadJsonError struct{}

func (e *CreateBadJsonError) Error() string {
	return "Bad json"
}

func (e *CreateBadJsonError) StatusCode() int {
	return 400
}

func (e *CreateBadJsonError) Response() *response.ErrorResponse {
	return &response.ErrorResponse{
		Message: "Need json body { 'name': string(max 64 char.), 'order_index': int(between 1 and 99), }",
	}
}

var ErrBadCreateJson = &CreateBadJsonError{}
