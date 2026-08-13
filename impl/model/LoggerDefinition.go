package model

type LoggerDefinition struct {
	Level     string   `yaml:"level"`
	Appenders []string `yaml:"appenders"`
	Additive  *bool    `yaml:"additive"`
}
