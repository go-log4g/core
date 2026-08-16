package impl

import (
	"os"
	"path/filepath"

	"github.com/go-jang/go/lang"
	"github.com/go-jang/go/util/optional"
	"github.com/go-log4g/core/impl/filter"
)

type FileAppender struct {
	*AbstractFileAppender

	handle *os.File
}

func NewFileAppender(file string, append bool, bufferSize int, immediateFlush bool, layout Layout, filter filter.Filter, statusLogger *StatusLogger) *FileAppender {
	handle := openLogFile(file, append)

	result := &FileAppender{
		handle: handle,
	}
	result.AbstractFileAppender = newAbstractFileAppender(file, bufferSize, immediateFlush, layout, filter, statusLogger, result)

	return result
}

// Implements fileAppenderDelegate
func (this *FileAppender) beforeAppend(data []byte) {
}

// Implements fileAppenderDelegate
func (this *FileAppender) write(data []byte) {
	written := optional.OfCommaErr(this.handle.Write(data)).OrElsePanic("failed to write data to log file %q", this.file)
	lang.Assert(written == len(data), "failed to write complete data to log file %q", this.file)
}

func openLogFile(file string, append bool) *os.File {
	dir := filepath.Dir(file)
	if dir != "." {
		lang.Assert(os.MkdirAll(dir, 0755) == nil, "failed to create log directory %q", dir)
	}

	flags := os.O_CREATE | os.O_WRONLY
	if append {
		flags |= os.O_APPEND
	} else {
		flags |= os.O_TRUNC
	}

	return optional.OfCommaErr(os.OpenFile(file, flags, 0644)).OrElsePanic("failed to open log file %q", file)
}
