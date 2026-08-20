package impl

import (
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/go-errr/go/err"
)

const statusErrorBurst = 3
const statusErrorInterval = 5 * time.Minute

type StatusLogger struct {
	mutex  sync.Mutex
	errors map[string]*statusErrorState
}

func NewStatusLogger() *StatusLogger {
	return &StatusLogger{
		errors: make(map[string]*statusErrorState),
	}
}

func (this *StatusLogger) Info(format string, args ...any) {
	_, _ = fmt.Fprintf(os.Stdout, "LOG4G INFO  "+format+"\n", args...)
}

func (this *StatusLogger) Warn(format string, args ...any) {
	_, _ = fmt.Fprintf(os.Stderr, "LOG4G WARN  "+format+"\n", args...)
}

func (this *StatusLogger) Error(format string, args ...any) {
	_, _ = fmt.Fprintf(os.Stderr, "LOG4G ERROR "+format+"\n", args...)
}

func (this *StatusLogger) ErrorThrottled(key string, e any, format string, args ...any) {
	message := fmt.Sprintf(format, args...)
	now := time.Now()

	this.mutex.Lock()
	defer this.mutex.Unlock()

	state := this.errors[key]
	if state == nil {
		state = &statusErrorState{}
		this.errors[key] = state
	} else if now.Sub(state.lastError) >= statusErrorInterval {
		state.count = 0
	}

	state.count++
	state.lastError = now
	allowed := state.count <= statusErrorBurst

	if allowed {
		_, _ = fmt.Fprintf(os.Stderr, "LOG4G ERROR %s: %s\n", message, err.PrintStackTrace(e))
	}
}

type statusErrorState struct {
	count     int
	lastError time.Time
}
