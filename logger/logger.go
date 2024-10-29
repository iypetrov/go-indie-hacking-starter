package logger

type Logger interface {
	Debug(msg string, args ...interface{})
	Info(msg string, args ...interface{})
	Warn(msg string, args ...interface{})
	Error(msg string, args ...interface{})
}

var mylog Logger

func Init() {
	mylog = New()
}

func Get() Logger {
	return mylog
}
