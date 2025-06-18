package privilege_error

import (
	"github.com/ebobola-dev/socially-app-go-server/internal/response"
)

type DeletingOwnerError struct{}

func (e *DeletingOwnerError) Error() string {
	return "Trying to delete owner privilege"
}

func (e *DeletingOwnerError) StatusCode() int {
	return 403
}

func (e *DeletingOwnerError) Response() *response.ErrorResponse {
	return &response.ErrorResponse{
		Message: "It is forbidden to delete the owner privilege",
	}
}

var ErrDeletingOwner = &DeletingOwnerError{}
