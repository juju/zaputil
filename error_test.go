package zaputil_test

import (
	"bytes"
	"io"
	"testing"

	qt "github.com/frankban/quicktest"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	errgo "gopkg.in/errgo.v1"

	"github.com/juju/zaputil"
)

func TestErrorJSONEncoder(t *testing.T) {
	c := qt.New(t)
	var buf bytes.Buffer
	logger := newJSONLogger(&buf)
	err := errgo.New("something")
	err = errgo.Mask(err)
	err = errgo.Notef(err, "an error")
	logger.Info("a message", zaputil.Error(err))
	c.Assert(buf.String(), qt.Matches, `\{"level":"info","ts":[0-9.]+,"msg":"a message","error":\{"msg":"an error: something","trace":\[\{"loc":".*zaputil/error_test.go:[0-9]+","msg":"an error"\},\{"loc":".*zaputil/error_test.go:[0-9]+"\},\{"loc":".*zaputil/error_test.go:[0-9]+","msg":"something"\}]\}\}\n`)
}

func TestConsoleEncoder(t *testing.T) {
	c := qt.New(t)
	var buf bytes.Buffer
	logger := newConsoleLogger(&buf)

	err := errgo.New("something")
	err = errgo.Mask(err)
	err = errgo.Notef(err, "an error")
	logger.Info("a message", zaputil.Error(err))
	c.Assert(buf.String(), qt.Matches, `[0-9.e+]+\tinfo\ta message\t\{"error": \{"msg":"an error: something","trace":\[\{"loc":".*zaputil/error_test.go:[0-9]+","msg":"an error"\},\{"loc":".*zaputil/error_test.go:[0-9]+"\},\{"loc":".*zaputil/error_test.go:[0-9]+","msg":"something"\}\]\}\}\n`)
}

func TestNilError(t *testing.T) {
	c := qt.New(t)
	var buf bytes.Buffer
	logger := newConsoleLogger(&buf)
	logger.Info("a message", zaputil.Error(nil))
	c.Assert(buf.String(), qt.Matches, `[0-9.e+]+\tinfo\ta message\n`)
}

func TestSimpleError(t *testing.T) {
	c := qt.New(t)
	var buf bytes.Buffer
	logger := newJSONLogger(&buf)
	logger.Info("a message", zaputil.Error(io.EOF))
	c.Assert(buf.String(), qt.Matches, `\{"level":"info","ts":[0-9.]+,"msg":"a message","error":\{"msg":"EOF"\}\}\n`)
}

var encoderConfig = zapcore.EncoderConfig{
	MessageKey:  "msg",
	LevelKey:    "level",
	TimeKey:     "ts",
	EncodeLevel: zapcore.LowercaseLevelEncoder,
	EncodeTime:  zapcore.EpochTimeEncoder,
}

func newJSONLogger(w io.Writer) *zap.Logger {
	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(encoderConfig),
		zapcore.AddSync(w),
		zapcore.InfoLevel,
	)
	return zap.New(core)
}

func newConsoleLogger(w io.Writer) *zap.Logger {
	return zap.New(
		zapcore.NewCore(
			zapcore.NewConsoleEncoder(encoderConfig),
			zapcore.AddSync(w),
			zapcore.InfoLevel,
		),
	)
}
