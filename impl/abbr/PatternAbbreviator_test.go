package abbr_test

import (
	"testing"

	"github.com/go-log4g/core/impl/abbr"
	"github.com/stretchr/testify/require"
)

func TestPatternAbbreviator(test *testing.T) {
	value := "org/apache/commons/test/Foo"

	require.Equal(test, "o/a/c/t/Foo", abbr.NewPatternAbbreviator(
		abbr.NewPatternAbbreviatorFragment(1),
	).Abbreviate(value))

	require.Equal(test, "o/ap/commons/test/Foo", abbr.NewPatternAbbreviator(
		abbr.NewPatternAbbreviatorFragment(1),
		abbr.NewPatternAbbreviatorFragment(2),
		abbr.NewPatternAbbreviatorFragment(-1),
	).Abbreviate(value))
}
