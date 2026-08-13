package impl

type MessagePatternConverter struct {
	AbstractPatternConverter
	formatter *SlogMessageFormatter
}

func NewMessagePatternConverter(formatting FormattingInfo) *MessagePatternConverter {
	return &MessagePatternConverter{
		AbstractPatternConverter: NewAbstractPatternConverter(formatting),
		formatter:                NewSlogMessageFormatter(),
	}
}

func (this *MessagePatternConverter) Append(result []byte, event *LogEvent) []byte {
	start := len(result)
	result = this.formatter.Append(result, event)
	return this.Format(result, start)
}
