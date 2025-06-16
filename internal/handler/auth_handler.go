package handler

import (
	"errors"
	"fmt"
	"time"

	auth_error "github.com/ebobola-dev/socially-app-go-server/internal/errors/auth"
	common_error "github.com/ebobola-dev/socially-app-go-server/internal/errors/common"
	"github.com/ebobola-dev/socially-app-go-server/internal/middleware"
	"gorm.io/gorm"

	"github.com/gofiber/fiber/v2"
)

type AuthHandler struct{}

func NewAuthHandler() IAuthHandler {
	return &AuthHandler{}
}

func (h *AuthHandler) Login(c *fiber.Ctx) error {
	s := middleware.GetAppScope(c)
	payload := struct {
		Username string `json:"username" validate:"required,username_length,username_charset,username_start_digit,username_start_dot"`
		Password string `json:"password" validate:"required,password"`
	}{}
	if err := c.BodyParser(&payload); err != nil {
		return common_error.ErrInvalidJSON
	}
	if err := s.Validate.Struct(payload); err != nil {
		return err
	}

	tx := middleware.GetTX(c)
	user, get_err := s.UserRepository.GetByUsername(tx, payload.Username)
	if get_err != nil {
		if errors.Is(get_err, gorm.ErrRecordNotFound) {
			return auth_error.NewInvalidLoginData(fmt.Sprintf("user not found (@%s)", payload.Username))
		}
		return get_err
	}
	passwordIsValid := s.HashService.CheckPasswordHash(payload.Password, user.Password)
	if !passwordIsValid {
		return auth_error.NewInvalidLoginData(fmt.Sprintf("wrong password (@%s)", payload.Username))
	}
	deviceId := middleware.GetDeviceId(c)
	access_token, refresh_token, jwt_err := s.JwtService.GenerateUserPair(user.ID, deviceId)
	if jwt_err != nil {
		return jwt_err
	}
	saved_rt, get_err := s.RefreshTokenRepository.GetByUIDAndDeviceID(tx, user.ID.String(), deviceId)
	if errors.Is(get_err, gorm.ErrRecordNotFound) {
		cr_err := s.RefreshTokenRepository.Create(tx, refresh_token)
		if cr_err != nil {
			return cr_err
		}
	} else if get_err != nil {
		return get_err
	} else {
		saved_rt.Value = refresh_token.Value
		saved_rt.ExpiresAt = refresh_token.ExpiresAt
		saved_rt.CreatedAt = time.Now().UTC()
		upd_err := s.RefreshTokenRepository.Update(tx, saved_rt)
		if upd_err != nil {
			return upd_err
		}
	}
	return c.JSON(fiber.Map{
		"access_token":  access_token,
		"refresh_token": refresh_token.Value,
		"user":          user,
	})
}

func (h *AuthHandler) Logout(c *fiber.Ctx) error {
	return common_error.ErrNotImplemented
}

func (h *AuthHandler) Refresh(c *fiber.Ctx) error {
	return common_error.ErrNotImplemented
}
