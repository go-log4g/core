package impl

import (
	"github.com/go-log4g/core/impl/model"
)

type AppenderFactory func(name string, definition *model.DynamicDefinition) Appender
