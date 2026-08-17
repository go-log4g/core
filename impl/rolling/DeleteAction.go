package rolling

import (
	"os"
	"path/filepath"
	"time"

	"github.com/go-jang/go/lang"
)

type DeleteAction struct {
	dir     string
	pattern string
	maxAge  time.Duration
}

func NewDeleteAction(filePattern string, maxAge time.Duration) *DeleteAction {
	dir, pattern := FilePatternMatcher(filePattern)

	return &DeleteAction{
		dir:     dir,
		pattern: pattern,
		maxAge:  maxAge,
	}
}

func (this *DeleteAction) Execute() {
	entries, e := os.ReadDir(this.dir)
	lang.Assert(e == nil, "failed to read rollover directory %q", this.dir)

	threshold := time.Now().Add(-this.maxAge)

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		match, e := filepath.Match(this.pattern, entry.Name())
		lang.Assert(e == nil, "invalid rollover file pattern %q", this.pattern)

		if !match {
			continue
		}

		info, e := entry.Info()
		lang.Assert(e == nil, "failed to stat rollover file %q", entry.Name())

		if info.ModTime().Before(threshold) {
			file := filepath.Join(this.dir, entry.Name())
			lang.Assert(os.Remove(file) == nil, "failed to delete expired rollover file %q", file)
		}
	}
}
