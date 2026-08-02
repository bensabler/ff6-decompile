package snesoam

import "testing"

// oamWith builds an OAM image with one sprite set, and everything else parked
// off-screen the way FF6 leaves unused slots.
func oamWith(idx int, x, y, tileLo, attr byte, xSign, large bool) []byte {
	o := make([]byte, Size)
	for i := 0; i < SpriteCount; i++ {
		o[i*4+1] = 240 // parked below the display
	}
	o[idx*4+0], o[idx*4+1], o[idx*4+2], o[idx*4+3] = x, y, tileLo, attr
	var bits byte
	if xSign {
		bits |= 1
	}
	if large {
		bits |= 2
	}
	o[LowTableSize+idx/4] |= bits << uint((idx%4)*2)
	return o
}

func TestDecodeAttributes(t *testing.T) {
	// attr = vhoopppN: vflip, hflip, priority 2, palette 7, tile bit 8.
	const attr = 0x80 | 0x40 | (2 << 4) | (7 << 1) | 1
	s, err := Decode(oamWith(53, 120, 81, 0xAC, attr, false, false), 3)
	if err != nil {
		t.Fatal(err)
	}
	g := s[53]
	if g.X != 120 || g.Y != 81 {
		t.Errorf("position (%d,%d), want (120,81)", g.X, g.Y)
	}
	if g.Tile != 0x1AC {
		t.Errorf("Tile = $%03X, want $1AC", g.Tile)
	}
	if g.Palette != 7 || g.Priority != 2 || !g.FlipH || !g.FlipV {
		t.Errorf("attrs: pal=%d pri=%d h=%v v=%v", g.Palette, g.Priority, g.FlipH, g.FlipV)
	}
	// Sprite palettes live in the upper half of CGRAM.
	if got := g.CGRAMBase(); got != 128+7*16 {
		t.Errorf("CGRAMBase = %d, want %d", got, 128+7*16)
	}
}

// The high table is where a hand-rolled reader goes wrong, so both bits get a
// test of their own.
func TestHighTableBits(t *testing.T) {
	// X sign: 200 with the sign set is -56, not 200.
	s, _ := Decode(oamWith(1, 200, 10, 0, 0, true, false), 3)
	if s[1].X != -56 {
		t.Errorf("X with sign = %d, want -56", s[1].X)
	}
	if s2, _ := Decode(oamWith(1, 200, 10, 0, 0, false, false), 3); s2[1].X != 200 {
		t.Errorf("X without sign = %d, want 200", s2[1].X)
	}

	// Size bit, against mode 3: small 16x16, large 32x32.
	small, _ := Decode(oamWith(2, 0, 0, 0, 0, false, false), 3)
	large, _ := Decode(oamWith(2, 0, 0, 0, 0, false, true), 3)
	if small[2].Size != (Dimensions{16, 16}) {
		t.Errorf("small in mode 3 = %v, want 16x16", small[2].Size)
	}
	if large[2].Size != (Dimensions{32, 32}) {
		t.Errorf("large in mode 3 = %v, want 32x32", large[2].Size)
	}

	// The bits must be read from the right slot within the packed byte.
	for _, idx := range []int{0, 1, 2, 3, 4, 127} {
		s, _ := Decode(oamWith(idx, 200, 10, 0, 0, true, true), 3)
		if s[idx].X != -56 || !s[idx].Large {
			t.Errorf("slot %d decoded x=%d large=%v", idx, s[idx].X, s[idx].Large)
		}
		for j := 0; j < SpriteCount; j++ {
			if j != idx && (s[j].X < 0 || s[j].Large) {
				t.Fatalf("slot %d leaked high bits into slot %d", idx, j)
			}
		}
	}
}

func TestSizesForMode(t *testing.T) {
	tests := []struct {
		mode         uint8
		small, large Dimensions
	}{
		{0, Dimensions{8, 8}, Dimensions{16, 16}},
		{3, Dimensions{16, 16}, Dimensions{32, 32}},
		{5, Dimensions{32, 32}, Dimensions{64, 64}},
		{6, Dimensions{16, 32}, Dimensions{32, 64}},
		{7, Dimensions{16, 32}, Dimensions{32, 32}},
	}
	for _, tt := range tests {
		s, l, err := SizesForMode(tt.mode)
		if err != nil || s != tt.small || l != tt.large {
			t.Errorf("mode %d = %v/%v (%v), want %v/%v", tt.mode, s, l, err, tt.small, tt.large)
		}
	}
	if _, _, err := SizesForMode(8); err == nil {
		t.Error("mode 8 should be rejected")
	}
}

// The name table is 16 tiles wide, so the tile below is +16 and not +cols.
// Getting this wrong builds a plausible-looking sprite from the wrong tiles.
func TestTileAtUsesA16WideNameTable(t *testing.T) {
	s := Sprite{Tile: 0x1AC, Size: Dimensions{16, 16}}
	tests := []struct {
		c, r int
		want int
	}{
		{0, 0, 0x1AC},
		{1, 0, 0x1AD},
		{0, 1, 0x1BC},
		{1, 1, 0x1BD},
	}
	for _, tt := range tests {
		if got := s.TileAt(tt.c, tt.r); got != tt.want {
			t.Errorf("TileAt(%d,%d) = $%03X, want $%03X", tt.c, tt.r, got, tt.want)
		}
	}
	// Names are 9 bits and wrap.
	if got := (Sprite{Tile: 0x1FF}).TileAt(1, 0); got != 0x000 {
		t.Errorf("name wrap = $%03X, want $000", got)
	}
}

func TestOnScreen(t *testing.T) {
	tests := []struct {
		name string
		x, y int
		want bool
	}{
		{"centre", 120, 81, true},
		{"parked below, as FF6 leaves unused slots", 0, 240, false},
		{"partly off the left", -8, 100, true},
		{"fully off the left", -16, 100, false},
		{"past the right edge", 256, 100, false},
		{"straddling the right edge", 250, 100, true},
	}
	for _, tt := range tests {
		s := Sprite{X: tt.x, Y: tt.y, Size: Dimensions{16, 16}}
		if got := s.OnScreen(); got != tt.want {
			t.Errorf("%s: OnScreen = %v, want %v", tt.name, got, tt.want)
		}
	}
}

func TestOnScreenFilter(t *testing.T) {
	// One sprite placed, 127 parked: the filter is what makes a capture
	// readable at all.
	s, _ := Decode(oamWith(53, 120, 81, 0, 0, false, false), 3)
	vis := OnScreen(s)
	if len(vis) != 1 || vis[0].Index != 53 {
		t.Errorf("got %d visible sprites, want 1 (index 53)", len(vis))
	}
}

func TestVRAMAddress(t *testing.T) {
	// The measured case: milestone 04's base is VRAM:$C000, and tile $1AC
	// sits at $F580.
	if got := VRAMAddress(0xC000, 0x1AC); got != 0xF580 {
		t.Errorf("VRAMAddress = $%04X, want $F580", got)
	}
	if got := VRAMAddress(0xC000, 0x000); got != 0xC000 {
		t.Errorf("VRAMAddress = $%04X, want $C000", got)
	}
	// Names wrap within the 64 KB space rather than running past it.
	if got := VRAMAddress(0xF000, 0x1FF); got >= 0x10000 {
		t.Errorf("VRAMAddress = $%X, should wrap into VRAM", got)
	}
}

func TestDecodeRejectsBadInput(t *testing.T) {
	for _, n := range []int{0, 512, 543, 545} {
		if _, err := Decode(make([]byte, n), 3); err == nil {
			t.Errorf("a %d-byte OAM should be rejected", n)
		}
	}
	if _, err := Decode(make([]byte, Size), 9); err == nil {
		t.Error("an out-of-range object mode should be rejected")
	}
}

func FuzzDecode(f *testing.F) {
	f.Add(make([]byte, Size), uint8(3))
	f.Fuzz(func(t *testing.T, b []byte, mode uint8) {
		s, err := Decode(b, mode)
		if err != nil {
			return
		}
		if len(s) != SpriteCount {
			t.Fatalf("got %d sprites, want %d", len(s), SpriteCount)
		}
		for i, sp := range s {
			if sp.Index != i {
				t.Fatalf("sprite %d reports index %d", i, sp.Index)
			}
			if sp.Tile < 0 || sp.Tile > 0x1FF {
				t.Fatalf("tile $%X out of the 9-bit range", sp.Tile)
			}
			if sp.Palette > 7 || sp.Priority > 3 {
				t.Fatalf("pal %d pri %d out of range", sp.Palette, sp.Priority)
			}
			if sp.X < -256 || sp.X > 255 {
				t.Fatalf("X %d out of range", sp.X)
			}
			for c := 0; c < 8; c++ {
				for r := 0; r < 8; r++ {
					if n := sp.TileAt(c, r); n < 0 || n > 0x1FF {
						t.Fatalf("TileAt(%d,%d) = $%X out of range", c, r, n)
					}
				}
			}
		}
	})
}
