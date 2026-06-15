//go:build windows

package recycle

import (
	"io/fs"
	"os"
)

func osReadDir(name string) ([]fs.DirEntry, error) { return os.ReadDir(name) }
func osIsNotExist(err error) bool                  { return os.IsNotExist(err) }
func osMkdirAll(path string, perm fs.FileMode) error {
	return os.MkdirAll(path, perm)
}
