//go:build windows

package recycle

import (
	"fmt"
	"io/fs"
	"path/filepath"
)

const nodeModulesDirName = "node_modules" // formaer skipDirName
const gitDirName = ".git"

// ClearOptions controls how ClearDirectory removes destination contents.
type ClearOptions struct {
	// RemoveNodeModules deletes node_modules folders (including nested ones).
	// Default: false (left untouched).
	RemoveNodeModules bool
	// CopyGit deletes the root .git folder in the destination before copying.
	// Default: false (left untouched).
	CopyGit bool
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

	return clearEntries(dir, entries, opts, true)
}

func clearEntries(dir string, entries []fs.DirEntry, opts ClearOptions, isRoot bool) error {
	for _, entry := range entries {
		target := filepath.Join(dir, entry.Name())

		if entry.IsDir() {
			if entry.Name() == nodeModulesDirName {
				if opts.RemoveNodeModules {
					if err := MoveToRecycleBin(target); err != nil {
						return fmt.Errorf("recycle bin: clear %q: %w", target, err)
					}
				}
				continue
			}

			if isRoot && entry.Name() == gitDirName {
				if opts.CopyGit {
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
				if err := clearEntries(target, childEntries, opts, false); err != nil {
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
	entries, err := osReadDir(root)
	if err != nil {
		return false, err
	}

	for _, entry := range entries {
		if entry.IsDir() && entry.Name() == nodeModulesDirName {
			return true, nil
		}
	}

	for _, entry := range entries {
		if !entry.IsDir() || entry.Name() == nodeModulesDirName {
			continue
		}
		found, err := containsNodeModules(filepath.Join(root, entry.Name()))
		if err != nil {
			return false, err
		}
		if found {
			return true, nil
		}
	}

	return false, nil
}
