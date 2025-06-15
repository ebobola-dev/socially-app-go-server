package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	ID          uuid.UUID `gorm:"type:char(36);primaryKey" json:"id"`
	Email       string    `gorm:"uniqueIndex;not null" json:"email_address"`
	Username    string    `gorm:"uniqueIndex;type:varchar(16);not null" json:"username"`
	Password    string    `gorm:"type:varchar(16), not null" json:"-"`
	Fullname    *string   `gorm:"type:varchar(32)" json:"fullname,omitempty"`
	AboutMe     *string   `gorm:"type:varchar(256)" json:"about_me,omitempty"`
	Gender      *Gender   `gorm:"type:enum('male','female')" json:"gender,omitempty"`
	DateOfBirth time.Time `gorm:"type:date;not null" json:"date_of_birth"`

	AvatarType *AvatarType `gorm:"type:enum('external','avatar1','avatar2', 'avatar3', 'avatar4', 'avatar5', 'avatar6', 'avatar7', 'avatar8', 'avatar9', 'avatar10');" json:"avatar_type,omitempty"`
	AvatarID   *uuid.UUID  `gorm:"type:char(36);uniqueIndex" json:"avatar_id,omitempty"`

	Privileges []Privilege `gorm:"many2many:user_privileges" json:"privileges"`

	LastSeen *time.Time `json:"last_seen,omitempty"`

	DeletedAt *time.Time `gorm:"index" json:"deleted_at,omitempty"`
	CreatedAt time.Time  `gorm:"autoCreateTime" json:"created_at"`
}

func (o *User) BeforeCreate(tx *gorm.DB) (err error) {
	o.ID = uuid.New()
	return
}
