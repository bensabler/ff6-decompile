// This file is an external test package on purpose.
//
// internal/content must never import internal/rom — that boundary is what
// makes "the demo binary cannot read a ROM" a structural fact rather than a
// convention, and `ff6lab audit` enforces it. A differential test, though,
// has to read both sides. package content_test is not compiled into
// internal/content, so it can import both without weakening the boundary.
package content_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/bensabler/ff6-decompile/internal/content"
	"github.com/bensabler/ff6-decompile/internal/game/hudfont"
	"github.com/bensabler/ff6-decompile/internal/graphics/tile2bpp"
	"github.com/bensabler/ff6-decompile/internal/rom"
)

// TestHUDFontMatchesROMBlock proves the archive path and the ROM path agree,
// tile for tile.
//
// This is the test whose absence let the HUD font extractor read the wrong
// ROM range for three days (DEMO-0001 deviation D1): the extractor, the
// manifest and the archive were all internally consistent with each other and
// all wrong about the ROM. Only a comparison against the ROM itself catches
// that class of error.
//
// It is opt-in via FF6_ROM and never runs in CI, because CI has no ROM.
func TestHUDFontMatchesROMBlock(t *testing.T) {
	romPath := os.Getenv(rom.EnvVar)
	if romPath == "" {
		t.Skipf("set %s to run the archive-vs-ROM differential", rom.EnvVar)
	}
	r, err := rom.Load(romPath)
	if err != nil {
		t.Fatalf("loading ROM: %v", err)
	}

	root := repoRoot(t)
	a, err := content.Open(root)
	if err != nil {
		t.Skipf("no generated archive: %v", err)
	}
	f, err := content.LoadHUDFont(a)
	if err != nil {
		t.Fatalf("loading the font from the archive: %v", err)
	}

	if a.ROMRevision() != r.SHA256() {
		t.Fatalf("archive was generated from ROM %s, supplied ROM is %s", a.ROMRevision(), r.SHA256())
	}

	var compared, blank int
	for i := 0; i < 256; i++ {
		b := byte(i)
		off := hudfont.ROMFileOffset(b)
		src, err := r.Slice(off, hudfont.TileBytes)
		if err != nil {
			t.Fatalf("glyph $%02X at ROMFILE:0x%06X: %v", b, off, err)
		}
		want, err := tile2bpp.Decode(src)
		if err != nil {
			t.Fatalf("glyph $%02X: %v", b, err)
		}
		got, ink := f.Glyph(b)
		if *got != want {
			t.Fatalf("glyph $%02X (VRAM $%03X, sheet index %d, ROMFILE:0x%06X) differs between the archive and the ROM",
				b, hudfont.VRAMTile(b), hudfont.SheetIndex(b), off)
		}
		compared++
		if !ink {
			blank++
		}
	}

	if compared != 256 {
		t.Fatalf("compared %d glyphs, want 256", compared)
	}
	// EXP-0049 classified 140 of the 256 as blank. If that count moves, the
	// block being read has changed and every glyph claim needs re-checking.
	if blank != 140 {
		t.Errorf("%d blank glyphs, EXP-0049 recorded 140; the block being read may have moved", blank)
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
