package model

type AppenderDefinition struct {
	Type           string           `yaml:"type"`
	Target         string           `yaml:"target"`
	File           string           `yaml:"file"`
	FilePattern    string           `yaml:"filePattern"`
	Append         *bool            `yaml:"append"`
	BufferSize     *int             `yaml:"bufferSize"`
	ImmediateFlush *bool            `yaml:"immediateFlush"`
	Layout         LayoutDefinition `yaml:"layout"`
	Filter         FilterDefinition `yaml:"filter"`

	Policies                PoliciesDefinition                `yaml:"policies"`
	DefaultRolloverStrategy DefaultRolloverStrategyDefinition `yaml:"defaultRolloverStrategy"`
}
