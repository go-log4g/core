package test

import (
	"os"
	"path/filepath"

	"github.com/go-jang/go/lang"
	"github.com/go-log4g/core/impl"
)

func init() {
	root := moduleRoot()

	impl.Initialize(
		filepath.Join(root, "config", "log4g-test.yaml"),
		filepath.Join(root, "log4g-test.yaml"),
		filepath.Join(root, "config", "log4g.yaml"),
		filepath.Join(root, "log4g.yaml"),
	)
}

func moduleRoot() string {
	dir, e := os.Getwd()
	lang.Assert(e == nil, "failed to get working directory")

	for {
		if _, e := os.Stat(filepath.Join(dir, "go.mod")); e == nil {
			return dir
		}

		parent := filepath.Dir(dir)
		lang.Assert(parent != dir, "cannot find go.mod from %q", dir)
		dir = parent
	}
}
