package log

import (
	"fmt"
	"os"
	"strings"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	logCfg zap.Config
	zapLvl zap.AtomicLevel
)

// var sentryClient *sentry.Client

// InitLogging initializes the logger.
// If the initialization sequence fails, the function will call os.Exit(1).
func InitLogging(logLvl, logEnv string) {
	if logEnv == "" {
		logEnv = "development"
	}

	if logLvl == "" {
		logLvl = "info"
	}

	lvl := decodeLogLevel(logLvl)
	zapLvl = zap.NewAtomicLevelAt(lvl)

	switch strings.ToLower(logEnv) {
	case "production":
		logCfg = newProductionConfig(zapLvl)
	case "development":
		logCfg = newDevelopmentConfig(zapLvl)

		// enable colored log output
		logCfg.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	default:
		fmt.Fprintf(os.Stderr, "Invalid log environment value %q provided. Exiting...\n", logEnv)
		os.Exit(1)
	}

	var err error

	logger, err = logCfg.Build(zap.AddStacktrace(zap.FatalLevel), zap.AddCallerSkip(1))
	if err != nil {
		panic(err)
	}
}

// SetLevel allows to set a new log level.
func SetLevel(newLevel string) {
	lvl := decodeLogLevel(newLevel)
	zapLvl.SetLevel(lvl)
}

// newProductionConfig is a reasonable production logging configuration.
// Logging is enabled at InfoLevel and above.
//
// It uses a JSON encoder, writes to standard error, and enables sampling.
// Stacktraces are automatically included on logs of ErrorLevel and above.
func newProductionConfig(lvl zap.AtomicLevel) zap.Config {
	return zap.Config{
		Level:       lvl,
		Development: false,
		Sampling: &zap.SamplingConfig{
			Initial:    100,
			Thereafter: 100,
		},
		Encoding:         "json",
		EncoderConfig:    zap.NewProductionEncoderConfig(),
		OutputPaths:      []string{"stderr"},
		ErrorOutputPaths: []string{"stderr"},
	}
}

// newDevelopmentConfig is a reasonable development logging configuration.
// Logging is enabled at DebugLevel and above.
//
// It enables development mode (which makes DPanicLevel logs panic), uses a
// console encoder, writes to standard error, and disables sampling.
// Stacktraces are automatically included on logs of WarnLevel and above.
func newDevelopmentConfig(lvl zap.AtomicLevel) zap.Config {
	return zap.Config{
		Level:            lvl,
		Development:      true,
		Encoding:         "console",
		EncoderConfig:    zap.NewDevelopmentEncoderConfig(),
		OutputPaths:      []string{"stderr"},
		ErrorOutputPaths: []string{"stderr"},
	}
}

// decodeLogLevel decodes a "log level string" into the zap.AtomicLevel.
func decodeLogLevel(l string) zapcore.Level {
	switch strings.ToLower(l) {
	case "trace":
		return zapcore.TraceLevel
	case "debug":
		return zapcore.DebugLevel
	case "info":
		return zapcore.InfoLevel
	case "warn":
		return zapcore.WarnLevel
	case "error":
		return zapcore.ErrorLevel
	case "fatal":
		return zapcore.FatalLevel
	default:
		fmt.Println("Unknown log level provided")
		os.Exit(1)
	}

	return zapcore.ErrorLevel
}

// logger is the logger instance.
var logger *zap.Logger

// traceProtol holds the trace protocol logger
// if `TraceProtocol` is ever used.
// This logger is used by any trace level function.
var traceProtocol *zap.Logger

// With creates a child logger and adds structured context to it. Fields added
// to the child don't affect the parent, and vice versa.
func With(fields ...zap.Field) *zap.Logger {
	return logger.With(fields...)
}

// Trace logs a message at TraceLevel. The message includes any fields passed
// at the log site, as well as any fields accumulated on the logger.
func Trace(msg string, fields ...zap.Field) {
	logger.Trace(msg, fields...)
}

// TraceProtocol logs a protocol message at TraceLevel.
// The message includes any fields passed at the log site, as well as any fields
// accumulated on the logger.
// Additionally, a new field called `trace_type` holding `protocol` as value will
// be passed to the logger.
func TraceProtocol(msg string, fields ...zap.Field) {
	if traceProtocol == nil {
		traceProtocol = With(zap.String("trace_type", "protocol"))
	}

	traceProtocol.Trace(msg, fields...)
}

// Debug logs a message at DebugLevel. The message includes any fields passed
// at the log site, as well as any fields accumulated on the logger.
func Debug(msg string, fields ...zap.Field) {
	logger.Debug(msg, fields...)
}

// Info logs a message at InfoLevel. The message includes any fields passed
// at the log site, as well as any fields accumulated on the logger.
func Info(msg string, fields ...zap.Field) {
	logger.Info(msg, fields...)
}

// Warn logs a message at WarnLevel. The message includes any fields passed
// at the log site, as well as any fields accumulated on the logger.
func Warn(msg string, fields ...zap.Field) {
	logger.Warn(msg, fields...)
}

// Error logs a message at ErrorLevel. The message includes any fields passed
// at the log site, as well as any fields accumulated on the logger.
func Error(msg string, fields ...zap.Field) {
	logger.Error(msg, fields...)
}

// DPanic logs a message at DPanicLevel. The message includes any fields
// passed at the log site, as well as any fields accumulated on the logger.
//
// If the logger is in development mode, it then panics (DPanic means
// "development panic"). This is useful for catching errors that are
// recoverable, but shouldn't ever happen.
func DPanic(msg string, fields ...zap.Field) {
	logger.DPanic(msg, fields...)
}

// Panic logs a message at PanicLevel. The message includes any fields passed
// at the log site, as well as any fields accumulated on the logger.
//
// The logger then panics, even if logging at PanicLevel is disabled.
func Panic(msg string, fields ...zap.Field) {
	logger.Panic(msg, fields...)
}

// Fatal logs a message at FatalLevel. The message includes any fields passed
// at the log site, as well as any fields accumulated on the logger.
//
// The logger then calls os.Exit(1), even if logging at FatalLevel is
// disabled.
func Fatal(msg string, fields ...zap.Field) {
	logger.Fatal(msg, fields...)
}

// Audit logs a message at AuditLevel. The message includes any fields passed
// at the log site, as well as any fields accumulated on the logger.
func Audit(msg string, fields ...zap.Field) {
	logger.Audit(msg, fields...)
}
