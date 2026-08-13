package impl_test

import (
	"testing"

	"github.com/go-log4g/core/impl"
	"github.com/stretchr/testify/require"
)

func TestParameterizedMessageFormatter(test *testing.T) {
	formatter := impl.NewParameterizedMessageFormatter()

	require.Equal(test, "Hello world", formatter.Format("Hello {}", "world"))
	require.Equal(test, "A 1 B 2", formatter.Format("A {} B {}", 1, 2))
	require.Equal(test, "A 1 B {}", formatter.Format("A {} B {}", 1))
	require.Equal(test, "A 1", formatter.Format("A {}", 1, 2))
	require.Equal(test, "A {}", formatter.Format(`A \{}`, 1))
	require.Equal(test, `A \1`, formatter.Format(`A \\{}`, 1))
}
