package impl

type FormattingInfo struct {
	LeftAlign    bool
	MinLength    int
	MaxLength    int
	LeftTruncate bool
	ZeroPad      bool
}

func NewDefaultFormattingInfo() FormattingInfo {
	return NewFormattingInfo(false, 0, 0, true, false)
}

func NewFormattingInfo(leftAlign bool, minLength, maxLength int, leftTruncate, zeroPad bool) FormattingInfo {
	return FormattingInfo{
		LeftAlign:    leftAlign,
		MinLength:    minLength,
		MaxLength:    maxLength,
		LeftTruncate: leftTruncate,
		ZeroPad:      zeroPad,
	}
}
