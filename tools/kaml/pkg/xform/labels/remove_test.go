package labels

import (
	"bytes"
	"context"
	"io/ioutil"
	"path/filepath"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	"sigs.k8s.io/cluster-addons/tools/kaml/pkg/xform"
	"sigs.k8s.io/kustomize/kyaml/fn/framework"
	"sigs.k8s.io/kustomize/kyaml/kio"
	"sigs.k8s.io/yaml"
)

func TestRemoveLabels(t *testing.T) {
	baseDir := "testdata/remove-label"
	dirs, err := ioutil.ReadDir(baseDir)
	if err != nil {
		t.Fatalf("failed to list %s: %v", baseDir, err)
	}

	for _, dir := range dirs {
		p := filepath.Join(baseDir, dir.Name())
		if !dir.IsDir() {
			t.Logf("skipping non-directory %v", p)
			continue
		}

		t.Run(dir.Name(), func(t *testing.T) {
			runTest(t, p, &RemoveLabel{})
		})
	}
}

func runTest(t *testing.T, p string, fn xform.Runnable) {
	commandPath := filepath.Join(p, "command.yaml")
	commandBytes, err := ioutil.ReadFile(commandPath)
	if err != nil {
		t.Fatalf("failed to read %q: %v", commandPath, err)
	}

	if err := yaml.Unmarshal(commandBytes, fn); err != nil {
		t.Fatalf("failed to parse %q: %v", commandPath, err)
	}

	inputPath := filepath.Join(p, "input.yaml")
	inputBytes, err := ioutil.ReadFile(inputPath)
	if err != nil {
		t.Fatalf("failed to read %q: %v", inputPath, err)
	}

	expectedPath := filepath.Join(p, "expected.yaml")
	expectedBytes, err := ioutil.ReadFile(expectedPath)
	if err != nil {
		t.Fatalf("failed to read %q: %v", expectedPath, err)
	}

	var out bytes.Buffer
	io := kio.ByteReadWriter{
		Reader: bytes.NewReader(inputBytes),
		Writer: &out,
	}

	items, err := io.Read()
	if err != nil {
		t.Fatalf("failed to parse yaml: %v", err)
	}

	resourceList := &framework.ResourceList{
		Items: items,
	}

	ctx := context.Background()
	if err := fn.Run(ctx, resourceList); err != nil {
		t.Fatalf("failed to run xform: %v", err)
	}

	if err := io.Write(resourceList.Items); err != nil {
		t.Fatalf("failed to write output: %v", err)
	}

	actual := strings.TrimSpace(string(out.String()))
	expected := strings.TrimSpace(string(expectedBytes))

	if actual == expected {
		return
	}

	t.Errorf("actual output did not match expected for %q", expectedPath)
	t.Logf("diff:\n%s", cmp.Diff(expected, actual))
}
