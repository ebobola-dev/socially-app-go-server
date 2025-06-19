package api_error

import "github.com/ebobola-dev/socially-app-go-server/internal/response"

type IApiError interface {
	error
	StatusCode() int
	Response() *response.ErrorResponse
}

type ApiError struct {
	ServerMessage string
	Code          int
	Resp          *response.ErrorResponse
}

func (e *ApiError) Error() string {
	return e.ServerMessage
}

func (e *ApiError) StatusCode() int {
	return e.Code
}

func (e *ApiError) Response() *response.ErrorResponse {
	return e.Resp
}

func NewUnexceptedErr(err error) IApiError {
	return &ApiError{
		ServerMessage: err.Error(),
		Code:          500,
		Resp: &response.ErrorResponse{
			Message: "Unexcepted server error",
		},
	}
}
