package model

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type OtpValue []int

type Otp struct {
	ID           uuid.UUID `gorm:"type:char(36); primaryKey" json:"id"`
	EmailAddress string    `gorm:"type:varchar(255); uniqueIndex" json:"email"`
	Value        OtpValue  `gorm:"type:json" json:"value"`
	CreatedAt    time.Time `gorm:"autoCreateTime" json:"created_at"`
}

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

// func (o *OtpValue) UnmarshalJSON(data []byte) error {
// 	var arr []int
// 	if err := json.Unmarshal(data, &arr); err != nil {
// 		return err
// 	}
// 	if len(arr) != 4 {
// 		return fmt.Errorf("otp_value must have 4 digits")
// 	}
// 	for _, digit := range arr {
// 		if digit < 0 || digit > 9 {
// 			return fmt.Errorf("otp_value digits must be between 0 and 9")
// 		}
// 	}
// 	copy(o[:], arr)
// 	return nil
// }

func (o *Otp) BeforeCreate(tx *gorm.DB) (err error) {
	if o.ID == uuid.Nil {
		o.ID = uuid.New()
	}
	return
}
