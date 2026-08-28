package core

import (
	"log/slog"

	"github.com/go-errr/go/err"
	"github.com/go-log4g/core/impl"
)

var configuration *impl.Configuration
var statusLogger *impl.StatusLogger

func init() {
	statusLogger = impl.NewStatusLogger()
	defer err.Recover(func(e any) {
		statusLogger.Error("Failed to configure Log4g: %v; using default configuration", e)
		install(impl.NewDefaultConfiguration(), statusLogger)
	})

	loader := impl.NewConfigurationLoader(statusLogger)
	definition := loader.Load()

	builder := impl.NewConfigurationBuilder(statusLogger)
	install(builder.Build(definition), statusLogger)
}

func install(value *impl.Configuration, statusLogger *impl.StatusLogger) {
	configuration = value

	callerResolver := impl.NewCallerResolver()
	callerContextResolver := impl.NewCallerContextResolver(callerResolver, configuration)
	callerContextCache := impl.NewCallerContextCache(callerContextResolver)

	slog.SetDefault(slog.New(impl.NewHandler(configuration, callerContextCache, statusLogger)))
}

func RegisterAppender(appenderType string, factory impl.AppenderFactory) {
	defer err.Recover(func(e any) {
		statusLogger.Error("Failed to configure Log4g: %v; using default configuration", e)
		install(impl.NewDefaultConfiguration(), statusLogger)
	})

	configuration.ActivatePendingAppenders(appenderType, factory)
}
