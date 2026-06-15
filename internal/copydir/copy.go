package copydir

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

const skipDirName = "node_modules"

// Copy copies src into dst, skipping node_modules directories and preserving file metadata.
func Copy(src, dst string) error {
	src = filepath.Clean(src)
	dst = filepath.Clean(dst)

	srcInfo, err := os.Lstat(src)
	if err != nil {
		return fmt.Errorf("source: %w", err)
	}
	if !srcInfo.IsDir() {
		return fmt.Errorf("source must be a directory: %s", src)
	}

	if samePath(src, dst) {
		return fmt.Errorf("source and destination must differ")
	}

	if isInside(src, dst) || isInside(dst, src) {
		return fmt.Errorf("source and destination cannot contain each other")
	}

	if err := os.MkdirAll(dst, srcInfo.Mode().Perm()); err != nil {
		return fmt.Errorf("create destination: %w", err)
	}

	return filepath.WalkDir(src, func(path string, entry fs.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}

		rel, err := filepath.Rel(src, path)
		if err != nil {
			return err
		}
		if rel == "." {
			return nil
		}

		if entry.IsDir() && entry.Name() == skipDirName {
			return filepath.SkipDir
		}

		target := filepath.Join(dst, rel)

		info, err := entry.Info()
		if err != nil {
			return err
		}

		mode := info.Mode()

		if mode&os.ModeSymlink != 0 {
			link, err := os.Readlink(path)
			if err != nil {
				return fmt.Errorf("read symlink %q: %w", path, err)
			}
			_ = os.Remove(target)
			if err := os.Symlink(link, target); err != nil {
				return fmt.Errorf("create symlink %q: %w", target, err)
			}
			return copyFileMetadata(path, target)
		}

		if info.IsDir() {
			if err := os.MkdirAll(target, mode.Perm()); err != nil {
				return fmt.Errorf("create directory %q: %w", target, err)
			}
			return copyFileMetadata(path, target)
		}

		if err := copyPlatformFile(path, target); err != nil {
			return fmt.Errorf("copy file %q: %w", path, err)
		}
		return nil
	})
}

func samePath(a, b string) bool {
	a = filepath.Clean(a)
	b = filepath.Clean(b)
	if a == b {
		return true
	}
	aa, errA := filepath.Abs(a)
	bb, errB := filepath.Abs(b)
	return errA == nil && errB == nil && aa == bb
}

func isInside(base, target string) bool {
	baseAbs, err := filepath.Abs(base)
	if err != nil {
		return false
	}
	targetAbs, err := filepath.Abs(target)
	if err != nil {
		return false
	}
	rel, err := filepath.Rel(baseAbs, targetAbs)
	if err != nil {
		return false
	}
	return rel != "." && !strings.HasPrefix(rel, "..")
}
