package impl

import (
	"os"

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
	result := make([]string, 0, 3)

	if path := os.Getenv("LOG4G_CONFIGURATION_FILE"); path != "" {
		result = append(result, path)
	}

	result = append(result,
		"config/log4g.yaml",
		"log4g.yaml",
	)

	return result
}
