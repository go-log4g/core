package impl_test

import (
	"testing"

	"github.com/go-log4g/core/impl"
	"github.com/stretchr/testify/require"
)

func TestAbstractPatternConverterFormatting(test *testing.T) {
	converter := impl.NewAbstractPatternConverter(impl.NewDefaultFormattingInfo())

	require.Equal(test, "INFO", string(converter.AppendFormatted(nil, "INFO")))

	converter.FormattingInfo.MinLength = 6
	require.Equal(test, "  INFO", string(converter.AppendFormatted(nil, "INFO")))

	converter.FormattingInfo.LeftAlign = true
	require.Equal(test, "INFO  ", string(converter.AppendFormatted(nil, "INFO")))

	converter.FormattingInfo.MaxLength = 3
	converter.FormattingInfo.LeftTruncate = true
	require.Equal(test, "NFO   ", string(converter.AppendFormatted(nil, "INFO")))

	converter.FormattingInfo.LeftTruncate = false
	require.Equal(test, "INF   ", string(converter.AppendFormatted(nil, "INFO")))

	converter.FormattingInfo = impl.NewDefaultFormattingInfo()
	converter.FormattingInfo.MinLength = 6
	converter.FormattingInfo.ZeroPad = true
	require.Equal(test, "00INFO", string(converter.AppendFormatted(nil, "INFO")))
}
