package impl

type LineSeparatorPatternConverter struct {
	AbstractPatternConverter
}

func NewLineSeparatorPatternConverter(formatting FormattingInfo) *LineSeparatorPatternConverter {
	return &LineSeparatorPatternConverter{
		AbstractPatternConverter: NewAbstractPatternConverter(formatting),
	}
}

func (this *LineSeparatorPatternConverter) Append(result []byte, event *LogEvent) []byte {
	return this.AppendFormatted(result, "\n")
}
