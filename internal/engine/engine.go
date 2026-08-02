// Package engine is the demo's scene machine and fixed-step loop.
//
// It has no third-party dependencies and does no I/O. A *host* (windowed or
// headless) owns the platform, and drives the machine through Tick and
// Render. That split is what lets the whole demo be tested without a display.
//
// # Determinism
//
// The machine is deterministic by construction, not by discipline:
//
//   - Tick advances exactly one logical frame. It reads no wall clock, keeps
//     no time accumulator, and starts no goroutines.
//   - Input arrives from a Source, so a test can supply a scripted one.
//   - Any randomness a scene needs must come from state the scene owns, so it
//     is captured by the frame count and the scene's own fields.
//
// The consequence a test can assert: the same machine plus the same input
// script produces the same per-frame framebuffer hash, however fast or slow
// the host calls Tick.
package engine

import (
	"github.com/bensabler/ff6-decompile/internal/graphics/framebuf"
	"github.com/bensabler/ff6-decompile/internal/platform/snespad"
)

// TicksPerSecond is the demo's logical update rate. The SNES NTSC frame rate
// is ~60.098 Hz; 60 is the integer rate hosts run at, and every timing claim
// the demo makes is in ticks, not seconds, so the difference never enters the
// simulation.
const TicksPerSecond = 60

// Source supplies one frame of controller state. Implementations must be
// deterministic given the same call sequence.
type Source interface {
	// Poll returns the buttons held for the frame about to be simulated.
	Poll(frame uint64) snespad.State
}

// SourceFunc adapts a function to Source.
type SourceFunc func(frame uint64) snespad.State

func (f SourceFunc) Poll(frame uint64) snespad.State { return f(frame) }

// NoInput is a Source that never presses anything.
var NoInput Source = SourceFunc(func(uint64) snespad.State { return 0 })

// Context is what a scene gets each frame.
type Context struct {
	// Frame counts ticks since the machine started, and is the demo's only
	// clock.
	Frame uint64
	// Input carries this frame's held buttons and transitions.
	Input snespad.Edges
	// Stack is the scene stack, for pushing, popping and replacing.
	Stack *Stack
}

// Scene is one screen or mode: a title, a field map, a battle, a dialogue
// window.
type Scene interface {
	// Update advances the scene one logical frame.
	Update(ctx *Context)
	// Draw composes the scene into the framebuffer. Draw must not mutate
	// scene state — a host may skip or repeat it.
	Draw(dst *framebuf.Indexed, pal *framebuf.Palette)
}

// TransparentScene is implemented by scenes that do not fill the screen, so
// the scene beneath keeps drawing. A dialogue window over a field map is the
// motivating case.
type TransparentScene interface {
	Scene
	// DrawBelow reports whether the scene under this one should draw first.
	DrawBelow() bool
}

// LiveScene is implemented by scenes that let the scene beneath keep
// updating. A battle's ATB continues while a submenu is open, which is
// exactly the ACTIVE/WAIT behaviour EXP-0044 established.
type LiveScene interface {
	Scene
	// UpdateBelow reports whether the scene under this one should update.
	UpdateBelow() bool
}

// Stack is a scene stack. The top scene always updates and draws; scenes
// below participate only if the scene above opts in.
type Stack struct {
	scenes  []Scene
	pending []func()
}

// NewStack returns a stack holding the given scene, bottom-first.
func NewStack(scenes ...Scene) *Stack {
	return &Stack{scenes: append([]Scene(nil), scenes...)}
}

// Len returns the number of scenes on the stack.
func (s *Stack) Len() int { return len(s.scenes) }

// Top returns the active scene, or nil if the stack is empty.
func (s *Stack) Top() Scene {
	if len(s.scenes) == 0 {
		return nil
	}
	return s.scenes[len(s.scenes)-1]
}

// Push, Pop and Replace are deferred to the end of the frame, so a scene can
// call them from inside its own Update without the stack changing underneath
// the in-progress traversal.
func (s *Stack) Push(sc Scene) {
	s.pending = append(s.pending, func() { s.scenes = append(s.scenes, sc) })
}

// Pop removes the top scene.
func (s *Stack) Pop() {
	s.pending = append(s.pending, func() {
		if len(s.scenes) > 0 {
			s.scenes = s.scenes[:len(s.scenes)-1]
		}
	})
}

// Replace swaps the top scene for sc.
func (s *Stack) Replace(sc Scene) {
	s.pending = append(s.pending, func() {
		if len(s.scenes) == 0 {
			s.scenes = append(s.scenes, sc)
			return
		}
		s.scenes[len(s.scenes)-1] = sc
	})
}

// applyPending runs deferred mutations in the order they were requested.
func (s *Stack) applyPending() {
	for _, f := range s.pending {
		f()
	}
	s.pending = s.pending[:0]
}

// updateDepth returns how many scenes from the top should update this frame.
func (s *Stack) updateDepth() int {
	n := len(s.scenes)
	if n == 0 {
		return 0
	}
	depth := 1
	for i := n - 1; i > 0; i-- {
		live, ok := s.scenes[i].(LiveScene)
		if !ok || !live.UpdateBelow() {
			break
		}
		depth++
	}
	return depth
}

// drawDepth returns how many scenes from the top should draw this frame.
func (s *Stack) drawDepth() int {
	n := len(s.scenes)
	if n == 0 {
		return 0
	}
	depth := 1
	for i := n - 1; i > 0; i-- {
		tr, ok := s.scenes[i].(TransparentScene)
		if !ok || !tr.DrawBelow() {
			break
		}
		depth++
	}
	return depth
}

// Machine drives a scene stack. The zero value is not usable; use New.
type Machine struct {
	stack  *Stack
	input  Source
	fb     *framebuf.Indexed
	pal    *framebuf.Palette
	frame  uint64
	prevIn snespad.State
}

// New returns a machine running the given scenes, bottom-first.
func New(input Source, pal *framebuf.Palette, scenes ...Scene) *Machine {
	if input == nil {
		input = NoInput
	}
	if pal == nil {
		pal = framebuf.GrayPalette()
	}
	return &Machine{
		stack: NewStack(scenes...),
		input: input,
		fb:    framebuf.New(),
		pal:   pal,
	}
}

// Frame returns the number of ticks simulated so far.
func (m *Machine) Frame() uint64 { return m.frame }

// Stack exposes the scene stack, for hosts that need to inspect it.
func (m *Machine) Stack() *Stack { return m.stack }

// Palette returns the machine's palette, which scenes may modify in place.
func (m *Machine) Palette() *framebuf.Palette { return m.pal }

// Done reports whether every scene has been popped, which is how the demo
// signals a clean exit.
func (m *Machine) Done() bool { return m.stack.Len() == 0 }

// Tick advances the simulation exactly one frame.
func (m *Machine) Tick() {
	cur := m.input.Poll(m.frame)
	ctx := &Context{
		Frame: m.frame,
		Input: snespad.NewEdges(m.prevIn, cur),
		Stack: m.stack,
	}

	// Snapshot the scenes to update before any of them mutate the stack, so
	// a push during Update takes effect next frame rather than mid-traversal.
	depth := m.stack.updateDepth()
	n := len(m.stack.scenes)
	active := append([]Scene(nil), m.stack.scenes[n-depth:]...)
	// Top-down: the active scene decides first whether those below run.
	for i := len(active) - 1; i >= 0; i-- {
		active[i].Update(ctx)
	}

	m.stack.applyPending()
	m.prevIn = cur
	m.frame++
}

// Render composes the current frame and returns the framebuffer and palette.
// It does not advance the simulation and may be called any number of times
// per Tick, or not at all.
func (m *Machine) Render() (*framebuf.Indexed, *framebuf.Palette) {
	m.fb.Fill(0)
	depth := m.stack.drawDepth()
	n := len(m.stack.scenes)
	// Bottom-up, so the topmost scene composes last.
	for i := n - depth; i < n; i++ {
		m.stack.scenes[i].Draw(m.fb, m.pal)
	}
	return m.fb, m.pal
}
