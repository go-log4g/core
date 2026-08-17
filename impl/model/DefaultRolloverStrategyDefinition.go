package model

type DefaultRolloverStrategyDefinition struct {
	Max    *int              `yaml:"max"`
	Delete *DeleteDefinition `yaml:"delete"`
}
