package extract

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"path"

	"github.com/bensabler/ff6-decompile/internal/audio/brr"
	"github.com/bensabler/ff6-decompile/internal/game/textenc"
	"github.com/bensabler/ff6-decompile/internal/graphics/tile2bpp"
	"github.com/bensabler/ff6-decompile/internal/rom"
)

// Confirmed ROM offsets (ROMFILE), each traced to its experiment record.
const (
	// spellNameBase: 7-byte [class icon][6-char name] records, 54 entries
	// (EXP-0026 located, EXP-0027 boundary Confirmed).
	spellNameBase   = 0x26F567
	spellNameStride = 7
	spellNameCount  = 54

	// spellDataBase: the global spell database, 14-byte records; byte 5 is
	// MP cost (Confirmed for id 45 by live deduction, EXP-0027).
	spellDataBase   = 0x046AC0
	spellDataStride = 14

	// esperNameBase: 8-byte names, 27 entries decoding cleanly (EXP-0027).
	esperNameBase   = 0x26F6E1
	esperNameStride = 8
	esperNameCount  = 27

	// hudFontBase: battle HUD 2bpp tile block, raw copy of VRAM tiles
	// $FF-$1FF (EXP-0023 / GFX-0001, ledger region ROM-0016).
	//
	// GFX-0001 states the provenance as the affine relation
	// "ROMFILE:0x046FC0 + 16*tileindex", where tileindex is the *VRAM* tile
	// number. 0x046FC0 is therefore the back-projected anchor for VRAM tile
	// $000 and is NOT the start of the block: it lands inside the attack-data
	// region (spellDataBase = 0x046AC0), with which the font shares no bytes.
	//
	// The block starts at the first tile that actually exists in the copy,
	// VRAM $FF: 0x046FC0 + 16*0xFF = 0x047FB0. Reading from the anchor instead
	// emits 255 tiles of attack-table bytes (DEMO-0001 deviation D1, fixed
	// 2026-08-02). hudFontFirstVRAMTile keeps the relation checkable.
	hudFontAnchor        = 0x046FC0
	hudFontFirstVRAMTile = 0xFF
	hudFontBase          = hudFontAnchor + hudFontFirstVRAMTile*tile2bpp.EncodedSize // 0x047FB0
	hudFontTiles         = 257

	// sfxPackBase: 288 bytes byte-identical to ARAM $4800-$491F, samples
	// 0-7; sample 5 (the battle confirm cue) starts at 0x051FA1
	// (EXP-0024 / AUD-0001).
	sfxPackBase   = 0x051EC9
	sfxPackLen    = 288
	sfxSample5Off = 0x051FA1

	// monster/formation tables (EXP-0028/0029/0030). Extracted as raw
	// evidence slices: the record *layouts* are only partly decoded, so
	// these are archived for analysis rather than rendered.
	monsterBase       = 0x0F0000
	monsterStride     = 32
	monsterSliceCount = 384
	formationBase     = 0x0F6200
	formationStride   = 15
	formationCount    = 576
	formationFlagBase = 0x0F5900
	formationFlagLen  = 576 * 4
)

func romSource(off, length int) string {
	return fmt.Sprintf("ROMFILE:0x%06X-0x%06X", off, off+length-1)
}

func jsonBytes(v any) ([]byte, error) {
	var buf bytes.Buffer
	enc := json.NewEncoder(&buf)
	enc.SetIndent("", "  ")
	enc.SetEscapeHTML(false)
	if err := enc.Encode(v); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// ---------------------------------------------------------------------
// data: spell names + records
// ---------------------------------------------------------------------

type spellExtractor struct{}

func (spellExtractor) ID() string       { return "spells" }
func (spellExtractor) Category() string { return "data" }
func (spellExtractor) Version() string  { return "1.0.0" }

// SpellRecord is the archived per-spell row: the name plus the raw record
// bytes. Names and numeric fields are policy-permitted in Git (the census
// inventory carries them); the archive keeps the full raw bytes too.
type SpellRecord struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	ClassIcon uint8  `json:"class_icon"`
	MPCost    uint8  `json:"mp_cost"`
	NameBytes string `json:"name_bytes"`
	DataBytes string `json:"data_bytes"`
}

// Spells decodes all 54 name+data records. Exported so the census
// migration and tests can reuse exactly what the extractor writes.
func Spells(r *rom.ROM) ([]SpellRecord, error) {
	out := make([]SpellRecord, 0, spellNameCount)
	for i := 0; i < spellNameCount; i++ {
		nb, err := r.Slice(spellNameBase+i*spellNameStride, spellNameStride)
		if err != nil {
			return nil, err
		}
		db, err := r.Slice(spellDataBase+i*spellDataStride, spellDataStride)
		if err != nil {
			return nil, err
		}
		name := textenc.DecodeFixed(nb[1:])
		if !textenc.Clean(nb[1:]) {
			return nil, fmt.Errorf("spell %d: name bytes contain unmapped glyphs (%s); refusing to guess", i, hex.EncodeToString(nb[1:]))
		}
		out = append(out, SpellRecord{
			ID:        i,
			Name:      name,
			ClassIcon: nb[0],
			MPCost:    db[5],
			NameBytes: hex.EncodeToString(nb),
			DataBytes: hex.EncodeToString(db),
		})
	}
	return out, nil
}

func (e spellExtractor) Extract(r *rom.ROM) ([]output, error) {
	recs, err := Spells(r)
	if err != nil {
		return nil, err
	}
	data, err := jsonBytes(map[string]any{
		"extractor":    e.ID(),
		"version":      e.Version(),
		"rom_revision": r.SHA256(),
		"name_table":   romSource(spellNameBase, spellNameCount*spellNameStride),
		"data_table":   romSource(spellDataBase, spellNameCount*spellDataStride),
		"record_count": len(recs),
		"records":      recs,
	})
	if err != nil {
		return nil, err
	}
	p := path.Join(ArchiveRoot, "raw", "spells.json")
	return []output{{
		asset: Asset{
			ID: "ASSET-DATA-0001", Name: "Spell name and record table (54 entries)",
			Category: "data", ROMRevision: r.SHA256(),
			ROMSource:   romSource(spellNameBase, spellNameCount*spellNameStride) + " + " + romSource(spellDataBase, spellNameCount*spellDataStride),
			ExtractorID: e.ID(), ExtractorVer: e.Version(),
			OutputPath: p, OutputFormat: "application/json",
			Properties: fmt.Sprintf("%d records; 7-byte name stride, 14-byte data stride", len(recs)),
			SHA256:     hashBytes(data), RuntimeConsumer: "magic menu list renderer; battle spell dispatch",
			CensusRefs:     []string{"CEN-MAGIC-0002", "CEN-MAGIC-0003"},
			ExperimentRefs: []string{"EXP-0026", "EXP-0027"},
			Verification:   "hash-verified against tracked manifest",
		},
		data: data,
	}}, nil
}

// ---------------------------------------------------------------------
// data: esper names
// ---------------------------------------------------------------------

type esperExtractor struct{}

func (esperExtractor) ID() string       { return "espers" }
func (esperExtractor) Category() string { return "data" }
func (esperExtractor) Version() string  { return "1.0.0" }

// EsperName is one archived esper name row.
type EsperName struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	NameBytes string `json:"name_bytes"`
}

// Espers decodes the 27-entry esper name table (EXP-0027).
func Espers(r *rom.ROM) ([]EsperName, error) {
	out := make([]EsperName, 0, esperNameCount)
	for i := 0; i < esperNameCount; i++ {
		nb, err := r.Slice(esperNameBase+i*esperNameStride, esperNameStride)
		if err != nil {
			return nil, err
		}
		if !textenc.Clean(nb) {
			return nil, fmt.Errorf("esper %d: unmapped glyphs (%s)", i, hex.EncodeToString(nb))
		}
		out = append(out, EsperName{ID: i, Name: textenc.DecodeFixed(nb), NameBytes: hex.EncodeToString(nb)})
	}
	return out, nil
}

func (e esperExtractor) Extract(r *rom.ROM) ([]output, error) {
	recs, err := Espers(r)
	if err != nil {
		return nil, err
	}
	data, err := jsonBytes(map[string]any{
		"extractor": e.ID(), "version": e.Version(), "rom_revision": r.SHA256(),
		"name_table":   romSource(esperNameBase, esperNameCount*esperNameStride),
		"record_count": len(recs), "records": recs,
	})
	if err != nil {
		return nil, err
	}
	p := path.Join(ArchiveRoot, "raw", "espers.json")
	return []output{{
		asset: Asset{
			ID: "ASSET-DATA-0002", Name: "Esper name table (27 entries)",
			Category: "data", ROMRevision: r.SHA256(),
			ROMSource:   romSource(esperNameBase, esperNameCount*esperNameStride),
			ExtractorID: e.ID(), ExtractorVer: e.Version(),
			OutputPath: p, OutputFormat: "application/json",
			Properties: fmt.Sprintf("%d records; 8-byte stride", len(recs)),
			SHA256:     hashBytes(data), RuntimeConsumer: "esper menu list renderer",
			CensusRefs: []string{"CEN-MAGIC-0002"}, ExperimentRefs: []string{"EXP-0027"},
			Verification: "hash-verified against tracked manifest",
		},
		data: data,
	}}, nil
}

// ---------------------------------------------------------------------
// graphics: battle HUD font sheet
// ---------------------------------------------------------------------

type hudFontExtractor struct{}

func (hudFontExtractor) ID() string       { return "hud-font" }
func (hudFontExtractor) Category() string { return "graphics" }

// Version 2.0.0: corrected the source range from the back-projected anchor
// 0x046FC0 to the true block start 0x047FB0 (DEMO-0001 deviation D1).
func (hudFontExtractor) Version() string { return "2.0.0" }

// hudPalette is the live BG-2bpp palette 0 recorded in EXP-0023
// (BGR555 $0000, $0000, $5294, $7FFF). Index 0 is rendered transparent.
var hudPalette = color.Palette{
	color.RGBA{0, 0, 0, 0},
	color.RGBA{0, 0, 0, 255},
	color.RGBA{165, 165, 165, 255},
	color.RGBA{255, 255, 255, 255},
}

func (e hudFontExtractor) Extract(r *rom.ROM) ([]output, error) {
	const cols = 16
	rows := (hudFontTiles + cols - 1) / cols
	img := image.NewPaletted(image.Rect(0, 0, cols*8, rows*8), hudPalette)
	for i := 0; i < hudFontTiles; i++ {
		src, err := r.Slice(hudFontBase+i*tile2bpp.EncodedSize, tile2bpp.EncodedSize)
		if err != nil {
			return nil, err
		}
		t, err := tile2bpp.Decode(src)
		if err != nil {
			return nil, fmt.Errorf("hud tile %d: %w", i, err)
		}
		ox, oy := (i%cols)*8, (i/cols)*8
		for y := 0; y < 8; y++ {
			for x := 0; x < 8; x++ {
				img.SetColorIndex(ox+x, oy+y, t[y][x])
			}
		}
	}
	// Deterministic PNG: default encoder settings, no timestamps.
	var buf bytes.Buffer
	if err := (&png.Encoder{CompressionLevel: png.DefaultCompression}).Encode(&buf, img); err != nil {
		return nil, err
	}
	data := buf.Bytes()
	p := path.Join(ArchiveRoot, "graphics", "hud-font-sheet.png")
	return []output{{
		asset: Asset{
			ID: "ASSET-GFX-0001", Name: "Battle HUD fixed 8x8 font sheet",
			Category: "graphics", ROMRevision: r.SHA256(),
			ROMSource:   romSource(hudFontBase, hudFontTiles*tile2bpp.EncodedSize),
			ExtractorID: e.ID(), ExtractorVer: e.Version(),
			OutputPath: p, OutputFormat: "image/png",
			Properties: fmt.Sprintf("%dx%d px; %d tiles of 8x8; 2bpp planar; palette BGR555 0000/0000/5294/7FFF (index 0 transparent); sheet index i = VRAM tile $%03X+i, row-major in %d columns",
				cols*8, rows*8, hudFontTiles, hudFontFirstVRAMTile, cols),
			SHA256: hashBytes(data), RuntimeConsumer: "battle HUD BG3 text layer (VRAM tiles $FF-$1FF); sheet index 0 = VRAM tile $FF",
			CensusRefs: []string{"CEN-GFX-0001"}, ExperimentRefs: []string{"EXP-0023"},
			Verification: "hash-verified against tracked manifest; ROM block proven byte-identical to live VRAM in EXP-0023",
		},
		data: data,
	}}, nil
}

// ---------------------------------------------------------------------
// audio: SFX sample pack (BRR + decoded PCM)
// ---------------------------------------------------------------------

type sfxPackExtractor struct{}

func (sfxPackExtractor) ID() string       { return "sfx-pack" }
func (sfxPackExtractor) Category() string { return "audio" }
func (sfxPackExtractor) Version() string  { return "1.0.0" }

func (e sfxPackExtractor) Extract(r *rom.ROM) ([]output, error) {
	pack, err := r.Slice(sfxPackBase, sfxPackLen)
	if err != nil {
		return nil, err
	}
	outs := []output{{
		asset: Asset{
			ID: "ASSET-AUD-0001", Name: "SFX BRR sample pack (samples 0-7)",
			Category: "audio", ROMRevision: r.SHA256(),
			ROMSource:   romSource(sfxPackBase, sfxPackLen),
			ExtractorID: e.ID(), ExtractorVer: e.Version(),
			OutputPath:   path.Join(ArchiveRoot, "audio", "sfx-pack.brr"),
			OutputFormat: "audio/x-brr",
			Properties:   fmt.Sprintf("%d bytes; 9-byte BRR blocks; byte-identical to ARAM $4800-$491F", sfxPackLen),
			SHA256:       hashBytes(pack), RuntimeConsumer: "S-DSP voices via sample directory ARAM $1B00",
			CensusRefs: []string{"CEN-AUDIO-0001"}, ExperimentRefs: []string{"EXP-0024"},
			Verification: "hash-verified against tracked manifest; ARAM<->ROM identity proven in EXP-0024",
		},
		data: pack,
	}}

	// Sample 5 (battle confirm cue): 2 BRR blocks, 18 bytes, decoded to
	// 16-bit mono PCM via the project's own decoder.
	s5, err := r.Slice(sfxSample5Off, 18)
	if err != nil {
		return nil, err
	}
	pcm, err := brr.Decode(s5)
	if err != nil {
		return nil, fmt.Errorf("decoding SFX sample 5: %w", err)
	}
	var raw bytes.Buffer
	for _, s := range pcm {
		if err := binary.Write(&raw, binary.LittleEndian, s); err != nil {
			return nil, err
		}
	}
	data := raw.Bytes()
	outs = append(outs, output{
		asset: Asset{
			ID: "ASSET-AUD-0002", Name: "Battle confirm cue (SFX sample 5), decoded PCM",
			Category: "audio", ROMRevision: r.SHA256(),
			ROMSource:   romSource(sfxSample5Off, 18),
			ExtractorID: e.ID(), ExtractorVer: e.Version(),
			OutputPath:   path.Join(ArchiveRoot, "audio", "sfx-05-confirm.pcm"),
			OutputFormat: "audio/L16 (raw, little-endian, mono)",
			Properties:   fmt.Sprintf("%d samples; 16-bit signed mono; source 18 bytes = 2 BRR blocks; nominal 32000 Hz S-DSP rate before pitch scaling", len(pcm)),
			SHA256:       hashBytes(data), RuntimeConsumer: "DSP voice 7 (SRCN 5) on battle command confirm",
			CensusRefs: []string{"CEN-AUDIO-0001"}, ExperimentRefs: []string{"EXP-0024"},
			Verification: "hash-verified against tracked manifest; decoded by internal/audio/brr",
		},
		data: data,
	})
	return outs, nil
}

// ---------------------------------------------------------------------
// raw: undecoded-but-located evidence slices
// ---------------------------------------------------------------------

type rawSliceExtractor struct{}

func (rawSliceExtractor) ID() string       { return "raw-tables" }
func (rawSliceExtractor) Category() string { return "raw" }
func (rawSliceExtractor) Version() string  { return "1.0.0" }

func (e rawSliceExtractor) Extract(r *rom.ROM) ([]output, error) {
	specs := []struct {
		id, name, file string
		off, length    int
		props          string
		census, exps   []string
	}{
		{"ASSET-RAW-0001", "Monster stat record table", "monsters.bin",
			monsterBase, monsterStride * monsterSliceCount,
			fmt.Sprintf("%d records x %d bytes; +$08 current/max HP, +$0A MP, +$01 power (Confirmed); remaining fields Unknown",
				monsterSliceCount, monsterStride),
			[]string{"CEN-MONSTER-0001", "CEN-MONSTER-0004"}, []string{"EXP-0028", "EXP-0030"}},
		{"ASSET-RAW-0002", "Battle formation table", "formations.bin",
			formationBase, formationStride * formationCount,
			fmt.Sprintf("%d records x %d bytes; bytes 2-7 monster ids ($FF empty); id x15 addressing (Confirmed)",
				formationCount, formationStride),
			[]string{"CEN-EVENT-0005", "CEN-EVENT-0006"}, []string{"EXP-0030", "EXP-0032", "EXP-0034", "EXP-0036"}},
		{"ASSET-RAW-0003", "Formation flags table", "formation-flags.bin",
			formationFlagBase, formationFlagLen,
			"4-byte stride, id x4 addressing; location Confirmed, semantics Unknown",
			[]string{"CEN-EVENT-0005"}, []string{"EXP-0030"}},
	}
	outs := make([]output, 0, len(specs))
	for _, s := range specs {
		data, err := r.Slice(s.off, s.length)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", s.id, err)
		}
		outs = append(outs, output{
			asset: Asset{
				ID: s.id, Name: s.name, Category: "raw", ROMRevision: r.SHA256(),
				ROMSource:   romSource(s.off, s.length),
				ExtractorID: e.ID(), ExtractorVer: e.Version(),
				OutputPath:   path.Join(ArchiveRoot, "raw", s.file),
				OutputFormat: "application/octet-stream",
				Properties:   s.props, SHA256: hashBytes(data),
				RuntimeConsumer: "battle init per-slot loader $C22C30 / staging copy $C2315C",
				CensusRefs:      s.census, ExperimentRefs: s.exps,
				Verification: "hash-verified against tracked manifest",
			},
			data: data,
		})
	}
	return outs, nil
}
