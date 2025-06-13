package logger

import (
	"fmt"
	"log"
	"os"
	"socially-app/internal/config"
	"time"

	"github.com/fatih/color"
	"github.com/rotisserie/eris"
)

type MyLogger struct {
	cfg         *config.Config
	infoLogger  *log.Logger
	errorLogger *log.Logger

	debugColor *color.Color
	infoColor  *color.Color
	warnColor  *color.Color
	errorColor *color.Color
	fatalColor *color.Color
}

func Create(cfg *config.Config) ILogger {
	return &MyLogger{
		cfg:         cfg,
		infoLogger:  log.New(os.Stdout, "", 0),
		errorLogger: log.New(os.Stderr, "", 0),

		debugColor: color.New(color.FgHiMagenta),
		infoColor:  color.New(color.FgCyan),
		warnColor:  color.New(color.FgYellow),
		errorColor: color.New(color.FgRed),
		fatalColor: color.New(color.FgHiRed, color.Bold),
	}
}

func timestamp() string {
	return time.Now().UTC().Format("02.01 15:04:05")
}

func (l *MyLogger) Debug(format string, args ...any) {
	if l.cfg.BuildType == config.Development {
		msg := fmt.Sprintf(format, args...)
		l.infoLogger.Println(l.debugColor.Sprintf("%s [DEBUG] %s", timestamp(), msg))
	}
}

func (l *MyLogger) Info(format string, args ...any) {
	msg := fmt.Sprintf(format, args...)
	l.infoLogger.Println(l.infoColor.Sprintf("%s [INFO] %s", timestamp(), msg))
}

func (l *MyLogger) Warn(format string, args ...any) {
	msg := fmt.Sprintf(format, args...)
	l.infoLogger.Println(l.warnColor.Sprintf("%s [WARN] %s", timestamp(), msg))
}

func (l *MyLogger) Error(err error) {
	stack := eris.ToString(err, true)
	l.errorLogger.Println(l.errorColor.Sprintf("%s [ERROR] %s", timestamp(), stack))
}

func (l *MyLogger) Fatal(err error) {
	stack := eris.ToString(err, true)
	l.errorLogger.Println(l.fatalColor.Sprintf("%s [FATAL] %s", timestamp(), stack))
	os.Exit(1)
}
