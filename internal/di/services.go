package scope

import (
	"github.com/ebobola-dev/socially-app-go-server/internal/config"
	"github.com/ebobola-dev/socially-app-go-server/internal/service/email"
	hash_s "github.com/ebobola-dev/socially-app-go-server/internal/service/hash"
	jwt_s "github.com/ebobola-dev/socially-app-go-server/internal/service/jwt"
)

type IServicesScope interface {
	GetEmailService() email.IEmailService
	GetJwtService() jwt_s.IJwtService
	GetHashService() hash_s.IHashService
}

type ServicesScope struct {
	email email.IEmailService
	jwt   jwt_s.IJwtService
	hash  hash_s.IHashService
}

func NewServicesScope(smtpCfg *config.SMTPConfig, jwtCfg *config.JWTConfig) IServicesScope {
	return &ServicesScope{
		email: email.NewEmailService(smtpCfg),
		jwt:   jwt_s.NewJwtService(jwtCfg),
		hash:  hash_s.NewHashService(),
	}
}

func (s *ServicesScope) GetEmailService() email.IEmailService {
	return s.email
}

func (s *ServicesScope) GetJwtService() jwt_s.IJwtService {
	return s.jwt
}

func (s *ServicesScope) GetHashService() hash_s.IHashService {
	return s.hash
}
