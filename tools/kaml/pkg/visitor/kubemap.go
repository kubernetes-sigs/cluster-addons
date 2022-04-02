package visitor

import (
	"fmt"

	"sigs.k8s.io/kustomize/kyaml/yaml"
)

type KubeMap struct {
	node *yaml.Node
}

func (m *KubeMap) Node() *yaml.Node {
	return m.node
}

func (m *KubeMap) getNode(k string) (*yaml.Node, bool, error) {
	node := m.node

	n := len(node.Content)
	if n%2 != 0 {
		return nil, false, fmt.Errorf("unexpected content length in MappingNode")
	}

	for i := 0; i < n; i += 2 {
		kNode := node.Content[i]
		ks, ok := AsString(kNode)
		if !ok {
			continue
		}

		if ks != k {
			continue
		}
		vNode := node.Content[i+1]
		return vNode, true, nil
	}

	return nil, false, nil
}

func (m *KubeMap) GetStringField(k string) (string, bool, error) {
	vNode, found, err := m.getNode(k)
	if !found || err != nil {
		return "", found, err
	}

	vs, ok := AsString(vNode)
	if !ok {
		return "", true, fmt.Errorf("field was not of type string")
	}
	return vs, true, nil
}

func (m *KubeMap) GetMapField(k string) (*KubeMap, bool, error) {
	vNode, found, err := m.getNode(k)
	if !found || err != nil {
		return nil, found, err
	}

	v, err := AsKubeMap(vNode)
	if err != nil {
		return nil, true, err
	}
	return v, true, nil
}

func (m *KubeMap) ExtractStringFields() (map[string]string, error) {
	node := m.node

	n := len(node.Content)
	if n%2 != 0 {
		return nil, fmt.Errorf("unexpected content length in MappingNode")
	}

	vals := make(map[string]string)

	for i := 0; i < n; i += 2 {
		kNode := node.Content[i]
		ks, ok := AsString(kNode)
		if !ok {
			continue
		}

		vNode := node.Content[i+1]
		vs, ok := AsString(vNode)
		if !ok {
			continue
		}
		vals[ks] = vs
	}

	return vals, nil
}

func (m *KubeMap) Set(k string, newValue string) error {
	node := m.node

	n := len(node.Content)
	if n%2 != 0 {
		return fmt.Errorf("unexpected content length in MappingNode")
	}

	for i := 0; i < n; i += 2 {
		kNode := node.Content[i]
		vNode := node.Content[i+1]
		ks, ok := AsString(kNode)
		if !ok {
			continue
		}

		if ks == k {
			if err := setNode(vNode, newValue); err != nil {
				return err
			}
			return nil
		}
	}

	var newContent []*yaml.Node
	for _, v := range node.Content {
		newContent = append(newContent, v)
	}

	newContent = append(newContent, newStringNode(k), newStringNode(newValue))
	m.node.Content = newContent
	return nil
}

func setNode(node *yaml.Node, value string) error {
	*node = yaml.Node{}
	node.Kind = yaml.ScalarNode
	node.Value = value
	return nil
}

func newStringNode(value string) *yaml.Node {
	node := &yaml.Node{
		Kind:  yaml.ScalarNode,
		Value: value,
	}
	return node
}

func AsKubeMap(v *yaml.Node) (*KubeMap, error) {
	switch v.Kind {
	case yaml.MappingNode:
		return &KubeMap{node: v}, nil
	default:
		return nil, fmt.Errorf("unexpected kind for Map, expected MappingNode, got %v", v.Kind)
	}
}
