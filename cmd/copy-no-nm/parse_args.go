//go:build windows

package main

import (
	console "copy-no-nm/internal/8-console"
	"strings"
)

type cliOptions struct {
	copyGit   bool
	check     bool
	reversePaths bool
}

type cliParseResult struct {
	options     cliOptions
	positionals []string
	unknown     []string
	args        []console.UsageArg
	help        bool
}

var knownBoolFlags = map[string]func(*cliOptions){
	"g":        func(o *cliOptions) { o.copyGit = true },
	"copy-git": func(o *cliOptions) { o.copyGit = true },
	"c":        func(o *cliOptions) { o.check = true },
	"check":    func(o *cliOptions) { o.check = true },
	"r":        func(o *cliOptions) { o.reversePaths = true },
	"reverse":  func(o *cliOptions) { o.reversePaths = true },
}

func parseCLI(args []string) cliParseResult {
	var result cliParseResult
	positionalIndex := 0

	for _, arg := range args {
		if !isOptionToken(arg) {
			label := positionalLabel(positionalIndex, result.options.reversePaths)
			result.positionals = append(result.positionals, arg)
			result.args = append(result.args, console.UsageArg{Label: label, Value: arg})
			positionalIndex++
			continue
		}

		name := optionName(arg)
		switch name {
		case "h", "help":
			result.help = true
			result.args = append(result.args, console.UsageArg{Label: "option", Value: arg})
		case "g", "copy-git", "c", "check", "r", "reverse":
			knownBoolFlags[name](&result.options)
			result.args = append(result.args, console.UsageArg{Label: "option", Value: arg})
		default:
			result.unknown = append(result.unknown, arg)
			result.args = append(result.args, console.UsageArg{Label: "unknown option", Value: arg})
		}
	}

	return result
}

func isOptionToken(arg string) bool {
	return strings.HasPrefix(arg, "-") && arg != "-"
}

func optionName(arg string) string {
	if strings.HasPrefix(arg, "--") {
		return strings.TrimPrefix(arg, "--")
	}
	return strings.TrimPrefix(arg, "-")
}

func positionalLabel(index int, reverse bool) string {
	switch index {
	case 0:
		if reverse {
			return "destination"
		}
		return "source"
	case 1:
		if reverse {
			return "source"
		}
		return "destination"
	default:
		return "extra argument"
	}
}
