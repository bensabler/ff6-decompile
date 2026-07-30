package tile4bpp

import "fmt"

const EncodedSize = 32

type Tile [8][8]uint8

func Decode(src []byte) (Tile, error) {
	var out Tile
	if len(src) < EncodedSize {
		return out, fmt.Errorf("decode 4bpp tile: need %d bytes, got %d", EncodedSize, len(src))
	}
	for y := 0; y < 8; y++ {
		p0 := src[y*2]
		p1 := src[y*2+1]
		p2 := src[16+y*2]
		p3 := src[16+y*2+1]
		for x := 0; x < 8; x++ {
			shift := uint(7 - x)
			out[y][x] =
				((p0 >> shift) & 1) |
					(((p1 >> shift) & 1) << 1) |
					(((p2 >> shift) & 1) << 2) |
					(((p3 >> shift) & 1) << 3)
		}
	}
	return out, nil
}
