package repository

import (
	"github.com/ebobola-dev/socially-app-go-server/internal/model"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type IPrivilegeRepository interface {
	GetByName(db *gorm.DB, name string) (*model.Privilege, error)
	GetByID(db *gorm.DB, ID uuid.UUID) (*model.Privilege, error)
	Create(db *gorm.DB, privilege *model.Privilege) error
	Update(tx *gorm.DB, privilege *model.Privilege) error
	Delete(db *gorm.DB, id uuid.UUID) error
	GetUsers(db *gorm.DB, privName string) ([]model.User, error)
}

type PrivilegeRepository struct{}

func NewPrivilegeRepository() IPrivilegeRepository {
	return &PrivilegeRepository{}
}

func (r *PrivilegeRepository) GetByName(db *gorm.DB, name string) (*model.Privilege, error) {
	var privilege model.Privilege
	err := db.Where("name = ?", name).First(&privilege).Error
	return &privilege, err
}

func (r *PrivilegeRepository) GetByID(db *gorm.DB, id uuid.UUID) (*model.Privilege, error) {
	var privilege model.Privilege
	err := db.Where("id = ?", id).First(&privilege).Error
	return &privilege, err
}

func (r *PrivilegeRepository) Create(db *gorm.DB, privilege *model.Privilege) error {
	if err := db.Create(privilege).Error; err != nil {
		return err
	}
	return nil
}

func (r *PrivilegeRepository) Update(tx *gorm.DB, privilege *model.Privilege) error {
	return tx.Save(privilege).Error
}

func (r *PrivilegeRepository) Delete(db *gorm.DB, id uuid.UUID) error {
	return db.Delete(&model.Privilege{}, "id = ?", id).Error
}

func (r *PrivilegeRepository) GetUsers(tx *gorm.DB, privName string) ([]model.User, error) {
	var users []model.User
	err := tx.
		Joins("JOIN user_privileges ON user_privileges.user_id = users.id").
		Joins("JOIN privileges ON privileges.id = user_privileges.privilege_id").
		Where("privileges.name = ?", privName).
		Where("users.deleted_at IS NULL").
		Find(&users).
		Error
	if err != nil {
		return nil, err
	}
	return users, nil
}
