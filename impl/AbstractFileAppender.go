package impl

import (
	"sync"
	"time"

	"github.com/go-errr/go/err"
	"github.com/go-jang/go/util/concurrent"
	"github.com/go-log4g/core/impl/filter"
)

const fileFlushInterval = 5 * time.Second

type fileAppenderDelegate interface {
	beforeAppend(data []byte)
	write(data []byte)
}

type AbstractFileAppender struct {
	file           string
	bufferSize     int
	immediateFlush bool
	layout         Layout
	filter         filter.Filter
	statusLogger   *StatusLogger
	delegate       fileAppenderDelegate

	writeLock *concurrent.FairLock
	mutex     sync.Mutex
	cond      *sync.Cond
	buffers   [2][]byte
	active    int
	flushing  int
	lastFlush time.Time
	pool      sync.Pool
}

func newAbstractFileAppender(file string, bufferSize int, immediateFlush bool, layout Layout, filter filter.Filter, statusLogger *StatusLogger, delegate fileAppenderDelegate) *AbstractFileAppender {
	result := &AbstractFileAppender{
		file:           file,
		bufferSize:     bufferSize,
		immediateFlush: immediateFlush,
		layout:         layout,
		filter:         filter,
		statusLogger:   statusLogger,
		delegate:       delegate,
		writeLock:      concurrent.NewFairLock(),
		flushing:       -1,
		lastFlush:      time.Now(),
	}

	result.pool.New = func() any {
		return make([]byte, 0, 256)
	}

	if !immediateFlush {
		result.buffers[0] = make([]byte, 0, bufferSize)
		result.buffers[1] = make([]byte, 0, bufferSize)
		result.cond = sync.NewCond(&result.mutex)
		go result.runWriter()
		go result.runPeriodicFlush()
	}

	return result
}

// Implements Appender
func (this *AbstractFileAppender) Append(event *LogEvent) {
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

	this.writeLock.Lock()
	defer this.writeLock.Unlock()

	this.delegate.beforeAppend(data)

	if this.immediateFlush {
		this.delegate.write(data)
	} else {
		this.writeBuffered(data)
	}
}

func (this *AbstractFileAppender) writeBuffered(data []byte) {
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

func (this *AbstractFileAppender) flushBuffered() {
	if this.immediateFlush {
		return
	}

	this.mutex.Lock()
	defer this.mutex.Unlock()

	for this.flushing != -1 {
		this.cond.Wait()
	}

	if len(this.buffers[this.active]) == 0 {
		return
	}

	this.swapBuffers()

	for this.flushing != -1 {
		this.cond.Wait()
	}
}

func (this *AbstractFileAppender) swapBuffers() {
	this.flushing = this.active
	this.active = 1 - this.active
	this.cond.Broadcast()
}

func (this *AbstractFileAppender) runWriter() {
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

func (this *AbstractFileAppender) writeBuffer(data []byte) {
	defer err.Recover(func(e any) {
		this.statusLogger.ErrorThrottled(this.file, e, "failed to flush log file %q", this.file)
	})

	this.delegate.write(data)
}

func (this *AbstractFileAppender) runPeriodicFlush() {
	ticker := time.NewTicker(fileFlushInterval)
	for range ticker.C {
		this.flushByTime()
	}
}

func (this *AbstractFileAppender) flushByTime() {
	this.mutex.Lock()
	defer this.mutex.Unlock()

	if this.flushing != -1 || len(this.buffers[this.active]) == 0 {
		return
	}

	if time.Since(this.lastFlush) >= fileFlushInterval {
		this.swapBuffers()
	}
}
