package impl

import "github.com/go-log4g/core/impl/abbr"

type LoggerPatternConverter struct {
	AbstractPatternConverter
	abbreviator abbr.NameAbbreviator
}

func NewLoggerPatternConverter(formatting FormattingInfo, abbreviator abbr.NameAbbreviator) *LoggerPatternConverter {
	return &LoggerPatternConverter{
		AbstractPatternConverter: NewAbstractPatternConverter(formatting),
		abbreviator:              abbreviator,
	}
}

func (this *LoggerPatternConverter) Append(result []byte, event *LogEvent) []byte {
	value := this.abbreviator.Abbreviate(event.CallerContext.Caller.Logger)
	return this.AppendFormatted(result, value)
}
