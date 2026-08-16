package rolling

import (
	"strings"

	"github.com/go-jang/go/lang"
	"github.com/go-log4g/core/impl/format"
)

func TimeUnitFromFilePattern(filePattern string) format.Unit {
	start := strings.LastIndex(filePattern, "%d{")
	lang.Assert(start >= 0, "filePattern %q must contain %%d{...}", filePattern)

	start += 3
	end := strings.IndexByte(filePattern[start:], '}')
	lang.Assert(end >= 0, "unterminated date pattern in filePattern %q", filePattern)

	_, unit := format.ConvertJavaPattern(filePattern[start : start+end])
	return unit
}
