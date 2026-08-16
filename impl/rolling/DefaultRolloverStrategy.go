package rolling

import (
	"os"
	"strings"
	"time"

	"github.com/go-jang/go/lang"
)

type DefaultRolloverStrategy struct {
	max int
}

func NewDefaultRolloverStrategy(max int) *DefaultRolloverStrategy {
	lang.Assert(max > 0, "max must be positive")
	return &DefaultRolloverStrategy{
		max: max,
	}
}

// Implements RolloverStrategy
func (this *DefaultRolloverStrategy) Rollover(file, filePattern string, at time.Time) {
	lang.Assert(strings.Contains(filePattern, "%i"), "filePattern must contain %%i")

	for i := 1; i <= this.max; i++ {
		target := FormatFilePattern(filePattern, i, at)

		if _, e := os.Stat(target); os.IsNotExist(e) {
			lang.Assert(os.Rename(file, target) == nil, "failed to rename log file %q to %q", file, target)
			return
		}
	}

	oldest := FormatFilePattern(filePattern, 1, at)
	lang.Assert(os.Remove(oldest) == nil, "failed to remove rollover file %q", oldest)

	for i := 2; i <= this.max; i++ {
		source := FormatFilePattern(filePattern, i, at)
		target := FormatFilePattern(filePattern, i-1, at)

		lang.Assert(os.Rename(source, target) == nil, "failed to rename rollover file %q to %q", source, target)
	}

	newest := FormatFilePattern(filePattern, this.max, at)
	lang.Assert(os.Rename(file, newest) == nil, "failed to rename log file %q to %q", file, newest)
}
