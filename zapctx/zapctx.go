// Copyright 2016 Canonical Ltd.

// Package zapctx provides support for associating zap loggers
// (see github.com/uber-go/zap) with contexts.
package zapctx

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"golang.org/x/net/context"
)

// LogLevel holds an AtomicLevel that can be used to change the logging
// level of Default.
var LogLevel = zap.NewAtomicLevel()

// Default holds the logger returned by Logger when there is no logger in
// the context. If replacing Default with a new Logger then consider
// using &LogLevel as the LevelEnabler so that SetLevel can still be used
// to dynamically change the logging level.
var Default = zap.New(
	zapcore.NewCore(
		zapcore.NewJSONEncoder(zapcore.EncoderConfig{
			MessageKey:  "msg",
			LevelKey:    "level",
			TimeKey:     "ts",
			EncodeLevel: zapcore.LowercaseLevelEncoder,
			EncodeTime:  zapcore.ISO8601TimeEncoder,
		}),
		os.Stdout,
		&LogLevel,
	),
)

// loggerKey holds the context key used for loggers.
type loggerKey struct{}

// WithLogger returns a new context derived from ctx that
// is associated with the given logger.
func WithLogger(ctx context.Context, logger *zap.Logger) context.Context {
	return context.WithValue(ctx, loggerKey{}, logger)
}

// WithFields returns a new context derived from ctx
// that has a logger that always logs the given fields.
func WithFields(ctx context.Context, fields ...zapcore.Field) context.Context {
	return WithLogger(ctx, Logger(ctx).With(fields...))
}

// WithLevel returns a new context derived from ctx
// that has a logger that only logs messages at or above
// the given level.
func WithLevel(ctx context.Context, level zapcore.Level) context.Context {
	return WithLogger(ctx, Logger(ctx).WithOptions(wrapCoreWithLevel(level)))
}

func wrapCoreWithLevel(level zapcore.Level) zap.Option {
	return zap.WrapCore(func(core zapcore.Core) zapcore.Core {
		return &coreWithLevel{
			Core:  core,
			level: level,
		}
	})
}

type coreWithLevel struct {
	zapcore.Core
	level zapcore.Level
}

func (c *coreWithLevel) Enabled(level zapcore.Level) bool {
	return level >= c.level && c.Core.Enabled(level)
}

func (c *coreWithLevel) Check(e zapcore.Entry, ce *zapcore.CheckedEntry) *zapcore.CheckedEntry {
	if !c.level.Enabled(e.Level) {
		return ce
	}
	return c.Core.Check(e, ce)
}

// Logger returns the logger associated with the given
// context. If there is no logger, it will return Default.
func Logger(ctx context.Context) *zap.Logger {
	if ctx == nil {
		panic("nil context passed to Logger")
	}
	if logger, _ := ctx.Value(loggerKey{}).(*zap.Logger); logger != nil {
		return logger
	}
	return Default
}

// Debug calls Logger(ctx).Debug(msg, fields...).
func Debug(ctx context.Context, msg string, fields ...zapcore.Field) {
	Logger(ctx).Debug(msg, fields...)
}

// Info calls Logger(ctx).Info(msg, fields...).
func Info(ctx context.Context, msg string, fields ...zapcore.Field) {
	Logger(ctx).Info(msg, fields...)
}

// Warn calls Logger(ctx).Warn(msg, fields...).
func Warn(ctx context.Context, msg string, fields ...zapcore.Field) {
	Logger(ctx).Warn(msg, fields...)
}

// Error calls Logger(ctx).Error(msg, fields...).
func Error(ctx context.Context, msg string, fields ...zapcore.Field) {
	Logger(ctx).Error(msg, fields...)
}
