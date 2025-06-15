package repository

import (
	"github.com/ebobola-dev/socially-app-go-server/internal/model"
	"gorm.io/gorm"
)

type IOtpRepository interface {
	GetByEmail(db *gorm.DB, email string) (*model.Otp, error)
	GetByID(db *gorm.DB, ID string) (*model.Otp, error)
	Create(db *gorm.DB, otp *model.Otp) error
	Update(tx *gorm.DB, otp *model.Otp) error
	Delete(db *gorm.DB, id string) error
}

type OtpRepository struct{}

func NewOtpRepository() IOtpRepository {
	return &OtpRepository{}
}

func (r *OtpRepository) GetByEmail(db *gorm.DB, email string) (*model.Otp, error) {
	var otp model.Otp
	err := db.First(&otp, "email_address = ?", email).Error
	return &otp, err
}

func (r *OtpRepository) GetByID(db *gorm.DB, ID string) (*model.Otp, error) {
	var otp model.Otp
	err := db.First(&otp, "id = ?", ID).Error
	return &otp, err
}

func (r *OtpRepository) Create(db *gorm.DB, otp *model.Otp) error {
	if err := db.Create(otp).Error; err != nil {
		return err
	}
	return nil
}

func (r *OtpRepository) Update(tx *gorm.DB, otp *model.Otp) error {
	return tx.Save(otp).Error
}

func (r *OtpRepository) Delete(db *gorm.DB, id string) error {
	return db.Delete(&model.Otp{}, "id = ?", id).Error
}
