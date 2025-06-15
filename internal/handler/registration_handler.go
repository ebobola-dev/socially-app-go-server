package handler

import (
	"errors"
	"fmt"
	"reflect"
	"strings"

	auth_error "github.com/ebobola-dev/socially-app-go-server/internal/errors/auth"
	common_error "github.com/ebobola-dev/socially-app-go-server/internal/errors/common"
	otp_error "github.com/ebobola-dev/socially-app-go-server/internal/errors/otp"
	"github.com/ebobola-dev/socially-app-go-server/internal/middleware"
	"github.com/ebobola-dev/socially-app-go-server/internal/model"
	"github.com/ebobola-dev/socially-app-go-server/internal/repository"
	"github.com/ebobola-dev/socially-app-go-server/internal/service/email"
	jwt_s "github.com/ebobola-dev/socially-app-go-server/internal/service/jwt"
	logger "github.com/ebobola-dev/socially-app-go-server/internal/util/logger"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"gorm.io/gorm"
)

type RegistrationHandler struct {
	validate      *validator.Validate
	log           logger.ILogger
	otpRepository repository.IOtpRepository
	emailService  email.IEmailService
	jwtService    jwt_s.IJwtService
}

func NewRegistrationHandler(
	log logger.ILogger,
	validator *validator.Validate,
	otpRepository repository.IOtpRepository,
	emailService email.IEmailService,
	jwtService jwt_s.IJwtService,
) IRegistrationHandler {
	return &RegistrationHandler{
		log:           log,
		validate:      validator,
		otpRepository: otpRepository,
		emailService:  emailService,
		jwtService:    jwtService,
	}
}

func (h *RegistrationHandler) Registration(c *fiber.Ctx) error {
	type request struct {
		Email string `json:"email" validate:"required,email"`
	}
	var payload request
	if err := c.BodyParser(&payload); err != nil {
		return common_error.ErrInvalidJSON
	}
	if err := h.validate.Struct(payload); err != nil {
		return err
	}
	h.log.Debug("Email: %s", payload.Email)

	tx := middleware.GetTX(c)
	otp, get_err := h.otpRepository.GetByEmail(tx, payload.Email)
	if get_err != nil && !errors.Is(get_err, gorm.ErrRecordNotFound) {
		return get_err
	} else if errors.Is(get_err, gorm.ErrRecordNotFound) {
		otp = &model.Otp{EmailAddress: payload.Email}
		cr_err := h.otpRepository.Create(tx, otp)
		if cr_err != nil {
			return cr_err
		}
	} else {
		if can_update, delta := otp.CanUpdate(); !can_update {
			return otp_error.NewCantUpdateOtpError(delta)
		}
		otp.RegenerateCode()
		upd_err := h.otpRepository.Update(tx, otp)
		if upd_err != nil {
			return upd_err
		}
	}

	email_err := h.emailService.Send(
		payload.Email,
		"Your OTP code for registration on Socially App",
		fmt.Sprintf("OTP code: %v\nThis code is valid for 15 minutes.", otp.Value),
	)
	if email_err != nil {
		return email_err
	}
	h.log.Debug("Generated otp: %v", otp.Value)
	return c.JSON(fiber.Map{
		"id":         otp.ID,
		"created_at": otp.CreatedAt,
	})
}

func (h *RegistrationHandler) VerifyOtp(c *fiber.Ctx) error {
	type request struct {
		OtpId string         `json:"otp_id" validate:"required,uuid4"`
		Value model.OtpValue `json:"value" validate:"required,otp_value"`
	}
	var payload request
	if err := c.BodyParser(&payload); err != nil {
		return common_error.ErrInvalidJSON
	}
	if err := h.validate.Struct(payload); err != nil {
		return err
	}
	tx := middleware.GetTX(c)
	otp, get_err := h.otpRepository.GetByID(tx, payload.OtpId)
	if get_err != nil && !errors.Is(get_err, gorm.ErrRecordNotFound) {
		return get_err
	} else if errors.Is(get_err, gorm.ErrRecordNotFound) {
		return common_error.NewRecordNotFoundError("OTP code")
	}
	if !otp.IsAlive() {
		del_err := h.otpRepository.Delete(tx, payload.OtpId)
		if del_err != nil {
			h.log.Error(del_err)
		} else {
			tx.Commit()
		}
		return otp_error.ErrIsOutdated
	}
	if !reflect.DeepEqual(payload.Value, otp.Value) {
		return otp_error.ErrIncorect
	}
	token, token_err := h.jwtService.GenerateRegistration(otp.EmailAddress)
	if token_err != nil {
		return token_err
	}
	del_err := h.otpRepository.Delete(tx, otp.ID.String())
	if del_err != nil {
		return del_err
	}
	return c.JSON(fiber.Map{
		"access_token": token,
	})
}

func (h *RegistrationHandler) CompleteRegistration(c *fiber.Ctx) error {
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
	data, token_err := h.jwtService.ValidateRegistration(token)
	if token_err != nil {
		if errors.Is(token_err, jwt.ErrTokenExpired) {
			return auth_error.ErrExpired
		}
		return auth_error.ErrInvalidToken
	}
	h.log.Debug("Email: %s", data.Email)

	payload := struct {
		Fullname    string `json:"fullname" validate:"omitempty,max=32"`
		DateOfBirth string `json:"date_of_birth" validate:"required,date"`
		Gender      string `json:"gender"   validate:"omitempty,gender"`
		AboutMe     string `json:"about_me" validate:"omitempty,max=256"`
		Username    string `json:"username" validate:"required,username_length,username_charset,username_start_digit,username_start_dot"`
		Password    string `json:"password" validate:"required,password"`
	}{}
	if err := c.BodyParser(&payload); err != nil {
		return common_error.ErrInvalidJSON
	}
	if err := h.validate.Struct(payload); err != nil {
		return err
	}
	return common_error.ErrNotImplemented
}
