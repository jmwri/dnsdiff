package log

import (
	"context"
	"go.uber.org/zap"
)

type logContextKey int

const (
	loggerKey logContextKey = iota
)

// WithFields adds the provided zap fields to the logger within the given context.
func WithFields(ctx context.Context, fields ...zap.Field) context.Context {
	return context.WithValue(ctx, loggerKey, fromContext(ctx).With(fields...))
}

// fromContext returns the logger from the given context, or returns the base logger.
func fromContext(ctx context.Context) *zap.Logger {
	if ctx == nil {
		return baseLogger
	}
	if ctxLogger, ok := ctx.Value(loggerKey).(*zap.Logger); ok {
		return ctxLogger
	}
	return baseLogger
}
