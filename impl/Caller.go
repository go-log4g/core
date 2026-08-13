package impl

type Caller struct {
	Logger string
	Method string
	File   string
	Line   int
}

func NewCaller(logger, method, file string, line int) *Caller {
	return &Caller{
		Logger: logger,
		Method: method,
		File:   file,
		Line:   line,
	}
}
