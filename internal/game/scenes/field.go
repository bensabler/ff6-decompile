package scenes

import (
	"fmt"

	"github.com/bensabler/ff6-decompile/internal/content"
	"github.com/bensabler/ff6-decompile/internal/engine"
	"github.com/bensabler/ff6-decompile/internal/graphics/framebuf"
	"github.com/bensabler/ff6-decompile/internal/platform/snespad"
)

// FieldTiles renders the Narshe field BG tileset: 256 real 8x8 tiles,
// extracted uncompressed from ROMFILE:0x208460, the span EXP-0050 proved
// byte-identical to VRAM:$0000-$1FFF on SCN-0001 milestones 02 and 04.
//
// # What this is, and what it is not
//
// It **is** the first authentic field graphics the demo has ever drawn. Every
// pixel comes from the operator's ROM.
//
// It is **not** a field map. A map needs a tilemap to say which tile goes
// where, and the header that selects the tileset in the first place — and
// EXP-0051 searched four pointer encodings for that header and found none. So
// this scene shows the tile block in tile-number order, which is an honest
// presentation of what the project actually has: the graphics, without the
// arrangement.
//
// Drawing an invented arrangement would look far more like FF6 and would be
// exactly the "route hardcoded as one bespoke cinematic" that the acceptance
// criteria prohibit.
type FieldTiles struct {
	tiles *content.Tileset
	font  *content.Font
	page  int
	frame uint64
}

// fieldPages is how many 16x8-tile pages the 256-tile block splits into for
// display, so tiles are shown at native size rather than scaled.
const fieldPages = 2

func NewFieldTiles(tiles *content.Tileset, font *content.Font) *FieldTiles {
	return &FieldTiles{tiles: tiles, font: font}
}

func (f *FieldTiles) Update(ctx *engine.Context) {
	f.frame = ctx.Frame
	if ctx.Input.JustPressed(snespad.Right) {
		f.page = (f.page + 1) % fieldPages
	}
	if ctx.Input.JustPressed(snespad.Left) {
		f.page = (f.page + fieldPages - 1) % fieldPages
	}
	if ctx.Input.JustPressed(snespad.Start) {
		ctx.Stack.Pop()
	}
}

func (f *FieldTiles) Draw(dst *framebuf.Indexed, _ *framebuf.Palette) {
	const (
		primary   = content.SubPalettePrimary
		secondary = content.SubPaletteDim
		cols      = 16
		rowsShown = 8
	)
	dst.Fill(0)

	text := func(x, y int, s string, pal content.SubPalette) {
		f.font.DrawString(dst, x, y, s, content.TextOptions{Palette: pal})
	}

	text(8, 6, "NARSHE FIELD TILESET", primary)
	text(8, 18, "ROMFILE 0x208460  EXP-0050", secondary)

	// The tiles, at native size, centred.
	first := f.page * cols * rowsShown
	originX := (framebuf.Width - cols*8) / 2
	originY := 40
	for i := 0; i < cols*rowsShown; i++ {
		n := first + i
		if n >= content.TilesetTiles {
			break
		}
		f.tiles.DrawTile(dst, originX+(i%cols)*8, originY+(i/cols)*8, n,
			framebuf.BlitOptions{PaletteBase: framebuf.GrayFieldPaletteBase})
	}

	y := originY + rowsShown*8 + 12
	text(8, y, fmt.Sprintf("TILES %d - %d OF %d", first, first+cols*rowsShown-1, content.TilesetTiles), secondary)
	y += 12
	text(8, y, "UNCOMPRESSED - NO DECODER NEEDED", secondary)
	y += 12
	// The honest labels. Neither is a detail.
	text(8, y, "NO TILEMAP - MAP HEADER UNLOCATED", secondary)
	y += 12
	text(8, y, "PALETTE IS PROJECT AUTHORED  D6", secondary)

	text(8, 208, "LEFT RIGHT PAGE   START EXIT", secondary)
}
