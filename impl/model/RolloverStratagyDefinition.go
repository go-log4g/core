package model

type RolloverStrategyDefinition struct {
	Type string `yaml:"type"`
	Max  *int   `yaml:"max"`
}
