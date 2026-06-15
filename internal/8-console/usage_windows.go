//go:build windows

package console

import (
	"fmt"
	"os"
	"strings"
)

// UsageOption describes one CLI flag.
type UsageOption struct {
	Flag        string
	Description string
}

// UsageHelp is the content shown for incorrect user input.
type UsageHelp struct {
	Message string
	Syntax  string
	Options []UsageOption
}

// PrintUsage shows a friendly yellow message, dim gray syntax and options, then waits for a key.
func PrintUsage(help UsageHelp) {
	fmt.Printf("%s%s%s\n\n", colorYellow, help.Message, colorReset)

	fmt.Printf("%sUsage:%s\n", colorDim, colorReset)
	fmt.Printf("  %s%s%s\n\n", colorDim, help.Syntax, colorReset)

	if len(help.Options) > 0 {
		fmt.Printf("%sOptions:%s\n", colorDim, colorReset)
		maxFlagLen := 0
		for _, opt := range help.Options {
			if len(opt.Flag) > maxFlagLen {
				maxFlagLen = len(opt.Flag)
			}
		}
		for _, opt := range help.Options {
			padding := strings.Repeat(" ", maxFlagLen-len(opt.Flag))
			fmt.Printf("  %s%s%s %s%s\n", colorDim, opt.Flag, padding, opt.Description, colorReset)
		}
		fmt.Println()
	}

	fmt.Print("Press any key to close...")
	waitForKey()
	os.Exit(1)
}
