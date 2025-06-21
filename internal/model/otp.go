package model

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"math/rand"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type OtpValue []int

func (o *OtpValue) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("OtpValue: failed to type assert value to []byte")
	}
	return json.Unmarshal(bytes, o)
}

func (o OtpValue) Value() (driver.Value, error) {
	return json.Marshal(o)
}

type Otp struct {
	ID           uuid.UUID `gorm:"type:char(36); primaryKey" json:"id"`
	EmailAddress string    `gorm:"type:varchar(255); uniqueIndex" json:"-"`
	Value        OtpValue  `gorm:"type:json" json:"-"`
	CreatedAt    time.Time `gorm:"not null;autoCreateTime(3)" json:"created_at"`
}

func (Otp) TableName() string {
	return "otp"
}

func (o *Otp) BeforeCreate(tx *gorm.DB) (err error) {
	if o.ID == uuid.Nil {
		o.ID = uuid.New()
	}
	if len(o.Value) == 0 {
		o.Value = GenerateOtpValue()
	}
	return
}

func GenerateOtpValue() OtpValue {
	return OtpValue{
		rand.Intn(10),
		rand.Intn(10),
		rand.Intn(10),
		rand.Intn(10),
	}
}

func (o *Otp) IsAlive() bool {
	return time.Since(o.CreatedAt) < 15*time.Minute
}

func (o *Otp) CanUpdate() (bool, int) {
	delta := time.Since(o.CreatedAt)
	seconds_delta := int(delta.Seconds())
	return delta > 1*time.Minute, seconds_delta
}

func (o *Otp) RegenerateCode() {
	o.Value = GenerateOtpValue()
	o.CreatedAt = time.Now().UTC()
}
