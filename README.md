# zapctx
--
    import "github.com/juju/zapctx"

Package zapctx provides support for associating zap loggers (see
github.com/uber-go/zap) with contexts.

## Usage

```go
var Default = zap.New(zap.NewJSONEncoder(), zap.Output(os.Stderr))
```
Default holds the logger returned by Logger when there is no logger in the
context.

#### func  Debug

```go
func Debug(ctx context.Context, msg string, fields ...zap.Field)
```
Debug calls Logger(ctx).Debug(msg, fields...).

#### func  Error

```go
func Error(ctx context.Context, msg string, fields ...zap.Field)
```
Error calls Logger(ctx).Error(msg, fields...).

#### func  Info

```go
func Info(ctx context.Context, msg string, fields ...zap.Field)
```
Info calls Logger(ctx).Info(msg, fields...).

#### func  Logger

```go
func Logger(ctx context.Context) zap.Logger
```
Logger returns the logger associated with the given context. If there is no
logger, it will return Default.

#### func  Warn

```go
func Warn(ctx context.Context, msg string, fields ...zap.Field)
```
Warn calls Logger(ctx).Warn(msg, fields...).

#### func  WithFields

```go
func WithFields(ctx context.Context, fields ...zap.Field) context.Context
```
WithFields returns a new context derived from ctx that has a logger that always
logs the given fields.

#### func  WithLogger

```go
func WithLogger(ctx context.Context, logger zap.Logger) context.Context
```
WithLogger returns a new context derived from ctx that is associated with the
given logger.
