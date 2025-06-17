package handler

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
	"time"

	auth_error "github.com/ebobola-dev/socially-app-go-server/internal/errors/auth"
	common_error "github.com/ebobola-dev/socially-app-go-server/internal/errors/common"
	otp_error "github.com/ebobola-dev/socially-app-go-server/internal/errors/otp"
	user_error "github.com/ebobola-dev/socially-app-go-server/internal/errors/user"
	"github.com/ebobola-dev/socially-app-go-server/internal/middleware"
	"github.com/ebobola-dev/socially-app-go-server/internal/model"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"gorm.io/gorm"
)

type RegistrationHandler struct {
}

func NewRegistrationHandler() IRegistrationHandler {
	return &RegistrationHandler{}
}

func (h *RegistrationHandler) Registration(c *fiber.Ctx) error {
	scope := middleware.GetAppScope(c)
	validate := scope.Validate
	log := scope.Log
	userRepository := scope.UserRepository
	otpRepository := scope.OtpRepository
	emailService := scope.EmailService

	payload := struct {
		Email string `json:"email" validate:"required,email"`
	}{}
	if err := c.BodyParser(&payload); err != nil {
		return common_error.ErrInvalidJSON
	}
	if err := validate.Struct(payload); err != nil {
		return err
	}
	log.Debug("Email: %s", payload.Email)

	tx := middleware.GetTX(c)

	if exists, ex_err := userRepository.ExistsByEmail(tx, payload.Email); ex_err != nil {
		return ex_err
	} else if exists {
		return user_error.NewAlreadyRegisteredError(payload.Email)
	}

	otp, get_err := otpRepository.GetByEmail(tx, payload.Email)
	if get_err != nil && !errors.Is(get_err, gorm.ErrRecordNotFound) {
		return get_err
	} else if errors.Is(get_err, gorm.ErrRecordNotFound) {
		otp = &model.Otp{EmailAddress: payload.Email}
		cr_err := otpRepository.Create(tx, otp)
		if cr_err != nil {
			return cr_err
		}
	} else {
		if can_update, delta := otp.CanUpdate(); !can_update {
			return otp_error.NewCantUpdateOtpError(delta)
		}
		otp.RegenerateCode()
		upd_err := otpRepository.Update(tx, otp)
		if upd_err != nil {
			return upd_err
		}
	}

	email_err := emailService.Send(
		payload.Email,
		"Your OTP code for registration on Socially App",
		fmt.Sprintf("OTP code: %v\nThis code is valid for 15 minutes.", otp.Value),
	)
	if email_err != nil {
		return email_err
	}
	log.Debug("Generated otp: %v", otp.Value)
	return c.JSON(fiber.Map{
		"id":         otp.ID,
		"created_at": otp.CreatedAt,
	})
}

func (h *RegistrationHandler) VerifyOtp(c *fiber.Ctx) error {
	scope := middleware.GetAppScope(c)
	validate := scope.Validate
	userRepository := scope.UserRepository
	otpRepository := scope.OtpRepository
	jwtService := scope.JwtService

	payload := struct {
		OtpId string         `json:"otp_id" validate:"required,uuid4"`
		Value model.OtpValue `json:"value" validate:"required,otp_value"`
	}{}
	if err := c.BodyParser(&payload); err != nil {
		return common_error.ErrInvalidJSON
	}
	if err := validate.Struct(payload); err != nil {
		return err
	}
	tx := middleware.GetTX(c)
	otp, get_err := otpRepository.GetByID(tx, payload.OtpId)
	if get_err != nil && !errors.Is(get_err, gorm.ErrRecordNotFound) {
		return get_err
	} else if errors.Is(get_err, gorm.ErrRecordNotFound) {
		return common_error.NewRecordNotFoundError("OTP code")
	}
	del_err := otpRepository.Delete(tx, otp.ID.String())
	if del_err != nil {
		return del_err
	}
	if exists, ex_err := userRepository.ExistsByEmail(tx, otp.EmailAddress); ex_err != nil {
		return ex_err
	} else if exists {
		return user_error.NewAlreadyRegisteredError(otp.EmailAddress)
	}
	if !otp.IsAlive() {
		tx.Commit()
		return otp_error.ErrIsOutdated
	}
	if !reflect.DeepEqual(payload.Value, otp.Value) {
		return otp_error.ErrIncorect
	}
	token, token_err := jwtService.GenerateRegistration(otp.EmailAddress)
	if token_err != nil {
		return token_err
	}
	return c.JSON(fiber.Map{
		"access_token": token,
	})
}

func (h *RegistrationHandler) CompleteRegistration(c *fiber.Ctx) error {
	scope := middleware.GetAppScope(c)
	validate := scope.Validate
	userRepository := scope.UserRepository
	jwtService := scope.JwtService
	hashService := scope.HashService

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
	reg_data, token_err := jwtService.ValidateRegistration(token)
	if token_err != nil {
		if errors.Is(token_err, jwt.ErrTokenExpired) {
			return auth_error.ErrExpired
		}
		return auth_error.ErrInvalidToken
	}

	tx := middleware.GetTX(c)

	if exists, ex_err := userRepository.ExistsByEmail(tx, reg_data.Email); ex_err != nil {
		return ex_err
	} else if exists {
		return user_error.NewAlreadyRegisteredError(reg_data.Email)
	}

	payload := struct {
		Fullname    *string `json:"fullname" validate:"omitempty,max=32"`
		DateOfBirth string  `json:"date_of_birth" validate:"required,date"`
		Gender      *string `json:"gender"   validate:"omitempty,gender"`
		AboutMe     *string `json:"about_me" validate:"omitempty,max=256"`
		Username    string  `json:"username" validate:"required,username_length,username_charset,username_start_digit,username_start_dot"`
		Password    string  `json:"password" validate:"required,password"`
	}{}
	if err := c.BodyParser(&payload); err != nil {
		return common_error.ErrInvalidJSON
	}
	if err := validate.Struct(payload); err != nil {
		return err
	}
	gender := model.GenderFromString(payload.Gender)
	dob, _ := time.Parse("02.01.2006", payload.DateOfBirth)
	hashed_password, hash_err := hashService.HashPassword(payload.Password)
	if hash_err != nil {
		return hash_err
	}

	new_user := &model.User{
		Email:       reg_data.Email,
		Fullname:    payload.Fullname,
		DateOfBirth: dob,
		Gender:      gender,
		AboutMe:     payload.AboutMe,
		Username:    payload.Username,
		Password:    hashed_password,
	}
	cr_err := userRepository.Create(tx, new_user)
	if cr_err != nil {
		return cr_err
	}

	return c.JSON(fiber.Map{
		"created_user": new_user,
	})
}
