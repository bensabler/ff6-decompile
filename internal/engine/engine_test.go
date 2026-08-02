package engine

import (
	"testing"
	"time"

	"github.com/bensabler/ff6-decompile/internal/graphics/framebuf"
	"github.com/bensabler/ff6-decompile/internal/platform/snespad"
)

// recorder is a scene that logs what happened to it, for ordering tests.
type recorder struct {
	name        string
	log         *[]string
	updateBelow bool
	drawBelow   bool
	live, trans bool
	onUpdate    func(ctx *Context)
	fillIdx     uint8
}

func (r *recorder) Update(ctx *Context) {
	*r.log = append(*r.log, "u:"+r.name)
	if r.onUpdate != nil {
		r.onUpdate(ctx)
	}
}

func (r *recorder) Draw(dst *framebuf.Indexed, _ *framebuf.Palette) {
	*r.log = append(*r.log, "d:"+r.name)
	if r.fillIdx != 0 {
		dst.Rect(0, 0, 4, 4, r.fillIdx)
	}
}

type liveRecorder struct{ *recorder }

func (l liveRecorder) UpdateBelow() bool { return l.updateBelow }

type transRecorder struct{ *recorder }

func (t transRecorder) DrawBelow() bool { return t.drawBelow }

type liveTransRecorder struct{ *recorder }

func (l liveTransRecorder) UpdateBelow() bool { return l.updateBelow }
func (l liveTransRecorder) DrawBelow() bool   { return l.drawBelow }

func TestStackUpdatesTopDownAndDrawsBottomUp(t *testing.T) {
	var log []string
	bottom := &recorder{name: "bottom", log: &log}
	top := liveTransRecorder{&recorder{name: "top", log: &log, updateBelow: true, drawBelow: true}}

	m := New(NoInput, nil, bottom, top)
	m.Tick()
	m.Render()

	want := []string{"u:top", "u:bottom", "d:bottom", "d:top"}
	if len(log) != len(want) {
		t.Fatalf("log = %v, want %v", log, want)
	}
	for i := range want {
		if log[i] != want[i] {
			t.Fatalf("log = %v, want %v", log, want)
		}
	}
}

func TestOpaqueSceneHidesThoseBelow(t *testing.T) {
	var log []string
	bottom := &recorder{name: "bottom", log: &log}
	top := &recorder{name: "top", log: &log} // neither Live nor Transparent

	m := New(NoInput, nil, bottom, top)
	m.Tick()
	m.Render()

	for _, e := range log {
		if e == "u:bottom" || e == "d:bottom" {
			t.Fatalf("an opaque top scene should suppress the one below; log = %v", log)
		}
	}
}

// TestLiveSceneMirrorsTheATBGate documents why LiveScene exists: EXP-0044
// established that a battle keeps advancing while the command window is open,
// but pauses under the ability list. Those are a live and an opaque scene.
func TestLiveSceneMirrorsTheATBGate(t *testing.T) {
	var log []string
	battle := &recorder{name: "battle", log: &log}
	commandWindow := liveRecorder{&recorder{name: "cmd", log: &log, updateBelow: true}}

	m := New(NoInput, nil, battle, commandWindow)
	m.Tick()
	found := false
	for _, e := range log {
		if e == "u:battle" {
			found = true
		}
	}
	if !found {
		t.Error("a live scene must let the battle below keep updating")
	}

	log = log[:0]
	abilityList := &recorder{name: "ability", log: &log}
	m2 := New(NoInput, nil, battle, abilityList)
	m2.Tick()
	for _, e := range log {
		if e == "u:battle" {
			t.Error("an opaque submenu must pause the battle below")
		}
	}
}

func TestStackMutationsAreDeferredToFrameEnd(t *testing.T) {
	var log []string
	pushed := &recorder{name: "pushed", log: &log}
	bottom := &recorder{name: "bottom", log: &log}
	bottom.onUpdate = func(ctx *Context) {
		if ctx.Frame == 0 {
			ctx.Stack.Push(pushed)
		}
	}

	m := New(NoInput, nil, bottom)
	m.Tick()
	// The push must not have taken effect mid-frame.
	for _, e := range log {
		if e == "u:pushed" {
			t.Fatal("a scene pushed during Update ran in the same frame")
		}
	}
	if m.Stack().Len() != 2 {
		t.Fatalf("stack length %d after the frame, want 2", m.Stack().Len())
	}
	m.Tick()
	found := false
	for _, e := range log {
		if e == "u:pushed" {
			found = true
		}
	}
	if !found {
		t.Error("the pushed scene should update on the next frame")
	}
}

func TestPopAndReplaceAndDone(t *testing.T) {
	var log []string
	a := &recorder{name: "a", log: &log}
	b := &recorder{name: "b", log: &log}

	m := New(NoInput, nil, a)
	if m.Done() {
		t.Error("a machine with a scene is not done")
	}
	m.Stack().Replace(b)
	m.Tick()
	if m.Stack().Top() != Scene(b) {
		t.Error("Replace did not swap the top scene")
	}
	m.Stack().Pop()
	m.Tick()
	if !m.Done() {
		t.Error("popping the last scene should leave the machine done")
	}
	// An empty machine must stay safe.
	m.Tick()
	m.Render()
}

func TestReplaceOnEmptyStackPushes(t *testing.T) {
	var log []string
	m := New(NoInput, nil)
	if !m.Done() {
		t.Fatal("a machine with no scenes starts done")
	}
	m.Stack().Replace(&recorder{name: "a", log: &log})
	m.Tick()
	if m.Stack().Len() != 1 {
		t.Errorf("Replace on an empty stack should push; length %d", m.Stack().Len())
	}
}

func TestInputEdges(t *testing.T) {
	script := map[uint64]snespad.State{
		0: 0,
		1: snespad.State(snespad.A),
		2: snespad.State(snespad.A),
		3: 0,
	}
	var seen []snespad.Edges
	sc := &recorder{name: "s", log: new([]string)}
	sc.onUpdate = func(ctx *Context) { seen = append(seen, ctx.Input) }

	m := New(SourceFunc(func(f uint64) snespad.State { return script[f] }), nil, sc)
	for i := 0; i < 4; i++ {
		m.Tick()
	}

	if seen[0].JustPressed(snespad.A) {
		t.Error("frame 0: A was not pressed")
	}
	if !seen[1].JustPressed(snespad.A) {
		t.Error("frame 1: A should register a fresh press")
	}
	if seen[2].JustPressed(snespad.A) {
		t.Error("frame 2: a held button is not a fresh press")
	}
	if !seen[2].IsHeld(snespad.A) {
		t.Error("frame 2: A should still be held")
	}
	if !seen[3].Released.Held(snespad.A) {
		t.Error("frame 3: A should register a release")
	}
}

// TestDeterminismUnderVaryingHostPace is the property the whole test strategy
// rests on: the simulation must depend on the tick sequence and nothing else.
// The second run inserts real delays between ticks and calls Render a varying
// number of times, both of which a windowed host does.
func TestDeterminismUnderVaryingHostPace(t *testing.T) {
	build := func() *Machine {
		sc := &counterScene{}
		return New(SourceFunc(func(f uint64) snespad.State {
			if f%7 == 0 {
				return snespad.State(snespad.A)
			}
			return snespad.State(snespad.Right)
		}), nil, sc)
	}

	const frames = 120
	first := make([][32]byte, frames)
	m := build()
	for i := 0; i < frames; i++ {
		m.Tick()
		fb, _ := m.Render()
		first[i] = fb.Sum256()
	}

	second := make([][32]byte, frames)
	m2 := build()
	for i := 0; i < frames; i++ {
		m2.Tick()
		time.Sleep(time.Microsecond) // a host's pace must not matter
		fb, _ := m2.Render()
		if i%3 == 0 {
			m2.Render() // an extra draw must not advance anything
		}
		second[i] = fb.Sum256()
	}

	for i := range first {
		if first[i] != second[i] {
			t.Fatalf("frame %d differs between runs: the simulation is not deterministic", i)
		}
	}
	if m.Frame() != m2.Frame() {
		t.Fatalf("frame counters diverged: %d vs %d", m.Frame(), m2.Frame())
	}
}

// counterScene draws something that changes every frame and responds to
// input, so the determinism test has real state to compare.
type counterScene struct {
	x, y  int
	frame uint64
}

func (c *counterScene) Update(ctx *Context) {
	c.frame = ctx.Frame
	dx, dy := ctx.Input.Direction()
	c.x += dx
	c.y += dy
	if ctx.Input.JustPressed(snespad.A) {
		c.x = 0
	}
}

func (c *counterScene) Draw(dst *framebuf.Indexed, _ *framebuf.Palette) {
	dst.Rect(c.x%framebuf.Width, c.y%framebuf.Height, 8, 8, uint8(c.frame%255)+1)
}

func TestRenderDoesNotAdvanceTheSimulation(t *testing.T) {
	sc := &counterScene{}
	m := New(NoInput, nil, sc)
	m.Render()
	m.Render()
	if m.Frame() != 0 {
		t.Errorf("Render advanced the frame counter to %d", m.Frame())
	}
}
