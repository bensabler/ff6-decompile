// Package framebuf provides the demo's render target: a palette-indexed
// SNES-resolution framebuffer.
//
// It lives under internal/graphics rather than internal/engine on purpose.
// Every FF6 compositor this project will grow — background layers, OAM,
// color math, windowing — belongs in internal/graphics and must be able to
// write into the render target. If the target lived in engine, graphics would
// have to import engine, which is an upward edge the layering forbids.
//
// # Why indexed rather than RGBA
//
// One uint8 per pixel is a lossless model of the SNES output stage: CGRAM is
// exactly 256 BGR555 entries. tile2bpp.Tile and tile4bpp.Tile are already
// [8][8]uint8 of palette indices, so blitting is an index copy with no color
// conversion in the hot path. Fades, brightness and color math are palette
// transforms, which stay possible. And a hash over the indices tests
// composition independently of color, where an RGBA hash conflates the two.
package framebuf

import (
	"crypto/sha256"
	"image"
	"image/color"

	"github.com/bensabler/ff6-decompile/internal/platform/bgr555"
)

// SNES active display, in pixels. 256x224 is the NTSC non-interlaced mode
// FF6 uses.
const (
	Width  = 256
	Height = 224
)

// Transparent is the palette index treated as "leave the destination alone"
// by BlitTile. It matches the SNES convention that color 0 of a background
// palette is transparent.
const Transparent uint8 = 0

// Palette is a CGRAM image: 256 BGR555 entries.
type Palette [256]bgr555.Color

// Indexed is a Width x Height framebuffer of palette indices, row-major.
type Indexed struct {
	Pix []uint8
}

// New returns a cleared framebuffer.
func New() *Indexed { return &Indexed{Pix: make([]uint8, Width*Height)} }

// Fill sets every pixel to idx.
func (f *Indexed) Fill(idx uint8) {
	for i := range f.Pix {
		f.Pix[i] = idx
	}
}

// Set writes one pixel, ignoring out-of-bounds coordinates.
func (f *Indexed) Set(x, y int, idx uint8) {
	if x < 0 || y < 0 || x >= Width || y >= Height {
		return
	}
	f.Pix[y*Width+x] = idx
}

// At reads one pixel. Out-of-bounds reads return Transparent rather than
// panicking, so drawing code can probe freely.
func (f *Indexed) At(x, y int) uint8 {
	if x < 0 || y < 0 || x >= Width || y >= Height {
		return Transparent
	}
	return f.Pix[y*Width+x]
}

// Rect fills an axis-aligned rectangle, clipped to the framebuffer.
func (f *Indexed) Rect(x, y, w, h int, idx uint8) {
	if w <= 0 || h <= 0 {
		return
	}
	x0, y0 := max(x, 0), max(y, 0)
	x1, y1 := min(x+w, Width), min(y+h, Height)
	for py := y0; py < y1; py++ {
		row := f.Pix[py*Width:]
		for px := x0; px < x1; px++ {
			row[px] = idx
		}
	}
}

// BlitOptions controls how a tile is drawn.
type BlitOptions struct {
	// PaletteBase is added to each non-transparent source index, selecting
	// a 4- or 16-color sub-palette within the 256-entry CGRAM image.
	//
	// It is the **first CGRAM entry of a sub-palette**, not an arbitrary
	// brightness offset: a multiple of 4 for 2bpp tiles, of 16 for 4bpp.
	// A 2bpp tile's ink values 1-3 land on PaletteBase+1..PaletteBase+3, so
	// a base that is not a multiple of 4 straddles two sub-palettes and
	// scatters one glyph's ink across colors that were never chosen
	// together. Use IsSubPaletteBase to check a value.
	PaletteBase uint8
	// FlipH and FlipV mirror the tile, as the SNES tilemap attribute bits
	// do.
	FlipH, FlipV bool
	// Opaque draws index 0 instead of skipping it.
	Opaque bool
}

// BlitTile draws an 8x8 tile of palette indices at (x, y), clipped to the
// framebuffer. It never panics for any coordinate.
func (f *Indexed) BlitTile(x, y int, tile *[8][8]uint8, o BlitOptions) {
	for ty := 0; ty < 8; ty++ {
		py := y + ty
		if py < 0 || py >= Height {
			continue
		}
		sy := ty
		if o.FlipV {
			sy = 7 - ty
		}
		row := f.Pix[py*Width:]
		for tx := 0; tx < 8; tx++ {
			px := x + tx
			if px < 0 || px >= Width {
				continue
			}
			sx := tx
			if o.FlipH {
				sx = 7 - tx
			}
			v := tile[sy][sx]
			if v == Transparent && !o.Opaque {
				continue
			}
			row[px] = v + o.PaletteBase
		}
	}
}

// Resolve writes the framebuffer into an RGBA image through a palette. dst
// must be Width x Height; Resolve allocates one if dst is nil.
func (f *Indexed) Resolve(dst *image.RGBA, p *Palette) *image.RGBA {
	if dst == nil || dst.Bounds().Dx() != Width || dst.Bounds().Dy() != Height {
		dst = image.NewRGBA(image.Rect(0, 0, Width, Height))
	}
	for i, idx := range f.Pix {
		c := p[idx]
		o := i * 4
		dst.Pix[o+0] = c.R
		dst.Pix[o+1] = c.G
		dst.Pix[o+2] = c.B
		dst.Pix[o+3] = 0xFF
	}
	return dst
}

// Paletted returns the framebuffer as an image.Paletted, for PNG encoding.
func (f *Indexed) Paletted(p *Palette) *image.Paletted {
	pal := make(color.Palette, len(p))
	for i, c := range p {
		pal[i] = color.RGBA{c.R, c.G, c.B, 0xFF}
	}
	img := image.NewPaletted(image.Rect(0, 0, Width, Height), pal)
	copy(img.Pix, f.Pix)
	return img
}

// Sum256 hashes the palette indices. This is the demo's golden-frame
// identity: it tests composition, and is deliberately independent of the
// palette so a color change and a layout change fail differently.
func (f *Indexed) Sum256() [32]byte { return sha256.Sum256(f.Pix) }

// IsSubPaletteBase reports whether base is the first entry of a sub-palette
// holding colors entries (4 for 2bpp tiles, 16 for 4bpp).
//
// The check exists because BlitTile cannot make it: it receives [8][8]uint8
// and has no idea what depth produced it. The caller knows, so the caller
// checks.
func IsSubPaletteBase(base uint8, colors int) bool {
	if colors <= 0 || colors > 256 {
		return false
	}
	return int(base)%colors == 0
}

// GrayPaletteDefined is the number of leading entries GrayPalette assigns a
// color to — two four-entry sub-palettes. Entries at or above it are black,
// so a tile drawn onto them disappears against a cleared background.
//
// This constant exists to be asserted against. Sum256 hashes indices and is
// deliberately palette-independent, which means the frame goldens cannot see
// a scene drawing onto undefined entries — and for two units, both of them
// did. Callers name their sub-palette through content.SubPalette; this is
// the bound that catches one that does not exist.
const GrayPaletteDefined = 8

// GrayPalette is a project-authored ramp for tests and for rendering before a
// real FF6 palette is loaded. It defines two four-entry sub-palettes,
// GrayPaletteDefined entries in total; everything above is black.
//
// It contains no ROM-derived color data. Entries 1-3 do coincide with the
// BG 2bpp palette 0 that EXP-0023 measured (BGR555 $0000/$5294/$7FFF) — but
// black, mid-gray and white is the obvious choice for a three-level ramp, and
// this palette is not evidence for anything. The real map and battle palettes
// are readiness row F2 and are still unlocated.
func GrayPalette() *Palette {
	var p Palette
	for i := range p {
		p[i] = bgr555.Color{R: 0, G: 0, B: 0}
	}
	// Sub-palette 0: index 0 transparent, then the font's three ink levels.
	p[1] = bgr555.Color{R: 0x00, G: 0x00, B: 0x00}
	p[2] = bgr555.Color{R: 0xA5, G: 0xA5, B: 0xA5}
	p[3] = bgr555.Color{R: 0xFF, G: 0xFF, B: 0xFF}
	// Sub-palette 1: the same shape, dimmed.
	p[5] = bgr555.Color{R: 0x00, G: 0x00, B: 0x00}
	p[6] = bgr555.Color{R: 0x52, G: 0x52, B: 0x52}
	p[7] = bgr555.Color{R: 0xA5, G: 0xA5, B: 0xA5}
	return &p
}
