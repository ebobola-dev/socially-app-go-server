package model

import (
	"time"
)

type UserPrivilege struct {
	UserID      string    `gorm:"type:char(36);primaryKey"`
	PrivilegeID string    `gorm:"type:char(36);primaryKey"`
	CreatedAt   time.Time `gorm:"precision:3;autoCreateTime"`

	User      *User      `gorm:"foreignKey:UserID;references:ID"`
	Privilege *Privilege `gorm:"foreignKey:PrivilegeID;references:ID"`
}

func (UserPrivilege) TableName() string {
	return "user_privileges"
}
