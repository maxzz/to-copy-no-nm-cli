//go:build windows

package recycle

import (
	"os"
	"path/filepath"
	"testing"
)

func TestContainsNodeModules(t *testing.T) {
	root := t.TempDir()

	if found, err := containsNodeModules(root); err != nil || found {
		t.Fatalf("empty dir: found=%v err=%v", found, err)
	}

	mustMkdir(t, filepath.Join(root, "src"))
	if found, err := containsNodeModules(root); err != nil || found {
		t.Fatalf("src only: found=%v err=%v", found, err)
	}

	mustMkdir(t, filepath.Join(root, "pkg", "node_modules"))
	if found, err := containsNodeModules(root); err != nil || !found {
		t.Fatalf("nested node_modules: found=%v err=%v", found, err)
	}
}

func mustMkdir(t *testing.T, path string) {
	t.Helper()
	if err := os.MkdirAll(path, 0o755); err != nil {
		t.Fatal(err)
	}
}
