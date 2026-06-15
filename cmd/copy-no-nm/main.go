//go:build windows

package main

import (
	"fmt"
	"os"
	"path/filepath"

	"copy-no-nm/internal/ascii"
	"copy-no-nm/internal/console"
	"copy-no-nm/internal/copydir"
	"copy-no-nm/internal/recycle"
)

func main() {
	if len(os.Args) != 3 {
		console.PrintError(fmt.Errorf("usage: copy-no-nm <source> <destination>"))
	}

	src, err := filepath.Abs(os.Args[1])
	if err != nil {
		console.PrintError(fmt.Errorf("invalid source path: %w", err))
	}

	dst, err := filepath.Abs(os.Args[2])
	if err != nil {
		console.PrintError(fmt.Errorf("invalid destination path: %w", err))
	}

	if err := recycle.ClearDirectory(dst); err != nil {
		console.PrintError(fmt.Errorf("clear destination: %w", err))
	}

	if err := copydir.Copy(src, dst); err != nil {
		console.PrintError(fmt.Errorf("copy failed: %w", err))
	}

	console.PrintSuccess(ascii.InspectorGadget())
}
