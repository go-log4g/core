package substitution_test

import (
	"testing"

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

func TestSubstitutorSelfReference(test *testing.T) {
	value := substitution.NewSubstitutor(map[string]string{
		"a": "${a}",
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
