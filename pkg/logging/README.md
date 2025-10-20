# Logging Package

This package provides a generic logging interface based on Logrus's `FieldLogger` interface.

## Overview

The `Logger` type is a type alias for `logrus.FieldLogger`, providing full compatibility with Logrus while maintaining flexibility for future implementations.

```go
type Logger = logrus.FieldLogger
```

## Features

- **Full Logrus Compatibility**: Direct integration with Logrus logger instances
- **Structured Logging**: Support for field-based structured logging via `WithField` and `WithFields`
- **Error Context**: Built-in error context support with `WithError`
- **Multiple Log Levels**: Debug, Info, Warn, Error, Fatal, and Panic levels
- **Flexible Formatting**: Support for formatted (`*f`), standard, and line (`*ln`) variants

## Usage

### Basic Usage

```go
import (
    "github.com/bdlilley/easygo/pkg/logging"
    "github.com/sirupsen/logrus"
)

// Create a Logrus logger
log := logrus.New()
log.SetLevel(logrus.InfoLevel)

// Use it as a Logger
var logger logging.Logger = log

logger.Info("application started")
```

### Structured Logging

```go
// Single field
logger.WithField("user_id", 12345).Info("user logged in")

// Multiple fields
logger.WithFields(logrus.Fields{
    "component": "auth",
    "action": "login",
    "ip": "192.168.1.1",
}).Info("authentication successful")
```

### Error Context

```go
err := someOperation()
if err != nil {
    logger.WithError(err).Error("operation failed")
}
```

### With AWS Client

```go
import (
    "context"
    "github.com/bdlilley/easygo"
    "github.com/sirupsen/logrus"
)

log := logrus.New()
log.SetLevel(logrus.DebugLevel)

client, err := easygo.NewAwsClient(context.Background(), &easygo.NewEGAwsClientArgs{
    Logger: log,
    Region: "us-west-2",
})
```

### Persistent Fields

Create a logger with fields that persist across all log entries:

```go
// Create a logger with default fields
logger := log.WithFields(logrus.Fields{
    "service": "my-service",
    "version": "1.0.0",
})

// All logs will include service and version
logger.Info("started")  // includes service and version fields
logger.Error("failed")   // includes service and version fields
```

## Available Methods

The `Logger` interface includes all methods from `logrus.FieldLogger`:

### Structured Logging
- `WithField(key string, value interface{}) *Entry`
- `WithFields(fields Fields) *Entry`
- `WithError(err error) *Entry`

### Log Levels
- `Debug(args ...interface{})`
- `Info(args ...interface{})`
- `Warn(args ...interface{})`
- `Error(args ...interface{})`
- `Fatal(args ...interface{})`
- `Panic(args ...interface{})`

### Formatted Logging
- `Debugf(format string, args ...interface{})`
- `Infof(format string, args ...interface{})`
- `Warnf(format string, args ...interface{})`
- `Errorf(format string, args ...interface{})`
- `Fatalf(format string, args ...interface{})`
- `Panicf(format string, args ...interface{})`

### Line-based Logging
- `Debugln(args ...interface{})`
- `Infoln(args ...interface{})`
- `Warnln(args ...interface{})`
- `Errorln(args ...interface{})`
- `Fatalln(args ...interface{})`
- `Panicln(args ...interface{})`

## Configuration Examples

### JSON Formatter
```go
log := logrus.New()
log.SetFormatter(&logrus.JSONFormatter{
    TimestampFormat: "2006-01-02 15:04:05",
    PrettyPrint:     false,
})
```

### Text Formatter
```go
log := logrus.New()
log.SetFormatter(&logrus.TextFormatter{
    FullTimestamp:   true,
    TimestampFormat: "2006-01-02 15:04:05",
})
```

### Log Level Configuration
```go
log := logrus.New()
log.SetLevel(logrus.DebugLevel) // Or InfoLevel, WarnLevel, ErrorLevel, etc.
```

## Why Logrus?

- **Industry Standard**: Widely used structured logging library in Go
- **Rich Feature Set**: Comprehensive logging capabilities out of the box
- **Flexible**: Multiple formatters, hooks, and customization options
- **Production Ready**: Battle-tested in production environments
- **Extensible**: Support for custom formatters and hooks

## See Also

- [Logrus Documentation](https://github.com/sirupsen/logrus)
- [Example Usage](example_test.go)
