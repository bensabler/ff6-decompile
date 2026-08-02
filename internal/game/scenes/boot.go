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
	tables  *content.BattleTables
	tiles   *content.Tileset
	battle  int
	frame   uint64
	cursorX int
	cursorY int
	pressed int
}

// NewBoot returns the boot scene.
func NewBoot(font *content.Font) *Boot {
	return &Boot{font: font, cursorX: 8, cursorY: 152}
}

// WithFieldTiles lets X push the field tileset scene. Without it the demo
// still runs; a missing asset must not stop it launching.
func (b *Boot) WithFieldTiles(tiles *content.Tileset) *Boot {
	b.tiles = tiles
	return b
}

// WithBattle lets A push a battle scene for the given formation. Without it
// the boot scene stands alone, so a missing battle table is not fatal.
func (b *Boot) WithBattle(tables *content.BattleTables, formationID int) *Boot {
	b.tables, b.battle = tables, formationID
	return b
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
		if b.tables != nil {
			if sc, err := NewBattle(b.font, b.tables, b.battle); err == nil {
				ctx.Stack.Push(sc)
			}
		}
	}
	if ctx.Input.JustPressed(snespad.X) && b.tiles != nil {
		ctx.Stack.Push(NewFieldTiles(b.tiles, b.font))
	}
	// Start exits, which is how the headless host learns a run finished
	// rather than merely timing out.
	if ctx.Input.JustPressed(snespad.Start) {
		ctx.Stack.Pop()
	}
}

func (b *Boot) Draw(dst *framebuf.Indexed, _ *framebuf.Palette) {
	// Text takes a sub-palette *slot*; Rect takes a direct palette *index*.
	// Conflating the two is what made this scene's bright text render black
	// on black: `white = 3` reads as an index, but reached the blitter as a
	// palette base and put the font's ink on entries nothing had defined.
	const (
		primary   = content.SubPalettePrimary
		secondary = content.SubPaletteDim
		borderIdx = 2 // a direct index: gray, entry 2 of sub-palette 0
		cursorIdx = 3 // a direct index: white
	)
	dst.Fill(0)

	// A border, so the 256x224 extent and the host's letterboxing are both
	// visible at a glance.
	dst.Rect(0, 0, framebuf.Width, 1, borderIdx)
	dst.Rect(0, framebuf.Height-1, framebuf.Width, 1, borderIdx)
	dst.Rect(0, 0, 1, framebuf.Height, borderIdx)
	dst.Rect(framebuf.Width-1, 0, 1, framebuf.Height, borderIdx)

	y := 16
	line := func(s string, pal content.SubPalette) {
		b.font.DrawString(dst, 16, y, s, content.TextOptions{Palette: pal})
		y += 12
	}

	line("FF6 RECONSTRUCTION", primary)
	line("DEMO-0001A", secondary)
	y += 8
	line("FONT FROM EXTRACTED ROM DATA", secondary)
	line("GLYPH MAP EXP-0049", secondary)
	y += 8
	line("ABCDEFGHIJKLMNOPQRSTUVWXYZ", primary)
	line("abcdefghijklmnopqrstuvwxyz", primary)
	line("0123456789 - ?", primary)
	y += 8

	// Live values, so a still frame shows the loop is actually running.
	line(fmt.Sprintf("FRAME %d", b.frame), secondary)
	line(fmt.Sprintf("A PRESSED %d", b.pressed), secondary)
	y += 4
	if b.tables != nil {
		line("A BATTLE  X FIELD TILES  START EXIT", secondary)
	} else {
		line("DPAD MOVES  START EXITS", secondary)
	}

	// The cursor the d-pad drives. A direct index, like the border.
	dst.Rect(b.cursorX, b.cursorY, 8, 8, cursorIdx)
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
