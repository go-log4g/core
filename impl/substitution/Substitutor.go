package substitution

import (
	"os"
	"regexp"
	"strings"

	"github.com/go-jang/go/lang"
	"github.com/go-jang/go/util/regex"
)

const regexVariable = `\$\{(?P<key>[^}]+)\}`
const RegexEnvVariablePattern = `{key:[^=\s]+}={value:.*}`
const RegexParameterPattern = `--?{key:[^=\s]+}\s*=?{value:.*}`

type Substitutor struct {
	properties  map[string]string
	environment map[string]string
	parameters  map[string]string
	processor   *regex.PatternProcessor
}

func NewSubstitutor(properties map[string]string) *Substitutor {
	result := &Substitutor{
		properties:  properties,
		environment: loadEnvironment(),
		parameters:  loadParameters(),
	}

	result.processor = regex.PatternProcessorOf(regexVariable)
	result.processor.OverrideResolve(func(match *regex.Match, super func(*regex.Match) any) any {
		return result.resolve(match.NamedGroup("key").Value())
	})

	return result
}

func (this *Substitutor) Substitute(value string) string {
	return this.processor.Process(value).(string)
}

func (this *Substitutor) resolve(key string) string {
	if name, ok := strings.CutPrefix(key, "env:"); ok {
		value, ok := this.environment[name]
		lang.Assert(ok, "undefined environment variable %q", name)
		return value
	}

	if name, ok := strings.CutPrefix(key, "arg:"); ok {
		value, ok := this.parameters[name]
		lang.Assert(ok, "undefined application parameter %q", name)
		return value
	}

	value, ok := this.properties[key]
	lang.Assert(ok, "undefined property %q", key)
	return value
}

func loadEnvironment() map[string]string {
	result := make(map[string]string)
	pattern := regexp.MustCompile(regex.NewPatternBuilder().Next(RegexEnvVariablePattern).Build())

	for _, keyValue := range os.Environ() {
		for _, indexes := range pattern.FindAllStringSubmatchIndex(keyValue, -1) {
			match := regex.MatchOf(pattern, keyValue, indexes)
			result[match.NamedGroup("key").Value()] = match.NamedGroup("value").Value()
		}
	}

	return result
}

func loadParameters() map[string]string {
	result := make(map[string]string)
	pattern := regexp.MustCompile(regex.NewPatternBuilder().Next(RegexParameterPattern).Build())

	for _, parameter := range os.Args[1:] {
		for _, indexes := range pattern.FindAllStringSubmatchIndex(parameter, -1) {
			match := regex.MatchOf(pattern, parameter, indexes)
			result[match.NamedGroup("key").Value()] = match.NamedGroup("value").Value()
		}
	}

	return result
}

func (this *Substitutor) RootLevel() string {
	if value, ok := this.parameters["log4g.level"]; ok {
		return value
	} else if value, ok := this.environment["LOG4G_LEVEL"]; ok {
		return value
	}
	return ""
}
