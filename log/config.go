package log

import (
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// 日志输出格式
const (
	// JSONFormat 以 JSON 形式输出
	JSONFormat = "json"
	// TextFormat 以 TEXT 形式输出
	TextFormat = "console"
)

type (
	// Level 类型别名，日志级别
	Level = zapcore.Level
	// Config 类型别名，日志配置
	Config = zap.Config
	// EncoderConfig 类型别名，棉麻配置
	EncoderConfig = zapcore.EncoderConfig
	// SamplingOption Sampling 配置参数
	SamplingOption struct {
		Tick       time.Duration
		Initial    int `json:"initial" yaml:"initial"`
		Thereafter int `json:"thereafter" yaml:"thereafter"`
	}
)

var defaultSamplingOption = &SamplingOption{
	Tick:       time.Second,
	Initial:    100,
	Thereafter: 100,
}

// NewDefaultEncoderConfig 创建默认编码配置信息
func NewDefaultEncoderConfig() EncoderConfig {
	return EncoderConfig{
		MessageKey:     "@message",
		LevelKey:       "@level",
		TimeKey:        "@timestamp",
		NameKey:        "@logger",
		CallerKey:      "@caller",
		StacktraceKey:  "@stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeTime:     zapcore.RFC3339TimeEncoder,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeDuration: zapcore.StringDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}
}

// NewDefaultConfig 默认日志配置
func NewDefaultConfig(format string, level Level) Config {
	return Config{
		Level:             zap.NewAtomicLevelAt(level),
		Development:       false,
		DisableCaller:     false,
		DisableStacktrace: true,
		Encoding:          format,
		EncoderConfig:     NewDefaultEncoderConfig(),
		OutputPaths:       []string{"stdout"},
		ErrorOutputPaths:  []string{"stderr"},
	}
}

// NewLoggerWithConfig 通过自定义配置创建 logger
func NewLoggerWithConfig(config Config) (*Logger, error) {
	logger, err := config.Build(zap.AddCallerSkip(CallerSkipOffset))
	if err != nil {
		return nil, err
	}
	return &Logger{
		level:  config.Level,
		logger: logger,
	}, nil
}
