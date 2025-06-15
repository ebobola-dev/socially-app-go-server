package scope

import "github.com/ebobola-dev/socially-app-go-server/internal/repository"

type IRepositoriesScope interface {
	GetOtpRepository() repository.IOtpRepository
	GetPrivilegeRepository() repository.IPrivilegeRepository
	GetUserRepository() repository.IUserRepository
}

type RepositoriesScope struct {
	otp       repository.IOtpRepository
	privilege repository.IPrivilegeRepository
	user      repository.IUserRepository
}

func NewRepositoriesScope() IRepositoriesScope {
	return &RepositoriesScope{
		otp:       repository.NewOtpRepository(),
		privilege: repository.NewPrivilegeRepository(),
		user:      repository.NewUserRepository(),
	}
}

func (s *RepositoriesScope) GetOtpRepository() repository.IOtpRepository {
	return s.otp
}

func (s *RepositoriesScope) GetPrivilegeRepository() repository.IPrivilegeRepository {
	return s.privilege
}

func (s *RepositoriesScope) GetUserRepository() repository.IUserRepository {
	return s.user
}
