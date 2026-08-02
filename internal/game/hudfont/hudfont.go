// Package hudfont models the relationship between the fixed-width text
// encoding and the battle HUD font tile block (EXP-0049).
//
// The relation is affine and Confirmed:
//
//	VRAM tile = $100 + encodedByte
//	ROMFILE   = 0x047FC0 + 16*encodedByte
//	sheet idx = encodedByte + 1   (the sheet starts at VRAM $FF)
//
// It was established by decoding the EXP-0023 battle HUD tilemap through the
// candidate relation: all 37 referenced tiles resolve to coherent game text
// ("Were-Rat", "Repo Man", "WEDGE", "VICKS", "MagiTek", "Item", "Row", "Tek",
// HP digits), which no incorrect mapping produces.
//
// This package carries the *mapping*, not the pixels. Tile bytes are
// ROM-derived and stay in the local archive; the tracked glyph table records
// each tile's SHA-256 so a decoded tile can be verified without the repository
// ever containing it (docs/legal/ASSET_POLICY.md).
package hudfont

import (
	"fmt"

	"github.com/bensabler/ff6-decompile/internal/game/textenc"
)

// Block geometry, from GFX-0001 / ROM-0016.
const (
	// BlockROMFile is the first byte of the 257-tile block, VRAM tile $FF.
	BlockROMFile = 0x047FB0
	// BlockTiles is the tile count, VRAM $FF-$1FF inclusive.
	BlockTiles = 257
	// FirstVRAMTile is the VRAM tile index of sheet index 0.
	FirstVRAMTile = 0xFF
	// GlyphBase is the VRAM tile index that encoded byte $00 maps to.
	GlyphBase = 0x100
	// TileBytes is the 2bpp encoded tile size.
	TileBytes = 16
)

// VRAMTile returns the VRAM tile index an encoded byte selects.
func VRAMTile(b byte) int { return GlyphBase + int(b) }

// SheetIndex returns the index of an encoded byte's tile within the extracted
// font sheet, whose element 0 is VRAM tile $FF.
func SheetIndex(b byte) int { return VRAMTile(b) - FirstVRAMTile }

// ROMFileOffset returns the ROMFILE offset of an encoded byte's tile.
func ROMFileOffset(b byte) int { return BlockROMFile + SheetIndex(b)*TileBytes }

// Classification records how much is known about one byte's tile.
type Classification string

const (
	// Character: the tile renders a glyph textenc names, corroborated in
	// rendered output. Confirmed.
	Character Classification = "character"
	// Blank: the tile is entirely palette index 0. Confirmed by decode.
	Blank Classification = "blank"
	// Structural: a non-character tile whose role is identified from
	// rendered context — the HP/ATB gauge pieces. Confirmed as a role,
	// Unknown as to exact semantics.
	Structural Classification = "structural"
	// Unidentified: a non-blank tile carrying some glyph this project has
	// not seen rendered in a known context. Deliberately not guessed.
	Unidentified Classification = "unidentified"
)

// Glyph is one row of the tracked glyph table.
type Glyph struct {
	EncodedByte    byte           `json:"encoded_byte"`
	VRAMTile       int            `json:"vram_tile"`
	SheetIndex     int            `json:"sheet_index"`
	ROMFileOffset  int            `json:"rom_file_offset"`
	Character      *string        `json:"character"`
	Classification Classification `json:"classification"`
	TileSHA256     string         `json:"tile_sha256"`
	Confidence     string         `json:"confidence"`
	Note           string         `json:"note,omitempty"`
}

// structuralNotes records the gauge tiles identified in EXP-0049 from their
// position in the HUD tilemap: they render the HP/ATB bar, not text.
var structuralNotes = map[byte]string{
	0xF9: "HP/ATB gauge left cap",
	0xFA: "HP/ATB gauge right cap",
	0xF0: "HP/ATB gauge segment, empty",
	0xF1: "HP/ATB gauge segment, partially filled",
	0xF8: "HP/ATB gauge segment, full",
}

// Classify reports what is known about an encoded byte, given whether its
// tile decoded to all-zero pixels. It never guesses a character.
func Classify(b byte, tileIsBlank bool) (Classification, *string, string) {
	if r, ok := textenc.Decode(b); ok && r != ' ' {
		s := string(r)
		return Character, &s, "confirmed"
	}
	if tileIsBlank {
		return Blank, nil, "confirmed"
	}
	if _, ok := structuralNotes[b]; ok {
		return Structural, nil, "confirmed"
	}
	return Unidentified, nil, "unknown"
}

// Note returns the recorded note for a byte, if any.
func Note(b byte) string { return structuralNotes[b] }

// Validate checks a decoded glyph table for the invariants the relation
// implies, so a regenerated table cannot drift from the model.
func Validate(glyphs []Glyph) error {
	if len(glyphs) != 256 {
		return fmt.Errorf("glyph table has %d rows, want 256 (one per encoded byte)", len(glyphs))
	}
	for i, g := range glyphs {
		if int(g.EncodedByte) != i {
			return fmt.Errorf("row %d describes byte $%02X; rows must be byte-ordered", i, g.EncodedByte)
		}
		if g.VRAMTile != VRAMTile(g.EncodedByte) {
			return fmt.Errorf("byte $%02X: vram_tile $%03X, relation gives $%03X", g.EncodedByte, g.VRAMTile, VRAMTile(g.EncodedByte))
		}
		if g.SheetIndex != SheetIndex(g.EncodedByte) {
			return fmt.Errorf("byte $%02X: sheet_index %d, relation gives %d", g.EncodedByte, g.SheetIndex, SheetIndex(g.EncodedByte))
		}
		if g.ROMFileOffset != ROMFileOffset(g.EncodedByte) {
			return fmt.Errorf("byte $%02X: rom_file_offset 0x%06X, relation gives 0x%06X", g.EncodedByte, g.ROMFileOffset, ROMFileOffset(g.EncodedByte))
		}
		if len(g.TileSHA256) != 64 {
			return fmt.Errorf("byte $%02X: tile_sha256 %q is not a full digest", g.EncodedByte, g.TileSHA256)
		}
		if g.Classification == Character && g.Character == nil {
			return fmt.Errorf("byte $%02X: classified character with no character", g.EncodedByte)
		}
		if g.Classification != Character && g.Character != nil {
			return fmt.Errorf("byte $%02X: classified %s but names character %q", g.EncodedByte, g.Classification, *g.Character)
		}
	}
	return nil
}
