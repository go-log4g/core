package rolling

type SizeBasedTriggeringPolicy struct {
	maxSize int64
}

func NewSizeBasedTriggeringPolicy(maxSize int64) *SizeBasedTriggeringPolicy {
	return &SizeBasedTriggeringPolicy{
		maxSize: maxSize,
	}
}

// Implements TriggeringPolicy
func (this *SizeBasedTriggeringPolicy) IsTriggered(context TriggeringContext) bool {
	return context.FileSize+int64(context.EventSize) > this.maxSize
}
