package common_error

import (
	"github.com/ebobola-dev/socially-app-go-server/internal/response"
)

type ConflictType string

const (
	AlreadyFollowingConflict ConflictType = "already_following"
	NotFollowingConflict     ConflictType = "not_following"
	FollowYourselfConflict   ConflictType = "follow_yourself"
)

type сonflictError struct {
	conflictType ConflictType
}

func (e *сonflictError) Error() string {
	switch e.conflictType {
	case AlreadyFollowingConflict:
		return "Already following"
	case NotFollowingConflict:
		return "Not following anyway"
	case FollowYourselfConflict:
		return "Trying to follow/unfollow yourself"
	default:
		return "Unknown conflict"
	}
}

func (e *сonflictError) StatusCode() int {
	return 409
}

func (e *сonflictError) Response() *response.ErrorResponse {
	var responseMessage string
	switch e.conflictType {
	case AlreadyFollowingConflict:
		responseMessage = "You are already following the target user"
	case NotFollowingConflict:
		responseMessage = "You are not following the target user anyway"
	case FollowYourselfConflict:
		responseMessage = "You can't follow/unfollow yourself"
	default:
		responseMessage = "Unknown conflic"
	}
	return &response.ErrorResponse{
		Message: responseMessage,
		Fields: map[string]string{
			"conflict_type": string(e.conflictType),
		},
	}
}

var ErrAlreadyFolowingConflict = &сonflictError{conflictType: AlreadyFollowingConflict}

var ErrNotFollowingConflict = &сonflictError{conflictType: NotFollowingConflict}

var ErrFollowYourselfConflict = &сonflictError{conflictType: FollowYourselfConflict}
