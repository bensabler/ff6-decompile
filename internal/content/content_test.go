package content

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"image"
	"image/color"
	"image/png"
	"os"
	"path/filepath"
	"testing"

	"github.com/bensabler/ff6-decompile/internal/assetmanifest"
	"github.com/bensabler/ff6-decompile/internal/game/hudfont"
	"github.com/bensabler/ff6-decompile/internal/graphics/framebuf"
)

// syntheticSheet builds a font sheet with the same geometry as the real one
// but project-authored pixels: glyph for byte b is a solid block of colour
// (b%3)+1, except b==0x20 which is left blank. No ROM bytes are involved, so
// everything derived from it is safe to assert on and safe to commit.
func syntheticSheet(t *testing.T) []byte {
	t.Helper()
	const cols = 16
	rows := (hudfont.BlockTiles + cols - 1) / cols
	pal := color.Palette{
		color.RGBA{0, 0, 0, 0},
		color.RGBA{0, 0, 0, 255},
		color.RGBA{165, 165, 165, 255},
		color.RGBA{255, 255, 255, 255},
	}
	img := image.NewPaletted(image.Rect(0, 0, cols*8, rows*8), pal)
	for i := 0; i < 256; i++ {
		b := byte(i)
		if b == 0x20 {
			continue // deliberately blank
		}
		sheet := hudfont.SheetIndex(b)
		ox, oy := (sheet%cols)*8, (sheet/cols)*8
		v := uint8(i%3) + 1
		for y := 0; y < 8; y++ {
			for x := 0; x < 8; x++ {
				img.SetColorIndex(ox+x, oy+y, v)
			}
		}
	}
	var buf bytes.Buffer
	if err := (&png.Encoder{CompressionLevel: png.DefaultCompression}).Encode(&buf, img); err != nil {
		t.Fatal(err)
	}
	return buf.Bytes()
}

// fixtureArchive writes a self-contained archive into a temp dir. It never
// reads the real archive, so these tests run in CI with no ROM present.
func fixtureArchive(t *testing.T, sheet []byte) string {
	t.Helper()
	root := t.TempDir()
	rel := filepath.Join(assetmanifest.ArchiveRoot, "graphics", "hud-font-sheet.png")
	full := filepath.Join(root, rel)
	if err := os.MkdirAll(filepath.Dir(full), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(full, sheet, 0o644); err != nil {
		t.Fatal(err)
	}
	sum := sha256.Sum256(sheet)
	man := assetmanifest.Manifest{
		SchemaVersion: "1.0",
		ROMRevision:   "0f51b4fca41b7fd509e4b8f9d543151f68efa5e97b08493e4b2a0c06f5d8d5e2",
		ArchiveRoot:   assetmanifest.ArchiveRoot,
		Assets: []assetmanifest.Asset{{
			ID: HUDFontAssetID, Name: "synthetic fixture font", Category: "graphics",
			ROMRevision: "0f51b4fca41b7fd509e4b8f9d543151f68efa5e97b08493e4b2a0c06f5d8d5e2",
			ROMSource:   "ROMFILE:0x047FB0-0x048FBF",
			ExtractorID: "hud-font", ExtractorVer: "2.0.0",
			OutputPath: filepath.ToSlash(rel), OutputFormat: "image/png",
			SHA256: hex.EncodeToString(sum[:]), Verification: "fixture",
		}},
	}
	raw, err := json.MarshalIndent(man, "", "  ")
	if err != nil {
		t.Fatal(err)
	}
	mp := filepath.Join(root, assetmanifest.Path)
	if err := os.MkdirAll(filepath.Dir(mp), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(mp, raw, 0o644); err != nil {
		t.Fatal(err)
	}
	return root
}

func TestOpenAndRead(t *testing.T) {
	sheet := syntheticSheet(t)
	root := fixtureArchive(t, sheet)

	a, err := Open(root)
	if err != nil {
		t.Fatalf("Open: %v", err)
	}
	if a.Root() != root {
		t.Errorf("Root() = %q, want %q", a.Root(), root)
	}
	got, err := a.Read(HUDFontAssetID)
	if err != nil {
		t.Fatalf("Read: %v", err)
	}
	if !bytes.Equal(got, sheet) {
		t.Error("Read returned different bytes than were written")
	}
}

func TestOpenReportsMissingArchive(t *testing.T) {
	_, err := Open(t.TempDir())
	if !errors.Is(err, ErrNoArchive) {
		t.Fatalf("err = %v, want ErrNoArchive", err)
	}
	// The message must tell the operator what to run; "no archive" alone is
	// not actionable.
	if !bytes.Contains([]byte(err.Error()), []byte("ff6lab extract all")) {
		t.Errorf("error should name the setup command:\n%v", err)
	}
}

func TestReadReportsMissingAsset(t *testing.T) {
	root := fixtureArchive(t, syntheticSheet(t))
	a, err := Open(root)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := a.Read("ASSET-DOES-NOT-EXIST"); !errors.Is(err, ErrAssetMissing) {
		t.Fatalf("err = %v, want ErrAssetMissing", err)
	}
	// A manifest entry whose file was deleted is also "missing", not "stale".
	if err := os.Remove(filepath.Join(root, assetmanifest.ArchiveRoot, "graphics", "hud-font-sheet.png")); err != nil {
		t.Fatal(err)
	}
	if _, err := a.Read(HUDFontAssetID); !errors.Is(err, ErrAssetMissing) {
		t.Fatalf("deleted file: err = %v, want ErrAssetMissing", err)
	}
}

func TestReadReportsStaleAsset(t *testing.T) {
	root := fixtureArchive(t, syntheticSheet(t))
	a, err := Open(root)
	if err != nil {
		t.Fatal(err)
	}
	p := filepath.Join(root, assetmanifest.ArchiveRoot, "graphics", "hud-font-sheet.png")
	if err := os.WriteFile(p, []byte("not a png"), 0o644); err != nil {
		t.Fatal(err)
	}
	if _, err := a.Read(HUDFontAssetID); !errors.Is(err, ErrAssetStale) {
		t.Fatalf("err = %v, want ErrAssetStale", err)
	}
}

func TestLoadHUDFont(t *testing.T) {
	root := fixtureArchive(t, syntheticSheet(t))
	a, err := Open(root)
	if err != nil {
		t.Fatal(err)
	}
	f, err := LoadHUDFont(a)
	if err != nil {
		t.Fatalf("LoadHUDFont: %v", err)
	}

	// Byte b's tile must be the fixture's block for b, which proves the
	// sheet-index relation was applied rather than a naive b-th cell.
	for _, b := range []byte{0x00, 0x80, 0xBF, 0xC4, 0xFF} {
		tile, ink := f.Glyph(b)
		if !ink {
			t.Errorf("glyph $%02X reports no ink, but the fixture filled it", b)
			continue
		}
		want := uint8(int(b)%3) + 1
		if tile[0][0] != want {
			t.Errorf("glyph $%02X = %d, want %d (sheet index %d)", b, tile[0][0], want, hudfont.SheetIndex(b))
		}
	}
	// The deliberately blank cell must report no ink.
	if _, ink := f.Glyph(0x20); ink {
		t.Error("glyph $20 was left blank in the fixture but reports ink")
	}
}

func TestLoadHUDFontRejectsNonPaletted(t *testing.T) {
	img := image.NewRGBA(image.Rect(0, 0, 128, 136))
	var buf bytes.Buffer
	if err := png.Encode(&buf, img); err != nil {
		t.Fatal(err)
	}
	root := fixtureArchive(t, buf.Bytes())
	a, err := Open(root)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := LoadHUDFont(a); err == nil {
		t.Fatal("an RGBA sheet should be rejected, not silently mis-decoded")
	}
}

func TestDrawString(t *testing.T) {
	root := fixtureArchive(t, syntheticSheet(t))
	a, _ := Open(root)
	f, err := LoadHUDFont(a)
	if err != nil {
		t.Fatal(err)
	}

	fb := framebuf.New()
	end := f.DrawString(fb, 0, 0, "AB", TextOptions{})
	if end != 16 {
		t.Errorf("DrawString returned x=%d, want 16", end)
	}
	if fb.At(0, 0) == 0 || fb.At(8, 0) == 0 {
		t.Error("two glyphs should have been drawn")
	}

	// An unverified rune must not draw a ROM tile. Without ShowMissing it
	// draws nothing at all; the cursor still advances.
	fb2 := framebuf.New()
	end = f.DrawString(fb2, 0, 0, "!", TextOptions{})
	if end != 8 {
		t.Errorf("an unmapped rune should still advance the cursor; x=%d", end)
	}
	for _, v := range fb2.Pix {
		if v != 0 {
			t.Fatal("an unmapped rune drew something without ShowMissing set")
		}
	}

	// With ShowMissing it draws the project-authored box.
	fb3 := framebuf.New()
	f.DrawString(fb3, 0, 0, "!", TextOptions{ShowMissing: true})
	if fb3.At(1, 1) == 0 {
		t.Error("ShowMissing should draw a box outline")
	}
	if fb3.At(3, 3) != 0 {
		t.Error("the missing-glyph box should be hollow, so it cannot pass for a glyph")
	}
}

func TestDrawStringSubPalette(t *testing.T) {
	root := fixtureArchive(t, syntheticSheet(t))
	a, _ := Open(root)
	f, _ := LoadHUDFont(a)

	// Slot 4 is entries $10-$13, so the ink lands at $10 + its level. The
	// arithmetic is what makes SubPalette a slot rather than an offset.
	const slot = SubPalette(4)
	if base := slot.Base(); base != 0x10 {
		t.Fatalf("SubPalette(4).Base() = $%02X, want $10", base)
	}

	fb := framebuf.New()
	f.DrawString(fb, 0, 0, "A", TextOptions{Palette: slot})
	want := uint8(int('A'-'A'+0x80)%3) + 1 + 0x10
	if got := fb.At(0, 0); got != want {
		t.Errorf("At(0,0) = $%02X, want $%02X", got, want)
	}
}

func TestSubPaletteBaseIsSlotTimesFour(t *testing.T) {
	tests := []struct {
		slot SubPalette
		want uint8
	}{
		{SubPalettePrimary, 0},
		{SubPaletteDim, 4},
		{2, 8},
		{3, 12},
		{63, 252},
	}
	for _, tt := range tests {
		if got := tt.slot.Base(); got != tt.want {
			t.Errorf("SubPalette(%d).Base() = %d, want %d", tt.slot, got, tt.want)
		}
	}
	// The two values the scenes used to pass as if they were brightness
	// levels. As slots they address entries far outside anything
	// framebuf.GrayPalette defines, which is what makes the mistake visible.
	if SubPalette(3).Base() != 12 || SubPalette(2).Base() != 8 {
		t.Error("the pre-fix values no longer land where the regression test expects")
	}
}

func TestDrawBytesSkipsPadding(t *testing.T) {
	root := fixtureArchive(t, syntheticSheet(t))
	a, _ := Open(root)
	f, _ := LoadHUDFont(a)

	fb := framebuf.New()
	// $FF is name-field padding. In the real font its tile is blank, and the
	// fixture fills it — so assert on the advance, which is the contract.
	if end := f.DrawBytes(fb, 0, 0, []byte{0x80, 0x81, 0xFF}, TextOptions{}); end != 24 {
		t.Errorf("DrawBytes returned x=%d, want 24", end)
	}
}

func TestMeasureString(t *testing.T) {
	if got := MeasureString("WEDGE"); got != 40 {
		t.Errorf("MeasureString(\"WEDGE\") = %d, want 40", got)
	}
	if got := MeasureString(""); got != 0 {
		t.Errorf("MeasureString(\"\") = %d, want 0", got)
	}
}
