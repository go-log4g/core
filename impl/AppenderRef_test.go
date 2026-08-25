package impl_test

import (
	"context"
	"log/slog"
	"testing"
	"time"

	"github.com/go-log4g/core/impl"
	"github.com/go-log4g/core/impl/filter"
	"github.com/stretchr/testify/require"
)

func TestAppenderRefAcceptLevel(test *testing.T) {
	level := slog.LevelWarn
	appenderRef := &impl.AppenderRef{
		Level: &level,
	}

	require.False(test, appenderRef.Accept(logEvent(slog.LevelInfo)))
	require.True(test, appenderRef.Accept(logEvent(slog.LevelWarn)))
	require.True(test, appenderRef.Accept(logEvent(slog.LevelError)))
}

func TestAppenderRefAcceptFilter(test *testing.T) {
	appenderRef := &impl.AppenderRef{
		Filter: filter.NewLevelRangeFilter(slog.LevelInfo, slog.LevelWarn, filter.Neutral, filter.Deny),
	}

	require.False(test, appenderRef.Accept(logEvent(slog.LevelDebug)))
	require.True(test, appenderRef.Accept(logEvent(slog.LevelInfo)))
	require.True(test, appenderRef.Accept(logEvent(slog.LevelWarn)))
	require.False(test, appenderRef.Accept(logEvent(slog.LevelError)))
}

func logEvent(level slog.Level) *impl.LogEvent {
	return impl.NewLogEvent(context.Background(), slog.NewRecord(time.Now(), level, "test", 0), &impl.CallerContext{}, nil, nil)
}
