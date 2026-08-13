package impl

import "sync"

type CallerContextCache struct {
	resolver *CallerContextResolver
	values   sync.Map
}

func NewCallerContextCache(resolver *CallerContextResolver) *CallerContextCache {
	return &CallerContextCache{
		resolver: resolver,
	}
}

func (this *CallerContextCache) Get(pc uintptr) *CallerContext {
	if value, ok := this.values.Load(pc); ok {
		return value.(*CallerContext)
	}

	callerContext := this.resolver.Resolve(pc)
	actual, _ := this.values.LoadOrStore(pc, callerContext)
	return actual.(*CallerContext)
}
