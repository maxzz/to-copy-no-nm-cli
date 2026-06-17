package progress

import "testing"

func TestPulseMatrixSpinnerCyclesFrames(t *testing.T) {
	spinner := NewPulseMatrixSpinner()

	for i := 0; i < len(PulseMatrixFrames)*2; i++ {
		want := PulseMatrixFrames[i%len(PulseMatrixFrames)]
		if got := spinner.Next(); got != want {
			t.Fatalf("frame %d: got %q, want %q", i, got, want)
		}
	}
}
