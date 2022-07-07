package logger

import (
	"context"
	"fmt"
	"os"

	"github.com/shivanshkc/rosenbridge/src/configs"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// zapLogger implements Logger using uber/zap package.
type zapLogger struct {
	client *zap.Logger
}

func (z *zapLogger) Debug(ctx context.Context, entry *Entry) {
	z.client.Debug("", zapFieldsFromEntry(ctx, entry)...)
}

func (z *zapLogger) Info(ctx context.Context, entry *Entry) {
	z.client.Info("", zapFieldsFromEntry(ctx, entry)...)
}

func (z *zapLogger) Warn(ctx context.Context, entry *Entry) {
	z.client.Warn("", zapFieldsFromEntry(ctx, entry)...)
}

func (z *zapLogger) Error(ctx context.Context, entry *Entry) {
	z.client.Error("", zapFieldsFromEntry(ctx, entry)...)
}

func (z *zapLogger) Fatal(ctx context.Context, entry *Entry) {
	z.client.Fatal("", zapFieldsFromEntry(ctx, entry)...)
}

func (z *zapLogger) Close() error {
	if err := z.client.Sync(); err != nil {
		return fmt.Errorf("failed to close the logger: %w", err)
	}
	// Closed successfully.
	return nil
}

// newZapLogger provides a new instance of zapLogger.
// Panic is allowed here because logger is crucial to the application.
func newZapLogger() *zapLogger {
	conf := configs.Get()

	// Converting the Log level from string to zapcore.Level.
	zapLevel, ok := zapLevelFromString(conf.Logger.Level)
	if !ok {
		panic(fmt.Sprintf("unknown log level: %s", conf.Logger.Level))
	}

	// This function lets Zap know which Level to log at.
	levelEnabler := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl >= zapLevel
	})

	// These are the various logging destinations. More destinations such as Kafka can be added here.
	syncers := []zapcore.WriteSyncer{os.Stdout}
	// Creating cores for the logger. See: https://pkg.go.dev/go.uber.org/zap#example-package-AdvancedConfiguration
	cores := make([]zapcore.Core, len(syncers))
	// Looping over all syncers and creating cores.
	for idx, dest := range syncers {
		encoder := zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig())
		cores[idx] = zapcore.NewCore(encoder, zapcore.Lock(dest), levelEnabler)
	}

	// zap.New call will give us a logger with the 'development' property as false.
	// We can turn it to true by using 'zap.NewDevelopment' instead of 'zap.New'.
	// If the development property is true, any DPanic level log calls will panic.
	client := zap.New(zapcore.NewTee(cores...))
	// Forming the zap implementation.
	return &zapLogger{client: client}
}

// zapLevelFromString converts a string level to zapcore.Level.
// If the conversion fails, the boolean is false.
func zapLevelFromString(level string) (zapcore.Level, bool) {
	switch level {
	case "debug":
		return zap.DebugLevel, true
	case "info":
		return zap.InfoLevel, true
	case "warn":
		return zap.WarnLevel, true
	case "error":
		return zap.ErrorLevel, true
	default:
		return zap.DPanicLevel, false
	}
}

// zapFieldsFromEntry processes the entry and converts it to a slice of Zap Fields.
func zapFieldsFromEntry(ctx context.Context, entry *Entry) []zap.Field {
	// Adding missing information to the log entry.
	entry.populate(ctx)
	// Creating zap fields out of the Entry.
	return []zap.Field{
		zap.Any("payload", entry.Payload),
		zap.Time("timestamp", entry.Timestamp),
		zap.Any("labels", entry.Labels),
		zap.Any("request", entry.Request),
		zap.Any("caller", entry.Caller),
		zap.String("trace", entry.Trace),
	}
}
