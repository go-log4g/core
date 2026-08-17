package rolling_test

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/go-jang/go/util/optional"
	"github.com/go-log4g/core/impl/rolling"
	"github.com/stretchr/testify/require"
)

func TestDefaultRolloverStrategy(test *testing.T) {
	dir := test.TempDir()
	file := filepath.Join(dir, "application.log")
	pattern := filepath.Join(dir, "application-%i.log")

	require.NoError(test, os.WriteFile(file, []byte("current"), 0644))
	require.NoError(test, os.WriteFile(filepath.Join(dir, "application-1.log"), []byte("one"), 0644))
	require.NoError(test, os.WriteFile(filepath.Join(dir, "application-2.log"), []byte("two"), 0644))
	require.NoError(test, os.WriteFile(filepath.Join(dir, "application-3.log"), []byte("three"), 0644))

	strategy := rolling.NewDefaultRolloverStrategy(3, nil)
	strategy.Rollover(file, pattern, time.Now())

	require.Equal(test, "two", string(optional.OfCommaErr(os.ReadFile(filepath.Join(dir, "application-1.log"))).OrElsePanic("cannot read file")))
	require.Equal(test, "three", string(optional.OfCommaErr(os.ReadFile(filepath.Join(dir, "application-2.log"))).OrElsePanic("cannot read file")))
	require.Equal(test, "current", string(optional.OfCommaErr(os.ReadFile(filepath.Join(dir, "application-3.log"))).OrElsePanic("cannot read file")))

	_, e := os.Stat(file)
	require.True(test, os.IsNotExist(e))
}

func TestDefaultRolloverStrategyNotFull(test *testing.T) {
	dir := test.TempDir()
	file := filepath.Join(dir, "application.log")
	pattern := filepath.Join(dir, "application-%i.log")

	require.NoError(test, os.WriteFile(file, []byte("current"), 0644))
	require.NoError(test, os.WriteFile(filepath.Join(dir, "application-1.log"), []byte("one"), 0644))
	require.NoError(test, os.WriteFile(filepath.Join(dir, "application-2.log"), []byte("two"), 0644))

	strategy := rolling.NewDefaultRolloverStrategy(3, nil)
	strategy.Rollover(file, pattern, time.Now())

	require.Equal(test, "one", string(optional.OfCommaErr(os.ReadFile(filepath.Join(dir, "application-1.log"))).OrElsePanic("cannot read file")))
	require.Equal(test, "two", string(optional.OfCommaErr(os.ReadFile(filepath.Join(dir, "application-2.log"))).OrElsePanic("cannot read file")))
	require.Equal(test, "current", string(optional.OfCommaErr(os.ReadFile(filepath.Join(dir, "application-3.log"))).OrElsePanic("cannot read file")))

	_, e := os.Stat(file)
	require.True(test, os.IsNotExist(e))
}
