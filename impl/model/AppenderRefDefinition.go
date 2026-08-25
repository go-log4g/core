package model

type AppenderRefDefinition struct {
	Ref    string           `yaml:"ref"`
	Level  string           `yaml:"level"`
	Filter FilterDefinition `yaml:"filter"`
}
