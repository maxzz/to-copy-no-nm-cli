package progress

// ChangeEntry describes one file difference found during comparison.
type ChangeEntry struct {
	Marker  rune
	RelPath string
}
