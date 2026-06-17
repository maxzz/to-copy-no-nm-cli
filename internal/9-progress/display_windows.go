//go:build windows

package progress

import (
	"fmt"
	"io"
	"os"
	"strings"
	"sync"
	"time"

	console "copy-no-nm/internal/8-console"
)

const (
	colorGray    = "\x1b[90m"
	colorYellow  = "\x1b[33m"
	colorGreen   = "\x1b[32m"
	colorRed     = "\x1b[31m"
	colorReset   = "\x1b[0m"
	changeLegend = "U = new in source, M = modified, D = extra in destination"
)

// FolderDisplay collects scan stats and prints a tree report when finished.
type FolderDisplay struct {
	out io.Writer

	mu sync.Mutex

	sourceRootLabel string
	scanRootLabel   string

	animFolder string
	animCount  int

	animStop chan struct{}
	animDone chan struct{}

	dirCounts  map[string]int
	totalFiles int
}

// NewFolderDisplay writes progress lines to stdout.
func NewFolderDisplay() *FolderDisplay {
	return &FolderDisplay{
		out:       os.Stdout,
		dirCounts: make(map[string]int),
	}
}

// SetSourceRootLabel marks which scan root supplies directory statistics.
func (d *FolderDisplay) SetSourceRootLabel(label string) {
	d.mu.Lock()
	d.sourceRootLabel = label
	d.mu.Unlock()
}

func (d *FolderDisplay) BeginScan(rootLabel string) {
	d.mu.Lock()
	d.scanRootLabel = rootLabel
	d.mu.Unlock()
}

func (d *FolderDisplay) RecordFile(relPath string) {
	d.mu.Lock()
	defer d.mu.Unlock()

	if d.scanRootLabel != d.sourceRootLabel {
		return
	}

	d.totalFiles++
	RecordSubtreeCounts(d.dirCounts, relPath)

	d.animFolder = AnimationLabel(relPath)
	d.animCount++
	if d.animStop == nil {
		d.startAnimationLocked()
	}
}

// Finish prints the tree report, totals, a legend, and waits for any key.
func (d *FolderDisplay) Finish(changes []ChangeEntry, _ string) {
	d.mu.Lock()
	d.stopAnimationLocked()
	totalFiles := d.totalFiles
	dirCounts := copyDirCounts(d.dirCounts)
	d.mu.Unlock()

	fmt.Fprint(d.out, "\n")
	d.printTreeReport(dirCounts, changes)

	fmt.Fprintf(d.out, "\nTotal: %d files in %d folders\n", totalFiles, CountTrackedFolders(dirCounts))
	if len(changes) > 0 {
		fmt.Fprintf(d.out, "%s\n", changeLegend)
	}
	fmt.Fprint(d.out, "\nPress any key to close the window")
	console.WaitForAnyKey()
}

func (d *FolderDisplay) printTreeReport(dirCounts map[string]int, changes []ChangeEntry) {
	report := BuildTreeReport(dirCounts, changes)
	hasRootFiles := len(report.RootChanges) > 0

	for i, node := range report.FirstLevel {
		d.printTopFolderLine(node.Name, node.FileCount)

		moreAfter := hasRootFiles || i < len(report.FirstLevel)-1
		if len(node.Children) > 0 {
			d.printSecondLevelBlock(node, moreAfter)
		} else if len(node.Changes) > 0 {
			d.printFileChangesAtDepth("", node.Changes, !moreAfter)
		}
	}

	if hasRootFiles {
		d.printFileChangesAtDepth("", report.RootChanges, true)
	}
}

func (d *FolderDisplay) printTopFolderLine(name string, count int) {
	fmt.Fprintf(d.out, "%s %s(%d files in all subfolders)%s\n", name, colorGray, count, colorReset)
}

func (d *FolderDisplay) printPrefixedFolderLine(prefix, name string, count int) {
	fmt.Fprintf(d.out, "%s%s %s(%d files in all subfolders)%s\n", prefix, name, colorGray, count, colorReset)
}

func (d *FolderDisplay) printSecondLevelBlock(node TreeNode, moreAfter bool) {
	children := node.Children

	for i, child := range children {
		isLastChild := i == len(children)-1
		branch := "├──"
		cont := "│   "
		if isLastChild {
			branch = "└──"
			cont = "    "
		}

		d.printPrefixedFolderLine(branch, child.Name, child.FileCount)
		d.printFileChangesAtDepth(cont, child.Changes, true)
	}

	if len(node.Changes) > 0 {
		d.printFileChangesAtDepth("", node.Changes, !moreAfter)
	}
}

func (d *FolderDisplay) printFileChangesAtDepth(cont string, changes []ChangeEntry, blockEnds bool) {
	for i, change := range changes {
		isLast := i == len(changes)-1
		branch := "├──"
		if isLast && blockEnds {
			branch = "└──"
		}
		fmt.Fprint(d.out, cont+branch+fileChangeText(change)+"\n")
	}
}

func fileChangeText(change ChangeEntry) string {
	color := colorForMarker(change.Marker)
	return fmt.Sprintf("File: %s%c %s%s", color, change.Marker, change.RelPath, colorReset)
}

func colorForMarker(marker rune) string {
	switch marker {
	case 'U':
		return colorGreen
	case 'M':
		return colorYellow
	case 'D':
		return colorRed
	default:
		return ""
	}
}

func copyDirCounts(src map[string]int) map[string]int {
	dst := make(map[string]int, len(src))
	for k, v := range src {
		dst[k] = v
	}
	return dst
}

func (d *FolderDisplay) startAnimationLocked() {
	d.animStop = make(chan struct{})
	d.animDone = make(chan struct{})

	stop := d.animStop
	done := d.animDone

	go func() {
		defer close(done)

		ticker := time.NewTicker(80 * time.Millisecond)
		defer ticker.Stop()

		spinner := NewPulseMatrixSpinner()
		for {
			select {
			case <-stop:
				return
			case <-ticker.C:
				d.mu.Lock()
				folder := d.animFolder
				count := d.animCount
				d.mu.Unlock()

				if folder == "" {
					continue
				}

				line := fmt.Sprintf("%s %3d %s", spinner.Next(), count, folder)
				d.mu.Lock()
				fmt.Fprint(d.out, "\r"+line+strings.Repeat(" ", 8))
				d.mu.Unlock()
			}
		}
	}()
}

func (d *FolderDisplay) stopAnimationLocked() {
	if d.animStop == nil {
		return
	}

	close(d.animStop)
	done := d.animDone
	d.animStop = nil
	d.animDone = nil

	d.mu.Unlock()
	<-done
	d.mu.Lock()

	fmt.Fprint(d.out, "\r"+strings.Repeat(" ", 40)+"\r")
}
