package model

type ConfigurationDefinition struct {
	Root      LoggerDefinition              `yaml:"root"`
	Loggers   map[string]LoggerDefinition   `yaml:"loggers"`
	Appenders map[string]AppenderDefinition `yaml:"appenders"`
}
