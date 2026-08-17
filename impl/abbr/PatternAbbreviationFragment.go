package abbr

type PatternAbbreviatorFragment struct {
	CharCount int
}

func NewPatternAbbreviatorFragment(charCount int) *PatternAbbreviatorFragment {
	return &PatternAbbreviatorFragment{
		CharCount: charCount,
	}
}
