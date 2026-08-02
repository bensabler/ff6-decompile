//go:build gui

// Package ebitenhost is the windowed host, and the only package in this
// repository permitted to import Ebitengine.
//
// The confinement is the point. `audit.CheckImportBoundaries` rejects an
// ebiten import anywhere else, so the dependency can be replaced by editing
// one package. Everything below — engine, content, game — stays
// dependency-free and testable without a display.
//
// # Why a build tag
//
// Measured against ebiten v2.9.9 on 2026-08-02: CGO_ENABLED=0 fails to build
// on both linux/amd64 and darwin/amd64 (only windows and js/wasm succeed
// without cgo). Without the tag, `go build ./...`, `go vet ./...` and
// `go test ./...` would all require system OpenGL and X11 headers, and CI
// could no longer run the demo's test surface on a bare container.
//
// # Not authoritative for tests
//
// Ebitengine may call Update more or fewer times than Draw, and pauses both
// when the window loses focus unless told otherwise. The simulation stays
// exact because Tick is 1:1 with Update, but frame-for-frame comparison
// belongs to internal/engine/headless.
package ebitenhost

import (
	"image"

	"github.com/hajimehoshi/ebiten/v2"

	"github.com/bensabler/ff6-decompile/internal/engine"
	"github.com/bensabler/ff6-decompile/internal/graphics/framebuf"
	"github.com/bensabler/ff6-decompile/internal/platform/snespad"
)

// keyBindings maps the keyboard to the SNES pad. The face-button layout
// follows the physical controller: B is the lower button (confirm in FF6), A
// the right one.
var keyBindings = map[ebiten.Key]snespad.Button{
	ebiten.KeyArrowUp:    snespad.Up,
	ebiten.KeyArrowDown:  snespad.Down,
	ebiten.KeyArrowLeft:  snespad.Left,
	ebiten.KeyArrowRight: snespad.Right,
	ebiten.KeyW:          snespad.Up,
	ebiten.KeyS:          snespad.Down,
	ebiten.KeyA:          snespad.Left,
	ebiten.KeyD:          snespad.Right,

	ebiten.KeyZ:          snespad.B,
	ebiten.KeyX:          snespad.A,
	ebiten.KeyC:          snespad.Y,
	ebiten.KeyV:          snespad.X,
	ebiten.KeyQ:          snespad.L,
	ebiten.KeyE:          snespad.R,
	ebiten.KeyEnter:      snespad.Start,
	ebiten.KeyBackspace:  snespad.Select,
	ebiten.KeyShiftRight: snespad.Select,
}

// Controls is the operator-facing key map, so the binding is documented in
// exactly one place.
const Controls = `  Arrows / WASD   D-pad
  Z               B (confirm)
  X               A
  C / V           Y / X
  Q / E           L / R
  Enter           Start
  Backspace       Select
  Esc             quit`

// keyboardSource reads the keyboard each frame.
//
// It ignores the frame argument deliberately: a live keyboard is not
// reproducible, which is why scenario tests use a scripted Source and the
// headless host instead.
type keyboardSource struct{}

func (keyboardSource) Poll(uint64) snespad.State {
	var s snespad.State
	for key, btn := range keyBindings {
		if ebiten.IsKeyPressed(key) {
			s = s.With(btn)
		}
	}
	return s
}

// NewKeyboardSource returns an engine.Source driven by the keyboard.
func NewKeyboardSource() engine.Source { return keyboardSource{} }

// game adapts engine.Machine to ebiten.Game. It holds no game state — the
// moment it needs any, the abstraction has failed.
type game struct {
	m      *engine.Machine
	rgba   *image.RGBA
	screen *ebiten.Image
}

func (g *game) Update() error {
	if ebiten.IsKeyPressed(ebiten.KeyEscape) || g.m.Done() {
		return ebiten.Termination
	}
	g.m.Tick()
	return nil
}

func (g *game) Draw(dst *ebiten.Image) {
	fb, pal := g.m.Render()
	g.rgba = fb.Resolve(g.rgba, pal)
	if g.screen == nil {
		g.screen = ebiten.NewImage(framebuf.Width, framebuf.Height)
	}
	g.screen.WritePixels(g.rgba.Pix)

	w, h := dst.Bounds().Dx(), dst.Bounds().Dy()
	scale := integerScale(w, h)

	// Floor the offsets so the image lands on whole device pixels; a half
	// pixel offset reintroduces the blur that FilterNearest exists to avoid.
	ox := (w - framebuf.Width*scale) / 2
	oy := (h - framebuf.Height*scale) / 2

	op := &ebiten.DrawImageOptions{Filter: ebiten.FilterNearest}
	op.GeoM.Scale(float64(scale), float64(scale))
	op.GeoM.Translate(float64(ox), float64(oy))
	dst.DrawImage(g.screen, op)
}

// LayoutF returns the size in device pixels, so the framebuffer is not
// resampled through a logical-pixel layer on a HiDPI display.
func (g *game) LayoutF(logicalW, logicalH float64) (float64, float64) {
	s := ebiten.Monitor().DeviceScaleFactor()
	return logicalW * s, logicalH * s
}

// Layout satisfies the interface; LayoutF supersedes it.
func (g *game) Layout(w, h int) (int, int) { return w, h }

// integerScale is the largest whole-number scale that fits, never below 1 —
// a fractional scale would resample the framebuffer and destroy the pixel
// grid the SNES output depends on.
func integerScale(w, h int) int {
	s := w / framebuf.Width
	if v := h / framebuf.Height; v < s {
		s = v
	}
	if s < 1 {
		s = 1
	}
	return s
}

// Run opens a window and drives the machine until the scene stack empties or
// the operator quits.
func Run(m *engine.Machine, title string) error {
	const initialScale = 3
	ebiten.SetWindowSize(framebuf.Width*initialScale, framebuf.Height*initialScale)
	ebiten.SetWindowTitle(title)
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)
	ebiten.SetTPS(engine.TicksPerSecond)
	// Keep simulating when the window loses focus: a demo that silently
	// freezes on alt-tab looks like a hang.
	ebiten.SetRunnableOnUnfocused(true)

	if err := ebiten.RunGame(&game{m: m}); err != nil && err != ebiten.Termination {
		return err
	}
	return nil
}
