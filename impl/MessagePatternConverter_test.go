package impl_test

import (
	"log/slog"
	"testing"
	"time"

	"github.com/go-log4g/core/impl"
	"github.com/stretchr/testify/require"
)

func TestMessagePatternConverterAttrs(test *testing.T) {
	record := slog.NewRecord(time.Now(), slog.LevelInfo, "User authenticated", 0)
	record.Add("userId", 123, "remote", "10.0.0.1")

	event := &impl.LogEvent{
		Record: record,
	}

	converter := impl.NewMessagePatternConverter(impl.NewDefaultFormattingInfo())

	require.Equal(test, "User authenticated userId=123 remote=10.0.0.1", string(converter.Append(nil, event)))
}
