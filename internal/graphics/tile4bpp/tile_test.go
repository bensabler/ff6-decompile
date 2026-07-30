package tile4bpp

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
