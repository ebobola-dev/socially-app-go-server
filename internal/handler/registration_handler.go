package handler

import (
	"errors"

	common_error "github.com/ebobola-dev/socially-app-go-server/internal/errors/common"
	otp_error "github.com/ebobola-dev/socially-app-go-server/internal/errors/otp"
	"github.com/ebobola-dev/socially-app-go-server/internal/middleware"
	"github.com/ebobola-dev/socially-app-go-server/internal/model"
	"github.com/ebobola-dev/socially-app-go-server/internal/repository"
	logger "github.com/ebobola-dev/socially-app-go-server/internal/util/logger"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type RegistrationHandler struct {
	validate      *validator.Validate
	log           logger.ILogger
	otpRepository repository.IOtpRepository
}

func NewRegistrationHandler(log logger.ILogger, validator *validator.Validate, otpRepository repository.IOtpRepository) IRegistrationHandler {
	return &RegistrationHandler{log: log, validate: validator, otpRepository: otpRepository}
}

func (h *RegistrationHandler) Registration(c *fiber.Ctx) error {
	type request struct {
		Email string `json:"email" validate:"required,email"`
	}
	var payload request
	if err := c.BodyParser(&payload); err != nil {
		return common_error.NewInvalidJSONError("Need JSON body with 'email' string field")
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
		h.log.Debug("Not found in database")
		otp := &model.Otp{EmailAddress: payload.Email}
		cr_err := h.otpRepository.Create(tx, otp)
		if cr_err != nil {
			return cr_err
		}
		h.log.Debug("Created otp: %s", otp)
		return c.JSON(fiber.Map{
			"id":         otp.ID,
			"created_at": otp.CreatedAt,
		})
	}
	if can_update, delta := otp.CanUpdate(); !can_update {
		return otp_error.NewCantUpdateOtpError(delta)
	}
	otp.RegenerateCode()
	upd_err := h.otpRepository.Update(tx, otp)
	if upd_err != nil {
		return upd_err
	}
	h.log.Debug("Regenerated otp: %s", otp)
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
		return common_error.NewInvalidJSONError("Need JSON body with 'otp_id' string(uuid) field and 'value'([4xint]) field")
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
		}
		return otp_error.NewOtdIsOutdatedError()
	}
	return common_error.NewNotImplementedError()
}

func (h *RegistrationHandler) CompleteRegistration(c *fiber.Ctx) error {
	payload := struct {
		//?
		Fullname    string `json:"email"`
		DateOfBirth string `json:"date_of_birth"`
		Gender      string `json:"gender"`
		AboutMe     string `json:"about_me"`
		Username    string `json:"username"`
		Password    string `json:"password"`
	}{}
	if err := c.BodyParser(&payload); err != nil {
		return common_error.NewInvalidJSONError("Invalid JSON")
	}
	return common_error.NewNotImplementedError()
}
