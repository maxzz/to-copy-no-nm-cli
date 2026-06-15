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
	gadget := ascii.InspectorGadget()

	if len(os.Args) != 3 {
		console.PrintError(fmt.Errorf("usage: copy-no-nm <source> <destination>"), gadget)
	}

	src, err := filepath.Abs(os.Args[1])
	if err != nil {
		console.PrintError(fmt.Errorf("invalid source path: %w", err), gadget)
	}

	dst, err := filepath.Abs(os.Args[2])
	if err != nil {
		console.PrintError(fmt.Errorf("invalid destination path: %w", err), gadget)
	}

	if err := recycle.ClearDirectory(dst); err != nil {
		console.PrintError(fmt.Errorf("clear destination: %w", err), gadget)
	}

	if err := copydir.Copy(src, dst); err != nil {
		console.PrintError(fmt.Errorf("copy failed: %w", err), gadget)
	}

	console.PrintSuccess(gadget)
}

//TODO: better error message instead of "Error: usage: copy-no-nm <source> <destination>"
//TODO: "Press any key to close..." accept only Enter but should be any
//TODO: better icon
//TODO: publish to npm
