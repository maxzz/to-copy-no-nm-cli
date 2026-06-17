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
	var copyGit bool
	flag.BoolVar(&copyGit, "copy-git", false, "copy the .git folder from the source root and clear the destination .git folder")
	flag.BoolVar(&copyGit, "g", false, "shorthand for --copy-git")
	flag.Usage = printUsage
	flag.Parse()

	console.PrintVersion(version)

	gadget := ascii.InspectorGadget()

	src, dst, err := resolveAndValidatePaths(flag.Args())
	if errors.Is(err, errUsage) {
		printUsageMessage("Please provide a source folder and a destination folder.")
	}

	if err != nil {
		console.PrintError(err, gadget)
	}

	if err := recycle.ClearDirectory(dst, recycle.ClearOptions{
		CopyGit: copyGit,
	}); err != nil {
		console.PrintError(fmt.Errorf("clear destination: %w", err), gadget)
	}

	if err := copydir.Copy(src, dst, copydir.CopyOptions{
		CopyGit: copyGit,
	}); err != nil {
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
		Options: usageOptions(),
	})
}

func usageOptions() []console.UsageOption {
	return []console.UsageOption{
		{
			Flag: "-g, --copy-git",
			Description: "Copy the .git folder from the source root and clear the destination .git folder " +
				"(default: off; source .git is not copied, destination .git is kept)",
		},
	}
}
