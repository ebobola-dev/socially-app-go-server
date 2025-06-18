package user_error

import (
	"github.com/ebobola-dev/socially-app-go-server/internal/response"
)

type NothingToUpdateProfileError struct{}

func (e *NothingToUpdateProfileError) Error() string {
	return "Nothing to update"
}

func (e *NothingToUpdateProfileError) StatusCode() int {
	return 400
}

func (e *NothingToUpdateProfileError) Response() *response.ErrorResponse {
	return &response.ErrorResponse{
		Message: "Nothing to update",
	}
}

var ErrNothingToUpdateProfile = &NothingToUpdateProfileError{}
