package model

type PoliciesDefinition struct {
	SizeBasedTriggeringPolicy *SizeBasedTriggeringPolicyDefinition `yaml:"sizeBasedTriggeringPolicy"`
	TimeBasedTriggeringPolicy *TimeBasedTriggeringPolicyDefinition `yaml:"timeBasedTriggeringPolicy"`
	OnStartupTriggeringPolicy *OnStartupTriggeringPolicyDefinition `yaml:"onStartupTriggeringPolicy"`
}
