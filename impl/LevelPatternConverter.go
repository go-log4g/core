package impl

type LevelPatternConverter struct {
	AbstractPatternConverter
}

func NewLevelPatternConverter(formattingInfo FormattingInfo) *LevelPatternConverter {
	return &LevelPatternConverter{
		AbstractPatternConverter: NewAbstractPatternConverter(formattingInfo),
	}
}

func (this *LevelPatternConverter) Append(result []byte, event *LogEvent) []byte {
	return this.AppendFormatted(result, event.Record.Level.String())
}
