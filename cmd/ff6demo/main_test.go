package main

import (
	"strings"
	"testing"

	"github.com/bensabler/ff6-decompile/internal/platform/snespad"
)

func TestParseInputScript(t *testing.T) {
	src, err := parseInputScript("A@10,Right@20-22,Start@100")
	if err != nil {
		t.Fatalf("parseInputScript: %v", err)
	}
	tests := []struct {
		frame uint64
		want  snespad.State
	}{
		{0, 0},
		{9, 0},
		{10, snespad.State(0).With(snespad.A)},
		{11, 0},
		{19, 0},
		{20, snespad.State(0).With(snespad.Right)},
		{22, snespad.State(0).With(snespad.Right)},
		{23, 0},
		{100, snespad.State(0).With(snespad.Start)},
	}
	for _, tt := range tests {
		if got := src.Poll(tt.frame); got != tt.want {
			t.Errorf("frame %d = %v, want %v", tt.frame, got, tt.want)
		}
	}
}

func TestParseInputScriptCombinesOverlaps(t *testing.T) {
	src, err := parseInputScript("Right@0-10,A@5-6")
	if err != nil {
		t.Fatal(err)
	}
	want := snespad.State(0).With(snespad.Right).With(snespad.A)
	if got := src.Poll(5); got != want {
		t.Errorf("overlapping spans should combine: frame 5 = %v, want %v", got, want)
	}
}

func TestParseInputScriptIsCaseInsensitive(t *testing.T) {
	src, err := parseInputScript("start@1, LEFT@2")
	if err != nil {
		t.Fatal(err)
	}
	if !src.Poll(1).Held(snespad.Start) || !src.Poll(2).Held(snespad.Left) {
		t.Error("button names should be case-insensitive")
	}
}

func TestParseInputScriptRejectsBadInput(t *testing.T) {
	for _, s := range []string{
		"A",           // no frame
		"Nope@1",      // unknown button
		"A@x",         // non-numeric frame
		"A@10-x",      // non-numeric end
		"A@20-10",     // reversed range
		"A@1,Bogus@2", // one bad entry poisons the script
	} {
		if _, err := parseInputScript(s); err == nil {
			t.Errorf("parseInputScript(%q) should have failed", s)
		}
	}
}

func TestParseInputScriptIgnoresEmptyEntries(t *testing.T) {
	src, err := parseInputScript("A@1,,  ,B@2")
	if err != nil {
		t.Fatalf("empty entries should be skipped, not rejected: %v", err)
	}
	if !src.Poll(1).Held(snespad.A) || !src.Poll(2).Held(snespad.B) {
		t.Error("surviving entries should still apply")
	}
}

// TestSafeOutDirKeepsFramesLocal guards the asset boundary: rendered frames
// are ROM-derived images, so the demo must not scatter them across the
// filesystem by default.
func TestSafeOutDirKeepsFramesLocal(t *testing.T) {
	for _, dir := range []string{"local_artifacts", "local_artifacts/demo-frames", "local_artifacts/x/y"} {
		if _, err := safeOutDir(dir, false); err != nil {
			t.Errorf("safeOutDir(%q) = %v, want it allowed", dir, err)
		}
	}
	for _, dir := range []string{"/tmp/frames", "docs", "../escape", "local_artifacts/../docs", ""} {
		if _, err := safeOutDir(dir, false); err == nil {
			t.Errorf("safeOutDir(%q) should be refused without -force-out", dir)
		}
	}
	// The override exists, and says so in the error it replaces.
	if _, err := safeOutDir("/tmp/frames", true); err != nil {
		t.Errorf("-force-out should allow any path: %v", err)
	}
	_, err := safeOutDir("/tmp/frames", false)
	if err == nil || !strings.Contains(err.Error(), "-force-out") {
		t.Errorf("the refusal should name the override flag: %v", err)
	}
}

func TestRunWindowedIsUnavailableWithoutTheGuiTag(t *testing.T) {
	err := runWindowed(nil)
	if err == nil {
		t.Skip("built with the gui tag")
	}
	if !strings.Contains(err.Error(), "-tags gui") || !strings.Contains(err.Error(), "-headless") {
		t.Errorf("the error should name both alternatives: %v", err)
	}
}
