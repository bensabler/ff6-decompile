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
