package repository

import (
	"strings"

	privilege_error "github.com/ebobola-dev/socially-app-go-server/internal/errors/privilege"
	"github.com/ebobola-dev/socially-app-go-server/internal/model"
	pagination "github.com/ebobola-dev/socially-app-go-server/internal/util/pagintation"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type IPrivilegeRepository interface {
	GetByName(db *gorm.DB, name string) (*model.Privilege, error)
	GetByID(db *gorm.DB, ID uuid.UUID) (*model.Privilege, error)
	Create(db *gorm.DB, privilege *model.Privilege) error
	Update(tx *gorm.DB, privilege *model.Privilege) error
	Delete(db *gorm.DB, id uuid.UUID) error
	GetUsers(db *gorm.DB, pagination pagination.Pagitation, privName string) ([]model.User, error)
	GetAll(tx *gorm.DB, pagination pagination.Pagitation) ([]model.Privilege, error)
}

type privilegeRepository struct{}

func NewPrivilegeRepository() IPrivilegeRepository {
	return &privilegeRepository{}
}

func (r *privilegeRepository) GetByName(tx *gorm.DB, name string) (*model.Privilege, error) {
	var privilege model.Privilege
	err := tx.Preload("Users").Where("name = ?", name).First(&privilege).Error
	return &privilege, err
}

func (r *privilegeRepository) GetByID(tx *gorm.DB, id uuid.UUID) (*model.Privilege, error) {
	var privilege model.Privilege
	err := tx.Preload("Users").Where("id = ?", id).First(&privilege).Error
	return &privilege, err
}

func (r *privilegeRepository) Create(tx *gorm.DB, privilege *model.Privilege) error {
	err := tx.Create(privilege).Error
	if err != nil {
		if strings.Contains(err.Error(), "Duplicate entry") {
			if strings.Contains(err.Error(), "name") {
				return privilege_error.NewDuplicateNameError(privilege.Name)
			}
			if strings.Contains(err.Error(), "order_index") {
				return privilege_error.NewDuplicateIndexError(privilege.OrderIndex)
			}
			return err
		}
		return err
	}
	return nil
}

func (r *privilegeRepository) Update(tx *gorm.DB, privilege *model.Privilege) error {
	return tx.Save(privilege).Error
}

func (r *privilegeRepository) Delete(db *gorm.DB, id uuid.UUID) error {
	result := db.Delete(&model.Privilege{}, "id = ?", id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func (r *privilegeRepository) GetUsers(tx *gorm.DB, pagination pagination.Pagitation, privName string) ([]model.User, error) {
	var users []model.User
	err := tx.
		Preload("Privileges").
		Joins("JOIN user_privileges ON user_privileges.user_id = users.id").
		Joins("JOIN privileges ON privileges.id = user_privileges.privilege_id").
		Where("privileges.name = ?", privName).
		Where("users.deleted_at IS NULL").
		Order("users.created_at DESC").
		Offset(pagination.Offset).
		Limit(pagination.Limit).
		Find(&users).
		Error
	if err != nil {
		return nil, err
	}
	return users, nil
}

func (r *privilegeRepository) GetAll(tx *gorm.DB, pagination pagination.Pagitation) ([]model.Privilege, error) {
	var privileges []model.Privilege
	err := tx.
		Preload("Users").
		Order("order_index DESC").
		Offset(pagination.Offset).
		Limit(pagination.Limit).
		Find(&privileges).Error

	return privileges, err
}
