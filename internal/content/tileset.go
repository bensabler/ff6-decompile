package content

import (
	"bytes"
	"fmt"
	"image"
	"image/png"

	"github.com/bensabler/ff6-decompile/internal/graphics/framebuf"
)

// NarsheTilesetAssetID is the manifest entry for the Narshe field BG tileset.
const NarsheTilesetAssetID = "ASSET-GFX-0002"

// TilesetTiles is how many 8x8 tiles a loaded block holds.
const TilesetTiles = 256

// Tileset is a block of 8x8 BG tiles, indexed by tile number.
//
// The tile number is the VRAM tile index the PPU would use, which is what a
// tilemap entry names. Nothing here knows which map uses which tileset —
// EXP-0051 established that the selecting structure is not a pointer table and
// has not been found.
type Tileset struct {
	tiles [TilesetTiles][8][8]uint8
	ink   [TilesetTiles]bool
}

// LoadNarsheTileset reads the Narshe field tileset from the archive.
func LoadNarsheTileset(a *Archive) (*Tileset, error) {
	data, err := a.Read(NarsheTilesetAssetID)
	if err != nil {
		return nil, err
	}
	img, err := png.Decode(bytes.NewReader(data))
	if err != nil {
		return nil, fmt.Errorf("decoding the field tileset: %w", err)
	}
	pal, ok := img.(*image.Paletted)
	if !ok {
		return nil, fmt.Errorf("the field tileset decoded as %T, want *image.Paletted", img)
	}
	return newTilesetFromSheet(pal)
}

// newTilesetFromSheet reads a 16-column sheet of 8x8 tiles.
func newTilesetFromSheet(pal *image.Paletted) (*Tileset, error) {
	const cols = 16
	ts := &Tileset{}
	for i := 0; i < TilesetTiles; i++ {
		ox, oy := (i%cols)*8, (i/cols)*8
		if ox+8 > pal.Bounds().Dx() || oy+8 > pal.Bounds().Dy() {
			return nil, fmt.Errorf("tile %d at (%d,%d) is outside the %dx%d sheet",
				i, ox, oy, pal.Bounds().Dx(), pal.Bounds().Dy())
		}
		for y := 0; y < 8; y++ {
			for x := 0; x < 8; x++ {
				v := pal.ColorIndexAt(pal.Bounds().Min.X+ox+x, pal.Bounds().Min.Y+oy+y)
				ts.tiles[i][y][x] = v
				if v != 0 {
					ts.ink[i] = true
				}
			}
		}
	}
	return ts, nil
}

// NewTileset builds a tileset from project-authored pixels, for tests.
func NewTileset(tiles *[TilesetTiles][8][8]uint8) *Tileset {
	ts := &Tileset{tiles: *tiles}
	for i := range ts.tiles {
		for y := range ts.tiles[i] {
			for x := range ts.tiles[i][y] {
				if ts.tiles[i][y][x] != 0 {
					ts.ink[i] = true
				}
			}
		}
	}
	return ts
}

// Tile returns one tile and whether it carries any ink.
func (t *Tileset) Tile(n int) (*[8][8]uint8, bool) {
	if n < 0 || n >= TilesetTiles {
		return nil, false
	}
	return &t.tiles[n], t.ink[n]
}

// DrawTile blits one tile at (x, y).
func (t *Tileset) DrawTile(dst *framebuf.Indexed, x, y, n int, o framebuf.BlitOptions) {
	tile, _ := t.Tile(n)
	if tile == nil {
		return
	}
	// Field tiles are opaque BG data, not sprites: index 0 is a real colour
	// on a background layer, so it must be drawn rather than skipped.
	o.Opaque = true
	dst.BlitTile(x, y, tile, o)
}

// DrawSheet lays the whole block out in a grid, which is what a tileset
// viewer shows: every tile, in tile-number order.
func (t *Tileset) DrawSheet(dst *framebuf.Indexed, x, y, cols int, o framebuf.BlitOptions) {
	if cols <= 0 {
		cols = 16
	}
	for i := 0; i < TilesetTiles; i++ {
		t.DrawTile(dst, x+(i%cols)*8, y+(i/cols)*8, i, o)
	}
}
