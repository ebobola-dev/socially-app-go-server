package model

import (
	"time"

	"github.com/google/uuid"
)

type UserPrivilege struct {
	UserID      uuid.UUID `gorm:"type:char(36);primaryKey"`
	PrivilegeID uuid.UUID `gorm:"type:char(36);primaryKey"`
	CreatedAt   time.Time `gorm:"not null;autoCreateTime(3)"`

	User      User      `gorm:"foreignKey:UserID;references:ID"`
	Privilege Privilege `gorm:"foreignKey:PrivilegeID;references:ID"`
}

func (UserPrivilege) TableName() string {
	return "user_privileges"
}

func (up *UserPrivilege) ToDto() UserPrivilegeDto {
	return UserPrivilegeDto{
		GrantedAt: up.CreatedAt,
		Privilege: up.Privilege.ToFullDto(),
	}
}

type UserPrivilegeDto struct {
	GrantedAt time.Time        `json:"granted_at"`
	Privilege FullPrivilegeDto `json:"privilege"`
}
