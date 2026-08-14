package impl

type Appender interface {
	Append(event *LogEvent)
}
