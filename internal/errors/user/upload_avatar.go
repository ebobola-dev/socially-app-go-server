package user_error

import (
	"fmt"

	api_error "github.com/ebobola-dev/socially-app-go-server/internal/errors"
	"github.com/ebobola-dev/socially-app-go-server/internal/response"
	size_util "github.com/ebobola-dev/socially-app-go-server/internal/util/size"
)

func NewAvatarTooLargeErr(receivedSize, maxAllowedSize int64) api_error.IApiError {
	return &api_error.ApiError{
		ServerMessage: fmt.Sprintf("Avatar too large: %s", size_util.BytesToString(receivedSize)),
		Code:          413,
		Resp: &response.ErrorResponse{
			Message: "Request Entity Too Large",
			Fields: map[string]string{
				"avatar": fmt.Sprintf("file too large, max allowed size: %s", size_util.BytesToString(maxAllowedSize)),
			},
		},
	}
}

func NewBadAvatarExtensionErr(receivedExt string, allowedExts []string) api_error.IApiError {
	return &api_error.ApiError{
		ServerMessage: fmt.Sprintf("Bad avatar extension: %s", receivedExt),
		Code:          400,
		Resp: &response.ErrorResponse{
			Message: "Validation error",
			Fields: map[string]string{
				"avatar": fmt.Sprintf("Bad extension, allowed: %v", allowedExts),
			},
		},
	}
}

func NewInvalidImageErr(errorMessage string) api_error.IApiError {
	return &api_error.ApiError{
		ServerMessage: errorMessage,
		Code:          400,
		Resp: &response.ErrorResponse{
			Message: "Validation error",
			Fields: map[string]string{
				"avatar": "invalid image",
			},
		},
	}
}
