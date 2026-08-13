package filter

type Result int

const (
	Neutral Result = iota
	Accept
	Deny
)
