package model

type AppenderDefinition struct {
	Type   string           `yaml:"type"`
	Target string           `yaml:"target"`
	File   string           `yaml:"file"`
	Layout LayoutDefinition `yaml:"layout"`
	Filter FilterDefinition `yaml:"filter"`
}
