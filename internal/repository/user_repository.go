package repository

import (
	"errors"
	"time"

	"github.com/ebobola-dev/socially-app-go-server/internal/model"
	pagination "github.com/ebobola-dev/socially-app-go-server/internal/util/pagination"
	"github.com/go-sql-driver/mysql"
	"github.com/google/uuid"
	"github.com/samber/lo"
	"gorm.io/gorm"
)

var (
	ErrSubActionYourself  = errors.New("can't follow yourself")
	ErrAlreadyFollowing   = errors.New("already following")
	ErrNotFollowingAnyway = errors.New("not following anyway")
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
	Follow(tx *gorm.DB, subscriberId, targetId uuid.UUID) (int64, error)
	Unfollow(tx *gorm.DB, followerID, targetID uuid.UUID) (int64, error)
	GetFollowers(tx *gorm.DB, options GetSubscriptionsOptions) ([]model.UserSubscription, error)
	GetFollowing(tx *gorm.DB, options GetSubscriptionsOptions) ([]model.UserSubscription, error)
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
	if err := tx.Preload("UserPrivileges.Privilege").First(&user, "id = ?", id).Error; err != nil {
		return nil, err
	}
	if options.CountSubscriptions {
		var count int64
		if err := tx.
			Model(&model.UserSubscription{}).
			Where("follower_id = ?", id).
			Count(&count).Error; err != nil {
			return nil, err
		}
		user.FollowingCount = count
		if err := tx.
			Model(&model.UserSubscription{}).
			Where("target_id  = ?", id).
			Count(&count).Error; err != nil {
			return nil, err
		}
		user.FollowersCount = count
	}
	return &user, nil
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

func (r *userRepository) Follow(tx *gorm.DB, subscriberId, targetId uuid.UUID) (int64, error) {
	if subscriberId == targetId {
		return 0, ErrSubActionYourself
	}
	var target model.User
	if err := tx.Select("id").First(&target, "id = ? AND deleted_at IS NULL", targetId).Error; err != nil {
		return 0, err
	}
	subscribtion := model.UserSubscription{
		FollowerID: subscriberId,
		TargetID:   targetId,
	}
	if err := tx.Create(&subscribtion).Error; err != nil {
		var mysqlErr *mysql.MySQLError
		if errors.As(err, &mysqlErr) && mysqlErr.Number == 1062 {
			return 0, ErrAlreadyFollowing
		}
		return 0, err
	}
	var newFollowersCount int64
	if err := tx.Model(&model.UserSubscription{}).
		Where("target_id = ?", targetId).
		Count(&newFollowersCount).Error; err != nil {
		return 0, err
	}
	return newFollowersCount, nil
}

func (r *userRepository) Unfollow(tx *gorm.DB, followerID, targetId uuid.UUID) (int64, error) {
	if followerID == targetId {
		return 0, ErrSubActionYourself
	}
	var target model.User
	if err := tx.Select("id").First(&target, "id = ? AND deleted_at IS NULL", targetId).Error; err != nil {
		return 0, err
	}
	res := tx.Delete(&model.UserSubscription{}, "follower_id = ? AND target_id = ?", followerID, targetId)
	if res.Error != nil {
		return 0, res.Error
	}
	if res.RowsAffected == 0 {
		return 0, ErrNotFollowingAnyway
	}
	var newFollowersCount int64
	if err := tx.Model(&model.UserSubscription{}).
		Where("target_id = ?", targetId).
		Count(&newFollowersCount).Error; err != nil {
		return 0, err
	}
	return newFollowersCount, nil
}

func (r *userRepository) GetFollowers(tx *gorm.DB, options GetSubscriptionsOptions) ([]model.UserSubscription, error) {
	var followers []model.UserSubscription
	var target model.User
	if err := tx.Select("id").First(&target, "id = ? AND deleted_at IS NULL", options.TargetUID).Error; err != nil {
		return nil, err
	}
	query := tx.
		Where("target_id = ?", options.TargetUID).
		Preload("Follower").
		Order("created_at DESC").
		Offset(options.Pagination.Offset).
		Limit(options.Pagination.Limit)

	if err := query.Find(&followers).Error; err != nil {
		return nil, err
	}

	return followers, nil
}

func (r *userRepository) GetFollowing(tx *gorm.DB, options GetSubscriptionsOptions) ([]model.UserSubscription, error) {
	var following []model.UserSubscription
	var target model.User
	if err := tx.Select("id").First(&target, "id = ? AND deleted_at IS NULL", options.TargetUID).Error; err != nil {
		return nil, err
	}
	query := tx.
		Where("follower_id = ?", options.TargetUID).
		Preload("Target").
		Order("created_at DESC").
		Offset(options.Pagination.Offset).
		Limit(options.Pagination.Limit)

	if err := query.Find(&following).Error; err != nil {
		return nil, err
	}

	return following, nil
}

type GetUserOptions struct {
	IncludeDeleted     bool
	CountSubscriptions bool
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

type GetSubscriptionsOptions struct {
	TargetUID  uuid.UUID
	Pagination pagination.Pagination
}
