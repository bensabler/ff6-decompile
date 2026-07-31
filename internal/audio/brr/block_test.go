package brr

import "testing"

func TestParseHeader(t *testing.T) {
	h, err := ParseHeader(0x0F)
	if err != nil {
		t.Fatal(err)
	}
	if h.Range != 0 || h.Filter != 3 || !h.Loop || !h.End {
		t.Fatalf("unexpected header: %+v", h)
	}
}

func FuzzParseHeader(f *testing.F) {
	f.Add(byte(0))
	f.Fuzz(func(t *testing.T, b byte) {
		_, _ = ParseHeader(b)
	})
}

// TestDecodeFilter0 hand-computes a one-block direct (filter 0) stream:
// range 5, nibbles 1,2,-1,0,... -> (n<<5)>>1 = 16, 32, -16, 0.
func TestDecodeFilter0(t *testing.T) {
	src := make([]byte, BlockSize)
	src[0] = 0x51 // range 5, filter 0, end
	src[1] = 0x12 // nibbles 1, 2
	src[2] = 0xF0 // nibbles -1, 0
	got, err := Decode(src)
	if err != nil {
		t.Fatal(err)
	}
	if len(got) != 16 {
		t.Fatalf("got %d samples, want 16", len(got))
	}
	want := []int16{16, 32, -16, 0}
	for i, w := range want {
		if got[i] != w {
			t.Fatalf("sample %d = %d, want %d", i, got[i], w)
		}
	}
}

// TestDecodeFilter1 hand-computes the first predictor: s = shifted +
// p1 + (-p1>>4). Nibbles 1,0,0,0 at range 5 give 16, 15, 14, 13.
func TestDecodeFilter1(t *testing.T) {
	src := make([]byte, BlockSize)
	src[0] = 0x55 // range 5, filter 1, end
	src[1] = 0x10
	got, err := Decode(src)
	if err != nil {
		t.Fatal(err)
	}
	want := []int16{16, 15, 14, 13}
	for i, w := range want {
		if got[i] != w {
			t.Fatalf("sample %d = %d, want %d", i, got[i], w)
		}
	}
}

func TestDecodeMultiBlockAndMissingEnd(t *testing.T) {
	two := make([]byte, BlockSize*2)
	two[0] = 0x00         // no end
	two[BlockSize] = 0x01 // end
	got, err := Decode(two)
	if err != nil || len(got) != 32 {
		t.Fatalf("two blocks: got %d samples, err %v", len(got), err)
	}
	if _, err := Decode(make([]byte, BlockSize)); err == nil {
		t.Fatal("stream without End flag should error")
	}
}

func FuzzDecode(f *testing.F) {
	seed := make([]byte, BlockSize)
	seed[0] = 0x01
	f.Add(seed)
	f.Fuzz(func(t *testing.T, data []byte) {
		_, _ = Decode(data)
	})
}
