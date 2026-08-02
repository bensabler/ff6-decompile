// Package snespad models the SNES controller as the hardware presents it.
//
// The bit order is the one the auto-joypad read latches into $4218/$4219, so
// a future differential against Mesen can compare the demo's pad state to the
// emulator's register without a translation layer in between.
package snespad

import "strings"

// Button is one controller bit. Values match the $4218/$4219 layout: the low
// byte holds A/X/L/R in bits 4-7, the high byte B/Y/Select/Start and the
// d-pad.
type Button uint16

const (
	R      Button = 1 << 4
	L      Button = 1 << 5
	X      Button = 1 << 6
	A      Button = 1 << 7
	Right  Button = 1 << 8
	Left   Button = 1 << 9
	Down   Button = 1 << 10
	Up     Button = 1 << 11
	Start  Button = 1 << 12
	Select Button = 1 << 13
	Y      Button = 1 << 14
	B      Button = 1 << 15
)

// All lists every button in a stable order, for iteration and formatting.
var All = []Button{Up, Down, Left, Right, A, B, X, Y, L, R, Start, Select}

func (b Button) String() string {
	switch b {
	case Up:
		return "Up"
	case Down:
		return "Down"
	case Left:
		return "Left"
	case Right:
		return "Right"
	case A:
		return "A"
	case B:
		return "B"
	case X:
		return "X"
	case Y:
		return "Y"
	case L:
		return "L"
	case R:
		return "R"
	case Start:
		return "Start"
	case Select:
		return "Select"
	default:
		return "?"
	}
}

// State is the set of buttons held on one frame.
type State uint16

// Held reports whether every button in mask is down.
func (s State) Held(mask Button) bool { return Button(s)&mask == mask }

// Any reports whether any button in mask is down.
func (s State) Any(mask Button) bool { return Button(s)&mask != 0 }

// With returns s with mask pressed.
func (s State) With(mask Button) State { return s | State(mask) }

// Without returns s with mask released.
func (s State) Without(mask Button) State { return s &^ State(mask) }

func (s State) String() string {
	if s == 0 {
		return "-"
	}
	var b strings.Builder
	for _, btn := range All {
		if s.Held(btn) {
			if b.Len() > 0 {
				b.WriteByte('+')
			}
			b.WriteString(btn.String())
		}
	}
	return b.String()
}

// Edges holds one frame of input together with the transitions since the
// previous frame. Menus need presses, movement needs holds, and computing
// both in one place keeps scenes from having to remember the last frame.
type Edges struct {
	Held     State
	Pressed  State // down this frame, up last frame
	Released State // up this frame, down last frame
}

// NewEdges derives the transitions between two frames.
func NewEdges(prev, cur State) Edges {
	return Edges{
		Held:     cur,
		Pressed:  cur &^ prev,
		Released: prev &^ cur,
	}
}

// JustPressed reports a fresh press of every button in mask.
func (e Edges) JustPressed(mask Button) bool { return Button(e.Pressed)&mask == mask }

// IsHeld reports whether every button in mask is currently down.
func (e Edges) IsHeld(mask Button) bool { return e.Held.Held(mask) }

// Direction returns the d-pad as a unit step, with opposing presses
// cancelling — the SNES d-pad is mechanically incapable of opposing presses,
// but a keyboard is not, and cancelling is more predictable than
// last-wins.
func (e Edges) Direction() (dx, dy int) {
	if e.Held.Any(Left) {
		dx--
	}
	if e.Held.Any(Right) {
		dx++
	}
	if e.Held.Any(Up) {
		dy--
	}
	if e.Held.Any(Down) {
		dy++
	}
	return dx, dy
}
