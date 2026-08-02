//go:build gui

package ebitenhost

import (
	"testing"

	"github.com/bensabler/ff6-decompile/internal/graphics/framebuf"
	"github.com/bensabler/ff6-decompile/internal/platform/snespad"
)

func TestIntegerScale(t *testing.T) {
	tests := []struct {
		name string
		w, h int
		want int
	}{
		{"exact 1x", framebuf.Width, framebuf.Height, 1},
		{"exact 3x", framebuf.Width * 3, framebuf.Height * 3, 3},
		{"width-limited", framebuf.Width * 2, framebuf.Height * 5, 2},
		{"height-limited", framebuf.Width * 5, framebuf.Height * 2, 2},
		{"fractional rounds down", framebuf.Width*2 + 100, framebuf.Height*2 + 100, 2},
		// Never fractional and never zero: a fractional scale resamples the
		// framebuffer and destroys the pixel grid.
		{"smaller than native", 100, 100, 1},
		{"zero", 0, 0, 1},
		{"negative", -10, -10, 1},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := integerScale(tt.w, tt.h); got != tt.want {
				t.Errorf("integerScale(%d,%d) = %d, want %d", tt.w, tt.h, got, tt.want)
			}
		})
	}
}

// TestKeyBindingsCoverEveryButton guards against a control the operator
// cannot reach.
func TestKeyBindingsCoverEveryButton(t *testing.T) {
	covered := map[snespad.Button]bool{}
	for _, b := range keyBindings {
		covered[b] = true
	}
	for _, b := range snespad.All {
		if !covered[b] {
			t.Errorf("no key is bound to %v", b)
		}
	}
}

// TestControlsDocumentsEveryBoundKey keeps the printed help honest.
func TestControlsDocumentsEveryBoundKey(t *testing.T) {
	if len(Controls) == 0 {
		t.Fatal("Controls is empty")
	}
	for _, want := range []string{"D-pad", "Start", "Select", "confirm", "Esc"} {
		if !contains(Controls, want) {
			t.Errorf("Controls should mention %q:\n%s", want, Controls)
		}
	}
}

func contains(s, sub string) bool {
	for i := 0; i+len(sub) <= len(s); i++ {
		if s[i:i+len(sub)] == sub {
			return true
		}
	}
	return false
}
