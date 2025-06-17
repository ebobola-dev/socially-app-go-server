package logger

type ILogger interface {
	Debug(format string, args ...any)
	Info(format string, args ...any)
	Warning(format string, args ...any)
	Error(format string, args ...any)
	Exception(err error)
	Fatal(err error)
}
