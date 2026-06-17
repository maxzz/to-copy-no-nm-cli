package progress

import (
	"path/filepath"
	"strings"
)

// BucketFolder maps a file path to its display folder: the root label or root/first-level subfolder.
func BucketFolder(rootLabel, relPath string) string {
	dir := filepath.Dir(relPath)
	if dir == "." {
		return rootLabel
	}
	first := strings.Split(filepath.ToSlash(dir), "/")[0]
	return filepath.Join(rootLabel, first)
}

// ChangeDisplayPath returns the file path shown under a bucket (relative to root or first-level subfolder).
func ChangeDisplayPath(relPath string) string {
	parts := strings.Split(filepath.ToSlash(relPath), "/")
	if len(parts) == 1 {
		return parts[0]
	}
	return filepath.Join(parts[1:]...)
}

// GroupChangesByBucket groups change entries by root or first-level folder buckets.
func GroupChangesByBucket(changes []ChangeEntry, rootLabel string) map[string][]ChangeEntry {
	byFolder := make(map[string][]ChangeEntry)
	for _, change := range changes {
		folder := BucketFolder(rootLabel, change.RelPath)
		byFolder[folder] = append(byFolder[folder], ChangeEntry{
			Marker:  change.Marker,
			RelPath: ChangeDisplayPath(change.RelPath),
		})
	}
	return byFolder
}
