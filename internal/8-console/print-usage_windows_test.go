//go:build windows

package console

import "testing"

func TestWrapText(t *testing.T) {
	lines := wrapText("one two three four five six seven eight", 20)
	if len(lines) < 2 {
		t.Fatalf("expected wrapping, got %v", lines)
	}
	for _, line := range lines {
		if len(line) > 20 {
			t.Fatalf("line exceeds width: %q (%d)", line, len(line))
		}
	}
}

func TestWrapTextEmpty(t *testing.T) {
	if lines := wrapText("   ", 80); lines != nil {
		t.Fatalf("expected nil, got %v", lines)
	}
}
