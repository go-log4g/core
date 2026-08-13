package impl_test

import (
	"log/slog"
	"testing"
	"time"

	"github.com/go-log4g/core/impl"
	"github.com/stretchr/testify/require"
)

func TestSlogMessageFormatter(test *testing.T) {
	record := slog.NewRecord(time.Now(), slog.LevelInfo, "User authenticated", 0)
	record.Add("userId", 123, "remote", "10.0.0.1")

	event := &impl.LogEvent{
		Record: record,
	}

	formatter := impl.NewSlogMessageFormatter()

	require.Equal(test, "User authenticated userId=123 remote=10.0.0.1", string(formatter.Append(nil, event)))
}

func TestSlogMessageFormatterGroups(test *testing.T) {
	record := slog.NewRecord(time.Now(), slog.LevelInfo, "Request", 0)
	record.Add("status", 200)

	event := &impl.LogEvent{
		Record: record,
		Attrs: []*impl.HandlerAttr{
			impl.NewHandlerAttr(slog.String("method", "GET"), []string{"http"}),
		},
		Groups: []string{"http"},
	}

	formatter := impl.NewSlogMessageFormatter()

	require.Equal(test, "Request http.method=GET http.status=200", string(formatter.Append(nil, event)))
}

func TestSlogMessageFormatterPreservesAttrGroups(test *testing.T) {
	record := slog.NewRecord(time.Now(), slog.LevelInfo, "Request", 0)
	record.Add("status", 200)

	event := &impl.LogEvent{
		Record: record,
		Attrs: []*impl.HandlerAttr{
			impl.NewHandlerAttr(slog.String("service", "auth"), nil),
			impl.NewHandlerAttr(slog.Int("id", 123), []string{"request"}),
		},
		Groups: []string{"request"},
	}

	formatter := impl.NewSlogMessageFormatter()

	require.Equal(test, "Request service=auth request.id=123 request.status=200", string(formatter.Append(nil, event)))
}
