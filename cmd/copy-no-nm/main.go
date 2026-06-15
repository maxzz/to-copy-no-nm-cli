//go:build windows

package main

import (
	"errors"
	"flag"
	"fmt"

	ascii "copy-no-nm/internal/8-result-ascii"
	"copy-no-nm/internal/8-console"
	"copy-no-nm/internal/2-copydir"
	"copy-no-nm/internal/1-recycle"
)

var version = "dev"

func main() {
	var removeNodeModules bool
	flag.BoolVar(&removeNodeModules, "remove-node-modules", false, "also delete node_modules folders (including nested) in the destination")
	flag.BoolVar(&removeNodeModules, "r", false, "shorthand for --remove-node-modules")
	flag.Usage = printUsage
	flag.Parse()

	console.PrintVersion("copy-no-nm", version)

	gadget := ascii.InspectorGadget()

	src, dst, err := resolveAndValidatePaths(flag.Args())
	if errors.Is(err, errUsage) {
		printUsageMessage("Please provide a source folder and a destination folder.")
	}

	if err != nil {
		console.PrintError(err, gadget)
	}

	if err := recycle.ClearDirectory(dst, recycle.ClearOptions{
		RemoveNodeModules: removeNodeModules,
	}); err != nil {
		console.PrintError(fmt.Errorf("clear destination: %w", err), gadget)
	}

	if err := copydir.Copy(src, dst); err != nil {
		console.PrintError(fmt.Errorf("copy failed: %w", err), gadget)
	}

	console.PrintSuccess(gadget)
}

func printUsage() {
	printUsageMessage("Copy a folder to another location while skipping node_modules during the copy.")
}

func printUsageMessage(message string) {
	console.PrintUsage(console.UsageHelp{
		Message: message,
		Syntax:  "copy-no-nm [options] <source> <destination>",
		Options: []console.UsageOption{
			{
				Flag:        "-r, --remove-node-modules",
				Description: "Also delete node_modules folders (including nested) in the destination before copying",
			},
		},
	})
}

//TODO: better icon
//TODO: publish to npm
