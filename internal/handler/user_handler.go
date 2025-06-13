package handler

import (
	logger "socially-app/internal/util"

	"github.com/gofiber/fiber/v2"
)

type UserHandler struct {
	log logger.ILogger
}

func NewUserHandler(log logger.ILogger) IUserHandler {
	return &UserHandler{log: log}
}

func (h *UserHandler) CheckUsername(c *fiber.Ctx) error {
	username := c.Query("username")

	h.log.Debug("Check username: %s", username)
	return c.JSON([]fiber.Map{
		{"username": username, "exists": false},
	})
}
