package impl

type Layout interface {
	Append(result []byte, event *LogEvent) []byte
}
