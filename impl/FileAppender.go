package impl

import (
	"bufio"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/go-errr/go/err"
	"github.com/go-jang/go/lang"
	"github.com/go-jang/go/util/optional"
	"github.com/go-log4g/core/impl/filter"
)

const fileFlushInterval = 5 * time.Second

type FileAppender struct {
	file           string
	append         bool
	bufferSize     int
	immediateFlush bool
	layout         Layout
	filter         filter.Filter
	statusLogger   *StatusLogger

	writer *bufio.Writer
	handle *os.File
	mutex  sync.Mutex
	dirty  bool
	pool   sync.Pool
}

func NewFileAppender(file string, append bool, bufferSize int, immediateFlush bool, layout Layout, filter filter.Filter, statusLogger *StatusLogger) *FileAppender {
	flags := os.O_CREATE | os.O_WRONLY
	if append {
		flags |= os.O_APPEND
	} else {
		flags |= os.O_TRUNC
	}

	dir := filepath.Dir(file)
	if dir != "." {
		lang.Assert(os.MkdirAll(dir, 0755) == nil, "failed to create log directory %q", dir)
	}

	handle := optional.OfCommaErr(os.OpenFile(file, flags, 0644)).OrElsePanic("Failed to open log file %q", file)

	appender := &FileAppender{
		file:           file,
		append:         append,
		bufferSize:     bufferSize,
		immediateFlush: immediateFlush,
		layout:         layout,
		filter:         filter,
		statusLogger:   statusLogger,
		handle:         handle,
		writer:         bufio.NewWriterSize(handle, bufferSize),
	}
	appender.pool.New = func() any {
		return make([]byte, 0, 256)
	}
	if !immediateFlush {
		go appender.runPeriodicFlush()
	}
	return appender
}

// Implements Appender
func (this *FileAppender) Append(event *LogEvent) {
	if this.filter != nil && this.filter.Filter(event.Record.Level) == filter.Deny {
		return
	}
	defer err.Recover(func(e any) {
		this.statusLogger.ErrorThrottled(this.file, e, "failed to append to log file %q", this.file)
	})

	data := this.pool.Get().([]byte)[:0]
	defer func() {
		if cap(data) <= maxPooledBufferCapacity {
			this.pool.Put(data)
		}
	}()

	data = this.layout.Append(data, event)

	this.mutex.Lock()
	defer this.mutex.Unlock()

	optional.OfCommaErr(this.writer.Write(data)).OrElsePanic("failed to write data to log file %q", this.file)
	if this.immediateFlush {
		lang.Assert(this.writer.Flush() == nil, "failed to flush log file %q", this.file)
	} else {
		this.dirty = true
	}
}

func (this *FileAppender) runPeriodicFlush() {
	ticker := time.NewTicker(fileFlushInterval)
	for range ticker.C {
		this.flushIfDirty()
	}
}

func (this *FileAppender) flushIfDirty() {
	defer err.Recover(func(e any) {
		this.statusLogger.ErrorThrottled(this.file, e, "Failed to flush log file %q", this.file)
	})

	this.mutex.Lock()
	defer this.mutex.Unlock()

	if this.dirty {
		lang.Assert(this.writer.Flush() == nil, "failed to flush log file %q", this.file)
		this.dirty = false
	}
}
