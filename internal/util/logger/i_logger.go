package logger

import "github.com/ebobola-dev/socially-app-go-server/internal/config"

type ILogger interface {
	Debug(format string, args ...any)
	Info(format string, args ...any)
	Warn(format string, args ...any)
	Error(err error)
	Fatal(err error)
	PrintConfig(cfg *config.Config)
}
