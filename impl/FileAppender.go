package impl

import (
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/go-errr/go/err"
	"github.com/go-jang/go/lang"
	"github.com/go-jang/go/util/concurrent"
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

	handle *os.File

	writeLock *concurrent.FairLock
	mutex     sync.Mutex
	cond      *sync.Cond

	buffers   [2][]byte
	active    int
	flushing  int
	lastFlush time.Time
	pool      sync.Pool
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
		writeLock:      concurrent.NewFairLock(),
		flushing:       -1,
		lastFlush:      time.Now(),
	}
	appender.pool.New = func() any {
		return make([]byte, 0, 256)
	}

	if !immediateFlush {
		appender.buffers[0] = make([]byte, 0, bufferSize)
		appender.buffers[1] = make([]byte, 0, bufferSize)
		appender.cond = sync.NewCond(&appender.mutex)
		go appender.runWriter()
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

	if this.immediateFlush {
		this.writeImmediate(data)
	} else {
		this.writeBuffered(data)
	}
}

func (this *FileAppender) writeImmediate(data []byte) {
	this.writeLock.Lock()
	defer this.writeLock.Unlock()

	written := optional.OfCommaErr(this.handle.Write(data)).OrElsePanic("failed to write data to log file %q", this.file)
	lang.Assert(written == len(data), "failed to write complete log event to file %q", this.file)
}

func (this *FileAppender) writeBuffered(data []byte) {
	this.writeLock.Lock()
	defer this.writeLock.Unlock()

	this.mutex.Lock()
	defer this.mutex.Unlock()

	for len(this.buffers[this.active])+len(data) > this.bufferSize && len(this.buffers[this.active]) > 0 {
		if this.flushing == -1 {
			this.swapBuffers()
		} else {
			this.cond.Wait()
		}
	}

	this.buffers[this.active] = append(this.buffers[this.active], data...)

	if len(this.buffers[this.active]) >= this.bufferSize {
		for this.flushing != -1 {
			this.cond.Wait()
		}
		this.swapBuffers()
	}
}

func (this *FileAppender) swapBuffers() {
	this.flushing = this.active
	this.active = 1 - this.active
	this.cond.Broadcast()
}

func (this *FileAppender) runWriter() {
	for {
		this.mutex.Lock()
		for this.flushing == -1 {
			this.cond.Wait()
		}

		index := this.flushing
		data := this.buffers[index]
		this.mutex.Unlock()

		this.writeBuffer(data)

		this.mutex.Lock()
		this.buffers[index] = this.buffers[index][:0]
		this.flushing = -1
		this.lastFlush = time.Now()
		this.cond.Broadcast()
		this.mutex.Unlock()
	}
}

func (this *FileAppender) writeBuffer(data []byte) {
	defer err.Recover(func(e any) {
		this.statusLogger.ErrorThrottled(this.file, e, "failed to flush log file %q", this.file)
	})

	written := optional.OfCommaErr(this.handle.Write(data)).OrElsePanic("failed to write data to log file %q", this.file)
	lang.Assert(written == len(data), "failed to write complete buffer to log file %q", this.file)
}

func (this *FileAppender) runPeriodicFlush() {
	ticker := time.NewTicker(fileFlushInterval)

	for range ticker.C {
		this.flushByTime()
	}
}

func (this *FileAppender) flushByTime() {
	this.mutex.Lock()
	defer this.mutex.Unlock()

	if this.flushing != -1 || len(this.buffers[this.active]) == 0 {
		return
	}

	if time.Since(this.lastFlush) >= fileFlushInterval {
		this.swapBuffers()
	}
}
