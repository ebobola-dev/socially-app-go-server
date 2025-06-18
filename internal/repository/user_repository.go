package repository

import (
	"time"

	"github.com/ebobola-dev/socially-app-go-server/internal/model"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type IUserRepository interface {
	GetByID(db *gorm.DB, id uuid.UUID, includeDeleted bool) (*model.User, error)
	GetByUsername(db *gorm.DB, username string) (*model.User, error)
	GetByEmail(db *gorm.DB, email string) (*model.User, error)
	Create(db *gorm.DB, user *model.User) error
	CreateWithPrivilege(tx *gorm.DB, user *model.User, privName string) error
	Update(tx *gorm.DB, user *model.User) error
	HardDelete(db *gorm.DB, id uuid.UUID) error
	ExistsByEmail(tx *gorm.DB, email string) (bool, error)
	ExistsByUsername(tx *gorm.DB, username string) (bool, error)
	AddPrivilege(tx *gorm.DB, userID uuid.UUID, privID uuid.UUID) error
	HasAnyPrivileges(tx *gorm.DB, userID uuid.UUID, privNames ...string) (bool, error)
	HasAllPrivileges(tx *gorm.DB, userID uuid.UUID, privNames ...string) (bool, error)
	RemovePrivilege(tx *gorm.DB, userId uuid.UUID, privName string) error
	SoftDelete(tx *gorm.DB, id uuid.UUID) error
}

type UserRepository struct{}

func NewUserRepository() IUserRepository {
	return &UserRepository{}
}

func (r *UserRepository) GetByID(db *gorm.DB, id uuid.UUID, includeDeleted bool) (*model.User, error) {
	var user model.User
	query := "id = ?"
	if !includeDeleted {
		query += " AND deleted_at IS NULL"
	}
	err := db.Preload("Privileges").First(&user, query, id).Error
	return &user, err
}

func (r *UserRepository) GetByUsername(db *gorm.DB, username string) (*model.User, error) {
	var user model.User
	err := db.Preload("Privileges").First(&user, "username = ? AND deleted_at IS NULL", username).Error
	return &user, err
}

func (r *UserRepository) GetByEmail(db *gorm.DB, email string) (*model.User, error) {
	var user model.User
	err := db.Preload("Privileges").First(&user, "email = ? AND deleted_at IS NULL", email).Error
	return &user, err
}

func (r *UserRepository) Create(tx *gorm.DB, user *model.User) error {
	if err := tx.Create(user).Error; err != nil {
		return err
	}
	return tx.Preload("Privileges").First(user, "id = ? AND deleted_at IS NULL", user.ID).Error
}

func (r *UserRepository) CreateWithPrivilege(tx *gorm.DB, user *model.User, privName string) error {
	if err := tx.Create(user).Error; err != nil {
		return err
	}

	var privilege model.Privilege
	if err := tx.Where("name = ?", privName).First(&privilege).Error; err != nil {
		return err
	}

	if err := tx.Model(user).Association("Privileges").Append(&privilege); err != nil {
		return err
	}

	return tx.Preload("Privileges").First(user, "id = ? AND deleted_at IS NULL", user.ID).Error
}

func (r *UserRepository) Update(tx *gorm.DB, user *model.User) error {
	return tx.Save(user).Error
}

func (r *UserRepository) HardDelete(db *gorm.DB, id uuid.UUID) error {
	result := db.Delete(&model.User{}, "id = ?", id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func (r *UserRepository) ExistsByEmail(tx *gorm.DB, email string) (bool, error) {
	var exists bool
	err := tx.
		Raw("SELECT EXISTS(SELECT 1 FROM users WHERE email = ? AND deleted_at IS NULL)", email).
		Scan(&exists).Error
	return exists, err
}

func (r *UserRepository) ExistsByUsername(tx *gorm.DB, username string) (bool, error) {
	var exists bool
	err := tx.
		Raw("SELECT EXISTS(SELECT 1 FROM users WHERE username = ? AND deleted_at IS NULL)", username).
		Scan(&exists).Error
	return exists, err
}

func (r *UserRepository) AddPrivilege(tx *gorm.DB, userID uuid.UUID, privID uuid.UUID) error {
	user := model.User{ID: userID}
	privilege := model.Privilege{ID: privID}
	err := tx.
		Model(&user).
		Association("Privileges").
		Append(&privilege)
	return err
}

func (r *UserRepository) HasAnyPrivileges(tx *gorm.DB, userID uuid.UUID, privNames ...string) (bool, error) {
	if len(privNames) == 0 {
		return true, nil
	}

	var count int64
	err := tx.
		Model(&model.Privilege{}).
		Joins("JOIN user_privileges ON user_privileges.privilege_id = privileges.id").
		Where("user_privileges.user_id = ? AND privileges.name IN ?", userID, privNames).
		Count(&count).
		Error

	return count > 0, err
}

func (r *UserRepository) HasAllPrivileges(tx *gorm.DB, userID uuid.UUID, privNames ...string) (bool, error) {
	if len(privNames) == 0 {
		return true, nil
	}

	var matchedCount int64
	err := tx.
		Model(&model.Privilege{}).
		Joins("JOIN user_privileges ON user_privileges.privilege_id = privileges.id").
		Where("user_privileges.user_id = ? AND privileges.name IN ?", userID, privNames).
		Count(&matchedCount).
		Error

	if err != nil {
		return false, err
	}

	return matchedCount == int64(len(privNames)), nil
}

func (r *UserRepository) RemovePrivilege(tx *gorm.DB, userId uuid.UUID, privName string) error {
	user := model.User{ID: userId}
	privilege := model.Privilege{Name: privName}
	err := tx.
		Model(&user).
		Association("Privileges").Delete(privilege)
	return err
}

func (r *UserRepository) SoftDelete(tx *gorm.DB, id uuid.UUID) error {
	var user model.User
	if err := tx.First(&user, "id = ? AND deleted_at IS NULL", id).Error; err != nil {
		return err
	}
	user.Email = ""
	user.Username = ""
	user.Fullname = nil
	user.AboutMe = nil
	user.Gender = nil
	user.DateOfBirth = time.Date(100, 1, 1, 0, 0, 0, 0, time.UTC)
	user.AvatarType = nil
	user.AvatarID = nil
	user.LastSeen = nil
	user.Privileges = []model.Privilege{}

	now := time.Now()
	user.DeletedAt = &now

	if err := tx.Model(&user).Association("Privileges").Clear(); err != nil {
		return err
	}

	return tx.Save(&user).Error
}
