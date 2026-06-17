//go:build windows

package main

import (
	"fmt"
	"os"

	recycle "copy-no-nm/internal/1-recycle"
	copydir "copy-no-nm/internal/2-copydir"
	checkdir "copy-no-nm/internal/3-check"
	console "copy-no-nm/internal/8-console"
	ascii "copy-no-nm/internal/8-result-ascii"
)

var version = "dev"

func main() {
	parsed := parseCLI(os.Args[1:])
	if parsed.help {
		printUsage()
	}

	console.PrintVersion(version)

	if len(parsed.unknown) > 0 || len(parsed.positionals) != 2 {
		printUsageMessage("Please provide a source folder and a destination folder.", parsed.args)
	}

	src, dst, err := resolveAndValidatePaths(parsed.positionals, parsed.options.check, parsed.options.reversePaths)
	if err != nil {
		console.PrintError(err)
	}

	if parsed.options.check {
		fileCount, err := checkdir.Compare(src, dst)
		if err != nil {
			console.PrintError(fmt.Errorf("check failed: %w", err))
		}
		console.PrintCheckSuccess(fileCount)
	}

	if err := recycle.ClearDirectory(dst, recycle.ClearOptions{
		CopyGit: parsed.options.copyGit,
	}); err != nil {
		console.PrintError(fmt.Errorf("clear destination: %w", err))
	}

	if err := copydir.Copy(src, dst, copydir.CopyOptions{
		CopyGit: parsed.options.copyGit,
	}); err != nil {
		console.PrintError(fmt.Errorf("copy failed: %w", err))
	}

	console.PrintSuccess(ascii.BuildOK())
}

func printUsage() {
	printUsageMessage("Copy a folder to another location while skipping node_modules during the copy.", nil)
}

func printUsageMessage(message string, args []console.UsageArg) {
	console.PrintUsage(console.UsageHelp{
		Message: message,
		Syntax:  "copy-no-nm [options] <source> <destination>",
		Options: usageOptions(),
		Args:    args,
	})
}

func usageOptions() []console.UsageOption {
	return []console.UsageOption{
		{
			Flag: "-r, --reverse",
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
