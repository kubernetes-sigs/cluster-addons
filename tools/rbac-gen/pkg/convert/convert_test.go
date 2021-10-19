package convert

import (
	"io/ioutil"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"golang.org/x/xerrors"

	"k8s.io/apimachinery/pkg/util/diff"
)

func TestRBACGen(t *testing.T) {
	testdir := "tests"

	files, err := ioutil.ReadDir(testdir)
	if err != nil {
		t.Fatalf("error reading dir %s: %v", testdir, err)
	}

	for _, f := range files {
		p := filepath.Join(testdir, f.Name())
		t.Logf("Filepath: %s", p)
		if f.IsDir() {
			t.Errorf("unexpected directory in tests directory: %s", p)
			continue
		}

		if strings.HasSuffix(p, "~") {
			// Ignore editor temp files (for sanity)
			t.Logf("ignoring editor temp file %s", p)
			continue
		}

		if !strings.HasSuffix(p, ".in.yaml") {
			if !strings.HasSuffix(p, ".out.yaml") {
				t.Errorf("unexpected file in tests directory: %s", p)
			}
			continue
		}

		b, err := ioutil.ReadFile(p)
		if err != nil {
			t.Errorf("error reading file %s: %v", p, err)
			continue
		}

		opt := BuildRoleOptions{
			Name:      "generated-role",
			Namespace: "kube-system",
		}
		actualYAML, err := ParseYAMLtoRole(string(b), opt)
		if err != nil {
			t.Errorf("error parsing YAML %s: %v", p, err)
			continue
		}

		expectedPath := strings.Replace(p, "in.yaml", "out.yaml", -1)
		var expectedYAML string

		{
			b, err := ioutil.ReadFile(expectedPath)
			if err != nil {
				t.Errorf("error reading file %s: %v", expectedPath, err)
				continue
			}
			expectedYAML = string(b)
		}

		if expectedYAML != actualYAML {
			if err := diffFiles(t, expectedPath, actualYAML); err != nil {
				t.Logf("failed to run system diff, falling back to string diff: %v", err)
				t.Logf("diff: %s", diff.StringDiff(actualYAML, expectedYAML))
			}

			t.Errorf("unexpected diff between actual and expected YAML. See previous output for details.")
			// TODO: Do we want to replace the out.yaml if HACK_AUTOFIX_EXPECTED_OUTPUT="true" is set?
			// t.Logf(`To regenerate the output based on this result,
			// rerun this test with HACK_AUTOFIX_EXPECTED_OUTPUT="true"`)
		}
	}
}

// Had to copy out this func from kubebuider-declarative-pattern. Could we simply export it?
func diffFiles(t *testing.T, expectedPath, actual string) error {
	t.Helper()
	writeTmp := func(content string) (string, error) {
		tmp, err := ioutil.TempFile("", "*.yaml")
		if err != nil {
			return "", err
		}
		defer func() {
			tmp.Close()
		}()
		if _, err := tmp.Write([]byte(content)); err != nil {
			return "", err
		}
		return tmp.Name(), nil
	}

	actualTmp, err := writeTmp(actual)
	if err != nil {
		return xerrors.Errorf("write actual yaml to temp file failed: %w", err)
	}
	t.Logf("Wrote actual to %s", actualTmp)

	// pls to use unified diffs, kthxbai?
	cmd := exec.Command("diff", "-u", expectedPath, actualTmp)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return xerrors.Errorf("set up stdout pipe from diff failed: %w", err)
	}

	if err := cmd.Start(); err != nil {
		return xerrors.Errorf("start command failed: %w", err)
	}

	diff, err := ioutil.ReadAll(stdout)
	if err != nil {
		return xerrors.Errorf("read from diff stdout failed: %w", err)
	}

	if err := cmd.Wait(); err != nil {
		exitErr, ok := err.(*exec.ExitError)
		if !ok {
			return xerrors.Errorf("wait for command to finish failed: %w", err)
		}
		t.Logf("Diff exited %s", exitErr)
	}

	expectedAbs, err := filepath.Abs(expectedPath)
	if err != nil {
		t.Logf("getting absolute path for %s failed: %s", expectedPath, err)
		expectedAbs = expectedPath
	}

	t.Logf("View diff: meld %s %s", expectedAbs, actualTmp)
	t.Logf("Diff: expected - + actual\n%s", diff)
	return nil
}
