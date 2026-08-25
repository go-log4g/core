package impl_test

import (
	"log/slog"
	"testing"

	"github.com/go-log4g/core/impl"
	"github.com/go-log4g/core/impl/model"
	"github.com/stretchr/testify/require"
)

func TestConfigurationBuilderFilters(test *testing.T) {
	statusLogger := impl.NewStatusLogger()
	builder := impl.NewConfigurationBuilder(statusLogger)

	definition := &model.ConfigurationDefinition{
		Root: model.LoggerDefinition{
			Level: "debug",
			AppenderRefs: []*model.AppenderRefDefinition{
				{Ref: "stdout"},
				{Ref: "stderr"},
			},
		},
		Appenders: map[string]model.AppenderDefinition{
			"stdout": {
				Type:   "console",
				Target: "stdout",
				Filter: model.FilterDefinition{
					Type:     "levelRange",
					MinLevel: "debug",
					MaxLevel: "warn",
				},
			},
			"stderr": {
				Type:   "console",
				Target: "stderr",
				Filter: model.FilterDefinition{
					Type:  "threshold",
					Level: "error",
				},
			},
		},
	}

	configuration := builder.Build(definition)

	require.Len(test, configuration.Appenders, 2)
	require.NotNil(test, configuration.Appenders["stdout"])
	require.NotNil(test, configuration.Appenders["stderr"])
}

func TestConfigurationBuilderRejectsInvalidFilter(test *testing.T) {
	statusLogger := impl.NewStatusLogger()
	builder := impl.NewConfigurationBuilder(statusLogger)

	definition := &model.ConfigurationDefinition{
		Root: model.LoggerDefinition{
			Level: "debug",
			AppenderRefs: []*model.AppenderRefDefinition{
				{Ref: "console"},
			},
		},
		Appenders: map[string]model.AppenderDefinition{
			"console": {
				Type:   "console",
				Target: "stdout",
				Filter: model.FilterDefinition{
					Type:     "levelRange",
					MinLevel: "error",
					MaxLevel: "debug",
				},
			},
		},
	}

	require.Panics(test, func() {
		builder.Build(definition)
	})
}

func TestConfigurationBuilderBuildsAppenderRefLevel(test *testing.T) {
	statusLogger := impl.NewStatusLogger()
	builder := impl.NewConfigurationBuilder(statusLogger)

	definition := &model.ConfigurationDefinition{
		Root: model.LoggerDefinition{
			Level: "debug",
			AppenderRefs: []*model.AppenderRefDefinition{
				{
					Ref:   "console",
					Level: "warn",
				},
			},
		},
		Appenders: map[string]model.AppenderDefinition{
			"console": {
				Type:   "console",
				Target: "stdout",
			},
		},
	}

	configuration := builder.Build(definition)
	resolved := configuration.Resolve("test")

	require.Len(test, resolved.AppenderRefs, 1)
	require.Equal(test, slog.LevelWarn, *resolved.AppenderRefs[0].Level)
	require.Same(test, configuration.Appenders["console"], resolved.AppenderRefs[0].Appender)
}

func TestConfigurationBuilderBuildsAppenderRefFilter(test *testing.T) {
	statusLogger := impl.NewStatusLogger()
	builder := impl.NewConfigurationBuilder(statusLogger)

	definition := &model.ConfigurationDefinition{
		Root: model.LoggerDefinition{
			Level: "debug",
			AppenderRefs: []*model.AppenderRefDefinition{
				{
					Ref: "console",
					Filter: model.FilterDefinition{
						Type:     "levelRange",
						MinLevel: "info",
						MaxLevel: "warn",
					},
				},
			},
		},
		Appenders: map[string]model.AppenderDefinition{
			"console": {
				Type:   "console",
				Target: "stdout",
			},
		},
	}

	configuration := builder.Build(definition)
	resolved := configuration.Resolve("test")
	appenderRef := resolved.AppenderRefs[0]

	require.True(test, appenderRef.Accept(logEvent(slog.LevelInfo)))
	require.True(test, appenderRef.Accept(logEvent(slog.LevelWarn)))
	require.False(test, appenderRef.Accept(logEvent(slog.LevelDebug)))
	require.False(test, appenderRef.Accept(logEvent(slog.LevelError)))
}

func TestConfigurationBuilderSubstitutesAppenderRef(test *testing.T) {
	statusLogger := impl.NewStatusLogger()
	builder := impl.NewConfigurationBuilder(statusLogger)

	definition := &model.ConfigurationDefinition{
		Properties: map[string]string{
			"appender": "console",
			"level":    "warn",
		},
		Root: model.LoggerDefinition{
			Level: "debug",
			AppenderRefs: []*model.AppenderRefDefinition{
				{
					Ref:   "${appender}",
					Level: "${level}",
				},
			},
		},
		Appenders: map[string]model.AppenderDefinition{
			"console": {
				Type:   "console",
				Target: "stdout",
			},
		},
	}

	configuration := builder.Build(definition)
	resolved := configuration.Resolve("test")

	require.Len(test, resolved.AppenderRefs, 1)
	require.Equal(test, slog.LevelWarn, *resolved.AppenderRefs[0].Level)
	require.Same(test, configuration.Appenders["console"], resolved.AppenderRefs[0].Appender)
}
