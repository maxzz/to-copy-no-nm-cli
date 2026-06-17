//go:build windows

package fullcopy

import (
	recycle "copy-no-nm/internal/1-recycle"
	copydir "copy-no-nm/internal/2-copydir"
	"copy-no-nm/internal/9-progress"
)

// Options controls full-copy behaviour.
type Options struct {
	// CopyGit copies the root .git folder from the source and clears destination .git.
	CopyGit bool
	// Reporter receives copy progress. Default: none.
	Reporter progress.Reporter
}

// Run clears the destination via the Recycle Bin, then copies every file from the source.
func Run(src, dst string, opts Options) error {
	if err := recycle.ClearDirectory(dst, recycle.ClearOptions{CopyGit: opts.CopyGit}); err != nil {
		return err
	}
	return copydir.Copy(src, dst, copydir.CopyOptions{
		CopyGit:  opts.CopyGit,
		Reporter: opts.Reporter,
	})
}
