package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Privilege struct {
	ID         uuid.UUID `gorm:"type:char(36); primaryKey"`
	Name       string    `gorm:"type:varchar(64); uniqueIndex"`
	OrderIndex int       `gorm:"not null;default:0; uniqueIndex"`
	CreatedAt  time.Time `gorm:"autoCreateTime"`

	Users      []User `gorm:"many2many:user_privileges"`
	UsersCount int    `gorm:"-"`
}

func (p *Privilege) BeforeCreate(tx *gorm.DB) (err error) {
	if p.ID == uuid.Nil {
		p.ID = uuid.New()
	}
	return
}

func (p *Privilege) ToShortDto() ShortPrivilegeDto {
	return ShortPrivilegeDto{
		Id:         p.ID,
		Name:       p.Name,
		OrderIndex: p.OrderIndex,
		CreatedAt:  p.CreatedAt,
	}
}

func (p *Privilege) ToFullDto() FullPrivilegeDto {
	return FullPrivilegeDto{
		Id:         p.ID,
		Name:       p.Name,
		OrderIndex: p.OrderIndex,
		CreatedAt:  p.CreatedAt,
		UsersCount: p.UsersCount,
	}
}

type PrivilegeDto interface{}

type ShortPrivilegeDto struct {
	Id         uuid.UUID `json:"id"`
	Name       string    `json:"name"`
	OrderIndex int       `json:"order_index"`
	CreatedAt  time.Time `json:"created_at"`
}

type FullPrivilegeDto struct {
	Id         uuid.UUID `json:"id"`
	Name       string    `json:"name"`
	OrderIndex int       `json:"order_index"`
	CreatedAt  time.Time `json:"created_at"`
	UsersCount int       `json:"users_count"`
}
