package repository

import (
	"github.com/ebobola-dev/socially-app-go-server/internal/model"
	"gorm.io/gorm"
)

type IUserRepository interface {
	GetByID(db *gorm.DB, ID string) (*model.User, error)
	GetByUsername(db *gorm.DB, username string) (*model.User, error)
	GetByEmail(db *gorm.DB, email string) (*model.User, error)
	Create(db *gorm.DB, user *model.User) error
	Update(tx *gorm.DB, user *model.User) error
	Delete(db *gorm.DB, id string) error
	ExistsByEmail(tx *gorm.DB, email string) (bool, error)
	ExistsByUsername(tx *gorm.DB, username string) (bool, error)
}

type UserRepository struct{}

func NewUserRepository() IUserRepository {
	return &UserRepository{}
}

func (r *UserRepository) GetByID(db *gorm.DB, ID string) (*model.User, error) {
	var user model.User
	err := db.First(&user, "id = ?", ID).First(&user).Error
	return &user, err
}

func (r *UserRepository) GetByUsername(db *gorm.DB, username string) (*model.User, error) {
	var user model.User
	err := db.First(&user, "username = ?", username).Error
	return &user, err
}

func (r *UserRepository) GetByEmail(db *gorm.DB, email string) (*model.User, error) {
	var user model.User
	err := db.First(&user, "email = ?", email).Error
	return &user, err
}

func (r *UserRepository) Create(db *gorm.DB, user *model.User) error {
	if err := db.Create(user).Error; err != nil {
		return err
	}
	return nil
}

func (r *UserRepository) Update(tx *gorm.DB, user *model.User) error {
	return tx.Save(user).Error
}

func (r *UserRepository) Delete(db *gorm.DB, id string) error {
	return db.Delete(&model.User{}, "id = ?", id).Error
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
