package substitution

import (
	"fmt"
	"strings"
)

type Substitutor struct {
	properties map[string]string
}

func NewSubstitutor(properties map[string]string) *Substitutor {
	return &Substitutor{
		properties: properties,
	}
}

func (this *Substitutor) Substitute(value string) string {
	return this.substitute(value, make(map[string]bool))
}

func (this *Substitutor) substitute(value string, resolving map[string]bool) string {
	for {
		start := strings.Index(value, "${")
		if start < 0 {
			return value
		}

		end := strings.IndexByte(value[start+2:], '}')
		if end < 0 {
			panic(fmt.Errorf("incomplete property substitution %q", value))
		}
		end += start + 2

		name := value[start+2 : end]

		if resolving[name] {
			panic(fmt.Errorf("recursive property substitution %q", name))
		}

		replacement, ok := this.properties[name]
		if !ok {
			panic(fmt.Errorf("undefined property %q", name))
		}

		resolving[name] = true
		replacement = this.substitute(replacement, resolving)
		delete(resolving, name)

		value = value[:start] + replacement + value[end+1:]
	}
}
