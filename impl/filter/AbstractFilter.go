package filter

type AbstractFilter struct {
	OnMatch    Result
	OnMismatch Result
}

func NewAbstractFilter(onMatch, onMismatch Result) AbstractFilter {
	return AbstractFilter{
		OnMatch:    onMatch,
		OnMismatch: onMismatch,
	}
}
