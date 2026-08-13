package model

type FilterDefinition struct {
	Type     string `yaml:"type"`
	Level    string `yaml:"level"`
	MinLevel string `yaml:"minLevel"`
	MaxLevel string `yaml:"maxLevel"`
}
