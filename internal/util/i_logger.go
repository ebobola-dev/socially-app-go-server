package logger

import "socially-app/internal/config"

type ILogger interface {
	Debug(format string, args ...any)
	Info(format string, args ...any)
	Warn(format string, args ...any)
	Error(err error)
	Fatal(err error)
	PrintConfig(cfg *config.Config)
}
