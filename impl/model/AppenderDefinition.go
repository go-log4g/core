package model

type AppenderDefinition struct {
	Type           string           `yaml:"type"`
	Target         string           `yaml:"target"`
	File           string           `yaml:"file"`
	Append         *bool            `yaml:"append"`
	BufferSize     *int             `yaml:"bufferSize"`
	ImmediateFlush *bool            `yaml:"immediateFlush"`
	Layout         LayoutDefinition `yaml:"layout"`
	Filter         FilterDefinition `yaml:"filter"`
}
