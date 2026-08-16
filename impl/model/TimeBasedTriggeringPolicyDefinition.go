package model

type TimeBasedTriggeringPolicyDefinition struct {
	Interval *int  `yaml:"interval"`
	Modulate *bool `yaml:"modulate"`
}
