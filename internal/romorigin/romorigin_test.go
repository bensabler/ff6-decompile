package romorigin

import (
	"bytes"
	"math/rand"
	"testing"
)

// noise builds distinctive, non-repeating bytes so probes are selective.
func noise(n int, seed int64) []byte {
	r := rand.New(rand.NewSource(seed))
	b := make([]byte, n)
	for i := range b {
		b[i] = byte(r.Intn(256))
	}
	return b
}

func TestTraceFindsACopiedSpan(t *testing.T) {
	rom := noise(8192, 1)
	// The image is a 1024-byte copy taken from ROM offset 2048.
	image := append([]byte(nil), rom[2048:3072]...)

	got := Trace(image, rom, DefaultOptions())
	if len(got) != 1 {
		t.Fatalf("got %d blocks, want 1: %+v", len(got), got)
	}
	want := Block{ImageOffset: 0, ROMOffset: 2048, Length: 1024}
	if got[0] != want {
		t.Errorf("got %+v, want %+v", got[0], want)
	}
}

func TestTraceReportsSeveralSpansInOrder(t *testing.T) {
	rom := noise(16384, 2)
	var image []byte
	image = append(image, rom[1000:1512]...) // 512 from 1000
	image = append(image, rom[9000:9256]...) // 256 from 9000

	got := Trace(image, rom, DefaultOptions())
	if len(got) != 2 {
		t.Fatalf("got %d blocks, want 2: %+v", len(got), got)
	}
	if got[0].ROMOffset != 1000 || got[0].Length != 512 {
		t.Errorf("first block %+v, want ROM 1000 len 512", got[0])
	}
	if got[1].ImageOffset != 512 || got[1].ROMOffset != 9000 || got[1].Length != 256 {
		t.Errorf("second block %+v, want image 512, ROM 9000, len 256", got[1])
	}
}

// The central negative case. A transformed region must not be reported as
// copied, because the whole point of the tool is to tell those apart.
func TestTraceDoesNotMatchTransformedData(t *testing.T) {
	rom := noise(8192, 3)
	image := append([]byte(nil), rom[2048:3072]...)
	for i := range image {
		image[i] ^= 0x5A // any transform at all
	}

	if got := Trace(image, rom, DefaultOptions()); len(got) != 0 {
		t.Errorf("transformed data reported as copied: %+v", got)
	}
}

// A single changed byte splits one run into two. This is recorded as a real
// limitation: "not verbatim" is a lead toward compression, never proof of it.
func TestOneChangedByteSplitsTheRun(t *testing.T) {
	rom := noise(8192, 4)
	image := append([]byte(nil), rom[1024:3072]...)
	image[1000] ^= 0xFF

	got := Trace(image, rom, DefaultOptions())
	if len(got) < 2 {
		t.Fatalf("want the run split by the edit, got %+v", got)
	}
	if Coverage(got) >= len(image) {
		t.Error("coverage should fall short of the image when a byte differs")
	}
}

func TestTraceSkipsUniformWindows(t *testing.T) {
	rom := noise(8192, 5)
	// Blank tiles: a long run of zeros matches almost anywhere, and
	// reporting it would be provenance noise.
	image := make([]byte, 512)

	if got := Trace(image, rom, DefaultOptions()); len(got) != 0 {
		t.Errorf("uniform image should produce no blocks, got %+v", got)
	}
}

func TestMinRunDropsShortMatches(t *testing.T) {
	rom := noise(8192, 6)
	image := append([]byte(nil), rom[100:180]...) // 80 bytes

	o := DefaultOptions()
	o.MinRun = 64
	if got := Trace(image, rom, o); len(got) != 1 {
		t.Errorf("an 80-byte run should survive MinRun 64, got %+v", got)
	}
	o.MinRun = 128
	if got := Trace(image, rom, o); len(got) != 0 {
		t.Errorf("an 80-byte run should be dropped by MinRun 128, got %+v", got)
	}
}

func TestLimitBoundsTheScan(t *testing.T) {
	rom := noise(8192, 7)
	image := append(append([]byte(nil), rom[0:512]...), rom[4096:4608]...)

	o := DefaultOptions()
	o.Limit = 512
	got := Trace(image, rom, o)
	for _, b := range got {
		if b.ImageOffset >= 512 {
			t.Errorf("block %+v starts past the limit", b)
		}
	}
}

func TestMergeJoinsContiguousRuns(t *testing.T) {
	in := []Block{
		{ImageOffset: 0, ROMOffset: 1000, Length: 256},
		{ImageOffset: 256, ROMOffset: 1256, Length: 256}, // continues
		{ImageOffset: 512, ROMOffset: 9000, Length: 128}, // does not
	}
	got := Merge(in)
	if len(got) != 2 {
		t.Fatalf("got %d blocks, want 2: %+v", len(got), got)
	}
	if got[0].Length != 512 || got[0].ROMOffset != 1000 {
		t.Errorf("merged block %+v, want ROM 1000 len 512", got[0])
	}

	// Contiguous in the image but not in the ROM must NOT merge: that is
	// two separate transfers landing side by side, which is exactly the
	// pattern an animated-tile table produces.
	split := []Block{
		{ImageOffset: 0, ROMOffset: 1000, Length: 128},
		{ImageOffset: 128, ROMOffset: 5000, Length: 128},
	}
	if got := Merge(split); len(got) != 2 {
		t.Errorf("image-contiguous but ROM-disjoint runs must stay separate, got %+v", got)
	}
}

func TestCoverage(t *testing.T) {
	tests := []struct {
		name   string
		blocks []Block
		want   int
	}{
		{"empty", nil, 0},
		{"one", []Block{{ImageOffset: 0, Length: 100}}, 100},
		{"two disjoint", []Block{{ImageOffset: 0, Length: 100}, {ImageOffset: 200, Length: 50}}, 150},
		{"overlapping", []Block{{ImageOffset: 0, Length: 100}, {ImageOffset: 50, Length: 100}}, 150},
		{"contained", []Block{{ImageOffset: 0, Length: 100}, {ImageOffset: 20, Length: 10}}, 100},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Coverage(tt.blocks); got != tt.want {
				t.Errorf("Coverage = %d, want %d", got, tt.want)
			}
		})
	}
}

func TestTraceHandlesDegenerateInput(t *testing.T) {
	rom := noise(1024, 8)
	for _, tt := range []struct {
		name         string
		image, romIn []byte
	}{
		{"empty image", nil, rom},
		{"empty rom", noise(256, 9), nil},
		{"both empty", nil, nil},
		{"image longer than rom", noise(4096, 10), rom},
	} {
		t.Run(tt.name, func(t *testing.T) {
			_ = Trace(tt.image, tt.romIn, DefaultOptions()) // must not panic
		})
	}

	// A probe larger than the image must terminate rather than spin.
	o := DefaultOptions()
	o.Probe = 4096
	if got := Trace(noise(64, 11), rom, o); len(got) != 0 {
		t.Errorf("probe larger than the image should find nothing, got %+v", got)
	}
}

func FuzzTrace(f *testing.F) {
	f.Add([]byte("the quick brown fox"), []byte("brown"))
	f.Add([]byte{}, []byte{0})
	f.Fuzz(func(t *testing.T, image, rom []byte) {
		o := DefaultOptions()
		o.Probe = 4
		o.MinRun = 4
		o.MinDistinct = 2
		for _, b := range Trace(image, rom, o) {
			if b.ImageOffset < 0 || b.ROMOffset < 0 || b.Length <= 0 {
				t.Fatalf("negative or empty block %+v", b)
			}
			if b.ImageOffset+b.Length > len(image) || b.ROMOffset+b.Length > len(rom) {
				t.Fatalf("block %+v runs past image (%d) or rom (%d)", b, len(image), len(rom))
			}
			if !bytes.Equal(image[b.ImageOffset:b.ImageOffset+b.Length],
				rom[b.ROMOffset:b.ROMOffset+b.Length]) {
				t.Fatalf("block %+v does not actually match", b)
			}
		}
	})
}
