package impl

type AbstractPatternConverter struct {
	FormattingInfo FormattingInfo
}

func NewAbstractPatternConverter(formatting FormattingInfo) AbstractPatternConverter {
	return AbstractPatternConverter{
		FormattingInfo: formatting,
	}
}

func (this *AbstractPatternConverter) AppendFormatted(result []byte, value string) []byte {
	start := len(result)
	result = append(result, value...)
	return this.Format(result, start)
}

func (this *AbstractPatternConverter) Format(result []byte, start int) []byte {
	length := len(result) - start

	if this.FormattingInfo.MaxLength > 0 && length > this.FormattingInfo.MaxLength {
		if this.FormattingInfo.LeftTruncate {
			remove := length - this.FormattingInfo.MaxLength
			copy(result[start:], result[start+remove:])
			result = result[:len(result)-remove]
		} else {
			result = result[:start+this.FormattingInfo.MaxLength]
		}

		length = this.FormattingInfo.MaxLength
	}

	padding := this.FormattingInfo.MinLength - length
	if padding <= 0 {
		return result
	}

	paddingByte := byte(' ')
	if this.FormattingInfo.ZeroPad {
		paddingByte = '0'
	}

	if this.FormattingInfo.LeftAlign {
		return this.appendPadding(result, padding, paddingByte)
	}

	result = append(result, make([]byte, padding)...)
	copy(result[start+padding:], result[start:len(result)-padding])

	for i := 0; i < padding; i++ {
		result[start+i] = paddingByte
	}

	return result
}

func (this *AbstractPatternConverter) appendPadding(result []byte, count int, value byte) []byte {
	for range count {
		result = append(result, value)
	}

	return result
}
