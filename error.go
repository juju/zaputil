package zaputil

import (
	"encoding/json"
	"fmt"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/errgo.v1"
)

// Error returns a field suitable for logging an error
// to a zap Logger along with its error trace.
// If err is nil, the field is a no-op.
// This is different from zap.Error because the logged
// error also includes error traceback information.
func Error(err error) zapcore.Field {
	if err == nil {
		return zap.Skip()
	}
	return zap.Reflect("error", errObject{err})
}

// errObject is the type stored in the zap.Field. It implements
// MarshalJSON and also formats decently as an error (for example when
// used by zap.TextEncoder).
type errObject struct {
	error
}

// MarshalJSON implements json.Marshaler.MarshalJSON to show the details
// of the error.
func (e errObject) MarshalJSON() ([]byte, error) {
	if e.error == nil {
		return []byte("null"), nil
	}
	return json.Marshal(jsonErr{
		Message: e.error.Error(),
		Trace:   errorTrace(e.error),
	})
}

// jsonErr is the actual type used for JSON-marshaling errors.
type jsonErr struct {
	Message string           `json:"msg"`
	Trace   []jsonTraceLevel `json:"trace,omitempty"`
}

// jsonTraceLevel represents one level of an error trace.
type jsonTraceLevel struct {
	Location string `json:"loc,omitempty"`
	Message  string `json:"msg,omitempty"`
}

// errorTrace returns the error's trace suitable for
// marshaling as JSON.
func errorTrace(err error) []jsonTraceLevel {
	trace := make([]jsonTraceLevel, 0, 10)
	for err != nil {
		var t jsonTraceLevel
		if err, ok := err.(errgo.Locationer); ok {
			if file, line := err.Location(); file != "" {
				t.Location = fmt.Sprintf("%s:%d", file, line)
			}
		}
		if werr, ok := err.(errgo.Wrapper); ok {
			t.Message = werr.Message()
			err = werr.Underlying()
		} else {
			if len(trace) == 0 && t.Location == "" {
				// There's no location or underlying error,
				// so the trace isn't adding anything we
				// won't already see.
				return nil
			}
			t.Message = err.Error()
			err = nil
		}
		trace = append(trace, t)
	}
	return trace
}
