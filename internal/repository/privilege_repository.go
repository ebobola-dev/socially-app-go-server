package repository

import (
	"strings"

	privilege_error "github.com/ebobola-dev/socially-app-go-server/internal/errors/privilege"
	"github.com/ebobola-dev/socially-app-go-server/internal/model"
	pagination "github.com/ebobola-dev/socially-app-go-server/internal/util/pagintation"
	"github.com/google/uuid"
	"github.com/samber/lo"
	"gorm.io/gorm"
)

type IPrivilegeRepository interface {
	GetByName(db *gorm.DB, name string, options GetPrivilegeOptions) (*model.Privilege, error)
	GetByID(db *gorm.DB, ID uuid.UUID, options GetPrivilegeOptions) (*model.Privilege, error)
	Create(db *gorm.DB, privilege *model.Privilege) error
	Update(tx *gorm.DB, privilege *model.Privilege) error
	Delete(db *gorm.DB, id uuid.UUID) error
	GetUsers(db *gorm.DB, pagination pagination.Pagination, privName string) ([]model.User, error)
	GetAll(tx *gorm.DB, options GetPrivilegesListOptions) ([]model.Privilege, error)
}

type privilegeRepository struct{}

func NewPrivilegeRepository() IPrivilegeRepository {
	return &privilegeRepository{}
}

func (r *privilegeRepository) GetByName(tx *gorm.DB, name string, options GetPrivilegeOptions) (*model.Privilege, error) {
	var privilege model.Privilege
	if err := tx.Where("name = ?", name).First(&privilege).Error; err != nil {
		return nil, err
	}
	if options.CountUsers {
		var count int64
		if err := tx.
			Table("user_privileges").
			Where("privilege_id = ?", privilege.ID).
			Count(&count).Error; err != nil {
			return nil, err
		}
		privilege.UsersCount = int(count)
	}
	return &privilege, nil
}

func (r *privilegeRepository) GetByID(tx *gorm.DB, id uuid.UUID, options GetPrivilegeOptions) (*model.Privilege, error) {
	var privilege model.Privilege
	if err := tx.Where("id = ?", id).First(&privilege).Error; err != nil {
		return nil, err
	}
	if options.CountUsers {
		var count int64
		if err := tx.
			Table("user_privileges").
			Where("privilege_id = ?", id).
			Count(&count).Error; err != nil {
			return nil, err
		}
		privilege.UsersCount = int(count)
	}
	return &privilege, nil
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

func (r *privilegeRepository) GetUsers(tx *gorm.DB, pagination pagination.Pagination, privName string) ([]model.User, error) {
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

func (r *privilegeRepository) GetAll(tx *gorm.DB, options GetPrivilegesListOptions) ([]model.Privilege, error) {
	var privileges []model.Privilege
	query := tx.Model(&model.Privilege{})
	if options.FilterUserId != uuid.Nil {
		query = query.
			Joins("JOIN user_privileges ON user_privileges.privilege_id = privileges.id").
			Where("user_privileges.user_id = ?", options.FilterUserId)
	}
	if err := query.
		Order("order_index DESC").
		Offset(options.Pagination.Offset).
		Limit(options.Pagination.Limit).
		Find(&privileges).Error; err != nil {
		return nil, err
	}
	if options.CountUsers && len(privileges) > 0 {
		type CountResult struct {
			PrivilegeID uuid.UUID
			Count       int
		}
		var results []CountResult
		if err := tx.
			Table("user_privileges").
			Select("privilege_id, COUNT(*) as count").
			Where("privilege_id IN ?", lo.Map(privileges, func(p model.Privilege, _ int) uuid.UUID {
				return p.ID
			})).
			Group("privilege_id").
			Find(&results).Error; err != nil {
			return nil, err
		}
		countMap := make(map[uuid.UUID]int)
		for _, r := range results {
			countMap[r.PrivilegeID] = r.Count
		}
		for i := range privileges {
			privileges[i].UsersCount = countMap[privileges[i].ID]
		}
	}

	return privileges, nil
}

type GetPrivilegeOptions struct {
	CountUsers bool
}

type GetPrivilegesListOptions struct {
	Pagination   pagination.Pagination
	CountUsers   bool
	FilterUserId uuid.UUID
}
