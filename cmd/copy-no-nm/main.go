//go:build windows

package main

import (
	"fmt"
	"os"

	ascii "copy-no-nm/internal/8-result-ascii"
	"copy-no-nm/internal/console"
	"copy-no-nm/internal/copydir"
	"copy-no-nm/internal/recycle"
)

func main() {
	gadget := ascii.InspectorGadget()

	src, dst, err := resolveAndValidatePaths(os.Args)
	if err != nil {
		console.PrintError(err, gadget)
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
