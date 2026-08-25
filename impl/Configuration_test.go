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

func (this *testAppender) Append(event *impl.LogEvent) {
	return
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

	configuration.Root = impl.NewLoggerConfig("", slog.LevelInfo, "console")

	service := impl.NewLoggerConfig("playground/internal/app/Service1", slog.LevelInfo, "audit")
	app := impl.NewLoggerConfig("playground/internal/app", slog.LevelInfo, "file")

	configuration.Loggers = []*impl.LoggerConfig{
		service,
		app,
	}

	resolved := configuration.Resolve("playground/internal/app/Service1")

	require.Len(test, resolved.AppenderRefs, 3)
	require.Same(test, audit, resolved.AppenderRefs[0].Appender)
	require.Same(test, file, resolved.AppenderRefs[1].Appender)
	require.Same(test, console, resolved.AppenderRefs[2].Appender)
}

func TestConfigurationResolveAppendersNonAdditive(test *testing.T) {
	configuration := impl.NewConfiguration()

	console := &testAppender{name: "console"}
	file := &testAppender{name: "file"}
	audit := &testAppender{name: "audit"}

	configuration.Appenders["console"] = console
	configuration.Appenders["file"] = file
	configuration.Appenders["audit"] = audit

	configuration.Root = impl.NewLoggerConfig("", slog.LevelInfo, "console")

	service := impl.NewLoggerConfig("playground/internal/app/Service1", slog.LevelInfo, "audit")
	service.Additive = false

	app := impl.NewLoggerConfig("playground/internal/app", slog.LevelInfo, "file")

	configuration.Loggers = []*impl.LoggerConfig{
		service,
		app,
	}

	resolved := configuration.Resolve("playground/internal/app/Service1")

	require.Len(test, resolved.AppenderRefs, 1)
	require.Same(test, audit, resolved.AppenderRefs[0].Appender)
}

func TestConfigurationResolveAppendersDeduplicates(test *testing.T) {
	configuration := impl.NewConfiguration()

	console := &testAppender{name: "console"}

	configuration.Appenders["console"] = console
	configuration.Root = impl.NewLoggerConfig("", slog.LevelInfo, "console")

	app := impl.NewLoggerConfig("playground/internal/app", slog.LevelInfo, "console")
	configuration.Loggers = []*impl.LoggerConfig{app}

	resolved := configuration.Resolve("playground/internal/app/Service1")

	require.Len(test, resolved.AppenderRefs, 1)
	require.Same(test, console, resolved.AppenderRefs[0].Appender)
}

func TestConfigurationResolveAppenderRefLevel(test *testing.T) {
	configuration := impl.NewConfiguration()

	file := &testAppender{name: "file"}
	configuration.Appenders["file"] = file

	info := slog.LevelInfo
	error := slog.LevelError

	service1 := impl.NewLoggerConfig("service1", slog.LevelDebug)
	service1.AppenderRefs = []*impl.AppenderRef{
		{Ref: "file", Level: &info},
	}

	service2 := impl.NewLoggerConfig("service2", slog.LevelDebug)
	service2.AppenderRefs = []*impl.AppenderRef{
		{Ref: "file", Level: &error},
	}

	configuration.Loggers = []*impl.LoggerConfig{
		service1,
		service2,
	}

	resolved1 := configuration.Resolve("service1")
	resolved2 := configuration.Resolve("service2")

	require.Len(test, resolved1.AppenderRefs, 1)
	require.Same(test, file, resolved1.AppenderRefs[0].Appender)
	require.Equal(test, slog.LevelInfo, *resolved1.AppenderRefs[0].Level)

	require.Len(test, resolved2.AppenderRefs, 1)
	require.Same(test, file, resolved2.AppenderRefs[0].Appender)
	require.Equal(test, slog.LevelError, *resolved2.AppenderRefs[0].Level)
}
