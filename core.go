package core

import (
	"log/slog"

	"github.com/go-log4g/core/impl"
)

func init() {
	statusLogger := impl.NewStatusLogger()

	defer func() {
		if recovered := recover(); recovered != nil {
			statusLogger.Error("Failed to configure log4g: %v; using default configuration", recovered)
			install(impl.NewDefaultConfiguration())
		}
	}()

	loader := impl.NewConfigurationLoader(statusLogger)
	definition := loader.Load()

	builder := impl.NewConfigurationBuilder(statusLogger)
	install(builder.Build(definition))
}

func install(configuration *impl.Configuration) {
	callerResolver := impl.NewCallerResolver()
	callerContextResolver := impl.NewCallerContextResolver(callerResolver, configuration)
	callerContextCache := impl.NewCallerContextCache(callerContextResolver)

	slog.SetDefault(slog.New(impl.NewHandler(configuration, callerContextCache)))
}
