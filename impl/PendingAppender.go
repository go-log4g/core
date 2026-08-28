package impl

import "github.com/go-log4g/core/impl/model"

type PendingAppender struct {
	Type       string
	Definition *model.DynamicDefinition
	Builder    *ConfigurationBuilder
}
