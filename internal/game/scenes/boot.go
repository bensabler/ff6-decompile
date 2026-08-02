// Package scenes holds the demo's scenes.
//
// A scene is one screen or mode. Scenes own game state and drawing; they know
// nothing about windows, files, or the platform.
package scenes

import (
	"fmt"

	"github.com/bensabler/ff6-decompile/internal/content"
	"github.com/bensabler/ff6-decompile/internal/engine"
	"github.com/bensabler/ff6-decompile/internal/graphics/framebuf"
	"github.com/bensabler/ff6-decompile/internal/platform/snespad"
)

// Boot is DEMO-0001A's proof scene: it renders real FF6 font data from the
// generated archive and responds to input.
//
// It is explicitly **not** an FF6 screen. FF6's title and opening are
// unlocated (readiness F14, B02), and drawing an invented approximation of
// them would be exactly the kind of scaffolding that later reads as parity.
// This screen reports what the demo can currently do, in the game's own
// glyphs, and says plainly what is not implemented.
type Boot struct {
	font    *content.Font
	frame   uint64
	cursorX int
	cursorY int
	pressed int
}

// NewBoot returns the boot scene.
func NewBoot(font *content.Font) *Boot {
	return &Boot{font: font, cursorX: 8, cursorY: 152}
}

func (b *Boot) Update(ctx *engine.Context) {
	b.frame = ctx.Frame

	// Movement proves the input path end to end: host -> Source -> Edges.
	dx, dy := ctx.Input.Direction()
	b.cursorX += dx
	b.cursorY += dy
	b.cursorX = clamp(b.cursorX, 0, framebuf.Width-8)
	b.cursorY = clamp(b.cursorY, 0, framebuf.Height-8)

	if ctx.Input.JustPressed(snespad.A) {
		b.pressed++
	}
	// Start exits, which is how the headless host learns a run finished
	// rather than merely timing out.
	if ctx.Input.JustPressed(snespad.Start) {
		ctx.Stack.Pop()
	}
}

func (b *Boot) Draw(dst *framebuf.Indexed, _ *framebuf.Palette) {
	const (
		white = 3
		gray  = 2
	)
	dst.Fill(0)

	// A border, so the 256x224 extent and the host's letterboxing are both
	// visible at a glance.
	dst.Rect(0, 0, framebuf.Width, 1, gray)
	dst.Rect(0, framebuf.Height-1, framebuf.Width, 1, gray)
	dst.Rect(0, 0, 1, framebuf.Height, gray)
	dst.Rect(framebuf.Width-1, 0, 1, framebuf.Height, gray)

	y := 16
	line := func(s string, pal uint8) {
		b.font.DrawString(dst, 16, y, s, content.TextOptions{PaletteBase: pal})
		y += 12
	}

	line("FF6 RECONSTRUCTION", white)
	line("DEMO-0001A", gray)
	y += 8
	line("FONT FROM EXTRACTED ROM DATA", gray)
	line("GLYPH MAP EXP-0049", gray)
	y += 8
	line("ABCDEFGHIJKLMNOPQRSTUVWXYZ", white)
	line("abcdefghijklmnopqrstuvwxyz", white)
	line("0123456789 - ?", white)
	y += 8

	// Live values, so a still frame shows the loop is actually running.
	line(fmt.Sprintf("FRAME %d", b.frame), gray)
	line(fmt.Sprintf("A PRESSED %d", b.pressed), gray)
	y += 4
	line("DPAD MOVES  START EXITS", gray)

	// The cursor the d-pad drives.
	dst.Rect(b.cursorX, b.cursorY, 8, 8, white)
}

func clamp(v, lo, hi int) int {
	if v < lo {
		return lo
	}
	if v > hi {
		return hi
	}
	return v
}
