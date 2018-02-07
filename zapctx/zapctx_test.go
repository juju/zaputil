package zapctx_test

import (
	"bytes"
	"io"

	"github.com/juju/testing"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"golang.org/x/net/context"
	gc "gopkg.in/check.v1"

	"github.com/juju/zaputil/zapctx"
)

type zapctxSuite struct {
	testing.CleanupSuite
}

var _ = gc.Suite(&zapctxSuite{})

func (*zapctxSuite) TestLogger(c *gc.C) {
	var buf bytes.Buffer
	logger := newLogger(&buf)
	ctx := zapctx.WithLogger(context.Background(), logger)
	zapctx.Logger(ctx).Info("hello")
	c.Assert(buf.String(), gc.Matches, `INFO\thello\n`)
}

func (s *zapctxSuite) TestDefaultLogger(c *gc.C) {
	var buf bytes.Buffer
	logger := newLogger(&buf)

	s.PatchValue(&zapctx.Default, logger)
	zapctx.Logger(context.Background()).Info("hello")
	c.Assert(buf.String(), gc.Matches, `INFO\thello\n`)
}

func (*zapctxSuite) TestWithFields(c *gc.C) {
	var buf bytes.Buffer
	logger := newLogger(&buf)

	ctx := zapctx.WithLogger(context.Background(), logger)
	ctx = zapctx.WithFields(ctx, zap.Int("foo", 999), zap.String("bar", "whee"))
	zapctx.Logger(ctx).Info("hello")
	c.Assert(buf.String(), gc.Matches, `INFO\thello\t\{"foo": 999, "bar": "whee"\}\n`)
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
