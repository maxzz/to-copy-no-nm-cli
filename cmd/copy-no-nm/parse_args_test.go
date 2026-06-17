//go:build windows

package main

import "testing"

func TestParseCLI_flagsAfterPositionals(t *testing.T) {
	result := parseCLI([]string{"src", "dst", "-c"})
	if !result.options.check {
		t.Fatal("expected -c to enable check")
	}
	if len(result.positionals) != 2 || result.positionals[0] != "src" || result.positionals[1] != "dst" {
		t.Fatalf("unexpected positionals: %v", result.positionals)
	}
}

func TestParseCLI_flagsBetweenPositionals(t *testing.T) {
	result := parseCLI([]string{"src", "-r", "dst"})
	if !result.options.reversePaths {
		t.Fatal("expected -r to enable reverse")
	}
	if len(result.positionals) != 2 || result.positionals[0] != "src" || result.positionals[1] != "dst" {
		t.Fatalf("unexpected positionals: %v", result.positionals)
	}
}

func TestParseCLI_unknownOption(t *testing.T) {
	result := parseCLI([]string{"src", "-x", "dst"})
	if len(result.unknown) != 1 || result.unknown[0] != "-x" {
		t.Fatalf("unexpected unknown flags: %v", result.unknown)
	}
}

func TestParseCLI_fullCopyFlag(t *testing.T) {
	result := parseCLI([]string{"--full", "src", "dst"})
	if !result.options.fullCopy {
		t.Fatal("expected --full to enable full copy")
	}
}

func TestParseCLI_describesArgsInOrder(t *testing.T) {
	result := parseCLI([]string{"src", "-c", "dst"})
	if len(result.args) != 3 {
		t.Fatalf("expected 3 arg lines, got %d", len(result.args))
	}
	if result.args[0].Label != "source" || result.args[0].Value != "src" {
		t.Fatalf("unexpected first arg line: %+v", result.args[0])
	}
	if result.args[1].Label != "option" || result.args[1].Value != "-c" {
		t.Fatalf("unexpected second arg line: %+v", result.args[1])
	}
	if result.args[2].Label != "destination" || result.args[2].Value != "dst" {
		t.Fatalf("unexpected third arg line: %+v", result.args[2])
	}
}
