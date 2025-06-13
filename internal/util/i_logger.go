package logger

type ILogger interface {
	Info(msg string)
	Warn(msg string)
	Error(err error)
	Fatal(err error)
}
