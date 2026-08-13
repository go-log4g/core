package impl

import (
	"fmt"
	"os"
)

type StatusLogger struct {
}

func NewStatusLogger() *StatusLogger {
	return &StatusLogger{}
}

func (this *StatusLogger) Error(format string, args ...any) {
	_, _ = fmt.Fprintf(os.Stderr, "LOG4G ERROR "+format+"\n", args...)
}

func (this *StatusLogger) Warn(format string, args ...any) {
	_, _ = fmt.Fprintf(os.Stderr, "LOG4G WARN  "+format+"\n", args...)
}

func (this *StatusLogger) Info(format string, args ...any) {
	_, _ = fmt.Fprintf(os.Stderr, "LOG4G INFO  "+format+"\n", args...)
}
