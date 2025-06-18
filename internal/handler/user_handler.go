package handler

import (
	"errors"
	"time"

	auth_error "github.com/ebobola-dev/socially-app-go-server/internal/errors/auth"
	common_error "github.com/ebobola-dev/socially-app-go-server/internal/errors/common"
	user_error "github.com/ebobola-dev/socially-app-go-server/internal/errors/user"
	"github.com/ebobola-dev/socially-app-go-server/internal/middleware"
	"github.com/ebobola-dev/socially-app-go-server/internal/model"
	"github.com/ebobola-dev/socially-app-go-server/internal/util/nullable"
	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/gofiber/fiber/v2"
)

type UserHandler struct{}

func NewUserHandler() IUserHandler {
	return &UserHandler{}
}

func (h *UserHandler) CheckUsername(c *fiber.Ctx) error {
	scope := middleware.GetAppScope(c)
	payload := struct {
		Username string `validate:"required,username_length,username_charset,username_start_digit,username_start_dot"`
	}{
		Username: c.Query("username"),
	}
	if err := scope.Validate.Struct(payload); err != nil {
		return err
	}

	tx := middleware.GetTX(c)
	exists, ex_err := scope.UserRepository.ExistsByUsername(tx, payload.Username)
	if ex_err != nil {
		return ex_err
	}
	return c.JSON(fiber.Map{
		"username": payload.Username,
		"exists":   exists,
	})
}

func (h *UserHandler) GetById(c *fiber.Ctx) error {
	s := middleware.GetAppScope(c)

	payload := struct {
		UserId string `validate:"required,uuid4"`
	}{
		UserId: c.Params("user_id"),
	}
	if err := s.Validate.Struct(payload); err != nil {
		return err
	}
	userId := uuid.MustParse(payload.UserId)
	tx := middleware.GetTX(c)
	user, get_err := s.UserRepository.GetByID(tx, userId, true)
	if get_err != nil && !errors.Is(get_err, gorm.ErrRecordNotFound) {
		return get_err
	} else if errors.Is(get_err, gorm.ErrRecordNotFound) {
		return common_error.NewRecordNotFoundError("User")
	}
	return c.JSON(user)
}

func (h *UserHandler) DeleteMyAccount(c *fiber.Ctx) error {
	s := middleware.GetAppScope(c)
	userId := middleware.GetUserId(c)
	tx := middleware.GetTX(c)

	//% Soft delete user
	if err := s.UserRepository.SoftDelete(tx, userId); errors.Is(err, gorm.ErrRecordNotFound) {
		return common_error.NewRecordNotFoundError("User")
	} else if err != nil {
		return err
	}

	//% Delete refresh tokens
	if _, err := s.RefreshTokenRepository.DeleteByUserId(tx, userId); errors.Is(err, gorm.ErrRecordNotFound) {
		return auth_error.ErrInvalidToken
	} else if err != nil {
		return err
	}

	tx.Commit()
	return auth_error.ErrAccountDeleted
}

func (h *UserHandler) Search(c *fiber.Ctx) error {
	tx := middleware.GetTX(c)
	s := middleware.GetAppScope(c)
	userId := middleware.GetUserId(c)
	pagination := middleware.GetPagination(c)

	pattern := c.Query("pattern")
	users, err := s.UserRepository.Search(tx, pagination, pattern)
	if err != nil {
		return err
	}
	filteredUsers := make([]model.User, 0, len(users))
	for _, user := range users {
		if user.ID != userId {
			filteredUsers = append(filteredUsers, user)
		}
	}
	return c.JSON(fiber.Map{
		"pagination": fiber.Map{
			"offset": pagination.Offset,
			"limit":  pagination.Limit,
		},
		"count":   len(filteredUsers),
		"pattern": pattern,
		"users":   filteredUsers,
	})
}

func (h *UserHandler) UpdateProfile(c *fiber.Ctx) error {
	s := middleware.GetAppScope(c)
	payload := struct {
		Fullname    *string                         `json:"fullname" validate:"omitempty,max=32"`
		Username    *string                         `json:"username" validate:"omitempty,username_length,username_charset,username_start_digit,username_start_dot"`
		Gender      nullable.Nullable[model.Gender] `json:"gender" validate:"omitempty,gender"`
		DateOfBirth *string                         `json:"date_of_birth" validate:"omitempty,datebt"`
		AboutMe     *string                         `json:"about_me" validate:"omitempty,max=256"`
	}{}
	if err := c.BodyParser(&payload); err != nil {
		return common_error.ErrInvalidJSON
	}
	if err := s.Validate.Struct(payload); err != nil {
		return err
	}
	tx := middleware.GetTX(c)
	userId := middleware.GetUserId(c)
	user, _ := s.UserRepository.GetByID(tx, userId, false)

	hasUpdates := false
	if payload.Fullname != nil && !nullable.StringEqual(payload.Fullname, user.Fullname) {
		user.Fullname = payload.Fullname
		hasUpdates = true
	}
	if payload.Username != nil && *payload.Username != user.Username {
		user.Username = *payload.Username
		hasUpdates = true
	}
	if payload.Gender.Present && !nullable.StringEqual((*string)(payload.Gender.Value), (*string)(user.Gender)) {
		user.Gender = payload.Gender.Value
		hasUpdates = true
	}
	if payload.DateOfBirth != nil {
		dob, _ := time.Parse("02.01.2006", *payload.DateOfBirth)
		if dob != user.DateOfBirth {
			user.DateOfBirth = dob
			hasUpdates = true
		}
	}
	if payload.AboutMe != nil && !nullable.StringEqual(payload.AboutMe, user.AboutMe) {
		user.AboutMe = payload.AboutMe
		hasUpdates = true
	}
	if !hasUpdates {
		return user_error.ErrNothingToUpdateProfile
	}
	if err := s.UserRepository.Update(tx, user); err != nil {
		return err
	}
	return c.JSON(fiber.Map{
		"updated_user": user,
	})
}
