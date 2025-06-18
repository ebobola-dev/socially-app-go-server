package user_error

import (
	"fmt"

	api_error "github.com/ebobola-dev/socially-app-go-server/internal/errors"
	"github.com/ebobola-dev/socially-app-go-server/internal/response"
)

type UsernameTakenError struct {
	msg  string
	code int
	resp *response.ErrorResponse
}

func (e *UsernameTakenError) Error() string {
	return e.msg
}

func (e *UsernameTakenError) StatusCode() int {
	return e.code
}

func (e *UsernameTakenError) Response() *response.ErrorResponse {
	return e.resp
}

func NewUsernameTakenError(username string) api_error.ApiError {
	return &UsernameTakenError{
		msg:  fmt.Sprintf("Username is taken: %s", username),
		code: 400,
		resp: &response.ErrorResponse{
			Message: "Username is already taken",
			Fields: map[string]string{
				"username": username,
			},
		},
	}
}
