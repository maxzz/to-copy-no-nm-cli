//go:build windows

package console

import (
	"fmt"
	"os"
	"strings"
)

const usageWrapWidth = 80

// UsageOption describes one CLI flag.
type UsageOption struct {
	Flag        string
	Description string
}

// UsageArg labels one token from the command line.
type UsageArg struct {
	Label string
	Value string
}

// UsageHelp is the content shown for incorrect user input.
type UsageHelp struct {
	Message string
	Syntax  string
	Options []UsageOption
	Args    []UsageArg
}

// PrintUsage shows a friendly yellow message, gray syntax, dim gray options, then waits for a key.
func PrintUsage(help UsageHelp) {
	fmt.Printf("%sUsage:%s\n", colorGray, colorReset)
	printWrappedLines(colorGray, "  ", help.Syntax, usageWrapWidth)
	fmt.Println()

	if len(help.Options) > 0 {
		fmt.Printf("%sOptions:%s\n", colorGray, colorReset)
		maxFlagLen := 0
		for _, opt := range help.Options {
			if len(opt.Flag) > maxFlagLen {
				maxFlagLen = len(opt.Flag)
			}
		}
		optionIndent := 2 + maxFlagLen + 1
		for _, opt := range help.Options {
			padding := strings.Repeat(" ", maxFlagLen-len(opt.Flag))
			prefix := "  " + opt.Flag + padding + " "
			printWrappedLines(colorGray, prefix, opt.Description, usageWrapWidth, strings.Repeat(" ", optionIndent))
		}
		fmt.Println()
	}

	if len(help.Args) > 0 {
		fmt.Printf("%sArguments:%s\n", colorGray, colorReset)
		for _, arg := range help.Args {
			fmt.Printf("  %s%s:%s %s\n", colorGray, arg.Label, colorReset, arg.Value)
		}
		fmt.Println()
	}

	printWrappedLines(colorYellow, "", help.Message, usageWrapWidth)
	fmt.Println()

	fmt.Print("Press any key to close...")
	waitForKey()
	os.Exit(1)
}

func printWrappedLines(color, firstIndent, text string, width int, continuationIndent ...string) {
	lines := wrapText(text, width-len(firstIndent))
	if len(lines) == 0 {
		return
	}

	contIndent := firstIndent
	if len(continuationIndent) > 0 {
		contIndent = continuationIndent[0]
	}

	fmt.Printf("%s%s%s%s\n", firstIndent, color, lines[0], colorReset)
	for _, line := range lines[1:] {
		fmt.Printf("%s%s%s%s\n", contIndent, color, line, colorReset)
	}
}

func wrapText(text string, width int) []string {
	text = strings.TrimSpace(text)
	if text == "" {
		return nil
	}
	if width < 1 {
		return []string{text}
	}

	words := strings.Fields(text)
	var lines []string
	var line strings.Builder

	flush := func() {
		if line.Len() > 0 {
			lines = append(lines, line.String())
			line.Reset()
		}
	}

	for _, word := range words {
		if line.Len() == 0 {
			if len(word) <= width {
				line.WriteString(word)
			} else {
				lines = append(lines, word)
			}
			continue
		}

		if line.Len()+1+len(word) <= width {
			line.WriteByte(' ')
			line.WriteString(word)
			continue
		}

		flush()
		if len(word) <= width {
			line.WriteString(word)
		} else {
			lines = append(lines, word)
		}
	}

	flush()
	return lines
}
