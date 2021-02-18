package log

import (
	"context"
)

const (
	// SampleCallerSkipOffset 通过内置 logger 输出日志， Caller + 1
	SampleCallerSkipOffset = 1
)

var (
	defaultLogger = NewLogger(JSONFormat, InfoLevel)
	simpleLogger  = defaultLogger.AddCallerSkip(SampleCallerSkipOffset)
)

var (
	// Debug 内置 logger 以 debug 等级输出日志
	Debug = simpleLogger.Debug
	// DebugF 通过内置 logger， 以 debug 等级格式化输出日志
	DebugF = simpleLogger.DebugF
	// DebugWithFields 通过内置 logger， 以 debug 等级输出日志，并附加 fields 信息
	DebugWithFields = simpleLogger.DebugWithField

	// Info 内置 logger 以 info 等级输出日志
	Info = simpleLogger.Info
	// InfoF 通过内置 logger， 以 info 等级格式化输出日志
	InfoF = simpleLogger.InfoF
	// InfoWithFields 通过内置 logger， 以 info 等级输出日志，并附加 fields 信息
	InfoWithFields = simpleLogger.InfoWithField

	// Warn 内置 logger 以 warn 等级输出日志
	Warn = simpleLogger.Warn
	// WarnF 通过内置 logger， 以 warn 等级格式化输出日志
	WarnF = simpleLogger.WarnF
	// WarnWithFields 通过内置 logger， 以 warn 等级输出日志，并附加 fields 信息
	WarnWithFields = simpleLogger.WarnWithField

	// Error 内置 logger 以 error 等级输出日志
	Error = simpleLogger.Error
	// ErrorF 通过内置 logger， 以 error 等级格式化输出日志
	ErrorF = simpleLogger.ErrorF
	// ErrorWithFields 通过内置 logger， 以 error 等级输出日志，并附加 fields 信息
	ErrorWithFields = simpleLogger.ErrorWithField
)

// SetLevel 设置默认 logger 输出级别
func SetLevel(level Level) {
	defaultLogger.SetLevel(level)
}

// DebugCtx 通过 ctx 获取 logger， 并以 debug 等级输出日志
func DebugCtx(ctx context.Context, message string) {
	loadLogger(ctx).Debug(message)
}

// DebugCtxF 通过 ctx 获取 logger， 并以 debug 等级格式化输出日志
func DebugCtxF(ctx context.Context, format string, args ...interface{}) {
	loadLogger(ctx).DebugF(format, args...)
}

// DebugCtxWithFields 通过 ctx 获取 logger， 以 debug 等级输出日志，并附加 fields 信息
func DebugCtxWithFields(ctx context.Context, message string, fields Fields) {
	loadLogger(ctx).DebugWithField(message, fields)
}

func InfoCtx(ctx context.Context, message string) {
	loadLogger(ctx).Info(message)
}
func InfoCtxF(ctx context.Context, format string, args ...interface{}) {
	loadLogger(ctx).InfoF(format, args...)
}

func InfoCtxWithFields(ctx context.Context, message string, fields Fields) {
	loadLogger(ctx).InfoWithField(message, fields)
}

func WarnCtx(ctx context.Context, message string) {
	loadLogger(ctx).Warn(message)
}

func WarnCtxF(ctx context.Context, format string, args ...interface{}) {
	loadLogger(ctx).WarnF(format, args...)
}

func WarnCtxWithFields(ctx context.Context, message string, fields Fields) {
	loadLogger(ctx).WarnWithField(message, fields)
}

func ErrorCtx(ctx context.Context, message string) {
	loadLogger(ctx).Error(message)
}

func ErrorCtxF(ctx context.Context, format string, args ...interface{}) {
	loadLogger(ctx).ErrorF(format, args...)
}

func ErrorCtxWithFields(ctx context.Context, message string, fields Fields) {
	loadLogger(ctx).ErrorWithField(message, fields)
}

func Panic(message string) {
	simpleLogger.Panic(message)
}
func PanicCtx(ctx context.Context, message string) {
	loadLogger(ctx).Panic(message)
}

// PanicF 通过内置 logger，以 panic 等级格式化输出日志
func PanicF(format string, args ...interface{}) {
	simpleLogger.PanicF(format, args...)
}

// PanicCtxF 通过 ctx 获取 logger，以 panic 等级格式化输出日志
func PanicCtxF(ctx context.Context, format string, args ...interface{}) {
	loadLogger(ctx).PanicF(format, args...)
}

// PanicWithFields 通过内置 logger，以 panic 等级输出日志，并 附加 fields 信息
func PanicWithFields(message string, fields Fields) {
	simpleLogger.PanicWithField(message, fields)
}

// PanicCtxWithFields 通过 ctx  获取logger，以 panic 等级输出日志，并 附加 fields 信息
func PanicCtxWithFields(ctx context.Context, message string, fields Fields) {
	loadLogger(ctx).PanicWithField(message, fields)
}

// Fatal 使用内置logger，以 fatal 等级输出日志
func Fatal(message string) {
	simpleLogger.Fatal(message)
}

// FatalCtx 通过 ctx 获取 logger，以 fatal 等级输出日志
func FatalCtx(ctx context.Context, message string) {
	loadLogger(ctx).Fatal(message)
}

// FatalF 通过 ctx 获取 logger，以 fatal 等级格式化输出日志
func FatalF(format string, args ...interface{}) {
	simpleLogger.FatalF(format, args...)
}

// FatalCtxF 通过 ctx 获取 logger，以 fatal 等级格式化输出日志
func FatalCtxF(ctx context.Context, format string, args ...interface{}) {
	loadLogger(ctx).FatalF(format, args...)
}

// FatalWithFields 通过 内置 logger，以 fatal 等级输出日志，并 附加 fields 信息
func FatalWithFields(message string, fields Fields) {
	simpleLogger.FatalWithField(message, fields)
}

// FatalCtxWithFields 通过 ctx 获取 logger，以 fatal 等级输出日志，并 附加 fields 信息
func FatalCtxWithFields(ctx context.Context, message string, fields Fields) {
	loadLogger(ctx).FatalWithField(message, fields)
}

func loadLogger(ctx context.Context) *Logger {
	if logger := extractLogger(ctx); logger != nil {
		return logger.AddCallerSkip(1)
	}
	return simpleLogger
}
