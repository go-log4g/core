package abbr

type NOPAbbreviator struct {
}

func NewNOPAbbreviator() *NOPAbbreviator {
	return &NOPAbbreviator{}
}

func (this *NOPAbbreviator) Abbreviate(value string) string {
	return value
}
