//go:build windows

package main

import (
	"errors"
	"flag"
	"fmt"

	recycle "copy-no-nm/internal/1-recycle"
	copydir "copy-no-nm/internal/2-copydir"
	checkdir "copy-no-nm/internal/3-check"
	console "copy-no-nm/internal/8-console"
	ascii "copy-no-nm/internal/8-result-ascii"
)

var version = "dev"

func main() {
	var copyGit bool
	var check bool
	var swapPaths bool
	flag.BoolVar(&copyGit, "copy-git", false, "copy the .git folder from the source root and clear the destination .git folder")
	flag.BoolVar(&copyGit, "g", false, "shorthand for --copy-git")
	flag.BoolVar(&check, "check", false, "verify source and destination match by file size and modification time")
	flag.BoolVar(&check, "c", false, "shorthand for --check")
	flag.BoolVar(&swapPaths, "swap", false, "treat the first argument as destination and the second as source")
	flag.BoolVar(&swapPaths, "s", false, "shorthand for --swap")
	flag.Usage = printUsage
	flag.Parse()

	console.PrintVersion(version)

	gadget := ascii.InspectorGadget()

	src, dst, err := resolveAndValidatePaths(flag.Args(), check, swapPaths)
	if errors.Is(err, errUsage) {
		printUsageMessage("Please provide a source folder and a destination folder.")
	}

	if err != nil {
		console.PrintError(err, gadget)
	}

	if check {
		fileCount, err := checkdir.Compare(src, dst)
		if err != nil {
			console.PrintError(fmt.Errorf("check failed: %w", err), gadget)
		}
		console.PrintCheckSuccess(fileCount)
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

	console.PrintSuccess(ascii.BuildOK())
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
			Flag: "-s, --swap",
			Description: "Treat the first argument as destination and the second as source " +
				"(default: off; normal order is source then destination)",
		},
		{
			Flag: "-c, --check",
			Description: "Verify source and destination are identical using file size and modification time " +
				"(default: off; excludes node_modules and .git; does not copy or clear)",
		},
		{
			Flag: "-g, --copy-git",
			Description: "Copy the .git folder from the source root and clear the destination .git folder " +
				"(default: off; source .git is not copied, destination .git is kept)",
		},
	}
}
