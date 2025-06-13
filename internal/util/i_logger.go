package logger

type ILogger interface {
	Debug(format string, args ...any)
	Info(format string, args ...any)
	Warn(format string, args ...any)
	Error(err error)
	Fatal(err error)
}
