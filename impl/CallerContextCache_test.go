package impl_test

import (
	"log/slog"
	"reflect"
	"testing"

	"github.com/go-log4g/core/impl"
	"github.com/stretchr/testify/require"
)

func TestCallerContextCache(test *testing.T) {
	config := impl.NewConfiguration()
	loggerConfig := impl.NewLoggerConfig("github.com/go-log4g/core/impl_test/callerTestService", slog.LevelDebug)
	config.Loggers = []*impl.LoggerConfig{loggerConfig}

	callerResolver := impl.NewCallerResolver()
	callerContextResolver := impl.NewCallerContextResolver(callerResolver, config)
	cache := impl.NewCallerContextCache(callerContextResolver)

	pc := reflect.ValueOf((*callerTestService).PointerMethod).Pointer()

	first := cache.Get(pc)
	second := cache.Get(pc)

	require.Same(test, first, second)
	require.Same(test, first.Caller, second.Caller)
	require.Same(test, first.Config, second.Config)
}
