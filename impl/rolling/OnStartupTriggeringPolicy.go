package rolling

type OnStartupTriggeringPolicy struct {
	minSize int64
}

func NewOnStartupTriggeringPolicy(minSize int64) *OnStartupTriggeringPolicy {
	return &OnStartupTriggeringPolicy{
		minSize: minSize,
	}
}

func (this *OnStartupTriggeringPolicy) IsTriggered(fileSize int64) bool {
	return fileSize >= this.minSize
}