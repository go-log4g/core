package impl_test

import (
	"context"
	"log/slog"
	"testing"
	"time"

	"github.com/go-log4g/core/impl"
	"github.com/stretchr/testify/require"
)

func TestPatternParserFormatting(test *testing.T) {
	statusLogger := impl.NewStatusLogger()
	parser := impl.NewPatternParser(statusLogger)

	record := slog.NewRecord(time.Now(), slog.LevelInfo, "", 0)
	event := impl.NewLogEvent(context.Background(), record, nil, nil, nil)

	require.Equal(test, "INFO", string(impl.NewPatternLayout("%p", parser).Append(nil, event)))
	require.Equal(test, " INFO", string(impl.NewPatternLayout("%5p", parser).Append(nil, event)))
	require.Equal(test, "INFO ", string(impl.NewPatternLayout("%-5p", parser).Append(nil, event)))
	require.Equal(test, "0INFO", string(impl.NewPatternLayout("%05p", parser).Append(nil, event)))
	require.Equal(test, "NFO", string(impl.NewPatternLayout("%.3p", parser).Append(nil, event)))
	require.Equal(test, "INF", string(impl.NewPatternLayout("%.-3p", parser).Append(nil, event)))
	require.Equal(test, "  NFO", string(impl.NewPatternLayout("%5.3p", parser).Append(nil, event)))
}

func TestPatternParserLoggerPrecision(test *testing.T) {
	statusLogger := impl.NewStatusLogger()
	parser := impl.NewPatternParser(statusLogger)

	caller := &impl.Caller{
		Logger: "playground/internal/app/Service1",
	}

	callerContext := &impl.CallerContext{
		Caller: caller,
	}

	event := &impl.LogEvent{
		CallerContext: callerContext,
	}

	require.Equal(test, "playground/internal/app/Service1", string(impl.NewPatternLayout("%c", parser).Append(nil, event)))
	require.Equal(test, "Service1", string(impl.NewPatternLayout("%c{0}", parser).Append(nil, event)))
	require.Equal(test, "Service1", string(impl.NewPatternLayout("%c{1}", parser).Append(nil, event)))
	require.Equal(test, "app/Service1", string(impl.NewPatternLayout("%c{2}", parser).Append(nil, event)))
	require.Equal(test, "internal/app/Service1", string(impl.NewPatternLayout("%c{3}", parser).Append(nil, event)))
	require.Equal(test, "internal/app/Service1", string(impl.NewPatternLayout("%c{-1}", parser).Append(nil, event)))
	require.Equal(test, "app/Service1", string(impl.NewPatternLayout("%c{-2}", parser).Append(nil, event)))
	require.Equal(test, "p/i/a/Service1", string(impl.NewPatternLayout("%c{1.}", parser).Append(nil, event)))
	require.Equal(test, "pl/in/ap/Service1", string(impl.NewPatternLayout("%c{2.}", parser).Append(nil, event)))
}

func TestPatternParserAliases(test *testing.T) {
	statusLogger := impl.NewStatusLogger()
	parser := impl.NewPatternParser(statusLogger)

	caller := &impl.Caller{
		Logger: "playground/internal/app/Service1",
		File:   "Service1.go",
		Line:   123,
		Method: "AfterPropertiesSet",
	}

	event := &impl.LogEvent{
		CallerContext: &impl.CallerContext{
			Caller: caller,
		},
	}

	require.Equal(test, "playground/internal/app/Service1", string(impl.NewPatternLayout("%logger", parser).Append(nil, event)))
	require.Equal(test, "Service1", string(impl.NewPatternLayout("%logger{1}", parser).Append(nil, event)))

	require.Equal(test, "Service1.go", string(impl.NewPatternLayout("%file", parser).Append(nil, event)))
	require.Equal(test, "123", string(impl.NewPatternLayout("%line", parser).Append(nil, event)))
	require.Equal(test, "AfterPropertiesSet", string(impl.NewPatternLayout("%method", parser).Append(nil, event)))
	require.Equal(test, "p/i/app/Service1", string(impl.NewPatternLayout("%c{1.2.*}", parser).Append(nil, event)))
	require.Equal(test, "p/internal/app/Service1", string(impl.NewPatternLayout("%c{1.3.*}", parser).Append(nil, event)))
	require.Equal(test, "p/i/a/Service1", string(impl.NewPatternLayout("%c{1.1.1.*}", parser).Append(nil, event)))
}
