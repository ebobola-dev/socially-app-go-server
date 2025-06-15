package scope

import "github.com/ebobola-dev/socially-app-go-server/internal/repository"

type IRepositoriesScope interface {
	GetOtpRepository() repository.IOtpRepository
}

type RepositoriesScope struct {
	otp repository.IOtpRepository
}

func NewRepositoriesScope() IRepositoriesScope {
	return &RepositoriesScope{
		otp: repository.NewOtpRepository(),
	}
}

func (s *RepositoriesScope) GetOtpRepository() repository.IOtpRepository {
	return s.otp
}
