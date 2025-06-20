package scope

import (
	"context"

	"github.com/ebobola-dev/socially-app-go-server/internal/config"
	"github.com/ebobola-dev/socially-app-go-server/internal/database"
	"github.com/ebobola-dev/socially-app-go-server/internal/repository"
	"github.com/ebobola-dev/socially-app-go-server/internal/service/email"
	"github.com/ebobola-dev/socially-app-go-server/internal/service/hash"
	jwt_service "github.com/ebobola-dev/socially-app-go-server/internal/service/jwt"
	minio_service "github.com/ebobola-dev/socially-app-go-server/internal/service/minio"
	"github.com/ebobola-dev/socially-app-go-server/internal/util/logger"
	"github.com/ebobola-dev/socially-app-go-server/internal/validation"
	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
)

type AppScope struct {
	Cfg                    *config.Config
	Log                    logger.ILogger
	Db                     *gorm.DB
	Validate               *validator.Validate
	OtpRepository          repository.IOtpRepository
	PrivilegeRepository    repository.IPrivilegeRepository
	UserRepository         repository.IUserRepository
	RefreshTokenRepository repository.IRefreshTokenRepository
	EmailService           email.IEmailService
	JwtService             jwt_service.IJwtService
	HashService            hash.IHashService
	MinioService           minio_service.IMinioService
}

func ConfigureAppScope(ctx context.Context) *AppScope {
	cfg := config.Initialize()
	return &AppScope{
		Cfg:                    cfg,
		Log:                    logger.Create(cfg),
		Db:                     database.Connect(cfg.Database),
		Validate:               validation.NewValidator(),
		OtpRepository:          repository.NewOtpRepository(),
		PrivilegeRepository:    repository.NewPrivilegeRepository(),
		UserRepository:         repository.NewUserRepository(),
		RefreshTokenRepository: repository.NewRefreshTokenRepository(),
		EmailService:           email.NewEmailService(cfg.SMTP),
		JwtService:             jwt_service.NewJwtService(cfg.JWT),
		HashService:            hash.NewHashService(),
		MinioService:           minio_service.NewMinioService(ctx, cfg.Minio),
	}
}
