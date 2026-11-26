package log4g

import (
	"context"
	"log/slog"
	"runtime"
	"time"

	"github.com/go-log4g/core/impl"
)

var messageFormatter = impl.NewParameterizedMessageFormatter()

func Trace(pattern string, args ...any) {
	log(context.Background(), impl.LevelTrace, pattern, args...)
}

func Debug(pattern string, args ...any) {
	log(context.Background(), slog.LevelDebug, pattern, args...)
}

func Info(pattern string, args ...any) {
	log(context.Background(), slog.LevelInfo, pattern, args...)
}

func Warn(pattern string, args ...any) {
	log(context.Background(), slog.LevelWarn, pattern, args...)
}

func Error(pattern string, args ...any) {
	log(context.Background(), slog.LevelError, pattern, args...)
}

func Fatal(pattern string, args ...any) {
	log(context.Background(), impl.LevelFatal, pattern, args...)
}

func TraceContext(ctx context.Context, pattern string, args ...any) {
	log(ctx, impl.LevelTrace, pattern, args...)
}

func DebugContext(ctx context.Context, pattern string, args ...any) {
	log(ctx, slog.LevelDebug, pattern, args...)
}

func InfoContext(ctx context.Context, pattern string, args ...any) {
	log(ctx, slog.LevelInfo, pattern, args...)
}

func WarnContext(ctx context.Context, pattern string, args ...any) {
	log(ctx, slog.LevelWarn, pattern, args...)
}

func ErrorContext(ctx context.Context, pattern string, args ...any) {
	log(ctx, slog.LevelError, pattern, args...)
}

func FatalContext(ctx context.Context, pattern string, args ...any) {
	log(ctx, impl.LevelFatal, pattern, args...)
}

func log(ctx context.Context, level slog.Level, pattern string, args ...any) {
	logger := slog.Default()
	if !logger.Enabled(ctx, level) {
		return
	}

	record := slog.NewRecord(time.Now(), level, messageFormatter.Format(pattern, args...), caller())
	_ = logger.Handler().Handle(ctx, record)
}

func TraceEnabled() bool {
	return enabled(context.Background(), impl.LevelTrace)
}

func DebugEnabled() bool {
	return enabled(context.Background(), slog.LevelDebug)
}

func InfoEnabled() bool {
	return enabled(context.Background(), slog.LevelInfo)
}

func WarnEnabled() bool {
	return enabled(context.Background(), slog.LevelWarn)
}

func ErrorEnabled() bool {
	return enabled(context.Background(), slog.LevelError)
}

func FatalEnabled() bool {
	return enabled(context.Background(), impl.LevelFatal)
}

func Enabled(level slog.Level) bool {
	return enabled(context.Background(), level)
}

func enabled(ctx context.Context, level slog.Level) bool {
	logger := slog.Default()
	if !logger.Enabled(ctx, level) {
		return false
	}

	if handler, ok := logger.Handler().(*impl.Handler); ok {
		return handler.IsEnabled(ctx, level, caller())
	}

	return true
}

func caller() uintptr {
	var pcs [1]uintptr
	runtime.Callers(4, pcs[:])
	return pcs[0]
}
