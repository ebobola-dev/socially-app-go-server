package scope

import (
	"github.com/ebobola-dev/socially-app-go-server/internal/config"
	"github.com/ebobola-dev/socially-app-go-server/internal/util/logger"
	"gorm.io/gorm"
)

type IAppScope interface {
	GetConfig() *config.Config
	GetLogger() logger.ILogger
	GetDB() *gorm.DB
}

type AppScope struct {
	cfg *config.Config
	log logger.ILogger
	db  *gorm.DB
}

func NewAppScope(cfg *config.Config, log logger.ILogger, db *gorm.DB) IAppScope {
	return &AppScope{
		cfg: cfg,
		log: log,
		db:  db,
	}
}

func (s *AppScope) GetConfig() *config.Config {
	return s.cfg
}

func (s *AppScope) GetLogger() logger.ILogger {
	return s.log
}

func (s *AppScope) GetDB() *gorm.DB {
	return s.db
}
