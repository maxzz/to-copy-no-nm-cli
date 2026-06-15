//go:build windows

package main

import (
	"flag"
	"fmt"

	ascii "copy-no-nm/internal/8-result-ascii"
	"copy-no-nm/internal/8-console"
	"copy-no-nm/internal/2-copydir"
	"copy-no-nm/internal/1-recycle"
)

func main() {
	removeNodeModules := flag.Bool(
		"remove-node-modules",
		false,
		"also delete node_modules folders (including nested) in the destination",
	)
	flag.Parse()

	gadget := ascii.InspectorGadget()

	src, dst, err := resolveAndValidatePaths(flag.Args())
	if err != nil {
		console.PrintError(err, gadget)
	}

	if err := recycle.ClearDirectory(dst, recycle.ClearOptions{
		RemoveNodeModules: *removeNodeModules,
	}); err != nil {
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
