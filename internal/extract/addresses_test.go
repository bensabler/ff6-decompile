package extract

import (
	"testing"

	"github.com/bensabler/ff6-decompile/internal/platform/snesaddr"
)

// TestExtractorOffsetsMatchTheirCPUAddresses asserts every hand-computed
// ROMFILE constant in extractors.go against the CPU address its experiment
// record cites. Each pair is two independent records of one fact, and the HUD
// font block is the standing proof that two such records can silently diverge
// (DEMO-0001 deviation D1). Needs no ROM.
func TestExtractorOffsetsMatchTheirCPUAddresses(t *testing.T) {
	tests := []struct {
		name string
		cpu  uint32
		file int
	}{
		{"spell names (EXP-0026/0027)", 0xE6F567, spellNameBase},
		{"spell/attack data (EXP-0019/0027)", 0xC46AC0, spellDataBase},
		{"esper names (EXP-0027)", 0xE6F6E1, esperNameBase},
		{"HUD font anchor (GFX-0001)", 0xC46FC0, hudFontAnchor},
		{"HUD font block start (ROM-0016)", 0xC47FB0, hudFontBase},
		{"SFX pack (EXP-0024/AUD-0001)", 0xC51EC9, sfxPackBase},
		{"SFX sample 5 (EXP-0024)", 0xC51FA1, sfxSample5Off},
		{"monster records (EXP-0028)", 0xCF0000, monsterBase},
		{"formation table (EXP-0030)", 0xCF6200, formationBase},
		{"formation flags (EXP-0030)", 0xCF5900, formationFlagBase},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			off, win, err := snesaddr.ROMFile(tt.cpu)
			if err != nil {
				t.Fatalf("ROMFile(ROMCPU:$%06X): %v", tt.cpu, err)
			}
			if off != tt.file {
				t.Errorf("ROMCPU:$%06X maps to ROMFILE:0x%06X, constant is ROMFILE:0x%06X",
					tt.cpu, off, tt.file)
			}
			if !win.Confirmed() {
				t.Errorf("resolves through %v, which carries no FF6 runtime evidence", win)
			}
		})
	}
}

// TestExtractorSpansStayInsideTheImage catches an off-by-one or a bad count
// before it reaches a ROM read, where it would surface as an opaque slice
// error during extraction.
func TestExtractorSpansStayInsideTheImage(t *testing.T) {
	tests := []struct {
		name  string
		start int
		size  int
	}{
		{"spell names", spellNameBase, spellNameCount * spellNameStride},
		{"spell data", spellDataBase, spellNameCount * spellDataStride},
		{"esper names", esperNameBase, esperNameCount * esperNameStride},
		{"HUD font", hudFontBase, hudFontTiles * 16},
		{"SFX pack", sfxPackBase, sfxPackLen},
		{"monsters", monsterBase, monsterStride * monsterSliceCount},
		{"formations", formationBase, formationStride * formationCount},
		{"formation flags", formationFlagBase, formationFlagLen},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.size <= 0 {
				t.Fatalf("span size %d", tt.size)
			}
			last := tt.start + tt.size - 1
			if _, err := snesaddr.ROMCPU(tt.start); err != nil {
				t.Errorf("start ROMFILE:0x%06X: %v", tt.start, err)
			}
			if _, err := snesaddr.ROMCPU(last); err != nil {
				t.Errorf("last byte ROMFILE:0x%06X: %v", last, err)
			}
		})
	}
}
