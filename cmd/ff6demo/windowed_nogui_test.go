//go:build !gui

package main

import (
	"strings"
	"testing"
)

// TestRunWindowedIsUnavailableWithoutTheGuiTag pins the message an operator
// sees when they run the default build and expect a window.
//
// It is tagged !gui deliberately. Under the gui tag runWindowed really does
// open a window, which cannot be done from a test binary's goroutine — the
// first attempt at this test crashed the gui suite with "NSWindow should only
// be instantiated on the main thread". The windowed path's testable parts
// live in internal/engine/ebitenhost instead.
func TestRunWindowedIsUnavailableWithoutTheGuiTag(t *testing.T) {
	err := runWindowed(nil)
	if err == nil {
		t.Fatal("the default build must not claim to open a window")
	}
	for _, want := range []string{"-tags gui", "-headless"} {
		if !strings.Contains(err.Error(), want) {
			t.Errorf("the error should name %q so the operator knows the way forward:\n%v", want, err)
		}
	}
}
