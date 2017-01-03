// Copyright 2016 Canonical Ltd.

// Package zapctx provides support for associating zap loggers
// (see github.com/uber-go/zap) with contexts.
package zapctx

import (
	"os"

	"github.com/uber-go/zap"
	"golang.org/x/net/context"
)

// Default holds the logger returned by Logger when
// there is no logger in the context.
var Default = zap.New(zap.NewJSONEncoder(), zap.Output(os.Stderr))

// loggerKey holds the context key used for loggers.
type loggerKey struct{}

// WithLogger returns a new context derived from ctx that
// is associated with the given logger.
func WithLogger(ctx context.Context, logger zap.Logger) context.Context {
	return context.WithValue(ctx, loggerKey{}, logger)
}

// WithFields returns a new context derived from ctx
// that has a logger that always logs the given fields.
func WithFields(ctx context.Context, fields ...zap.Field) context.Context {
	return WithLogger(ctx, Logger(ctx).With(fields...))
}

// Logger returns the logger associated with the given
// context. If there is no logger, it will return Default.
func Logger(ctx context.Context) zap.Logger {
	if ctx == nil {
		panic("nil context passed to Logger")
	}
	if logger, _ := ctx.Value(loggerKey{}).(zap.Logger); logger != nil {
		return logger
	}
	return Default
}

// Debug calls Logger(ctx).Debug(msg, fields...).
func Debug(ctx context.Context, msg string, fields ...zap.Field) {
	Logger(ctx).Debug(msg, fields...)
}

// Info calls Logger(ctx).Info(msg, fields...).
func Info(ctx context.Context, msg string, fields ...zap.Field) {
	Logger(ctx).Info(msg, fields...)
}

// Warn calls Logger(ctx).Warn(msg, fields...).
func Warn(ctx context.Context, msg string, fields ...zap.Field) {
	Logger(ctx).Warn(msg, fields...)
}

// Error calls Logger(ctx).Error(msg, fields...).
func Error(ctx context.Context, msg string, fields ...zap.Field) {
	Logger(ctx).Error(msg, fields...)
}
