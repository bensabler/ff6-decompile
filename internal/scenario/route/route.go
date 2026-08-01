// Package route models the scheduled input routes that drive the SCN-0001
// golden playthrough (docs/scenarios/SCN-0001-opening-to-whelk.md).
//
// A route is an ordered list of legs. Each leg holds one input for as long
// as its completion condition is unmet. Advancement is state-driven rather
// than duration-driven: walk legs finish when a player-position byte
// reaches its target, and the remaining legs finish on a battle edge, a
// map-value change, or an elapsed-frame settle. Timeouts exist only to
// detect divergence — a timed-out leg fails the route and names itself, it
// is never retried or corrected (EXP-0036).
//
// This package is the tracked model of the route the Mesen probe executes
// (mesen/probes/EXP-0036.lua); it holds the sequencing rules so they can be
// tested off-emulator. It is not an emulator and does not run input.
//
// Position bytes are the EXP-0035 candidates WRAM:+$00AF (X tile) and
// WRAM:+$00B0 (Y tile), both Strong hypothesis. The map-context byte
// WRAM:+$1EA5 is a Tentative candidate and is named accordingly.
package route

import "fmt"

// WRAM offsets the route observes. Recorded for traceability; this package
// does not model the SNES address space.
const (
	// AddrPlayerTileX is the player's X tile position (EXP-0035).
	AddrPlayerTileX = 0x00AF
	// AddrPlayerTileY is the player's Y tile position (EXP-0035).
	AddrPlayerTileY = 0x00B0
	// AddrMapContextCandidate is the candidate map-context byte. Its
	// semantics are unresolved: it separates opening, Narshe-exterior and
	// mines states, but a visually exterior scripted state also reads the
	// opening value (EXP-0035, EXP-0036).
	AddrMapContextCandidate = 0x1EA5
)

// Direction is a held d-pad input. The empty Direction means no input.
type Direction string

const (
	None  Direction = ""
	Up    Direction = "up"
	Down  Direction = "down"
	Left  Direction = "left"
	Right Direction = "right"
)

// Axis selects which position byte a walk leg compares.
type Axis int

const (
	AxisX Axis = iota
	AxisY
)

// Compare is the leg's completion comparison. Walk legs are
// overshoot-tolerant: a leg that overshoots into a wall still satisfies its
// intent, so targets are inequalities rather than equalities.
type Compare int

const (
	// AtLeast completes when the observed value is >= Target (moving
	// right or down).
	AtLeast Compare = iota
	// AtMost completes when the observed value is <= Target (moving left
	// or up).
	AtMost
)

// Until is the completion condition for legs that do not track a position.
type Until int

const (
	// UntilPosition is the default: the leg's Axis/Compare/Target apply.
	UntilPosition Until = iota
	// UntilBattleStart completes on a battle-entry edge.
	UntilBattleStart
	// UntilBattleEnd completes on a battle-end edge.
	UntilBattleEnd
	// UntilElapsed completes once Duration frames have passed. Used where
	// no observed state distinguishes "done" — dismissing the reward and
	// shaft dialogue windows.
	UntilElapsed
	// UntilMapChange completes when the map-context candidate reaches
	// MapValue.
	UntilMapChange
)

// Leg is one step of a route: hold Input until the condition is met.
type Leg struct {
	Num   int
	Input Direction
	// Pulse taps the A button for the leg's duration, in addition to any
	// held Input. A leg may hold a direction, pulse A, or both: EXP-0036
	// found that the guard trigger fires on walking *into* it, so the leg
	// that starts the fifth battle must keep holding its direction while
	// tapping A to dismiss the dialogue.
	Pulse bool

	Until    Until
	Axis     Axis
	Compare  Compare
	Target   uint8
	Duration int
	MapValue uint8

	// Timeout in frames. Exceeding it fails the route.
	Timeout int

	// ExpectX, ExpectY are the coordinates the leg is expected to end on,
	// recorded from the verified run. Zero means "not asserted".
	ExpectX, ExpectY uint8
	// Note is a project-authored description of what the leg does.
	Note string
}

// State is the observed machine state a route is evaluated against.
type State struct {
	Frame      int
	X, Y       uint8
	MapContext uint8
	InBattle   bool
	// BattleStarted and BattleEnded are edges: true only on the frame the
	// transition is observed.
	BattleStarted bool
	BattleEnded   bool
}

// Outcome reports what a step did.
type Outcome int

const (
	// Continue means the current leg is still running.
	Continue Outcome = iota
	// Advanced means the leg completed and the next leg is now current.
	Advanced
	// Done means the final leg completed.
	Done
	// Failed means the current leg exceeded its timeout.
	Failed
)

// Route is an ordered list of legs plus the frame at which leg 1 starts.
type Route struct {
	Name       string
	StartFrame int
	// NeutralFrames is held with no input between legs so the previous
	// direction is released before the next is applied.
	NeutralFrames int
	Legs          []Leg
}

// Runner tracks progress through a Route.
type Runner struct {
	route Route

	index        int // 0-based index of the current leg
	legStart     int // frame the current leg began
	neutralUntil int
	started      bool
	done         bool
	failedAt     int // leg number that timed out, 0 if none
}

// NewRunner returns a Runner positioned before the first leg.
func NewRunner(r Route) *Runner { return &Runner{route: r} }

// Current returns the leg in progress and whether one is active.
func (rn *Runner) Current() (Leg, bool) {
	if !rn.started || rn.done || rn.failedAt != 0 || rn.index >= len(rn.route.Legs) {
		return Leg{}, false
	}
	return rn.route.Legs[rn.index], true
}

// Done reports whether every leg completed.
func (rn *Runner) Done() bool { return rn.done }

// FailedLeg returns the leg number that timed out, or 0 if none has.
func (rn *Runner) FailedLeg() int { return rn.failedAt }

// Input returns the input to apply for the given state. It is None during
// the neutral gap between legs, while a battle owns the screen on a walk
// leg (the route fights instead of walking), and once the route has ended.
func (rn *Runner) Input(s State) Direction {
	leg, ok := rn.Current()
	if !ok || s.Frame < rn.neutralUntil {
		return None
	}
	if s.InBattle {
		return None
	}
	return leg.Input
}

// Step advances the runner against one observed state and reports what
// happened. Callers invoke it once per frame from StartFrame onward.
func (rn *Runner) Step(s State) Outcome {
	if rn.done || rn.failedAt != 0 {
		return Continue
	}
	if !rn.started {
		if s.Frame < rn.route.StartFrame {
			return Continue
		}
		rn.started = true
		rn.index = 0
		rn.legStart = s.Frame
		if len(rn.route.Legs) == 0 {
			rn.done = true
			return Done
		}
		return Continue
	}
	if s.Frame < rn.neutralUntil {
		return Continue
	}

	leg := rn.route.Legs[rn.index]

	// A position-tracking leg cannot make progress while a battle owns the
	// screen — the position bytes are not field-meaningful there (EXP-0036
	// observed them reading 00,00 at battle end) — and must not time out
	// for time spent fighting. Legs waiting on a battle edge still run.
	blocked := s.InBattle && leg.Until == UntilPosition
	if blocked {
		rn.legStart = s.Frame
		return Continue
	}

	if legSatisfied(leg, s, rn.legStart) {
		if rn.index == len(rn.route.Legs)-1 {
			rn.done = true
			return Done
		}
		rn.index++
		rn.neutralUntil = s.Frame + rn.route.NeutralFrames
		rn.legStart = rn.neutralUntil
		return Advanced
	}

	if leg.Timeout > 0 && s.Frame-rn.legStart > leg.Timeout {
		rn.failedAt = leg.Num
		return Failed
	}
	return Continue
}

func legSatisfied(leg Leg, s State, legStart int) bool {
	switch leg.Until {
	case UntilBattleStart:
		return s.BattleStarted
	case UntilBattleEnd:
		return s.BattleEnded
	case UntilElapsed:
		return s.Frame-legStart >= leg.Duration
	case UntilMapChange:
		return s.MapContext == leg.MapValue
	default:
		v := s.X
		if leg.Axis == AxisY {
			v = s.Y
		}
		if leg.Compare == AtLeast {
			return v >= leg.Target
		}
		return v <= leg.Target
	}
}

// Validate checks a route for encoding mistakes that would make a run
// uninterpretable: out-of-order or duplicated leg numbers, walk legs with
// no direction, position legs with no timeout, and pulse legs that expect a
// held direction.
func (r Route) Validate() error {
	if len(r.Legs) == 0 {
		return fmt.Errorf("route %q has no legs", r.Name)
	}
	for i, leg := range r.Legs {
		if leg.Num != i+1 {
			return fmt.Errorf("leg at index %d has Num %d, want %d", i, leg.Num, i+1)
		}
		if leg.Timeout <= 0 {
			return fmt.Errorf("leg %d has no timeout", leg.Num)
		}
		if leg.Input == None && !leg.Pulse {
			return fmt.Errorf("leg %d drives no input", leg.Num)
		}
		if leg.Until == UntilElapsed && leg.Duration <= 0 {
			return fmt.Errorf("leg %d waits on elapsed frames but has no duration", leg.Num)
		}
		if leg.Until == UntilElapsed && !leg.Pulse {
			return fmt.Errorf("leg %d waits on elapsed frames but taps nothing", leg.Num)
		}
		if leg.Until == UntilPosition && leg.Input == None {
			return fmt.Errorf("leg %d tracks position but holds no direction", leg.Num)
		}
	}
	return nil
}
