//go:build windows

package main

import (
	"fmt"
	"os"
	"path/filepath"

	checkdir "copy-no-nm/internal/3-check"
	syncdir "copy-no-nm/internal/4-syncdir"
	fullcopy "copy-no-nm/internal/5-fullcopy"
	progress "copy-no-nm/internal/9-progress"
	console "copy-no-nm/internal/8-console"
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

	srcLabel := filepath.Base(src)

	if parsed.options.check {
		display := progress.NewFolderDisplay()
		display.SetSourceRootLabel(srcLabel)
		result, err := checkdir.Compare(src, dst, display)
		if err != nil {
			console.PrintError(fmt.Errorf("check failed: %w", err))
		}
		display.Finish(result.Changes, srcLabel, progress.OperationCheck)
		os.Exit(0)
	}

	display := progress.NewFolderDisplay()
	display.SetSourceRootLabel(srcLabel)

	operation := progress.OperationSync
	if parsed.options.fullCopy {
		operation = progress.OperationFullCopy
		if err := fullcopy.Run(src, dst, fullcopy.Options{
			CopyGit:  parsed.options.copyGit,
			Reporter: display,
		}); err != nil {
			console.PrintError(fmt.Errorf("full copy failed: %w", err))
		}
	} else {
		if err := syncdir.Sync(src, dst, syncdir.SyncOptions{
			CopyGit:  parsed.options.copyGit,
			Reporter: display,
		}); err != nil {
			console.PrintError(fmt.Errorf("sync failed: %w", err))
		}
	}

	display.Finish(nil, srcLabel, operation)
	os.Exit(0)
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
			Flag: "-f, --full",
			Description: "Copy every file after clearing the destination via the Recycle Bin " +
				"(default: off; sync mode copies only new or changed files and removes extras)",
		},
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
