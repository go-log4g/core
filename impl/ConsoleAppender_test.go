package impl_test

import (
	"bytes"
	"log/slog"
	"testing"
	"time"

	"github.com/go-log4g/core/impl"
	"github.com/go-log4g/core/impl/filter"
	"github.com/stretchr/testify/require"
)

func TestConsoleAppenderFilter(test *testing.T) {
	buffer := &bytes.Buffer{}
	parser := impl.NewPatternParser(impl.NewStatusLogger())
	layout := impl.NewPatternLayout("%m%n", parser)

	appender := impl.NewConsoleAppender(
		buffer,
		layout,
		filter.NewThresholdFilter(slog.LevelError, filter.Neutral, filter.Deny),
	)

	info := slog.NewRecord(time.Now(), slog.LevelInfo, "Info", 0)
	errorRecord := slog.NewRecord(time.Now(), slog.LevelError, "Error", 0)

	require.NoError(test, appender.Append(&impl.LogEvent{Record: info}))
	require.Empty(test, buffer.String())

	require.NoError(test, appender.Append(&impl.LogEvent{Record: errorRecord}))
	require.Equal(test, "Error\n", buffer.String())
}
