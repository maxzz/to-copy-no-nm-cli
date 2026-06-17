//go:build windows

package progress

import (
	"fmt"
	"io"
	"os"
	"sort"
	"strings"
	"sync"
	"time"

	console "copy-no-nm/internal/8-console"
)

const (
	spinnerPad     = " "
	changeIndent   = "        "
	colorRed       = "\x1b[31m"
	colorReset     = "\x1b[0m"
	changeLegend   = "U = new in source, M = modified, D = extra in destination"
)

type folderRow struct {
	name  string
	count int
}

// FolderDisplay prints one line per folder with a pulse-matrix spinner on the active line.
type FolderDisplay struct {
	out io.Writer

	mu sync.Mutex

	sourceRootLabel string
	scanRootLabel   string

	folderName string
	fileCount  int

	animStop chan struct{}
	animDone chan struct{}

	completedRows      []folderRow
	sourceFolderCounts map[string]int

	totalFiles  int
	folderCount int
}

// NewFolderDisplay writes progress lines to stdout.
func NewFolderDisplay() *FolderDisplay {
	return &FolderDisplay{
		out:                os.Stdout,
		sourceFolderCounts: make(map[string]int),
	}
}

// SetSourceRootLabel marks which scan root supplies folder counts for the change report.
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

func (d *FolderDisplay) BeginFolder(folderPath string) {
	d.mu.Lock()
	defer d.mu.Unlock()

	d.finalizeActiveFolderLocked()
	d.folderName = folderPath
	d.fileCount = 0
	d.startAnimationLocked()
}

func (d *FolderDisplay) UpdateFileCount(count int) {
	d.mu.Lock()
	d.fileCount = count
	d.mu.Unlock()
}

func (d *FolderDisplay) CompleteFolder() {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.finalizeActiveFolderLocked()
}

// Finish prints comparison differences, totals, a legend, and waits for any key.
func (d *FolderDisplay) Finish(changes []ChangeEntry, srcRootLabel string) {
	d.mu.Lock()
	d.finalizeActiveFolderLocked()
	totalFiles := d.totalFiles
	totalFolders := d.folderCount
	d.mu.Unlock()

	if len(changes) > 0 {
		d.printChangeReport(changes, srcRootLabel)
	}

	fmt.Fprintf(d.out, "\nTotal: %d files in %d folders\n", totalFiles, totalFolders)
	if len(changes) > 0 {
		fmt.Fprintf(d.out, "%s\n", changeLegend)
	}
	fmt.Fprint(d.out, "\nPress any key to close the window")
	console.WaitForAnyKey()
}

func (d *FolderDisplay) printChangeReport(changes []ChangeEntry, srcRootLabel string) {
	byFolder := GroupChangesByBucket(changes, srcRootLabel)

	folders := make([]string, 0, len(byFolder))
	for folder := range byFolder {
		folders = append(folders, folder)
	}
	sort.Strings(folders)

	fmt.Fprint(d.out, "\n")
	for _, folder := range folders {
		entries := byFolder[folder]
		sort.Slice(entries, func(i, j int) bool {
			if entries[i].Marker != entries[j].Marker {
				return entries[i].Marker < entries[j].Marker
			}
			return entries[i].RelPath < entries[j].RelPath
		})

		count := d.folderCountForReport(folder)
		fmt.Fprint(d.out, colorRed+d.formatLine(spinnerPad, folder, count)+colorReset+"\n")
		for _, entry := range entries {
			fmt.Fprintf(d.out, "%s%c  %s\n", changeIndent, entry.Marker, entry.RelPath)
		}
	}
}

func (d *FolderDisplay) folderCountForReport(folder string) int {
	d.mu.Lock()
	count, ok := d.sourceFolderCounts[folder]
	d.mu.Unlock()
	if ok {
		return count
	}
	return 0
}

func (d *FolderDisplay) finalizeActiveFolderLocked() {
	d.stopAnimationLocked()

	if d.folderName == "" {
		return
	}

	line := d.formatLine(spinnerPad, d.folderName, d.fileCount)
	fmt.Fprint(d.out, line+"\n")
	d.completedRows = append(d.completedRows, folderRow{name: d.folderName, count: d.fileCount})
	if d.scanRootLabel == d.sourceRootLabel {
		d.sourceFolderCounts[d.folderName] = d.fileCount
	}
	d.totalFiles += d.fileCount
	d.folderCount++
	d.folderName = ""
	d.fileCount = 0
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
				name := d.folderName
				count := d.fileCount
				d.mu.Unlock()

				if name == "" {
					continue
				}

				line := d.formatLine(spinner.Next(), name, count)
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
}

func (d *FolderDisplay) formatLine(spinnerChar, folder string, count int) string {
	return fmt.Sprintf("%s %3d %s", spinnerChar, count, folder)
}
