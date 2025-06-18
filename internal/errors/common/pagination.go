package common_error

import (
	"github.com/ebobola-dev/socially-app-go-server/internal/response"
)

type InvalidPagintationError struct{}

func (e *InvalidPagintationError) Error() string {
	return "Invalid pagintation parameters"
}

func (e *InvalidPagintationError) StatusCode() int {
	return 400
}

func (e *InvalidPagintationError) Response() *response.ErrorResponse {
	return &response.ErrorResponse{
		Message: "Invalid pagintation parameters, need 'offset': int > 0, 'limit': int > 0, in query parameters",
	}
}

var ErrInvalidPagintation = &InvalidPagintationError{}
