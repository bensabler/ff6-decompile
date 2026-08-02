// Package snesoam decodes the SNES object attribute table.
//
// OAM is 544 bytes: 128 four-byte entries, then a 32-byte high table holding
// two more bits per sprite — the X sign and the size select. Splitting one
// sprite's attributes across two tables is the detail that makes a hand-rolled
// reader get X positions and sizes wrong, so this package does both halves.
//
// # Sizes
//
// A sprite's size bit selects between two dimensions, and which two depends on
// the PPU's object mode (OBSEL bits 5-7), not on the sprite. A sprite is
// therefore not interpretable without the mode, which is why Decode takes it.
package snesoam

import "fmt"

// Size constants.
const (
	// Size is the OAM image size in bytes: 128 four-byte entries plus the
	// 32-byte high table.
	Size = 544
	// SpriteCount is how many sprites OAM describes.
	SpriteCount = 128
	// LowTableSize is the four-byte-per-sprite region.
	LowTableSize = 512
)

// Dimensions is a sprite's pixel size.
type Dimensions struct{ W, H int }

func (d Dimensions) String() string { return fmt.Sprintf("%dx%d", d.W, d.H) }

// Tiles returns how many 8x8 tiles the sprite occupies, across and down.
func (d Dimensions) Tiles() (cols, rows int) { return d.W / 8, d.H / 8 }

// sizeModes is the OBSEL size table: for each object mode, the dimensions the
// size bit selects between.
var sizeModes = [8][2]Dimensions{
	{{8, 8}, {16, 16}},
	{{8, 8}, {32, 32}},
	{{8, 8}, {64, 64}},
	{{16, 16}, {32, 32}},
	{{16, 16}, {64, 64}},
	{{32, 32}, {64, 64}},
	{{16, 32}, {32, 64}},
	{{16, 32}, {32, 32}},
}

// SizesForMode returns the small and large dimensions an object mode selects.
func SizesForMode(mode uint8) (small, large Dimensions, err error) {
	if mode > 7 {
		return Dimensions{}, Dimensions{}, fmt.Errorf("snesoam: object mode %d out of range 0-7", mode)
	}
	m := sizeModes[mode]
	return m[0], m[1], nil
}

// Sprite is one decoded OAM entry.
type Sprite struct {
	// Index is the sprite's slot, 0-127. Slot order is priority order.
	Index int
	// X is the signed horizontal position. The high table supplies bit 8 as
	// a sign, so a sprite can sit partly off the left edge.
	X int
	// Y is the vertical position, 0-255.
	Y int
	// Tile is the 9-bit name, the sprite's top-left 8x8 tile.
	Tile int
	// Palette is 0-7, selecting a sixteen-colour group in the sprite half
	// of CGRAM (entries 128-255).
	Palette uint8
	// Priority is 0-3.
	Priority uint8
	// FlipH and FlipV mirror the whole sprite, not individual tiles.
	FlipH, FlipV bool
	// Large is the size bit: false selects the mode's small dimensions.
	Large bool
	// Size is the resolved pixel size.
	Size Dimensions
}

// CGRAMBase returns the first CGRAM entry of the sprite's palette. Sprite
// palettes occupy the upper half of CGRAM, sixteen colours each.
func (s Sprite) CGRAMBase() int { return 128 + int(s.Palette)*16 }

// OnScreen reports whether the sprite falls within the 256x224 active display.
//
// This is a filter, not hardware truth: the SNES draws whatever OAM says, and
// games park unused sprites off the bottom rather than disabling them. FF6
// leaves most of OAM parked, so without this filter every capture reads as 128
// active sprites.
func (s Sprite) OnScreen() bool {
	const w, h = 256, 224
	return s.X+s.Size.W > 0 && s.X < w && s.Y < h && s.Y+s.Size.H > 0
}

// TileAt returns the tile name at column c, row r within the sprite.
//
// The sprite's tiles are laid out in the name table as a rectangle, and the
// name table is **16 tiles wide**, so the row below a tile is +16 and not
// +cols. Getting that wrong assembles a sprite from the wrong tiles while
// still producing a plausible-looking image, which is why it is a method
// rather than left to callers.
func (s Sprite) TileAt(c, r int) int {
	return (s.Tile + r*16 + c) & 0x1FF
}

// Decode reads the OAM image. mode is the PPU object mode (OBSEL bits 5-7),
// which the size bit is interpreted against.
func Decode(oam []byte, mode uint8) ([]Sprite, error) {
	if len(oam) != Size {
		return nil, fmt.Errorf("snesoam: OAM is %d bytes, want %d", len(oam), Size)
	}
	small, large, err := SizesForMode(mode)
	if err != nil {
		return nil, err
	}

	out := make([]Sprite, SpriteCount)
	high := oam[LowTableSize:]
	for i := 0; i < SpriteCount; i++ {
		e := oam[i*4 : i*4+4]
		bits := (high[i/4] >> uint((i%4)*2)) & 0x3
		xSign := bits&1 != 0
		isLarge := bits&2 != 0

		x := int(e[0])
		if xSign {
			x -= 256
		}
		attr := e[3]
		s := Sprite{
			Index:    i,
			X:        x,
			Y:        int(e[1]),
			Tile:     int(e[2]) | int(attr&1)<<8,
			Palette:  (attr >> 1) & 7,
			Priority: (attr >> 4) & 3,
			FlipH:    attr&0x40 != 0,
			FlipV:    attr&0x80 != 0,
			Large:    isLarge,
			Size:     small,
		}
		if isLarge {
			s.Size = large
		}
		out[i] = s
	}
	return out, nil
}

// OnScreen filters a decoded table to the sprites the display shows.
func OnScreen(sprites []Sprite) []Sprite {
	var out []Sprite
	for _, s := range sprites {
		if s.OnScreen() {
			out = append(out, s)
		}
	}
	return out
}

// VRAMAddress returns the byte offset of a sprite tile within VRAM.
//
// base is the object name base address in bytes (the PPU register holds it as
// a word address; mesenstate doubles it). Names wrap within the 64 KB space.
func VRAMAddress(base, tile int) int {
	const bytesPerTile = 32
	return (base + tile*bytesPerTile) & 0xFFFF
}
