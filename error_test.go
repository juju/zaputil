package zaputil_test

import (
	"bytes"
	"io"

	"github.com/uber-go/zap"
	gc "gopkg.in/check.v1"
	errgo "gopkg.in/errgo.v1"

	"github.com/juju/zaputil"
)

type zaputilSuite struct{}

var _ = gc.Suite(&zaputilSuite{})

func (*zaputilSuite) TestErrorJSONEncoder(c *gc.C) {
	var buf bytes.Buffer
	logger := zap.New(zap.NewJSONEncoder(), zap.Output(zap.AddSync(&buf)))

	err := errgo.New("something")
	err = errgo.Mask(err)
	err = errgo.Notef(err, "an error")
	logger.Info("a message", zaputil.Error(err))
	c.Assert(buf.String(), gc.Matches, `\{"level":"info","ts":[0-9.]+,"msg":"a message","error":\{"msg":"an error: something","trace":\[\{"loc":".*zaputil/error_test.go:[0-9]+","msg":"an error"\},\{"loc":".*zaputil/error_test.go:[0-9]+"\},\{"loc":".*zaputil/error_test.go:[0-9]+","msg":"something"\}]\}\}\n`)
}

func (*zaputilSuite) TestTextEncoder(c *gc.C) {
	var buf bytes.Buffer
	logger := zap.New(zap.NewTextEncoder(), zap.Output(zap.AddSync(&buf)))

	err := errgo.New("something")
	err = errgo.Mask(err)
	err = errgo.Notef(err, "an error")
	logger.Info("a message", zaputil.Error(err))
	c.Assert(buf.String(), gc.Matches, `\[I\] .* a message error=an error: something\n`)
}

func (*zaputilSuite) TestNilError(c *gc.C) {
	var buf bytes.Buffer
	logger := zap.New(zap.NewTextEncoder(), zap.Output(zap.AddSync(&buf)))

	logger.Info("a message", zaputil.Error(nil))
	c.Assert(buf.String(), gc.Matches, `\[I\] .* a message\n`)
}

func (*zaputilSuite) TestSimpleError(c *gc.C) {
	var buf bytes.Buffer
	logger := zap.New(zap.NewJSONEncoder(), zap.Output(zap.AddSync(&buf)))

	logger.Info("a message", zaputil.Error(io.EOF))
	c.Assert(buf.String(), gc.Matches, `\{"level":"info","ts":[0-9.]+,"msg":"a message","error":\{"msg":"EOF"\}\}\n`)
}
