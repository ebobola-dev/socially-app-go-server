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
	GetById(c *fiber.Ctx) error
}

type IAuthHandler interface {
	Login(c *fiber.Ctx) error
	Logout(c *fiber.Ctx) error
	Refresh(c *fiber.Ctx) error
}
