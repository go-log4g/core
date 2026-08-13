package impl

type CallerContextResolver struct {
	callerResolver *CallerResolver
	configuration  *Configuration
}

func NewCallerContextResolver(callerResolver *CallerResolver, configuration *Configuration) *CallerContextResolver {
	return &CallerContextResolver{
		callerResolver: callerResolver,
		configuration:  configuration,
	}
}

func (this *CallerContextResolver) Resolve(pc uintptr) *CallerContext {
	caller := this.callerResolver.Resolve(pc)
	config := this.configuration.Resolve(caller.Logger)

	return &CallerContext{
		Caller: caller,
		Config: config,
	}
}
