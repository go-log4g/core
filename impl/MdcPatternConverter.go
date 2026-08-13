package impl

import (
	"sort"

	"github.com/go-log4g/core/mdc"
)

type MdcPatternConverter struct {
	AbstractPatternConverter
	key string
}

func NewMdcPatternConverter(formatting FormattingInfo, key string) *MdcPatternConverter {
	return &MdcPatternConverter{
		AbstractPatternConverter: NewAbstractPatternConverter(formatting),
		key:                      key,
	}
}

func (this *MdcPatternConverter) Append(result []byte, event *LogEvent) []byte {
	start := len(result)

	if this.key != "" {
		if value, ok := mdc.Get(event.Context, this.key); ok {
			result = append(result, value...)
		}
		return this.Format(result, start)
	}

	values := mdc.Values(event.Context)
	keys := make([]string, 0, len(values))

	for key := range values {
		keys = append(keys, key)
	}

	sort.Strings(keys)

	result = append(result, '{')
	for index, key := range keys {
		if index > 0 {
			result = append(result, ", "...)
		}
		result = append(result, key...)
		result = append(result, '=')
		result = append(result, values[key]...)
	}
	result = append(result, '}')

	return this.Format(result, start)
}
