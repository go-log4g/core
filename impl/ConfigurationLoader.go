package impl

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/go-jang/go/lang"
	"github.com/go-jang/go/util/optional"
	"github.com/go-jang/go/util/stream"
	"github.com/go-log4g/core/impl/model"
	"gopkg.in/yaml.v3"
)

type ConfigurationLoader struct {
	statusLogger *StatusLogger
}

func NewConfigurationLoader(statusLogger *StatusLogger) *ConfigurationLoader {
	return &ConfigurationLoader{
		statusLogger: statusLogger,
	}
}

func (this *ConfigurationLoader) Load() *model.ConfigurationDefinition {
	for _, path := range this.paths() {
		data, e := os.ReadFile(path)

		if os.IsNotExist(e) {
			continue
		}
		if e != nil {
			this.statusLogger.Error("Cannot read configuration %s: %v", path, e)
			this.statusLogger.Warn("Using default configuration")
			return nil
		}

		result := &model.ConfigurationDefinition{}

		if e := yaml.Unmarshal(data, result); e != nil {
			this.statusLogger.Error("Cannot parse configuration %s: %v", path, e)
			this.statusLogger.Warn("Using default configuration")
			return nil
		}

		this.statusLogger.Info("Loaded configuration %s", path)
		return result
	}

	this.statusLogger.Warn("No configuration found, using default configuration")
	return nil
}

func (this *ConfigurationLoader) paths() []string {
	result := make([]string, 0, 5)
	if path := substitutor.ConfigurationFile(); path != "" {
		result = append(result, path)
	}

	if isTest() {
		root := moduleRoot()
		return append(result,
			filepath.Join(root, "config", "log4g-test.yaml"),
			filepath.Join(root, "log4g-test.yaml"),
			filepath.Join(root, "config", "log4g.yaml"),
			filepath.Join(root, "log4g.yaml"),
		)
	}

	return append(result,
		"config/log4g.yaml",
		"log4g.yaml",
	)
}

func isTest() bool {
	return stream.From(os.Args[1:]).
		Filter(func(s string) bool { return strings.HasPrefix(s, "-test.timeout=") }).
		FindFirst().Present()
}

func moduleRoot() string {
	dir := optional.OfCommaErr(os.Getwd()).OrElsePanic("Cannot get working directory")
	for {
		if _, e := os.Stat(filepath.Join(dir, "go.mod")); e == nil {
			return dir
		}
		parent := filepath.Dir(dir)
		lang.Assert(parent != dir, "Cannot find go.mod from %q", dir)
		dir = parent
	}
}
