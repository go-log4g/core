package impl

import (
	"io"
	"sync"

	"github.com/go-log4g/core/impl/filter"
)

const maxPooledBufferCapacity = 64 * 1024

type ConsoleAppender struct {
	writer io.Writer
	layout Layout
	filter filter.Filter
	pool   sync.Pool
}

func NewConsoleAppender(writer io.Writer, layout Layout, filter filter.Filter) *ConsoleAppender {
	result := &ConsoleAppender{
		writer: writer,
		layout: layout,
		filter: filter,
	}
	result.pool.New = func() any {
		return make([]byte, 0, 256)
	}
	return result
}

func (this *ConsoleAppender) Append(event *LogEvent) error {
	if this.filter != nil && this.filter.Filter(event.Record.Level) == filter.Deny {
		return nil
	}

	data := this.pool.Get().([]byte)[:0]
	defer func() {
		if cap(data) <= maxPooledBufferCapacity {
			this.pool.Put(data)
		}
	}()

	data = this.layout.Append(data, event)
	_, e := this.writer.Write(data)
	return e
}
