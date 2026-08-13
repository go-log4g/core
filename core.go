package core

import (
	"log/slog"

	"github.com/go-log4g/core/impl"
)

func init() {
	statusLogger := impl.NewStatusLogger()

	loader := impl.NewConfigurationLoader(statusLogger)
	definition := loader.Load()

	builder := impl.NewConfigurationBuilder(statusLogger)
	configuration := builder.Build(definition)

	callerResolver := impl.NewCallerResolver()
	callerContextResolver := impl.NewCallerContextResolver(callerResolver, configuration)
	callerContextCache := impl.NewCallerContextCache(callerContextResolver)

	slog.SetDefault(slog.New(impl.NewHandler(configuration, callerContextCache)))
}
