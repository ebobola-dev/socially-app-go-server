package handler

import (
	"errors"

	common_error "github.com/ebobola-dev/socially-app-go-server/internal/errors/common"
	"github.com/ebobola-dev/socially-app-go-server/internal/middleware"
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
	tx := middleware.GetTX(c)
	user, get_err := s.UserRepository.GetByID(tx, payload.UserId)
	if get_err != nil && !errors.Is(get_err, gorm.ErrRecordNotFound) {
		return get_err
	} else if errors.Is(get_err, gorm.ErrRecordNotFound) {
		return common_error.NewRecordNotFoundError("User")
	}
	return c.JSON(user)
}
