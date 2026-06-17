package progress

// Reporter receives folder-level progress updates during long-running operations.
type Reporter interface {
	BeginFolder(folderPath string)
	UpdateFileCount(count int)
	CompleteFolder()
}
