package checkdir

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"time"

	"copy-no-nm/internal/9-progress"
)

const (
	skipNodeModules = "node_modules"
	skipGit         = ".git"
)

type fileSignature struct {
	size    int64
	modTime time.Time
	symlink bool
	target  string
}

// Compare checks that files under src and dst match by size and modification time.
// Directories named node_modules or .git are excluded at any depth.
// Pass nil for reporter when no progress output is needed.
func Compare(src, dst string, reporter progress.Reporter) (int, error) {
	src = filepath.Clean(src)
	dst = filepath.Clean(dst)

	if reporter == nil {
		reporter = progress.NopReporter{}
	}

	srcFiles, err := collectFiles(src, reporter)
	if err != nil {
		return 0, fmt.Errorf("scan source: %w", err)
	}

	dstFiles, err := collectFiles(dst, reporter)
	if err != nil {
		return 0, fmt.Errorf("scan destination: %w", err)
	}

	for rel, srcSig := range srcFiles {
		dstSig, ok := dstFiles[rel]
		if !ok {
			return 0, fmt.Errorf("missing in destination: %s", rel)
		}
		if !signaturesEqual(srcSig, dstSig) {
			return 0, signatureMismatchError(rel, srcSig, dstSig)
		}
	}

	for rel := range dstFiles {
		if _, ok := srcFiles[rel]; !ok {
			return 0, fmt.Errorf("extra in destination: %s", rel)
		}
	}

	return len(srcFiles), nil
}

func collectFiles(root string, reporter progress.Reporter) (map[string]fileSignature, error) {
	files := make(map[string]fileSignature)
	rootLabel := filepath.Base(root)

	var currentFolder string
	var folderFileCount int

	completeFolder := func() {
		if currentFolder == "" {
			return
		}
		reporter.CompleteFolder()
		currentFolder = ""
		folderFileCount = 0
	}

	err := filepath.WalkDir(root, func(path string, entry fs.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}

		if entry.IsDir() {
			if entry.Name() == skipNodeModules || entry.Name() == skipGit {
				return filepath.SkipDir
			}
			return nil
		}

		rel, err := filepath.Rel(root, path)
		if err != nil {
			return err
		}

		folderLabel := folderDisplayName(rootLabel, filepath.Dir(rel))
		if folderLabel != currentFolder {
			completeFolder()
			currentFolder = folderLabel
			reporter.BeginFolder(currentFolder)
		}
		folderFileCount++
		reporter.UpdateFileCount(folderFileCount)

		info, err := entry.Info()
		if err != nil {
			return err
		}

		sig := fileSignature{
			size:    info.Size(),
			modTime: info.ModTime(),
		}

		if info.Mode()&os.ModeSymlink != 0 {
			target, err := os.Readlink(path)
			if err != nil {
				return fmt.Errorf("read symlink %q: %w", path, err)
			}
			sig.symlink = true
			sig.target = target
		}

		files[rel] = sig
		return nil
	})
	if err != nil {
		return nil, err
	}

	completeFolder()
	return files, nil
}

func folderDisplayName(rootLabel, relDir string) string {
	if relDir == "." {
		return rootLabel
	}
	return filepath.Join(rootLabel, relDir)
}

func signaturesEqual(a, b fileSignature) bool {
	if a.symlink || b.symlink {
		return a.symlink && b.symlink && a.target == b.target
	}
	if a.size != b.size {
		return false
	}
	return a.modTime.Equal(b.modTime)
}

func signatureMismatchError(rel string, src, dst fileSignature) error {
	if src.symlink || dst.symlink {
		return fmt.Errorf("different symlink %s: source %q, destination %q", rel, src.target, dst.target)
	}
	return fmt.Errorf(
		"different file %s: source size=%d mtime=%s, destination size=%d mtime=%s",
		rel,
		src.size,
		src.modTime.Format(time.RFC3339),
		dst.size,
		dst.modTime.Format(time.RFC3339),
	)
}
