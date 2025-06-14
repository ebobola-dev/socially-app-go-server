package handler

import (
	"github.com/ebobola-dev/socially-app-go-server/internal/model"
	"github.com/ebobola-dev/socially-app-go-server/internal/response"
	logger "github.com/ebobola-dev/socially-app-go-server/internal/util/logger"
	"github.com/ebobola-dev/socially-app-go-server/internal/validation"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type RegistrationHandler struct {
	log logger.ILogger
}

func NewRegistrationHandler(log logger.ILogger) IRegistrationHandler {
	return &RegistrationHandler{log: log}
}

func (h *RegistrationHandler) Registration(c *fiber.Ctx) error {
	type request struct {
		Email string `json:"email" validate:"required,email"`
	}
	var payload request
	if err := c.BodyParser(&payload); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.ErrorResponse{
			Message: "Need JSON body with 'email' string field",
		}.ToJSON())
	}
	validate := validator.New()
	if err := validate.Struct(payload); err != nil {
		errResp := response.ParseValidationErrors(err)
		return c.Status(fiber.StatusBadRequest).JSON(errResp.ToJSON())
	}
	return c.SendStatus(fiber.StatusNotImplemented)
}

func (h *RegistrationHandler) VerifyOtp(c *fiber.Ctx) error {
	type request struct {
		OtpId string         `json:"otp_id" validate:"required,uuid4"`
		Value model.OtpValue `json:"value" validate:"required,otp_value"`
	}
	var payload request
	if err := c.BodyParser(&payload); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.ErrorResponse{
			Message: "Need JSON body with 'otp_id' string(uuid) field and 'value'([4xint]) field",
		}.ToJSON())
	}
	validate := validator.New()
	validation.RegisterCustomValidators(validate)
	if err := validate.Struct(payload); err != nil {
		errResp := response.ParseValidationErrors(err)
		return c.Status(fiber.StatusBadRequest).JSON(errResp.ToJSON())
	}
	return c.SendStatus(fiber.StatusNotImplemented)
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
		return fiber.NewError(fiber.StatusBadRequest, "Invalid JSON")
	}

	return c.SendStatus(fiber.StatusNotImplemented)
}
