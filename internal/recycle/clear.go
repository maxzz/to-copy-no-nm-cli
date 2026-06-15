//go:build windows

package recycle

import (
	"fmt"
	"io/fs"
	"path/filepath"
)

const skipDirName = "node_modules"

// ClearOptions controls how ClearDirectory removes destination contents.
type ClearOptions struct {
	// RemoveNodeModules deletes node_modules folders (including nested ones).
	// By default they are left untouched.
	RemoveNodeModules bool
}

// ClearDirectory clears destination contents before copying.
// Subfolders without any nested node_modules are removed as a whole.
// Subfolders that contain node_modules are processed with the same rules as the root.
func ClearDirectory(dir string, opts ClearOptions) error {
	entries, err := osReadDir(dir)
	if err != nil {
		if osIsNotExist(err) {
			return osMkdirAll(dir, 0o755)
		}
		return err
	}

	return clearEntries(dir, entries, opts)
}

func clearEntries(dir string, entries []fs.DirEntry, opts ClearOptions) error {
	for _, entry := range entries {
		target := filepath.Join(dir, entry.Name())

		if entry.IsDir() {
			if entry.Name() == skipDirName {
				if opts.RemoveNodeModules {
					if err := MoveToRecycleBin(target); err != nil {
						return fmt.Errorf("recycle bin: clear %q: %w", target, err)
					}
				}
				continue
			}

			hasNodeModules, err := containsNodeModules(target)
			if err != nil {
				return fmt.Errorf("scan %q: %w", target, err)
			}

			if hasNodeModules {
				childEntries, err := osReadDir(target)
				if err != nil {
					return fmt.Errorf("read %q: %w", target, err)
				}
				if err := clearEntries(target, childEntries, opts); err != nil {
					return err
				}
				continue
			}

			if err := MoveToRecycleBin(target); err != nil {
				return fmt.Errorf("recycle bin: clear %q: %w", target, err)
			}
			continue
		}

		if err := MoveToRecycleBin(target); err != nil {
			return fmt.Errorf("recycle bin: clear %q: %w", target, err)
		}
	}

	return nil
}

func containsNodeModules(root string) (bool, error) {
	found := false
	err := filepath.WalkDir(root, func(path string, entry fs.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		if entry.IsDir() && entry.Name() == skipDirName {
			found = true
			return fs.SkipAll
		}
		return nil
	})
	return found, err
}
