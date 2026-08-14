package core

import (
	"log/slog"

	"github.com/go-errr/go/err"
	"github.com/go-log4g/core/impl"
)

func init() {
	statusLogger := impl.NewStatusLogger()
	defer err.Recover(func(e any) {
		statusLogger.Error("Failed to configure Log4g: %v; using default configuration", e)
		install(impl.NewDefaultConfiguration(), statusLogger)
	})

	loader := impl.NewConfigurationLoader(statusLogger)
	definition := loader.Load()

	builder := impl.NewConfigurationBuilder(statusLogger)
	install(builder.Build(definition), statusLogger)
}

func install(configuration *impl.Configuration, statusLogger *impl.StatusLogger) {
	callerResolver := impl.NewCallerResolver()
	callerContextResolver := impl.NewCallerContextResolver(callerResolver, configuration)
	callerContextCache := impl.NewCallerContextCache(callerContextResolver)

	slog.SetDefault(slog.New(impl.NewHandler(configuration, callerContextCache, statusLogger)))
}
