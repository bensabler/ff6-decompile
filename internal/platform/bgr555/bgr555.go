package bgr555

type Color struct {
	R uint8
	G uint8
	B uint8
}

func Decode(v uint16) Color {
	return Color{
		R: expand5(uint8(v & 0x1F)),
		G: expand5(uint8((v >> 5) & 0x1F)),
		B: expand5(uint8((v >> 10) & 0x1F)),
	}
}

func expand5(v uint8) uint8 {
	return (v << 3) | (v >> 2)
}
