package handler

import (
	"errors"

	common_error "github.com/ebobola-dev/socially-app-go-server/internal/errors/common"
	privilege_error "github.com/ebobola-dev/socially-app-go-server/internal/errors/privilege"
	"github.com/ebobola-dev/socially-app-go-server/internal/middleware"
	"github.com/ebobola-dev/socially-app-go-server/internal/model"
	"github.com/ebobola-dev/socially-app-go-server/internal/repository"
	"github.com/ebobola-dev/socially-app-go-server/internal/util/pagination"
	"github.com/google/uuid"
	"github.com/samber/lo"
	"gorm.io/gorm"

	"github.com/gofiber/fiber/v2"
)

type privilegeHandler struct{}

func NewPrivilegeHandler() IPrivilegeHandler {
	return &privilegeHandler{}
}

func (h *privilegeHandler) GetAll(c *fiber.Ctx) error {
	s := middleware.GetAppScope(c)
	pag := middleware.GetPagination(c)
	tx := middleware.GetTX(c)
	privileges, err := s.PrivilegeRepository.GetAll(tx, repository.GetPrivilegesListOptions{
		Pagination: pag,
		CountUsers: true,
	})
	if err != nil {
		return err
	}
	return c.JSON(struct {
		Count      int                      `json:"count"`
		Pagination pagination.Pagination    `json:"pagination"`
		Privileges []model.FullPrivilegeDto `json:"privileges"`
	}{
		Count:      len(privileges),
		Pagination: pag,
		Privileges: lo.Map(privileges, func(privilege model.Privilege, _ int) model.FullPrivilegeDto {
			return privilege.ToFullDto()
		}),
	})
}

func (h *privilegeHandler) GetUsers(c *fiber.Ctx) error {
	s := middleware.GetAppScope(c)
	privName := c.Query("privilege")
	tx := middleware.GetTX(c)

	privilege, err := s.PrivilegeRepository.GetByName(tx, privName, repository.GetPrivilegeOptions{
		CountUsers: true,
	})
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return common_error.NewRecordNotFoundErr("Privilege")
	} else if err != nil {
		return err
	}

	pag := middleware.GetPagination(c)
	users, err := s.PrivilegeRepository.GetUsers(tx, pag, privName)
	if err != nil {
		return err
	}
	return c.JSON(struct {
		Count      int                    `json:"count"`
		Pagination pagination.Pagination  `json:"pagination"`
		Privilege  model.FullPrivilegeDto `json:"privilege"`
		Users      []model.ShortUserDto   `json:"users"`
	}{
		Count:      len(users),
		Pagination: pag,
		Privilege:  privilege.ToFullDto(),
		Users: lo.Map(users, func(user model.User, _ int) model.ShortUserDto {
			return user.ToShortDto()
		}),
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
		"created_privilege": newPrivilege.ToFullDto(),
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
	if privilege, err := s.PrivilegeRepository.GetByID(tx, privilegeId, repository.GetPrivilegeOptions{
		CountUsers: true,
	}); errors.Is(err, gorm.ErrRecordNotFound) {
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
