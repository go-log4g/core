package impl

import "strconv"

type LinePatternConverter struct {
	AbstractPatternConverter
}

func NewLinePatternConverter(formatting FormattingInfo) *LinePatternConverter {
	return &LinePatternConverter{
		AbstractPatternConverter: NewAbstractPatternConverter(formatting),
	}
}

func (this *LinePatternConverter) Append(result []byte, event *LogEvent) []byte {
	start := len(result)
	result = strconv.AppendInt(result, int64(event.CallerContext.Caller.Line), 10)
	return this.Format(result, start)
}
