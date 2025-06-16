package repository

import (
	"github.com/ebobola-dev/socially-app-go-server/internal/model"
	"gorm.io/gorm"
)

type IRefreshTokenRepository interface {
	GetByID(db *gorm.DB, ID string) (*model.RefreshToken, error)
	GetByUIDAndDeviceID(db *gorm.DB, userId, deviceId string) (*model.RefreshToken, error)
	Create(db *gorm.DB, user *model.RefreshToken) error
	Update(tx *gorm.DB, user *model.RefreshToken) error
	Delete(db *gorm.DB, id string) error
}

type RefreshTokenRepository struct{}

func NewRefreshTokenRepository() IRefreshTokenRepository {
	return &RefreshTokenRepository{}
}

func (r *RefreshTokenRepository) GetByID(db *gorm.DB, ID string) (*model.RefreshToken, error) {
	var refreshToken model.RefreshToken
	err := db.First(&refreshToken, "id = ?", ID).Error
	return &refreshToken, err
}

func (r *RefreshTokenRepository) GetByUIDAndDeviceID(db *gorm.DB, userId, deviceId string) (*model.RefreshToken, error) {
	var refreshToken model.RefreshToken
	err := db.First(&refreshToken, "user_id = ? AND device_id = ?", userId, deviceId).Error
	return &refreshToken, err
}

func (r *RefreshTokenRepository) Create(db *gorm.DB, refreshToken *model.RefreshToken) error {
	if err := db.Create(refreshToken).Error; err != nil {
		return err
	}
	return nil
}

func (r *RefreshTokenRepository) Update(tx *gorm.DB, refreshToken *model.RefreshToken) error {
	return tx.Save(refreshToken).Error
}

func (r *RefreshTokenRepository) Delete(db *gorm.DB, id string) error {
	return db.Delete(&model.RefreshToken{}, "id = ?", id).Error
}
