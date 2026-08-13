package impl

import "strings"

type DynamicWordAbbreviator struct {
	charCount   int
	retainCount int
}

func NewDynamicWordAbbreviator(charCount, retainCount int) *DynamicWordAbbreviator {
	return &DynamicWordAbbreviator{
		charCount:   charCount,
		retainCount: retainCount,
	}
}

func (this *DynamicWordAbbreviator) Abbreviate(value string) string {
	componentCount := strings.Count(value, "/") + 1
	if componentCount <= this.retainCount {
		return value
	}

	abbreviateCount := componentCount - this.retainCount
	var result strings.Builder
	start := 0

	for range abbreviateCount {
		end := strings.IndexByte(value[start:], '/')
		if end < 0 {
			return value
		}
		end += start

		length := min(this.charCount, end-start)
		result.WriteString(value[start : start+length])
		result.WriteByte('/')
		start = end + 1
	}

	result.WriteString(value[start:])
	return result.String()
}
