package progress

// PulseMatrixFrames is the rotating braille sequence used for activity feedback.
var PulseMatrixFrames = []string{"⠂", "⠒", "⠖", "⠶", "⡶", "⣶", "⣶", "⣾", "⣿"}

// Spinner cycles through a fixed frame sequence.
type Spinner struct {
	frames []string
	index  int
}

// NewPulseMatrixSpinner returns a spinner using PulseMatrixFrames.
func NewPulseMatrixSpinner() *Spinner {
	frames := make([]string, len(PulseMatrixFrames))
	copy(frames, PulseMatrixFrames)
	return &Spinner{frames: frames}
}

// Next returns the current frame and advances to the next one.
func (s *Spinner) Next() string {
	frame := s.frames[s.index]
	s.index = (s.index + 1) % len(s.frames)
	return frame
}
