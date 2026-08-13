package substitution_test

import (
	"regexp"
	"testing"

	"github.com/go-jang/go/util/regex"
	"github.com/go-log4g/core/impl/substitution"
	"github.com/stretchr/testify/require"
)

func TestSubstitutor(test *testing.T) {
	value := substitution.NewSubstitutor(map[string]string{
		"pattern":   "${timestamp} %-5p %c.%M:%L - %m%n",
		"timestamp": "%d{yyyy-MM-dd HH:mm:ss.SSS}{UTC}",
	})

	require.Equal(test, "Hello world", value.Substitute("Hello world"))
	require.Equal(test, "%d{yyyy-MM-dd HH:mm:ss.SSS}{UTC}", value.Substitute("${timestamp}"))
	require.Equal(test, "%d{yyyy-MM-dd HH:mm:ss.SSS}{UTC} %-5p %c.%M:%L - %m%n", value.Substitute("${pattern}"))
	require.Equal(test, "before %d{yyyy-MM-dd HH:mm:ss.SSS}{UTC} after", value.Substitute("before ${timestamp} after"))
}

func TestSubstitutorUndefinedProperty(test *testing.T) {
	value := substitution.NewSubstitutor(map[string]string{
		"pattern": "${unknown}",
	})

	require.Panics(test, func() {
		value.Substitute("${pattern}")
	})
}

func TestSubstitutorRecursiveProperty(test *testing.T) {
	value := substitution.NewSubstitutor(map[string]string{
		"a": "${b}",
		"b": "${a}",
	})

	require.Panics(test, func() {
		value.Substitute("${a}")
	})
}

func TestSubstitutorMultipleProperties(test *testing.T) {
	value := substitution.NewSubstitutor(map[string]string{
		"date":  "%d{HH:mm:ss}",
		"level": "%-5p",
	})

	require.Equal(test, "%d{HH:mm:ss} %-5p %m%n", value.Substitute("${date} ${level} %m%n"))
}

func TestEnvironmentPattern(test *testing.T) {
	pattern := regexp.MustCompile(regex.NewPatternBuilder().Next(substitution.RegexEnvVariablePattern).Build())

	match := regex.MatchOf(pattern, "LOG_LEVEL=debug", pattern.FindStringSubmatchIndex("LOG_LEVEL=debug"))

	require.Equal(test, "LOG_LEVEL", match.NamedGroup("key").Value())
	require.Equal(test, "debug", match.NamedGroup("value").Value())
}

func TestEnvironmentPatternValueWithEquals(test *testing.T) {
	pattern := regexp.MustCompile(regex.NewPatternBuilder().Next(substitution.RegexEnvVariablePattern).Build())

	match := regex.MatchOf(pattern, "URL=a=b=c", pattern.FindStringSubmatchIndex("URL=a=b=c"))

	require.Equal(test, "URL", match.NamedGroup("key").Value())
	require.Equal(test, "a=b=c", match.NamedGroup("value").Value())
}

func TestParameterPattern(test *testing.T) {
	pattern := regexp.MustCompile(regex.NewPatternBuilder().Next(substitution.RegexParameterPattern).Build())

	match := regex.MatchOf(pattern, "--log.level=debug", pattern.FindStringSubmatchIndex("--log.level=debug"))

	require.Equal(test, "log.level", match.NamedGroup("key").Value())
	require.Equal(test, "debug", match.NamedGroup("value").Value())
}

func TestParameterPatternSingleDash(test *testing.T) {
	pattern := regexp.MustCompile(regex.NewPatternBuilder().Next(substitution.RegexParameterPattern).Build())

	match := regex.MatchOf(pattern, "-log.level=debug", pattern.FindStringSubmatchIndex("-log.level=debug"))

	require.Equal(test, "log.level", match.NamedGroup("key").Value())
	require.Equal(test, "debug", match.NamedGroup("value").Value())
}
