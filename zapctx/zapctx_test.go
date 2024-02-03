package zapctx_test

import (
	"bytes"
	"context"
	"io"
	"testing"

	qt "github.com/frankban/quicktest"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/juju/zaputil/zapctx"
)

func TestLogger(t *testing.T) {
	c := qt.New(t)
	var buf bytes.Buffer
	logger := newLogger(&buf)
	ctx := zapctx.WithLogger(context.Background(), logger)
	zapctx.Logger(ctx).Info("hello")
	c.Assert(buf.String(), qt.Matches, `INFO\thello\n`)
}

func TestGetLoggerOrDefault(t *testing.T) {
	c := qt.New(t)
	var buf bytes.Buffer
	logger := newLogger(&buf)
	ctx := zapctx.WithLogger(context.Background(), logger)
	zapctx.GetLoggerOrDefault(ctx).Info("hello")
	c.Assert(buf.String(), qt.Matches, `INFO\thello\n`)
}

func TestGetLogger(t *testing.T) {
	c := qt.New(t)
	var buf bytes.Buffer
	logger := newLogger(&buf)
	ctx := zapctx.WithLogger(context.Background(), logger)
	zapctx.GetLoggerOrDefault(ctx).Info("hello")
	c.Assert(buf.String(), qt.Matches, `INFO\thello\n`)

	c.Assert(func() { zapctx.GetLogger(context.Background()) }, qt.PanicMatches, "logger not found in the context")
}

func TestGetLoggerOrError(t *testing.T) {
	c := qt.New(t)
	var buf bytes.Buffer
	logger := newLogger(&buf)
	ctx := zapctx.WithLogger(context.Background(), logger)
	zapctx.GetLoggerOrDefault(ctx).Info("hello")
	c.Assert(buf.String(), qt.Matches, `INFO\thello\n`)

	loggerOrError, err := zapctx.GetLoggerOrError(context.Background())
	c.Assert(loggerOrError, qt.IsNil)
	c.Assert(err, qt.ErrorMatches, "logger not found in the context")
}

func TestDefaultLogger(t *testing.T) {
	c := qt.New(t)
	var buf bytes.Buffer
	logger := newLogger(&buf)
	ctx := zapctx.WithLogger(context.Background(), logger)
	zapctx.Logger(ctx).Info("hello")
	c.Assert(buf.String(), qt.Matches, `INFO\thello\n`)
}
func TestDefaultGetLoggerOrDefault(t *testing.T) {
	c := qt.New(t)
	var buf bytes.Buffer
	logger := newLogger(&buf)
	ctx := zapctx.WithLogger(context.Background(), logger)
	zapctx.GetLoggerOrDefault(ctx).Info("hello")
	c.Assert(buf.String(), qt.Matches, `INFO\thello\n`)
}

func TestWithFields(t *testing.T) {
	c := qt.New(t)
	var buf bytes.Buffer
	logger := newLogger(&buf)

	ctx := zapctx.WithLogger(context.Background(), logger)
	ctx = zapctx.WithFields(ctx, zap.Int("foo", 999), zap.String("bar", "whee"))
	zapctx.GetLoggerOrDefault(ctx).Info("hello")
	c.Assert(buf.String(), qt.Matches, `INFO\thello\t\{"foo": 999, "bar": "whee"\}\n`)
}

func TestWithLevel(t *testing.T) {
	c := qt.New(t)
	var buf bytes.Buffer
	logger := newLogger(&buf)

	ctx := zapctx.WithLogger(context.Background(), logger)
	ctx1 := zapctx.WithLevel(ctx, zap.WarnLevel)
	zapctx.Info(ctx, "one")
	zapctx.Info(ctx1, "should not appear")
	zapctx.Warn(ctx1, "two")
	zapctx.Error(ctx1, "three")
	c.Assert(buf.String(), qt.Matches, `INFO\tone\nWARN\ttwo\nERROR\tthree\n`)
}

func TestMultistageSetupA(t *testing.T) {
	c := qt.New(t)
	var buf bytes.Buffer
	logger := newLogger(&buf)

	ctx := zapctx.WithLogger(context.Background(), logger)
	ctx = zapctx.WithLevel(ctx, zapcore.WarnLevel)
	ctx = zapctx.WithFields(ctx, zap.String("foo", "bar"))
	zapctx.Info(ctx, "one")
	zapctx.Warn(ctx, "two")
	c.Assert(buf.String(), qt.Matches, `WARN\ttwo\t{\"foo\": \"bar\"}\n`)
}

func TestMultistageSetupB(t *testing.T) {
	c := qt.New(t)
	var buf bytes.Buffer
	logger := newLogger(&buf)

	ctx := zapctx.WithLogger(context.Background(), logger)
	ctx = zapctx.WithFields(ctx, zap.String("foo", "bar"))
	ctx = zapctx.WithLevel(ctx, zapcore.WarnLevel)
	zapctx.Info(ctx, "one")
	zapctx.Warn(ctx, "two")
	c.Assert(buf.String(), qt.Matches, `WARN\ttwo\t{\"foo\": \"bar\"}\n`)
}

func newLogger(w io.Writer) *zap.Logger {
	config := zapcore.EncoderConfig{
		MessageKey:  "msg",
		LevelKey:    "level",
		EncodeLevel: zapcore.CapitalLevelEncoder,
	}
	core := zapcore.NewCore(
		zapcore.NewConsoleEncoder(config),
		zapcore.AddSync(w),
		zapcore.InfoLevel,
	)
	return zap.New(core)
}
