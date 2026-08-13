package impl

import (
	"io"
	"sync"
)

const maxPooledBufferCapacity = 64 * 1024

type ConsoleAppender struct {
	writer io.Writer
	layout Layout
	pool   sync.Pool
}

func NewConsoleAppender(writer io.Writer, layout Layout) *ConsoleAppender {
	result := &ConsoleAppender{
		writer: writer,
		layout: layout,
	}
	result.pool.New = func() any {
		return make([]byte, 0, 256)
	}
	return result
}

func (this *ConsoleAppender) Append(event *LogEvent) error {
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
