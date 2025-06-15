package repository

import (
	"github.com/ebobola-dev/socially-app-go-server/internal/model"
	"gorm.io/gorm"
)

type IPrivilegeRepository interface {
	GetByName(db *gorm.DB, name string) (*model.Privilege, error)
	GetByID(db *gorm.DB, ID string) (*model.Privilege, error)
	Create(db *gorm.DB, privilege *model.Privilege) error
	Update(tx *gorm.DB, privilege *model.Privilege) error
	Delete(db *gorm.DB, id string) error
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

func (r *PrivilegeRepository) GetByID(db *gorm.DB, ID string) (*model.Privilege, error) {
	var privilege model.Privilege
	err := db.Where("id = ?", ID).First(&privilege).Error
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

func (r *PrivilegeRepository) Delete(db *gorm.DB, id string) error {
	return db.Delete(&model.Privilege{}, "id = ?", id).Error
}
