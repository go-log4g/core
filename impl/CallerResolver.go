package impl

import (
	"runtime"
	"strconv"
	"strings"
)

type CallerResolver struct {
}

func NewCallerResolver() *CallerResolver {
	return &CallerResolver{}
}

func (this *CallerResolver) Resolve(pc uintptr) *Caller {
	if pc == 0 {
		return NewCaller("", "", "", 0)
	}

	frame, _ := runtime.CallersFrames([]uintptr{pc}).Next()
	logger, method := this.resolveFunction(frame.Function)

	return NewCaller(logger, method, frame.File, frame.Line)
}

func (this *CallerResolver) resolveFunction(function string) (string, string) {
	slash := strings.LastIndex(function, "/")
	start := slash + 1

	dot := strings.Index(function[start:], ".")
	if dot < 0 {
		return function, ""
	}

	dot += start

	packageName := function[:dot]
	symbol := function[dot+1:]

	return this.resolveSymbol(packageName, symbol)
}

func (this *CallerResolver) resolveSymbol(packageName, symbol string) (string, string) {
	symbol = this.resolveClosure(symbol)

	lastDot := strings.LastIndexByte(symbol, '.')
	if lastDot < 0 {
		return packageName, symbol
	}

	method := symbol[lastDot+1:]
	owner := symbol[:lastDot]

	previousDot := strings.LastIndexByte(owner, '.')
	if previousDot >= 0 {
		owner = owner[previousDot+1:]
	}

	if strings.HasPrefix(owner, "(*") && strings.HasSuffix(owner, ")") {
		owner = owner[2 : len(owner)-1]
		return packageName + "/" + owner, method
	}

	if strings.HasPrefix(owner, "(") && strings.HasSuffix(owner, ")") {
		owner = owner[1 : len(owner)-1]
		return packageName + "/" + owner, method
	}

	if previousDot < 0 {
		return packageName + "/" + owner, method
	}

	return packageName, method
}

func (this *CallerResolver) resolveClosure(symbol string) string {
	for {
		lastDot := strings.LastIndexByte(symbol, '.')
		if lastDot < 0 {
			return symbol
		}

		part := symbol[lastDot+1:]
		if _, e := strconv.Atoi(part); e == nil {
			symbol = symbol[:lastDot]
			continue
		}

		if strings.HasPrefix(part, "func") {
			if _, e := strconv.Atoi(part[4:]); e == nil {
				return symbol[:lastDot]
			}
		}

		return symbol
	}
}
