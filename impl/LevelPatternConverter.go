package impl

import "log/slog"

type LevelPatternConverter struct {
	AbstractPatternConverter
}

func NewLevelPatternConverter(formattingInfo FormattingInfo) *LevelPatternConverter {
	return &LevelPatternConverter{
		AbstractPatternConverter: NewAbstractPatternConverter(formattingInfo),
	}
}

func (this *LevelPatternConverter) Append(result []byte, event *LogEvent) []byte {
	return this.AppendFormatted(result, this.levelName(event.Record.Level))
}

func (this *LevelPatternConverter) levelName(level slog.Level) string {
	switch level {
	case LevelTrace:
		return "TRACE"
	case slog.LevelDebug:
		return "DEBUG"
	case slog.LevelInfo:
		return "INFO"
	case slog.LevelWarn:
		return "WARN"
	case slog.LevelError:
		return "ERROR"
	case LevelFatal:
		return "FATAL"
	default:
		return level.String()
	}
}
