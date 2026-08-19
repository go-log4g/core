package core

import (
	"github.com/go-log4g/core/impl"
)

func init() {
	impl.Initialize(
		"config/log4g.yaml",
		"log4g.yaml",
	)
}
