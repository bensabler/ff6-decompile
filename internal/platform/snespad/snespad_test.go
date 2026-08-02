package snespad

import "testing"

// TestBitLayoutMatchesHardware pins the $4218/$4219 order. A future
// differential against Mesen reads that register directly, so a reordering
// here would silently break the comparison rather than fail a build.
func TestBitLayoutMatchesHardware(t *testing.T) {
	want := map[Button]uint16{
		R: 0x0010, L: 0x0020, X: 0x0040, A: 0x0080,
		Right: 0x0100, Left: 0x0200, Down: 0x0400, Up: 0x0800,
		Start: 0x1000, Select: 0x2000, Y: 0x4000, B: 0x8000,
	}
	for b, bits := range want {
		if uint16(b) != bits {
			t.Errorf("%v = $%04X, want $%04X", b, uint16(b), bits)
		}
	}
	// The low nibble is unused by the controller and must stay clear.
	var all Button
	for _, b := range All {
		all |= b
	}
	if all&0x000F != 0 {
		t.Errorf("buttons occupy bits 0-3, which the joypad register does not use: $%04X", all)
	}
	if len(All) != 12 {
		t.Errorf("All lists %d buttons, want 12", len(All))
	}
}

func TestStateHeldAndAny(t *testing.T) {
	s := State(0).With(A).With(Up)
	if !s.Held(A) || !s.Held(Up) {
		t.Error("With did not set the buttons")
	}
	if !s.Held(A | Up) {
		t.Error("Held with a multi-button mask should require all of them")
	}
	if s.Held(A | B) {
		t.Error("Held must not be satisfied by a partial match")
	}
	if !s.Any(A | B) {
		t.Error("Any should be satisfied by a partial match")
	}
	if s.Without(A).Held(A) {
		t.Error("Without did not clear the button")
	}
}

func TestStateString(t *testing.T) {
	if got := State(0).String(); got != "-" {
		t.Errorf("empty state = %q, want %q", got, "-")
	}
	// Order follows All, not bit order, so transcripts read consistently.
	if got := State(0).With(B).With(Up).String(); got != "Up+B" {
		t.Errorf("state string = %q, want %q", got, "Up+B")
	}
}

func TestEdges(t *testing.T) {
	e := NewEdges(0, State(0).With(A))
	if !e.JustPressed(A) || !e.IsHeld(A) || e.Released != 0 {
		t.Errorf("fresh press: %+v", e)
	}
	e = NewEdges(State(0).With(A), State(0).With(A))
	if e.JustPressed(A) || !e.IsHeld(A) {
		t.Errorf("held: %+v", e)
	}
	e = NewEdges(State(0).With(A), 0)
	if e.JustPressed(A) || e.IsHeld(A) || !e.Released.Held(A) {
		t.Errorf("release: %+v", e)
	}
}

func TestDirectionCancelsOpposites(t *testing.T) {
	tests := []struct {
		name   string
		held   State
		dx, dy int
	}{
		{"none", 0, 0, 0},
		{"left", State(0).With(Left), -1, 0},
		{"right", State(0).With(Right), 1, 0},
		{"up", State(0).With(Up), 0, -1},
		{"down", State(0).With(Down), 0, 1},
		{"diagonal", State(0).With(Up).With(Right), 1, -1},
		// A d-pad cannot physically do these; a keyboard can. Cancelling is
		// more predictable than last-wins.
		{"left+right cancel", State(0).With(Left).With(Right), 0, 0},
		{"up+down cancel", State(0).With(Up).With(Down), 0, 0},
		{"all four cancel", State(0).With(Up).With(Down).With(Left).With(Right), 0, 0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dx, dy := NewEdges(0, tt.held).Direction()
			if dx != tt.dx || dy != tt.dy {
				t.Errorf("Direction() = (%d,%d), want (%d,%d)", dx, dy, tt.dx, tt.dy)
			}
		})
	}
}
