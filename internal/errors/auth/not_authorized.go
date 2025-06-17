package auth_error

import (
	"fmt"

	api_error "github.com/ebobola-dev/socially-app-go-server/internal/errors"
	"github.com/ebobola-dev/socially-app-go-server/internal/response"
)

type NotAuthorizedError struct {
	msg  string
	code int
	resp *response.ErrorResponse
}

func (e *NotAuthorizedError) Error() string {
	return e.msg
}

func (e *NotAuthorizedError) StatusCode() int {
	return e.code
}

func (e *NotAuthorizedError) Response() *response.ErrorResponse {
	return e.resp
}

var ErrMissingHeader = &NotAuthorizedError{
	msg:  "Missing header",
	code: 401,
	resp: &response.ErrorResponse{
		Message: "You are not authorized",
		Fields: map[string]string{
			"type": "header_missing",
		},
	},
}

var ErrWrongFormat = &NotAuthorizedError{
	msg:  "Wrong Format",
	code: 401,
	resp: &response.ErrorResponse{
		Message: "You are not authorized",
		Fields: map[string]string{
			"type": "wrong_format",
		},
	},
}

var ErrNoToken = &NotAuthorizedError{
	msg:  "Token is empty",
	code: 401,
	resp: &response.ErrorResponse{
		Message: "You are not authorized",
		Fields: map[string]string{
			"type": "no_token",
		},
	},
}

var ErrExpired = &NotAuthorizedError{
	msg:  "Token expired",
	code: 401,
	resp: &response.ErrorResponse{
		Message: "You are not authorized",
		Fields: map[string]string{
			"type": "token_expired",
		},
	},
}

var ErrInvalidToken = &NotAuthorizedError{
	msg:  "Invalid token",
	code: 401,
	resp: &response.ErrorResponse{
		Message: "You are not authorized",
		Fields: map[string]string{
			"type": "invalid_token",
		},
	},
}

var ErrLoggetOut = &NotAuthorizedError{
	msg:  "Logget out",
	code: 401,
	resp: &response.ErrorResponse{
		Message: "You are not authorized",
		Fields: map[string]string{
			"type": "logged_out",
		},
	},
}

var ErrAccountDeleted = &NotAuthorizedError{
	msg:  "Account was deleted",
	code: 401,
	resp: &response.ErrorResponse{
		Message: "You are not authorized",
		Fields: map[string]string{
			"type": "account_deleted",
		},
	},
}

func NewUserNotFoundError(userId string) api_error.ApiError {
	return &NotAuthorizedError{
		msg:  fmt.Sprintf("!User with id in token not found! (%s)", userId),
		code: 401,
		resp: &response.ErrorResponse{
			Message: "You are not authorized",
			Fields: map[string]string{
				"type": "user not found",
			},
		},
	}
}
