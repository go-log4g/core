package impl

type PatternLayout struct {
	pattern    string
	converters []PatternConverter
}

func NewPatternLayout(pattern string, parser *PatternParser) *PatternLayout {
	return &PatternLayout{
		pattern:    pattern,
		converters: parser.Parse(pattern),
	}
}

func (this *PatternLayout) Pattern() string {
	return this.pattern
}

func (this *PatternLayout) Append(result []byte, event *LogEvent) []byte {
	for _, converter := range this.converters {
		result = converter.Append(result, event)
	}
	return result
}
