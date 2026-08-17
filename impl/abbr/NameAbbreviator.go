package abbr

type NameAbbreviator interface {
	Abbreviate(value string) string
}
