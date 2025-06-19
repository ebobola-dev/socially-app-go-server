package handler

import (
	"errors"

	common_error "github.com/ebobola-dev/socially-app-go-server/internal/errors/common"
	privilege_error "github.com/ebobola-dev/socially-app-go-server/internal/errors/privilege"
	"github.com/ebobola-dev/socially-app-go-server/internal/middleware"
	"github.com/ebobola-dev/socially-app-go-server/internal/model"
	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/gofiber/fiber/v2"
)

type privilegeHandler struct{}

func NewPrivilegeHandler() IPrivilegeHandler {
	return &privilegeHandler{}
}

func (h *privilegeHandler) GetAll(c *fiber.Ctx) error {
	s := middleware.GetAppScope(c)
	pagintation := middleware.GetPagination(c)
	tx := middleware.GetTX(c)
	privileges, err := s.PrivilegeRepository.GetAll(tx, pagintation)
	if err != nil {
		return err
	}
	return c.JSON(fiber.Map{
		"pagination": fiber.Map{
			"offset": pagintation.Offset,
			"limit":  pagintation.Limit,
		},
		"count":      len(privileges),
		"privileges": privileges,
	})
}

func (h *privilegeHandler) GetUsers(c *fiber.Ctx) error {
	s := middleware.GetAppScope(c)
	privName := c.Query("privilege")
	tx := middleware.GetTX(c)

	privilege, err := s.PrivilegeRepository.GetByName(tx, privName)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return common_error.NewRecordNotFoundErr("Privilege")
	} else if err != nil {
		return err
	}

	pagintation := middleware.GetPagination(c)
	users, err := s.PrivilegeRepository.GetUsers(tx, pagintation, privName)
	if err != nil {
		return err
	}
	return c.JSON(fiber.Map{
		"pagination": fiber.Map{
			"offset": pagintation.Offset,
			"limit":  pagintation.Limit,
		},
		"privilege": privilege,
		"count":     len(users),
		"users":     users,
	})
}

func (h *privilegeHandler) Create(c *fiber.Ctx) error {
	s := middleware.GetAppScope(c)
	payload := struct {
		Name       string `json:"name" validate:"required,min=1,max=64"`
		OrderIndex int    `json:"order_index" validate:"required,gt=0,lt=100"`
	}{}
	if err := c.BodyParser(&payload); err != nil {
		return privilege_error.ErrBadCreateJson
	}
	if err := s.Validate.Struct(payload); err != nil {
		return err
	}
	newPrivilege := model.Privilege{
		Name:       payload.Name,
		OrderIndex: payload.OrderIndex,
	}
	tx := middleware.GetTX(c)
	if err := s.PrivilegeRepository.Create(tx, &newPrivilege); err != nil {
		return err
	}
	return c.JSON(fiber.Map{
		"created_privilege": newPrivilege,
	})
}

func (h *privilegeHandler) Delete(c *fiber.Ctx) error {
	s := middleware.GetAppScope(c)
	payload := struct {
		Id string `json:"id" validate:"required,uuid4"`
	}{
		Id: c.Query("id"),
	}
	if err := s.Validate.Struct(payload); err != nil {
		return err
	}
	privilegeId := uuid.MustParse(payload.Id)
	tx := middleware.GetTX(c)
	if privilege, err := s.PrivilegeRepository.GetByID(tx, privilegeId); errors.Is(err, gorm.ErrRecordNotFound) {
		return common_error.NewRecordNotFoundErr("Privilege")
	} else if err != nil {
		return err
	} else if privilege.OrderIndex == 100 {
		return privilege_error.ErrDeletingOwner
	}
	if err := s.PrivilegeRepository.Delete(tx, privilegeId); errors.Is(err, gorm.ErrRecordNotFound) {
		return common_error.NewRecordNotFoundErr("Privilege")
	} else if err != nil {
		return err
	}

	return c.SendStatus(200)
}
