package convert

import (
	"bytes"

	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/yaml"
)

// ToYAML convert objects to a YAML multidoc string
func ToYAML(objects []runtime.Object) ([]byte, error) {
	var buf bytes.Buffer

	for i, obj := range objects {
		if i != 0 {
			buf.WriteString("\n---\n\n")
		}

		b, err := yaml.Marshal(obj)
		if err != nil {
			return nil, err
		}
		buf.Write(b)
	}

	return buf.Bytes(), nil
}
