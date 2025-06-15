package scope

import "github.com/ebobola-dev/socially-app-go-server/internal/repository"

type IRepoistoriesScope interface {
	GetOtpRepository() repository.IOtpRepository
}

type RepositoriesScope struct {
	otp repository.IOtpRepository
}

func NewRepositoriesScope() IRepoistoriesScope {
	return &RepositoriesScope{
		otp: repository.NewOtpRepository(),
	}
}

func (s *RepositoriesScope) GetOtpRepository() repository.IOtpRepository {
	return s.otp
}
