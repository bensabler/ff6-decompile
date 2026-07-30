package brr

import "fmt"

const BlockSize = 9

type Header struct {
	Range  uint8
	Filter uint8
	Loop   bool
	End    bool
}

func ParseHeader(b byte) (Header, error) {
	h := Header{
		Range:  b >> 4,
		Filter: (b >> 2) & 0x03,
		Loop:   b&0x02 != 0,
		End:    b&0x01 != 0,
	}
	if h.Range > 12 {
		return h, fmt.Errorf("unsupported BRR range %d", h.Range)
	}
	return h, nil
}
