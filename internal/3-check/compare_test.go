package checkdir

import (
	"os"
	"path/filepath"
	"testing"

	"copy-no-nm/internal/9-progress"
)

func TestCompareMatchingTrees(t *testing.T) {
	src := t.TempDir()
	dst := t.TempDir()

	mirrorFile(t, filepath.Join(src, "a.txt"), filepath.Join(dst, "a.txt"), "hello")

	result, err := Compare(src, dst, nil)
	if err != nil {
		t.Fatal(err)
	}
	if result.SourceFileCount != 1 {
		t.Fatalf("expected 1 file, got %d", result.SourceFileCount)
	}
	if len(result.Changes) != 0 {
		t.Fatalf("expected no changes, got %v", result.Changes)
	}
}

func TestCompareSkipsNodeModulesAndGit(t *testing.T) {
	src := t.TempDir()
	dst := t.TempDir()

	mirrorFile(t, filepath.Join(src, "keep.txt"), filepath.Join(dst, "keep.txt"), "ok")
	writeFile(t, filepath.Join(src, "node_modules", "pkg", "index.js"), "skip")
	writeFile(t, filepath.Join(dst, "node_modules", "pkg", "index.js"), "different")
	writeFile(t, filepath.Join(src, ".git", "HEAD"), "skip")
	writeFile(t, filepath.Join(dst, ".git", "HEAD"), "different")

	result, err := Compare(src, dst, nil)
	if err != nil {
		t.Fatal(err)
	}
	if result.SourceFileCount != 1 {
		t.Fatalf("expected 1 file, got %d", result.SourceFileCount)
	}
	if len(result.Changes) != 0 {
		t.Fatalf("expected no changes, got %v", result.Changes)
	}
}

func TestCompareDetectsAddFile(t *testing.T) {
	src := t.TempDir()
	dst := t.TempDir()

	writeFile(t, filepath.Join(src, "only-src.txt"), "x")

	result, err := Compare(src, dst, nil)
	if err != nil {
		t.Fatal(err)
	}
	if len(result.Changes) != 1 {
		t.Fatalf("expected 1 change, got %v", result.Changes)
	}
	if result.Changes[0].Marker != progress.MarkerAdd || result.Changes[0].RelPath != "only-src.txt" {
		t.Fatalf("unexpected change: %+v", result.Changes[0])
	}
}

func TestCompareDetectsModifiedFile(t *testing.T) {
	src := t.TempDir()
	dst := t.TempDir()

	writeFile(t, filepath.Join(src, "a.txt"), "source")
	writeFile(t, filepath.Join(dst, "a.txt"), "dest")

	result, err := Compare(src, dst, nil)
	if err != nil {
		t.Fatal(err)
	}
	if len(result.Changes) != 1 {
		t.Fatalf("expected 1 change, got %v", result.Changes)
	}
	if result.Changes[0].Marker != progress.MarkerModify {
		t.Fatalf("expected modified change, got %+v", result.Changes[0])
	}
}

func TestCompareDetectsDeletedFile(t *testing.T) {
	src := t.TempDir()
	dst := t.TempDir()

	writeFile(t, filepath.Join(dst, "only-dst.txt"), "x")

	result, err := Compare(src, dst, nil)
	if err != nil {
		t.Fatal(err)
	}
	if len(result.Changes) != 1 {
		t.Fatalf("expected 1 change, got %v", result.Changes)
	}
	if result.Changes[0].Marker != progress.MarkerDelete || result.Changes[0].RelPath != "only-dst.txt" {
		t.Fatalf("unexpected change: %+v", result.Changes[0])
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
