package handler

import (
	"errors"

	common_error "github.com/ebobola-dev/socially-app-go-server/internal/errors/common"
	"github.com/ebobola-dev/socially-app-go-server/internal/middleware"
	"github.com/ebobola-dev/socially-app-go-server/internal/repository"
	logger "github.com/ebobola-dev/socially-app-go-server/internal/util/logger"
	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"

	"github.com/gofiber/fiber/v2"
)

type UserHandler struct {
	log            logger.ILogger
	validate       *validator.Validate
	userRepository repository.IUserRepository
}

func NewUserHandler(log logger.ILogger, validate *validator.Validate, userRepository repository.IUserRepository) IUserHandler {
	return &UserHandler{
		log:            log,
		validate:       validate,
		userRepository: userRepository,
	}
}

func (h *UserHandler) CheckUsername(c *fiber.Ctx) error {
	payload := struct {
		Username string `validate:"required,username_length,username_charset,username_start_digit,username_start_dot"`
	}{
		Username: c.Query("username"),
	}
	if err := h.validate.Struct(payload); err != nil {
		return err
	}

	tx := middleware.GetTX(c)
	exists, ex_err := h.userRepository.ExistsByUsername(tx, payload.Username)
	if ex_err != nil {
		return ex_err
	}
	return c.JSON(fiber.Map{
		"username": payload.Username,
		"exists":   exists,
	})
}

func (h *UserHandler) GetById(c *fiber.Ctx) error {
	payload := struct {
		UserId string `validate:"required,uuid4"`
	}{
		UserId: c.Params("user_id"),
	}
	if err := h.validate.Struct(payload); err != nil {
		return err
	}
	tx := middleware.GetTX(c)
	user, get_err := h.userRepository.GetByID(tx, payload.UserId)
	if get_err != nil && !errors.Is(get_err, gorm.ErrRecordNotFound) {
		return get_err
	} else if errors.Is(get_err, gorm.ErrRecordNotFound) {
		return common_error.NewRecordNotFoundError("User")
	}
	return c.JSON(user)
}
