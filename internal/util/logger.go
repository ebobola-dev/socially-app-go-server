package logger

import (
	"log"
	"os"
	"time"

	"github.com/fatih/color"
	"github.com/rotisserie/eris"
)

type MyLogger struct {
	infoLogger  *log.Logger
	errorLogger *log.Logger

	infoColor  *color.Color
	warnColor  *color.Color
	errorColor *color.Color
	fatalColor *color.Color
}

func Create() ILogger {
	return &MyLogger{
		infoLogger:  log.New(os.Stdout, "", 0),
		errorLogger: log.New(os.Stderr, "", 0),
		infoColor:   color.New(color.FgCyan),
		warnColor:   color.New(color.FgYellow),
		errorColor:  color.New(color.FgRed),
		fatalColor:  color.New(color.FgHiRed, color.Bold),
	}
}

func timestamp() string {
	return time.Now().UTC().Format("02.01 15:04:05")
}

func (l *MyLogger) Info(msg string) {
	l.infoLogger.Println(l.infoColor.Sprintf("%s [INFO] %s", timestamp(), msg))
}

func (l *MyLogger) Warn(msg string) {
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
