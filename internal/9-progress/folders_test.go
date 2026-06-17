package progress

import (
	"path/filepath"
	"testing"
)

func TestBucketFolder(t *testing.T) {
	tests := []struct {
		rel  string
		want string
	}{
		{"readme.md", "pmac"},
		{`src\main.ts`, "pmac\\src"},
		{`src\components\ui\btn.tsx`, "pmac\\src"},
	}

	for _, tc := range tests {
		got := BucketFolder("pmac", tc.rel)
		if got != tc.want {
			t.Errorf("BucketFolder(%q): got %q, want %q", tc.rel, got, tc.want)
		}
	}
}

func TestBuildTreeReport(t *testing.T) {
	dirCounts := map[string]int{
		".vscode":                          1,
		"packages":                         203,
		filepath.Join("packages", "shared-types"): 1,
		filepath.Join("packages", "template"):     1,
		filepath.Join("packages", "utility"):      1,
	}

	changes := []ChangeEntry{
		{Marker: 'M', RelPath: "README.md"},
		{Marker: 'U', RelPath: "pnpm-workspace.yaml"},
		{Marker: 'M', RelPath: `packages\utility\some other filename`},
	}

	report := BuildTreeReport(dirCounts, changes)

	if len(report.FirstLevel) != 2 {
		t.Fatalf("expected 2 first-level dirs, got %d", len(report.FirstLevel))
	}
	if report.FirstLevel[0].Name != ".vscode" || report.FirstLevel[0].FileCount != 1 {
		t.Fatalf("unexpected .vscode node: %+v", report.FirstLevel[0])
	}

	packages := report.FirstLevel[1]
	if packages.Name != "packages" || packages.FileCount != 203 {
		t.Fatalf("unexpected packages node: %+v", packages)
	}
	if len(packages.Children) != 3 {
		t.Fatalf("expected 3 package children, got %d", len(packages.Children))
	}

	utility := packages.Children[2]
	if utility.Name != "utility" || len(utility.Changes) != 1 {
		t.Fatalf("unexpected utility node: %+v", utility)
	}
	if utility.Changes[0].RelPath != "some other filename" {
		t.Fatalf("utility change path: got %q", utility.Changes[0].RelPath)
	}

	if len(report.RootChanges) != 2 {
		t.Fatalf("expected 2 root changes, got %d", len(report.RootChanges))
	}
}

func TestRecordSubtreeCounts(t *testing.T) {
	counts := map[string]int{}
	RecordSubtreeCounts(counts, `packages\utility\foo.ts`)
	if counts["packages"] != 1 || counts[filepath.Join("packages", "utility")] != 1 {
		t.Fatalf("unexpected counts: %#v", counts)
	}
}
