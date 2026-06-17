package progress

import (
	"testing"
)

func TestBucketFolder(t *testing.T) {
	tests := []struct {
		rel    string
		want   string
	}{
		{"readme.md", "pmac"},
		{`src\main.ts`, "pmac\\src"},
		{`src\components\ui\btn.tsx`, "pmac\\src"},
		{`lib\util.go`, "pmac\\lib"},
	}

	for _, tc := range tests {
		got := BucketFolder("pmac", tc.rel)
		if got != tc.want {
			t.Errorf("BucketFolder(%q): got %q, want %q", tc.rel, got, tc.want)
		}
	}
}

func TestChangeDisplayPath(t *testing.T) {
	tests := []struct {
		rel  string
		want string
	}{
		{"readme.md", "readme.md"},
		{`src\main.ts`, "main.ts"},
		{`src\components\ui\btn.tsx`, `components\ui\btn.tsx`},
	}

	for _, tc := range tests {
		got := ChangeDisplayPath(tc.rel)
		if got != tc.want {
			t.Errorf("ChangeDisplayPath(%q): got %q, want %q", tc.rel, got, tc.want)
		}
	}
}

func TestGroupChangesByBucketDeduplicatesRootFiles(t *testing.T) {
	changes := []ChangeEntry{
		{Marker: 'U', RelPath: "new.txt"},
		{Marker: 'M', RelPath: `src\a.ts`},
		{Marker: 'U', RelPath: `src\deep\b.ts`},
	}

	byFolder := GroupChangesByBucket(changes, "pmac")

	if len(byFolder["pmac"]) != 1 {
		t.Fatalf("root bucket: got %d entries, want 1", len(byFolder["pmac"]))
	}
	if byFolder["pmac"][0].RelPath != "new.txt" {
		t.Fatalf("root entry: got %q", byFolder["pmac"][0].RelPath)
	}

	src := byFolder["pmac\\src"]
	if len(src) != 2 {
		t.Fatalf("src bucket: got %d entries, want 2", len(src))
	}
}
