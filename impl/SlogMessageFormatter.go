package impl

import (
	"fmt"
	"log/slog"
	"strconv"
	"time"
)

type SlogMessageFormatter struct {
}

func NewSlogMessageFormatter() *SlogMessageFormatter {
	return &SlogMessageFormatter{}
}

func (this *SlogMessageFormatter) Append(result []byte, event *LogEvent) []byte {
	result = append(result, event.Record.Message...)

	for _, handlerAttr := range event.Attrs {
		result = this.appendAttr(result, handlerAttr.Groups, handlerAttr.Attr)
	}

	event.Record.Attrs(func(attr slog.Attr) bool {
		result = this.appendAttr(result, event.Groups, attr)
		return true
	})

	return result
}

func (this *SlogMessageFormatter) appendAttr(result []byte, groups []string, attr slog.Attr) []byte {
	attr.Value = attr.Value.Resolve()
	if attr.Equal(slog.Attr{}) {
		return result
	}
	if attr.Value.Kind() == slog.KindGroup {
		return this.appendGroup(result, groups, attr)
	}

	result = append(result, ' ')
	for _, group := range groups {
		result = append(result, group...)
		result = append(result, '.')
	}
	result = append(result, attr.Key...)
	result = append(result, '=')
	return this.appendValue(result, attr.Value)
}

func (this *SlogMessageFormatter) appendGroup(result []byte, groups []string, attr slog.Attr) []byte {
	if attr.Key != "" {
		groups = append(groups, attr.Key)
	}
	for _, child := range attr.Value.Group() {
		result = this.appendAttr(result, groups, child)
	}
	return result
}

func (this *SlogMessageFormatter) appendValue(result []byte, value slog.Value) []byte {
	switch value.Kind() {
	case slog.KindString:
		return append(result, value.String()...)
	case slog.KindInt64:
		return strconv.AppendInt(result, value.Int64(), 10)
	case slog.KindUint64:
		return strconv.AppendUint(result, value.Uint64(), 10)
	case slog.KindFloat64:
		return strconv.AppendFloat(result, value.Float64(), 'g', -1, 64)
	case slog.KindBool:
		return strconv.AppendBool(result, value.Bool())
	case slog.KindDuration:
		return append(result, value.Duration().String()...)
	case slog.KindTime:
		return value.Time().AppendFormat(result, time.RFC3339Nano)
	default:
		return fmt.Append(result, value.Any())
	}
}
