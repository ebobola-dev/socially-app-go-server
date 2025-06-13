package handler

import (
	"github.com/gofiber/fiber/v2"
)

type IRegistrationHandler interface {
	Registration(c *fiber.Ctx) error
	VerifyOtp(c *fiber.Ctx) error
	CompleteRegistration(c *fiber.Ctx) error
}

type IUserHandler interface {
	CheckUsername(c *fiber.Ctx) error
}
