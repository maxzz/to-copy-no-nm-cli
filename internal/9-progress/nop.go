package progress

// NopReporter discards progress updates. Use in tests or when no UI is needed.
type NopReporter struct{}

func (NopReporter) BeginFolder(string)    {}
func (NopReporter) UpdateFileCount(int)   {}
func (NopReporter) CompleteFolder()       {}
