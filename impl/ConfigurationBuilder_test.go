package impl_test

import (
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
			Level:     "debug",
			Appenders: []string{"stdout", "stderr"},
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
			Level:     "debug",
			Appenders: []string{"console"},
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
