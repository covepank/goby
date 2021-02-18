package log

import (
	"context"
)

// Key 内部key类型
type Key string

const (
	// ContextLoggerKey 在 context 中存储 logger 的 key
	ContextLoggerKey Key = "context.logger.key"
)

// WithLogger 将 logger 存储到指定 context 中
func WithLogger(ctx context.Context, logger *Logger) context.Context {
	return context.WithValue(ctx, ContextLoggerKey, logger)
}

// ExtractLogger 从 ctx 中提取 logger
func ExtractLogger(ctx context.Context) *Logger {
	if logger := extractLogger(ctx); logger != nil {
		return logger
	}
	return defaultLogger.SetName("default")
}

func extractLogger(ctx context.Context) *Logger {
	if l := ctx.Value(ContextLoggerKey); l != nil {
		if logger, ok := l.(*Logger); ok {
			return logger
		}
	}

	return nil
}
