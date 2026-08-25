package model

type LoggerDefinition struct {
	Level        string                   `yaml:"level"`
	AppenderRefs []*AppenderRefDefinition `yaml:"appenderRefs"`
	Additive     *bool                    `yaml:"additive"`
}
