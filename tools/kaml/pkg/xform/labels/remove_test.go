package labels

import (
	"io/ioutil"
	"path/filepath"
	"testing"

	"sigs.k8s.io/cluster-addons/tools/kaml/pkg/testutils"
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
			testutils.RunGoldenTest(t, p, &RemoveLabel{})
		})
	}
}
