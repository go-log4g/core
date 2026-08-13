package log4g

import (
	"context"
	"log/slog"
	"runtime"
	"time"

	"github.com/go-log4g/core/impl"
)

var messageFormatter = impl.NewParameterizedMessageFormatter()

func Debug(pattern string, args ...any) {
	log(context.Background(), slog.LevelDebug, caller(), pattern, args...)
}

func Info(pattern string, args ...any) {
	log(context.Background(), slog.LevelInfo, caller(), pattern, args...)
}

func Warn(pattern string, args ...any) {
	log(context.Background(), slog.LevelWarn, caller(), pattern, args...)
}

func Error(pattern string, args ...any) {
	log(context.Background(), slog.LevelError, caller(), pattern, args...)
}

func DebugContext(ctx context.Context, pattern string, args ...any) {
	log(ctx, slog.LevelDebug, caller(), pattern, args...)
}

func InfoContext(ctx context.Context, pattern string, args ...any) {
	log(ctx, slog.LevelInfo, caller(), pattern, args...)
}

func WarnContext(ctx context.Context, pattern string, args ...any) {
	log(ctx, slog.LevelWarn, caller(), pattern, args...)
}

func ErrorContext(ctx context.Context, pattern string, args ...any) {
	log(ctx, slog.LevelError, caller(), pattern, args...)
}

func caller() uintptr {
	var pcs [1]uintptr
	runtime.Callers(3, pcs[:])
	return pcs[0]
}

func log(ctx context.Context, level slog.Level, pc uintptr, pattern string, args ...any) {
	logger := slog.Default()
	if !logger.Enabled(ctx, level) {
		return
	}

	record := slog.NewRecord(time.Now(), level, messageFormatter.Format(pattern, args...), pc)
	_ = logger.Handler().Handle(ctx, record)
}
