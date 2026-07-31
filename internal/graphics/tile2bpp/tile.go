// Package tile2bpp decodes SNES 2bpp planar tiles (mode-1 BG3 format).
// The FF6 battle HUD text layer uses this format: EXP-0023 captured the
// BG3 chr bank live (chr word base $5000) and proved tiles $FF-$1FF
// byte-identical to ROMFILE:0x046FC0 + 16*index — every glyph the HUD
// tilemap references sits in that raw-copy block.
package tile2bpp

import "fmt"

// EncodedSize is one 2bpp tile: 8 rows x 2 interleaved plane bytes.
const EncodedSize = 16

// Tile holds 8x8 palette indices (0-3).
type Tile [8][8]uint8

// Decode converts one 16-byte planar tile. Row y is encoded as byte
// pair (plane0, plane1) at src[2y], src[2y+1]; pixel x takes bit 7-x
// of each plane.
func Decode(src []byte) (Tile, error) {
	var out Tile
	if len(src) < EncodedSize {
		return out, fmt.Errorf("decode 2bpp tile: need %d bytes, got %d", EncodedSize, len(src))
	}
	for y := 0; y < 8; y++ {
		p0 := src[y*2]
		p1 := src[y*2+1]
		for x := 0; x < 8; x++ {
			shift := uint(7 - x)
			out[y][x] = ((p0 >> shift) & 1) | (((p1 >> shift) & 1) << 1)
		}
	}
	return out, nil
}
