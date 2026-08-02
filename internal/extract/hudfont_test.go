package extract

import (
	"testing"

	"github.com/bensabler/ff6-decompile/internal/census"
	"github.com/bensabler/ff6-decompile/internal/graphics/tile2bpp"
)

// hudFontLedgerRegion is the ROM-ownership ledger entry that owns the battle
// HUD font block. The extractor and the ledger are two independent records of
// the same fact, and until 2026-08-02 they disagreed: the extractor read from
// the back-projected anchor 0x046FC0 while the ledger recorded the real block
// at 0x047FB0, emitting 255 tiles of attack-table bytes (DEMO-0001 deviation
// D1). Nothing caught it, because nothing compared them.
const hudFontLedgerRegion = "ROM-0016"

// TestHUDFontMatchesROMLedger is the regression test for D1. It needs no ROM:
// it asserts the extractor's source span is exactly the span the ledger claims
// the font occupies.
func TestHUDFontMatchesROMLedger(t *testing.T) {
	_, regions, err := census.Load(repoRoot(t))
	if err != nil {
		t.Fatalf("loading ROM regions: %v", err)
	}

	var region *census.Region
	for i := range regions.Regions {
		if regions.Regions[i].ID == hudFontLedgerRegion {
			region = &regions.Regions[i]
			break
		}
	}
	if region == nil {
		t.Fatalf("ledger region %s not found in manifests/rom-regions.json", hudFontLedgerRegion)
	}

	wantSize := hudFontTiles * tile2bpp.EncodedSize
	if region.Size != wantSize {
		t.Errorf("ledger %s size %d, extractor reads %d tiles x %d = %d bytes",
			hudFontLedgerRegion, region.Size, hudFontTiles, tile2bpp.EncodedSize, wantSize)
	}
	if hudFontBase != region.Start {
		t.Errorf("extractor hudFontBase = ROMFILE:0x%06X, ledger %s starts at ROMFILE:0x%06X",
			hudFontBase, hudFontLedgerRegion, region.Start)
	}
}

// TestHUDFontAnchorRelation pins the affine relation GFX-0001 states, so the
// block start stays derivable rather than becoming a bare magic number. The
// anchor back-projects to VRAM tile $000 and deliberately falls inside the
// attack-data region; only tiles $FF-$1FF are real font data.
func TestHUDFontAnchorRelation(t *testing.T) {
	if got := hudFontAnchor + hudFontFirstVRAMTile*tile2bpp.EncodedSize; got != hudFontBase {
		t.Errorf("anchor relation gives 0x%06X, hudFontBase is 0x%06X", got, hudFontBase)
	}
	if hudFontAnchor >= hudFontBase {
		t.Errorf("anchor 0x%06X should precede the block start 0x%06X", hudFontAnchor, hudFontBase)
	}
	// The anchor lands inside the attack/spell data region. If a future edit
	// ever makes the anchor a plausible block start, that coincidence should
	// fail loudly rather than silently look reasonable.
	if hudFontAnchor <= spellDataBase {
		t.Errorf("anchor 0x%06X no longer sits above the attack-data base 0x%06X; re-verify GFX-0001",
			hudFontAnchor, spellDataBase)
	}
	// VRAM tiles $FF-$1FF inclusive is 257 tiles.
	if want := 0x1FF - hudFontFirstVRAMTile + 1; hudFontTiles != want {
		t.Errorf("hudFontTiles = %d, VRAM $%03X-$1FF spans %d tiles",
			hudFontTiles, hudFontFirstVRAMTile, want)
	}
}
