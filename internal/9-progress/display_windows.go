//go:build windows

package progress

import (
	"fmt"
	"io"
	"os"
	"sync"
	"time"

	console "copy-no-nm/internal/8-console"
)

const spinnerPad = " "

// FolderDisplay prints one line per folder with a pulse-matrix spinner on the active line.
type FolderDisplay struct {
	out io.Writer

	mu sync.Mutex

	folderName string
	fileCount  int

	animStop chan struct{}
	animDone chan struct{}

	totalFiles   int
	folderCount  int
}

// NewFolderDisplay writes progress lines to stdout.
func NewFolderDisplay() *FolderDisplay {
	return &FolderDisplay{out: os.Stdout}
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

// Finish prints totals, prompts for a key press, and waits for any key.
func (d *FolderDisplay) Finish() {
	d.mu.Lock()
	d.finalizeActiveFolderLocked()
	totalFiles := d.totalFiles
	totalFolders := d.folderCount
	d.mu.Unlock()

	fmt.Fprintf(d.out, "\nTotal: %d files in %d folders\n\n", totalFiles, totalFolders)
	fmt.Fprint(d.out, "Press any key to close the window")
	console.WaitForAnyKey()
}

func (d *FolderDisplay) finalizeActiveFolderLocked() {
	d.stopAnimationLocked()

	if d.folderName == "" {
		return
	}

	fmt.Fprint(d.out, d.formatLine(spinnerPad, d.folderName, d.fileCount)+"\n")
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
				fmt.Fprint(d.out, "\r"+line)
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
	return fmt.Sprintf("%s %s  %d files", spinnerChar, folder, count)
}
