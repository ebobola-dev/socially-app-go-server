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
	DeleteMyAccount(c *fiber.Ctx) error
	Search(c *fiber.Ctx) error
	UpdateProfile(c *fiber.Ctx) error
	UpdatePassword(c *fiber.Ctx) error
	UpdateAvatar(c *fiber.Ctx) error
	DeleteAvatar(c *fiber.Ctx) error
	GetPrivileges(c *fiber.Ctx) error
	Follow(c *fiber.Ctx) error
	Unfollow(c *fiber.Ctx) error
	GetFollowers(c *fiber.Ctx) error
	GetFollowing(c *fiber.Ctx) error
}

type IAuthHandler interface {
	Login(c *fiber.Ctx) error
	Refresh(c *fiber.Ctx) error
	Logout(c *fiber.Ctx) error
}

type IPrivilegeHandler interface {
	GetAll(c *fiber.Ctx) error
	GetUsers(c *fiber.Ctx) error
	Create(c *fiber.Ctx) error
	Delete(c *fiber.Ctx) error
}

type IMediaHandler interface {
	Get(c *fiber.Ctx) error
}
