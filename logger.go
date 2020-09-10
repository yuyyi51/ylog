package ylog

import (
	"fmt"
	"io"
	"os"
	"path"
	"runtime"
	"strings"
	"time"
)

type ILogger interface {
	HitLevel(level LogLevel) bool
	Trace(format string, args ...interface{})
	Debug(format string, args ...interface{})
	Info(format string, args ...interface{})
	Notice(format string, args ...interface{})
	Warn(format string, args ...interface{})
	Error(format string, args ...interface{})
	Panic(format string, args ...interface{})
	Fatal(format string, args ...interface{})
	AddLogChain(logger ILogger)
}

type Logger struct {
	logLevel       LogLevel
	logWriter      *logWriter
	skipStackLevel int
	chain          ILogger
}

func LogLevelToString(level LogLevel) string {
	if level <= LogLevelInvalid || level >= LogLevelNone {
		return ""
	}
	switch level {
	case LogLevelTrace:
		return "TRACE"
	case LogLevelDebug:
		return "DEBUG"
	case LogLevelInfo:
		return "INFO"
	case LogLevelNotice:
		return "NOTICE"
	case LogLevelWarn:
		return "WARN"
	case LogLevelError:
		return "ERROR"
	case LogLevelPanic:
		return "PANIC"
	case LogLevelFatal:
		return "FATAL"
	}
	return ""
}

func StringToLogLevel(level string) LogLevel {
	level = strings.ToLower(level)
	switch level {
	case "none":
		return LogLevelNone
	case "trace":
		return LogLevelTrace
	case "debug":
		return LogLevelDebug
	case "info":
		return LogLevelInfo
	case "notice":
		return LogLevelNotice
	case "warn":
		return LogLevelWarn
	case "error":
		return LogLevelError
	case "panic":
		return LogLevelPanic
	case "fatal":
		return LogLevelFatal
	}
	return LogLevelNone
}

func NewFileLogger(path, prefix string, level LogLevel, skip int) (ILogger, error) {
	fileWriter, err := newFileWriter(path, prefix)
	if err != nil {
		return nil, err
	}
	writer := newLogWriter(fileWriter)
	writer.startRun()
	return &Logger{
		logLevel:       level,
		skipStackLevel: skip,
		logWriter:      writer,
		chain:          nil,
	}, nil
}

func NewConsoleLogger(level LogLevel, skip int) (ILogger, error) {
	writer := newLogWriter(os.Stdout)
	writer.startRun()
	return &Logger{
		logLevel:       level,
		skipStackLevel: skip,
		logWriter:      writer,
		chain:          nil,
	}, nil
}

func NewWriterLogger(w io.Writer, level LogLevel, skip int) (ILogger, error) {
	writer := newLogWriter(w)
	writer.startRun()
	return &Logger{
		logLevel:       level,
		skipStackLevel: skip,
		logWriter:      writer,
		chain:          nil,
	}, nil
}

func (l *Logger) AddLogChain(next ILogger) {
	if l.chain == nil {
		l.chain = next
	} else {
		l.chain.AddLogChain(next)
	}
}

func (l *Logger) log(format string, level LogLevel, args ...interface{}) {
	_, file, line, _ := runtime.Caller(2 + l.skipStackLevel)
	file = path.Base(file)
	obj := &logObject{
		file:    file,
		logTime: time.Now(),
		level:   level,
		line:    line,
		format:  format,
		args:    args,
	}
	select {
	case l.logWriter.objectQueue <- obj:
	default:
	}

	//format = fmt.Sprintf("%s %s %s:%d %s\n", time.Now().Format(time.RFC3339Nano), LogLevelToString(level), file, line, format)
	//fmt.Printf(format, args...)
}

func (l *Logger) logSync(format string, level LogLevel, args ...interface{}) {
	_, file, line, _ := runtime.Caller(2 + l.skipStackLevel)
	file = path.Base(file)
	obj := &logObject{
		file:    file,
		logTime: time.Now(),
		level:   level,
		line:    line,
		format:  format,
		args:    args,
	}
	l.logWriter.forceWriteObject(obj)
}

func (l *Logger) HitLevel(level LogLevel) bool {
	return l.logLevel <= level
}

func (l *Logger) Trace(format string, args ...interface{}) {
	if l.HitLevel(LogLevelTrace) {
		l.log(format, LogLevelTrace, args...)
	}
	if l.chain != nil {
		l.chain.Trace(format, args...)
	}
}

func (l *Logger) Debug(format string, args ...interface{}) {
	if l.HitLevel(LogLevelDebug) {
		l.log(format, LogLevelDebug, args...)
	}
	if l.chain != nil {
		l.chain.Debug(format, args...)
	}
}

func (l *Logger) Info(format string, args ...interface{}) {
	if l.HitLevel(LogLevelInfo) {
		l.log(format, LogLevelInfo, args...)
	}
	if l.chain != nil {
		l.chain.Info(format, args...)
	}
}

func (l *Logger) Notice(format string, args ...interface{}) {
	if l.HitLevel(LogLevelNotice) {
		l.log(format, LogLevelNotice, args...)
	}
	if l.chain != nil {
		l.chain.Notice(format, args...)
	}
}

func (l *Logger) Warn(format string, args ...interface{}) {
	if l.HitLevel(LogLevelWarn) {
		l.log(format, LogLevelWarn, args...)
	}
	if l.chain != nil {
		l.chain.Warn(format, args...)
	}
}

func (l *Logger) Error(format string, args ...interface{}) {
	if l.HitLevel(LogLevelError) {
		l.log(format, LogLevelError, args...)
	}
	if l.chain != nil {
		l.chain.Error(format, args...)
	}
}

func (l *Logger) Panic(format string, args ...interface{}) {
	if l.HitLevel(LogLevelPanic) {
		l.logSync(format, LogLevelPanic, args...)
	}
	panic(fmt.Sprintf(format, args...))
	// won't effect to log chain, for panic will stop this function
	// and maybe whole process
}

func (l *Logger) Fatal(format string, args ...interface{}) {
	if l.HitLevel(LogLevelFatal) {
		l.logSync(format, LogLevelFatal, args...)
	}
	os.Exit(-1)
	// won't effect to log chain, for fatal will stop process
}
