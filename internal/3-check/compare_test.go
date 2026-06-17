package checkdir

import (
	"os"
	"path/filepath"
	"testing"
)

func TestCompareMatchingTrees(t *testing.T) {
	src := t.TempDir()
	dst := t.TempDir()

	mirrorFile(t, filepath.Join(src, "a.txt"), filepath.Join(dst, "a.txt"), "hello")

	count, err := Compare(src, dst)
	if err != nil {
		t.Fatal(err)
	}
	if count != 1 {
		t.Fatalf("expected 1 file, got %d", count)
	}
}

func TestCompareSkipsNodeModulesAndGit(t *testing.T) {
	src := t.TempDir()
	dst := t.TempDir()

	writeFile(t, filepath.Join(src, "keep.txt"), "ok")
	writeFile(t, filepath.Join(dst, "keep.txt"), "ok")
	writeFile(t, filepath.Join(src, "node_modules", "pkg", "index.js"), "skip")
	writeFile(t, filepath.Join(dst, "node_modules", "pkg", "index.js"), "different")
	writeFile(t, filepath.Join(src, ".git", "HEAD"), "skip")
	writeFile(t, filepath.Join(dst, ".git", "HEAD"), "different")

	count, err := Compare(src, dst)
	if err != nil {
		t.Fatal(err)
	}
	if count != 1 {
		t.Fatalf("expected 1 file, got %d", count)
	}
}

func TestCompareDetectsMissingFile(t *testing.T) {
	src := t.TempDir()
	dst := t.TempDir()

	writeFile(t, filepath.Join(src, "only-src.txt"), "x")

	_, err := Compare(src, dst)
	if err == nil {
		t.Fatal("expected error for missing destination file")
	}
}

func writeFile(t *testing.T, path, content string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
}

func mirrorFile(t *testing.T, srcPath, dstPath, content string) {
	t.Helper()
	writeFile(t, srcPath, content)
	info, err := os.Stat(srcPath)
	if err != nil {
		t.Fatal(err)
	}
	writeFile(t, dstPath, content)
	if err := os.Chtimes(dstPath, info.ModTime(), info.ModTime()); err != nil {
		t.Fatal(err)
	}
}
