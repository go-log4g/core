package impl

import "strings"

type MaxElementAbbreviator struct {
	precision int
}

func NewMaxElementAbbreviator(precision int) *MaxElementAbbreviator {
	return &MaxElementAbbreviator{
		precision: precision,
	}
}

func (this *MaxElementAbbreviator) Abbreviate(value string) string {
	if this.precision > 0 {
		return this.keepRight(value, this.precision)
	}

	if this.precision < 0 {
		return this.removeLeft(value, -this.precision)
	}

	return this.keepRight(value, 1)
}

func (this *MaxElementAbbreviator) keepRight(value string, count int) string {
	index := len(value)

	for range count {
		next := strings.LastIndexByte(value[:index], '/')
		if next < 0 {
			return value
		}
		index = next
	}

	return value[index+1:]
}

func (this *MaxElementAbbreviator) removeLeft(value string, count int) string {
	index := 0

	for range count {
		next := strings.IndexByte(value[index:], '/')
		if next < 0 {
			return value
		}
		index += next + 1
	}

	return value[index:]
}
