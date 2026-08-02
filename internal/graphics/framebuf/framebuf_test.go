package framebuf

import (
	"testing"
)

// checker is a synthetic tile: no ROM bytes.
func checker() *[8][8]uint8 {
	var t [8][8]uint8
	for y := 0; y < 8; y++ {
		for x := 0; x < 8; x++ {
			t[y][x] = uint8((x + y) % 4)
		}
	}
	return &t
}

// corners marks each corner with a distinct non-zero index so flips are
// unambiguous.
func corners() *[8][8]uint8 {
	var t [8][8]uint8
	t[0][0], t[0][7], t[7][0], t[7][7] = 1, 2, 3, 4
	return &t
}

func TestNewIsCleared(t *testing.T) {
	f := New()
	if len(f.Pix) != Width*Height {
		t.Fatalf("Pix length %d, want %d", len(f.Pix), Width*Height)
	}
	for i, v := range f.Pix {
		if v != 0 {
			t.Fatalf("pixel %d = %d, want 0", i, v)
		}
	}
}

func TestSetAtBoundsAreSafe(t *testing.T) {
	f := New()
	// Out-of-bounds writes are dropped, not wrapped: a wrapped write would
	// silently corrupt the opposite edge.
	for _, p := range [][2]int{{-1, 0}, {0, -1}, {Width, 0}, {0, Height}, {-99, -99}, {9999, 9999}} {
		f.Set(p[0], p[1], 9)
	}
	for i, v := range f.Pix {
		if v != 0 {
			t.Fatalf("out-of-bounds Set touched pixel %d", i)
		}
	}
	for _, p := range [][2]int{{-1, 0}, {Width, Height}} {
		if got := f.At(p[0], p[1]); got != Transparent {
			t.Errorf("At(%d,%d) = %d, want Transparent", p[0], p[1], got)
		}
	}
	f.Set(3, 4, 7)
	if got := f.At(3, 4); got != 7 {
		t.Errorf("At(3,4) = %d, want 7", got)
	}
}

func TestRectClips(t *testing.T) {
	tests := []struct {
		name                   string
		x, y, w, h             int
		wantSet                int
		probeX, probeY, probeV int
	}{
		{"interior", 10, 10, 4, 5, 20, 11, 12, 6},
		{"straddles left", -2, 0, 4, 1, 2, 0, 0, 6},
		{"straddles top", 0, -2, 1, 4, 2, 0, 0, 6},
		{"straddles right", Width - 2, 0, 4, 1, 2, Width - 1, 0, 6},
		{"straddles bottom", 0, Height - 2, 1, 4, 2, 0, Height - 1, 6},
		{"entirely outside", -50, -50, 10, 10, 0, 0, 0, 0},
		{"zero width", 5, 5, 0, 5, 0, 0, 0, 0},
		{"negative size", 5, 5, -3, -3, 0, 0, 0, 0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := New()
			f.Rect(tt.x, tt.y, tt.w, tt.h, 6)
			n := 0
			for _, v := range f.Pix {
				if v != 0 {
					n++
				}
			}
			if n != tt.wantSet {
				t.Errorf("set %d pixels, want %d", n, tt.wantSet)
			}
			if tt.wantSet > 0 {
				if got := f.At(tt.probeX, tt.probeY); got != uint8(tt.probeV) {
					t.Errorf("At(%d,%d) = %d, want %d", tt.probeX, tt.probeY, got, tt.probeV)
				}
			}
		})
	}
}

func TestBlitTileTransparency(t *testing.T) {
	f := New()
	f.Fill(9)
	f.BlitTile(0, 0, corners(), BlitOptions{})
	// Index 0 must not overwrite: only the four corners change.
	if got := f.At(1, 1); got != 9 {
		t.Errorf("transparent pixel overwrote destination: At(1,1) = %d, want 9", got)
	}
	if got := f.At(0, 0); got != 1 {
		t.Errorf("At(0,0) = %d, want 1", got)
	}

	f2 := New()
	f2.Fill(9)
	f2.BlitTile(0, 0, corners(), BlitOptions{Opaque: true})
	if got := f2.At(1, 1); got != 0 {
		t.Errorf("Opaque blit should write index 0: At(1,1) = %d, want 0", got)
	}
}

func TestBlitTileFlips(t *testing.T) {
	tests := []struct {
		name           string
		o              BlitOptions
		wantTL, wantTR uint8
		wantBL, wantBR uint8
	}{
		{"none", BlitOptions{}, 1, 2, 3, 4},
		{"h", BlitOptions{FlipH: true}, 2, 1, 4, 3},
		{"v", BlitOptions{FlipV: true}, 3, 4, 1, 2},
		{"hv", BlitOptions{FlipH: true, FlipV: true}, 4, 3, 2, 1},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := New()
			f.BlitTile(0, 0, corners(), tt.o)
			got := [4]uint8{f.At(0, 0), f.At(7, 0), f.At(0, 7), f.At(7, 7)}
			want := [4]uint8{tt.wantTL, tt.wantTR, tt.wantBL, tt.wantBR}
			if got != want {
				t.Errorf("corners = %v, want %v", got, want)
			}
		})
	}
}

func TestBlitTilePaletteBase(t *testing.T) {
	f := New()
	f.BlitTile(0, 0, corners(), BlitOptions{PaletteBase: 0x30})
	if got := f.At(0, 0); got != 0x31 {
		t.Errorf("At(0,0) = $%02X, want $31 (1 + base $30)", got)
	}
	// Transparent pixels must not have the base added, or a skipped pixel
	// would become a visible palette-base-colored one.
	if got := f.At(1, 1); got != 0 {
		t.Errorf("At(1,1) = $%02X, want 0: transparency is decided before the base is added", got)
	}
}

func TestBlitTileClipsAtEveryEdge(t *testing.T) {
	for _, p := range [][2]int{
		{-4, -4}, {Width - 4, -4}, {-4, Height - 4}, {Width - 4, Height - 4},
		{-8, 0}, {Width, 0}, {0, -8}, {0, Height},
		{-1000, -1000}, {1000, 1000},
	} {
		f := New()
		f.BlitTile(p[0], p[1], checker(), BlitOptions{Opaque: true})
		// The only assertion that matters: it did not panic, and it wrote
		// nothing outside the buffer (guaranteed by Pix's length).
		_ = f
	}
}

func TestResolveAndPaletted(t *testing.T) {
	f := New()
	f.Set(0, 0, 3)
	f.Set(1, 0, 2)
	p := GrayPalette()

	img := f.Resolve(nil, p)
	if got := img.RGBAAt(0, 0); got.R != 0xFF || got.A != 0xFF {
		t.Errorf("Resolve pixel (0,0) = %v, want opaque white", got)
	}
	if got := img.RGBAAt(1, 0); got.R != 0xA5 {
		t.Errorf("Resolve pixel (1,0) = %v, want the mid ramp entry", got)
	}
	// Reusing a correctly sized destination must not reallocate.
	again := f.Resolve(img, p)
	if again != img {
		t.Error("Resolve reallocated a correctly sized destination")
	}

	pal := f.Paletted(p)
	if got := pal.ColorIndexAt(0, 0); got != 3 {
		t.Errorf("Paletted index (0,0) = %d, want 3", got)
	}
	if got := pal.ColorIndexAt(5, 5); got != 0 {
		t.Errorf("Paletted index (5,5) = %d, want 0", got)
	}
}

func TestSum256TracksCompositionNotPalette(t *testing.T) {
	a, b := New(), New()
	if a.Sum256() != b.Sum256() {
		t.Fatal("identical framebuffers hashed differently")
	}
	b.Set(100, 100, 1)
	if a.Sum256() == b.Sum256() {
		t.Fatal("a changed pixel did not change the hash")
	}
	// The hash is over indices, so it is independent of the palette by
	// construction; this documents that intent.
	c := New()
	c.Set(100, 100, 1)
	if c.Sum256() != b.Sum256() {
		t.Fatal("same indices must hash the same regardless of how they are resolved")
	}
}

func FuzzBlitTile(f *testing.F) {
	f.Add(0, 0, false, false, false, uint8(0))
	f.Add(-1000, 1000, true, true, true, uint8(255))
	f.Fuzz(func(t *testing.T, x, y int, fh, fv, op bool, base uint8) {
		fb := New()
		fb.BlitTile(x, y, checker(), BlitOptions{PaletteBase: base, FlipH: fh, FlipV: fv, Opaque: op})
		if len(fb.Pix) != Width*Height {
			t.Fatal("BlitTile resized the framebuffer")
		}
	})
}

func FuzzRect(f *testing.F) {
	f.Add(0, 0, 1, 1, uint8(1))
	f.Add(-99999, -99999, 99999, 99999, uint8(7))
	f.Fuzz(func(t *testing.T, x, y, w, h int, idx uint8) {
		fb := New()
		fb.Rect(x, y, w, h, idx)
		if len(fb.Pix) != Width*Height {
			t.Fatal("Rect resized the framebuffer")
		}
	})
}
