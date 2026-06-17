//go:build windows

package syncdir

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestSyncCopiesNewAndChangedFiles(t *testing.T) {
	src := t.TempDir()
	dst := t.TempDir()

	writeFile(t, filepath.Join(src, "keep.txt"), "same")
	writeFile(t, filepath.Join(dst, "keep.txt"), "same")
	writeFile(t, filepath.Join(src, "new.txt"), "new")
	writeFile(t, filepath.Join(dst, "old.txt"), "remove me")
	writeFile(t, filepath.Join(dst, "change.txt"), "old")

	srcChangePath := filepath.Join(src, "change.txt")
	writeFile(t, srcChangePath, "new")
	time.Sleep(10 * time.Millisecond)
	if err := os.Chtimes(srcChangePath, time.Now(), time.Now()); err != nil {
		t.Fatal(err)
	}

	mirrorFile(t, filepath.Join(src, "keep.txt"), filepath.Join(dst, "keep.txt"), "same")

	if err := Sync(src, dst, SyncOptions{}); err != nil {
		t.Fatal(err)
	}

	assertFileContent(t, filepath.Join(dst, "new.txt"), "new")
	assertFileContent(t, filepath.Join(dst, "change.txt"), "new")
	assertNoFile(t, filepath.Join(dst, "old.txt"))

	count, err := compareFileCount(src, dst)
	if err != nil {
		t.Fatal(err)
	}
	if count != 3 {
		t.Fatalf("expected 3 synced files, got %d", count)
	}
}

func TestSyncSkipsUnchangedFiles(t *testing.T) {
	src := t.TempDir()
	dst := t.TempDir()

	mirrorFile(t, filepath.Join(src, "a.txt"), filepath.Join(dst, "a.txt"), "hello")

	before, err := os.Stat(filepath.Join(dst, "a.txt"))
	if err != nil {
		t.Fatal(err)
	}

	if err := Sync(src, dst, SyncOptions{}); err != nil {
		t.Fatal(err)
	}

	after, err := os.Stat(filepath.Join(dst, "a.txt"))
	if err != nil {
		t.Fatal(err)
	}
	if !before.ModTime().Equal(after.ModTime()) {
		t.Fatal("unchanged file should not be recopied")
	}
}

func TestSyncPreservesNodeModules(t *testing.T) {
	src := t.TempDir()
	dst := t.TempDir()

	writeFile(t, filepath.Join(src, "app.js"), "src")
	writeFile(t, filepath.Join(dst, "app.js"), "src")
	writeFile(t, filepath.Join(dst, "node_modules", "pkg", "index.js"), "keep me")

	if err := Sync(src, dst, SyncOptions{}); err != nil {
		t.Fatal(err)
	}

	assertFileContent(t, filepath.Join(dst, "node_modules", "pkg", "index.js"), "keep me")
}

func compareFileCount(src, dst string) (int, error) {
	srcFiles, _, err := collectTree(src, SyncOptions{})
	if err != nil {
		return 0, err
	}
	dstFiles, _, err := collectTree(dst, SyncOptions{})
	if err != nil {
		return 0, err
	}
	if len(srcFiles) != len(dstFiles) {
		return 0, fmt.Errorf("file count mismatch: source=%d destination=%d", len(srcFiles), len(dstFiles))
	}
	for rel, srcSig := range srcFiles {
		dstSig, ok := dstFiles[rel]
		if !ok || !signaturesEqual(srcSig, dstSig) {
			return 0, fmt.Errorf("mismatch for %s", rel)
		}
	}
	return len(srcFiles), nil
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

func assertFileContent(t *testing.T, path, want string) {
	t.Helper()
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read %s: %v", path, err)
	}
	if string(data) != want {
		t.Fatalf("file %s: got %q, want %q", path, string(data), want)
	}
}

func assertNoFile(t *testing.T, path string) {
	t.Helper()
	if _, err := os.Stat(path); !os.IsNotExist(err) {
		t.Fatalf("expected %s to be removed, err=%v", path, err)
	}
}
