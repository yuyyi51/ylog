package ylog

type LogLevel int

const (
	LogLevelInvalid LogLevel = iota
	LogLevelTrace
	LogLevelDebug
	LogLevelInfo
	LogLevelNotice
	LogLevelWarn
	LogLevelError
	LogLevelPanic
	LogLevelFatal
	LogLevelNone
)
