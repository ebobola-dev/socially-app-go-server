package model

import (
	"time"

	"github.com/google/uuid"
	"github.com/samber/lo"
	"gorm.io/gorm"
)

type User struct {
	ID          uuid.UUID `gorm:"type:char(36);primaryKey"`
	Email       string    `gorm:"uniqueIndex;not null"`
	Username    string    `gorm:"uniqueIndex;type:varchar(16);not null"`
	Password    string    `gorm:"type:char(60), not null"`
	Fullname    *string   `gorm:"type:varchar(32)"`
	AboutMe     *string   `gorm:"type:varchar(256)"`
	Gender      *Gender   `gorm:"type:enum('male','female')"`
	DateOfBirth time.Time `gorm:"type:date;not null"`

	AvatarType *AvatarType `gorm:"type:enum('external','avatar1','avatar2', 'avatar3', 'avatar4', 'avatar5', 'avatar6', 'avatar7', 'avatar8', 'avatar9', 'avatar10');"`
	AvatarID   *uuid.UUID  `gorm:"type:char(36);uniqueIndex"`

	Privileges     []Privilege         `gorm:"many2many:user_privileges"`
	UserPrivileges []UserPrivilege     `gorm:"foreignKey:UserID"`
	Following      []*UserSubscription `gorm:"foreignKey:FollowerID"`
	Followers      []*UserSubscription `gorm:"foreignKey:TargetID"`

	LastSeen *time.Time `gorm:"default:null"`

	DeletedAt *time.Time `gorm:"index"`
	CreatedAt time.Time  `gorm:"autoCreateTime(3)"`

	FollowersCount int64 `gorm:"-"`
	FollowingCount int64 `gorm:"-"`
}

func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
	if u.ID == uuid.Nil {
		u.ID = uuid.New()
	}
	return
}

func (u *User) ToShortDto() ShortUserDto {
	return ShortUserDto{
		Id:         u.ID,
		Username:   u.Username,
		Fullname:   u.Fullname,
		AvatarType: u.AvatarType,
		AvatarId:   u.AvatarID,
		CreatedAt:  u.CreatedAt,
		DeletedAt:  u.DeletedAt,
	}
}

func (u *User) ToFullDto(safe bool) FullUserDto {
	jsonView := FullUserDto{
		Id:             u.ID,
		Username:       u.Username,
		Fullname:       u.Fullname,
		AvatarType:     u.AvatarType,
		AvatarId:       u.AvatarID,
		CreatedAt:      u.CreatedAt,
		DeletedAt:      u.DeletedAt,
		AboutMe:        u.AboutMe,
		Gender:         u.Gender,
		DateOfBirth:    u.DateOfBirth.Format(time.DateOnly),
		LastSeen:       u.LastSeen,
		FollowersCount: u.FollowersCount,
		FollowingCount: u.FollowingCount,
	}
	if safe {
		jsonView.Email = &u.Email
	}
	jsonView.Privileges = lo.Map(u.UserPrivileges, func(up UserPrivilege, _ int) string {
		return up.Privilege.Name
	})
	return jsonView
}

func (u *User) ToDto(options SerializeUserOptions) UserDto {
	if options.Short {
		return u.ToShortDto()
	}
	return u.ToFullDto(options.Safe)
}

type SerializeUserOptions struct {
	Safe  bool
	Short bool
}

type UserDto interface{}

type ShortUserDto struct {
	Id         uuid.UUID   `json:"id"`
	Username   string      `json:"username"`
	Fullname   *string     `json:"fullname"`
	AvatarType *AvatarType `json:"avatar_type"`
	AvatarId   *uuid.UUID  `json:"avatar_id"`
	CreatedAt  time.Time   `json:"created_at"`
	DeletedAt  *time.Time  `json:"deleted_at"`
}

type FullUserDto struct {
	Id             uuid.UUID   `json:"id"`
	Email          *string     `json:"email,omitempty"`
	Username       string      `json:"username"`
	Fullname       *string     `json:"fullname"`
	Gender         *Gender     `json:"gender"`
	DateOfBirth    string      `json:"date_of_birth"`
	AvatarType     *AvatarType `json:"avatar_type"`
	AvatarId       *uuid.UUID  `json:"avatar_id"`
	AboutMe        *string     `json:"about_me"`
	LastSeen       *time.Time  `json:"last_seen"`
	CreatedAt      time.Time   `json:"created_at"`
	DeletedAt      *time.Time  `json:"deleted_at"`
	FollowersCount int64       `json:"followers_count"`
	FollowingCount int64       `json:"following_count"`
	Privileges     []string    `json:"privileges"`
}
