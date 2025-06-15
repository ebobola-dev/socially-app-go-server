package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Privilege struct {
	ID         uuid.UUID `gorm:"type:char(36); primaryKey" json:"id"`
	Name       string    `gorm:"type:varchar(64); uniqueIndex" json:"name"`
	OrderIndex int       `gorm:"not null;default:0"  json:"order_index"`
	CreatedAt  time.Time `gorm:"autoCreateTime" json:"created_at"`
}

func (o *Privilege) BeforeCreate(tx *gorm.DB) (err error) {
	o.ID = uuid.New()
	return
}
