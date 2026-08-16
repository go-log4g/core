package impl

import (
	"os"
	"time"

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
}

func NewRollingFileAppender(file, filePattern string, append bool, bufferSize int, immediateFlush bool, layout Layout, filter filter.Filter, statusLogger *StatusLogger, policy rolling.TriggeringPolicy, strategy rolling.RolloverStrategy) *RollingFileAppender {
	handle := openLogFile(file, append)
	info := optional.OfCommaErr(handle.Stat()).OrElsePanic("failed to stat log file %q", file)

	result := &RollingFileAppender{
		filePattern: filePattern,
		handle:      handle,
		policy:      policy,
		strategy:    strategy,
		size:        info.Size(),
		fileTime:    lang.If(info.Size() > 0, info.ModTime(), time.Now()),
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

	if this.size > 0 && this.policy.IsTriggered(context) {
		this.flushBuffered()
		this.rollover(this.fileTime)
		this.fileTime = now
	}

	this.size += int64(len(data))
}

// Implements fileAppenderDelegate
func (this *RollingFileAppender) write(data []byte) {
	written := optional.OfCommaErr(this.handle.Write(data)).OrElsePanic("failed to write data to log file %q", this.file)
	lang.Assert(written == len(data), "failed to write complete data to log file %q", this.file)
}

func (this *RollingFileAppender) rollover(at time.Time) {
	lang.Assert(this.handle.Close() == nil, "failed to close log file %q", this.file)

	this.strategy.Rollover(this.file, this.filePattern, at)

	this.handle = openLogFile(this.file, false)
	this.size = 0
}
