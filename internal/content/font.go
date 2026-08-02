package content

import (
	"bytes"
	"fmt"
	"image"
	"image/png"

	"github.com/bensabler/ff6-decompile/internal/game/hudfont"
	"github.com/bensabler/ff6-decompile/internal/game/textenc"
	"github.com/bensabler/ff6-decompile/internal/graphics/framebuf"
)

// HUDFontAssetID is the manifest entry for the battle HUD font sheet.
const HUDFontAssetID = "ASSET-GFX-0001"

// GlyphWidth and GlyphHeight are the fixed cell size of the HUD font. It is
// a fixed-width font; the variable-width dialogue font is a different,
// currently unlocated system (CEN-GFX-0004).
const (
	GlyphWidth  = 8
	GlyphHeight = 8
)

// Font is the battle HUD font, ready to draw.
//
// Tiles are indexed by the *encoding byte*, via the EXP-0049 relation, so
// callers never handle sheet indices. present[b] records whether that byte's
// tile carries ink, which is what lets DrawString distinguish "a space" from
// "a glyph this project has not identified".
type Font struct {
	tiles   [256][8][8]uint8
	present [256]bool
}

// LoadHUDFont reads the font sheet from the archive.
//
// The sheet is an indexed PNG, which the round-trip through our own encoder
// preserves exactly. Decoding it back to indices is therefore lossless, and
// avoids adding a second on-disk format purely for the renderer's
// convenience.
func LoadHUDFont(a *Archive) (*Font, error) {
	data, err := a.Read(HUDFontAssetID)
	if err != nil {
		return nil, err
	}
	img, err := png.Decode(bytes.NewReader(data))
	if err != nil {
		return nil, fmt.Errorf("decoding the HUD font sheet: %w", err)
	}
	pal, ok := img.(*image.Paletted)
	if !ok {
		return nil, fmt.Errorf("the HUD font sheet decoded as %T, want *image.Paletted; the extractor writes an indexed PNG", img)
	}

	const cols = 16
	f := &Font{}
	for i := 0; i < 256; i++ {
		b := byte(i)
		sheet := hudfont.SheetIndex(b)
		ox, oy := (sheet%cols)*GlyphWidth, (sheet/cols)*GlyphHeight
		if ox+GlyphWidth > pal.Bounds().Dx() || oy+GlyphHeight > pal.Bounds().Dy() {
			return nil, fmt.Errorf("glyph $%02X wants sheet cell %d at (%d,%d), outside the %dx%d sheet",
				b, sheet, ox, oy, pal.Bounds().Dx(), pal.Bounds().Dy())
		}
		ink := false
		for y := 0; y < GlyphHeight; y++ {
			for x := 0; x < GlyphWidth; x++ {
				v := pal.ColorIndexAt(pal.Bounds().Min.X+ox+x, pal.Bounds().Min.Y+oy+y)
				f.tiles[i][y][x] = v
				if v != 0 {
					ink = true
				}
			}
		}
		f.present[i] = ink
	}
	return f, nil
}

// NewFont builds a font from tiles indexed by encoding byte, however those
// tiles were obtained.
//
// Tests use it to build a font from project-authored pixels, so frame goldens
// can be committed without containing ROM-derived graphics.
func NewFont(tiles *[256][8][8]uint8) *Font {
	f := &Font{tiles: *tiles}
	for i := range f.tiles {
		for y := range f.tiles[i] {
			for x := range f.tiles[i][y] {
				if f.tiles[i][y][x] != 0 {
					f.present[i] = true
				}
			}
		}
	}
	return f
}

// Glyph returns the tile for an encoding byte. ok is false when the tile is
// blank, which for an unmapped byte means there is nothing to draw.
func (f *Font) Glyph(b byte) (*[8][8]uint8, bool) {
	return &f.tiles[b], f.present[b]
}

// SubPalette selects one of the 64 four-color sub-palettes a 2bpp tile can
// use within the 256-entry CGRAM image. It is a **slot number**, not a
// palette index and not a byte offset: slot n covers entries 4n..4n+3, and
// the font's ink values 1-3 land on 4n+1..4n+3.
//
// The type exists because the units were the defect. TextOptions used to take
// a raw uint8 base, and both scenes passed 3 for "white" and 2 for "gray" —
// reading the field as a brightness level. Those are not sub-palette bases,
// and the ink went to entries the palette had never defined, so every bright
// string in the demo rendered black on black for two units. Expressed as a
// slot the same mistake is caught: slot 3 is entries 12-15, far outside what
// framebuf.GrayPalette defines, and TestScenesDrawOnlyDefinedPaletteEntries
// fails.
type SubPalette uint8

// Sub-palettes framebuf.GrayPalette defines. Real FF6 palettes are readiness
// row F2 and still unlocated; when they land, these names stay and only the
// colors behind them change.
const (
	// SubPalettePrimary is entries 0-3: transparent, black, gray, white.
	SubPalettePrimary SubPalette = 0
	// SubPaletteDim is entries 4-7, the same shape one step darker, for
	// secondary text that should read as subordinate rather than absent.
	SubPaletteDim SubPalette = 1
)

// Base returns the first CGRAM entry of the sub-palette, which is what
// framebuf.BlitOptions.PaletteBase wants.
func (s SubPalette) Base() uint8 { return uint8(s) * 4 }

// TextOptions controls text rendering.
type TextOptions struct {
	// Palette selects the four-color sub-palette the glyph's ink resolves
	// through. The zero value is SubPalettePrimary.
	Palette SubPalette
	// ShowMissing draws a project-authored box for runes with no verified
	// glyph, instead of skipping them. Off by default so ordinary text
	// renders cleanly; on for debugging a decode.
	ShowMissing bool
}

// DrawString renders s at (x, y) with the fixed 8-pixel advance.
//
// Runes the encoding has no verified value for are **not** substituted with
// a ROM tile — that would put an invented glyph on screen and make a decode
// error look like content. They advance the cursor, and with ShowMissing set
// they draw a box this project authored.
//
// It returns the x coordinate just past the last cell.
func (f *Font) DrawString(dst *framebuf.Indexed, x, y int, s string, o TextOptions) int {
	for _, r := range s {
		b, ok := textenc.Encode(r)
		switch {
		case !ok:
			if o.ShowMissing {
				drawMissingBox(dst, x, y, o.Palette.Base())
			}
		default:
			if tile, ink := f.Glyph(b); ink {
				dst.BlitTile(x, y, tile, framebuf.BlitOptions{PaletteBase: o.Palette.Base()})
			}
		}
		x += GlyphWidth
	}
	return x
}

// DrawBytes renders already-encoded bytes, for text read straight out of a
// ROM-derived table. Bytes whose tiles are blank draw nothing, which is
// correct for the padding values $FE and $FF.
func (f *Font) DrawBytes(dst *framebuf.Indexed, x, y int, src []byte, o TextOptions) int {
	for _, b := range src {
		if tile, ink := f.Glyph(b); ink {
			dst.BlitTile(x, y, tile, framebuf.BlitOptions{PaletteBase: o.Palette.Base()})
		} else if o.ShowMissing {
			if _, named := textenc.Decode(b); !named {
				drawMissingBox(dst, x, y, o.Palette.Base())
			}
		}
		x += GlyphWidth
	}
	return x
}

// MeasureString returns the pixel width of s in this fixed-width font.
func MeasureString(s string) int { return len([]rune(s)) * GlyphWidth }

// drawMissingBox draws a hollow 6x6 box. It is project-authored, so it can
// never be mistaken for FF6 content.
func drawMissingBox(dst *framebuf.Indexed, x, y int, base uint8) {
	const idx = 1
	c := idx + base
	dst.Rect(x+1, y+1, 6, 1, c)
	dst.Rect(x+1, y+6, 6, 1, c)
	dst.Rect(x+1, y+1, 1, 6, c)
	dst.Rect(x+6, y+1, 1, 6, c)
}
