package convert

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

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

		t.Run(f.Name(), func(t *testing.T) {
			ctx := context.Background()

			b, err := ioutil.ReadFile(p)
			if err != nil {
				t.Fatalf("error reading file %s: %v", p, err)
			}

			opt := BuildRoleOptions{
				Name:      "generated-role",
				Namespace: "kube-system",
			}
			actualObjects, err := BuildRole(ctx, string(b), opt)
			if err != nil {
				t.Fatalf("error building role %s: %v", p, err)
			}

			actualYAML, err := ToYAML(actualObjects)
			if err != nil {
				t.Fatalf("error converting to YAML %s: %v", p, err)
			}

			expectedPath := strings.Replace(p, "in.yaml", "out.yaml", -1)

			CheckGoldenFile(t, expectedPath, string(actualYAML))
		})
	}
}

func CheckGoldenFile(t *testing.T, expectedPath, actual string) {
	t.Helper()

	var expected string
	{
		b, err := ioutil.ReadFile(expectedPath)
		if err != nil {
			t.Fatalf("error reading file %s: %v", expectedPath, err)
		}
		expected = string(b)
	}

	delta, err := runDiffCommand(actual, expected)
	if err != nil {
		t.Logf("failed to run system diff, falling back to string diff: %v", err)
		delta = diff.StringDiff(actual, expected)
	}

	if delta == "" {
		return
	}

	t.Errorf("Diff: expected - + actual\n%s", delta)
	t.Errorf("unexpected diff between actual and expected YAML. See previous output for details.")

	if os.Getenv("HACK_AUTOFIX_EXPECTED_OUTPUT") != "" {
		if err := os.WriteFile(expectedPath, []byte(actual), 0644); err != nil {
			t.Errorf("failed to write expected file %q: %w", expectedPath, err)
		}
	} else {
		t.Logf(`To regenerate the output based on this result, rerun this test with HACK_AUTOFIX_EXPECTED_OUTPUT=1`)
	}
}

func runDiffCommand(actual, expected string) (string, error) {
	tempDir, err := os.MkdirTemp("", "")
	if err != nil {
		return "", fmt.Errorf("failed to create temp dir: %w", err)
	}

	defer os.RemoveAll(tempDir)

	actualPath := filepath.Join(tempDir, "actual.yaml")
	if err := ioutil.WriteFile(actualPath, []byte(actual), 0644); err != nil {
		return "", fmt.Errorf("failed to write file %q: %w", actualPath, err)
	}

	expectedPath := filepath.Join(tempDir, "expected.yaml")
	if err := ioutil.WriteFile(expectedPath, []byte(expected), 0644); err != nil {
		return "", fmt.Errorf("failed to write file %q: %w", expectedPath, err)
	}

	// pls to use unified diffs, kthxbai?
	cmd := exec.Command("diff", "-u", expectedPath, actualPath)
	var stdout bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		// exit code 1 means there was a diff
		if cmd.ProcessState.ExitCode() == 1 {
			return stdout.String(), nil
		}

		return "", fmt.Errorf("failed to run diff: %w", err)
	}

	return stdout.String(), nil
}
