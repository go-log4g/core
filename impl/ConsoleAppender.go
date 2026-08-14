package impl

import (
	"io"
	"sync"

	"github.com/go-log4g/core/impl/filter"
)

const maxPooledBufferCapacity = 64 * 1024

type ConsoleAppender struct {
	writer       io.Writer
	layout       Layout
	filter       filter.Filter
	statusLogger *StatusLogger

	pool sync.Pool
}

func NewConsoleAppender(writer io.Writer, layout Layout, filter filter.Filter, statusLogger *StatusLogger) *ConsoleAppender {
	result := &ConsoleAppender{
		writer:       writer,
		layout:       layout,
		filter:       filter,
		statusLogger: statusLogger,
	}
	result.pool.New = func() any {
		return make([]byte, 0, 256)
	}
	return result
}

// Implements Appender
func (this *ConsoleAppender) Append(event *LogEvent) {
	if this.filter != nil && this.filter.Filter(event.Record.Level) == filter.Deny {
		return
	}

	data := this.pool.Get().([]byte)[:0]
	defer func() {
		if cap(data) <= maxPooledBufferCapacity {
			this.pool.Put(data)
		}
	}()

	data = this.layout.Append(data, event)
	if _, e := this.writer.Write(data); e != nil {
		this.statusLogger.ErrorThrottled("console", e, "failed to append to console")
	}
}
