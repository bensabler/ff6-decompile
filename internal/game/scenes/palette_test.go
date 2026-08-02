package scenes

import (
	"testing"

	"github.com/bensabler/ff6-decompile/internal/content"
	"github.com/bensabler/ff6-decompile/internal/graphics/framebuf"
)

// The assertion this file exists to make.
//
// framebuf.BlitTile computes dst = ink + PaletteBase, so PaletteBase is the
// first entry of a sub-palette, not a brightness offset. Both scenes used to
// pass 3 for "white" and 2 for "gray". The HUD font's ink values are 1-3, so
// a base of 3 put every pixel of every bright string on palette entries 4-6,
// and GrayPalette defined only 0-3. Every "white" string in the shipped demo
// rendered black on black; every "gray" string drew one of its three ink
// levels and dropped the other two.
//
// Nothing caught it for two units. Indexed.Sum256 hashes palette indices and
// is deliberately palette-independent — that is the right design for a
// composition golden, and it is exactly why the goldens could not see this.
// A frame can hash correctly and be blank on screen.
//
// So the assertion has to be made against the palette, once per scene.

// drawnIndices returns the set of palette indices a frame actually contains.
func drawnIndices(pix []uint8) map[uint8]int {
	seen := make(map[uint8]int)
	for _, v := range pix {
		seen[v]++
	}
	return seen
}

func TestScenesDrawOnlyDefinedPaletteEntries(t *testing.T) {
	tables := syntheticTables(t)
	battle, err := NewBattle(syntheticFont(), tables, 1)
	if err != nil {
		t.Fatal(err)
	}

	// A synthetic 4bpp tileset: ink across the full 0-15 range, which is
	// what makes the field scene exercise the 16-colour sub-palette.
	var tiles [content.TilesetTiles][8][8]uint8
	for i := range tiles {
		for y := 0; y < 8; y++ {
			for x := 0; x < 8; x++ {
				tiles[i][y][x] = uint8((i + x + y) % 16)
			}
		}
	}

	cases := []struct {
		name  string
		scene interface {
			Draw(*framebuf.Indexed, *framebuf.Palette)
		}
	}{
		{"boot", NewBoot(syntheticFont())},
		{"battle", battle},
		{"field tiles", NewFieldTiles(content.NewTileset(&tiles), syntheticFont())},
	}

	pal := framebuf.GrayPalette()
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			fb := framebuf.New()
			tc.scene.Draw(fb, pal)

			for idx, n := range drawnIndices(fb.Pix) {
				if int(idx) >= framebuf.GrayPaletteDefined {
					t.Errorf("index %d drawn on %d pixels, but GrayPalette defines only entries 0-%d; "+
						"those pixels resolve to black and are invisible. A PaletteBase is a "+
						"sub-palette base (a multiple of 4 for 2bpp), not a brightness level",
						idx, n, framebuf.GrayPaletteDefined-1)
				}
			}
		})
	}
}

// TestMostDrawnInkIsVisible asserts the ratio the defect actually moved:
// of the pixels a scene draws, how many resolve to a color different from
// the background.
//
// It is a ratio rather than a count because the count is a layout golden and
// would need regenerating on every text change, while the ratio is a property
// of the palette mapping.
//
// Measured on the boot scene with the synthetic font, both runs drawing an
// identical 6406-pixel ink mask:
//
//	before the fix: 2238 of 6406 visible (35%) — indices 4, 5 and 6 carried
//	                4168 pixels and the palette defined only 0-3
//	after:          4306 of 6406 visible (67%)
//
// Ink that stays invisible is legitimate: the font's darkest level is the
// glyph outline, and it is black in the sub-palette EXP-0023 measured too.
// What is not legitimate is two thirds of the ink disappearing.
func TestMostDrawnInkIsVisible(t *testing.T) {
	pal := framebuf.GrayPalette()
	bg := pal[0]

	fb := framebuf.New()
	NewBoot(syntheticFont()).Draw(fb, pal)

	drawn, visible := 0, 0
	for _, idx := range fb.Pix {
		if idx == 0 {
			continue
		}
		drawn++
		if pal[idx] != bg {
			visible++
		}
	}
	if drawn == 0 {
		t.Fatal("the boot scene drew nothing")
	}

	// Half. The pre-fix value was 35%, the post-fix value 67%, so this
	// discriminates between them without pinning either.
	if visible*2 < drawn {
		t.Errorf("only %d of %d drawn pixels resolve to a non-background color (%d%%); want over half. "+
			"Ink is landing on palette entries that hold the background color",
			visible, drawn, visible*100/drawn)
	}
}

func TestSubPaletteBases(t *testing.T) {
	tests := []struct {
		name   string
		base   uint8
		colors int
		want   bool
	}{
		{"2bpp slot 0", 0, 4, true},
		{"2bpp slot 1", 4, 4, true},
		{"2bpp slot 63", 252, 4, true},
		{"the old white", 3, 4, false},
		{"the old gray", 2, 4, false},
		{"4bpp slot 0", 0, 16, true},
		{"4bpp slot 1", 16, 16, true},
		{"4bpp misaligned", 4, 16, false},
		{"zero colors", 0, 0, false},
		{"negative colors", 0, -4, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := framebuf.IsSubPaletteBase(tt.base, tt.colors); got != tt.want {
				t.Errorf("IsSubPaletteBase(%d, %d) = %v, want %v", tt.base, tt.colors, got, tt.want)
			}
		})
	}

	// The constants the scenes use must satisfy the contract they document.
	for _, sp := range []content.SubPalette{content.SubPalettePrimary, content.SubPaletteDim} {
		base := sp.Base()
		if !framebuf.IsSubPaletteBase(base, 4) {
			t.Errorf("sub-palette base %d is not a valid 2bpp base", base)
		}
		if int(base)+4 > framebuf.GrayPaletteDefined {
			t.Errorf("sub-palette base %d runs past GrayPalette's %d defined entries",
				base, framebuf.GrayPaletteDefined)
		}
	}
}
