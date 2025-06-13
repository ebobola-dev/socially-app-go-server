package handler

import (
	logger "github.com/ebobola-dev/socially-app-go-server/internal/util/logger"

	"github.com/gofiber/fiber/v2"
)

type RegistrationHandler struct {
	log logger.ILogger
}

func NewRegistrationHandler(log logger.ILogger) IRegistrationHandler {
	return &RegistrationHandler{log: log}
}

func (h *RegistrationHandler) Registration(c *fiber.Ctx) error {
	payload := struct {
		Email string `json:"email"`
	}{}
	if err := c.BodyParser(&payload); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid JSON")
	}
	return c.SendStatus(fiber.StatusNotImplemented)
}

func (h *RegistrationHandler) VerifyOtp(c *fiber.Ctx) error {
	payload := struct {
		OtpId    string `json:"otp_id"`
		OtpValue string `json:"otp_value"`
	}{}
	if err := c.BodyParser(&payload); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid JSON")
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
