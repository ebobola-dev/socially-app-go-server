package repository

import (
	"github.com/ebobola-dev/socially-app-go-server/internal/model"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type IOtpRepository interface {
	GetByEmail(db *gorm.DB, email string) (*model.Otp, error)
	GetByID(db *gorm.DB, id uuid.UUID) (*model.Otp, error)
	Create(db *gorm.DB, otp *model.Otp) error
	Update(tx *gorm.DB, otp *model.Otp) error
	Delete(db *gorm.DB, id uuid.UUID) error
}

type otpRepository struct{}

func NewOtpRepository() IOtpRepository {
	return &otpRepository{}
}

func (r *otpRepository) GetByEmail(db *gorm.DB, email string) (*model.Otp, error) {
	var otp model.Otp
	err := db.First(&otp, "email_address = ?", email).Error
	return &otp, err
}

func (r *otpRepository) GetByID(db *gorm.DB, id uuid.UUID) (*model.Otp, error) {
	var otp model.Otp
	err := db.First(&otp, "id = ?", id).Error
	return &otp, err
}

func (r *otpRepository) Create(db *gorm.DB, otp *model.Otp) error {
	if err := db.Create(otp).Error; err != nil {
		return err
	}
	return nil
}

func (r *otpRepository) Update(tx *gorm.DB, otp *model.Otp) error {
	return tx.Save(otp).Error
}

func (r *otpRepository) Delete(db *gorm.DB, id uuid.UUID) error {
	result := db.Delete(&model.Otp{}, "id = ?", id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}
