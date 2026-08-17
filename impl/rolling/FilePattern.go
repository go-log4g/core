package rolling

import (
	"fmt"
	"path/filepath"
	"strings"
	"time"

	"github.com/go-jang/go/lang"
	"github.com/go-log4g/core/impl/format"
)

func FormatFilePattern(pattern string, index int, at time.Time) string {
	result := strings.ReplaceAll(pattern, "%i", fmt.Sprint(index))

	for {
		start := strings.Index(result, "%d{")
		if start < 0 {
			return result
		}

		patternStart := start + 3
		end := strings.IndexByte(result[patternStart:], '}')
		lang.Assert(end >= 0, "unterminated date pattern in filePattern %q", pattern)

		end += patternStart
		datePattern := result[patternStart:end]
		goPattern, _ := format.ConvertJavaPattern(datePattern)
		value := at.Format(goPattern)
		result = result[:start] + value + result[end+1:]
	}
}

func FilePatternMatcher(filePattern string) (string, string) {
	dir := filepath.Dir(filePattern)
	name := filepath.Base(filePattern)

	for {
		start := strings.Index(name, "%d{")
		if start < 0 {
			break
		}

		end := strings.IndexByte(name[start+3:], '}')
		lang.Assert(end >= 0, "unterminated date pattern in filePattern %q", filePattern)

		end += start + 3
		name = name[:start] + "*" + name[end+1:]
	}

	name = strings.ReplaceAll(name, "%i", "*")
	return dir, name
}
