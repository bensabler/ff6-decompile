package route

import "testing"

// walkLeg and pulseLeg build the two leg shapes the SCN-0001 route uses.
func walkLeg(num int, dir Direction, axis Axis, cmp Compare, target uint8) Leg {
	return Leg{Num: num, Input: dir, Until: UntilPosition, Axis: axis,
		Compare: cmp, Target: target, Timeout: 900}
}

func pulseLeg(num int, until Until, dur int) Leg {
	return Leg{Num: num, Pulse: true, Until: until, Duration: dur, Timeout: 1800}
}

// drive steps the runner frame by frame over the supplied states and
// returns the outcome sequence and the leg numbers that were current.
func drive(rn *Runner, states []State) (outs []Outcome, legs []int) {
	for _, s := range states {
		leg, ok := rn.Current()
		if ok {
			legs = append(legs, leg.Num)
		} else {
			legs = append(legs, 0)
		}
		outs = append(outs, rn.Step(s))
	}
	return outs, legs
}

func TestLegSequencing(t *testing.T) {
	r := Route{
		Name: "two-walks", StartFrame: 10, NeutralFrames: 2,
		Legs: []Leg{
			walkLeg(1, Right, AxisX, AtLeast, 0x1B),
			walkLeg(2, Up, AxisY, AtMost, 0x27),
		},
	}
	rn := NewRunner(r)

	// Before StartFrame nothing runs.
	if got := rn.Step(State{Frame: 5, X: 0x1A, Y: 0x2A}); got != Continue {
		t.Fatalf("pre-start Step = %v, want Continue", got)
	}
	if _, ok := rn.Current(); ok {
		t.Fatal("Current() active before StartFrame")
	}

	// Arm at StartFrame, then leg 1 runs until X reaches its target.
	rn.Step(State{Frame: 10, X: 0x1A, Y: 0x2A})
	if leg, ok := rn.Current(); !ok || leg.Num != 1 {
		t.Fatalf("after arming Current() = %v/%v, want leg 1 active", leg.Num, ok)
	}
	if got := rn.Step(State{Frame: 11, X: 0x1A, Y: 0x2A}); got != Continue {
		t.Fatalf("leg 1 unmet = %v, want Continue", got)
	}
	if got := rn.Step(State{Frame: 12, X: 0x1B, Y: 0x2A}); got != Advanced {
		t.Fatalf("leg 1 met = %v, want Advanced", got)
	}
	if leg, _ := rn.Current(); leg.Num != 2 {
		t.Fatalf("current leg = %d, want 2", leg.Num)
	}

	// The neutral gap holds no input and does not evaluate leg 2.
	if got := rn.Input(State{Frame: 13, X: 0x1B, Y: 0x27}); got != None {
		t.Fatalf("input during neutral gap = %q, want None", got)
	}
	if got := rn.Step(State{Frame: 13, X: 0x1B, Y: 0x27}); got != Continue {
		t.Fatalf("step during neutral gap = %v, want Continue", got)
	}

	// Past the gap, leg 2's already-satisfied target finishes the route.
	if got := rn.Step(State{Frame: 14, X: 0x1B, Y: 0x27}); got != Done {
		t.Fatalf("final leg met = %v, want Done", got)
	}
	if !rn.Done() {
		t.Fatal("Done() = false after final leg")
	}
	if got := rn.Input(State{Frame: 15}); got != None {
		t.Fatalf("input after Done = %q, want None", got)
	}
}

func TestDirectionChanges(t *testing.T) {
	// The zigzag: up, sidestep left, up again. Each leg must hold only its
	// own direction, and only outside the neutral gap.
	r := Route{
		Name: "zigzag", StartFrame: 0, NeutralFrames: 2,
		Legs: []Leg{
			walkLeg(1, Up, AxisY, AtMost, 0x21),
			walkLeg(2, Left, AxisX, AtMost, 0x1D),
			walkLeg(3, Up, AxisY, AtMost, 0x1E),
		},
	}
	rn := NewRunner(r)
	rn.Step(State{Frame: 0, X: 0x1E, Y: 0x25})

	if got := rn.Input(State{Frame: 1, X: 0x1E, Y: 0x25}); got != Up {
		t.Fatalf("leg 1 input = %q, want up", got)
	}
	if got := rn.Step(State{Frame: 1, X: 0x1E, Y: 0x21}); got != Advanced {
		t.Fatalf("leg 1 = %v, want Advanced", got)
	}
	// Gap: no input even though leg 2 is current.
	if got := rn.Input(State{Frame: 2, X: 0x1E, Y: 0x21}); got != None {
		t.Fatalf("gap input = %q, want None", got)
	}
	if got := rn.Input(State{Frame: 3, X: 0x1E, Y: 0x21}); got != Left {
		t.Fatalf("leg 2 input = %q, want left", got)
	}
	if got := rn.Step(State{Frame: 3, X: 0x1D, Y: 0x21}); got != Advanced {
		t.Fatalf("leg 2 = %v, want Advanced", got)
	}
	if got := rn.Input(State{Frame: 6, X: 0x1D, Y: 0x21}); got != Up {
		t.Fatalf("leg 3 input = %q, want up", got)
	}
}

func TestCoordinateCompletionIsOvershootTolerant(t *testing.T) {
	tests := []struct {
		name    string
		leg     Leg
		x, y    uint8
		wantEnd bool
	}{
		{"right exact", walkLeg(1, Right, AxisX, AtLeast, 0x1B), 0x1B, 0, true},
		{"right overshoot", walkLeg(1, Right, AxisX, AtLeast, 0x1B), 0x1C, 0, true},
		{"right short", walkLeg(1, Right, AxisX, AtLeast, 0x1B), 0x1A, 0, false},
		{"up exact", walkLeg(1, Up, AxisY, AtMost, 0x16), 0, 0x16, true},
		{"up overshoot", walkLeg(1, Up, AxisY, AtMost, 0x16), 0, 0x15, true},
		{"up short", walkLeg(1, Up, AxisY, AtMost, 0x16), 0, 0x17, false},
		{"left overshoot", walkLeg(1, Left, AxisX, AtMost, 0x1D), 0x1C, 0, true},
		{"wrong axis ignored", walkLeg(1, Up, AxisY, AtMost, 0x16), 0x00, 0x20, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := Route{StartFrame: 0, Legs: []Leg{tt.leg}}
			rn := NewRunner(r)
			rn.Step(State{Frame: 0})
			got := rn.Step(State{Frame: 1, X: tt.x, Y: tt.y})
			if (got == Done) != tt.wantEnd {
				t.Fatalf("Step = %v, wantEnd = %v", got, tt.wantEnd)
			}
		})
	}
}

func TestTimeoutFailsAndNamesLeg(t *testing.T) {
	r := Route{
		Name: "stuck", StartFrame: 0, NeutralFrames: 0,
		Legs: []Leg{
			walkLeg(1, Right, AxisX, AtLeast, 0x1B),
			{Num: 2, Input: Up, Until: UntilPosition, Axis: AxisY,
				Compare: AtMost, Target: 0x10, Timeout: 5},
		},
	}
	rn := NewRunner(r)
	rn.Step(State{Frame: 0, X: 0x1A})
	if got := rn.Step(State{Frame: 1, X: 0x1B}); got != Advanced {
		t.Fatalf("leg 1 = %v, want Advanced", got)
	}
	// Leg 2 never reaches Y <= 0x10; it must fail rather than hang.
	var last Outcome
	for f := 2; f <= 12; f++ {
		last = rn.Step(State{Frame: f, X: 0x1B, Y: 0x20})
		if last == Failed {
			break
		}
	}
	if last != Failed {
		t.Fatalf("stuck leg outcome = %v, want Failed", last)
	}
	if got := rn.FailedLeg(); got != 2 {
		t.Fatalf("FailedLeg = %d, want 2 (the earliest divergent leg)", got)
	}
	if rn.Done() {
		t.Fatal("Done() = true after a failure")
	}
	// A failed route stays failed and stops driving input.
	if got := rn.Step(State{Frame: 13, X: 0x1B, Y: 0x10}); got != Continue {
		t.Fatalf("post-failure Step = %v, want Continue (inert)", got)
	}
	if got := rn.Input(State{Frame: 13}); got != None {
		t.Fatalf("post-failure input = %q, want None", got)
	}
}

func TestUnexpectedDivergenceDoesNotAdvance(t *testing.T) {
	// Walking right, the party instead drifts left (a wrong-turn or
	// scripted shove). The leg must neither advance nor silently accept it.
	r := Route{StartFrame: 0, Legs: []Leg{walkLeg(1, Right, AxisX, AtLeast, 0x1E)}}
	rn := NewRunner(r)
	rn.Step(State{Frame: 0, X: 0x1B})
	for f, x := range []uint8{0x1A, 0x19, 0x18} {
		if got := rn.Step(State{Frame: f + 1, X: x}); got != Continue {
			t.Fatalf("divergent frame %d (X=%02X) = %v, want Continue", f+1, x, got)
		}
	}
	if leg, ok := rn.Current(); !ok || leg.Num != 1 {
		t.Fatal("route left leg 1 despite moving away from the target")
	}
}

func TestBattleInterruptionAndResume(t *testing.T) {
	// A battle during a walk leg: the route must stop walking, must not
	// time the leg out for time spent fighting, and must resume afterwards.
	r := Route{
		Name: "interrupted", StartFrame: 0, NeutralFrames: 0,
		Legs: []Leg{
			{Num: 1, Input: Right, Until: UntilPosition, Axis: AxisX,
				Compare: AtLeast, Target: 0x1E, Timeout: 10},
		},
	}
	rn := NewRunner(r)
	rn.Step(State{Frame: 0, X: 0x1B})

	// 50 frames of battle — well past the 10-frame timeout.
	for f := 1; f <= 50; f++ {
		s := State{Frame: f, X: 0x1B, InBattle: true}
		if got := rn.Step(s); got != Continue {
			t.Fatalf("frame %d during battle = %v, want Continue", f, got)
		}
		if got := rn.Input(s); got != None {
			t.Fatalf("frame %d input during battle = %q, want None", f, got)
		}
	}
	if rn.FailedLeg() != 0 {
		t.Fatalf("leg timed out while a battle was active (FailedLeg = %d)", rn.FailedLeg())
	}
	// Control returns; the leg resumes and completes.
	if got := rn.Input(State{Frame: 51, X: 0x1B}); got != Right {
		t.Fatalf("input after battle = %q, want right", got)
	}
	if got := rn.Step(State{Frame: 51, X: 0x1E}); got != Done {
		t.Fatalf("leg after battle = %v, want Done", got)
	}
}

func TestWalkAndPulseLegHoldsDirectionUntilBattleEdge(t *testing.T) {
	// The trigger leg: keep walking into the guard while tapping A. The
	// direction must be held every frame outside battle, and the leg must
	// only end on the battle-start edge — not on any coordinate.
	r := Route{
		Name: "trigger", StartFrame: 0, NeutralFrames: 0,
		Legs: []Leg{
			{Num: 1, Input: Right, Pulse: true, Until: UntilBattleStart, Timeout: 500},
		},
	}
	rn := NewRunner(r)
	rn.Step(State{Frame: 0, X: 0x1E, Y: 0x27})

	for f := 1; f <= 20; f++ {
		s := State{Frame: f, X: 0x1E, Y: 0x27}
		if got := rn.Input(s); got != Right {
			t.Fatalf("frame %d input = %q, want right held", f, got)
		}
		if got := rn.Step(s); got != Continue {
			t.Fatalf("frame %d = %v, want Continue (no battle yet)", f, got)
		}
	}
	if got := rn.Step(State{Frame: 21, X: 0x1E, Y: 0x27, BattleStarted: true, InBattle: true}); got != Done {
		t.Fatalf("battle-start edge = %v, want Done", got)
	}
}

func TestBattleEdgeLegs(t *testing.T) {
	r := Route{
		Name: "battle-pair", StartFrame: 0, NeutralFrames: 0,
		Legs: []Leg{
			pulseLeg(1, UntilBattleStart, 0),
			pulseLeg(2, UntilBattleEnd, 0),
		},
	}
	rn := NewRunner(r)
	rn.Step(State{Frame: 0})

	if got := rn.Step(State{Frame: 1}); got != Continue {
		t.Fatalf("no battle yet = %v, want Continue", got)
	}
	if got := rn.Step(State{Frame: 2, BattleStarted: true, InBattle: true}); got != Advanced {
		t.Fatalf("battle-start edge = %v, want Advanced", got)
	}
	// Pulse legs keep evaluating during the battle (they are how the route
	// fights), unlike walk legs.
	if got := rn.Step(State{Frame: 3, InBattle: true}); got != Continue {
		t.Fatalf("mid-battle = %v, want Continue", got)
	}
	if got := rn.Step(State{Frame: 4, BattleEnded: true}); got != Done {
		t.Fatalf("battle-end edge = %v, want Done", got)
	}
}

func TestElapsedAndMapChangeLegs(t *testing.T) {
	r := Route{
		Name: "settle-then-transition", StartFrame: 0, NeutralFrames: 0,
		Legs: []Leg{
			pulseLeg(1, UntilElapsed, 5),
			{Num: 2, Input: Up, Until: UntilWatchedByte, WatchValue: 0x0D, Timeout: 100},
		},
	}
	rn := NewRunner(r)
	rn.Step(State{Frame: 0})

	if got := rn.Step(State{Frame: 4}); got != Continue {
		t.Fatalf("before elapsed = %v, want Continue", got)
	}
	if got := rn.Step(State{Frame: 5}); got != Advanced {
		t.Fatalf("at elapsed = %v, want Advanced", got)
	}
	if got := rn.Step(State{Frame: 6, WatchByte: 0x05}); got != Continue {
		t.Fatalf("exterior map value = %v, want Continue", got)
	}
	if got := rn.Step(State{Frame: 7, WatchByte: 0x0D}); got != Done {
		t.Fatalf("mines map value = %v, want Done", got)
	}
}

func TestValidate(t *testing.T) {
	good := walkLeg(1, Right, AxisX, AtLeast, 0x1B)
	tests := []struct {
		name    string
		route   Route
		wantErr bool
	}{
		{"valid walk", Route{Name: "r", Legs: []Leg{good}}, false},
		{"valid pulse", Route{Name: "r", Legs: []Leg{pulseLeg(1, UntilElapsed, 60)}}, false},
		{"empty", Route{Name: "r"}, true},
		{"misnumbered", Route{Name: "r", Legs: []Leg{walkLeg(2, Up, AxisY, AtMost, 1)}}, true},
		{"no timeout", Route{Name: "r", Legs: []Leg{
			{Num: 1, Input: Up, Until: UntilPosition}}}, true},
		{"position leg without direction", Route{Name: "r", Legs: []Leg{
			{Num: 1, Until: UntilPosition, Timeout: 10}}}, true},
		{"walk and pulse together is legal", Route{Name: "r", Legs: []Leg{
			{Num: 1, Input: Right, Pulse: true, Until: UntilBattleStart, Timeout: 10}}}, false},
		{"drives nothing", Route{Name: "r", Legs: []Leg{
			{Num: 1, Until: UntilBattleStart, Timeout: 10}}}, true},
		{"pulse elapsed without duration", Route{Name: "r", Legs: []Leg{
			{Num: 1, Pulse: true, Until: UntilElapsed, Timeout: 10}}}, true},
		{"elapsed leg that taps nothing", Route{Name: "r", Legs: []Leg{
			{Num: 1, Input: Up, Until: UntilElapsed, Duration: 5, Timeout: 10}}}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.route.Validate()
			if (err != nil) != tt.wantErr {
				t.Fatalf("Validate() error = %v, wantErr = %v", err, tt.wantErr)
			}
		})
	}
}

func TestMinesRouteEncodingIsValid(t *testing.T) {
	r := MinesRoute()
	if err := r.Validate(); err != nil {
		t.Fatalf("MinesRoute() invalid: %v", err)
	}
	if got, want := len(r.Legs), 17; got != want {
		t.Fatalf("MinesRoute() has %d legs, want %d", got, want)
	}
	// The final leg must detect the transition by position, not by the
	// map-context byte: EXP-0036 showed +$1EA5 reaches $0D while the party
	// is still on the exterior, which completes a map-value leg instantly
	// and skips the step into the shaft.
	transition := r.Legs[15]
	if transition.Until != UntilPosition || transition.Axis != AxisX ||
		transition.Compare != AtLeast || transition.Target != 0x26 {
		t.Fatalf("transition leg = %+v, want a position leg on X >= $26", transition)
	}
	last := r.Legs[len(r.Legs)-1]
	if last.Until != UntilPosition || last.Axis != AxisY ||
		last.Compare != AtMost || last.Target != 0x1C {
		t.Fatalf("final leg = %+v, want a position leg on Y <= $1C", last)
	}
	if last.ExpectX != 0x26 || last.ExpectY != 0x1C {
		t.Fatalf("final leg expects (%02X,%02X), want (26,1C)", last.ExpectX, last.ExpectY)
	}
}
