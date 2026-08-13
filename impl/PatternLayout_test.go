package impl_test

import (
	"context"
	"testing"

	"github.com/go-log4g/core/impl"
	"github.com/go-log4g/core/mdc"
	"github.com/stretchr/testify/require"
)

func TestPatternLayoutMdc(test *testing.T) {
	statusLogger := impl.NewStatusLogger()
	parser := impl.NewPatternParser(statusLogger)

	ctx := context.Background()
	ctx = mdc.Put(ctx, "requestId", "123")
	ctx = mdc.Put(ctx, "userId", "456")

	event := &impl.LogEvent{
		Context: ctx,
	}

	layout := impl.NewPatternLayout("%X{requestId} %X", parser)

	require.Equal(test, "123 {requestId=123, userId=456}", string(layout.Append(nil, event)))
}
