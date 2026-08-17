package impl

import (
	"os"
	"sync/atomic"
	"time"

	"github.com/go-errr/go/err"
	"github.com/go-jang/go/lang"
	"github.com/go-jang/go/util/optional"
	"github.com/go-log4g/core/impl/filter"
	"github.com/go-log4g/core/impl/rolling"
)

type RollingFileAppender struct {
	*AbstractFileAppender

	filePattern string
	fileTime    time.Time
	handle      *os.File
	policy      rolling.TriggeringPolicy
	strategy    rolling.RolloverStrategy
	size        int64
	failed      atomic.Bool
}

func NewRollingFileAppender(file string, filePattern string, append bool, bufferSize int, immediateFlush bool, layout Layout, filter filter.Filter, statusLogger *StatusLogger,
	startupPolicy *rolling.OnStartupTriggeringPolicy, triggerPolicy rolling.TriggeringPolicy, strategy rolling.RolloverStrategy) *RollingFileAppender {

	fileTime, size := startupRollover(file, filePattern, startupPolicy, strategy)
	handle := openLogFile(file, append)

	result := &RollingFileAppender{
		filePattern: filePattern,
		fileTime:    fileTime,
		handle:      handle,
		policy:      triggerPolicy,
		strategy:    strategy,
		size:        size,
	}
	result.AbstractFileAppender = newAbstractFileAppender(file, bufferSize, immediateFlush, layout, filter, statusLogger, result)

	return result
}

// Implements fileAppenderDelegate
func (this *RollingFileAppender) beforeAppend(data []byte) {
	now := time.Now()
	context := rolling.TriggeringContext{
		FileSize:  this.size,
		EventSize: len(data),
		Time:      now,
	}

	if this.size > 0 && this.policy != nil && this.policy.IsTriggered(context) {
		this.flushBuffered()
		this.rollover(this.fileTime)
	}

	this.size += int64(len(data))
}

// Implements fileAppenderDelegate
func (this *RollingFileAppender) write(data []byte) {
	if this.failed.Load() {
		return
	}
	written := optional.OfCommaErr(this.handle.Write(data)).OrElsePanic("failed to write data to log file %q", this.file)
	lang.Assert(written == len(data), "failed to write complete data to log file %q", this.file)
}

func startupRollover(file string, filePattern string, startupPolicy *rolling.OnStartupTriggeringPolicy, strategy rolling.RolloverStrategy) (time.Time, int64) {
	fileTime := time.Now()

	info, e := os.Stat(file)
	if os.IsNotExist(e) {
		return fileTime, 0
	}
	lang.Assert(e == nil, "failed to stat log file %q", file)

	if startupPolicy != nil && startupPolicy.IsTriggered(info.Size()) {
		strategy.Rollover(file, filePattern, info.ModTime())
		return time.Now(), 0
	}

	return info.ModTime(), info.Size()
}

func (this *RollingFileAppender) rollover(at time.Time) {
	defer err.Recover(func(e any) {
		this.failed.Store(true)
		panic(e)
	})

	lang.Assert(this.handle.Close() == nil, "failed to close log file %q", this.file)

	this.strategy.Rollover(this.file, this.filePattern, at)

	this.handle = openLogFile(this.file, false)
	this.size = 0
	this.fileTime = time.Now()
}
