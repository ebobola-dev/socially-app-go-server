package repository

import (
	"time"

	"github.com/ebobola-dev/socially-app-go-server/internal/model"
	pagination "github.com/ebobola-dev/socially-app-go-server/internal/util/pagintation"
	"github.com/google/uuid"
	"github.com/samber/lo"
	"gorm.io/gorm"
)

type IUserRepository interface {
	GetByID(tx *gorm.DB, id uuid.UUID, options GetUserOptions) (*model.User, error)
	GetByUsername(tx *gorm.DB, username string) (*model.User, error)
	GetByEmail(tx *gorm.DB, email string) (*model.User, error)
	Create(tx *gorm.DB, user *model.User) error
	CreateWithPrivilege(tx *gorm.DB, user *model.User, privName string) error
	Update(tx *gorm.DB, user *model.User) error
	HardDelete(tx *gorm.DB, id uuid.UUID) error
	ExistsByEmail(tx *gorm.DB, email string) (bool, error)
	ExistsByUsername(tx *gorm.DB, username string) (bool, error)
	AddPrivilege(tx *gorm.DB, userID uuid.UUID, privID uuid.UUID) error
	HasAnyPrivileges(tx *gorm.DB, userID uuid.UUID, privNames ...string) (bool, error)
	HasAllPrivileges(tx *gorm.DB, userID uuid.UUID, privNames ...string) (bool, error)
	RemovePrivilege(tx *gorm.DB, userId uuid.UUID, privName string) error
	SoftDelete(tx *gorm.DB, id uuid.UUID) error
	Search(tx *gorm.DB, options SearchUsersOptions) ([]model.User, error)
	GetPrivileges(tx *gorm.DB, opts GetUserPrivilegesOptions) ([]model.UserPrivilege, error)
}

type userRepository struct{}

func NewUserRepository() IUserRepository {
	return &userRepository{}
}

func (r *userRepository) GetByID(tx *gorm.DB, id uuid.UUID, options GetUserOptions) (*model.User, error) {
	var user model.User
	if !options.IncludeDeleted {
		tx = tx.Where("deleted_at IS NULL")
	}
	err := tx.Preload("UserPrivileges.Privilege").First(&user, "id = ?", id).Error
	return &user, err
}

func (r *userRepository) GetByUsername(tx *gorm.DB, username string) (*model.User, error) {
	var user model.User
	err := tx.Preload("UserPrivileges.Privilege").First(&user, "username = ? and deleted_at IS NULL", username).Error
	return &user, err
}

func (r *userRepository) GetByEmail(tx *gorm.DB, email string) (*model.User, error) {
	var user model.User
	err := tx.Preload("UserPrivileges.Privilege").First(&user, "email = ? and deleted_at IS NULL", email).Error
	return &user, err
}

func (r *userRepository) Create(tx *gorm.DB, user *model.User) error {
	if err := tx.Create(user).Error; err != nil {
		return err
	}
	return tx.Preload("UserPrivileges.Privilege").First(user, "id = ? AND deleted_at IS NULL", user.ID).Error
}

func (r *userRepository) CreateWithPrivilege(tx *gorm.DB, user *model.User, privName string) error {
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

	return tx.Preload("UserPrivileges.Privilege").First(user, "id = ? AND deleted_at IS NULL", user.ID).Error
}

func (r *userRepository) Update(tx *gorm.DB, user *model.User) error {
	return tx.Save(user).Error
}

func (r *userRepository) HardDelete(tx *gorm.DB, id uuid.UUID) error {
	result := tx.Delete(&model.User{}, "id = ?", id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func (r *userRepository) ExistsByEmail(tx *gorm.DB, email string) (bool, error) {
	var exists bool
	err := tx.
		Raw("SELECT EXISTS(SELECT 1 FROM users WHERE email = ? AND deleted_at IS NULL)", email).
		Scan(&exists).Error
	return exists, err
}

func (r *userRepository) ExistsByUsername(tx *gorm.DB, username string) (bool, error) {
	var exists bool
	err := tx.
		Raw("SELECT EXISTS(SELECT 1 FROM users WHERE username = ? AND deleted_at IS NULL)", username).
		Scan(&exists).Error
	return exists, err
}

func (r *userRepository) AddPrivilege(tx *gorm.DB, userID uuid.UUID, privID uuid.UUID) error {
	user := model.User{ID: userID}
	privilege := model.Privilege{ID: privID}
	err := tx.
		Model(&user).
		Association("Privileges").
		Append(&privilege)
	return err
}

func (r *userRepository) HasAnyPrivileges(tx *gorm.DB, userID uuid.UUID, privNames ...string) (bool, error) {
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

func (r *userRepository) HasAllPrivileges(tx *gorm.DB, userID uuid.UUID, privNames ...string) (bool, error) {
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

func (r *userRepository) RemovePrivilege(tx *gorm.DB, userId uuid.UUID, privName string) error {
	user := model.User{ID: userId}
	privilege := model.Privilege{Name: privName}
	err := tx.
		Model(&user).
		Association("Privileges").Delete(privilege)
	return err
}

func (r *userRepository) SoftDelete(tx *gorm.DB, id uuid.UUID) error {
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

func (r *userRepository) Search(
	tx *gorm.DB,
	options SearchUsersOptions,
) ([]model.User, error) {
	var users []model.User
	searchPattern := "%" + options.Pattern + "%"

	tx = tx.Model(&model.User{})

	if !options.IncludeDeleted {
		tx = tx.Where("deleted_at IS NULL")
	}
	if options.IgnoreId != uuid.Nil {
		tx = tx.Where("id <> ?", options.IgnoreId)
	}

	if err := tx.
		Where("(username LIKE ? OR fullname LIKE ?)", searchPattern, searchPattern).
		Order("created_at DESC").
		Offset(options.Pagination.Offset).
		Limit(options.Pagination.Limit).
		Find(&users).Error; err != nil {
		return nil, err
	}

	return users, nil
}

func (r *userRepository) GetPrivileges(tx *gorm.DB, opts GetUserPrivilegesOptions) ([]model.UserPrivilege, error) {
	var userPrivileges []model.UserPrivilege

	query := tx.Model(&model.UserPrivilege{}).
		Preload("Privilege").
		Where("user_id = ?", opts.UserID).
		Order("privileges.order_index DESC").
		Offset(opts.Pagination.Offset).
		Limit(opts.Pagination.Limit).
		Joins("JOIN privileges ON privileges.id = user_privileges.privilege_id")

	if err := query.Find(&userPrivileges).Error; err != nil {
		return nil, err
	}

	if opts.CountUsers && len(userPrivileges) > 0 {
		type CountResult struct {
			PrivilegeID uuid.UUID
			Count       int
		}
		var results []CountResult
		if err := tx.
			Table("user_privileges").
			Select("privilege_id, COUNT(*) as count").
			Where("privilege_id IN ?", lo.Map(userPrivileges, func(up model.UserPrivilege, _ int) uuid.UUID {
				return up.PrivilegeID
			})).
			Group("privilege_id").
			Find(&results).Error; err != nil {
			return nil, err
		}

		countMap := make(map[uuid.UUID]int)
		for _, r := range results {
			countMap[r.PrivilegeID] = r.Count
		}
		for i := range userPrivileges {
			userPrivileges[i].Privilege.UsersCount = countMap[userPrivileges[i].PrivilegeID]
		}
	}

	return userPrivileges, nil
}

type GetUserOptions struct {
	IncludeDeleted bool
}

type SearchUsersOptions struct {
	Pagination     pagination.Pagination
	Pattern        string
	IncludeDeleted bool
	IgnoreId       uuid.UUID
}

type GetUserPrivilegesOptions struct {
	Pagination pagination.Pagination
	UserID     uuid.UUID
	CountUsers bool
}
