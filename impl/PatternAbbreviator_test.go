package impl_test

import (
	"testing"

	"github.com/go-log4g/core/impl"
	"github.com/stretchr/testify/require"
)

func TestPatternAbbreviator(test *testing.T) {
	value := "org/apache/commons/test/Foo"

	require.Equal(test, "o/a/c/t/Foo", impl.NewPatternAbbreviator(
		impl.NewPatternAbbreviatorFragment(1),
	).Abbreviate(value))

	require.Equal(test, "o/ap/commons/test/Foo", impl.NewPatternAbbreviator(
		impl.NewPatternAbbreviatorFragment(1),
		impl.NewPatternAbbreviatorFragment(2),
		impl.NewPatternAbbreviatorFragment(-1),
	).Abbreviate(value))
}
