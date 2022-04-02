package visitor

import "fmt"

type KubeObject struct {
	obj *KubeMap
}

func (o *KubeObject) GetKind() string {
	s, _, _ := o.obj.GetStringField("kind")
	return s
}

func (o *KubeObject) GetAPIVersion() string {
	s, _, _ := o.obj.GetStringField("apiVersion")
	return s
}

func (o *KubeObject) GetName() string {
	metadata := o.GetMetadata()
	if metadata == nil {
		return ""
	}
	s, _, _ := metadata.GetStringField("name")
	return s
}

func (o *KubeObject) GetMetadata() *KubeMap {
	m, _, _ := o.obj.GetMapField("metadata")
	return m
}

func (o *KubeObject) SetName(name string) error {
	metadata := o.GetMetadata()
	if metadata == nil {
		return fmt.Errorf("metadata not found")
	}
	if err := metadata.Set("name", name); err != nil {
		return err
	}
	return nil
}
