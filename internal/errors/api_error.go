package api_error

import "github.com/ebobola-dev/socially-app-go-server/internal/response"

type ApiError interface {
	error
	StatusCode() int
	Response() *response.ErrorResponse
}
