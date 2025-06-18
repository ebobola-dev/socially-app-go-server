package database

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/ebobola-dev/socially-app-go-server/internal/config"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func Connect(cfg *config.DatabaseConfig) *gorm.DB {
	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?parseTime=true&charset=utf8mb4&loc=UTC",
		cfg.User, cfg.Password, cfg.Host, cfg.Name)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.New(log.New(os.Stdout, "\r\n", log.LstdFlags), logger.Config{
			SlowThreshold:             200 * time.Millisecond,
			LogLevel:                  logger.Warn,
			IgnoreRecordNotFoundError: true,
			Colorful:                  true,
		}),
	})
	if err != nil {
		log.Fatal("failed to connect to database:", err)
	}
	return db
}
