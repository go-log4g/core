package format

import (
	"math"
	"strconv"
	"strings"

	"github.com/go-jang/go/lang"
)

var fileSizeUnits = []struct {
	suffix     string
	multiplier int64
}{
	{"KB", 1024},
	{"K", 1024},
	{"MB", 1024 * 1024},
	{"M", 1024 * 1024},
	{"GB", 1024 * 1024 * 1024},
	{"G", 1024 * 1024 * 1024},
	{"TB", 1024 * 1024 * 1024 * 1024},
	{"T", 1024 * 1024 * 1024 * 1024},
	{"B", 1},
}

var javaTimePatternReplacements = []struct {
	pattern string
	layout  string
	unit    Unit
}{
	{"SSS", "000", Second},
	{"ss", "05", Second},
	{"mm", "04", Minute},
	{"HH", "15", Hour},
	{"dd", "02", Day},
	{"MM", "01", Month},
	{"yyyy", "2006", Year},
	{"yy", "06", Year},
}

func ParseFileSize(value string, defaultValue int64) int64 {
	original := value
	value = strings.ToUpper(strings.TrimSpace(value))

	if value == "" {
		return defaultValue
	}

	multiplier := int64(1)
	for _, unit := range fileSizeUnits {
		if strings.HasSuffix(value, unit.suffix) {
			multiplier = unit.multiplier
			value = strings.TrimSpace(strings.TrimSuffix(value, unit.suffix))
			break
		}
	}

	size, e := strconv.ParseFloat(value, 64)
	lang.Assert(e == nil && size > 0, "invalid file size %q", original)
	lang.Assert(size <= float64(math.MaxInt64)/float64(multiplier), "file size %q is too large", original)

	return int64(size * float64(multiplier))
}

func ConvertJavaPattern(pattern string) (string, Unit) {
	result := pattern
	smallestUnit := Year

	for _, replacement := range javaTimePatternReplacements {
		if strings.Contains(result, replacement.pattern) {
			result = strings.ReplaceAll(result, replacement.pattern, replacement.layout)
			if replacement.unit > smallestUnit {
				smallestUnit = replacement.unit
			}
		}
	}

	for _, ch := range result {
		lang.Assert(ch < 'A' || ch > 'Z' && ch < 'a' || ch > 'z', "unsupported time pattern %q", pattern)
	}

	return result, smallestUnit
}
