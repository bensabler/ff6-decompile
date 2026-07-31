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

// Decode decodes a BRR stream from its first block through the block
// carrying the End flag, returning the 16-bit PCM samples (16 per
// block). Prediction state starts at zero, matching a voice keyed on
// at the sample start. Standard S-DSP semantics: nibble<<range then
// >>1, filter prediction from the two previous outputs, clamp to
// 16-bit, clip to 15-bit (sign-preserving wrap).
func Decode(src []byte) ([]int16, error) {
	if len(src) < BlockSize {
		return nil, fmt.Errorf("decode BRR: need at least %d bytes, got %d", BlockSize, len(src))
	}
	var out []int16
	var p1, p2 int32
	for off := 0; ; off += BlockSize {
		if off+BlockSize > len(src) {
			return nil, fmt.Errorf("decode BRR: stream ended at byte %d without an End block", len(src))
		}
		h, err := ParseHeader(src[off])
		if err != nil {
			return nil, fmt.Errorf("decode BRR block at byte %d: %w", off, err)
		}
		for i := 0; i < 16; i++ {
			nib := src[off+1+i/2]
			var s int32
			if i%2 == 0 {
				s = int32(nib >> 4)
			} else {
				s = int32(nib & 0x0F)
			}
			if s >= 8 {
				s -= 16
			}
			s = (s << h.Range) >> 1
			switch h.Filter {
			case 1:
				s += p1 + (-p1 >> 4)
			case 2:
				s += 2*p1 + (-(3 * p1) >> 5) - p2 + (p2 >> 4)
			case 3:
				s += 2*p1 + (-(13 * p1) >> 6) - p2 + ((3 * p2) >> 4)
			}
			if s > 32767 {
				s = 32767
			} else if s < -32768 {
				s = -32768
			}
			s = (s << 17) >> 17 // 15-bit clip, sign preserved
			p2, p1 = p1, s
			out = append(out, int16(s))
		}
		if h.End {
			return out, nil
		}
	}
}
