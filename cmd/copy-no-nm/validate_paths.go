//go:build windows

package main

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

var errUsage = errors.New("usage")

func resolveAndValidatePaths(args []string, destinationMustExist bool) (src, dst string, err error) {
	if len(args) != 2 {
		return "", "", errUsage
	}

	src, err = filepath.Abs(args[0])
	if err != nil {
		return "", "", fmt.Errorf("invalid source path: %w", err)
	}

	dst, err = filepath.Abs(args[1])
	if err != nil {
		return "", "", fmt.Errorf("invalid destination path: %w", err)
	}

	if samePath(src, dst) {
		return "", "", fmt.Errorf("source and destination cannot be the same")
	}

	if err := requireDirectory("source", src, true); err != nil {
		return "", "", err
	}

	if destinationMustExist {
		if err := requireDirectory("destination", dst, true); err != nil {
			return "", "", err
		}
	} else if err := requireOrCreateDirectory(dst); err != nil {
		return "", "", err
	}

	return src, dst, nil
}

func requireOrCreateDirectory(path string) error {
	info, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			if err := os.MkdirAll(path, 0o755); err != nil {
				return fmt.Errorf("create destination: %w", err)
			}
			return nil
		}
		return fmt.Errorf("destination: %w", err)
	}

	if !info.IsDir() {
		return fmt.Errorf("destination is not a directory: %s", path)
	}

	return nil
}

func requireDirectory(label, path string, mustExist bool) error {
	info, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			if mustExist {
				return fmt.Errorf("%s does not exist: %s", label, path)
			}
			return nil
		}
		return fmt.Errorf("%s: %w", label, err)
	}

	if !info.IsDir() {
		return fmt.Errorf("%s is not a directory: %s", label, path)
	}

	return nil
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
