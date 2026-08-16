package rolling

type TriggeringPolicy interface {
	IsTriggered(context TriggeringContext) bool
}