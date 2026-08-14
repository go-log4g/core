package impl_test

import (
	"log/slog"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/go-log4g/core/impl"
	"github.com/stretchr/testify/require"
)

func TestFileAppender(test *testing.T) {
	file := filepath.Join(os.TempDir(), "go-log4g-TestFileAppender.log")
	_ = os.Remove(file)

	parser := impl.NewPatternParser(impl.NewStatusLogger())
	layout := impl.NewPatternLayout("%m%n", parser)
	appender := impl.NewFileAppender(file, true, 8192, true, layout, nil, nil)

	record := slog.NewRecord(time.Now(), slog.LevelInfo, "Hello", 0)
	appender.Append(&impl.LogEvent{Record: record})

	data, e := os.ReadFile(file)
	require.NoError(test, e)
	require.Equal(test, "Hello\n", string(data))
}

func TestFileAppenderAppend(test *testing.T) {
	file := filepath.Join(os.TempDir(), "go-log4g-TestFileAppenderAppend.log")
	_ = os.Remove(file)

	require.NoError(test, os.WriteFile(file, []byte("Before\n"), 0644))

	parser := impl.NewPatternParser(impl.NewStatusLogger())
	layout := impl.NewPatternLayout("%m%n", parser)
	appender := impl.NewFileAppender(file, true, 8192, true, layout, nil, nil)

	record := slog.NewRecord(time.Now(), slog.LevelInfo, "After", 0)
	appender.Append(&impl.LogEvent{Record: record})

	data, e := os.ReadFile(file)
	require.NoError(test, e)
	require.Equal(test, "Before\nAfter\n", string(data))
}
