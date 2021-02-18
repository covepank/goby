package log

import (
	"fmt"
	"io"
	"net/http"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const (
	DebugLevel = zap.DebugLevel
	InfoLevel  = zap.InfoLevel
	WarnLevel  = zap.WarnLevel
	ErrorLevel = zap.ErrorLevel
	PanicLevel = zap.PanicLevel
	FatalLevel = zap.FatalLevel
)

const (
	IdKey    = "@logId"
	ErrorKey = "@error"
	ScopeKey = "@scope"
)

// CallerSkipOffset 调用栈便宜，用于输出日志位置
const CallerSkipOffset = 2

type (
	Logger struct {
		level  zap.AtomicLevel
		logger *zap.Logger
	}
	Fields map[string]interface{}
)

func NewLogger(format string, level Level) *Logger {
	config := NewDefaultConfig(format, level)
	config.Level = zap.NewAtomicLevelAt(level)
	logger, _ := NewLoggerWithConfig(config)
	return logger

}

// SetLevel 设置当前 logger 以所有派生 logger 日志等级
func (l *Logger) SetLevel(level Level) {
	l.level.SetLevel(level)
}

// Sync 刷新数据
// flushing any buffered log entries, take care to call sync before exiting.
func (l *Logger) Sync() {
	_ = l.logger.Sync()
}

// SetLevelHandler 设置 Web handler, 动态改变日志输出等级
func (l *Logger) SetLevelHandler(w http.ResponseWriter, r *http.Request) {
	l.level.ServeHTTP(w, r)
}

// SetName 派生 Logger 并设置 Logger Name
func (l *Logger) SetName(name string) *Logger {
	return &Logger{
		level:  l.level,
		logger: l.logger.Named(name),
	}
}

// WithLogId 写入 LogID
func (l *Logger) WithLogId(id string) *Logger {
	logger := l.logger.With(zap.String(IdKey, id))
	return &Logger{
		level:  l.level,
		logger: logger,
	}
}

// NewNamespace 派生 Logger 并创建一个 NameSpace
func (l *Logger) NewNamespace(name string) *Logger {
	logger := l.logger.With(zap.Namespace(name))

	return &Logger{
		level:  l.level,
		logger: logger,
	}
}

// WithOutFile 派生一个Logger，将日志写入指定文件列表
func (l *Logger) WithOutFile(outputPaths []string, level Level) (*Logger, func(), error) {
	sink, closeOut, err := zap.Open(outputPaths...)
	if err != nil {
		return nil, nil, err
	}
	return l.WithOutput(sink, level, TextFormat), closeOut, nil

}

// WithOutput 派生一个Logger，附加 writer， 收集日志信息
func (l *Logger) WithOutput(writer io.Writer, level Level, format string) *Logger {
	enc := func() zapcore.Encoder {
		if format == JSONFormat {
			return zapcore.NewJSONEncoder(NewDefaultEncoderConfig())
		}
		return zapcore.NewConsoleEncoder(NewDefaultEncoderConfig())
	}()

	logger := l.logger.WithOptions(zap.WrapCore(func(core zapcore.Core) zapcore.Core {
		return zapcore.NewTee(core, zapcore.NewCore(enc, zapcore.AddSync(writer), level))
	}))

	return &Logger{
		level:  l.level,
		logger: logger,
	}
}

// WithError 派生一个 Logger，并附加 Error 信息
func (l *Logger) WithError(err error) *Logger {
	if err != nil {
		logger := l.logger.With(zap.NamedError(ErrorKey, err))
		return &Logger{
			level:  l.level,
			logger: logger,
		}
	}
	return l
}

// WithScope 派生一个 Logger，并记录 Scope
func (l *Logger) WithScope(scope string) *Logger {
	logger := l.logger.With(zap.String(ScopeKey, scope))
	return &Logger{
		level:  l.level,
		logger: logger,
	}
}

// 设置 Key-Value 对
func (l *Logger) WithField(key string, value interface{}) *Logger {
	logger := l.logger.With(zap.Any(key, value))
	return &Logger{
		level:  l.level,
		logger: logger,
	}
}

// 数据 Debug 日志
func (l *Logger) Debug(message string) {
	l.log(zapcore.DebugLevel, message, nil, nil)
}

func (l *Logger) DebugF(format string, args ...interface{}) {
	l.log(zapcore.DebugLevel, format, args, nil)
}

// 输出自定义 Key-Value 对
// 性能优于 WithField
func (l *Logger) DebugWithField(message string, fields Fields) {
	l.log(zapcore.DebugLevel, message, nil, fields)
}

// 输出 Info 日志
func (l *Logger) Info(message string) {
	l.log(zapcore.InfoLevel, message, nil, nil)
}

func (l *Logger) InfoF(format string, args ...interface{}) {
	l.log(zapcore.InfoLevel, format, args, nil)
}

func (l *Logger) InfoWithField(message string, fields Fields) {
	l.log(zapcore.InfoLevel, message, nil, fields)
}

// 输出 Warn 日志
func (l *Logger) Warn(message string) {
	l.log(zapcore.WarnLevel, message, nil, nil)
}

func (l *Logger) WarnF(format string, args ...interface{}) {
	l.log(zapcore.WarnLevel, format, args, nil)
}

func (l *Logger) WarnWithField(message string, fields Fields) {
	l.log(zapcore.WarnLevel, message, nil, fields)
}

// 输出 Error 日志
func (l *Logger) Error(message string) {
	l.log(zapcore.ErrorLevel, message, nil, nil)
}

func (l *Logger) ErrorF(format string, args ...interface{}) {
	l.log(zapcore.ErrorLevel, format, args, nil)
}

func (l *Logger) ErrorWithField(message string, fields Fields) {
	l.log(zapcore.ErrorLevel, message, nil, fields)
}

// Panic 输出 panic 日志
func (l *Logger) Panic(message string) {
	l.log(zapcore.PanicLevel, message, nil, nil)
}

func (l *Logger) PanicF(format string, args ...interface{}) {
	l.log(zapcore.PanicLevel, format, args, nil)
}

func (l *Logger) PanicWithField(message string, fields Fields) {
	l.log(zapcore.PanicLevel, message, nil, fields)
}

// Fatal 输出 Fatal 日志
func (l *Logger) Fatal(message string) {
	l.log(zapcore.FatalLevel, message, nil, nil)
}

func (l *Logger) FatalF(format string, args ...interface{}) {
	l.log(zapcore.FatalLevel, format, args, nil)
}

func (l *Logger) FatalWithField(message string, fields Fields) {
	l.log(zapcore.FatalLevel, message, nil, fields)
}

func (l *Logger) log(level zapcore.Level, format string, args []interface{}, fields Fields) {
	if !l.logger.Core().Enabled(level) {
		return
	}

	message := format
	if len(args) > 0 {
		message = fmt.Sprintf(format, args...)
	}

	f := make([]zap.Field, 0, len(fields))
	for key := range fields {
		f = append(f, zap.Any(key, fields[key]))
	}

	if ce := l.logger.Check(level, message); ce != nil {
		ce.Write(f...)
	}
}

// Write impl io.Writer
func (l *Logger) Write(msg []byte) (int, error) {
	l.Info(string(msg))
	return len(msg), nil
}

// AddCallerSkip 创建新的 logger 并且增加 CallerSkip
func (l *Logger) AddCallerSkip(skip int) *Logger {
	return &Logger{
		level:  l.level,
		logger: l.logger.WithOptions(zap.AddCallerSkip(skip)),
	}
}

func (l *Logger) WithSampling(opts *SamplingOption) *Logger {
	if opts == nil {
		opts = defaultSamplingOption
	}

	logger := l.logger.WithOptions(zap.WrapCore(func(core zapcore.Core) zapcore.Core {
		return zapcore.NewSamplerWithOptions(core, opts.Tick, opts.Initial, opts.Thereafter)
	}))

	return &Logger{
		level:  l.level,
		logger: logger,
	}
}
