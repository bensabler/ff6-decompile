package tile2bpp

import "testing"

func TestDecodeZero(t *testing.T) {
	got, err := Decode(make([]byte, EncodedSize))
	if err != nil {
		t.Fatal(err)
	}
	if got[0][0] != 0 || got[7][7] != 0 {
		t.Fatalf("zero tile decoded nonzero")
	}
}

// TestDecodePlanes uses a synthetic tile exercising both planes:
// row 0 has plane0=$FF plane1=$00 (all pixels 1), row 1 the reverse
// (all pixels 2), row 2 both ($FF,$FF -> all 3), row 3 alternating
// plane0=$AA plane1=$55 (pixels 1,2,1,2...).
func TestDecodePlanes(t *testing.T) {
	src := make([]byte, EncodedSize)
	src[0], src[1] = 0xFF, 0x00
	src[2], src[3] = 0x00, 0xFF
	src[4], src[5] = 0xFF, 0xFF
	src[6], src[7] = 0xAA, 0x55
	got, err := Decode(src)
	if err != nil {
		t.Fatal(err)
	}
	for x := 0; x < 8; x++ {
		if got[0][x] != 1 {
			t.Fatalf("row0 x%d = %d, want 1", x, got[0][x])
		}
		if got[1][x] != 2 {
			t.Fatalf("row1 x%d = %d, want 2", x, got[1][x])
		}
		if got[2][x] != 3 {
			t.Fatalf("row2 x%d = %d, want 3", x, got[2][x])
		}
		want := uint8(1)
		if x%2 == 1 {
			want = 2
		}
		if got[3][x] != want {
			t.Fatalf("row3 x%d = %d, want %d", x, got[3][x], want)
		}
	}
	if got[4][0] != 0 || got[7][7] != 0 {
		t.Fatalf("untouched rows decoded nonzero")
	}
}

func TestDecodeShort(t *testing.T) {
	if _, err := Decode(make([]byte, EncodedSize-1)); err == nil {
		t.Fatal("expected error")
	}
}

func FuzzDecode(f *testing.F) {
	f.Add(make([]byte, EncodedSize))
	f.Fuzz(func(t *testing.T, data []byte) {
		_, _ = Decode(data)
	})
}
