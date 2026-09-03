package log4g_test

import (
	"context"
	"log/slog"
	"runtime"
	"strings"
	"testing"

	"github.com/go-log4g/core/impl"
	"github.com/go-log4g/core/log4g"
	"github.com/go-log4g/core/mdc"
	"github.com/stretchr/testify/require"
)

type captureHandler struct {
	ctx    context.Context
	record slog.Record
}

func (this *captureHandler) Enabled(ctx context.Context, level slog.Level) bool {
	return true
}

func (this *captureHandler) Handle(ctx context.Context, record slog.Record) error {
	this.ctx = ctx
	this.record = record
	return nil
}

func (this *captureHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return this
}

func (this *captureHandler) WithGroup(name string) slog.Handler {
	return this
}

func TestInfoCaller(test *testing.T) {
	handler := &captureHandler{}
	original := slog.Default()
	slog.SetDefault(slog.New(handler))
	defer slog.SetDefault(original)

	log4g.Info("Hello {}", "world")

	frame, _ := runtime.CallersFrames([]uintptr{handler.record.PC}).Next()

	require.True(test, strings.HasSuffix(frame.File, "log4g_test.go"))
	require.Equal(test, "github.com/go-log4g/core/log4g_test.TestInfoCaller", frame.Function)
	require.Equal(test, "Hello world", handler.record.Message)
}

func TestLevels(test *testing.T) {
	handler := &captureHandler{}
	original := slog.Default()
	slog.SetDefault(slog.New(handler))
	defer slog.SetDefault(original)

	log4g.Trace("Trace")
	require.Equal(test, impl.LevelTrace, handler.record.Level)
	require.Equal(test, "Trace", handler.record.Message)

	log4g.Debug("Debug")
	require.Equal(test, slog.LevelDebug, handler.record.Level)
	require.Equal(test, "Debug", handler.record.Message)

	log4g.Info("Info")
	require.Equal(test, slog.LevelInfo, handler.record.Level)
	require.Equal(test, "Info", handler.record.Message)

	log4g.Warn("Warn")
	require.Equal(test, slog.LevelWarn, handler.record.Level)
	require.Equal(test, "Warn", handler.record.Message)

	log4g.Error("Error")
	require.Equal(test, slog.LevelError, handler.record.Level)
	require.Equal(test, "Error", handler.record.Message)

	log4g.Fatal("Fatal")
	require.Equal(test, impl.LevelFatal, handler.record.Level)
	require.Equal(test, "Fatal", handler.record.Message)
}

func TestContextLevels(test *testing.T) {
	handler := &captureHandler{}
	original := slog.Default()
	slog.SetDefault(slog.New(handler))
	defer slog.SetDefault(original)

	ctx := context.Background()

	log4g.TraceContext(ctx, "Trace")
	require.Equal(test, impl.LevelTrace, handler.record.Level)
	require.Equal(test, "Trace", handler.record.Message)

	log4g.DebugContext(ctx, "Debug")
	require.Equal(test, slog.LevelDebug, handler.record.Level)
	require.Equal(test, "Debug", handler.record.Message)

	log4g.InfoContext(ctx, "Info")
	require.Equal(test, slog.LevelInfo, handler.record.Level)
	require.Equal(test, "Info", handler.record.Message)

	log4g.WarnContext(ctx, "Warn")
	require.Equal(test, slog.LevelWarn, handler.record.Level)
	require.Equal(test, "Warn", handler.record.Message)

	log4g.ErrorContext(ctx, "Error")
	require.Equal(test, slog.LevelError, handler.record.Level)
	require.Equal(test, "Error", handler.record.Message)

	log4g.FatalContext(ctx, "Fatal")
	require.Equal(test, impl.LevelFatal, handler.record.Level)
	require.Equal(test, "Fatal", handler.record.Message)
}

func TestInfoContext(test *testing.T) {
	handler := &captureHandler{}
	original := slog.Default()
	slog.SetDefault(slog.New(handler))
	defer slog.SetDefault(original)

	type contextKey string
	ctx := context.WithValue(context.Background(), contextKey("requestId"), "123")

	log4g.InfoContext(ctx, "Hello {}", "world")

	require.Same(test, ctx, handler.ctx)
	require.Equal(test, "Hello world", handler.record.Message)

	frame, _ := runtime.CallersFrames([]uintptr{handler.record.PC}).Next()

	require.True(test, strings.HasSuffix(frame.File, "log4g_test.go"))
	require.Equal(test, "github.com/go-log4g/core/log4g_test.TestInfoContext", frame.Function)
	require.Same(test, ctx, handler.ctx)
	require.Equal(test, "Hello world", handler.record.Message)
}

func TestInfoContextMdc(test *testing.T) {
	handler := &captureHandler{}
	original := slog.Default()
	slog.SetDefault(slog.New(handler))
	defer slog.SetDefault(original)

	ctx := context.Background()
	ctx = mdc.Put(ctx, "requestId", "123")

	log4g.InfoContext(ctx, "Hello")

	value, ok := mdc.Get(handler.ctx, "requestId")

	require.True(test, ok)
	require.Equal(test, "123", value)
}
