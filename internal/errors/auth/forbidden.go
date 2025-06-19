package auth_error

import (
	"fmt"

	api_error "github.com/ebobola-dev/socially-app-go-server/internal/errors"
	"github.com/ebobola-dev/socially-app-go-server/internal/response"
)

type ForbiddenError struct {
	msg  string
	resp *response.ErrorResponse
}

func (e *ForbiddenError) Error() string {
	return e.msg
}

func (e *ForbiddenError) StatusCode() int {
	return 403
}

func (e *ForbiddenError) Response() *response.ErrorResponse {
	return e.resp
}

func NewForbidden(serverMsg, responseMsg string) api_error.IApiError {
	return &ForbiddenError{
		msg: serverMsg,
		resp: &response.ErrorResponse{
			Message: responseMsg,
		},
	}
}

func NewNoAnyPrivilegeError(requiredPrivileges ...string) api_error.IApiError {
	return &ForbiddenError{
		msg: "No one of necessary privileges",
		resp: &response.ErrorResponse{
			Message: fmt.Sprintf("Required one of privileges: %v", requiredPrivileges),
		},
	}
}

func NewNoAllPrivilegeError(requiredPrivileges ...string) api_error.IApiError {
	return &ForbiddenError{
		msg: "No privileges",
		resp: &response.ErrorResponse{
			Message: fmt.Sprintf("Required privileges: %v", requiredPrivileges),
		},
	}
}
