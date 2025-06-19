package common_error

import (
	api_error "github.com/ebobola-dev/socially-app-go-server/internal/errors"
	"github.com/ebobola-dev/socially-app-go-server/internal/response"
)

type BadRequestError struct {
	serverMessage   string
	responseMessage string
}

func (e *BadRequestError) Error() string {
	return e.serverMessage
}

func (e *BadRequestError) StatusCode() int {
	return 400
}

func (e *BadRequestError) Response() *response.ErrorResponse {
	return &response.ErrorResponse{
		Message: e.responseMessage,
	}
}

// ? Optional fields
// ? First parameter - response message
// ? Second parameter - server message
// ? If first present, but no second -> server message = response message (first parameter)
func NewBadRequestErr(messages ...string) api_error.IApiError {
	responseMessage := "No message"
	serverMessage := "No message"
	if len(messages) >= 1 {
		responseMessage = messages[0]
		serverMessage = messages[0]
	}
	if len(messages) >= 2 {
		serverMessage = messages[1]
	}
	return &BadRequestError{
		serverMessage:   serverMessage,
		responseMessage: responseMessage,
	}
}
