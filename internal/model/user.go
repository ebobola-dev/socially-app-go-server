package model

import (
	"reflect"
	"slices"
	"strings"
	"time"

	string_utils "github.com/ebobola-dev/socially-app-go-server/internal/util/strings"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	ID          uuid.UUID `gorm:"type:char(36);primaryKey" serializer:"short"`
	Email       string    `gorm:"uniqueIndex;not null" serializer:"safe"`
	Username    string    `gorm:"uniqueIndex;type:varchar(16);not null" serializer:"short"`
	Password    string    `gorm:"type:char(60), not null"`
	Fullname    *string   `gorm:"type:varchar(32)" serializer:"short"`
	AboutMe     *string   `gorm:"type:varchar(256)" serializer:""`
	Gender      *Gender   `gorm:"type:enum('male','female')" serializer:""`
	DateOfBirth time.Time `gorm:"type:date;not null" serializer:""`

	AvatarType *AvatarType `gorm:"type:enum('external','avatar1','avatar2', 'avatar3', 'avatar4', 'avatar5', 'avatar6', 'avatar7', 'avatar8', 'avatar9', 'avatar10');" serializer:"short"`
	AvatarID   *uuid.UUID  `gorm:"type:char(36);uniqueIndex" serializer:"short"`

	Privileges []Privilege `gorm:"many2many:user_privileges"`

	LastSeen *time.Time `serializer:""`

	DeletedAt *time.Time `gorm:"index" serializer:"short"`
	CreatedAt time.Time  `gorm:"autoCreateTime" serializer:"short"`
}

func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
	if u.ID == uuid.Nil {
		u.ID = uuid.New()
	}
	return
}

func (u *User) ToJson(options SerializeUserOptions) map[string]interface{} {
	val := reflect.ValueOf(u)
	if val.Kind() == reflect.Pointer {
		val = val.Elem()
	}
	typ := val.Type()

	out := map[string]interface{}{}

	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)
		tag, ok := field.Tag.Lookup("serializer")
		if !ok {
			continue
		}
		flags := strings.Split(tag, ",")
		if slices.Contains(flags, "safe") && !options.Safe {
			continue
		}
		if !slices.Contains(flags, "short") && options.Short {
			continue
		}

		fieldName := string_utils.ToSnakeCase(field.Name)
		if !slices.Contains([]string{"safe", "short"}, flags[0]) && flags[0] != "" {
			fieldName = flags[0]
		}
		out[fieldName] = val.Field(i).Interface()
	}
	privileges := make([]string, len(u.Privileges))
	for i, privilege := range u.Privileges {
		privileges[i] = privilege.Name
	}
	if !options.Short {
		out["privileges"] = privileges
	}
	return out
}

type SerializeUserOptions struct {
	Safe  bool
	Short bool
}
