package impl

type MethodPatternConverter struct {
	AbstractPatternConverter
}

func NewMethodPatternConverter(formatting FormattingInfo) *MethodPatternConverter {
	return &MethodPatternConverter{
		AbstractPatternConverter: NewAbstractPatternConverter(formatting),
	}
}

func (this *MethodPatternConverter) Append(result []byte, event *LogEvent) []byte {
	return this.AppendFormatted(result, event.CallerContext.Caller.Method)
}
