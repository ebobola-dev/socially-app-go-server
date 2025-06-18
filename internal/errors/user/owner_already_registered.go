package user_error

import (
	"github.com/ebobola-dev/socially-app-go-server/internal/response"
)

type OwnerRegisteredError struct{}

func (e *OwnerRegisteredError) Error() string {
	return "Owner already registered"
}

func (e *OwnerRegisteredError) StatusCode() int {
	return 400
}

func (e *OwnerRegisteredError) Response() *response.ErrorResponse {
	return &response.ErrorResponse{
		Message: "Owner already registered",
	}
}

var ErrOwnerAlreadyRegistered = &OwnerRegisteredError{}
