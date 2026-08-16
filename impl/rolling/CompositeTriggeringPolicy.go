package rolling

type CompositeTriggeringPolicy struct {
	policies []TriggeringPolicy
}

func NewCompositeTriggeringPolicy(policies ...TriggeringPolicy) *CompositeTriggeringPolicy {
	return &CompositeTriggeringPolicy{
		policies: policies,
	}
}

// Implements TriggeringPolicy
func (this *CompositeTriggeringPolicy) IsTriggered(context TriggeringContext) bool {
	triggered := false

	for _, policy := range this.policies {
		if policy.IsTriggered(context) {
			triggered = true
		}
	}

	return triggered
}
