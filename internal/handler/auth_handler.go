package handler

import (
	"errors"
	"fmt"
	"strings"
	"time"

	auth_error "github.com/ebobola-dev/socially-app-go-server/internal/errors/auth"
	common_error "github.com/ebobola-dev/socially-app-go-server/internal/errors/common"
	"github.com/ebobola-dev/socially-app-go-server/internal/middleware"
	"github.com/golang-jwt/jwt/v5"
	"gorm.io/gorm"

	"github.com/gofiber/fiber/v2"
)

type authHandler struct{}

func NewAuthHandler() IAuthHandler {
	return &authHandler{}
}

func (h *authHandler) Login(c *fiber.Ctx) error {
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
	saved_rt, get_err := s.RefreshTokenRepository.GetByUIDAndDeviceID(tx, user.ID, deviceId)
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

func (h *authHandler) Refresh(c *fiber.Ctx) error {
	s := middleware.GetAppScope(c)

	authHeaders := c.Request().Header.Peek("Authorization")
	if len(authHeaders) == 0 {
		return auth_error.ErrMissingHeader
	}
	headerValue := string(authHeaders)
	parts := strings.SplitN(headerValue, " ", 2)
	if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
		return auth_error.ErrWrongFormat
	}
	token := parts[1]
	if token == "" {
		return auth_error.ErrNoToken
	}
	userData, token_err := s.JwtService.ValidateUserRefresh(token)
	if token_err != nil {
		if errors.Is(token_err, jwt.ErrTokenExpired) {
			return auth_error.ErrExpired
		}
		return auth_error.ErrInvalidToken
	}
	userId := userData.ID

	tx := middleware.GetTX(c)

	user, getUErr := s.UserRepository.GetByID(tx, userId, false)
	if errors.Is(getUErr, gorm.ErrRecordNotFound) {
		return auth_error.ErrInvalidToken
	} else if getUErr != nil {
		return getUErr
	}

	deviceId := middleware.GetDeviceId(c)
	savedRefresh, getTErr := s.RefreshTokenRepository.GetByValue(tx, token)
	if errors.Is(getTErr, gorm.ErrRecordNotFound) {
		return auth_error.ErrInvalidToken
	} else if getTErr != nil {
		return getTErr
	}
	newAccessToken, newRefreshToken, jwtErr := s.JwtService.GenerateUserPair(userId, deviceId)
	if jwtErr != nil {
		return jwtErr
	}
	savedRefresh.Value = newRefreshToken.Value
	savedRefresh.ExpiresAt = newRefreshToken.ExpiresAt
	savedRefresh.CreatedAt = time.Now().UTC()
	upd_err := s.RefreshTokenRepository.Update(tx, savedRefresh)
	if upd_err != nil {
		return upd_err
	}

	return c.JSON(fiber.Map{
		"access_token":  newAccessToken,
		"refresh_token": newRefreshToken.Value,
		"user":          user,
	})
}

func (h *authHandler) Logout(c *fiber.Ctx) error {
	s := middleware.GetAppScope(c)
	tx := middleware.GetTX(c)
	deviceId := middleware.GetDeviceId(c)
	userId := middleware.GetUserId(c)
	if delErr := s.RefreshTokenRepository.DeleteByUIDAndDeviceID(tx, userId, deviceId); errors.Is(delErr, gorm.ErrRecordNotFound) {
		return auth_error.ErrInvalidToken
	} else if delErr != nil {
		return delErr
	}
	tx.Commit()
	return auth_error.ErrLoggetOut
}
