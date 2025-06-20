package common_error

import (
	"fmt"

	api_error "github.com/ebobola-dev/socially-app-go-server/internal/errors"
	"github.com/ebobola-dev/socially-app-go-server/internal/response"
)

type RecordNotFoundError struct {
	serverMessage   string
	responseMessage string
}

func (e *RecordNotFoundError) Error() string {
	return e.serverMessage
}

func (e *RecordNotFoundError) StatusCode() int {
	return 404
}

func (e *RecordNotFoundError) Response() *response.ErrorResponse {
	return &response.ErrorResponse{
		Message: e.responseMessage,
	}
}

func NewRecordNotFoundErr(recordName string) api_error.IApiError {
	return &RecordNotFoundError{
		serverMessage:   fmt.Sprintf("%s not found", recordName),
		responseMessage: fmt.Sprintf("%s not found", recordName),
	}
}

func NewMinioNotFoundErr(bucket, path string) api_error.IApiError {
	return &RecordNotFoundError{
		serverMessage:   fmt.Sprintf("File not found in minio bucket[%s]: %s", bucket, path),
		responseMessage: "File not found",
	}
}
