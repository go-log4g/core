package model

import "gopkg.in/yaml.v3"

type DynamicDefinition struct {
	Node *yaml.Node
}

func (this *DynamicDefinition) UnmarshalYAML(node *yaml.Node) error {
	this.Node = node
	return nil
}
