package impl

type LiteralPatternConverter struct {
	value string
}

func NewLiteralPatternConverter(value string) *LiteralPatternConverter {
	return &LiteralPatternConverter{
		value: value,
	}
}

func (this *LiteralPatternConverter) Append(result []byte, event *LogEvent) []byte {
	return append(result, this.value...)
}
