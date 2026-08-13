package impl

type PatternConverter interface {
	Append(result []byte, event *LogEvent) []byte
}
