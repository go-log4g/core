package impl

import "strings"

type PatternAbbreviator struct {
	fragments []*PatternAbbreviatorFragment
}

func NewPatternAbbreviator(fragments ...*PatternAbbreviatorFragment) *PatternAbbreviator {
	return &PatternAbbreviator{
		fragments: fragments,
	}
}

func (this *PatternAbbreviator) Abbreviate(value string) string {
	var result strings.Builder
	start := 0
	fragmentIndex := 0

	for {
		end := strings.IndexByte(value[start:], '/')
		if end < 0 {
			result.WriteString(value[start:])
			break
		}

		end += start
		fragment := this.fragments[min(fragmentIndex, len(this.fragments)-1)]

		if fragment.CharCount < 0 || end-start <= fragment.CharCount {
			result.WriteString(value[start:end])
		} else {
			result.WriteString(value[start : start+fragment.CharCount])
		}

		result.WriteByte('/')
		start = end + 1
		fragmentIndex++
	}

	return result.String()
}
