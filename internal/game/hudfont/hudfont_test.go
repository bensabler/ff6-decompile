package hudfont

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func TestRelation(t *testing.T) {
	tests := []struct {
		b     byte
		vram  int
		sheet int
		off   int
	}{
		{0x80, 0x180, 129, 0x0487C0}, // 'A'
		{0x81, 0x181, 130, 0x0487D0}, // 'B'
		{0x99, 0x199, 154, 0x048950}, // 'Z'
		{0x9A, 0x19A, 155, 0x048960}, // 'a'
		{0xB4, 0x1B4, 181, 0x048B00}, // '0'
		{0xBF, 0x1BF, 192, 0x048BB0}, // '?'
		{0xC4, 0x1C4, 197, 0x048C00}, // '-'
		{0x00, 0x100, 1, 0x047FC0},
		{0xFF, 0x1FF, 256, 0x048FB0},
	}
	for _, tt := range tests {
		if got := VRAMTile(tt.b); got != tt.vram {
			t.Errorf("VRAMTile($%02X) = $%03X, want $%03X", tt.b, got, tt.vram)
		}
		if got := SheetIndex(tt.b); got != tt.sheet {
			t.Errorf("SheetIndex($%02X) = %d, want %d", tt.b, got, tt.sheet)
		}
		if got := ROMFileOffset(tt.b); got != tt.off {
			t.Errorf("ROMFileOffset($%02X) = 0x%06X, want 0x%06X", tt.b, got, tt.off)
		}
	}
}

// TestGlyphRangeFitsTheBlock checks the encoding cannot address a tile outside
// the 257-tile block. Sheet index 0 (VRAM $FF) is deliberately unreachable: it
// is in the block but below the glyph range.
func TestGlyphRangeFitsTheBlock(t *testing.T) {
	for i := 0; i < 256; i++ {
		idx := SheetIndex(byte(i))
		if idx < 0 || idx >= BlockTiles {
			t.Fatalf("byte $%02X -> sheet index %d, outside the %d-tile block", i, idx, BlockTiles)
		}
		off := ROMFileOffset(byte(i))
		if off < BlockROMFile || off+TileBytes > BlockROMFile+BlockTiles*TileBytes {
			t.Fatalf("byte $%02X -> ROMFILE:0x%06X, outside the block", i, off)
		}
	}
	if SheetIndex(0x00) != 1 {
		t.Error("byte $00 should map to sheet index 1; index 0 is VRAM $FF, below the glyph range")
	}
	if SheetIndex(0xFF) != BlockTiles-1 {
		t.Error("byte $FF should map to the last tile of the block")
	}
}

func TestClassifyNeverGuesses(t *testing.T) {
	// A character byte classifies as a character.
	if c, ch, conf := Classify(0x80, false); c != Character || ch == nil || *ch != "A" || conf != "confirmed" {
		t.Errorf("Classify($80) = %v, %v, %q", c, ch, conf)
	}
	// A blank tile is blank, and names no character.
	if c, ch, _ := Classify(0xEE, true); c != Blank || ch != nil {
		t.Errorf("Classify($EE, blank) = %v, %v", c, ch)
	}
	// A gauge tile is structural, and still names no character.
	if c, ch, conf := Classify(0xF9, false); c != Structural || ch != nil || conf != "confirmed" {
		t.Errorf("Classify($F9) = %v, %v, %q", c, ch, conf)
	}
	// Everything else stays unknown rather than being guessed from shape.
	for _, b := range []byte{0xC5, 0xCC, 0xD5, 0xE0} {
		c, ch, conf := Classify(b, false)
		if c != Unidentified || ch != nil || conf != "unknown" {
			t.Errorf("Classify($%02X) = %v, %v, %q; unidentified glyphs must not be guessed", b, c, ch, conf)
		}
	}
	// $FE/$FF decode to ' ' but their tiles are blank; they must not be
	// recorded as characters, or a renderer would draw a space glyph.
	for _, b := range []byte{0xFE, 0xFF} {
		if c, _, _ := Classify(b, true); c != Blank {
			t.Errorf("Classify($%02X, blank) = %v, want Blank", b, c)
		}
	}
}

func TestValidateRejectsDrift(t *testing.T) {
	good := make([]Glyph, 256)
	for i := range good {
		b := byte(i)
		good[i] = Glyph{
			EncodedByte:    b,
			VRAMTile:       VRAMTile(b),
			SheetIndex:     SheetIndex(b),
			ROMFileOffset:  ROMFileOffset(b),
			Classification: Unidentified,
			TileSHA256:     "0000000000000000000000000000000000000000000000000000000000000000",
			Confidence:     "unknown",
		}
	}
	if err := Validate(good); err != nil {
		t.Fatalf("a consistent table was rejected: %v", err)
	}

	bad := func(mutate func([]Glyph)) []Glyph {
		c := append([]Glyph(nil), good...)
		mutate(c)
		return c
	}
	cases := map[string][]Glyph{
		"wrong length":      good[:255],
		"out of order":      bad(func(g []Glyph) { g[5].EncodedByte = 6 }),
		"bad vram tile":     bad(func(g []Glyph) { g[5].VRAMTile = 0x999 }),
		"bad sheet index":   bad(func(g []Glyph) { g[5].SheetIndex = 999 }),
		"bad rom offset":    bad(func(g []Glyph) { g[5].ROMFileOffset = 0 }),
		"short hash":        bad(func(g []Glyph) { g[5].TileSHA256 = "abc" }),
		"character no rune": bad(func(g []Glyph) { g[5].Classification = Character }),
		"rune no character": bad(func(g []Glyph) { s := "A"; g[5].Character = &s }),
	}
	for name, tbl := range cases {
		if err := Validate(tbl); err == nil {
			t.Errorf("%s: Validate accepted an inconsistent table", name)
		}
	}
}

// TestTrackedGlyphMapIsConsistent validates the committed table. It needs no
// ROM: the table carries hashes, and this checks the relation and the schema.
func TestTrackedGlyphMapIsConsistent(t *testing.T) {
	p := filepath.Join(repoRoot(t), "data", "graphics", "hud-font-glyphs.json")
	raw, err := os.ReadFile(p)
	if err != nil {
		t.Skipf("glyph table not generated: %v", err)
	}
	var doc struct {
		Glyphs []Glyph `json:"glyphs"`
	}
	if err := json.Unmarshal(raw, &doc); err != nil {
		t.Fatalf("parsing %s: %v", p, err)
	}
	if err := Validate(doc.Glyphs); err != nil {
		t.Fatalf("tracked glyph table is inconsistent: %v", err)
	}

	// The characters the table names must agree with textenc, or the two
	// records of the encoding have diverged.
	var chars int
	for _, g := range doc.Glyphs {
		if g.Classification != Character {
			continue
		}
		chars++
		want, _, _ := Classify(g.EncodedByte, false)
		if want != Character {
			t.Errorf("byte $%02X is a character in the table but not in textenc", g.EncodedByte)
		}
	}
	if chars != 64 {
		t.Errorf("table names %d characters, want 64 (26+26+10+'-'+'?')", chars)
	}
}

func repoRoot(t *testing.T) string {
	t.Helper()
	dir, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	for i := 0; i < 6; i++ {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir
		}
		dir = filepath.Dir(dir)
	}
	t.Fatal("could not locate module root")
	return ""
}
