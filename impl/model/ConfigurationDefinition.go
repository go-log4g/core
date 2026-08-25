package model

type ConfigurationDefinition struct {
	Properties map[string]string             `yaml:"properties"`
	Appenders  map[string]*DynamicDefinition `yaml:"appenders"`
	Root       LoggerDefinition              `yaml:"root"`
	Loggers    map[string]LoggerDefinition   `yaml:"loggers"`
}
