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

type Privilege struct {
	ID         uuid.UUID `gorm:"type:char(36); primaryKey" serializer:"short"`
	Name       string    `gorm:"type:varchar(64); uniqueIndex" serializer:"short"`
	OrderIndex int       `gorm:"not null;default:0; uniqueIndex"  serializer:"short"`
	CreatedAt  time.Time `gorm:"autoCreateTime" serializer:"short"`

	Users      []User `gorm:"many2many:user_privileges"`
	UsersCount int    `gorm:"-"`
}

func (p *Privilege) BeforeCreate(tx *gorm.DB) (err error) {
	if p.ID == uuid.Nil {
		p.ID = uuid.New()
	}
	return
}

func (p *Privilege) ToJson(options SerializePrivilegeOptions) map[string]interface{} {
	val := reflect.ValueOf(p)
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
		if !slices.Contains(flags, "short") && options.Short {
			continue
		}

		fieldName := string_utils.ToSnakeCase(field.Name)
		if flags[0] != "short" && flags[0] != "" {
			fieldName = flags[0]
		}
		out[fieldName] = val.Field(i).Interface()
	}
	if !options.Short {
		out["users_count"] = p.UsersCount
	}
	return out
}

type SerializePrivilegeOptions struct {
	Short bool
}
