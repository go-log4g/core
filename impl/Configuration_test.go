package impl_test

import (
	"log/slog"
	"testing"

	"github.com/go-log4g/core/impl"
	"github.com/stretchr/testify/require"
)

type testAppender struct {
	name string
}

func (this *testAppender) Append(event *impl.LogEvent) error {
	return nil
}

func TestConfigurationResolveLoggerConfig(test *testing.T) {
	configuration := impl.NewConfiguration()
	configuration.Root = impl.NewLoggerConfig("", slog.LevelInfo)
	configuration.Loggers = []*impl.LoggerConfig{
		impl.NewLoggerConfig("playground/internal/app/Service1", slog.LevelWarn),
		impl.NewLoggerConfig("playground/internal/app", slog.LevelDebug),
		impl.NewLoggerConfig("github.com/go-beans/go", slog.LevelError),
	}

	require.Equal(test, slog.LevelInfo, configuration.Resolve("something/else").Level)
	require.Equal(test, slog.LevelDebug, configuration.Resolve("playground/internal/app").Level)
	require.Equal(test, slog.LevelDebug, configuration.Resolve("playground/internal/app/Service2").Level)
	require.Equal(test, slog.LevelWarn, configuration.Resolve("playground/internal/app/Service1").Level)
	require.Equal(test, slog.LevelError, configuration.Resolve("github.com/go-beans/go/ioc/ApplicationContext").Level)
	require.Equal(test, slog.LevelInfo, configuration.Resolve("playground/internal/application").Level)
}

func TestConfigurationResolveLevel(test *testing.T) {
	configuration := impl.NewConfiguration()
	configuration.Root = impl.NewLoggerConfig("", slog.LevelInfo)
	configuration.Loggers = []*impl.LoggerConfig{
		impl.NewLoggerConfig("playground/internal/app/Service1", slog.LevelWarn),
		impl.NewLoggerConfig("playground/internal/app", slog.LevelDebug),
	}

	require.Equal(test, slog.LevelWarn, configuration.Resolve("playground/internal/app/Service1").Level)
	require.Equal(test, slog.LevelDebug, configuration.Resolve("playground/internal/app/Service2").Level)
	require.Equal(test, slog.LevelInfo, configuration.Resolve("something/else").Level)
}

func TestConfigurationResolveAppenders(test *testing.T) {
	configuration := impl.NewConfiguration()

	console := &testAppender{name: "console"}
	file := &testAppender{name: "file"}
	audit := &testAppender{name: "audit"}

	configuration.Appenders["console"] = console
	configuration.Appenders["file"] = file
	configuration.Appenders["audit"] = audit

	configuration.Root.Appenders = []string{"console"}

	service := impl.NewLoggerConfig("playground/internal/app/Service1", slog.LevelInfo, "audit")
	app := impl.NewLoggerConfig("playground/internal/app", slog.LevelInfo, "file")

	configuration.Loggers = []*impl.LoggerConfig{
		service,
		app,
	}

	resolved := configuration.Resolve("playground/internal/app/Service1")

	require.Len(test, resolved.Appenders, 3)
	require.Same(test, audit, resolved.Appenders[0])
	require.Same(test, file, resolved.Appenders[1])
	require.Same(test, console, resolved.Appenders[2])
}

func TestConfigurationResolveAppendersNonAdditive(test *testing.T) {
	configuration := impl.NewConfiguration()

	console := &testAppender{name: "console"}
	file := &testAppender{name: "file"}
	audit := &testAppender{name: "audit"}

	configuration.Appenders["console"] = console
	configuration.Appenders["file"] = file
	configuration.Appenders["audit"] = audit

	configuration.Root.Appenders = []string{"console"}

	service := impl.NewLoggerConfig("playground/internal/app/Service1", slog.LevelInfo, "audit")
	service.Additive = false

	app := impl.NewLoggerConfig("playground/internal/app", slog.LevelInfo, "file")

	configuration.Loggers = []*impl.LoggerConfig{
		service,
		app,
	}

	resolved := configuration.Resolve("playground/internal/app/Service1")

	require.Len(test, resolved.Appenders, 1)
	require.Same(test, audit, resolved.Appenders[0])
}

func TestConfigurationResolveAppendersDeduplicates(test *testing.T) {
	configuration := impl.NewConfiguration()

	console := &testAppender{name: "console"}

	configuration.Appenders["console"] = console
	configuration.Root.Appenders = []string{"console"}

	app := impl.NewLoggerConfig("playground/internal/app", slog.LevelInfo, "console")
	configuration.Loggers = []*impl.LoggerConfig{app}

	resolved := configuration.Resolve("playground/internal/app/Service1")

	require.Len(test, resolved.Appenders, 1)
	require.Same(test, console, resolved.Appenders[0])
}
