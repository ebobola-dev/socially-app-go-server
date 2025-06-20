package model

import (
	"time"

	"github.com/google/uuid"
)

type UserPrivilege struct {
	UserID      uuid.UUID `gorm:"type:char(36);primaryKey"`
	PrivilegeID uuid.UUID `gorm:"type:char(36);primaryKey"`
	GrantedAt   time.Time `gorm:"column:created_at;precision:3;autoCreateTime"`

	User      User      `gorm:"foreignKey:UserID;references:ID"`
	Privilege Privilege `gorm:"foreignKey:PrivilegeID;references:ID"`
}

func (UserPrivilege) TableName() string {
	return "user_privileges"
}

func (up *UserPrivilege) ToJsonPrivilege() map[string]interface{} {
	jsonPrivilege := up.Privilege.ToJson(SerializePrivilegeOptions{})
	jsonPrivilege["granted_at"] = up.GrantedAt
	return jsonPrivilege
}
