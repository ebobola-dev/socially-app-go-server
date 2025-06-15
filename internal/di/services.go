package scope

import (
	"github.com/ebobola-dev/socially-app-go-server/internal/config"
	"github.com/ebobola-dev/socially-app-go-server/internal/service/email"
	jwt_service "github.com/ebobola-dev/socially-app-go-server/internal/service/jwt"
)

type IServicesScope interface {
	GetEmailService() email.IEmailService
	GetJwtService() jwt_service.IJwtService
}

type ServicesScope struct {
	email email.IEmailService
	jwt   jwt_service.IJwtService
}

func NewServicesScope(smtpCfg *config.SMTPConfig, jwtCfg *config.JWTConfig) IServicesScope {
	return &ServicesScope{
		email: email.NewEmailService(smtpCfg),
		jwt:   jwt_service.NewJwtService(jwtCfg),
	}
}

func (s *ServicesScope) GetEmailService() email.IEmailService {
	return s.email
}

func (s *ServicesScope) GetJwtService() jwt_service.IJwtService {
	return s.jwt
}
