package impl_test

import (
	"testing"

	"github.com/go-log4g/core/impl"
	"github.com/go-log4g/core/impl/abbr"
	"github.com/stretchr/testify/require"
)

func TestLoggerPatternConverterPrecision(test *testing.T) {
	caller := &impl.Caller{
		Logger: "playground/internal/app/Service1",
	}

	callerContext := &impl.CallerContext{
		Caller: caller,
	}

	event := &impl.LogEvent{
		CallerContext: callerContext,
	}

	formatting := impl.NewDefaultFormattingInfo()

	converter := impl.NewLoggerPatternConverter(formatting, abbr.NewNOPAbbreviator())
	require.Equal(test, "playground/internal/app/Service1", string(converter.Append(nil, event)))

	converter = impl.NewLoggerPatternConverter(formatting, impl.NewMaxElementAbbreviator(0))
	require.Equal(test, "Service1", string(converter.Append(nil, event)))

	converter = impl.NewLoggerPatternConverter(formatting, impl.NewMaxElementAbbreviator(1))
	require.Equal(test, "Service1", string(converter.Append(nil, event)))

	converter = impl.NewLoggerPatternConverter(formatting, impl.NewMaxElementAbbreviator(2))
	require.Equal(test, "app/Service1", string(converter.Append(nil, event)))

	converter = impl.NewLoggerPatternConverter(formatting, impl.NewMaxElementAbbreviator(3))
	require.Equal(test, "internal/app/Service1", string(converter.Append(nil, event)))

	converter = impl.NewLoggerPatternConverter(formatting, impl.NewMaxElementAbbreviator(-1))
	require.Equal(test, "internal/app/Service1", string(converter.Append(nil, event)))

	converter = impl.NewLoggerPatternConverter(formatting, impl.NewMaxElementAbbreviator(-2))
	require.Equal(test, "app/Service1", string(converter.Append(nil, event)))

	converter = impl.NewLoggerPatternConverter(formatting, abbr.NewPatternAbbreviator(
		abbr.NewPatternAbbreviatorFragment(1),
	))
	require.Equal(test, "p/i/a/Service1", string(converter.Append(nil, event)))

	converter = impl.NewLoggerPatternConverter(formatting, abbr.NewPatternAbbreviator(
		abbr.NewPatternAbbreviatorFragment(2),
	))
	require.Equal(test, "pl/in/ap/Service1", string(converter.Append(nil, event)))
}
