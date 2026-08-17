package model

type DeleteDefinition struct {
	MaxAge       string `yaml:"maxAge"`
	MaxFiles     *int   `yaml:"maxFiles"`
	MaxTotalSize string `yaml:"maxTotalSize"`
}