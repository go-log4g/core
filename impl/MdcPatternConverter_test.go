package impl_test

import (
	"context"
	"testing"

	"github.com/go-log4g/core/impl"
	"github.com/go-log4g/core/mdc"
	"github.com/stretchr/testify/require"
)

func TestMdcPatternConverter(test *testing.T) {
	ctx := context.Background()
	ctx = mdc.Put(ctx, "requestId", "123")
	ctx = mdc.Put(ctx, "userId", "456")

	event := &impl.LogEvent{
		Context: ctx,
	}

	formatting := impl.NewDefaultFormattingInfo()

	converter := impl.NewMdcPatternConverter(formatting, "requestId")
	require.Equal(test, "123", string(converter.Append(nil, event)))

	converter = impl.NewMdcPatternConverter(formatting, "")
	require.Equal(test, "{requestId=123, userId=456}", string(converter.Append(nil, event)))
}
