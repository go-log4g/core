package impl

type LoggerPatternConverter struct {
	AbstractPatternConverter
	abbreviator NameAbbreviator
}

func NewLoggerPatternConverter(formatting FormattingInfo, abbreviator NameAbbreviator) *LoggerPatternConverter {
	return &LoggerPatternConverter{
		AbstractPatternConverter: NewAbstractPatternConverter(formatting),
		abbreviator:              abbreviator,
	}
}

func (this *LoggerPatternConverter) Append(result []byte, event *LogEvent) []byte {
	value := this.abbreviator.Abbreviate(event.CallerContext.Caller.Logger)
	return this.AppendFormatted(result, value)
}
