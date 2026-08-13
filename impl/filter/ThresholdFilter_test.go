package filter_test

import (
	"log/slog"
	"testing"

	"github.com/go-log4g/core/impl/filter"
	"github.com/stretchr/testify/require"
)

func TestThresholdFilter(test *testing.T) {
	value := filter.NewThresholdFilter(slog.LevelError, filter.Neutral, filter.Deny)

	require.Equal(test, filter.Deny, value.Filter(slog.LevelDebug))
	require.Equal(test, filter.Deny, value.Filter(slog.LevelInfo))
	require.Equal(test, filter.Deny, value.Filter(slog.LevelWarn))
	require.Equal(test, filter.Neutral, value.Filter(slog.LevelError))
}
